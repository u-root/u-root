package stboot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"syscall"

	"github.com/u-root/u-root/pkg/storage"
)

type NetVars struct {
	HostIP         string `json:"host_ip"`
	HostNetmask    string `json:"netmask"`
	DefaultGateway string `json:"gateway"`
	DNSServer      string `json:"dns"`

	HostPrivKey string `json:"host_priv_key"`
	HostPupKey  string `json:"host_pub_key"`

	BootstrapURL    string `json:"bootstrap_url"`
	SignaturePubKey string `json:"signature_pub_key"`

	MinimalAmountSignatures int `json:"minimal-amount-signatures"`
}

const (
	netVarsPath = "netvars.json"
)

// FindNetVars mounts all possible devices with every possible file system and looks for netvars.json
func FindNetVars() (NetVars, error) {
	var vars NetVars
	vars, err := findNetVarsInInitramfs(netVarsPath)
	if err == nil {
		return vars, nil
	}
	devices, err := storage.GetBlockStats()
	if err != nil {
		log.Fatal(err)
	}

	filesystems, err := storage.GetSupportedFilesystems()
	if err != nil {
		log.Fatal(err)
	}
	var mounted []storage.Mountpoint

	mounted = make([]storage.Mountpoint, 0)
	for _, dev := range devices {
		devname := path.Join("/dev", dev.Name)
		mountpath := path.Join("/mnt", dev.Name)
		if mountpoint, err := storage.Mount(devname, mountpath, filesystems); err != nil {
			debug("Failed to mount %s on %s: %v", devname, mountpath, err)
		} else {
			mounted = append(mounted, *mountpoint)
		}
	}
	defer func() {
		// clean up
		for _, mountpoint := range mounted {
			syscall.Unmount(mountpoint.Path, syscall.MNT_DETACH)
		}
	}()

	var data []byte
	for _, mountpoint := range mounted {
		path := path.Join(mountpoint.Path, "netvars.json")
		log.Printf("Trying to read %s", path)
		data, err = ioutil.ReadFile(path)
		if err == nil {
			break
		}
		log.Printf("cannot open %s: %v", path, err)
	}

	if err = json.Unmarshal(data, &vars); err != nil {
		return vars, fmt.Errorf("unable to get netvars.json at all: %v", err)
	}

	return vars, nil
}

// FindNetVarsInInitramfs looks for netvars.json in a given path in the root file system
func findNetVarsInInitramfs(path string) (NetVars, error) {
	var vars NetVars
	if _, err := os.Stat(path); os.IsNotExist(err) != false {
		return vars, fmt.Errorf("Path not found: %v", err)
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return vars, fmt.Errorf("cant open netvars.json in path: %v", err)
	}
	if err = json.Unmarshal(file, &vars); err != nil {
		return vars, fmt.Errorf("cant parse data from netvars.json")
	}
	return vars, nil
}
