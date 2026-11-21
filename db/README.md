# Database benchmarks

Compare matching inserts and indexed queries. **PostgreSQL is the reference for
networked server comparisons with MariaDB.** Embedded SQLite is shown separately
because its execution and networking costs differ. The configured durability
settings also differ; these are workload measurements, not equal-durability rankings.

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

**Workloads:** [Insert a fixed value](#insert-a-fixed-value) · [Generate and insert a UUID](#generate-and-insert-a-uuid) · [Indexed lookup in 1 million rows](#indexed-lookup-in-1-million-rows) · [Read and order 100,000 rows](#read-and-order-100000-rows) · [Filter and order 80,000 matches](#filter-and-order-80000-matches) · [Embedded SQLite costs](#embedded-sqlite-costs)

### Insert a fixed value

Prepare and execute one insert of a fixed string; the table grows during timing.

Baseline: **PostgreSQL = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

Both clients use local container TCP connections. Server memory is excluded; configured durability and caching differ between engines.

| Implementation | Time/query, µs ↓ | Speed vs baseline ↑ | Allocated/query, B ↓ | Allocs/query ↓ |
| --- | ---: | ---: | ---: | ---: |
| **PostgreSQL (baseline)** | 292.71 | 1.00× | 1,149 | 26 |
| MariaDB | 269.95 | 1.08× | 824 | 20 |

MariaDB has the lowest recorded median time: 1.08× the speed of PostgreSQL, with 0.72× its allocated bytes per operation.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| PostgreSQL | 292.62–296.34 µs | 1.12–1.12 KiB | 25–26 | 1.00× |
| MariaDB | 268.87–271.78 µs | 824–824 B | 20 | 0.72× |

</details>

### Generate and insert a UUID

Generate a UUID, then prepare and execute its insert; the table grows during timing.

Baseline: **PostgreSQL = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

Both clients use local container TCP connections. Server memory is excluded; configured durability and caching differ between engines.

| Implementation | Time/query, µs ↓ | Speed vs baseline ↑ | Allocated/query, B ↓ | Allocs/query ↓ |
| --- | ---: | ---: | ---: | ---: |
| **PostgreSQL (baseline)** | 298.32 | 1.00× | 1,216 | 28 |
| MariaDB | 289.54 | 1.03× | 888 | 22 |

MariaDB has the lowest recorded median time: 1.03× the speed of PostgreSQL, with 0.73× its allocated bytes per operation. The timing ranges of MariaDB and PostgreSQL overlap; their order is inconclusive.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| PostgreSQL | 291.96–299.79 µs | 1.19–1.19 KiB | 28 | 1.00× |
| MariaDB | 277.76–292.12 µs | 888–888 B | 22 | 0.73× |

</details>

### Indexed lookup in 1 million rows

Look up and verify one indexed value in a prepared one-million-row table. Population is excluded.

Baseline: **PostgreSQL = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

Both clients use local container TCP connections. Server memory is excluded; configured durability and caching differ between engines.

| Implementation | Time/query, µs ↓ | Speed vs baseline ↑ | Allocated/query, B ↓ | Allocs/query ↓ |
| --- | ---: | ---: | ---: | ---: |
| **PostgreSQL (baseline)** | 229.42 | 1.00× | 1,008 | 27 |
| MariaDB | 216.85 | 1.06× | 800 | 24 |

MariaDB has the lowest recorded median time: 1.06× the speed of PostgreSQL, with 0.79× its allocated bytes per operation.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| PostgreSQL | 224.26–233.41 µs | 1,008–1,015 B | 27 | 1.00× |
| MariaDB | 213.64–217.78 µs | 800–801 B | 24 | 0.79× |

</details>

### Read and order 100,000 rows

Read 100,000 rows ordered by value, scan them, and verify their ordering and count.

Baseline: **PostgreSQL = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

Both clients use local container TCP connections. Server memory is excluded; configured durability and caching differ between engines.

| Implementation | Time/query, ms ↓ | Speed vs baseline ↑ | Allocated/query, MiB ↓ | Allocs/query ↓ |
| --- | ---: | ---: | ---: | ---: |
| **PostgreSQL (baseline)** | 36.31 | 1.00× | 7.63 | 300,016 |
| MariaDB | 18.6 | 1.95× | 8.39 | 300,016 |

MariaDB has the lowest recorded median time: 1.95× the speed of PostgreSQL, with 1.10× its allocated bytes per operation.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| PostgreSQL | 36.09–36.49 ms | 7.63–7.63 MiB | 300,016 | 1.00× |
| MariaDB | 18.3–18.82 ms | 8.39–8.39 MiB | 300,016 | 1.10× |

</details>

### Filter and order 80,000 matches

Filter a prepared 100,000-row table by a fixed timestamp cutoff, scan the 80,000 matching rows, and verify their order and count.

Baseline: **PostgreSQL = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

Both clients use local container TCP connections. Server memory is excluded; configured durability and caching differ between engines.

| Implementation | Time/query, ms ↓ | Speed vs baseline ↑ | Allocated/query, MiB ↓ | Allocs/query ↓ |
| --- | ---: | ---: | ---: | ---: |
| **PostgreSQL (baseline)** | 43.48 | 1.00× | 11.6 | 480,027 |
| MariaDB | 40.58 | 1.07× | 10.38 | 400,022 |

MariaDB has the lowest recorded median time: 1.07× the speed of PostgreSQL, with 0.89× its allocated bytes per operation.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| PostgreSQL | 43.23–44.01 ms | 11.6–11.6 MiB | 480,027–480,030 | 1.00× |
| MariaDB | 40.54–40.77 ms | 10.38–10.38 MiB | 400,022 | 0.89× |

</details>

### Embedded SQLite costs

The same insert and read workloads use an embedded engine and a private persistent database file. Read fixture population is excluded; scans include count and ordering checks.

No equivalent reference measured; absolute costs are shown. Median of 3 samples; GOMAXPROCS=16.

SQLite has no client/server TCP round trip. Its costs are shown separately from the networked database comparison; the five rows also perform different operations.

| Operation / variant | Time/query ↓ | Allocated/query ↓ | Allocs/query ↓ |
| --- | ---: | ---: | ---: |
| Insert a fixed value | 19.77 µs | 840 B | 15 |
| Generate and insert a UUID | 22.22 µs | 904 B | 17 |
| Indexed lookup in 1 million rows | 7.26 µs | 800 B | 21 |
| Read and order 100,000 rows | 23.93 ms | 7.63 MiB | 300,017 |
| Filter and order 80,000 matches | 65.17 ms | 12.21 MiB | 480,024 |

<details>
<summary>Observed sample ranges</summary>

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| Insert a fixed value | 19.77–19.92 µs | 840–840 B | 15 |
| Generate and insert a UUID | 22.21–22.35 µs | 904–904 B | 17 |
| Indexed lookup in 1 million rows | 7.25–7.27 µs | 800–800 B | 21 |
| Read and order 100,000 rows | 23.89–23.94 ms | 7.63–7.63 MiB | 300,016–300,017 |
| Filter and order 80,000 matches | 64.41–65.6 ms | 12.21–12.21 MiB | 480,024 |

</details>

**Recorded samples:** [db/sqlite](../results/db-sqlite.txt) · [db/mariadb](../results/db-mariadb.txt) · [db/postgresql](../results/db-postgresql.txt)

<!-- benchmark-results:end -->

## Methodology

Each table has a string value and insertion time, with indexes on both columns.
Every sample resets its table. Read fixtures are populated in a transaction before
`testing.B.Loop` starts timing; setup and previous samples cannot add measured rows.

| Workload | Timed work per operation |
| --- | --- |
| Insert a fixed value | Prepare and execute one insert with a fixed string. |
| Insert a random value | Generate a UUID, then prepare and execute one insert. |
| Indexed lookup in 1 million rows | Find and verify one value in a prepared table. Population is excluded. |
| Read 100,000 sorted rows | Read all rows ordered by value, scan them, and verify order and count. |
| Read 80,000 matching rows | Filter a 100,000-row fixture using a fixed time cutoff, order by timestamp, scan, and verify the matches. |

Insert tables grow during the measured loop. Read fixtures use independently
generated UUIDs and a fixed timestamp distribution. Go allocations exclude the
MariaDB and PostgreSQL server processes.

### Engine settings

Testcontainers starts MariaDB 13.0.2 and PostgreSQL 19.0 on dynamic local TCP
ports and removes them afterward. Exact image pins and options are in the
[MariaDB setup](mariadb/mariadb_test.go) and
[PostgreSQL setup](postgresql/postgresql_test.go).

| Engine | Durability and memory configuration |
| --- | --- |
| MariaDB | `innodb-flush-log-at-trx-commit=0`, `sync-binlog=0`, 1 GiB buffer pool. |
| PostgreSQL | `fsync=off`, `synchronous_commit=off`, 1 GiB shared buffers. |
| SQLite | WAL, `synchronous=OFF`, memory temporary storage, exclusive locking, 2 GiB page-cache limit, 8 GiB mmap limit. |

[SQLite setup](sqlite/sqlite.go) applies its PRAGMAs to one retained connection
and uses a private temporary database. The [full runner](../scripts/benchmark.bash)
puts temporary files under `/var/tmp` or `BENCH_STORAGE_ROOT`, avoiding an
unintended in-memory measurement when `/tmp` is tmpfs. Compare the recorded
filesystem when comparing runs. The SQLite engine comes from the pinned
`go-sqlite3` dependency.

## Running

Requires Go 1.27.4+, a C compiler for SQLite, and a local Docker engine.
From the repository root:

```sh
go test -p=1 -run '^$' -bench . -benchmem -benchtime=1s -count=3 -timeout=2h ./db/...
```

Use the full runner to reproduce its persistent temporary-storage placement.
