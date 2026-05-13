# Validator Benchmarks

Comparison of Go validation libraries across common real-world scenarios.

## Libraries

| Library | Version | Approach |
|---|---|---|
| [ward](https://github.com/rah-0/ward) | v0.0.5 | Typed, reflection-free, explicit rule functions |
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

## Initialization and reuse model

Each library is benchmarked in its own recommended production pattern:

| Library | Model | Rationale |
|---|---|---|
| ward | New instance per request | Not safe to share across goroutines — `Field` holds a `*T` to the source value |
| go-playground/validator | Singleton, reused across requests | Thread-safe by design; `validator.New()` builds reflection caches that are meant to be shared |
| ozzo-validation | Stateless — no instance | Purely functional API, nothing to initialise or reuse |

## Raw Results

```
goos: linux
goarch: amd64
cpu: AMD Ryzen 9 5950X 16-Core Processor

BenchmarkOzzoSingleFieldValid-8                       	 1656590	       719.9 ns/op	      64 B/op	       3 allocs/op
BenchmarkOzzoSingleFieldInvalid-8                     	 9663243	       123.4 ns/op	      64 B/op	       3 allocs/op
BenchmarkOzzoMultiFieldAllValid-8                     	  124765	      9098 ns/op	    1666 B/op	      24 allocs/op
BenchmarkOzzoMultiFieldSomeInvalid-8                  	  665569	      1759 ns/op	    1616 B/op	      24 allocs/op
BenchmarkOzzoFullCycleSingleFieldValid-8              	 1654305	       727.5 ns/op	      64 B/op	       3 allocs/op
BenchmarkOzzoFullCycleSingleFieldInvalid-8            	 9029662	       130.0 ns/op	      64 B/op	       3 allocs/op
BenchmarkOzzoFullCycleMultiFieldAllValid-8            	  128338	      9110 ns/op	    1672 B/op	      24 allocs/op
BenchmarkOzzoFullCycleMultiFieldSomeInvalid-8         	  106894	     11143 ns/op	    8985 B/op	     118 allocs/op
BenchmarkPlaygroundSingleFieldValid-8                 	 1460930	       814.4 ns/op	     105 B/op	       6 allocs/op
BenchmarkPlaygroundSingleFieldInvalid-8               	 1882658	       625.9 ns/op	     352 B/op	      13 allocs/op
BenchmarkPlaygroundMultiFieldAllValid-8               	  622957	      1956 ns/op	     527 B/op	      11 allocs/op
BenchmarkPlaygroundMultiFieldSomeInvalid-8            	  626252	      1643 ns/op	    1153 B/op	      22 allocs/op
BenchmarkPlaygroundFullCycleSingleFieldValid-8        	 1466439	       819.7 ns/op	     105 B/op	       6 allocs/op
BenchmarkPlaygroundFullCycleSingleFieldInvalid-8      	 1897406	       636.9 ns/op	     352 B/op	      13 allocs/op
BenchmarkPlaygroundFullCycleMultiFieldAllValid-8      	  601244	      2053 ns/op	     528 B/op	      11 allocs/op
BenchmarkPlaygroundFullCycleMultiFieldSomeInvalid-8   	  688426	      1704 ns/op	    1153 B/op	      22 allocs/op
BenchmarkWardSingleFieldValid-8                       	 3365926	       351.0 ns/op	     216 B/op	       9 allocs/op
BenchmarkWardSingleFieldInvalid-8                     	 2148957	       555.2 ns/op	     352 B/op	      15 allocs/op
BenchmarkWardMultiFieldAllValid-8                     	  645504	      1737 ns/op	    1014 B/op	      22 allocs/op
BenchmarkWardMultiFieldSomeInvalid-8                  	  546021	      2154 ns/op	    1779 B/op	      42 allocs/op
BenchmarkWardMultiFieldStopOnFail-8                   	 1000000	      1019 ns/op	     992 B/op	      27 allocs/op
BenchmarkWardFullCycleSingleFieldValid-8              	 3266635	       344.4 ns/op	     216 B/op	       9 allocs/op
BenchmarkWardFullCycleSingleFieldInvalid-8            	 2133907	       560.0 ns/op	     352 B/op	      15 allocs/op
BenchmarkWardFullCycleMultiFieldAllValid-8            	  660826	      1726 ns/op	    1015 B/op	      22 allocs/op
BenchmarkWardFullCycleMultiFieldSomeInvalid-8         	  563469	      2172 ns/op	    1780 B/op	      42 allocs/op
```

## Speed (ns/op, lower is better)

### Single Field Valid

| Library | ns/op | Relative |
|---|---|---|
| ward | 351 | **1.00x** ⚡ |
| playground | 814 | 2.32x |
| ozzo | 720 | 2.05x |

### Single Field Invalid

| Library | ns/op | Relative |
|---|---|---|
| ozzo | 123 | **1.00x** ⚡ |
| playground | 626 | 5.09x |
| ward | 555 | 4.51x |

### Multi Field All Valid

| Library | ns/op | Relative |
|---|---|---|
| ward | 1,737 | **1.00x** ⚡ |
| playground | 1,956 | 1.13x |
| ozzo | 9,098 | 5.24x |

### Multi Field Some Invalid

| Library | ns/op | Relative |
|---|---|---|
| playground | 1,643 | **1.00x** ⚡ |
| ward | 2,154 | 1.31x |
| ozzo | 1,759 | 1.07x |

## Full Cycle Speed (validate + inspect results, ns/op, lower is better)

### Single Field Valid

| Library | ns/op | Relative |
|---|---|---|
| ward | 344 | **1.00x** ⚡ |
| ozzo | 728 | 2.12x |
| playground | 820 | 2.38x |

### Single Field Invalid

| Library | ns/op | Relative |
|---|---|---|
| ozzo | 130 | **1.00x** ⚡ |
| playground | 637 | 4.90x |
| ward | 560 | 4.31x |

### Multi Field All Valid

| Library | ns/op | Relative |
|---|---|---|
| ward | 1,726 | **1.00x** ⚡ |
| playground | 2,053 | 1.19x |
| ozzo | 9,110 | 5.28x |

### Multi Field Some Invalid

| Library | ns/op | B/op | allocs/op | Relative |
|---|---|---|---|---|
| playground | 1,704 | 1,153 | 22 | **1.00x** ⚡ |
| ward | 2,172 | 1,780 | 42 | 1.27x |
| ozzo | 11,143 | 8,985 | 118 | 6.54x |

## Memory (B/op, lower is better)

### Single Field Valid

| Library | B/op | allocs/op |
|---|---|---|
| ozzo | 64 | 3 |
| playground | 105 | 6 |
| ward | 216 | 9 |

### Multi Field All Valid

| Library | B/op | allocs/op |
|---|---|---|
| playground | 527 | 11 |
| ozzo | 1,666 | 24 |
| ward | 1,014 | 22 |

## Visual Comparison

### Speed — Single Field (ns/op, lower is better)

```
Valid:
ward        ████████████████ 351 ns/op
ozzo        ████████████████████████████████ 720 ns/op
playground  ████████████████████████████████████ 814 ns/op

Invalid:
ozzo        █████ 123 ns/op
ward        ████████████████████ 555 ns/op
playground  ███████████████████████ 626 ns/op
```

### Speed — Multi Field (ns/op, lower is better)

```
All Valid:
ward        ████████ 1,737 ns/op
playground  █████████ 1,956 ns/op
ozzo        ████████████████████████████████████████ 9,098 ns/op

Some Invalid:
playground  ████████ 1,643 ns/op
ozzo        █████████ 1,759 ns/op
ward        ██████████ 2,154 ns/op
```

### Memory — Multi Field All Valid (B/op, lower is better)

```
playground  ████████████ 527 B/op, 11 allocs
ward        ████████████████████████ 1,014 B/op, 22 allocs
ozzo        ████████████████████████████████████████ 1,666 B/op, 24 allocs
```

### Full Cycle Speed — Single Field (ns/op, lower is better)

```
Valid:
ward        ████████████████ 344 ns/op
ozzo        ████████████████████████████████ 728 ns/op
playground  ████████████████████████████████████ 820 ns/op

Invalid:
ozzo        █████ 130 ns/op
ward        ████████████████████ 560 ns/op
playground  ███████████████████████ 637 ns/op
```

### Full Cycle Speed — Multi Field (ns/op, lower is better)

```
All Valid:
ward        ████████ 1,726 ns/op
playground  █████████ 2,053 ns/op
ozzo        ████████████████████████████████████████ 9,110 ns/op

Some Invalid:
playground  ████████ 1,704 ns/op
ward        █████████ 2,172 ns/op
ozzo        ████████████████████████████████████████ 11,143 ns/op
```

### Full Cycle Memory — Multi Field Some Invalid (B/op, lower is better)

```
playground  █████ 1,153 B/op, 22 allocs
ward        ████████ 1,780 B/op, 42 allocs
ozzo        ████████████████████████████████████████ 8,985 B/op, 118 allocs
```

## Notes

**Happy path (all valid)** is where ward wins most decisively — only failing rules allocate a `Result`, so the valid path costs nearly nothing per rule regardless of how many rules run. Ward is **2x faster** than ozzo and playground on single-field valid, and comparable to playground on multi-field.

**Failure path some-invalid:** playground and ozzo edge ward on multi-field some-invalid because ward constructs field structs per request. However, ozzo's full-cycle some-invalid cost (11,143 ns, 8,985 B/op, 118 allocs) is dramatically higher when results are inspected — its `validation.Errors` map allocates heavily per entry. Playground and ward remain flat on result inspection since their failure structures are plain slices.

**StopOnFail** (ward only) halts at the first failing field, bringing multi-field invalid cost down to single-field territory (1,019 ns vs 2,154 ns).

**Ozzo multi-field all-valid** is the weakest scenario for ozzo — it always evaluates all fields eagerly (8,928 ns vs ward's 1,737 ns).
