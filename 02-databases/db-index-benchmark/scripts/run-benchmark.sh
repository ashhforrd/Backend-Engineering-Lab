#!/usr/bin/env bash

set -euo pipefail

runs="${1:-5}"
unindexed_file="$(mktemp)"
indexed_file="$(mktemp)"

cleanup() {
  rm -f "$unindexed_file" "$indexed_file"
}

trap cleanup EXIT

for run in $(seq 1 "$runs"); do
  echo "Running benchmark $run/$runs..." >&2
  output="$(mvn -q compile exec:java)"

  printf '%s\n' "$output" | awk '
    /=== WITHOUT INDEX ===/ { section = "without"; next }
    /=== WITH INDEX ===/ { section = "with"; next }
    section == "without" && /Execution Time:/ { print $(NF - 1); exit }
  ' >> "$unindexed_file"

  printf '%s\n' "$output" | awk '
    /=== WITH INDEX ===/ { section = "with"; next }
    section == "with" && /Execution Time:/ { print $(NF - 1); exit }
  ' >> "$indexed_file"
done

summarize() {
  local label="$1"
  local file="$2"

  awk -v label="$label" '
    NR == 1 { min = max = $1 }
    {
      if ($1 < min) min = $1
      if ($1 > max) max = $1
      sum += $1
    }
    END {
      printf "%s: min=%.3f ms avg=%.3f ms max=%.3f ms runs=%d\n",
        label, min, sum / NR, max, NR
    }
  ' "$file"
}

summarize "without index" "$unindexed_file"
summarize "with index" "$indexed_file"
