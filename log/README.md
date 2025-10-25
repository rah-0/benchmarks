# Go Logging Libraries Performance Benchmark

Comprehensive performance comparison of Go logging libraries with structured JSON output.

## Libraries Tested

- **[nabu](https://github.com/rah-0/nabu)** v1.0.0 - Structured logging with error chain tracking and UUID correlation
- **[zerolog](https://github.com/rs/zerolog)** v1.34.0 - Zero allocation JSON logger
- **[zap](https://github.com/uber-go/zap)** v1.27.0 - Uber's blazing fast, structured logger
- **[logrus](https://github.com/sirupsen/logrus)** v1.9.3 - Structured logger with hooks
- **[slog](https://pkg.go.dev/log/slog)** (stdlib) - Go 1.21+ standard library structured logging

## Benchmark Results

### 1. Simple Message Logging

Basic log message without any fields.

| Library  | ns/op   | B/op | allocs/op | Relative Speed |
|----------|---------|------|-----------|----------------|
| zerolog  | 135     | 0    | 0         | **1.00x** ⚡   |
| zap      | 276     | 0    | 0         | 2.04x          |
| slog     | 405     | 0    | 0         | 3.00x          |
| logrus   | 1,318   | 873  | 19        | 9.76x          |
| nabu     | 1,817   | 576  | 8         | 13.46x         |

### 2. Message with Fields

Log message with 3 structured fields (userID, action, ip).

| Library  | ns/op   | B/op  | allocs/op | Relative Speed |
|----------|---------|-------|-----------|----------------|
| zerolog  | 182     | 0     | 0         | **1.00x** ⚡   |
| zap      | 493     | 192   | 1         | 2.71x          |
| slog     | 672     | 0     | 0         | 3.69x          |
| nabu     | 2,263   | 824   | 10        | 12.43x         |
| logrus   | 2,454   | 1,810 | 29        | 13.48x         |

### 3. Error Logging

Logging an error with a message.

| Library  | ns/op   | B/op  | allocs/op | Relative Speed |
|----------|---------|-------|-----------|----------------|
| zerolog  | 155     | 0     | 0         | **1.00x** ⚡   |
| zap      | 379     | 64    | 1         | 2.45x          |
| slog     | 564     | 0     | 0         | 3.64x          |
| logrus   | 2,018   | 1,685 | 26        | 13.02x         |
| nabu     | 2,693   | 1,032 | 12        | 17.37x         |

### 4. Error with Fields

Error logging with 3 additional structured fields.

| Library  | ns/op   | B/op  | allocs/op | Relative Speed |
|----------|---------|-------|-----------|----------------|
| zerolog  | 196     | 0     | 0         | **1.00x** ⚡   |
| zap      | 548     | 256   | 1         | 2.80x          |
| slog     | 765     | 0     | 0         | 3.90x          |
| logrus   | 3,006   | 2,374 | 35        | 15.34x         |
| nabu     | 3,106   | 1,345 | 14        | 15.85x         |

### 5. Error Chain (3 levels)

Logging a chain of 3 related errors with correlation.
- **nabu**: Automatic UUID correlation (built-in)
- **others**: Manual trace_id field for correlation

| Library  | ns/op   | B/op  | allocs/op | Relative Speed |
|----------|---------|-------|-----------|----------------|
| zerolog  | 523     | 0     | 0         | **1.00x** ⚡   |
| zap      | 1,245   | 272   | 4         | 2.38x          |
| slog     | 1,925   | 64    | 4         | 3.68x          |
| logrus   | 6,857   | 5,616 | 85        | 13.11x         |
| nabu     | 8,155   | 2,938 | 32        | 15.59x         |

## Visual Performance Comparison

### Speed Comparison (Lower is Better)

```
Simple Message:
zerolog  ████ 135 ns/op
zap      ████████ 276 ns/op
slog     ████████████ 405 ns/op
logrus   ███████████████████████████████████████ 1,318 ns/op
nabu     ████████████████████████████████████████████████████ 1,817 ns/op

With Fields:
zerolog  █ 182 ns/op
zap      ███ 493 ns/op
slog     ████ 672 ns/op
nabu     ████████████ 2,263 ns/op
logrus   █████████████ 2,454 ns/op

Error Logging:
zerolog  █ 155 ns/op
zap      ██ 379 ns/op
slog     ███ 564 ns/op
logrus   █████████████ 2,018 ns/op
nabu     █████████████████ 2,693 ns/op

Error Chain (3 levels with correlation):
zerolog  ██ 523 ns/op
zap      █████ 1,245 ns/op
slog     ███████ 1,925 ns/op
logrus   ██████████████████████████ 6,857 ns/op
nabu     ███████████████████████████████ 8,155 ns/op
```

### Memory Allocations (Lower is Better)

```
Simple Message:
zerolog  ░ 0 B/op, 0 allocs
zap      ░ 0 B/op, 0 allocs
slog     ░ 0 B/op, 0 allocs
nabu     ███ 576 B/op, 8 allocs
logrus   █████ 873 B/op, 19 allocs

Error Chain (with correlation):
zerolog  ░ 0 B/op, 0 allocs
slog     █ 64 B/op, 4 allocs
zap      ██ 272 B/op, 4 allocs
nabu     ████████ 2,938 B/op, 32 allocs
logrus   ███████████████ 5,616 B/op, 85 allocs
```
