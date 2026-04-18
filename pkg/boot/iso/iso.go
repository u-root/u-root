// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package iso helps to boot directly from iso files on removable drives.
// Works best with "Live" iso files.
package iso

import (
	"context"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/grub"
	"github.com/u-root/u-root/pkg/curl"
	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/mount/block"
	"github.com/u-root/u-root/pkg/mount/loop"
	"github.com/u-root/u-root/pkg/ulog"

	"golang.org/x/sys/unix"
)

var (
	probeLoopbackFiles = []string{
		"boot/grub/loopback.cfg",
		"boot/grub2/loopback.cfg",
	}

	// Allow mock
	blockPath = "/sys/class/block"
)

func mountISOFile(path string, mountPool *mount.Pool) (string, error) {
	dir, err := os.MkdirTemp("", "iso-mount-")
	if err != nil {
		return "", err
	}

	lo, err := loop.New(path, "iso9660", "")
	if err != nil {
		unix.Rmdir(dir)
		return "", err
	}

	mp, err := lo.Mount(dir, mount.ReadOnly)
	if err != nil {
		lo.Free()
		unix.Rmdir(dir)
		return "", err
	}
	mountPool.Add(mp)

	// This won't actually free the loop device once it's mounted.
	// Instead it will result in setting AUTOCLEAR=1, and the loop
	// devices will remove itself once unmounted.
	lo.Free()
	return dir, nil
}

// Allow mock
var mountISO = mountISOFile

func isRemovable(b *block.BlockDev) (bool, error) {
	// resolve full path, like /sys/device/.../sda/sda2
	p, err := filepath.EvalSymlinks(filepath.Join(blockPath, b.Name))
	if err != nil {
		return false, err
	}
	// if partition, parent is one above
	if _, err := os.Stat(filepath.Join(p, "partition")); err == nil {
		p = filepath.Join(p, "..")
	}
	data, err := os.ReadFile(filepath.Join(p, "removable"))
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(data)) == "1", nil
}

func parseImg(img boot.OSImage, isoloc, fsuuid string) (boot.OSImage, bool) {
	if li, ok := img.(*boot.LinuxImage); ok {
		// grubMapper expands vars set by grub parser
		// NOTE grub parser does not process any logic,
		// and just uses the last var "set"
		grubMapper := func(k string) string {
			if v, ok := li.Env[k]; ok {
				return v
			}
			return fmt.Sprintf("${%s}", k)
		}

		// manualMapper expands known vars used for loopback
		imgParsed := false
		manualMapper := func(s string) string {
			switch s {
			case "iso_path": // happy path
				imgParsed = true
				return fmt.Sprintf("/%s", isoloc)

			// Arch Linux
			case "archiso_img_dev_uuid":
				return fsuuid
			case "archiso_platform": // used in li.Name
				return "BIOS"
			}
			return fmt.Sprintf("${%s}", s)
		}

		li.Name = fmt.Sprintf("[%s] %s", isoloc, os.Expand(li.Name, manualMapper))
		li.Cmdline = os.Expand(li.Cmdline, grubMapper)
		li.Cmdline = os.Expand(li.Cmdline, manualMapper)
		if imgParsed {
			return img, true
		}
		// Other tests for ISO, where var sub is not used OR var not defined using "set"
		if strings.Contains(li.Cmdline, "root=live:CDLABEL=") || // Fedora, Centos, etc.
			strings.Contains(li.Cmdline, "inst.stage2=hd:LABEL=") ||
			strings.Contains(li.Cmdline, "boot=casper") { // PopOS
			li.Cmdline += " " + fmt.Sprintf("iso-scan/filename=/%s", isoloc)
			return img, true
		}
	}
	return nil, false
}

// ParseISOFiles scans a device for .iso files, mounts and adds them to the Pool.
// Then parses any loopback.cfg or grub.cfg files for boot entries, while
// substituting/adding cmdline vars to boot directly from the iso.
// Skips entries it can't process.
func ParseISOFiles(l ulog.Logger, mountDir string, dev *block.BlockDev, mountPool *mount.Pool) ([]boot.OSImage, error) {
	var images []boot.OSImage

	removable, err := isRemovable(dev)
	if err != nil {
		return nil, err
	}
	if !removable {
		l.Printf("[iso] skip non removable device: %s", dev)
		return nil, nil
	}

	err = filepath.WalkDir(mountDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				return nil
			}
			return err
		}
		if d.IsDir() || strings.ToLower(filepath.Ext(path)) != ".iso" {
			return nil
		}

		dir, err := mountISO(path, mountPool)
		if err != nil {
			l.Printf("[iso] could not mount %s, err: %v", path, err)
			return nil
		}

		isoloc, err := filepath.Rel(mountDir, path)
		if err != nil {
			return nil
		}

		root := &url.URL{
			Scheme: "file",
			Path:   dir,
		}
		// first try loopback.cfg
		var isoImgs []boot.OSImage
		for _, cfgfile := range probeLoopbackFiles {
			isoImgs, err = grub.ParseConfigFile(context.Background(), curl.DefaultSchemes, cfgfile, root, nil, nil)
			if err == nil && len(isoImgs) > 0 {
				break
			}
		}
		// failing that, try the usual grub method
		if len(isoImgs) == 0 {
			isoImgs, err = grub.ParseLocalConfig(context.Background(), dir, nil, nil)
			if err != nil {
				return nil
			}
		}

		imgFound := false
		for _, img := range isoImgs {
			if i, ok := parseImg(img, isoloc, dev.FsUUID); ok {
				imgFound = true
				images = append(images, i)
			}
		}
		if !imgFound {
			l.Printf("[iso] could not find boot entry for %s", filepath.Base(path))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return images, nil
}
