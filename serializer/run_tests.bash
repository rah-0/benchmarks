#!/bin/bash
# Function to run benchmarks and generate profiling outputs
run_benchmarks() {
  local go_file="$1"      # First parameter: the Go file (e.g., "bytes.go")
  local test_file="$2"    # Second parameter: the test file (e.g., "bytes_test.go")
  local profile_name="$3" # Third parameter: name for pprof output files (e.g., "custom_profile")
  local bench_time="5s"

  if [[ -z "$go_file" || -z "$test_file" || -z "$profile_name" ]]; then
    echo "Usage: run_benchmarks <go_file> <test_file> <profile_name>"
    return 1
  fi

  # Construct profiling output file paths based on profile_name
  local mem_profile="/tmp/${profile_name}_mem.out"
  local cpu_profile="/tmp/${profile_name}_cpu.out"
  local mem_svg="./pprof_svg/${profile_name}_mem.svg"
  local cpu_svg="./pprof_svg/${profile_name}_cpu.svg"

  # Ensure the output directory exists
  mkdir -p ./pprof_svg

  # Run the benchmark and generate profiling outputs
  go test -bench=. -benchmem ./data.go ./"${go_file}" ./"${test_file}" -benchtime="${bench_time}" -timeout=0 -memprofile="${mem_profile}" -cpuprofile="${cpu_profile}" \
    && go tool pprof -svg -output="${mem_svg}" "${mem_profile}" \
    && go tool pprof -svg -output="${cpu_svg}" "${cpu_profile}"

  # Check for errors
  if [[ $? -eq 0 ]]; then
    echo "Benchmarks and profiles generated successfully."
    echo "Memory profile saved to: ${mem_svg}"
    echo "CPU profile saved to: ${cpu_svg}"
  else
    echo "An error occurred during benchmarks or profile generation."
  fi
}

rm -rf ./pprof_svg/*

run_benchmarks "bytes.go" "bytes_test.go"     "bytes"
run_benchmarks "bytes.go" "bytes_opt_test.go" "bytes_opt"
run_benchmarks "json.go"  "json_test.go"      "json"
run_benchmarks "proto.go" "proto_test.go"     "proto"

rm serializer.test
