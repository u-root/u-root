// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vmtest

import (
	"reflect"
	"testing"
)

func TestFields(t *testing.T) {
	for _, test := range []struct {
		data string
		want []string
	}{
		{
			data: "foo bar baz",
			want: []string{"foo", "bar", "baz"},
		},
		{
			data: "foo\tbar\nbaz",
			want: []string{"foo", "bar", "baz"},
		},
		{
			data: `"foo bar" baz`,
			want: []string{`"foo bar"`, "baz"},
		},
		{
			data: `"foo\tbar" baz`,
			want: []string{`"foo\tbar"`, "baz"},
		},
		{
			data: `foo "bar baz"`,
			want: []string{"foo", `"bar baz"`},
		},
		{
			data: `foo "bar' baz`,
			want: []string{"foo", `"bar' baz`},
		},
		{
			data: `foo bar "baz`,
			want: []string{"foo", "bar", `"baz`},
		},
	} {
		t.Run(test.data, func(t *testing.T) {
			got := fields(test.data)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("fields() got %v, want %v", got, test.want)
			}
		})
	}
}
