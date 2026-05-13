# Benchmark Results — ward v0.0.3

```
goos: linux
goarch: amd64
pkg: github.com/rah-0/benchmarks/validator
cpu: AMD Ryzen 9 5950X 16-Core Processor

BenchmarkOzzoSingleFieldValid-8              	 5044948	       714.1 ns/op	      64 B/op	       3 allocs/op
BenchmarkOzzoSingleFieldInvalid-8            	28502269	       125.2 ns/op	      64 B/op	       3 allocs/op
BenchmarkOzzoMultiFieldAllValid-8            	  397057	      8990 ns/op	    1667 B/op	      24 allocs/op
BenchmarkOzzoMultiFieldSomeInvalid-8         	 2034433	      1718 ns/op	    1616 B/op	      24 allocs/op
BenchmarkPlaygroundSingleFieldValid-8        	 4472498	       817.9 ns/op	     105 B/op	       6 allocs/op
BenchmarkPlaygroundSingleFieldInvalid-8      	 5362477	       620.6 ns/op	     352 B/op	      13 allocs/op
BenchmarkPlaygroundMultiFieldAllValid-8      	 1831135	      1937 ns/op	     527 B/op	      11 allocs/op
BenchmarkPlaygroundMultiFieldSomeInvalid-8   	 2127007	      1626 ns/op	    1153 B/op	      22 allocs/op
BenchmarkWardSingleFieldValid-8              	15776647	       230.5 ns/op	      88 B/op	       5 allocs/op
BenchmarkWardSingleFieldInvalid-8            	 8930869	       421.9 ns/op	     216 B/op	      10 allocs/op
BenchmarkWardMultiFieldAllValid-8            	 3169851	      1140 ns/op	     234 B/op	       6 allocs/op
BenchmarkWardMultiFieldSomeInvalid-8         	 2601549	      1377 ns/op	     877 B/op	      22 allocs/op
BenchmarkWardMultiFieldStopOnFail-8          	 8634157	       421.8 ns/op	     216 B/op	      10 allocs/op
```

## v0.0.2 → v0.0.3

| Benchmark | v0.0.2 ns/op | v0.0.3 ns/op | Δ |
|---|---|---|---|
| WardSingleFieldValid | 225.5 | 230.5 | +2% (noise) |
| WardSingleFieldInvalid | 398.1 | 421.9 | +6% (noise) |
| WardMultiFieldAllValid | 1110 | 1140 | +3% (noise) |
| WardMultiFieldSomeInvalid | 1316 | 1377 | +5% (noise) |
| WardMultiFieldStopOnFail | 412.9 | 421.8 | +2% (noise) |
