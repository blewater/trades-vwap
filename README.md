# Trades -> vwap
#### A VWAP (volume-weighted average price) calculator for streaming Coinbase trades

#### A single consumer with multiple VWAP producers design. 
A thread listens to the Coinbase socket connection and fans-out trade messages to the workers' routines pool. The thread streams the unmarshalled Json trade message to the messaging Queue for the Goroutines to process. There is Go's memory pool employed too. More of that is in the next section.
Goroutines receive in async and random fashion an available queued trade ticker to generate a VWAP result. 
By doing so, it consults the matching memory Product queue of VWAP results, so the newly generated VWAP value reflects being chronologically at the top of the last 200 VWAP data points generated before it.

![diagram](./vwap.drawio.png)

#### About memory poooling
Any long-running streaming service unavoidably puts enormous memory pressure on memory-managed languages, i.e., Go. The constant creation and disposal of temporary objects quickly fill up the avaiable memory heap resulting in intermittent activation of the garbage collection language runtime. While Go strives not to impact performance heavily, there is still a penalty, and an ill-designed service may render its runtime container unstable, i.e., `LXC`. To address that constant HEAP pressure, Go offers a **memory pool** for recycling temp objects and is employed in this service for big.Decimal, and other trade structs. While it appeared that float64 offers sufficient precision for the incoming trade values, it seemed more appropriate to employ them. A testing algorithm using float64 data types is also included for documentation purposes

#### Go routines
While Go offers lightweight threads in the fashion of Erlang, they still occur overhead, e.g., 4k stack each thus, a throttling design should be employed. A known straightforward, efficient pattern is thread pools. Launching to a specific limit at the service's bootstrap time, they scale with sufficient processing bandwidth to a much higher number of incoming requests. 

If the host allowed multiple connections from the same client IP, it would enable input processing parallelism. Since this is not the case here, it is still feasible to achieve a degree of parallelism later in the pipeline (as the included benchmark test shows) by queueing the ingested trade messages for the thread pool to process.

#### Go's sync.Map
Based on the language documentation, the `sync.Map` is suitable for disporportionate number of reads vs writing which is the case here.

#### Benchmarking
In the `makefile`, there is a target to benchmark 100 transactions (listening, processing, and ingesting the VWAP results queue)  in several pool sizes, i.e., 1, 2, 3, 5, 10, 100, 200. The benchmarks can only highlight trends since the transactions are asynchronous and heavily influenced by the day traffic. 

#### In Ubuntu 20.04 
**A run employing memory pool**
`make bench`
``` shell
Benchmark_100_VWAP_Trx_1Thread-16                      1        18692163055 ns/op        3954464 B/op      43279 allocs/op
Benchmark_100_VWAP_Trx_2Threads-16                     1        20447973426 ns/op         306000 B/op       4890 allocs/op
Benchmark_100_VWAP_Trx_3Threads-16                     1        14296614771 ns/op         314472 B/op       5058 allocs/op
Benchmark_100_VWAP_Trx_5Threads-16                     1        23358468964 ns/op         310888 B/op       4971 allocs/op
Benchmark_100_VWAP_Trx_10Threads-16                    1        17694583925 ns/op         318744 B/op       5129 allocs/op
Benchmark_100_VWAP_Trx_100Threads-16                   1        13594751784 ns/op         411320 B/op       6004 allocs/op
Benchmark_100_VWAP_Trx_200Threads-16                   1        14380694155 ns/op         482472 B/op       6487 allocs/op```

```
**A run without a memory pool**
```
Benchmark_100_VWAP_Trx_1Thread-16                      1        16615332934 ns/op
Benchmark_100_VWAP_Trx_2Threads-16                     1        26071586811 ns/op
Benchmark_100_VWAP_Trx_3Threads-16                     1        29935633754 ns/op
Benchmark_100_VWAP_Trx_5Threads-16                     1        23910722054 ns/op
Benchmark_100_VWAP_Trx_10Threads-16                    1        27848430412 ns/op
Benchmark_100_VWAP_Trx_100Threads-16                   1        22826719074 ns/op
Benchmark_100_VWAP_Trx_200Threads-16                   1        22012305333 ns/op
```
Benchmarking also generates profiling information.
Socket I/O and json marshalling present an opportunity for optimization as is more evident in the profile_cpu.svg diagram below. 
```sh
λ go tool pprof workflow.test profile_cpu.out
File: workflow.test
Type: cpu
Time: Nov 15, 2021 at 2:25pm (EET)
Duration: 173.26s, Total samples = 60ms (0.035%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top10
Showing nodes accounting for 60ms, 100% of 60ms total
Showing top 10 nodes out of 40
      flat  flat%   sum%        cum   cum%
      10ms 16.67% 16.67%       20ms 33.33%  internal/poll.(*FD).Read
      10ms 16.67% 33.33%       10ms 16.67%  runtime.(*randomEnum).next
      10ms 16.67% 50.00%       10ms 16.67%  runtime.bgscavenge.func1
      10ms 16.67% 66.67%       10ms 16.67%  runtime.epollwait
      10ms 16.67% 83.33%       10ms 16.67%  runtime.madvise
      10ms 16.67%   100%       10ms 16.67%  syscall.Syscall
         0     0%   100%       20ms 33.33%  bufio.(*Reader).Peek
         0     0%   100%       20ms 33.33%  bufio.(*Reader).fill
         0     0%   100%       20ms 33.33%  bytes.(*Buffer).ReadFrom
         0     0%   100%       20ms 33.33%  crypto/tls.(*Conn).Read
```
#### The generated profile_cpu.svg
![call graph](./profile_cpu.svg)

#### Logging, instrumentation and observations.
While this service only includes logging, I selected Uber's zap logger for its environment configuration awareness and efficient logging without employing reflection.