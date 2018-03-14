// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
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

	// State of the service
	CurEssid        string
	ConnectingEssid string
	NearbyWifis     []wifi.WifiOption
	WifiWorker      wifi.Wifi
)

func init() {
	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = cmd
		defUsage()
	}
}

func getState() State {
	return State{NearbyWifis, ConnectingEssid, CurEssid}
}

func scanWifi() error {
	o, err := WifiWorker.ScanWifi()
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	NearbyWifis = o
	return nil
}

// To prevent race condition, there should be only one
// goroutine running connectWifiArbitrator at any one time
func connectWifiArbitrator() {
	// Accepted Routine = routine that is allowed to change the state
	var acceptedRoutineId []byte
	for req := range ConnectReqChan {
		switch {
		// Accepted routine returns its result
		case bytes.Equal(req.routineID, acceptedRoutineId):
			if req.success {
				CurEssid = req.essid
			} else {
				essid, err := WifiWorker.ScanCurrentWifi()
				if err != nil {
					log.Fatalf("error: %v", err)
				}
				CurEssid = essid
			}
			acceptedRoutineId = nil
			ConnectingEssid = ""
			req.c <- nil // Neccessary for testing
		// The requesting routine wins and can change the state
		case ConnectingEssid == "":
			acceptedRoutineId = req.routineID
			ConnectingEssid = req.essid
			req.c <- nil
		// The requesting routine loses
		default:
			req.c <- fmt.Errorf("Service is trying to connect to %s", ConnectingEssid)
		}
	}
}

// To prevent race condition, there should be only one
// goroutine running refreshNotifier at any one time
func refreshNotifier() {
	refreshing := false
	workDone := make(chan bool, 1)
	pool := make(chan RefreshReqChanMsg, DefaultBufferSize)

	// Pooler
	for {
		select {
		case r := <-RefreshReqChan:
			if !refreshing {
				refreshing = true
				// Notifier
				go func(p chan RefreshReqChanMsg) {
					err := scanWifi()
					workDone <- true
					for req := range p {
						req.c <- err
					}
				}(pool)
			}
			pool <- r
		case <-workDone:
			close(pool)
			refreshing = false
			pool = make(chan RefreshReqChanMsg, DefaultBufferSize)
		}
	}
}

func main() {
	flag.Parse()

	// Start a Server with Stub data
	// for manual end-to-end testing
	if *test {
		NearbyWifis = []wifi.WifiOption{
			{"Stub1", wifi.NoEnc},
			{"Stub2", wifi.WpaPsk},
			{"Stub3", wifi.WpaEap},
			{"Stub4", wifi.NotSupportedProto},
		}
		WifiWorker = wifi.StubWifiWorker{
			ScanInterfacesOut:  nil,
			ScanWifiOut:        NearbyWifis,
			ScanCurrentWifiOut: CurEssid,
		}
		go connectWifiArbitrator()
		startServer()
	}

	w, err := wifi.NewWorker(*iface)
	if err != nil {
		log.Fatal(err)
	}
	WifiWorker = w

	if *list {
		wifiOpts, err := WifiWorker.ScanWifi()
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
		interfaces, err := WifiWorker.ScanInterfaces()
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
		go scanWifi()
		go connectWifiArbitrator()
		startServer()
		return
	}

	if err := WifiWorker.Connect(a...); err != nil {
		log.Fatalf("error: %v", err)
	}
}
