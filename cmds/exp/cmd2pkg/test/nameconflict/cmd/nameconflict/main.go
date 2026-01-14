// Copyright 2015-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	// defaultlog declares itself as `package deflog`.

	// anotherlog makes sure that a package can be imported twice with a different name.
	anotherlog "github.com/u-root/gobusybox/test/nameconflict/pkg/defaultlog"

	// Create a conflict with the self-registering package import.
	bbmain "flag"
)

var something = bbmain.String("someflag", "", "")

// log will conflict with `import "log"` in order to read
//
// var log *log.Logger
var log = deflog.Default()
var log2 = anotherlog.Default()

// should conflict with init being rewritten.
func busyboxInit0() {
	fmt.Println("busyboxInit0")
}

// should be rewritten as busyboxInit2 because of name conflict with busyboxInit0 and busyboxInit1.
func init() {
	fmt.Println("init")
}

func main() {
	busyboxInit0()
	busyboxInit1()
	registeredInit()
	registeredMain()
}
