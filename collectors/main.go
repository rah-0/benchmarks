// Command collectors compares external log collectors through a checked HTTP sink.
package main

import (
	"bufio"
	"bytes"
	"context"
	"debug/buildinfo"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type options struct {
	Records     int
	Sizes       string
	Repeat      int
	Collectors  string
	StorageRoot string
	CPU         string
	Timeout     time.Duration
	Outage      time.Duration
	Binary      string
}

type sample struct {
	Collector            string   `json:"collector"`
	Image                string   `json:"image,omitempty"`
	ImageID              string   `json:"image_id,omitempty"`
	Sample               int      `json:"sample"`
	Records              int      `json:"records"`
	RecordBytes          int      `json:"record_bytes"`
	InputBytes           int      `json:"input_bytes"`
	Seconds              float64  `json:"seconds"`
	RecordsPerSecond     float64  `json:"records_per_second"`
	InputMiBPerSecond    float64  `json:"input_mib_per_second"`
	ProducerSeconds      float64  `json:"producer_seconds"`
	HTTPRequests         int      `json:"http_requests"`
	RetryRequests        int      `json:"retry_requests"`
	Duplicates           int      `json:"duplicates"`
	CPUSeconds           *float64 `json:"cpu_seconds,omitempty"`
	PeakRSSBytes         *int64   `json:"peak_rss_bytes,omitempty"`
	OutageSeconds        float64  `json:"outage_seconds"`
	RecoveryDrainSeconds float64  `json:"recovery_drain_seconds,omitempty"`
}

func main() {
	var opts options
	flag.IntVar(&opts.Records, "records", 2000, "Measured records per sample")
	flag.StringVar(&opts.Sizes, "sizes", "256,1024", "Comma-separated NDJSON record sizes, including newline")
	flag.IntVar(&opts.Repeat, "repeat", 3, "Samples per collector and size")
	flag.StringVar(&opts.Collectors, "collectors", "http-control,slogx-collector,vector,fluent-bit", "Comma-separated collectors; http-control checks receiver capacity without a collector")
	flag.StringVar(&opts.StorageRoot, "storage-root", "/var/tmp", "Existing persistent filesystem directory for temporary buffers")
	flag.StringVar(&opts.CPU, "cpu", "", "Optional collector CPU affinity, using Docker cpuset syntax")
	flag.DurationVar(&opts.Timeout, "timeout", 5*time.Minute, "Maximum startup, delivery and drain time per sample")
	flag.DurationVar(&opts.Outage, "outage", 0, "Return HTTP 503 for this duration after warmup")
	flag.Parse()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	if err := run(ctx, opts, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, opts options, output io.Writer) error {
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		return errors.New("collector comparisons require Linux/amd64 Docker images, host networking and /proc")
	}
	if opts.Records < 1 || opts.Repeat < 1 || opts.Timeout <= 0 || opts.Outage < 0 || opts.Outage >= opts.Timeout {
		return errors.New("records, repeat and timeout must be positive; outage must be nonnegative and shorter than timeout")
	}
	var sizes []int
	for item := range strings.SplitSeq(opts.Sizes, ",") {
		size, err := strconv.Atoi(strings.TrimSpace(item))
		if err != nil || size < 256 || size > 16<<10 {
			return errors.New("sizes must be integers between 256 and 16384 bytes")
		}
		sizes = append(sizes, size)
	}
	storageRoot, err := filepath.Abs(opts.StorageRoot)
	if err != nil {
		return err
	}
	info, err := os.Stat(storageRoot)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("storage-root must be an existing directory: %s", storageRoot)
	}
	specs, err := collectorSpecs("http://127.0.0.1:1/api/default/bench/_json")
	if err != nil {
		return err
	}
	byName := make(map[string]collectorSpec)
	for _, spec := range specs {
		byName[spec.Name] = spec
	}
	names := strings.Split(opts.Collectors, ",")
	images := make(map[string]string)
	selected := make(map[string]bool)
	for _, name := range names {
		if selected[name] {
			return fmt.Errorf("duplicate collector %q", name)
		}
		selected[name] = true
		if name == "http-control" {
			continue
		}
		spec, ok := byName[name]
		if !ok {
			return fmt.Errorf("unknown collector %q", name)
		}
		images[name] = spec.Image
	}
	var binaryDigest, collectorGo string
	if selected["slogx-collector"] {
		fmt.Fprintln(os.Stderr, "Preparing slogx-collector")
		path, digest, cleanup, err := prepareSlogx(ctx)
		if err != nil {
			return err
		}
		defer cleanup()
		opts.Binary = path
		binaryDigest = digest
		info, err := buildinfo.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read collector build information: %w", err)
		}
		collectorGo = info.GoVersion
	}
	encoder := json.NewEncoder(output)
	meta := map[string]any{
		"type": "environment", "date": time.Now().UTC().Format(time.RFC3339),
		"go": runtime.Version(), "os": runtime.GOOS, "arch": runtime.GOARCH,
		"collector_cpu": opts.CPU, "records": opts.Records, "sizes": sizes,
		"repeat": opts.Repeat, "outage_seconds": opts.Outage.Seconds(),
		"images": images, "collectors": names,
		"container_manager": "testcontainers-go v0.44.0",
	}
	if binaryDigest != "" {
		meta["slogx_module"] = slogxModule
		meta["slogx_binary_sha256"] = binaryDigest
		meta["slogx_go"] = collectorGo
	}
	if kernel, err := os.ReadFile("/proc/sys/kernel/osrelease"); err == nil {
		meta["kernel"] = strings.TrimSpace(string(kernel))
	}
	if status, err := os.ReadFile("/proc/self/status"); err == nil {
		for line := range strings.SplitSeq(string(status), "\n") {
			if value, ok := strings.CutPrefix(line, "Cpus_allowed_list:"); ok {
				meta["harness_cpus"] = strings.TrimSpace(value)
			}
		}
	}
	if cpu, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		for line := range strings.SplitSeq(string(cpu), "\n") {
			if key, value, ok := strings.Cut(line, ":"); ok && strings.TrimSpace(key) == "model name" {
				meta["cpu"] = strings.TrimSpace(value)
				break
			}
		}
	}
	if filesystem, err := exec.CommandContext(ctx, "findmnt", "-n", "-o", "FSTYPE", "--target", storageRoot).Output(); err == nil {
		meta["storage_filesystem"] = strings.TrimSpace(string(filesystem))
	}
	if err := encoder.Encode(meta); err != nil {
		return err
	}
	for _, size := range sizes {
		expected, lines, err := makeRecords(opts.Records, size)
		if err != nil {
			return err
		}
		// Rotate the order between repetitions to reduce fixed-order bias.
		for repetition := 0; repetition < opts.Repeat; repetition++ {
			for offset := range names {
				name := names[(offset+repetition)%len(names)]
				fmt.Fprintf(os.Stderr, "%s: %d records x %d bytes, sample %d/%d\n", name, opts.Records, size, repetition+1, opts.Repeat)
				sampleCtx, cancel := context.WithTimeout(ctx, opts.Timeout)
				result, err := runSample(sampleCtx, opts, storageRoot, name, expected, lines)
				cancel()
				if err != nil {
					return fmt.Errorf("%s, %d bytes, sample %d: %w", name, size, repetition+1, err)
				}
				result.Sample = repetition + 1
				if err := encoder.Encode(result); err != nil {
					return err
				}
			}
		}
	}
	return encoder.Encode(map[string]any{"type": "complete", "samples": len(sizes) * opts.Repeat * len(names)})
}

func runSample(ctx context.Context, opts options, storageRoot, name string, expected []record, lines [][]byte) (sample, error) {
	var result sample
	sink := newReceiver(expected)
	server := httptest.NewUnstartedServer(sink)
	server.Config.ReadTimeout = 15 * time.Second
	server.Config.WriteTimeout = 15 * time.Second
	server.Start()
	defer server.Close()
	endpoint := server.URL + "/api/default/bench/_json"
	var proc *collectorProcess
	var before processResources
	var image string
	client := &http.Client{Timeout: 15 * time.Second}
	defer client.CloseIdleConnections()
	if name == "http-control" {
		if err := postRecords(ctx, client, endpoint, lines[:1]); err != nil {
			return result, err
		}
	} else {
		specs, err := collectorSpecs(endpoint)
		if err != nil {
			return result, err
		}
		var spec collectorSpec
		for _, candidate := range specs {
			if candidate.Name == name {
				spec = candidate
			}
		}
		image = spec.Image
		spec.Binary = opts.Binary
		if name != "slogx-collector" {
			spec.Binary = ""
		}
		dir, err := os.MkdirTemp(storageRoot, "collector-benchmark-")
		if err != nil {
			return result, err
		}
		defer os.RemoveAll(dir)
		proc, err = startCollector(ctx, spec, dir, opts.CPU)
		if err != nil {
			return result, err
		}
		defer proc.Close()
		result.ImageID = proc.ImageID
	}
	warmCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	if proc != nil {
		written := make(chan error, 1)
		go func() {
			_, err := proc.Input.Write(lines[0])
			written <- err
		}()
		var err error
		select {
		case err = <-written:
		case <-warmCtx.Done():
			err = warmCtx.Err()
		}
		if err != nil {
			cancel()
			return result, proc.failure("write warmup", err)
		}
	}
	_, err := sink.wait(warmCtx, true)
	cancel()
	if err != nil {
		if proc != nil {
			return result, proc.failure("warmup", err)
		}
		return result, fmt.Errorf("warmup: %w", err)
	}
	if proc != nil {
		before, err = proc.Resources()
		if err != nil {
			return result, err
		}
	}
	start := time.Now()
	sink.begin(start, opts.Outage)
	type writeResult struct {
		elapsed time.Duration
		err     error
	}
	waitCtx, stopWait := context.WithCancel(ctx)
	defer stopWait()
	writeDone := make(chan writeResult, 1)
	go func() {
		var writeErr error
		if proc == nil {
			batchSize := min(500, (4<<20)/len(lines[1]))
			for offset := 1; offset < len(lines); offset += batchSize {
				if writeErr = postRecords(waitCtx, client, endpoint, lines[offset:min(offset+batchSize, len(lines))]); writeErr != nil {
					break
				}
			}
		} else {
			writer := bufio.NewWriterSize(proc.Input, 64<<10)
			for _, line := range lines[1:] {
				if _, writeErr = writer.Write(line); writeErr != nil {
					break
				}
			}
			if writeErr == nil {
				writeErr = writer.Flush()
			}
		}
		writeDone <- writeResult{time.Since(start), writeErr}
		if writeErr != nil {
			stopWait()
		}
	}()
	stats, waitErr := sink.wait(waitCtx, false)
	if waitErr != nil {
		select {
		case outcome := <-writeDone:
			if outcome.err != nil {
				waitErr = fmt.Errorf("write input: %w", outcome.err)
			}
		default:
		}
		if proc != nil {
			return result, proc.failure("delivery", waitErr)
		}
		return result, waitErr
	}
	var outcome writeResult
	select {
	case outcome = <-writeDone:
	case <-ctx.Done():
		return result, ctx.Err()
	}
	if outcome.err != nil {
		return result, fmt.Errorf("write input: %w", outcome.err)
	}
	if proc != nil {
		after, err := proc.Resources()
		if err != nil {
			return result, err
		}
		cpu := after.CPUSeconds - before.CPUSeconds
		result.CPUSeconds = &cpu
		result.PeakRSSBytes = &after.PeakRSSBytes
		stopCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
		err = proc.Stop(stopCtx)
		cancel()
		if err != nil {
			return result, err
		}
		stats, err = sink.wait(ctx, false)
		if err != nil {
			return result, err
		}
	}
	result.Collector = name
	result.Image = image
	result.Records = opts.Records
	result.RecordBytes = len(lines[1])
	result.InputBytes = opts.Records * result.RecordBytes
	result.Seconds = stats.Finished.Sub(start).Seconds()
	result.RecordsPerSecond = float64(opts.Records) / result.Seconds
	result.InputMiBPerSecond = float64(result.InputBytes) / (1 << 20) / result.Seconds
	result.ProducerSeconds = outcome.elapsed.Seconds()
	result.HTTPRequests = stats.Requests
	result.RetryRequests = stats.RetryRequests
	result.Duplicates = stats.Duplicates
	result.OutageSeconds = opts.Outage.Seconds()
	if opts.Outage > 0 {
		result.RecoveryDrainSeconds = max(0, stats.Finished.Sub(start.Add(opts.Outage)).Seconds())
		if stats.RetryRequests == 0 {
			return result, errors.New("outage was not exercised: no requests received HTTP 503")
		}
	}
	return result, nil
}

// The direct sender provides a receiver-capacity control, not a collector score.
func postRecords(ctx context.Context, client *http.Client, endpoint string, lines [][]byte) error {
	body := bytes.Join(lines, nil)
	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/x-ndjson")
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		_, readErr := io.Copy(io.Discard, resp.Body)
		closeErr := resp.Body.Close()
		if err := errors.Join(readErr, closeErr); err != nil {
			return err
		}
		if resp.StatusCode == http.StatusOK {
			return nil
		}
		if resp.StatusCode != http.StatusServiceUnavailable {
			return fmt.Errorf("receiver HTTP status %d", resp.StatusCode)
		}
		timer := time.NewTimer(time.Second)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
}
