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

## **Package API**
Here is a summary of how this package can be used.

The package defines two data types:  
* `Var`: An variable/object to be used inside a transaction. As an implementation detail, a **Var** is also the basic unit that a lock is asscoiated with.
* `Trx`: A bookkeeper data structure for a transaction. Load and store operations in a transaction are associated with a `trx` instance for bookkeeping. Users don't need to initialize any `Trx` instances; they are created internally by the STM.

And a few functions:  
* `func NewVar(val interface{})`: Create a transactional variable with the value `val`.
* `func (trx *Trx) Load(x *Var)`: Read a transactional variable `x`'s value in a transaction.
* `func (trx *Trx) Store(x *Var, val interface{})`: Write the value `val` to transactional variable `x`.
* `func Atomically(trans func(trx *Trx) interface{})`: Execute the transaction `trans` atomically. `trans` is a callback function that must have  a `Trx` type as its parameter. The `Atomically` function will return the return value of `trans`.
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

<div style="page-break-after: always;"></div>

## **TL2 Algorithm & Implementation**

### **1. Global Version-Clock**

This STM package is an implementation of the [Transactional Locking II (TL2)](https://perso.telecom-paristech.fr/kuznetso/INF346-2015/papers/tl2.pdf) algorithm. Therefore, building an in-depth understanding of this algorithm is a key part of this project.

I would like to briefly introduce TL2 and talk about some important implementation details along the way.

TL2 is based upon a **global version-clock**, which is its central means of detecting and resolving conflicts. The version-clock is incremented atomically to avoid race conditions.

I will talk about when the version-clock is incremented later, but the incrementation operation is quite rare so this global variable will not be a big performance bottleneck.

In my implementation, the global version-clock is simply a global uint64 variable:
```go
var globalClock uint64
```

### **2. Variable/Object Version Number & Lock**

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
And lock acquisition is done with Go's built-in CAS:
```go
func (v *Var) acquireLock() bool {
	return atomic.CompareAndSwapUint64(&v.wlock, 0, 1) {
}
```

### **3. Read-set and Write-set**

The transaction bookkeeper has to maintain records of `Load` and `Store` operations inside the transaction that it is managing. These operations are organized into two collections call `read-set` and `write-set`. This is a common design choice in many STM implementations.

In my implementation, the `read-set` and `write-set` are implemented as two linked-lists inside a `Trx` struct:
```go
type Trx struct {
    // ...
	// Read-set and write-set as linked lists
	rset, wset *list.List
    // ...
}
```

### **4. Algorithm Workflow**

Now that we have the necessary data structures, we can briefly describe the workflow of this algorithm.

For each transaction, perform the following steps:
* Sample the global version-clock and store it as read-version number (`rv`). `rv` is later used to detect conflicts.
* Run through a speculative execution - logically execute the `Load` and `Store` operations in the transaction. In implementation, it is just to place the `Load` and `Store` requests into the transaction bookkeeper's `read-set` and `write-set` in order.  
Validation are performed during the proess: If a `Var`'s version number > `rv`, which means this `Var` has been modified by another transaction; or the `Var`'s write-lock is locked, which means it is being modified by some other transaction, then we spotted a read-write conflict.  
In case of a conflict, the entire transaction aborts and retries.
* Lock the `write-set`. To update the values in every variable of `write-set`, we need to acquires all of their locks. If not all the locks are acquired, the transaction aborts.
* Increment the global version-clock. This is where we finally increment the global version-clock to notify other transactions that we are about to make an update to some variables. We record the updated version-clock as write-version number (`wv`).
* Validate the read-set again to ensure no changes are made to the variables between now and the validation during speculative exectuion.
* Commit and release the locks and update the version number of all the variables in the write-set to be `wv`.

As you can see, specific rules are set up to avoid read-write conflict and write-write conflicts.

For the read-only transaction, things are much simplier - no `read-set` or `write-set` are needed and the only work to be done is to perform validations.

<div style="page-break-after: always;"></div>

In my implementation, the transaction is placed inside a abort-retry loop that only exits upon seccess:
```go
func Atomically(trans func(trx *Trx) interface{}) interface{} {
	// Initialize transaction bookkeeper with corresponding mode
	trx := NewTrx(false)
    var res interface{} // Result of the transaction to be returned
    
	// Keep retrying until the transaction is committed
	for !trx.success {
        // ...
    }
    
    return res
}
```
And the implementation closely follows the workflow described above.

If you are interested in the details, you can go visit the git repository from the link at the top. The code is well-commented.

<div style="page-break-after: always;"></div>

## **Testing & Conflict Dection**

Now let's create two test cases to help us unedrstand how the conflict dections works.

> **Case 1**  
The value of Var `x` is intialized as 0, the value of Var `y` is initialized as 0.  
`Thread 1` tries to increment `x` twice in a transaction, in two separate `Store` operations.  
`Thread 2` tries to load the value of `x` and store it to `y`.

If both transactions are executed atomically, the value of `y` will be either 0 or 2. If the value of `y` is 1, then their is a read-write conflict.

> **Case 2**  
20 threads are trying to increment a Var `x` 2 times.  
The increment itself is not atomic - need a `Load` op and a `Store` op. We place them inside a transaction to make this atomic.  
We repeat this process 1000 times.

If the STM works fine, the final value of the variable should be 2 * 20 * 1000 = 40000. Otherwise, we have write-write conflicts.

How does TL2 detect these conflicts?

First, TL2 employs a speculative exectuion and commit-time lock aquisition strategy. Updates to the variables are not visible to other threads until the transaction commits.  

Second, if the Var `x` is read by `Thread 1`, then written by `Thread 2`, `x` will have a greater version number than `Thread 1`'s `rv`. And this will be detected during the second validation process.

Third, during the `write-set` lock acquisition step, a transaction must be able to acquire all the locks in its `write-set`. This is to ensure that no other transaction are trying to modify the value of these variables at the same time. This effectively avoids write-write conflicts.

These are the basic testing I have performed on my implementation. Thanks to Go's built-in testing facilities, testing process is made very smooth.

<div style="page-break-after: always;"></div>

## **Benchmarking**

The following benchmarking are performed with Go's built-in benchmarking facilities, which are great.

### **1. Read-Only Transaction API Performance**

As I have described, there is a special API `ReadOnlyAtomically` that should be able to execute read-only transactions very fast. 

Here are the benchmarking results of running a read-only transaction using `Atomically` and `ReadOnlyAtomically`, with [1, 2, 4, 8] CPU cores:
```
BenchmarkReadOnlyTrx                    198 ns/op
BenchmarkReadOnlyTrx-2                  123 ns/op            
BenchmarkReadOnlyTrx-4                  80.2 ns/op           
BenchmarkReadOnlyTrx-8                  69.3 ns/op 

BenchmarkReadOnlyTrxRO                  97.2 ns/op           
BenchmarkReadOnlyTrxRO-2                57.5 ns/op           
BenchmarkReadOnlyTrxRO-4                31.9 ns/op           
BenchmarkReadOnlyTrxRO-8                23.8 ns/op  
```
Indeed, there is a significant improvement in the performance of `ReadOnlyAtomically`.

### **2. Performance Comparision with Others' Implementations**

It is interesting to see how my implementation compares with others' work. I find two other implementations from GitHub:
* [decillion's stm implementation for Go](https://github.com/decillion/go-stm)
* [lukechampine and anacrolix's stm implementation for Go](https://github.com/anacrolix/stm)

Among the two of them, decillion is using TL2 algorithm and lukechampine is not. 

<div style="page-break-after: always;"></div>

Benchmarking the performance with [decillion's benchmarking code](https://github.com/decillion/go-stm/blob/master/benchmark/bench_test.go):

```
BenchmarkRead90Write10Trx                225 ns/op
BenchmarkRead90Write10Trx-2              277 ns/op
BenchmarkRead90Write10Trx-4              344 ns/op
BenchmarkRead90Write10Trx-8              529 ns/op

BenchmarkRead90Write10Decillion          168 ns/op
BenchmarkRead90Write10Decillion-2        112 ns/op
BenchmarkRead90Write10Decillion-4        75.1 ns/op
BenchmarkRead90Write10Decillion-8        72.5 ns/op

BenchmarkRead90Write10Lukechampine       336 ns/op
BenchmarkRead90Write10Lukechampine-2     528 ns/op
BenchmarkRead90Write10Lukechampine-4     591 ns/op
BenchmarkRead90Write10Lukechampine-8     1458 ns/op
```
As you can see from the results, TL2 is quite efficient. However, my implementation is not as good as decellion's and his/her implementation scales very well.

### **3. Performance under Different Workloads**

I come up with three different workloads:
* 90% Read 10% Writes
* 70% Read 30% Writes
* 50% Read 50% Writes

Running the benchmarks:

```
BenchmarkRead90Write10Trx                225 ns/op
BenchmarkRead90Write10Trx-2              277 ns/op
BenchmarkRead90Write10Trx-4              344 ns/op
BenchmarkRead90Write10Trx-8              529 ns/op

enchmarkRead70Write30Trx                 282 ns/op
BenchmarkRead70Write30Trx-2              346 ns/op
BenchmarkRead70Write30Trx-4              432 ns/op
BenchmarkRead70Write30Trx-8              656 ns/op

BenchmarkRead50Write50Trx                334 ns/op
BenchmarkRead50Write50Trx-2              391 ns/op
BenchmarkRead50Write50Trx-4              490 ns/op
BenchmarkRead50Write50Trx-8              764 ns/op
```
More writes we have, the higher the chance of conflict there will be, so it is exepected to be slower. 

However, I think the implementation performs quite well in handling large amounts of write operartions, as the performs didn't degrade very much.

## **References**
* Dave Dice, Ori Shalev, Nir Shavit, [Transactional Locking II](https://perso.telecom-paristech.fr/kuznetso/INF346-2015/papers/tl2.pdf)
* Maurice Herlihy, J. Eliot B. Moss, [Transactional Memory: Architectural Support for Lock-Free Data Structures](http://cs.brown.edu/~mph/HerlihyM93/herlihy93transactional.pdf)
* Dave Dice, Nir Shavit, [What Really Makes Transactions Faster?](http://people.csail.mit.edu/shanir/publications/TRANSACT06.pdf)
* Nir Shavit, Dan Touitou, [Software Transactional Memory](https://groups.csail.mit.edu/tds/papers/Shavit/ShavitTouitou-podc95.pdf)
* Simon Peyton Jones, [Beautiful Concurrency](https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/beautiful.pdf)
* Testing and benchmarking code of the package are inspired by the work of [decillion](https://github.com/decillion/go-stm).
