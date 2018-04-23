package main

import (
	"os"
	"path"
	"strings"
)

var (
	GrubPaths = []string{
		// grub2
		"boot/grub2/grub.cfg",
		"boot/grub2.cfg",
		"grub2/grub.cfg",
		"grub2.cfg",
		// grub legacy
		"boot/grub/grub.cfg",
		"boot/grub.cfg",
		"grub/grub.cfg",
		"grub.cfg",
	}
)

// ParseGrubCfg parses the content of a grub.cfg and returns a list of
// BootConfig structures, one for each menuentry, in the same order as they
// appear in grub.cfg. All opened kernel and initrd files are relative to
// basedir.
func ParseGrubCfg(grubcfg string, basedir string) []BootConfig {
	// This parser sucks. It's not even a parser, it just looks for lines
	// starting with menuentry, linux or initrd.
	// TODO use a parser, e.g. https://github.com/alecthomas/participle
	bootconfigs := make([]BootConfig, 0)
	inMenuEntry := false
	var cfg *BootConfig
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
			}
			inMenuEntry = true
			cfg = new(BootConfig)
		} else if inMenuEntry {
			// otherwise look for kernel and initramfs configuration
			if len(sline) < 2 {
				// surely not a valid linux or initrd directive, skip it
				continue
			}
			if sline[0] == "linux" || sline[0] == "linux16" {
				kernel := sline[1]
				cmdline := strings.Join(sline[2:], " ")
				fullpath := path.Join(basedir, kernel)
				fd, err := os.Open(fullpath)
				if err != nil {
					debug("error opening kernel file %s: %v", fullpath, err)
				}
				cfg.Kernel = fd
				cfg.KernelName = kernel
				cfg.Cmdline = cmdline
			} else if sline[0] == "initrd" || sline[0] == "initrd16" {
				initrd := sline[1]
				fullpath := path.Join(basedir, initrd)
				fd, err := os.Open(fullpath)
				if err != nil {
					debug("error opening initrd file %s: %v", fullpath, err)
				}
				cfg.Initrd = fd
				cfg.InitrdName = initrd
			}
		}
	}
	// append last kernel config if it wasn't already
	if inMenuEntry && cfg.IsValid() {
		bootconfigs = append(bootconfigs, *cfg)
	}
	return bootconfigs
}
