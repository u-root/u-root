package checker

import (
	"errors"
	"fmt"
	"net"
)

// InterfaceExists returns a Checker that verifies if an interface is present on
// the system
func InterfaceExists(ifname string) Checker {
	return func() error {
		_, err := net.InterfaceByName(ifname)
		return err
	}
}

func addresses(ifname string) ([]net.IP, error) {
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		return nil, err
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}
	iplist := make([]net.IP, 0)
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok {
			return nil, errors.New("not a net.IPNet")
		}
		iplist = append(iplist, ipnet.IP)
	}
	return iplist, nil
}

// InterfaceHasLinkLocalAddress returns a Checker that verifies if an interface
// has a configured link-local address.
func InterfaceHasLinkLocalAddress(ifname string) Checker {
	return func() error {
		addrs, err := addresses(ifname)
		if err != nil {
			return err
		}
		for _, addr := range addrs {
			if addr.IsLinkLocalUnicast() {
				return nil
			}
		}
		return fmt.Errorf("no link local addresses for interface %s", ifname)
	}
}

// InterfaceHasGlobalAddresses returns a Checker that verifies if an interface has
// at least one global address.
func InterfaceHasGlobalAddresses(ifname string) Checker {
	return func() error {
		addrs, err := addresses(ifname)
		if err != nil {
			return err
		}
		for _, addr := range addrs {
			if addr.IsGlobalUnicast() {
				return nil
			}
		}
		return fmt.Errorf("no unicast global addresses for interface %s", ifname)
	}
}

// InterfaceRemediate returns a Remediator that tries to fix a missing
// interface issue.
func InterfaceRemediate(ifname string) Remediator {
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
