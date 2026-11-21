package main

import (
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// The manifest identifies measured code and configuration, excluding reporting
// tools and their tests. It uses repository-relative paths on every platform.
func writeManifest(root, destination string) error {
	var paths []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if path != root {
				switch entry.Name() {
				case ".git", ".idea", "pprof_svg", "results":
					return filepath.SkipDir
				}
				if path == filepath.Join(root, "scripts") {
					return filepath.SkipDir
				}
			}
			return nil
		}
		if !entry.Type().IsRegular() {
			return nil
		}
		extension := filepath.Ext(path)
		if extension == ".go" || extension == ".toml" || extension == ".conf" || entry.Name() == "go.mod" || entry.Name() == "go.sum" {
			relative, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			paths = append(paths, filepath.ToSlash(relative))
		}
		return nil
	})
	if err != nil {
		return err
	}
	slices.SortFunc(paths, func(a, b string) int {
		return slices.Compare(strings.Split(a, "/"), strings.Split(b, "/"))
	})
	var result strings.Builder
	for _, relative := range paths {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(relative)))
		if err != nil {
			return err
		}
		fmt.Fprintf(&result, "%x  %s\n", sha256.Sum256(data), relative)
	}
	return os.WriteFile(destination, []byte(result.String()), 0o600)
}
