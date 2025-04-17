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
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/ulog"
)

const (
	blsEntriesDir  = "loader/entries"
	blsEntriesDir2 = "boot/loader/entries"
	// Set a higher default rank for BLS. It should be booted prior to the
	// other local images.
	blsDefaultRank = 1
)

// ScanBLSEntries scans the filesystem root for valid BLS entries.
// This function skips over invalid or unreadable entries in an effort
// to return everything that is bootable. map variables is the parsed result
// from Grub parser that should be used by BLS parser, pass nil if there's none.
func ScanBLSEntries(l ulog.Logger, fsRoot string, variables map[string]string, grubDefaultSavedEntry string) ([]boot.OSImage, error) {
	entriesDir := filepath.Join(fsRoot, blsEntriesDir)

	files, err := filepath.Glob(filepath.Join(entriesDir, "*.conf"))
	if err != nil || len(files) == 0 {
		// Try blsEntriesDir2
		entriesDir = filepath.Join(fsRoot, blsEntriesDir2)
		files, err = filepath.Glob(filepath.Join(entriesDir, "*.conf"))
		if err != nil || len(files) == 0 {
			return nil, fmt.Errorf("no BootLoaderSpec entries found: %w", err)
		}
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
		identifier := strings.TrimSuffix(filepath.Base(f), ".conf")

		// If the config file name is the same as the Grub default option, pass true for grubDefaultFlag
		var img boot.OSImage
		var err error
		if strings.Compare(identifier, grubDefaultSavedEntry) == 0 {
			img, err = parseBLSEntry(f, fsRoot, variables, true)
		} else {
			img, err = parseBLSEntry(f, fsRoot, variables, false)
		}
		if err != nil {
			l.Printf("BootLoaderSpec skipping entry %s: %v", f, err)
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

func getGrubvalue(variables map[string]string, key string) (string, error) {
	if variables == nil {
		// Only return error for nil variables map.
		return "", fmt.Errorf("variables map is nil")
	}
	if val, ok := variables[key]; ok && len(val) > 0 {
		return val, nil
	}
	return "", nil
}

func parseLinuxImage(vals map[string]string, fsRoot string, variables map[string]string, grubDefaultFlag bool) (boot.OSImage, error) {
	linux := &boot.LinuxImage{}
	var cmdlines []string
	var tokens []string
	var value string
	for key, val := range vals {
		switch key {
		case "linux":
			f, err := os.Open(filePath(fsRoot, val))
			if err != nil {
				return nil, err
			}
			linux.Kernel = f

		// TODO: initrd may be specified more than once.
		// TODO: For now only process the first token, the rest are ignored, e.g. '$tuned_initrd'.
		case "initrd":
			tokens = strings.Split(val, " ")
			f, err := os.Open(filePath(fsRoot, tokens[0]))
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
			tokens = strings.Split(val, " ")
			var err error
			for _, w := range tokens {
				switch w {
				// TODO: GRUB/BLS parser should also get kernelopts from grubenv file
				case "$kernelopts":
					if value, err = getGrubvalue(variables, "kernelopts"); err != nil {
						return nil, fmt.Errorf("variables map is nil for $kernelopts")
					}
					if value == "" {
						// If it's not found, fallback to look for default_kernelopts
						log.Printf("kernelopts is empty, look for default_kernelopts\n")
						if value, _ = getGrubvalue(variables, "default_kernelopts"); value == "" {
							return nil, fmt.Errorf("no valid kernelopts is found")
						}
					}
					cmdlines = append(cmdlines, value)
				case "$tuned_params":
					if value, err = getGrubvalue(variables, "tuned_params"); err != nil {
						return nil, fmt.Errorf("variables map is nil for $tuned_params")
					}
					cmdlines = append(cmdlines, value)
				default:
					cmdlines = append(cmdlines, w)
				}
			}
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
	// If this is the default option, increase the BootRank by 1
	// when os.LookupEnv("BLS_BOOT_RANK") doesn't exist so it's not affected.
	if val, exist := os.LookupEnv("BLS_BOOT_RANK"); exist {
		if rank, err := strconv.Atoi(val); err == nil {
			linux.BootRank = rank
		}
	} else {
		if grubDefaultFlag {
			linux.BootRank = blsDefaultRank + 1
		} else {
			linux.BootRank = blsDefaultRank
		}
	}

	return linux, nil
}

// parseBLSEntry takes a Type #1 BLS entry and the directory of entries, and
// returns a LinuxImage.
// An error is returned if the syntax is wrong or required keys are missing.
func parseBLSEntry(entryPath, fsRoot string, variables map[string]string, grubDefaultFlag bool) (boot.OSImage, error) {
	vals, err := parseConf(entryPath)
	if err != nil {
		return nil, fmt.Errorf("error parsing config in %s: %w", entryPath, err)
	}

	var img boot.OSImage
	err = fmt.Errorf("neither linux, efi, nor multiboot present in BootLoaderSpec config")
	if _, ok := vals["linux"]; ok {
		img, err = parseLinuxImage(vals, fsRoot, variables, grubDefaultFlag)
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
