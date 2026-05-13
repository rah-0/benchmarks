# Benchmark Results — ward v0.0.2

```
goos: linux
goarch: amd64
pkg: github.com/rah-0/benchmarks/validator
cpu: AMD Ryzen 9 5950X 16-Core Processor

BenchmarkOzzoSingleFieldValid-8              	 4972308	       717.9 ns/op	      64 B/op	       3 allocs/op
BenchmarkOzzoSingleFieldInvalid-8            	29044957	       122.9 ns/op	      64 B/op	       3 allocs/op
BenchmarkOzzoMultiFieldAllValid-8            	  389754	      8995 ns/op	    1663 B/op	      24 allocs/op
BenchmarkOzzoMultiFieldSomeInvalid-8         	 2100700	      1709 ns/op	    1616 B/op	      24 allocs/op
BenchmarkPlaygroundSingleFieldValid-8        	 4399880	       809.6 ns/op	     104 B/op	       6 allocs/op
BenchmarkPlaygroundSingleFieldInvalid-8      	 5749041	       623.8 ns/op	     352 B/op	      13 allocs/op
BenchmarkPlaygroundMultiFieldAllValid-8      	 1873872	      1935 ns/op	     527 B/op	      11 allocs/op
BenchmarkPlaygroundMultiFieldSomeInvalid-8   	 2228694	      1627 ns/op	    1153 B/op	      22 allocs/op
BenchmarkWardSingleFieldValid-8              	15896882	       225.5 ns/op	      88 B/op	       5 allocs/op
BenchmarkWardSingleFieldInvalid-8            	 9096153	       398.1 ns/op	     216 B/op	      10 allocs/op
BenchmarkWardMultiFieldAllValid-8            	 3230227	      1110 ns/op	     234 B/op	       6 allocs/op
BenchmarkWardMultiFieldSomeInvalid-8         	 2716425	      1316 ns/op	     876 B/op	      22 allocs/op
BenchmarkWardMultiFieldStopOnFail-8          	 8820512	       412.9 ns/op	     216 B/op	      10 allocs/op
```
