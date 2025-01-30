package encoding

import (
	"encoding/ascii85"
)

func EncodeBase85(input []byte) []byte {
	encoded := make([]byte, ascii85.MaxEncodedLen(len(input)))
	n := ascii85.Encode(encoded, input)
	return encoded[:n]
}
