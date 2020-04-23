// trx_test.go - Tests for STM implementation.
// Author: Yanxing Pan (yanxingp)
// Date: April 2020
package stm

import (
	"sync"
	"testing"
	"time"
)

func TestBasicIncrement(t *testing.T) {
	old := 1
	v := NewVar(old)

	atomicIncrement := func(trx *Trx) interface{} {
		val := trx.Load(v).(int)
		trx.Store(v, val+1)
		nval := trx.Load(v).(int)
		return nval
	}

	nval := Atomically(atomicIncrement)
	if nval != old+1 {
		t.Errorf("nval = %v; want %v", nval, old+1)
	}
}

func TestBasicConcurrency(t *testing.T) {
	niter := 100

	for i := 0; i < niter; i++ {
		x := NewVar(0)
		y := NewVar(0)

		thread1 := func(trx *Trx) interface{} {
			trx.Store(x, 1)
			trx.Store(x, 2)
			return trx.Load(x)
		}

		thread2 := func(trx *Trx) interface{} {
			xval := trx.Load(x)
			trx.Store(y, xval)
			return trx.Load(y)
		}

		var yval int

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			Atomically(thread1)
		}()

		go func() {
			defer wg.Done()
			time.Sleep(20 * time.Nanosecond)
			yval = Atomically(thread2).(int)
		}()

		wg.Wait()
		// fmt.Printf("x = %v, y = %v\n", xval, yval)

		if yval != 0 && yval != 2 {
			t.Errorf("y = %v; should be 0 or 2\n", yval)
		}
	}
}

func TestConcurentIncrement(t *testing.T) {
	ng := 20
	niter := 1000

	v := NewVar(0)

	atomicGet := func(trx *Trx) interface{} {
		val := trx.Load(v).(int)
		return val
	}

	atomicIncrement := func(trx *Trx) interface{} {
		val := trx.Load(v).(int)
		trx.Store(v, val+1)
		val = trx.Load(v).(int)
		trx.Store(v, val+1)

		nval := trx.Load(v).(int)
		return nval
	}

	var wg sync.WaitGroup
	wg.Add(ng)

	for i := 0; i < ng; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < niter; j++ {
				Atomically(atomicIncrement)
			}
		}()
	}

	wg.Wait()
	nval := Atomically(atomicGet).(int)

	// Non-TM. Will cause race.
	vn := 0

	var wg2 sync.WaitGroup
	wg2.Add(ng)

	increment := func() interface{} {
		val := vn
		vn = val + 1
		val = vn
		vn = val + 1

		nval := vn
		return nval
	}

	for i := 0; i < ng; i++ {
		go func() {
			defer wg2.Done()
			for j := 0; j < niter; j++ {
				increment()
			}
		}()
	}

	wg2.Wait()
	// fmt.Printf("After each of %v goroutines did %v iterations\n", ng, niter)
	// fmt.Printf("Using TM: nval = %v\tNot using TM: nval = %v\n", nval, vn)

	// Check
	if nval != ng*niter*2 {
		t.Errorf("After each of %v goroutines did %v iterations, nval is %v; should be %v\n", ng, niter, nval, ng*niter*2)
	}

}

// Decellion's increment test
// https://github.com/decillion/go-stm/blob/master/tran_test.go
func TestDecillion(t *testing.T) {
	iter := 1000

	x := NewVar(0)
	y := NewVar(0)

	inc := func(trx *Trx) interface{} {
		trx.Store(x, trx.Load(x).(int)+1)
		trx.Store(y, trx.Load(y).(int)+1)
		trx.Store(x, trx.Load(x).(int)+1)
		trx.Store(y, trx.Load(y).(int)+1)
		return nil
	}

	read := func(trx *Trx) interface{} {
		var curr [2]int
		curr[0] = trx.Load(x).(int)
		curr[1] = trx.Load(y).(int)
		return curr
	}

	wg := sync.WaitGroup{}

	for i := 0; i < 2*iter; i++ {
		wg.Add(1)
		j := i
		go func() {
			if j%2 == 0 {
				Atomically(inc)
			} else {
				Atomically(read)
			}
			wg.Done()
		}()
	}

	wg.Wait()

	curr := Atomically(read).([2]int)
	currX, currY := curr[0], curr[1]

	if currX != iter*2 || currY != iter*2 {
		t.Errorf("want: (x, y) = (%v, %v) got: (x, y) = (%v, %v)",
			iter*2, iter*2, currX, currY)
	}
}
