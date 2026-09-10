package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const slogxModule = "github.com/rah-0/slogx-collector@v0.6.0"

// prepareSlogx builds the pinned module outside the measured interval. The
// caller copies the binary into its runtime container and releases it afterward.
func prepareSlogx(ctx context.Context) (path, digest string, cleanup func(), err error) {
	dir, err := os.MkdirTemp("", "collector-benchmark-build-")
	if err != nil {
		return "", "", nil, fmt.Errorf("create collector build directory: %w", err)
	}
	cleanup = func() { _ = os.RemoveAll(dir) }
	defer func() {
		if err != nil {
			cleanup()
		}
	}()

	buildCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	command := exec.CommandContext(buildCtx, "go", "install", "-trimpath", slogxModule)
	command.Dir = dir
	command.Env = collectorBuildEnv(os.Environ(), dir)
	var output diagnosticTail
	command.Stdout, command.Stderr = &output, &output
	if err := command.Run(); err != nil {
		if buildCtx.Err() != nil {
			err = buildCtx.Err()
		}
		return "", "", cleanup, fmt.Errorf("build %s: %w\n%s", slogxModule, err, output.String())
	}
	path = filepath.Join(dir, "slogx-collector")
	binary, err := os.ReadFile(path)
	if err != nil {
		return "", "", cleanup, fmt.Errorf("read built collector: %w", err)
	}
	sum := sha256.Sum256(binary)
	return path, fmt.Sprintf("%x", sum), cleanup, nil
}

func collectorBuildEnv(environment []string, binDir string) []string {
	var result []string
	for _, entry := range environment {
		key, _, _ := strings.Cut(entry, "=")
		switch key {
		case "GOBIN", "CGO_ENABLED", "GOOS", "GOARCH", "GOWORK":
			continue
		}
		result = append(result, entry)
	}
	return append(result, "GOBIN="+binDir, "CGO_ENABLED=0", "GOOS=linux", "GOARCH=amd64", "GOWORK=off")
}
