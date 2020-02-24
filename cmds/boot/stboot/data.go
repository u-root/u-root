// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"path/filepath"
)

type data interface {
	get(filename string) ([]byte, error)
}

type initramfsData struct{}

func (*initramfsData) get(filename string) ([]byte, error) {
	f := filepath.Join("/root", filename)
	return ioutil.ReadFile(f)
}
