package main

import (
	"maps"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"testing"
)

type rendererInvalidCase struct {
	Name         string
	Groups       []groupDef
	Measurements map[measurementKey]measurement
}

type rendererRatioCase struct {
	Name                   string
	Numerator, Denominator float64
	Want                   string
}

type rendererAnchorCase struct {
	Title, Want string
}

type rendererSample struct {
	Time, Allocated, Allocs float64
}

func renderTestRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func renderTestMeasurement(name string, times []float64, allocated float64, cpu int) measurement {
	samples := make([]map[string]float64, len(times))
	for index, value := range times {
		samples[index] = map[string]float64{"ns/op": value, "B/op": allocated, "allocs/op": 0}
	}
	return measurement{Name: name, CPU: cpu, Samples: samples}
}

func renderTestGroup() groupDef {
	return groupDef{
		Title: "One operation", Description: "Measure matching operations.", Baseline: "Reference",
		Rows: []rowDef{{Label: "Reference", Package: "fixture", Name: "BenchmarkReference"},
			{Label: "Candidate", Package: "fixture", Name: "BenchmarkCandidate"}},
	}
}

func renderTestValues() map[measurementKey]measurement {
	return map[measurementKey]measurement{
		{Package: "fixture", Name: "BenchmarkReference"}: renderTestMeasurement("BenchmarkReference", []float64{10, 20, 30}, 0, 1),
		{Package: "fixture", Name: "BenchmarkCandidate"}: renderTestMeasurement("BenchmarkCandidate", []float64{5, 10, 15}, 0, 1),
	}
}

func TestRenderedGoReportsPreserveRecordedPresentation(t *testing.T) {
	root := renderTestRoot(t)
	for _, suite := range suites {
		if len(suite.Packages) == 0 {
			continue
		}
		t.Run(suite.Name, func(t *testing.T) {
			original, err := os.ReadFile(filepath.Join(root, suite.Path))
			if err != nil {
				t.Fatal(err)
			}
			_, rest, found := strings.Cut(string(original), startMarker)
			if !found {
				t.Fatal("missing opening results marker")
			}
			body, _, found := strings.Cut(rest, endMarker)
			if !found {
				t.Fatal("missing closing results marker")
			}
			want := startMarker + body + endMarker
			got, err := renderCategory(root, suite.Name)
			if err != nil {
				t.Fatal(err)
			}
			if want != got {
				wantLines, gotLines := strings.Split(want, "\n"), strings.Split(got, "\n")
				for index := range min(len(wantLines), len(gotLines)) {
					if wantLines[index] != gotLines[index] {
						t.Fatalf("presentation changed at line %d\nwant: %s\ngot:  %s", index+1, wantLines[index], gotLines[index])
					}
				}
				t.Fatalf("presentation line count changed: want %d, got %d", len(wantLines), len(gotLines))
			}
		})
	}
}

func TestPresentationCoversEveryRecordedCaseAndSample(t *testing.T) {
	root := renderTestRoot(t)
	paths, err := filepath.Glob(filepath.Join(root, "results", "*.txt"))
	if err != nil {
		t.Fatal(err)
	}
	packagePattern := regexp.MustCompile(`(?m)^pkg: github\.com/rah-0/benchmarks/(.+)$`)
	benchmarkPattern := regexp.MustCompile(`(?m)^(Benchmark\S+)\s+\d+\s+`)
	cpuSuffix := regexp.MustCompile(`-\d+$`)
	expected := make(map[measurementKey]int)
	for _, path := range paths {
		contents, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		pkg := packagePattern.FindStringSubmatch(string(contents))
		if pkg == nil {
			continue
		}
		for _, match := range benchmarkPattern.FindAllStringSubmatch(string(contents), -1) {
			expected[measurementKey{Package: pkg[1], Name: cpuSuffix.ReplaceAllString(match[1], "")}]++
		}
	}
	actual := make(map[measurementKey]int)
	for category, definitions := range groups {
		values, _, err := loadCategory(root, category)
		if err != nil {
			t.Fatal(err)
		}
		if err := validateGroups(category, values); err != nil {
			t.Fatal(err)
		}
		for key, value := range values {
			actual[key] += len(value.Samples)
		}
		for _, group := range definitions {
			if _, err := groupContext(group, values); err != nil {
				t.Fatal(err)
			}
		}
	}
	if len(expected) == 0 || !reflect.DeepEqual(actual, expected) {
		t.Fatal("presentation inventory differs from the raw benchmark cases or sample counts")
	}
	total := 0
	for _, count := range actual {
		total += count
	}
	t.Logf("covered %d benchmark cases and %d recorded samples", len(actual), total)
}

func TestUnmappedMissingAndAmbiguousComparisonsAreRejected(t *testing.T) {
	group, values := renderTestGroup(), renderTestValues()
	unknown := maps.Clone(values)
	unknown[measurementKey{Package: "fixture", Name: "BenchmarkUnknown"}] = renderTestMeasurement("BenchmarkUnknown", []float64{1, 2, 3}, 0, 1)
	missing := maps.Clone(values)
	delete(missing, measurementKey{Package: "fixture", Name: "BenchmarkCandidate"})
	missingBaseline := group
	missingBaseline.Baseline = "Absent"
	duplicateLabel := group
	duplicateLabel.Rows = []rowDef{group.Rows[0], {Label: "Reference", Package: "fixture", Name: "BenchmarkCandidate"}}
	duplicateValue := group
	duplicateValue.Rows = []rowDef{group.Rows[0], {Label: "Second label", Package: "fixture", Name: "BenchmarkReference"}}
	duplicateCase := group
	duplicateCase.Rows = nil
	duplicateCase.Cases = []caseDef{{Label: "Same input", Rows: group.Rows}, {Label: "Same input", Rows: group.Rows}}
	emptyCases := group
	emptyCases.Cases = []caseDef{}
	for _, test := range []rendererInvalidCase{
		{Name: "unmapped measurement", Groups: []groupDef{group}, Measurements: unknown},
		{Name: "missing measurement", Groups: []groupDef{group}, Measurements: missing},
		{Name: "missing baseline", Groups: []groupDef{missingBaseline}, Measurements: values},
		{Name: "duplicate baseline", Groups: []groupDef{duplicateLabel}, Measurements: values},
		{Name: "duplicate measurement", Groups: []groupDef{duplicateValue}, Measurements: values},
		{Name: "duplicate input", Groups: []groupDef{duplicateCase}, Measurements: values},
		{Name: "duplicate group", Groups: []groupDef{group, group}, Measurements: values},
		{Name: "empty sweep", Groups: []groupDef{emptyCases}, Measurements: values},
	} {
		t.Run(test.Name, func(t *testing.T) {
			if err := validateGroupDefinitions("fixture", test.Groups, test.Measurements); err == nil {
				t.Fatal("invalid comparison was accepted")
			}
		})
	}
}

func TestSpeedUsesRatioOfMedians(t *testing.T) {
	reference := renderTestMeasurement("BenchmarkReference", []float64{10, 100, 1000}, 32, 1)
	candidate := renderTestMeasurement("BenchmarkCandidate", []float64{20, 25, 100}, 64, 1)
	output := strings.Join(comparisonTable(renderTestGroup(), []resolvedRow{
		{Label: "Reference", Measurement: reference}, {Label: "Candidate", Measurement: candidate},
	}), "\n")
	if !strings.Contains(output, "| Candidate | 25 | 4.00× | 64 |") {
		t.Fatal("candidate speed must be reference median 100 divided by candidate median 25")
	}
	if got := allocatedComparison(64, 32); got != "2.00×" {
		t.Fatalf("allocation ratio points in the wrong direction: %s", got)
	}
}

func TestZeroAndSmallRatiosRemainMeaningful(t *testing.T) {
	for _, test := range []rendererRatioCase{
		{Name: "zero denominator", Numerator: 1, Denominator: 0, Want: "—"},
		{Name: "actual zero", Numerator: 0, Denominator: 1, Want: "0.00×"},
		{Name: "tiny allocation", Numerator: 314, Denominator: 41_741_273, Want: "<0.01×"},
		{Name: "tiny speed", Numerator: 0.067, Denominator: 544, Want: "<0.01×"},
		{Name: "large speed", Numerator: 9_978_787, Denominator: 328.9, Want: "30,339.88×"},
	} {
		t.Run(test.Name, func(t *testing.T) {
			if got := ratio(test.Numerator, test.Denominator); got != test.Want {
				t.Fatalf("want %s, got %s", test.Want, got)
			}
		})
	}
	if got := allocatedComparison(64, 0); got != "+64 B" {
		t.Fatalf("zero reference should use absolute allocation difference: %s", got)
	}
	if got := allocatedComparison(0, 0); got != "same (0 B)" {
		t.Fatalf("zero allocations should remain explicit: %s", got)
	}
}

func TestExtremeRangesKeepValuesReadable(t *testing.T) {
	reference := renderTestMeasurement("BenchmarkReference", []float64{10_000_000, 10_000_000, 10_000_000}, 64*1024*1024, 1)
	candidate := renderTestMeasurement("BenchmarkCandidate", []float64{500, 500, 500}, 64, 1)
	text := strings.Join(comparisonTable(renderTestGroup(), []resolvedRow{
		{Label: "Reference", Measurement: reference}, {Label: "Candidate", Measurement: candidate},
	}), "\n")
	for _, value := range []string{"10 ms", "500 ns", "64 MiB", "64 B"} {
		if !strings.Contains(text, value) {
			t.Errorf("missing readable value %q", value)
		}
	}
	for _, rawMagnitude := range []string{"10,000,000", "67,108,864"} {
		if strings.Contains(text, rawMagnitude) {
			t.Errorf("unreadable raw magnitude %q remains", rawMagnitude)
		}
	}
}

func TestDifferentCPUAndSampleCountsCannotBeCompared(t *testing.T) {
	group, values := renderTestGroup(), renderTestValues()
	key := measurementKey{Package: "fixture", Name: "BenchmarkCandidate"}
	candidate := values[key]
	candidate.CPU = 2
	values[key] = candidate
	if _, err := resolveRows(group, group.Rows, values); err == nil || !strings.Contains(err.Error(), "GOMAXPROCS") {
		t.Fatalf("expected differing CPU settings to be rejected, got %v", err)
	}
	candidate.CPU = 1
	candidate.Samples = candidate.Samples[:2]
	values[key] = candidate
	if _, err := groupContext(group, values); err == nil || !strings.Contains(err.Error(), "different samples") {
		t.Fatalf("expected differing sample counts to be rejected, got %v", err)
	}
}

func TestIncompatibleWorkloadsHaveNoSpeedClaim(t *testing.T) {
	root := renderTestRoot(t)
	for _, category := range []string{"crypto", "protocol"} {
		values, _, err := loadCategory(root, category)
		if err != nil {
			t.Fatal(err)
		}
		for _, group := range groups[category] {
			t.Run(category+"/"+group.Title, func(t *testing.T) {
				if group.Baseline != "" {
					t.Fatal("incompatible workload has an invented baseline")
				}
				text, err := renderGroup(group, values)
				if err != nil {
					t.Fatal(err)
				}
				if !strings.Contains(text, "No equivalent reference measured") {
					t.Fatal("missing comparability boundary")
				}
				for _, invalid := range []string{"Speed vs baseline", "lowest recorded median", "×"} {
					if strings.Contains(text, invalid) {
						t.Errorf("incompatible workloads contain %q", invalid)
					}
				}
			})
		}
	}
}

func TestRecordedSortingExampleAndObservedSpread(t *testing.T) {
	// A fixed recorded fixture verifies arithmetic independently of future reruns.
	goSort := measurement{Name: "BenchmarkSortInts/Asc_Int_1000000", CPU: 16}
	parsort := measurement{Name: "BenchmarkParsortInts/Asc_Int_1000000", CPU: 16}
	for _, sample := range []rendererSample{
		{Time: 60_364_617, Allocated: 8_003_607, Allocs: 1},
		{Time: 59_796_753, Allocated: 8_003_584, Allocs: 1},
		{Time: 60_154_184, Allocated: 8_003_607, Allocs: 1},
	} {
		goSort.Samples = append(goSort.Samples, map[string]float64{"ns/op": sample.Time, "B/op": sample.Allocated, "allocs/op": sample.Allocs})
	}
	for _, sample := range []rendererSample{
		{Time: 14_872_250, Allocated: 40_110_987, Allocs: 89},
		{Time: 14_812_719, Allocated: 40_110_991, Allocs: 89},
		{Time: 14_736_572, Allocated: 40_110_990, Allocs: 89},
	} {
		parsort.Samples = append(parsort.Samples, map[string]float64{"ns/op": sample.Time, "B/op": sample.Allocated, "allocs/op": sample.Allocs})
	}
	if goSort.median("ns/op") != 60_154_184 || parsort.median("ns/op") != 14_812_719 || parsort.median("B/op") != 40_110_990 {
		t.Fatal("wrong fixture medians")
	}
	if low, high := parsort.bounds("ns/op"); low != 14_736_572 || high != 14_872_250 {
		t.Fatalf("wrong observed range: %v–%v", low, high)
	}
	if got := ratio(goSort.median("ns/op"), parsort.median("ns/op")); got != "4.06×" {
		t.Fatalf("wrong speed ratio: %s", got)
	}
	if got := allocatedComparison(parsort.median("B/op"), goSort.median("B/op")); got != "5.01×" {
		t.Fatalf("wrong allocation ratio: %s", got)
	}
	group := groupDef{Baseline: "Go", Operation: "sort"}
	resolved := []resolvedRow{{Label: "Go", Measurement: goSort}, {Label: "Parsort", Measurement: parsort}}
	if text := strings.Join(comparisonTable(group, resolved), "\n"); !strings.Contains(text, "| Parsort | 14.81 | 4.06× | 38.25 | 89 |") {
		t.Fatal("sorting comparison lost its recorded cost/speed tradeoff")
	}
	if text := strings.Join(rangesTable(group, resolved), "\n"); !strings.Contains(text, "14.74–14.87 ms") {
		t.Fatal("observed sorting range is missing")
	}
}

func TestHeadingAnchorsFollowGitHubPunctuation(t *testing.T) {
	for _, test := range []rendererAnchorCase{
		{Title: "Shared counter: 90% updates, 10% reads", Want: "shared-counter-90-updates-10-reads"},
		{Title: "Counter request/reply reads", Want: "counter-requestreply-reads"},
		{Title: "Input / output (µs)", Want: "input--output-µs"},
		{Title: "A_B and C++", Want: "a_b-and-c"},
	} {
		t.Run(test.Title, func(t *testing.T) {
			if got := headingAnchor(test.Title); got != test.Want {
				t.Fatalf("want %q, got %q", test.Want, got)
			}
		})
	}
}

func TestReplacementPreservesUnrelatedMarkdown(t *testing.T) {
	original := "# Fixture\n\n" + startMarker + "\nOld results\n" + endMarker + "\n\nKeep the methodology.\n"
	rendered := startMarker + "\nNew results\n" + endMarker
	updated, err := replaceResults(original, rendered)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(updated, "# Fixture\n\n") || !strings.HasSuffix(updated, "\n\nKeep the methodology.\n") {
		t.Fatal("replacing the results changed unrelated Markdown")
	}
	if again, err := replaceResults(updated, rendered); err != nil || again != updated {
		t.Fatal("replacing an identical block was not idempotent")
	}
	for _, invalid := range []string{
		"No markers", original + startMarker, original + endMarker, endMarker + "\n" + startMarker,
	} {
		if _, err := replaceResults(invalid, rendered); err == nil {
			t.Fatal("invalid marker structure was accepted")
		}
	}
}

func TestReferenceFirstAndSortingDoesNotMutateRows(t *testing.T) {
	group, values := renderTestGroup(), renderTestValues()
	slices.Reverse(group.Rows)
	before := slices.Clone(group.Rows)
	resolved, err := resolveRows(group, group.Rows, values)
	if err != nil {
		t.Fatal(err)
	}
	if resolved[0].Label != group.Baseline || !slices.Equal(group.Rows, before) {
		t.Fatal("baseline ordering modified the comparison definitions")
	}
}
