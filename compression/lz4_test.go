package compression

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	serializer2 "github.com/rah-0/benchmarks/serializer"
	"github.com/rah-0/benchmarks/util/testutil"
)

func TestLZ4_2KB(t *testing.T) {
	runLZ4TestRandomData(t, 1)
}

func TestLZ4_4KB(t *testing.T) {
	runLZ4TestRandomData(t, 2)
}

func TestLZ4_8KB(t *testing.T) {
	runLZ4TestRandomData(t, 3)
}

func TestLZ4_16KB(t *testing.T) {
	runLZ4TestRandomData(t, 4)
}

func TestLZ4_32KB(t *testing.T) {
	runLZ4TestRandomData(t, 5)
}

func TestLZ4_64KB(t *testing.T) {
	runLZ4TestRandomData(t, 6)
}

func TestLZ4_128KB(t *testing.T) {
	runLZ4TestRandomData(t, 7)
}

func TestLZ4_256KB(t *testing.T) {
	runLZ4TestRandomData(t, 8)
}

func TestLZ4_512KB(t *testing.T) {
	runLZ4TestRandomData(t, 9)
}

func TestLZ4_1MB(t *testing.T) {
	runLZ4TestRandomData(t, 10)
}

func TestLZ4_10MB(t *testing.T) {
	runLZ4TestRandomData(t, 11)
}

func TestLZ4_1Small(t *testing.T) {
	runLZ4TestStructsSmall(t, 1)
}

func TestLZ4_100Small(t *testing.T) {
	runLZ4TestStructsSmall(t, 100)
}

func TestLZ4_10000Small(t *testing.T) {
	runLZ4TestStructsSmall(t, 10000)
}

func TestLZ4_100000Small(t *testing.T) {
	runLZ4TestStructsSmall(t, 100000)
}

func runLZ4TestRandomData(t *testing.T, sizeType int) {
	originalData, err := testutil.GenerateMessage(sizeType)
	if err != nil {
		t.Fatalf("Failed to generate message of size type %d: %v", sizeType, err)
	}

	compressedData, err := CompressLZ4(originalData)
	if err != nil {
		t.Fatalf("Failed to compress data for size type %d: %v", sizeType, err)
	}

	decompressedData, err := DecompressLZ4(compressedData)
	if err != nil {
		t.Fatalf("Failed to decompress data for size type %d: %v", sizeType, err)
	}

	assert.Equal(t, originalData, decompressedData, "Decompressed data should match the original for size type %d", sizeType)

	dataSize := len(originalData)
	dataSizeCompressed := len(compressedData)
	fmt.Printf("Encoded size: %d bytes, Compressed size: %d bytes, Gain %.2f%%\n",
		dataSize, dataSizeCompressed, testutil.PercentDifference(dataSize, dataSizeCompressed))
}

func runLZ4TestStructsSmall(t *testing.T, size int) {
	originalData := serializer2.GenerateRandomPersons(size)
	err := serializer2.EncodeBytes(serializer2.GobEnc, originalData)
	if err != nil {
		t.Fatalf("Failed to encode data %d: %v", size, err)
	}
	encodedData := serializer2.GobBuf.Bytes()

	compressedData, err := CompressLZ4(encodedData)
	if err != nil {
		t.Fatalf("Failed to compress data for size %d: %v", size, err)
	}

	decompressedData, err := DecompressLZ4(compressedData)
	if err != nil {
		t.Fatalf("Failed to decompress data for size %d: %v", size, err)
	}

	assert.Equal(t, encodedData, decompressedData, "Decompressed data should match the original for size %d", size)

	// Output size comparison
	dataSize := len(encodedData)
	dataSizeCompressed := len(compressedData)
	fmt.Printf("Encoded size: %d bytes, Compressed size: %d bytes, Gain %.2f%%\n",
		dataSize, dataSizeCompressed, testutil.PercentDifference(dataSize, dataSizeCompressed))

	serializer2.GobBuf.Reset()
}

func BenchmarkCompressDecompressLZ4(b *testing.B) {
	data, err := testutil.GenerateMessage(10) // 10 corresponds to 1MB
	if err != nil {
		b.Fatalf("Failed to generate message of size type %d: %v", 10, err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		compressedData, err := CompressLZ4(data)
		if err != nil {
			b.Fatalf("Failed to Compress with LZ4: %v", err)
		}

		_, err = DecompressLZ4(compressedData)
		if err != nil {
			b.Fatalf("Failed to Decompress with LZ4: %v", err)
		}
	}
}
