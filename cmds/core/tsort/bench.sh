#!/usr/bin/env bash

# Copyright 2012-2026 the u-root Authors. All rights reserved
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# bench.sh runs Go and hyperfine benchmarks against two given Git-commit-based
# versions of tsort. The hyperfine benchmarks also compare these versions of
# tsort with uutils/coreutils' [2] tsort.
#
# This script accepts two arguments, each containing a Git commit pointing at
# two different versions of tsort.
#
# Usage
#
#     ./cmds/core/tsort/bench.sh <commit-before> <commit-after>
#
# Note: This script assumes that hyperfine [1] is installed and that
# uutils/coreutils is installed on the PATH with a "uu" prefix such their tsort
# is named `uutsort`.
#
# [1] https://github.com/sharkdp/hyperfine
# [2] https://github.com/uutils/coreutils

set -euo pipefail

git switch --detach "$1"
go build -o ./tsort-before ./cmds/core/tsort
printf "\nRunning warmup Go benchmarks for ./tsort-before...\n"
go test -run=XXX -bench=Tsort -benchmem -count=10 ./cmds/core/tsort/...
printf "\nRunning real Go benchmarks for ./tsort-before...\n"
go test -run=XXX -bench=Tsort -benchmem -count=20 ./cmds/core/tsort/... | tee tsort-bench-before.txt
git switch -

git switch --detach "$2"
go build -o ./tsort-after ./cmds/core/tsort
printf "\nRunning warmup Go benchmarks for ./tsort-after...\n"
go test -run=XXX -bench=Tsort -benchmem -count=10 ./cmds/core/tsort/...
printf "\nRunning real Go benchmarks for ./tsort-after...\n"
go test -run=XXX -bench=Tsort -benchmem -count=20 ./cmds/core/tsort/... | tee tsort-bench-after.txt
git switch -

go run golang.org/x/perf/cmd/benchstat@latest tsort-bench-before.txt tsort-bench-after.txt | tee tsort-bench-comparison.txt

go run ./cmds/core/tsort/generate_graph_fixtures.go
printf "\nRunning Hyperfine benchmarks for ./tsort-before, ./tsort-after and uutsort on an acyclic graph...\n"
hyperfine --warmup 15 --runs 50 --shell=none --export-markdown=acyclic.md \
    "./tsort-before ./some-random-acyclic-graph.txt" \
    "./tsort-after ./some-random-acyclic-graph.txt" \
    "uutsort ./some-random-acyclic-graph.txt"
printf "\nRunning Hyperfine benchmarks for ./tsort-before, ./tsort-after and uutsort on a cyclic graph...\n"
hyperfine --warmup 15 --runs 50 --shell=none --export-markdown=cyclic.md --ignore-failure \
    "./tsort-before ./some-random-cyclic-graph.txt" \
    "./tsort-after ./some-random-cyclic-graph.txt" \
    "uutsort ./some-random-cyclic-graph.txt"

rm ./some-random-acyclic-graph.txt
rm ./some-random-cyclic-graph.txt
rm ./tsort-before
rm ./tsort-after
