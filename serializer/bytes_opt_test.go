package serializer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rah-0/benchmarks/util/testutil"
)

func TestEncodeDecodeBytesOpt1Small(t *testing.T) {
	testEncodeDecodeOptBytesPersons(t, 1)
}

func TestEncodeDecodeBytesOpt100Small(t *testing.T) {
	testEncodeDecodeOptBytesPersons(t, 100)
}

func TestEncodeDecodeBytesOpt10000Small(t *testing.T) {
	testEncodeDecodeOptBytesPersons(t, 10000)
}

func TestEncodeDecodeBytesOpt1000000Small(t *testing.T) {
	testEncodeDecodeOptBytesPersons(t, 1000000)
}

func testEncodeDecodeOptBytesPersons(t *testing.T, count int) {
	persons := GenerateRandomPersons(count)

	for _, original := range persons {
		err := EncodeBytes(GobEnc, original)
		assert.NoError(t, err, "GOB serialization should not fail")
		assert.NotZero(t, GobBuf.Len(), "Serialized data should not be empty")
	}

	for _, original := range persons {
		decoded := &Person{}
		err := DecodeBytes(GobDec, decoded)
		assert.NoError(t, err, "GOB deserialization should not fail")

		assert.Equal(t, original.Name, decoded.Name, "Names should match")
		assert.Equal(t, original.Age, decoded.Age, "Ages should match")
		assert.Equal(t, original.Surname, decoded.Surname, "Surnames should match")
	}

	GobBuf.Reset()
}

func TestEncodeDecodeBytesOpt1Unreal(t *testing.T) {
	testEncodeDecodeOptBytesUnreals(t, 1)
}

func TestEncodeDecodeBytesOpt10Unreal(t *testing.T) {
	testEncodeDecodeOptBytesUnreals(t, 10)
}

func TestEncodeDecodeBytesOpt100Unreal(t *testing.T) {
	testEncodeDecodeOptBytesUnreals(t, 100)
}

func TestEncodeDecodeBytesOpt1000Unreal(t *testing.T) {
	testEncodeDecodeOptBytesUnreals(t, 1000)
}

func testEncodeDecodeOptBytesUnreals(t *testing.T, count int) {
	unreals := GenerateRandomUnreals(count)

	for _, original := range unreals {
		err := EncodeBytes(GobEnc, original)
		assert.NoError(t, err, "GOB serialization should not fail")
		assert.NotZero(t, GobBuf.Len(), "Serialized data should not be empty")
	}

	for _, original := range unreals {
		decoded := &Unreal{}
		err := DecodeBytes(GobDec, decoded)
		assert.NoError(t, err, "GOB deserialization should not fail")

		assert.Equal(t, original.Prop1, decoded.Prop1, "Prop1 should match")
		assert.Equal(t, original.Prop2, decoded.Prop2, "Prop2 should match")
		assert.Equal(t, original.Prop3, decoded.Prop3, "Prop3 should match")
		assert.Equal(t, original.Prop4, decoded.Prop4, "Prop4 should match")
		assert.Equal(t, original.Prop5, decoded.Prop5, "Prop5 should match")
		assert.Equal(t, original.Prop6, decoded.Prop6, "Prop6 should match")
		assert.Equal(t, original.Prop7, decoded.Prop7, "Prop7 should match")
		assert.Equal(t, original.Prop8, decoded.Prop8, "Prop8 should match")
		assert.Equal(t, original.Prop9, decoded.Prop9, "Prop9 should match")
		assert.Equal(t, original.Prop10, decoded.Prop10, "Prop10 should match")
		assert.Equal(t, original.Prop11, decoded.Prop11, "Prop11 should match")
		assert.Equal(t, original.Prop12, decoded.Prop12, "Prop12 should match")
		assert.Equal(t, original.Prop13, decoded.Prop13, "Prop13 should match")
		assert.Equal(t, original.Prop14, decoded.Prop14, "Prop14 should match")
		assert.Equal(t, original.Prop15, decoded.Prop15, "Prop15 should match")
		assert.Equal(t, original.Prop16, decoded.Prop16, "Prop16 should match")
		assert.Equal(t, original.Prop17, decoded.Prop17, "Prop17 should match")
		assert.Equal(t, original.Prop18, decoded.Prop18, "Prop18 should match")
		assert.Equal(t, original.Prop19, decoded.Prop19, "Prop19 should match")
		assert.Equal(t, original.Prop20, decoded.Prop20, "Prop20 should match")
		assert.Equal(t, original.Prop21, decoded.Prop21, "Prop21 should match")
		assert.Equal(t, original.Prop22, decoded.Prop22, "Prop22 should match")
		assert.Equal(t, original.Prop23, decoded.Prop23, "Prop23 should match")
		assert.Equal(t, original.Prop24, decoded.Prop24, "Prop24 should match")
		assert.Equal(t, original.Prop25, decoded.Prop25, "Prop25 should match")
		assert.Equal(t, original.Prop26, decoded.Prop26, "Prop26 should match")
		assert.Equal(t, original.Prop27, decoded.Prop27, "Prop27 should match")
		assert.Equal(t, original.Prop28, decoded.Prop28, "Prop28 should match")
		assert.Equal(t, original.Prop29, decoded.Prop29, "Prop29 should match")
		assert.Equal(t, original.Prop30, decoded.Prop30, "Prop30 should match")
		assert.Equal(t, original.Prop31, decoded.Prop31, "Prop31 should match")
		assert.Equal(t, original.Prop32, decoded.Prop32, "Prop32 should match")
		assert.Equal(t, original.Prop33, decoded.Prop33, "Prop33 should match")
		assert.Equal(t, original.Prop34, decoded.Prop34, "Prop34 should match")
		assert.Equal(t, original.Prop35, decoded.Prop35, "Prop35 should match")
		assert.Equal(t, original.Prop36, decoded.Prop36, "Prop36 should match")
		assert.Equal(t, original.Prop37, decoded.Prop37, "Prop37 should match")
		assert.Equal(t, original.Prop38, decoded.Prop38, "Prop38 should match")
		assert.Equal(t, original.Prop39, decoded.Prop39, "Prop39 should match")
		assert.Equal(t, original.Prop40, decoded.Prop40, "Prop40 should match")
		assert.Equal(t, original.Prop41, decoded.Prop41, "Prop41 should match")
		assert.Equal(t, original.Prop42, decoded.Prop42, "Prop42 should match")
		assert.Equal(t, original.Prop43, decoded.Prop43, "Prop43 should match")
		assert.Equal(t, original.Prop44, decoded.Prop44, "Prop44 should match")
		assert.Equal(t, original.Prop45, decoded.Prop45, "Prop45 should match")
		assert.Equal(t, original.Prop46, decoded.Prop46, "Prop46 should match")
		assert.Equal(t, original.Prop47, decoded.Prop47, "Prop47 should match")
		assert.Equal(t, original.Prop48, decoded.Prop48, "Prop48 should match")
		assert.Equal(t, original.Prop49, decoded.Prop49, "Prop49 should match")
		assert.Equal(t, original.Prop50, decoded.Prop50, "Prop50 should match")
		assert.Equal(t, original.Prop51, decoded.Prop51, "Prop51 should match")
		assert.Equal(t, original.Prop52, decoded.Prop52, "Prop52 should match")
		assert.Equal(t, original.Prop53, decoded.Prop53, "Prop53 should match")
		assert.Equal(t, original.Prop54, decoded.Prop54, "Prop54 should match")
		assert.Equal(t, original.Prop55, decoded.Prop55, "Prop55 should match")
		assert.Equal(t, original.Prop56, decoded.Prop56, "Prop56 should match")
		assert.Equal(t, original.Prop57, decoded.Prop57, "Prop57 should match")
		assert.Equal(t, original.Prop58, decoded.Prop58, "Prop58 should match")
		assert.Equal(t, original.Prop59, decoded.Prop59, "Prop59 should match")
		assert.Equal(t, original.Prop60, decoded.Prop60, "Prop60 should match")
		assert.Equal(t, original.Prop61, decoded.Prop61, "Prop61 should match")
		assert.Equal(t, original.Prop62, decoded.Prop62, "Prop62 should match")
		assert.Equal(t, original.Prop63, decoded.Prop63, "Prop63 should match")
		assert.Equal(t, original.Prop64, decoded.Prop64, "Prop64 should match")
		assert.Equal(t, original.Prop65, decoded.Prop65, "Prop65 should match")
		assert.Equal(t, original.Prop66, decoded.Prop66, "Prop66 should match")
		assert.Equal(t, original.Prop67, decoded.Prop67, "Prop67 should match")
		assert.Equal(t, original.Prop68, decoded.Prop68, "Prop68 should match")
		assert.Equal(t, original.Prop69, decoded.Prop69, "Prop69 should match")
		assert.Equal(t, original.Prop70, decoded.Prop70, "Prop70 should match")
		assert.Equal(t, original.Prop71, decoded.Prop71, "Prop71 should match")
		assert.Equal(t, original.Prop72, decoded.Prop72, "Prop72 should match")
		assert.Equal(t, original.Prop73, decoded.Prop73, "Prop73 should match")
		assert.Equal(t, original.Prop74, decoded.Prop74, "Prop74 should match")
		assert.Equal(t, original.Prop75, decoded.Prop75, "Prop75 should match")
		assert.Equal(t, original.Prop76, decoded.Prop76, "Prop76 should match")
		assert.Equal(t, original.Prop77, decoded.Prop77, "Prop77 should match")
		assert.Equal(t, original.Prop78, decoded.Prop78, "Prop78 should match")
		assert.Equal(t, original.Prop79, decoded.Prop79, "Prop79 should match")
		assert.Equal(t, original.Prop80, decoded.Prop80, "Prop80 should match")
		assert.Equal(t, original.Prop81, decoded.Prop81, "Prop81 should match")
		assert.Equal(t, original.Prop82, decoded.Prop82, "Prop82 should match")
		assert.Equal(t, original.Prop83, decoded.Prop83, "Prop83 should match")
		assert.Equal(t, original.Prop84, decoded.Prop84, "Prop84 should match")
		assert.Equal(t, original.Prop85, decoded.Prop85, "Prop85 should match")
		assert.Equal(t, original.Prop86, decoded.Prop86, "Prop86 should match")
		assert.Equal(t, original.Prop87, decoded.Prop87, "Prop87 should match")
		assert.Equal(t, original.Prop88, decoded.Prop88, "Prop88 should match")
		assert.Equal(t, original.Prop89, decoded.Prop89, "Prop89 should match")
		assert.Equal(t, original.Prop90, decoded.Prop90, "Prop90 should match")
		assert.Equal(t, original.Prop91, decoded.Prop91, "Prop91 should match")
		assert.Equal(t, original.Prop92, decoded.Prop92, "Prop92 should match")
		assert.Equal(t, original.Prop93, decoded.Prop93, "Prop93 should match")
		assert.Equal(t, original.Prop94, decoded.Prop94, "Prop94 should match")
		assert.Equal(t, original.Prop95, decoded.Prop95, "Prop95 should match")
		assert.Equal(t, original.Prop96, decoded.Prop96, "Prop96 should match")
		assert.Equal(t, original.Prop97, decoded.Prop97, "Prop97 should match")
		assert.Equal(t, original.Prop98, decoded.Prop98, "Prop98 should match")
		assert.Equal(t, original.Prop99, decoded.Prop99, "Prop99 should match")
		assert.Equal(t, original.Prop100, decoded.Prop100, "Prop100 should match")
	}

	GobBuf.Reset()
}

func BenchmarkEncodeDecodeBytesOpt1Small(b *testing.B) {
	runEncodeDecodeBytesOptBenchmark(b, 1)
}

func BenchmarkEncodeDecodeBytesOpt100Small(b *testing.B) {
	runEncodeDecodeBytesOptBenchmark(b, 100)
}

func BenchmarkEncodeDecodeBytesOpt10000Small(b *testing.B) {
	runEncodeDecodeBytesOptBenchmark(b, 10000)
}

func BenchmarkEncodeDecodeBytesOpt1000000Small(b *testing.B) {
	runEncodeDecodeBytesOptBenchmark(b, 1000000)
}

func BenchmarkEncodeDecodeBytesOpt1Unreal(b *testing.B) {
	runEncodeDecodeBytesOptUnrealBenchmark(b, 1)
}

func BenchmarkEncodeDecodeBytesOpt10Unreals(b *testing.B) {
	runEncodeDecodeBytesOptUnrealBenchmark(b, 10)
}

func BenchmarkEncodeDecodeBytesOpt100Unreals(b *testing.B) {
	runEncodeDecodeBytesOptUnrealBenchmark(b, 100)
}

func BenchmarkEncodeDecodeBytesOpt1000Unreals(b *testing.B) {
	runEncodeDecodeBytesOptUnrealBenchmark(b, 1000)
}

func runEncodeDecodeBytesOptBenchmark(b *testing.B, count int) {
	defer testutil.RecoverBenchHandler(b)

	data := GenerateRandomPersons(count)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := EncodeBytes(GobEnc, data)
		if err != nil {
			b.Fatalf("Failed to Encode Gob: %v", err)
		}
		var decoded []Person
		err = DecodeBytes(GobDec, &decoded)
		if err != nil {
			b.Fatalf("Failed to Decode Gob: %v", err)
		}
	}

	GobBuf.Reset()
}

func runEncodeDecodeBytesOptUnrealBenchmark(b *testing.B, count int) {
	defer testutil.RecoverBenchHandler(b)

	data := GenerateRandomUnreals(count)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := EncodeBytes(GobEnc, data)
		if err != nil {
			b.Fatalf("Failed to Encode Gob: %v", err)
		}
		var decoded []Unreal
		err = DecodeBytes(GobDec, &decoded)
		if err != nil {
			b.Fatalf("Failed to Decode Gob: %v", err)
		}
	}

	GobBuf.Reset()
}
