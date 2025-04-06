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

#run_benchmark "mariadb1" "${time}" "BenchmarkMariaDBSingleInsertFixedData" "./db/mariadb"
#run_benchmark "postgres1" "${time}" "BenchmarkPostgresSingleInsertFixedData" "./db/postgresql"
#run_benchmark "sqlite1" "${time}" "BenchmarkSQLiteSingleInsertFixedData" "./db/sqlite"
#run_benchmark "hyperion1" "${time}" "BenchmarkHyperionSingleInsertFixedData" "./db/hyperion"

#run_benchmark "mariadb2" "${time}" "BenchmarkMariaDBSingleInsertRandomData" "./db/mariadb"
#run_benchmark "postgres2" "${time}" "BenchmarkPostgresSingleInsertRandomData" "./db/postgresql"
#run_benchmark "sqlite2" "${time}" "BenchmarkSQLiteSingleInsertRandomData" "./db/sqlite"
#run_benchmark "hyperion2" "${time}" "BenchmarkHyperionSingleInsertRandomData" "./db/hyperion"

#run_benchmark "mariadb3" "${time}" "BenchmarkMariaDBInsert1MilAndFindMiddle" "./db/mariadb"
#run_benchmark "postgres3" "${time}" "BenchmarkPostgresInsert1MilAndFindMiddle" "./db/postgresql"
#run_benchmark "sqlite3" "${time}" "BenchmarkSQLiteInsert1MilAndFindMiddle" "./db/sqlite"
#run_benchmark "hyperion3" "${time}" "BenchmarkHyperionInsert1MilAndFindMiddle" "./db/hyperion"

#run_benchmark "mariadb4" "${time}" "BenchmarkMariaDBInsert100kAndSort" "./db/mariadb"
#run_benchmark "postgres4" "${time}" "BenchmarkPostgresInsert100kAndSort" "./db/postgresql"
#run_benchmark "sqlite4" "${time}" "BenchmarkSQLiteInsert100kAndSort" "./db/sqlite"
#run_benchmark "hyperion4" "${time}" "BenchmarkHyperionInsert100kAndSort" "./db/hyperion"

#run_benchmark "mariadb5" "${time}" "BenchmarkMariaDBInsert100kAndQueryOlderThan15Min" "./db/mariadb"
#run_benchmark "postgres5" "${time}" "BenchmarkPostgresInsert100kAndQueryOlderThan15Min" "./db/postgresql"
#run_benchmark "sqlite5" "${time}" "BenchmarkSQLiteInsert100kAndQueryOlderThan15Min" "./db/sqlite"
#run_benchmark "hyperion5" "${time}" "BenchmarkHyperionInsert100kAndQueryOlderThan15Min" "./db/hyperion"


rm ./*.test > /dev/null 2>&1
