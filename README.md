# Benchmarks

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Compare Go libraries and application workloads by execution time, throughput,
and allocation cost. Each category groups matching work around a named reference
and shows the tradeoffs alongside the relative performance.

Requires **Go 1.27.4 or later**. Service and hardware requirements vary by suite.

## Explore the results

| Category | Workloads | Reference |
| --- | --- | --- |
| [Cache](cache/README.md#results) | Redis and Memcached set/get by payload size | Memcached |
| [Collectors](collectors/README.md#results) | Log delivery and recovery after an outage | Fluent Bit |
| [Compression](compression/README.MD#results) | Brotli, Gzip, and LZ4 round trips | Gzip |
| [Cryptography](crypto/README.md#results) | Create and read signed or encrypted tokens | Separate operations |
| [Databases](db/README.md#results) | Inserts, indexed lookup, sorting, and filtering | PostgreSQL; embedded SQLite separately |
| [Encoding](encoding/README.md#results) | Base64 and ASCII85 across input sizes | Base64 |
| [Iteration](iteration/README.md#results) | Full scan, early close, first row, and parallel callers | Direct slice or slice adapter |
| [Logging](log/README.md#results) | JSON messages, fields, errors, and slogx features | Native `slog` for shared scenarios |
| [Meta](meta/README.MD#results) | Reflection, sorting, and synchronization | Direct access, `sort.Ints`, or matching primitive |
| [Protocols](protocol/README.md#results) | Local HTTP and TCP payload scaling | Separate completion boundaries |
| [Serialization](serializer/README.md#results) | JSON, Gob, and Protobuf by record shape/count | JSON; warmed streams separately |
| [Validation](validator/README.md#results) | Valid/invalid fields and forms, error processing | Playground for matching scenarios |

## Reading the comparisons

The reference is fixed for each workload and appears as **1.00×**. It is not
chosen from whichever implementation happened to be fastest in a run.

| Relative metric | Calculation using medians | Direction |
| --- | --- | --- |
| Speed vs reference | Reference time ÷ implementation time | Higher is better: 2.00× takes half the time; 0.50× takes twice the time. |
| Throughput vs reference | Implementation throughput ÷ reference throughput | Higher is better. |
| Allocated bytes or RSS vs reference | Implementation value ÷ reference value | Lower is better. |

Ratios compare the same workload, input size, and concurrency. Different
completion, reuse, or correctness guarantees remain in separate groups; a
missing equivalent reference is not replaced with an artificial ratio. Zero
allocation references use absolute differences; other undefined ratios appear
as **—** rather than dividing by zero.
Tiny positive ratios appear as **<0.01×** when two decimals would round them to
zero. Units stay consistent within a table unless its values span several
orders of magnitude; those cells show their units explicitly.

**Allocated bytes/op** count cumulative Go heap allocations. They do not measure
retained memory, peak memory, or memory in external servers. Collector **peak
RSS** measures the collector process instead; its boundaries are documented in
that suite. Parallel time/op describes aggregate throughput, not individual-call
latency.

Size sweeps show scaling against the same reference at each size. Expand the
details for allocation costs and observed ranges, or follow the relative links
to the recorded samples, module versions, and environment. The current runs have
three samples per case and five for logging. Their medians and ranges describe
those observations; close values do not establish a significant performance
difference.

Before/after version comparisons require matching hardware, workload, toolchain,
and configuration, except for the variable intentionally changed. A ratio
between libraries in the current run is not evidence of a version improvement.

## Running

Run the complete matrix on a Linux/amd64 machine with Docker, Bash, a C compiler,
`taskset`, and enough memory for the 1 GiB payload cases:

```sh
evo go deps update
go test -p=1 -count=1 ./...
go -C collectors test -count=1 ./...
bash scripts/benchmark.bash
bash scripts/results.bash render
```

The benchmark runner executes suites sequentially and records benchmark samples,
toolchain, hardware, and module versions in `results/`. Database and cache
suites start disposable Testcontainers services on dynamically assigned ports.
The collector suite uses persistent storage under `/var/tmp`. See its
[methodology](collectors/README.md) for durability and measurement boundaries.

Published Go samples retain measurement headers, benchmark rows, and completion
lines. Docker connection details, transient identifiers, and other diagnostics
are omitted. Failed Go benchmark runs retain captured diagnostics in a private
temporary directory outside the checkout. Local filesystem module replacements
prevent publication because their paths and contents cannot identify a
reproducible public dependency.

Check saved samples and generated documentation without changing files:

```sh
bash scripts/results.bash sanitize --check results
bash scripts/results.bash render --check
sha256sum --check results/source.sha256
```

`BENCHTIME` and `BENCHCOUNT` default to `1s` and `3`; logging uses five samples.
`LOG_CPU`, `HARNESS_CPU`, and `COLLECTOR_CPU` default to `0`, `2`, and `0`.
Choose available CPUs, with the collector and harness on separate physical cores.
Measurements exclude race and coverage instrumentation.

The recorded [source checksums](results/source.sha256) identify the benchmark Go source,
module files, and collector configuration used for these results. Workload labels
and references are explicit in the [comparison definitions](scripts/benchmarks/groups.go);
new cases need a defined comparison before their documentation can be generated.
The [script guide](scripts/README.md) describes the Go reporting tools and Bash runner.

## Support

If these benchmarks are useful to you, consider supporting the project.

[![Buy Me A Coffee](https://cdn.buymeacoffee.com/buttons/default-orange.png)](https://www.buymeacoffee.com/rah.0)
