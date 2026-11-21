package main

import (
	"bytes"
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type collectorFieldTest struct {
	Name  string
	Field string
	Value any
}

type collectorRatioTest struct {
	Numerator   float64
	Denominator float64
	Want        string
}

func collectorFixture(outage float64) []map[string]any {
	images := make(map[string]string)
	for i, name := range collectorNames {
		images[name] = "example/" + name + ":1.0.0@sha256:" + strings.Repeat(string(rune('1'+i)), 64)
	}
	names := []string{"http-control", "slogx-collector", "vector", "fluent-bit"}
	env := map[string]any{
		"type": "environment", "date": "2026-01-01T00:00:00Z",
		"go": "go1.27.4", "os": "linux", "arch": "amd64", "collector_cpu": "0",
		"records": 100, "sizes": []int{256}, "repeat": 3, "outage_seconds": outage,
		"images": images, "collectors": names,
		"container_manager":   "testcontainers-go v0.45.0",
		"slogx_module":        "github.com/rah-0/slogx-collector@v0.7.0",
		"slogx_binary_sha256": strings.Repeat("a", 64), "slogx_go": "go1.27.4",
		"kernel": "example-kernel", "harness_cpus": "2", "cpu": "Example CPU",
		"storage_filesystem": "ext4",
	}
	rates := map[string]float64{"fluent-bit": 100, "slogx-collector": 200, "vector": 400, "http-control": 1000}
	recovery := map[string]float64{"fluent-bit": 2, "slogx-collector": 0.25, "vector": 0.125, "http-control": 0.05}
	rss := map[string]int64{"fluent-bit": 16, "slogx-collector": 8, "vector": 64}
	entries := []map[string]any{env}
	for repetition, factor := range []float64{0.9, 1, 1.1} {
		for _, name := range names {
			seconds := 100 / rates[name] * factor
			if outage != 0 {
				seconds = outage + recovery[name]*factor
			}
			sample := map[string]any{
				"collector": name, "sample": repetition + 1, "records": 100, "record_bytes": 256,
				"input_bytes": 25600, "seconds": seconds, "records_per_second": 100 / seconds,
				"input_mib_per_second": 25600.0 / 1048576 / seconds, "producer_seconds": 0.01,
				"http_requests": 4, "retry_requests": 0, "duplicates": 0,
				"outage_seconds": outage,
			}
			if name != "http-control" {
				sample["image"] = images[name]
				sample["image_id"] = strings.SplitN(images[name], "@", 2)[1]
				sample["cpu_seconds"] = 0.01
				sample["peak_rss_bytes"] = rss[name] * 1048576
			}
			if outage != 0 {
				sample["recovery_drain_seconds"] = recovery[name] * factor
				sample["retry_requests"] = 2
			}
			entries = append(entries, sample)
		}
	}
	return append(entries, map[string]any{"type": "complete", "samples": len(entries) - 1})
}

func writeCollectorFixture(t *testing.T, root, name string, entries []map[string]any) string {
	t.Helper()
	var data bytes.Buffer
	encoder := json.NewEncoder(&data)
	for _, entry := range entries {
		if err := encoder.Encode(entry); err != nil {
			t.Fatal(err)
		}
	}
	path := filepath.Join(root, "results", name+".jsonl")
	if err := os.WriteFile(path, data.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func collectorFixtureRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "results"), 0o700); err != nil {
		t.Fatal(err)
	}
	writeCollectorFixture(t, root, "healthy", collectorFixture(0))
	writeCollectorFixture(t, root, "outage", collectorFixture(2))
	return root
}

func TestCollectorBaselineAndMetricDirections(t *testing.T) {
	text, err := renderCollectors(collectorFixtureRoot(t))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"| **Fluent Bit — baseline** | 100 | 1.00× | 16.00 | 1.00× |",
		"| slogx-collector | 200 | 2.00× | 8.00 | 0.50× |",
		"| Vector | 400 | 4.00× | 64.00 | 4.00× |",
		"| slogx-collector | 0.250 | 8.00× |",
		"| Vector | 0.125 | 16.00× |",
		"Observed ranges", "Persistence differs", "excluded from collector speed ratios",
		"[Healthy samples](../results/healthy.jsonl)", "../README.md#reading-the-comparisons",
		"Vector has the lowest measured median recovery time.",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("missing expected presentation %q", want)
		}
	}
	if strings.Index(text, "| **Fluent Bit — baseline**") >= strings.Index(text, "| slogx-collector") {
		t.Error("baseline is not the first collector")
	}
	if strings.Count(text, startMarker) != 1 || strings.Count(text, endMarker) != 1 {
		t.Error("result markers must occur exactly once")
	}
}

func TestCollectorRenderingPreservesSamples(t *testing.T) {
	root := collectorFixtureRoot(t)
	before := make(map[string][]byte)
	for _, name := range []string{"healthy", "outage"} {
		path := filepath.Join(root, "results", name+".jsonl")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		before[path] = data
	}
	if _, err := renderCollectors(root); err != nil {
		t.Fatal(err)
	}
	for path, want := range before {
		got, err := os.ReadFile(path)
		if err != nil || !bytes.Equal(got, want) {
			t.Error("rendering changed a sample file")
		}
	}
}

func TestCollectorOverlappingRecoveryRanges(t *testing.T) {
	root := collectorFixtureRoot(t)
	entries := collectorFixture(2)
	for _, sample := range entries[1 : len(entries)-1] {
		if sample["collector"] == "vector" && sample["sample"] == 3 {
			sample["recovery_drain_seconds"] = 0.3
			sample["seconds"] = 2.3
			sample["records_per_second"] = 100.0 / 2.3
			sample["input_mib_per_second"] = 25600.0 / 1048576 / 2.3
		}
	}
	writeCollectorFixture(t, root, "outage", entries)
	text, err := renderCollectors(root)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(text, "Its observed range overlaps slogx-collector's") || !strings.Contains(text, "do not establish a consistent ordering between them") {
		t.Error("overlapping recovery ranges did not qualify their ranking")
	}
}

func TestCollectorRejectsIncompleteDuplicateAndUnknownMatrices(t *testing.T) {
	root := collectorFixtureRoot(t)
	incomplete := collectorFixture(0)
	incomplete = append(incomplete[:len(incomplete)-2], incomplete[len(incomplete)-1])
	duplicate := collectorFixture(0)
	duplicate[1] = duplicate[2]
	unknown := collectorFixture(0)
	unknown[0]["collectors"].([]string)[0] = "unmapped-collector"
	for i, entries := range [][]map[string]any{incomplete, duplicate, unknown} {
		if _, err := parseCollectors(writeCollectorFixture(t, root, "invalid", entries)); err == nil {
			t.Errorf("accepted invalid matrix %d", i)
		}
	}
}

func TestCollectorRejectsInvalidSampleFields(t *testing.T) {
	for _, test := range []collectorFieldTest{
		{"boolean count", "sample", true}, {"negative count", "duplicates", -1},
		{"non-finite metric", "cpu_seconds", json.Number("1e309")},
		{"fractional RSS", "peak_rss_bytes", 1.5}, {"invalid collector type", "collector", []string{}},
		{"unknown metric", "unmapped_metric", 1}, {"null required field", "cpu_seconds", nil},
		{"missing zero-valued field", "cpu_seconds", "delete"},
	} {
		t.Run(test.Name, func(t *testing.T) {
			root := collectorFixtureRoot(t)
			entries := collectorFixture(0)
			if action, ok := test.Value.(string); ok && action == "delete" {
				delete(entries[2], test.Field)
			} else {
				entries[2][test.Field] = test.Value
			}
			if _, err := parseCollectors(writeCollectorFixture(t, root, "invalid", entries)); err == nil {
				t.Fatal("accepted invalid sample field")
			}
		})
	}
}

func TestCollectorRejectsImageChanges(t *testing.T) {
	root := collectorFixtureRoot(t)
	differentReference := collectorFixture(0)
	differentReference[2]["image"] = "example/other:1.0.0@sha256:" + strings.Repeat("1", 64)
	differentRuntime := collectorFixture(0)
	differentRuntime[2]["image_id"] = "sha256:" + strings.Repeat("f", 64)
	for _, entries := range [][]map[string]any{differentReference, differentRuntime} {
		if _, err := parseCollectors(writeCollectorFixture(t, root, "invalid", entries)); err == nil {
			t.Error("accepted an inconsistent sample image")
		}
	}
	entries := collectorFixture(2)
	for _, sample := range entries[1 : len(entries)-1] {
		if sample["collector"] == "vector" {
			sample["image_id"] = "sha256:" + strings.Repeat("f", 64)
		}
	}
	writeCollectorFixture(t, root, "outage", entries)
	if _, err := renderCollectors(root); err == nil {
		t.Error("accepted different runtime images across healthy/outage runs")
	}
}

func TestCollectorRejectsDifferentEnvironments(t *testing.T) {
	root := collectorFixtureRoot(t)
	entries := collectorFixture(2)
	entries[0]["collector_cpu"] = "1"
	writeCollectorFixture(t, root, "outage", entries)
	if _, err := renderCollectors(root); err == nil {
		t.Error("accepted different environments")
	}
}

func TestCollectorRequiresOutageRetriesAndConsistentDuration(t *testing.T) {
	root := collectorFixtureRoot(t)
	for _, test := range []collectorFieldTest{{"missing retry", "retry_requests", 0}, {"wrong duration", "recovery_drain_seconds", 999}} {
		entries := collectorFixture(2)
		entries[1][test.Field] = test.Value
		if _, err := parseCollectors(writeCollectorFixture(t, root, "invalid", entries)); err == nil {
			t.Errorf("accepted %s", test.Name)
		}
	}
	entries := collectorFixture(0)
	entries[1]["retry_requests"] = 1
	if _, err := parseCollectors(writeCollectorFixture(t, root, "invalid", entries)); err == nil {
		t.Error("accepted HTTP503 retries in a healthy sample")
	}
}

func TestCollectorZeroRecoveryHasNoInfiniteSpeedup(t *testing.T) {
	root := collectorFixtureRoot(t)
	entries := collectorFixture(2)
	for _, sample := range entries[1 : len(entries)-1] {
		if sample["collector"] == "slogx-collector" {
			delete(sample, "recovery_drain_seconds")
			sample["seconds"] = 2
			sample["records_per_second"] = 50
			sample["input_mib_per_second"] = 25600.0 / 1048576 / 2
		}
	}
	writeCollectorFixture(t, root, "outage", entries)
	text, err := renderCollectors(root)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(text, "| slogx-collector | 0.000 | — |") {
		t.Error("zero recovery did not suppress the speed ratio")
	}
}

func TestCollectorThroughputMustMatchDuration(t *testing.T) {
	root := collectorFixtureRoot(t)
	entries := collectorFixture(0)
	entries[1]["records_per_second"] = entries[1]["records_per_second"].(float64) * 2
	if _, err := parseCollectors(writeCollectorFixture(t, root, "invalid", entries)); err == nil {
		t.Error("accepted throughput inconsistent with duration")
	}
	entries = collectorFixture(0)
	entries[1]["seconds"] = math.SmallestNonzeroFloat64
	if _, err := parseCollectors(writeCollectorFixture(t, root, "invalid", entries)); err == nil {
		t.Error("accepted a duration that overflows derived throughput")
	}
}

func TestCollectorDuplicateTotalCannotOverflow(t *testing.T) {
	root := collectorFixtureRoot(t)
	entries := collectorFixture(0)
	entries[1]["duplicates"], entries[2]["duplicates"] = int64(math.MaxInt64), int64(1)
	writeCollectorFixture(t, root, "healthy", entries)
	if _, err := renderCollectors(root); err == nil {
		t.Error("accepted a duplicate total that overflows its count")
	}
}

func TestCollectorTinyAndUndefinedRatios(t *testing.T) {
	for _, test := range []collectorRatioTest{{1, 10000, "<0.01×"}, {0, 10000, "0.00×"}, {10000, 1, "10,000.00×"}, {1, 0, "—"}} {
		if got := collectorRatio(test.Numerator, test.Denominator); got != test.Want {
			t.Errorf("ratio(%g, %g) = %q, want %q", test.Numerator, test.Denominator, got, test.Want)
		}
	}
	if collectorNumber(0.0001, 2) != "0.0001" {
		t.Error("a positive metric rounded to zero")
	}
}

func TestCollectorRejectsMalformedEnvironmentAndOverflowCounts(t *testing.T) {
	for _, test := range []collectorFieldTest{
		{"unknown field", "unexpected", "value"}, {"markup", "cpu", "CPU | injected column"},
		{"null count", "records", nil}, {"duplicate collectors", "collectors", []string{"vector", "vector", "fluent-bit", "http-control"}},
		{"duplicate sizes", "sizes", []int{256, 256}}, {"empty sizes", "sizes", []int{}},
		{"invalid digest", "slogx_binary_sha256", "invalid"}, {"large repeat", "repeat", math.MaxInt},
	} {
		t.Run(test.Name, func(t *testing.T) {
			root := collectorFixtureRoot(t)
			entries := collectorFixture(0)
			entries[0][test.Field] = test.Value
			if _, err := parseCollectors(writeCollectorFixture(t, root, "invalid", entries)); err == nil {
				t.Fatal("accepted malformed environment")
			}
		})
	}
}

func TestCollectorCurrentMarkdownMatchesRecordedData(t *testing.T) {
	root := filepath.Join("..", "..")
	readme, err := os.ReadFile(filepath.Join(root, "collectors", "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(readme)
	start, end := strings.Index(text, startMarker), strings.Index(text, endMarker)
	if start < 0 || end < start {
		t.Fatal("missing result markers")
	}
	want := text[start : end+len(endMarker)]
	got, err := renderCollectors(root)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		gotLines, wantLines := strings.Split(got, "\n"), strings.Split(want, "\n")
		for i := range min(len(gotLines), len(wantLines)) {
			if gotLines[i] != wantLines[i] {
				t.Fatalf("collector output differs at line %d:\ngot:  %s\nwant: %s", i+1, gotLines[i], wantLines[i])
			}
		}
		t.Fatalf("collector output differs in length: got %d lines, want %d", len(gotLines), len(wantLines))
	}
}
