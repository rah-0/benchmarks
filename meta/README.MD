# Reflection, sorting, and synchronization

Compare field access, integer sorting, and shared-state operations in separate
workloads. **Direct access and Go `sort.Ints` are the reflection and sorting
references.** Synchronization references depend on the operation and its
completion guarantees; the [contention guide](CONTENTION.md) explains them.

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

- [Sort random integers](#sort-random-integers)
- [Set one string field](#set-one-string-field)
- [Read one string field](#read-one-string-field)
- [Synchronous shared-counter increments](#synchronous-shared-counter-increments)
- [Shared-counter reads](#shared-counter-reads)
- [Shared counter: 90% updates, 10% reads](#shared-counter-90-updates-10-reads)
- [Counter increment submission by channel capacity](#counter-increment-submission-by-channel-capacity)
- [Counter request/reply reads by channel capacity](#counter-requestreply-reads-by-channel-capacity)
- [Channel counter: 90% submissions, 10% reads](#channel-counter-90-submissions-10-reads)
- [Single-producer, single-consumer transfer](#single-producer-single-consumer-transfer)
- [Concurrent map stores](#concurrent-map-stores)
- [Concurrent map loads](#concurrent-map-loads)
- [Sharded counter updates](#sharded-counter-updates)
- [Read a sharded total](#read-a-sharded-total)
- [Buffered local increments](#buffered-local-increments)
- [MPSC counter submission](#mpsc-counter-submission)
- [Approximate counter update costs](#approximate-counter-update-costs)
- [Read an approximate counter](#read-an-approximate-counter)
- [Four-producer queue transfer](#four-producer-queue-transfer)
- [Semaphore acquire and release](#semaphore-acquire-and-release)
- [Start goroutines and wait](#start-goroutines-and-wait)
- [Pool get and put](#pool-get-and-put)
- [Once initialization paths](#once-initialization-paths)

</details>

### Sort random integers

One operation allocates and copies the input slice, then sorts it in ascending order. Each implementation receives a separately generated input of the same size and distribution.

Baseline: **Go sort.Ints = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

Libraries retain their default worker policies; Sorty uses MaxGor=3. GOMAXPROCS does not mean every implementation uses that many workers.

| Input / workload | Go sort.Ints time (1.00×) | Sorty speed ↑ | Pargo speed ↑ | Parsort speed ↑ |
| --- | ---: | ---: | ---: | ---: |
| 10,000 integers | 395.44 µs | 1.77× | 1.16× | 1.33× |
| 100,000 integers | 5.11 ms | 2.67× | 2.16× | 2.90× |
| 1,000,000 integers | 60.15 ms | 2.92× | 2.88× | 4.06× |
| 10,000,000 integers | 681.71 ms | 3.15× | 3.14× | 4.42× |

Allocation cost accompanies the speed comparison:

| Input / workload | Go sort.Ints allocated | Sorty allocated vs baseline ↓ | Pargo allocated vs baseline ↓ | Parsort allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| 10,000 integers | 80 KiB | 1.01× | 1.04× | 5.04× |
| 100,000 integers | 784.02 KiB | 1.00× | 1.03× | 5.09× |
| 1,000,000 integers | 7.63 MiB | 1.00× | 1.04× | 5.01× |
| 10,000,000 integers | 76.3 MiB | 1.00× | 1.04× | 5.00× |

The lowest median shifts from Sorty at 10,000 integers to Parsort at 10,000,000 integers.

At 10,000,000 integers, Parsort has the lowest recorded median time: 4.42× the speed of Go sort.Ints, with 5.00× its allocated bytes per operation.

<details>
<summary>Absolute measurements, allocation counts, and observed ranges</summary>

#### 10,000 integers

| Implementation | Time/sort, µs ↓ | Speed vs baseline ↑ | Allocated/sort, KiB ↓ | Allocs/sort ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Go sort.Ints (baseline)** | 395.44 | 1.00× | 80 | 1 |
| Sorty | 223.94 | 1.77× | 80.49 | 10 |
| Pargo | 342.35 | 1.16× | 82.88 | 103 |
| Parsort | 296.93 | 1.33× | 402.88 | 89 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Go sort.Ints | 393.41–396.28 µs | 80–80 KiB | 1 | 1.00× |
| Sorty | 208.53–225.13 µs | 80.42–80.5 KiB | 9–10 | 1.01× |
| Pargo | 326.5–355.91 µs | 82.2–82.88 KiB | 79–103 | 1.04× |
| Parsort | 296.02–297.29 µs | 402.88–402.88 KiB | 89 | 5.04× |

#### 100,000 integers

| Implementation | Time/sort, ms ↓ | Speed vs baseline ↑ | Allocated/sort, KiB ↓ | Allocs/sort ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Go sort.Ints (baseline)** | 5.11 | 1.00× | 784.02 | 1 |
| Sorty | 1.92 | 2.67× | 784.74 | 15 |
| Pargo | 2.36 | 2.16× | 811.25 | 960 |
| Parsort | 1.76 | 2.90× | 3,986.89 | 89 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Go sort.Ints | 5.07–5.11 ms | 784.01–784.03 KiB | 1 | 1.00× |
| Sorty | 1.88–2 ms | 784.7–784.83 KiB | 15–17 | 1.00× |
| Pargo | 2.34–2.68 ms | 811.21–811.38 KiB | 959–967 | 1.03× |
| Parsort | 1.7–1.85 ms | 3.89–3.89 MiB | 89 | 5.09× |

#### 1,000,000 integers

| Implementation | Time/sort, ms ↓ | Speed vs baseline ↑ | Allocated/sort, MiB ↓ | Allocs/sort ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Go sort.Ints (baseline)** | 60.15 | 1.00× | 7.63 | 1 |
| Sorty | 20.59 | 2.92× | 7.64 | 48 |
| Pargo | 20.87 | 2.88× | 7.9 | 9,815 |
| Parsort | 14.81 | 4.06× | 38.25 | 89 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Go sort.Ints | 59.8–60.36 ms | 7.63–7.63 MiB | 1 | 1.00× |
| Sorty | 18.97–20.76 ms | 7.63–7.64 MiB | 46–57 | 1.00× |
| Pargo | 20.55–21.05 ms | 7.9–7.9 MiB | 9,641–9,827 | 1.04× |
| Parsort | 14.74–14.87 ms | 38.25–38.25 MiB | 89 | 5.01× |

#### 10,000,000 integers

| Implementation | Time/sort, ms ↓ | Speed vs baseline ↑ | Allocated/sort, MiB ↓ | Allocs/sort ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Go sort.Ints (baseline)** | 681.71 | 1.00× | 76.3 | 1 |
| Sorty | 216.07 | 3.15× | 76.3 | 130 |
| Pargo | 216.94 | 3.14× | 79.02 | 98,145 |
| Parsort | 154.33 | 4.42× | 381.52 | 89 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Go sort.Ints | 674.17–686.17 ms | 76.3–76.3 MiB | 1 | 1.00× |
| Sorty | 212.51–224.57 ms | 76.3–76.3 MiB | 83–161 | 1.00× |
| Pargo | 212.39–219.85 ms | 79.01–79.03 MiB | 97,916–98,146 | 1.04× |
| Parsort | 152.68–156.52 ms | 381.52–381.52 MiB | 89 | 5.00× |


</details>

### Set one string field

Set one string field through its direct method or a prepared reflect.Value.

Baseline: **Direct method = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

Field lookup and reflect.Value preparation happen before timing; this measures field access only.

| Implementation | Time/access, ns ↓ | Speed vs baseline ↑ | Allocated/access, B ↓ | Allocs/access ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Direct method (baseline)** | 0.89 | 1.00× | 0 | 0 |
| Reflection | 2.21 | 0.40× | 0 | 0 |

Direct method has the lowest recorded median time in this comparison.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Direct method | 0.89–0.89 ns | 0–0 B | 0 | same (0 B) |
| Reflection | 2.21–2.21 ns | 0–0 B | 0 | same (0 B) |

</details>

### Read one string field

Read one string field through its direct method or a prepared reflect.Value.

Baseline: **Direct method = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

Field lookup and reflect.Value preparation happen before timing; this measures field access only.

| Implementation | Time/access, ns ↓ | Speed vs baseline ↑ | Allocated/access, B ↓ | Allocs/access ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Direct method (baseline)** | 0.66 | 1.00× | 0 | 0 |
| Reflection | 0.88 | 0.75× | 0 | 0 |

Direct method has the lowest recorded median time in this comparison.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Direct method | 0.66–0.66 ns | 0–0 B | 0 | same (0 B) |
| Reflection | 0.88–0.88 ns | 0–0 B | 0 | same (0 B) |

</details>

### Synchronous shared-counter increments

Parallel workers increment one exact shared count and wait for each update to complete.

Baseline: **Atomic = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

Parallel time per operation is elapsed time divided by all workers' operations, not individual-call latency. CAS backoff includes scheduler yields.

| Implementation | Time/increment, ns ↓ | Speed vs baseline ↑ | Allocated/increment, B ↓ | Allocs/increment ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Atomic (baseline)** | 12.57 | 1.00× | 0 | 0 |
| Mutex | 17.76 | 0.71× | 0 | 0 |
| RWMutex | 15.72 | 0.80× | 0 | 0 |
| CAS spin | 87.95 | 0.14× | 0 | 0 |
| CAS with backoff | 3.86 | 3.25× | 0 | 0 |
| Flat combining | 149.3 | 0.08× | 12 | 0 |
| RCU copy-on-write | 175 | 0.07× | 81 | 10 |

CAS with backoff has the lowest recorded median time: 3.25× the speed of Atomic; both allocate 0 B per operation.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Atomic | 12.53–12.57 ns | 0–0 B | 0 | same (0 B) |
| Mutex | 17.72–17.81 ns | 0–0 B | 0 | same (0 B) |
| RWMutex | 15.69–15.76 ns | 0–0 B | 0 | same (0 B) |
| CAS spin | 87.57–88.09 ns | 0–0 B | 0 | same (0 B) |
| CAS with backoff | 3.85–3.87 ns | 0–0 B | 0 | same (0 B) |
| Flat combining | 142.5–190.1 ns | 12–13 B | 0 | +12 B |
| RCU copy-on-write | 169.1–178.3 ns | 79–84 B | 9–10 | +81 B |

</details>

### Shared-counter reads

Parallel workers read a shared counter without concurrent writers.

Baseline: **Atomic = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

Parallel time per operation is elapsed time divided by all workers' operations, not individual-call latency. CAS backoff includes scheduler yields.

| Implementation | Time/read, ns ↓ | Speed vs baseline ↑ | Allocated/read, B ↓ | Allocs/read ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Atomic (baseline)** | 0.067 | 1.00× | 0 | 0 |
| Mutex | 15.24 | <0.01× | 0 | 0 |
| RWMutex | 26.42 | <0.01× | 0 | 0 |
| CAS spin | 0.067 | 1.00× | 0 | 0 |
| CAS with backoff | 0.0673 | 1.00× | 0 | 0 |
| RCU pointer read | 0.0831 | 0.81× | 0 | 0 |

CAS spin has the lowest recorded median time: 1.00× the speed of Atomic; both allocate 0 B per operation. The timing ranges of CAS spin and Atomic overlap; their order is inconclusive.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Atomic | 0.0668–0.0675 ns | 0–0 B | 0 | same (0 B) |
| Mutex | 15.22–15.3 ns | 0–0 B | 0 | same (0 B) |
| RWMutex | 26.39–26.45 ns | 0–0 B | 0 | same (0 B) |
| CAS spin | 0.0666–0.0673 ns | 0–0 B | 0 | same (0 B) |
| CAS with backoff | 0.0669–0.0674 ns | 0–0 B | 0 | same (0 B) |
| RCU pointer read | 0.0741–0.0832 ns | 0–0 B | 0 | same (0 B) |

</details>

### Shared counter: 90% updates, 10% reads

Each worker repeats nine completed increments and one read. One reported operation is one increment or read, not a ten-operation batch.

Baseline: **Atomic = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

Parallel time per operation is elapsed time divided by all workers' operations, not individual-call latency. CAS backoff includes scheduler yields.

| Implementation | Time/op, ns ↓ | Speed vs baseline ↑ | Allocated/op, B ↓ | Allocs/op ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Atomic (baseline)** | 11.55 | 1.00× | 0 | 0 |
| Mutex | 16.9 | 0.68× | 0 | 0 |
| RWMutex | 15.22 | 0.76× | 0 | 0 |
| CAS spin | 78.29 | 0.15× | 0 | 0 |
| CAS with backoff | 5.77 | 2.00× | 0 | 0 |

CAS with backoff has the lowest recorded median time: 2.00× the speed of Atomic; both allocate 0 B per operation.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Atomic | 11.51–11.56 ns | 0–0 B | 0 | same (0 B) |
| Mutex | 16.89–17.13 ns | 0–0 B | 0 | same (0 B) |
| RWMutex | 15.17–15.22 ns | 0–0 B | 0 | same (0 B) |
| CAS spin | 77.83–78.4 ns | 0–0 B | 0 | same (0 B) |
| CAS with backoff | 5.77–5.77 ns | 0–0 B | 0 | same (0 B) |

</details>

### Counter increment submission by channel capacity

Parallel workers submit increment requests to one owner goroutine; return does not confirm that an increment has been applied.

Baseline: **Unbuffered = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

These owner-goroutine operations have different completion boundaries from atomic and locked counter calls.

| Implementation | Time/op, ns ↓ | Speed vs baseline ↑ | Allocated/op, B ↓ | Allocs/op ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Unbuffered (baseline)** | 417.6 | 1.00× | 0 | 0 |
| Capacity 1 | 342.9 | 1.22× | 0 | 0 |
| Capacity 256 | 82.26 | 5.08× | 0 | 0 |
| Capacity 4,096 | 88.43 | 4.72× | 0 | 0 |

Capacity 256 has the lowest recorded median time: 5.08× the speed of Unbuffered; both allocate 0 B per operation.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Unbuffered | 408.9–419.4 ns | 0–0 B | 0 | same (0 B) |
| Capacity 1 | 342.8–345.6 ns | 0–0 B | 0 | same (0 B) |
| Capacity 256 | 81.2–85.16 ns | 0–0 B | 0 | same (0 B) |
| Capacity 4,096 | 87.06–88.79 ns | 0–0 B | 0 | same (0 B) |

</details>

### Counter request/reply reads by channel capacity

Parallel workers send read requests to an owner goroutine and wait for their replies.

Baseline: **Unbuffered = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

These owner-goroutine operations have different completion boundaries from atomic and locked counter calls.

| Implementation | Time/op, ns ↓ | Speed vs baseline ↑ | Allocated/op, B ↓ | Allocs/op ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Unbuffered (baseline)** | 544 | 1.00× | 0 | 0 |
| Capacity 1 | 544.7 | 1.00× | 0 | 0 |
| Capacity 256 | 502.7 | 1.08× | 0 | 0 |
| Capacity 4,096 | 484.3 | 1.12× | 0 | 0 |

Capacity 4,096 has the lowest recorded median time: 1.12× the speed of Unbuffered; both allocate 0 B per operation. The timing ranges of Capacity 4,096 and Capacity 256 overlap; their order is inconclusive.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Unbuffered | 538.2–549.7 ns | 0–0 B | 0 | same (0 B) |
| Capacity 1 | 542.2–549.7 ns | 0–0 B | 0 | same (0 B) |
| Capacity 256 | 478.6–511.1 ns | 0–0 B | 0 | same (0 B) |
| Capacity 4,096 | 479.5–485.3 ns | 0–0 B | 0 | same (0 B) |

</details>

### Channel counter: 90% submissions, 10% reads

Each worker repeats nine increment submissions and one request/reply read. One operation is one submission or read.

Baseline: **Unbuffered = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

These owner-goroutine operations have different completion boundaries from atomic and locked counter calls.

| Implementation | Time/op, ns ↓ | Speed vs baseline ↑ | Allocated/op, B ↓ | Allocs/op ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Unbuffered (baseline)** | 426.8 | 1.00× | 0 | 0 |
| Capacity 1 | 379 | 1.13× | 0 | 0 |
| Capacity 256 | 142.8 | 2.99× | 0 | 0 |
| Capacity 4,096 | 144.9 | 2.95× | 0 | 0 |

Capacity 256 has the lowest recorded median time: 2.99× the speed of Unbuffered; both allocate 0 B per operation.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Unbuffered | 425.9–430.1 ns | 0–0 B | 0 | same (0 B) |
| Capacity 1 | 374.7–383.8 ns | 0–0 B | 0 | same (0 B) |
| Capacity 256 | 140.5–142.9 ns | 0–0 B | 0 | same (0 B) |
| Capacity 4,096 | 143.3–144.9 ns | 0–0 B | 0 | same (0 B) |

</details>

### Single-producer, single-consumer transfer

One operation transfers a value between one producer and one consumer. Both implementations retry full/empty operations with scheduler yields.

Baseline: **Channel = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

A ring buffer reserves one slot, so its usable capacity is one less than the listed channel capacity.

| Input / workload | Channel time (1.00×) | Ring buffer speed ↑ |
| --- | ---: | ---: |
| 256 slots | 45.14 ns | 1.22× |
| 4,096 slots | 45.05 ns | 1.64× |
| 65,536 slots | 45.12 ns | 1.62× |

Allocation cost accompanies the speed comparison:

| Input / workload | Channel allocated | Ring buffer allocated vs baseline ↓ |
| --- | ---: | ---: |
| 256 slots | 0 B | same (0 B) |
| 4,096 slots | 0 B | same (0 B) |
| 65,536 slots | 0 B | same (0 B) |

At 65,536 slots, Ring buffer has the lowest recorded median time: 1.62× the speed of Channel; both allocate 0 B per operation.

<details>
<summary>Absolute measurements, allocation counts, and observed ranges</summary>

#### 256 slots

| Implementation | Time/transfer, ns ↓ | Speed vs baseline ↑ | Allocated/transfer, B ↓ | Allocs/transfer ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Channel (baseline)** | 45.14 | 1.00× | 0 | 0 |
| Ring buffer | 37.08 | 1.22× | 0 | 0 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Channel | 44.88–45.33 ns | 0–0 B | 0 | same (0 B) |
| Ring buffer | 35.8–39.25 ns | 0–0 B | 0 | same (0 B) |

#### 4,096 slots

| Implementation | Time/transfer, ns ↓ | Speed vs baseline ↑ | Allocated/transfer, B ↓ | Allocs/transfer ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Channel (baseline)** | 45.05 | 1.00× | 0 | 0 |
| Ring buffer | 27.49 | 1.64× | 0 | 0 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Channel | 44.82–45.15 ns | 0–0 B | 0 | same (0 B) |
| Ring buffer | 26.82–33.38 ns | 0–0 B | 0 | same (0 B) |

#### 65,536 slots

| Implementation | Time/transfer, ns ↓ | Speed vs baseline ↑ | Allocated/transfer, B ↓ | Allocs/transfer ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Channel (baseline)** | 45.12 | 1.00× | 0 | 0 |
| Ring buffer | 27.92 | 1.62× | 0 | 0 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Channel | 44.86–45.39 ns | 0–0 B | 0 | same (0 B) |
| Ring buffer | 26.75–32.66 ns | 0–0 B | 0 | same (0 B) |


</details>

### Concurrent map stores

Parallel workers store values across 1,024 shared string keys; timing includes first insertion and subsequent replacement.

Baseline: **Mutex map = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

Results describe this shared-key workload; time per access is an aggregate-throughput measure.

| Implementation | Time/access, ns ↓ | Speed vs baseline ↑ | Allocated/access, B ↓ | Allocs/access ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Mutex map (baseline)** | 44.29 | 1.00× | 0 | 0 |
| sync.Map | 40.17 | 1.10× | 72 | 2 |

sync.Map has the lowest recorded median time: 1.10× the speed of Mutex map, adding 72 B allocated per operation.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Mutex map | 43.66–44.73 ns | 0–0 B | 0 | same (0 B) |
| sync.Map | 39.88–40.59 ns | 72–72 B | 2 | +72 B |

</details>

### Concurrent map loads

Parallel workers load values across 1,024 shared string keys populated before timing.

Baseline: **Mutex map = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

Results describe this shared-key workload; time per access is an aggregate-throughput measure.

| Implementation | Time/access, ns ↓ | Speed vs baseline ↑ | Allocated/access, B ↓ | Allocs/access ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Mutex map (baseline)** | 37.31 | 1.00× | 0 | 0 |
| sync.Map | 2.77 | 13.46× | 0 | 0 |

sync.Map has the lowest recorded median time: 13.46× the speed of Mutex map; both allocate 0 B per operation.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Mutex map | 37.26–37.34 ns | 0–0 B | 0 | same (0 B) |
| sync.Map | 2.77–2.77 ns | 0–0 B | 0 | same (0 B) |

</details>

### Sharded counter updates

Parallel workers increment their selected padded atomic shard or slot.

Baseline: **8 shards = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

Workers are assigned by ID, not pinned to physical CPUs. Reading the sum is separate work and is not an atomic snapshot during concurrent writes.

| Implementation | Time/increment, ns ↓ | Speed vs baseline ↑ | Allocated/increment, B ↓ | Allocs/increment ↓ |
| --- | ---: | ---: | ---: | ---: |
| **8 shards (baseline)** | 1.53 | 1.00× | 0 | 0 |
| One shard per logical CPU | 0.24 | 6.27× | 0 | 0 |
| GOMAXPROCS slots | 0.24 | 6.28× | 0 | 0 |

GOMAXPROCS slots has the lowest recorded median time: 6.28× the speed of 8 shards; both allocate 0 B per operation. The timing ranges of GOMAXPROCS slots and One shard per logical CPU overlap; their order is inconclusive.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| 8 shards | 1.53–1.58 ns | 0–0 B | 0 | same (0 B) |
| One shard per logical CPU | 0.24–0.24 ns | 0–0 B | 0 | same (0 B) |
| GOMAXPROCS slots | 0.24–0.25 ns | 0–0 B | 0 | same (0 B) |

</details>

### Read a sharded total

Parallel readers sum padded atomic slots without concurrent writers.

Baseline: **8 shards = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

Larger shard arrays change the amount of read work. These measurements do not establish snapshot consistency with active writers.

| Implementation | Time/sum, ns ↓ | Speed vs baseline ↑ | Allocated/sum, B ↓ | Allocs/sum ↓ |
| --- | ---: | ---: | ---: | ---: |
| **8 shards (baseline)** | 0.65 | 1.00× | 0 | 0 |
| GOMAXPROCS slots | 1.21 | 0.54× | 0 | 0 |

8 shards has the lowest recorded median time in this comparison.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| 8 shards | 0.65–0.65 ns | 0–0 B | 0 | same (0 B) |
| GOMAXPROCS slots | 1.21–1.21 ns | 0–0 B | 0 | same (0 B) |

</details>

### Buffered local increments

Workers accumulate increments locally and publish them at the configured threshold, then flush their remainders before finishing.

Baseline: **Flush every 16 = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

A larger threshold changes publication frequency and visibility. Reads omit unflushed increments.

| Implementation | Time/increment, ns ↓ | Speed vs baseline ↑ | Allocated/increment, B ↓ | Allocs/increment ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Flush every 16 (baseline)** | 2.15 | 1.00× | 0 | 0 |
| Flush every 64 | 0.72 | 2.98× | 0 | 0 |
| Flush every 256 | 0.32 | 6.67× | 0 | 0 |

Flush every 256 has the lowest recorded median time: 6.67× the speed of Flush every 16; both allocate 0 B per operation.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Flush every 16 | 2.14–2.15 ns | 0–0 B | 0 | same (0 B) |
| Flush every 64 | 0.72–0.72 ns | 0–0 B | 0 | same (0 B) |
| Flush every 256 | 0.32–0.32 ns | 0–0 B | 0 | same (0 B) |

</details>

### MPSC counter submission

Parallel producers submit increments through a buffered channel to one aggregator.

Baseline: **Capacity 256 = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

Submission does not wait for the aggregate total to be published. These are not completed counter-update timings.

| Implementation | Time/submission, ns ↓ | Speed vs baseline ↑ | Allocated/submission, B ↓ | Allocs/submission ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Capacity 256 (baseline)** | 118.3 | 1.00× | 0 | 0 |
| Capacity 4,096 | 92.17 | 1.28× | 0 | 0 |

Capacity 4,096 has the lowest recorded median time: 1.28× the speed of Capacity 256; both allocate 0 B per operation.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Capacity 256 | 117.7–118.5 ns | 0–0 B | 0 | same (0 B) |
| Capacity 4,096 | 91.35–92.54 ns | 0–0 B | 0 | same (0 B) |

</details>

### Approximate counter update costs

Workers sample updates with a local pseudo-random state; an accepted update adds the configured sampling rate.

No equivalent reference measured; absolute costs are shown. Median of 3 samples; GOMAXPROCS=16.

Sampling rates trade execution cost for estimation error. Estimator accuracy was not measured, so no relative-speed ranking is assigned.

| Operation / variant | Time/attempt, ns ↓ | Allocated/attempt, B ↓ | Allocs/attempt ↓ |
| --- | ---: | ---: | ---: |
| Sample 1 in 10 | 3.37 | 0 | 0 |
| Sample 1 in 100 | 0.78 | 0 | 0 |
| Sample 1 in 1,000 | 0.35 | 0 | 0 |

<details>
<summary>Observed sample ranges</summary>

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| Sample 1 in 10 | 3.36–3.37 ns | 0–0 B | 0 |
| Sample 1 in 100 | 0.78–0.78 ns | 0–0 B | 0 |
| Sample 1 in 1,000 | 0.35–0.35 ns | 0–0 B | 0 |

</details>

### Read an approximate counter

Parallel readers load the current estimate from a counter configured to sample one update in 100; there are no concurrent writers.

No equivalent reference measured; absolute costs are shown. Median of 3 samples; GOMAXPROCS=16.

This reads an estimate; no accuracy-equivalent alternative was measured.

| Operation / variant | Time/read, ns ↓ | Allocated/read, B ↓ | Allocs/read ↓ |
| --- | ---: | ---: | ---: |
| Sample rate 100 | 0.0671 | 0 | 0 |

<details>
<summary>Observed sample ranges</summary>

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| Sample rate 100 | 0.0668–0.0673 ns | 0–0 B | 0 |

</details>

### Four-producer queue transfer

Four producers publish messages and one consumer receives all of them through a Disruptor-style ring. One operation is one complete transfer.

Baseline: **256 slots = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

Timing includes goroutine startup, retries, and completion waiting. This is separate from channel submission-only measurements.

| Implementation | Time/transfer, µs ↓ | Speed vs baseline ↑ | Allocated/transfer, B ↓ | Allocs/transfer ↓ |
| --- | ---: | ---: | ---: | ---: |
| **256 slots (baseline)** | 2.45 | 1.00× | 0 | 0 |
| 4,096 slots | 2.46 | 1.00× | 0 | 0 |
| 65,536 slots | 2.45 | 1.00× | 0 | 0 |

65,536 slots has the lowest recorded median time: 1.00× the speed of 256 slots; both allocate 0 B per operation. The timing ranges of 65,536 slots and 256 slots overlap; their order is inconclusive.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| 256 slots | 2.45–2.46 µs | 0–0 B | 0 | same (0 B) |
| 4,096 slots | 2.45–2.46 µs | 0–0 B | 0 | same (0 B) |
| 65,536 slots | 2.45–2.45 µs | 0–0 B | 0 | same (0 B) |

</details>

### Semaphore acquire and release

Parallel workers acquire and release a channel-based semaphore without any protected work between the calls.

Baseline: **Capacity 1 = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

Capacity changes the permitted concurrency. This isolates semaphore overhead, not useful protected-work throughput.

| Implementation | Time/acquire/release, ns ↓ | Speed vs baseline ↑ | Allocated/acquire/release, B ↓ | Allocs/acquire/release ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Capacity 1 (baseline)** | 438.9 | 1.00× | 0 | 0 |
| Capacity 4 | 292.8 | 1.50× | 0 | 0 |
| Capacity 16 | 45.95 | 9.55× | 0 | 0 |
| Capacity 64 | 44.27 | 9.91× | 0 | 0 |

Capacity 64 has the lowest recorded median time: 9.91× the speed of Capacity 1; both allocate 0 B per operation.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Capacity 1 | 437.1–440.5 ns | 0–0 B | 0 | same (0 B) |
| Capacity 4 | 292.5–293.9 ns | 0–0 B | 0 | same (0 B) |
| Capacity 16 | 45.21–46.09 ns | 0–0 B | 0 | same (0 B) |
| Capacity 64 | 43.64–44.62 ns | 0–0 B | 0 | same (0 B) |

</details>

### Start goroutines and wait

One operation creates the listed number of goroutines with empty work and waits for all of them to finish.

No equivalent reference measured; absolute costs are shown. Median of 3 samples; GOMAXPROCS=16.

Rows perform different amounts of startup work; the time is per complete batch, not per goroutine.

| Operation / variant | Time/batch, µs ↓ | Allocated/batch, B ↓ | Allocs/batch ↓ |
| --- | ---: | ---: | ---: |
| 1 goroutine | 1.23 | 40 | 2 |
| 4 goroutines | 3.32 | 112 | 5 |
| 16 goroutines | 10.31 | 400 | 17 |
| 64 goroutines | 32.76 | 1,552 | 65 |

<details>
<summary>Observed sample ranges</summary>

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| 1 goroutine | 1.23–1.23 µs | 40–40 B | 2 |
| 4 goroutines | 3.31–3.32 µs | 112–112 B | 5 |
| 16 goroutines | 10.31–10.36 µs | 400–400 B | 17 |
| 64 goroutines | 32.73–32.98 µs | 1.52–1.52 KiB | 65 |

</details>

### Pool get and put

Get a reusable buffer from sync.Pool and return it, serially or across parallel workers.

No equivalent reference measured; absolute costs are shown. Median of 3 samples; GOMAXPROCS=16.

The serial case measures a sequential operation; the parallel case reports elapsed time divided by total worker operations. Their ratio is not per-call speedup.

| Operation / variant | Time/get/put, ns ↓ | Allocated/get/put, B ↓ | Allocs/get/put ↓ |
| --- | ---: | ---: | ---: |
| Serial | 11 | 0 | 0 |
| Parallel | 1.5 | 0 | 0 |

<details>
<summary>Observed sample ranges</summary>

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| Serial | 10.99–11.01 ns | 0–0 B | 0 |
| Parallel | 1.5–1.51 ns | 0–0 B | 0 |

</details>

### Once initialization paths

Initialize and read a fresh instance, reuse one shared instance across workers, or reuse an instance initialized before timing.

No equivalent reference measured; absolute costs are shown. Median of 3 samples; GOMAXPROCS=16.

The shared instance initializes only once across the benchmark. These initialization states and operation boundaries are different.

| Operation / variant | Time/op, ns ↓ | Allocated/op, B ↓ | Allocs/op ↓ |
| --- | ---: | ---: | ---: |
| Fresh instance, serial | 18.49 | 24 | 1 |
| Shared instance, parallel | 0.26 | 0 | 0 |
| Preinitialized, parallel | 0.26 | 0 | 0 |

<details>
<summary>Observed sample ranges</summary>

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| Fresh instance, serial | 18.43–18.65 ns | 24–24 B | 1 |
| Shared instance, parallel | 0.26–0.26 ns | 0–0 B | 0 |
| Preinitialized, parallel | 0.26–0.26 ns | 0–0 B | 0 |

</details>

**Recorded samples:** [meta](../results/meta.txt)

<!-- benchmark-results:end -->

## Methodology

### Reflection

Read or write one string field through a direct method or a prepared
`reflect.Value`. `FieldByName` runs before timing, so these cases measure access
rather than field lookup. Both paths use `testing.B.Loop`.

### Sorting

Sort 10,000 to 10,000,000 random integers with Go `sort.Ints`, Sorty v2, Pargo,
and Parsort. **Allocation and input copying are included in every operation.**
Each implementation generates its own input before timing, with the same size
and distribution.

Sorty uses `SortSlice` and its default `MaxGor=3`, including the caller. Library
defaults determine worker counts; GOMAXPROCS is not a promise that every sorter
uses that many workers.

### Synchronization

Counters, channels, queues, maps, and runtime primitives perform different work.
Synchronous updates, deferred publication, approximate counts, and initialization
states remain distinct in the result groups. See [operation definitions and
comparison boundaries](CONTENTION.md).

Parallel elapsed time is wall time divided by the operations completed across
workers. It describes aggregate throughput, not the latency of one call.

The implementations and measurement loops are in [reflect_test.go](reflect_test.go),
[sort_test.go](sort_test.go), and [contention_test.go](contention_test.go).

## Running

Requires Go 1.27.4+. From the repository root, without competing benchmark work:

```sh
GOMAXPROCS=16 go test -p=1 -run '^$' -bench . -benchmem -benchtime=1s -count=3 ./meta
```

The [full runner](../scripts/benchmark.bash) records the environment and dependency
versions. Regenerate the comparisons from recorded samples with:

```sh
bash scripts/results.bash render meta
```
