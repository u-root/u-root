// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
//  created by Manoel Vilela (manoel_vilela@engineer.com)

package main

import (
	"io/ioutil"
	"os"
	"testing"
)

type testCase struct {
	name string
	f    mktempflags
	arg  string
	err  func(error) bool
}

func TestRemove(t *testing.T) {
	var (
		nilerr    = func(err error) bool { return err == nil }
		testCases = []testCase{
			{
				name: "no flags",
				err:  nilerr,
			},
			{
				name: "q",
				f:    mktempflags{q: true},
				err:  nilerr,
			},
			{
				name: "d",
				f:    mktempflags{d: true},
				err:  nilerr,
			},
			{
				name: "p",
				f:    mktempflags{prefix: "hithere"},
				err:  nilerr,
			},
		}
	)

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			// Always test with no dir specified and with one specified
			flags = tc.f
			n, err := mktemp()
			if err != nil {
				t.Error(err)
			}
			t.Logf("%v: name %v", tc, n)
			d, err := ioutil.TempDir(os.TempDir(), tc.f.dir)
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(d)
			tc.f.dir = d
			flags.dir = d
			n, err = mktemp()
			if err != nil {
				t.Error(err)
			}
			t.Logf("%v: name %v", tc, n)
		})
	}
}
