package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/boot/menu"
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
var subMenu []menu.Entry

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

var majorOS = []string{
	"Ubuntu",
	"Fedora",
	"Pop",
	"Debian",
	"Manjaro",
	"Arch",
	"Mint"}

var OSCommandline = map[string]string{
	"Ubuntu":  "ip=dhcp boot=casper initrd=initrd netboot=url url=%s",
	"Fedora":  "root=live:%s ip=dhcp ro rd.live.image rd.lvm=0 rd.luks=0 rd.md=0 rd.dm=0 initrd=initrd",
	"Pop":     "ip=dhcp boot=casper netboot=url url=%s",
	"Debian":  "auto=true priority=critical preseed/url=%s initrd=initrd.gz",
	"Manjaro": "ip=dhcp net.ifnames=0 miso_http_srv=%s nouveau.modeset=1 i915.modeset=1 radeon.modeset=1 driver=free tz=UTC lang=en_US keytable=us",
	"Arch":    "ip=dhcp archiso_http_srv=%s archisobasedir=arch verify=y net.ifnames=0",
	"Mint":    "ip=dhcp boot=casper netboot=http fetch=%s",
}

// Label - Menu Function Label
func (o OSEndpoint) Label() string {
	return o.Name
}

// Load - Load data into kexec
func (o OSEndpoint) Load() error {
	if o.onlyLabel == true {
		subMenu = nil
		if o.Name == "Other" {
			// Load all other OS's
		} else {
			// Now we load everything with this label
			for _, value := range OSEndpoints {
				if value.OS == o.Name {
					subMenu = append(subMenu, value)
				}
			}
		}
		menu.ShowMenuAndLoad(os.Stdin, subMenu...)

		return nil
	}

	tmpPath := "/tmp/" + strings.ReplaceAll(o.Name, " ", "") + "/"
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
	fmt.Printf("With Kernel at %s\n", tmpPath+"vmlinuz")
	fmt.Printf("With Initrd at %s\n", tmpPath+"initrd")
	fmt.Printf("Commandline: %s\n", o.Commandline)
	// Load Kernel and initrd
	err = kexec.FileLoad(vmlinuz, initrd, o.Commandline)
	// Load KExec kernel and initrd - init cmdline

	cmd := exec.Command("kexec", "-e")

	err = cmd.Run()
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

// Edit - edit something
func (o OSEndpoint) Edit(func(cmdline string) string) {
	return
}

// Endpoints - map for OS Endpoints
type Endpoints struct {
	Endpoints map[string]Endpoint `yaml:"endpoints"`
}

var kernelList []Endpoint
var filesystemList []Endpoint

var OSEndpoints []OSEndpoint

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

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

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

	var e Endpoints
	if err := yaml.Unmarshal(content, &e); err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(e)

	tmp := make(map[string]struct{})

	// Sort entries into either Kernel or Distros
	// File Systems could also contain a Kernel directly!

	for key, value := range e.Endpoints {
		value.Name = key
		if value.Os != "" {
			tmp[value.Os] = struct{}{}
		}
		if strings.Contains(key, "kernel") {
			kernelList = append(kernelList, value)
		} else {
			filesystemList = append(filesystemList, value)
		}
	}

	OSEntriesInMenu := make([]string, 0, len(tmp))
	for k := range tmp {
		OSEntriesInMenu = append(OSEntriesInMenu, strings.Title(k))
	}

	for _, entry := range filesystemList {
		var OSEntry OSEndpoint

		OSEntry.RawName = entry.Name
		OSEntry.Name = strings.Title(entry.Os) + " " + strings.Title(entry.Version) + " " + strings.Title(entry.Flavor)
		OSEntry.OS = strings.Title(entry.Os)

		files := entry.Files

		if entry.Name == entry.Kernel || entry.Kernel == "" {
			OSEntry.Vmlinuz = githubBaseURL + entry.Path + "vmlinuz"
			OSEntry.Initrd = githubBaseURL + entry.Path + "initrd"
		} else if entry.Kernel != "" {
			// Search for Kernel in Kernel List
			for _, value := range kernelList {
				if value.Name == entry.Kernel {
					OSEntry.Vmlinuz = githubBaseURL + value.Path + "vmlinuz"
					OSEntry.Initrd = githubBaseURL + value.Path + "initrd"
				}
			}
		}
		if indexOf("vmlinuz", files) != -1 {
			files = remove(files, indexOf("vmlinuz", files))
		}
		if indexOf("initrd", files) != -1 {
			files = remove(files, indexOf("initrd", files))
		}
		if len(files) != 0 {
			OSEntry.Filesystem = githubBaseURL + entry.Path + files[0]
		}
		OSEntry.Commandline = "earlyprintk=ttyS0,115200 console=ttyS0,115200 " + fmt.Sprintf(OSCommandline[OSEntry.OS], OSEntry.Filesystem)

		OSEndpoints = append(OSEndpoints, OSEntry)
	}

	// Only Add Major OS first
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
	entry := OSEndpoint{
		Name:      "Other",
		onlyLabel: true,
	}
	bootMenu = append(bootMenu, entry)
	bootMenu = append(bootMenu, menu.Reboot{})
	bootMenu = append(bootMenu, menu.StartShell{})

	menu.ShowMenuAndLoad(os.Stdin, bootMenu...)
}
