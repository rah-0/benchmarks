# Iterator benchmarks

Compare ways to adapt a push source to `RowIter.Next` and `Close`.
**Full scans use the direct slice iterator as the reference; other matching
workloads use the slice adapter.** Full scans, early termination, first-row
access, and concurrent callers have separate comparisons.

[Results](#results) · [Methodology](#methodology) · [Run](#running)

<!-- benchmark-results:start -->
## Results

**2026-11-21 · go1.27.4 · linux/amd64 · AMD Ryzen 9 5900HX with Radeon Graphics**
GOMAXPROCS: 16; samples per case: 3.

Tables show medians. **Speed = baseline time ÷ candidate time**: 2.00× is twice
the baseline speed; 0.50× is half. Allocated bytes and allocation counts describe
Go heap allocation per operation, not retained memory or process RSS.
These samples are descriptive; observed min–max ranges appear in the details
and are not confidence intervals. A zero denominator is shown as **—**.

[Environment](../results/environment.txt) · [Versions](../results/modules.txt) · [Reading the comparisons](../README.md#reading-the-comparisons)

<details>
<summary>Jump to a workload</summary>

- [Read all 100,000 rows](#read-all-100000-rows)
- [Stop after 10 of 100,000 rows](#stop-after-10-of-100000-rows)
- [First row from 1 million rows](#first-row-from-1-million-rows)
- [Consume 100,000 rows with consumer work](#consume-100000-rows-with-consumer-work)
- [Consume 100,000 rows with producer work](#consume-100000-rows-with-producer-work)
- [Read 100,000 rows across ten ranges](#read-100000-rows-across-ten-ranges)
- [Open and close an empty range](#open-and-close-an-empty-range)
- [Parallel scans of 10,000 rows](#parallel-scans-of-10000-rows)

</details>

### Read all 100,000 rows

One operation creates an iterator, consumes every row, and closes it.

Baseline: **Direct slice = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

The pull adapter uses iter.Pull; the direct-channel reference uses a producer goroutine and explicit channel handoff.

| Implementation | Time/scan, ms ↓ | Speed vs baseline ↑ | Allocated/scan, MiB ↓ | Allocs/scan ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Direct slice (baseline)** | 3.79 | 1.00× | 6.2 | 199,769 |
| Slice adapter | 3.76 | 1.01× | 6.2 | 199,771 |
| Pull adapter | 8.96 | 0.42× | 2.29 | 199,753 |
| Direct channel | 58.32 | 0.07× | 4.61 | 299,751 |

Slice adapter has the lowest recorded median time: 1.01× the speed of Direct slice, with 1.00× its allocated bytes per operation. The timing ranges of Slice adapter and Direct slice overlap; their order is inconclusive.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Direct slice | 3.75–3.8 ms | 6.2–6.2 MiB | 199,769 | 1.00× |
| Slice adapter | 3.71–3.78 ms | 6.2–6.2 MiB | 199,771 | 1.00× |
| Pull adapter | 8.94–9 ms | 2.29–2.29 MiB | 199,753 | 0.37× |
| Direct channel | 58.28–59.95 ms | 4.61–4.62 MiB | 299,751–299,752 | 0.74× |

</details>

### Stop after 10 of 100,000 rows

Create an iterator, read ten rows, and close it.

Baseline: **Slice adapter = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

The slice adapter preloads the requested range; the pull adapter produces rows on demand. Allocation totals do not measure peak retained memory.

| Implementation | Time/scan ↓ | Speed vs baseline ↑ | Allocated/scan ↓ | Allocs/scan ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Slice adapter (baseline)** | 1.28 ms | 1.00× | 3.91 MiB | 37 |
| Pull adapter | 1.13 µs | 1,133.72× | 456 B | 19 |

Pull adapter has the lowest recorded median time: 1,133.72× the speed of Slice adapter, with <0.01× its allocated bytes per operation.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Slice adapter | 1.24–1.34 ms | 3.91–3.91 MiB | 37 | 1.00× |
| Pull adapter | 1.13–1.13 µs | 456–456 B | 19 | <0.01× |

</details>

### First row from 1 million rows

Create an iterator, read one row, and close it. This is the complete first-row lifecycle, including close.

Baseline: **Slice adapter = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

The slice adapter preloads the requested range; the pull adapter produces rows on demand. Allocation totals do not measure peak retained memory.

| Implementation | Time/scan ↓ | Speed vs baseline ↑ | Allocated/scan ↓ | Allocs/scan ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Slice adapter (baseline)** | 9.98 ms | 1.00× | 39.81 MiB | 38 |
| Pull adapter | 328.9 ns | 30,339.88× | 314 B | 10 |

Pull adapter has the lowest recorded median time: 30,339.88× the speed of Slice adapter, with <0.01× its allocated bytes per operation.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Slice adapter | 8.34–10.2 ms | 39.81–39.81 MiB | 38 | 1.00× |
| Pull adapter | 328.4–329.5 ns | 314–314 B | 10 | <0.01× |

</details>

### Consume 100,000 rows with consumer work

Create, drain, and close an iterator while running the workload's synthetic consumer loop for each row.

Baseline: **Slice adapter = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

The slice adapter preloads the requested range; the pull adapter produces rows on demand. Allocation totals do not measure peak retained memory.

| Implementation | Time/scan, ms ↓ | Speed vs baseline ↑ | Allocated/scan, MiB ↓ | Allocs/scan ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Slice adapter (baseline)** | 4.94 | 1.00× | 6.2 | 199,771 |
| Pull adapter | 10.36 | 0.48× | 2.29 | 199,753 |

Slice adapter has the lowest recorded median time in this comparison.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Slice adapter | 4.92–4.99 ms | 6.2–6.2 MiB | 199,771 | 1.00× |
| Pull adapter | 10.3–10.39 ms | 2.29–2.29 MiB | 199,753 | 0.37× |

</details>

### Consume 100,000 rows with producer work

Create, drain, and close an iterator whose source runs the workload's synthetic producer loop for each row.

Baseline: **Slice adapter = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

The slice adapter preloads the requested range; the pull adapter produces rows on demand. Allocation totals do not measure peak retained memory.

| Implementation | Time/scan, ms ↓ | Speed vs baseline ↑ | Allocated/scan, MiB ↓ | Allocs/scan ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Slice adapter (baseline)** | 5.34 | 1.00× | 6.2 | 199,773 |
| Pull adapter | 10.84 | 0.49× | 2.29 | 199,755 |

Slice adapter has the lowest recorded median time in this comparison.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Slice adapter | 5.28–5.34 ms | 6.2–6.2 MiB | 199,773 | 1.00× |
| Pull adapter | 10.83–10.94 ms | 2.29–2.29 MiB | 199,755 | 0.37× |

</details>

### Read 100,000 rows across ten ranges

Create an iterator over ten consecutive ranges, consume all rows, and close it.

Baseline: **Slice adapter = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

The slice adapter preloads the requested range; the pull adapter produces rows on demand. Allocation totals do not measure peak retained memory.

| Implementation | Time/scan, ms ↓ | Speed vs baseline ↑ | Allocated/scan, MiB ↓ | Allocs/scan ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Slice adapter (baseline)** | 4.29 | 1.00× | 6.2 | 199,781 |
| Pull adapter | 9.41 | 0.46× | 2.29 | 199,773 |

Slice adapter has the lowest recorded median time in this comparison.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Slice adapter | 4.28–4.3 ms | 6.2–6.2 MiB | 199,781 | 1.00× |
| Pull adapter | 9.39–9.44 ms | 2.29–2.29 MiB | 199,773 | 0.37× |

</details>

### Open and close an empty range

Create an iterator for an empty source, request its end, and close it.

Baseline: **Slice adapter = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

The slice adapter preloads the requested range; the pull adapter produces rows on demand. Allocation totals do not measure peak retained memory.

| Implementation | Time/scan, ns ↓ | Speed vs baseline ↑ | Allocated/scan, B ↓ | Allocs/scan ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Slice adapter (baseline)** | 152.5 | 1.00× | 600 | 6 |
| Pull adapter | 265.6 | 0.57× | 296 | 9 |

Slice adapter has the lowest recorded median time in this comparison.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Slice adapter | 151.5–154.3 ns | 600–600 B | 6 | 1.00× |
| Pull adapter | 264.4–266.2 ns | 296–296 B | 9 | 0.49× |

</details>

### Parallel scans of 10,000 rows

Each parallel operation creates an iterator, reads 10,000 rows, and closes it; parallelism is four times GOMAXPROCS.

Baseline: **Slice adapter = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

The slice adapter preloads the requested range; the pull adapter produces rows on demand. Allocation totals do not measure peak retained memory.

| Implementation | Time/scan, µs ↓ | Speed vs baseline ↑ | Allocated/scan, KiB ↓ | Allocs/scan ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Slice adapter (baseline)** | 114.51 | 1.00× | 581.23 | 19,762 |
| Pull adapter | 258.58 | 0.44× | 232.7 | 19,753 |

Slice adapter has the lowest recorded median time in this comparison.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Slice adapter | 114.05–116.4 µs | 581.23–581.23 KiB | 19,762 | 1.00× |
| Pull adapter | 258.29–259.41 µs | 232.7–232.7 KiB | 19,753 | 0.40× |

</details>

**Recorded samples:** [iteration](../results/iteration.txt)

<!-- benchmark-results:end -->

## Methodology

The source contains integers and exposes `AscendRange`; each returned `Row` is
constructed by the adapter.

| Implementation | Behavior |
| --- | --- |
| Slice adapter | `PullSlice` copies the complete range into a slice before serving rows, even if the consumer stops early. |
| Pull adapter | `PullChan` uses `iter.Pull` to request values on demand. Its benchmark name does not imply an explicit Go channel. |
| Direct slice | `Baseline_SliceDirect` provides the direct iterator used as the full-scan reference. |
| Direct channel | `Baseline_ChanDirect` uses a producer goroutine, channel handoff, cancellation, and waiting during `Close`. |

| Workload | One operation |
| --- | --- |
| Full scan | Construct an adapter and consume 100,000 rows. |
| Early close | Consume ten rows from a 100,000-row range, then close. |
| First row | Construct an adapter, request one row from a one-million-row range, and close. |
| Producer/consumer work | Scan with additional synthetic work on the producer or consumer side. |
| Multiple ranges | Consume ten consecutive ranges. |
| Empty range | Construct and close an empty iterator. |
| Parallel callers | Perform the adapter workload across concurrent benchmark workers. |

The first-row result includes construction and cleanup; it is not an isolated
`Next` timing. Parallel elapsed time per operation describes aggregate throughput
and includes scheduling costs.

Streaming avoids preloading the complete range. Allocated bytes count cumulative
allocations during an operation, however, and do not measure retained or peak memory.

## Running

Requires Go 1.27.4+. From the repository root:

```sh
go test -p=1 -run '^$' -bench . -benchmem -benchtime=1s -count=3 ./iteration
```
