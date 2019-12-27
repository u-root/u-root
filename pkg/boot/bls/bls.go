// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bls parses systemd Boot Loader Spec config files.
//
// See spec at https://systemd.io/BOOT_LOADER_SPECIFICATION. Only Type #1 BLS
// entries are supported at the moment, while Type #2 EFI entries are left
// unimplemented awaiting EFI boot support in u-root/LinuxBoot.
package bls

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/ulog"
)

const (
	blsEntriesDir = "loader/entries"
)

// ScanBLSEntries scans the filesystem root for valid BLS entries.
// This function skips over invalid or unreadable entries in an effort
// to return everything that is bootable.
func ScanBLSEntries(log ulog.Logger, fsRoot string) ([]boot.OSImage, error) {
	entriesDir := filepath.Join(fsRoot, blsEntriesDir)

	files, err := filepath.Glob(filepath.Join(entriesDir, "*.conf"))
	if err != nil {
		return nil, fmt.Errorf("no BootLoaderSpec entries found: %w", err)
	}

	// TODO: Rank entries by version or machine-id attribute as suggested
	// in the spec (but not mandated, surprisingly).
	var imgs []boot.OSImage
	for _, f := range files {
		entry, err := parseBLSEntry(f, entriesDir)
		if err != nil {
			log.Printf("BootLoaderSpec skipping entry %s: %v", f, err)
			continue
		}
		imgs = append(imgs, entry)
	}
	return imgs, nil
}

type entry struct {
	dir   string
	name  string
	vals  map[string]string
	image boot.OSImage
}

func parseEntry(entryPath string) (*entry, error) {
	f, err := os.Open(entryPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dir, name := filepath.Split(entryPath)
	e := &entry{
		dir:  dir,
		name: name,
		vals: make(map[string]string),
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimSpace(line)

		sline := strings.SplitN(line, " ", 2)
		if len(sline) != 2 {
			continue
		}
		e.vals[sline[0]] = strings.TrimSpace(sline[1])
	}
	return e, nil
}

func parseLinuxImage(e *entry, baseDir string) (boot.OSImage, error) {
	linux := &boot.LinuxImage{}

	var cmdlines []string
	for key, val := range e.vals {
		switch key {
		case "linux":
			f, err := os.Open(filepath.Join(baseDir, val))
			if err != nil {
				return nil, err
			}
			linux.Kernel = f

		// TODO: initrd may be specified more than once.
		case "initrd":
			f, err := os.Open(filepath.Join(baseDir, val))
			if err != nil {
				return nil, err
			}
			linux.Initrd = f

		case "devicetree":
			// Explicitly return an error rather than ignore this,
			// because the intended kernel likely won't boot
			// correctly if we silently ignore this attribute.
			return nil, fmt.Errorf("devicetree attribute unsupported for Linux entries")

		// options may appear more than once.
		case "options":
			cmdlines = append(cmdlines, val)
		}
	}

	// Spec says kernel is required.
	if linux.Kernel == nil {
		return nil, fmt.Errorf("malformed Linux config: linux keyword missing")
	}

	var name []string
	if title, ok := e.vals["title"]; ok && len(title) > 0 {
		name = append(name, title)
	}
	if version, ok := e.vals["version"]; ok && len(version) > 0 {
		name = append(name, version)
	}
	// If both title and version were empty, so will this.
	linux.Name = strings.Join(name, " ")
	linux.Cmdline = strings.Join(cmdlines, " ")
	return linux, nil
}

// parseBLSEntry takes a Type #1 BLS entry and the directory of entries, and
// returns a LinuxImage.
// An error is returned if the syntax is wrong or required keys are missing.
func parseBLSEntry(entryPath, entriesDir string) (boot.OSImage, error) {
	baseDir := filepath.Dir(entriesDir)

	e, err := parseEntry(entryPath)
	if err != nil {
		return nil, fmt.Errorf("error parsing config in %s: %w", entryPath, err)
	}

	var img boot.OSImage
	err = fmt.Errorf("neither linux, efi, nor multiboot present in BootLoaderSpec config")
	if _, ok := e.vals["linux"]; ok {
		img, err = parseLinuxImage(e, baseDir)
	} else if _, ok := e.vals["multiboot"]; ok {
		err = fmt.Errorf("multiboot not yet supported")
	} else if _, ok := e.vals["efi"]; ok {
		err = fmt.Errorf("EFI not yet supported")
	}
	if err != nil {
		return nil, fmt.Errorf("error parsing config in %s: %w", entryPath, err)
	}
	return img, nil
}
