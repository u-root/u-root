package diskboot

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"syscall"
)

// Config contains boot entries for a single configuration file (grub, syslinux, etc.)
type Config struct {
	MountPath    string
	ConfigPath   string
	Entries      []Entry
	DefaultEntry int
}

// EntryType dictates the method by which kexec should use to load the new kernel
type EntryType int

// EntryType can be either Elf or Multiboot
const (
	Elf EntryType = iota
	Multiboot
)

// Module represents a path to a binary along with arguments for its execution
// The path in the module is relative to the mount path
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

// KexecLoad calls the appropriate kexec load routines based on the type of Entry
func (e *Entry) KexecLoad() error {
	switch e.Type {
	case Multiboot:
		// TODO: implement using kexec_load syscall
		return syscall.ENOSYS
	case Elf:
		// TODO: implement using kexec_file_load syscall
		// e.Module[0].Path is kernel
		// e.Module[0].Params is kernel parameters
		// e.Module[1].Path is initrd
		return syscall.ENOSYS
	}
	return nil
}

type location struct {
	Path string
	Type parserState
}

// TODO: change to search and autodetect format
var (
	locations = []location{
		{"boot/grub/grub.cfg", grub},
		{"isolinux/isolinux.cfg", syslinux},
	}
)

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

		config := ParseConfig(mountPath, configPath, lines)
		if config != nil {
			configs = append(configs, config)
		}
	}

	return configs
}

func loadSyslinuxLines(configPath string, contents []byte) []string {
	var newLines []string
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		fields := strings.Fields(trimmedLine)
		if len(fields) == 2 && strings.ToUpper(fields[0]) == "INCLUDE" {
			includePath := filepath.Join(filepath.Dir(configPath), fields[1])
			includeContents, err := ioutil.ReadFile(includePath)
			if err != nil {
				// TODO: log error
				continue
			}
			includeLines := loadSyslinuxLines(includePath, includeContents)
			newLines = append(newLines, includeLines...)
		} else {
			newLines = append(newLines, line)
		}
	}
	return newLines
}
