// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build plan9

package main

import "io/ioutil"

func Sethostname(n string) error {
	return ioutil.WriteFile("#c/sysname", []byte(n), 0644)
}
