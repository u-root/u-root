package main

import (
	"fmt"
	"net"
)

func interfaceExists(ifname string) func() error {
	return func() error {
		_, err := net.InterfaceByName(ifname)
		return err
	}
}

func interfaceRemediate(ifname string) func() error {
	return func() error {
		// TODO implement driver checking logic
		dmesg, err := getDmesg()
		if err != nil {
			return fmt.Errorf("cannot read dmesg to look for NIC driver information: %v", err)
		}
		lines := grep(dmesg, ifname)
		if len(lines) == 0 {
			return fmt.Errorf("no trace of %s in dmesg", ifname)
		}
		// TODO should this be returned as a string to the caller?
		fmt.Printf("  found %d references to %s in dmesg\n", len(lines), ifname)
		return nil
	}
}
