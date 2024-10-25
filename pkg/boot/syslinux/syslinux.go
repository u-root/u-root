// Copyright 2017-2020 the u-root Authors. All rights reserved
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
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/multiboot"
	"github.com/u-root/u-root/pkg/curl"
	"github.com/u-root/uio/uio"
)

func probeIsolinuxFiles() []string {
	files := make([]string, 0, 10)
	// search order from the syslinux wiki
	// http://wiki.syslinux.org/wiki/index.php?title=Config
	dirs := []string{
		"boot/isolinux",
		"isolinux",
		"boot/syslinux",
		"extlinux",
		"syslinux",
		"",
	}
	confs := []string{
		"isolinux.cfg",
		"extlinux.conf",
		"syslinux.cfg",
	}
	for _, dir := range dirs {
		for _, conf := range confs {
			if dir == "" {
				files = append(files, conf)
			} else {
				files = append(files, filepath.Join(dir, conf))
			}
		}
	}
	return files
}

// ParseLocalConfig treats diskDir like a mount point on the local file system
// and finds an isolinux config under there.
func ParseLocalConfig(ctx context.Context, diskDir string) ([]boot.OSImage, error) {
	rootdir := &url.URL{
		Scheme: "file",
		Path:   diskDir,
	}

	for _, relname := range probeIsolinuxFiles() {
		dir, name := filepath.Split(relname)

		// "When booting, the initial working directory for SYSLINUX /
		// ISOLINUX will be the directory containing the initial
		// configuration file."
		//
		// https://wiki.syslinux.org/wiki/index.php?title=Config#Working_directory
		imgs, err := ParseConfigFile(ctx, curl.DefaultSchemes, name, rootdir, dir)
		if curl.IsURLError(err) {
			continue
		}
		return imgs, err
	}
	return nil, fmt.Errorf("no valid syslinux config found on %s", diskDir)
}

// ParseConfigFile parses a Syslinux configuration as specified in
// http://www.syslinux.org/wiki/index.php?title=Config
//
// Currently, only the APPEND, INCLUDE, KERNEL, LABEL, DEFAULT, and INITRD
// directives are partially supported.
//
// `s` is used to fetch any files that must be parsed or provided.
//
// rootdir is the partition mount point that syslinux is operating under.
// Parsed absolute paths will be interpreted relative to the rootdir.
//
// wd is a directory within rootdir that is the current working directory.
// Parsed relative paths will be interpreted relative to rootdir + "/" + wd.
//
// For PXE clients, rootdir will be the the URL without the path, and wd the
// path component of the URL (e.g. rootdir = http://foobar.com, wd =
// barfoo/pxelinux.cfg/).
func ParseConfigFile(ctx context.Context, s curl.Schemes, configFile string, rootdir *url.URL, wd string) ([]boot.OSImage, error) {
	p := newParser(rootdir, wd, s)
	if err := p.appendFile(ctx, configFile); err != nil {
		return nil, err
	}

	// Assign the right label to display to users.
	for label, displayLabel := range p.menuLabel {
		if e, ok := p.linuxEntries[label]; ok {
			e.Name = displayLabel
		}
		if e, ok := p.mbEntries[label]; ok {
			e.Name = displayLabel
		}
	}

	// Intended order:
	//
	// 1. nerfDefaultEntry
	// 2. defaultEntry
	// 3. labels in order they appeared in config
	if len(p.labelOrder) == 0 {
		return nil, nil
	}
	if len(p.defaultEntry) > 0 {
		p.labelOrder = append([]string{p.defaultEntry}, p.labelOrder...)
	}
	if len(p.nerfDefaultEntry) > 0 {
		p.labelOrder = append([]string{p.nerfDefaultEntry}, p.labelOrder...)
	}
	p.labelOrder = dedupStrings(p.labelOrder)

	var images []boot.OSImage
	for _, label := range p.labelOrder {
		if img, ok := p.linuxEntries[label]; ok && img.Kernel != nil {
			images = append(images, img)
		}
		if img, ok := p.mbEntries[label]; ok && img.Kernel != nil {
			images = append(images, img)
		}
	}
	return images, nil
}

func dedupStrings(list []string) []string {
	var newList []string
	seen := make(map[string]struct{})
	for _, s := range list {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			newList = append(newList, s)
		}
	}
	return newList
}

type parser struct {
	// linuxEntries is a map of label name -> label configuration.
	linuxEntries map[string]*boot.LinuxImage
	mbEntries    map[string]*boot.MultibootImage

	// labelOrder is the order of label entries in linuxEntries.
	labelOrder []string

	// menuLabel are human-readable labels defined by the "menu label" directive.
	menuLabel map[string]string

	defaultEntry     string
	nerfDefaultEntry string

	// parser internals.
	globalAppend string
	scope        scope
	curEntry     string
	wd           string
	rootdir      *url.URL
	schemes      curl.Schemes
}

type scope uint8

const (
	scopeGlobal scope = iota
	scopeEntry
)

// newParser returns a new PXE parser using working directory `wd`
// and schemes `s`.
//
// If a path encountered in a configuration file is relative instead of a full
// URL, `wd` is used as the "working directory" of that relative path; the
// resulting URL is roughly `wd.String()/path`.
//
// `s` is used to get files referred to by URLs.
func newParser(rootdir *url.URL, wd string, s curl.Schemes) *parser {
	return &parser{
		linuxEntries: make(map[string]*boot.LinuxImage),
		mbEntries:    make(map[string]*boot.MultibootImage),
		scope:        scopeGlobal,
		wd:           wd,
		rootdir:      rootdir,
		schemes:      s,
		menuLabel:    make(map[string]string),
	}
}

func parseURL(name string, rootdir *url.URL, wd string) (*url.URL, error) {
	u, err := url.Parse(name)
	if err != nil {
		return nil, fmt.Errorf("could not parse URL %q: %w", name, err)
	}

	// If it parsed, but it didn't have a Scheme or Host, use the working
	// directory's values.
	if len(u.Scheme) == 0 && rootdir != nil {
		u.Scheme = rootdir.Scheme

		if len(u.Host) == 0 {
			// If this is not there, it was likely just a path.
			u.Host = rootdir.Host

			// Absolute file names don't get the parent
			// directories, just the host and scheme.
			//
			// "All (paths to) file names inside the configuration
			// file are relative to the Working Directory, unless
			// preceded with a slash."
			//
			// https://wiki.syslinux.org/wiki/index.php?title=Config#Working_directory
			if path.IsAbs(name) {
				u.Path = path.Join(rootdir.Path, path.Clean(u.Path))
			} else {
				u.Path = path.Join(rootdir.Path, wd, path.Clean(u.Path))
			}
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
	u, err := parseURL(url, c.rootdir, c.wd)
	if err != nil {
		return nil, err
	}

	return c.schemes.LazyFetch(u)
}

// getFileWithoutCache gets the file at `url` without caching.
func (c *parser) getFileWithoutCache(surl string) (io.Reader, error) {
	u, err := parseURL(surl, c.rootdir, c.wd)
	if err != nil {
		return nil, fmt.Errorf("could not parse URL %q: %w", surl, err)
	}
	return c.schemes.LazyFetchWithoutCache(u)
}

// appendFile parses the config file downloaded from `url` and adds it to `c`.
func (c *parser) appendFile(ctx context.Context, url string) error {
	u, err := parseURL(url, c.rootdir, c.wd)
	if err != nil {
		return err
	}

	r, err := c.schemes.Fetch(ctx, u)
	if err != nil {
		return err
	}
	config, err := uio.ReadAll(r)
	if err != nil {
		return err
	}
	log.Printf("Got config file %s:\n%s\n", r, string(config))
	return c.append(ctx, string(config))
}

// Append parses `config` and adds the respective configuration to `c`.
func (c *parser) append(ctx context.Context, config string) error {
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
			c.defaultEntry = arg

		case "nerfdefault":
			c.nerfDefaultEntry = arg

		case "include":
			if err := c.appendFile(ctx, arg); curl.IsURLError(err) {
				log.Printf("failed to parse %s: %v", arg, err)
				// Means we didn't find the file. Just ignore
				// it.
				// TODO(hugelgupf): plumb a logger through here.
				continue
			} else if err != nil {
				return err
			}

		case "menu":
			opt := strings.Fields(arg)
			if len(opt) < 1 {
				continue
			}
			switch strings.ToLower(opt[0]) {
			case "label":
				// Note that "menu label" only changes the
				// displayed label, not the identifier for this
				// entry.
				//
				// We track these separately because "menu
				// label" directives may happen before we know
				// whether this is a Linux or Multiboot entry.
				c.menuLabel[c.curEntry] = strings.Join(opt[1:], " ")

			case "default":
				// Are we in label scope?
				//
				// "Only valid after a LABEL statement" -syslinux wiki.
				if c.scope == scopeEntry {
					c.defaultEntry = c.curEntry
				}
			}

		case "label":
			// We forever enter label scope.
			c.scope = scopeEntry
			c.curEntry = arg
			c.linuxEntries[c.curEntry] = &boot.LinuxImage{
				Cmdline: c.globalAppend,
				Name:    c.curEntry,
			}
			c.labelOrder = append(c.labelOrder, c.curEntry)

		case "kernel":
			// I hate special cases like these, but we aren't gonna
			// implement syslinux modules.
			if arg == "mboot.c32" {
				// Prepare for a multiboot kernel.
				delete(c.linuxEntries, c.curEntry)
				c.mbEntries[c.curEntry] = &boot.MultibootImage{
					Name: c.curEntry,
				}
			}
			fallthrough

		case "linux":
			if e, ok := c.linuxEntries[c.curEntry]; ok {
				k, err := c.getFile(arg)
				if err != nil {
					return err
				}
				e.Kernel = k
			}

		case "initrd":
			if e, ok := c.linuxEntries[c.curEntry]; ok {
				// TODO: append "initrd=$arg" to the cmdline.
				//
				// For how this interacts with global appends,
				// read
				// https://wiki.syslinux.org/wiki/index.php?title=Directives/append
				// Multiple initrds are comma-separated
				var initrds []io.Reader
				for _, f := range strings.Split(arg, ",") {
					i, err := c.getFileWithoutCache(f)
					if err != nil {
						return err
					}
					initrds = append(initrds, i)
				}
				e.Initrd = boot.CatInitrdsWithFileCache(initrds...)
			}

		case "fdt":
			// TODO: fdtdir support
			//
			// The logic in u-boot is quite obscure and replies on soc/board names to select the right dtb file.
			// https://gitlab.com/u-boot/u-boot/-/blob/master/boot/pxe_utils.c#L634
			// Can be implemented based on data in /proc/device-tree/compatible

			if e, ok := c.linuxEntries[c.curEntry]; ok {
				dtb, err := c.getFile(arg)
				if err != nil {
					return err
				}
				e.DTB = dtb
			}

		case "append":
			switch c.scope {
			case scopeGlobal:
				c.globalAppend = arg

			case scopeEntry:
				if e, ok := c.mbEntries[c.curEntry]; ok {
					modules := strings.Split(arg, "---")
					// The first module is special -- the kernel.
					if len(modules) > 0 {
						kernel := strings.Fields(modules[0])
						if len(kernel) == 0 {
							return fmt.Errorf("no kernel specified by %v", modules[0])
						}
						k, err := c.getFile(kernel[0])
						if err != nil {
							return err
						}
						e.Kernel = k
						if len(kernel) > 1 {
							e.Cmdline = strings.Join(kernel[1:], " ")
						}
						modules = modules[1:]
					}
					for _, cmdline := range modules {
						m := strings.Fields(cmdline)
						if len(m) == 0 {
							continue
						}
						file, err := c.getFile(m[0])
						if err != nil {
							return err
						}
						e.Modules = append(e.Modules, multiboot.Module{
							Cmdline: strings.TrimSpace(cmdline),
							Module:  file,
						})
					}
				}
				if e, ok := c.linuxEntries[c.curEntry]; ok {
					if arg == "-" {
						e.Cmdline = ""
					} else {
						// Yes, we explicitly _override_, not
						// concatenate. If a specific append
						// directive is present, a global
						// append directive is ignored.
						//
						// Also, "If you enter multiple APPEND
						// statements in a single LABEL entry,
						// only the last one will be used".
						//
						// https://wiki.syslinux.org/wiki/index.php?title=Directives/append
						e.Cmdline = arg
					}
				}
			}
		}
	}

	// Go through all labels and download the initrds.
	for _, label := range c.linuxEntries {
		// If the initrd was set via the INITRD directive, don't
		// overwrite that.
		//
		// TODO(hugelgupf): Is this really what syslinux does? Does
		// INITRD trump cmdline? Does it trump global? What if both the
		// directive and cmdline initrd= are set? Does it depend on the
		// order in the config file? (My current best guess: order.)
		//
		// Answer: Normally, the INITRD directive appends to the
		// cmdline, and the _last_ effective initrd= parameter is used
		// for loading initrd files.
		if label.Initrd != nil {
			continue
		}

		for _, opt := range strings.Fields(label.Cmdline) {
			optkv := strings.Split(opt, "=")
			if len(optkv) != 2 || optkv[0] != "initrd" {
				continue
			}

			i, err := c.getFile(optkv[1])
			if err != nil {
				return err
			}
			label.Initrd = i
		}
	}
	return nil
}
