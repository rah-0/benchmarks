package encoding

import (
	"encoding/base64"
)

func EncodeBase64(input []byte) []byte {
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(input)))
	base64.StdEncoding.Encode(encoded, input)
	return encoded
}
