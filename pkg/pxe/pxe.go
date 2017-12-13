// Package pxe aims to implement the PXE specification.
//
// See http://www.pix.net/software/pxeboot/archive/pxespec.pdf
package pxe

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/url"
	"path"
	"strings"
)

var (
	// ErrConfigNotFound is returned when no suitable configuration file
	// was found by AppendFile.
	ErrConfigNotFound = errors.New("configuration file was not found")

	// ErrDefaultEntryNotFound is returned when the configuration file
	// names a default label that is not part of the configuration.
	ErrDefaultEntryNotFound = errors.New("default label not found in configuration")
)

// Entry encapsulates a Syslinux "label" config directive.
type Entry struct {
	// Kernel is the kernel for this label.
	Kernel io.Reader

	// Initrd is the initial ramdisk for this label.
	Initrd io.Reader

	// Cmdline is the list of kernel command line parameters.
	Cmdline string
}

// Config encapsulates a parsed Syslinux configuration file.
//
// See http://www.syslinux.org/wiki/index.php?title=Config for the
// configuration file specification.
//
// TODO: Tear apart parser internals from Config.
type Config struct {
	// Entries is a map of label name -> label configuration.
	Entries map[string]*Entry

	// DefaultEntry is the default label key to use.
	//
	// If DefaultEntry is non-empty, the label is guaranteed to exist in
	// `Entries`.
	DefaultEntry string

	// Parser internals.
	globalAppend string
	scope        scope
	curEntry     string
	wd           *url.URL
	schemes      Schemes
}

type scope uint8

const (
	scopeGlobal scope = iota
	scopeEntry
)

// NewConfig returns a new PXE parser using working directory `wd` and default
// schemes.
//
// See NewConfigWithSchemes for more details.
func NewConfig(wd *url.URL) *Config {
	return NewConfigWithSchemes(wd, DefaultSchemes)
}

// NewConfigWithSchemes returns a new PXE parser using working directory `wd`
// and schemes `s`.
//
// If a path encountered in a configuration file is relative instead of a full
// URI, `wd` is used as the "working directory" of that relative path; the
// resulting URI is roughly `wd.String()/path`.
//
// `s` is used to get files referred to by URIs.
func NewConfigWithSchemes(wd *url.URL, s Schemes) *Config {
	return &Config{
		Entries: make(map[string]*Entry),
		scope:   scopeGlobal,
		wd:      wd,
		schemes: s,
	}
}

// FindConfigFile probes for config files based on the Mac and IP given.
func (c *Config) FindConfigFile(mac net.HardwareAddr, ip net.IP) error {
	for _, relname := range probeFiles(mac, ip) {
		if err := c.AppendFile(relname); err != ErrConfigNotFound {
			return err
		}
	}
	return fmt.Errorf("no valid pxelinux config found")
}

// ParseConfigFile parses a PXE/Syslinux configuration as specified in
// http://www.syslinux.org/wiki/index.php?title=Config
//
// Currently, only the APPEND, INCLUDE, KERNEL, LABEL, DEFAULT, and INITRD
// directives are partially supported.
//
// `wd` is the default scheme, host, and path for any files named as a
// relative path. The default path for config files is assumed to be
// `wd.Path`/pxelinux.cfg/.
func ParseConfigFile(uri string, wd *url.URL) (*Config, error) {
	c := NewConfig(wd)
	if err := c.AppendFile(uri); err != nil {
		return nil, err
	}
	return c, nil
}

// AppendFile parses the config file downloaded from `uri` and adds it to `c`.
func (c *Config) AppendFile(uri string) error {
	cfgWd := *c.wd
	// The default location for looking for configuration is
	// CWD/pxelinux.cfg/, so if `uri` is just a relative path, this will be
	// the default directory.
	cfgWd.Path = path.Join(cfgWd.Path, "pxelinux.cfg")

	r, err := c.schemes.GetFile(uri, &cfgWd)
	if err != nil {
		return ErrConfigNotFound
	}
	config, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return c.Append(string(config))
}

// Append parses `config` and adds the respective configuration to `c`.
func (c *Config) Append(config string) error {
	// Here's a shitty parser.
	for _, line := range strings.Split(config, "\n") {
		// This is stupid. There should be a FieldsN(...).
		kv := strings.Fields(line)
		if len(kv) <= 1 {
			continue
		}
		directive := strings.ToLower(kv[0])
		var arg string
		if len(kv) == 2 {
			arg = kv[1]
		} else {
			arg = strings.Join(kv[1:], " ")
		}

		switch directive {
		case "default":
			c.DefaultEntry = arg

		case "include":
			if err := c.AppendFile(arg); err == ErrConfigNotFound {
				// TODO: What to do in this case?
				// Return for now. Test this with pxelinux.
				return err
			} else if err != nil {
				return err
			}

		case "label":
			// We forever enter label scope.
			c.scope = scopeEntry
			c.curEntry = arg
			c.Entries[c.curEntry] = &Entry{}
			c.Entries[c.curEntry].Cmdline = c.globalAppend

		case "kernel":
			k, err := c.schemes.LazyGetFile(arg, c.wd)
			if err != nil {
				return err
			}
			c.Entries[c.curEntry].Kernel = k

		case "initrd":
			i, err := c.schemes.LazyGetFile(arg, c.wd)
			if err != nil {
				return err
			}
			c.Entries[c.curEntry].Initrd = i

		case "append":
			switch c.scope {
			case scopeGlobal:
				c.globalAppend = arg

			case scopeEntry:
				if arg == "-" {
					c.Entries[c.curEntry].Cmdline = ""
				} else {
					c.Entries[c.curEntry].Cmdline = arg
				}
			}
		}
	}

	// Go through all labels and download the initrds.
	for _, label := range c.Entries {
		// If the initrd was set via the INITRD directive, don't
		// overwrite that.
		//
		// TODO(hugelgupf): Is this really what syslinux does? Does
		// INITRD trump cmdline? Does it trump global? What if both the
		// directive and cmdline initrd= are set? Does it depend on the
		// order in the config file? (My current best guess: order.)
		if label.Initrd != nil {
			continue
		}

		for _, opt := range strings.Fields(label.Cmdline) {
			optkv := strings.Split(opt, "=")
			if optkv[0] != "initrd" {
				continue
			}

			i, err := c.schemes.LazyGetFile(optkv[1], c.wd)
			if err != nil {
				return err
			}
			label.Initrd = i
		}
	}

	if len(c.DefaultEntry) > 0 {
		if _, ok := c.Entries[c.DefaultEntry]; !ok {
			return ErrDefaultEntryNotFound
		}
	}
	return nil

}

func probeFiles(ethernetMac net.HardwareAddr, ip net.IP) []string {
	files := make([]string, 0, 10)
	// Skipping client UUID. Figure that out later.

	// MAC address.
	files = append(files, fmt.Sprintf("01-%s", strings.ToLower(strings.Replace(ethernetMac.String(), ":", "-", -1))))

	// IP address in upper case hex, chopping one letter off at a time.
	ipf := strings.ToUpper(hex.EncodeToString(ip))
	for n := len(ipf); n >= 1; n-- {
		files = append(files, ipf[:n])
	}
	files = append(files, "default")
	return files
}
