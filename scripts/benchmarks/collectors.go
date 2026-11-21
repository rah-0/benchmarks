package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const collectorBaseline = "fluent-bit"

var collectorNames = []string{collectorBaseline, "slogx-collector", "vector"}
var collectorLabels = map[string]string{
	collectorBaseline: "Fluent Bit", "slogx-collector": "slogx-collector", "vector": "Vector",
}
var collectorEnvironmentFields = strings.Fields("type date go os arch collector_cpu records sizes repeat outage_seconds images collectors container_manager slogx_module slogx_binary_sha256 slogx_go kernel harness_cpus cpu storage_filesystem")
var collectorSampleFields = strings.Fields("collector sample records record_bytes input_bytes seconds records_per_second input_mib_per_second producer_seconds http_requests retry_requests duplicates outage_seconds")
var collectorProcessFields = strings.Fields("image image_id cpu_seconds peak_rss_bytes")
var collectorDigest = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
var collectorBinaryDigest = regexp.MustCompile(`^[0-9a-f]{64}$`)

type collectorEnvironment struct {
	Type              string            `json:"type"`
	Date              string            `json:"date"`
	Go                string            `json:"go"`
	OS                string            `json:"os"`
	Arch              string            `json:"arch"`
	CollectorCPU      string            `json:"collector_cpu"`
	Records           int64             `json:"records"`
	Sizes             []int64           `json:"sizes"`
	Repeat            int               `json:"repeat"`
	OutageSeconds     float64           `json:"outage_seconds"`
	Images            map[string]string `json:"images"`
	Collectors        []string          `json:"collectors"`
	ContainerManager  string            `json:"container_manager"`
	SlogxModule       string            `json:"slogx_module"`
	SlogxBinarySHA256 string            `json:"slogx_binary_sha256"`
	SlogxGo           string            `json:"slogx_go"`
	Kernel            string            `json:"kernel"`
	HarnessCPUs       string            `json:"harness_cpus"`
	CPU               string            `json:"cpu"`
	StorageFilesystem string            `json:"storage_filesystem"`
}

type collectorSample struct {
	Collector            string  `json:"collector"`
	Sample               int     `json:"sample"`
	Records              int64   `json:"records"`
	RecordBytes          int64   `json:"record_bytes"`
	InputBytes           int64   `json:"input_bytes"`
	Seconds              float64 `json:"seconds"`
	RecordsPerSecond     float64 `json:"records_per_second"`
	InputMiBPerSecond    float64 `json:"input_mib_per_second"`
	ProducerSeconds      float64 `json:"producer_seconds"`
	HTTPRequests         int64   `json:"http_requests"`
	RetryRequests        int64   `json:"retry_requests"`
	Duplicates           int64   `json:"duplicates"`
	OutageSeconds        float64 `json:"outage_seconds"`
	Image                string  `json:"image"`
	ImageID              string  `json:"image_id"`
	CPUSeconds           float64 `json:"cpu_seconds"`
	PeakRSSBytes         int64   `json:"peak_rss_bytes"`
	RecoveryDrainSeconds float64 `json:"recovery_drain_seconds"`
}

type collectorCompletion struct {
	Type    string `json:"type"`
	Samples int    `json:"samples"`
}

type collectorGroupKey struct {
	Name string
	Size int64
}

type collectorSampleKey struct {
	Name   string
	Size   int64
	Sample int
}

type collectorReport struct {
	Environment collectorEnvironment
	Groups      map[collectorGroupKey][]collectorSample
	ImageIDs    map[string]string
	Count       int
}

type collectorMetricDef struct {
	Name   string
	Scale  float64
	Places int
}

func collectorKnown(name string) bool {
	return name == "http-control" || slices.Contains(collectorNames, name)
}

func collectorRequireFields(raw map[string]json.RawMessage, required, optional []string) error {
	if raw == nil {
		return fmt.Errorf("unexpected or missing fields")
	}
	for _, name := range required {
		if _, ok := raw[name]; !ok {
			return fmt.Errorf("unexpected or missing fields")
		}
	}
	for name, value := range raw {
		if (!slices.Contains(required, name) && !slices.Contains(optional, name)) || bytes.Equal(bytes.TrimSpace(value), []byte("null")) {
			return fmt.Errorf("unexpected or missing fields")
		}
	}
	return nil
}

func collectorDecode(raw map[string]json.RawMessage, target any) error {
	data, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func collectorNonnegative(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0) && value >= 0
}

func collectorClose(left, right, relative, absolute float64) bool {
	if math.IsInf(left, 0) || math.IsInf(right, 0) || math.IsNaN(left) || math.IsNaN(right) {
		return false
	}
	return math.Abs(left-right) <= max(absolute, relative*max(math.Abs(left), math.Abs(right)))
}

// parseCollectors requires the full declared matrix and validates the evidence
// used in ratios before any documentation is generated.
func parseCollectors(path string) (collectorReport, error) {
	var report collectorReport
	data, err := os.ReadFile(path)
	if err != nil {
		return report, err
	}
	where := filepath.Base(path)
	fail := func(message string) (collectorReport, error) {
		return collectorReport{}, fmt.Errorf("%s: %s", where, message)
	}
	var entries []map[string]json.RawMessage
	for _, line := range bytes.Split(data, []byte("\n")) {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var entry map[string]json.RawMessage
		if err := json.Unmarshal(line, &entry); err != nil {
			return fail("invalid JSON record")
		}
		entries = append(entries, entry)
	}
	if len(entries) < 3 {
		return fail("missing environment, samples, or completion")
	}
	rawEnvironment, rawCompletion := entries[0], entries[len(entries)-1]
	if err := collectorRequireFields(rawEnvironment, collectorEnvironmentFields, nil); err != nil {
		return fail(err.Error())
	}
	if err := collectorRequireFields(rawCompletion, []string{"type", "samples"}, nil); err != nil {
		return fail(err.Error())
	}
	var env collectorEnvironment
	var completion collectorCompletion
	if collectorDecode(rawEnvironment, &env) != nil || collectorDecode(rawCompletion, &completion) != nil {
		return fail("invalid environment or completion field type")
	}
	if env.Type != "environment" || completion.Type != "complete" {
		return fail("missing environment or completion record")
	}
	if env.Repeat < 1 || env.Records < 1 || !collectorNonnegative(env.OutageSeconds) {
		return fail("invalid workload counts or outage duration")
	}
	seenNames := make(map[string]bool)
	for _, name := range env.Collectors {
		if !collectorKnown(name) || seenNames[name] {
			return fail("collector matrix must contain the three named collectors and HTTP control")
		}
		seenNames[name] = true
	}
	if len(seenNames) != 4 {
		return fail("collector matrix must contain the three named collectors and HTTP control")
	}
	seenSizes := make(map[int64]bool)
	for _, size := range env.Sizes {
		if size < 256 || size > 16384 || seenSizes[size] {
			return fail("invalid or repeated record size")
		}
		seenSizes[size] = true
	}
	if len(seenSizes) == 0 {
		return fail("invalid or repeated record size")
	}
	nonStrings := []string{"records", "sizes", "repeat", "outage_seconds", "images", "collectors"}
	for name, value := range rawEnvironment {
		if slices.Contains(nonStrings, name) {
			continue
		}
		var text string
		if json.Unmarshal(value, &text) != nil || strings.ContainsAny(text, "\r\n|`<>") {
			return fail("invalid environment field " + name)
		}
	}
	if !collectorBinaryDigest.MatchString(env.SlogxBinarySHA256) || !strings.HasPrefix(env.SlogxModule, "github.com/rah-0/slogx-collector@v") {
		return fail("missing collector build identity")
	}
	if len(env.Images) != len(collectorNames) {
		return fail("image map differs from the collector matrix")
	}
	for _, name := range collectorNames {
		reference, ok := env.Images[name]
		if !ok {
			return fail("image map differs from the collector matrix")
		}
		at := strings.LastIndexByte(reference, '@')
		if at < 0 || !collectorDigest.MatchString(reference[at+1:]) || strings.ContainsAny(reference, "\r\n|`<> ") {
			return fail("collector image must have an explicit digest")
		}
	}
	sampleCount := len(entries) - 2
	// Check with division before multiplication so a malformed repeat cannot
	// overflow or allocate an unbounded expected-sample map.
	if completion.Samples != sampleCount || sampleCount/4/len(env.Sizes) != env.Repeat || sampleCount%(4*len(env.Sizes)) != 0 {
		return fail("incomplete sample matrix")
	}
	report = collectorReport{Environment: env, Groups: make(map[collectorGroupKey][]collectorSample), ImageIDs: make(map[string]string), Count: sampleCount}
	seen := make(map[collectorSampleKey]bool)
	for _, raw := range entries[1 : len(entries)-1] {
		var name string
		if json.Unmarshal(raw["collector"], &name) != nil || !collectorKnown(name) {
			return fail("unmapped collector sample")
		}
		required := slices.Clone(collectorSampleFields)
		if name != "http-control" {
			required = append(required, collectorProcessFields...)
		}
		var optional []string
		if env.OutageSeconds > 0 {
			optional = []string{"recovery_drain_seconds"}
		}
		if err := collectorRequireFields(raw, required, optional); err != nil {
			return fail(err.Error())
		}
		var sample collectorSample
		if collectorDecode(raw, &sample) != nil {
			return fail("invalid sample field type")
		}
		if sample.Sample < 1 || sample.Records < 1 || sample.RecordBytes < 1 || sample.InputBytes < 1 || sample.HTTPRequests < 1 || sample.RetryRequests < 0 || sample.Duplicates < 0 {
			return fail("invalid sample count")
		}
		key := collectorSampleKey{name, sample.RecordBytes, sample.Sample}
		if !seenSizes[sample.RecordBytes] || sample.Sample > env.Repeat || seen[key] {
			return fail("unexpected or duplicate sample")
		}
		seen[key] = true
		if sample.Records != env.Records || sample.InputBytes/sample.RecordBytes != env.Records || sample.InputBytes%sample.RecordBytes != 0 || sample.OutageSeconds != env.OutageSeconds {
			return fail("sample workload differs from its environment")
		}
		metrics := []float64{sample.Seconds, sample.RecordsPerSecond, sample.InputMiBPerSecond, sample.ProducerSeconds, sample.OutageSeconds}
		if name != "http-control" {
			metrics = append(metrics, sample.CPUSeconds, float64(sample.PeakRSSBytes))
			if sample.Image != env.Images[name] || !collectorDigest.MatchString(sample.ImageID) {
				return fail("sample has a different or invalid collector image")
			}
			if imageID, found := report.ImageIDs[name]; found && imageID != sample.ImageID {
				return fail("runtime image changed between samples")
			}
			report.ImageIDs[name] = sample.ImageID
		}
		if env.OutageSeconds > 0 {
			if sample.RetryRequests < 1 {
				return fail("outage sample did not encounter HTTP 503")
			}
			// Missing recovery is the Go encoder's representation of zero.
			metrics = append(metrics, sample.RecoveryDrainSeconds)
		} else if sample.RetryRequests != 0 {
			return fail("healthy sample contains HTTP 503 retries")
		}
		if slices.ContainsFunc(metrics, func(value float64) bool { return !collectorNonnegative(value) }) || sample.Seconds <= 0 {
			return fail("invalid sample metric")
		}
		if !collectorClose(sample.RecordsPerSecond, float64(sample.Records)/sample.Seconds, 1e-9, 0) || !collectorClose(sample.InputMiBPerSecond, float64(sample.InputBytes)/1048576/sample.Seconds, 1e-9, 0) {
			return fail("throughput disagrees with the recorded duration")
		}
		if env.OutageSeconds > 0 && !collectorClose(sample.RecoveryDrainSeconds, max(0, sample.Seconds-env.OutageSeconds), 1e-7, 1e-9) {
			return fail("recovery duration disagrees with the recorded outage")
		}
		group := collectorGroupKey{name, sample.RecordBytes}
		report.Groups[group] = append(report.Groups[group], sample)
	}
	return report, nil
}

func collectorMetric(sample collectorSample, metric string) float64 {
	switch metric {
	case "records_per_second":
		return sample.RecordsPerSecond
	case "peak_rss_bytes":
		return float64(sample.PeakRSSBytes)
	case "input_mib_per_second":
		return sample.InputMiBPerSecond
	case "cpu_seconds":
		return sample.CPUSeconds
	case "producer_seconds":
		return sample.ProducerSeconds
	case "http_requests":
		return float64(sample.HTTPRequests)
	case "recovery_drain_seconds":
		return sample.RecoveryDrainSeconds
	case "retry_requests":
		return float64(sample.RetryRequests)
	default:
		panic("unknown collector metric: " + metric)
	}
}

func collectorMedian(samples []collectorSample, metric string) float64 {
	values := make([]float64, len(samples))
	for i, sample := range samples {
		values[i] = collectorMetric(sample, metric)
	}
	slices.Sort(values)
	middle := len(values) / 2
	if len(values)%2 == 0 {
		return (values[middle-1] + values[middle]) / 2
	}
	return values[middle]
}

func collectorGrouped(value string) string {
	integer, fraction, hasFraction := strings.Cut(value, ".")
	sign := ""
	if strings.HasPrefix(integer, "-") {
		sign, integer = "-", integer[1:]
	}
	for i := len(integer) - 3; i > 0; i -= 3 {
		integer = integer[:i] + "," + integer[i:]
	}
	if hasFraction {
		return sign + integer + "." + fraction
	}
	return sign + integer
}

func collectorNumber(value float64, places int) string {
	formatted := strconv.FormatFloat(value, 'f', places, 64)
	rounded, _ := strconv.ParseFloat(formatted, 64)
	if value != 0 && rounded == 0 {
		return strconv.FormatFloat(value, 'g', 2, 64)
	}
	return collectorGrouped(formatted)
}

func collectorRatio(numerator, denominator float64) string {
	if denominator == 0 {
		return "—"
	}
	value := numerator / denominator
	if 0 < value && value < 0.005 {
		return "<0.01×"
	}
	return collectorNumber(value, 2) + "×"
}

func collectorBounds(samples []collectorSample, metric string) (float64, float64) {
	low, high := collectorMetric(samples[0], metric), collectorMetric(samples[0], metric)
	for _, sample := range samples[1:] {
		value := collectorMetric(sample, metric)
		low, high = min(low, value), max(high, value)
	}
	return low, high
}

func collectorObservedRange(samples []collectorSample, metric string, scale float64, places int) string {
	low, high := collectorBounds(samples, metric)
	return collectorNumber(low*scale, places) + "–" + collectorNumber(high*scale, places)
}

func collectorLabel(name string) string {
	if name == collectorBaseline {
		return "**Fluent Bit — baseline**"
	}
	return collectorLabels[name]
}

func renderCollectors(root string) (string, error) {
	healthy, err := parseCollectors(filepath.Join(root, "results", "healthy.jsonl"))
	if err != nil {
		return "", err
	}
	outage, err := parseCollectors(filepath.Join(root, "results", "outage.jsonl"))
	if err != nil {
		return "", err
	}
	env, outageEnv := healthy.Environment, outage.Environment
	left, right := env, outageEnv
	left.Date, right.Date = "", ""
	left.OutageSeconds, right.OutageSeconds = 0, 0
	if !reflect.DeepEqual(left, right) {
		return "", fmt.Errorf("collector runs have different workloads or environments")
	}
	if env.OutageSeconds != 0 || outageEnv.OutageSeconds <= 0 {
		return "", fmt.Errorf("collector comparisons require healthy and receiver-outage matrices")
	}
	if !reflect.DeepEqual(healthy.ImageIDs, outage.ImageIDs) {
		return "", fmt.Errorf("collector runtime image differs between healthy and outage runs")
	}
	versions := make(map[string]string)
	for _, name := range []string{"vector", "fluent-bit"} {
		reference, _, _ := strings.Cut(env.Images[name], "@")
		versions[name] = reference[strings.LastIndexByte(reference, ':')+1:]
	}
	lines := []string{
		startMarker, "## Results", "",
		fmt.Sprintf("**Baseline: Fluent Bit %s.** Compared with slogx-collector %s and Vector %s.", versions["fluent-bit"], env.SlogxModule[strings.LastIndexByte(env.SlogxModule, '@')+1:], versions["vector"]),
		fmt.Sprintf("Each sample delivers **%s records**; cells are medians of **%d samples**. Record sizes include the newline.", collectorGrouped(strconv.FormatInt(env.Records, 10)), env.Repeat),
		"", "**Persistence differs:** slogx delivers synced journal records; Vector can deliver",
		"before disk synchronization; Fluent Bit synchronizes filesystem chunks.",
		"The ratios compare these configured pipelines, including their distinct batching costs.",
		"", "Higher throughput and recovery-speed ratios are better; lower RSS ratios are better.",
		"Ratios use unrounded medians. Observed ranges appear in the details; these samples",
		"describe measured bursts and do not establish statistical significance.",
		"See [reading the comparisons](../README.md#reading-the-comparisons) for the shared conventions.",
		"", "### Healthy delivery", "",
		"Time runs from the first producer write to the final validated HTTP acknowledgement.",
		"RSS is collector-process lifetime peak memory, including startup and warmup.",
	}
	sizes := slices.Clone(env.Sizes)
	slices.Sort(sizes)
	for _, size := range sizes {
		base := healthy.Groups[collectorGroupKey{collectorBaseline, size}]
		baseRate, baseRSS := collectorMedian(base, "records_per_second"), collectorMedian(base, "peak_rss_bytes")
		lines = append(lines, "", "#### "+collectorGrouped(strconv.FormatInt(size, 10))+"-byte records", "",
			"| Collector | Records/s ↑ | Throughput vs baseline ↑ | Peak RSS (MiB) ↓ | RSS vs baseline ↓ |",
			"| --- | ---: | ---: | ---: | ---: |")
		fastest, smallest := collectorBaseline, collectorBaseline
		for _, name := range collectorNames {
			samples := healthy.Groups[collectorGroupKey{name, size}]
			rate, rss := collectorMedian(samples, "records_per_second"), collectorMedian(samples, "peak_rss_bytes")
			lines = append(lines, fmt.Sprintf("| %s | %s | %s | %s | %s |", collectorLabel(name), collectorNumber(rate, 0), collectorRatio(rate, baseRate), collectorNumber(rss/1048576, 2), collectorRatio(rss, baseRSS)))
			if rate > collectorMedian(healthy.Groups[collectorGroupKey{fastest, size}], "records_per_second") {
				fastest = name
			}
			if rss < collectorMedian(healthy.Groups[collectorGroupKey{smallest, size}], "peak_rss_bytes") {
				smallest = name
			}
		}
		lines = append(lines, "", collectorLabels[fastest]+" has the highest measured median delivery rate; "+collectorLabels[smallest]+" has the lowest measured median peak RSS.",
			"", "<details>", "<summary>Observed ranges and supporting measurements</summary>", "",
			"| Collector | Records/s range | Peak RSS range (MiB) |", "| --- | ---: | ---: |")
		for _, name := range collectorNames {
			samples := healthy.Groups[collectorGroupKey{name, size}]
			lines = append(lines, fmt.Sprintf("| %s | %s | %s |", collectorLabels[name], collectorObservedRange(samples, "records_per_second", 1, 0), collectorObservedRange(samples, "peak_rss_bytes", 1.0/1048576, 2)))
		}
		lines = append(lines, "", "| Collector | Input MiB/s | CPU/burst (ms) | Producer (ms) | HTTP requests |", "| --- | ---: | ---: | ---: | ---: |")
		metrics := []collectorMetricDef{{"input_mib_per_second", 1, 2}, {"cpu_seconds", 1000, 2}, {"producer_seconds", 1000, 2}, {"http_requests", 1, 0}}
		for _, name := range collectorNames {
			samples := healthy.Groups[collectorGroupKey{name, size}]
			values := make([]string, 0, len(metrics))
			for _, metric := range metrics {
				values = append(values, collectorNumber(collectorMedian(samples, metric.Name)*metric.Scale, metric.Places))
			}
			lines = append(lines, "| "+collectorLabels[name]+" | "+strings.Join(values, " | ")+" |")
		}
		lines = append(lines, "", "CPU counters have 10 ms resolution on this host; zero is below one observable",
			"counter increment. Producer completion does not establish ingestion or durability.", "", "</details>")
	}

	lines = append(lines, "", "### Recovery after a "+collectorNumber(outageEnv.OutageSeconds, 0)+"-second receiver outage", "",
		"Recovery runs from receiver availability to the final unique-record acknowledgement.",
		"**Recovery speed = Fluent Bit recovery time / collector recovery time.**",
		"It includes retry scheduling and drain time; it does not measure crash recovery.",
		"A dash means a zero-duration result cannot provide a finite speed ratio.")
	for _, size := range sizes {
		baseRecovery := collectorMedian(outage.Groups[collectorGroupKey{collectorBaseline, size}], "recovery_drain_seconds")
		lines = append(lines, "", "#### "+collectorGrouped(strconv.FormatInt(size, 10))+"-byte records", "",
			"| Collector | Recovery (s) ↓ | Recovery speed vs baseline ↑ |", "| --- | ---: | ---: |")
		for _, name := range collectorNames {
			recovery := collectorMedian(outage.Groups[collectorGroupKey{name, size}], "recovery_drain_seconds")
			speed := "—"
			if baseRecovery != 0 {
				speed = collectorRatio(baseRecovery, recovery)
			}
			lines = append(lines, fmt.Sprintf("| %s | %s | %s |", collectorLabel(name), collectorNumber(recovery, 3), speed))
		}
		order := slices.Clone(collectorNames)
		slices.SortStableFunc(order, func(a, b string) int {
			left := collectorMedian(outage.Groups[collectorGroupKey{a, size}], "recovery_drain_seconds")
			right := collectorMedian(outage.Groups[collectorGroupKey{b, size}], "recovery_drain_seconds")
			if left < right {
				return -1
			}
			if left > right {
				return 1
			}
			return 0
		})
		first, second := order[0], order[1]
		summary := collectorLabels[first] + " has the lowest measured median recovery time."
		_, firstHigh := collectorBounds(outage.Groups[collectorGroupKey{first, size}], "recovery_drain_seconds")
		secondLow, _ := collectorBounds(outage.Groups[collectorGroupKey{second, size}], "recovery_drain_seconds")
		if firstHigh >= secondLow {
			summary += " Its observed range overlaps " + collectorLabels[second] + "'s; these samples do not establish a consistent ordering between them."
		}
		lines = append(lines, "", summary, "", "<details>", "<summary>Observed recovery ranges and retry activity</summary>", "",
			"| Collector | Recovery range (s) | Producer (ms) | HTTP 503 requests |", "| --- | ---: | ---: | ---: |")
		for _, name := range collectorNames {
			samples := outage.Groups[collectorGroupKey{name, size}]
			lines = append(lines, fmt.Sprintf("| %s | %s | %s | %s |", collectorLabels[name], collectorObservedRange(samples, "recovery_drain_seconds", 1, 3), collectorNumber(collectorMedian(samples, "producer_seconds")*1000, 2), collectorNumber(collectorMedian(samples, "retry_requests"), 0)))
		}
		lines = append(lines, "", "</details>")
	}

	lines = append(lines, "", "### Direct receiver control", "",
		"The control bypasses collector processes, stdin attachment, and disk buffering.",
		"It checks receiver capacity and is excluded from collector speed ratios.", "",
		"| Record size | Records/s ↑ | Input MiB/s ↑ | Outage recovery (s) ↓ |", "| --- | ---: | ---: | ---: |")
	for _, size := range sizes {
		samples := healthy.Groups[collectorGroupKey{"http-control", size}]
		lines = append(lines, fmt.Sprintf("| %s B | %s | %s | %s |", collectorGrouped(strconv.FormatInt(size, 10)), collectorNumber(collectorMedian(samples, "records_per_second"), 0), collectorNumber(collectorMedian(samples, "input_mib_per_second"), 2), collectorNumber(collectorMedian(outage.Groups[collectorGroupKey{"http-control", size}], "recovery_drain_seconds"), 3)))
	}
	var duplicates int64
	for _, report := range []collectorReport{healthy, outage} {
		for _, samples := range report.Groups {
			for _, sample := range samples {
				if sample.Duplicates > math.MaxInt64-duplicates {
					return "", fmt.Errorf("collector duplicate total exceeds the supported count")
				}
				duplicates += sample.Duplicates
			}
		}
	}
	collectorCPU := env.CollectorCPU
	if collectorCPU == "" {
		collectorCPU = "unrestricted"
	}
	lines = append(lines, "", fmt.Sprintf("**Delivery checks:** %d healthy and %d outage samples completed; %s duplicate records. Every outage sample encountered HTTP 503 before successful delivery.", healthy.Count, healthy.Count, collectorGrouped(strconv.FormatInt(duplicates, 10))),
		"", "<details>", "<summary>Measurement environment and source evidence</summary>", "",
		"Healthy run: "+env.Date+"; outage run: "+outageEnv.Date+".",
		fmt.Sprintf("%s, %s/%s, %s; kernel %s.", env.Go, env.OS, env.Arch, env.CPU, env.Kernel),
		fmt.Sprintf("Collector CPU: `%s`; harness CPUs: `%s`; buffer filesystem: `%s`.", collectorCPU, env.HarnessCPUs, env.StorageFilesystem),
		"Container manager: "+env.ContainerManager+"; slogx build: "+env.SlogxGo+".",
		"", "[Healthy samples](../results/healthy.jsonl) · [Outage samples](../results/outage.jsonl) ·",
		"[Host environment](../results/environment.txt) · [Module versions](../results/collectors-modules.txt) ·",
		"[Source checksums](../results/source.sha256)", "",
		"The sample files retain configured image digests, resolved runtime image IDs,",
		"collector module and binary identities, all measurements, and completion records.",
		"", "</details>", "", endMarker)
	return strings.Join(lines, "\n"), nil
}
