package serializer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rah-0/benchmarks/util"
	"github.com/rah-0/benchmarks/util/testutil"
)

func TestEncodeDecodeJson1Small(t *testing.T) {
	testEncodeDecodeJsonPersons(t, 1)
}

func TestEncodeDecodeJson100Small(t *testing.T) {
	testEncodeDecodeJsonPersons(t, 100)
}

func TestEncodeDecodeJson10000Small(t *testing.T) {
	testEncodeDecodeJsonPersons(t, 10000)
}

func TestEncodeDecodeJson1000000Small(t *testing.T) {
	testEncodeDecodeJsonPersons(t, 1000000)
}

func testEncodeDecodeJsonPersons(t *testing.T, count int) {
	persons := util.GenerateRandomPersons(count)
	encodedData := make([][]byte, count)

	for i, original := range persons {
		data, err := EncodeJson(&original)
		assert.NoError(t, err, "JSON serialization should not fail")
		assert.NotEmpty(t, data, "Serialized data should not be empty")
		encodedData[i] = data
	}

	for i, data := range encodedData {
		decoded := &util.Person{}
		err := DecodeJson(data, decoded)
		assert.NoError(t, err, "JSON deserialization should not fail")

		assert.Equal(t, persons[i].Name, decoded.Name, "Names should match")
		assert.Equal(t, persons[i].Age, decoded.Age, "Ages should match")
		assert.Equal(t, persons[i].Surname, decoded.Surname, "Surnames should match")
	}
}

func TestEncodeDecodeJson1Unreal(t *testing.T) {
	testEncodeDecodeJsonUnreals(t, 1)
}

func TestEncodeDecodeJson10Unreals(t *testing.T) {
	testEncodeDecodeJsonUnreals(t, 10)
}

func TestEncodeDecodeJson100Unreals(t *testing.T) {
	testEncodeDecodeJsonUnreals(t, 100)
}

func TestEncodeDecodeJson1000Unreals(t *testing.T) {
	testEncodeDecodeJsonUnreals(t, 1000)
}

func testEncodeDecodeJsonUnreals(t *testing.T, count int) {
	unreals := util.GenerateRandomUnreals(count)
	encodedData := make([][]byte, count)

	for i, original := range unreals {
		data, err := EncodeJson(&original)
		assert.NoError(t, err, "JSON serialization should not fail")
		assert.NotEmpty(t, data, "Serialized data should not be empty")
		encodedData[i] = data
	}

	for i, data := range encodedData {
		decoded := &util.Unreal{}
		err := DecodeJson(data, decoded)
		assert.NoError(t, err, "JSON deserialization should not fail")

		assert.Equal(t, unreals[i].Prop1, decoded.Prop1, "Prop1 should match")
		assert.Equal(t, unreals[i].Prop2, decoded.Prop2, "Prop2 should match")
		assert.Equal(t, unreals[i].Prop3, decoded.Prop3, "Prop3 should match")
		assert.Equal(t, unreals[i].Prop4, decoded.Prop4, "Prop4 should match")
		assert.Equal(t, unreals[i].Prop5, decoded.Prop5, "Prop5 should match")
		assert.Equal(t, unreals[i].Prop6, decoded.Prop6, "Prop6 should match")
		assert.Equal(t, unreals[i].Prop7, decoded.Prop7, "Prop7 should match")
		assert.Equal(t, unreals[i].Prop8, decoded.Prop8, "Prop8 should match")
		assert.Equal(t, unreals[i].Prop9, decoded.Prop9, "Prop9 should match")
		assert.Equal(t, unreals[i].Prop10, decoded.Prop10, "Prop10 should match")
		assert.Equal(t, unreals[i].Prop11, decoded.Prop11, "Prop11 should match")
		assert.Equal(t, unreals[i].Prop12, decoded.Prop12, "Prop12 should match")
		assert.Equal(t, unreals[i].Prop13, decoded.Prop13, "Prop13 should match")
		assert.Equal(t, unreals[i].Prop14, decoded.Prop14, "Prop14 should match")
		assert.Equal(t, unreals[i].Prop15, decoded.Prop15, "Prop15 should match")
		assert.Equal(t, unreals[i].Prop16, decoded.Prop16, "Prop16 should match")
		assert.Equal(t, unreals[i].Prop17, decoded.Prop17, "Prop17 should match")
		assert.Equal(t, unreals[i].Prop18, decoded.Prop18, "Prop18 should match")
		assert.Equal(t, unreals[i].Prop19, decoded.Prop19, "Prop19 should match")
		assert.Equal(t, unreals[i].Prop20, decoded.Prop20, "Prop20 should match")
		assert.Equal(t, unreals[i].Prop21, decoded.Prop21, "Prop21 should match")
		assert.Equal(t, unreals[i].Prop22, decoded.Prop22, "Prop22 should match")
		assert.Equal(t, unreals[i].Prop23, decoded.Prop23, "Prop23 should match")
		assert.Equal(t, unreals[i].Prop24, decoded.Prop24, "Prop24 should match")
		assert.Equal(t, unreals[i].Prop25, decoded.Prop25, "Prop25 should match")
		assert.Equal(t, unreals[i].Prop26, decoded.Prop26, "Prop26 should match")
		assert.Equal(t, unreals[i].Prop27, decoded.Prop27, "Prop27 should match")
		assert.Equal(t, unreals[i].Prop28, decoded.Prop28, "Prop28 should match")
		assert.Equal(t, unreals[i].Prop29, decoded.Prop29, "Prop29 should match")
		assert.Equal(t, unreals[i].Prop30, decoded.Prop30, "Prop30 should match")
		assert.Equal(t, unreals[i].Prop31, decoded.Prop31, "Prop31 should match")
		assert.Equal(t, unreals[i].Prop32, decoded.Prop32, "Prop32 should match")
		assert.Equal(t, unreals[i].Prop33, decoded.Prop33, "Prop33 should match")
		assert.Equal(t, unreals[i].Prop34, decoded.Prop34, "Prop34 should match")
		assert.Equal(t, unreals[i].Prop35, decoded.Prop35, "Prop35 should match")
		assert.Equal(t, unreals[i].Prop36, decoded.Prop36, "Prop36 should match")
		assert.Equal(t, unreals[i].Prop37, decoded.Prop37, "Prop37 should match")
		assert.Equal(t, unreals[i].Prop38, decoded.Prop38, "Prop38 should match")
		assert.Equal(t, unreals[i].Prop39, decoded.Prop39, "Prop39 should match")
		assert.Equal(t, unreals[i].Prop40, decoded.Prop40, "Prop40 should match")
		assert.Equal(t, unreals[i].Prop41, decoded.Prop41, "Prop41 should match")
		assert.Equal(t, unreals[i].Prop42, decoded.Prop42, "Prop42 should match")
		assert.Equal(t, unreals[i].Prop43, decoded.Prop43, "Prop43 should match")
		assert.Equal(t, unreals[i].Prop44, decoded.Prop44, "Prop44 should match")
		assert.Equal(t, unreals[i].Prop45, decoded.Prop45, "Prop45 should match")
		assert.Equal(t, unreals[i].Prop46, decoded.Prop46, "Prop46 should match")
		assert.Equal(t, unreals[i].Prop47, decoded.Prop47, "Prop47 should match")
		assert.Equal(t, unreals[i].Prop48, decoded.Prop48, "Prop48 should match")
		assert.Equal(t, unreals[i].Prop49, decoded.Prop49, "Prop49 should match")
		assert.Equal(t, unreals[i].Prop50, decoded.Prop50, "Prop50 should match")
		assert.Equal(t, unreals[i].Prop51, decoded.Prop51, "Prop51 should match")
		assert.Equal(t, unreals[i].Prop52, decoded.Prop52, "Prop52 should match")
		assert.Equal(t, unreals[i].Prop53, decoded.Prop53, "Prop53 should match")
		assert.Equal(t, unreals[i].Prop54, decoded.Prop54, "Prop54 should match")
		assert.Equal(t, unreals[i].Prop55, decoded.Prop55, "Prop55 should match")
		assert.Equal(t, unreals[i].Prop56, decoded.Prop56, "Prop56 should match")
		assert.Equal(t, unreals[i].Prop57, decoded.Prop57, "Prop57 should match")
		assert.Equal(t, unreals[i].Prop58, decoded.Prop58, "Prop58 should match")
		assert.Equal(t, unreals[i].Prop59, decoded.Prop59, "Prop59 should match")
		assert.Equal(t, unreals[i].Prop60, decoded.Prop60, "Prop60 should match")
		assert.Equal(t, unreals[i].Prop61, decoded.Prop61, "Prop61 should match")
		assert.Equal(t, unreals[i].Prop62, decoded.Prop62, "Prop62 should match")
		assert.Equal(t, unreals[i].Prop63, decoded.Prop63, "Prop63 should match")
		assert.Equal(t, unreals[i].Prop64, decoded.Prop64, "Prop64 should match")
		assert.Equal(t, unreals[i].Prop65, decoded.Prop65, "Prop65 should match")
		assert.Equal(t, unreals[i].Prop66, decoded.Prop66, "Prop66 should match")
		assert.Equal(t, unreals[i].Prop67, decoded.Prop67, "Prop67 should match")
		assert.Equal(t, unreals[i].Prop68, decoded.Prop68, "Prop68 should match")
		assert.Equal(t, unreals[i].Prop69, decoded.Prop69, "Prop69 should match")
		assert.Equal(t, unreals[i].Prop70, decoded.Prop70, "Prop70 should match")
		assert.Equal(t, unreals[i].Prop71, decoded.Prop71, "Prop71 should match")
		assert.Equal(t, unreals[i].Prop72, decoded.Prop72, "Prop72 should match")
		assert.Equal(t, unreals[i].Prop73, decoded.Prop73, "Prop73 should match")
		assert.Equal(t, unreals[i].Prop74, decoded.Prop74, "Prop74 should match")
		assert.Equal(t, unreals[i].Prop75, decoded.Prop75, "Prop75 should match")
		assert.Equal(t, unreals[i].Prop76, decoded.Prop76, "Prop76 should match")
		assert.Equal(t, unreals[i].Prop77, decoded.Prop77, "Prop77 should match")
		assert.Equal(t, unreals[i].Prop78, decoded.Prop78, "Prop78 should match")
		assert.Equal(t, unreals[i].Prop79, decoded.Prop79, "Prop79 should match")
		assert.Equal(t, unreals[i].Prop80, decoded.Prop80, "Prop80 should match")
		assert.Equal(t, unreals[i].Prop81, decoded.Prop81, "Prop81 should match")
		assert.Equal(t, unreals[i].Prop82, decoded.Prop82, "Prop82 should match")
		assert.Equal(t, unreals[i].Prop83, decoded.Prop83, "Prop83 should match")
		assert.Equal(t, unreals[i].Prop84, decoded.Prop84, "Prop84 should match")
		assert.Equal(t, unreals[i].Prop85, decoded.Prop85, "Prop85 should match")
		assert.Equal(t, unreals[i].Prop86, decoded.Prop86, "Prop86 should match")
		assert.Equal(t, unreals[i].Prop87, decoded.Prop87, "Prop87 should match")
		assert.Equal(t, unreals[i].Prop88, decoded.Prop88, "Prop88 should match")
		assert.Equal(t, unreals[i].Prop89, decoded.Prop89, "Prop89 should match")
		assert.Equal(t, unreals[i].Prop90, decoded.Prop90, "Prop90 should match")
		assert.Equal(t, unreals[i].Prop91, decoded.Prop91, "Prop91 should match")
		assert.Equal(t, unreals[i].Prop92, decoded.Prop92, "Prop92 should match")
		assert.Equal(t, unreals[i].Prop93, decoded.Prop93, "Prop93 should match")
		assert.Equal(t, unreals[i].Prop94, decoded.Prop94, "Prop94 should match")
		assert.Equal(t, unreals[i].Prop95, decoded.Prop95, "Prop95 should match")
		assert.Equal(t, unreals[i].Prop96, decoded.Prop96, "Prop96 should match")
		assert.Equal(t, unreals[i].Prop97, decoded.Prop97, "Prop97 should match")
		assert.Equal(t, unreals[i].Prop98, decoded.Prop98, "Prop98 should match")
		assert.Equal(t, unreals[i].Prop99, decoded.Prop99, "Prop99 should match")
		assert.Equal(t, unreals[i].Prop100, decoded.Prop100, "Prop100 should match")
	}
}

func BenchmarkEncodeDecodeJson1Small(b *testing.B) {
	runEncodeDecodeJsonBenchmark(b, 1)
}

func BenchmarkEncodeDecodeJson100Small(b *testing.B) {
	runEncodeDecodeJsonBenchmark(b, 100)
}

func BenchmarkEncodeDecodeJson10000Small(b *testing.B) {
	runEncodeDecodeJsonBenchmark(b, 10000)
}

func BenchmarkEncodeDecodeJson1000000Small(b *testing.B) {
	runEncodeDecodeJsonBenchmark(b, 1000000)
}

func runEncodeDecodeJsonBenchmark(b *testing.B, count int) {
	defer testutil.RecoverBenchHandler(b)

	data := util.GenerateRandomPersons(count)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encodedData, err := EncodeJson(data)
		if err != nil {
			b.Fatalf("Failed to Encode JSON: %v", err)
		}

		var decoded []util.Person
		err = DecodeJson(encodedData, &decoded)
		if err != nil {
			b.Fatalf("Failed to Decode JSON: %v", err)
		}
	}
}

func BenchmarkEncodeDecodeJson1Unreal(b *testing.B) {
	runEncodeDecodeJsonUnrealBenchmark(b, 1)
}

func BenchmarkEncodeDecodeJson10Unreals(b *testing.B) {
	runEncodeDecodeJsonUnrealBenchmark(b, 10)
}

func BenchmarkEncodeDecodeJson100Unreals(b *testing.B) {
	runEncodeDecodeJsonUnrealBenchmark(b, 100)
}

func BenchmarkEncodeDecodeJson1000Unreals(b *testing.B) {
	runEncodeDecodeJsonUnrealBenchmark(b, 1000)
}

func runEncodeDecodeJsonUnrealBenchmark(b *testing.B, count int) {
	defer testutil.RecoverBenchHandler(b)

	data := util.GenerateRandomUnreals(count)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encodedData, err := EncodeJson(data)
		if err != nil {
			b.Fatalf("Failed to Encode JSON: %v", err)
		}

		var decoded []util.Unreal
		err = DecodeJson(encodedData, &decoded)
		if err != nil {
			b.Fatalf("Failed to Decode JSON: %v", err)
		}
	}
}
