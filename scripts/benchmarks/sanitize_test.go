package main

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const sanitizerSample = "goos: linux\n" +
	"goarch: amd64\n" +
	"pkg: example.com/benchmarks/cache\n" +
	"cpu: Example CPU\n" +
	"BenchmarkExample/1024-8 \t 123456\t 12.345 ns/op\t 64 B/op\t 2 allocs/op\n" +
	"PASS\n" +
	"ok  \texample.com/benchmarks/cache\t1.234s\n"

const sanitizerDiagnostics = "Docker Host: unix:///home/example/.docker/run/docker.sock\n" +
	"diagnostic token=fixture-only-do-not-publish\n" +
	"server listening on 192.0.2.10:12345\n"

func writeSanitizeFixture(t *testing.T, directory, name, content string) string {
	t.Helper()
	path := filepath.Join(directory, name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func readSanitizeFixture(t *testing.T, path string) []byte {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return content
}

func writeSanitizeMatrix(t *testing.T, directory, content string) {
	t.Helper()
	for _, name := range goResultNames() {
		writeSanitizeFixture(t, directory, name, content)
	}
	for _, name := range moduleResultNames {
		writeSanitizeFixture(t, directory, name, "example.com/benchmarks\nexample.com/library v1.2.3\n")
	}
}

func snapshotSanitizeFiles(t *testing.T, directory string) map[string]string {
	t.Helper()
	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatal(err)
	}
	result := make(map[string]string, len(entries))
	for _, entry := range entries {
		result[entry.Name()] = string(readSanitizeFixture(t, filepath.Join(directory, entry.Name())))
	}
	return result
}

func TestSanitizeDiagnosticsPreservesMeasurementBytes(t *testing.T) {
	original := sanitizerDiagnostics + strings.ReplaceAll(sanitizerSample, "PASS\n", sanitizerDiagnostics+"PASS\n") + sanitizerDiagnostics
	path := writeSanitizeFixture(t, t.TempDir(), "cache.txt", original)
	got, err := sanitizedGoOutput(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != sanitizerSample {
		t.Fatalf("sanitization changed measurement bytes: %q", got)
	}
	if string(readSanitizeFixture(t, path)) != original {
		t.Fatal("filtering unexpectedly modified the original file")
	}
}

func TestSanitizeRejectsFailedAndIncompleteRuns(t *testing.T) {
	directory := t.TempDir()
	for _, content := range []string{
		sanitizerSample + "FAIL\n",
		sanitizerSample + "--- FAIL: TestExample (0.01s)\n",
		sanitizerSample + "panic: example failure\n",
		sanitizerSample + "fatal error: example failure\n",
		strings.ReplaceAll(sanitizerSample, "PASS\n", ""),
		strings.Split(sanitizerSample, "ok  \t")[0],
		sanitizerSample + "PASS\n",
		sanitizerSample + string([]byte{0xff}),
	} {
		path := writeSanitizeFixture(t, directory, "cache.txt", content)
		if _, err := sanitizedGoOutput(path); err == nil {
			t.Error("accepted failed, incomplete, duplicated, or invalid output")
		}
	}
}

func TestSanitizePreservesCRLFAndMetricFormatting(t *testing.T) {
	sample := strings.ReplaceAll(sanitizerSample, "\n", "\r\n")
	path := writeSanitizeFixture(t, t.TempDir(), "cache.txt", sanitizerDiagnostics+sample)
	got, err := sanitizedGoOutput(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != sample {
		t.Fatalf("line endings or metric bytes changed: %q", got)
	}
}

func TestSanitizeAcceptsPublishedModuleReplacements(t *testing.T) {
	path := writeSanitizeFixture(t, t.TempDir(), "modules.txt",
		"example.com/benchmarks\n"+
			"example.com/library v1.2.3 => example.org/fork v1.2.4\n"+
			"example.com/pseudo v0.0.0-20260101000000-abcdefabcdef\n"+
			"example.com/legacy v2.0.0+incompatible\n")
	if err := validateModules(path); err != nil {
		t.Fatal(err)
	}
}

func TestSanitizeRejectsLocalModuleReplacements(t *testing.T) {
	directory := t.TempDir()
	for _, replacement := range []string{
		"../library", "./library", "/home/example/library", "~/library",
		`C:\Users\example\library`, "C:/Users/example/library",
		"../library v1.2.3", "/home/example/library v1.2.3", "~/library v1.2.3",
		"example.org/fork", "file:///home/example/library",
	} {
		t.Run(replacement, func(t *testing.T) {
			path := writeSanitizeFixture(t, directory, "modules.txt", "example.com/library v1.2.3 => "+replacement+"\n")
			if err := validateModules(path); err == nil {
				t.Fatal("accepted a local or unversioned replacement")
			}
		})
	}
}

func TestSanitizeValidatesAllInputsBeforeWriting(t *testing.T) {
	directory := t.TempDir()
	writeSanitizeMatrix(t, directory, sanitizerDiagnostics+sanitizerSample)
	names := goResultNames()
	writeSanitizeFixture(t, directory, names[len(names)-1], sanitizerSample+"FAIL\n")
	before := snapshotSanitizeFiles(t, directory)
	if _, err := sanitizeDirectory(directory, false); err == nil {
		t.Fatal("accepted a failed matrix")
	}
	if !reflect.DeepEqual(before, snapshotSanitizeFiles(t, directory)) {
		t.Fatal("modified files before completing validation")
	}
}

func TestSanitizeCheckIsReadOnlyAndIdempotent(t *testing.T) {
	directory := t.TempDir()
	writeSanitizeMatrix(t, directory, sanitizerDiagnostics+sanitizerSample)
	before := snapshotSanitizeFiles(t, directory)
	if _, err := sanitizeDirectory(directory, true); err == nil {
		t.Fatal("check accepted diagnostic output")
	}
	if !reflect.DeepEqual(before, snapshotSanitizeFiles(t, directory)) {
		t.Fatal("check modified the matrix")
	}
	changed, err := sanitizeDirectory(directory, false)
	if err != nil || changed != len(goResultNames()) {
		t.Fatalf("first sanitization changed %d files, error %v", changed, err)
	}
	for _, check := range []bool{false, true} {
		changed, err := sanitizeDirectory(directory, check)
		if err != nil || changed != 0 {
			t.Fatalf("sanitization was not idempotent: changed %d, error %v", changed, err)
		}
	}
	for _, name := range goResultNames() {
		if !bytes.Equal(readSanitizeFixture(t, filepath.Join(directory, name)), []byte(sanitizerSample)) {
			t.Fatalf("wrong measurement output in %s", name)
		}
	}
}

func TestSanitizeRejectsUnexpectedArtifacts(t *testing.T) {
	directory := t.TempDir()
	writeSanitizeMatrix(t, directory, sanitizerSample)
	for _, name := range []string{"debug.txt", "debug.log", ".env"} {
		t.Run(name, func(t *testing.T) {
			path := writeSanitizeFixture(t, directory, name, sanitizerDiagnostics)
			if _, err := sanitizeDirectory(directory, false); err == nil {
				t.Error("accepted an unrecognized artifact")
			}
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
		})
	}
	if err := os.Mkdir(filepath.Join(directory, "diagnostics"), 0o700); err != nil {
		t.Fatal(err)
	}
	if _, err := sanitizeDirectory(directory, false); err == nil {
		t.Fatal("accepted an unexpected directory")
	}
}

func TestSanitizeRejectsSymlinksAndNonFiles(t *testing.T) {
	directory := t.TempDir()
	writeSanitizeMatrix(t, directory, sanitizerSample)
	names := goResultNames()
	path := filepath.Join(directory, names[0])
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(names[1], path); err != nil {
		t.Fatal(err)
	}
	if _, err := sanitizeDirectory(directory, false); err == nil {
		t.Error("accepted a symlink to a known artifact")
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(path, 0o700); err != nil {
		t.Fatal(err)
	}
	if _, err := sanitizeDirectory(directory, false); err == nil {
		t.Fatal("accepted a directory in place of a known artifact")
	}
}
