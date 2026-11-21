# Synchronization workloads and comparison boundaries

Use this guide to interpret the [measured comparisons](README.MD#results).
The loops are in [contention_test.go](contention_test.go); implementations are in
[contention.go](contention.go).

**Match completion and visibility guarantees before comparing speed.** Updating
a shared value, enqueueing an update, and producing an estimate are different
operations, even when each benchmark is described as a counter increment.

| Workload | Meaningful comparison |
| --- | --- |
| Synchronous shared counter | Atomic, lock, and CAS operations that finish the update before returning. |
| Deferred update submission | Submission to an owner goroutine or buffered aggregator; publication may happen later. |
| Reduced contention | Sharding, local accumulation, and sampling, with their read and visibility costs stated. |
| Queue transfer | One value transferred between the stated number of producers and consumers. |
| Shared-key map access | The same operation and populated key set. |
| Runtime primitive | Matching work and initialization state. |

`RunParallel` distributes operations across workers according to GOMAXPROCS.
Elapsed time divided by the total operations describes **aggregate throughput**,
not per-call latency. Tiny parallel time/op values should be read in that sense.

## Shared counters

The three access patterns are update-only, read-only, and nine updates followed
by one read per worker. Atomic access is the reference for equivalent synchronous
operations; channel submission has its own completion boundary.

| Implementation | Operation and semantics |
| --- | --- |
| Atomic | An atomic addition or load of one shared counter. |
| Mutex | An addition or read protected by `sync.Mutex`. |
| RWMutex | Writers take the exclusive lock; readers take the shared lock. |
| CAS Spin | Retry a compare-and-swap until the addition succeeds. |
| CAS Backoff | Retry compare-and-swap with increasing `runtime.Gosched` calls after failures. Timing includes scheduler yields. |
| Channel | Send an increment request to one owner goroutine, or send a read request and wait for its reply. |

Channel buffer sizes are 0, 1, 256, and 4,096. An increment measures submission;
the owner may apply it later. A read includes a request/reply exchange. Each
worker reuses its reply channel. These completion boundaries differ from the
atomic and locked updates.

## Reducing contention

Update cases measure increments. Read cases have no concurrent writers.
The implementations exchange write contention for memory, read cost, delayed
visibility, or approximation.

| Implementation | Operation and semantics |
| --- | --- |
| Sharded | Each worker updates a selected padded atomic shard. Reads sum all shards; a sum is not an atomic snapshot during concurrent writes. |
| PerCPU | An array of GOMAXPROCS padded slots, selected by worker ID. Workers are not pinned to physical CPUs. Reads sum all slots. |
| LocalBuffered | Each worker accumulates increments locally, flushes at a threshold, and flushes its remainder when finished. Reads omit unflushed increments. |
| Approx | A worker-local xorshift state samples updates at rates 10, 100, or 1,000. Accepted updates add the sampling rate. This measures an estimate, not an exact count. |
| FlatCombining | Queue an operation, then compete to apply a batch under a lock until that operation completes. |
| MPSC aggregator | Submit increments through a buffered channel to one consumer, which publishes accumulated totals. Submission does not wait for publication. |
| RCU | Allocate an updated counter state and replace the shared pointer with compare-and-swap; reads load the current pointer. Go's garbage collector reclaims old states. |

Sharded cases use eight shards or `runtime.NumCPU()` shards. Buffered cases use
thresholds 16, 64, and 256; MPSC channels use capacities 256 and 4,096. The
approximate cases use a fixed seed per worker and measure execution cost rather
than estimator accuracy. Read costs and visibility guarantees should be
considered alongside update throughput.

## Queue transfers

The single-producer/single-consumer cases compare a channel with a ring buffer.
Each operation transfers one value. Both retry non-blocking operations
with `runtime.Gosched` when full or empty. Channel capacities are 256, 4,096, and
65,536; the corresponding ring buffers reserve one slot to distinguish empty
from full, so their usable capacities are one smaller.

The Disruptor MPSC case divides all measured messages across four producers,
including any remainder, and consumes every message before returning. One
reported operation is one publish/consume transfer. Producers coordinate through
the contiguous commit cursor; the consumer observes published values in order.
Timing includes goroutine startup, retries, and completion waiting.

## Concurrent maps

Store-only and load-only cases compare `sync.Map` with a mutex-protected reference.
Workers repeatedly access the same set of 1,024 string keys. Load cases
populate both maps before timing; store benchmarks include first insertion and
subsequent replacement. Results describe this shared-key workload, not arbitrary
map sizes or access patterns.

## Runtime primitives

| Workload | One operation |
| --- | --- |
| Semaphore | Acquire and release a channel-based semaphore with capacity 1, 4, 16, or 64. No protected workload runs between them. |
| Spawn and wait | Create 1, 4, 16, or 64 goroutines with empty work and wait for all of them. |
| Pool access | Get a reusable buffer from `sync.Pool` and return it, serially or across workers. |
| Fresh initialization | Create a fresh instance, initialize it, and read its value. |
| Shared initialization | Call initialization and read from one shared instance. Initialization happens once; subsequent operations use the initialized path. |
| Already initialized | Initialize before timing, then repeatedly call initialization and read across workers. |

These cases perform different amounts of work. Compare matching operations and
their allocation counts instead of treating every row as an interchangeable
synchronization strategy.
