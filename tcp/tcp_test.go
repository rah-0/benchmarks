package main

import (
	"testing"

	"github.com/rah-0/benchmarks/util/testutil"
)

func TestMain(m *testing.M) {
	testutil.TestMainWrapper(testutil.TestConfig{
		M: m,
		LoadResources: func() error {
			if err := tcpStartServer(); err != nil {
				return err
			}

			client = &TCPClient{}
			if err := client.connect(); err != nil {
				return err
			}

			return nil
		},
		UnloadResources: func() error {
			if err := client.Close(); err != nil {
				return err
			}

			if err := server.Close(); err != nil {
				return err
			}

			return nil
		},
	})
}

func BenchmarkTCPClientSendSingleMessage(b *testing.B) {
	defer testutil.RecoverBenchHandler(b)

	msg := []byte("something")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := client.send(msg)
		if err != nil {
			b.Fatalf("Sending failed on iteration %d: %v", i, err)
		}
	}
}

func BenchmarkTCPClientSendSingleMessage2KB(b *testing.B)   { benchmarkTCP(b, 1) }
func BenchmarkTCPClientSendSingleMessage4KB(b *testing.B)   { benchmarkTCP(b, 2) }
func BenchmarkTCPClientSendSingleMessage8KB(b *testing.B)   { benchmarkTCP(b, 3) }
func BenchmarkTCPClientSendSingleMessage16KB(b *testing.B)  { benchmarkTCP(b, 4) }
func BenchmarkTCPClientSendSingleMessage32KB(b *testing.B)  { benchmarkTCP(b, 5) }
func BenchmarkTCPClientSendSingleMessage64KB(b *testing.B)  { benchmarkTCP(b, 6) }
func BenchmarkTCPClientSendSingleMessage128KB(b *testing.B) { benchmarkTCP(b, 7) }
func BenchmarkTCPClientSendSingleMessage256KB(b *testing.B) { benchmarkTCP(b, 8) }
func BenchmarkTCPClientSendSingleMessage512KB(b *testing.B) { benchmarkTCP(b, 9) }
func BenchmarkTCPClientSendSingleMessage1MB(b *testing.B)   { benchmarkTCP(b, 10) }
func BenchmarkTCPClientSendSingleMessage10MB(b *testing.B)  { benchmarkTCP(b, 11) }
func BenchmarkTCPClientSendSingleMessage100MB(b *testing.B) { benchmarkTCP(b, 12) }
func BenchmarkTCPClientSendSingleMessage1GB(b *testing.B)   { benchmarkTCP(b, 13) }

// Helper function for TCP benchmark
func benchmarkTCP(b *testing.B, sizeType int) {
	defer testutil.RecoverBenchHandler(b)

	// Generate a message of the given size
	message, err := testutil.GenerateMessage(sizeType)
	if err != nil {
		b.Fatalf("Failed to generate message: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	// Benchmark sending the message
	for i := 0; i < b.N; i++ {
		err = client.send(message)
		if err != nil {
			b.Fatalf("Sending failed on iteration %d: %v", i, err)
		}
	}
}
