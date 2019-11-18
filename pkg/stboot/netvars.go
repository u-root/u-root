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

	BootstrapURL    string `json:"bootstrap_url"`
	SignaturePubKey string `json:"signature_pub_key"`

	MinimalAmountSignatures int `json:"minimal-amount-signatures"`
}

const (
	netVarsPath = "netvars.json"
)

// FindNetVarsOnPartition mounts all possible devices with every possible
// file system and looks for netvars.json partition root
func FindNetVarsOnPartition() (NetVars, error) {
	var vars NetVars
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
		return vars, fmt.Errorf("unable to parse netvars.json: %v", err)
	}

	return vars, nil
}

// FindNetVarsInInitramfs looks for netvars.json at a given path inside
// the initramfs file system. The netvars.json is
// expected to be at the root of the file system.
func FindNetVarsInInitramfs() (NetVars, error) {
	var vars NetVars
	if _, err := os.Stat("/netvars.json"); os.IsNotExist(err) != false {
		return vars, fmt.Errorf("netvars.json not found: %v", err)
	}
	file, err := ioutil.ReadFile("/netvars.json")
	if err != nil {
		return vars, fmt.Errorf("cant open netvars.json: %v", err)
	}
	if err = json.Unmarshal(file, &vars); err != nil {
		return vars, fmt.Errorf("cant parse data from netvars.json")
	}
	return vars, nil
}
