# Results

```
MariaDB
SingleInsertFixedData-8	        193053 ns/op        752 B/op        21 allocs/op # innodb-flush-log-at-trx-commit is set to 0
SingleInsertRandomData-8        195988 ns/op	    816 B/op	    23 allocs/op # innodb-flush-log-at-trx-commit is set to 0
SingleInsertFixedData-8         312530 ns/op        752 B/op        21 allocs/op
SingleInsertRandomData-8        312220 ns/op        816 B/op        23 allocs/op
Insert1MilAndFindMiddle-8       164306697 ns/op     987 B/op	    26 allocs/op

PostgreSQL
SingleInsertFixedData-8         392730 ns/op	    861 B/op	    17 allocs/op
SingleInsertRandomData-8        407166 ns/op	    925 B/op	    19 allocs/op
Insert1MilAndFindMiddle-8       53271012 ns/op      864 B/op        22 allocs/op

SQLite
SingleInsertFixedData-8         687551 ns/op        736 B/op        18 allocs/op
SingleInsertRandomData-8        717496 ns/op        800 B/op        20 allocs/op
Insert1MilAndFindMiddle-8       63418844 ns/op      728 B/op        23 allocs/op
```

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
fsync = on
synchronous_commit = on       
wal_buffers = 16MB
commit_delay = 0
checkpoint_completion_target = 0.9
checkpoint_timeout = 15min
max_wal_size = 8GB          
wal_writer_delay = 200ms

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


