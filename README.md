[![Go Report Card](https://goreportcard.com/badge/github.com/rah-0/benchmarks)](https://goreportcard.com/report/github.com/rah-0/benchmarks)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

<a href="https://www.buymeacoffee.com/rah.0" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/arial-orange.png" alt="Buy Me A Coffee" height="50"></a>

# Benchmarks
These benchmarks are executed on a Virtual Machine.  
The aim of benching is to see what's faster while using the least amount of memory possible and also fewer allocs.

## About the main system
**CPU**
- AMD Ryzen 9 5950X
- Precision Boost Overdrive (PBO) Enabled
- Maximum Clocks (all cores): 5GHz
- Voltage Average: 1.40V

**RAM**
- DDR4
- 64GB (4x16GB) 
- XMP Enabled
- System Memory Multiplier: 36
- FCLK Frequency: 1800MHz
- CL19
- Voltage: 1.40V

## About the VM
CPU: 8 threads  
RAM: 16GB  
OS: Debian 12, Kernel: 6.1.0-23-amd64  
NIC: 1GB


# Categories
These are the tested categories so far:
- [Cache](https://github.com/rah-0/benchmarks/tree/master/cache)
- [Compression](https://github.com/rah-0/benchmarks/tree/master/compression)
- [DB](https://github.com/rah-0/benchmarks/tree/master/db)
- [Encoding](https://github.com/rah-0/benchmarks/tree/master/encoding)
- [Meta](https://github.com/rah-0/benchmarks/tree/master/meta)
- [Protocol](https://github.com/rah-0/benchmarks/tree/master/protocol)
- [Serializer](https://github.com/rah-0/benchmarks/tree/master/serializer)

# ☕ Support
If **benchmarks** helped you optimize your Go projects, consider showing your appreciation. Your support fuels further benchmarking adventures!

[![Buy Me A Coffee](https://cdn.buymeacoffee.com/buttons/default-orange.png)](https://www.buymeacoffee.com/rah.0)~~
