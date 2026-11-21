# Benchmark scripts

Run these commands from the repository root. Bash runs the benchmark suites;
Go maintains their samples, comparison groups, and Markdown reports. The tools
use the standard library and require no additional scripting runtime.

```sh
# Run the complete matrix, then refresh every category report.
bash scripts/benchmark.bash
bash scripts/results.bash render

# Refresh one category, or check the saved reports without changing files.
bash scripts/results.bash render meta
bash scripts/results.bash render --check

# Verify that saved samples contain only publishable measurement output.
bash scripts/results.bash sanitize --check results
```

`results.bash` invokes `go run ./scripts/benchmarks` from the repository root.
The suite runner builds that tool once before measurement and uses it to record
source checksums and sanitize completed output before publication.

## Maintaining comparisons

[groups.go](benchmarks/groups.go) defines each workload's display names,
baseline, input sizes, and operation boundaries. Add every new Go benchmark to
an appropriate group; unmapped cases, missing baselines, and incompatible CPU
settings stop report generation. Collector comparisons and their validation
live in [collectors.go](benchmarks/collectors.go).

The renderer reads `results/` and replaces only marked result blocks in the
category READMEs. It validates every requested category before writing any
file. Sanitization preserves measurement lines and rejects failed runs, local
module replacements, and unexpected artifacts. Diagnostic files stay outside
the repository.

The `manifest <file>` command hashes benchmark Go source, module files, and
collector configuration using relative paths. Reporting tools under `scripts/`
are excluded because they do not participate in measured operations.

To verify the reporting tools:

```sh
go test -count=1 -race -cover -covermode=atomic ./scripts/benchmarks
go vet ./scripts/benchmarks
```
