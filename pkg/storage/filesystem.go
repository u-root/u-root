package storage

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
)

// Mountpoint holds mount point information for a given device
type Mountpoint struct {
	DeviceName string
	Path       string
	FsType     string
}

// GetSupportedFilesystems returns the supported file systems for block devices,
func GetSupportedFilesystems() ([]string, error) {
	fd, err := os.Open("/proc/filesystems")
	if err != nil {
		return nil, err
	}
	filesystems := make([]string, 0)
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), "\t")
		if fields[0] == "" {
			filesystems = append(filesystems, fields[1])
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return filesystems, nil
}

// Mount tries to mount a block device on the given mountpoint, trying in order
// the provided file system types. It returns a Mountpoint structure, or an error
// if the device could not be mounted. If the mount point does not exist, it will
// be created.
func Mount(devname, mountpath string, filesystems []string) (*Mountpoint, error) {
	if err := os.MkdirAll(mountpath, 0744); err != nil {
		return nil, err
	}
	for _, fstype := range filesystems {
		log.Printf(" * trying %s on %s", fstype, devname)
		// MS_RDONLY should be enough. See mount(2)
		flags := uintptr(syscall.MS_RDONLY)
		// no options
		data := ""
		if err := syscall.Mount(devname, mountpath, fstype, flags, data); err != nil {
			log.Printf("    failed with %v", err)
			continue
		}
		log.Printf(" * mounted %s on %s with filesystem type %s", devname, mountpath, fstype)
		return &Mountpoint{DeviceName: devname, Path: mountpath, FsType: fstype}, nil
	}
	return nil, fmt.Errorf("no suitable filesystem type found to mount %s", devname)
}
