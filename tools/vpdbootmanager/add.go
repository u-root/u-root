// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/u-root/u-root/pkg/boot/systembooter"
	"github.com/u-root/u-root/pkg/vpd"
)

var dryRun = false

func add(entrytype string, args []string) error {
	var (
		entry  systembooter.Booter
		vpdDir string
		err    error
	)
	switch entrytype {
	case "netboot":
		if len(args) < 2 {
			return fmt.Errorf("you need to pass method and MAC address")
		}
		entry, vpdDir, err = parseNetbootFlags(args[0], args[1], args[2:])
		if err != nil {
			return err
		}
	case "localboot":
		if len(args) < 1 {
			return fmt.Errorf("you need to provide method")
		}
		entry, vpdDir, err = parseLocalbootFlags(args[0], args[1:])
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown entry type")
	}
	if dryRun {
		b, err := json.Marshal(entry)
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stderr, "Using -dryrun, will not write any variable. Content of boot entry:")
		fmt.Println(string(b))
		return nil
	}
	return addBootEntry(entry, vpdDir)
}

func parseLocalbootFlags(method string, args []string) (*systembooter.LocalBooter, string, error) {
	cfg := &systembooter.LocalBooter{
		Type:   "localboot",
		Method: method,
	}
	var flagVpdDir string
	flg := flag.NewFlagSet("localboot", flag.ExitOnError)
	flg.StringVar(&cfg.KernelArgs, "kernel-args", "", "additional kernel args")
	flg.StringVar(&cfg.Initramfs, "ramfs", "", "path of ramfs to be used for kexec'ing into the target kernel.")
	flg.StringVar(&flagVpdDir, "vpd-dir", vpd.DefaultVpdDir, "VPD dir to use")
	flg.BoolVar(&dryRun, "dryrun", false, "only print values that would be set")

	switch method {
	case "grub":
		if err := flg.Parse(args); err != nil {
			return nil, flagVpdDir, err
		}
	case "path":
		if len(args) < 2 {
			return nil, "", fmt.Errorf("you need to pass DeviceGUID and Kernel path")
		}
		cfg.DeviceGUID = args[0]
		cfg.Kernel = args[1]
		if err := flg.Parse(args[2:]); err != nil {
			return nil, flagVpdDir, err
		}
	default:
		return nil, flagVpdDir, fmt.Errorf("method needs to be grub or path")
	}
	return cfg, flagVpdDir, nil
}

func parseNetbootFlags(method, mac string, args []string) (*systembooter.NetBooter, string, error) {
	if method != "dhcpv4" && method != "dhcpv6" {
		return nil, "", fmt.Errorf("method needs to be either dhcpv4 or dhcpv6")
	}

	_, err := net.ParseMAC(mac)
	if err != nil {
		return nil, "", err
	}

	cfg := &systembooter.NetBooter{
		Type:   "netboot",
		Method: method,
		MAC:    mac,
	}
	var flagVpdDir string

	flg := flag.NewFlagSet("netboot", flag.ExitOnError)
	overrideURL := flg.String("override-url", "", "an optional URL used to override the boot file URL used")
	retries := flg.Int("retries", -1, "the number of times a DHCP request should be retried if failed.")
	flg.BoolVar(&dryRun, "dryrun", false, "only print values that would be set")
	flg.StringVar(&flagVpdDir, "vpd-dir", vpd.DefaultVpdDir, "VPD dir to use")
	if err := flg.Parse(args); err != nil {
		return nil, "", err
	}

	if *overrideURL != "" {
		cfg.OverrideURL = overrideURL
	}

	if *retries != -1 {
		cfg.Retries = retries
	}

	return cfg, flagVpdDir, nil
}

func addBootEntry(cfg systembooter.Booter, vpdDir string) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	vpdReader := vpd.NewReader()
	vpdReader.VpdDir = vpdDir
	for i := 1; i < vpd.MaxBootEntry; i++ {
		key := fmt.Sprintf("Boot%04d", i)
		if _, err := vpdReader.Get(key, false); err != nil {
			if os.IsNotExist(err) {
				if err := vpdReader.Set(key, data, false); err != nil {
					return err
				}
				return nil
			}
			return err
		}
	}
	return errors.New("maximum number of boot entries already set")
}

func set(key string, value string) error {
	return vpd.FlashromRWVpdSet(key, []byte(value), false)
}

func remove(key string) error {
	return vpd.FlashromRWVpdSet(key, []byte("dummy"), true)
}

func dump() error {
	return vpd.FlashromVpdDump()
}
