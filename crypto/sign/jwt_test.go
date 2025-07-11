package sign

import (
	"testing"
)

func BenchmarkJWTGeneration(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := generateJWT("user123")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJWTVerification(b *testing.B) {
	tokenStr, err := generateJWT("user123")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := verifyJWT(tokenStr)
		if err != nil {
			b.Fatal(err)
		}
	}
}
