package diskboot

import (
	"log"
	"path/filepath"
	"strings"
)

type parserState int

const (
	search   parserState = iota // searching for a valid entry
	grub                        // building a grub entry
	syslinux                    // building a syslinux entry
)

type parser struct {
	state  parserState
	config *Config
	entry  *Entry
}

func (p *parser) parseSearch(line string) {
	trimmedLine := strings.TrimSpace(line)
	f := strings.Fields(trimmedLine)
	if len(f) == 0 {
		return
	}

	newEntry := false
	var name string

	// TODO: add name to the entry
	switch f[0] {
	case "menuentry", "MENUENTRY":
		p.state = grub
		newEntry = true
		names := strings.Split(strings.Replace(trimmedLine, "'", "\"", -1), "\"")
		if len(names) > 1 {
			name = names[1]
		}
	case "label", "LABEL":
		p.state = syslinux
		newEntry = true
		name = trimmedLine[6:]
	}

	if newEntry {
		p.entry = &Entry{
			Name: name,
			Type: Elf,
		}
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
		// TODO: fix to remove "--nounzip"
		p.entry.Modules = append(p.entry.Modules, NewModule(f[1], f[2:]))
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
	switch f[0] {
	case "label", "LABEL": // new entry, finish this one and start a new one
		p.finishEntry()
		p.parseSearch(line)
	case "menu", "MENU":
		if len(f) > 1 && (strings.ToUpper(f[1]) == "LABEL") {
			p.entry.Name = strings.Replace(strings.Join(f[2:], " "), "^", "", -1)
		}
	case "linux", "LINUX", "kernel", "KERNEL":
		p.parseSyslinuxKernel(val)
	case "initrd", "INITRD":
		p.entry.Modules = append(p.entry.Modules, NewModule(f[1], nil))
	case "append", "APPEND":
		p.parseSyslinuxAppend(val)
	}
}

func (p *parser) parseSyslinuxKernel(val string) {
	if strings.HasSuffix(val, "mboot.c32") {
		p.entry.Type = Multiboot
	} else if strings.HasSuffix(val, ".c32") {
		// skip this entry
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

func (p *parser) parseLines(lines []string) error {
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
	return nil
}

// ParseConfig attemps to construct a valid boot Config from the location
// and lines contents passed in.
func ParseConfig(mountPath, configPath string, lines []string) *Config {
	parser := &parser{
		config: &Config{
			MountPath:    mountPath,
			ConfigPath:   configPath,
			DefaultEntry: -1,
		},
	}
	err := parser.parseLines(lines)
	if err != nil {
		return nil
	}
	return parser.config
}
