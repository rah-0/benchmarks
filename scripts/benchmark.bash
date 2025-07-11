#!/bin/bash
cd ..
go install github.com/rah-0/testmark@latest

# Function to run benchmarks on a specific function and generate profiling outputs
run_benchmark() {
  local profile_name="$1"  # Name for pprof output files (e.g., "custom_profile")
  local bench_time="$2"    # Benchmark duration (e.g., "5s", "60s")
  local bench_target="$3"  # Specific benchmark function (e.g., "BenchmarkMyFunction")
  local bench_dir="$4"     # Directory for go test

  local mem_profile="./pprof_svg/${profile_name}_mem.out"
  local cpu_profile="./pprof_svg/${profile_name}_cpu.out"
  local mem_svg="./pprof_svg/${profile_name}_mem.svg"
  local cpu_svg="./pprof_svg/${profile_name}_cpu.svg"

  mkdir -p ./pprof_svg

  go test -run=^$ -bench="^${bench_target}$" -benchmem "${bench_dir}" -benchtime="${bench_time}" -timeout=0 \
    -memprofile="${mem_profile}" -cpuprofile="${cpu_profile}" | testmark | grep -E '^Benchmark'

  go tool pprof -svg -output="${mem_svg}" "${mem_profile}" >/dev/null 2>&1
  go tool pprof -svg -output="${cpu_svg}" "${cpu_profile}" >/dev/null 2>&1
}

rm -rf ./pprof_svg/*

# Example usage:
time="60s"

run_benchmark "jwt-gen" "${time}" "BenchmarkJWTGeneration" "./crypto/sign"
run_benchmark "jwt-verify" "${time}" "BenchmarkJWTVerification" "./crypto/sign"
run_benchmark "aes-enc" "${time}" "BenchmarkGenerateEncryptedToken" "./crypto/encrypt"
run_benchmark "aes-dec" "${time}" "BenchmarkDecryptEncryptedToken" "./crypto/encrypt"

rm ./*.test > /dev/null 2>&1
