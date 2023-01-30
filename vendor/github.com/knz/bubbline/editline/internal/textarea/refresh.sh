#!/usr/bin/env bash
#
# This file re-generates the sources in this directory from the
# upstream bubbles repository.
#
set -euxo pipefail

git clone --depth=1 https://github.com/charmbracelet/bubbles
cp bubbles/textarea/textarea.go textarea.go.orig
cp textarea.go.orig textarea.go
patch -p0 <textarea.go.diff
rm -rf bubbles
diff -uNr textarea.go.orig textarea.go >textarea.go.diff

