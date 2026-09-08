// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uflag

import (
	"slices"
	"testing"
)

func TestArgvs(t *testing.T) {
	for _, tt := range []struct {
		argv []string
	}{
		{
			argv: []string{"--append=\"foobar\nfoobaz\"", "--haha"},
		},
		{
			argv: []string{"oh damn", "--append=\"foobar foobaz\"", "--haha"},
		},
		{
			argv: []string{},
		},
	} {
		got := FileToArgv(ArgvToFile(tt.argv))
		if !slices.Equal(got, tt.argv) {
			t.Errorf("FileToArgv(ArgvToFile(%#v)) = %#v, wanted original value back", tt.argv, got)
		}
	}
}
