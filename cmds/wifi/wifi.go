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

	"github.com/u-root/u-root/pkg/wpa/passphrase"
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

func generateConfig(essid string, id string, pass string) (conf []byte, err error) {
	switch {
	case essid != "" && id != "" && pass != "":
		conf = []byte(fmt.Sprintf(eap, essid, id, pass))
	case essid != "" && id == "" && pass != "":
		conf, err = passphrase.Run(essid, pass)
		if err != nil {
			return nil, fmt.Errorf("essid: %v, pass: %v : %v", essid, pass, err)
		}
	case essid != "" && id == "" && pass == "":
		conf = []byte(fmt.Sprintf(nopassphrase, essid))
	default:
		return nil, fmt.Errorf("Invalid Argument: essid: %v, id: %v, pass: %v", essid, id, pass)
	}
	return
}

func main() {
	go startServer()

	// Service
	var (
		iface = flag.String("i", "wlan0", "interface to use")
		conf  []byte
		err   error
	)

	flag.Parse()
	a := flag.Args()

	switch {
	case len(a) == 3:
		conf, err = generateConfig(a[0], a[2], a[1])
	case len(a) == 2:
		conf, err = generateConfig(a[0], "", a[1])
	case len(a) == 1:
		conf, err = generateConfig(a[0], "", "")
	case len(a) == 0:
		// Experimental Part
		// if len(a) = 0, can use the web interface to get user's input
		msg := <-ServerToServiceChan
		conf, err = generateConfig(msg.essid, msg.id, msg.pass)
		stubMsg := ServiceToServerMessage{
			essid: msg.essid,
		}
		ServiceToServerChan <- stubMsg
		_ = <-ServerToServiceChan // (Experimental) So we can see the result of the page load
	default:
		flag.Usage()
		os.Exit(1)
	}

	if err != nil {
		log.Fatalf("error: %v", err)
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
