// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"os/exec"

	"github.com/u-root/u-root/pkg/sos"
)

var htmlRoot = flag.String("html", "/etc/sos/html", "Path for root of SOS html files")

func main() {
	if o, err := exec.Command("ip", "link", "set", "dev", "lo", "up").CombinedOutput(); err != nil {
		log.Fatalf("ip link set dev lo: %v (%v)", string(o), err)
	}
	sos.StartServer(sos.NewSosService())
}
