// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dt

import (
	"errors"
	"fmt"
	"io"
)

// PrintDTS prints the FDT in the .dts format.
// TODO: not yet implemented
//
//nolint:staticcheck
func (fdt *FDT) PrintDTS(f io.Writer) error {
	return errors.New("not yet implemented")
}

// String implements String() for an FDT
func (fdt *FDT) String() string {
	return fmt.Sprintf("%#04x %s", fdt.Header, fdt.RootNode)
}

func (p *Property) String() string {
	var more string
	l := len(p.Value)
	if l > 64 {
		more = "..."
		l = 64
	}
	return fmt.Sprintf("%s[%#02x]%q{%#x}%s", p.Name, len(p.Value), p.Value[:l], p.Value[:l], more)
}

func (n *Node) String() string {
	var s string
	var indent string
	n.Walk(func(n *Node) error {
		i := indent
		indent += "\t"
		s += fmt.Sprintf("%s%s: [", indent, n.Name)
		for _, p := range n.Properties {
			s += fmt.Sprintf("%s, ", &p)
		}
		s += "]\n"
		indent = i
		return nil
	})
	return s
}
