#!/usr/bin/env bash
set -euo pipefail
umask 077

cd "$(dirname "${BASH_SOURCE[0]}")/.."
output_dir="${1:-results}"
mkdir -p "$output_dir"
output_dir="$(cd "$output_dir" && pwd)"
destination="$output_dir"
# Keep unrestricted diagnostics outside the checkout. Only sanitized, complete
# matrices are copied to the publication directory.
output_dir="$(mktemp -d /tmp/benchmark-results.XXXXXXXX)"
bench_time="${BENCHTIME:-1s}"
bench_count="${BENCHCOUNT:-3}"
export GOGC=100 GOMEMLIMIT=off
storage_dir="$(mktemp -d "${BENCH_STORAGE_ROOT:-/var/tmp}/benchmarks.XXXXXXXX")"
trap 'rm -rf -- "$storage_dir"; if [[ -d "$output_dir" ]]; then printf "Private diagnostics retained in %s\n" "$output_dir" >&2; fi' EXIT
export TMPDIR="$storage_dir"
report_tool="$storage_dir/benchmark-results"
go build -trimpath -o "$report_tool" ./scripts/benchmarks

{
  date -u +'%Y-%m-%dT%H:%M:%SZ'
  go version
  uname -srmo
  command -v lscpu >/dev/null && lscpu
  go env GOOS GOARCH GOAMD64 CGO_ENABLED GOEXPERIMENT GOTOOLCHAIN
  printf 'BENCHTIME=%s\nBENCHCOUNT=%s\nGOMAXPROCS=%s\nGOGC=%s\nGOMEMLIMIT=%s\n' \
    "$bench_time" "$bench_count" "${GOMAXPROCS:-default}" "$GOGC" "$GOMEMLIMIT"
  docker version --format 'Docker server: {{.Server.Version}}'
  findmnt -T "$storage_dir" -no FSTYPE
} > "$output_dir/environment.txt"
go list -m all > "$output_dir/modules.txt"
go -C collectors list -m all > "$output_dir/collectors-modules.txt"
"$report_tool" manifest "$output_dir/source.sha256"

# One package at a time keeps suites from competing for CPU, memory, or disk.
# Database and cache TestMain functions provision and release their containers.
for package in cache compression crypto/encrypt crypto/sign db/mariadb \
  db/postgresql db/sqlite encoding iteration meta protocol/http protocol/tcp \
  serializer validator; do
  printf 'Benchmarking %s\n' "$package"
  go test -p=1 -run '^$' -bench . -benchmem -benchtime="$bench_time" \
    -count="$bench_count" -timeout=2h "./$package" \
    > "$output_dir/${package//\//-}.txt" 2>&1
done

# Logging is serial; retain the established five-sample, one-CPU methodology.
printf 'Benchmarking log\n'
GOMAXPROCS=1 taskset -c "${LOG_CPU:-0}" \
  go test -p=1 -run '^$' -bench . -benchmem -benchtime="$bench_time" \
    -count=5 -cpu=1 -timeout=30m ./log > "$output_dir/log.txt" 2>&1

# Collector buffers must use persistent storage, not the output directory's FS.
printf 'Benchmarking collectors: healthy delivery\n'
taskset -c "${HARNESS_CPU:-2}" go -C collectors run . \
  -cpu="${COLLECTOR_CPU:-0}" > "$output_dir/healthy.jsonl"
printf 'Benchmarking collectors: two-second outage\n'
taskset -c "${HARNESS_CPU:-2}" go -C collectors run . \
  -cpu="${COLLECTOR_CPU:-0}" -outage=2s > "$output_dir/outage.jsonl"
"$report_tool" sanitize "$output_dir"
mv -- "$output_dir"/* "$destination/"
rmdir -- "$output_dir"
printf 'Complete. Results: %s\n' "$destination"
