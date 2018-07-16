// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uroot

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/u-root/u-root/pkg/golang"
)

func TestResolvePackagePaths(t *testing.T) {
	defaultEnv := golang.Default()
	gopath1, err := filepath.Abs("test/gopath1")
	if err != nil {
		t.Fatalf("failure to set up test: %v", err)
	}
	gopath2, err := filepath.Abs("test/gopath2")
	if err != nil {
		t.Fatalf("failure to set up test: %v", err)
	}
	gopath1Env := defaultEnv
	gopath1Env.GOPATH = gopath1
	gopath2Env := defaultEnv
	gopath2Env.GOPATH = gopath2
	everythingEnv := defaultEnv
	everythingEnv.GOPATH = gopath1 + ":" + gopath2
	foopath, err := filepath.Abs("test/gopath1/src/foo")
	if err != nil {
		t.Fatalf("failure to set up test: %v", err)
	}

	for _, tc := range []struct {
		env      golang.Environ
		in       []string
		expected []string
		wantErr  bool
	}{
		// Nonexistent Package
		{
			env:      defaultEnv,
			in:       []string{"fakepackagename"},
			expected: nil,
			wantErr:  true,
		},
		// Single go package import
		{
			env: defaultEnv,
			in:  []string{"github.com/u-root/u-root/cmds/ls"},
			// We expect the full URL format because that's the path in our default GOPATH
			expected: []string{"github.com/u-root/u-root/cmds/ls"},
			wantErr:  false,
		},
		// Single package directory relative to working dir
		{
			env:      defaultEnv,
			in:       []string{"test/gopath1/src/foo"},
			expected: []string{"github.com/u-root/u-root/pkg/uroot/test/gopath1/src/foo"},
			wantErr:  false,
		},
		// Single package directory with absolute path
		{
			env:      defaultEnv,
			in:       []string{foopath},
			expected: []string{"github.com/u-root/u-root/pkg/uroot/test/gopath1/src/foo"},
			wantErr:  false,
		},
		// Single package directory relative to GOPATH
		{
			env: gopath1Env,
			in:  []string{"foo"},
			expected: []string{
				"foo",
			},
			wantErr: false,
		},
		// Package directory glob
		{
			env: defaultEnv,
			in:  []string{"test/gopath2/src/mypkg*"},
			expected: []string{
				"github.com/u-root/u-root/pkg/uroot/test/gopath2/src/mypkga",
				"github.com/u-root/u-root/pkg/uroot/test/gopath2/src/mypkgb",
			},
			wantErr: false,
		},
		// GOPATH glob
		{
			env: gopath2Env,
			in:  []string{"mypkg*"},
			expected: []string{
				"mypkga",
				"mypkgb",
			},
			wantErr: false,
		},
		// Single ambiguous package - exists in both GOROOT and GOPATH
		{
			env: gopath1Env,
			in:  []string{"os"},
			expected: []string{
				"os",
			},
			wantErr: false,
		},
		// Packages from different gopaths
		{
			env: everythingEnv,
			in:  []string{"foo", "mypkga"},
			expected: []string{
				"foo",
				"mypkga",
			},
			wantErr: false,
		},
		// Same package specified twice
		{
			env: defaultEnv,
			in:  []string{"test/gopath2/src/mypkga", "test/gopath2/src/mypkga"},
			// TODO: This returns the package twice. Is this preferred?
			expected: []string{
				"github.com/u-root/u-root/pkg/uroot/test/gopath2/src/mypkga",
				"github.com/u-root/u-root/pkg/uroot/test/gopath2/src/mypkga",
			},
			wantErr: false,
		},
	} {
		t.Run(fmt.Sprintf("%q", tc.in), func(t *testing.T) {
			out, err := ResolvePackagePaths(tc.env, tc.in)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ResolvePackagePaths(%#v, %v) err != nil is %v, want %v\nerr is %v",
					tc.env, tc.in, err != nil, tc.wantErr, err)
			}
			if !reflect.DeepEqual(out, tc.expected) {
				t.Errorf("ResolvePackagePaths(%#v, %v) = %v; want %v",
					tc.env, tc.in, out, tc.expected)
			}
		})
	}
}
