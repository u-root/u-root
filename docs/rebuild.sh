#!/bin/sh

# extract '# Desription' section from root README.md
csplit -f prep_ ../README.md '/^# Description/' > /dev/null
csplit -f desc_ prep_01 '/^#/' '{1}' > /dev/null
# change headline to 'u-root' and fix relative links to point to GitHub
sed \
  -e 's/# Description/\n# u-root/' \
  -e 's#(\(cmds\|pkg\)#(https://github.com/u-root/u-root/tree/master/\1#g' \
  desc_01 > description.md
rm desc_* prep_*

# fetch pandoc-uikit template
_TEMPLATE=pandoc-uikit-master
[ -d "$_TEMPLATE" ] ||
  curl -L https://github.com/diversen/pandoc-uikit/archive/master.tar.gz | tar -xzf -

# cat it all and pipe into pandoc :)
cat header.md description.md index.md | pandoc --metadata title="u-root" --toc \
  -o index.html --template="$_TEMPLATE"/template.html -
