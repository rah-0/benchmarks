# Summary

## Brotli
- Achieves the highest compression ratios, especially with repeated data
- Good all-rounder: will compress either small or big amount of data 
- Extremely slow and memory intensive

## Gzip
- Compressing small items can even yield a bigger data output
- Very high compression ratios for repeated data
- Not so good at random data

## LZ4
- The fastest but it can deal with repeated data only, it's quite bad at completely random data
- Compressing small items can even yield a bigger data output

# Benchmark Results

The data sample is 1mb of random data.

| Algorithm | Time (ms/op) | Memory Allocated (MB/op)  | Allocations/op |
|-----------|--------------|---------------------------|----------------|
| Brotli    | 1114.68 ms   | 131.4 MB                  | 119            |
| Gzip      | 25.99 ms     | 7.14 MB                   | 57             |
| LZ4       | 2.80 ms      | 10.0 MB                   | 38             |

# Compression Results (Random Data)
| Size       | Algorithm | Encoded Size | Compressed Size | Gain (%)  |
|------------|-----------|--------------|-----------------|-----------|
| **2KB**    | Brotli    | 2048 bytes   | 1709 bytes      | 16.55%    |
|            | Gzip      | 2048 bytes   | 1734 bytes      | 15.33%    |
|            | LZ4       | 2048 bytes   | 2067 bytes      | -0.93%    |
| **4KB**    | Brotli    | 4096 bytes   | 3391 bytes      | 17.21%    |
|            | Gzip      | 4096 bytes   | 3419 bytes      | 16.53%    |
|            | LZ4       | 4096 bytes   | 4115 bytes      | -0.46%    |
| **8KB**    | Brotli    | 8192 bytes   | 6776 bytes      | 17.29%    |
|            | Gzip      | 8192 bytes   | 6798 bytes      | 17.02%    |
|            | LZ4       | 8192 bytes   | 8211 bytes      | -0.23%    |
| **16KB**   | Brotli    | 16384 bytes  | 13526 bytes     | 17.44%    |
|            | Gzip      | 16384 bytes  | 13575 bytes     | 17.14%    |
|            | LZ4       | 16384 bytes  | 16403 bytes     | -0.12%    |
| **32KB**   | Brotli    | 32768 bytes  | 27064 bytes     | 17.41%    |
|            | Gzip      | 32768 bytes  | 27107 bytes     | 17.28%    |
|            | LZ4       | 32768 bytes  | 32787 bytes     | -0.06%    |
| **64KB**   | Brotli    | 65536 bytes  | 54126 bytes     | 17.41%    |
|            | Gzip      | 65536 bytes  | 54211 bytes     | 17.28%    |
|            | LZ4       | 65536 bytes  | 65555 bytes     | -0.03%    |
| **128KB**  | Brotli    | 131072 bytes | 108264 bytes    | 17.40%    |
|            | Gzip      | 131072 bytes | 108350 bytes    | 17.34%    |
|            | LZ4       | 131072 bytes | 131091 bytes    | -0.01%    |
| **256KB**  | Brotli    | 262144 bytes | 216549 bytes    | 17.39%    |
|            | Gzip      | 262144 bytes | 216699 bytes    | 17.34%    |
|            | LZ4       | 262144 bytes | 262163 bytes    | -0.01%    |
| **512KB**  | Brotli    | 524288 bytes | 433146 bytes    | 17.38%    |
|            | Gzip      | 524288 bytes | 433376 bytes    | 17.34%    |
|            | LZ4       | 524288 bytes | 524307 bytes    | -0.00%    |
| **1MB**    | Brotli    | 1048576 bytes| 866449 bytes    | 17.37%    |
|            | Gzip      | 1048576 bytes| 866921 bytes    | 17.32%    |
|            | LZ4       | 1048576 bytes| 1048595 bytes   | -0.00%    |
| **10MB**   | Brotli    | 10485760 bytes| 8665188 bytes  | 17.36%    |
|            | Gzip      | 10485760 bytes| 8667954 bytes  | 17.34%    |
|            | LZ4       | 10485760 bytes| 10485787 bytes | -0.00%    |

# Compression Results (Small structs with random data)
| Size         | Algorithm | Encoded Size  | Compressed Size | Gain (%)  |
|--------------|-----------|---------------|-----------------|-----------|
| **1Small**   | Brotli    | 147 bytes     | 130 bytes       | 11.56%    |
|              | Gzip      | 84 bytes      | 93 bytes        | -10.71%   |
|              | LZ4       | 85 bytes      | 104 bytes       | -22.35%   |
| **100Small** | Brotli    | 7943 bytes    | 4088 bytes      | 48.53%    |
|              | Gzip      | 7946 bytes    | 4435 bytes      | 44.19%    |
|              | LZ4       | 7952 bytes    | 7549 bytes      | 5.07%     |
| **10KSmall** | Brotli    | 793441 bytes  | 393685 bytes    | 50.38%    |
|              | Gzip      | 793320 bytes  | 426396 bytes    | 46.25%    |
|              | LZ4       | 793333 bytes  | 601696 bytes    | 24.16%    |
| **100KSmall**| Brotli    | 7933926 bytes | 3908434 bytes   | 50.74%    |
|              | Gzip      | 7934106 bytes | 4261133 bytes   | 46.29%    |
|              | LZ4       | 7934005 bytes | 5984094 bytes   | 24.58%    |

# Compression Results (Repeated Data)
| Size       | Algorithm | Encoded Size | Compressed Size | Gain (%)  |
|------------|-----------|--------------|-----------------|-----------|
| **2KB**    | Brotli    | 2048 bytes   | 24 bytes        | 98.83%    |
|            | Gzip      | 2048 bytes   | 56 bytes        | 97.27%    |
|            | LZ4       | 2048 bytes   | 56 bytes        | 97.27%    |
| **4KB**    | Brotli    | 4096 bytes   | 25 bytes        | 99.39%    |
|            | Gzip      | 4096 bytes   | 63 bytes        | 98.46%    |
|            | LZ4       | 4096 bytes   | 64 bytes        | 98.44%    |
| **8KB**    | Brotli    | 8192 bytes   | 25 bytes        | 99.69%    |
|            | Gzip      | 8192 bytes   | 71 bytes        | 99.13%    |
|            | LZ4       | 8192 bytes   | 80 bytes        | 99.02%    |
| **16KB**   | Brotli    | 16384 bytes  | 25 bytes        | 99.85%    |
|            | Gzip      | 16384 bytes  | 87 bytes        | 99.47%    |
|            | LZ4       | 16384 bytes  | 113 bytes       | 99.31%    |
| **32KB**   | Brotli    | 32768 bytes  | 25 bytes        | 99.92%    |
|            | Gzip      | 32768 bytes  | 117 bytes       | 99.64%    |
|            | LZ4       | 32768 bytes  | 177 bytes       | 99.46%    |
| **64KB**   | Brotli    | 65536 bytes  | 25 bytes        | 99.96%    |
|            | Gzip      | 65536 bytes  | 181 bytes       | 99.72%    |
|            | LZ4       | 65536 bytes  | 305 bytes       | 99.53%    |
| **128KB**  | Brotli    | 131072 bytes | 26 bytes        | 99.98%    |
|            | Gzip      | 131072 bytes | 308 bytes       | 99.77%    |
|            | LZ4       | 131072 bytes | 562 bytes       | 99.57%    |
| **256KB**  | Brotli    | 262144 bytes | 26 bytes        | 99.99%    |
|            | Gzip      | 262144 bytes | 560 bytes       | 99.79%    |
|            | LZ4       | 262144 bytes | 1076 bytes      | 99.59%    |
| **512KB**  | Brotli    | 524288 bytes | 26 bytes        | 100.00%   |
|            | Gzip      | 524288 bytes | 1071 bytes      | 99.80%    |
|            | LZ4       | 524288 bytes | 2104 bytes      | 99.60%    |
| **1MB**    | Brotli    | 1048576 bytes| 26 bytes        | 100.00%   |
|            | Gzip      | 1048576 bytes| 2087 bytes      | 99.80%    |
|            | LZ4       | 1048576 bytes| 4160 bytes      | 99.60%    |
| **10MB**   | Brotli    | 10485760 bytes| 39 bytes       | 100.00%   |
|            | Gzip      | 10485760 bytes| 20403 bytes    | 99.81%    |
|            | LZ4       | 10485760 bytes| 41236 bytes    | 99.61%    |



