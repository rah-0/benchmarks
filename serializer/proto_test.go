package serializer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	pb2 "github.com/rah-0/benchmarks/serializer/pb"
	"github.com/rah-0/benchmarks/util"
	"github.com/rah-0/benchmarks/util/testutil"
)

func TestEncodeDecodePerson1Small(t *testing.T) {
	testEncodeDecodePersons(t, 1)
}

func TestEncodeDecodeProto100Small(t *testing.T) {
	testEncodeDecodePersons(t, 100)
}

func TestEncodeDecodeProto10000Small(t *testing.T) {
	testEncodeDecodePersons(t, 10000)
}

func TestEncodeDecodeProto1000000Small(t *testing.T) {
	testEncodeDecodePersons(t, 1000000)
}

func testEncodeDecodePersons(t *testing.T, count int) {
	persons := util.GenerateRandomPersons(count)
	encodedData := make([][]byte, count)

	for i, original := range persons {
		data, err := EncodeProto(&pb2.Person{
			Name:    original.Name,
			Age:     int32(original.Age),
			Surname: original.Surname,
		})
		assert.NoError(t, err, "Serialization of Person should not fail")
		assert.NotEmpty(t, data, "Serialized data should not be empty")
		encodedData[i] = data
	}

	for i, data := range encodedData {
		decoded := &pb2.Person{}
		err := DecodeProto(data, decoded)
		assert.NoError(t, err, "Deserialization of Person should not fail")

		assert.Equal(t, persons[i].Name, decoded.Name, "Names should match")
		assert.Equal(t, persons[i].Age, int(decoded.Age), "Ages should match")
		assert.Equal(t, persons[i].Surname, decoded.Surname, "Surnames should match")
	}
}

func TestEncodeDecodeUnreal1(t *testing.T) {
	testEncodeDecodeUnreals(t, 1)
}

func TestEncodeDecodeUnreal10(t *testing.T) {
	testEncodeDecodeUnreals(t, 10)
}

func TestEncodeDecodeUnreal100(t *testing.T) {
	testEncodeDecodeUnreals(t, 100)
}

func TestEncodeDecodeUnreal1000(t *testing.T) {
	testEncodeDecodeUnreals(t, 1000)
}

func testEncodeDecodeUnreals(t *testing.T, count int) {
	unreals := util.GenerateRandomUnreals(count)
	encodedData := make([][]byte, count)

	for i, original := range unreals {
		data, err := EncodeProto(&pb2.Unreal{
			Prop1:   original.Prop1,
			Prop2:   original.Prop2,
			Prop3:   original.Prop3,
			Prop4:   original.Prop4,
			Prop5:   original.Prop5,
			Prop6:   original.Prop6,
			Prop7:   original.Prop7,
			Prop8:   original.Prop8,
			Prop9:   original.Prop9,
			Prop10:  original.Prop10,
			Prop11:  original.Prop11,
			Prop12:  original.Prop12,
			Prop13:  original.Prop13,
			Prop14:  original.Prop14,
			Prop15:  original.Prop15,
			Prop16:  original.Prop16,
			Prop17:  original.Prop17,
			Prop18:  original.Prop18,
			Prop19:  original.Prop19,
			Prop20:  original.Prop20,
			Prop21:  original.Prop21,
			Prop22:  original.Prop22,
			Prop23:  original.Prop23,
			Prop24:  original.Prop24,
			Prop25:  original.Prop25,
			Prop26:  original.Prop26,
			Prop27:  original.Prop27,
			Prop28:  original.Prop28,
			Prop29:  original.Prop29,
			Prop30:  original.Prop30,
			Prop31:  original.Prop31,
			Prop32:  original.Prop32,
			Prop33:  original.Prop33,
			Prop34:  original.Prop34,
			Prop35:  original.Prop35,
			Prop36:  original.Prop36,
			Prop37:  original.Prop37,
			Prop38:  original.Prop38,
			Prop39:  original.Prop39,
			Prop40:  original.Prop40,
			Prop41:  original.Prop41,
			Prop42:  original.Prop42,
			Prop43:  original.Prop43,
			Prop44:  original.Prop44,
			Prop45:  original.Prop45,
			Prop46:  original.Prop46,
			Prop47:  original.Prop47,
			Prop48:  original.Prop48,
			Prop49:  original.Prop49,
			Prop50:  original.Prop50,
			Prop51:  original.Prop51,
			Prop52:  original.Prop52,
			Prop53:  original.Prop53,
			Prop54:  original.Prop54,
			Prop55:  original.Prop55,
			Prop56:  original.Prop56,
			Prop57:  original.Prop57,
			Prop58:  original.Prop58,
			Prop59:  original.Prop59,
			Prop60:  original.Prop60,
			Prop61:  original.Prop61,
			Prop62:  original.Prop62,
			Prop63:  original.Prop63,
			Prop64:  original.Prop64,
			Prop65:  original.Prop65,
			Prop66:  original.Prop66,
			Prop67:  original.Prop67,
			Prop68:  original.Prop68,
			Prop69:  original.Prop69,
			Prop70:  original.Prop70,
			Prop71:  original.Prop71,
			Prop72:  original.Prop72,
			Prop73:  original.Prop73,
			Prop74:  original.Prop74,
			Prop75:  original.Prop75,
			Prop76:  original.Prop76,
			Prop77:  original.Prop77,
			Prop78:  original.Prop78,
			Prop79:  original.Prop79,
			Prop80:  original.Prop80,
			Prop81:  original.Prop81,
			Prop82:  original.Prop82,
			Prop83:  original.Prop83,
			Prop84:  original.Prop84,
			Prop85:  original.Prop85,
			Prop86:  original.Prop86,
			Prop87:  original.Prop87,
			Prop88:  original.Prop88,
			Prop89:  original.Prop89,
			Prop90:  original.Prop90,
			Prop91:  original.Prop91,
			Prop92:  original.Prop92,
			Prop93:  original.Prop93,
			Prop94:  original.Prop94,
			Prop95:  original.Prop95,
			Prop96:  original.Prop96,
			Prop97:  original.Prop97,
			Prop98:  original.Prop98,
			Prop99:  original.Prop99,
			Prop100: original.Prop100,
		})
		assert.NoError(t, err, "Serialization of Unreal should not fail")
		assert.NotEmpty(t, data, "Serialized data should not be empty")
		encodedData[i] = data
	}

	for i, data := range encodedData {
		decoded := &pb2.Unreal{}
		err := DecodeProto(data, decoded)
		assert.NoError(t, err, "Deserialization of Unreal should not fail")

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

func BenchmarkEncodeDecodeProto1Small(b *testing.B) {
	runEncodeDecodeProtoBenchmark(b, 1)
}

func BenchmarkEncodeDecodeProto100Small(b *testing.B) {
	runEncodeDecodeProtoBenchmark(b, 100)
}

func BenchmarkEncodeDecodeProto10000Small(b *testing.B) {
	runEncodeDecodeProtoBenchmark(b, 10000)
}

func BenchmarkEncodeDecodeProto1000000Small(b *testing.B) {
	runEncodeDecodeProtoBenchmark(b, 1000000)
}

func runEncodeDecodeProtoBenchmark(b *testing.B, count int) {
	defer testutil.RecoverBenchHandler(b)

	data := util.GenerateRandomPersons(count)
	protoData := make([]*pb2.Person, len(data))
	for i, p := range data {
		protoData[i] = &pb2.Person{
			Name:    p.Name,
			Age:     int32(p.Age),
			Surname: p.Surname,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, person := range protoData {
			encodedData, err := EncodeProto(person)
			if err != nil {
				b.Fatalf("Failed to Encode Protobuf: %v", err)
			}

			var decoded pb2.Person
			err = DecodeProto(encodedData, &decoded)
			if err != nil {
				b.Fatalf("Failed to Decode Protobuf: %v", err)
			}
		}
	}
}

func BenchmarkEncodeDecodeProto1Unreal(b *testing.B) {
	runEncodeDecodeProtoUnrealBenchmark(b, 1)
}

func BenchmarkEncodeDecodeProto10Unreals(b *testing.B) {
	runEncodeDecodeProtoUnrealBenchmark(b, 10)
}

func BenchmarkEncodeDecodeProto100Unreals(b *testing.B) {
	runEncodeDecodeProtoUnrealBenchmark(b, 100)
}

func BenchmarkEncodeDecodeProto1000Unreals(b *testing.B) {
	runEncodeDecodeProtoUnrealBenchmark(b, 1000)
}

func runEncodeDecodeProtoUnrealBenchmark(b *testing.B, count int) {
	defer testutil.RecoverBenchHandler(b)

	data := util.GenerateRandomUnreals(count)
	protoData := make([]*pb2.Unreal, len(data))
	for i, u := range data {
		protoData[i] = &pb2.Unreal{
			Prop1:   u.Prop1,
			Prop2:   u.Prop2,
			Prop3:   u.Prop3,
			Prop4:   u.Prop4,
			Prop5:   u.Prop5,
			Prop6:   u.Prop6,
			Prop7:   u.Prop7,
			Prop8:   u.Prop8,
			Prop9:   u.Prop9,
			Prop10:  u.Prop10,
			Prop11:  u.Prop11,
			Prop12:  u.Prop12,
			Prop13:  u.Prop13,
			Prop14:  u.Prop14,
			Prop15:  u.Prop15,
			Prop16:  u.Prop16,
			Prop17:  u.Prop17,
			Prop18:  u.Prop18,
			Prop19:  u.Prop19,
			Prop20:  u.Prop20,
			Prop21:  u.Prop21,
			Prop22:  u.Prop22,
			Prop23:  u.Prop23,
			Prop24:  u.Prop24,
			Prop25:  u.Prop25,
			Prop26:  u.Prop26,
			Prop27:  u.Prop27,
			Prop28:  u.Prop28,
			Prop29:  u.Prop29,
			Prop30:  u.Prop30,
			Prop31:  u.Prop31,
			Prop32:  u.Prop32,
			Prop33:  u.Prop33,
			Prop34:  u.Prop34,
			Prop35:  u.Prop35,
			Prop36:  u.Prop36,
			Prop37:  u.Prop37,
			Prop38:  u.Prop38,
			Prop39:  u.Prop39,
			Prop40:  u.Prop40,
			Prop41:  u.Prop41,
			Prop42:  u.Prop42,
			Prop43:  u.Prop43,
			Prop44:  u.Prop44,
			Prop45:  u.Prop45,
			Prop46:  u.Prop46,
			Prop47:  u.Prop47,
			Prop48:  u.Prop48,
			Prop49:  u.Prop49,
			Prop50:  u.Prop50,
			Prop51:  u.Prop51,
			Prop52:  u.Prop52,
			Prop53:  u.Prop53,
			Prop54:  u.Prop54,
			Prop55:  u.Prop55,
			Prop56:  u.Prop56,
			Prop57:  u.Prop57,
			Prop58:  u.Prop58,
			Prop59:  u.Prop59,
			Prop60:  u.Prop60,
			Prop61:  u.Prop61,
			Prop62:  u.Prop62,
			Prop63:  u.Prop63,
			Prop64:  u.Prop64,
			Prop65:  u.Prop65,
			Prop66:  u.Prop66,
			Prop67:  u.Prop67,
			Prop68:  u.Prop68,
			Prop69:  u.Prop69,
			Prop70:  u.Prop70,
			Prop71:  u.Prop71,
			Prop72:  u.Prop72,
			Prop73:  u.Prop73,
			Prop74:  u.Prop74,
			Prop75:  u.Prop75,
			Prop76:  u.Prop76,
			Prop77:  u.Prop77,
			Prop78:  u.Prop78,
			Prop79:  u.Prop79,
			Prop80:  u.Prop80,
			Prop81:  u.Prop81,
			Prop82:  u.Prop82,
			Prop83:  u.Prop83,
			Prop84:  u.Prop84,
			Prop85:  u.Prop85,
			Prop86:  u.Prop86,
			Prop87:  u.Prop87,
			Prop88:  u.Prop88,
			Prop89:  u.Prop89,
			Prop90:  u.Prop90,
			Prop91:  u.Prop91,
			Prop92:  u.Prop92,
			Prop93:  u.Prop93,
			Prop94:  u.Prop94,
			Prop95:  u.Prop95,
			Prop96:  u.Prop96,
			Prop97:  u.Prop97,
			Prop98:  u.Prop98,
			Prop99:  u.Prop99,
			Prop100: u.Prop100,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, unreal := range protoData {
			encodedData, err := EncodeProto(unreal)
			if err != nil {
				b.Fatalf("Failed to Encode Protobuf: %v", err)
			}

			var decoded pb2.Unreal
			err = DecodeProto(encodedData, &decoded)
			if err != nil {
				b.Fatalf("Failed to Decode Protobuf: %v", err)
			}
		}
	}
}
