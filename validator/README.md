# Validation benchmarks

Compare Ward v1.1.1, Playground v10.30.6, and Ozzo v4.5.0 on matching fixtures.
**Playground is the reference for each shared validation scenario.** Validation,
error processing, and stopping after the first failed field are separate workloads.

[Results](#results) · [Methodology](#methodology) · [Run](#running)

<!-- benchmark-results:start -->
## Results

**2026-11-21 · go1.27.4 · linux/amd64 · AMD Ryzen 9 5900HX with Radeon Graphics**
GOMAXPROCS: 16; samples per case: 3.

Tables show medians. **Speed = baseline time ÷ candidate time**: 2.00× is twice
the baseline speed; 0.50× is half. Allocated bytes and allocation counts describe
Go heap allocation per operation, not retained memory or process RSS.
These samples are descriptive; observed min–max ranges appear in the details
and are not confidence intervals. A zero denominator is shown as **—**.

[Environment](../results/environment.txt) · [Versions](../results/modules.txt) · [Reading the comparisons](../README.md#reading-the-comparisons)

<details>
<summary>Jump to a workload</summary>

- [One valid email](#one-valid-email)
- [One invalid email](#one-invalid-email)
- [Five valid fields](#five-valid-fields)
- [Five invalid fields](#five-invalid-fields)
- [One valid email: validate and inspect](#one-valid-email-validate-and-inspect)
- [One invalid email: validate and inspect](#one-invalid-email-validate-and-inspect)
- [Five valid fields: validate and inspect](#five-valid-fields-validate-and-inspect)
- [Five invalid fields: validate and inspect](#five-invalid-fields-validate-and-inspect)
- [Ward: stop at the first failure](#ward-stop-at-the-first-failure)

</details>

### One valid email

Validate one valid email.

Baseline: **Playground = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

Playground reuses its validator and caches; Ward constructs validators per operation; Ozzo uses its stateless API. Native email and URL checks can differ beyond these fixtures.

| Implementation | Time/validation, ns ↓ | Speed vs baseline ↑ | Allocated/validation, B ↓ | Allocs/validation ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Playground (baseline)** | 790.7 | 1.00× | 105 | 6 |
| Ward | 278.5 | 2.84× | 216 | 9 |
| Ozzo | 99.17 | 7.97× | 72 | 3 |

Ozzo has the lowest recorded median time: 7.97× the speed of Playground, with 0.69× its allocated bytes per operation.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Playground | 790.5–794.5 ns | 104–105 B | 6 | 1.00× |
| Ward | 277.2–279 ns | 216–216 B | 9 | 2.06× |
| Ozzo | 98.77–102.5 ns | 72–72 B | 3 | 0.69× |

</details>

### One invalid email

Validate one invalid email and retain the native failure result.

Baseline: **Playground = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

Playground reuses its validator and caches; Ward constructs validators per operation; Ozzo uses its stateless API. Native email and URL checks can differ beyond these fixtures.

| Implementation | Time/validation, ns ↓ | Speed vs baseline ↑ | Allocated/validation, B ↓ | Allocs/validation ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Playground (baseline)** | 587.9 | 1.00× | 368 | 14 |
| Ward | 456.4 | 1.29× | 368 | 16 |
| Ozzo | 97.73 | 6.02× | 72 | 3 |

Ozzo has the lowest recorded median time: 6.02× the speed of Playground, with 0.20× its allocated bytes per operation.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Playground | 587.1–588 ns | 368–368 B | 14 | 1.00× |
| Ward | 455.2–458.1 ns | 368–368 B | 16 | 1.00× |
| Ozzo | 97.61–98 ns | 72–72 B | 3 | 0.20× |

</details>

### Five valid fields

Validate five valid string fields with the declared matching fixture rules.

Baseline: **Playground = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

Playground reuses its validator and caches; Ward constructs validators per operation; Ozzo uses its stateless API. Native email and URL checks can differ beyond these fixtures.

| Implementation | Time/validation, µs ↓ | Speed vs baseline ↑ | Allocated/validation, B ↓ | Allocs/validation ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Playground (baseline)** | 1.98 | 1.00× | 316 | 7 |
| Ward | 1.51 | 1.31× | 1,025 | 22 |
| Ozzo | 8.31 | 0.24× | 1,712 | 25 |

Ward has the lowest recorded median time: 1.31× the speed of Playground, with 3.24× its allocated bytes per operation.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Playground | 1.98–1.98 µs | 316–316 B | 7 | 1.00× |
| Ward | 1.51–1.52 µs | 1–1 KiB | 22 | 3.24× |
| Ozzo | 8.28–8.36 µs | 1.66–1.67 KiB | 25 | 5.42× |

</details>

### Five invalid fields

Validate the invalid fixture through all five fields; every field fails.

Baseline: **Playground = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

Playground reuses its validator and caches; Ward constructs validators per operation; Ozzo uses its stateless API. Native email and URL checks can differ beyond these fixtures.

| Implementation | Time/validation, µs ↓ | Speed vs baseline ↑ | Allocated/validation, KiB ↓ | Allocs/validation ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Playground (baseline)** | 1.62 | 1.00× | 1.42 | 26 |
| Ward | 1.32 | 1.23× | 1.64 | 42 |
| Ozzo | 1.34 | 1.21× | 1.62 | 25 |

Ward has the lowest recorded median time: 1.23× the speed of Playground, with 1.15× its allocated bytes per operation.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Playground | 1.62–1.63 µs | 1.42–1.42 KiB | 26 | 1.00× |
| Ward | 1.32–1.33 µs | 1.64–1.64 KiB | 42 | 1.15× |
| Ozzo | 1.34–1.35 µs | 1.62–1.62 KiB | 25 | 1.14× |

</details>

### One valid email: validate and inspect

Validate one valid email. Check the success result.

Baseline: **Playground = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

All three cases validate and inspect success, with their declared initialization/reuse policies.

| Implementation | Time/validation, ns ↓ | Speed vs baseline ↑ | Allocated/validation, B ↓ | Allocs/validation ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Playground (baseline)** | 788.6 | 1.00× | 104 | 6 |
| Ward | 281.9 | 2.80× | 216 | 9 |
| Ozzo | 99.52 | 7.92× | 72 | 3 |

Ozzo has the lowest recorded median time: 7.92× the speed of Playground, with 0.69× its allocated bytes per operation.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Playground | 785.1–791.6 ns | 104–105 B | 6 | 1.00× |
| Ward | 281.2–282 ns | 216–216 B | 9 | 2.08× |
| Ozzo | 98.93–102 ns | 72–72 B | 3 | 0.69× |

</details>

### One invalid email: validate and inspect

Validate one invalid email and retain the native failure result. Inspect the library's native failure representation.

No equivalent reference measured; absolute costs are shown. Median of 3 samples; GOMAXPROCS=16.

Ozzo formats error text; Playground and Ward inspect field/rule identifiers. These consumer operations are not equivalent, so no speed ratio is assigned.

| Operation / variant | Time/validation, ns ↓ | Allocated/validation, B ↓ | Allocs/validation ↓ |
| --- | ---: | ---: | ---: |
| Playground | 590.4 | 368 | 14 |
| Ward | 460.5 | 368 | 16 |
| Ozzo | 98.96 | 72 | 3 |

<details>
<summary>Observed sample ranges</summary>

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| Playground | 589.3–595.8 ns | 368–368 B | 14 |
| Ward | 458.6–461.2 ns | 368–368 B | 16 |
| Ozzo | 98.71–103.4 ns | 72–72 B | 3 |

</details>

### Five valid fields: validate and inspect

Validate five valid string fields with the declared matching fixture rules. Check the success result.

Baseline: **Playground = 1.00×**. Median of 3 samples; GOMAXPROCS=16.

All three cases validate and inspect success, with their declared initialization/reuse policies.

| Implementation | Time/validation, µs ↓ | Speed vs baseline ↑ | Allocated/validation, B ↓ | Allocs/validation ↓ |
| --- | ---: | ---: | ---: | ---: |
| **Playground (baseline)** | 1.97 | 1.00× | 315 | 7 |
| Ward | 1.51 | 1.31× | 1,024 | 22 |
| Ozzo | 8.37 | 0.24× | 1,708 | 25 |

Ward has the lowest recorded median time: 1.31× the speed of Playground, with 3.25× its allocated bytes per operation.

<details>
<summary>Observed ranges and allocation comparison</summary>

| Implementation | Time min–max | Allocated min–max | Allocs min–max | Allocated vs baseline ↓ |
| --- | ---: | ---: | ---: | ---: |
| Playground | 1.97–1.98 µs | 314–315 B | 7 | 1.00× |
| Ward | 1.51–1.51 µs | 1,023–1,025 B | 22 | 3.25× |
| Ozzo | 8.28–8.39 µs | 1.66–1.67 KiB | 25 | 5.42× |

</details>

### Five invalid fields: validate and inspect

Validate the invalid fixture through all five fields; every field fails. Inspect the library's native failure representation.

No equivalent reference measured; absolute costs are shown. Median of 3 samples; GOMAXPROCS=16.

Ozzo formats error text; Playground and Ward inspect field/rule identifiers. These consumer operations are not equivalent, so no speed ratio is assigned.

| Operation / variant | Time/validation, µs ↓ | Allocated/validation, KiB ↓ | Allocs/validation ↓ |
| --- | ---: | ---: | ---: |
| Playground | 1.65 | 1.42 | 26 |
| Ward | 1.33 | 1.64 | 42 |
| Ozzo | 8.38 | 8.89 | 119 |

<details>
<summary>Observed sample ranges</summary>

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| Playground | 1.65–1.66 µs | 1.42–1.42 KiB | 26 |
| Ward | 1.32–1.33 µs | 1.64–1.64 KiB | 42 |
| Ozzo | 8.35–8.4 µs | 8.89–8.89 KiB | 119 |

</details>

### Ward: stop at the first failure

Validate the shared invalid fixture and stop after its first failing field, the email.

No equivalent reference measured; absolute costs are shown. Median of 3 samples; GOMAXPROCS=16.

This checks less of the form than the five-field failure comparison. No equivalent stop-on-first-field alternative was measured.

| Operation / variant | Time/validation, ns ↓ | Allocated/validation, B ↓ | Allocs/validation ↓ |
| --- | ---: | ---: | ---: |
| Ward, stop on failure | 749.9 | 1,008 | 28 |

<details>
<summary>Observed sample ranges</summary>

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| Ward, stop on failure | 749.4–752.2 ns | 1,008–1,008 B | 28 |

</details>

**Recorded samples:** [validator](../results/validator.txt)

<!-- benchmark-results:end -->

## Methodology

### Fields and rules

All five inputs are strings. Playground and Ozzo use explicit username, password,
and age predicates matching Ward. Each library retains its native email and URL
checks, which may differ on inputs beyond these fixtures.

| Field | Rules |
| --- | --- |
| Email | Required; native email validation. |
| Username | Required; 5–50 runes; ASCII letters, digits, `_`, `.`, `@`, and `-`. |
| Password | Required; 8–150 runes; an ASCII lowercase letter, uppercase letter, digit, and character from `@$!%*?&_-=`. |
| Age | Required decimal string from 0 through 999999999; no sign or extra leading zeroes. |
| Website | Required; native URL validation. |

[Behavioral tests](outcomes_test.go) check field outcomes, failure counts, and rule
boundaries. The comparisons establish matching fixture outcomes, not identical
validation standards for every possible string.

### Validation and error processing

| Workload | One operation |
| --- | --- |
| Valid or invalid email | Validate one email field. |
| Valid form | Validate five valid fields. |
| Invalid form | Validate the shared fixture in which all five fields fail. |
| Stop after first failure | Ward stops after its first failing field. |
| Validation and result processing | Validate, then inspect success or native errors. Ozzo formats error strings; the others inspect field and rule identifiers. |

Each library stops at the first failed rule within a field. Shared form cases
continue through all fields. Ward's stop-on-failure result performs less work
and is not an alternative implementation of the complete form validation case.
The legacy `MultiFieldSomeInvalid` name represents the all-five-invalid fixture.

### Initialization and reuse

| Library | Measured setup |
| --- | --- |
| Ward | Construct a validator and fields per operation, pointing at that operation's values. |
| Playground | Reuse a validator and registered predicates prepared before timing; reflection caches remain reusable. |
| Ozzo | Use its stateless API and assemble native errors per operation. |

All cases use `testing.B.Loop` and retain or inspect their results. Initialization
and error-processing differences are part of the measured API use.

## Running

Requires Go 1.27.4+. From the repository root:

```sh
go test -p=1 -run '^$' -bench . -benchmem -benchtime=1s -count=3 ./validator
```
