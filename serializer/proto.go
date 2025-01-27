package serializer

import (
	"google.golang.org/protobuf/proto"
)

func EncodeProto[T proto.Message](message T) ([]byte, error) {
	return proto.Marshal(message)
}

func DecodeProto[T proto.Message](data []byte, message T) error {
	return proto.Unmarshal(data, message)
}
