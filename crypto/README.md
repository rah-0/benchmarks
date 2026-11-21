# Token benchmarks

Measure token creation and reading, including JSON and claim handling, for
HS256 JWTs and AES-256-GCM encrypted payloads. **Each operation is reported
separately:** signing and encryption provide different guarantees, so these
results have no interchangeable implementation baseline.

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

**Workloads:** [Signed JWT lifecycle](#signed-jwt-lifecycle) · [Encrypted token lifecycle](#encrypted-token-lifecycle)

### Signed JWT lifecycle

Create an HS256 JWT, or parse and verify a prepared JWT including its claims and expiry. Creation includes timestamp calls.

No equivalent reference measured; absolute costs are shown. Median of 3 samples; GOMAXPROCS=16.

Creation and verification perform different work. No equivalent implementation reference was measured.

| Operation / variant | Time/token, µs ↓ | Allocated/token, KiB ↓ | Allocs/token ↓ |
| --- | ---: | ---: | ---: |
| Create and sign | 8.04 | 2.28 | 44 |
| Parse and verify | 8.54 | 2.13 | 38 |

<details>
<summary>Observed sample ranges</summary>

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| Create and sign | 8.02–8.05 µs | 2.28–2.28 KiB | 44 |
| Parse and verify | 8.53–8.56 µs | 2.13–2.14 KiB | 38 |

</details>

### Encrypted token lifecycle

Create an AES-256-GCM token with JSON encoding and a fresh nonce, or decrypt a prepared token and decode its payload and expiry.

No equivalent reference measured; absolute costs are shown. Median of 3 samples; GOMAXPROCS=16.

Encryption and JWT signing provide different guarantees. These lifecycle costs are kept separate from the JWT results.

| Operation / variant | Time/token, µs ↓ | Allocated/token, KiB ↓ | Allocs/token ↓ |
| --- | ---: | ---: | ---: |
| Encode and encrypt | 5.79 | 1.47 | 7 |
| Decrypt and decode | 3.43 | 1.35 | 4 |

<details>
<summary>Observed sample ranges</summary>

| Operation / variant | Time min–max | Allocated min–max | Allocs min–max |
| --- | ---: | ---: | ---: |
| Encode and encrypt | 5.78–5.81 µs | 1.47–1.47 KiB | 7 |
| Decrypt and decode | 3.42–3.43 µs | 1.35–1.35 KiB | 4 |

</details>

**Recorded samples:** [crypto/encrypt](../results/crypto-encrypt.txt) · [crypto/sign](../results/crypto-sign.txt)

<!-- benchmark-results:end -->

## Methodology

| Operation | Timed work |
| --- | --- |
| Create JWT | Build claims, encode the token, and sign with HS256. |
| Verify JWT | Parse a prepared token, verify its signature, and check claims and expiry. |
| Encrypt payload | Encode JSON, construct AES-256-GCM, generate a nonce, and encrypt. |
| Decrypt payload | Construct AES-256-GCM, decrypt a prepared token, decode JSON, and check expiry. |

Key initialization and the tokens used for verification/decryption are prepared
before timing. Token creation includes timestamp calls. These fixtures measure
application-level operations; they are not a complete authentication system.

## Running

Requires Go 1.27.4+. From the repository root:

```sh
go test -p=1 -run '^$' -bench . -benchmem -benchtime=1s -count=3 ./crypto/...
```
