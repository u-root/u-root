// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file exists to make gobusybox and u-root happy.
// It is replaced by the tinygobb tool.

//go:build tinygo

package main

func runone(c *Command) error {
	return nil
}
