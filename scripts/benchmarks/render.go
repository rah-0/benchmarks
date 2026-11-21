package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

type displayUnit struct {
	Scale float64
	Label string
}

type resolvedRow struct {
	Label       string
	Measurement measurement
}

type resolvedCase struct {
	Label string
	Rows  []resolvedRow
}

var timeUnits = []displayUnit{{1e9, "s"}, {1e6, "ms"}, {1e3, "µs"}, {1, "ns"}}
var byteUnits = []displayUnit{{1024 * 1024 * 1024, "GiB"}, {1024 * 1024, "MiB"}, {1024, "KiB"}, {1, "B"}}

var reportVersionPattern = regexp.MustCompile(`(?m)^go version (\S+)`)
var reportDatePattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}`)

// groupedDecimal keeps the fractional part untouched while grouping whole digits.
func groupedDecimal(value string) string {
	whole, fraction, hasFraction := strings.Cut(value, ".")
	sign := ""
	if strings.HasPrefix(whole, "-") {
		sign, whole = "-", whole[1:]
	}
	var output strings.Builder
	output.WriteString(sign)
	for index, digit := range whole {
		if index > 0 && (len(whole)-index)%3 == 0 {
			output.WriteByte(',')
		}
		output.WriteRune(digit)
	}
	if hasFraction {
		output.WriteByte('.')
		output.WriteString(fraction)
	}
	return output.String()
}

func number(value float64) string {
	if value == 0 {
		return "0"
	}
	if math.Abs(value) < 0.1 {
		return strconv.FormatFloat(value, 'g', 3, 64)
	}
	formatted := groupedDecimal(strconv.FormatFloat(value, 'f', 2, 64))
	return strings.TrimSuffix(strings.TrimRight(formatted, "0"), ".")
}

func unitFor(values []float64, units []displayUnit) displayUnit {
	smallest := math.Inf(1)
	for _, value := range values {
		if value > 0 && value < smallest {
			smallest = value
		}
	}
	if math.IsInf(smallest, 1) {
		smallest = 0
	}
	for _, unit := range units {
		if smallest >= unit.Scale {
			return unit
		}
	}
	return units[len(units)-1]
}

func quantity(value float64, units []displayUnit) string {
	unit := unitFor([]float64{value}, units)
	return number(value/unit.Scale) + " " + unit.Label
}

func ratio(numerator, denominator float64) string {
	if denominator <= 0 {
		return "—"
	}
	value := numerator / denominator
	if value > 0 && value < 0.005 {
		return "<0.01×"
	}
	return groupedDecimal(strconv.FormatFloat(value, 'f', 2, 64)) + "×"
}

func spansOrdersOfMagnitude(values []float64) bool {
	low, high := math.Inf(1), 0.0
	for _, value := range values {
		if value > 0 {
			low = min(low, value)
			high = max(high, value)
		}
	}
	return high > 0 && high/low >= 1000
}

func allocatedComparison(value, baseline float64) string {
	if baseline > 0 {
		return ratio(value, baseline)
	}
	if value > 0 {
		return "+" + quantity(value, byteUnits)
	}
	return "same (0 B)"
}

func markdownTable(headings []string, rows [][]string) []string {
	line := func(cells []string) string {
		escaped := make([]string, len(cells))
		for index, cell := range cells {
			escaped[index] = strings.ReplaceAll(cell, "|", `\|`)
		}
		return "| " + strings.Join(escaped, " | ") + " |"
	}
	separators := make([]string, len(headings))
	for index := range separators {
		separators[index] = "---:"
	}
	separators[0] = "---"
	lines := []string{line(headings), line(separators)}
	for _, row := range rows {
		lines = append(lines, line(row))
	}
	return lines
}

func headingAnchor(title string) string {
	var result strings.Builder
	for _, char := range strings.ToLower(title) {
		switch {
		case char == ' ':
			result.WriteByte('-')
		case unicode.IsLetter(char), unicode.IsNumber(char), unicode.IsSpace(char), char == '_', char == '-':
			result.WriteRune(char)
		}
	}
	return result.String()
}

func groupCases(group groupDef) []caseDef {
	if group.Cases != nil {
		return group.Cases
	}
	return []caseDef{{Rows: group.Rows}}
}

func resolveRows(group groupDef, rows []rowDef, measurements map[measurementKey]measurement) ([]resolvedRow, error) {
	if len(rows) == 0 {
		return nil, fmt.Errorf("%s: empty or duplicate implementation labels", group.Title)
	}
	labels := make(map[string]bool)
	resolved := make([]resolvedRow, 0, len(rows))
	var cpu int
	for index, row := range rows {
		if labels[row.Label] {
			return nil, fmt.Errorf("%s: empty or duplicate implementation labels", group.Title)
		}
		labels[row.Label] = true
		value, exists := measurements[measurementKey{Package: row.Package, Name: row.Name}]
		if !exists {
			return nil, fmt.Errorf("%s: missing measurement for %s", group.Title, row.Name)
		}
		if index > 0 && value.CPU != cpu {
			return nil, fmt.Errorf("%s: comparing different GOMAXPROCS settings", group.Title)
		}
		cpu = value.CPU
		resolved = append(resolved, resolvedRow{Label: row.Label, Measurement: value})
	}
	if group.Baseline != "" {
		if !labels[group.Baseline] {
			return nil, fmt.Errorf("%s: missing declared baseline", group.Title)
		}
		slices.SortStableFunc(resolved, func(a, b resolvedRow) int {
			if a.Label == group.Baseline && b.Label != group.Baseline {
				return -1
			}
			if a.Label != group.Baseline && b.Label == group.Baseline {
				return 1
			}
			return 0
		})
	}
	return resolved, nil
}

func validateGroups(category string, measurements map[measurementKey]measurement) error {
	definitions, exists := groups[category]
	if !exists {
		return fmt.Errorf("%s: missing presentation definitions", category)
	}
	return validateGroupDefinitions(category, definitions, measurements)
}

func validateGroupDefinitions(category string, definitions []groupDef, measurements map[measurementKey]measurement) error {
	used := make(map[measurementKey]bool)
	titles := make(map[string]bool)
	for _, group := range definitions {
		if titles[group.Title] {
			return fmt.Errorf("%s: duplicate workload title", category)
		}
		titles[group.Title] = true
		caseLabels := make(map[string]bool)
		var variants []string
		groupedCases := groupCases(group)
		if len(groupedCases) == 0 {
			return fmt.Errorf("%s: empty workload", group.Title)
		}
		for index, currentCase := range groupedCases {
			if caseLabels[currentCase.Label] {
				return fmt.Errorf("%s: duplicate input label", group.Title)
			}
			caseLabels[currentCase.Label] = true
			current := make([]string, len(currentCase.Rows))
			for rowIndex, row := range currentCase.Rows {
				current[rowIndex] = row.Label
			}
			if index > 0 && !slices.Equal(current, variants) {
				return fmt.Errorf("%s: sweep implementations differ", group.Title)
			}
			variants = current
			for _, row := range currentCase.Rows {
				key := measurementKey{Package: row.Package, Name: row.Name}
				if used[key] {
					return fmt.Errorf("%s: duplicate presentation of %s", category, row.Name)
				}
				if _, exists := measurements[key]; !exists {
					return fmt.Errorf("%s: missing measurement for %s", category, row.Name)
				}
				used[key] = true
			}
			if _, err := resolveRows(group, currentCase.Rows, measurements); err != nil {
				return err
			}
		}
	}
	var unmatched []measurementKey
	for key := range measurements {
		if !used[key] {
			unmatched = append(unmatched, key)
		}
	}
	if len(unmatched) > 0 {
		slices.SortFunc(unmatched, func(a, b measurementKey) int {
			if order := strings.Compare(a.Package, b.Package); order != 0 {
				return order
			}
			return strings.Compare(a.Name, b.Name)
		})
		names := make([]string, len(unmatched))
		for index, key := range unmatched {
			names[index] = key.Name
		}
		return fmt.Errorf("%s: benchmarks need presentation rules: %s", category, strings.Join(names, ", "))
	}
	return nil
}

func groupContext(group groupDef, measurements map[measurementKey]measurement) (string, error) {
	count, cpu := -1, -1
	for _, currentCase := range groupCases(group) {
		resolved, err := resolveRows(group, currentCase.Rows, measurements)
		if err != nil {
			return "", err
		}
		for _, row := range resolved {
			value := row.Measurement
			if count >= 0 && (len(value.Samples) != count || value.CPU != cpu) {
				return "", fmt.Errorf("%s: separate groups are required for different samples or CPU settings", group.Title)
			}
			count, cpu = len(value.Samples), value.CPU
		}
	}
	if count < 1 {
		return "", fmt.Errorf("%s: missing workload samples", group.Title)
	}
	reference := "No equivalent reference measured; absolute costs are shown."
	if group.Baseline != "" {
		reference = fmt.Sprintf("Baseline: **%s = 1.00×**.", group.Baseline)
	}
	return fmt.Sprintf("%s Median of %d samples; GOMAXPROCS=%d.", reference, count, cpu), nil
}

func baselineFor(group groupDef, rows []resolvedRow) (measurement, bool) {
	for _, row := range rows {
		if group.Baseline != "" && row.Label == group.Baseline {
			return row.Measurement, true
		}
	}
	return measurement{}, false
}

func comparisonTable(group groupDef, resolved []resolvedRow) []string {
	baseline, hasBaseline := baselineFor(group, resolved)
	times, allocations := make([]float64, len(resolved)), make([]float64, len(resolved))
	for index, row := range resolved {
		times[index] = row.Measurement.median("ns/op")
		allocations[index] = row.Measurement.median("B/op")
	}
	timeUnit, byteUnit := unitFor(times, timeUnits), unitFor(allocations, byteUnits)
	mixedTime, mixedBytes := spansOrdersOfMagnitude(times), spansOrdersOfMagnitude(allocations)
	operation := group.Operation
	if operation == "" || operation == "operation" {
		operation = "op"
	}
	timeHeading, byteHeading := "Time/"+operation, "Allocated/"+operation
	if !mixedTime {
		timeHeading += ", " + timeUnit.Label
	}
	if !mixedBytes {
		byteHeading += ", " + byteUnit.Label
	}
	rowHeading := "Operation / variant"
	if hasBaseline {
		rowHeading = "Implementation"
	}
	headings := []string{rowHeading, timeHeading + " ↓"}
	if hasBaseline {
		headings = append(headings, "Speed vs baseline ↑")
	}
	headings = append(headings, byteHeading+" ↓", "Allocs/"+operation+" ↓")
	var rows [][]string
	for _, current := range resolved {
		name, value := current.Label, current.Measurement
		if current.Label == group.Baseline {
			name = "**" + name + " (baseline)**"
		}
		time, allocated := value.median("ns/op"), value.median("B/op")
		formattedTime, formattedBytes := number(time/timeUnit.Scale), number(allocated/byteUnit.Scale)
		if mixedTime {
			formattedTime = quantity(time, timeUnits)
		}
		if mixedBytes {
			formattedBytes = quantity(allocated, byteUnits)
		}
		row := []string{name, formattedTime}
		if hasBaseline {
			row = append(row, ratio(baseline.median("ns/op"), time))
		}
		row = append(row, formattedBytes, number(value.median("allocs/op")))
		rows = append(rows, row)
	}
	return markdownTable(headings, rows)
}

func rangesTable(group groupDef, resolved []resolvedRow) []string {
	baseline, hasBaseline := baselineFor(group, resolved)
	rowHeading := "Operation / variant"
	if hasBaseline {
		rowHeading = "Implementation"
	}
	headings := []string{rowHeading, "Time min–max", "Allocated min–max", "Allocs min–max"}
	if hasBaseline {
		headings = append(headings, "Allocated vs baseline ↓")
	}
	var rows [][]string
	for _, current := range resolved {
		value := current.Measurement
		row := []string{current.Label}
		for _, metric := range []string{"ns/op", "B/op"} {
			low, high := value.bounds(metric)
			units := timeUnits
			if metric == "B/op" {
				units = byteUnits
			}
			unit := unitFor([]float64{low, high}, units)
			row = append(row, number(low/unit.Scale)+"–"+number(high/unit.Scale)+" "+unit.Label)
		}
		low, high := value.bounds("allocs/op")
		countRange := number(low)
		if low != high {
			countRange += "–" + number(high)
		}
		row = append(row, countRange)
		if hasBaseline {
			row = append(row, allocatedComparison(value.median("B/op"), baseline.median("B/op")))
		}
		rows = append(rows, row)
	}
	return markdownTable(headings, rows)
}

func byMedianTime(a, b resolvedRow) int {
	aTime, bTime := a.Measurement.median("ns/op"), b.Measurement.median("ns/op")
	if aTime < bTime {
		return -1
	}
	if aTime > bTime {
		return 1
	}
	return 0
}

func takeaway(group groupDef, resolved []resolvedRow, caseLabel string) string {
	if group.Baseline == "" || len(resolved) < 2 {
		return ""
	}
	baseline, _ := baselineFor(group, resolved)
	ordered := slices.Clone(resolved)
	slices.SortStableFunc(ordered, byMedianTime)
	label, best := ordered[0].Label, ordered[0].Measurement
	prefix := ""
	if caseLabel != "" {
		prefix = "At " + caseLabel + ", "
	}
	var text string
	if label == group.Baseline {
		text = prefix + label + " has the lowest recorded median time in this comparison."
	} else {
		text = prefix + label + " has the lowest recorded median time: " +
			ratio(baseline.median("ns/op"), best.median("ns/op")) + " the speed of " + group.Baseline
		switch {
		case baseline.median("B/op") > 0:
			text += ", with " + allocatedComparison(best.median("B/op"), baseline.median("B/op")) + " its allocated bytes per operation."
		case best.median("B/op") > 0:
			text += ", adding " + quantity(best.median("B/op"), byteUnits) + " allocated per operation."
		default:
			text += "; both allocate 0 B per operation."
		}
	}
	secondLabel, second := ordered[1].Label, ordered[1].Measurement
	_, bestHigh := best.bounds("ns/op")
	secondLow, _ := second.bounds("ns/op")
	if bestHigh >= secondLow {
		text += " The timing ranges of " + label + " and " + secondLabel + " overlap; their order is inconclusive."
	}
	return text
}

func renderGroup(group groupDef, measurements map[measurementKey]measurement) (string, error) {
	context, err := groupContext(group, measurements)
	if err != nil {
		return "", err
	}
	lines := []string{"### " + group.Title, "", group.Description, "", context}
	if group.Note != "" {
		lines = append(lines, "", group.Note)
	}
	lines = append(lines, "")
	var resolvedCases []resolvedCase
	for _, current := range groupCases(group) {
		rows, err := resolveRows(group, current.Rows, measurements)
		if err != nil {
			return "", err
		}
		resolvedCases = append(resolvedCases, resolvedCase{Label: current.Label, Rows: rows})
	}
	if group.Cases == nil {
		resolved := resolvedCases[0].Rows
		lines = append(lines, comparisonTable(group, resolved)...)
		if summary := takeaway(group, resolved, ""); summary != "" {
			lines = append(lines, "", summary)
		}
		detailLabel := "Observed sample ranges"
		if group.Baseline != "" {
			detailLabel = "Observed ranges and allocation comparison"
		}
		lines = append(lines, "", "<details>", "<summary>"+detailLabel+"</summary>", "")
		lines = append(lines, rangesTable(group, resolved)...)
	} else {
		lines = append(lines, renderSweep(group, resolvedCases)...)
	}
	lines = append(lines, "", "</details>")
	return strings.Join(lines, "\n"), nil
}

func renderSweep(group groupDef, resolvedCases []resolvedCase) []string {
	var lines []string
	first, last := resolvedCases[0], resolvedCases[len(resolvedCases)-1]
	if group.Baseline != "" {
		headings := []string{"Input / workload", group.Baseline + " time (1.00×)"}
		for _, row := range first.Rows {
			if row.Label != group.Baseline {
				headings = append(headings, row.Label+" speed ↑")
			}
		}
		var rows [][]string
		for _, current := range resolvedCases {
			baseline, _ := baselineFor(group, current.Rows)
			row := []string{current.Label, quantity(baseline.median("ns/op"), timeUnits)}
			for _, candidate := range current.Rows {
				if candidate.Label != group.Baseline {
					row = append(row, ratio(baseline.median("ns/op"), candidate.Measurement.median("ns/op")))
				}
			}
			rows = append(rows, row)
		}
		lines = append(lines, markdownTable(headings, rows)...)
		lines = append(lines, "", "Allocation cost accompanies the speed comparison:", "")
		headings = []string{"Input / workload", group.Baseline + " allocated"}
		for _, row := range first.Rows {
			if row.Label != group.Baseline {
				headings = append(headings, row.Label+" allocated vs baseline ↓")
			}
		}
		rows = nil
		for _, current := range resolvedCases {
			baseline, _ := baselineFor(group, current.Rows)
			row := []string{current.Label, quantity(baseline.median("B/op"), byteUnits)}
			for _, candidate := range current.Rows {
				if candidate.Label != group.Baseline {
					row = append(row, allocatedComparison(candidate.Measurement.median("B/op"), baseline.median("B/op")))
				}
			}
			rows = append(rows, row)
		}
		lines = append(lines, markdownTable(headings, rows)...)
		firstBest, lastBest := slices.MinFunc(first.Rows, byMedianTime).Label, slices.MinFunc(last.Rows, byMedianTime).Label
		if firstBest != lastBest {
			lines = append(lines, "", "The lowest median shifts from "+firstBest+" at "+first.Label+" to "+lastBest+" at "+last.Label+".")
		}
		if summary := takeaway(group, last.Rows, last.Label); summary != "" {
			lines = append(lines, "", summary)
		}
	} else {
		headings := []string{"Input / workload", "Implementation", "Time/op ↓", "Allocated/op ↓", "Allocs/op ↓"}
		singleImplementation := len(first.Rows) == 1
		if singleImplementation {
			headings = slices.Delete(headings, 1, 2)
		}
		var rows [][]string
		for _, current := range resolvedCases {
			for _, item := range current.Rows {
				value := item.Measurement
				row := []string{current.Label, item.Label, quantity(value.median("ns/op"), timeUnits), quantity(value.median("B/op"), byteUnits), number(value.median("allocs/op"))}
				if singleImplementation {
					row = slices.Delete(row, 1, 2)
				}
				rows = append(rows, row)
			}
		}
		lines = append(lines, markdownTable(headings, rows)...)
	}
	lines = append(lines, "", "<details>", "<summary>Absolute measurements, allocation counts, and observed ranges</summary>", "")
	for _, current := range resolvedCases {
		lines = append(lines, "#### "+current.Label, "")
		if group.Baseline != "" {
			lines = append(lines, comparisonTable(group, current.Rows)...)
			lines = append(lines, "")
		}
		lines = append(lines, rangesTable(group, current.Rows)...)
		lines = append(lines, "")
	}
	return lines
}

func joinedInts(set map[int]bool) string {
	values := make([]int, 0, len(set))
	for value := range set {
		values = append(values, value)
	}
	slices.Sort(values)
	labels := make([]string, len(values))
	for index, value := range values {
		labels[index] = strconv.Itoa(value)
	}
	return strings.Join(labels, ", ")
}

func renderCategory(root, category string) (string, error) {
	measurements, env, err := loadCategory(root, category)
	if err != nil {
		return "", err
	}
	if err := validateGroups(category, measurements); err != nil {
		return "", err
	}
	environmentData, err := os.ReadFile(filepath.Join(root, "results", "environment.txt"))
	if err != nil {
		return "", err
	}
	version := reportVersionPattern.FindStringSubmatch(string(environmentData))
	date := reportDatePattern.FindString(string(environmentData))
	if version == nil || date == "" {
		return "", fmt.Errorf("environment.txt: missing measurement date or Go version")
	}
	counts, cpus := make(map[int]bool), make(map[int]bool)
	for _, value := range measurements {
		counts[len(value.Samples)] = true
		cpus[value.CPU] = true
	}
	lines := []string{
		startMarker, "## Results", "",
		fmt.Sprintf("**%s · %s · %s/%s · %s**", date, version[1], env.OS, env.Arch, env.CPU),
		fmt.Sprintf("GOMAXPROCS: %s; samples per case: %s.", joinedInts(cpus), joinedInts(counts)),
		"",
		"Tables show medians. **Speed = baseline time ÷ candidate time**: 2.00× is twice",
		"the baseline speed; 0.50× is half. Allocated bytes and allocation counts describe",
		"Go heap allocation per operation, not retained memory or process RSS.",
		"These samples are descriptive; observed min–max ranges appear in the details",
		"and are not confidence intervals. A zero denominator is shown as **—**.",
		"",
		"[Environment](../results/environment.txt) · [Versions](../results/modules.txt) · [Reading the comparisons](../README.md#reading-the-comparisons)",
		"",
	}
	definitions := groups[category]
	navigation := make([]string, len(definitions))
	for index, group := range definitions {
		navigation[index] = "[" + group.Title + "](#" + headingAnchor(group.Title) + ")"
	}
	if len(navigation) > 6 {
		lines = append(lines, "<details>", "<summary>Jump to a workload</summary>", "")
		for _, link := range navigation {
			lines = append(lines, "- "+link)
		}
		lines = append(lines, "", "</details>")
	} else {
		lines = append(lines, "**Workloads:** "+strings.Join(navigation, " · "))
	}
	for _, group := range definitions {
		block, err := renderGroup(group, measurements)
		if err != nil {
			return "", err
		}
		lines = append(lines, "", block)
	}
	var links []string
	for _, suite := range suites {
		if suite.Name == category {
			for _, pkg := range suite.Packages {
				links = append(links, "["+pkg+"](../results/"+strings.ReplaceAll(pkg, "/", "-")+".txt)")
			}
			break
		}
	}
	lines = append(lines, "", "**Recorded samples:** "+strings.Join(links, " · "), "", endMarker)
	return strings.Join(lines, "\n"), nil
}

func replaceResults(original, rendered string) (string, error) {
	if strings.Count(original, startMarker) != 1 || strings.Count(original, endMarker) != 1 {
		return "", fmt.Errorf("README must contain exactly one benchmark-results block")
	}
	before, rest, _ := strings.Cut(original, startMarker)
	_, after, found := strings.Cut(rest, endMarker)
	if !found {
		return "", fmt.Errorf("README benchmark-results markers are reversed")
	}
	return before + rendered + after, nil
}
