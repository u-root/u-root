// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

const cmd = "wifi [options] essid passphrase"

var (
	iface = flag.String("i", "wlan0", "interface to use")
	essid string
	pass  string
)

func init() {
	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = cmd
		defUsage()
	}
}

func main() {
	flag.Parse()
	a := flag.Args()
	if len(a) != 2 {
		flag.Usage()
	}
	essid, pass = a[0], a[1]

	conf, err := exec.Command("wpa_passphrase", essid, pass).CombinedOutput()
	if err != nil {
		log.Fatalf("%v %v: %v", essid, pass, err)
	}

	if err := ioutil.WriteFile("/tmp/wifi.conf", []byte(conf), 0444); err != nil {
		log.Fatalf("/tmp/wifi.conf: %v", err)
	}

	// There's no telling how long the supplicant will take, but on the other hand,
	// it's been almost instantaneous. But, further, it needs to keep running.
	go func() {
		if o, err := exec.Command("wpa_supplicant", "-i"+*iface, "-c/tmp/wifi.conf").CombinedOutput(); err != nil {
			log.Fatalf("wpa_supplicant: %v (%v)", o, err)
		}
	}()

	cmd := exec.Command("dhclient", "-ipv4=true", "-ipv6=false", "-verbose", *iface)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("%v: %v", cmd, err)
	}

}
