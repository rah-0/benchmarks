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

BenchmarkOzzoSingleFieldValid-8                       	16983620	       709.2 ns/op	      64 B/op	       3 allocs/op
BenchmarkOzzoSingleFieldInvalid-8                     	98540638	       119.9 ns/op	      64 B/op	       3 allocs/op
BenchmarkOzzoMultiFieldAllValid-8                     	 1357358	      8828 ns/op	    1669 B/op	      24 allocs/op
BenchmarkOzzoMultiFieldSomeInvalid-8                  	 7121955	      1707 ns/op	    1616 B/op	      24 allocs/op
BenchmarkOzzoFullCycleSingleFieldValid-8              	17016690	       709.5 ns/op	      64 B/op	       3 allocs/op
BenchmarkOzzoFullCycleSingleFieldInvalid-8            	92854027	       129.2 ns/op	      64 B/op	       3 allocs/op
BenchmarkOzzoFullCycleMultiFieldAllValid-8            	 1330212	      8871 ns/op	    1667 B/op	      24 allocs/op
BenchmarkOzzoFullCycleMultiFieldSomeInvalid-8         	 1000000	     10587 ns/op	    8985 B/op	     118 allocs/op
BenchmarkPlaygroundSingleFieldValid-8                 	15043352	       801.0 ns/op	     105 B/op	       6 allocs/op
BenchmarkPlaygroundSingleFieldInvalid-8               	19458655	       616.7 ns/op	     352 B/op	      13 allocs/op
BenchmarkPlaygroundMultiFieldAllValid-8               	 6310456	      1922 ns/op	     527 B/op	      11 allocs/op
BenchmarkPlaygroundMultiFieldSomeInvalid-8            	 7218932	      1625 ns/op	    1153 B/op	      22 allocs/op
BenchmarkPlaygroundFullCycleSingleFieldValid-8        	14817615	       806.5 ns/op	     105 B/op	       6 allocs/op
BenchmarkPlaygroundFullCycleSingleFieldInvalid-8      	19336328	       626.6 ns/op	     352 B/op	      13 allocs/op
BenchmarkPlaygroundFullCycleMultiFieldAllValid-8      	 6217468	      1900 ns/op	     527 B/op	      11 allocs/op
BenchmarkPlaygroundFullCycleMultiFieldSomeInvalid-8   	 7228701	      1655 ns/op	    1153 B/op	      22 allocs/op
BenchmarkWardSingleFieldValid-8                       	52729874	       227.3 ns/op	      88 B/op	       5 allocs/op
BenchmarkWardSingleFieldInvalid-8                     	30465769	       396.6 ns/op	     216 B/op	      10 allocs/op
BenchmarkWardMultiFieldAllValid-8                     	10905783	      1103 ns/op	     234 B/op	       6 allocs/op
BenchmarkWardMultiFieldSomeInvalid-8                  	 8913338	      1334 ns/op	     877 B/op	      22 allocs/op
BenchmarkWardMultiFieldStopOnFail-8                   	29362284	       408.5 ns/op	     216 B/op	      10 allocs/op
BenchmarkWardFullCycleSingleFieldValid-8              	52909069	       225.4 ns/op	      88 B/op	       5 allocs/op
BenchmarkWardFullCycleSingleFieldInvalid-8            	29764695	       399.9 ns/op	     216 B/op	      10 allocs/op
BenchmarkWardFullCycleMultiFieldAllValid-8            	10852663	      1104 ns/op	     235 B/op	       6 allocs/op
BenchmarkWardFullCycleMultiFieldSomeInvalid-8         	 8940127	      1354 ns/op	     877 B/op	      22 allocs/op
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
| ward | 1,334 | **1.00x** ⚡ |
| playground | 1,625 | 1.22x |
| ozzo | 1,707 | 1.28x |

## Full Cycle Speed (validate + inspect results, ns/op, lower is better)

### Single Field Valid

| Library | ns/op | Relative |
|---|---|---|
| ward | 225 | **1.00x** ⚡ |
| ozzo | 709 | 3.15x |
| playground | 806 | 3.58x |

### Single Field Invalid

| Library | ns/op | Relative |
|---|---|---|
| ozzo | 129 | **1.00x** ⚡ |
| ward | 399 | 3.09x |
| playground | 626 | 4.85x |

### Multi Field All Valid

| Library | ns/op | Relative |
|---|---|---|
| ward | 1,104 | **1.00x** ⚡ |
| playground | 1,900 | 1.72x |
| ozzo | 8,871 | 8.03x |

### Multi Field Some Invalid

| Library | ns/op | B/op | allocs/op | Relative |
|---|---|---|---|---|
| ward | 1,354 | 877 | 22 | **1.00x** ⚡ |
| playground | 1,655 | 1,153 | 22 | 1.22x |
| ozzo | 10,587 | 8,985 | 118 | **7.82x** |

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

### Full Cycle Speed — Single Field (ns/op, lower is better)

```
Valid:
ward        ███████████ 225 ns/op
ozzo        ███████████████████████████████████ 709 ns/op
playground  ████████████████████████████████████████ 806 ns/op

Invalid:
ozzo        █████ 129 ns/op
ward        ████████████████ 399 ns/op
playground  █████████████████████████ 626 ns/op
```

### Full Cycle Speed — Multi Field (ns/op, lower is better)

```
All Valid:
ward        █████ 1,104 ns/op
playground  █████████ 1,900 ns/op
ozzo        ████████████████████████████████████████ 8,871 ns/op

Some Invalid:
ward        █████ 1,354 ns/op
playground  ██████ 1,655 ns/op
ozzo        ████████████████████████████████████████ 10,587 ns/op
```

### Full Cycle Memory — Multi Field Some Invalid (B/op, lower is better)

```
ward        ████ 877 B/op, 22 allocs
playground  █████ 1,153 B/op, 22 allocs
ozzo        ████████████████████████████████████████ 8,985 B/op, 118 allocs
```

## Notes

**Happy path (all valid)** is where ward wins most decisively — only failing rules allocate a `Result`, so the valid path costs nearly nothing per rule regardless of how many rules run.

**Failure path** allocates one `Result` per failing rule by design — failures carry field name, rule ID and args back to the caller.

**StopOnFail** (ward only) halts validation at the first failing field, cutting the failure path cost to match the single-field case.

**Ozzo full cycle multi invalid** is the biggest surprise: iterating its `validation.Errors` map allocates heavily (8,985 B/op, 118 allocs) because the map itself and each error entry are heap-allocated. This is the realistic production cost — validation without inspecting the errors is not a real scenario. Ward and playground are unaffected by result inspection since their failure structures are flat slices.
