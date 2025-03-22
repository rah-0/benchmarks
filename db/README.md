# DB's

MariaDB is installed in the VM itself to avoid networking delays.

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
MariaDB
BenchmarkMariaDBSingleInsertFixedData-8         312530 ns/op        752 B/op        21 allocs/op
BenchmarkMariaDBSingleInsertRandomData-8        312220 ns/op        816 B/op        23 allocs/op
Note: if innodb-flush-log-at-trx-commit is set to 0 the results are:
BenchmarkMariaDBSingleInsertFixedData-8	        193053 ns/op        752 B/op        21 allocs/op
BenchmarkMariaDBSingleInsertRandomData-8        195988 ns/op	    816 B/op	    23 allocs/op

SQLite
BenchmarkSQLiteSingleInsertFixedData-8          687551 ns/op        736 B/op        18 allocs/op
BenchmarkSQLiteSingleInsertRandomData-8         717496 ns/op        800 B/op        20 allocs/op
```
