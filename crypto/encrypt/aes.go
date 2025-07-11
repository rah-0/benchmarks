package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

var aesKey []byte

func initAESKey() error {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return err
	}
	aesKey = key
	return nil
}

type TokenPayload struct {
	Sub string `json:"sub"`
	Exp int64  `json:"exp"`
	Iat int64  `json:"iat"`
}

func generateEncryptedToken(userID string) ([]byte, error) {
	payload := TokenPayload{
		Sub: userID,
		Exp: time.Now().Add(15 * time.Minute).Unix(),
		Iat: time.Now().Unix(),
	}

	plain, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plain, nil)
	return ciphertext, nil
}

func decryptEncryptedToken(ciphertext []byte) (*TokenPayload, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ct := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plain, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return nil, err
	}

	var payload TokenPayload
	if err = json.Unmarshal(plain, &payload); err != nil {
		return nil, err
	}

	if time.Now().Unix() > payload.Exp {
		return nil, fmt.Errorf("token expired")
	}

	return &payload, nil
}
