#!/bin/bash

if [ -z "$1" ]; then
	echo usage $0 path/to/grub_source_dir/
	exit 1
fi

grubsrcdir=$1

for f in $(grep -rnl grub-shell-tester $grubsrcdir/tests/*); do cp $f .; done
