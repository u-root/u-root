// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

func downloadFile(filepath string, url string) error {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer closeIO(resp.Body, &err)

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer closeFile(out, &err)

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// Allows error handling on deferred file.Close()
func closeFile(f *os.File, err *error) {
	e := f.Close()
	switch *err {
	case nil:
		*err = e
	default:
		if e != nil {
			log.Println("Error:", e)
		}
	}
}

// Allows error handling on deferred io.Close()
func closeIO(c io.Closer, err *error) {
	e := c.Close()
	switch *err {
	case nil:
		*err = e
	default:
		if e != nil {
			log.Println("Error:", e)
		}
	}
}

func configureDHCPNetwork() error {
	if *verbose {
		log.Printf("Trying to configure network configuration dynamically...")
	}

	link, err := findNetworkInterface(*ifName)
	if err != nil {
		return err
	}

	var links []netlink.Link
	links = append(links, link)

	var level dhclient.LogLevel

	config := dhclient.Config{
		Timeout:  dhcpTimeout,
		Retries:  dhcpTries,
		LogLevel: level,
	}

	r := dhclient.SendRequests(context.TODO(), links, true, false, config, 20*time.Second)
	for result := range r {
		if result.Err == nil {
			return result.Lease.Configure()
		}
		log.Printf("dhcp response error: %v", result.Err)
	}
	return errors.New("no valid DHCP configuration recieved")
}

func findNetworkInterface(ifName string) (netlink.Link, error) {
	if ifName != "" {
		if *verbose {
			log.Printf("Try using %s", ifName)
		}
		link, err := netlink.LinkByName(ifName)
		if err == nil {
			return link, nil
		}
		log.Print(err)
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	if len(ifaces) == 0 {
		return nil, errors.New("no network interface found")
	}

	var ifnames []string
	for _, iface := range ifaces {
		ifnames = append(ifnames, iface.Name)
		// skip loopback
		if iface.Flags&net.FlagLoopback != 0 || iface.HardwareAddr.String() == "" {
			continue
		}
		if *verbose {
			log.Printf("Try using %s", iface.Name)
		}
		link, err := netlink.LinkByName(iface.Name)
		if err == nil {
			return link, nil
		}
		log.Print(err)
	}

	return nil, fmt.Errorf("could not find a non-loopback network interface with hardware address in any of %v", ifnames)
}
