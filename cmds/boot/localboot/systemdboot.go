// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/boot/jsonboot"
	"github.com/u-root/u-root/pkg/mount/block"
)

// ParseSystemdbootCfg will try to parse a single file based on the Boot Loader Specification
// TODO It is not 100% compliant and only parses the bare minium options we need in order to sucessfully kexec to Linux.
func parseSystemdbootCfg(devices block.BlockDevices, config string, basedir string) jsonboot.BootConfig {
	var cfg jsonboot.BootConfig

	for _, line := range strings.Split(config, "\n") {
		debug(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		key := strings.SplitN(line, " ", 2)
		if len(key) < 2 {
			// this should not happen
			log.Println("incorrectly formatted")
			continue
		}
		value := strings.TrimLeft(key[1], " ")
		switch key[0] {
		case "linux":
			cfg.Kernel = basedir + value
		case "initrd":
			cfg.Initramfs = basedir + value
		case "options":
			cfg.KernelArgs = value
		case "title":
			cfg.Name = value
		}
	}
	return cfg
}

// Looks into basedir + /loader/entries to search for Boot Loader specification entries.
// It will try to parse all available entries and return them as BootConfig configurations
func ScanSystemdbootConfigs(devices block.BlockDevices, basedir string) []jsonboot.BootConfig {
	files, err := os.ReadDir(basedir + "/loader/entries")
	if err != nil {
		log.Printf("error scanning %s/loader/entries directory: %v\n", basedir, err)
		return nil
	}
	log.Printf("found systemd boot config files: %v\n", files)

	bootconfigs := make([]jsonboot.BootConfig, 0)
	for _, file := range files {
		debug("filename %s, isdir: %v\n", file.Name(), file.IsDir())
		if !file.IsDir() {
			filepath := basedir + "/loader/entries/" + file.Name()
			data, err := os.ReadFile(filepath)
			if err != nil {
				log.Printf("error reading file: %s", filepath)
				continue
			}
			cfg := parseSystemdbootCfg(devices, string(data), basedir)
			bootconfigs = append(bootconfigs, cfg)
		}
	}
	return bootconfigs
}
