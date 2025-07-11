# Benchmark Results

## Sign/Verify vs Encrypt/Decrypt

| Operation   | ns/op     | B/op     | allocs/op | CPU time   | Memory Use   | 
|-------------|-----------|----------|-----------|------------|--------------|
| JWT Sign    | 2,994 ns  | 2,457 B  | 38        | 2.99 µs    | ~2.4 KB      |
| JWT Verify  | 4,126 ns  | 2,536 B  | 50        | 4.13 µs    | ~2.5 KB      |
| AES Encrypt | 1,541 ns  | 1,472 B  | 6         | 1.54 µs    | ~1.4 KB      |
| AES Decrypt | 1,669 ns  | 1,600 B  | 9         | 1.67 µs    | ~1.6 KB      |
