#!/usr/bin/env python3
#
# Copyright 2016-2019 the u-root Authors. All rights reserved
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import os
import re
import sys

type_name = sys.argv[1]

enum_values = []

for line in sys.stdin:
    parts = line.strip().split(None, 1)
    if len(parts) != 2 or not re.match(r"[0-9A-Fa-f]{2}h", parts[0]):
        continue   # some junk, skip
    v = int(parts[0][:-1], 16)
    descr = parts[1]
    ident_words = []
    for w in descr.split():
        w = re.sub(r"[^a-zA-Z0-9_]", "", w)
        if w[0].islower():
            w = w.capitalize()
        ident_words.append(w)
    ident = type_name + "".join(ident_words)
    enum_values.append((v, ident, descr))

print("""\
// %s values are defined in DSP0134 x.x.x
const (""" % type_name)
i = 0
for v, ident, descr in enum_values:
    i += 1
    if i == 1:
        print("\t%s %s = 0x%02x // %s" % (ident, type_name, v, descr))
    else:
        print("\t%s = 0x%02x // %s" % (ident, v, descr))
print("""\
)

func (v %s) String() string {
\tswitch v {""" % type_name)

for v, ident, descr in enum_values:
    print('\tcase %s: return "%s"' % (ident, descr))

print("""\
\t}
\treturn fmt.Sprintf("%d", v)
}""")
