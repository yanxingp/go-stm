// benchmark_test.go - Benchmarking the performance of stm implementation 
// under different workloads and in comparision with the mutex lock.
// Insipred by work of decellion 
// (https://github.com/decillion/go-stm/blob/master/benchmark/bench_test.go)
// Author: Yanxing Pan (yanxingp)
// Date: April 2020

package stm

import (
	"sync"
	"sync/atomic"
	"testing"
)

// BenchmarkRead - All reads no contention.
func BenchmarkReadOnlyMutex(b *testing.B) {
	x, y, z := 1, 2, 3
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mutex := &sync.Mutex{}
			mutex.Lock()
			_ = x
			_ = y
			_ = z
			mutex.Unlock()
		}
	})
}

func BenchmarkReadOnlyTrx(b *testing.B) {
	x := NewVar(0)
	y := NewVar(0)

	load := func(trx *Trx) interface{} {
		trx.Load(x)
		trx.Load(y)
		return nil
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Atomically(load)
		}
	})
}

func BenchmarkReadOnlyTrxRO(b *testing.B) {
	x := NewVar(0)
	y := NewVar(0)

	load := func(trx *Trx) interface{} {
		trx.Load(x)
		trx.Load(y)
		return nil
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ReadOnlyAtomically(load)
		}
	})
}

func BenchmarkRead90Write10Mutex(b *testing.B) {
	var x, y int
	mu := &sync.RWMutex{}
	b.ResetTimer()

	var counter uint64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := atomic.AddUint64(&counter, 1)
			if id%10 == 0 {
				mu.Lock()
				x++
				y++
				mu.Unlock()
			} else {
				mu.RLock()
				_ = x
				_ = y
				mu.RUnlock()
			}
		}
	})
}

func BenchmarkRead90Write10Trx(b *testing.B) {
	x := NewVar(0)
	y := NewVar(0)

	inc := func(trx *Trx) interface{} {
		trx.Store(x, trx.Load(x).(int)+1)
		trx.Store(y, trx.Load(y).(int)+1)
		return nil
	}
	load := func(trx *Trx) interface{} {
		trx.Load(x)
		trx.Load(y)
		return nil
	}
	b.ResetTimer()

	var counter uint64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := atomic.AddUint64(&counter, 1)
			if id%10 == 0 {
				Atomically(inc)
			} else {
				Atomically(load)
			}
		}
	})
}

func BenchmarkRead70Write30Mutex(b *testing.B) {
	var x, y int
	mu := &sync.RWMutex{}
	b.ResetTimer()

	var counter uint64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := atomic.AddUint64(&counter, 1)
			if id%10 < 3 {
				mu.Lock()
				x++
				y++
				mu.Unlock()
			} else {
				mu.RLock()
				_ = x
				_ = y
				mu.RUnlock()
			}
		}
	})
}

func BenchmarkRead70Write30Trx(b *testing.B) {
	x := NewVar(0)
	y := NewVar(0)

	inc := func(trx *Trx) interface{} {
		trx.Store(x, trx.Load(x).(int)+1)
		trx.Store(y, trx.Load(y).(int)+1)
		return nil
	}
	load := func(trx *Trx) interface{} {
		trx.Load(x)
		trx.Load(y)
		return nil
	}
	b.ResetTimer()

	var counter uint64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := atomic.AddUint64(&counter, 1)
			if id%10 < 3 {
				Atomically(inc)
			} else {
				Atomically(load)
			}
		}
	})
}

func BenchmarkRead50Write50Trx(b *testing.B) {
	x := NewVar(0)
	y := NewVar(0)

	inc := func(trx *Trx) interface{} {
		trx.Store(x, trx.Load(x).(int)+1)
		trx.Store(y, trx.Load(y).(int)+1)
		return nil
	}
	load := func(trx *Trx) interface{} {
		trx.Load(x)
		trx.Load(y)
		return nil
	}
	b.ResetTimer()

	var counter uint64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := atomic.AddUint64(&counter, 1)
			if id%10 < 5 {
				Atomically(inc)
			} else {
				Atomically(load)
			}
		}
	})
}
