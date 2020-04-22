# **go-stm: Software Transactionl Memory in Go**
Team: Yanxing Pan  
AndrewID: yanxingp 

## **URL**
[https://github.com/yanxingp/go-stm](https://github.com/yanxingp/go-stm)

## **Summary**
I would like to implement a simple software transactional memory as a Go package, and benchmark its performance under different workloads.

## **Background**
[Transactional memory (TM)](https://en.wikipedia.org/wiki/Transactional_memory) is introduced to us in [15-418 Lecture 19](http://www.cs.cmu.edu/~418/lectures/19_transactionalmem.pdf) as an alternative to the standard synchronization method that uses locks. [Software transactional memory (STM)](https://en.wikipedia.org/wiki/Software_transactional_memory), which is TM implemented in software, has been in research since late 1990's and has recently started to appear in products and various programming languages (like Haskell and Clojure).

THe [Go language](https://golang.org/) is built for programming concurrency. It has powerful built-in support for concurrency. The basic building block of concurrency in Go is *goroutine*, which is a light-weight thread of execution managed by the Go rontime. Instead of relying on shared memory, Go encourages programmers to use *channels* for synchronization. These features make Go an interesting choice for implementing STM.

I plan to base my implementation upon the [TL2 (Transactional Locking II)](https://perso.telecom-paristech.fr/kuznetso/INF346-2015/papers/tl2.pdf) algorithm. TL2 is a highly efficient STM algorithm that overcomes the two major drawbacks in STM implementation design: Closed-Memory and Need for Specialized Runtime. TL2 uses a global version clocks and two-phase locking scheme, and employs commit-time lock acquisition.

## **Challenge**
Implementing STM is solving a concurrency problem, which is inherently difficult. A simple mistake may take a huge amount of time to debug.

From a personal perspective, since the topic of this project is an extension to what we learned in class, there will be challenges in learning new materials. I will need to read a few papers and understand them thoroughly before I can start implementing the system.

## **Resources**
### **Literature**
The central piece of information I will need for this project is the TL2 paper:
* [Transactional Locking II](https://perso.telecom-paristech.fr/kuznetso/INF346-2015/papers/tl2.pdf)

Besides, there are a lot of publications on TM and STM. Some of them:
* [Transactional Memory: Architectural Support for Lock-Free Data Structures](http://cs.brown.edu/~mph/HerlihyM93/herlihy93transactional.pdf)
* [What Really Makes Transactions Faster?](http://people.csail.mit.edu/shanir/publications/TRANSACT06.pdf)
* [Software Transactional Memory](https://groups.csail.mit.edu/tds/papers/Shavit/ShavitTouitou-podc95.pdf)
* [Beautiful Concurrency](https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/beautiful.pdf)

Also, there are many excellent online blog posts that talk about STM.

### **Software**
The Haskell and Clojure programming languages have built-in support for STM, which may provide inspiration on the API design of my own package:
* [Haskell's stm package](https://hackage.haskell.org/package/stm)
* [Clojure's Refs and Transactions](https://clojure.org/reference/refs)

Also, there are open-source implementations on STM for reference:
* [pveentjer's stm implementation for JVM](https://github.com/pveentjer/Multiverse)
* [lukechampine and anacrolix's stm implementation for Go](https://github.com/anacrolix/stm)
* [nbronson's stm implementation for Scala](https://github.com/nbronson/scala-stm)
* [decillion's stm implementation for Go](https://github.com/decillion/go-stm)

## **Goals and Deliverables**
The goals for STM implementation:
* Implement a correct and robust STM package.
* Avoid race, deadlock or priority inversion conditions.
* Achieve a performance that should be better than course-grained locking.
* Make the API of the STM package clean and easy to use.

The deliverables will be:
* STM implementation code and benchmark code
* Final report that contains detailed description of software design and implementation, as well as benchmark results.
* A Google Slide for presentation.

## **Platform Choice**
Since this project is not about gaining computation speedup with parallelism, platform choice is not very important.

Both development and benchmark will be on my own laptop. Should the need for larger machines arise, I will use the GHC machines.

## **Schedule**
From the day this proposal is submitted, there will be about three weeks to finish this project.  
I plan to divide them into 6 half-weeks and plan accordingly:

|Half-Week|Plan|
|---------|----|
|1|Read and understand TL2 and STM paper. <br>Think about system architecture.|
|2|Read other material to reinforce understanding of STM.<br>Build system skeleton and design API.|
|3|Implement core part of the systems.|
|4|Implement core part of the systems.<br>Prepare checkpoint report.|
|5|Debug.<br>|
|6|Design and run benchmarks.<br>Prepare final report.|