# **go-stm: Software Transactionl Memory in Go**
Team: Yanxing Pan  
AndrewID: yanxingp 

## **URL**
[https://github.com/yanxingp/go-stm](https://github.com/yanxingp/go-stm)

## **Summary**
I want to implement a simple software transactional memory as a Go package, and benchmark its performance under different workloads.

## **Background**
[Transactional memory (TM)](https://en.wikipedia.org/wiki/Transactional_memory) is introduced to us in [15-418 Lecture 19](http://www.cs.cmu.edu/~418/lectures/19_transactionalmem.pdf) as an alternative to the classic synchronization method that uses locks. [Software transactional memory (STM)](https://en.wikipedia.org/wiki/Software_transactional_memory), which is TM implemented in software, has been in research since late 1990's and has recently started to appear in products and various programming languages (like Haskell and Clojure).

THe [Go language](https://golang.org/) is built to program concurrency, for which it has powerful built-in support. The basic building block of concurrency in Go is *goroutine*, which is a light-weight thread of execution managed by the Go rontime. Instead of relying on shared memory, Go encourages programmers to use *channels* for synchronization. These features make Go an interesting choice for implementing STM.

I plan to base my implementation upon the [TL2 (Transactional Locking II)](https://perso.telecom-paristech.fr/kuznetso/INF346-2015/papers/tl2.pdf) algorithm. TL2 is a highly efficient STM algorithm that overcomes the two major drawbacks in STM implementation design: Closed-Memory and Need for Specialized Runtime. TL2 uses a global version clocks and two-phase locking scheme, and employs commit-time lock acquisition.

## **Approach**

### **Package API**
Here is a summary of how this package can be used.

The package defines two data types:  
* `Var`: An variable/object to be used inside a transaction. As an implementation detail, a **Var** is also the basic unit that a lock is asscoiated with.
* `Trx`: A bookkeeper data structure for a transaction. Load and store operations in a transaction are associated with a `trx` instance for bookkeeping. Users don't need to initialize any `Trx` instances; they are created internally by the STM.

And a few functions:  
* `func NewVar(val interface{})`: Create a transactional variable with the value `val`.
* `func (trx *Trx) Load(x *Var)`: Read a transactional variable `x`'s value in a transaction.
* `func (trx *Trx) Store(x *Var, val interface{})`: Write the value `val` to transactional variable `x`.
* `func Atomically(trans func(trx *Trx) interface{})`: Execute the transaction `trans` atomically. `transact` is a callback function that must have  a `Trx` type as its parameter. The `Atomically` function will return the return value of `trans`.
* `func ReadOnlyAtomically(trans func(trx *Trx) interface{})`: A read-only version of `Atomically` to achieve higher performance. Only `Load` operations are allowed in the callback `trans`; otherwise, the program panics.

The API is very simple but is able to demonstate the capability of STM.

Here is an example of how the package can be used:
```go
func BasicUsage() {
    // Build transactional variables
	old := 1
    v := NewVar(old)
    
    // Define the transaction
	atomicIncrement := func(trx *Trx) interface{} {
        // Perform read & write ops inside the transaction
		val := trx.Load(v).(int)
		trx.Store(v, val+1)
		nval := trx.Load(v).(int)
		return nval
    }
    
    // Execute the transaction atomically
    nval := Atomically(atomicIncrement)
    
	fmt.Println(nval)
}
```
### **TL2 Algorithm & Implementation**

**1. Global Version-Clock**

This STM package is an implementation of the [Transactional Locking II (TL2)](https://perso.telecom-paristech.fr/kuznetso/INF346-2015/papers/tl2.pdf) algorithm. Therefore, building an in-depth understanding of this algorithm is a key part of this project.

I would like to briefly introduce TL2 and talk about some important implementation details along the way.

TL2 is based upon a **global version-clock**, which is its central means of detecting and resolving conflicts. The version-clock is incremented atomically to avoid race conditions.

I will talk about when the version-clock is incremented later, but the incrementation operation is quite rare so this global variable will not be a big performance bottleneck.

In my implementation, the global version-clock is simply a global uint64 variable:
```go
var globalClock uint64
```

**2. Variable/Object Version Number & Lock**

The primary difference between a transactional variable (`Var`) with an ordinary variable is that, for each `Var` there is a version number and a write-lock associated with this `Var`.

The version number is sampled from the global version-clock when the `Var` is created or updated, and is used to detect conflicts. The write-locks are used to synchronize write operations on the data.

In my implementation, the `Var` type is simply:
```go
type Var struct {
	wlock uint64 // Write lock
	ver uint64 // Version number
	data interface{} // Address of the data
}
```

**3. Read-set and Write-set**

The transaction bookkeeper has to maintain records of `Load` and `Store` operations inside the transaction that it is managing. These operations are organized into two collections call `read-set` and `write-set`. This is a common design choice in many STM implementations.

**My Implementation:**

The `read-set` and `write-set` are implemented as two linked-lists inside a `Trx` struct:
```go
type Trx struct {
    ...
	// Read-set and write-set as linked lists
	rset, wset *list.List
    ...
}
```

**4. Algorithm Workflow & Conflict Dection**
Now that we have the necessary data structures, we can briefly describe the workflow of this algorithm:

## **Results**

## **References**
