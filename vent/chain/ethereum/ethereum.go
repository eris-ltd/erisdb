package ethereum

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hyperledger/burrow/event"
	"github.com/hyperledger/burrow/event/query"
	"github.com/hyperledger/burrow/logging"
	"github.com/hyperledger/burrow/rpc/web3"
	"github.com/hyperledger/burrow/vent/chain"

	"github.com/hyperledger/burrow/rpc/web3/ethclient"

	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/execution/errors"
	"github.com/hyperledger/burrow/execution/exec"
	"github.com/hyperledger/burrow/rpc"
	"github.com/hyperledger/burrow/rpc/rpcevents"
	"github.com/hyperledger/burrow/vent/types"
	"google.golang.org/grpc/connectivity"
)

type Chain struct {
	client  rpc.Client
	filter  *chain.Filter
	chainID string
	version string
	logger  *logging.Logger
}

var _ chain.Chain = (*Chain)(nil)

// We rely on this failing if the chain is not an Ethereum Chain
func New(client rpc.Client, filter *chain.Filter, logger *logging.Logger) (*Chain, error) {
	chainID, err := ethclient.NetVersion(client)
	if err != nil {
		return nil, fmt.Errorf("could not get Ethereum ChainID: %w", err)
	}
	version, err := ethclient.Web3ClientVersion(client)
	if err != nil {
		return nil, fmt.Errorf("could not get Ethereum node version: %w", err)
	}
	return &Chain{
		client:  client,
		filter:  filter,
		chainID: chainID,
		version: version,
		logger:  logger,
	}, nil
}

func (c *Chain) StatusMessage(ctx context.Context, lastProcessedHeight uint64) []interface{} {
	// TODO: more info is available from web3
	return []interface{}{
		"msg", "status",
		"chain_type", "Ethereum",
		"last_processed_height", lastProcessedHeight,
	}
}

func (c *Chain) GetABI(ctx context.Context, address crypto.Address) (string, error) {
	// Unsupported by Ethereum
	return "", nil
}

func (c *Chain) GetVersion() string {
	return c.version
}

func (c *Chain) GetChainID() string {
	return c.chainID
}

func (c *Chain) ConsumeBlocks(ctx context.Context, in *rpcevents.BlockRange, consumer func(chain.Block) error) error {
	return Consume(c.client, c.filter, in, c.logger, consumer)
}

func (c *Chain) Connectivity() connectivity.State {
	// TODO: better connectivity information
	_, err := ethclient.EthSyncing(c.client)
	if err != nil {
		return connectivity.TransientFailure
	}
	return connectivity.Ready
}

func (c *Chain) Close() error {
	// just a http.Client - nothing to free
	return nil
}

type Block struct {
	Height       uint64
	Transactions []chain.Transaction
}

func NewEthereumBlock(log *Event) *Block {
	return &Block{
		Height:       log.Height,
		Transactions: []chain.Transaction{NewEthereumTransaction(log)},
	}
}

var _ chain.Block = (*Block)(nil)

func (b *Block) GetHeight() uint64 {
	return b.Height
}

func (b *Block) GetTxs() []chain.Transaction {
	return b.Transactions
}

func (b *Block) GetMetadata(columns types.SQLColumnNames) (map[string]interface{}, error) {
	return map[string]interface{}{
		columns.Height: fmt.Sprintf("%v", b.Height),
	}, nil
}

func (b *Block) appendTransaction(log *Event) {
	b.Transactions = append(b.Transactions, &Transaction{
		Index:  uint64(len(b.Transactions)),
		Hash:   log.TransactionHash,
		Events: []chain.Event{log},
	})
}

func (b *Block) appendEvent(log *Event) {
	tx := b.Transactions[len(b.Transactions)-1].(*Transaction)
	log.Index = uint64(len(tx.Events))
	tx.Events = append(tx.Events, log)
}

type Transaction struct {
	Height uint64
	Index  uint64
	Hash   binary.HexBytes
	Events []chain.Event
}

func NewEthereumTransaction(log *Event) *Transaction {
	return &Transaction{
		Height: log.Height,
		Index:  0,
		Hash:   log.TransactionHash,
		Events: []chain.Event{log},
	}
}

func (tx *Transaction) GetHash() binary.HexBytes {
	return tx.Hash
}

func (tx *Transaction) GetIndex() uint64 {
	return tx.Index
}

func (tx *Transaction) GetEvents() []chain.Event {
	return tx.Events
}

func (tx *Transaction) GetException() *errors.Exception {
	// Ethereum does not retain an log from reverted transactions
	return nil
}

func (tx *Transaction) GetOrigin() *chain.Origin {
	// Origin refers to a previous dumped chain which is not a concept in Ethereum
	return nil
}

func (tx *Transaction) GetMetadata(columns types.SQLColumnNames) (map[string]interface{}, error) {
	return map[string]interface{}{
		columns.Height:  tx.Height,
		columns.TxHash:  tx.Hash.String(),
		columns.TxIndex: tx.Index,
		columns.TxType:  exec.TypeLog.String(),
	}, nil
}

var _ chain.Transaction = (*Transaction)(nil)

type Event struct {
	exec.LogEvent
	Height uint64
	// Index of event in entire block (what ethereum provides us with
	IndexInBlock uint64
	// Index of event in transaction
	Index           uint64
	TransactionHash binary.HexBytes
}

var _ chain.Event = (*Event)(nil)

func NewEthereumEvent(log *ethclient.EthLog) (*Event, error) {
	d := new(web3.HexDecoder)
	topics := make([]binary.Word256, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = binary.LeftPadWord256(d.Bytes(t))
	}
	txHash := d.Bytes(log.TransactionHash)
	return &Event{
		LogEvent: exec.LogEvent{
			Topics:  topics,
			Address: d.Address(log.Address),
			Data:    d.Bytes(log.Data),
		},
		Height:          d.Uint64(log.BlockNumber),
		IndexInBlock:    d.Uint64(log.LogIndex),
		TransactionHash: txHash,
	}, d.Err()
}

func (ev *Event) GetIndex() uint64 {
	return ev.Index
}

func (ev *Event) GetTransactionHash() binary.HexBytes {
	return ev.TransactionHash
}

func (ev *Event) GetAddress() crypto.Address {
	return ev.Address
}

func (ev *Event) GetTopics() []binary.Word256 {
	return ev.Topics
}

func (ev *Event) GetData() []byte {
	return ev.Data
}

func (ev *Event) Get(key string) (value interface{}, ok bool) {
	switch key {
	case event.EventTypeKey:
		return exec.TypeLog, true
	}
	v, ok := ev.LogEvent.Get(key)
	if ok {
		return v, ok
	}
	return query.GetReflect(reflect.ValueOf(ev), key)
}
