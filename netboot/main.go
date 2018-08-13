package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/insomniacslk/dhcp/iana"
	"github.com/insomniacslk/dhcp/netboot"
	"github.com/u-root/u-root/pkg/kexec"
)

var (
	useV4              = flag.Bool("4", false, "Get a DHCPv4 lease")
	useV6              = flag.Bool("6", true, "Get a DHCPv6 lease")
	ifname             = flag.String("i", "eth0", "Interface to send packets through")
	dryRun             = flag.Bool("dryrun", false, "Do everything except assigning IP addresses, changing DNS, and kexec")
	doDebug            = flag.Bool("d", false, "Print debug output")
	skipDHCP           = flag.Bool("skip-dhcp", false, "Skip DHCP and rely on SLAAC for network configuration. This requires -netboot-url")
	overrideNetbootURL = flag.String("netboot-url", "", "Override the netboot URL normally obtained via DHCP")
	readTimeout        = flag.Int("timeout", 3, "Read timeout in seconds")
	dhcpRetries        = flag.Int("retries", 3, "Number of times a DHCP request is retried")
	userClass          = flag.String("userclass", "", "Override DHCP User Class option")
)

const (
	interfaceUpTimeout = 30 * time.Second
)

var banner = `

 _________________________________
< Net booting is so hot right now >
 ---------------------------------
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||

`

func main() {
	flag.Parse()
	if *skipDHCP && *overrideNetbootURL == "" {
		log.Fatal("-skip-dhcp requires -netboot-url")
	}
	debug := func(string, ...interface{}) {}
	if *doDebug {
		debug = log.Printf
	}
	log.Print(banner)

	if !*useV6 && !*useV4 {
		log.Fatal("At least one of DHCPv6 and DHCPv4 is required")
	}
	// DHCPv6
	if *useV6 {
		log.Printf("Trying to obtain a DHCPv6 lease on %s", *ifname)
		log.Printf("Waiting for network interface %s to come up", *ifname)
		start := time.Now()
		_, err := netboot.IfUp(*ifname, interfaceUpTimeout)
		if err != nil {
			log.Fatalf("DHCPv6: IfUp failed: %v", err)
		}
		debug("Interface %s is up after %v", *ifname, time.Since(start))
		var (
			netconf  *netboot.NetConf
			bootfile string
		)
		if *skipDHCP {
			log.Print("Skipping DHCP")
		} else {
			// send a netboot request via DHCP
			modifiers := []dhcpv6.Modifier{
				dhcpv6.WithArchType(iana.EFI_X86_64),
			}
			if *userClass != "" {
				modifiers = append(modifiers, dhcpv6.WithUserClass([]byte(*userClass)))
			}
			conversation, err := netboot.RequestNetbootv6(*ifname, time.Duration(*readTimeout)*time.Second, *dhcpRetries, modifiers...)
			for _, m := range conversation {
				debug(m.Summary())
			}
			if err != nil {
				log.Fatalf("DHCPv6: netboot request for interface %s failed: %v", *ifname, err)
			}
			// get network configuration and boot file
			netconf, bootfile, err = netboot.ConversationToNetconf(conversation)
			if err != nil {
				log.Fatalf("DHCPv6: failed to extract network configuration for %s: %v", *ifname, err)
			}
			debug("DHCPv6: network configuration: %+v", netconf)
			if !*dryRun {
				// Set up IP addresses
				log.Printf("DHCPv6: configuring network interface %s", *ifname)
				if err = netboot.ConfigureInterface(*ifname, netconf); err != nil {
					log.Fatalf("DHCPv6: cannot configure IPv6 addresses on interface %s: %v", *ifname, err)
				}
				// Set up DNS
			}
			if *overrideNetbootURL != "" {
				bootfile = *overrideNetbootURL
			}
			log.Printf("DHCPv6: boot file for interface %s is %s", *ifname, bootfile)
		}
		if *overrideNetbootURL != "" {
			bootfile = *overrideNetbootURL
		}
		debug("DHCPv6: boot file URL is %s", bootfile)
		// check for supported schemes
		if !strings.HasPrefix(bootfile, "http://") {
			log.Fatal("DHCPv6: can only handle http scheme")
		}

		log.Printf("DHCPv6: fetching boot file URL: %s", bootfile)
		resp, err := http.Get(bootfile)
		if err != nil {
			log.Fatalf("DHCPv6: http.Get of %s failed: %v", bootfile, err)
		}
		// FIXME this will not be called if something fails after this point
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			log.Fatalf("Status code is not 200 OK: %d", resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("DHCPv6: cannot read boot file from the network: %v", err)
		}
		u, err := url.Parse(bootfile)
		if err != nil {
			log.Fatalf("DHCPv6: cannot parse URL %s: %v", bootfile, err)
		}
		// extract file name component
		if strings.HasSuffix(u.Path, "/") {
			log.Fatalf("Invalid file path, cannot end with '/': %s", u.Path)
		}
		filename := filepath.Base(u.Path)
		if filename == "." || filename == "" {
			log.Fatalf("Invalid empty file name extracted from file path %s", u.Path)
		}
		if err = ioutil.WriteFile(filename, body, 0400); err != nil {
			log.Fatalf("DHCPv6: cannot write to file %s: %v", filename, err)
		}
		debug("DHCPv6: saved boot file to %s", filename)
		if !*dryRun {
			log.Printf("DHCPv6: kexec'ing into %s", filename)
			kernel, err := os.OpenFile(filename, os.O_RDONLY, 0)
			if err != nil {
				log.Fatalf("DHCPv6: cannot open file %s: %v", filename, err)
			}
			if err = kexec.FileLoad(kernel, nil /* ramfs */, "" /* cmdline */); err != nil {
				log.Fatalf("DHCPv6: kexec.FileLoad failed: %v", err)
			}
			if err = kexec.Reboot(); err != nil {
				log.Fatalf("DHCPv6: kexec.Reboot failed: %v", err)
			}
		}
	}
	// DHCPv4
	if *useV4 {
		log.Printf("Trying to obtain a DHCPv4 lease on %s", *ifname)
		_, err := netboot.IfUp(*ifname, interfaceUpTimeout)
		if err != nil {
			log.Fatalf("DHCPv4: IfUp failed: %v", err)
		}
		debug("DHCPv4: interface %s is up", *ifname)
		if *skipDHCP {
			log.Print("Skipping DHCP")
		} else {
			log.Print("DHCPv4: sending request")
			client := dhcpv4.NewClient()
			// TODO add options to request to netboot
			conversation, err := client.Exchange(*ifname, nil)
			for _, m := range conversation {
				debug(m.Summary())
			}
			if err != nil {
				log.Fatalf("DHCPv4: Exchange failed: %v", err)
			}
			// TODO configure the network and DNS
			// TODO extract the next server and boot file and fetch it
			// TODO kexec into the NBP
		}
	}

}
