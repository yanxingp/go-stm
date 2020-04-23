# **Checkpoint Report** 
**Project: go-stm: Software Transactional Memory in Go**

Team: Yanxing Pan  
AndrewID: yanxingp 

## **Project Summary**
Implement a simple software transactional memory (STM) system as a Go package, and benchmark its performance under different workloads.

## **Current Progress**
Progress so far:
1. Read papers on STM algorithms &#x2611;
2. Design & implement STM based on TL2 algorithm &#x2611;
3. Debug concurrency &#x2611;
4. Run benchmark on different workloads & analyze results
5. Prepare report and presentation

I have finished the implementation of the STM package based on [the TL2 algorithm]((https://perso.telecom-paristech.fr/kuznetso/INF346-2015/papers/tl2.pdf)), written a few tests on the most common use cases and debugged the program to pass the tests.

Project progress has been quite smooth, going ahead of the schedule:
|Half-Week|Plan|Done|
|---------|----|----|
|1|Read and understand TL2 and STM paper. <br>Think about system architecture.|Yes|
|2|Read other material to reinforce understanding of STM.<br>Build system skeleton and design API.|Yes|
|3|Implement core part of the systems.|Yes|
|4|Implement core part of the systems.<br>Prepare checkpoint report.|Yes|
|5|Debug.<br>|In part|
|6|Design and run benchmarks.<br>Prepare final report.|No|

The implementation has been more straightforward than I expected, and the debugging process has been quite smooth (or perhaps some more troublesome bugs are still hidden).

## **Work to Do**
Here are the things that I plan to do next:
* Write more test cases for more complex scenarios and debug
* Inspect the implementation and find opportunities to optimize
* Design & run benchmarks for different workloads
* Prepare report and presentation

Here is a renewed schedule for the new objectives:
|Half-Week|Plan|
|---------|----|
|4|Thorough testing & debugging, optimization|
|5|Design & run benchmarks|
|6|Prepare report & presentation|

## **Preliminary Results**
Code can be found in the project repository:  
[https://github.com/yanxingp/go-stm](https://github.com/yanxingp/go-stm)

Here is the result from one of the test cases: concurrent increments.  
The test works as follow:
1. One single variable is concurrently incremented by multiple threads. 
2. Each thread will run a number of iterations. 
3. Each iteration is an atomic transaction, where the variable will be incremented two times.

If there is no atomicity, race condition will occur and the final value of the variable will be less than **2 * num_thread * num_iteration**.  
With transactional memory, there will be no race and the result should be correct:
```
pyx-mac:stm pyx$ go test
After each of 20 goroutines did 1000 iterations
Using TM: nval = 40000  Not using TM: nval = 29480
PASS
ok      _/Users/pyx/Desktop/15-618/final-project/go-stm/stm     6.024s
```

## **Plan for Presentation**
I plan to prepare a slide and talk about:
1. 20-second introduction to STM and Go
2. Basic usage of STM with my package
3. Brief introduction to TL2 algorithm
4. Some important implementaion details (e.g. Conflict handling)
5. Benchmarking design and results
6. Performance comparsion with other synchronization methods
7. Conclusion
