package main

import (
	"io"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/moby/moby/client"
)

func TestCollectorInputEOFKeepsOutputOpen(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	conn, err := net.DialTimeout("tcp", listener.Addr().String(), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	peer, err := listener.Accept()
	if err != nil {
		t.Fatal(err)
	}
	defer peer.Close()
	deadline := time.Now().Add(5 * time.Second)
	if err := conn.SetDeadline(deadline); err != nil {
		t.Fatal(err)
	}
	if err := peer.SetDeadline(deadline); err != nil {
		t.Fatal(err)
	}
	response := client.NewHijackedResponse(conn, "")
	input := &collectorInput{response: &response}
	if _, err := io.WriteString(input, "record\n"); err != nil {
		t.Fatal(err)
	}
	for range 2 {
		if err := input.Close(); err != nil {
			t.Fatal(err)
		}
	}
	data, err := io.ReadAll(peer)
	if err != nil || string(data) != "record\n" {
		t.Fatalf("input=%q, error=%v", data, err)
	}
	if _, err := io.WriteString(peer, "drained"); err != nil {
		t.Fatal(err)
	}
	data = make([]byte, len("drained"))
	if _, err := io.ReadFull(response.Reader, data); err != nil || string(data) != "drained" {
		t.Fatalf("output after EOF=%q, error=%v", data, err)
	}
}

func TestProcessStatAllowsParenthesesInCommand(t *testing.T) {
	stat := "123 (worker name (batch)) R 1 2 3 4 5 6 7 8 9 10 125 25 0 0 20 0 4 0 999 1 2\n"
	cpu, start, err := parseProcessStat([]byte(stat))
	if err != nil || cpu != 150 || start != 999 {
		t.Fatalf("CPU=%d start=%d error=%v", cpu, start, err)
	}
}

func TestProcessStatRejectsIncompleteAndInvalidCounters(t *testing.T) {
	valid := "123 (worker) R 1 2 3 4 5 6 7 8 9 10 125 25 0 0 20 0 4 0 999"
	for _, stat := range []string{
		"",
		"123 worker R 1 2 3",
		"123 (worker) R 1 2 3",
		strings.Replace(valid, "125", "-1", 1),
		strings.Replace(valid, " 25 ", " invalid ", 1),
		strings.Replace(valid, "999", "invalid", 1),
		strings.Replace(valid, "125", "18446744073709551615", 1),
	} {
		if _, _, err := parseProcessStat([]byte(stat)); err == nil {
			t.Errorf("accepted invalid process stat %q", stat)
		}
	}
}

func TestPeakRSSUsesHighWaterMark(t *testing.T) {
	status := "Name:\tcollector\nVmPeak:\t9999 kB\nVmHWM:\t2048 kB\nVmRSS:\t1024 kB\n"
	rss, err := parsePeakRSS([]byte(status))
	if err != nil || rss != 2<<20 {
		t.Fatalf("RSS=%d error=%v", rss, err)
	}
	for _, status := range []string{"VmRSS: 1 kB\n", "VmHWM: 1 MB\n", "VmHWM: -1 kB\n", "VmHWM: 9007199254740992 kB\n"} {
		if _, err := parsePeakRSS([]byte(status)); err == nil {
			t.Errorf("accepted invalid process status %q", status)
		}
	}
}

func TestDiagnosticTailKeepsNewestBytes(t *testing.T) {
	var tail diagnosticTail
	want := strings.Repeat("b", diagnosticLimit-3) + "end"
	for _, data := range []string{strings.Repeat("a", diagnosticLimit+10), strings.Repeat("b", diagnosticLimit), "end"} {
		if n, err := tail.Write([]byte(data)); err != nil || n != len(data) {
			t.Fatalf("write count=%d error=%v", n, err)
		}
	}
	if got := tail.String(); got != want {
		t.Fatalf("diagnostics differ: length=%d, want=%d", len(got), len(want))
	}
}

func TestDiagnosticTailConcurrentReadersAndWriters(t *testing.T) {
	var tail diagnosticTail
	var workers sync.WaitGroup
	for range 4 {
		workers.Go(func() {
			for range 100 {
				tail.Write([]byte(strings.Repeat("record\n", 100)))
				if len(tail.String()) > diagnosticLimit {
					t.Error("diagnostic tail exceeded its limit")
				}
			}
		})
	}
	workers.Wait()
	if len(tail.String()) != diagnosticLimit {
		t.Fatalf("diagnostic length=%d", len(tail.String()))
	}
}
