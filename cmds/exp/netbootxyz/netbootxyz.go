// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/boot/menu"
	"gopkg.in/yaml.v2"

	// To build the dependencies of this package with TinyGo, we need to include
	// the cpuid package, since tinygo does not support the asm code in the
	// cpuid package. The cpuid package will use the tinygo bridge to get the
	// CPU information. For further information see
	// github.com/u-root/cpuid/cpuid_amd64_tinygo_bridge.go
	_ "github.com/u-root/cpuid"
)

var (
	githubBaseURL = "https://github.com/netbootxyz"
	netbootxyzURL = "https://raw.githubusercontent.com/netbootxyz/netboot.xyz/development/endpoints.yml"

	verbose = flag.Bool("v", false, "Verbose output")
	noLoad  = flag.Bool("no-load", false, "Get DHCP response, print chosen boot configuration, but do not download + exec it")
	noExec  = flag.Bool("no-exec", false, "Download boot configuration, but do not exec it")
	ifName  = flag.String("i", "eth0", "Interface to send packets through")

	bootMenu []menu.Entry
	subMenu  []menu.Entry

	kernelList     []Endpoint
	filesystemList []Endpoint

	OSEndpoints []OSEndpoint

	majorOS = []string{
		"Ubuntu",
		"Pop",
		"Mint",
		"Debian",
		"Fedora",
		"Gentoo",
	}

	OSCommandline = map[string]string{
		"Ubuntu": "boot=casper netboot=url url=%s",
		"Mint":   "boot=casper netboot=url url=%s",
		"Pop":    "boot=casper netboot=url url=%s",
		"Debian": "boot=live fetch=%s",
		"Fedora": "root=live:%s ro rd.live.image rd.lvm=0 rd.luks=0 rd.md=0 rd.dm=0",
		"Gentoo": "root=/dev/ram0 init=/linuxrc loop=/image.squashfs looptype=squashfs cdroot=1 real_root=/ fetch=%s",
	}

	banner = `

  __________________________________
< Netbootxyz is even hotter nowadays >
  ----------------------------------
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||

`
)

const (
	dhcpTimeout = 5 * time.Second
	dhcpTries   = 3
)

// Endpoint - YAML Endpoint
type Endpoint struct {
	Name    string
	Path    string   `yaml:"path"`
	Os      string   `yaml:"os"`
	Version string   `yaml:"version"`
	Files   []string `yaml:"files"`
	Flavor  string   `yaml:"flavor"`
	Kernel  string   `yaml:"kernel"`
}

// Endpoints - Map for OS Endpoints
type Endpoints struct {
	Endpoints map[string]Endpoint `yaml:"endpoints"`
}

// OSEndpoint - Parsed version of Endpoint
type OSEndpoint struct {
	Name        string
	RawName     string
	Vmlinuz     string
	Initrd      string
	Filesystem  string
	Version     string
	Commandline string
	OS          string
	onlyLabel   bool
}

// Label - Menu Function Label
func (o OSEndpoint) Label() string {
	return o.Name
}

// Load - Load data into kexec
func (o OSEndpoint) Load() error {
	if o.onlyLabel {
		subMenu = nil
		if o.Name == "Other" {
			// Load all other OS's
			for _, value := range OSEndpoints {
				_, found := OSCommandline[value.OS]
				if !found {
					subMenu = append(subMenu, value)
				}
			}
		} else {
			// Now we load everything with this label
			for _, value := range OSEndpoints {
				if value.OS == o.Name {
					subMenu = append(subMenu, value)
				}
			}
		}
		menu.ShowMenuAndLoad(true, subMenu...)

		return nil
	}

	if *noLoad {
		fmt.Printf("Selected %s\n", o.Name)
		fmt.Printf("Commandline: %s\n", o.Commandline)
	} else {
		tmpPath := "/tmp/" + strings.ReplaceAll(o.Name, " ", "") + "/"
		err := os.Mkdir(tmpPath, 0o666)
		if err != nil {
			return err
		}
		fmt.Printf("Download to %s\n", tmpPath)

		err = downloadFile(tmpPath+"vmlinuz", o.Vmlinuz)
		if err != nil {
			return err
		}

		err = downloadFile(tmpPath+"initrd", o.Initrd)
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

		if *noExec {
			fmt.Println("Loaded Kernel and Initrd")
			fmt.Printf("With Kernel at %s\n", tmpPath+"vmlinuz")
			fmt.Printf("With Initrd at %s\n", tmpPath+"initrd")
			fmt.Printf("Commandline: %s\n", o.Commandline)
		} else {
			fmt.Println("Loading Kernel and Initrd into kexec")
			fmt.Printf("With Kernel at %s\n", tmpPath+"vmlinuz")
			fmt.Printf("With Initrd at %s\n", tmpPath+"initrd")
			fmt.Printf("Commandline: %s\n", o.Commandline)

			// Load Kernel and initrd
			if err = kexec.FileLoad(vmlinuz, initrd, o.Commandline); err != nil {
				return err
			}

			// Load KExec kernel and initrd - init cmdline
			return kexec.Reboot()
		}
		return err
	}
	return nil
}

// Exec - Execute new kernel
func (o OSEndpoint) Exec() error {
	return nil
}

// IsDefault - Default Configuration
func (o OSEndpoint) IsDefault() bool {
	return false
}

// Edit - Edit something
func (o OSEndpoint) Edit(func(cmdline string) string) {}

// indexOf - Returns index of an element in an array
func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 // not found.
}

// remove element from array
func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func main() {
	// Print Banner and parse arguments
	fmt.Print(banner)
	time.Sleep(2 * time.Second)
	flag.Parse()

	// Get an IP address via DHCP
	err := configureDHCPNetwork()
	if err != nil {
		fmt.Printf("Error while getting IP : %v\n", err)
	}

	// Set up HTTP client
	config := &tls.Config{InsecureSkipVerify: false}
	tr := &http.Transport{TLSClientConfig: config}
	client := &http.Client{Transport: tr}

	// Fetch NetBoot(dot)xyz endpoints
	req, err := http.NewRequest(http.MethodGet, netbootxyzURL, nil)
	if err != nil {
		fmt.Printf("New Request Error : %v\n", err)
	}
	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error : %v\n", err)
	}
	content, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Parse YAML
	var e Endpoints
	err = yaml.Unmarshal(content, &e)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Sort entries into either Kernel or Distros
	// File Systems could also contain a Kernel directly!
	tmp := make(map[string]struct{})
	for key, value := range e.Endpoints {
		value.Name = key
		if value.Os != "" {
			tmp[value.Os] = struct{}{}
		}
		if strings.Contains(key, "kernel") {
			// Endpoint contains kernel and initrd
			kernelList = append(kernelList, value)
		} else {
			// Endpoint contains filesystem
			filesystemList = append(filesystemList, value)
		}
	}

	// Store keys from tmp in OSEntriesInMenu
	OSEntriesInMenu := make([]string, 0, len(tmp))
	for k := range tmp {
		OSEntriesInMenu = append(OSEntriesInMenu, strings.Title(k))
	}

	for _, entry := range filesystemList {
		// Define menu entry and fill with data
		var OSEntry OSEndpoint
		OSEntry.RawName = entry.Name
		OSEntry.Name = strings.Title(entry.Os) + " " + strings.Title(entry.Version) + " " + strings.Title(entry.Flavor)
		OSEntry.OS = strings.Title(entry.Os)

		files := entry.Files
		// Set kernel and initrd if found in endpoint
		if entry.Name == entry.Kernel || entry.Kernel == "" {
			OSEntry.Vmlinuz = githubBaseURL + entry.Path + "vmlinuz"
			OSEntry.Initrd = githubBaseURL + entry.Path + "initrd"
		} else if entry.Kernel != "" {
			// Search for corresponding kernel in the kernel list
			for _, value := range kernelList {
				if value.Name == entry.Kernel {
					OSEntry.Vmlinuz = githubBaseURL + value.Path + "vmlinuz"
					OSEntry.Initrd = githubBaseURL + value.Path + "initrd"
					break
				}
			}
		}
		// Remove already saved kernel and initrd from "files"
		if indexOf("vmlinuz", files) != -1 {
			files = remove(files, indexOf("vmlinuz", files))
		}
		if indexOf("initrd", files) != -1 {
			files = remove(files, indexOf("initrd", files))
		}
		// Add filesystem entry
		if len(files) != 0 {
			OSEntry.Filesystem = githubBaseURL + entry.Path + files[0]
		}
		// Set specific cmdline if defined or resort to generic cmdline
		_, found := OSCommandline[OSEntry.OS]
		if found {
			OSEntry.Commandline = "earlyprintk=ttyS0,115200 console=ttyS0,115200 ip=dhcp initrd=initrd " +
				fmt.Sprintf(OSCommandline[OSEntry.OS], OSEntry.Filesystem)
		} else {
			OSEntry.Commandline = "earlyprintk=ttyS0,115200 console=ttyS0,115200 ip=dhcp initrd=initrd " +
				fmt.Sprintf("boot=casper netboot=url url=%s", OSEntry.Filesystem)
		}
		// Store each fully configured endpoint
		OSEndpoints = append(OSEndpoints, OSEntry)
	}
	// Only add major OS first
	for _, value := range OSEntriesInMenu {
		entry := OSEndpoint{
			Name:      value,
			onlyLabel: true,
		}
		for _, val := range majorOS {
			if val == value {
				bootMenu = append(bootMenu, entry)
			}
		}
	}
	// Group non-major distributions here
	entry := OSEndpoint{
		Name:      "Other",
		onlyLabel: true,
	}
	// Fill menu with remaining options
	bootMenu = append(bootMenu, entry)
	bootMenu = append(bootMenu, menu.Reboot{})
	bootMenu = append(bootMenu, menu.StartShell{})
	menu.SetInitialTimeout(90 * time.Second)
	menu.ShowMenuAndLoad(true, bootMenu...)
}
