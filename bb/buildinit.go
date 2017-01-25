// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"
	"os/exec"
)

func buildinit() {
	e := os.Environ()
	e = append(e, "CGO_ENABLED=0")
	e = append(e, "GO15VENDOREXPERIMENT=1")
	cmd := exec.Command("go", "build", "-o", "init", ".")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = config.Bbsh
	cmd.Env = e

	if err := cmd.Run(); err != nil {
		log.Fatalf("%v\n", err)
	}
}
