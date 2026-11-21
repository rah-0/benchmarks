# Logging benchmarks

Compare JSON logging with **standard-library `slog` as the reference** for each
shared scenario. The peers are zerolog v1.35.1, zap v1.28.0, logrus v1.10.2, and
[slogx](https://github.com/rah-0/slogx) v0.3.1. Separate comparisons show the cost
of slogx context enrichment, spans, and structured errors.

[Results](#results) · [Methodology](#methodology) · [Run](#running)

<!-- benchmark-results:start -->
## Results

**2026-11-21 · go1.27.4 · linux/amd64 · AMD Ryzen 9 5900HX with Radeon Graphics**
GOMAXPROCS: 1; samples per case: 5.

Tables show medians. **Speed = baseline time ÷ candidate time**: 2.00× is twice
the baseline speed; 0.50× is half. Allocated bytes and allocation counts describe
Go heap allocation per operation, not retained memory or process RSS.
These samples are descriptive; observed min–max ranges appear in the details
and are not confidence intervals. A zero denominator is shown as **—**.

[Environment](../results/environment.txt) · [Versions](../results/modules.txt) · [Reading the comparisons](../README.md#reading-the-comparisons)

<details>
<summary>Jump to a workload</summary>

- [Simple INFO message](#simple-info-message)
- [INFO message with three fields](#info-message-with-three-fields)
- [Plain error](#plain-error)
- [Error with three fields](#error-with-three-fields)
- [Three correlated error records](#three-correlated-error-records)
- [slogx context enrichment](#slogx-context-enrichment)
- [slogx recording spans](#slogx-recording-spans)
- [slogx non-recording spans](#slogx-non-recording-spans)
- [slogx structured error lifecycle](#slogx-structured-error-lifecycle)

</details>

### Simple INFO message

Encode one INFO record with a timestamp and no attributes. All output goes to io.Discard.

Baseline: **slog = 1.00×**. Median of 5 samples; GOMAXPROCS=1.

Each logger keeps its native JSON field names and timestamp format. Source locations are disabled.

| Implementation | Time/op, µs ↓ | Speed vs baseline ↑ | Allocated/op, B ↓ | Allocs/op ↓ |
| --- | ---: | ---: | ---: | ---: |
| **slog (baseline)** | 2.81 | 1.00× | 0 | 0 |
| slogx | 2.81 | 1.00× | 0 | 0 |
| zerolog | 2.41 | 1.17× | 0 | 0 |
| zap | 2.55 | 1.11× | 0 | 0 |
| logrus | 4.14 | 0.68× | 752 | 18 |

zerolog has the lowest recorded median time: 1.17× the speed of slog; both allocate 0 B per operation.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| slog | 2.81–2.81 µs | 0–0 B | 0 | same (0 B) |
| slogx | 2.81–2.81 µs | 0–0 B | 0 | same (0 B) |
| zerolog | 2.41–2.42 µs | 0–0 B | 0 | same (0 B) |
| zap | 2.54–2.55 µs | 0–0 B | 0 | same (0 B) |
| logrus | 4.12–4.15 µs | 752–752 B | 18 | +752 B |

</details>

### INFO message with three fields

Encode one INFO record with a timestamp and three structured attributes. All output goes to io.Discard.

Baseline: **slog = 1.00×**. Median of 5 samples; GOMAXPROCS=1.

Each logger keeps its native JSON field names and timestamp format. Source locations are disabled.

| Implementation | Time/op, µs ↓ | Speed vs baseline ↑ | Allocated/op, B ↓ | Allocs/op ↓ |
| --- | ---: | ---: | ---: | ---: |
| **slog (baseline)** | 3.04 | 1.00× | 0 | 0 |
| slogx | 3.07 | 0.99× | 0 | 0 |
| zerolog | 2.41 | 1.26× | 0 | 0 |
| zap | 2.88 | 1.06× | 192 | 1 |
| logrus | 5.62 | 0.54× | 1,584 | 29 |

zerolog has the lowest recorded median time: 1.26× the speed of slog; both allocate 0 B per operation.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| slog | 3.03–3.05 µs | 0–0 B | 0 | same (0 B) |
| slogx | 3.07–3.07 µs | 0–0 B | 0 | same (0 B) |
| zerolog | 2.41–2.42 µs | 0–0 B | 0 | same (0 B) |
| zap | 2.88–2.89 µs | 192–192 B | 1 | +192 B |
| logrus | 5.6–5.63 µs | 1.55–1.55 KiB | 29 | +1.55 KiB |

</details>

### Plain error

Encode one ERROR record carrying an existing plain error. All output goes to io.Discard.

Baseline: **slog = 1.00×**. Median of 5 samples; GOMAXPROCS=1.

Each logger keeps its native JSON field names and timestamp format. Source locations are disabled.

| Implementation | Time/op, µs ↓ | Speed vs baseline ↑ | Allocated/op, B ↓ | Allocs/op ↓ |
| --- | ---: | ---: | ---: | ---: |
| **slog (baseline)** | 2.89 | 1.00× | 0 | 0 |
| slogx | 2.95 | 0.98× | 0 | 0 |
| zerolog | 2.41 | 1.20× | 0 | 0 |
| zap | 2.66 | 1.09× | 64 | 1 |
| logrus | 4.81 | 0.60× | 1,536 | 25 |

zerolog has the lowest recorded median time: 1.20× the speed of slog; both allocate 0 B per operation.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| slog | 2.88–2.91 µs | 0–0 B | 0 | same (0 B) |
| slogx | 2.95–2.96 µs | 0–0 B | 0 | same (0 B) |
| zerolog | 2.41–2.41 µs | 0–0 B | 0 | same (0 B) |
| zap | 2.65–2.66 µs | 64–64 B | 1 | +64 B |
| logrus | 4.8–4.82 µs | 1.5–1.5 KiB | 25 | +1.5 KiB |

</details>

### Error with three fields

Encode one ERROR record carrying an existing plain error and three attributes. All output goes to io.Discard.

Baseline: **slog = 1.00×**. Median of 5 samples; GOMAXPROCS=1.

Each logger keeps its native JSON field names and timestamp format. Source locations are disabled.

| Implementation | Time/op, µs ↓ | Speed vs baseline ↑ | Allocated/op, B ↓ | Allocs/op ↓ |
| --- | ---: | ---: | ---: | ---: |
| **slog (baseline)** | 3.22 | 1.00× | 0 | 0 |
| slogx | 3.22 | 1.00× | 0 | 0 |
| zerolog | 2.52 | 1.27× | 0 | 0 |
| zap | 2.97 | 1.08× | 256 | 1 |
| logrus | 6.23 | 0.52× | 2,080 | 35 |

zerolog has the lowest recorded median time: 1.27× the speed of slog; both allocate 0 B per operation.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| slog | 3.22–3.22 µs | 0–0 B | 0 | same (0 B) |
| slogx | 3.22–3.23 µs | 0–0 B | 0 | same (0 B) |
| zerolog | 2.52–2.53 µs | 0–0 B | 0 | same (0 B) |
| zap | 2.96–2.97 µs | 256–256 B | 1 | +256 B |
| logrus | 6.22–6.25 µs | 2.03–2.03 KiB | 35 | +2.03 KiB |

</details>

### Three correlated error records

Encode three ERROR records sharing a prepared trace ID; the first record carries an existing plain error. Identifier generation is excluded. All output goes to io.Discard.

Baseline: **slog = 1.00×**. Median of 5 samples; GOMAXPROCS=1.

Each logger keeps its native JSON field names and timestamp format. Source locations are disabled.

| Implementation | Time/op, µs ↓ | Speed vs baseline ↑ | Allocated/op, B ↓ | Allocs/op ↓ |
| --- | ---: | ---: | ---: | ---: |
| **slog (baseline)** | 9.24 | 1.00× | 0 | 0 |
| slogx | 9.18 | 1.01× | 0 | 0 |
| zerolog | 7.24 | 1.28× | 0 | 0 |
| zap | 8.11 | 1.14× | 256 | 3 |
| logrus | 15.02 | 0.62× | 5,056 | 78 |

zerolog has the lowest recorded median time: 1.28× the speed of slog; both allocate 0 B per operation.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| slog | 9.24–9.25 µs | 0–0 B | 0 | same (0 B) |
| slogx | 9.18–9.19 µs | 0–0 B | 0 | same (0 B) |
| zerolog | 7.24–7.25 µs | 0–0 B | 0 | same (0 B) |
| zap | 8.1–8.16 µs | 256–256 B | 3 | +256 B |
| logrus | 14.99–15.04 µs | 4.94–4.94 KiB | 78 | +4.94 KiB |

</details>

### slogx context enrichment

Encode one INFO record with a call-site attribute using a prepared context. Context construction is excluded.

Baseline: **No context fields = 1.00×**. Median of 5 samples; GOMAXPROCS=1.

Ratios describe the additional enrichment cost; these cases intentionally emit different fields.

| Implementation | Time/record, µs ↓ | Speed vs baseline ↑ | Allocated/record, B ↓ | Allocs/record ↓ |
| --- | ---: | ---: | ---: | ---: |
| **No context fields (baseline)** | 2.82 | 1.00× | 0 | 0 |
| 3 context fields | 3.22 | 0.88× | 0 | 0 |
| 10 context fields | 3.99 | 0.71× | 624 | 2 |
| Trace and span identity | 3.22 | 0.88× | 0 | 0 |

No context fields has the lowest recorded median time in this comparison.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| No context fields | 2.82–2.83 µs | 0–0 B | 0 | same (0 B) |
| 3 context fields | 3.22–3.22 µs | 0–0 B | 0 | same (0 B) |
| 10 context fields | 3.99–4 µs | 624–624 B | 2 | +624 B |
| Trace and span identity | 3.22–3.23 µs | 0–0 B | 0 | same (0 B) |

</details>

### slogx recording spans

Start and end a span with two attributes, including identifier generation, context allocation, and JSON emission.

No equivalent reference measured; absolute costs are shown. Median of 5 samples; GOMAXPROCS=1.

Root and child spans have different parent-context work. No equivalent alternative implementation was measured.

| Operation / variant | Time/span, µs ↓ | Allocated/span, KiB ↓ | Allocs/span ↓ |
| --- | ---: | ---: | ---: |
| Root span | 7.7 | 1.05 | 13 |
| Child of a prepared parent | 7.94 | 1.7 | 13 |

<details>
<summary>Observed sample ranges</summary>

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| Root span | 7.69–7.71 µs | 1.05–1.05 KiB | 13 |
| Child of a prepared parent | 7.92–7.94 µs | 1.7–1.7 KiB | 13 |

</details>

### slogx non-recording spans

Start and end a span without emitting its JSON record. Identifier generation, context allocation, and lifecycle work remain included.

No equivalent reference measured; absolute costs are shown. Median of 5 samples; GOMAXPROCS=1.

These are distinct suppression and parent-context paths, kept separate from recording spans.

| Operation / variant | Time/span, µs ↓ | Allocated/span, B ↓ | Allocs/span ↓ |
| --- | ---: | ---: | ---: |
| Root, level disabled | 5.05 | 496 | 6 |
| Child, level disabled | 4.96 | 464 | 5 |
| Child, parent unsampled | 4.96 | 464 | 5 |

<details>
<summary>Observed sample ranges</summary>

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| Root, level disabled | 5.04–5.05 µs | 496–496 B | 6 |
| Child, level disabled | 4.96–4.97 µs | 464–464 B | 5 |
| Child, parent unsampled | 4.95–4.96 µs | 464–464 B | 5 |

</details>

### slogx structured error lifecycle

Construct three attributed error wrappers, or encode a prepared three-layer error in one ERROR record.

No equivalent reference measured; absolute costs are shown. Median of 5 samples; GOMAXPROCS=1.

Construction excludes logging; logging excludes construction. Neither row is a baseline for the other.

| Operation / variant | Time/op, ns ↓ | Allocated/op, B ↓ | Allocs/op ↓ |
| --- | ---: | ---: | ---: |
| Construct three wrappers | 209.2 | 336 | 6 |
| Log a prepared wrapped error | 3,947 | 384 | 3 |

<details>
<summary>Observed sample ranges</summary>

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| Construct three wrappers | 209.1–211.8 ns | 336–336 B | 6 |
| Log a prepared wrapped error | 3.93–3.96 µs | 384–384 B | 3 |

</details>

**Recorded samples:** [log](../results/log.txt)

<!-- benchmark-results:end -->

## Methodology

All loggers encode JSON to `io.Discard`. Timings include encoding and exclude
terminal, file, pipe, network, and collector work. Loggers are prepared before
timing; all cases use `testing.B.Loop`.

Timestamps are enabled and source locations are disabled. Native slog and slogx
share JSON handler settings, the INFO threshold, arguments, and timestamp format.
Slogx uses `slogx.New`, including its context handler; it excludes the custom
timestamp default of `slogx.NewDefault`. Other libraries retain their own field
names and timestamp formats, so output bytes are not identical across all peers.

### Shared scenarios

| Scenario | One operation |
| --- | --- |
| Simple message | One INFO record without attributes. |
| Structured fields | One INFO record with three attributes. |
| Error | One ERROR record carrying an existing plain error. |
| Error with fields | One ERROR record with the plain error and three attributes. |
| Correlated error chain | Three ERROR records sharing a fixed 32-character trace ID; the first carries the plain error. ID generation is excluded. |

Shared error scenarios use native logger methods. The three-record operation
must be compared with the same operation in the other libraries. Slogx's nested
error representation is measured separately because it performs different work.

### Slogx features

| Feature | One operation |
| --- | --- |
| Context enrichment | Log one INFO record with a call-site attribute and a prepared context containing zero, three, or ten attributes, or a valid trace/span identity. Context construction is excluded. |
| Spans | Start and end one span with two attributes. Recording cases emit JSON; non-recording cases still create IDs and contexts. Child cases reuse a prepared remote parent. |
| Build structured error | Construct three attributed wrappers around an existing error, excluding logging. |
| Log structured error | Encode a prepared three-layer structured error, excluding wrapper construction. |

These are serial workloads. They do not measure concurrent handler contention,
disabled ordinary log calls, context attachment, custom timestamps, source
locations, text decoration, or JSON decoding. Reported time is an average per
operation, not individual-record tail latency.

## Running

Requires Go 1.27.4+. The [full runner](../scripts/benchmark.bash) uses one CPU,
five samples, `GOGC=100`, and `GOMEMLIMIT=off`, with no competing benchmark suite.
From the repository root, choose an available CPU and reproduce those settings:

```sh
GOMAXPROCS=1 GOGC=100 GOMEMLIMIT=off taskset -c 0 \
  go test -p=1 -run '^$' -bench . -benchmem \
  -benchtime=1s -count=5 -cpu=1 ./log
```

Use the scenario filters for focused runs:

```sh
# Shared scenarios.
go test -run '^$' -bench '^Benchmark(Zerolog|Zap|Logrus|Slog|Slogx)_' -benchmem ./log

# Slogx feature costs.
go test -run '^$' -bench '^BenchmarkSlogxFeatures_' -benchmem ./log
```

A race-instrumented execution check verifies behavior rather than performance:

```sh
go test -count=1 -race -cover -covermode=atomic -run '^$' -bench . -benchtime=1x ./log
```

The [correlation example](correlation_demo.go) prints the shared three-record
scenario for all five libraries:

```sh
go run ./log/correlation_demo.go
```
