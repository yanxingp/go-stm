# **go-stm: Software Transactional Memory in Go**
Team: Yanxing Pan  
AndrewID: yanxingp

## **Idea**
I would like to implement a Go package on software transactional memory and benchmark its performance under various workloads.

## **Motivation**
Transactional memory (TM) is an interesting parallel computing topic and is introduced to us in 15-418 Lecture 19. Software transactional memory (STM), which is TM implemented in software, has been around in research environment since late 1990's and has recently started to appear in products and various programming languages (like Haskell and Clojure).

By implementing a simple STM package, I have the opportunity to take a deeper dive into TM and STM, and sharpen my skill in programming concurrency. By running benchmarks on different workloads and compare STM's performance with old-school locking, I can gain a deeper understanding on the challenges posed by synchronization in concurrent environments, and possibly some insights on how to tackle them.

I choose the Go programming language because it has built-in support for concurrency. Also I am learning Go this semester and programming with Go is very enjoyable.

## **Resources**
### **Literature**
There are a lot of publications on TM and STM. Some of them:
* [Transactional Memory: Architectural Support for Lock-Free Data Structures](http://cs.brown.edu/~mph/HerlihyM93/herlihy93transactional.pdf)
* [Software Transactional Memory](https://groups.csail.mit.edu/tds/papers/Shavit/ShavitTouitou-podc95.pdf)
* [Beautiful Concurrency](https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/beautiful.pdf)

Also, there are many excellent online blog posts that talk about STM.

### **Software**
The Haskell and Clojure programming languages have built-in support for STM, which may provide inspiration on the API design of my own package:
* [Haskell's stm package](https://hackage.haskell.org/package/stm)
* [Clojure's Refs and Transactions](https://clojure.org/reference/refs)

Also, there are open-source implementations on STM:
* [pveentjer's stm implementation for JVM](https://github.com/pveentjer/Multiverse)
* [lukechampine and anacrolix's stm implementation for Go](https://github.com/anacrolix/stm)
* [nbronson's stm implementation for Scala](https://github.com/nbronson/scala-stm)
* [decillion's stm implementation for Go](https://github.com/decillion/go-stm)

Some of them are actually implemented in Go, which is exactly what I am going to do here. But I don't plan to copy their code. I want to save the fun for myself and only use their code for reference, and take this as an opportunity to learn from more experienced Go programmers.

## **Workload**
[decillion's stm implementation](https://github.com/decillion/go-stm) contains about 900 lines of Go code.
[lukechampine and anacrolix's stm implementation](https://github.com/anacrolix/stm) contains about 2200 lines because their version supports a richer API.

I think the size of this project is within the reach for one person in about three weeks.

## **Challenges**
For now, I think the major challenges might be in implementing concurrency and understanding STM. I will have to learn more about STM to fully understand the challenges.

## **Deliverables**
The final product will be software and a report.  
The software will be a Go package of my stm implmentation as well as benchmarking code.
The report will describe the design and implementaion of the package in detail, as well as benchmarking results, some analysis on those result and conclusion drawn from the analysis.

## **Platform**
The development will be on my own laptop, which is a multi-core processor machine and that is good enough.
The benchmarking will be on the GHC machine, which has more cores and can perform the benchmarks at a larger scale.

