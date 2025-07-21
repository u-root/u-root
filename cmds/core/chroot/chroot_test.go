// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"reflect"
	"testing"

	"github.com/hugelgupf/vmtest/guest"
)

func TestUserSpecSet(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "1000:1001",
			input: "1000:1001",
		},
		{
			name:  "test:1000",
			input: "test:1000",
			want:  fmt.Sprintf("strconv.ParseUint: parsing %q: invalid syntax", "test"),
		},
		{
			name:  "1000:test",
			input: "1000:test",
			want:  fmt.Sprintf("strconv.ParseUint: parsing %q: invalid syntax", "test"),
		},
		{
			name:  "1000:1001:",
			input: "1000:1001:",
			want:  fmt.Sprintf("expected user spec flag to be %q separated values received %s", ":", "1000:1001:"),
		},
		{
			name:  ":1000",
			input: ":1000",
			want:  fmt.Sprintf("strconv.ParseUint: parsing %q: invalid syntax", ""),
		},
		{
			name:  "1000:",
			input: "1000:",
			want:  fmt.Sprintf("expected user spec flag to be %q separated values received %s", ":", "1000:"),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var user userSpec
			if got := user.Set(tt.input); got != nil {
				if got.Error() != tt.want {
					t.Errorf("user.Set()= %q, want: %q", got.Error(), tt.want)
				}
			} else if user.uid != 1000 || user.gid != 1001 {
				t.Errorf("Expected uid 1000, gid 1001, got: uid %d, gid %d", user.uid, user.gid)
			}
		})
	}
}

func TestUserSpecGet(t *testing.T) {
	var user userSpec
	want := user
	got := user.Get()
	if got != want {
		t.Errorf("user.Get() = %v, want: %v", got, want)
	}
}

func TestGroupsSet(t *testing.T) {
	for _, tt := range []struct {
		name     string
		input    string
		expected []uint32
		want     string
	}{
		{
			name:     "1000",
			input:    "1000",
			expected: []uint32{1000},
		},
		{
			name:     "1000,1001",
			input:    "1000,1001",
			expected: []uint32{1000, 1001},
		},
		{
			name:     "1000,1001,",
			input:    "1000,1001,",
			expected: []uint32{},
			want:     fmt.Sprintf("strconv.ParseUint: parsing %q: invalid syntax", ""),
		},
		{
			name:     "test,1001",
			input:    "test,1001",
			expected: []uint32{},
			want:     fmt.Sprintf("strconv.ParseUint: parsing %q: invalid syntax", "test"),
		},
		{
			name:     ",1000",
			input:    ",1000",
			expected: []uint32{},
			want:     fmt.Sprintf("strconv.ParseUint: parsing %q: invalid syntax", ""),
		},
		{
			name:     "1000,",
			input:    "1000,",
			expected: []uint32{},
			want:     fmt.Sprintf("strconv.ParseUint: parsing %q: invalid syntax", ""),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var groups groupsSpec
			if err := groups.Set(tt.input); err != nil {
				if err.Error() != tt.want {
					t.Errorf("group.Set() = %q, want: %q", err.Error(), tt.want)
				}
			} else if len(groups.groups) != len(tt.expected) {
				t.Errorf("len(groups.groups) = %q, want: %q", len(groups.groups), len(tt.expected))
			} else {
				for index, group := range groups.groups {
					if tt.expected[index] != group {
						t.Errorf("expected[index] = %q, want: %q", tt.expected[index], group)
					}
				}
			}
		})
	}
}

func TestGroupsSpecGet(t *testing.T) {
	var groups groupsSpec
	want := groups
	got := groups.Get()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Get() = %q, want: %q", got, want)
	}
}

func TestGroupsString(t *testing.T) {
	var groups groupsSpec
	groups.groups = make([]uint32, 4)
	groups.groups[0] = 1000
	groups.groups[1] = 1001
	groups.groups[2] = 456666
	groups.groups[3] = 56758
	want := "1000,1001,456666,56758"
	got := groups.String()
	if got != want {
		t.Errorf("groups.String() = %q, want: %q", got, want)
	}
}

func TestChroot(t *testing.T) {
	guest.SkipIfNotInVM(t)

	for _, tt := range []struct {
		name          string
		args          []string
		skipchdirFlag bool
		want          string
	}{
		{
			name: "print defaults",
			args: []string{},
			want: defaults,
		},
		{
			name: "error in isRoot 1",
			args: []string{"/bin/sh"},
			want: "chdir /bin/sh: not a directory",
		},
		{
			name: "no error",
			args: []string{"/bin"},
			want: "fork/exec /bin/sh: no such file or directory",
		},
		{
			name:          "skipchdirFlag = true",
			args:          []string{"/bin/"},
			skipchdirFlag: true,
			want:          "the -s option is only permitted when newroot is the old / directory",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			skipchdirFlag = tt.skipchdirFlag
			buf := &bytes.Buffer{}
			flag.CommandLine.SetOutput(buf)
			if got := chroot(buf, tt.args...); got != nil {
				if got.Error() != tt.want {
					t.Errorf("chroot() = %q, want: %q", got.Error(), tt.want)
				}
			} else {
				if buf.String() != tt.want {
					t.Errorf("chroot() = %q, want: %q", buf.String(), tt.want)
				}
			}
		})
	}
}
