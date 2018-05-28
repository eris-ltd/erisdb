package abci

import (
	"fmt"
	"sync"
	"time"

	bcm "github.com/hyperledger/burrow/blockchain"
	"github.com/hyperledger/burrow/consensus/tendermint/codes"
	"github.com/hyperledger/burrow/execution"
	"github.com/hyperledger/burrow/logging"
	"github.com/hyperledger/burrow/logging/structure"
	"github.com/hyperledger/burrow/project"
	"github.com/hyperledger/burrow/txs"
	"github.com/pkg/errors"
	abci_types "github.com/tendermint/abci/types"
	"github.com/tendermint/go-wire"
)

const responseInfoName = "Burrow"

type App struct {
	// State
	blockchain    bcm.MutableBlockchain
	checker       execution.BatchExecutor
	committer     execution.BatchCommitter
	mempoolLocker sync.Locker
	// We need to cache these from BeginBlock for when we need actually need it in Commit
	block *abci_types.RequestBeginBlock
	// Utility
	txDecoder txs.Decoder
	// Logging
	logger *logging.Logger
}

var _ abci_types.Application = &App{}

func NewApp(blockchain bcm.MutableBlockchain,
	checker execution.BatchExecutor,
	committer execution.BatchCommitter,
	logger *logging.Logger) *App {
	return &App{
		blockchain: blockchain,
		checker:    checker,
		committer:  committer,
		txDecoder:  txs.NewGoWireCodec(),
		logger:     logger.WithScope("abci.NewApp").With(structure.ComponentKey, "ABCI_App"),
	}
}

// Provide the Mempool lock. When provided we will attempt to acquire this lock in a goroutine during the Commit. We
// will keep the checker cache locked until we are able to acquire the mempool lock which signals the end of the commit
// and possible recheck on Tendermint's side.
func (app *App) SetMempoolLocker(mempoolLocker sync.Locker) {
	app.mempoolLocker = mempoolLocker
}

func (app *App) Info(info abci_types.RequestInfo) abci_types.ResponseInfo {
	tip := app.blockchain.Tip()
	return abci_types.ResponseInfo{
		Data:             responseInfoName,
		Version:          project.History.CurrentVersion().String(),
		LastBlockHeight:  int64(tip.LastBlockHeight()),
		LastBlockAppHash: tip.AppHashAfterLastBlock(),
	}
}

func (app *App) SetOption(option abci_types.RequestSetOption) (respSetOption abci_types.ResponseSetOption) {
	respSetOption.Log = "SetOption not supported"
	respSetOption.Code = codes.UnsupportedRequestCode
	return
}

func (app *App) Query(reqQuery abci_types.RequestQuery) (respQuery abci_types.ResponseQuery) {
	respQuery.Log = "Query not supported"
	respQuery.Code = codes.UnsupportedRequestCode
	return
}

func (app *App) CheckTx(txBytes []byte) abci_types.ResponseCheckTx {
	tx, err := app.txDecoder.DecodeTx(txBytes)
	if err != nil {
		app.logger.TraceMsg("CheckTx decoding error",
			"tag", "CheckTx",
			structure.ErrorKey, err)
		return abci_types.ResponseCheckTx{
			Code: codes.EncodingErrorCode,
			Log:  fmt.Sprintf("Encoding error: %s", err),
		}
	}
	receipt := txs.GenerateReceipt(app.blockchain.ChainID(), tx)

	err = app.checker.Execute(tx)
	if err != nil {
		app.logger.TraceMsg("CheckTx execution error",
			structure.ErrorKey, err,
			"tag", "CheckTx",
			"tx_hash", receipt.TxHash,
			"creates_contract", receipt.CreatesContract)
		return abci_types.ResponseCheckTx{
			Code: codes.EncodingErrorCode,
			Log:  fmt.Sprintf("CheckTx could not execute transaction: %s, error: %v", tx, err),
		}
	}

	receiptBytes := wire.BinaryBytes(receipt)
	app.logger.TraceMsg("CheckTx success",
		"tag", "CheckTx",
		"tx_hash", receipt.TxHash,
		"creates_contract", receipt.CreatesContract)
	return abci_types.ResponseCheckTx{
		Code: codes.TxExecutionSuccessCode,
		Log:  "CheckTx success - receipt in data",
		Data: receiptBytes,
	}
}

func (app *App) InitChain(chain abci_types.RequestInitChain) (respInitChain abci_types.ResponseInitChain) {
	// Could verify agreement on initial validator set here
	return
}

func (app *App) BeginBlock(block abci_types.RequestBeginBlock) (respBeginBlock abci_types.ResponseBeginBlock) {
	nextHeigth := uint64(block.Header.GetHeight() + 1)
	app.blockchain.EvaluateSortition(nextHeigth, block.Hash)

	app.block = &block
	return
}

func (app *App) DeliverTx(txBytes []byte) abci_types.ResponseDeliverTx {
	tx, err := app.txDecoder.DecodeTx(txBytes)
	if err != nil {
		app.logger.TraceMsg("DeliverTx decoding error",
			"tag", "DeliverTx",
			structure.ErrorKey, err)
		return abci_types.ResponseDeliverTx{
			Code: codes.EncodingErrorCode,
			Log:  fmt.Sprintf("Encoding error: %s", err),
		}
	}

	receipt := txs.GenerateReceipt(app.blockchain.ChainID(), tx)
	err = app.committer.Execute(tx)
	if err != nil {
		app.logger.TraceMsg("DeliverTx execution error",
			structure.ErrorKey, err,
			"tag", "DeliverTx",
			"tx_hash", receipt.TxHash,
			"creates_contract", receipt.CreatesContract)
		return abci_types.ResponseDeliverTx{
			Code: codes.TxExecutionErrorCode,
			Log:  fmt.Sprintf("DeliverTx could not execute transaction: %s, error: %s", tx, err),
		}
	}

	app.logger.TraceMsg("DeliverTx success",
		"tag", "DeliverTx",
		"tx_hash", receipt.TxHash,
		"creates_contract", receipt.CreatesContract)
	receiptBytes := wire.BinaryBytes(receipt)
	return abci_types.ResponseDeliverTx{
		Code: codes.TxExecutionSuccessCode,
		Log:  "DeliverTx success - receipt in data",
		Data: receiptBytes,
	}
}

func (app *App) EndBlock(reqEndBlock abci_types.RequestEndBlock) (respEndBlock abci_types.ResponseEndBlock) {
	// Validator mutation goes here
	validatorSet := app.blockchain.Tip().ValidatorSet()
	validatorSet.AdjustPower(reqEndBlock.GetHeight())
	validators := validatorSet.Validators()
	setLeavers := validatorSet.SetLeavers()

	updates := make([]abci_types.Validator, len(validators)+len(setLeavers))

	for i, validator := range validators {
		updates[i].Power = validator.Power()
		updates[i].PubKey = validator.PublicKey().Bytes()
	}

	for i, validator := range setLeavers {
		updates[i+len(validators)].Power = 0
		updates[i+len(validators)].PubKey = validator.PublicKey().Bytes()
	}

	respEndBlock.ValidatorUpdates = updates
	return
}

func (app *App) Commit() abci_types.ResponseCommit {
	tip := app.blockchain.Tip()
	app.logger.InfoMsg("Committing block",
		"tag", "Commit",
		structure.ScopeKey, "Commit()",
		"block_height", app.block.Header.Height,
		"block_hash", app.block.Hash,
		"block_time", app.block.Header.Time,
		"num_txs", app.block.Header.NumTxs,
		"last_block_time", tip.LastBlockTime(),
		"last_block_hash", tip.LastBlockHash())

	// Lock the checker while we reset it and possibly while recheckTxs replays transactions
	app.checker.Lock()
	defer func() {
		// Tendermint may replay transactions to the check cache during a recheck, which happens after we have returned
		// from Commit(). The mempool is locked by Tendermint for the duration of the commit phase; during Commit() and
		// the subsequent mempool.Update() so we schedule an acquisition of the mempool lock in a goroutine in order to
		// 'observe' the mempool unlock event that happens later on. By keeping the checker read locked during that
		// period we can ensure that anything querying the checker (such as service.MempoolAccounts()) will block until
		// the full Tendermint commit phase has completed.
		if app.mempoolLocker != nil {
			go func() {
				// we won't get this until after the commit and we will acquire strictly after this commit phase has
				// ended (i.e. when Tendermint's BlockExecutor.Commit() returns
				app.mempoolLocker.Lock()
				// Prevent any mempool getting relocked while we unlock - we could just unlock immediately but if a new
				// commit starts gives goroutines blocked on checker a chance to progress before the next commit phase
				defer app.mempoolLocker.Unlock()
				app.checker.Unlock()
			}()
		} else {
			// If we have not be provided with access to the mempool lock
			app.checker.Unlock()
		}
	}()

	appHash, err := app.committer.Commit()
	if err != nil {
		panic(errors.Wrap(err, "Could not commit transactions in block to execution state"))

	}

	// Commit to our blockchain state
	err = app.blockchain.CommitBlock(time.Unix(int64(app.block.Header.Time), 0), app.block.Hash, appHash)
	if err != nil {
		panic(errors.Wrap(err, "could not commit block to blockchain state"))
	}

	err = app.checker.Reset()
	if err != nil {
		panic(errors.Wrap(err, "could not reset check cache during commit"))
	}

	// Perform a sanity check our block height
	if app.blockchain.LastBlockHeight() != uint64(app.block.Header.Height) {
		app.logger.InfoMsg("Burrow block height disagrees with Tendermint block height",
			structure.ScopeKey, "Commit()",
			"burrow_height", app.blockchain.LastBlockHeight(),
			"tendermint_height", app.block.Header.Height)

		panic(fmt.Errorf("burrow has recorded a block height of %v, "+
			"but Tendermint reports a block height of %v, and the two should agree",
			app.blockchain.LastBlockHeight(), app.block.Header.Height))
	}
	return abci_types.ResponseCommit{
		Data: appHash,
	}
}
