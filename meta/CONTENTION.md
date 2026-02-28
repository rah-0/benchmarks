# Contention Benchmarks

`GOMAXPROCS=8`, `benchtime=10s`. All results use a package-level sink variable to prevent dead-code elimination.

---

## A. Shared-State Protection

Use when you have one piece of shared state and need correctness first. Simple API, exact semantics.

- **Atomic** — best for trivial numeric updates (counters, gauges).
- **Mutex** — general-purpose; protects arbitrary state.
- **RWMutex** — read-heavy state with occasional writes.
- **CAS Spin / CAS Backoff** — educational / special-case; CAS Backoff yields via `Gosched()` between retries which reduces cache-line bouncing but is scheduler-yielding backoff, not comparable to `atomic.Add` as pure work.

### Write-Only (Inc-Only)

| Algorithm | ns/op | B/op | allocs/op | Benchmark |
|-----------|------:|-----:|----------:|-----------|
| CAS Backoff | 4.06 | 0 | 0 | `Counters_IncOnly/CAS_Backoff` |
| Atomic | 11.28 | 0 | 0 | `Counters_IncOnly/Atomic` |
| RWMutex | 23.29 | 0 | 0 | `Counters_IncOnly/RWMutex` |
| Mutex | 29.38 | 0 | 0 | `Counters_IncOnly/Mutex` |
| CAS Spin | 49.87 | 0 | 0 | `Counters_IncOnly/CAS_Spin` |

### Read-Only (Get-Only)

| Algorithm | ns/op | B/op | allocs/op | Benchmark |
|-----------|------:|-----:|----------:|-----------|
| Atomic | 1.78 | 0 | 0 | `Counters_GetOnly/Atomic` |
| CAS Backoff | 1.78 | 0 | 0 | `Counters_GetOnly/CAS_Backoff` |
| CAS Spin | 1.79 | 0 | 0 | `Counters_GetOnly/CAS_Spin` |
| RWMutex | 26.62 | 0 | 0 | `Counters_GetOnly/RWMutex` |
| Mutex | 31.70 | 0 | 0 | `Counters_GetOnly/Mutex` |

### Mixed (90% Write / 10% Read)

| Algorithm | ns/op | B/op | allocs/op | Benchmark |
|-----------|------:|-----:|----------:|-----------|
| CAS Backoff | 4.48 | 0 | 0 | `Counters_Mixed90_10/CAS_Backoff` |
| Atomic | 13.52 | 0 | 0 | `Counters_Mixed90_10/Atomic` |
| RWMutex | 14.04 | 0 | 0 | `Counters_Mixed90_10/RWMutex` |
| Mutex | 30.42 | 0 | 0 | `Counters_Mixed90_10/Mutex` |
| CAS Spin | 34.76 | 0 | 0 | `Counters_Mixed90_10/CAS_Spin` |

---

## B. Single-Owner / Actor Model

Use when one goroutine owns the state and others send operations via channels. Good for complex operations, cancellation, sequencing. Not for max throughput.

### Write-Only (Inc-Only)

| Buffer | ns/op | B/op | allocs/op | Benchmark |
|-------:|------:|-----:|----------:|-----------|
| 4096 | 70.28 | 0 | 0 | `Counters_IncOnly/Channel_buf4096` |
| 256 | 72.25 | 0 | 0 | `Counters_IncOnly/Channel_buf256` |
| 1 | 178.3 | 0 | 0 | `Counters_IncOnly/Channel_buf1` |
| 0 | 212.0 | 0 | 0 | `Counters_IncOnly/Channel_buf0` |

### Read-Only (Get-Only, request/reply round-trip)

| Buffer | ns/op | B/op | allocs/op | Benchmark |
|-------:|------:|-----:|----------:|-----------|
| 4096 | 349.4 | 0 | 0 | `Counters_GetOnly/Channel_buf4096` |
| 256 | 355.0 | 0 | 0 | `Counters_GetOnly/Channel_buf256` |
| 0 | 372.3 | 0 | 0 | `Counters_GetOnly/Channel_buf0` |
| 1 | 393.8 | 0 | 0 | `Counters_GetOnly/Channel_buf1` |

### Mixed (90% Write / 10% Read)

| Buffer | ns/op | B/op | allocs/op | Benchmark |
|-------:|------:|-----:|----------:|-----------|
| 256 | 94.93 | 0 | 0 | `Counters_Mixed90_10/Channel_buf256` |
| 4096 | 94.24 | 0 | 0 | `Counters_Mixed90_10/Channel_buf4096` |
| 1 | 203.4 | 0 | 0 | `Counters_Mixed90_10/Channel_buf1` |
| 0 | 224.8 | 0 | 0 | `Counters_Mixed90_10/Channel_buf0` |

---

## C. Contention Reduction

Use when you have very hot writes and can pay in memory, read cost, or precision.

- **Sharded / PerCPU** — exact writes, near-zero contention; reads cost O(shards).
- **LocalBuffered** — eventual consistency; goroutine-local accumulator flushed at threshold.
- **Approx** — fastest; probabilistic sampling via xorshift PRNG, trades accuracy.
- **FlatCombining** — niche; one combiner applies batched ops. Helps when ops are heavier than `++`.
- **MPSC Aggregator** — batching via channel; good when you already have an event pipeline.
- **RCU** — read-mostly structures; writes allocate a new copy (expensive).

### Write Throughput (Inc-Only)

| Algorithm | ns/op | B/op | allocs/op | Notes | Benchmark |
|-----------|------:|-----:|----------:|-------|-----------|
| Striped counter (N=8) | 0.222 | 0 | 0 | 8 padded shards | `Reduction_IncOnly/Sharded_8` |
| Striped counter (N=NumCPU) | 0.248 | 0 | 0 | One shard per CPU | `Reduction_IncOnly/Sharded_NumCPU` |
| Per-GOMAXPROCS counter | 0.252 | 0 | 0 | Sized to `GOMAXPROCS(0)` | `Reduction_IncOnly/PerCPU` |
| Probabilistic (1-in-1000) | 0.313 | 0 | 0 | xorshift sample | `Reduction_IncOnly/Approx_rate1000` |
| Local buffered (t=256) | 0.314 | 0 | 0 | Flush every 256 | `Reduction_IncOnly/LocalBuffered_t256` |
| Local buffered (t=64) | 0.757 | 0 | 0 | | `Reduction_IncOnly/LocalBuffered_t64` |
| Probabilistic (1-in-100) | 0.777 | 0 | 0 | | `Reduction_IncOnly/Approx_rate100` |
| Local buffered (t=16) | 2.35 | 0 | 0 | More frequent flushes | `Reduction_IncOnly/LocalBuffered_t16` |
| Probabilistic (1-in-10) | 3.61 | 0 | 0 | | `Reduction_IncOnly/Approx_rate10` |
| MPSC aggregator (buf=256) | 88.27 | 0 | 0 | Channel + aggregator goroutine | `Reduction_IncOnly/MPSC_buf256` |
| Flat combining | 91.71 | 9 | 0 | TryLock + combiner bottleneck | `Reduction_IncOnly/FlatCombining` |
| MPSC aggregator (buf=4096) | 106.5 | 0 | 0 | | `Reduction_IncOnly/MPSC_buf4096` |
| RCU (copy-on-write) | 126.6 | 29 | 3 | Allocates per write | `Reduction_IncOnly/RCU` |

### Read Throughput (Get-Only)

| Algorithm | ns/op | B/op | allocs/op | Notes | Benchmark |
|-----------|------:|-----:|----------:|-------|-----------|
| Striped counter (N=8) | 1.36 | 0 | 0 | Sums 8 shards | `Reduction_GetOnly/Sharded_8` |
| Per-GOMAXPROCS counter | 1.36 | 0 | 0 | Sums all slots | `Reduction_GetOnly/PerCPU` |
| RCU (copy-on-write) | 1.81 | 0 | 0 | Single pointer load | `Reduction_GetOnly/RCU` |
| Probabilistic (rate=100) | 1.82 | 0 | 0 | Single atomic load | `Reduction_GetOnly/Approx_rate100` |

---

## D. Message Passing / Queues

Use when the problem is moving data between goroutines, not protecting a counter.

- **SPSC Ring Buffer** — lock-free, cache-friendly single-producer single-consumer.
- **Channel** — Go runtime channel; same SPSC role, higher overhead.
- **Disruptor MPSC** — sequence-based ring, many producers → one consumer. Sensitive to commit-ordering contention on the contiguous cursor.

### SPSC: Ring Buffer vs Channel

Both use non-blocking send + `Gosched()` backpressure.

| Implementation | Size | ns/op | B/op | allocs/op | Benchmark |
|----------------|-----:|------:|-----:|----------:|-----------|
| Ring Buffer | 256 | 16.56 | 0 | 0 | `SPSC/RingBuffer_size256` |
| Ring Buffer | 65536 | 26.10 | 0 | 0 | `SPSC/RingBuffer_size65536` |
| Ring Buffer | 4096 | 29.55 | 0 | 0 | `SPSC/RingBuffer_size4096` |
| Channel | 65536 | 41.21 | 0 | 0 | `SPSC/Chan_buf65536` |
| Channel | 4096 | 41.66 | 0 | 0 | `SPSC/Chan_buf4096` |
| Channel | 256 | 46.04 | 0 | 0 | `SPSC/Chan_buf256` |

Ring buffer is ~2× faster than channels. Both zero-alloc.

### Disruptor MPSC (4 producers, 1 consumer)

| Ring Size | ns/op | B/op | allocs/op | Benchmark |
|----------:|------:|-----:|----------:|-----------|
| 256 | 2,619 | 0 | 0 | `DisruptorMPSC/size256` |
| 4096 | 2,713 | 0 | 0 | `DisruptorMPSC/size4096` |
| 65536 | 2,781 | 0 | 0 | `DisruptorMPSC/size65536` |

ns/op includes 4 producer publishes + 4 consumer reads per iteration. The contiguous commit CAS is the bottleneck.

---

## E. Concurrent Maps

Use when you need shared key/value access.

- **sync.Map** — optimized for read-mostly, many goroutines. Allocates on writes.
- **MutexMap** — `sync.Mutex` + `map`. Predictable; better for write-heavy.

| Operation | sync.Map | MutexMap | Benchmark |
|-----------|------:|------:|-----------|
| StoreOnly | 83.13 ns/op (72 B, 3 allocs) | 43.71 ns/op (0 B, 0 allocs) | `Map_StoreOnly/*` |
| LoadOnly | 9.77 ns/op (0 B, 0 allocs) | 49.62 ns/op (0 B, 0 allocs) | `Map_LoadOnly/*` |

sync.Map is 5× faster for reads but 2× slower (and allocates) for writes.

---

## F. Runtime Primitives

Common building blocks — not contention mechanisms per se, but useful for benchmarking coordination costs.

### Semaphore (channel-based)

| Max Concurrency | ns/op | B/op | allocs/op | Benchmark |
|----------------:|------:|-----:|----------:|-----------|
| 1 | 107.2 | 0 | 0 | `Semaphore/max1` |
| 4 | 198.7 | 0 | 0 | `Semaphore/max4` |
| 16 | 30.98 | 0 | 0 | `Semaphore/max16` |
| 64 | 31.80 | 0 | 0 | `Semaphore/max64` |

Contention spikes when max concurrency is near the number of goroutines competing. Once capacity exceeds demand, acquire/release is cheap.

### Goroutine Spawn + WaitGroup

| Workers | ns/op | B/op | allocs/op | Benchmark |
|--------:|------:|-----:|----------:|-----------|
| 1 | 347 | 40 | 2 | `SpawnAndWait/workers1` |
| 4 | 891 | 112 | 5 | `SpawnAndWait/workers4` |
| 16 | 3,305 | 400 | 17 | `SpawnAndWait/workers16` |
| 64 | 13,588 | 1,552 | 65 | `SpawnAndWait/workers64` |

~210 ns/goroutine (creation + scheduling + barrier).

### sync.Pool

| Mode | ns/op | B/op | allocs/op | Benchmark |
|------|------:|-----:|----------:|-----------|
| Serial | 10.86 | 0 | 0 | `Pool/Serial` |
| Parallel | 1.41 | 0 | 0 | `Pool/Parallel` |

Per-P fast path scales perfectly under parallelism.

### sync.Once

| Mode | ns/op | B/op | allocs/op | Benchmark |
|------|------:|-----:|----------:|-----------|
| Serial (fresh init) | 26.24 | 24 | 1 | `Once/Serial` |
| Contended (shared init) | 2.98 | 0 | 0 | `Once/Contended` |
| Fast path (already init) | 2.94 | 0 | 0 | `Once/FastPath_AfterInit` |

After first init, `Do()` is a single atomic load (~3 ns).