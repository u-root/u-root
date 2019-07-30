#!/usr/bin/env python3
#
# Copyright 2016-2019 the u-root Authors. All rights reserved
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import os
import re
import sys

type_name = sys.argv[1]

length_values = {"BYTE": "uint8", "WORD": "uint16", "DWORD": "uint32", "QWORD": "uint64"}

fields = []
enums = []

def ProcessRow(row):
    if not row:
        return
    parts = row.strip().split(None, 1)
    off = int(parts[0][:-1], 16)
    if off < 4:  # Header
        return
    rem = parts[1]
    ident, go_type, comment = "", "", ""
    print(rem)
    for w in rem.split():
        if not go_type:
            len_col = False
            for lv, ft in length_values.items():
                if lv in w:
                    # This is the length column which marks the end of identifier
                    # and is our first indication of type
                    # (may be corrected later for enums and strings).
                    go_type = ft
                    len_col = True
            if len_col:
                continue
            # Words before the type comprise the identifier.
            if w.endswith("+"):
                continue  # Version, skip it.
            # Sanitize the identifier.
            w = re.sub(r"[^a-zA-Z0-9_]", "", w)
            if w and w[0].islower():
                w = w.capitalize()
            ident += w
        elif not comment:    # Value type column.
            if w == "STRING":
                go_type = "string"
            elif w == "ENUM":
                enums.append((ident, go_type))
                go_type = ident
            elif w in ("Bit", "Field"):
                if go_type != ident:
                    enums.append((ident, go_type))
                go_type = ident
            elif w == "Varies":
                pass
            else:
                comment = w
        else:
            if len(comment) < 50:
                comment += " " + w
            elif not comment.endswith(" ..."):
                comment += " ..."
    fields.append((off, ident, go_type, comment))


row = ""
for line in sys.stdin:
    if re.search(r"^[0-9A-Fa-f]{2}h", line):
        ProcessRow(row)
        row = line.strip()
    else:
        row += " " + line.strip()

ProcessRow(row)


print("""\
// %s is defined in DSP0134 x.x.
type %s struct {""" % (type_name,type_name))
i = 0
for off, ident, go_type, comment in fields:
    i += 1
    if i == 1:
        print("\tTable")
    print("\t%s\t%s\t// %02Xh" % (ident, go_type, off))
print("}")

for n, t in enums:
    print("// %s is defined in DSP0134 x.x.x." % n)
    print("type %s %s" % (n, t))
