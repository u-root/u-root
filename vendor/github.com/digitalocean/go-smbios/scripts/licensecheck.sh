#!/bin/bash

# Verify that the correct license block is present in all Go source
# files.
EXPECTED=$(cat ./scripts/license.txt)

# Scan each Go source file for license.
EXIT=0
GOFILES=$(find . -name "*.go")

for FILE in $GOFILES; do
	BLOCK=$(head -n 14 $FILE)

	if [ "$BLOCK" != "$EXPECTED" ]; then
		echo "file missing license: $FILE"
		EXIT=1
	fi
done

exit $EXIT
