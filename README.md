# Benchmarks
These benchmarks are executed on a Virtual Machine.  
MariaDB is installed in the VM itself to avoid networking delays.

## About the system
CPU: AMD Ryzen 9 5950X  
RAM: 64GB (4x16GB) 4266MHz DDR4 CL19  

## About the VM
CPU: 8 threads  
RAM: 16GB  
OS: Debian 12, Kernel: 6.1.0-23-amd64  
NIC: 1GB


# Results
The aim of benching is to see what's faster while using the least amount of memory possible for each category.

---
## Protocols (sending only)

| Message Size       | Protocol | Time per Op     | Memory per Op | Allocations per Op |
|--------------------|----------|-----------------|---------------|--------------------|
| **Single Message** |          |                 |               |                    |
|                    | TCP      | 4.32 µs         | 0 B           | 0                  |
|                    | HTTP     | 63.52 µs        | 6.07 KB       | 67                 |
| **2KB**            |          |                 |               |                    |
|                    | TCP      | 5.14 µs         | 0 B           | 0                  |
|                    | HTTP     | 73.73 µs        | 10.45 KB      | 72                 |
| **4KB**            |          |                 |               |                    |
|                    | TCP      | 5.61 µs         | 0 B           | 0                  |
|                    | HTTP     | 90.24 µs        | 17.96 KB      | 78                 |
| **8KB**            |          |                 |               |                    |
|                    | TCP      | 5.79 µs         | 0 B           | 0                  |
|                    | HTTP     | 110.38 µs       | 44.98 KB      | 81                 |
| **16KB**           |          |                 |               |                    |
|                    | TCP      | 6.33 µs         | 0 B           | 0                  |
|                    | HTTP     | 136.30 µs       | 83.47 KB      | 84                 |
| **32KB**           |          |                 |               |                    |
|                    | TCP      | 12.68 µs        | 0 B           | 0                  |
|                    | HTTP     | 196.97 µs       | 197.12 KB     | 89                 |
| **64KB**           |          |                 |               |                    |
|                    | TCP      | 27.72 µs        | 0 B           | 0                  |
|                    | HTTP     | 267.97 µs       | 329.33 KB     | 93                 |
| **128KB**          |          |                 |               |                    |
|                    | TCP      | 52.94 µs        | 0 B           | 0                  |
|                    | HTTP     | 382.92 µs       | 558.91 KB     | 97                 |
| **256KB**          |          |                 |               |                    |
|                    | TCP      | 104.00 µs       | 0 B           | 0                  |
|                    | HTTP     | 719.29 µs       | 1.19 MB       | 105                |
| **512KB**          |          |                 |               |                    |
|                    | TCP      | 203.25 µs       | 0 B           | 0                  |
|                    | HTTP     | 1.46 ms         | 2.49 MB       | 113                |
| **1MB**            |          |                 |               |                    |
|                    | TCP      | 411.23 µs       | 0 B           | 0                  |
|                    | HTTP     | 3.45 ms         | 5.07 MB       | 118                |
| **10MB**           |          |                 |               |                    |
|                    | TCP      | 3.93 ms         | 0 B           | 0                  |
|                    | HTTP     | 18.55 ms        | 49.89 MB      | 130                |
| **100MB**          |          |                 |               |                    |
|                    | TCP      | 39.15 ms        | 0 B           | 0                  |
|                    | HTTP     | 148.33 ms       | 586.89 MB     | 139                |
| **1GB**            |          |                 |               |                    |
|                    | TCP      | 400.58 ms       | 0 B           | 0                  |
|                    | HTTP     | 1.16 s          | 5.34 GB       | 149                |


## Serializers

| Size                   | Encoding Type      | Time per Op     | Memory per Op   | Allocations per Op |
|------------------------|--------------------|-----------------|-----------------|--------------------|
| **1Small-8**           |                    |                 |                 |                    |
|                        | GOB                | 16.44 µs        | 9.18 KB         | 229                |
|                        | GOB (Optimized)    | 774.7 ns        | 376 B           | 8                  |
|                        | JSON               | 1.70 µs         | 536 B           | 11                 |
|                        | Protobuf           | 365.1 ns        | 256 B           | 4                  |
| **100Small-8**         |                    |                 |                 |                    |
|                        | GOB                | 59.67 µs        | 71.40 KB        | 437                |
|                        | GOB (Optimized)    | 23.04 µs        | 21.51 KB        | 206                |
|                        | JSON               | 110.95 µs       | 32.37 KB        | 216                |
|                        | Protobuf           | 36.52 µs        | 25.00 KB        | 400                |
| **10000Small-8**       |                    |                 |                 |                    |
|                        | GOB                | 4.00 ms         | 6.76 MB         | 20,254             |
|                        | GOB (Optimized)    | 1.82 ms         | 2.05 MB         | 20,006             |
|                        | JSON               | 11.40 ms        | 4.53 MB         | 20,034             |
|                        | Protobuf           | 3.56 ms         | 2.44 MB         | 40,000             |
| **1000000Small-8**     |                    |                 |                 |                    |
|                        | GOB                | 309.24 ms       | 1.11 GB         | 2,000,289          |
|                        | GOB (Optimized)    | 233.99 ms       | 592.67 MB       | 2,000,019          |
|                        | JSON               | 965.17 ms       | 611.46 MB       | 2,000,088          |
|                        | Protobuf           | 305.25 ms       | 244.12 MB       | 4,000,002          |
| **1Unreal-8**          |                    |                 |                 |                    |
|                        | GOB                | 104.47 µs       | 55.38 KB        | 1,014              |
|                        | GOB (Optimized)    | 5.54 µs         | 10.57 KB        | 106                |
|                        | JSON               | 42.52 µs        | 11.54 KB        | 109                |
|                        | Protobuf           | 8.67 µs         | 10.43 KB        | 102                |
| **10Unreals-8**        |                    |                 |                 |                    |
|                        | GOB                | 265.22 µs       | 320.96 KB       | 1,922              |
|                        | GOB (Optimized)    | 49.53 µs        | 102.97 KB       | 1,006              |
|                        | JSON               | 424.56 µs       | 153.83 KB       | 1,014              |
|                        | Protobuf           | 85.58 µs        | 104.38 KB       | 1,020              |
| **100Unreals-8**       |                    |                 |                 |                    |
|                        | GOB                | 1.56 ms         | 3.26 MB         | 10,931             |
|                        | GOB (Optimized)    | 487.03 µs       | 1005.86 KB      | 10,006             |
|                        | JSON               | 4.38 ms         | 1.49 MB         | 10,020             |
|                        | Protobuf           | 781.50 µs       | 1.02 MB         | 10,200             |
| **1000Unreals-8**      |                    |                 |                 |                    |
|                        | GOB                | 13.77 ms        | 33.52 MB        | 100,941            |
|                        | GOB (Optimized)    | 5.04 ms         | 9.73 MB         | 100,006            |
|                        | JSON               | 43.05 ms        | 19.68 MB        | 100,034            |
|                        | Protobuf           | 6.63 ms         | 10.19 MB        | 102,000            |

---

## DB's

### MariaDB
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

innodb-flush-log-at-trx-commit=1
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

net-read-timeout=3600
net-write-timeout=3600
```

```
BenchmarkMariaDBSingleInsertFixedData-8           226855            312530 ns/op             752 B/op         21 allocs/op
BenchmarkMariaDBSingleInsertRandomData-8          229183            312220 ns/op             816 B/op         23 allocs/op
```
