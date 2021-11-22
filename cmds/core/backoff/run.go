// Copyright 2012-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/cenkalti/backoff/v4"
)

func runit(timeout string, c string, a ...string) error {
	if c == "" {
		return fmt.Errorf("no command passed")
	}
	ctx := context.Background()
	if len(timeout) != 0 {
		d, err := time.ParseDuration(timeout)
		if err != nil {
			return err
		}
		cx, cancel := context.WithTimeout(context.Background(), d)
		defer cancel()
		ctx = cx
	}
	b := backoff.WithContext(backoff.NewExponentialBackOff(), ctx)
	f := func() error {
		cmd := exec.Command(c, a...)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		return cmd.Run()
	}

	return backoff.Retry(f, b)
}
