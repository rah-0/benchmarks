# Iterator Benchmarks: Slice vs Channel (Pull Adapters)

This repo explores different ways of adapting a **push-only source** (`AscendRange`) into a common `RowIter` API (`Next() / Close()`).

Two main adapters are benchmarked:

- **Pull-SLICE**  
  Preloads the entire range into a slice, then serves rows from memory.  
  Pros: very fast `Next()` calls once loaded.  
  Cons: O(n) memory, high latency to first row, wasteful if iteration is aborted early.

- **Pull-CHAN**  
  Wraps the push source with `iter.Pull` (channel-based handoff). Rows are produced on demand in lockstep with the consumer.  
  Pros: bounded memory, very low latency to first row, cheap early termination.  
  Cons: per-row synchronization overhead (slower in full scans).

There are also baseline implementations (`SliceRowIter`, `ChanRowIter`) included for comparison.

---

## Benchmark Results

| Benchmark                                   | ns/op   | B/op    | allocs/op | Notes                                |
|---------------------------------------------|---------|---------|-----------|--------------------------------------|
| **Full Scan (100k rows)**                   |         |         |           |                                      |
| PullChan_Full_100k                          | 10.4 ms | 2.3 MiB | 199k      | Streaming, bounded memory            |
| PullSlice_Full_100k                         | 5.6 ms  | 6.2 MiB | 199k      | Preloaded, faster steady-state       |
| Baseline_ChanDirect_Full_100k               | 28.2 ms | 4.6 MiB | 299k      | Naive channel bridge (slowest)       |
| Baseline_SliceDirect_Full_100k              | 5.6 ms  | 6.2 MiB | 199k      | Same as PullSlice (direct preload)   |
| **Early Close (10 rows of 100k)**           |         |         |           |                                      |
| PullChan_EarlyClose10_100k                  | 1.38 µs | 600 B   | 25        | Immediate stop, minimal work         |
| PullSlice_EarlyClose10_100k                 | 1.66 ms | 3.9 MiB | 37        | Paid preload cost despite early stop |
| **Time To First Row (1M rows)**             |         |         |           |                                      |
| PullChan_TTFR_1M                            | 593 ns  | 456 B   | 16        | Yields first row instantly           |
| PullSlice_TTFR_1M                           | 8.16 ms | 39 MiB  | 38        | Must preload 1M rows first           |
| **Slow Consumer (100k rows, 50 work/row)**  |         |         |           |                                      |
| PullChan_SlowConsumerWork50_100k            | 12.2 ms | 2.3 MiB | 199k      | Slight sync overhead                 |
| PullSlice_SlowConsumerWork50_100k           | 8.6 ms  | 6.2 MiB | 199k      | Faster once preloaded                |
| **Slow Producer (100k rows, 50 work/item)** |         |         |           |                                      |
| PullChan_SlowProducerWork50_100k            | 13.8 ms | 2.3 MiB | 199k      | Producer cost dominates              |
| PullSlice_SlowProducerWork50_100k           | 7.1 ms  | 6.2 MiB | 199k      | Still faster steady-state            |
| **Multi-Range (10×10k ranges, total 100k)** |         |         |           |                                      |
| PullChan_MultiRange_100k_10parts            | 11.1 ms | 2.3 MiB | 199k      | Streaming across ranges              |
| PullSlice_MultiRange_100k_10parts           | 6.1 ms  | 6.2 MiB | 199k      | Full preload again                   |
| **Empty Range**                             |         |         |           |                                      |
| PullChan_EmptyRange                         | 490 ns  | 440 B   | 15        | Overhead of setup                    |
| PullSlice_EmptyRange                        | 327 ns  | 600 B   | 6         | Slightly cheaper                     |

---

## Key Conclusions

- **Latency to first row & early-close:**  
  Channels dominate (µs vs ms) because rows are streamed on demand. Slice preloading wastes time/memory if you don’t consume all.

- **Full scans (consume all rows):**  
  Slice is ~2× faster CPU but ~3× higher memory. Good if you always need everything.

- **Memory use:**  
  Channel stays bounded (~2–3 MiB for 100k rows). Slice grows O(n).

- **Baseline channel bridge** is worst (slowest, most allocs). The `iter.Pull` wrapper is the correct streaming adapter.
