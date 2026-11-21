# Cache benchmarks

Compare Memcached and Redis for one set, one get, and a returned-value check.
**Memcached is the reference at each payload size.** Measurements include the
client, local container networking, and server work.

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

**Workloads:** [Set and get by payload size](#set-and-get-by-payload-size)

### Set and get by payload size

One operation sets a prepared value, gets it back, and checks the returned bytes over local container TCP connections.

Baseline: **Memcached = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

Allocated bytes belong to the Go client process; they exclude server memory. Redis persistence is disabled.

| Input / workload | Memcached time (1.00×) | Redis speed ↑ |
| --- | ---: | ---: |
| 2 KiB | 154.28 µs | 0.73× |
| 4 KiB | 173.68 µs | 0.79× |
| 8 KiB | 179.53 µs | 0.84× |
| 16 KiB | 191.21 µs | 0.86× |
| 32 KiB | 208.65 µs | 0.85× |
| 64 KiB | 237.48 µs | 0.77× |
| 128 KiB | 291.65 µs | 0.84× |
| 256 KiB | 396.55 µs | 0.92× |
| 512 KiB | 793.38 µs | 1.33× |

Allocation cost accompanies the speed comparison:

| Input / workload | Memcached allocated | Redis allocated vs baseline ↓ |
| --- | ---: | ---: |
| 2 KiB | 6.5 KiB | 0.42× |
| 4 KiB | 13.01 KiB | 0.40× |
| 8 KiB | 25.52 KiB | 0.38× |
| 16 KiB | 50.3 KiB | 0.37× |
| 32 KiB | 104.36 KiB | 0.39× |
| 64 KiB | 200.48 KiB | 0.36× |
| 128 KiB | 392.65 KiB | 0.35× |
| 256 KiB | 777.08 KiB | 0.34× |
| 512 KiB | 1.51 MiB | 0.34× |

The lowest median shifts from Memcached at 2 KiB to Redis at 512 KiB.

At 512 KiB, Redis has the lowest recorded median time: 1.33× the speed of Memcached, with 0.34× its allocated bytes per operation.

<details>
<summary>Absolute measurements, allocation counts, and observed ranges</summary>

#### 2 KiB

| Implementation | Time/set/get, µs ↓ | Speed vs baseline ↑ | Allocated/set/get, KiB ↓ | Allocs/set/get ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Memcached (baseline)** | 154.28 | 1.00× | 6.5 | 10 |
| Redis | 212.48 | 0.73× | 2.74 | 15 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Memcached | 153.29–154.97 µs | 6.48–6.5 KiB | 10 | 1.00× |
| Redis | 211.91–212.55 µs | 2.74–2.74 KiB | 15 | 0.42× |

#### 4 KiB

| Implementation | Time/set/get, µs ↓ | Speed vs baseline ↑ | Allocated/set/get, KiB ↓ | Allocs/set/get ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Memcached (baseline)** | 173.68 | 1.00× | 13.01 | 10 |
| Redis | 220.56 | 0.79× | 5.24 | 15 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Memcached | 173.62–175.72 µs | 13.01–13.01 KiB | 10 | 1.00× |
| Redis | 220.02–224.83 µs | 5.24–5.24 KiB | 15 | 0.40× |

#### 8 KiB

| Implementation | Time/set/get, µs ↓ | Speed vs baseline ↑ | Allocated/set/get, KiB ↓ | Allocs/set/get ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Memcached (baseline)** | 179.53 | 1.00× | 25.52 | 10 |
| Redis | 213.96 | 0.84× | 9.74 | 15 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Memcached | 178.66–180.4 µs | 25.52–25.52 KiB | 10 | 1.00× |
| Redis | 213.31–217.87 µs | 9.74–9.74 KiB | 15 | 0.38× |

#### 16 KiB

| Implementation | Time/set/get, µs ↓ | Speed vs baseline ↑ | Allocated/set/get, KiB ↓ | Allocs/set/get ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Memcached (baseline)** | 191.21 | 1.00× | 50.3 | 10 |
| Redis | 222.3 | 0.86× | 18.49 | 15 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Memcached | 190.1–192.25 µs | 50.3–50.3 KiB | 10 | 1.00× |
| Redis | 222.18–222.74 µs | 18.49–18.49 KiB | 15 | 0.37× |

#### 32 KiB

| Implementation | Time/set/get, µs ↓ | Speed vs baseline ↑ | Allocated/set/get, KiB ↓ | Allocs/set/get ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Memcached (baseline)** | 208.65 | 1.00× | 104.36 | 10 |
| Redis | 246.29 | 0.85× | 40.49 | 15 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Memcached | 207.62–209.16 µs | 104.36–104.36 KiB | 10 | 1.00× |
| Redis | 245.31–246.49 µs | 40.49–40.49 KiB | 15 | 0.39× |

#### 64 KiB

| Implementation | Time/set/get, µs ↓ | Speed vs baseline ↑ | Allocated/set/get, KiB ↓ | Allocs/set/get ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Memcached (baseline)** | 237.48 | 1.00× | 200.48 | 10 |
| Redis | 309.82 | 0.77× | 72.49 | 15 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Memcached | 237.14–237.63 µs | 200.48–200.48 KiB | 10 | 1.00× |
| Redis | 308.36–309.88 µs | 72.49–72.49 KiB | 15 | 0.36× |

#### 128 KiB

| Implementation | Time/set/get, µs ↓ | Speed vs baseline ↑ | Allocated/set/get, KiB ↓ | Allocs/set/get ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Memcached (baseline)** | 291.65 | 1.00× | 392.65 | 10 |
| Redis | 347.82 | 0.84× | 136.49 | 15 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Memcached | 287.6–294.6 µs | 392.64–392.67 KiB | 10–11 | 1.00× |
| Redis | 347.03–350.44 µs | 136.49–136.49 KiB | 15 | 0.35× |

#### 256 KiB

| Implementation | Time/set/get, µs ↓ | Speed vs baseline ↑ | Allocated/set/get, KiB ↓ | Allocs/set/get ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Memcached (baseline)** | 396.55 | 1.00× | 777.08 | 11 |
| Redis | 432.28 | 0.92× | 264.49 | 15 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Memcached | 392.74–401.09 µs | 777.07–777.08 KiB | 11 | 1.00× |
| Redis | 429.69–433.44 µs | 264.49–264.49 KiB | 15 | 0.34× |

#### 512 KiB

| Implementation | Time/set/get, µs ↓ | Speed vs baseline ↑ | Allocated/set/get, KiB ↓ | Allocs/set/get ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Memcached (baseline)** | 793.38 | 1.00× | 1,545.81 | 13 |
| Redis | 596.42 | 1.33× | 520.49 | 15 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Memcached | 774.24–849.91 µs | 1.51–1.51 MiB | 13 | 1.00× |
| Redis | 590.43–602.79 µs | 520.49–520.49 KiB | 15 | 0.34× |


</details>

**Recorded samples:** [cache](../results/cache.txt)

<!-- benchmark-results:end -->

## Methodology

Each operation reuses a key and a prepared random printable-ASCII payload.
Input generation is excluded. Payloads range from 2 KiB to 512 KiB; the legacy
`KB` suffix in benchmark names means KiB. Requests run serially.

Go allocations describe the client process. They exclude Redis and Memcached
memory, so a lower allocation count is not a measurement of server memory use.

[Test setup](main_test.go) starts disposable, digest-pinned Redis 8.10.3 and
Memcached 1.6.45 containers on dynamically assigned ports. Redis snapshots and
append-only persistence are disabled. Memcached permits items up to 16 MiB so
correctness tests can include larger values. Containers are removed afterward.

## Running

Requires Go 1.27.4+ and a local Docker engine. From the repository root:

```sh
go test -p=1 -run '^$' -bench . -benchmem -benchtime=1s -count=3 ./cache
```
