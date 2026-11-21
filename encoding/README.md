# Encoding benchmarks

Compare standard Base64 and ASCII85 encoding across payload sizes.
**Base64 is the reference at each size.** Each operation allocates an output
buffer and encodes a prepared byte slice; decoding is excluded.

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

**Workloads:** [Encode by input size](#encode-by-input-size)

### Encode by input size

One operation allocates an output buffer and encodes a prepared printable-ASCII input. Decoding and input generation are excluded.

Baseline: **Base64 = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

The formats have different output expansion. Allocated bytes include buffer capacity and are not a measurement of encoded length.

| Input / workload | Base64 time (1.00×) | ASCII85 speed ↑ |
| --- | ---: | ---: |
| 2 KiB | 1.59 µs | 0.43× |
| 4 KiB | 3.15 µs | 0.43× |
| 8 KiB | 6.25 µs | 0.43× |
| 16 KiB | 12.27 µs | 0.41× |
| 32 KiB | 24.73 µs | 0.40× |
| 64 KiB | 48.7 µs | 0.38× |
| 128 KiB | 96.83 µs | 0.37× |
| 256 KiB | 194.46 µs | 0.42× |
| 512 KiB | 390.76 µs | 0.41× |
| 1 MiB | 774.62 µs | 0.42× |
| 10 MiB | 7.02 ms | 0.40× |
| 100 MiB | 68.9 ms | 0.40× |
| 1 GiB | 708.85 ms | 0.40× |

Allocation cost accompanies the speed comparison:

| Input / workload | Base64 allocated | ASCII85 allocated vs baseline ↓ |
| --- | ---: | ---: |
| 2 KiB | 3 KiB | 0.88× |
| 4 KiB | 6 KiB | 0.88× |
| 8 KiB | 12 KiB | 0.83× |
| 16 KiB | 24 KiB | 0.83× |
| 32 KiB | 48 KiB | 0.83× |
| 64 KiB | 88 KiB | 0.91× |
| 128 KiB | 176 KiB | 0.91× |
| 256 KiB | 344 KiB | 0.93× |
| 512 KiB | 688 KiB | 0.93× |
| 1 MiB | 1.34 MiB | 0.94× |
| 10 MiB | 13.34 MiB | 0.94× |
| 100 MiB | 133.34 MiB | 0.94× |
| 1 GiB | 1.33 GiB | 0.94× |

At 1 GiB, Base64 has the lowest recorded median time in this comparison.

<details>
<summary>Absolute measurements, allocation counts, and observed ranges</summary>

#### 2 KiB

| Implementation | Time/encoding, µs ↓ | Speed vs baseline ↑ | Allocated/encoding, KiB ↓ | Allocs/encoding ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Base64 (baseline)** | 1.59 | 1.00× | 3 | 1 |
| ASCII85 | 3.69 | 0.43× | 2.62 | 1 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Base64 | 1.58–1.59 µs | 3–3 KiB | 1 | 1.00× |
| ASCII85 | 3.68–3.71 µs | 2.62–2.62 KiB | 1 | 0.88× |

#### 4 KiB

| Implementation | Time/encoding, µs ↓ | Speed vs baseline ↑ | Allocated/encoding, KiB ↓ | Allocs/encoding ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Base64 (baseline)** | 3.15 | 1.00× | 6 | 1 |
| ASCII85 | 7.42 | 0.43× | 5.25 | 1 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Base64 | 3.15–3.15 µs | 6–6 KiB | 1 | 1.00× |
| ASCII85 | 7.4–7.46 µs | 5.25–5.25 KiB | 1 | 0.88× |

#### 8 KiB

| Implementation | Time/encoding, µs ↓ | Speed vs baseline ↑ | Allocated/encoding, KiB ↓ | Allocs/encoding ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Base64 (baseline)** | 6.25 | 1.00× | 12 | 1 |
| ASCII85 | 14.6 | 0.43× | 10 | 1 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Base64 | 6.24–6.26 µs | 12–12 KiB | 1 | 1.00× |
| ASCII85 | 14.54–14.64 µs | 10–10 KiB | 1 | 0.83× |

#### 16 KiB

| Implementation | Time/encoding, µs ↓ | Speed vs baseline ↑ | Allocated/encoding, KiB ↓ | Allocs/encoding ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Base64 (baseline)** | 12.27 | 1.00× | 24 | 1 |
| ASCII85 | 29.95 | 0.41× | 20 | 1 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Base64 | 12.26–12.36 µs | 24–24 KiB | 1 | 1.00× |
| ASCII85 | 29.9–30.13 µs | 20–20 KiB | 1 | 0.83× |

#### 32 KiB

| Implementation | Time/encoding, µs ↓ | Speed vs baseline ↑ | Allocated/encoding, KiB ↓ | Allocs/encoding ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Base64 (baseline)** | 24.73 | 1.00× | 48 | 1 |
| ASCII85 | 61.08 | 0.40× | 40 | 1 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Base64 | 24.72–24.74 µs | 48–48 KiB | 1 | 1.00× |
| ASCII85 | 60.69–61.13 µs | 40–40 KiB | 1 | 0.83× |

#### 64 KiB

| Implementation | Time/encoding, µs ↓ | Speed vs baseline ↑ | Allocated/encoding, KiB ↓ | Allocs/encoding ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Base64 (baseline)** | 48.7 | 1.00× | 88 | 1 |
| ASCII85 | 127.84 | 0.38× | 80 | 1 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Base64 | 48.65–48.8 µs | 88–88 KiB | 1 | 1.00× |
| ASCII85 | 126.99–128.15 µs | 80–80 KiB | 1 | 0.91× |

#### 128 KiB

| Implementation | Time/encoding, µs ↓ | Speed vs baseline ↑ | Allocated/encoding, KiB ↓ | Allocs/encoding ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Base64 (baseline)** | 96.83 | 1.00× | 176 | 1 |
| ASCII85 | 260.05 | 0.37× | 160 | 1 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Base64 | 96.48–96.96 µs | 176–176 KiB | 1 | 1.00× |
| ASCII85 | 259.1–264.98 µs | 160–160 KiB | 1 | 0.91× |

#### 256 KiB

| Implementation | Time/encoding, µs ↓ | Speed vs baseline ↑ | Allocated/encoding, KiB ↓ | Allocs/encoding ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Base64 (baseline)** | 194.46 | 1.00× | 344 | 1 |
| ASCII85 | 464.12 | 0.42× | 320 | 1 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Base64 | 193.81–194.57 µs | 344–344 KiB | 1 | 1.00× |
| ASCII85 | 462.01–466.2 µs | 320–320.01 KiB | 1 | 0.93× |

#### 512 KiB

| Implementation | Time/encoding, µs ↓ | Speed vs baseline ↑ | Allocated/encoding, KiB ↓ | Allocs/encoding ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Base64 (baseline)** | 390.76 | 1.00× | 688 | 1 |
| ASCII85 | 946.23 | 0.41× | 640.01 | 1 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Base64 | 389.97–391.88 µs | 688–688 KiB | 1 | 1.00× |
| ASCII85 | 943.88–948.76 µs | 640.01–640.01 KiB | 1 | 0.93× |

#### 1 MiB

| Implementation | Time/encoding, µs ↓ | Speed vs baseline ↑ | Allocated/encoding, MiB ↓ | Allocs/encoding ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Base64 (baseline)** | 774.62 | 1.00× | 1.34 | 1 |
| ASCII85 | 1,841.29 | 0.42× | 1.25 | 1 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Base64 | 773.07–780.66 µs | 1.34–1.34 MiB | 1 | 1.00× |
| ASCII85 | 1.84–1.85 ms | 1.25–1.25 MiB | 1 | 0.94× |

#### 10 MiB

| Implementation | Time/encoding, ms ↓ | Speed vs baseline ↑ | Allocated/encoding, MiB ↓ | Allocs/encoding ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Base64 (baseline)** | 7.02 | 1.00× | 13.34 | 1 |
| ASCII85 | 17.48 | 0.40× | 12.5 | 1 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Base64 | 7–7.02 ms | 13.34–13.34 MiB | 1 | 1.00× |
| ASCII85 | 17.34–17.56 ms | 12.5–12.5 MiB | 1 | 0.94× |

#### 100 MiB

| Implementation | Time/encoding, ms ↓ | Speed vs baseline ↑ | Allocated/encoding, MiB ↓ | Allocs/encoding ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Base64 (baseline)** | 68.9 | 1.00× | 133.34 | 1 |
| ASCII85 | 172.94 | 0.40× | 125 | 1 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Base64 | 68.84–68.93 ms | 133.34–133.34 MiB | 1 | 1.00× |
| ASCII85 | 172.01–173.39 ms | 125–125 MiB | 1 | 0.94× |

#### 1 GiB

| Implementation | Time/encoding, ms ↓ | Speed vs baseline ↑ | Allocated/encoding, GiB ↓ | Allocs/encoding ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Base64 (baseline)** | 708.85 | 1.00× | 1.33 | 2 |
| ASCII85 | 1,761.19 | 0.40× | 1.25 | 3 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Base64 | 708.44–710.28 ms | 1.33–1.33 GiB | 2 | 1.00× |
| ASCII85 | 1.76–1.79 s | 1.25–1.25 GiB | 1–3 | 0.94× |


</details>

**Recorded samples:** [encoding](../results/encoding.txt)

<!-- benchmark-results:end -->

## Methodology

Inputs contain random printable ASCII and are generated before timing. Sizes use
binary units: legacy benchmark suffixes `KB`, `MB`, and `GB` mean KiB, MiB, and GiB.

Allocated bytes include the output allocation; they are not a measurement of the
encoded message length. Output expansion is not reported as a separate metric.
The 100 MiB and 1 GiB cases also retain their input, so their peak memory exceeds
the per-operation allocation shown in the tables.

## Running

Requires Go 1.27.4+ and enough memory for the large payloads. No external service
is required. From the repository root:

```sh
go test -p=1 -run '^$' -bench . -benchmem -benchtime=1s -count=3 ./encoding
```
