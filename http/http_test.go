package http

import (
	"testing"

	"github.com/rah-0/benchmarks/util/testutil"
)

func TestMain(m *testing.M) {
	testutil.TestMainWrapper(testutil.TestConfig{
		M: m,
		LoadResources: func() error {
			go httpServerStart()

			if err := NewHTTPClient(); err != nil {
				return err
			}

			return nil
		},
	})
}

func BenchmarkHTTPClientSendSingleMessage(b *testing.B) {
	msg := []byte("something")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := client.sendBytes(msg)
		if err != nil {
			b.Fatalf("Failed to send message: %v", err)
		}
	}
}

func BenchmarkHTTPClientSendSingleMessage2KB(b *testing.B)   { benchmarkHTTP(b, 1) }
func BenchmarkHTTPClientSendSingleMessage4KB(b *testing.B)   { benchmarkHTTP(b, 2) }
func BenchmarkHTTPClientSendSingleMessage8KB(b *testing.B)   { benchmarkHTTP(b, 3) }
func BenchmarkHTTPClientSendSingleMessage16KB(b *testing.B)  { benchmarkHTTP(b, 4) }
func BenchmarkHTTPClientSendSingleMessage32KB(b *testing.B)  { benchmarkHTTP(b, 5) }
func BenchmarkHTTPClientSendSingleMessage64KB(b *testing.B)  { benchmarkHTTP(b, 6) }
func BenchmarkHTTPClientSendSingleMessage128KB(b *testing.B) { benchmarkHTTP(b, 7) }
func BenchmarkHTTPClientSendSingleMessage256KB(b *testing.B) { benchmarkHTTP(b, 8) }
func BenchmarkHTTPClientSendSingleMessage512KB(b *testing.B) { benchmarkHTTP(b, 9) }
func BenchmarkHTTPClientSendSingleMessage1MB(b *testing.B)   { benchmarkHTTP(b, 10) }
func BenchmarkHTTPClientSendSingleMessage10MB(b *testing.B)  { benchmarkHTTP(b, 11) }
func BenchmarkHTTPClientSendSingleMessage100MB(b *testing.B) { benchmarkHTTP(b, 12) }
func BenchmarkHTTPClientSendSingleMessage1GB(b *testing.B)   { benchmarkHTTP(b, 13) }

func benchmarkHTTP(b *testing.B, sizeType int) {
	defer testutil.RecoverBenchHandler(b)

	msg, err := testutil.GenerateMessage(sizeType)
	if err != nil {
		b.Fatalf("Failed to generate message: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err = client.sendBytes(msg)
		if err != nil {
			b.Fatalf("Sending failed on iteration %d: %v", i, err)
		}
	}
}
