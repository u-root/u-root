// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"

	"github.com/u-root/u-root/pkg/uroot/test/bar"
)

var assignedTwice = "foo"

var assignWithoutType = "bla"

var aWT1, aWT2 = "foo", "bar"

var declOnly1, declOnly2 string

var (
	groupedDecl1               string
	groupedDecl2, groupedDecl3 string
)

var (
	groupedDeclOnlyIntf any
	nonConstantAssign   = fmt.Errorf("foo")
)

var nil1 any = nil

var (
	f1 func() string
	f2 = debug
	f3 = f1
)

func debug() string {
	return "hahaha"
}

type (
	someStuff  any
	someStruct struct{}
)

var (
	_ someStuff = &someStruct{}
	_           = "assign to no name"
)

func init() {
	groupedDecl1 = "foo"
}

func init() {
	groupedDecl2 = "urgh"
}

func init() {
	assignedTwice = "bar"
}

func main() {
	if err := verify(); err != nil {
		log.Fatalln(err)
	}
}

func verify() error {
	for _, tt := range []struct {
		name  string
		thing *string
		want  string
	}{
		{
			name:  "assignWithoutType",
			thing: &assignWithoutType,
			want:  "bla",
		},
		{
			name:  "aWT1",
			thing: &aWT1,
			want:  "foo",
		},
		{
			name:  "aWT2",
			thing: &aWT2,
			want:  "bar",
		},
		{
			name:  "declOnly1",
			thing: &declOnly1,
			want:  "",
		},
		{
			name:  "declOnly2",
			thing: &declOnly2,
			want:  "",
		},
		{
			name:  "groupedDecl1",
			thing: &groupedDecl1,
			want:  "foo",
		},
		{
			name:  "groupedDecl2",
			thing: &groupedDecl2,
			want:  "urgh",
		},
		{
			name:  "groupedDecl3",
			thing: &groupedDecl3,
			want:  "",
		},
	} {
		if got := *tt.thing; got != tt.want {
			return fmt.Errorf("%s is %s, want %s", tt.name, got, tt.want)
		}
	}

	if f1 != nil {
		return fmt.Errorf("f1 is non-nil, want nil")
	}
	if got := f2(); got != "hahaha" {
		return fmt.Errorf("f2 should return hahaha, but got %q", got)
	}
	if f3 != nil {
		return fmt.Errorf("f3 is non-nil, want nil")
	}

	if nil1 != any(nil) {
		return fmt.Errorf("nil1 is %v, want nil interface", nil1)
	}

	// Test unused method elimination.
	var b bar.Interface
	b = bar.Bar{}
	b.UsedInterfaceMethod()

	return nil
}
