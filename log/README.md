# Go logging benchmarks

Requires Go 1.27.1 or later.

The shared scenarios compare JSON logging with zerolog v1.34.0, zap v1.27.0,
logrus v1.9.3, native `log/slog`, and
[slogx](https://github.com/rah-0/slogx) v0.3.1.
Additional scenarios measure slogx context enrichment, spans, and structured errors.

## Running

From the repository root:

```sh
go test -run '^$' -bench . -benchmem -benchtime=1s -count=5 ./log

# Only the shared scenarios, including native slog and slogx:
go test -run '^$' -bench '^Benchmark(Zerolog|Zap|Logrus|Slog|Slogx)_' -benchmem ./log

# Only slogx feature costs:
go test -run '^$' -bench '^BenchmarkSlogxFeatures_' -benchmem ./log
```

For an execution check with race detection, rather than performance measurements:

```sh
go test -count=1 -race -cover -covermode=atomic -run '^$' -bench . -benchtime=1x ./log
```

The [correlation example](correlation_demo.go) prints the shared three-record
scenario for all five libraries:

```sh
go run ./log/correlation_demo.go
```

## What is measured

All loggers encode JSON to `io.Discard`; results include encoding but exclude
terminal, file, pipe, network, and collector work. Loggers are created before
timing, and every benchmark uses `testing.B.Loop`.

The shared scenarios use each library's normal structured logging API, with
timestamps enabled and source locations disabled. Native slog and slogx use the
same JSON handler settings, INFO threshold, arguments, and native timestamp
format. The slogx cases use `slogx.New`, so they include its context handler but
exclude the custom timestamp default of `slogx.NewDefault`.
Other libraries retain their own JSON field names and timestamp formats; output
is not byte-for-byte identical across all libraries.

| Shared scenario | One operation |
| --- | --- |
| `SimpleMessage` | One INFO record without attributes. |
| `WithFields` | One INFO record with three attributes. |
| `Error` | One ERROR record carrying an existing plain error. |
| `ErrorWithFields` | One ERROR record with an existing plain error and three attributes. |
| `ErrorChain` | Three ERROR records sharing a fixed 32-character trace ID; the first carries the plain error. Identifier generation is excluded. |

The shared error scenarios use native logger methods. Slogx's nested error
representation is measured separately because it performs different work.

| Slogx feature | One operation |
| --- | --- |
| `Context` | One INFO record with a call-site attribute, using a prepared context containing zero, three, or ten attributes, or a valid trace/span identity. Context construction is excluded. |
| `Span` | Start and end one span with two attributes. Recording cases include JSON emission; non-recording cases still generate IDs and contexts. Child cases reuse a prepared remote parent. |
| `StructuredError/WrapThreeLayers` | Construct three attributed wrappers around an existing error; excludes logging. |
| `StructuredError/LogThreeLayers` | Encode one ERROR record containing a prepared three-layer structured error; excludes wrapper construction. |

These are serial benchmarks. They do not measure concurrent handler contention,
disabled ordinary log calls, context attachment, custom timestamps, source
locations, text decoration, or JSON decoding.

## Measurements

Measured on 2026-09-10 with Go 1.27.1, linux/amd64, on an AMD Ryzen AI 9 HX 370,
Debian 13, Linux 6.12.105+deb13-amd64. This is a development host.
The run pinned the process to CPU 0, used `GOMAXPROCS=1`, `GOGC=100`,
`GOMEMLIMIT=off`, and the toolchain's default experiments.
No other test or benchmark runs were started concurrently.

```sh
GOGC=100 GOMEMLIMIT=off taskset -c 0 \
  go test -p=1 -run '^$' -bench . -benchmem \
  -benchtime=1s -count=5 -cpu=1 ./log > log/benchmark.txt
```

Choose an available CPU for `taskset`, or omit it on platforms without that
command. Performance measurements exclude race and coverage instrumentation.
Tables report the median of five samples from the
[complete benchmark output](benchmark.txt). Times are averages per operation,
not individual-record tail latency. Small timing differences on this host
should not be treated as universal rankings.


### Simple message

| Library | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: |
| zerolog | 90.74 | 0 | 0 |
| zap | 172.1 | 0 | 0 |
| slog | 290.1 | 0 | 0 |
| slogx | 304.2 | 0 | 0 |
| logrus | 1258 | 784 | 21 |

### Message with three fields

| Library | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: |
| zerolog | 120.8 | 0 | 0 |
| zap | 367.5 | 192 | 1 |
| slog | 484.9 | 0 | 0 |
| slogx | 495.6 | 0 | 0 |
| logrus | 2388 | 1608 | 32 |

### Plain error

| Library | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: |
| zerolog | 111.6 | 0 | 0 |
| zap | 295.3 | 64 | 1 |
| slog | 386 | 0 | 0 |
| slogx | 367.4 | 0 | 0 |
| logrus | 1969 | 1563 | 28 |

### Plain error with three fields

| Library | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: |
| zerolog | 137.4 | 0 | 0 |
| zap | 419 | 256 | 1 |
| slog | 513.7 | 0 | 0 |
| slogx | 498.6 | 0 | 0 |
| logrus | 2851 | 2128 | 38 |

### Three correlated records

| Library | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: |
| zerolog | 330.5 | 0 | 0 |
| zap | 888.1 | 256 | 3 |
| slog | 1494 | 0 | 0 |
| slogx | 1433 | 0 | 0 |
| logrus | 5924 | 5137 | 87 |

### Slogx context enrichment

| Case | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: |
| `NoAttrs` | 337.5 | 0 | 0 |
| `Attrs3` | 484.3 | 0 | 0 |
| `Attrs10` | 1067 | 624 | 2 |
| `TraceOnly` | 527.3 | 0 | 0 |

### Slogx span lifecycle

| Case | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: |
| `RootRecording` | 2191 | 1072 | 13 |
| `ChildRecording` | 2225 | 1744 | 13 |
| `RootLevelDisabled` | 361 | 496 | 6 |
| `ChildLevelDisabled` | 309.7 | 464 | 5 |
| `ChildUnsampled` | 312.7 | 464 | 5 |

### Slogx structured errors

| Case | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: |
| `WrapThreeLayers` | 149.1 | 336 | 6 |
| `LogThreeLayers` | 1053 | 384 | 3 |

In this run, slogx and native slog have the same allocation counts in all five
shared scenarios. Slogx's context, span, and structured-error cases report their
additional work separately; their timings are not direct cross-library comparisons.
