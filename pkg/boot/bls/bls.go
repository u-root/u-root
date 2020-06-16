// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bls parses systemd Boot Loader Spec config files.
//
// See spec at https://systemd.io/BOOT_LOADER_SPECIFICATION. Only Type #1 BLS
// entries are supported at the moment, while Type #2 EFI entries are left
// unimplemented awaiting EFI boot support in u-root/LinuxBoot.
//
// This package also supports the systemd-boot loader.conf as described in
// https://www.freedesktop.org/software/systemd/man/loader.conf.html. Only the
// "default" keyword is implemented.
package bls

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/ulog"
)

const (
	blsEntriesDir = "loader/entries"
)

func cutConf(s string) string {
	if strings.HasSuffix(s, ".conf") {
		return s[:len(s)-6]
	}
	return s
}

// ScanBLSEntries scans the filesystem root for valid BLS entries.
// This function skips over invalid or unreadable entries in an effort
// to return everything that is bootable.
func ScanBLSEntries(log ulog.Logger, fsRoot string) ([]boot.OSImage, error) {
	entriesDir := filepath.Join(fsRoot, blsEntriesDir)

	files, err := filepath.Glob(filepath.Join(entriesDir, "*.conf"))
	if err != nil {
		return nil, fmt.Errorf("no BootLoaderSpec entries found: %w", err)
	}

	// loader.conf is not in the real spec; it's an implementation detail
	// of systemd-boot. It is specified in
	// https://www.freedesktop.org/software/systemd/man/loader.conf.html
	loaderConf, err := parseConf(filepath.Join(fsRoot, "loader", "loader.conf"))
	if err != nil {
		// loader.conf is optional.
		loaderConf = make(map[string]string)
	}

	// TODO: Rank entries by version or machine-id attribute as suggested
	// in the spec (but not mandated, surprisingly).
	imgs := make(map[string]boot.OSImage)
	for _, f := range files {
		identifier := cutConf(filepath.Base(f))

		img, err := parseBLSEntry(f, fsRoot)
		if err != nil {
			log.Printf("BootLoaderSpec skipping entry %s: %v", f, err)
			continue
		}
		imgs[identifier] = img
	}

	return sortImages(loaderConf, imgs), nil
}

func sortImages(loaderConf map[string]string, imgs map[string]boot.OSImage) []boot.OSImage {
	// rankedImages = sort(default-images) + sort(remaining images)
	var rankedImages []boot.OSImage

	pattern, ok := loaderConf["default"]
	if !ok {
		// All images are default.
		pattern = "*"
	}

	var defaultIdents []string
	var otherIdents []string

	// Find default and non-default identifiers.
	for ident := range imgs {
		ok, err := filepath.Match(pattern, ident)
		if err != nil && ok {
			defaultIdents = append(defaultIdents, ident)
		} else {
			otherIdents = append(otherIdents, ident)
		}
	}

	// Sort them in the order we want them.
	sort.Sort(sort.Reverse(sort.StringSlice(defaultIdents)))
	sort.Sort(sort.Reverse(sort.StringSlice(otherIdents)))

	// Add images to rankedImages in that sorted order, defaults first.
	for _, ident := range defaultIdents {
		rankedImages = append(rankedImages, imgs[ident])
	}
	for _, ident := range otherIdents {
		rankedImages = append(rankedImages, imgs[ident])
	}
	return rankedImages
}

func parseConf(entryPath string) (map[string]string, error) {
	f, err := os.Open(entryPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	vals := make(map[string]string)

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
		vals[sline[0]] = strings.TrimSpace(sline[1])
	}
	return vals, nil
}

// The spec says "$BOOT/loader/ is the directory containing all files needed
// for Type #1 entries", but that's bullshit. Relative file names are indeed in
// the $BOOT/loader/ directory, but absolute path names are in $BOOT, as
// evidenced by the entries that kernel-install installs on Fedora 32.
func filePath(fsRoot, value string) string {
	if !filepath.IsAbs(value) {
		return filepath.Join(fsRoot, "loader", value)
	}
	return filepath.Join(fsRoot, value)
}

func parseLinuxImage(vals map[string]string, fsRoot string) (boot.OSImage, error) {
	linux := &boot.LinuxImage{}

	var cmdlines []string
	for key, val := range vals {
		switch key {
		case "linux":
			f, err := os.Open(filePath(fsRoot, val))
			if err != nil {
				return nil, err
			}
			linux.Kernel = f

		// TODO: initrd may be specified more than once.
		case "initrd":
			f, err := os.Open(filePath(fsRoot, val))
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
	if title, ok := vals["title"]; ok && len(title) > 0 {
		name = append(name, title)
	}
	if version, ok := vals["version"]; ok && len(version) > 0 {
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
func parseBLSEntry(entryPath, fsRoot string) (boot.OSImage, error) {
	vals, err := parseConf(entryPath)
	if err != nil {
		return nil, fmt.Errorf("error parsing config in %s: %w", entryPath, err)
	}

	var img boot.OSImage
	err = fmt.Errorf("neither linux, efi, nor multiboot present in BootLoaderSpec config")
	if _, ok := vals["linux"]; ok {
		img, err = parseLinuxImage(vals, fsRoot)
	} else if _, ok := vals["multiboot"]; ok {
		err = fmt.Errorf("multiboot not yet supported")
	} else if _, ok := vals["efi"]; ok {
		err = fmt.Errorf("EFI not yet supported")
	}
	if err != nil {
		return nil, fmt.Errorf("error parsing config in %s: %w", entryPath, err)
	}
	return img, nil
}
