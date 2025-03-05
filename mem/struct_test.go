package mem

import (
	"testing"
)

type PersonV1 struct {
	ID   string
	Name string
	Age  int
}

type PersonV2 struct {
	ID   string
	Name string
	Age  int
}

func ConvertV1ToV2(v1 PersonV1) PersonV2 {
	return PersonV2{
		ID:   v1.ID,
		Name: v1.Name,
		Age:  v1.Age,
	}
}

func BenchmarkStructConversion(b *testing.B) {
	v1 := PersonV1{ID: "123", Name: "John Doe", Age: 30}

	for i := 0; i < b.N; i++ {
		_ = ConvertV1ToV2(v1)
	}
}
