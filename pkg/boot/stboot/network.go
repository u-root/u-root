package stboot

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/client4"
	"github.com/insomniacslk/dhcp/netboot"
)

const (
	eth = "eth0"
	//BootFilePath       = "root/stboot.zip"
	rootCACertPath     = "/root/LetsEncrypt_Authority_X3.pem"
	entropyAvail       = "/proc/sys/kernel/random/entropy_avail"
	interfaceUpTimeout = 10 * time.Second
)

// ConfigureStaticNetwork sets up your eth interface from netvars.json
func ConfigureStaticNetwork(vars HostVars, doDebug bool) error {
	//setup ip
	log.Printf("Setup network configuration with IP: " + vars.HostIP)
	cmd := exec.Command("ip", "addr", "add", vars.HostIP, "dev", eth)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error executing %v: %v", cmd, err)
		return err
	}
	cmd = exec.Command("ip", "link", "set", eth, "up")
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error executing %v: %v", cmd, err)
		return err
	}
	cmd = exec.Command("ip", "route", "add", "default", "via", vars.DefaultGateway, "dev", eth)
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error executing %v: %v", cmd, err)
		return err
	}

	if doDebug {
		cmd = exec.Command("ip", "addr")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Printf("Error executing %v: %v", cmd, err)
		}
		cmd = exec.Command("ip", "route")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Printf("Error executing %v: %v", cmd, err)
		}
	}

	return nil
}

// ConfigureDHCPNetwork configures DHCP on eth0
func ConfigureDHCPNetwork() error {

	log.Printf("Trying to configure network configuration dynamically..")
	attempts := 10
	var conversation []*dhcpv4.DHCPv4

	_, err := netboot.IfUp(eth, interfaceUpTimeout)
	if err != nil {
		log.Println("Enabling eth0 failed.")
		return fmt.Errorf("Ifup failed: %v", err)
	}
	if attempts < 1 {
		attempts = 1
	}

	client := client4.NewClient()
	for attempt := 0; attempt < attempts; attempt++ {
		log.Printf("Attempt to get DHCP lease %d of %d for interface %s", attempt+1, attempts, eth)
		conversation, err = client.Exchange(eth)

		if err != nil && attempt < attempts {
			log.Printf("Error: %v", err)
			continue
		}
		break
	}

	if conversation[3] == nil {
		return fmt.Errorf("Gateway is null")
	}
	netbootConfig, err := netboot.GetNetConfFromPacketv4(conversation[3])

	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	err = netboot.ConfigureInterface(eth, netbootConfig)

	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	// Some manual shit - for now
	cmd := exec.Command("ip", "route", "add", "default", "via", netbootConfig.Routers[0].String()+"/24", "dev", eth)
	if err := cmd.Run(); err != nil {
		log.Printf("Error executing %v: %v", cmd, err)
		return err
	}

	return nil
}

// DownloadFromHTTPS downloads the stboot.zip file
// to a specific destination via HTTPS.
func DownloadFromHTTPS(url string, destination string) error {

	roots := x509.NewCertPool()
	if err := loadHTTPSCertificate(roots); err != nil {
		return fmt.Errorf("Failed to verify root certificate: %v", err)
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
	es := strings.TrimSpace(string(e))
	entr, err := strconv.Atoi(es) // XXX: Insecure?
	if err != nil {
		return fmt.Errorf("Cannot evaluate entropy, %v", err)
	}
	log.Printf("Available kernel entropy: %d", entr)
	if entr < 128 {
		log.Print("WARNING: low entropy!")
		log.Printf("%s : %d", entropyAvail, entr)
	}
	// get remote boot bundle
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
		return fmt.Errorf("failed to write boot config file: %v", err)
	}

	return nil
}

// loadHTTPSCertificate loads the certificate needed
// for HTTPS and verifies it.
func loadHTTPSCertificate(roots *x509.CertPool) error {
	// load CA certificate
	log.Printf("Load %s as CA certificate", rootCACertPath)
	rootCertBytes, err := ioutil.ReadFile(rootCACertPath)
	if err != nil {
		return fmt.Errorf("Failed to read CA root certificate file: %v", err)
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
