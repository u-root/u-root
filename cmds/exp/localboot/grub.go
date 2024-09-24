// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"encoding/hex"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"github.com/u-root/u-root/pkg/boot/jsonboot"
	"github.com/u-root/u-root/pkg/mount/block"
)

// List of directories where to recursively look for grub config files. The root dorectory
// of each mountpoint, these folders inside the mountpoint and all subfolders
// of these folders are searched
var (
	GrubSearchDirectories = []string{
		"boot",
		"EFI",
		"efi",
		"grub",
		"grub2",
	}
)

// Limits rekursive search of grub files. It is the maximum directory depth
// that is searched through. Since on efi partitions grub files reside usually
// at /boot/efi/EFI/distro/ , 4 might be a good choice.
const searchDepth = 4

type grubVersion int

var (
	grubV1 grubVersion = 1
	grubV2 grubVersion = 2
)

func isGrubSearchDir(dirname string) bool {
	for _, dir := range GrubSearchDirectories {
		if dirname == dir {
			return true
		}
	}
	return false
}

// ParseGrubCfg parses the content of a grub.cfg and returns a list of
// BootConfig structures, one for each menuentry, in the same order as they
// appear in grub.cfg. All opened kernel and initrd files are relative to
// basedir.
func ParseGrubCfg(ver grubVersion, devices block.BlockDevices, grubcfg string, basedir string) []jsonboot.BootConfig {
	// This parser sucks. It's not even a parser, it just looks for lines
	// starting with menuentry, linux or initrd.
	// TODO use a parser, e.g. https://github.com/alecthomas/participle
	if ver != grubV1 && ver != grubV2 {
		log.Printf("Warning: invalid GRUB version: %d", ver)
		return nil
	}
	kernelBasedir := basedir
	bootconfigs := make([]jsonboot.BootConfig, 0)
	inMenuEntry := false
	var cfg *jsonboot.BootConfig
	for _, line := range strings.Split(grubcfg, "\n") {
		// remove all leading spaces as they are not relevant for the config
		// line
		line = strings.TrimLeft(line, " ")
		sline := strings.Fields(line)
		if len(sline) == 0 {
			continue
		}
		if sline[0] == "menuentry" {
			// if a "menuentry", start a new boot config
			if cfg != nil {
				// save the previous boot config, if any
				if cfg.IsValid() {
					// only consider valid boot configs, i.e. the ones that have
					// both kernel and initramfs
					bootconfigs = append(bootconfigs, *cfg)
				}
				// reset kernelBaseDir
				kernelBasedir = basedir
			}
			inMenuEntry = true
			cfg = new(jsonboot.BootConfig)
			name := ""
			if len(sline) > 1 {
				name = strings.Join(sline[1:], " ")
				name = unquoteGrubString(name)
				name = strings.Split(name, "--")[0]
			}
			cfg.Name = name
		} else if inMenuEntry {
			// check if location of kernel is at an other partition
			// see https://www.gnu.org/software/grub/manual/grub/html_node/search.html
			if sline[0] == "search" {
				for _, str1 := range sline {
					if str1 == "--set=root" {
						log.Printf("Kernel seems to be on an other partition then the grub.cfg file")
						for _, str2 := range sline {
							if isValidFsUUID(str2) {
								kernelFsUUID := str2
								log.Printf("fs-uuid: %s", kernelFsUUID)
								partitions := devices.FilterFSUUID(kernelFsUUID)
								if len(partitions) == 0 {
									log.Printf("WARNING: No partition found with filesystem UUID:'%s' to load kernel from!", kernelFsUUID) // TODO throw error ?
									continue
								}
								if len(partitions) > 1 {
									log.Printf("WARNING: more than one partition found with the given filesystem UUID. Using the first one")
								}
								dev := partitions[0]
								kernelBasedir = path.Dir(kernelBasedir)
								kernelBasedir = path.Join(kernelBasedir, dev.Name)
								log.Printf("Kernel is on: %s", dev.Name)
							}
						}
					}
				}
			}
			// otherwise look for kernel and initramfs configuration
			if len(sline) < 2 {
				// surely not a valid linux or initrd directive, skip it
				continue
			}
			if sline[0] == "linux" || sline[0] == "linux16" || sline[0] == "linuxefi" {
				kernel := sline[1]
				cmdline := strings.Join(sline[2:], " ")
				cmdline = unquoteGrubString(cmdline)
				cfg.Kernel = path.Join(kernelBasedir, kernel)
				cfg.KernelArgs = cmdline
			} else if sline[0] == "initrd" || sline[0] == "initrd16" || sline[0] == "initrdefi" {
				initrd := sline[1]
				cfg.Initramfs = path.Join(kernelBasedir, initrd)
			} else if sline[0] == "multiboot" || sline[0] == "multiboot2" {
				multiboot := sline[1]
				cmdline := strings.Join(sline[2:], " ")
				cmdline = unquoteGrubString(cmdline)
				cfg.Multiboot = path.Join(kernelBasedir, multiboot)
				cfg.MultibootArgs = cmdline
			} else if sline[0] == "module" || sline[0] == "module2" {
				module := sline[1]
				cmdline := strings.Join(sline[2:], " ")
				cmdline = unquoteGrubString(cmdline)
				module = path.Join(kernelBasedir, module)
				if cmdline != "" {
					module = module + " " + cmdline
				}
				cfg.Modules = append(cfg.Modules, module)
			}
		}
	}

	// append last kernel config if it wasn't already
	if inMenuEntry && cfg.IsValid() {
		bootconfigs = append(bootconfigs, *cfg)
	}
	return bootconfigs
}

func isValidFsUUID(uuid string) bool {
	for _, h := range strings.Split(uuid, "-") {
		if _, err := hex.DecodeString(h); err != nil {
			return false
		}
	}
	return true
}

func unquoteGrubString(text string) string {
	// unquote the string to prevent special characters used by GRUB
	// from being passed thru kexec
	// https://www.gnu.org/software/grub/manual/grub/grub.html#Quoting
	// TODO unquote everything, not just \$
	return strings.Replace(text, `\$`, "$", -1)
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

// ScanGrubConfigs looks for grub2 and grub legacy config files in the known
// locations and returns a list of boot configurations.
func ScanGrubConfigs(devices block.BlockDevices, basedir string) []jsonboot.BootConfig {
	bootconfigs := make([]jsonboot.BootConfig, 0)
	err := filepath.Walk(basedir, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
		currentPath, _, _ = transform.String(t, currentPath)
		if info.IsDir() {
			if path.Dir(currentPath) == basedir && !isGrubSearchDir(path.Base(currentPath)) {
				debug("Skip %s: not significant", currentPath)
				// skip irrelevant toplevel directories
				return filepath.SkipDir
			}
			p, err := filepath.Rel(basedir, currentPath)
			if err != nil {
				return err
			}
			depth := len(strings.Split(p, string(os.PathSeparator)))
			if depth > searchDepth {
				debug("Skip %s, depth limit", currentPath)
				// skip
				return filepath.SkipDir
			}
			debug("Step into %s", currentPath)
			// continue
			return nil
		}
		cfgname := info.Name()
		var ver grubVersion
		switch cfgname {
		case "grub.cfg":
			ver = grubV1
		case "grub2.cfg":
			ver = grubV2
		default:
			return nil
		}
		log.Printf("Parsing %s", currentPath)
		data, err := os.ReadFile(currentPath)
		if err != nil {
			return err
		}
		cfgs := ParseGrubCfg(ver, devices, string(data), basedir)
		bootconfigs = append(bootconfigs, cfgs...)
		return nil
	})
	if err != nil {
		log.Printf("filepath.Walk error: %v", err)
	}
	return bootconfigs
}
