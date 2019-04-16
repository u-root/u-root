package main

import (
	"errors"
	"fmt"
	"net"
)

func interfaceExists(ifname string) Checker {
	return func() error {
		_, err := net.InterfaceByName(ifname)
		return err
	}
}

func interfaceHasLinkLocalAddress(ifname string) Checker {
	return func() error {
		iface, err := net.InterfaceByName(ifname)
		if err != nil {
			return err
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return err
		}
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok {
				return errors.New("not a net.IPNet")
			}
			if ipnet.IP.IsLinkLocalUnicast() {
				return nil
			}
		}
		return fmt.Errorf("no link local addresses for interface %s", ifname)
	}
}

func interfaceRemediate(ifname string) Remediator {
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
