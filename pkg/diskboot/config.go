// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package diskboot

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/u-root/u-root/pkg/cmdline"
	"github.com/u-root/u-root/pkg/kexec"
)

// Config contains boot entries for a single configuration file
// (grub, syslinux, etc.)
type Config struct {
	MountPath    string
	ConfigPath   string
	Entries      []Entry
	DefaultEntry int
}

// EntryType dictates the method by which kexec should use to load
// the new kernel
type EntryType int

// EntryType can be either Elf or Multiboot
const (
	Elf EntryType = iota
	Multiboot
)

// Module represents a path to a binary along with arguments for its
// xecution. The path in the module is relative to the mount path
type Module struct {
	Path   string
	Params string
}

func (m Module) String() string {
	return fmt.Sprintf("|'%v' (%v)|", m.Path, m.Params)
}

// NewModule constructs a module for a boot entry
func NewModule(path string, args []string) Module {
	return Module{
		Path:   path,
		Params: strings.Join(args, " "),
	}
}

// Entry contains the necessary info to kexec into a new kernel
type Entry struct {
	Name    string
	Type    EntryType
	Modules []Module
}

// KexecLoad calls the appropriate kexec load routines based on the
// type of Entry
func (e *Entry) KexecLoad(mountPath string, filterCmdline cmdline.Filter, dryrun bool) error {
	switch e.Type {
	case Multiboot:
		// TODO: implement using kexec_load syscall
		return syscall.ENOSYS
	case Elf:
		// TODO: implement using kexec_file_load syscall
		// e.Module[0].Path is kernel
		// e.Module[0].Params is kernel parameters
		// e.Module[1].Path is initrd
		if len(e.Modules) < 1 {
			return fmt.Errorf("missing kernel")
		}
		var ramfs *os.File
		kernelPath := filepath.Join(mountPath, e.Modules[0].Path)
		log.Print("Kernel Path:", kernelPath)
		kernel, err := os.OpenFile(kernelPath, os.O_RDONLY, 0)
		commandline := e.Modules[0].Params
		if filterCmdline != nil {
			commandline = filterCmdline.Update(commandline)
		}

		log.Print("Kernel Params:", commandline)
		if err != nil {
			return fmt.Errorf("failed to load kernel: %v", err)
		}
		if len(e.Modules) > 1 {
			ramfsPath := filepath.Join(mountPath, e.Modules[1].Path)
			log.Print("Ramfs Path:", ramfsPath)
			ramfs, err = os.OpenFile(ramfsPath, os.O_RDONLY, 0)
			if err != nil {
				return fmt.Errorf("failed to load ramfs: %v", err)
			}
		}
		if !dryrun {
			return kexec.FileLoad(kernel, ramfs, commandline)
		}
	}
	return nil
}

type location struct {
	Path string
	Type parserState
}

var (
	locations = []location{
		{"boot/grub/grub.cfg", grub},
		{"grub/grub.cfg", grub},
		{"grub2/grub.cfg", grub},
		// following entries from the syslinux wiki
		// TODO: add priorities override (top over bottom)
		{"boot/isolinux/isolinux.cfg", syslinux},
		{"isolinux/isolinux.cfg", syslinux},
		{"isolinux.cfg", syslinux},
		{"boot/syslinux/syslinux.cfg", syslinux},
		{"syslinux/syslinux.cfg", syslinux},
		{"syslinux.cfg", syslinux},
	}
)

// TODO: add iso handling along with iso_path variable replacement

// FindConfigs searching the path for valid boot configuration files
// and returns a Config for each valid instance found.
func FindConfigs(mountPath string) []*Config {
	var configs []*Config

	for _, location := range locations {
		configPath := filepath.Join(mountPath, location.Path)
		contents, err := ioutil.ReadFile(configPath)
		if err != nil {
			// TODO: log error
			continue
		}

		var lines []string
		if location.Type == syslinux {
			lines = loadSyslinuxLines(configPath, contents)
		} else {
			lines = strings.Split(string(contents), "\n")
		}

		configs = append(configs, ParseConfig(mountPath, configPath, lines))
	}

	return configs
}

func loadSyslinuxLines(configPath string, contents []byte) []string {
	// TODO: just parse includes inline with syslinux specific parser
	var newLines, includeLines []string
	menuKernel := false

	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(strings.TrimSpace(line))
		includeDir := filepath.Dir(configPath)
		if len(fields) == 2 && strings.ToUpper(fields[0]) == "INCLUDE" {
			includeLines = loadSyslinuxInclude(includeDir, fields[1])
		} else if len(fields) == 3 &&
			strings.ToUpper(fields[0]) == "MENU" &&
			strings.ToUpper(fields[1]) == "INCLUDE" {
			includeLines = loadSyslinuxInclude(includeDir, fields[2])
		} else if len(fields) > 1 &&
			strings.ToUpper(fields[0]) == "APPEND" &&
			menuKernel {
			includeLines = []string{}
			for _, includeFile := range fields[1:] {
				includeLines = append(includeLines,
					loadSyslinuxInclude(includeDir, includeFile)...)
			}
		} else {
			if len(fields) > 0 && strings.ToUpper(fields[0]) == "LABEL" {
				menuKernel = false
			} else if len(fields) == 2 &&
				strings.ToUpper(fields[0]) == "KERNEL" &&
				(strings.ToUpper(fields[1]) == "VESAMENU.C32" ||
					strings.ToUpper(fields[1]) == "MENU.C32") {
				menuKernel = true
			}
			includeLines = []string{line}
		}
		newLines = append(newLines, includeLines...)
	}
	return newLines
}

func loadSyslinuxInclude(includePath, includeFile string) []string {
	path := filepath.Join(includePath, includeFile)
	includeContents, err := ioutil.ReadFile(path)
	if err != nil {
		// TODO: log error
		return nil
	}
	return loadSyslinuxLines(path, includeContents)
}
