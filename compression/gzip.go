package compression

import (
	"bytes"
	"compress/gzip"
)

func CompressGzip(input []byte) ([]byte, error) {
	var output bytes.Buffer

	writer, err := gzip.NewWriterLevel(&output, gzip.BestCompression)
	if err != nil {
		return nil, err
	}

	if _, err = writer.Write(input); err != nil {
		return nil, err
	}

	if err = writer.Close(); err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}

func DecompressGzip(input []byte) ([]byte, error) {
	var output bytes.Buffer

	reader, err := gzip.NewReader(bytes.NewReader(input))
	if err != nil {
		return nil, err
	}

	if _, err = output.ReadFrom(reader); err != nil {
		return nil, err
	}

	if err = reader.Close(); err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}
