// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package diskboot

import (
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

type parserState int

const (
	search   parserState = iota // searching for a valid entry
	grub                        // building a grub entry
	syslinux                    // building a syslinux entry
)

type parser struct {
	state        parserState
	config       *Config
	entry        *Entry
	defaultName  string
	defaultIndex int
}

func (p *parser) parseSearch(line string) {
	trimmedLine := strings.TrimSpace(line)
	f := strings.Fields(trimmedLine)
	if len(f) == 0 {
		return
	}

	newEntry := false
	var name string

	switch strings.ToUpper(f[0]) {
	case "MENUENTRY": // grub
		p.state = grub
		newEntry = true
		repNames := strings.Replace(trimmedLine, "'", "\"", -1)
		names := strings.Split(repNames, "\"")
		if len(names) > 1 {
			name = names[1]
		}
	case "SET": // grub
		if len(f) > 1 {
			p.parseSearchHandleSet(f[1])
		}
	case "LABEL": // syslinux
		p.state = syslinux
		newEntry = true
		name = trimmedLine[6:]
		if name == "" {
			name = "linux" // default for syslinux
		}
	case "DEFAULT": // syslinux
		if len(f) > 1 {
			label := strings.Join(f[1:], " ")
			if !strings.HasSuffix(label, ".c32") {
				p.defaultName = label
			}
		}
	}

	if newEntry {
		p.entry = &Entry{
			Name: name,
			Type: Elf,
		}
	}
}

func (p *parser) parseSearchHandleSet(val string) {
	val = strings.Replace(val, "\"", "", -1)
	expr := strings.Split(val, "=")
	if len(expr) != 2 {
		return
	}

	switch expr[0] {
	case "default":
		index, err := strconv.Atoi(expr[1])
		if err == nil {
			p.defaultIndex = index
		}
	default:
		// TODO: handle variables when grub conditionals are implemented
		return
	}
}

func (p *parser) parseGrubEntry(line string) {
	trimmedLine := strings.TrimSpace(line)
	f := strings.Fields(trimmedLine)
	if len(f) == 0 {
		return
	}

	switch f[0] {
	case "}":
		p.finishEntry()
	case "multiboot":
		p.entry.Type = Multiboot
		p.entry.Modules = append(p.entry.Modules, NewModule(f[1], f[2:]))
	case "module":
		var filteredParams []string
		for _, param := range f {
			if param != "--nounzip" {
				filteredParams = append(filteredParams, param)
			}
		}
		p.entry.Modules = append(p.entry.Modules,
			NewModule(filteredParams[1], filteredParams[2:]))
	case "linux":
		p.entry.Modules = append(p.entry.Modules, NewModule(f[1], f[2:]))
	case "initrd":
		p.entry.Modules = append(p.entry.Modules, NewModule(f[1], nil))
	}
}

func (p *parser) parseSyslinuxEntry(line string) {
	trimmedLine := strings.TrimSpace(line)
	if len(trimmedLine) == 0 {
		p.finishEntry()
		return
	}

	f := strings.Fields(trimmedLine)
	val := strings.Join(f[1:], " ")
	switch strings.ToUpper(f[0]) {
	case "LABEL": // new entry, finish this one and start a new one
		p.finishEntry()
		p.parseSearch(line)
	case "MENU":
		if len(f) > 1 {
			switch strings.ToUpper(f[1]) {
			case "LABEL":
				tempName := strings.Join(f[2:], " ")
				p.entry.Name = strings.Replace(tempName, "^", "", -1)
			case "DEFAULT":
				p.defaultName = p.entry.Name
			}
		}
	case "LINUX", "KERNEL":
		p.parseSyslinuxKernel(val)
	case "INITRD":
		p.entry.Modules = append(p.entry.Modules, NewModule(f[1], nil))
	case "APPEND":
		p.parseSyslinuxAppend(val)
	}
}

func (p *parser) parseSyslinuxKernel(val string) {
	if strings.HasSuffix(val, "mboot.c32") {
		p.entry.Type = Multiboot
	} else if strings.HasSuffix(val, ".c32") {
		// skip this entry - not valid for kexec
		p.state = search
	} else {
		p.entry.Modules = append(p.entry.Modules, NewModule(val, nil))
	}
}

func (p *parser) parseSyslinuxAppend(val string) {
	if p.entry.Type == Multiboot {
		// split params by "---" for each module
		modules := strings.Split(val, " --- ")
		for _, module := range modules {
			moduleFields := strings.Fields(module)
			p.entry.Modules = append(p.entry.Modules,
				NewModule(moduleFields[0], moduleFields[1:]))
		}
	} else {
		if len(p.entry.Modules) == 0 {
			// TODO: log error
			return
		}
		p.entry.Modules[0].Params = val
	}
}

func (p *parser) finishEntry() {
	// skip empty entries
	if len(p.entry.Modules) == 0 {
		return
	}

	// try to fix up initrd from kernel params
	if len(p.entry.Modules) == 1 && p.entry.Type == Elf {
		var initrd string
		var newParams []string

		params := strings.Fields(p.entry.Modules[0].Params)
		for _, param := range params {
			if strings.HasPrefix(param, "initrd=") {
				initrd = param[7:]
			} else {
				newParams = append(newParams, param)
			}
		}
		if initrd != "" {
			p.entry.Modules = append(p.entry.Modules, NewModule(initrd, nil))
			p.entry.Modules[0].Params = strings.Join(newParams, " ")
		}
	}

	appendPath, err := filepath.Rel(p.config.MountPath, filepath.Dir(p.config.ConfigPath))
	if err != nil {
		log.Fatal("Config file path not relative to mount path")
	}
	for i, module := range p.entry.Modules {
		if !strings.HasPrefix(module.Path, "/") {
			module.Path = filepath.Join("/"+appendPath, module.Path)
		}
		module.Path = filepath.Clean(module.Path)
		p.entry.Modules[i] = module
	}

	p.state = search
	p.config.Entries = append(p.config.Entries, *p.entry)
}

func (p *parser) parseLines(lines []string) {
	p.state = search

	for _, line := range lines {
		switch p.state {
		case search:
			p.parseSearch(line)
		case grub:
			p.parseGrubEntry(line)
		case syslinux:
			p.parseSyslinuxEntry(line)
		}
	}

	if p.state == syslinux {
		p.finishEntry()
	}
}

// ParseConfig attempts to construct a valid boot Config from the location
// and lines contents passed in.
func ParseConfig(mountPath, configPath string, lines []string) *Config {
	p := &parser{
		config: &Config{
			MountPath:    mountPath,
			ConfigPath:   configPath,
			DefaultEntry: -1,
		},
		defaultIndex: -1,
	}
	p.parseLines(lines)

	if p.defaultName != "" {
		for i, entry := range p.config.Entries {
			if entry.Name == p.defaultName {
				p.config.DefaultEntry = i
			}
		}
	}
	if p.defaultIndex >= 0 && len(p.config.Entries) > p.defaultIndex {
		p.config.DefaultEntry = p.defaultIndex
	}

	return p.config
}
