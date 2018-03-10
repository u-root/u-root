package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/u-root/dhcp4/dhcp4client"
	"github.com/u-root/u-root/pkg/dhclient"
	"github.com/u-root/u-root/pkg/kexec"
	"github.com/u-root/u-root/pkg/pxe"
	"github.com/vishvananda/netlink"
)

var (
	verbose = flag.Bool("v", true, "print all kinds of things out, more than Chris wants")
	dryRun  = flag.Bool("dry-run", false, "download kernel, but don't kexec it")
	debug   = func(string, ...interface{}) {}
)

func copyToFile(r io.Reader) (*os.File, error) {
	f, err := ioutil.TempFile("", "nerf-netboot")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		return nil, err
	}
	if err := f.Sync(); err != nil {
		return nil, err
	}
	// Re-open the file read-only so we don't get ETXTBSY from kexec.
	readOnlyFile, err := os.Open(f.Name())
	if err != nil {
		return nil, err
	}
	return readOnlyFile, nil
}

func attemptDHCPLease(iface netlink.Link, timeout time.Duration, retry int) (*dhclient.Packet4, error) {
	if _, err := dhclient.IfUp(iface.Attrs().Name); err != nil {
		return nil, err
	}

	client, err := dhcp4client.New(iface,
		dhcp4client.WithTimeout(timeout),
		dhcp4client.WithRetry(retry))
	if err != nil {
		return nil, err
	}

	p, err := client.Request()
	if err != nil {
		return nil, err
	}
	return dhclient.NewPacket4(p), nil
}

func Netboot() error {
	ifs, err := netlink.LinkList()
	if err != nil {
		return err
	}

	for _, iface := range ifs {
		// TODO: Do 'em all in parallel.
		if iface.Attrs().Name != "eth0" {
			continue
		}

		log.Printf("Attempting to get DHCP lease on %s", iface.Attrs().Name)
		packet, err := attemptDHCPLease(iface, 30*time.Second, 5)
		if err != nil {
			log.Printf("No lease on %s: %v", iface.Attrs().Name, err)
			continue
		}
		log.Printf("Got lease on %s", iface.Attrs().Name)
		if err := dhclient.Configure4(iface, packet.P); err != nil {
			log.Printf("shit: %v", err)
			continue
		}

		// We may have to make this DHCPv6 and DHCPv4-specific anyway.
		// Only tested with v4 right now; and assuming the uri points
		// to a pxelinux.0.
		//
		// Or rather, we need to make this option-specific. DHCPv6 has
		// options for passing a kernel and cmdline directly. v4
		// usually just passes a pxelinux.0. But what about an initrd?
		uri, err := packet.Boot()
		if err != nil {
			log.Printf("Got DHCP lease, but no valid PXE information.")
			continue
		}

		log.Printf("Boot URI: %v", uri)

		wd := &url.URL{
			Scheme: uri.Scheme,
			Host:   uri.Host,
			Path:   path.Dir(uri.Path),
		}
		pc := pxe.NewConfig(wd)
		if err := pc.FindConfigFile(iface.Attrs().HardwareAddr, packet.Lease().IP); err != nil {
			return fmt.Errorf("failed to parse pxelinux config: %v", err)
		}

		label := pc.Entries[pc.DefaultEntry]
		log.Printf("Got configuration: %v", label)

		k, err := copyToFile(label.Kernel)
		if err != nil {
			return err
		}
		defer k.Close()

		var i *os.File
		if label.Initrd != nil {
			i, err = copyToFile(label.Initrd)
			if err != nil {
				return err
			}
			defer i.Close()
		}

		if *dryRun {
			log.Printf("Kernel: %s", k.Name())
			if i != nil {
				log.Printf("Initrd: %s", i.Name())
			}
			log.Printf("Command line: %s", label.Cmdline)
		} else {
			if err := kexec.FileLoad(k, i, label.Cmdline); err != nil {
				return err
			}
			return kexec.Reboot()
		}
	}
	return nil
}

func main() {
	flag.Parse()
	if *verbose {
		debug = log.Printf
	}

	if err := Netboot(); err != nil {
		log.Fatal(err)
	}
}
