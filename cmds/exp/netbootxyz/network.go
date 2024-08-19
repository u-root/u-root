// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

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
	"strconv"
	"strings"
	"time"

	"github.com/u-root/u-root/pkg/dhclient"
	"github.com/vishvananda/netlink"
)

// Counts number of bytes written
// Conforms to io.Writer interface
type ProgressCounter struct {
	Downloaded    uint64
	Total         uint64
	PreviousRatio int
	Writer        io.Writer
}

func (counter *ProgressCounter) Write(p []byte) (int, error) {
	n := len(p)
	counter.Downloaded += uint64(n)
	counter.PrintProgress()
	return n, nil
}

func (counter *ProgressCounter) PrintProgress() {
	ratio := int(float64(counter.Downloaded) / float64(counter.Total) * 100)

	// Only print every 5% to avoid spamming the serial port and making it look weird
	if ratio%5 == 0 && ratio != counter.PreviousRatio {
		// Clear the line by using a character return to go back to the start and
		// remove the remaining characters by filling it with spaces
		fmt.Fprintf(counter.Writer, "\r%s", strings.Repeat(" ", 50))

		fmt.Fprintf(counter.Writer, "\rDownloading... %s out of %s (%d%%)", bytesToHuman(counter.Downloaded), bytesToHuman(counter.Total), ratio)
		counter.PreviousRatio = ratio
	}

	if counter.Downloaded == counter.Total {
		fmt.Fprintf(counter.Writer, "\n")
	}
}

func bytesToHuman(bytes uint64) string {
	const unit = 1000 // Instead of 1024 so we'll get MB instead of MiB
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div := int64(unit)
	exponent := 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exponent++
	}
	return fmt.Sprintf("%4.1f %cB", float64(bytes)/float64(div), "kMGTPE"[exponent])
}

func downloadFile(filepath string, url string) error {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: false}

	log.Printf("Downloading file %s from %s\n", filepath, url)

	headResp, err := http.Head(url)
	if err != nil {
		return err
	}

	defer headResp.Body.Close()

	// Get size for progress indicator
	size, err := strconv.ParseUint(headResp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return err
	}

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
	counter := &ProgressCounter{Total: size, PreviousRatio: 0, Writer: os.Stdout}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))

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
