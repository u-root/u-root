// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base64

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type failer struct{}

// Write implements io.Writer, and always fails with os.ErrInvalid
func (failer) Write([]byte) (int, error) {
	return -1, os.ErrInvalid
}

func TestBase64(t *testing.T) {
	tests := []struct {
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
			out: []byte(`REVTQ1JJUFRJT04KICAgICAgIEJhc2U2NCBlbmNvZGUgb3IgZGVjb2RlIEZJTEUsIG9yIHN0YW5kYXJkIGlucHV0LCB0byBzdGFuZGFyZCBvdXRwdXQuCgogICAgICAgV2l0aCBubyBGSUxFLCBvciB3aGVuIEZJTEUgaXMgLSwgcmVhZCBzdGFuZGFyZCBpbnB1dC4KCiAgICAgICBNYW5kYXRvcnkgYXJndW1lbnRzIHRvIGxvbmcgb3B0aW9ucyBhcmUgbWFuZGF0b3J5IGZvciBzaG9ydCBvcHRpb25zIHRvby4K
`),
		},
	}
	d := t.TempDir()
	for _, tt := range tests {
		nin := filepath.Join(d, "in")
		if err := os.WriteFile(nin, tt.in, 0o666); err != nil {
			t.Fatalf(`WriteFile(%q, %v, 0666): %v != nil`, nin, tt.in, err)
		}
		nout := filepath.Join(d, "out")
		if err := os.WriteFile(nout, tt.out, 0o666); err != nil {
			t.Fatalf(`WriteFile(%q, %v, 0666): %v != nil`, nout, tt.out, err)
		}

		// Loop over encodes, then loop over decodes
		for _, n := range [][]string{{nin}, {}} {
			t.Run(fmt.Sprintf("run with file name %q", n), func(t *testing.T) {
				cmd := New()
				var o bytes.Buffer
				cmd.SetIO(bytes.NewBuffer(tt.in), &o, &bytes.Buffer{})
				if err := cmd.Run(n...); err != nil {
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
				cmd := New()
				var o bytes.Buffer
				cmd.SetIO(bytes.NewBuffer(tt.out), &o, &bytes.Buffer{})
				if err := cmd.Run(append([]string{"-d"}, n...)...); err != nil {
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
		cmd := New()
		if err := cmd.Run(n); err == nil {
			t.Errorf("run(%q): nil != an error", n)
		}
	})

	// Try with a bad length
	t.Run("bad data", func(t *testing.T) {
		cmd := New()
		bad := bytes.NewBuffer([]byte{'t'})
		var o bytes.Buffer
		cmd.SetIO(bad, &o, &bytes.Buffer{})
		if err := cmd.Run("-d"); err == nil {
			t.Errorf(`run("-d"): nil != an error`)
		}
	})
}

func TestBadWriter(t *testing.T) {
	cmd := New()
	cmd.SetIO(bytes.NewBufferString("hi there"), failer{}, &bytes.Buffer{})
	if err := cmd.Run(); !errors.Is(err, os.ErrInvalid) {
		t.Errorf(`Run(): got %v, want %v`, err, os.ErrInvalid)
	}
}

func TestBadUsage(t *testing.T) {
	tests := []struct {
		args []string
		err  error
	}{
		{args: []string{"x", "y"}, err: errBadUsage},
	}

	for _, tt := range tests {
		cmd := New()
		if err := cmd.Run(tt.args...); !errors.Is(err, tt.err) {
			t.Errorf(`Run(%q): got %v, want %v`, tt.args, err, tt.err)
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
			cmd := New()

			// encode first
			var encoded bytes.Buffer
			cmd.SetIO(strings.NewReader(tt.input), &encoded, &bytes.Buffer{})
			err := cmd.Run()
			if err != nil {
				t.Fatalf("encoding failed: %v", err)
			}

			// then decode
			cmd2 := New()
			var decoded bytes.Buffer
			cmd2.SetIO(bytes.NewReader(encoded.Bytes()), &decoded, &bytes.Buffer{})
			err = cmd2.Run("-d")
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

func TestRunContext(t *testing.T) {
	cmd := New()
	var o bytes.Buffer
	cmd.SetIO(strings.NewReader("hello"), &o, &bytes.Buffer{})

	ctx := context.Background()
	err := cmd.RunContext(ctx, []string{}...)
	if err != nil {
		t.Errorf("RunContext failed: %v", err)
	}

	if o.Len() == 0 {
		t.Error("RunContext produced no output")
	}
}
