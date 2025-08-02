// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package termios

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestFS9(t *testing.T) {
	d := t.TempDir()
	for _, p := range []string{"fd", "dev"} {
		dd := filepath.Join(d, p)
		if err := os.MkdirAll(dd, 0o777); err != nil {
			t.Fatalf("mkdir %q: got %v, want nil", dd, err)
		}
	}
	wctl := filepath.Join(d, "dev", "wctl")
	window := []byte("          12          8          80         96  rawoff hold")
	devcons := filepath.Join(d, "dev", "cons")
	ctl := devcons + "ctl"
	for _, tt := range []struct {
		name     string
		make     func() error
		contents string
		f        func() error
		err      error
	}{
		{
			name: "failed open", make: func() error { return nil }, err: os.ErrNotExist, f: func() error {
				_, err := consctl(d, 0)
				return err
			},
		},
		{
			name: "bad sd0ctl", make: func() error { return os.WriteFile(filepath.Join(d, "fd/sd0ctl"), nil, 0o666) }, err: os.ErrNotExist, f: func() error {
				_, err := consctl(d, 0)
				return err
			},
		},
		{
			name: "good sd0ctl", make: func() error { return os.WriteFile(filepath.Join(d, "fd/0ctl"), []byte("a b "+devcons), 0o666) }, err: nil, f: func() error {
				n, err := consctl(d, 0)
				if err != nil {
					return err
				}
				if n != devcons {
					return fmt.Errorf("%q is not = %q:%w", n, devcons, os.ErrInvalid)
				}
				return nil
			},
		},
		{
			name: "no /dev/consctl", make: func() error { return nil }, err: os.ErrNotExist, f: func() error {
				_, err := consctlFile(d, 0)
				return err
			},
		},
		{
			name: "good /dev/consctl", make: func() error { return os.WriteFile(ctl, nil, 0o666) }, err: nil, f: func() error {
				_, err := consctlFile(d, 0)
				return err
			},
		},
		{
			name: "no wctl", make: func() error { return nil }, err: os.ErrNotExist, f: func() error {
				_, _, err := readWinSize(wctl)
				return err
			},
		},
		{
			name: "bad wctl", make: func() error { return os.WriteFile(wctl, []byte("12 33 11 "), 0o666) }, err: io.ErrUnexpectedEOF, f: func() error {
				_, _, err := readWinSize(wctl)
				return err
			},
		},
		{
			name: "good wctl", make: func() error { return os.WriteFile(wctl, window, 0o666) }, err: nil, f: func() error {
				r, c, err := readWinSize(wctl)
				if err != nil {
					return err
				}
				if c != 80-12 || r != 96-8 {
					return fmt.Errorf("dimension is (%d,%d) should be (%d,%d): %w", r, c, 96-8, 80-12, os.ErrInvalid)
				}
				return nil
			},
		},
	} {
		if err := tt.make(); err != nil {
			t.Errorf("%s: make(): got %v, want nil", tt.name, err)
			continue
		}
		if err := tt.f(); !errors.Is(err, tt.err) {
			t.Errorf("%s: got %v, want %v", tt.name, err, tt.err)
		}
	}
}
