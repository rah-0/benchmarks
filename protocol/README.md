# Protocols (sending only)

> [!WARNING]
> These tests were performed locally, not sending over a real wire.

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
