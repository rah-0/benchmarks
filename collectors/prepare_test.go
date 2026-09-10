package main

import (
	"slices"
	"strings"
	"testing"
)

func TestCollectorBuildEnv(t *testing.T) {
	environment := []string{
		"PATH=/usr/bin", "HOME=/home/example", "GOCACHE=/cache", "GOFLAGS=-buildvcs=false",
		"GOBIN=/global/bin", "GOBIN=/other/bin", "CGO_ENABLED=1",
		"GOOS=windows", "GOARCH=arm64", "GOWORK=/workspace/go.work",
		"EXAMPLE=GOBIN=unchanged",
	}
	original := slices.Clone(environment)
	result := collectorBuildEnv(environment, "/tmp/build space")
	if !slices.Equal(environment, original) {
		t.Fatal("modified the inherited environment")
	}
	want := map[string]string{
		"PATH": "/usr/bin", "HOME": "/home/example", "GOCACHE": "/cache", "GOFLAGS": "-buildvcs=false",
		"GOBIN": "/tmp/build space", "CGO_ENABLED": "0", "GOOS": "linux", "GOARCH": "amd64", "GOWORK": "off",
		"EXAMPLE": "GOBIN=unchanged",
	}
	seen := make(map[string]bool)
	for _, entry := range result {
		key, value, ok := strings.Cut(entry, "=")
		if !ok || seen[key] {
			t.Fatalf("invalid or duplicate environment entry %q", entry)
		}
		if expected, ok := want[key]; !ok || value != expected {
			t.Fatalf("unexpected environment entry %q", entry)
		}
		seen[key] = true
	}
	if len(seen) != len(want) {
		t.Fatalf("lost environment entries: got %d, want %d", len(seen), len(want))
	}
}
