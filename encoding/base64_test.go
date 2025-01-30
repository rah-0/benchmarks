package encoding

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rah-0/benchmarks/util/testutil"
)

func TestBase64_2KB(t *testing.T)   { runBase64TestRandomData(t, 1) }
func TestBase64_4KB(t *testing.T)   { runBase64TestRandomData(t, 2) }
func TestBase64_8KB(t *testing.T)   { runBase64TestRandomData(t, 3) }
func TestBase64_16KB(t *testing.T)  { runBase64TestRandomData(t, 4) }
func TestBase64_32KB(t *testing.T)  { runBase64TestRandomData(t, 5) }
func TestBase64_64KB(t *testing.T)  { runBase64TestRandomData(t, 6) }
func TestBase64_128KB(t *testing.T) { runBase64TestRandomData(t, 7) }
func TestBase64_256KB(t *testing.T) { runBase64TestRandomData(t, 8) }
func TestBase64_512KB(t *testing.T) { runBase64TestRandomData(t, 9) }
func TestBase64_1MB(t *testing.T)   { runBase64TestRandomData(t, 10) }
func TestBase64_10MB(t *testing.T)  { runBase64TestRandomData(t, 11) }
func TestBase64_100MB(t *testing.T) { runBase64TestRandomData(t, 12) }
func TestBase64_1GB(t *testing.T)   { runBase64TestRandomData(t, 13) }

func runBase64TestRandomData(t *testing.T, sizeType int) {
	originalData, err := testutil.GenerateMessage(sizeType)
	if err != nil {
		t.Fatalf("Failed to generate message of size type %d: %v", sizeType, err)
	}

	encodedData := EncodeBase64(originalData)

	dataSize := len(originalData)
	dataSizeEncoded := len(encodedData)
	sizeIncrease := testutil.PercentDifference(dataSize, dataSizeEncoded)

	fmt.Printf("Original size: %d bytes, Encoded size: %d bytes, Increase: %.2f%%\n",
		dataSize, dataSizeEncoded, sizeIncrease)

	assert.Greater(t, dataSizeEncoded, dataSize, "Encoded data should be larger than the original")
}

func BenchmarkBase64_2KB(b *testing.B)   { runBase64Benchmark(b, 1) }
func BenchmarkBase64_4KB(b *testing.B)   { runBase64Benchmark(b, 2) }
func BenchmarkBase64_8KB(b *testing.B)   { runBase64Benchmark(b, 3) }
func BenchmarkBase64_16KB(b *testing.B)  { runBase64Benchmark(b, 4) }
func BenchmarkBase64_32KB(b *testing.B)  { runBase64Benchmark(b, 5) }
func BenchmarkBase64_64KB(b *testing.B)  { runBase64Benchmark(b, 6) }
func BenchmarkBase64_128KB(b *testing.B) { runBase64Benchmark(b, 7) }
func BenchmarkBase64_256KB(b *testing.B) { runBase64Benchmark(b, 8) }
func BenchmarkBase64_512KB(b *testing.B) { runBase64Benchmark(b, 9) }
func BenchmarkBase64_1MB(b *testing.B)   { runBase64Benchmark(b, 10) }
func BenchmarkBase64_10MB(b *testing.B)  { runBase64Benchmark(b, 11) }
func BenchmarkBase64_100MB(b *testing.B) { runBase64Benchmark(b, 12) }
func BenchmarkBase64_1GB(b *testing.B)   { runBase64Benchmark(b, 13) }

func runBase64Benchmark(b *testing.B, sizeType int) {
	data, err := testutil.GenerateMessage(sizeType)
	if err != nil {
		b.Fatalf("Failed to generate message of size type %d: %v", sizeType, err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = EncodeBase64(data)
	}
}
