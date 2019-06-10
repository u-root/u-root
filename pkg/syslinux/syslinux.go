// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package syslinux implements a syslinux config file parser.
//
// See http://www.syslinux.org/wiki/index.php?title=Config for general syslinux
// config features.
//
// Currently, only the APPEND, INCLUDE, KERNEL, LABEL, DEFAULT, and INITRD
// directives are partially supported.
package syslinux

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/uio"
	"github.com/u-root/u-root/pkg/urlfetch"
)

var (
	// ErrDefaultEntryNotFound is returned when the configuration file
	// names a default label that is not part of the configuration.
	ErrDefaultEntryNotFound = errors.New("default label not found in configuration")
)

// Config encapsulates a parsed Syslinux configuration file.
//
// See http://www.syslinux.org/wiki/index.php?title=Config for the
// configuration file specification.
type Config struct {
	// Entries is a map of label name -> label configuration.
	Entries map[string]*boot.LinuxImage

	// DefaultEntry is the default label key to use.
	//
	// If DefaultEntry is non-empty, the label is guaranteed to exist in
	// `Entries`.
	DefaultEntry string
}

// ParseConfigFile parses a Syslinux configuration as specified in
// http://www.syslinux.org/wiki/index.php?title=Config
//
// Currently, only the APPEND, INCLUDE, KERNEL, LABEL, DEFAULT, and INITRD
// directives are partially supported.
//
// urlfetch.DefaultSchemes is used to fetch any files that must be parsed or
// provided.
//
// `wd` is the default scheme, host, and path for any files named as a
// relative path - e.g. kernel, include, and initramfs paths are requested
// relative to the wd. The default path for config files is assumed to be
// `wd.Path`/pxelinux.cfg/.
func ParseConfigFile(url string, wd *url.URL) (*Config, error) {
	p := newParser(wd)
	if err := p.appendFile(url); err != nil {
		return nil, err
	}
	return p.config, nil
}

type parser struct {
	config *Config

	// parser internals.
	globalAppend string
	scope        scope
	curEntry     string
	wd           *url.URL
	schemes      urlfetch.Schemes
}

type scope uint8

const (
	scopeGlobal scope = iota
	scopeEntry
)

// newParser returns a new PXE parser using working directory `wd` and default
// schemes.
//
// See newParserWithSchemes for more details.
func newParser(wd *url.URL) *parser {
	return newParserWithSchemes(wd, urlfetch.DefaultSchemes)
}

// newParserWithSchemes returns a new PXE parser using working directory `wd`
// and schemes `s`.
//
// If a path encountered in a configuration file is relative instead of a full
// URL, `wd` is used as the "working directory" of that relative path; the
// resulting URL is roughly `wd.String()/path`.
//
// `s` is used to get files referred to by URLs.
func newParserWithSchemes(wd *url.URL, s urlfetch.Schemes) *parser {
	return &parser{
		config: &Config{
			Entries: make(map[string]*boot.LinuxImage),
		},
		scope:   scopeGlobal,
		wd:      wd,
		schemes: s,
	}
}

func parseURL(surl string, wd *url.URL) (*url.URL, error) {
	u, err := url.Parse(surl)
	if err != nil {
		return nil, fmt.Errorf("could not parse URL %q: %v", surl, err)
	}

	if len(u.Scheme) == 0 {
		u.Scheme = wd.Scheme

		if len(u.Host) == 0 {
			// If this is not there, it was likely just a path.
			u.Host = wd.Host
			u.Path = filepath.Join(wd.Path, filepath.Clean(u.Path))
		}
	}
	return u, nil
}

// getFile parses `url` relative to the config's working directory and returns
// an io.Reader for the requested url.
//
// If url is just a relative path and not a full URL, c.wd is used as the
// "working directory" of that relative path; the resulting URL is roughly
// path.Join(wd.String(), url).
func (c *parser) getFile(url string) (io.ReaderAt, error) {
	u, err := parseURL(url, c.wd)
	if err != nil {
		return nil, err
	}

	return c.schemes.LazyFetch(u)
}

// appendFile parses the config file downloaded from `url` and adds it to `c`.
func (c *parser) appendFile(url string) error {
	r, err := c.getFile(url)
	if err != nil {
		return err
	}
	config, err := uio.ReadAll(r)
	if err != nil {
		return err
	}
	log.Printf("Got config file %s:\n%s\n", r, string(config))
	return c.append(string(config))
}

// Append parses `config` and adds the respective configuration to `c`.
func (c *parser) append(config string) error {
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
			c.config.DefaultEntry = arg

		case "include":
			if err := c.appendFile(arg); urlfetch.IsURLError(err) {
				// Means we didn't find the file. Just ignore
				// it.
				// TODO(hugelgupf): plumb a logger through here.
				continue
			} else if err != nil {
				return err
			}

		case "label":
			// We forever enter label scope.
			c.scope = scopeEntry
			c.curEntry = arg
			c.config.Entries[c.curEntry] = &boot.LinuxImage{}
			c.config.Entries[c.curEntry].Cmdline = c.globalAppend

		case "kernel":
			k, err := c.getFile(arg)
			if err != nil {
				return err
			}
			c.config.Entries[c.curEntry].Kernel = k

		case "initrd":
			i, err := c.getFile(arg)
			if err != nil {
				return err
			}
			c.config.Entries[c.curEntry].Initrd = i

		case "append":
			switch c.scope {
			case scopeGlobal:
				c.globalAppend = arg

			case scopeEntry:
				if arg == "-" {
					c.config.Entries[c.curEntry].Cmdline = ""
				} else {
					c.config.Entries[c.curEntry].Cmdline = arg
				}
			}
		}
	}

	// Go through all labels and download the initrds.
	for _, label := range c.config.Entries {
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

			i, err := c.getFile(optkv[1])
			if err != nil {
				return err
			}
			label.Initrd = i
		}
	}

	if len(c.config.DefaultEntry) > 0 {
		if _, ok := c.config.Entries[c.config.DefaultEntry]; !ok {
			return ErrDefaultEntryNotFound
		}
	}
	return nil

}
