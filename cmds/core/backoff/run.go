// Copyright 2012-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/cenkalti/backoff/v4"
)

var (
	ErrNoCmd = fmt.Errorf("no command passed")
)

func runit(timeout string, c string, a ...string) error {
	if c == "" {
		return ErrNoCmd
	}
	b := backoff.NewExponentialBackOff()
	if len(timeout) != 0 {
		d, err := time.ParseDuration(timeout)
		if err != nil {
			return err
		}
		v("Set timeout to %v", d)
		b.MaxElapsedTime = d
	}
	f := func() error {
		cmd := exec.Command(c, a...)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		err := cmd.Run()
		v("%q %q:%v", c, a, err)
		return err
	}

	return backoff.Retry(f, b)
}
