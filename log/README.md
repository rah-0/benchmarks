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
| zerolog  | 131     | 0    | 0         | **1.00x** ⚡   |
| zap      | 297     | 0    | 0         | 2.27x          |
| slog     | 415     | 0    | 0         | 3.17x          |
| logrus   | 1,402   | 873  | 19        | 10.72x         |
| nabu     | 1,606   | 320  | 5         | 12.28x         |

### 2. Message with Fields

Log message with 3 structured fields (userID, action, ip).

| Library  | ns/op   | B/op  | allocs/op | Relative Speed |
|----------|---------|-------|-----------|----------------|
| zerolog  | 177     | 0     | 0         | **1.00x** ⚡   |
| zap      | 517     | 192   | 1         | 2.92x          |
| slog     | 730     | 0     | 0         | 4.12x          |
| nabu     | 2,067   | 568   | 7         | 11.68x         |
| logrus   | 2,669   | 1,810 | 29        | 15.08x         |

### 3. Error Logging

Logging an error with a message.

| Library  | ns/op   | B/op  | allocs/op | Relative Speed |
|----------|---------|-------|-----------|----------------|
| zerolog  | 162     | 0     | 0         | **1.00x** ⚡   |
| zap      | 401     | 64    | 1         | 2.47x          |
| slog     | 578     | 0     | 0         | 3.57x          |
| logrus   | 2,268   | 1,685 | 26        | 14.00x         |
| nabu     | 2,906   | 1,032 | 12        | 17.94x         |

### 4. Error with Fields

Error logging with 3 additional structured fields.

| Library  | ns/op   | B/op  | allocs/op | Relative Speed |
|----------|---------|-------|-----------|----------------|
| zerolog  | 208     | 0     | 0         | **1.00x** ⚡   |
| zap      | 581     | 256   | 1         | 2.79x          |
| slog     | 861     | 0     | 0         | 4.14x          |
| nabu     | 3,353   | 1,345 | 14        | 16.12x         |
| logrus   | 3,401   | 2,374 | 35        | 16.35x         |

### 5. Error Chain (3 levels)

Logging a chain of 3 related errors.

| Library  | ns/op   | B/op  | allocs/op | Relative Speed |
|----------|---------|-------|-----------|----------------|
| zerolog  | 491     | 0     | 0         | **1.00x** ⚡   |
| zap      | 1,236   | 192   | 3         | 2.52x          |
| slog     | 1,755   | 0     | 0         | 3.57x          |
| logrus   | 6,864   | 5,056 | 78        | 13.98x         |
| nabu     | 9,020   | 2,938 | 32        | 18.37x         |

## Visual Performance Comparison

### Speed Comparison (Lower is Better)

```
Simple Message:
zerolog  ████ 131 ns/op
zap      █████████ 297 ns/op
slog     █████████████ 415 ns/op
logrus   ████████████████████████████████████████████ 1,402 ns/op
nabu     ██████████████████████████████████████████████████ 1,606 ns/op

With Fields:
zerolog  █ 177 ns/op
zap      ███ 517 ns/op
slog     ████ 730 ns/op
nabu     ████████████ 2,067 ns/op
logrus   ███████████████ 2,669 ns/op

Error Logging:
zerolog  █ 162 ns/op
zap      ██ 401 ns/op
slog     ███ 578 ns/op
logrus   ██████████████ 2,268 ns/op
nabu     ██████████████████ 2,906 ns/op

Error Chain (3 levels):
zerolog  ██ 491 ns/op
zap      █████ 1,236 ns/op
slog     ███████ 1,755 ns/op
logrus   ██████████████████████████ 6,864 ns/op
nabu     ███████████████████████████████████ 9,020 ns/op
```

### Memory Allocations (Lower is Better)

```
Simple Message:
zerolog  ░ 0 B/op, 0 allocs
zap      ░ 0 B/op, 0 allocs
slog     ░ 0 B/op, 0 allocs
nabu     ██ 320 B/op, 5 allocs
logrus   ██████ 873 B/op, 19 allocs

Error Chain:
zerolog  ░ 0 B/op, 0 allocs
zap      █ 192 B/op, 3 allocs
slog     ░ 0 B/op, 0 allocs
nabu     ████████ 2,938 B/op, 32 allocs
logrus   ██████████████ 5,056 B/op, 78 allocs
```
