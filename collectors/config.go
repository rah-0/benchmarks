package main

import (
	_ "embed"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Runtime images are pinned to Linux/amd64 manifests.
const (
	slogxImage     = "debian:trixie-slim@sha256:7378260c0a098d27ac206225db5e7d92b8c7c84ea7afccc86eeb513756bbd392"
	vectorImage    = "timberio/vector:0.62.0-debian@sha256:6ac7a9272e529e4d111ad30493a5a35c290b3d7e1c33f2f2df306c5fad1e64d3"
	fluentBitImage = "cr.fluentbit.io/fluent/fluent-bit:5.2.0@sha256:de3d2860ffaf6a7a3d9f26c80d184b49eb8e08d649730acac4ec2d9cfb671c5b"
)

//go:embed config/vector.toml
var vectorConfig string

//go:embed config/fluent-bit.conf
var fluentBitConfig string

func collectorSpecs(endpoint string) ([]collectorSpec, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("parse collector endpoint: %w", err)
	}
	if u.Scheme != "http" || u.Hostname() == "" || u.User != nil || u.Fragment != "" || strings.ContainsAny(endpoint, "\r\n\t") {
		return nil, fmt.Errorf("collector endpoint must be an HTTP URL without credentials, fragments, or control characters")
	}
	port, err := strconv.Atoi(u.Port())
	if err != nil || port < 1 || port > 65535 {
		return nil, fmt.Errorf("collector endpoint must include a port between 1 and 65535")
	}
	fluentConfig := strings.NewReplacer(
		"__HOST__", u.Hostname(),
		"__PORT__", strconv.Itoa(port),
		"__URI__", u.RequestURI(),
	).Replace(fluentBitConfig)

	return []collectorSpec{
		{
			Name:       "slogx-collector",
			Image:      slogxImage,
			Entrypoint: []string{"/slogx-collector"},
			Args: []string{
				"-destination", "openobserve",
				"-endpoint", endpoint,
				"-journal-dir", "/data",
				"-timestamp-field=",
				"-batch-size", "500",
				"-batch-bytes", "4194304",
				"-flush-interval", "100ms",
				"-retry-interval", "1s",
				"-max-retry-interval", "1s",
				"-shutdown-timeout", "10s",
			},
		},
		{
			Name:   "vector",
			Image:  vectorImage,
			Config: strings.ReplaceAll(vectorConfig, "__ENDPOINT__", strconv.Quote(u.String())),
			Args:   []string{"--config-toml", "/etc/collector-bench.conf", "--graceful-shutdown-limit-secs", "10", "--quiet"},
		},
		{
			Name:   "fluent-bit",
			Image:  fluentBitImage,
			Config: fluentConfig,
			Args:   []string{"--config", "/etc/collector-bench.conf"},
		},
	}, nil
}
