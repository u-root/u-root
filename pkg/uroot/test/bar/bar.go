// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bar

import "fmt"

type Interface interface {
	UsedInterfaceMethod()
	UnusedInterfaceMethod()
}

type Bar struct{}

func (Bar) UsedInterfaceMethod() {
	fmt.Println("one")
}

// This method is unused but cannot be eliminated yet.
// https://github.com/golang/go/issues/38685
func (Bar) UnusedInterfaceMethod() {
	fmt.Println("two")
}

// This method is unused and should be eliminated.
func (Bar) UnusedNonInterfaceMethod() {
	fmt.Println("three")
}
