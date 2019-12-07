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

// HostVars contains contains platform-specific data
type HostVars struct {
	HostIP         string `json:"host_ip"`
	HostNetmask    string `json:"netmask"`
	DefaultGateway string `json:"gateway"`
	DNSServer      string `json:"dns"`

	BootstrapURL string `json:"bootstrap_url"`

	MinimalSignaturesMatch int `json:"minimal_signatures_match"`
}

// FindHostVarsOnPartition mounts all possible devices with every possible
// file system and looks for hostvars.json at root of partition
func FindHostVarsOnPartition() (HostVars, error) {
	var vars HostVars
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
			fmt.Printf("Failed to mount %s on %s: %v", devname, mountpath, err)
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
	var file string
	for _, mountpoint := range mounted {
		file = path.Join(mountpoint.Path, HostVarsName)
		log.Printf("Trying to read %s", file)
		data, err = ioutil.ReadFile(file)
		if err == nil {
			break
		}
		log.Printf("cannot open %s: %v", file, err)
	}

	if err = json.Unmarshal(data, &vars); err != nil {
		return vars, fmt.Errorf("unable to parse %s: %v", file, err)
	}

	return vars, nil
}

// FindHostVarsInInitramfs looks for netvars.json at a given path inside
// the initramfs file system. The hostvars.json is
// expected to be in /etc.
func FindHostVarsInInitramfs() (HostVars, error) {
	var vars HostVars
	file := path.Join("etc/", HostVarsName)
	if _, err := os.Stat(file); os.IsNotExist(err) != false {
		return vars, fmt.Errorf("%s not found: %v", file, err)
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return vars, fmt.Errorf("cant open %s: %v", file, err)
	}
	if err = json.Unmarshal(data, &vars); err != nil {
		return vars, fmt.Errorf("cant parse data from %s", file)
	}
	return vars, nil
}
