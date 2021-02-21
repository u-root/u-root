package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/u-root/u-root/pkg/dhclient"
	"github.com/vishvananda/netlink"
)

func DownloadFile(filepath string, url string) error {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func configureDHCPNetwork() error {
	log.Printf("Trying to configure network configuration dynamically...")

	link, err := findNetworkInterface()
	if err != nil {
		return err
	}

	var links []netlink.Link
	links = append(links, link)

	var level dhclient.LogLevel

	config := dhclient.Config{
		Timeout:  6 * time.Second,
		Retries:  4,
		LogLevel: level,
	}

	r := dhclient.SendRequests(context.TODO(), links, true, false, config, 30*time.Second)
	for result := range r {
		if result.Err == nil {
			return result.Lease.Configure()
		} else {
			log.Printf("dhcp response error: %v", result.Err)
		}
	}
	return errors.New("no valid DHCP configuration recieved")
}

func findNetworkInterface() (netlink.Link, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	if len(ifaces) == 0 {
		return nil, errors.New("No network interface found")
	}

	var ifnames []string
	for _, iface := range ifaces {
		ifnames = append(ifnames, iface.Name)
		// skip loopback
		if iface.Flags&net.FlagLoopback != 0 || iface.HardwareAddr.String() == "" {
			continue
		}
		log.Printf("Try using %s", iface.Name)
		link, err := netlink.LinkByName(iface.Name)
		if err == nil {
			return link, nil
		}
		log.Print(err)
	}

	return nil, fmt.Errorf("Could not find a non-loopback network interface with hardware address in any of %v", ifnames)
}
