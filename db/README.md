# Results

All database engines are configured to minimize durability guarantees and maximize raw throughput, ensuring fair, in-memory-heavy benchmarking without I/O bottlenecks or network delays. All databases are installed directly on the same VM.

Results are ordered from the fastest CPU bench to slowest per category.

| Category                     | Database | MEM               | CPU               |
|------------------------------|----------|-------------------|-------------------|
| Single Insert (Fixed)        |          |                   |                   |
|                              | SQLite   | 784 B             | 16 µs             |
|                              | Hyperion | 480 B             | 106 µs            |
|                              | MariaDB  | 800 B             | 197 µs            |
|                              | Postgres | 912 B             | 275 µs            |
| Single Insert (Random)       |          |                   |                   |
|                              | SQLite   | 848 B             | 46 µs             |
|                              | Hyperion | 544 B             | 106 µs            |
|                              | MariaDB  | 864 B             | 199 µs            |
|                              | Postgres | 976 B             | 290 µs            |
| Insert 1M + Find Middle      |          |                   |                   |
|                              | SQLite   | 728 B             | 6 µs              |
|                              | Hyperion | 888 B             | 100 µs            |
|                              | MariaDB  | 792 B             | 222 µs            |
|                              | Postgres | 864 B             | 201 µs            |
| Insert 100K + Order + Check  |          |                   |                   |
|                              | MariaDB  | 8MiB 402KiB 442B  | 18ms 392µs 175ns  | 
|                              | Postgres | 7MiB 645KiB 194B  | 18ms 897µs 904ns  |  
|                              | SQLite   | 7MiB 645KiB 229B  | 35ms 931µs 155ns  |
|                              | Hyperion | 29MiB 978KiB 328B | 203ms 496µs 979ns |
| Insert 100K + Order + Filter |          |                   |                   | 
|                              | Postgres | 12MiB 212KiB 953B | 35ms 884µs 984ns  |
|                              | MariaDB  | 11MiB 611KiB 920B | 45ms 475µs 517ns  |
|                              | SQLite   | 29MiB 930KiB 899B | 130ms 278µs 472ns | 
|                              | Hyperion | 26MiB 579KiB 297B | 152ms 145µs 757ns |

# Configs and details
DB's are installed in the VM itself to avoid networking delays.

## MariaDB
Version: 11.4.2-MariaDB-deb12  
Config:
```server.cnf
default-storage-engine=INNODB
sql-mode="STRICT_TRANS_TABLES,NO_ENGINE_SUBSTITUTION,NO_BACKSLASH_ESCAPES"

character-set-server=utf8mb4
collation-server=utf8mb4_unicode_ci
character-set-client-handshake=utf8mb4
default-time-zone=+00:00

join-buffer-size=512M
max-allowed-packet=2G
sort-buffer-size=128M
table-definition-cache=2000
max-connections=3000
tmp-table-size=512M

innodb-flush-log-at-trx-commit=0
innodb-log-buffer-size=512M
innodb-buffer-pool-size=16G
innodb-buffer-pool-instances=8
innodb-thread-concurrency=8
innodb-log-files-in-group=4
innodb-log-file-size=8G
innodb-write-io-threads=4
innodb-read-io-threads=4
innodb-autoextend-increment=256
innodb-old-blocks-time=500
innodb-file-per-table=ON
```

## PostgreSQL
Version: psql 15.12 (Debian 15.12-0+deb12u2)  
Config: 
```
# Encoding and collation
client_encoding = 'UTF8'
timezone = 'UTC'

# Memory settings 
work_mem = 128MB               
maintenance_work_mem = 512MB
shared_buffers = 16GB         
effective_cache_size = 48GB   

# Write behavior
wal_level = replica
fsync = off
synchronous_commit = off       
wal_buffers = 16MB
commit_delay = 0
checkpoint_completion_target = 0.9
checkpoint_timeout = 15min
max_wal_size = 8GB          
wal_writer_delay = 1000ms

# Temporary tables and buffers
temp_buffers = 64MB
temp_file_limit = 2GB
max_files_per_process = 1000

# Parallelism and concurrency
max_worker_processes = 8
max_parallel_workers = 8
max_parallel_workers_per_gather = 4
parallel_leader_participation = on
```

## SQLite
Version: 3.40.1 2022-12-28 14:03:47 df5c253c0b3dd24916e4ec7cf77d3db5294cc9fd45ae7b9c5e82ad8197f3alt1  
Config:
```
PRAGMA synchronous = OFF
PRAGMA journal_mode = WAL
PRAGMA cache_size = -2097152
PRAGMA temp_store = MEMORY
PRAGMA locking_mode = EXCLUSIVE
PRAGMA mmap_size = 8589934592
```
