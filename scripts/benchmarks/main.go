// Command benchmarks maintains the published benchmark samples and comparisons.
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var errStale = errors.New("generated documentation or samples need updating")

type documentUpdate struct {
	Path, Content string
}

func main() {
	root, err := repositoryRoot()
	if err == nil {
		err = run(os.Args[1:], root, os.Stdout)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func repositoryRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		data, err := os.ReadFile(filepath.Join(dir, "go.mod"))
		if err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				if strings.TrimSpace(line) == "module github.com/rah-0/benchmarks" {
					return dir, nil
				}
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("run the benchmark tools from this repository or use scripts/results.bash")
		}
		dir = parent
	}
}

func run(args []string, root string, out io.Writer) error {
	usage := "usage: scripts/results.bash render [--check] [category ...] | sanitize [--check] <directory> | manifest <file>"
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		if len(args) == 0 {
			return errors.New(usage)
		}
		_, err := fmt.Fprintln(out, usage)
		return err
	}
	command := args[0]
	check := false
	var operands []string
	for _, argument := range args[1:] {
		if argument == "--check" {
			check = true
		} else if strings.HasPrefix(argument, "-") {
			return fmt.Errorf("unknown option %q", argument)
		} else {
			operands = append(operands, argument)
		}
	}
	switch command {
	case "render":
		return renderDocuments(root, operands, check, out)
	case "sanitize":
		if len(operands) != 1 {
			return errors.New(usage)
		}
		changed, err := sanitizeDirectory(operands[0], check)
		if err != nil {
			return err
		}
		count := 0
		for _, suite := range suites {
			count += len(suite.Packages)
		}
		_, err = fmt.Fprintf(out, "Validated %d Go result files; sanitized %d.\n", count, changed)
		return err
	case "manifest":
		if len(operands) != 1 || check {
			return errors.New(usage)
		}
		return writeManifest(root, operands[0])
	default:
		return fmt.Errorf("unknown command %q; %s", command, usage)
	}
}

// Validate every selected report before replacing any documentation.
func renderDocuments(root string, categories []string, check bool, out io.Writer) error {
	if len(categories) == 0 {
		for _, suite := range suites {
			categories = append(categories, suite.Name)
		}
	}
	for _, name := range []string{"environment.txt", "modules.txt"} {
		if _, err := os.Stat(filepath.Join(root, "results", name)); err != nil {
			return fmt.Errorf("missing results/%s: %w", name, err)
		}
	}
	var updates []documentUpdate
	seen := make(map[string]bool)
	for _, category := range categories {
		suite, ok := findSuite(category)
		if !ok || seen[category] {
			return fmt.Errorf("unknown or duplicate category %q", category)
		}
		seen[category] = true
		path := filepath.Join(root, filepath.FromSlash(suite.Path))
		original, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var rendered string
		if category == "collectors" {
			if _, err := os.Stat(filepath.Join(root, "results", "collectors-modules.txt")); err != nil {
				return fmt.Errorf("missing results/collectors-modules.txt: %w", err)
			}
			rendered, err = renderCollectors(root)
		} else {
			rendered, err = renderCategory(root, category)
		}
		if err != nil {
			return err
		}
		updated, err := replaceResults(string(original), rendered)
		if err != nil {
			return fmt.Errorf("%s: %w", suite.Path, err)
		}
		if updated != string(original) {
			updates = append(updates, documentUpdate{Path: suite.Path, Content: updated})
		}
	}
	for _, update := range updates {
		verb := "stale"
		if !check {
			if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(update.Path)), []byte(update.Content), 0o644); err != nil {
				return err
			}
			verb = "updated"
		}
		if _, err := fmt.Fprintf(out, "%s: %s\n", verb, update.Path); err != nil {
			return err
		}
	}
	if check && len(updates) > 0 {
		return errStale
	}
	return nil
}
