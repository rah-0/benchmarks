# Log collector benchmarks

Compare `slogx-collector`, Vector, and Fluent Bit by feeding the same synthetic
NDJSON through stdin and checking delivery to a local HTTP receiver. The suite
measures fixed bursts with filesystem buffering enabled. **The collectors have
different persistence guarantees; these are not equal-durability rankings.**

## Collectors and configuration

| Collector | Version | Persistent buffering | HTTP batching |
| --- | --- | --- | --- |
| [slogx-collector](https://github.com/rah-0/slogx-collector) | 0.6.0 | Disk journal; sync at 500 records, 4 MiB of framed data, or the 1 s timer. Only synced records are delivered. Checkpoint every 16 successful batches or on the same timer. | Up to 500 records or 4 MiB; 100 ms partial flush. |
| [Vector](https://vector.dev/releases/0.58.0/) | 0.58.0 | Disk buffer, 512 MiB capacity, blocking when full; synchronization interval approximately 500 ms. Delivery can start before disk synchronization. | Up to 500 records or 4 MiB; 100 ms partial flush. |
| [Fluent Bit](https://fluentbit.io/announcements/) | 5.1.2 | Filesystem chunks with `storage.sync=full` and checksums. Synchronization is chunk-based. | 100 ms flush; native chunk boundaries, without an equivalent event-count limit. |

The [Vector configuration](config/vector.toml),
[Fluent Bit configuration](config/fluent-bit.conf), and
[collector arguments and pinned images](config.go) accompany the runner.
Runtime images are pinned to Linux/amd64 manifests. Testcontainers pulls missing
images and manages the containers. The [provisioner](prepare.go) builds
`github.com/rah-0/slogx-collector@v0.6.0` with the local Go compiler, `-trimpath`,
and CGO disabled, then Testcontainers copies the static binary into the pinned
Debian Bookworm runtime. Preparation happens before measurement. Results record
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
its pinned [disk buffer implementation](https://github.com/vectordotdev/vector/blob/v0.58.0/lib/vector-buffers/src/variants/disk_v2/writer.rs),
and Fluent Bit's pinned
[chunk management](https://github.com/fluent/fluent-bit/blob/v5.1.2/src/flb_input_chunk.c)
and [synchronization implementation](https://github.com/fluent/fluent-bit/blob/v5.1.2/lib/chunkio/src/cio_file_unix.c).
Retry settings follow the
[Fluent Bit scheduler reference](https://docs.fluentbit.io/manual/administration/scheduling-and-retries).
Disk buffering does not make successful writes to a producer pipe durable.

## Running

Requires Go 1.27.1+, a local Linux/amd64 Docker engine, access to its containers'
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
taskset -c 2 go -C collectors run . -cpu=0 > collectors/results/healthy.jsonl
taskset -c 2 go -C collectors run . -cpu=0 -outage=2s > collectors/results/outage.jsonl
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

## Measurement boundaries

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

## Tests

Unit tests exercise fixture generation, receipt validation, duplicate handling,
outages, cancellation, build-environment isolation, stdin EOF with output still
open, process-counter parsing, and bounded diagnostics without starting Docker
containers:

```sh
go -C collectors test -count=1 -race -cover -covermode=atomic ./...
go -C collectors vet ./...
```

## Results

Measured on 2026-09-10 with slogx-collector 0.6.0, Go 1.27.1, Testcontainers for
Go 0.44.0, and Docker 29.7.1 on Linux/amd64, an AMD Ryzen 9 5900HX, Linux
6.12.100+deb13-amd64, and ext4 on a Samsung SSD 980 PRO 1TB. The collector uses
CPU 0; the harness uses CPU 2, on a separate physical core. Other test and
benchmark runs were not started concurrently.

The tables report medians of three samples, each delivering 2,000 records.
Raw samples, image identities, and environment metadata are retained in
[healthy.jsonl](results/healthy.jsonl) and [outage.jsonl](results/outage.jsonl).

### Healthy delivery

| Record bytes | Collector | Records/s | Input MiB/s | CPU ms | Peak RSS MiB | Producer ms | Requests |
| ---: | --- | ---: | ---: | ---: | ---: | ---: | ---: |
| 256 | slogx-collector | 41045.7 | 10.021 | 20 | 12.2 | 11.074 | 4 |
| 256 | vector | 56550.4 | 13.806 | 30 | 75.8 | 0.889 | 4 |
| 256 | fluent-bit | 18433.0 | 4.500 | 10 | 23.9 | 3.483 | 1 |
| 1024 | slogx-collector | 35478.9 | 34.647 | 30 | 15.0 | 35.212 | 4 |
| 1024 | vector | 46181.6 | 45.099 | 30 | 81.6 | 5.351 | 4 |
| 1024 | fluent-bit | 17315.8 | 16.910 | 20 | 27.2 | 11.787 | 1 |

Slogx-collector delivered the 2,000-record bursts in a median 48.7 ms at 256 bytes
and 56.4 ms at 1,024 bytes. Both bursts reach the 500-record journal commit
threshold before the one-second timer fires.

### Two-second receiver outage

| Record bytes | Collector | Recovery to final acknowledgement, s | Producer ms | HTTP 503 requests |
| ---: | --- | ---: | ---: | ---: |
| 256 | slogx-collector | 0.047 | 11.211 | 2 |
| 256 | vector | 0.043 | 0.804 | 2 |
| 256 | fluent-bit | 1.602 | 3.175 | 2 |
| 1024 | slogx-collector | 0.063 | 34.953 | 2 |
| 1024 | vector | 0.055 | 4.917 | 2 |
| 1024 | fluent-bit | 2.008 | 11.155 | 2 |

### Direct receiver control

| Record bytes | Healthy records/s | Healthy input MiB/s | Outage recovery to final acknowledgement, s |
| ---: | ---: | ---: | ---: |
| 256 | 274302.7 | 66.968 | 0.010 |
| 1024 | 157654.3 | 153.959 | 0.016 |

All 48 samples completed with zero missing, altered, or duplicate records.
Every outage sample observed HTTP 503 before successful delivery.
The direct receiver control was faster than every collector in both healthy
workloads. These burst results retain each collector's distinct batching and
persistence costs; compare them together with the configuration table above.
