package mock

// Basic imports
import (
	"bytes"
	"encoding/hex"
	"net/http"
	"os"
	"runtime"
	"testing"

	"github.com/eris-ltd/eris-db/account"
	"github.com/eris-ltd/eris-db/config"
	core "github.com/eris-ltd/eris-db/core"
	core_types "github.com/eris-ltd/eris-db/core/types"
	"github.com/eris-ltd/eris-db/rpc"
	"github.com/eris-ltd/eris-db/server"
	td "github.com/eris-ltd/eris-db/test/testdata/testdata"
	"github.com/eris-ltd/eris-db/txs"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/log15"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log15.Root().SetHandler(log15.LvlFilterHandler(
		log15.LvlWarn,
		log15.StreamHandler(os.Stdout, log15.TerminalFormat()),
	))
	gin.SetMode(gin.ReleaseMode)
}

type MockSuite struct {
	suite.Suite
	baseDir      string
	serveProcess *server.ServeProcess
	codec        rpc.Codec
	sUrl         string
	testData     *td.TestData
}

func (this *MockSuite) SetupSuite() {
	gin.SetMode(gin.ReleaseMode)
	// Load the supporting objects.
	testData := td.LoadTestData()
	pipe := NewMockPipe(testData)
	codec := &core_types.TCodec{}
	evtSubs := core.NewEventSubscriptions(pipe.Events())
	// The server
	restServer := core.NewRestServer(codec, pipe, evtSubs)
	sConf := config.DefaultServerConfig()
	sConf.Bind.Port = 31402
	// Create a server process.
	proc := server.NewServeProcess(&sConf, restServer)
	err := proc.Start()
	if err != nil {
		panic(err)
	}
	this.serveProcess = proc
	this.codec = core_types.NewTCodec()
	this.testData = testData
	this.sUrl = "http://localhost:31402"
}

func (this *MockSuite) TearDownSuite() {
	sec := this.serveProcess.StopEventChannel()
	this.serveProcess.Stop(0)
	<-sec
}

// ********************************************* Accounts *********************************************

func (this *MockSuite) TestGetAccounts() {
	resp := this.get("/accounts")
	ret := &core_types.AccountList{}
	errD := this.codec.Decode(ret, resp.Body)
	this.NoError(errD)
	this.Equal(ret, this.testData.GetAccounts.Output)
}

func (this *MockSuite) TestGetAccount() {
	addr := hex.EncodeToString(this.testData.GetAccount.Input.Address)
	resp := this.get("/accounts/" + addr)
	ret := &account.Account{}
	errD := this.codec.Decode(ret, resp.Body)
	this.NoError(errD)
	this.Equal(ret, this.testData.GetAccount.Output)
}

func (this *MockSuite) TestGetStorage() {
	addr := hex.EncodeToString(this.testData.GetStorage.Input.Address)
	resp := this.get("/accounts/" + addr + "/storage")
	ret := &core_types.Storage{}
	errD := this.codec.Decode(ret, resp.Body)
	this.NoError(errD)
	this.Equal(ret, this.testData.GetStorage.Output)
}

func (this *MockSuite) TestGetStorageAt() {
	addr := hex.EncodeToString(this.testData.GetStorageAt.Input.Address)
	key := hex.EncodeToString(this.testData.GetStorageAt.Input.Key)
	resp := this.get("/accounts/" + addr + "/storage/" + key)
	ret := &core_types.StorageItem{}
	errD := this.codec.Decode(ret, resp.Body)
	this.NoError(errD)
	this.Equal(ret, this.testData.GetStorageAt.Output)
}

// ********************************************* Blockchain *********************************************

func (this *MockSuite) TestGetBlockchainInfo() {
	resp := this.get("/blockchain")
	ret := &core_types.BlockchainInfo{}
	errD := this.codec.Decode(ret, resp.Body)
	this.NoError(errD)
	this.Equal(ret, this.testData.GetBlockchainInfo.Output)
}

func (this *MockSuite) TestGetChainId() {
	resp := this.get("/blockchain/chain_id")
	ret := &core_types.ChainId{}
	errD := this.codec.Decode(ret, resp.Body)
	this.NoError(errD)
	this.Equal(ret, this.testData.GetChainId.Output)
}

func (this *MockSuite) TestGetGenesisHash() {
	resp := this.get("/blockchain/genesis_hash")
	ret := &core_types.GenesisHash{}
	errD := this.codec.Decode(ret, resp.Body)
	this.NoError(errD)
	this.Equal(ret, this.testData.GetGenesisHash.Output)
}

func (this *MockSuite) TestLatestBlockHeight() {
	resp := this.get("/blockchain/latest_block_height")
	ret := &core_types.LatestBlockHeight{}
	errD := this.codec.Decode(ret, resp.Body)
	this.NoError(errD)
	this.Equal(ret, this.testData.GetLatestBlockHeight.Output)
}

func (this *MockSuite) TestBlocks() {
	resp := this.get("/blockchain/blocks")
	ret := &core_types.Blocks{}
	errD := this.codec.Decode(ret, resp.Body)
	this.NoError(errD)
	this.Equal(ret, this.testData.GetBlocks.Output)
}

// ********************************************* Consensus *********************************************

func (this *MockSuite) TestGetConsensusState() {
	resp := this.get("/consensus")
	ret := &core_types.ConsensusState{}
	errD := this.codec.Decode(ret, resp.Body)
	this.NoError(errD)
	ret.StartTime = ""
	this.Equal(ret, this.testData.GetConsensusState.Output)
}

func (this *MockSuite) TestGetValidators() {
	resp := this.get("/consensus/validators")
	ret := &core_types.ValidatorList{}
	errD := this.codec.Decode(ret, resp.Body)
	this.NoError(errD)
	this.Equal(ret, this.testData.GetValidators.Output)
}

// ********************************************* NameReg *********************************************

func (this *MockSuite) TestGetNameRegEntry() {
	resp := this.get("/namereg/" + this.testData.GetNameRegEntry.Input.Name)
	ret := &txs.NameRegEntry{}
	errD := this.codec.Decode(ret, resp.Body)
	this.NoError(errD)
	this.Equal(ret, this.testData.GetNameRegEntry.Output)
}

func (this *MockSuite) TestGetNameRegEntries() {
	resp := this.get("/namereg")
	ret := &core_types.ResultListNames{}
	errD := this.codec.Decode(ret, resp.Body)
	this.NoError(errD)
	this.Equal(ret, this.testData.GetNameRegEntries.Output)
}

// ********************************************* Network *********************************************

func (this *MockSuite) TestGetNetworkInfo() {
	resp := this.get("/network")
	ret := &core_types.NetworkInfo{}
	errD := this.codec.Decode(ret, resp.Body)
	this.NoError(errD)
	this.Equal(ret, this.testData.GetNetworkInfo.Output)
}

func (this *MockSuite) TestGetClientVersion() {
	resp := this.get("/network/client_version")
	ret := &core_types.ClientVersion{}
	errD := this.codec.Decode(ret, resp.Body)
	this.NoError(errD)
	this.Equal(ret, this.testData.GetClientVersion.Output)
}

func (this *MockSuite) TestGetMoniker() {
	resp := this.get("/network/moniker")
	ret := &core_types.Moniker{}
	errD := this.codec.Decode(ret, resp.Body)
	this.NoError(errD)
	this.Equal(ret, this.testData.GetMoniker.Output)
}

func (this *MockSuite) TestIsListening() {
	resp := this.get("/network/listening")
	ret := &core_types.Listening{}
	errD := this.codec.Decode(ret, resp.Body)
	this.NoError(errD)
	this.Equal(ret, this.testData.IsListening.Output)
}

func (this *MockSuite) TestGetListeners() {
	resp := this.get("/network/listeners")
	ret := &core_types.Listeners{}
	errD := this.codec.Decode(ret, resp.Body)
	this.NoError(errD)
	this.Equal(ret, this.testData.GetListeners.Output)
}

func (this *MockSuite) TestGetPeers() {
	resp := this.get("/network/peers")
	ret := []*core_types.Peer{}
	errD := this.codec.Decode(ret, resp.Body)
	this.NoError(errD)
	this.Equal(ret, this.testData.GetPeers.Output)
}

/*
func (this *MockSuite) TestGetPeer() {
	addr := this.testData.GetPeer.Input.Address
	resp := this.get("/network/peer/" + addr)
	ret := []*core_types.Peer{}
	errD := this.codec.Decode(ret, resp.Body)
	this.NoError(errD)
	this.Equal(ret, this.testData.GetPeers.Output)
}
*/

// ********************************************* Transactions *********************************************

func (this *MockSuite) TestTransactCreate() {
	resp := this.postJson("/unsafe/txpool", this.testData.TransactCreate.Input)
	ret := &core_types.Receipt{}
	errD := this.codec.Decode(ret, resp.Body)
	this.NoError(errD)
	this.Equal(ret, this.testData.TransactCreate.Output)
}

func (this *MockSuite) TestTransact() {
	resp := this.postJson("/unsafe/txpool", this.testData.Transact.Input)
	ret := &core_types.Receipt{}
	errD := this.codec.Decode(ret, resp.Body)
	this.NoError(errD)
	this.Equal(ret, this.testData.Transact.Output)
}

func (this *MockSuite) TestTransactNameReg() {
	resp := this.postJson("/unsafe/namereg/txpool", this.testData.TransactNameReg.Input)
	ret := &core_types.Receipt{}
	errD := this.codec.Decode(ret, resp.Body)
	this.NoError(errD)
	this.Equal(ret, this.testData.TransactNameReg.Output)
}

func (this *MockSuite) TestGetUnconfirmedTxs() {
	resp := this.get("/txpool")
	ret := &core_types.UnconfirmedTxs{}
	errD := this.codec.Decode(ret, resp.Body)
	this.NoError(errD)
	this.Equal(ret, this.testData.GetUnconfirmedTxs.Output)
}

func (this *MockSuite) TestCallCode() {
	resp := this.postJson("/codecalls", this.testData.CallCode.Input)
	ret := &core_types.Call{}
	errD := this.codec.Decode(ret, resp.Body)
	this.NoError(errD)
	this.Equal(ret, this.testData.CallCode.Output)
}

func (this *MockSuite) TestCall() {
	resp := this.postJson("/calls", this.testData.Call.Input)
	ret := &core_types.Call{}
	errD := this.codec.Decode(ret, resp.Body)
	this.NoError(errD)
	this.Equal(ret, this.testData.CallCode.Output)
}

// ********************************************* Utilities *********************************************

func (this *MockSuite) get(endpoint string) *http.Response {
	resp, errG := http.Get(this.sUrl + endpoint)
	this.NoError(errG)
	this.Equal(200, resp.StatusCode)
	return resp
}

func (this *MockSuite) postJson(endpoint string, v interface{}) *http.Response {
	bts, errE := this.codec.EncodeBytes(v)
	this.NoError(errE)
	resp, errP := http.Post(this.sUrl+endpoint, "application/json", bytes.NewBuffer(bts))
	this.NoError(errP)
	this.Equal(200, resp.StatusCode)
	return resp
}

// ********************************************* Entrypoint *********************************************

func TestMockSuite(t *testing.T) {
	suite.Run(t, &MockSuite{})
}
