# Summary

## Redis
Version: Redis server v=7.0.15 sha=00000000:0 malloc=jemalloc-5.3.0 bits=64 build=5f65bccf5d6cab79

## Memcached
Version: 1.6.18

## Results
| Size (KB) | Type      | Avg Time (µs/op) | Avg Mem (KB/op) | Allocations/op |
|-----------|-----------|------------------|-----------------|----------------|
| 2         | Memcache  | 148.47           | 6.79            | 16             |
| 2         | Redis     | 149.68           | 2.78            | 17             |
| 4         | Memcache  | 157.56           | 13.45           | 16             |
| 4         | Redis     | 162.86           | 5.34            | 17             |
| 8         | Memcache  | 162.90           | 26.26           | 16             |
| 8         | Redis     | 163.82           | 9.95            | 17             |
| 16        | Memcache  | 169.30           | 51.63           | 16             |
| 16        | Redis     | 170.42           | 18.91           | 17             |
| 32        | Memcache  | 186.44           | 106.97          | 16             |
| 32        | Redis     | 183.07           | 41.44           | 17             |
| 64        | Memcache  | 215.77           | 205.37          | 16             |
| 64        | Redis     | 201.41           | 74.20           | 17             |
| 128       | Memcache  | 268.19           | 402.18          | 17             |
| 128       | Redis     | 235.08           | 139.74          | 17             |
| 256       | Memcache  | 409.56           | 795.71          | 18             |
| 256       | Redis     | 294.38           | 270.81          | 17             |
| 512       | Memcache  | 830.66           | 1582.79         | 20             |
| 512       | Redis     | 451.52           | 532.96          | 17             |
