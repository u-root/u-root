// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/u-root/u-root/pkg/boot/stboot"
	"github.com/u-root/u-root/pkg/dhclient"
	"github.com/vishvananda/netlink"
)

func configureStaticNetwork(vars stboot.HostVars) error {
	log.Printf("Setup network configuration with IP: " + vars.HostIP)
	addr, err := netlink.ParseAddr(vars.HostIP)
	if err != nil {
		return fmt.Errorf("Error parsing HostIP string to CIDR format address: %v", err)
	}

	link, err := findNetworkInterface()
	if err != nil {
		return err
	}

	if err = netlink.AddrAdd(link, addr); err != nil {
		return fmt.Errorf("Error adding address: %v", err)
	}

	if err = netlink.LinkSetUp(link); err != nil {
		return fmt.Errorf("Error bringing up interface: %v", err)
	}

	gateway, err := netlink.ParseAddr(vars.DefaultGateway)
	if err != nil {
		return fmt.Errorf("Error parsing GatewayIP string to CIDR format address: %v", err)
	}

	r := &netlink.Route{LinkIndex: link.Attrs().Index, Gw: gateway.IPNet.IP}
	if err = netlink.RouteAdd(r); err != nil {
		return fmt.Errorf("Error setting default gateway: %v", err)
	}

	return nil
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
	if *doDebug {
		level = 1
	} else {
		level = 0
	}
	config := dhclient.Config{
		Timeout:  interfaceUpTimeout,
		Retries:  4,
		LogLevel: level,
	}

	r := dhclient.SendRequests(context.TODO(), links, true, false, config, 30*time.Second)
	for result := range r {
		if result.Err == nil {
			return result.Lease.Configure()
		} else if *doDebug {
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
		if *doDebug {
			log.Printf("Found interface %s", iface.Name)
			log.Printf("    MTU: %d Hardware Addr: %s", iface.MTU, iface.HardwareAddr.String())
			log.Printf("    Flags: %v", iface.Flags)
		}
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

func downloadFromHTTPS(url string, destination string) error {
	roots := x509.NewCertPool()
	if err := loadHTTPSCertificate(roots); err != nil {
		return fmt.Errorf("Failed to load root certificate: %v", err)
	}

	// setup https client
	client := http.Client{
		Transport: (&http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: (&tls.Config{
				RootCAs: roots,
			}),
		}),
	}

	// check available kernel entropy
	e, err := ioutil.ReadFile(entropyAvail)
	if err != nil {
		return fmt.Errorf("Cannot evaluate entropy, %v", err)
	}
	es := strings.TrimSpace(string(e))
	entr, err := strconv.Atoi(es)
	if err != nil {
		return fmt.Errorf("Cannot evaluate entropy, %v", err)
	}
	log.Printf("Available kernel entropy: %d", entr)
	if entr < 128 {
		log.Print("WARNING: low entropy!")
		log.Printf("%s : %d", entropyAvail, entr)
	}
	// get remote boot bundle
	log.Print("Downloading bootball ...")
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("non-200 HTTP status: %d", resp.StatusCode)
	}
	f, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("failed create boot config file: %v", err)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save bootball: %v", err)
	}

	return nil
}

// loadHTTPSCertificate loads the certificate needed
// for HTTPS and verifies it.
func loadHTTPSCertificate(roots *x509.CertPool) error {
	log.Printf("Load %s as CA certificate", rootCACertPath)
	rootCertBytes, err := ioutil.ReadFile(rootCACertPath)
	if err != nil {
		return err
	}
	rootCertPem, _ := pem.Decode(rootCertBytes)
	if rootCertPem.Type != "CERTIFICATE" {
		return fmt.Errorf("Failed decoding certificate: Certificate is of the wrong type. PEM Type is: %s", rootCertPem.Type)
	}
	ok := roots.AppendCertsFromPEM([]byte(rootCertBytes))
	if !ok {
		return fmt.Errorf("Error parsing CA root certificate")
	}
	return nil
}
