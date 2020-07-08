// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/u-root/u-root/pkg/checker"
)

// fixmynetboot is a troubleshooting tool that can help you identify issues that
// won't let your system boot over the network.
// NOTE: this is a DEMO tool. It's here only to show you how to write your own
// checks and remediations. Don't use it in production.

var emergencyShellBanner = `
**************************************************************************
** Interface checks failed, see the output above to debug the issue.     *
** Entering the emergency shell, where you can run "fixmynetboot" again, *
** or any other LinuxBoot command.                                       *
**************************************************************************
`

var (
	doEmergencyShell = flag.Bool("shell", false, "Run emergency shell if checks fail")
)

func checkInterface(ifname string) error {
	checklist := []checker.Check{
		{
			Name:        fmt.Sprintf("%s exists", ifname),
			Run:         checker.InterfaceExists(ifname),
			Remediate:   checker.InterfaceRemediate(ifname),
			StopOnError: true,
		},
		{
			Name:        fmt.Sprintf("%s link speed", ifname),
			Run:         checker.LinkSpeed(ifname, 400000),
			Remediate:   nil,
			StopOnError: false},
		{
			Name:        fmt.Sprintf("%s link autoneg", ifname),
			Run:         checker.LinkAutoneg(ifname, true),
			Remediate:   nil,
			StopOnError: false,
		},
		{
			Name:        fmt.Sprintf("%s has link-local", ifname),
			Run:         checker.InterfaceHasLinkLocalAddress(ifname),
			Remediate:   nil,
			StopOnError: true,
		},
		{
			Name:        fmt.Sprintf("%s has global addresses", ifname),
			Run:         checker.InterfaceHasGlobalAddresses("eth0"),
			Remediate:   nil,
			StopOnError: true,
		},
	}

	return checker.Run(checklist)
}

func getNonLoopbackInterfaces() ([]string, error) {
	var interfaces []string
	allInterfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range allInterfaces {
		if iface.Flags&net.FlagLoopback == 0 {
			interfaces = append(interfaces, iface.Name)
		}
	}
	return interfaces, nil
}

func main() {
	flag.Parse()
	var (
		interfaces []string
		err        error
	)
	ifname := flag.Arg(0)
	if ifname == "" {
		interfaces, err = getNonLoopbackInterfaces()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		interfaces = []string{ifname}
	}

	for _, ifname := range interfaces {
		if err := checkInterface(ifname); err != nil {
			if !*doEmergencyShell {
				log.Fatal(err)
			}
			if err := checker.EmergencyShell(emergencyShellBanner)(); err != nil {
				log.Fatal(err)
			}
		}
	}
}
