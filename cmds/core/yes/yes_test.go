// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}

func TestYes(t *testing.T) {
	count := uint64(5)
	for _, tt := range []struct {
		name     string
		in       []string
		expected string
	}{
		{
			name:     "noParameterCloseTest",
			in:       []string{},
			expected: "y\ny\ny\ny\ny\n",
		},
		{
			name:     "oneParameterCloseTest",
			in:       []string{"hi"},
			expected: "hi\nhi\nhi\nhi\nhi\n",
		},
		{
			name:     "fourParameterCloseTest",
			in:       []string{"hi", "how", "are", "you"},
			expected: "hi how are you\nhi how are you\nhi how are you\nhi how are you\nhi how are you\n",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := runYes(&buf, count, tt.in...); err != nil {
				t.Error(err)
			}
			if buf.String() != tt.expected {
				t.Errorf("%s does not match expected %s", buf.String(), tt.expected)
			}
		})
	}
}
