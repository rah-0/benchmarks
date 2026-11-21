#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "${BASH_SOURCE[0]}")/.."
go test -run '^$' -bench . -benchmem -benchtime="${BENCHTIME:-1s}" \
  -count="${BENCHCOUNT:-3}" ./serializer "$@"
