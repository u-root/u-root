// Copyright 2017-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package grub implements a grub config file parser.
//
// See the grub manual https://www.gnu.org/software/grub/manual/grub/ for
// a reference of the configuration format
// In particular the following pages:
// - https://www.gnu.org/software/grub/manual/grub/html_node/Shell_002dlike-scripting.html
// - https://www.gnu.org/software/grub/manual/grub/html_node/Commands.html
//
// Currently, only the linux[16|efi], initrd[16|efi], menuentry and set
// directives are partially supported.
package grub

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/multiboot"
	"github.com/u-root/u-root/pkg/curl"
	"github.com/u-root/u-root/pkg/shlex"
	"github.com/u-root/u-root/pkg/uio"
)

var probeGrubFiles = []string{
	"boot/grub/grub.cfg",
	"grub/grub.cfg",
	"grub2/grub.cfg",
	"boot/grub2/grub.cfg",
}

// ParseLocalConfig looks for a GRUB config in the disk partition mounted at
// diskDir and parses out OSes to boot.
//
// This... is at best crude, at worst totally wrong, since we fundamentally
// assume that the kernels we boot are only on this one partition. But so is
// this whole parser.
func ParseLocalConfig(ctx context.Context, diskDir string) ([]boot.OSImage, error) {
	wd := &url.URL{
		Scheme: "file",
		Path:   diskDir,
	}

	// This is a hack. GRUB should stop caring about URLs at least in the
	// way we use them today, because GRUB has additional path encoding
	// methods. Sorry.
	//
	// Normally, stuff like this will be in EFI/BOOT/grub.cfg, but some
	// distro's have their own directory in this EFI namespace. Just check
	// 'em all.
	files, err := filepath.Glob(filepath.Join(diskDir, "EFI", "*", "grub.cfg"))
	if err != nil {
		log.Printf("[grub] Could not glob for %s/EFI/*/grub.cfg: %v", diskDir, err)
	}
	var relNames []string
	for _, file := range files {
		base, err := filepath.Rel(diskDir, file)
		if err == nil {
			relNames = append(relNames, base)
		}
	}

	for _, relname := range append(relNames, probeGrubFiles...) {
		c, err := ParseConfigFile(ctx, curl.DefaultSchemes, relname, wd)
		if curl.IsURLError(err) {
			continue
		}
		return c, err
	}
	return nil, fmt.Errorf("no valid grub config found")
}

// ParseConfigFile parses a grub configuration as specified in
// https://www.gnu.org/software/grub/manual/grub/
//
// Currently, only the linux[16|efi], initrd[16|efi], menuentry and set
// directives are partially supported.
//
// `wd` is the default scheme, host, and path for any files named as a
// relative path - e.g. kernel, include, and initramfs paths are requested
// relative to the wd.
func ParseConfigFile(ctx context.Context, s curl.Schemes, configFile string, wd *url.URL) ([]boot.OSImage, error) {
	p := newParser(wd, s)
	if err := p.appendFile(ctx, configFile); err != nil {
		return nil, err
	}

	// Don't add entries twice.
	//
	// Multiple labels can refer to the same image, so we have to dedup by pointer.
	seenLinux := make(map[*boot.LinuxImage]struct{})
	seenMB := make(map[*boot.MultibootImage]struct{})

	if len(p.defaultEntry) > 0 {
		p.labelOrder = append([]string{p.defaultEntry}, p.labelOrder...)
	}

	var images []boot.OSImage
	for _, label := range p.labelOrder {
		if img, ok := p.linuxEntries[label]; ok {
			if _, ok := seenLinux[img]; !ok {
				images = append(images, img)
				seenLinux[img] = struct{}{}
			}
		}

		if img, ok := p.mbEntries[label]; ok {
			if _, ok := seenMB[img]; !ok {
				images = append(images, img)
				seenMB[img] = struct{}{}
			}
		}
	}
	return images, nil
}

type parser struct {
	linuxEntries map[string]*boot.LinuxImage
	mbEntries    map[string]*boot.MultibootImage

	labelOrder   []string
	defaultEntry string

	W io.Writer

	// parser internals.
	numEntry int

	// curEntry is the current entry number as a string.
	curEntry string

	// curLabel is the last parsed label from a "menuentry".
	curLabel string

	wd      *url.URL
	schemes curl.Schemes
}

// newParser returns a new grub parser using working directory `wd`
// and schemes `s`.
//
// If a path encountered in a configuration file is relative instead of a full
// URL, `wd` is used as the "working directory" of that relative path; the
// resulting URL is roughly `wd.String()/path`.
//
// `s` is used to get files referred to by URLs.
func newParser(wd *url.URL, s curl.Schemes) *parser {
	return &parser{
		linuxEntries: make(map[string]*boot.LinuxImage),
		mbEntries:    make(map[string]*boot.MultibootImage),
		wd:           wd,
		schemes:      s,
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
func (c *parser) appendFile(ctx context.Context, url string) error {
	u, err := parseURL(url, c.wd)
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
	if len(config) > 500 {
		// Avoid flooding the console on real systems
		// TODO: do we want to pass a verbose flag or a logger?
		log.Printf("[grub] Got config file %s", r)
	} else {
		log.Printf("[grub] Got config file %s:\n%s\n", r, string(config))
	}
	return c.append(ctx, string(config))
}

// CmdlineQuote quotes the command line as grub-core/lib/cmdline.c does
func cmdlineQuote(args []string) string {
	q := make([]string, len(args))
	for i, s := range args {
		s = strings.Replace(s, `\`, `\\`, -1)
		s = strings.Replace(s, `'`, `\'`, -1)
		s = strings.Replace(s, `"`, `\"`, -1)
		if strings.ContainsRune(s, ' ') {
			s = `"` + s + `"`
		}
		q[i] = s
	}
	return strings.Join(q, " ")
}

// append parses `config` and adds the respective configuration to `c`.
//
// NOTE: This parser has outlived its usefulness already, given that it doesn't
// even understand the {} scoping in GRUB. But let's get the tests to pass, and
// then we can do a rewrite.
func (c *parser) append(ctx context.Context, config string) error {
	// Here's a shitty parser.
	for _, line := range strings.Split(config, "\n") {
		kv := shlex.Argv(line)
		if len(kv) < 1 {
			continue
		}
		directive := strings.ToLower(kv[0])
		// Used by tests (allow no parameters here)
		if c.W != nil && directive == "echo" {
			fmt.Fprintf(c.W, "echo:%#v\n", kv[1:])
		}

		if len(kv) <= 1 {
			continue
		}
		arg := kv[1]

		switch directive {
		case "set":
			vals := strings.SplitN(arg, "=", 2)
			if len(vals) == 2 {
				//TODO handle vars? bootVars[vals[0]] = vals[1]
				//log.Printf("grubvar: %s=%s", vals[0], vals[1])
				if vals[0] == "default" {
					c.defaultEntry = vals[1]
				}
			}

		case "configfile":
			// TODO test that
			if err := c.appendFile(ctx, arg); err != nil {
				return err
			}

		case "menuentry":
			c.curEntry = strconv.Itoa(c.numEntry)
			c.curLabel = arg
			c.numEntry++
			c.labelOrder = append(c.labelOrder, c.curEntry, c.curLabel)

		case "linux", "linux16", "linuxefi":
			k, err := c.getFile(arg)
			if err != nil {
				return err
			}
			// from grub manual: "Any initrd must be reloaded after using this command" so we can replace the entry
			entry := &boot.LinuxImage{
				Name:    c.curLabel,
				Kernel:  k,
				Cmdline: cmdlineQuote(kv[2:]),
			}
			c.linuxEntries[c.curEntry] = entry
			c.linuxEntries[c.curLabel] = entry

		case "initrd", "initrd16", "initrdefi":
			if e, ok := c.linuxEntries[c.curEntry]; ok {
				i, err := c.getFile(arg)
				if err != nil {
					return err
				}
				e.Initrd = i
			}

		case "multiboot":
			// TODO handle --quirk-* arguments ? (change parsing)
			k, err := c.getFile(arg)
			if err != nil {
				return err
			}
			// from grub manual: "Any initrd must be reloaded after using this command" so we can replace the entry
			entry := &boot.MultibootImage{
				Name:    c.curLabel,
				Kernel:  k,
				Cmdline: cmdlineQuote(kv[2:]),
			}
			c.mbEntries[c.curEntry] = entry
			c.mbEntries[c.curLabel] = entry

		case "module":
			// TODO handle --nounzip arguments ? (change parsing)
			if e, ok := c.mbEntries[c.curEntry]; ok {
				// The only allowed arg
				cmdline := kv[1:]
				if arg == "--nounzip" {
					arg = kv[2]
					cmdline = kv[2:]
				}

				m, err := c.getFile(arg)
				if err != nil {
					return err
				}
				// TODO: Lasy tryGzipFilter(m)
				mod := multiboot.Module{
					Module:  m,
					Cmdline: cmdlineQuote(cmdline),
				}
				e.Modules = append(e.Modules, mod)
			}
		}
	}
	return nil

}
