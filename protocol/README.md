# Local HTTP and TCP benchmarks

Measure sending prepared messages to a server in the same Go process over
loopback. **HTTP and TCP have separate scaling results:** HTTP waits for a
response after the body is read; TCP stops after writing and flushing bytes.
Those completion boundaries do not support a direct protocol speedup claim.

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

**Workloads:** [HTTP payload scaling](#http-payload-scaling) · [TCP payload scaling](#tcp-payload-scaling)

### HTTP payload scaling

POST a prepared body to a loopback server and wait for HTTP 200 after the handler reads the body.

No equivalent reference measured; absolute costs are shown. Median of 3 samples; GOMAXPROCS=16.

Completion includes server acknowledgement. These costs are not normalized against TCP writes, which have no application acknowledgement. Go allocations can include the server goroutines in the same process.

| Input / workload | Time/op ↓ | Allocated/op ↓ | Allocs/op ↓ |
| --- | ---: | ---: | ---: |
| 9 bytes | 80.64 µs | 6.3 KiB | 67 |
| 2 KiB | 83.86 µs | 10.61 KiB | 74 |
| 4 KiB | 93.22 µs | 16.69 KiB | 79 |
| 8 KiB | 98.87 µs | 28.87 KiB | 80 |
| 16 KiB | 110.42 µs | 59.08 KiB | 84 |
| 32 KiB | 128.46 µs | 123.96 KiB | 87 |
| 64 KiB | 145.47 µs | 182.32 KiB | 89 |
| 128 KiB | 191.69 µs | 347.42 KiB | 94 |
| 256 KiB | 294.11 µs | 674.67 KiB | 100 |
| 512 KiB | 480.12 µs | 1.08 MiB | 104 |
| 1 MiB | 857.61 µs | 2.2 MiB | 111 |
| 10 MiB | 5.47 ms | 22.45 MiB | 117 |
| 100 MiB | 48.42 ms | 240.57 MiB | 124 |
| 1 GiB | 508.82 ms | 2.56 GiB | 132 |

<details>
<summary>Absolute measurements, allocation counts, and observed ranges</summary>

#### 9 bytes

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| HTTP | 80.61–80.97 µs | 6.26–6.34 KiB | 67 |

#### 2 KiB

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| HTTP | 83.62–84.2 µs | 10.59–10.66 KiB | 74 |

#### 4 KiB

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| HTTP | 92.81–93.33 µs | 16.66–16.7 KiB | 79 |

#### 8 KiB

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| HTTP | 98.67–98.88 µs | 28.87–28.88 KiB | 80 |

#### 16 KiB

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| HTTP | 110.41–110.61 µs | 59.01–59.1 KiB | 84 |

#### 32 KiB

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| HTTP | 128.37–128.48 µs | 123.82–123.96 KiB | 87 |

#### 64 KiB

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| HTTP | 145–145.92 µs | 182.08–182.48 KiB | 89 |

#### 128 KiB

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| HTTP | 191.23–191.8 µs | 347.33–347.43 KiB | 93–94 |

#### 256 KiB

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| HTTP | 292.01–295.02 µs | 674.42–674.87 KiB | 99–100 |

#### 512 KiB

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| HTTP | 476.84–481.28 µs | 1.08–1.08 MiB | 104 |

#### 1 MiB

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| HTTP | 849.94–915.05 µs | 2.2–2.2 MiB | 111 |

#### 10 MiB

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| HTTP | 5.44–5.49 ms | 22.45–22.45 MiB | 117 |

#### 100 MiB

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| HTTP | 48.21–48.66 ms | 240.57–240.58 MiB | 124 |

#### 1 GiB

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| HTTP | 507.45–518.85 ms | 2.56–2.56 GiB | 130–132 |


</details>

### TCP payload scaling

Write and flush a prepared payload to a loopback socket while a server drains the stream.

No equivalent reference measured; absolute costs are shown. Median of 3 samples; GOMAXPROCS=16.

Completion does not confirm server consumption of the full message. These costs are not normalized against HTTP request/response completion. Go allocations can include the server goroutines in the same process.

| Input / workload | Time/op ↓ | Allocated/op ↓ | Allocs/op ↓ |
| --- | ---: | ---: | ---: |
| 9 bytes | 6.53 µs | 0 B | 0 |
| 2 KiB | 6.87 µs | 0 B | 0 |
| 4 KiB | 7.04 µs | 0 B | 0 |
| 8 KiB | 9.32 µs | 0 B | 0 |
| 16 KiB | 8.38 µs | 0 B | 0 |
| 32 KiB | 16.7 µs | 0 B | 0 |
| 64 KiB | 33.41 µs | 0 B | 0 |
| 128 KiB | 66.65 µs | 0 B | 0 |
| 256 KiB | 133.88 µs | 0 B | 0 |
| 512 KiB | 268.72 µs | 0 B | 0 |
| 1 MiB | 538.92 µs | 0 B | 0 |
| 10 MiB | 5.46 ms | 0 B | 0 |
| 100 MiB | 53.29 ms | 0 B | 0 |
| 1 GiB | 549.04 ms | 0 B | 0 |

<details>
<summary>Absolute measurements, allocation counts, and observed ranges</summary>

#### 9 bytes

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| TCP | 6.5–6.54 µs | 0–0 B | 0 |

#### 2 KiB

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| TCP | 6.86–6.87 µs | 0–0 B | 0 |

#### 4 KiB

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| TCP | 7.04–7.18 µs | 0–0 B | 0 |

#### 8 KiB

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| TCP | 9.27–9.34 µs | 0–0 B | 0 |

#### 16 KiB

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| TCP | 8.28–8.51 µs | 0–0 B | 0 |

#### 32 KiB

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| TCP | 16.5–16.83 µs | 0–0 B | 0 |

#### 64 KiB

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| TCP | 33.29–33.41 µs | 0–0 B | 0 |

#### 128 KiB

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| TCP | 66.55–67.41 µs | 0–0 B | 0 |

#### 256 KiB

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| TCP | 133.49–134.27 µs | 0–0 B | 0 |

#### 512 KiB

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| TCP | 268.2–268.94 µs | 0–0 B | 0 |

#### 1 MiB

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| TCP | 536.25–540.62 µs | 0–0 B | 0 |

#### 10 MiB

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| TCP | 5.4–5.49 ms | 0–0 B | 0 |

#### 100 MiB

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| TCP | 53.15–54.32 ms | 0–0 B | 0 |

#### 1 GiB

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| TCP | 544.26–553.57 ms | 0–0 B | 0 |


</details>

**Recorded samples:** [protocol/http](../results/protocol-http.txt) · [protocol/tcp](../results/protocol-tcp.txt)

<!-- benchmark-results:end -->

## Methodology

| Protocol | Completion boundary |
| --- | --- |
| HTTP | POST the body, then wait for HTTP 200 after the handler reads it completely. |
| TCP | Write and flush bytes to the socket while the server drains the stream, without an application acknowledgement. |

A completed TCP operation does not prove the server has consumed the complete
message. Go allocation measurements can include server goroutines because the
servers run inside the benchmark process.

The smallest case sends the fixed string `something`; other cases use prepared
random printable-ASCII payloads. Input generation is excluded. Legacy benchmark
suffixes `KB`, `MB`, and `GB` mean KiB, MiB, and GiB.

Tests use dynamically assigned local ports. They do not measure physical-network
traffic, TLS, remote hosts, or production service processing.

## Running

Requires Go 1.27.4+. From the repository root:

```sh
go test -p=1 -run '^$' -bench . -benchmem -benchtime=1s -count=3 ./protocol/...
```
