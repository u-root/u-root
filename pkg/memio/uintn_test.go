// Copyright 2012-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package memio

import (
	"fmt"
	"reflect"
	"testing"
)

// Slice of strings taken from io_test.go `var tests`
var wants = []string{
	"0x12",
	"0x1234",
	"0x12345678",
	"0x1234567890abcdef",
	"0x48656c6c6f",
}

func TestString(t *testing.T) {
	for i, tt := range tests {
		t.Run(fmt.Sprintf(tt.name), func(t *testing.T) {
			got := tt.writeData.String()
			want := wants[i]
			if !reflect.DeepEqual(want, got) {
				t.Errorf("Got: %v, want: %v", got, want)
			}
		},
		)
	}
}
