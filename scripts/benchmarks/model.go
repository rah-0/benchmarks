package main

import (
	"slices"
)

const (
	startMarker = "<!-- benchmark-results:start -->"
	endMarker   = "<!-- benchmark-results:end -->"
)

type suiteDef struct {
	Name, Path string
	Packages   []string
}

var suites = []suiteDef{
	{"cache", "cache/README.md", []string{"cache"}},
	{"collectors", "collectors/README.md", nil},
	{"compression", "compression/README.MD", []string{"compression"}},
	{"crypto", "crypto/README.md", []string{"crypto/encrypt", "crypto/sign"}},
	{"db", "db/README.md", []string{"db/sqlite", "db/mariadb", "db/postgresql"}},
	{"encoding", "encoding/README.md", []string{"encoding"}},
	{"iteration", "iteration/README.md", []string{"iteration"}},
	{"log", "log/README.md", []string{"log"}},
	{"meta", "meta/README.MD", []string{"meta"}},
	{"protocol", "protocol/README.md", []string{"protocol/http", "protocol/tcp"}},
	{"serializer", "serializer/README.md", []string{"serializer"}},
	{"validator", "validator/README.md", []string{"validator"}},
}

type rowDef struct {
	Label, Package, Name string
}

type caseDef struct {
	Label string
	Rows  []rowDef
}

type groupDef struct {
	Title, Description, Baseline, Operation, Note string
	Rows                                          []rowDef
	Cases                                         []caseDef
}

type measurementKey struct {
	Package, Name string
}

type measurement struct {
	Name    string
	CPU     int
	Samples []map[string]float64
}

func (m measurement) median(metric string) float64 {
	values := make([]float64, len(m.Samples))
	for i, sample := range m.Samples {
		values[i] = sample[metric]
	}
	slices.Sort(values)
	middle := len(values) / 2
	if len(values)%2 == 0 {
		return (values[middle-1] + values[middle]) / 2
	}
	return values[middle]
}

func (m measurement) bounds(metric string) (float64, float64) {
	low, high := m.Samples[0][metric], m.Samples[0][metric]
	for _, sample := range m.Samples[1:] {
		low, high = min(low, sample[metric]), max(high, sample[metric])
	}
	return low, high
}

type environment struct {
	OS, Arch, CPU string
}

func findSuite(name string) (suiteDef, bool) {
	for _, suite := range suites {
		if suite.Name == name {
			return suite, true
		}
	}
	return suiteDef{}, false
}
