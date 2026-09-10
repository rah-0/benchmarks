package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/client"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type collectorSpec struct {
	Name       string
	Image      string
	Config     string
	Args       []string
	Binary     string
	Entrypoint []string
}

type processResources struct {
	CPUSeconds   float64
	PeakRSSBytes int64
}

type collectorProcess struct {
	Input   io.WriteCloser
	ImageID string

	container  testcontainers.Container
	client     *testcontainers.DockerClient
	attachment *client.HijackedResponse
	done       chan struct{}
	streamErr  error
	logs       diagnosticTail
	pid        int
	clockTick  float64
	startTick  uint64
	closeOnce  sync.Once
	closeErr   error
}

// startCollector starts a pinned image with private, bind-mounted storage.
// Testcontainers pulls missing images and copies files before starting the
// process. Initialization and process discovery precede measurement.
func startCollector(ctx context.Context, spec collectorSpec, storageDir, cpu string) (_ *collectorProcess, result error) {
	storageDir, err := filepath.Abs(storageDir)
	if err != nil {
		return nil, err
	}
	dataDir := filepath.Join(storageDir, "data")
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		return nil, err
	}
	clockOutput, err := exec.CommandContext(ctx, "getconf", "CLK_TCK").Output()
	if err != nil {
		return nil, fmt.Errorf("read process clock frequency: %w", err)
	}
	clockTick, err := strconv.ParseUint(strings.TrimSpace(string(clockOutput)), 10, 64)
	if err != nil || clockTick == 0 {
		return nil, fmt.Errorf("invalid process clock frequency %q", clockOutput)
	}
	p := &collectorProcess{clockTick: float64(clockTick)}
	defer func() {
		if result != nil {
			result = errors.Join(result, p.cleanup())
		}
	}()
	req := testcontainers.ContainerRequest{
		Name:       "collector-bench-" + rand.Text(),
		Image:      spec.Image,
		Cmd:        spec.Args,
		Entrypoint: spec.Entrypoint,
		Env:        map[string]string{"GOMAXPROCS": "1", "VECTOR_THREADS": "1"},
		Files: []testcontainers.ContainerFile{{
			Reader:            strings.NewReader(spec.Config),
			ContainerFilePath: "/etc/collector-bench.conf",
			FileMode:          0o444,
		}},
		ConfigModifier: func(config *container.Config) {
			config.User = fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid())
			config.OpenStdin, config.AttachStdin, config.StdinOnce = true, true, true
			config.AttachStdout, config.AttachStderr = true, true
			config.Tty = false
		},
		HostConfigModifier: func(config *container.HostConfig) {
			config.NetworkMode = "host"
			config.CpusetCpus = cpu
			config.Mounts = append(config.Mounts, mount.Mount{Type: mount.TypeBind, Source: dataDir, Target: "/data"})
		},
	}
	if spec.Binary != "" {
		req.Files = append(req.Files, testcontainers.ContainerFile{
			HostFilePath: spec.Binary, ContainerFilePath: "/slogx-collector", FileMode: 0o555,
		})
		if len(req.Entrypoint) == 0 {
			req.Entrypoint = []string{"/slogx-collector"}
		}
	}
	p.container, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          false,
		Logger:           log.New(&p.logs, "", 0),
	})
	if err != nil {
		return nil, p.failure("create "+spec.Name, err)
	}
	p.client, err = testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		return nil, p.failure("create attachment client", err)
	}
	attached, err := p.client.ContainerAttach(ctx, p.container.GetContainerID(), client.ContainerAttachOptions{
		Stream: true, Stdin: true, Stdout: true, Stderr: true,
	})
	if err != nil {
		return nil, p.failure("attach "+spec.Name, err)
	}
	// The hijacked stream is explicitly owned until cleanup. Closing only its
	// write half delivers EOF while preserving stdout and stderr during drain.
	p.attachment = &attached.HijackedResponse
	if _, ok := p.attachment.Conn.(client.CloseWriter); !ok {
		return nil, errors.New("Docker attachment does not support stdin EOF without closing output")
	}
	p.Input = &collectorInput{response: p.attachment}
	p.done = make(chan struct{})
	go func() {
		_, p.streamErr = stdcopy.StdCopy(&p.logs, &p.logs, p.attachment.Reader)
		close(p.done)
	}()
	if err := p.container.Start(ctx); err != nil {
		return nil, p.failure("start "+spec.Name, err)
	}
	inspect, err := p.container.Inspect(ctx)
	if err != nil {
		return nil, p.failure("inspect "+spec.Name, err)
	}
	if inspect.State == nil || !inspect.State.Running || inspect.State.Pid <= 0 {
		return nil, p.failure("start "+spec.Name, errors.New("collector exited before readiness"))
	}
	p.pid, p.ImageID = inspect.State.Pid, inspect.Image
	stat, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", p.pid))
	if err != nil {
		return nil, p.failure("read collector process", err)
	}
	_, p.startTick, err = parseProcessStat(stat)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// Resources reports CPU for the process and its threads and lifetime peak RSS.
// CPU differences between two snapshots exclude startup and cleanup work.
func (p *collectorProcess) Resources() (processResources, error) {
	stat, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", p.pid))
	if err != nil {
		return processResources{}, err
	}
	cpuTicks, startTick, err := parseProcessStat(stat)
	if err != nil {
		return processResources{}, err
	}
	if startTick != p.startTick {
		return processResources{}, errors.New("collector process identity changed")
	}
	status, err := os.ReadFile(fmt.Sprintf("/proc/%d/status", p.pid))
	if err != nil {
		return processResources{}, err
	}
	peakRSS, err := parsePeakRSS(status)
	if err != nil {
		return processResources{}, err
	}
	return processResources{CPUSeconds: float64(cpuTicks) / p.clockTick, PeakRSSBytes: peakRSS}, nil
}

// Stop closes input, requires graceful exit before ctx expires, and removes only
// this sample's container. Every failure includes the bounded diagnostic tail.
func (p *collectorProcess) Stop(ctx context.Context) error {
	defer p.Close()
	if err := p.Input.Close(); err != nil && !errors.Is(err, os.ErrClosed) {
		return p.failure("close collector input", err)
	}
	if err := wait.ForExit().WithPollInterval(25*time.Millisecond).WaitUntilReady(ctx, p.container); err != nil {
		return p.failure("drain collector", err)
	}
	inspect, err := p.container.Inspect(ctx)
	if err != nil {
		return p.failure("inspect collector exit", err)
	}
	if inspect.State == nil {
		return p.failure("collector exit", errors.New("container state is missing"))
	}
	if inspect.State.Running || inspect.State.ExitCode != 0 {
		return p.failure("collector exit", fmt.Errorf("running=%t, exit code=%d", inspect.State.Running, inspect.State.ExitCode))
	}
	select {
	case <-ctx.Done():
		return p.failure("drain collector output", ctx.Err())
	case <-p.done:
	}
	if p.streamErr != nil {
		return p.failure("read collector output", p.streamErr)
	}
	return p.cleanup()
}

// Close is safe after partial startup and after Stop. It cannot remove another
// sample's container or any pre-existing Docker resources.
func (p *collectorProcess) Close() {
	p.cleanup()
}

func (p *collectorProcess) cleanup() error {
	p.closeOnce.Do(func() {
		if p.attachment != nil {
			// Closing both halves unblocks any input writer and log reader even
			// when startup failed or the sample context has been cancelled.
			p.attachment.Close()
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if p.container != nil {
			if err := p.container.Terminate(ctx, testcontainers.StopTimeout(0)); err != nil {
				p.closeErr = p.failure("terminate collector", err)
			}
		}
		if p.done != nil {
			select {
			case <-p.done:
			case <-ctx.Done():
				p.closeErr = errors.Join(p.closeErr, errors.New("collector output did not stop after closing attachment"))
			}
		}
		if p.client != nil {
			p.closeErr = errors.Join(p.closeErr, p.client.Close())
		}
	})
	return p.closeErr
}

type collectorInput struct {
	response  *client.HijackedResponse
	closeOnce sync.Once
	closeErr  error
}

func (input *collectorInput) Write(data []byte) (int, error) {
	return input.response.Conn.Write(data)
}

func (input *collectorInput) Close() error {
	input.closeOnce.Do(func() { input.closeErr = input.response.CloseWrite() })
	return input.closeErr
}

func (p *collectorProcess) failure(action string, err error) error {
	return fmt.Errorf("%s: %w\n%s", action, err, p.logs.String())
}

func parseProcessStat(data []byte) (cpuTicks, startTick uint64, result error) {
	// comm can contain spaces and parentheses; its final ')' ends field two.
	_, rest, found := strings.Cut(string(data), "(")
	end := strings.LastIndexByte(rest, ')')
	if !found || end < 0 {
		return 0, 0, errors.New("process stat has no command name")
	}
	fields := strings.Fields(rest[end+1:])
	if len(fields) < 20 {
		return 0, 0, errors.New("process stat is missing CPU or start fields")
	}
	user, err := strconv.ParseUint(fields[11], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parse process user ticks: %w", err)
	}
	system, err := strconv.ParseUint(fields[12], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parse process system ticks: %w", err)
	}
	start, err := strconv.ParseUint(fields[19], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parse process start ticks: %w", err)
	}
	if user > ^uint64(0)-system {
		return 0, 0, errors.New("process CPU ticks overflow")
	}
	return user + system, start, nil
}

func parsePeakRSS(data []byte) (int64, error) {
	for line := range strings.SplitSeq(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 || fields[0] != "VmHWM:" {
			continue
		}
		if len(fields) != 3 || fields[2] != "kB" {
			return 0, errors.New("invalid process peak RSS units")
		}
		value, err := strconv.ParseUint(fields[1], 10, 53)
		if err != nil {
			return 0, fmt.Errorf("parse process peak RSS: %w", err)
		}
		return int64(value) * 1024, nil
	}
	return 0, errors.New("process status is missing peak RSS")
}

const diagnosticLimit = 64 << 10

type diagnosticTail struct {
	mu   sync.Mutex
	data []byte
}

func (t *diagnosticTail) Write(data []byte) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	n := len(data)
	if n >= diagnosticLimit {
		t.data = append(t.data[:0], data[n-diagnosticLimit:]...)
	} else {
		if excess := len(t.data) + n - diagnosticLimit; excess > 0 {
			copy(t.data, t.data[excess:])
			t.data = t.data[:len(t.data)-excess]
		}
		t.data = append(t.data, data...)
	}
	return n, nil
}

func (t *diagnosticTail) String() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return string(t.data)
}
