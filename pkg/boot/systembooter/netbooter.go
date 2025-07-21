// Copyright 2017-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package systembooter

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/u-root/u-root/pkg/ulog"
)

// NetBooter implements the Booter interface for booting over DHCPv6.
// See NewNetBooterDHCPv6 for details on the fields.
type NetBooter struct {
	Type           string  `json:"type"`
	Method         string  `json:"method"`
	MAC            string  `json:"mac"`
	OverrideURL    *string `json:"override_url,omitempty"`
	Retries        *int    `json:"retries,omitempty"`
	DebugOnFailure bool    `json:"debug_on_failure,omitempty"`
}

// NewNetBooter parses a boot entry config and returns a Booter instance, or an
// error if any
func NewNetBooter(config []byte, l ulog.Logger) (Booter, error) {
	// The configuration format for a NetBooterDHCPv6 entry is a JSON with the
	// following structure:
	// {
	//     "type": "netboot",
	//     "method": "<method>",
	//     "mac": "<mac_addr>",
	//     "override_url": "<url>",
	//     "retries": <num_retries>,
	//     "debug_on_failure": <true|false>
	// }
	//
	// `type` is always set to "netboot".
	// `method` is one of "dhcpv6", "slaac" or "dhcpv4".
	// `mac` is the MAC address of the interface to use to boot.
	// `override_url` is an optional URL used to override the boot file URL used
	//   to fetch the network boot program. This field becomes mandatory if
	//   `method` is set to "slaac".
	// `retries` is the number of times a DHCP request should be retried if
	//   failed. If unspecified, it will use the underlying `netboot` program's
	//   default.
	// `debug_on_failure` is an optional boolean that will signal a request for
	//   a debugging attempt if netboot fails.
	//
	// An example configuration is:
	// {
	//     "type": "netboot",
	//     "method": "dhcpv6",
	//     "mac": "aa:bb:cc:dd:ee:ff",
	//     "override_url": "http://[fe80::face:booc]:8080/path/to/boot/file"
	// }
	//
	// Note that the optional `override_url` in the example above will override
	// whatever URL is returned in the OPT_BOOTFILE_URL for DHCPv6, or TFTP server
	// name + bootfile URL in case of DHCPv4.
	//
	// Additional options may be added in the future.
	l.Printf("Trying NetBooter...")
	l.Printf("Config: %s", string(config))
	nb := NetBooter{}
	if err := json.Unmarshal(config, &nb); err != nil {
		return nil, err
	}
	l.Printf("NetBooter: %+v", nb)
	if nb.Type != "netboot" {
		return nil, fmt.Errorf("%w:%q", errWrongType, nb.Type)
	}
	return &nb, nil
}

// Boot will run the boot procedure. In the case of NetBooter, it will call the
// `fbnetboot` command
func (nb *NetBooter) Boot(debugEnabled bool) error {
	var bootcmd []string
	l := ulog.Null
	if debugEnabled {
		bootcmd = []string{"fbnetboot", "-v", "-userclass", "linuxboot"}
		l = ulog.Log
	} else {
		bootcmd = []string{"fbnetboot", "-userclass", "linuxboot"}
	}
	if nb.OverrideURL != nil {
		bootcmd = append(bootcmd, "-netboot-url", *nb.OverrideURL)
	}
	if nb.Retries != nil {
		bootcmd = append(bootcmd, "-retries", strconv.Itoa(*nb.Retries))
	}
	if nb.Method == "dhcpv6" {
		bootcmd = append(bootcmd, []string{"-6=true", "-4=false"}...)
	} else if nb.Method == "dhcpv4" {
		bootcmd = append(bootcmd, []string{"-6=false", "-4=true"}...)
	} else {
		return fmt.Errorf("netboot: unknown method %s", nb.Method)
	}
	if nb.DebugOnFailure {
		bootcmd = append(bootcmd, "-fix")
	}
	l.Printf("Executing command: %v", bootcmd)
	cmd := exec.Command(bootcmd[0], bootcmd[1:]...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error executing %v: %w", cmd, err)
	}
	return nil
}

// TypeName returns the name of the booter type
func (nb *NetBooter) TypeName() string {
	return nb.Type
}
