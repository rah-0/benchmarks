# Validator Benchmarks

Comparison of Go validation libraries across common real-world scenarios.

## Libraries

| Library | Version | Approach |
|---|---|---|
| [ward](https://github.com/rah-0/ward) | v0.0.4 | Typed, reflection-free, explicit rule functions |
| [go-playground/validator](https://github.com/go-playground/validator) | v10 | Struct tags + reflection |
| [ozzo-validation](https://github.com/go-ozzo/ozzo-validation) | v4 | Fluent rules, reflection-based |

## Input

All libraries validate the same `UserForm` struct with equivalent rules:

```go
type UserForm struct {
    Email    string // required, valid email
    Username string // required, 5–50 chars, valid username chars
    Password string // required, 8–150 chars, lower+upper+digit+special
    Age      string // required, non-negative integer
    Website  string // required, valid URL
}
```

## Scenarios

- **SingleFieldValid** — one field (`Email`), all rules pass
- **SingleFieldInvalid** — one field (`Email`), rules fail
- **MultiFieldAllValid** — all 5 fields, all rules pass (happy path)
- **MultiFieldSomeInvalid** — all 5 fields, 3 fields fail, all rules run
- **MultiFieldStopOnFail** — all 5 fields, 3 fail, stop on first failing field (ward only)

## Raw Results

```
goos: linux
goarch: amd64
cpu: AMD Ryzen 9 5950X 16-Core Processor

BenchmarkOzzoSingleFieldValid-8              	16693197	       724.4 ns/op	      64 B/op	       3 allocs/op
BenchmarkOzzoSingleFieldInvalid-8            	100000000	       122.2 ns/op	      64 B/op	       3 allocs/op
BenchmarkOzzoMultiFieldAllValid-8            	 1340925	      8922 ns/op	    1667 B/op	      24 allocs/op
BenchmarkOzzoMultiFieldSomeInvalid-8         	 7192389	      1669 ns/op	    1616 B/op	      24 allocs/op
BenchmarkPlaygroundSingleFieldValid-8        	14818170	       807.7 ns/op	     104 B/op	       6 allocs/op
BenchmarkPlaygroundSingleFieldInvalid-8      	18958045	       616.3 ns/op	     352 B/op	      13 allocs/op
BenchmarkPlaygroundMultiFieldAllValid-8      	 6250140	      1913 ns/op	     527 B/op	      11 allocs/op
BenchmarkPlaygroundMultiFieldSomeInvalid-8   	 7486387	      1602 ns/op	    1153 B/op	      22 allocs/op
BenchmarkWardSingleFieldValid-8              	51375868	       234.8 ns/op	      88 B/op	       5 allocs/op
BenchmarkWardSingleFieldInvalid-8            	29662713	       401.0 ns/op	     216 B/op	      10 allocs/op
BenchmarkWardMultiFieldAllValid-8            	10639156	      1111 ns/op	     234 B/op	       6 allocs/op
BenchmarkWardMultiFieldSomeInvalid-8         	 8997855	      1333 ns/op	     877 B/op	      22 allocs/op
BenchmarkWardMultiFieldStopOnFail-8          	29967513	       402.7 ns/op	     216 B/op	      10 allocs/op
```

## Speed (ns/op, lower is better)

### Single Field Valid

| Library | ns/op | Relative |
|---|---|---|
| ward | 234 | **1.00x** ⚡ |
| ozzo | 724 | 3.09x |
| playground | 807 | 3.44x |

### Single Field Invalid

| Library | ns/op | Relative |
|---|---|---|
| ozzo | 122 | **1.00x** ⚡ |
| ward | 401 | 3.29x |
| playground | 616 | 5.05x |

### Multi Field All Valid

| Library | ns/op | Relative |
|---|---|---|
| ward | 1,111 | **1.00x** ⚡ |
| playground | 1,913 | 1.72x |
| ozzo | 8,922 | 8.03x |

### Multi Field Some Invalid

| Library | ns/op | Relative |
|---|---|---|
| ward | 1,333 | **1.00x** ⚡ |
| playground | 1,602 | 1.20x |
| ozzo | 1,669 | 1.25x |

## Memory (B/op, lower is better)

### Single Field Valid

| Library | B/op | allocs/op |
|---|---|---|
| ozzo | 64 | 3 |
| ward | 88 | 5 |
| playground | 104 | 6 |

### Multi Field All Valid

| Library | B/op | allocs/op |
|---|---|---|
| ward | 234 | 6 |
| playground | 527 | 11 |
| ozzo | 1,667 | 24 |

## Visual Comparison

### Speed — Single Field (ns/op, lower is better)

```
Valid:
ward        ████████████ 234 ns/op
ozzo        ████████████████████████████████████ 724 ns/op
playground  ████████████████████████████████████████ 807 ns/op

Invalid:
ozzo        █████ 122 ns/op
ward        ████████████████ 401 ns/op
playground  █████████████████████████ 616 ns/op
```

### Speed — Multi Field (ns/op, lower is better)

```
All Valid:
ward        █████ 1,111 ns/op
playground  █████████ 1,913 ns/op
ozzo        ████████████████████████████████████████ 8,922 ns/op

Some Invalid:
ward        ████████████████████████████████ 1,333 ns/op
playground  ██████████████████████████████████████ 1,602 ns/op
ozzo        ████████████████████████████████████████ 1,669 ns/op
```

### Memory — Multi Field All Valid (B/op, lower is better)

```
ward        ██████ 234 B/op, 6 allocs
playground  █████████████ 527 B/op, 11 allocs
ozzo        ████████████████████████████████████████ 1,667 B/op, 24 allocs
```

## Notes

**Happy path (all valid)** is where ward wins most decisively — only failing rules allocate a `Result`, so the valid path costs nearly nothing per rule regardless of how many rules run.

**Failure path** allocates one `Result` per failing rule by design — failures carry field name, rule ID and args back to the caller.

**StopOnFail** (ward only) halts validation at the first failing field, cutting the failure path cost to match the single-field case.
