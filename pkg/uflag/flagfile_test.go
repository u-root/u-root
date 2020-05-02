// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uflag

import (
	"reflect"
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
		// Accept nil for []string{} by checking len == 0.
		if !(len(tt.argv) == 0 && len(got) == 0) && !reflect.DeepEqual(got, tt.argv) {
			t.Errorf("FileToArgv(ArgvToFile(%#v)) = %#v, wanted original value back", tt.argv, got)
		}
	}
}
