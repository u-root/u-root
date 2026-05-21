#!/usr/bin/env bash

# Copyright 2012-2026 the u-root Authors. All rights reserved
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# hyperfine_bench.sh runs Hyperfine [1] benchmarks against the version of tsort
# at the current Git commit, the system tsort and uutils/coreutils' [2] tsort.
# This gives a point of comparison of the performance of tsort versus other
# implementations.
#
# Note: This script assumes that a system tsort and Hyperfine [1] are installed
# on the PATH and that uutils/coreutils is installed on the PATH with a "uu"
# prefix such their tsort is named `uutsort`.
#
# [1] https://github.com/sharkdp/hyperfine
# [2] https://github.com/uutils/coreutils

set -euo pipefail

awk 'BEGIN {
    srand(1)
    nodeCount = 10000
    for (i = 0; i < 100 * nodeCount; i++) {
        x = int(rand() * (nodeCount + 1))
        y = int(rand() * (nodeCount + 1))
        if (x < y) print x, y; else print y, x
    }
}' > some-random-acyclic-graph.txt
trap 'rm some-random-acyclic-graph.txt' EXIT

awk 'BEGIN {
    srand(1)
    nodeCount = 200
    # Produces a cyclic graph with a fixed RNG seed and through
    # sheer probability.
    for (i = 0; i < 100 * nodeCount; i++) {
        x = int(rand() * (nodeCount + 1))
        y = int(rand() * (nodeCount + 1))
        print x, y
    }
}' > some-random-cyclic-graph.txt
trap 'rm some-random-acyclic-graph.txt some-random-cyclic-graph.txt' EXIT

go build -o ./tsort ./cmds/core/tsort
trap 'rm tsort some-random-acyclic-graph.txt some-random-cyclic-graph.txt' EXIT

printf "\nRunning Hyperfine benchmarks for ./tsort, tsort and uutsort on an acyclic graph...\n"
hyperfine --warmup 15 --runs 50 --shell=none --export-markdown=acyclic.md \
    "./tsort ./some-random-acyclic-graph.txt" \
    "tsort ./some-random-acyclic-graph.txt" \
    "uutsort ./some-random-acyclic-graph.txt"
printf "\nRunning Hyperfine benchmarks for ./tsort, tsort and uutsort on a cyclic graph...\n"
hyperfine --warmup 15 --runs 50 --shell=none --export-markdown=cyclic.md --ignore-failure \
    "./tsort ./some-random-cyclic-graph.txt" \
    "tsort ./some-random-cyclic-graph.txt" \
    "uutsort ./some-random-cyclic-graph.txt"
