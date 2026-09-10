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
	slogxImage     = "debian:bookworm-slim@sha256:5ae3c39ebd15e229dcedd5cee596b2497182493d41ff162e824ba13fc1b2b867"
	vectorImage    = "timberio/vector:0.58.0-debian@sha256:9a68d40824ea5285a716bbe1872dba3cafdc7899deae4821d0c414416f7a8648"
	fluentBitImage = "cr.fluentbit.io/fluent/fluent-bit:5.1.2@sha256:71cda445290efc2d45d565c12e0de4b15aa0182510276ad50c317aae1423d7ee"
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
