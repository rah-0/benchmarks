# Serializers

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
