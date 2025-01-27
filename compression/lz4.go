package compression

import (
	"bytes"

	"github.com/pierrec/lz4/v4"
)

func CompressLZ4(data []byte) ([]byte, error) {
	var output bytes.Buffer

	writer := lz4.NewWriter(&output)
	writer.Apply(
		lz4.CompressionLevelOption(lz4.Level9),
	)

	if _, err := writer.Write(data); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}

func DecompressLZ4(data []byte) ([]byte, error) {
	var output bytes.Buffer

	reader := lz4.NewReader(bytes.NewReader(data))
	if _, err := output.ReadFrom(reader); err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}
