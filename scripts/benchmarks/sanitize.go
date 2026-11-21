package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

var (
	sanitizerBenchmark  = regexp.MustCompile(`^Benchmark\S+[ \t]+\d+[ \t]+.+$`)
	sanitizerCompletion = regexp.MustCompile(`^ok[ \t]+\S+[ \t]+\d+(?:\.\d+)?s[ \t]*$`)
	sanitizerVersion    = regexp.MustCompile(`^v\d+\.\d+\.\d+(?:-[0-9A-Za-z.-]+)?(?:\+[0-9A-Za-z.-]+)?$`)
	sanitizerModulePath = regexp.MustCompile(`^[0-9A-Za-z._~/-]+$`)
	moduleResultNames   = []string{"modules.txt", "collectors-modules.txt"}
)

// goResultNames follows the closed publication matrix, including nested packages.
func goResultNames() []string {
	var names []string
	for _, suite := range suites {
		for _, pkg := range suite.Packages {
			names = append(names, strings.ReplaceAll(pkg, "/", "-")+".txt")
		}
	}
	return names
}

// sanitizedGoOutput validates the unfiltered run and retains measurement bytes.
// Other output is private diagnostic data, not a published benchmark result.
func sanitizedGoOutput(path string) ([]byte, error) {
	original, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if !utf8.Valid(original) {
		return nil, fmt.Errorf("%s: invalid or failed benchmark output", filepath.Base(path))
	}
	if _, err := parseGoOutput(path); err != nil {
		// Parser errors may quote the offending line. Keep that content private.
		return nil, fmt.Errorf("%s: invalid or failed benchmark output", filepath.Base(path))
	}
	var sanitized bytes.Buffer
	passes, completions := 0, 0
	for raw := range bytes.SplitAfterSeq(original, []byte("\n")) {
		line := bytes.TrimRight(raw, "\r\n")
		pass := bytes.Equal(line, []byte("PASS"))
		complete := sanitizerCompletion.Match(line)
		if pass {
			passes++
		}
		if complete {
			completions++
		}
		header := bytes.HasPrefix(line, []byte("goos: ")) || bytes.HasPrefix(line, []byte("goarch: ")) ||
			bytes.HasPrefix(line, []byte("pkg: ")) || bytes.HasPrefix(line, []byte("cpu: "))
		if header || sanitizerBenchmark.Match(line) || pass || complete {
			sanitized.Write(raw)
		}
	}
	if passes != 1 || completions != 1 {
		return nil, fmt.Errorf("%s: missing or ambiguous successful completion", filepath.Base(path))
	}
	return sanitized.Bytes(), nil
}

func validModulePath(value string) bool {
	if strings.HasPrefix(value, "~") || !sanitizerModulePath.MatchString(value) {
		return false
	}
	for part := range strings.SplitSeq(value, "/") {
		if part == "" || part == "." || part == ".." {
			return false
		}
	}
	return true
}

// validateModules rejects private filesystem replacements instead of erasing
// their paths and claiming a reproducible public dependency snapshot.
func validateModules(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	name := filepath.Base(path)
	if len(data) == 0 {
		return fmt.Errorf("%s: empty module snapshot", name)
	}
	if !utf8.Valid(data) {
		return fmt.Errorf("%s: invalid module snapshot text", name)
	}
	lines := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	for index, line := range lines {
		fields := strings.Fields(line)
		source := fields
		if replacementAt := slices.Index(fields, "=>"); replacementAt >= 0 {
			source = fields[:replacementAt]
			replacement := fields[replacementAt+1:]
			if len(replacement) != 2 || !validModulePath(replacement[0]) || !sanitizerVersion.MatchString(replacement[1]) {
				return fmt.Errorf("%s:%d: local or invalid module replacement; use a published module version", name, index+1)
			}
		}
		if (len(source) != 1 && len(source) != 2) || !validModulePath(source[0]) ||
			(len(source) == 2 && !sanitizerVersion.MatchString(source[1])) {
			return fmt.Errorf("%s:%d: invalid module snapshot entry", name, index+1)
		}
	}
	return nil
}

type sanitizedFile struct {
	Path string
	Data []byte
}

// sanitizeDirectory is an output allowlist, not a general-purpose secret
// scanner. Environment and collector metadata have dedicated exporters.
// Every input is validated before any file is rewritten.
func sanitizeDirectory(directory string, check bool) (int, error) {
	goNames := goResultNames()
	allowed := make(map[string]bool, len(goNames)+len(moduleResultNames)+4)
	for _, name := range goNames {
		allowed[name] = true
	}
	for _, name := range moduleResultNames {
		allowed[name] = true
	}
	for _, name := range []string{"environment.txt", "healthy.jsonl", "outage.jsonl", "source.sha256"} {
		allowed[name] = true
	}
	entries, err := os.ReadDir(directory)
	if err != nil {
		return 0, err
	}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return 0, err
		}
		if !allowed[entry.Name()] || !info.Mode().IsRegular() {
			return 0, fmt.Errorf("unexpected artifact in results directory")
		}
	}
	for _, name := range moduleResultNames {
		if err := validateModules(filepath.Join(directory, name)); err != nil {
			return 0, err
		}
	}
	var updates []sanitizedFile
	for _, name := range goNames {
		path := filepath.Join(directory, name)
		sanitized, err := sanitizedGoOutput(path)
		if err != nil {
			return 0, err
		}
		original, err := os.ReadFile(path)
		if err != nil {
			return 0, err
		}
		if !bytes.Equal(original, sanitized) {
			updates = append(updates, sanitizedFile{Path: path, Data: sanitized})
		}
	}
	if check && len(updates) > 0 {
		names := make([]string, len(updates))
		for index, update := range updates {
			names[index] = filepath.Base(update.Path)
		}
		return 0, fmt.Errorf("benchmark diagnostics remain in: %s", strings.Join(names, ", "))
	}
	for _, update := range updates {
		if err := os.WriteFile(update.Path, update.Data, 0o600); err != nil {
			return 0, err
		}
	}
	return len(updates), nil
}
