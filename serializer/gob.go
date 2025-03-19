package serializer

import (
	"bytes"
	"encoding/gob"
)

var (
	GobBuf = &bytes.Buffer{}
	GobEnc = gob.NewEncoder(GobBuf)
	GobDec = gob.NewDecoder(GobBuf)
)

func EncodeBytes(enc *gob.Encoder, x any) error {
	if err := enc.Encode(x); err != nil {
		return err
	}
	return nil
}

func DecodeBytes(dec *gob.Decoder, x any) error {
	return dec.Decode(x)
}
