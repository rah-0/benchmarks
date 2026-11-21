# Log collector benchmarks

Compare slogx-collector, Vector, and Fluent Bit delivering fixed NDJSON bursts
through stdin to a checked local HTTP receiver. Each pipeline uses filesystem
buffering; the results retain its native persistence and batching costs.

[Results](#results) · [Methodology](#methodology) ·
[Configuration](#collectors-and-configuration) · [Run the benchmarks](#running)

<!-- benchmark-results:start -->
## Results

**Baseline: Fluent Bit 5.2.0.** Compared with slogx-collector v0.7.0 and Vector 0.62.0-debian.
Each sample delivers **2,000 records**; cells are medians of **3 samples**. Record sizes include the newline.

**Persistence differs:** slogx delivers synced journal records; Vector can deliver
before disk synchronization; Fluent Bit synchronizes filesystem chunks.
The ratios compare these configured pipelines, including their distinct batching costs.

Higher throughput and recovery-speed ratios are better; lower RSS ratios are better.
Ratios use unrounded medians. Observed ranges appear in the details; these samples
describe measured bursts and do not establish statistical significance.
See [reading the comparisons](../README.md#reading-the-comparisons) for the shared conventions.

### Healthy delivery

Time runs from the first producer write to the final validated HTTP acknowledgement.
RSS is collector-process lifetime peak memory, including startup and warmup.

#### 256-byte records

| Collector | Records/s ↑ | Throughput vs baseline ↑ | Peak RSS (MiB) ↓ | RSS vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Fluent Bit — baseline** | 18,522 | 1.00× | 24.16 | 1.00× |
| slogx-collector | 39,893 | 2.15× | 12.33 | 0.51× |
| Vector | 60,407 | 3.26× | 101.14 | 4.19× |

Vector has the highest measured median delivery rate; slogx-collector has the lowest measured median peak RSS.

<details>
<summary>Observed ranges and supporting measurements</summary>

| Collector | Records/s range | Peak RSS range (MiB) |
| --- | ---: | ---: |
| Fluent Bit | 18,483–18,637 | 24.09–25.92 |
| slogx-collector | 38,052–39,976 | 10.34–12.34 |
| Vector | 59,134–60,422 | 101.03–103.12 |

| Collector | Input MiB/s | CPU/burst (ms) | Producer (ms) | HTTP requests |
| --- | ---: | ---: | ---: | ---: |
| Fluent Bit | 4.52 | 10.00 | 3.69 | 1 |
| slogx-collector | 9.74 | 20.00 | 11.13 | 4 |
| Vector | 14.75 | 30.00 | 1.55 | 4 |

CPU counters have 10 ms resolution on this host; zero is below one observable
counter increment. Producer completion does not establish ingestion or durability.

</details>

#### 1,024-byte records

| Collector | Records/s ↑ | Throughput vs baseline ↑ | Peak RSS (MiB) ↓ | RSS vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Fluent Bit — baseline** | 17,238 | 1.00× | 27.11 | 1.00× |
| slogx-collector | 33,822 | 1.96× | 13.00 | 0.48× |
| Vector | 49,064 | 2.85× | 107.03 | 3.95× |

Vector has the highest measured median delivery rate; slogx-collector has the lowest measured median peak RSS.

<details>
<summary>Observed ranges and supporting measurements</summary>

| Collector | Records/s range | Peak RSS range (MiB) |
| --- | ---: | ---: |
| Fluent Bit | 16,708–17,356 | 27.08–27.83 |
| slogx-collector | 27,610–35,501 | 12.88–13.02 |
| Vector | 47,214–49,165 | 106.71–111.40 |

| Collector | Input MiB/s | CPU/burst (ms) | Producer (ms) | HTTP requests |
| --- | ---: | ---: | ---: | ---: |
| Fluent Bit | 16.83 | 20.00 | 10.52 | 1 |
| slogx-collector | 33.03 | 20.00 | 36.72 | 4 |
| Vector | 47.91 | 30.00 | 6.23 | 4 |

CPU counters have 10 ms resolution on this host; zero is below one observable
counter increment. Producer completion does not establish ingestion or durability.

</details>

### Recovery after a 2-second receiver outage

Recovery runs from receiver availability to the final unique-record acknowledgement.
**Recovery speed = Fluent Bit recovery time / collector recovery time.**
It includes retry scheduling and drain time; it does not measure crash recovery.
A dash means a zero-duration result cannot provide a finite speed ratio.

#### 256-byte records

| Collector | Recovery (s) ↓ | Recovery speed vs baseline ↑ |
| --- | ---: | ---: |
| **Fluent Bit — baseline** | 1.303 | 1.00× |
| slogx-collector | 0.049 | 26.79× |
| Vector | 0.042 | 31.09× |

Vector has the lowest measured median recovery time.

<details>
<summary>Observed recovery ranges and retry activity</summary>

| Collector | Recovery range (s) | Producer (ms) | HTTP 503 requests |
| --- | ---: | ---: | ---: |
| Fluent Bit | 1.200–2.001 | 3.80 | 2 |
| slogx-collector | 0.049–0.049 | 11.04 | 2 |
| Vector | 0.039–0.043 | 0.59 | 2 |

</details>

#### 1,024-byte records

| Collector | Recovery (s) ↓ | Recovery speed vs baseline ↑ |
| --- | ---: | ---: |
| **Fluent Bit — baseline** | 1.310 | 1.00× |
| slogx-collector | 0.063 | 20.92× |
| Vector | 0.052 | 25.05× |

Vector has the lowest measured median recovery time. Its observed range overlaps slogx-collector's; these samples do not establish a consistent ordering between them.

<details>
<summary>Observed recovery ranges and retry activity</summary>

| Collector | Recovery range (s) | Producer (ms) | HTTP 503 requests |
| --- | ---: | ---: | ---: |
| Fluent Bit | 1.208–2.010 | 10.10 | 2 |
| slogx-collector | 0.042–0.063 | 35.06 | 2 |
| Vector | 0.051–0.052 | 5.21 | 2 |

</details>

### Direct receiver control

The control bypasses collector processes, stdin attachment, and disk buffering.
It checks receiver capacity and is excluded from collector speed ratios.

| Record size | Records/s ↑ | Input MiB/s ↑ | Outage recovery (s) ↓ |
| --- | ---: | ---: | ---: |
| 256 B | 340,225 | 83.06 | 0.010 |
| 1,024 B | 157,686 | 153.99 | 0.019 |

**Delivery checks:** 24 healthy and 24 outage samples completed; 0 duplicate records. Every outage sample encountered HTTP 503 before successful delivery.

<details>
<summary>Measurement environment and source evidence</summary>

Healthy run: 2026-11-21T17:28:07Z; outage run: 2026-11-21T17:28:45Z.
go1.27.4, linux/amd64, AMD Ryzen 9 5900HX with Radeon Graphics; kernel 6.12.100+deb13-amd64.
Collector CPU: `0`; harness CPUs: `2`; buffer filesystem: `ext4`.
Container manager: testcontainers-go v0.45.0; slogx build: go1.27.4.

[Healthy samples](../results/healthy.jsonl) · [Outage samples](../results/outage.jsonl) ·
[Host environment](../results/environment.txt) · [Module versions](../results/collectors-modules.txt) ·
[Source checksums](../results/source.sha256)

The sample files retain configured image digests, resolved runtime image IDs,
collector module and binary identities, all measurements, and completion records.

</details>

<!-- benchmark-results:end -->

## Methodology

Each sample first delivers a separate warmup record through the complete
pipeline. Timing then starts before a 64 KiB buffered producer writes the
pre-generated records. Stdin remains open until the receiver validates all
expected IDs and fields and flushes its final successful HTTP acknowledgement.
Closing stdin and waiting for successful collector shutdown happens afterward.

The receiver rejects malformed JSON, unexpected fields, changed values, and
out-of-range IDs. Missing records, invalid payloads, shutdown errors, or timeouts
fail the sample. Duplicate deliveries are counted separately. Only validated
unique records contribute to throughput. Output ends with a `complete` record
only when the entire requested matrix succeeds; incomplete files are not results.

| Metric | Meaning |
| --- | --- |
| `seconds`, `records_per_second` | First producer write through final unique-record acknowledgement at the receiver. Includes pipe blocking, collection, buffering, encoding, HTTP transport, and receiver work. |
| `input_mib_per_second` | Original NDJSON bytes, including newlines, divided by elapsed time. |
| `producer_seconds` | Time to write and flush all input to Docker's attached stdin stream; this does not confirm ingestion or durability. For the HTTP control, this includes its request/response work. |
| `cpu_seconds` | Collector process and thread CPU counter difference from immediately before input to the snapshot after final receipt. Excludes harness and Docker daemon CPU. Small scheduling gaps surround the measured interval. |
| `peak_rss_bytes` | Collector process lifetime high-water RSS through the final snapshot, including startup, warmup, and resident file mappings. Excludes Docker overhead and filesystem cache outside those mappings. |
| `http_requests`, `retry_requests`, `duplicates` | Measured-phase requests, requests rejected with HTTP 503, and extra accepted copies of already-seen IDs. |
| `recovery_drain_seconds` | Time from the receiver becoming available to final unique-record acknowledgement. Includes any remaining input/journal work and retry scheduling. |

CPU counters have the host's `CLK_TCK` resolution, 10 ms on the measured host.
A reported zero means less than one observable counter increment, not free work.
Throughput ends before final collector checkpointing and shutdown; their success
is verified but their duration is excluded.

The `http-control` writes directly to the same checked receiver with the same
500-record/4 MiB batch limits. It has no collector, stdin attachment, or disk buffer.
Use it to identify receiver limits, not as a collector performance score.

With `-outage=2s`, the warmed receiver rejects requests for two seconds while
input continues. Every sample must actually encounter HTTP 503 and subsequently
deliver every record. This checks transient retry and drain behavior, not
process-crash recovery or power-loss durability.

These finite bursts are sensitive to flush intervals and batch boundaries.
They do not establish maximum sustained throughput, concurrent-source scaling,
individual-record tail latency, disk-full behavior, real backend ingestion
capacity, or OTLP trace conversion performance.

## Collectors and configuration

| Collector | Version | Persistent buffering | HTTP batching |
| --- | --- | --- | --- |
| [slogx-collector](https://github.com/rah-0/slogx-collector) | 0.7.0 | Disk journal; sync at 500 records, 4 MiB of framed data, or the 1 s timer. Only synced records are delivered. Checkpoint every 16 successful batches or on the same timer. | Up to 500 records or 4 MiB; 100 ms partial flush. |
| [Vector](https://vector.dev/releases/0.62.0/) | 0.62.0 | Disk buffer, 512 MiB capacity, blocking when full; synchronization interval approximately 500 ms. Delivery can start before disk synchronization. | Up to 500 records or 4 MiB; 100 ms partial flush. |
| [Fluent Bit](https://fluentbit.io/announcements/v5.2.0/) | 5.2.0 | Filesystem chunks with `storage.sync=full` and checksums. Synchronization is chunk-based. | 100 ms flush; native chunk boundaries, without an equivalent event-count limit. |

The [Vector configuration](config/vector.toml),
[Fluent Bit configuration](config/fluent-bit.conf), and
[collector arguments and pinned images](config.go) accompany the runner.
Runtime images are pinned to Linux/amd64 manifests. Testcontainers pulls missing
images and manages the containers. The [provisioner](prepare.go) builds
`github.com/rah-0/slogx-collector@v0.7.0` with the local Go compiler, `-trimpath`,
and CGO disabled, then Testcontainers copies the static binary into the pinned
Debian Trixie runtime. Preparation happens before measurement. Results record
the configured images, resolved runtime image IDs, harness and collector Go
versions, the pinned collector module, and the collector binary's SHA256.

Slogx also syncs pending journal writes on EOF or graceful shutdown and checkpoints
successful deliveries before exit. Segment rotation can sync earlier. Journal
commits are separate from the 100 ms HTTP partial flush. A crash can lose records
since the last successful journal sync and replay acknowledged batches since the
last checkpoint. Segments are removed only after their acknowledgement checkpoint
is durable.

All outputs use one worker or request at a time, no compression or TLS, and
one-second retry delays. Vector retry jitter is disabled; Fluent Bit's retry
base and cap are both one second. Scheduling can still affect actual retry time.
The runner sets `GOMAXPROCS=1` and `VECTOR_THREADS=1` inside containers.

Slogx timestamp mapping is disabled. Vector uses the metadata namespace so its
source metadata does not alter the JSON body. Fluent Bit omits its additional
timestamp key. Slogx sends JSON arrays; the peers send NDJSON. The receiver
accepts both encodings and returns successful record counts for the OpenObserve
acknowledgement contract. It does not run an OpenObserve backend.

The persistence distinction follows the
[Vector HTTP sink reference](https://vector.dev/docs/reference/configuration/sinks/http/),
its pinned [disk buffer implementation](https://github.com/vectordotdev/vector/blob/v0.62.0/lib/vector-buffers/src/variants/disk_v2/writer.rs),
and Fluent Bit's pinned
[chunk management](https://github.com/fluent/fluent-bit/blob/v5.2.0/src/flb_input_chunk.c)
and [synchronization implementation](https://github.com/fluent/fluent-bit/blob/v5.2.0/lib/chunkio/src/cio_file_unix.c).
Retry settings follow the
[Fluent Bit scheduler reference](https://docs.fluentbit.io/manual/administration/scheduling-and-retries).
Disk buffering does not make successful writes to a producer pipe durable.

## Running

Requires Go 1.27.4+, a local Linux/amd64 Docker engine, access to its containers'
host `/proc` entries, and `getconf`. Docker Desktop, remote Docker engines, and
other architectures are outside this runner's supported environment.

This directory is a separate Go module containing the Testcontainers dependencies.
Run it from the repository root with `go -C collectors`. Go downloads module
dependencies, and the runner builds the collector and pulls missing images
automatically:

```sh
# Run all collectors and the direct HTTP control, three times at each size.
go -C collectors run .

# Reproduce the measured CPU placement and retain raw results.
mkdir -p results
taskset -c 2 go -C collectors run . -cpu=0 > results/healthy.jsonl
taskset -c 2 go -C collectors run . -cpu=0 -outage=2s > results/outage.jsonl
```

Choose available CPUs on the machine. The example places the producer, HTTP
receiver, and container stdin attachment on CPU 2 and each collector on CPU 0.
The attachment uses the Docker API within the harness process. Only one collector
instance runs at a time; preparation failures stop the run.

Use `-records`, `-sizes`, `-repeat`, `-collectors`, `-cpu`, `-timeout`, and
`-outage` to select workloads. The default is 2,000 records at 256 and 1,024 bytes,
with three repetitions. Sizes include the newline and must be 256–16,384 bytes.
For example:

```sh
go -C collectors run . -collectors=slogx-collector,vector,fluent-bit \
  -records=20000 -sizes=1024 -repeat=3 -timeout=10m
```

`-storage-root` defaults to `/var/tmp`. Set it to an existing directory on the
persistent filesystem being measured. Every sample gets a fresh private buffer
directory bind-mounted into its container; sample containers and directories
are removed after completion or failure. The temporary collector binary is
removed when the run ends; pulled images remain available for subsequent runs.
Using tmpfs measures different persistence behavior.

## Tests

Unit tests exercise fixture generation, receipt validation, duplicate handling,
outages, cancellation, build-environment isolation, stdin EOF with output still
open, process-counter parsing, and bounded diagnostics without starting Docker
containers:

```sh
go -C collectors test -count=1 -race -cover -covermode=atomic ./...
go -C collectors vet ./...
```
