// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type failer struct {
}

// Write implements io.Writer, and always fails with os.ErrInvalid
func (failer) Write([]byte) (int, error) {
	return -1, os.ErrInvalid
}

func TestBase64(t *testing.T) {
	var tests = []struct {
		in   []byte
		out  []byte
		args []string
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
	d := t.TempDir()
	for _, tt := range tests {
		nin := filepath.Join(d, "in")
		if err := os.WriteFile(nin, tt.in, 0666); err != nil {
			t.Fatalf(`WriteFile(%q, %v, 0666): %v != nil`, nin, tt.in, err)
		}
		nout := filepath.Join(d, "out")
		if err := os.WriteFile(nout, tt.out, 0666); err != nil {
			t.Fatalf(`WriteFile(%q, %v, 0666): %v != nil`, nout, tt.out, err)
		}

		// Loop over encodes, then loop over decodes
		for _, n := range [][]string{{nin}, {}} {
			t.Run(fmt.Sprintf("run with file name %q", n), func(t *testing.T) {
				var o bytes.Buffer
				// n.b. the bytes.NewBuffer is ignored in all but one case ...
				if err := run(bytes.NewBuffer(tt.in), &o, false, n...); err != nil {
					t.Errorf("Encode: got %v, want nil", err)
					return
				}
				if !bytes.Equal(o.Bytes(), tt.out) {
					t.Errorf("Encode: %q != %q", o.Bytes(), tt.out)
				}
			})
		}

		for _, n := range [][]string{{nout}, {}} {
			t.Run(fmt.Sprintf("run with file name %q", n), func(t *testing.T) {
				var o bytes.Buffer
				// n.b. the bytes.NewBuffer is ignored in all but one case ...
				if err := run(bytes.NewBuffer(tt.out), &o, true, n...); err != nil {
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
		if err := run(nil, nil, false, n); err == nil {
			t.Errorf("run(%q, nil, nil, false): nil != an error", n)
		}
	})

	// Try with a bad length
	t.Run("bad data", func(t *testing.T) {
		var bad = bytes.NewBuffer([]byte{'t'})
		var o bytes.Buffer
		// n.b. the bytes.NewBuffer is ignored in all but one case ...
		if err := run(bad, &o, true); err == nil {
			t.Errorf(`run("", zero-length buffer, zero-length-buffer, false): nil != an error`)
		}
	})

}

func TestBadWriter(t *testing.T) {
	if err := run(bytes.NewBufferString("hi there"), failer{}, false); !errors.Is(err, os.ErrInvalid) {
		t.Errorf(`bytes.NewBufferString("hi there"), failer{}, false): got %v, want %v`, err, os.ErrInvalid)
	}
}
func TestBadUsage(t *testing.T) {
	var tests = []struct {
		args []string
		err  error
	}{
		{args: []string{"x", "y"}, err: errBadUsage},
	}

	for _, tt := range tests {
		if err := run(nil, nil, false, tt.args...); !errors.Is(err, tt.err) {
			t.Errorf(`run(nil, nil, false, %q): got %v, want %v`, tt.args, err, tt.err)
		}
	}
}

func TestDo(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{

		{
			name:  "single character",
			input: "a",
		},
		{
			name:  "four bytes",
			input: "abcd",
		},
		{
			name:  "five bytes",
			input: "abcde",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// encode first
			var encoded bytes.Buffer
			err := do(strings.NewReader(tt.input), &encoded, false)
			if err != nil {
				t.Fatalf("encoding failed: %v", err)
			}

			// then decode
			var decoded bytes.Buffer
			err = do(bytes.NewReader(encoded.Bytes()), &decoded, true)
			if err != nil {
				t.Fatalf("decoding failed: %v", err)
			}

			d := decoded.String()
			if d != tt.input {
				t.Errorf("encode/decode failed:\noriginal: %q\ndecoded:  %q", tt.input, d)
			}
		})
	}
}
