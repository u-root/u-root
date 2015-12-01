// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
//  created by Manoel Vilela (manoel_vilela@engineer.com)

package main

import (
	"fmt"
	"testing"
	"time"
)

// without any flag
func TestDateNoFlags(t *testing.T) {
	fmt.Println("== Printing date with default location (no flags)...")
	msg, err := date()
	fmt.Printf("%v\n", msg)
	if err != nil {
		t.Error(err)
	}
	_, err = time.Parse(time.UnixDate, msg)
	if err != nil {
		t.Error(err)
	}

}

// using u flag
func TestDateUniversal(t *testing.T) {
	flags.universal = true
	fmt.Println("== Printing date with UTC (using -u flag)...")
	msg, err := date()
	fmt.Printf("%v\n", msg)
	if err != nil {
		t.Error(err)
	}
	_, err = time.Parse(time.UnixDate, msg)
	if err != nil {
		t.Error(err)
	}
}
