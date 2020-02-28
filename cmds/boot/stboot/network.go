// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
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

	"github.com/u-root/u-root/pkg/dhclient"
	"github.com/vishvananda/netlink"
)

type netConf struct {
	HostIP         string `json:"host_ip"`
	HostNetmask    string `json:"netmask"`
	DefaultGateway string `json:"gateway"`
	DNSServer      string `json:"dns"`
}

func getNetConf() (netConf, error) {
	data := initramfsData{}
	bytes, err := data.get(networkFile)
	var net netConf
	if err != nil {
		return net, err
	}
	err = json.Unmarshal(bytes, &net)
	if err != nil {
		return net, err
	}
	return net, nil
}

func configureStaticNetwork(nc netConf) error {
	log.Printf("Setup network configuration with IP: " + nc.HostIP)
	addr, err := netlink.ParseAddr(nc.HostIP)
	if err != nil {
		return fmt.Errorf("Error parsing HostIP string to CIDR format address: %v", err)
	}

	links, err := findNetworkInterfaces()
	if err != nil {
		return err
	}

	for _, link := range links {

		if err = netlink.AddrAdd(link, addr); err != nil {
			if *doDebug {
				log.Printf("%s: IP config failed: %v", link.Attrs().Name, err)
			}
			continue
		}

		if err = netlink.LinkSetUp(link); err != nil {
			if *doDebug {
				log.Printf("%s: IP config failed: %v", link.Attrs().Name, err)
			}
			continue
		}

		gateway, err := netlink.ParseAddr(nc.DefaultGateway)
		if err != nil {
			if *doDebug {
				log.Printf("%s: IP config failed: %v", link.Attrs().Name, err)
			}
			continue
		}

		r := &netlink.Route{LinkIndex: link.Attrs().Index, Gw: gateway.IPNet.IP}
		if err = netlink.RouteAdd(r); err != nil {
			if *doDebug {
				log.Printf("%s: IP config failed: %v", link.Attrs().Name, err)
			}
		} else {
			log.Printf("%s: IP configuration successful", link.Attrs().Name)
			return nil
		}
	}
	return errors.New("IP configuration failed")
}

func configureDHCPNetwork() error {
	log.Printf("Trying to configure network configuration dynamically...")

	links, err := findNetworkInterfaces()
	if err != nil {
		return err
	}

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

	r := dhclient.SendRequests(context.TODO(), links, true, false, config)
	for result := range r {
		if result.Err != nil {
			if *doDebug {
				log.Printf("%s: DHCP response error: %v", result.Interface.Attrs().Name, result.Err)
			}
			continue
		}
		err = result.Lease.Configure()
		if err != nil {
			if *doDebug {
				log.Printf("%s: DHCP configuration error: %v", result.Interface.Attrs().Name, err)
			}
		} else {
			log.Printf("%s: DHCP successful", result.Interface.Attrs().Name)
			return nil
		}
	}
	return errors.New("DHCP configuration failed")
}

func findNetworkInterfaces() ([]netlink.Link, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	if len(ifaces) == 0 {
		return nil, errors.New("No network interface found")
	}

	var links []netlink.Link
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
		link, err := netlink.LinkByName(iface.Name)
		if err != nil {
			log.Print(err)
		}
		links = append(links, link)
	}

	if len(links) <= 0 {
		return nil, fmt.Errorf("Could not find a non-loopback network interface with hardware address in any of %v", ifnames)
	}

	return links, nil
}

func downloadFromHTTPS(url string, destination string) error {
	roots, err := loadHTTPSCertificates()
	if err != nil {
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
	if entr < 128 {
		log.Print("WARNING: low entropy!")
		log.Printf("%s : %d", entropyAvail, entr)
	}
	// get remote boot bundle
	log.Printf("Downloading from %s", url)
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
func loadHTTPSCertificates() (*x509.CertPool, error) {
	roots := x509.NewCertPool()
	data := initramfsData{}
	bytes, err := data.get(httpsRootsFile)
	if err != nil {
		return roots, err
	}
	ok := roots.AppendCertsFromPEM(bytes)
	if !ok {
		return roots, fmt.Errorf("Error parsing %s", httpsRootsFile)
	}
	return roots, nil
}
