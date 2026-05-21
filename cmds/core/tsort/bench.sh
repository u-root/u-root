#!/usr/bin/env bash

# Copyright 2012-2026 the u-root Authors. All rights reserved
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# bench.sh runs Go benchmarks against two given Git-commit-based versions of
# tsort.
#
# This script accepts two arguments, each containing a Git commit pointing at
# two different versions of tsort.
#
# Usage
#
#     ./cmds/core/tsort/bench.sh <commit-before> <commit-after>

set -euo pipefail

current_branch_or_commit="$(git rev-parse --abbrev-ref HEAD)"
if [ "$current_branch_or_commit" = "HEAD" ]; then
    current_branch_or_commit="$(git rev-parse HEAD)"
fi
trap 'git checkout "$current_branch_or_commit"' EXIT

git checkout "$1"
go build -o ./tsort-before ./cmds/core/tsort
trap 'rm ./tsort-before; git checkout "$current_branch_or_commit"' EXIT
printf "\nRunning warmup Go benchmarks for ./tsort-before...\n"
go test -run=XXX -bench=Tsort -benchmem -count=2 ./cmds/core/tsort/...
printf "\nRunning real Go benchmarks for ./tsort-before...\n"
go test -run=XXX -bench=Tsort -benchmem -count=10 ./cmds/core/tsort/... | tee tsort-bench-before.txt

git checkout "$2"
go build -o ./tsort-after ./cmds/core/tsort
trap 'rm ./tsort-after; rm ./tsort-before; git checkout "$current_branch_or_commit"' EXIT
printf "\nRunning warmup Go benchmarks for ./tsort-after...\n"
go test -run=XXX -bench=Tsort -benchmem -count=2 ./cmds/core/tsort/...
printf "\nRunning real Go benchmarks for ./tsort-after...\n"
go test -run=XXX -bench=Tsort -benchmem -count=10 ./cmds/core/tsort/... | tee tsort-bench-after.txt

go run golang.org/x/perf/cmd/benchstat@latest tsort-bench-before.txt tsort-bench-after.txt | tee tsort-bench-comparison.txt
