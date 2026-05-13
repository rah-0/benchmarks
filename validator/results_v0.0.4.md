# Benchmark Results — ward v0.0.4

```
goos: linux
goarch: amd64
pkg: github.com/rah-0/benchmarks/validator
cpu: AMD Ryzen 9 5950X 16-Core Processor
benchtime: 10s

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

## v0.0.3 → v0.0.4

| Benchmark | v0.0.3 ns/op | v0.0.4 ns/op | Δ |
|---|---|---|---|
| WardSingleFieldValid | 230.5 | 234.8 | +2% (noise) |
| WardSingleFieldInvalid | 421.9 | 401.0 | -5% (noise) |
| WardMultiFieldAllValid | 1140 | 1111 | -3% (noise) |
| WardMultiFieldSomeInvalid | 1377 | 1333 | -3% (noise) |
| WardMultiFieldStopOnFail | 421.8 | 402.7 | -5% (noise) |
