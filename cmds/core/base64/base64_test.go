// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestBase64(t *testing.T) {
	var tests = []struct {
		in  []byte
		out []byte
	}{
		{
			in: []byte(`DESCRIPTION
       Base64 encode or decode FILE, or standard input, to standard output.

       With no FILE, or when FILE is -, read standard input.

       Mandatory arguments to long options are mandatory for short options too.
`),
			out: []byte(`REVTQ1JJUFRJT04KICAgICAgIEJhc2U2NCBlbmNvZGUgb3IgZGVjb2RlIEZJTEUsIG9yIHN0YW5kYXJkIGlucHV0LCB0byBzdGFuZGFyZCBvdXRwdXQuCgogICAgICAgV2l0aCBubyBGSUxFLCBvciB3aGVuIEZJTEUgaXMgLSwgcmVhZCBzdGFuZGFyZCBpbnB1dC4KCiAgICAgICBNYW5kYXRvcnkgYXJndW1lbnRzIHRvIGxvbmcgb3B0aW9ucyBhcmUgbWFuZGF0b3J5IGZvciBzaG9ydCBvcHRpb25zIHRvby4K`),
		},
	}
	d, err := ioutil.TempDir("", "base64")
	if err != nil {
		t.Fatalf(`TempDir("", "base64"): %v != nil`, err)
	}
	for _, tt := range tests {
		nin := filepath.Join(d, "in")
		if err := ioutil.WriteFile(nin, tt.in, 0666); err != nil {
			t.Fatalf(`WriteFile(%q, %v, 0666): %v != nil`, nin, tt.in, err)
		}
		nout := filepath.Join(d, "out")
		if err := ioutil.WriteFile(nout, tt.out, 0666); err != nil {
			t.Fatalf(`WriteFile(%q, %v, 0666): %v != nil`, nout, tt.out, err)
		}

		// Loop over encodes, then loop over decodes
		for _, n := range []string{"", "-", nin} {
			t.Run(fmt.Sprintf("run with file name %q", n), func(t *testing.T) {
				var o bytes.Buffer
				// n.b. the bytes.NewBuffer is ignored in all but one case ...
				if err := run(n, bytes.NewBuffer(tt.in), &o, false); err != nil {
					t.Errorf("Encode: got %v, want nil", err)
					return
				}
				if !bytes.Equal(o.Bytes(), tt.out) {
					t.Errorf("Encode: %q != %q", o.Bytes(), tt.out)
				}
			})
		}

		for _, n := range []string{"", "-", nout} {
			t.Run(fmt.Sprintf("run with file name %q", n), func(t *testing.T) {
				var o bytes.Buffer
				// n.b. the bytes.NewBuffer is ignored in all but one case ...
				if err := run(n, bytes.NewBuffer(tt.out), &o, true); err != nil {
					t.Errorf("Decode: got %v, want nil", err)
					return
				}
				if !bytes.Equal(o.Bytes(), tt.in) {
					t.Errorf("Decode: %q != %q", o.Bytes(), tt.out)
				}
			})
		}
	}
	// Try opening a file we know does not exist.
	n := filepath.Join(d, "nosuchfile")
	t.Run(fmt.Sprintf("bad file %q", n), func(t *testing.T) {
		// n.b. the bytes.NewBuffer is ignored in all but one case ...
		if err := run(n, nil, nil, false); err == nil {
			t.Errorf("run(%q, nil, nil, false): nil != an error", n)
		}
	})

	// Try with a bad length
	t.Run("bad data", func(t *testing.T) {
		var bad = bytes.NewBuffer([]byte{'t'})
		var o bytes.Buffer
		// n.b. the bytes.NewBuffer is ignored in all but one case ...
		if err := run("", bad, &o, true); err == nil {
			t.Errorf(`run("", zero-length buffer, zero-length-buffer, false): nil != an error`)
		}
	})

}
