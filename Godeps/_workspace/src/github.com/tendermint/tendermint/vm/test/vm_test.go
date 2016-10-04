package vm

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
	"time"

	"errors"

	. "github.com/eris-ltd/eris-db/Godeps/_workspace/src/github.com/tendermint/tendermint/common"
	"github.com/eris-ltd/eris-db/Godeps/_workspace/src/github.com/tendermint/tendermint/events"
	ptypes "github.com/eris-ltd/eris-db/Godeps/_workspace/src/github.com/tendermint/tendermint/permission/types"
	"github.com/eris-ltd/eris-db/Godeps/_workspace/src/github.com/tendermint/tendermint/types"
	. "github.com/eris-ltd/eris-db/Godeps/_workspace/src/github.com/tendermint/tendermint/vm"
	"github.com/stretchr/testify/assert"
)

func init() {
	SetDebug(true)
}

func newAppState() *FakeAppState {
	fas := &FakeAppState{
		accounts: make(map[string]*Account),
		storage:  make(map[string]Word256),
	}
	// For default permissions
	fas.accounts[ptypes.GlobalPermissionsAddress256.String()] = &Account{
		Permissions: ptypes.DefaultAccountPermissions,
	}
	return fas
}

func newParams() Params {
	return Params{
		BlockHeight: 0,
		BlockHash:   Zero256,
		BlockTime:   0,
		GasLimit:    0,
	}
}

func makeBytes(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}

// Runs a basic loop
func TestVM(t *testing.T) {
	ourVm := NewVM(newAppState(), newParams(), Zero256, nil)

	// Create accounts
	account1 := &Account{
		Address: Int64ToWord256(100),
	}
	account2 := &Account{
		Address: Int64ToWord256(101),
	}

	var gas int64 = 100000
	N := []byte{0x0f, 0x0f}
	// Loop N times
	code := []byte{0x60, 0x00, 0x60, 0x20, 0x52, 0x5B, byte(0x60 + len(N) - 1)}
	code = append(code, N...)
	code = append(code, []byte{0x60, 0x20, 0x51, 0x12, 0x15, 0x60, byte(0x1b + len(N)), 0x57, 0x60, 0x01, 0x60, 0x20, 0x51, 0x01, 0x60, 0x20, 0x52, 0x60, 0x05, 0x56, 0x5B}...)
	start := time.Now()
	output, err := ourVm.Call(account1, account2, code, []byte{}, 0, &gas)
	fmt.Printf("Output: %v Error: %v\n", output, err)
	fmt.Println("Call took:", time.Since(start))
	if err != nil {
		t.Fatal(err)
	}
}

func TestJumpErr(t *testing.T) {
	ourVm := NewVM(newAppState(), newParams(), Zero256, nil)

	// Create accounts
	account1 := &Account{
		Address: Int64ToWord256(100),
	}
	account2 := &Account{
		Address: Int64ToWord256(101),
	}

	var gas int64 = 100000
	code := []byte{0x60, 0x10, 0x56} // jump to position 16, a clear failure
	var output []byte
	var err error
	ch := make(chan struct{})
	go func() {
		output, err = ourVm.Call(account1, account2, code, []byte{}, 0, &gas)
		ch <- struct{}{}
	}()
	tick := time.NewTicker(time.Second * 2)
	select {
	case <-tick.C:
		t.Fatal("VM ended up in an infinite loop from bad jump dest (it took too long!)")
	case <-ch:
		if err == nil {
			t.Fatal("Expected invalid jump dest err")
		}
	}
}

// Tests the code for a subcurrency contract compiled by serpent
func TestSubcurrency(t *testing.T) {
	st := newAppState()
	// Create accounts
	account1 := &Account{
		Address: LeftPadWord256(makeBytes(20)),
	}
	account2 := &Account{
		Address: LeftPadWord256(makeBytes(20)),
	}
	st.accounts[account1.Address.String()] = account1
	st.accounts[account2.Address.String()] = account2

	ourVm := NewVM(st, newParams(), Zero256, nil)

	var gas int64 = 1000
	code_parts := []string{"620f42403355",
		"7c0100000000000000000000000000000000000000000000000000000000",
		"600035046315cf268481141561004657",
		"6004356040526040515460605260206060f35b63693200ce81141561008757",
		"60043560805260243560a052335460c0523360e05260a05160c05112151561008657",
		"60a05160c0510360e0515560a0516080515401608051555b5b505b6000f3"}
	code, _ := hex.DecodeString(strings.Join(code_parts, ""))
	fmt.Printf("Code: %x\n", code)
	data, _ := hex.DecodeString("693200CE0000000000000000000000004B4363CDE27C2EB05E66357DB05BC5C88F850C1A0000000000000000000000000000000000000000000000000000000000000005")
	output, err := ourVm.Call(account1, account2, code, data, 0, &gas)
	fmt.Printf("Output: %v Error: %v\n", output, err)
	if err != nil {
		t.Fatal(err)
	}
}

// Test sending tokens from a contract to another account
func TestSendCall(t *testing.T) {
	fakeAppState := newAppState()
	ourVm := NewVM(fakeAppState, newParams(), Zero256, nil)

	// Create accounts
	account1 := &Account{
		Address: Int64ToWord256(100),
	}
	account2 := &Account{
		Address: Int64ToWord256(101),
	}
	account3 := &Account{
		Address: Int64ToWord256(102),
	}

	// account1 will call account2 which will trigger CALL opcode to account3
	addr := account3.Address.Postfix(20)
	contractCode := callContractCode(addr)

	//----------------------------------------------
	// account2 has insufficient balance, should fail
	_, err := runVMWaitError(ourVm, account1, account2, addr, contractCode, 1000)
	assert.Error(t, err, "Expected insufficient balance error")

	//----------------------------------------------
	// give account2 sufficient balance, should pass
	account2.Balance = 100000
	_, err = runVMWaitError(ourVm, account1, account2, addr, contractCode, 1000)
	assert.NoError(t, err, "Should have sufficient balance")

	//----------------------------------------------
	// insufficient gas, should fail

	account2.Balance = 100000
	_, err = runVMWaitError(ourVm, account1, account2, addr, contractCode, 100)
	assert.Error(t, err, "Expected insufficient gas error")
}

// This test was introduced to cover an issues exposed in our handling of the
// gas limit passed from caller to callee on various forms of CALL
// this ticket gives some background: https://github.com/eris-ltd/eris-pm/issues/212
// The idea of this test is to implement a simple DelegateCall in EVM code
// We first run the DELEGATECALL with _just_ enough gas expecting a simple return,
// and then run it with 1 gas unit less, expecting a failure
func TestDelegateCallGas(t *testing.T) {
	appState := newAppState()
	ourVm := NewVM(appState, newParams(), Zero256, nil)

	inOff := 0
	inSize := 0 // no call data
	retOff := 0
	retSize := 32
	calleeReturnValue := int64(20)

	// DELEGATECALL(retSize, refOffset, inSize, inOffset, addr, gasLimit)
	// 6 pops
	delegateCallCost := GasStackOp * 6
	// 1 push
	gasCost := GasStackOp
	// 2 pops, 1 push
	subCost := GasStackOp * 3
	pushCost := GasStackOp

	costBetweenGasAndDelegateCall := gasCost + subCost + delegateCallCost + pushCost

	// Do a simple operation using 1 gas unit
	calleeAccount, calleeAddress := makeAccountWithCode(appState, "callee",
		bytecode(PUSH1, calleeReturnValue, return1()))

	// Here we split up the caller code so we can make a DELEGATE call with
	// different amounts of gas. The value we sandwich in the middle is the amount
	// we subtract from the available gas (that the caller has available), so:
	// code := bytecode(callerCodePrefix, <amount to subtract from GAS> , callerCodeSuffix)
	// gives us the code to make the call
	callerCodePrefix := bytecode(PUSH1, retSize, PUSH1, retOff, PUSH1, inSize,
		PUSH1, inOff, PUSH20, calleeAddress, PUSH1)
	callerCodeSuffix := bytecode(GAS, SUB, DELEGATECALL, returnWord())

	// Perform a delegate call
	callerAccount, _ := makeAccountWithCode(appState, "caller",
		bytecode(callerCodePrefix,
			// Give just enough gas to make the DELEGATECALL
			costBetweenGasAndDelegateCall,
			callerCodeSuffix))

	// Should pass
	output, err := runVMWaitError(ourVm, callerAccount, calleeAccount, calleeAddress,
		callerAccount.Code, 100)
	assert.NoError(t, err, "Should have sufficient funds for call")
	assert.Equal(t, Int64ToWord256(calleeReturnValue).Bytes(), output)

	callerAccount.Code = bytecode(callerCodePrefix,
		// Shouldn't be enough gas to make call
		costBetweenGasAndDelegateCall-1,
		callerCodeSuffix)

	// Should fail
	_, err = runVMWaitError(ourVm, callerAccount, calleeAccount, calleeAddress,
		callerAccount.Code, 100)
	assert.Error(t, err, "Should have insufficient funds for call")
}

// Store the top element of the stack (which is a 32-byte word) in memory
// and return it. Useful for a simple return value.
func return1() []byte {
	return bytecode(PUSH1, 0, MSTORE, returnWord())
}

func returnWord() []byte {
	// PUSH1 => return size, PUSH1 => return offset, RETURN
	return bytecode(PUSH1, 32, PUSH1, 0, RETURN)
}

func makeAccountWithCode(appState AppState, name string,
	code []byte) (*Account, []byte) {
	account := &Account{
		Address: LeftPadWord256([]byte(name)),
		Balance: 9999999,
		Code:    code,
		Nonce:   0,
	}
	account.Code = code
	appState.UpdateAccount(account)
	// Sanity check
	address := new([20]byte)
	for i, b := range account.Address.Postfix(20) {
		address[i] = b
	}
	return account, address[:]
}

// Subscribes to an AccCall, runs the vm, returns the output any direct exception
// and then waits for any exceptions transmitted by EventData in the AccCall
// event (in the case of no direct error from call we will block waiting for
// at least 1 AccCall event)
func runVMWaitError(ourVm *VM, caller, callee *Account, subscribeAddr,
	contractCode []byte, gas int64) (output []byte, err error) {
	eventCh := make(chan types.EventData)
	output, err = runVM(eventCh, ourVm, caller, callee, subscribeAddr,
		contractCode, gas)
	if err != nil {
		return
	}
	msg := <-eventCh
	var errString string
	switch ev := msg.(type) {
	case types.EventDataTx:
		errString = ev.Exception
	case types.EventDataCall:
		errString = ev.Exception
	}

	if errString != "" {
		err = errors.New(errString)
	}
	return
}

// Subscribes to an AccCall, runs the vm, returns the output and any direct
// exception
func runVM(eventCh chan types.EventData, ourVm *VM, caller, callee *Account,
	subscribeAddr, contractCode []byte, gas int64) ([]byte, error) {

	// we need to catch the event from the CALL to check for exceptions
	evsw := events.NewEventSwitch()
	evsw.Start()
	fmt.Printf("subscribe to %x\n", subscribeAddr)
	evsw.AddListenerForEvent("test", types.EventStringAccCall(subscribeAddr),
		func(msg types.EventData) {
			eventCh <- msg
		})
	evc := events.NewEventCache(evsw)
	ourVm.SetFireable(evc)
	start := time.Now()
	output, err := ourVm.Call(caller, callee, contractCode, []byte{}, 0, &gas)
	fmt.Printf("Output: %v Error: %v\n", output, err)
	fmt.Println("Call took:", time.Since(start))
	go func() { evc.Flush() }()
	return output, err
}

// this is code to call another contract (hardcoded as addr)
func callContractCode(addr []byte) []byte {
	gas1, gas2 := byte(0x1), byte(0x1)
	value := byte(0x69)
	inOff, inSize := byte(0x0), byte(0x0) // no call data
	retOff, retSize := byte(0x0), byte(0x20)
	// this is the code we want to run (send funds to an account and return)
	return bytecode(PUSH1, retSize, PUSH1, retOff, PUSH1, inSize, PUSH1,
		inOff, PUSH1, value, PUSH20, addr, PUSH2, gas1, gas2, CALL, PUSH1, retSize,
		PUSH1, retOff, RETURN)
}

func TestBytecode(t *testing.T) {
	assert.Equal(t,
		bytecode(1, 2, 3, 4, 5, 6),
		bytecode(1, 2, 3, bytecode(4, 5, 6)))
	assert.Equal(t,
		bytecode(1, 2, 3, 4, 5, 6, 7, 8),
		bytecode(1, 2, 3, bytecode(4, bytecode(5), 6), 7, 8))
	assert.Equal(t,
		bytecode(PUSH1, 2),
		bytecode(byte(PUSH1), 0x02))
	assert.Equal(t,
		[]byte{},
		bytecode(bytecode(bytecode())))

	contractAccount := &Account{Address: Int64ToWord256(102)}
	addr := contractAccount.Address.Postfix(20)
	gas1, gas2 := byte(0x1), byte(0x1)
	value := byte(0x69)
	inOff, inSize := byte(0x0), byte(0x0) // no call data
	retOff, retSize := byte(0x0), byte(0x20)
	contractCodeBytecode := bytecode(PUSH1, retSize, PUSH1, retOff, PUSH1, inSize, PUSH1,
		inOff, PUSH1, value, PUSH20, addr, PUSH2, gas1, gas2, CALL, PUSH1, retSize,
		PUSH1, retOff, RETURN)
	contractCode := []byte{0x60, retSize, 0x60, retOff, 0x60, inSize, 0x60, inOff, 0x60, value, 0x73}
	contractCode = append(contractCode, addr...)
	contractCode = append(contractCode, []byte{0x61, gas1, gas2, 0xf1, 0x60, 0x20, 0x60, 0x0, 0xf3}...)
	assert.Equal(t, contractCode, contractCodeBytecode)
}

func TestConcat(t *testing.T) {
	assert.Equal(t,
		[]byte{0x01, 0x02, 0x03, 0x04},
		concat([]byte{0x01, 0x02}, []byte{0x03, 0x04}))
}

// Convenience function to allow us to mix bytes, ints, and OpCodes that
// represent bytes in an EVM assembly code to make assembly more readable.
// Also allows us to splice together assembly
// fragments because any []byte arguments are flattened in the result.
func bytecode(bytelikes ...interface{}) []byte {
	bytes := make([]byte, len(bytelikes))
	for i, bytelike := range bytelikes {
		switch b := bytelike.(type) {
		case byte:
			bytes[i] = b
		case OpCode:
			bytes[i] = byte(b)
		case int:
			bytes[i] = byte(b)
			if int(bytes[i]) != b {
				panic(fmt.Sprintf("The int %v does not fit inside a byte", b))
			}
		case int64:
			bytes[i] = byte(b)
			if int64(bytes[i]) != b {
				panic(fmt.Sprintf("The int64 %v does not fit inside a byte", b))
			}
		case []byte:
			// splice
			return concat(bytes[:i], b, bytecode(bytelikes[i+1:]...))
		default:
			panic("Only byte-like codes (and []byte sequences) can be used to form bytecode")
		}
	}
	return bytes
}

func concat(bss ...[]byte) []byte {
	offset := 0
	for _, bs := range bss {
		offset += len(bs)
	}
	bytes := make([]byte, offset)
	offset = 0
	for _, bs := range bss {
		for i, b := range bs {
			bytes[offset+i] = b
		}
		offset += len(bs)
	}
	return bytes
}

/*
	// infinite loop
	code := []byte{0x5B, 0x60, 0x00, 0x56}
	// mstore
	code := []byte{0x60, 0x00, 0x60, 0x20}
	// mstore, mload
	code := []byte{0x60, 0x01, 0x60, 0x20, 0x52, 0x60, 0x20, 0x51}
*/
