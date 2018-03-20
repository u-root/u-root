// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/u-root/u-root/pkg/wifi"
)

const (
	cmd = "wifi [options] essid [passphrase] [identity]"
)

var (
	// flags
	iface = flag.String("i", "wlan0", "interface to use")
	list  = flag.Bool("l", false, "list all nearby WiFi")
	show  = flag.Bool("s", false, "list interfaces allowed with WiFi extension")
	test  = flag.Bool("test", false, "set up a test server")

	Worker          wifi.Wifi
	Service         WifiService
	NearbyWifisStub = []wifi.WifiOption{
		{"Stub1", wifi.NoEnc},
		{"Stub2", wifi.WpaPsk},
		{"Stub3", wifi.WpaEap},
		{"Stub4", wifi.NotSupportedProto},
	}
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

	// Start a Server with Stub data
	// for manual end-to-end testing
	if *test {
		Worker = wifi.StubWifiWorker{
			ScanInterfacesOut:  nil,
			ScanWifiOut:        NearbyWifisStub,
			ScanCurrentWifiOut: "",
		}
		Service = NewWifiService(Worker)
		Service.Start()
		startServer()
		return
	}

	Worker, err := wifi.NewWorker(*iface)
	if err != nil {
		log.Fatal(err)
	}

	if *list {
		wifiOpts, err := Worker.ScanWifi()
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		for _, wifiOpt := range wifiOpts {
			switch wifiOpt.AuthSuite {
			case wifi.NoEnc:
				fmt.Printf("%s: No Passphrase\n", wifiOpt.Essid)
			case wifi.WpaPsk:
				fmt.Printf("%s: WPA-PSK (only passphrase)\n", wifiOpt.Essid)
			case wifi.WpaEap:
				fmt.Printf("%s: WPA-EAP (passphrase and identity)\n", wifiOpt.Essid)
			case wifi.NotSupportedProto:
				fmt.Printf("%s: Not a supported protocol\n", wifiOpt.Essid)
			}
		}
		return
	}

	if *show {
		interfaces, err := Worker.ScanInterfaces()
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		for _, i := range interfaces {
			fmt.Println(i)
		}
		return
	}

	a := flag.Args()
	if len(a) > 3 {
		flag.Usage()
		os.Exit(1)
	}

	// Experimental Part
	if len(a) == 0 {
		if o, err := exec.Command("ip", "link", "set", "dev", "lo", "up").CombinedOutput(); err != nil {
			log.Fatalf("ip link set dev lo: %v (%v)", string(o), err)
		}
		Service = NewWifiService(Worker)
		Service.Start()
		Service.RefreshReqChan <- make(chan error, 1)
		startServer()
		return
	}

	if err := Worker.Connect(a...); err != nil {
		log.Fatalf("error: %v", err)
	}
}
