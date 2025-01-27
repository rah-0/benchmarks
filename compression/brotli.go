package compression

import (
	"bytes"

	"github.com/andybalholm/brotli"
)

func CompressBrotli(input []byte) ([]byte, error) {
	var output bytes.Buffer

	writer := brotli.NewWriterLevel(&output, brotli.BestCompression)
	_, err := writer.Write(input)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}

func DecompressBrotli(input []byte) ([]byte, error) {
	var output bytes.Buffer

	reader := brotli.NewReader(bytes.NewReader(input))
	_, err := output.ReadFrom(reader)
	if err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}
