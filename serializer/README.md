# Serialization benchmarks

Compare encoding and decoding of small `Person` records and 100-field `Unreal`
records. **JSON is the reference for matching record shapes and counts.**
Warmed Gob streams are shown separately because they amortize type metadata
across operations.

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

**Workloads:** [Small records: fresh round trips](#small-records-fresh-round-trips) · [Small records: warmed Gob stream](#small-records-warmed-gob-stream) · [Large records: fresh round trips](#large-records-fresh-round-trips) · [Large records: warmed Gob stream](#large-records-warmed-gob-stream)

### Small records: fresh round trips

One operation serializes and deserializes the listed number of generated Person records. Input generation and Protobuf type conversion are excluded.

Baseline: **JSON = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

JSON and Gob process a slice; Protobuf processes each message separately. The ratios compare these implemented workload choices, including their different framing.

| Input / workload | JSON time (1.00×) | Gob, fresh stream speed ↑ | Protobuf speed ↑ |
| --- | ---: | ---: | ---: |
| 1 record | 1.06 µs | 0.09× | 4.00× |
| 100 records | 54.53 µs | 1.68× | 2.06× |
| 10,000 records | 5.4 ms | 2.91× | 2.04× |
| 1,000,000 records | 636.28 ms | 2.63× | 2.30× |

Allocation cost accompanies the speed comparison:

| Input / workload | JSON allocated | Gob, fresh stream allocated vs baseline ↓ | Protobuf allocated vs baseline ↓ |
| --- | ---: | ---: | ---: |
| 1 record | 232 B | 41.21× | 1.10× |
| 100 records | 27.65 KiB | 2.59× | 0.90× |
| 10,000 records | 3.76 MiB | 1.80× | 0.65× |
| 1,000,000 records | 984.84 MiB | 1.16× | 0.25× |

The lowest median shifts from Protobuf at 1 record to Gob, fresh stream at 1,000,000 records.

At 1,000,000 records, Gob, fresh stream has the lowest recorded median time: 2.63× the speed of JSON, with 1.16× its allocated bytes per operation.

<details>
<summary>Absolute measurements, allocation counts, and observed ranges</summary>

#### 1 record

| Implementation | Time/round trip, ns ↓ | Speed vs baseline ↑ | Allocated/round trip, B ↓ | Allocs/round trip ↓ |
| --- | ---: | ---: | ---: | ---: |
| **JSON (baseline)** | 1,056 | 1.00× | 232 | 5 |
| Gob, fresh stream | 11,911 | 0.09× | 9,560 | 200 |
| Protobuf | 263.8 | 4.00× | 256 | 4 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| JSON | 1.05–1.06 µs | 232–232 B | 5 | 1.00× |
| Gob, fresh stream | 11.85–11.91 µs | 9.34–9.34 KiB | 200 | 41.21× |
| Protobuf | 263.6–264.8 ns | 256–256 B | 4 | 1.10× |

#### 100 records

| Implementation | Time/round trip, µs ↓ | Speed vs baseline ↑ | Allocated/round trip, KiB ↓ | Allocs/round trip ↓ |
| --- | ---: | ---: | ---: | ---: |
| **JSON (baseline)** | 54.53 | 1.00× | 27.65 | 125 |
| Gob, fresh stream | 32.47 | 1.68× | 71.57 | 408 |
| Protobuf | 26.5 | 2.06× | 25 | 400 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| JSON | 54.47–54.73 µs | 26.95–27.65 KiB | 110–125 | 1.00× |
| Gob, fresh stream | 32.19–33.29 µs | 71.57–71.57 KiB | 408 | 2.59× |
| Protobuf | 26.46–26.5 µs | 25–25 KiB | 400 | 0.90× |

#### 10,000 records

| Implementation | Time/round trip, ms ↓ | Speed vs baseline ↑ | Allocated/round trip, MiB ↓ | Allocs/round trip ↓ |
| --- | ---: | ---: | ---: | ---: |
| **JSON (baseline)** | 5.4 | 1.00× | 3.76 | 20,023 |
| Gob, fresh stream | 1.85 | 2.91× | 6.76 | 20,225 |
| Protobuf | 2.65 | 2.04× | 2.44 | 40,000 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| JSON | 5.39–5.41 ms | 3.76–3.76 MiB | 20,023 | 1.00× |
| Gob, fresh stream | 1.85–1.88 ms | 6.76–6.76 MiB | 20,225 | 1.80× |
| Protobuf | 2.64–2.65 ms | 2.44–2.44 MiB | 40,000 | 0.65× |

#### 1,000,000 records

| Implementation | Time/round trip, ms ↓ | Speed vs baseline ↑ | Allocated/round trip, MiB ↓ | Allocs/round trip ↓ |
| --- | ---: | ---: | ---: | ---: |
| **JSON (baseline)** | 636.28 | 1.00× | 984.84 | 2,000,111 |
| Gob, fresh stream | 242.38 | 2.63× | 1,137.61 | 2,000,259 |
| Protobuf | 276.44 | 2.30× | 244.14 | 4,000,000 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| JSON | 613.29–637.34 ms | 691.47–984.84 MiB | 2,000,079–2,000,111 | 1.00× |
| Gob, fresh stream | 235.6–265.25 ms | 1.11–1.11 GiB | 2,000,259 | 1.16× |
| Protobuf | 274.91–277.84 ms | 244.14–244.14 MiB | 4,000,000 | 0.25× |


</details>

### Small records: warmed Gob stream

One operation encodes and decodes the listed number of Person records using a reused buffer, encoder, and decoder after stream metadata has been warmed.

No equivalent reference measured; absolute costs are shown. Median of 3 samples; GOMAXPROCS=16.

Stream metadata is amortized across operations. These bytes are not independently decodable standalone messages, so this panel is separate from fresh round trips.

| Input / workload | Time/op ↓ | Allocated/op ↓ | Allocs/op ↓ |
| --- | ---: | ---: | ---: |
| 1 record | 649 ns | 376 B | 8 |
| 100 records | 17.47 µs | 21.51 KiB | 206 |
| 10,000 records | 1.63 ms | 2.06 MiB | 20,006 |
| 1,000,000 records | 205.68 ms | 592.69 MiB | 2,000,019 |

<details>
<summary>Absolute measurements, allocation counts, and observed ranges</summary>

#### 1 record

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| Gob, warmed stream | 647.8–654 ns | 376–376 B | 8 |

#### 100 records

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| Gob, warmed stream | 17.44–17.7 µs | 21.51–21.51 KiB | 206 |

#### 10,000 records

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| Gob, warmed stream | 1.61–1.63 ms | 2.06–2.06 MiB | 20,006 |

#### 1,000,000 records

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| Gob, warmed stream | 194.2–207.42 ms | 592.69–592.69 MiB | 2,000,019 |


</details>

### Large records: fresh round trips

One operation serializes and deserializes the listed number of generated 100-field Unreal records. Input generation and Protobuf type conversion are excluded.

Baseline: **JSON = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

JSON and Gob process a slice; Protobuf processes each message separately. The ratios compare these implemented workload choices, including their different framing.

| Input / workload | JSON time (1.00×) | Gob, fresh stream speed ↑ | Protobuf speed ↑ |
| --- | ---: | ---: | ---: |
| 1 record | 18.38 µs | 0.26× | 3.44× |
| 10 records | 187.63 µs | 1.44× | 3.49× |
| 100 records | 1.83 ms | 2.57× | 3.42× |
| 1,000 records | 18.85 ms | 2.92× | 3.54× |

Allocation cost accompanies the speed comparison:

| Input / workload | JSON allocated | Gob, fresh stream allocated vs baseline ↓ | Protobuf allocated vs baseline ↓ |
| --- | ---: | ---: | ---: |
| 1 record | 7.98 KiB | 6.68× | 1.31× |
| 10 records | 145.56 KiB | 2.19× | 0.72× |
| 100 records | 1.36 MiB | 2.39× | 0.75× |
| 1,000 records | 16.5 MiB | 2.03× | 0.62× |

At 1,000 records, Protobuf has the lowest recorded median time: 3.54× the speed of JSON, with 0.62× its allocated bytes per operation.

<details>
<summary>Absolute measurements, allocation counts, and observed ranges</summary>

#### 1 record

| Implementation | Time/round trip, µs ↓ | Speed vs baseline ↑ | Allocated/round trip, KiB ↓ | Allocs/round trip ↓ |
| --- | ---: | ---: | ---: | ---: |
| **JSON (baseline)** | 18.38 | 1.00× | 7.98 | 35 |
| Gob, fresh stream | 71.04 | 0.26× | 53.26 | 694 |
| Protobuf | 5.34 | 3.44× | 10.44 | 102 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| JSON | 18.29–18.55 µs | 7.79–8.07 KiB | 31–37 | 1.00× |
| Gob, fresh stream | 70.78–71.24 µs | 53.26–53.26 KiB | 694 | 6.68× |
| Protobuf | 5.34–5.35 µs | 10.44–10.44 KiB | 102 | 1.31× |

#### 10 records

| Implementation | Time/round trip, µs ↓ | Speed vs baseline ↑ | Allocated/round trip, KiB ↓ | Allocs/round trip ↓ |
| --- | ---: | ---: | ---: | ---: |
| **JSON (baseline)** | 187.63 | 1.00× | 145.56 | 992 |
| Gob, fresh stream | 130.69 | 1.44× | 318.95 | 1,602 |
| Protobuf | 53.7 | 3.49× | 104.38 | 1,020 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| JSON | 187.57–187.87 µs | 145.56–145.63 KiB | 989–992 | 1.00× |
| Gob, fresh stream | 130.44–130.83 µs | 318.95–318.95 KiB | 1,602 | 2.19× |
| Protobuf | 53.52–53.77 µs | 104.38–104.38 KiB | 1,020 | 0.72× |

#### 100 records

| Implementation | Time/round trip, µs ↓ | Speed vs baseline ↑ | Allocated/round trip, MiB ↓ | Allocs/round trip ↓ |
| --- | ---: | ---: | ---: | ---: |
| **JSON (baseline)** | 1,828.68 | 1.00× | 1.36 | 10,012 |
| Gob, fresh stream | 712.07 | 2.57× | 3.26 | 10,611 |
| Protobuf | 535.29 | 3.42× | 1.02 | 10,200 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| JSON | 1.83–1.83 ms | 1.36–1.37 MiB | 10,012 | 1.00× |
| Gob, fresh stream | 711.12–713.67 µs | 3.26–3.26 MiB | 10,611 | 2.39× |
| Protobuf | 533.8–535.43 µs | 1.02–1.02 MiB | 10,200 | 0.75× |

#### 1,000 records

| Implementation | Time/round trip, ms ↓ | Speed vs baseline ↑ | Allocated/round trip, MiB ↓ | Allocs/round trip ↓ |
| --- | ---: | ---: | ---: | ---: |
| **JSON (baseline)** | 18.85 | 1.00× | 16.5 | 100,020 |
| Gob, fresh stream | 6.47 | 2.92× | 33.52 | 100,621 |
| Protobuf | 5.33 | 3.54× | 10.19 | 102,000 |

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| JSON | 18.82–18.89 ms | 16.28–16.69 MiB | 100,019–100,020 | 1.00× |
| Gob, fresh stream | 6.45–6.48 ms | 33.52–33.52 MiB | 100,621 | 2.03× |
| Protobuf | 5.32–5.33 ms | 10.19–10.19 MiB | 102,000 | 0.62× |


</details>

### Large records: warmed Gob stream

One operation encodes and decodes the listed number of 100-field Unreal records using a reused buffer, encoder, and decoder after stream metadata has been warmed.

No equivalent reference measured; absolute costs are shown. Median of 3 samples; GOMAXPROCS=16.

Stream metadata is amortized across operations. These bytes are not independently decodable standalone messages, so this panel is separate from fresh round trips.

| Input / workload | Time/op ↓ | Allocated/op ↓ | Allocs/op ↓ |
| --- | ---: | ---: | ---: |
| 1 record | 5.27 µs | 10.57 KiB | 106 |
| 10 records | 48.21 µs | 103.01 KiB | 1,006 |
| 100 records | 466.39 µs | 1,004.88 KiB | 10,006 |
| 1,000 records | 4.68 ms | 9.73 MiB | 100,006 |

<details>
<summary>Absolute measurements, allocation counts, and observed ranges</summary>

#### 1 record

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| Gob, warmed stream | 5.23–5.28 µs | 10.57–10.57 KiB | 106 |

#### 10 records

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| Gob, warmed stream | 48–48.44 µs | 103.01–103.01 KiB | 1,006 |

#### 100 records

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| Gob, warmed stream | 465.2–466.53 µs | 1,004.88–1,004.88 KiB | 10,006 |

#### 1,000 records

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| Gob, warmed stream | 4.68–4.68 ms | 9.73–9.73 MiB | 100,006 |


</details>

**Recorded samples:** [serializer](../results/serializer.txt)

<!-- benchmark-results:end -->

## Methodology

Each operation includes both serialization and deserialization. Input generation
and conversion into Protobuf message types occur before timing. Sizes in the
result groups count records, not bytes.

| Format | Encoding and reuse model |
| --- | --- |
| JSON | Standard-library JSON encodes and decodes a slice. |
| Gob | A fresh buffer, encoder, and decoder encode and decode the slice. |
| Protobuf | Each message is encoded and decoded separately within the operation. |
| Warmed Gob stream | Reuses a warmed encoder, decoder, and buffer; stream type metadata is amortized. |

Framing differs between formats: Protobuf processes messages individually,
whereas JSON and Gob process slices. Ratios compare these implemented workloads,
not identical wire formats. Warmed Gob results do not represent independent,
standalone messages.

Large record-count cases can allocate substantially more bytes than the retained
input size. Encoded message size is not measured here.

## Running

Requires Go 1.27.4+. From the repository root:

```sh
go test -p=1 -run '^$' -bench . -benchmem -benchtime=1s -count=3 ./serializer
```
