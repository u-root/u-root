#!/bin/bash
# this wrapper calls $1 instead of the buildin for echo while running $2
echocmd=$1
enable -n echo
mkdir testdata/bin
echo "#!/bin/bash" > testdata/bin/echo
echo "$echocmd \"\$@\"" >> testdata/bin/echo
chmod +x testdata/bin/echo
export PATH=$(realpath testdata/bin):$PATH

source $2
