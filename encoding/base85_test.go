package encoding

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rah-0/benchmarks/util/testutil"
)

func TestBase85_2KB(t *testing.T)   { runBase85TestRandomData(t, 1) }
func TestBase85_4KB(t *testing.T)   { runBase85TestRandomData(t, 2) }
func TestBase85_8KB(t *testing.T)   { runBase85TestRandomData(t, 3) }
func TestBase85_16KB(t *testing.T)  { runBase85TestRandomData(t, 4) }
func TestBase85_32KB(t *testing.T)  { runBase85TestRandomData(t, 5) }
func TestBase85_64KB(t *testing.T)  { runBase85TestRandomData(t, 6) }
func TestBase85_128KB(t *testing.T) { runBase85TestRandomData(t, 7) }
func TestBase85_256KB(t *testing.T) { runBase85TestRandomData(t, 8) }
func TestBase85_512KB(t *testing.T) { runBase85TestRandomData(t, 9) }
func TestBase85_1MB(t *testing.T)   { runBase85TestRandomData(t, 10) }
func TestBase85_10MB(t *testing.T)  { runBase85TestRandomData(t, 11) }
func TestBase85_100MB(t *testing.T) { runBase85TestRandomData(t, 12) }
func TestBase85_1GB(t *testing.T)   { runBase85TestRandomData(t, 13) }

func runBase85TestRandomData(t *testing.T, sizeType int) {
	originalData, err := testutil.GenerateMessage(sizeType)
	if err != nil {
		t.Fatalf("Failed to generate message of size type %d: %v", sizeType, err)
	}

	encodedData := EncodeBase85(originalData)

	dataSize := len(originalData)
	dataSizeEncoded := len(encodedData)
	sizeIncrease := testutil.PercentDifference(dataSize, dataSizeEncoded)

	fmt.Printf("Original size: %d bytes, Encoded size: %d bytes, Increase: %.2f%%\n",
		dataSize, dataSizeEncoded, sizeIncrease)

	assert.Greater(t, dataSizeEncoded, dataSize, "Encoded data should be larger than the original")
}

func BenchmarkBase85_2KB(b *testing.B)   { runBase85Benchmark(b, 1) }
func BenchmarkBase85_4KB(b *testing.B)   { runBase85Benchmark(b, 2) }
func BenchmarkBase85_8KB(b *testing.B)   { runBase85Benchmark(b, 3) }
func BenchmarkBase85_16KB(b *testing.B)  { runBase85Benchmark(b, 4) }
func BenchmarkBase85_32KB(b *testing.B)  { runBase85Benchmark(b, 5) }
func BenchmarkBase85_64KB(b *testing.B)  { runBase85Benchmark(b, 6) }
func BenchmarkBase85_128KB(b *testing.B) { runBase85Benchmark(b, 7) }
func BenchmarkBase85_256KB(b *testing.B) { runBase85Benchmark(b, 8) }
func BenchmarkBase85_512KB(b *testing.B) { runBase85Benchmark(b, 9) }
func BenchmarkBase85_1MB(b *testing.B)   { runBase85Benchmark(b, 10) }
func BenchmarkBase85_10MB(b *testing.B)  { runBase85Benchmark(b, 11) }
func BenchmarkBase85_100MB(b *testing.B) { runBase85Benchmark(b, 12) }
func BenchmarkBase85_1GB(b *testing.B)   { runBase85Benchmark(b, 13) }

func runBase85Benchmark(b *testing.B, sizeType int) {
	data, err := testutil.GenerateMessage(sizeType)
	if err != nil {
		b.Fatalf("Failed to generate message of size type %d: %v", sizeType, err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = EncodeBase85(data)
	}
}
