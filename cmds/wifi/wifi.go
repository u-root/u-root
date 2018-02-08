// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

const (
	cmd          = "wifi [options] essid [passphrase] [identity]"
	nopassphrase = `network={
		ssid="%s"
		proto=RSN
		key_mgmt=NONE
	}`
	eap = `network={
		ssid="%s"
		key_mgmt=WPA-EAP
		identity="%s"
		password="%s"
	}`
)

func init() {
	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = cmd
		defUsage()
	}
}

func main() {
	var (
		iface = flag.String("i", "wlan0", "interface to use")
		essid string
		conf  []byte
	)

	flag.Parse()
	a := flag.Args()

	switch {
	case len(a) == 3:
		essid = a[0]
		conf = []byte(fmt.Sprintf(eap, essid, a[2], a[1]))
	case len(a) == 2:
		essid = a[0]
		pass := a[1]
		o, err := exec.Command("wpa_passphrase", essid, pass).CombinedOutput()
		if err != nil {
			log.Fatalf("%v %v: %v", essid, pass, err)
		}
		conf = o
	case len(a) == 1:
		essid = a[0]
		conf = []byte(fmt.Sprintf(nopassphrase, essid))
	default:
		flag.Usage()
		os.Exit(1)
	}

	if err := ioutil.WriteFile("/tmp/wifi.conf", conf, 0444); err != nil {
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
