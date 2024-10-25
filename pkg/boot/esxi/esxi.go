// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package esxi contains an ESXi boot config parser for disks and CDROMs.
//
// For CDROMs, it parses the boot.cfg found in the root directory and tries to
// boot from it.
//
// For disks, there may be multiple boot partitions:
//
// - Locates both <device>5/boot.cfg and <device>6/boot.cfg.
//
// - If parsable, chooses partition with bootstate=(0|2|empty) and greater
// updated=N.
//
// Sometimes, an ESXi partition can contain a valid boot.cfg, but not actually
// any of the named modules. Hence it is important to try fully loading ESXi
// into memory, and only then falling back to the other partition.
//
// Only boots partitions with bootstate=0, bootstate=2, bootstate=(empty) will
// boot at all.
//
// Most of the parsing logic in this package comes from
// https://github.com/vmware/esx-boot/blob/master/safeboot/bootbank.c
package esxi

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/sys/unix"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/multiboot"
	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/mount/gpt"
	"github.com/u-root/uio/uio"
)

func partNo(device string, number int) (string, error) {
	var name string
	if unicode.IsDigit(rune(device[len(device)-1])) {
		name = fmt.Sprintf("%sp%d", device, number)
	} else {
		name = fmt.Sprintf("%s%d", device, number)
	}
	if _, err := os.Stat(name); err != nil {
		return "", err
	}
	return name, nil
}

// LoadDisk loads the right ESXi multiboot kernel from partitions 5 or 6 of the
// given device.
//
// The kernels are returned in the priority order according to the bootstate
// and updated values in their boot configurations.
//
// The caller should try loading all returned images in order, as some of them
// may not be valid.
//
// device5 and device6 will be mounted at temporary directories.
func LoadDisk(device string) ([]*boot.MultibootImage, []*mount.MountPoint, error) {
	opts5, mp5, err5 := mountPartition(device, 5)
	opts6, mp6, err6 := mountPartition(device, 6)
	if err5 != nil && err6 != nil {
		return nil, nil, fmt.Errorf("could not mount or read either partition 5 (%w) or partition 6 (%w)", err5, err6)
	}
	var mps []*mount.MountPoint
	if mp5 != nil {
		mps = append(mps, mp5)
	}
	if mp6 != nil {
		mps = append(mps, mp6)
	}

	imgs, err := getImages(device, opts5, opts6)
	if err != nil {
		for _, mp := range mps {
			mp.Unmount(mount.MNT_DETACH)
		}
		return nil, nil, err
	}
	return imgs, mps, nil
}

func getImages(device string, opts5, opts6 *options) ([]*boot.MultibootImage, error) {
	var (
		img5, img6 *boot.MultibootImage
		err5, err6 error
	)
	if opts5 != nil {
		name, _ := partNo(device, 5)
		img5, err5 = getBootImage(*opts5, device, 5, name)
	}
	if opts6 != nil {
		name, _ := partNo(device, 6)
		img6, err6 = getBootImage(*opts6, device, 6, name)
	}
	if img5 == nil && img6 == nil {
		return nil, fmt.Errorf("could not read boot configs on partition 5 (%w) or partition 6 (%w)", err5, err6)
	}

	if img5 != nil && img6 != nil {
		if opts6.updated > opts5.updated {
			return []*boot.MultibootImage{img6, img5}, nil
		}
		return []*boot.MultibootImage{img5, img6}, nil
	} else if img5 != nil {
		return []*boot.MultibootImage{img5}, nil
	}
	return []*boot.MultibootImage{img6}, nil
}

// LoadCDROM loads an ESXi multiboot kernel from a CDROM at device.
//
// device will be mounted at mountPoint.
func LoadCDROM(device string) (*boot.MultibootImage, *mount.MountPoint, error) {
	mountPoint, err := os.MkdirTemp("", "esxi-mount-")
	if err != nil {
		return nil, nil, err
	}
	mp, err := mount.Mount(device, mountPoint, "iso9660", "", unix.MS_RDONLY|unix.MS_NOATIME)
	if err != nil {
		os.RemoveAll(mountPoint)
		return nil, nil, err
	}

	opts, err := parse(filepath.Join(mountPoint, "boot.cfg"))
	if err != nil {
		mp.Unmount(mount.MNT_DETACH)
		os.RemoveAll(mountPoint)
		return nil, nil, fmt.Errorf("cannot parse config from %s: %w", device, err)
	}
	img, err := getBootImage(opts, "", 0, device)
	if err != nil {
		mp.Unmount(mount.MNT_DETACH)
		os.RemoveAll(mountPoint)
		return nil, nil, err
	}
	return img, mp, nil
}

// LoadConfig loads an ESXi configuration from configFile.
func LoadConfig(configFile string) (*boot.MultibootImage, error) {
	opts, err := parse(configFile)
	if err != nil {
		return nil, fmt.Errorf("cannot parse config at %s: %w", configFile, err)
	}
	return getBootImage(opts, "", 0, fmt.Sprintf("config file %s", configFile))
}

func mountPartition(parentdev string, partition int) (*options, *mount.MountPoint, error) {
	dev, err := partNo(parentdev, partition)
	if err != nil {
		return nil, nil, err
	}
	base := filepath.Base(dev)
	mountPoint, err := os.MkdirTemp("", fmt.Sprintf("%s-", base))
	if err != nil {
		return nil, nil, err
	}
	mp, err := mount.Mount(dev, mountPoint, "vfat", "", unix.MS_RDONLY|unix.MS_NOATIME)
	if err != nil {
		os.RemoveAll(mountPoint)
		return nil, nil, err
	}

	configFile := filepath.Join(mountPoint, "boot.cfg")
	opts, err := parse(configFile)
	if err != nil {
		mp.Unmount(mount.MNT_DETACH)
		os.RemoveAll(mountPoint)
		return nil, nil, fmt.Errorf("cannot parse config at %s: %w", configFile, err)
	}
	return &opts, mp, nil
}

// lazyOpenModules assigns modules to be opened as files.
//
// Each module is a path followed by optional command-line arguments, e.g.
// []string{"./module arg1 arg2", "./module2 arg3 arg4"}.
func lazyOpenModules(mods []module) multiboot.Modules {
	modules := make([]multiboot.Module, 0, len(mods))
	for _, m := range mods {
		modules = append(modules, multiboot.Module{
			Cmdline: m.cmdline,
			Module:  uio.NewLazyFile(m.path),
		})
	}
	return modules
}

func getBootImage(opts options, device string, partition int, name string) (*boot.MultibootImage, error) {
	// Only valid and upgrading are bootable partitions.
	//
	// We are supposed to support the following two state transitions (only
	// one transition every boot!):
	//
	// upgrading -> dirty
	// dirty -> invalid
	//
	// A validly booted system will set its own bootstate to "valid" from
	// "dirty".
	//
	// We currently don't support writing the state back to disk, which is
	// fine in our manual testing.
	if opts.bootstate != bootValid && opts.bootstate != bootUpgrading {
		return nil, fmt.Errorf("boot state %d invalid", opts.bootstate)
	}

	if len(device) > 0 {
		if err := opts.addUUID(device, partition); err != nil {
			return nil, fmt.Errorf("cannot add boot uuid of %s: %w", device, err)
		}
	}

	return &boot.MultibootImage{
		Name:    fmt.Sprintf("%s from %s", opts.title, name),
		Kernel:  uio.NewLazyFile(opts.kernel),
		Cmdline: opts.args,
		Modules: lazyOpenModules(opts.modules),
	}, nil
}

type module struct {
	path    string
	cmdline string
}

type options struct {
	title     string
	kernel    string
	args      string
	modules   []module
	updated   int
	bootstate bootstate
}

type bootstate int

// From safeboot.c
const (
	bootValid     bootstate = 0
	bootUpgrading bootstate = 1
	bootDirty     bootstate = 2
	bootInvalid   bootstate = 3
)

// So tests can replace this and don't have to have actual block devices.
var getBlockSize = gpt.GetBlockSize

func getUUID(device string, partition int) (string, error) {
	device = strings.TrimRight(device, "/")
	blockSize, err := getBlockSize(device)
	if err != nil {
		return "", err
	}

	dev, err := partNo(device, partition)
	if err != nil {
		return "", err
	}

	f, err := os.Open(dev)
	if err != nil {
		return "", err
	}

	// Boot uuid is stored in the second block of the disk
	// in the following format:
	//
	// VMWARE FAT16    <uuid>
	// <---128 bit----><128 bit>
	data := make([]byte, uuidSize)
	n, err := f.ReadAt(data, int64(blockSize))
	if err != nil {
		return "", err
	}
	if n != uuidSize {
		return "", io.ErrUnexpectedEOF
	}

	if magic := string(data[:len(uuidMagic)]); magic != uuidMagic {
		return "", fmt.Errorf("bad uuid magic %q, want %q", magic, uuidMagic)
	}

	uuid := hex.EncodeToString(data[len(uuidMagic):])
	return fmt.Sprintf("bootUUID=%s", uuid), nil
}

func (o *options) addUUID(device string, partition int) error {
	uuid, err := getUUID(device, partition)
	if err != nil {
		return err
	}
	o.args += " " + uuid
	return nil
}

const (
	comment = '#'
	sep     = "---"

	uuidMagic = "VMWARE FAT16    "
	uuidSize  = 32
)

func parse(configFile string) (options, error) {
	dir := filepath.Dir(configFile)

	f, err := os.Open(configFile)
	if err != nil {
		return options{}, err
	}
	defer f.Close()

	// An empty or missing updated value is always 0, so we can let the
	// ints be initialized to 0.
	//
	// see esx-boot/bootlib/parse.c:parse_config_file.
	opt := options{
		title: "VMware ESXi",
		// Default value taken from
		// esx-boot/safeboot/bootbank.c:bank_scan.
		bootstate: bootInvalid,
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 || line[0] == comment {
			continue
		}

		tokens := strings.SplitN(line, "=", 2)
		if len(tokens) != 2 {
			return opt, fmt.Errorf("bad line %q", line)
		}
		key := strings.TrimSpace(tokens[0])
		val := strings.TrimSpace(tokens[1])
		switch key {
		case "title":
			opt.title = val

		case "kernel":
			opt.kernel = filepath.Join(dir, val)

			// The kernel cmdline is expected to have the filename
			// first, as in cmdlines[0] here:
			// https://github.com/vmware/esx-boot/blob/1380fc86cffdfb83448e2913ae11f6b7f248cf23/mboot/mutiboot.c#L870
			//
			// Note that the kernel is module 0 in the esx-boot
			// code base, but it doesn't get loaded like that into
			// the info structure; see -- so don't panic like I did
			// when you read that!
			// https://github.com/vmware/esx-boot/blob/1380fc86cffdfb83448e2913ae11f6b7f248cf23/mboot/mutiboot.c#L578
			opt.args = val + " " + opt.args

		case "kernelopt":
			opt.args += val

		case "updated":
			if len(val) == 0 {
				// Explicitly setting to 0, as in
				// esx-boot/bootlib/parse.c:parse_config_file,
				// in case this value is specified twice.
				opt.updated = 0
			} else {
				n, err := strconv.Atoi(val)
				if err != nil {
					return options{}, err
				}
				opt.updated = n
			}
		case "bootstate":
			if len(val) == 0 {
				// Explicitly setting to valid, as in
				// esx-boot/bootlib/parse.c:parse_config_file,
				// in case this value is specified twice.
				opt.bootstate = bootValid
			} else {
				n, err := strconv.Atoi(val)
				if err != nil {
					return options{}, err
				}
				if n < 0 || n > 3 {
					opt.bootstate = bootInvalid
				} else {
					opt.bootstate = bootstate(n)
				}
			}
		case "modules":
			for _, tok := range strings.Split(val, sep) {
				// Each module is "filename arg0 arg1 arg2" and
				// the filename is relative to the directory
				// the module is in.
				tok = strings.TrimSpace(tok)
				if len(tok) > 0 {
					entry := strings.Fields(tok)
					opt.modules = append(opt.modules, module{
						path:    filepath.Join(dir, entry[0]),
						cmdline: tok,
					})
				}
			}
		}
	}

	err = scanner.Err()
	return opt, err
}
