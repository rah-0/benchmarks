package serializer

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rah-0/benchmarks/util/testutil"
)

func TestEncodeDecodeBytes(t *testing.T) {
	original := GenerateRandomPersons(1)[0]

	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	dec := gob.NewDecoder(buf)

	err := EncodeBytes(enc, original)
	assert.NoError(t, err, "GOB serialization should not fail")

	decoded := &Person{}
	err = DecodeBytes(dec, decoded)
	assert.NoError(t, err, "GOB deserialization should not fail")

	assert.Equal(t, original.Name, decoded.Name, "Names should match")
	assert.Equal(t, original.Age, decoded.Age, "Ages should match")
	assert.Equal(t, original.Surname, decoded.Surname, "Surnames should match")
}

func BenchmarkEncodeDecodeBytes1Small(b *testing.B) {
	runEncodeDecodeBytesBenchmark(b, 1)
}

func BenchmarkEncodeDecodeBytes100Small(b *testing.B) {
	runEncodeDecodeBytesBenchmark(b, 100)
}

func BenchmarkEncodeDecodeBytes10000Small(b *testing.B) {
	runEncodeDecodeBytesBenchmark(b, 10000)
}

func BenchmarkEncodeDecodeBytes1000000Small(b *testing.B) {
	runEncodeDecodeBytesBenchmark(b, 1000000)
}

func runEncodeDecodeBytesBenchmark(b *testing.B, count int) {
	defer testutil.RecoverBenchHandler(b)

	data := GenerateRandomPersons(count)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := &bytes.Buffer{}
		var decoded []Person

		enc := gob.NewEncoder(buf)
		err := EncodeBytes(enc, data)
		if err != nil {
			b.Fatalf("Failed to Encode Gob: %v", err)
		}

		dec := gob.NewDecoder(buf)
		err = DecodeBytes(dec, &decoded)
		if err != nil {
			b.Fatalf("Failed to Decode Gob: %v", err)
		}
	}
}

func BenchmarkEncodeDecodeBytes1Unreal(b *testing.B) {
	runEncodeDecodeBytesUnrealBenchmark(b, 1)
}

func BenchmarkEncodeDecodeBytes10Unreals(b *testing.B) {
	runEncodeDecodeBytesUnrealBenchmark(b, 10)
}

func BenchmarkEncodeDecodeBytes100Unreals(b *testing.B) {
	runEncodeDecodeBytesUnrealBenchmark(b, 100)
}

func BenchmarkEncodeDecodeBytes1000Unreals(b *testing.B) {
	runEncodeDecodeBytesUnrealBenchmark(b, 1000)
}

func runEncodeDecodeBytesUnrealBenchmark(b *testing.B, count int) {
	defer testutil.RecoverBenchHandler(b)

	data := GenerateRandomUnreals(count)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := &bytes.Buffer{}
		var decoded []Unreal

		enc := gob.NewEncoder(buf)
		err := EncodeBytes(enc, data)
		if err != nil {
			b.Fatalf("Failed to Encode Gob: %v", err)
		}

		dec := gob.NewDecoder(buf)
		err = DecodeBytes(dec, &decoded)
		if err != nil {
			b.Fatalf("Failed to Decode Gob: %v", err)
		}
	}
}
