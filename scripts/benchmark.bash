#!/bin/bash
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
time="10s"

# # Full consume
# run_benchmark "BenchmarkPullChan_Full_100k"            "$time" "BenchmarkPullChan_Full_100k"            "./iteration"
# run_benchmark "BenchmarkPullSlice_Full_100k"           "$time" "BenchmarkPullSlice_Full_100k"           "./iteration"
# run_benchmark "BenchmarkBaseline_ChanDirect_Full_100k" "$time" "BenchmarkBaseline_ChanDirect_Full_100k" "./iteration"
# run_benchmark "BenchmarkBaseline_SliceDirect_Full_100k""$time" "BenchmarkBaseline_SliceDirect_Full_100k" "./iteration"
# # Early close
# run_benchmark "BenchmarkPullChan_EarlyClose10_100k"    "$time" "BenchmarkPullChan_EarlyClose10_100k"    "./iteration"
# run_benchmark "BenchmarkPullSlice_EarlyClose10_100k"   "$time" "BenchmarkPullSlice_EarlyClose10_100k"   "./iteration"
# # Time to first row
# run_benchmark "BenchmarkPullChan_TTFR_1M"              "$time" "BenchmarkPullChan_TTFR_1M"              "./iteration"
# run_benchmark "BenchmarkPullSlice_TTFR_1M"             "$time" "BenchmarkPullSlice_TTFR_1M"             "./iteration"
# # Slow consumer
# run_benchmark "BenchmarkPullChan_SlowConsumerWork50_100k"  "$time" "BenchmarkPullChan_SlowConsumerWork50_100k"  "./iteration"
# run_benchmark "BenchmarkPullSlice_SlowConsumerWork50_100k" "$time" "BenchmarkPullSlice_SlowConsumerWork50_100k" "./iteration"
# # Slow producer
# run_benchmark "BenchmarkPullChan_SlowProducerWork50_100k"  "$time" "BenchmarkPullChan_SlowProducerWork50_100k"  "./iteration"
# run_benchmark "BenchmarkPullSlice_SlowProducerWork50_100k" "$time" "BenchmarkPullSlice_SlowProducerWork50_100k" "./iteration"
# # Multi-range
# run_benchmark "BenchmarkPullChan_MultiRange_100k_10parts"  "$time" "BenchmarkPullChan_MultiRange_100k_10parts"  "./iteration"
# run_benchmark "BenchmarkPullSlice_MultiRange_100k_10parts" "$time" "BenchmarkPullSlice_MultiRange_100k_10parts" "./iteration"
# # Empty ranges
# run_benchmark "BenchmarkPullChan_EmptyRange"           "$time" "BenchmarkPullChan_EmptyRange"           "./iteration"
# run_benchmark "BenchmarkPullSlice_EmptyRange"          "$time" "BenchmarkPullSlice_EmptyRange"          "./iteration"
# # Parallel
# run_benchmark "BenchmarkPullChan_Parallel_10k"         "$time" "BenchmarkPullChan_Parallel_10k"         "./iteration"
# run_benchmark "BenchmarkPullSlice_Parallel_10k"        "$time" "BenchmarkPullSlice_Parallel_10k"        "./iteration"

# ============================================================================
# LOGGING LIBRARY BENCHMARKS
# ============================================================================

# Simple message logging
run_benchmark "BenchmarkNabu_SimpleMessage"    "$time" "BenchmarkNabu_SimpleMessage"    "./log"
run_benchmark "BenchmarkZerolog_SimpleMessage" "$time" "BenchmarkZerolog_SimpleMessage" "./log"
run_benchmark "BenchmarkZap_SimpleMessage"     "$time" "BenchmarkZap_SimpleMessage"     "./log"
run_benchmark "BenchmarkLogrus_SimpleMessage"  "$time" "BenchmarkLogrus_SimpleMessage"  "./log"
run_benchmark "BenchmarkSlog_SimpleMessage"    "$time" "BenchmarkSlog_SimpleMessage"    "./log"

# Message with fields
run_benchmark "BenchmarkNabu_WithFields"    "$time" "BenchmarkNabu_WithFields"    "./log"
run_benchmark "BenchmarkZerolog_WithFields" "$time" "BenchmarkZerolog_WithFields" "./log"
run_benchmark "BenchmarkZap_WithFields"     "$time" "BenchmarkZap_WithFields"     "./log"
run_benchmark "BenchmarkLogrus_WithFields"  "$time" "BenchmarkLogrus_WithFields"  "./log"
run_benchmark "BenchmarkSlog_WithFields"    "$time" "BenchmarkSlog_WithFields"    "./log"

# Error logging
run_benchmark "BenchmarkNabu_Error"    "$time" "BenchmarkNabu_Error"    "./log"
run_benchmark "BenchmarkZerolog_Error" "$time" "BenchmarkZerolog_Error" "./log"
run_benchmark "BenchmarkZap_Error"     "$time" "BenchmarkZap_Error"     "./log"
run_benchmark "BenchmarkLogrus_Error"  "$time" "BenchmarkLogrus_Error"  "./log"
run_benchmark "BenchmarkSlog_Error"    "$time" "BenchmarkSlog_Error"    "./log"

# Error with fields
run_benchmark "BenchmarkNabu_ErrorWithFields"    "$time" "BenchmarkNabu_ErrorWithFields"    "./log"
run_benchmark "BenchmarkZerolog_ErrorWithFields" "$time" "BenchmarkZerolog_ErrorWithFields" "./log"
run_benchmark "BenchmarkZap_ErrorWithFields"     "$time" "BenchmarkZap_ErrorWithFields"     "./log"
run_benchmark "BenchmarkLogrus_ErrorWithFields"  "$time" "BenchmarkLogrus_ErrorWithFields"  "./log"
run_benchmark "BenchmarkSlog_ErrorWithFields"    "$time" "BenchmarkSlog_ErrorWithFields"    "./log"

# Error chain
run_benchmark "BenchmarkNabu_ErrorChain"    "$time" "BenchmarkNabu_ErrorChain"    "./log"
run_benchmark "BenchmarkZerolog_ErrorChain" "$time" "BenchmarkZerolog_ErrorChain" "./log"
run_benchmark "BenchmarkZap_ErrorChain"     "$time" "BenchmarkZap_ErrorChain"     "./log"
run_benchmark "BenchmarkLogrus_ErrorChain"  "$time" "BenchmarkLogrus_ErrorChain"  "./log"
run_benchmark "BenchmarkSlog_ErrorChain"    "$time" "BenchmarkSlog_ErrorChain"    "./log"

rm ./*.test > /dev/null 2>&1
