package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"testing"
	"time"

	"github.com/rah-0/benchmarks/util/testutil"
)

func TestMain(m *testing.M) {
	testutil.TestMainWrapper(testutil.TestConfig{
		M: m,
		LoadResources: func() error {
			return initAESKey()
		},
		UnloadResources: func() error {
			return nil
		},
	})
}

func TestGenerateAndDecryptEncryptedToken(t *testing.T) {
	token, err := generateEncryptedToken("user123")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	payload, err := decryptEncryptedToken(token)
	if err != nil {
		t.Fatalf("Failed to decrypt token: %v", err)
	}

	if payload.Sub != "user123" {
		t.Errorf("Expected sub=user123, got %s", payload.Sub)
	}

	if time.Now().Unix() > payload.Exp {
		t.Errorf("Token already expired")
	}
}

func TestExpiredToken(t *testing.T) {
	expired := TokenPayload{
		Sub: "user123",
		Exp: time.Now().Add(-1 * time.Hour).Unix(),
		Iat: time.Now().Unix(),
	}

	plain, err := json.Marshal(expired)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	block, _ := aes.NewCipher(aesKey)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)
	token := gcm.Seal(nonce, nonce, plain, nil)

	_, err = decryptEncryptedToken(token)
	if err == nil || err.Error() != "token expired" {
		t.Errorf("Expected token expired error, got: %v", err)
	}
}

func BenchmarkGenerateEncryptedToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := generateEncryptedToken("user123")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecryptEncryptedToken(b *testing.B) {
	token, err := generateEncryptedToken("user123")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := decryptEncryptedToken(token)
		if err != nil {
			b.Fatal(err)
		}
	}
}
