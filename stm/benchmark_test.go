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

	de "github.com/decillion/go-stm"
	lu "github.com/lukechampine/stm"
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

func BenchmarkRead90Write10Decillion(b *testing.B) {
	x := de.New(0)
	y := de.New(0)

	inc := func(rec *de.TRec) interface{} {
		rec.Store(x, rec.Load(x).(int)+1)
		rec.Store(y, rec.Load(y).(int)+1)
		return nil
	}
	load := func(rec *de.TRec) interface{} {
		rec.Load(x)
		rec.Load(y)
		return nil
	}
	b.ResetTimer()

	var counter uint64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := atomic.AddUint64(&counter, 1)
			if id%10 == 0 {
				de.Atomically(inc)
			} else {
				de.Atomically(load)
			}
		}
	})
}

func BenchmarkRead90Write10Lukechampine(b *testing.B) {
	x := lu.NewVar(0)
	y := lu.NewVar(0)

	inc := func(tx *lu.Tx) {
		tx.Set(x, tx.Get(x).(int)+1)
		tx.Set(y, tx.Get(y).(int)+1)
	}
	load := func(tx *lu.Tx) {
		tx.Get(x)
		tx.Get(y)
	}
	b.ResetTimer()

	var counter uint64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := atomic.AddUint64(&counter, 1)
			if id%10 == 0 {
				lu.Atomically(inc)
			} else {
				lu.Atomically(load)
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
