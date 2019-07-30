#!/usr/bin/env python3
#
# Copyright 2016-2019 the u-root Authors. All rights reserved
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import os
import re
import sys

type_name = sys.argv[1]

bits, extra = [], []

for line in sys.stdin:
    parts = line.strip().split(None, 2)
    if len(parts) != 3 or parts[0] != "Bit":
        extra.append("// %s" % line.strip())
        continue   # some junk, skip
    v = int(parts[1], 10)
    descr = parts[2]
    ident_words = []
    for w in descr.split():
        w = re.sub(r"[^a-zA-Z0-9_]", "", w)
        if w and w[0].islower():
            w = w.capitalize()
        ident_words.append(w)
    ident = type_name + "".join(ident_words)
    bits.append((v, ident, descr))

print("""\
// %s fields are defined in DSP0134 x.x.x
const (""" % )
i = 0
for v, ident, descr in bits:
    i += 1
    if i == 1:
        print("\t%s %s = (1 << %d) // %s" % (ident, type_name, v, descr))
    else:
        print("\t%s = (1 << %d) // %s" % (ident, v, descr))
print("""\
)

func (v %s) String() string {
\tvar lines []string""" % type_name)

for v, ident, descr in bits:
    print('\tif (v & %s != 0) { lines = append(lines, "\\t\\t%s") }' % (ident, descr))

print("""\
\treturn strings.Join(lines, "\\n")
}""")

print("\n".join(extra))
