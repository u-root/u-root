#!/usr/bin/env bash

# Copyright 2012-2024 the u-root Authors. All rights reserved
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
#     ./cmds/extra/tsort/bench.sh <commit-before> <commit-after>
#
# Note: This script assumes that hyperfine [1] is installed and that
# uutils/coreutils is installed on the PATH with a "uu" prefix such their tsort
# is named `uutsort`.
#
# [1] https://github.com/sharkdp/hyperfine
# [2] https://github.com/uutils/coreutils

set -euo pipefail

git switch --detach "$1"
go build -o ./tsort-before ./cmds/extra/tsort
echo "Running warmup benchmarks..."
go test -run=XXX -bench=Tsort -benchmem -count=10 ./cmds/extra/tsort/...
echo "Running real benchmarks..."
go test -run=XXX -bench=Tsort -benchmem -count=20 ./cmds/extra/tsort/... | tee tsort-bench-before.txt

git switch --detach "$2"
go build -o ./tsort-after ./cmds/extra/tsort
echo "Running warmup benchmarks..."
go test -run=XXX -bench=Tsort -benchmem -count=10 ./cmds/extra/tsort/...
echo "Running real benchmarks..."
go test -run=XXX -bench=Tsort -benchmem -count=20 ./cmds/extra/tsort/... | tee tsort-bench-after.txt

go run golang.org/x/perf/cmd/benchstat@latest tsort-bench-before.txt tsort-bench-after.txt | tee tsort-bench-comparison.txt

go run ./cmds/extra/tsort/generate_graph_fixtures.go
hyperfine --warmup 5 --shell=none --export-markdown=acyclic.md \
    "./tsort-before ./some-random-acyclic-graph.txt" \
    "./tsort-after ./some-random-acyclic-graph.txt" \
    "uutsort ./some-random-acyclic-graph.txt"
hyperfine --warmup 5 --shell=none --export-markdown=cyclic.md --ignore-failure \
    "./tsort-before ./some-random-cyclic-graph.txt" \
    "./tsort-after ./some-random-cyclic-graph.txt" \
    "uutsort ./some-random-cyclic-graph.txt"

rm ./some-random-acyclic-graph.txt
rm ./some-random-cyclic-graph.txt
rm ./tsort-before
rm ./tsort-after
