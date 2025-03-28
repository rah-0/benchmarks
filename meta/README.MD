# Reflection vs Direct Method Access

This benchmark measures the performance of setting and getting a struct field using direct method access versus using reflection

## Results

| Operation            | Time (ns/op) | Reflection Slowdown (x) |
|----------------------|--------------|-------------------------|
| SetName (Method)     | 0.4402       | 1x (Baseline)           |
| GetName (Method)     | 0.2199       | 1x (Baseline)           |
| SetName (Reflection) | 2.204        | ~5x slower              |
| GetName (Reflection) | 0.4398       | ~2x slower              |

## Observations
- **Reflection is ~5x slower** for setting a field compared to direct method access
- **Reflection is ~2x slower** for retrieving a field


# Sorting

## ⏱️ Benchmark Relative Performance (% of Go `sort.Ints`) – Higher is better

| Size        | Base (Go) | Sorty (%) | Pargo (%) | Parsort (%) |
|-------------|-----------|-----------|-----------|-------------|
| 10,000      | 100%      | 113.3%    | 65.7%     | 112.2%      |
| 100,000     | 100%      | 181.8%    | 150.7%    | 166.4%      |
| 1,000,000   | 100%      | 284.0%    | 245.4%    | 377.3%      |
| 10,000,000  | 100%      | 230.6%    | 233.7%    | 434.5%      |

### 📌 Notes:
- Values > 100% mean the alternative is **faster** than Go's sort.Ints.
- Parsort shows the most consistent and significant gains, especially with large inputs.
- Pargo suffers from high allocation overhead, which affects small-size performance.
- Sorty performs better than Go for medium-sized inputs but lags at small sizes.
