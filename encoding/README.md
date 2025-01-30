# Benchmark Results

| Size (KB) | Encoder   | ns/op     | B/op      | allocs/op  |
|-----------|:----------|:----------|:----------|------------|
| 2 KB      | Base64    | 5.22 µs   | 3.07 KB   | 1          |
| 2 KB      | Base85    | 9.97 µs   | 2.69 KB   | 1          |
| 4 KB      | Base64    | 10.33 µs  | 6.14 KB   | 1          |
| 4 KB      | Base85    | 19.90 µs  | 5.38 KB   | 1          |
| 8 KB      | Base64    | 20.81 µs  | 12.29 KB  | 1          |
| 8 KB      | Base85    | 39.13 µs  | 10.24 KB  | 1          |
| 16 KB     | Base64    | 41.57 µs  | 24.58 KB  | 1          |
| 16 KB     | Base85    | 77.83 µs  | 20.48 KB  | 1          |
| 32 KB     | Base64    | 81.07 µs  | 49.15 KB  | 1          |
| 32 KB     | Base85    | 158.14 µs | 40.96 KB  | 1          |
| 64 KB     | Base64    | 171.59 µs | 90.11 KB  | 1          |
| 64 KB     | Base85    | 317.44 µs | 81.92 KB  | 1          |
| 128 KB    | Base64    | 340.39 µs | 180.23 KB | 1          |
| 128 KB    | Base85    | 636.96 µs | 163.84 KB | 1          |
| 256 KB    | Base64    | 684.27 µs | 352.26 KB | 1          |
| 256 KB    | Base85    | 1.28 ms   | 327.68 KB | 1          |
| 512 KB    | Base64    | 1.39 ms   | 704.52 KB | 1          |
| 512 KB    | Base85    | 2.58 ms   | 655.37 KB | 1          |
| 1 MB      | Base64    | 2.79 ms   | 1.40 MB   | 2          |
| 1 MB      | Base85    | 5.32 ms   | 1.31 MB   | 1          |
| 10 MB     | Base64    | 22.82 ms  | 13.98 MB  | 2          |
| 10 MB     | Base85    | 47.19 ms  | 13.11 MB  | 2          |
| 100 MB    | Base64    | 224.19 ms | 139.81 MB | 2          |
| 100 MB    | Base85    | 467.38 ms | 131.07 MB | 2          |
| 1 GB      | Base64    | 2.23 s    | 1.43 GB   | 3          |
| 1 GB      | Base85    | 4.76 s    | 1.34 GB   | 3          |
