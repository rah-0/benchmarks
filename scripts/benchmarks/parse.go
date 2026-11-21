package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	failedRun    = regexp.MustCompile(`(?m)^(?:FAIL\b|--- FAIL:|panic:|fatal error:)`)
	completedRun = regexp.MustCompile(`(?m)^ok\s+\S+\s+`)
	benchmarkRow = regexp.MustCompile(`^(Benchmark\S+)\s+(\d+)\s+(.+)$`)
	cpuSuffix    = regexp.MustCompile(`^(Benchmark.*?)(?:-(\d+))?$`)
)

// parseGoOutput validates the original run before presentation or sanitization.
func parseGoOutput(path string) (map[string][]map[string]float64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if !utf8.Valid(data) {
		return nil, fmt.Errorf("%s: benchmark output is not UTF-8", filepath.Base(path))
	}
	text := string(data)
	if failedRun.MatchString(text) {
		return nil, fmt.Errorf("%s: failed benchmark run", filepath.Base(path))
	}
	if !completedRun.MatchString(text) {
		return nil, fmt.Errorf("%s: missing successful package completion", filepath.Base(path))
	}
	results := make(map[string][]map[string]float64)
	for _, line := range strings.Split(text, "\n") {
		match := benchmarkRow.FindStringSubmatch(strings.TrimSuffix(line, "\r"))
		if match == nil {
			continue
		}
		iterations, err := strconv.ParseUint(match[2], 10, 64)
		fields := strings.Fields(match[3])
		if err != nil || iterations == 0 || len(fields)%2 != 0 {
			return nil, fmt.Errorf("%s: invalid benchmark row", filepath.Base(path))
		}
		sample := make(map[string]float64)
		for i := 0; i < len(fields); i += 2 {
			value, err := strconv.ParseFloat(fields[i], 64)
			_, duplicate := sample[fields[i+1]]
			if err != nil || math.IsNaN(value) || math.IsInf(value, 0) || value < 0 || duplicate {
				return nil, fmt.Errorf("%s: invalid benchmark metric", filepath.Base(path))
			}
			sample[fields[i+1]] = value
		}
		if _, ok := sample["ns/op"]; !ok {
			return nil, fmt.Errorf("%s: missing ns/op", filepath.Base(path))
		}
		if previous := results[match[1]]; len(previous) > 0 && !sameMetrics(previous[0], sample) {
			return nil, fmt.Errorf("%s: inconsistent benchmark metrics", filepath.Base(path))
		}
		results[match[1]] = append(results[match[1]], sample)
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("%s: no benchmark measurements", filepath.Base(path))
	}
	return results, nil
}

func sameMetrics(a, b map[string]float64) bool {
	if len(a) != len(b) {
		return false
	}
	for key := range a {
		if _, ok := b[key]; !ok {
			return false
		}
	}
	return true
}

func loadCategory(root, category string) (map[measurementKey]measurement, environment, error) {
	suite, ok := findSuite(category)
	if !ok || len(suite.Packages) == 0 {
		return nil, environment{}, fmt.Errorf("unknown Go category %q", category)
	}
	measurements := make(map[measurementKey]measurement)
	var common environment
	for index, pkg := range suite.Packages {
		path := filepath.Join(root, "results", strings.ReplaceAll(pkg, "/", "-")+".txt")
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, environment{}, err
		}
		headers := make(map[string]string)
		for _, field := range []string{"goos", "goarch", "pkg", "cpu"} {
			pattern := regexp.MustCompile(`(?m)^` + field + `: (.+)$`)
			matches := pattern.FindAllStringSubmatch(string(data), -1)
			if len(matches) != 1 {
				return nil, environment{}, fmt.Errorf("%s: expected one %s header", filepath.Base(path), field)
			}
			headers[field] = strings.TrimSpace(matches[0][1])
		}
		if headers["pkg"] != "github.com/rah-0/benchmarks/"+pkg {
			return nil, environment{}, fmt.Errorf("%s: wrong benchmark package", filepath.Base(path))
		}
		current := environment{OS: headers["goos"], Arch: headers["goarch"], CPU: headers["cpu"]}
		if index > 0 && current != common {
			return nil, environment{}, fmt.Errorf("%s: packages have different measurement environments", category)
		}
		common = current
		parsed, err := parseGoOutput(path)
		if err != nil {
			return nil, environment{}, err
		}
		for name, samples := range parsed {
			match := cpuSuffix.FindStringSubmatch(name)
			key := measurementKey{Package: pkg, Name: match[1]}
			if _, exists := measurements[key]; exists {
				return nil, environment{}, fmt.Errorf("%s: multiple CPU settings for %s; publish separate runs", filepath.Base(path), match[1])
			}
			if !sameMetrics(samples[0], map[string]float64{"ns/op": 0, "B/op": 0, "allocs/op": 0}) {
				return nil, environment{}, fmt.Errorf("%s: new metrics need explicit presentation rules", filepath.Base(path))
			}
			cpu := 1
			if match[2] != "" {
				cpu, err = strconv.Atoi(match[2])
				if err != nil || cpu < 1 {
					return nil, environment{}, fmt.Errorf("%s: invalid CPU setting", filepath.Base(path))
				}
			}
			measurements[key] = measurement{Name: name, CPU: cpu, Samples: samples}
		}
	}
	return measurements, common, nil
}
