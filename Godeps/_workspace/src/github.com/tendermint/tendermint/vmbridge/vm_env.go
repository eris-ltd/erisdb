// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// This file was forked from go-ethereum and modified to suit erisdb

package vmbridge

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	mintcommon "github.com/eris-ltd/eris-db/Godeps/_workspace/src/github.com/tendermint/tendermint/common"
	"github.com/eris-ltd/eris-db/Godeps/_workspace/src/github.com/tendermint/tendermint/events"
	mintvm "github.com/eris-ltd/eris-db/Godeps/_workspace/src/github.com/tendermint/tendermint/vm"
)

type VMEnv struct {
	state  *Database
	params mintvm.Params
	origin common.Address
	value  *big.Int
	evc    events.Fireable

	depth int
	typ   vm.Type
	// structured logging
	logs []vm.StructLog
}

func NewEnv(appState mintvm.AppState, params mintvm.Params, origin mintcommon.Word256, value int64) *VMEnv {
	return &VMEnv{
		state:  NewDatabase(appState),
		params: params,
		origin: common.BytesToAddress(origin.Postfix(20)),
		value:  big.NewInt(value),
		typ:    vm.StdVmTy,
	}
}

func (env *VMEnv) SetFireable(evc events.Fireable) {
	env.evc = evc
}

func (self *VMEnv) Origin() common.Address   { return self.origin }
func (self *VMEnv) BlockNumber() *big.Int    { return big.NewInt(self.params.BlockHeight) }
func (self *VMEnv) Coinbase() common.Address { return common.Address{} } // could be proposer one day
func (self *VMEnv) Time() *big.Int           { return big.NewInt(self.params.BlockTime) }
func (self *VMEnv) Difficulty() *big.Int     { return big.NewInt(0) }
func (self *VMEnv) GasLimit() *big.Int       { return big.NewInt(self.params.GasLimit) }
func (self *VMEnv) Value() *big.Int          { return self.value }
func (self *VMEnv) Db() vm.Database          { return self.state }
func (self *VMEnv) Depth() int               { return self.depth }
func (self *VMEnv) SetDepth(i int)           { self.depth = i }
func (self *VMEnv) VmType() vm.Type          { return self.typ }
func (self *VMEnv) SetVmType(t vm.Type)      { self.typ = t }

//---------
// need the chain

func (self *VMEnv) GetHash(n uint64) common.Hash {
	// TODO
	/*
		for block := self.chain.GetBlock(self.header.ParentHash); block != nil; block = self.chain.GetBlock(block.ParentHash()) {
			if block.NumberU64() == n {
				return block.Hash()
			}
		}
	*/

	return common.Hash{}
}

//------------------------------
// logs logs logs

func (self *VMEnv) AddLog(log *vm.Log) {
	// TODO: something about logs?!
	// self.state.AddLog(log)

	// TODO: fire event
	/*
		if vm.evc != nil {
			eventID := types.EventStringLogEvent(callee.Address.Postfix(20))
			fmt.Printf("eventID: %s\n", eventID)
			log := types.EventDataLog{
				callee.Address,
				topics,
				data,
				vm.params.BlockHeight,
			}
			vm.evc.FireEvent(eventID, log)
		}
	*/
}

func (self *VMEnv) StructLogs() []vm.StructLog {
	return self.logs
}

func (self *VMEnv) AddStructLog(log vm.StructLog) {
	self.logs = append(self.logs, log)
}

//-------------------------------
// transfers

func (self *VMEnv) CanTransfer(from common.Address, balance *big.Int) bool {
	return self.state.GetBalance(from).Cmp(balance) >= 0
}

func (self *VMEnv) Transfer(from, to vm.Account, amount *big.Int) {
	from.SubBalance(amount)
	to.AddBalance(amount)
}

//-------------------------------

func (self *VMEnv) Call(me vm.ContractRef, addr common.Address, data []byte, gas, price, value *big.Int) ([]byte, error) {
	return Call(self, me, addr, data, gas, price, value)
}
func (self *VMEnv) CallCode(me vm.ContractRef, addr common.Address, data []byte, gas, price, value *big.Int) ([]byte, error) {
	return CallCode(self, me, addr, data, gas, price, value)
}

func (self *VMEnv) DelegateCall(me vm.ContractRef, addr common.Address, data []byte, gas, price *big.Int) ([]byte, error) {
	return DelegateCall(self, me, addr, data, gas, price)
}

func (self *VMEnv) Create(me vm.ContractRef, data []byte, gas, price, value *big.Int) ([]byte, common.Address, error) {
	return Create(self, me, data, gas, price, value)
}

//----------------------------------
// snapshots not needed

func (self *VMEnv) MakeSnapshot() vm.Database {
	// not used.
	// could just return a copy of the underlying state.
	// we use maps for intermediate state
	return nil
}

func (self *VMEnv) SetSnapshot(copy vm.Database) {
	// we should never actually need to call this
	// see MakeSnapshot.
}
