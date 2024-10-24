// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !tinygo

package main

import (
	"fmt"
	"os"
)

func runone(c *Command) error {
	forkAttr(c)
	if err := c.Start(); err != nil {
		return fmt.Errorf("%w: Path %v", err, os.Getenv("PATH"))
	}
	if err := c.Wait(); err != nil {
		return fmt.Errorf("wait: %w", err)
	}
	return nil
}
