// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dt

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

// sysfs binary attributes may return a single page even for a larger Read.
type chunkedReadSeeker struct {
	io.ReadSeeker
	max int
}

func (r chunkedReadSeeker) Read(p []byte) (int, error) {
	if len(p) > r.max {
		p = p[:r.max]
	}
	return r.ReadSeeker.Read(p)
}

func TestReadFDTShortStringsRead(t *testing.T) {
	root := NewNode("")
	for i := range 400 {
		root.Properties = append(root.Properties,
			PropertyString(fmt.Sprintf("property-%04d-unique", i), "value"))
	}
	original := &FDT{
		Header:   Header{Magic: Magic, Version: 17, LastCompVersion: 16},
		RootNode: root,
	}
	var raw bytes.Buffer
	if _, err := original.Write(&raw); err != nil {
		t.Fatal(err)
	}
	if original.Header.SizeDtStrings <= 4096 {
		t.Fatal("test FDT strings block is too short")
	}
	parsed, err := ReadFDT(chunkedReadSeeker{bytes.NewReader(raw.Bytes()), 4096})
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed.RootNode.Properties) != len(root.Properties) {
		t.Fatalf("got %d properties, want %d", len(parsed.RootNode.Properties), len(root.Properties))
	}
	for i, prop := range parsed.RootNode.Properties {
		if prop.Name != root.Properties[i].Name {
			t.Fatalf("property %d name = %q, want %q", i, prop.Name, root.Properties[i].Name)
		}
	}
}
