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
