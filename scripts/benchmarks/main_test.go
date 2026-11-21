package main

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func cliFixture(t *testing.T) (string, string) {
	t.Helper()
	source, err := repositoryRoot()
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	for _, directory := range []string{"results", "compression"} {
		if err := os.Mkdir(filepath.Join(root, directory), 0o700); err != nil {
			t.Fatal(err)
		}
	}
	for _, name := range []string{"environment.txt", "modules.txt", "compression.txt"} {
		data, err := os.ReadFile(filepath.Join(source, "results", name))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, "results", name), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	original := "# Fixture\n\n" + startMarker + "\nOld results\n" + endMarker + "\n\nKeep this suffix.\n"
	if err := os.WriteFile(filepath.Join(root, "compression/README.MD"), []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}
	return root, original
}

func TestRenderCommandCheckAndIdempotence(t *testing.T) {
	root, original := cliFixture(t)
	path := filepath.Join(root, "compression/README.MD")
	var output bytes.Buffer
	if err := run([]string{"render", "compression", "--check"}, root, &output); !errors.Is(err, errStale) {
		t.Fatalf("check returned %v, want stale", err)
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != original {
		t.Fatal("check modified the document", err)
	}
	if !strings.Contains(output.String(), "stale: compression/README.MD") {
		t.Fatal("missing stale document report")
	}
	output.Reset()
	if err := run([]string{"render", "compression"}, root, &output); err != nil {
		t.Fatal(err)
	}
	first, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(first, []byte("# Fixture\n\n")) || !bytes.HasSuffix(first, []byte("\n\nKeep this suffix.\n")) {
		t.Fatal("renderer changed content outside the result markers")
	}
	output.Reset()
	if err := run([]string{"render", "--check", "compression"}, root, &output); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"render", "compression"}, root, &output); err != nil {
		t.Fatal(err)
	}
	second, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(first, second) || output.Len() != 0 {
		t.Fatal("generation is not idempotent", err)
	}
}

func TestRenderCommandValidatesAllBeforeWriting(t *testing.T) {
	root, original := cliFixture(t)
	var output bytes.Buffer
	if err := run([]string{"render", "compression", "unknown-category"}, root, &output); err == nil {
		t.Fatal("accepted an unmapped category")
	}
	data, err := os.ReadFile(filepath.Join(root, "compression/README.MD"))
	if err != nil || string(data) != original {
		t.Fatal("updated earlier reports before validating every category", err)
	}
}

func TestCommandArguments(t *testing.T) {
	var output bytes.Buffer
	for _, args := range [][]string{nil, {"unknown"}, {"render", "--unknown"}, {"sanitize"}, {"manifest", "--check", "file"}} {
		if err := run(args, t.TempDir(), &output); err == nil {
			t.Fatalf("accepted invalid arguments %q", args)
		}
	}
	if err := run([]string{"--help"}, "", &output); err != nil || !strings.Contains(output.String(), "usage:") {
		t.Fatal("help is unavailable", err)
	}
}

func TestManifestSelectsMeasuredFilesAndRelativePaths(t *testing.T) {
	root := t.TempDir()
	selected := map[string]string{
		"main.go": "package main\n", "go.mod": "module example.com/benchmark\n",
		"collectors/config/vector.toml": "fixture = true\n",
		"collectors/config.go":          "package collectors\n",
		"fixture/scripts/input.go":      "package fixture\n",
	}
	ignored := []string{"scripts/report.go", "results/sample.go", ".git/config.go", "pprof_svg/sample.go", "README.md"}
	for name, content := range selected {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	for _, name := range ignored {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("unmeasured fixture\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	destination := filepath.Join(root, "source.sha256")
	if err := writeManifest(root, destination); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(destination)
	if err != nil {
		t.Fatal(err)
	}
	var expected strings.Builder
	for _, name := range []string{"collectors/config/vector.toml", "collectors/config.go", "fixture/scripts/input.go", "go.mod", "main.go"} {
		fmt.Fprintf(&expected, "%x  %s\n", sha256.Sum256([]byte(selected[name])), name)
	}
	if string(data) != expected.String() || strings.Contains(string(data), root) {
		t.Fatal("manifest differs from the selected, sorted, relative-path source inventory")
	}
}
