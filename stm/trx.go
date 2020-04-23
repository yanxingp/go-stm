// trx.go - Golang implementation of software transactional memory (STM).
// Author: Yanxing Pan (yanxingp)
// Date: April 2020
//
// Using lazy-versioning and a hybrid optimistic conflict detection.

package stm

import (
	"container/list"
	"sync/atomic"
)

// Global version-clock
var globalClock uint64

// Var - A STM variable (object).
type Var struct {
	// Address of the object's write-lock
	// 1 - locked, 0 - unlocked
	wlock uint64
	// Version number
	ver uint64
	// Address of data.
	// (Is this correct?)
	data interface{}
}

// acquireLock - Try to acquire the write-lock of a variable.
// Return value indicates whether the lock is successfully acquired.
func (v *Var) acquireLock() bool {
	// Make this bounded to avoid deadlock
	stop := 10000000
	for i := 0; i < stop; i++ {
		if atomic.CompareAndSwapUint64(&v.wlock, 0, 1) {
			return true
		}
	}
	return false
}

// releaseLock - Release the write-lock of a variable.
func (v *Var) releaseLock() {
	v.wlock = 0
}

// NewVar - Create a STM variable.
func NewVar(val interface{}) *Var {
	v := &Var{}
	v.data = val
	v.wlock = 0
	v.ver = globalClock
	return v
}

// Trx - A bookkeeper for the information of a transaction.
type Trx struct {
	// Indicate if this transaction is read-only
	isReadOnly bool

	// The current value of the global version clock
	rv uint64

	// Read-set and write-set as linked lists
	rset, wset *list.List

	success bool
}

// Read-set data type (load requests)
type rreq struct {
	v *Var
}

// Read-set data type (store requests)
type wreq struct {
	v   *Var
	val interface{}
}

// Load - Transactional load, return the value of STM variable
func (trx *Trx) Load(x *Var) interface{} {
	// If the transaction already aborted, simply skip the rest
	if !trx.success {
		return x.data
	}

	// If lock is not free, or variable's version number is greater,
	// abort and retry
	locked, ver := x.wlock, x.ver
	if locked == 1 || ver > trx.rv {
		trx.success = false
		return x.data
	}

	// Life is easier if the transaction is read-only:
	// No need to construct the read-set, simply validate version number
	if trx.isReadOnly {
		return x.data
	}

	// If it's a write transaction:
	// Append to read-set
	trx.rset.PushBack(rreq{v: x})

	// Check if the read address is in write-set,
	// to avoid read-write conflict
	for e := trx.wset.Front(); e != nil; e = e.Next() {
		if e.Value.(wreq).v == x {
			// Read from write-set directly
			// (Using type assertion here. Not sure if it's right..)
			return e.Value.(wreq).val
		}
	}

	return x.data
}

// Store - Transactional set, change the value of STM variable x to v.
func (trx *Trx) Store(x *Var, val interface{}) {
	// If the transaction already aborted, simply skip the rest
	if !trx.success {
		return
	}

	if trx.isReadOnly {
		// Doing store operation in read-only transactions.
		// The programmer is misusing the semantics.
		panic("Error: Store operation in read-only transaction.")
	}

	// Check if the variable is already in write-set
	found := false
	for e := trx.wset.Front(); e != nil; e = e.Next() {
		v := e.Value.(wreq).v
		if v == x {
			e.Value = wreq{v: v, val: val}
			found = true
			break
		}
	}

	// If not in write-set, append to write-set
	if !found {
		trx.wset.PushBack(wreq{v: x, val: val})
	}

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

		// Debug
		// fmt.Printf("write-set length = %v, read-set length = %v\n", trx.wset.Len(), trx.rset.Len())

		// If no conflict detected during speculative execution, continue
		// Otherwise, abort and retry the transaction
		if trx.success {
			// Lock the entire write-set
			lockedSet := list.New()
			acquired := true
			for e := trx.wset.Front(); e != nil; e = e.Next() {
				v := e.Value.(wreq).v

				// In case any of these lock is not successfully acquired
				// the transaction fails
				if !v.acquireLock() {
					acquired = false
					break
				}

				lockedSet.PushBack(v)
			}

			// If failed, need to release every acquired lock
			if !acquired {
				for e := lockedSet.Front(); e != nil; e = e.Next() {
					v := e.Value.(*Var)
					v.releaseLock()
				}
				trx.success = false
				// And retry the transaction
				continue
			}

			// Increment global version-lock
			// and store the new value as wv (write-version)
			wv := atomic.AddUint64(&globalClock, 1)

			// Validate the read-set
			// If rv + 1 = wv, this validation can be skipped
			// because no other transaction could have intervened
			if wv != trx.rv+1 {
				valid := true
				for e := trx.rset.Front(); e != nil; e = e.Next() {
					v := e.Value.(rreq).v
					locked, ver := v.wlock, v.ver
					if locked == 1 || ver > trx.rv {
						valid = false
						break
					}
				}

				// If failed, need to release every acquired lock
				if !valid {
					for e := lockedSet.Front(); e != nil; e = e.Next() {
						v := e.Value.(*Var)
						v.releaseLock()
					}
					trx.success = false
					// And retry the transaction
					continue
				}
			}

			// Commit and release the lock
			for e := trx.wset.Front(); e != nil; e = e.Next() {
				v := e.Value.(wreq).v
				val := e.Value.(wreq).val

				// Store the new value to the location
				v.data = val
				// Release the lock
				//  and change the variable's version number to wv
				v.releaseLock()
				v.ver = wv
			}

			// A transaction is successfully committed.
			trx.success = true
		}
	}

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
