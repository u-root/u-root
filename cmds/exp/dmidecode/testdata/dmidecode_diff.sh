#!/bin/bash
#
# Diff against dmidecode(8) output:
#
# DMIDECODE=./dmidecode DMIDECODE_ORIG=/usr/sbin/dmidecode testdata/dmidecode_diff.sh testdata/*.bin

set -e

if [ -z "$DMIDECODE" ]; then
  DMIDECODE="../demidecode"
fi

if [ -z "$DMIDECODE_ORIG" ]; then
  DMIDECODE_ORIG=$(which dmidecode)
fi

for bf in $*; do
  bn="${bf%.bin}"
  tf="$bn.txt"
  otf="$bn.orig.txt"
  df="$bn.orig.diff"
  $DMIDECODE --from-dump "$bf" > "$tf"
  $DMIDECODE_ORIG --from-dump "$bf" > "$otf"
  diff -u --label="$otf" "$otf" --label="$tf" "$tf" > "$df" || true
done
