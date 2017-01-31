#!/bin/bash

# Measure the performance of building all the Go commands under various GOGC
# values. The output is in a csv format. The first argument is passed to `time`
# as the format string. It can be one of:
#
# - %e: real time (s)
# - %U: user time (s)
# - %M: max RSS (KiB)
#
# Example: ./build_perf.sh %M > build_perf.csv

# TODO: refactor this into Go
FORMAT=${1-%U}
CMDS_PATH="$GOPATH/src/github.com/u-root/u-root/cmds"
PACKAGES=$(ls "$CMDS_PATH")

# header
printf 'GOGC,'
echo -n $PACKAGES | tr ' ' ','
echo

for GOGC in $(seq 50 50 2000); do
  printf "$GOGC"
  for PACKAGE in $PACKAGES; do
    cd "$CMDS_PATH/$PACKAGE"
    printf ,
    GOGC=$GOGC /usr/bin/time -f "$FORMAT" go build 2>&1 | tr -d '\n'
  done
  echo
done
