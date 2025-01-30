package compression

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rah-0/benchmarks/serializer"
	"github.com/rah-0/benchmarks/util"
	"github.com/rah-0/benchmarks/util/testutil"
)

func TestBrotli_2KB(t *testing.T) {
	runBrotliTestRandomData(t, 1)
}

func TestBrotli_4KB(t *testing.T) {
	runBrotliTestRandomData(t, 2)
}

func TestBrotli_8KB(t *testing.T) {
	runBrotliTestRandomData(t, 3)
}

func TestBrotli_16KB(t *testing.T) {
	runBrotliTestRandomData(t, 4)
}

func TestBrotli_32KB(t *testing.T) {
	runBrotliTestRandomData(t, 5)
}

func TestBrotli_64KB(t *testing.T) {
	runBrotliTestRandomData(t, 6)
}

func TestBrotli_128KB(t *testing.T) {
	runBrotliTestRandomData(t, 7)
}

func TestBrotli_256KB(t *testing.T) {
	runBrotliTestRandomData(t, 8)
}

func TestBrotli_512KB(t *testing.T) {
	runBrotliTestRandomData(t, 9)
}

func TestBrotli_1MB(t *testing.T) {
	runBrotliTestRandomData(t, 10)
}

func TestBrotli_10MB(t *testing.T) {
	runBrotliTestRandomData(t, 11)
}

func TestBrotli_1Small(t *testing.T) {
	runBrotliTestStructsSmall(t, 1)
}

func TestBrotli_100Small(t *testing.T) {
	runBrotliTestStructsSmall(t, 100)
}

func TestBrotli_10000Small(t *testing.T) {
	runBrotliTestStructsSmall(t, 10000)
}

func TestBrotli_100000Small(t *testing.T) {
	runBrotliTestStructsSmall(t, 100000)
}

func runBrotliTestRandomData(t *testing.T, sizeType int) {
	originalData, err := testutil.GenerateMessage(sizeType)
	if err != nil {
		t.Fatalf("Failed to generate message of size type %d: %v", sizeType, err)
	}

	compressedData, err := CompressBrotli(originalData)
	if err != nil {
		t.Fatalf("Failed to compress data for size type %d: %v", sizeType, err)
	}

	decompressedData, err := DecompressBrotli(compressedData)
	if err != nil {
		t.Fatalf("Failed to decompress data for size type %d: %v", sizeType, err)
	}

	assert.Equal(t, originalData, decompressedData, "Decompressed data should match the original for size type %d", sizeType)

	dataSize := len(originalData)
	dataSizeCompressed := len(compressedData)
	fmt.Printf("Encoded size: %d bytes, Compressed size: %d bytes, Gain %.2f%%\n",
		dataSize, dataSizeCompressed, testutil.PercentDifference(dataSize, dataSizeCompressed))
}

func runBrotliTestStructsSmall(t *testing.T, size int) {
	originalData := util.GenerateRandomPersons(size)
	err := serializer.EncodeBytes(serializer.GobEnc, originalData)
	if err != nil {
		t.Fatalf("Failed to encode data %d: %v", size, err)
	}
	encodedData := serializer.GobBuf.Bytes()

	compressedData, err := CompressBrotli(encodedData)
	if err != nil {
		t.Fatalf("Failed to compress data for size %d: %v", size, err)
	}

	decompressedData, err := DecompressBrotli(compressedData)
	if err != nil {
		t.Fatalf("Failed to decompress data for type %d: %v", size, err)
	}

	assert.Equal(t, encodedData, decompressedData, "Decompressed data should match the original for size %d", size)

	dataSize := len(encodedData)
	dataSizeCompressed := len(compressedData)
	fmt.Printf("Encoded size: %d bytes, Compressed size: %d bytes, Gain %.2f%%\n",
		dataSize, dataSizeCompressed, testutil.PercentDifference(dataSize, dataSizeCompressed))

	serializer.GobBuf.Reset()
}

func BenchmarkCompressDecompressBrotli(b *testing.B) {
	data, err := testutil.GenerateMessage(10)
	if err != nil {
		b.Fatalf("Failed to generate message of size type %d: %v", 10, err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		compressedData, err := CompressBrotli(data)
		if err != nil {
			b.Fatalf("Failed to Compress with Brotli : %v", err)
		}

		_, err = DecompressBrotli(compressedData)
		if err != nil {
			b.Fatalf("Failed to Decompress with Brotli : %v", err)
		}
	}
}
