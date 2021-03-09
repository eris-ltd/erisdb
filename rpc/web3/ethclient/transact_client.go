package ethclient

import (
	"context"
	"fmt"

	"github.com/hyperledger/burrow/acm"
	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/execution/exec"
	"github.com/hyperledger/burrow/rpc"
	"github.com/hyperledger/burrow/rpc/web3"
	"github.com/hyperledger/burrow/txs"
	"github.com/hyperledger/burrow/txs/payload"
	"google.golang.org/grpc"
)

const BasicGasLimit = 21000

// Provides a partial implementation of the GRPC-generated TransactClient suitable for testing Vent on Ethereum
type TransactClient struct {
	client   rpc.Client
	chainID  string
	accounts []acm.AddressableSigner
}

func NewTransactClient(client rpc.Client) *TransactClient {
	return &TransactClient{
		client: client,
	}
}

func (cli *TransactClient) WithAccounts(signers ...acm.AddressableSigner) *TransactClient {
	return &TransactClient{
		client:   cli.client,
		accounts: append(cli.accounts, signers...),
	}
}

func (cli *TransactClient) CallTxSync(ctx context.Context, tx *payload.CallTx,
	opts ...grpc.CallOption) (*exec.TxExecution, error) {

	var signer acm.AddressableSigner

	for _, sa := range cli.accounts {
		if sa.GetAddress() == tx.Input.Address {
			signer = sa
			break
		}
	}

	err := cli.completeTx(tx)
	if err != nil {
		return nil, fmt.Errorf("could not set values on transaction")
	}

	var txHash string
	if signer == nil {
		txHash, err = cli.SendTransaction(tx)
	} else {
		txHash, err = cli.SendRawTransaction(tx, signer)
	}
	if err != nil {
		return nil, fmt.Errorf("could not send ethereum transaction: %w", err)
	}

	fmt.Printf("Waiting for tranasaction %s to be confirmed...\n", txHash)
	receipt, err := AwaitTransaction(ctx, cli.client, txHash)
	if err != nil {
		return nil, err
	}

	d := new(web3.HexDecoder)

	header := &exec.TxHeader{
		TxType: payload.TypeCall,
		TxHash: d.Bytes(receipt.TransactionHash),
		Height: d.Uint64(receipt.BlockNumber),
		Index:  d.Uint64(receipt.TransactionIndex),
	}

	// Attempt to provide sufficient return values to satisfy Vent's needs.
	return &exec.TxExecution{
		TxHeader: header,
		Receipt: &txs.Receipt{
			TxType:          header.TxType,
			TxHash:          header.TxHash,
			CreatesContract: receipt.ContractAddress != "",
			ContractAddress: d.Address(receipt.ContractAddress),
		},
	}, d.Err()
}

func (cli *TransactClient) SendTransaction(tx *payload.CallTx) (string, error) {
	var to string
	if tx.Address != nil {
		to = web3.HexEncoder.Address(*tx.Address)
	}

	var nonce string
	if tx.Input.Sequence != 0 {
		nonce = web3.HexEncoder.Uint64OmitEmpty(tx.Input.Sequence)
	}

	param := &EthSendTransactionParam{
		From:     web3.HexEncoder.Address(tx.Input.Address),
		To:       to,
		Gas:      web3.HexEncoder.Uint64OmitEmpty(tx.GasLimit),
		GasPrice: web3.HexEncoder.Uint64OmitEmpty(tx.GasPrice),
		Value:    web3.HexEncoder.Uint64OmitEmpty(tx.Input.Amount),
		Data:     web3.HexEncoder.BytesTrim(tx.Data),
		Nonce:    nonce,
	}

	return EthSendTransaction(cli.client, param)
}

func (cli *TransactClient) SendRawTransaction(tx *payload.CallTx, signer acm.AddressableSigner) (string, error) {
	chainID, err := cli.GetChainID()
	if err != nil {
		return "", err
	}
	txEnv := txs.Enclose(chainID, tx)

	txEnv.Encoding = txs.Envelope_RLP

	err = txEnv.Sign(signer)
	if err != nil {
		return "", fmt.Errorf("could not sign Ethereum transaction: %w", err)
	}

	rawTx, err := txs.EthRawTxFromEnvelope(txEnv)
	if err != nil {
		return "", fmt.Errorf("could not generate Ethereum raw transaction: %w", err)
	}

	bs, err := rawTx.Marshal()
	if err != nil {
		return "", fmt.Errorf("could not marshal Ethereum raw transaction: %w", err)
	}

	return EthSendRawTransaction(cli.client, web3.HexEncoder.BytesTrim(bs))
}

func (cli *TransactClient) GetChainID() (string, error) {
	if cli.chainID == "" {
		var err error
		cli.chainID, err = NetVersion(cli.client)
		if err != nil {
			return "", fmt.Errorf("TransactClient could not get ChainID: %w", err)
		}
	}
	return cli.chainID, nil
}

func (cli *TransactClient) GetGasPrice() (uint64, error) {
	gasPrice, err := EthGasPrice(cli.client)
	if err != nil {
		return 0, fmt.Errorf("could not get gas price: %w", err)
	}
	d := new(web3.HexDecoder)
	return d.Uint64(gasPrice), d.Err()
}

func (cli *TransactClient) GetTransactionCount(address crypto.Address) (uint64, error) {
	count, err := EthGetTransactionCount(cli.client, address)
	if err != nil {
		return 0, fmt.Errorf("could not get transaction acount for address %s: %w", address, err)
	}
	d := new(web3.HexDecoder)
	return d.Uint64(count), d.Err()
}

func (cli *TransactClient) completeTx(tx *payload.CallTx) error {
	if tx.GasLimit == 0 {
		tx.GasLimit = BasicGasLimit
	}
	var err error
	if tx.GasPrice == 0 {
		tx.GasPrice, err = cli.GetGasPrice()
		if err != nil {
			return err
		}
	}
	if tx.Input.Sequence == 0 {
		tx.Input.Sequence, err = cli.GetTransactionCount(tx.Input.Address)
		if err != nil {
			return err
		}
	}
	return nil
}
