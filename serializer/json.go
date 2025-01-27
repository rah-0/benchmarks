package serializer

import (
	"encoding/json"
)

func EncodeJson(x any) ([]byte, error) {
	return json.Marshal(x)
}

func DecodeJson(data []byte, x any) error {
	return json.Unmarshal(data, x)
}
