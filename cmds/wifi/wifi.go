// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

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

var (
	// flags
	iface = flag.String("i", "wlan0", "interface to use")
	list  = flag.Bool("l", false, "list all nearby WiFi")
	show  = flag.Bool("s", false, "list interfaces allowed with WiFi extension")
	test  = flag.Bool("test", false, "set up a test server")

	// RegEx for parsing iwlist output
	CellRE       = regexp.MustCompile("(?m)^\\s*Cell")
	EssidRE      = regexp.MustCompile("(?m)^\\s*ESSID.*")
	EncKeyOptRE  = regexp.MustCompile("(?m)^\\s*Encryption key:(on|off)$")
	Wpa2RE       = regexp.MustCompile("(?m)^\\s*IE: IEEE 802.11i/WPA2 Version 1$")
	AuthSuitesRE = regexp.MustCompile("(?m)^\\s*Authentication Suites .*$")

	// RegEx for parsing iwconfig output
	IwconfigRE = regexp.MustCompile("(?m)^[a-zA-Z0-9]+\\s*IEEE 802.11.*$")

	// State of the service
	CurEssid        string
	ConnectingEssid string
	NearbyWifis     []WifiOption

	// State for the Test Server
	StubNearbyWifis = []WifiOption{
		{"Stub1", NoEnc},
		{"Stub2", WpaPsk},
		{"Stub3", WpaEap},
		{"Stub4", NotSupportedProto},
	}
)

func init() {
	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = cmd
		defUsage()
	}
}

func scanWifi() error {
	o, err := exec.Command("iwlist", *iface, "scanning").CombinedOutput()
	if err != nil {
		return fmt.Errorf("iwlist: %v (%v)", string(o), err)
	}
	NearbyWifis = parseIwlistOut(o)
	return nil
}

func getState() State {
	return State{NearbyWifis, ConnectingEssid, CurEssid}
}

/*
 * Assumptions:
 *	1) Cell, essid, and encryption key option are 1:1 match
 *	2) We only support IEEE 802.11i/WPA2 Version 1
 *	3) Each WiFi only support (1) authentication suites (based on observations)
 */

func parseIwlistOut(o []byte) []WifiOption {
	cells := CellRE.FindAllIndex(o, -1)
	essids := EssidRE.FindAll(o, -1)
	encKeyOpts := EncKeyOptRE.FindAll(o, -1)

	if cells == nil {
		return nil
	}

	var res []WifiOption
	knownEssids := make(map[string]bool)

	// Assemble all the WiFi options
	for i := 0; i < len(cells); i++ {
		essid := strings.Trim(strings.Split(string(essids[i]), ":")[1], "\"\n")
		if knownEssids[essid] {
			continue
		}
		knownEssids[essid] = true
		encKeyOpt := strings.Trim(strings.Split(string(encKeyOpts[i]), ":")[1], "\n")
		if encKeyOpt == "off" {
			res = append(res, WifiOption{essid, NoEnc})
			continue
		}
		// Find the proper Authentication Suites
		start, end := cells[i][0], len(o)
		if i != len(cells)-1 {
			end = cells[i+1][0]
		}
		// Narrow down the scope when looking for WPA Tag
		wpa2SearchArea := o[start:end]
		l := Wpa2RE.FindIndex(wpa2SearchArea)
		if l == nil {
			res = append(res, WifiOption{essid, NotSupportedProto})
			continue
		}
		// Narrow down the scope when looking for Authorization Suites
		authSearchArea := wpa2SearchArea[l[0]:]
		authSuites := strings.Trim(strings.Split(string(AuthSuitesRE.Find(authSearchArea)), ":")[1], "\n ")
		switch authSuites {
		case "PSK":
			res = append(res, WifiOption{essid, WpaPsk})
		case "802.1x":
			res = append(res, WifiOption{essid, WpaEap})
		default:
			res = append(res, WifiOption{essid, NotSupportedProto})
		}
	}
	return res
}

func scanInterfaces() ([]string, error) {
	o, err := exec.Command("iwconfig").CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("iwconfig: %v (%v)", string(o), err)
	}
	return parseIwconfig(o), nil
}

func parseIwconfig(o []byte) (res []string) {
	interfaces := IwconfigRE.FindAll(o, -1)
	for _, i := range interfaces {
		res = append(res, strings.Split(string(i), " ")[0])
	}
	return
}

func generateConfig(a ...string) (conf []byte, err error) {
	// format of a: [essid, pass, id]
	switch {
	case len(a) == 3:
		conf = []byte(fmt.Sprintf(eap, a[0], a[2], a[1]))
	case len(a) == 2:
		conf, err = passphrase.Run(a[0], a[1])
		if err != nil {
			return nil, fmt.Errorf("essid: %v, pass: %v : %v", a[0], a[1], err)
		}
	case len(a) == 1:
		conf = []byte(fmt.Sprintf(nopassphrase, a[0]))
	default:
		return nil, fmt.Errorf("generateConfig needs 1, 2, or 3 args")
	}
	return
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
				// Depending on where the failure happens,
				// we might or might not be connected to a WiFi
				o, _ := exec.Command("iwgetid", "-r").CombinedOutput()
				CurEssid = strings.Trim(string(o), " \n")
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

func connectWifi(a ...string) error {
	// format of a: [essid, pass, id]
	conf, err := generateConfig(a...)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile("/tmp/wifi.conf", conf, 0444); err != nil {
		return fmt.Errorf("/tmp/wifi.conf: %v", err)
	}

	c := make(chan error, 2)

	// There's no telling how long the supplicant will take, but on the other hand,
	// it's been almost instantaneous. But, further, it needs to keep running.
	go func() {
		cmd := exec.Command("wpa_supplicant", "-i"+*iface, "-c/tmp/wifi.conf")
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr //For an easier time debugging
		if err := cmd.Run(); err != nil {
			c <- fmt.Errorf("wpa_supplicant error: %v", err)
		} else {
			c <- nil
		}
	}()

	go func() {
		cmd := exec.Command("dhclient", "-ipv4=true", "-ipv6=false", "-verbose", *iface)
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr //For an easier time debugging
		if err := cmd.Run(); err != nil {
			c <- fmt.Errorf("dhclient error: %v", err)
		} else {
			c <- nil
		}
	}()

	if errWpaSupplicant, errDhClient := <-c, <-c; errWpaSupplicant != nil || errDhClient != nil {
		return fmt.Errorf("%v \n %v", errWpaSupplicant, errDhClient)
	}
	return nil
}

func main() {
	flag.Parse()

	if *list {
		if o, err := exec.Command("ip", "link", "set", "dev", *iface).CombinedOutput(); err != nil {
			log.Fatalf("ip link set dev %v: %v (%v)", *iface, string(o), err)
		}
		if err := scanWifi(); err != nil {
			log.Fatalf("error: %v", err)
		}
		for _, wifiOpt := range NearbyWifis {
			switch wifiOpt.AuthSuite {
			case NoEnc:
				fmt.Printf("%s: No Passphrase\n", wifiOpt.Essid)
			case WpaPsk:
				fmt.Printf("%s: WPA-PSK (only passphrase)\n", wifiOpt.Essid)
			case WpaEap:
				fmt.Printf("%s: WPA-EAP (passphrase and identity)\n", wifiOpt.Essid)
			case NotSupportedProto:
				fmt.Printf("%s: Not a supported protocol\n", wifiOpt.Essid)
			}
		}
		return
	}

	if *show {
		interfaces, err := scanInterfaces()
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		for _, i := range interfaces {
			fmt.Println(i)
		}
		return
	}

	if *test {
		NearbyWifis = StubNearbyWifis
		go connectWifiArbitrator()
		startServer()
	}

	a := flag.Args()

	// Experimental Part
	if len(a) == 0 {
		if o, err := exec.Command("ip", "link", "set", "dev", "lo").CombinedOutput(); err != nil {
			log.Fatalf("ip link set dev lo: %v (%v)", string(o), err)
		}
		if o, err := exec.Command("ip", "link", "set", "dev", *iface).CombinedOutput(); err != nil {
			log.Fatalf("ip link set dev %v: %v (%v)", *iface, string(o), err)
		}
		go scanWifi()
		go connectWifiArbitrator()
		startServer()
		return
	}

	if err := connectWifi(a...); err != nil {
		flag.Usage()
		log.Fatalf("error: %v", err)
	}
}
