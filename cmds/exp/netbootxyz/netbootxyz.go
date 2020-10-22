package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/boot/menu"
	"github.com/u-root/u-root/pkg/dhclient"
	"github.com/vishvananda/netlink"
	"gopkg.in/yaml.v2"
)

var githubBaseURL = "https://github.com/netbootxyz"

var (
	ifName        = "^e.*"
	verbose       = flag.Bool("v", false, "Verbose output")
	netbootxyzURL = "https://raw.githubusercontent.com/netbootxyz/netboot.xyz/development/endpoints.yml"
	noLoad        = flag.Bool("no-load", false, "get DHCP response, print chosen boot configuration, but do not download + exec it")
	noExec        = flag.Bool("no-exec", false, "download boot configuration, but do not exec it")
	ifname        = flag.String("i", "eth0", "Interface to send packets through")
)

var bootMenu []menu.Entry

const (
	dhcpTimeout = 5 * time.Second
	dhcpTries   = 3
)

// Endpoint - YAML Endpoint
type Endpoint struct {
	Path    string   `yaml:"path"`
	Os      string   `yaml:"os"`
	Version string   `yaml:"version"`
	Files   []string `yaml:"files"`
	Flavor  string   `yaml:"flavor"`
	Kernel  string   `yaml:"kernel"`
}

type OSEndpoint struct {
	Name        string
	Vmlinuz     string
	Initrd      string
	Filesystem  string
	Version     string
	Commandline string
	OS          string
}

// Label - Menu Function Label
func (o OSEndpoint) Label() string {
	return o.Name
}

// Load - Load data into kexec
func (o OSEndpoint) Load() error {
	tmpPath := "/tmp/" + o.Name + "/"
	err := os.Mkdir(tmpPath, 0666)
	if err != nil {
		return err
	}
	fmt.Printf("Download to %s\n", tmpPath)

	err = DownloadFile(tmpPath+"vmlinuz", o.Vmlinuz)
	if err != nil {
		return err
	}

	err = DownloadFile(tmpPath+"initrd", o.Initrd)
	if err != nil {
		return err
	}

	vmlinuz, err := os.Open(tmpPath + "vmlinuz")
	if err != nil {
		return err
	}

	initrd, err := os.Open(tmpPath + "initrd")
	if err != nil {
		return err
	}

	fmt.Println("Loading Kernel and Initrd into kexec")
	fmt.Printf("With Kernel at %s\n", tmpPath+"vmliuz")
	fmt.Printf("With Initrd at %s\n", tmpPath+"initrd")
	fmt.Printf("Commandline: %s\n", o.Commandline)
	// Load Kernel and initrd
	err = kexec.FileLoad(vmlinuz, initrd, o.Commandline)
	// Load KExec kernel and initrd - init cmdline
	return err
}

// Exec - execute new kernel
func (o OSEndpoint) Exec() error {
	// Execute
	return nil
}

// IsDefault - Default Configuration
func (o OSEndpoint) IsDefault() bool {
	return false
}

// Endpoints - map for OS Endpoints
type Endpoints struct {
	Endpoints map[string]Endpoint
}

var OSEndpoints []OSEndpoint

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

// UnmarshalYAML - Custom unmarshal YAML function
func (e *Endpoints) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var endpoints map[string]Endpoint
	if err := unmarshal(&endpoints); err != nil {
		// Here we expect an error because a boolean cannot be converted to a
		// a MajorVersion
		if _, ok := err.(*yaml.TypeError); !ok {
			return err
		}
	}
	e.Endpoints = endpoints
	return nil
}

var banner = `

 _________________________________
< Netbootxyz is even hoter nowadays >
 ---------------------------------
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||

`

func main() {

	fmt.Print(banner)
	time.Sleep(2 * time.Second)

	flag.Parse()

	configureDHCPNetwork()

	config := &tls.Config{
		InsecureSkipVerify: true,
	}

	tr := &http.Transport{TLSClientConfig: config}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(http.MethodGet, netbootxyzURL, nil)
	if err != nil {
		fmt.Printf("New Request Error : %v\n", err)
	}
	response, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error : %v\n", err)
	}

	content, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("YAML: %v \n", string(content))

	var e map[string]Endpoints
	if err := yaml.Unmarshal(content, &e); err != nil {
		fmt.Println(err.Error())
	}

	for key, value := range e["endpoints"].Endpoints {
		if !strings.Contains(key, "ubuntu") {
			continue
		}
		fmt.Printf("OS: %s\n", key)
		fmt.Printf("Kernel: %s\n", value.Kernel)
		fmt.Printf("Files: %v\n", value.Files)

		var vmlinuz, initrd, filesystem string

		// If kernel is empty - or we do point to ourselves
		if (value.Kernel == "") || (value.Kernel == key) {
			vmlinuz = githubBaseURL + value.Path + "vmlinuz"
			initrd = githubBaseURL + value.Path + "initrd"
			// TODO: Check for real FS
			filesystem = githubBaseURL + value.Path + "filesystem.squashfs"
		} else {
			vmlinuz = value.Kernel
			initrd = value.Kernel
		}

		filesystem = githubBaseURL + "/ubuntu-squash/releases/download/18.04.5-0dd1e29b/filesystem.squashfs"

		OSEndpoint := OSEndpoint{
			Name:        key,
			Vmlinuz:     vmlinuz,
			Initrd:      initrd,
			Filesystem:  filesystem,
			Version:     value.Version,
			Commandline: "earlyprintk=ttyS0,115200 console=ttyS0,115200 ip=dhcp boot=casper initrd=initrd netboot=url url=" + filesystem, // Function which gets this
			OS:          value.Os,
		}

		OSEndpoints = append(OSEndpoints, OSEndpoint)
		bootMenu = append(bootMenu, OSEndpoint)

		fmt.Printf("\n%v\n", OSEndpoint)
	}

	fmt.Printf("\n\n%v\n", OSEndpoints)

	bootMenu = append(bootMenu, menu.Reboot{})
	bootMenu = append(bootMenu, menu.StartShell{})

	// Boot does not return.
	menu.ShowMenuAndLoad(os.Stdin, bootMenu...)
}
