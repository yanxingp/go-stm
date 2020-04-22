// trx.go - Golang implementation of software transactional memory (STM).
// Author: Yanxing Pan (yanxingp)
// Date: April 2020
//
// Using lazy-versioning and optimistic conflict detection.

package stm

import (
	"container/list"
)

type wlock uint64

type version uint64

// Global version-clock
var globalClock version

// Var - A STM variable (object).
type Var struct {
	// Address of the object's write-lock
	lock wlock
	// Address of data.
	// (Is this correct?)
	data interface{}
}

// decodeWlock and encodeWlock - Helper functions to use with wlock.
func decodeWlock(lock wlock) (bool, version) {
	return (lock >> 63) == 1, version(uint64(lock) & (1<<63 - 1))
}

func encodeWlock(locked bool, ver version) wlock {
	if (ver >> 63) == 1 {
		panic("Version number exceeded limit.")
	}
	if locked {
		return wlock(1<<63 | uint64(ver))
	}
	return wlock(ver)
}

// NewVar - Create a STM variable.
func NewVar(val interface{}) Var {
	v := Var{}
	v.data = val
	v.lock = encodeWlock(false, globalClock)
	return v
}

// Trx - A bookkeeper for the information of a transaction.
type Trx struct {
	// Indicate if this transaction is read-only
	isReadOnly bool

	// The current value of the global version clock
	rv version

	// Read-set and write-set as linked lists
	rset, wset *list.List

	success bool
}

// Read-set data type (load requests)
type lreq struct {
	v *Var
}

// Read-set data type (store requests)
type sreq struct {
	v   *Var
	val interface{}
}

// Load - Transactional load, return the value of STM variable
func (trx *Trx) Load(x *Var) interface{} {
	// If the transaction already aborted, simply skip the rest
	if !trx.success {
		return nil
	}

	// If lock is not free, or variable's version number is greater,
	// abort and retry
	locked, ver := decodeWlock(x.lock)
	if locked || ver > trx.rv {
		trx.success = false
		return nil
	}

	// Life is easier if the transaction is read-only
	// No need to construct the read-set, simply validate version number
	if trx.isReadOnly {
		return x.data
	}

	// If it's a write transaction

	// Append to read-set
	trx.rset.PushBack(lreq{v: x})

	// Check if the read address is in write-set,
	// to avoid read-write conflict
	for e := trx.wset.Front(); e != nil; e = e.Next() {
		if e == x {

		}
	}

	// TODO

	return nil
}

// Store - Transactional set, change the value of STM variable x to v.
func (trx *Trx) Store(x *Var, v interface{}) {
	// If the transaction already aborted, simply skip the rest
	if !trx.success {
		return
	}

	if trx.isReadOnly {
		// Doing store operation in read-only transactions.
		// The programmer is misusing the semantics.
		panic("Error: Store operation in read-only transaction.")
	}

	// TODO

}

// Atomically - Execute a transaction and return its return value
func Atomically(trans func(trx *Trx) interface{}) interface{} {
	// Initialize transaction bookkeeper with corresponding mode
	trx := NewTrx(false)
	var res interface{}

	// Keep retrying until the transaction is commited
	for !trx.success {
		// Reset bookkeeper's states
		trx.init()

		// Start speculative execution,
		// which means carry out instructions but doesn't modify memory values
		res = trans(trx)

		// If no conflict detected during speculative execution, continue
		// Otherwise, abort and retry the transaction
		if trx.success {
			// TODO
		}
	}

	// TODO

	return res
}

// ReadOnlyAtomically - Execute a read-only transaction
// Should be much more efficient
func ReadOnlyAtomically(trans func(trx *Trx) interface{}) interface{} {
	// Initialize transaction bookkeeper with corresponding mode
	trx := NewTrx(true)
	var res interface{}

	// Keep retrying until the transaction is commited
	for !trx.success {
		// Reset bookkeeper's states
		trx.init()

		// Start speculative execution,
		// which means carry out instructions but doesn't modify memory values
		res = trans(trx)

		// If no conflict detected during speculative execution, transaction commits
		// Otherwise, abort and retry the transaction
	}

	return res
}

// NewTrx - Return a transaction bookkeeper
func NewTrx(isReadOnly bool) *Trx {
	trx := &Trx{}
	trx.isReadOnly = isReadOnly
	trx.rset, trx.wset = list.New(), list.New()

	// Make sure the loop can start
	trx.success = false

	return trx
}

// init - Initialize or reset the transaction bookkeeper's states
func (trx *Trx) init() {
	trx.rset.Init()
	trx.wset.Init()

	// Sample the global version clock
	trx.rv = globalClock

	trx.success = true
}
