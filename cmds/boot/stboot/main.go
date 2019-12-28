package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/client4"
	"github.com/insomniacslk/dhcp/netboot"
	"github.com/u-root/u-root/pkg/boot/stboot"
	"github.com/vishvananda/netlink"
)

var debug = func(string, ...interface{}) {}

var (
	dryRun  = flag.Bool("dryrun", false, "Do everything except booting the loaded kernel")
	doDebug = flag.Bool("d", false, "Print debug output")
)

const (
	eth                = "eth0"
	rootCACertPath     = "/root/LetsEncrypt_Authority_X3.pem"
	entropyAvail       = "/proc/sys/kernel/random/entropy_avail"
	interfaceUpTimeout = 10 * time.Second
)

var banner = `
  _____ _______   _____   ____   ____________
 / ____|__   __|  |  _ \ / __ \ / __ \__   __|
| (___    | |     | |_) | |  | | |  | | | |   
 \___ \   | |     |  _ <| |  | | |  | | | |   
 ____) |  | |     | |_) | |__| | |__| | | |   
|_____/   |_|     |____/ \____/ \____/  |_|   

`

var check = `           
           //\\
OS is     //  \\
valid    //   //
        //   //
 //\\  //   //
//  \\//   //
\\        //
 \\      //
  \\    //
   \\__//
`

func main() {
	log.SetPrefix("stboot: ")

	flag.Parse()
	if *doDebug {
		debug = log.Printf
	}
	log.Print(banner)

	vars, err := stboot.FindHostVarsInInitramfs()
	if err != nil {
		log.Fatalf("Cant find Netvars at all: %v", err)
	}

	if *doDebug {
		str, _ := json.MarshalIndent(vars, "", "  ")
		log.Printf("Host variables: %s", str)
	}

	debug("Configuring network interfaces")

	if vars.HostIP != "" {
		err = configureStaticNetwork(vars, *doDebug)
	} else {
		err = configureDHCPNetwork()
	}

	if err != nil {
		log.Fatalf("Can not set up IO: %v", err)
	}

	ballPath := path.Join("root/", stboot.BallName)
	url, err := url.Parse(vars.BootstrapURL)
	if err != nil {
		log.Fatalf("Invalid bootstrap URL: %v", err)
	}
	url.Path = path.Join(url.Path, stboot.BallName)
	err = downloadFromHTTPS(url.String(), ballPath)
	if err != nil {
		log.Fatalf("Downloading bootball from %s failed: %v", url, err)
	}

	ball, err := stboot.BootBallFromArchie(ballPath)
	if err != nil {
		log.Fatal("Cannot open bootball")
	}

	// Just choose the first Bootconfig for now
	log.Printf("Pick the first boot configuration")
	var index = 0
	bc, err := ball.GetBootConfigByIndex(index)
	if err != nil {
		log.Fatalf("Cannot get boot configuration %d: %v", index, err)
	}

	if *doDebug {
		str, _ := json.MarshalIndent(*bc, "", "  ")
		log.Printf("Bootconfig: %s", str)
	}

	n, err := ball.VerifyBootconfigByName(bc.Name)
	if err != nil {
		log.Fatalf("Bootconfig %d seems to be not trustworthy: %v", index, err)
	}
	if n < vars.MinimalSignaturesMatch {
		log.Fatalf("Did not found enough valid signatures. %d valid, %d required", n, vars.MinimalSignaturesMatch)
	}
	log.Printf("Bootconfig '%s' passed verification", bc.Name)
	log.Print(check)

	if *dryRun {
		debug("Dryrun mode: will not boot")
		return
	}

	log.Println("Starting up new kernel.")

	if err := bc.Boot(); err != nil {
		log.Printf("Failed to boot kernel %s: %v", bc.Kernel, err)
	}
	// if we reach this point, no boot configuration succeeded
	log.Print("No boot configuration succeeded")

	return
}

func configureStaticNetwork(vars stboot.HostVars, doDebug bool) error {
	log.Printf("Setup network configuration with IP: " + vars.HostIP)
	addr, err := netlink.ParseAddr(vars.HostIP)
	if err != nil {
		return fmt.Errorf("Error parsing HostIP string to CIDR format address: %v", err)
	}

	iface, err := netlink.LinkByName(eth)
	if err = netlink.AddrAdd(iface, addr); err != nil {
		return fmt.Errorf("Error retrieving interface by name: %v", err)
	}

	if err = netlink.LinkSetUp(iface); err != nil {
		return fmt.Errorf("Error bringing up interface:%v with error: %v", eth, err)
	}

	gateway, err := netlink.ParseAddr(vars.DefaultGateway)
	if err != nil {
		return fmt.Errorf("Error parsing GatewayIP string to CIDR format address: %v", err)
	}

	r := &netlink.Route{LinkIndex: iface.Attrs().Index, Gw: gateway.IPNet.IP}
	if err = netlink.RouteAdd(r); err != nil {
		return fmt.Errorf("Error setting default gateway: %v", err)
	}

	return nil
}

func configureDHCPNetwork() error {

	log.Printf("Trying to configure network configuration dynamically...")
	attempts := 10
	var conversation []*dhcpv4.DHCPv4

	_, err := netboot.IfUp(eth, interfaceUpTimeout)
	if err != nil {
		return fmt.Errorf("Ifup with %s failed: %v", eth, err)
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
		return err
	}

	err = netboot.ConfigureInterface(eth, netbootConfig)

	if err != nil {
		return err
	}

	iface, err := netlink.LinkByName(eth)
	if err != nil {
		return fmt.Errorf("Error getting Link by Name: %v", err)
	}
	gateway, err := netlink.ParseAddr(netbootConfig.Routers[0].String() + "/24")
	if err != nil {
		return fmt.Errorf("Error parsing GatewayIP string to CIDR format address: %v", err)
	}
	r := &netlink.Route{LinkIndex: iface.Attrs().Index, Gw: gateway.IPNet.IP}
	if err = netlink.RouteAdd(r); err != nil {
		return fmt.Errorf("Error setting default gateway: %v", err)
	}

	return nil
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
