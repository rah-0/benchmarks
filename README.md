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
- [Compression](https://github.com/rah-0/benchmarks/tree/master/compression)
- [DB](https://github.com/rah-0/benchmarks/tree/master/db)
- [Protocol](https://github.com/rah-0/benchmarks/tree/master/protocol)
- [Serializer](https://github.com/rah-0/benchmarks/tree/master/serializer)
