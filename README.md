# Benchmarks

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Go benchmarks comparing execution time, throughput, memory use, and allocations
across libraries and common application workloads.

Requires **Go 1.27.1 or later**. See each category for results and methodology;
requirements and measurement environments vary by suite.

## Categories

| Category | Benchmarks |
| --- | --- |
| [Cache](cache) | Redis and Memcached |
| [Collectors](collectors) | slogx-collector, Vector, and Fluent Bit delivery and outage recovery |
| [Compression](compression) | Brotli, Gzip, and LZ4 |
| [Cryptography](crypto) | JWT signing and verification; AES encryption and decryption |
| [Databases](db) | SQLite, MariaDB, and PostgreSQL |
| [Encoding](encoding) | Base64 and Base85 |
| [Iteration](iteration) | Slice and channel iterator adapters |
| [Logging](log) | Structured logging, context enrichment, spans, and errors |
| [Meta](meta) | Reflection, sorting, and synchronization |
| [Protocols](protocol) | Local HTTP and TCP transport |
| [Serialization](serializer) | JSON, Gob, and Protocol Buffers |
| [Validation](validator) | Ward, go-playground/validator, and ozzo-validation |

## Support

If these benchmarks are useful to you, consider supporting the project.

[![Buy Me A Coffee](https://cdn.buymeacoffee.com/buttons/default-orange.png)](https://www.buymeacoffee.com/rah.0)
