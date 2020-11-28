// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
)

// printfWrapper is a function suitable to be used for checks and remediations.
// It just wraps `fmt.Printf` and returns a `nil` error, since the CheckRunner
// and RemediationRunner type requires returning an error.
func printfWrapper(fmtstr string, args ...interface{}) error {
	fmt.Printf(fmtstr, args...)
	return nil
}
