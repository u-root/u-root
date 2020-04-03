// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/u-root/u-root/pkg/dhclient"
	"github.com/vishvananda/netlink"
)

const (
	entropyAvail       = "/proc/sys/kernel/random/entropy_avail"
	interfaceUpTimeout = 6 * time.Second
)

type netConf struct {
	HostIP         string `json:"host_ip"`
	DefaultGateway string `json:"gateway"`
	DNSServer      string `json:"dns"`
}

func getNetConf() (netConf, error) {
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
	info("Setup network configuration with IP: " + nc.HostIP)
	addr, err := netlink.ParseAddr(nc.HostIP)
	if err != nil {
		return fmt.Errorf("error parsing HostIP string to CIDR format address: %v", err)
	}
	gateway, err := netlink.ParseAddr(nc.DefaultGateway)
	if err != nil {
		return fmt.Errorf("error parsing DefaultGateway string to CIDR format address: %v", err)
	}
	if nc.DNSServer != "" {
		dns := net.ParseIP(nc.DNSServer)
		if dns == nil {
			return fmt.Errorf("cannot parse DNSServer string %s", nc.DNSServer)
		}
		resolvconf := fmt.Sprintf("nameserver %s\n", dns.String())
		if err = ioutil.WriteFile("/etc/resolv.conf", []byte(resolvconf), 0644); err != nil {
			return fmt.Errorf("could not write DNS servers to resolv.conf: %v", err)
		}
	}

	links, err := findNetworkInterfaces()
	if err != nil {
		return err
	}

	for _, link := range links {

		if err = netlink.AddrAdd(link, addr); err != nil {
			debug("%s: IP config failed: %v", link.Attrs().Name, err)
			continue
		}

		if err = netlink.LinkSetUp(link); err != nil {
			debug("%s: IP config failed: %v", link.Attrs().Name, err)
			continue
		}

		if err != nil {
			debug("%s: IP config failed: %v", link.Attrs().Name, err)
			continue
		}

		r := &netlink.Route{LinkIndex: link.Attrs().Index, Gw: gateway.IPNet.IP}
		if err = netlink.RouteAdd(r); err != nil {
			debug("%s: IP config failed: %v", link.Attrs().Name, err)
			continue
		}

		info("%s: IP configuration successful", link.Attrs().Name)
		return nil
	}
	return errors.New("IP configuration failed for all interfaces")
}

func configureDHCPNetwork() error {
	info("Trying to configure network configuration dynamically...")

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
			debug("%s: DHCP response error: %v", result.Interface.Attrs().Name, result.Err)
			continue
		}
		err = result.Lease.Configure()
		if err != nil {
			debug("%s: DHCP configuration error: %v", result.Interface.Attrs().Name, err)
		} else {
			info("%s: DHCP successful", result.Interface.Attrs().Name)
			return nil
		}
	}
	return errors.New("DHCP configuration failed")
}

func findNetworkInterfaces() ([]netlink.Link, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	if len(interfaces) == 0 {
		return nil, errors.New("no network interface found")
	}

	var links []netlink.Link
	var ifnames []string
	for _, i := range interfaces {
		debug("Found interface %s", i.Name)
		debug("    MTU: %d Hardware Addr: %s", i.MTU, i.HardwareAddr.String())
		debug("    Flags: %v", i.Flags)
		ifnames = append(ifnames, i.Name)
		// skip loopback
		if i.Flags&net.FlagLoopback != 0 || bytes.Compare(i.HardwareAddr, nil) == 0 {
			continue
		}
		link, err := netlink.LinkByName(i.Name)
		if err != nil {
			debug("%v", err)
		}
		links = append(links, link)
	}

	if len(links) <= 0 {
		return nil, fmt.Errorf("could not find a non-loopback network interface with hardware address in any of %v", ifnames)
	}

	return links, nil
}

func tryDownload(urls []string, file string) (dest string, err error) {
	dest = filepath.Join("/root", file)
	for _, rawurl := range urls {
		url, err := url.Parse(rawurl)
		if err != nil {
			debug("Skip %s : %v", rawurl, err)
			continue
		}

		url.Path = path.Join(url.Path, file)
		err = download(url.String(), dest)
		if err != nil {
			debug("%v", err)
			continue
		}
		return dest, nil
	}
	return "", fmt.Errorf("cannot find %s on provisioning servers", file)
}

func download(url string, destination string) error {
	roots, err := loadHTTPSCertificates()
	if err != nil {
		return fmt.Errorf("failed to load root certificate: %v", err)
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
		return fmt.Errorf("cannot evaluate entropy, %v", err)
	}
	es := strings.TrimSpace(string(e))
	entr, err := strconv.Atoi(es)
	if err != nil {
		return fmt.Errorf("cannot evaluate entropy, %v", err)
	}
	if entr < 128 {
		debug("WARNING: low entropy!")
		debug("%s : %d", entropyAvail, entr)
	}
	// get remote boot bundle
	info("Downloading %s", url)
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("HTTPS client: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTPS client: %d", resp.StatusCode)
	}
	f, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("Download: %v", err)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return fmt.Errorf("Download: %v", err)
	}

	return nil
}

// loadHTTPSCertificate loads the certificate needed
// for HTTPS and verifies it.
func loadHTTPSCertificates() (*x509.CertPool, error) {
	roots := x509.NewCertPool()
	bytes, err := data.get(httpsRootsFile)
	if err != nil {
		return roots, err
	}

	ok := roots.AppendCertsFromPEM(bytes)
	if !ok {
		return roots, fmt.Errorf("error parsing %s", httpsRootsFile)
	}
	return roots, nil
}

func hostHWAddr() (net.HardwareAddr, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return net.HardwareAddr{}, err
	}
	for _, i := range interfaces {
		if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
			return i.HardwareAddr, nil
		}
	}
	return net.HardwareAddr{}, fmt.Errorf("cannot find out hardware address")
}
