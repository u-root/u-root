// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wifi

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/u-root/u-root/pkg/wpa/passphrase"
)

const (
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
	// RegEx for parsing iwlist output
	cellRE       = regexp.MustCompile("(?m)^\\s*Cell")
	essidRE      = regexp.MustCompile("(?m)^\\s*ESSID.*")
	encKeyOptRE  = regexp.MustCompile("(?m)^\\s*Encryption key:(on|off)$")
	wpa2RE       = regexp.MustCompile("(?m)^\\s*IE: IEEE 802.11i/WPA2 Version 1$")
	authSuitesRE = regexp.MustCompile("(?m)^\\s*Authentication Suites .*$")
)

type SecProto int

const (
	NoEnc SecProto = iota
	WpaPsk
	WpaEap
	NotSupportedProto
)

// IWLWorker implements the WiFi interface using the Intel Wireless LAN commands
type IWLWorker struct {
	Interface string
}

func NewIWLWorker(i string) (WiFi, error) {
	if o, err := exec.Command("ip", "link", "set", "dev", i, "up").CombinedOutput(); err != nil {
		return &IWLWorker{""}, fmt.Errorf("ip link set dev %v up: %v (%v)", i, string(o), err)
	}
	return &IWLWorker{i}, nil
}

func (w *IWLWorker) Scan() ([]Option, error) {
	o, err := exec.Command("iwlist", w.Interface, "scanning").CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("iwlist: %v (%v)", string(o), err)
	}
	return parseIwlistOut(o), nil
}

/*
 * Assumptions:
 *	1) Cell, essid, and encryption key option are 1:1 match
 *	2) We only support IEEE 802.11i/WPA2 Version 1
 *	3) Each Wifi only support (1) authentication suites (based on observations)
 */

func parseIwlistOut(o []byte) []Option {
	cells := cellRE.FindAllIndex(o, -1)
	essids := essidRE.FindAll(o, -1)
	encKeyOpts := encKeyOptRE.FindAll(o, -1)

	if cells == nil {
		return nil
	}

	var res []Option
	knownEssids := make(map[string]bool)

	// Assemble all the Wifi options
	for i := 0; i < len(cells); i++ {
		essid := strings.Trim(strings.Split(string(essids[i]), ":")[1], "\"\n")
		if knownEssids[essid] {
			continue
		}
		knownEssids[essid] = true
		encKeyOpt := strings.Trim(strings.Split(string(encKeyOpts[i]), ":")[1], "\n")
		if encKeyOpt == "off" {
			res = append(res, Option{essid, NoEnc})
			continue
		}
		// Find the proper Authentication Suites
		start, end := cells[i][0], len(o)
		if i != len(cells)-1 {
			end = cells[i+1][0]
		}
		// Narrow down the scope when looking for WPA Tag
		wpa2SearchArea := o[start:end]
		l := wpa2RE.FindIndex(wpa2SearchArea)
		if l == nil {
			res = append(res, Option{essid, NotSupportedProto})
			continue
		}
		// Narrow down the scope when looking for Authorization Suites
		authSearchArea := wpa2SearchArea[l[0]:]
		authSuites := strings.Trim(strings.Split(string(authSuitesRE.Find(authSearchArea)), ":")[1], "\n ")
		switch authSuites {
		case "PSK":
			res = append(res, Option{essid, WpaPsk})
		case "802.1x":
			res = append(res, Option{essid, WpaEap})
		default:
			res = append(res, Option{essid, NotSupportedProto})
		}
	}
	return res
}

func (w *IWLWorker) GetID() (string, error) {
	o, err := exec.Command("iwgetid", "-r").CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.Trim(string(o), " \n"), nil
}

func (w *IWLWorker) Connect(a ...string) error {
	// format of a: [essid, pass, id]
	conf, err := generateConfig(a...)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile("/tmp/wifi.conf", conf, 0444); err != nil {
		return fmt.Errorf("/tmp/wifi.conf: %v", err)
	}

	// Each request has a 30 second window to make a connection
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	c := make(chan error, 1)

	// There's no telling how long the supplicant will take, but on the other hand,
	// it's been almost instantaneous. But, further, it needs to keep running.
	go func() {
		cmd := exec.CommandContext(ctx, "wpa_supplicant", "-i"+w.Interface, "-c/tmp/wifi.conf")
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr //For an easier time debugging
		cmd.Run()
	}()

	// dhclient might never return on incorrect passwords or identity
	go func() {
		cmd := exec.CommandContext(ctx, "dhclient", "-ipv4=true", "-ipv6=false", "-verbose", w.Interface)
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr //For an easier time debugging
		if err := cmd.Run(); err != nil {
			c <- err
		} else {
			c <- nil
		}
	}()

	select {
	case err := <-c:
		return err
	case <-ctx.Done():
		return fmt.Errorf("Connection timeout")
	}
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
