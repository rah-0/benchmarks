package mem

import (
	"bytes"
	"io"
	"testing"
)

func BenchmarkCopy(b *testing.B) {
	compressedData := bytes.Repeat([]byte("Hello, World!"), 1000)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		br := bytes.NewReader(compressedData)
		var output bytes.Buffer
		_, err := io.Copy(&output, br)
		if err != nil {
			b.Fatal(err)
		}
		result := output.Bytes()
		_ = result
	}
}

func BenchmarkReadAll(b *testing.B) {
	compressedData := bytes.Repeat([]byte("Hello, World!"), 1000)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		br := bytes.NewReader(compressedData)
		result, err := io.ReadAll(br)
		if err != nil {
			b.Fatal(err)
		}
		_ = result
	}
}
