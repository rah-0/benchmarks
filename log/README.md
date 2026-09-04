# Go Logging Libraries Performance Benchmark

Comprehensive performance comparison of Go logging libraries with structured JSON output.

## Libraries Tested

- **[zerolog](https://github.com/rs/zerolog)** v1.34.0 - Zero allocation JSON logger
- **[zap](https://github.com/uber-go/zap)** v1.27.0 - Uber's blazing fast, structured logger
- **[logrus](https://github.com/sirupsen/logrus)** v1.9.3 - Structured logger with hooks
- **[slog](https://pkg.go.dev/log/slog)** (stdlib) - Go 1.21+ standard library structured logging

## Benchmark Results

### 1. Simple Message Logging

Basic log message without any fields.

| Library  | ns/op   | B/op | allocs/op | Relative Speed |
|----------|---------|------|-----------|----------------|
| zerolog  | 132     | 0    | 0         | **1.00x** ⚡   |
| zap      | 256     | 0    | 0         | 1.94x          |
| slog     | 406     | 0    | 0         | 3.07x          |
| logrus   | 1,424   | 873  | 19        | 10.78x         |

### 2. Message with Fields

Log message with 3 structured fields (userID, action, ip).

| Library  | ns/op   | B/op  | allocs/op | Relative Speed |
|----------|---------|-------|-----------|----------------|
| zerolog  | 190     | 0     | 0         | **1.00x** ⚡   |
| zap      | 494     | 192   | 1         | 2.60x          |
| slog     | 668     | 0     | 0         | 3.51x          |
| logrus   | 2,689   | 1,812 | 29        | 14.15x         |

### 3. Error Logging

Logging an error with a message.

| Library  | ns/op   | B/op  | allocs/op | Relative Speed |
|----------|---------|-------|-----------|----------------|
| zerolog  | 167     | 0     | 0         | **1.00x** ⚡   |
| zap      | 392     | 64    | 1         | 2.35x          |
| slog     | 557     | 0     | 0         | 3.34x          |
| logrus   | 2,237   | 1,686 | 26        | 13.40x         |

### 4. Error with Fields

Error logging with 3 additional structured fields.

| Library  | ns/op   | B/op  | allocs/op | Relative Speed |
|----------|---------|-------|-----------|----------------|
| zerolog  | 212     | 0     | 0         | **1.00x** ⚡   |
| zap      | 567     | 256   | 1         | 2.67x          |
| slog     | 794     | 0     | 0         | 3.74x          |
| logrus   | 3,424   | 2,376 | 35        | 16.15x         |

### 5. Error Chain (3 levels)

Logging three related errors with a shared manual `trace_id` field.

| Library  | ns/op   | B/op  | allocs/op | Relative Speed |
|----------|---------|-------|-----------|----------------|
| zerolog  | 526     | 0     | 0         | **1.00x** ⚡   |
| zap      | 1,306   | 272   | 4         | 2.48x          |
| slog     | 1,936   | 64    | 4         | 3.68x          |
| logrus   | 7,557   | 5,621 | 85        | 14.37x         |

## Visual Performance Comparison

### Speed Comparison (Lower is Better)

```
Simple Message:
zerolog  ████ 132 ns/op
zap      ███████ 256 ns/op
slog     ███████████ 406 ns/op
logrus   ██████████████████████████████████████ 1,424 ns/op

With Fields:
zerolog  █ 190 ns/op
zap      ███ 494 ns/op
slog     ████ 668 ns/op
logrus   ██████████████ 2,689 ns/op

Error Logging:
zerolog  █ 167 ns/op
zap      ██ 392 ns/op
slog     ███ 557 ns/op
logrus   █████████████ 2,237 ns/op

Error Chain (3 levels with correlation):
zerolog  ██ 526 ns/op
zap      █████ 1,306 ns/op
slog     ███████ 1,936 ns/op
logrus   ██████████████████████████████ 7,557 ns/op
```

### Memory Allocations (Lower is Better)

```
Simple Message:
zerolog  ░ 0 B/op, 0 allocs
zap      ░ 0 B/op, 0 allocs
slog     ░ 0 B/op, 0 allocs
logrus   █████ 873 B/op, 19 allocs

Error Chain (with correlation):
zerolog  ░ 0 B/op, 0 allocs
slog     █ 64 B/op, 4 allocs
zap      ██ 272 B/op, 4 allocs
logrus   ███████████████ 5,621 B/op, 85 allocs
```
