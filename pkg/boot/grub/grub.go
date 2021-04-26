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
// See parser.append function for list of commands that are supported.
package grub

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/multiboot"
	"github.com/u-root/u-root/pkg/curl"
	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/mount/block"
	"github.com/u-root/u-root/pkg/shlex"
	"github.com/u-root/u-root/pkg/uio"
)

var probeGrubFiles = []string{
	"boot/grub/grub.cfg",
	"grub/grub.cfg",
	"grub2/grub.cfg",
	"boot/grub2/grub.cfg",
}

// Grub syntax for OpenSUSE/Fedora/RHEL has some undocumented quirks. You
// won't find it on the master branch, but instead look at the rhel and fedora
// branches for these commits:
//
// * https://github.com/rhboot/grub2/commit/7e6775e6d4a8de9baf3f4676d4e021cc2f5dd761
// * https://github.com/rhboot/grub2/commit/0c26c6f7525737962d1389ebdfbb918f52d1b3b6
//
// They add a special case to not escape hex sequences:
//
//     grub> echo hello \xff \xfg
//     hello \xff xfg
//
// Their default installations depend on this functionality.
var hexEscape = regexp.MustCompile(`\\x[0-9a-fA-F]{2}`)
var anyEscape = regexp.MustCompile(`\\.{0,3}`)

// mountFlags are the flags this grub interpreter uses to mount partitions.
var mountFlags = uintptr(mount.ReadOnly)

// absFileScheme creates a file:/// scheme with an absolute path. Technically,
// file schemes must be absolute paths and Go makes that assumption.
func absFileScheme(path string) (*url.URL, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return &url.URL{
		Scheme: "file",
		Path:   path,
	}, nil
}

// ParseLocalConfig looks for a GRUB config in the disk partition mounted at
// diskDir and parses out OSes to boot.
func ParseLocalConfig(ctx context.Context, diskDir string, devices block.BlockDevices, mountPool *mount.Pool) ([]boot.OSImage, error) {
	root, err := absFileScheme(diskDir)
	if err != nil {
		return nil, err
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
		c, err := ParseConfigFile(ctx, curl.DefaultSchemes, relname, root, devices, mountPool)
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
// See parser.append function for list of commands that are supported.
//
// `root` is the default scheme, host, and path for any files named as a
// relative path - e.g. kernel and initramfs paths are requested relative to
// the root.
func ParseConfigFile(ctx context.Context, s curl.Schemes, configFile string, root *url.URL, devices block.BlockDevices, mountPool *mount.Pool) ([]boot.OSImage, error) {
	p := newParser(root, devices, mountPool, s)
	if err := p.appendFile(ctx, configFile); err != nil {
		return nil, err
	}

	// Don't add entries twice.
	//
	// Multiple labels can refer to the same image, so we have to dedup by pointer.
	seenLinux := make(map[*boot.LinuxImage]struct{})
	seenMB := make(map[*boot.MultibootImage]struct{})

	if defaultEntry, ok := p.variables["default"]; ok {
		p.labelOrder = append([]string{defaultEntry}, p.labelOrder...)
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

	labelOrder []string

	W io.Writer

	// parser internals.
	numEntry int
	// Special variables:
	//   * default: Default boot option.
	//   * root: Root "partition" as a URL.
	variables map[string]string

	// curEntry is the current entry number as a string.
	curEntry string

	// curLabel is the last parsed label from a "menuentry".
	curLabel string

	devices   block.BlockDevices
	mountPool *mount.Pool
	schemes   curl.Schemes
}

// newParser returns a new grub parser using `root` and schemes `s`.
//
// We are going off script here by using URLs instead of grub's device syntax.
//
// Typically, the default value for root should be the mount point containing
// the grub config, for example: "file:///tmp/sda1/". Kernel and initramfs
// files are opened relative to this path.
//
// Some grub configs may set a different local root. For this, all partitions
// must be mounted beforehand and made available to grub through `mounts`.
//
// For example, if the grub config contains `search --by-label LINUX`, this
// resolves to the device node "/dev/disk/by-partlabel/LINUX". This grub parser
// looks through mounts for a matching device number.
func newParser(root *url.URL, devices block.BlockDevices, mountPool *mount.Pool, s curl.Schemes) *parser {
	return &parser{
		linuxEntries: make(map[string]*boot.LinuxImage),
		mbEntries:    make(map[string]*boot.MultibootImage),
		variables: map[string]string{
			"root": root.String(),
		},
		devices:   devices,
		mountPool: mountPool,
		schemes:   s,
	}
}

func parseURL(surl string, root string) (*url.URL, error) {
	u, err := url.Parse(surl)
	if err != nil {
		return nil, fmt.Errorf("could not parse URL %q: %v", surl, err)
	}
	ru, err := url.Parse(root)
	if err != nil {
		return nil, fmt.Errorf("could not parse URL %q: %v", root, err)
	}

	if len(u.Scheme) == 0 {
		u.Scheme = ru.Scheme

		if len(u.Host) == 0 {
			// If this is not there, it was likely just a path.
			u.Host = ru.Host
			u.Path = filepath.Join(ru.Path, filepath.Clean(u.Path))
		}
	}
	return u, nil
}

// getFile parses `url` relative to the current root and returns an io.Reader
// for the requested url.
//
// If url is just a relative path and not a full URL, c.root is used for the
// relative path; the resulting URL is roughly path.Join(root, url).
func (c *parser) getFile(url string) (io.ReaderAt, error) {
	u, err := parseURL(url, c.variables["root"])
	if err != nil {
		return nil, err
	}

	return c.schemes.LazyFetch(u)
}

// appendFile parses the config file downloaded from `url` and adds it to `c`.
func (c *parser) appendFile(ctx context.Context, url string) error {
	u, err := parseURL(url, c.variables["root"])
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
		// Replace \ with \\ unless it matches \xXX
		s = anyEscape.ReplaceAllStringFunc(s, func(match string) string {
			if hexEscape.MatchString(match) {
				return match
			}
			return strings.Replace(match, `\`, `\\`, -1)
		})
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
		// Add extra backslash for OpenSUSE/Fedora/RHEL use case. shlex
		// will convert it back to a single backslash.
		line = hexEscape.ReplaceAllString(line, `\\$0`)
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
		case "search.file", "search.fs_label", "search.fs_uuid":
			// Alias to regular search directive.
			kv = append(
				[]string{"search", map[string]string{
					"search.file":     "--file",
					"search.fs_label": "--fs-label",
					"search.fs_uuid":  "--fs-uuid",
				}[directive]},
				kv[1:]...,
			)
			fallthrough
		case "search":
			// Parses a line with this format:
			//   search [--file|--label|--fs-uuid] [--set [var]] [--no-floppy] name
			fs := pflag.NewFlagSet("grub.search", pflag.ContinueOnError)
			searchUUID := fs.BoolP("fs-uuid", "u", false, "")
			searchLabel := fs.BoolP("fs-label", "l", false, "")
			searchFile := fs.BoolP("file", "f", false, "")
			setVar := fs.String("set", "root", "")
			// Ignored flags
			fs.String("no-floppy", "", "ignored")
			fs.String("hint", "", "ignored")
			fs.SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
				// Everything that begins with "hint" is ignored.
				if strings.HasPrefix(name, "hint") {
					name = "hint"
				}
				return pflag.NormalizedName(name)
			})

			if err := fs.Parse(kv[1:]); err != nil || fs.NArg() != 1 {
				log.Printf("Warning: Grub parser could not parse %q", kv)
				continue
			}
			searchName := fs.Arg(0)
			if *searchUUID && *searchLabel || *searchUUID && *searchFile || *searchLabel && *searchFile {
				log.Printf("Warning: Grub parser found more than one search option in %q, skipping line", line)
				continue
			}
			if !*searchUUID && !*searchLabel && !*searchFile {
				// defaults to searchUUID
				*searchUUID = true
			}

			switch {
			case *searchUUID:
				d := c.devices.FilterFSUUID(searchName)
				if len(d) != 1 {
					log.Printf("Error: Expected 1 device with UUID %q, found %d", searchName, len(d))
					continue
				}
				mp, err := c.mountPool.Mount(d[0], mountFlags)
				if err != nil {
					log.Printf("Error: Could not mount %v: %v", d[0], err)
					continue
				}
				setVal, err := absFileScheme(mp.Path)
				if err != nil {
					continue
				}
				c.variables[*setVar] = setVal.String()
			case *searchLabel:
				d, err := c.devices.FilterPartLabel(searchName)
				if err != nil {
					log.Printf("Error: Could not search label %q: %v", searchName, err)
					continue
				}
				if len(d) != 1 {
					log.Printf("Error: Expected 1 device with label %q, found %d", searchName, len(d))
					continue
				}
				mp, err := c.mountPool.Mount(d[0], mountFlags)
				if err != nil {
					log.Printf("Error: Could not mount %v: %v", d[0], err)
					continue
				}
				setVal, err := absFileScheme(mp.Path)
				if err != nil {
					continue
				}
				c.variables[*setVar] = setVal.String()
			case *searchFile:
				// Make sure searchName stays in mountpoint. Remove "../" components.
				cleanPath, err := filepath.Rel("/", filepath.Clean(filepath.Join("/", searchName)))
				if err != nil {
					log.Printf("Error: Could not clean path %q: %v", searchName, err)
					continue
				}
				// Search through all the devices for the file.
				for _, d := range c.devices {
					mp, err := c.mountPool.Mount(d, mountFlags)
					if err != nil {
						log.Printf("Warning: Could not mount %v: %v", mp, err)
						continue
					}
					file := filepath.Join(mp.Path, cleanPath)
					if _, err := os.Stat(file); err == nil {
						setVal, err := absFileScheme(mp.Path)
						if err != nil {
							continue
						}
						c.variables[*setVar] = setVal.String()
						break
					}
				}
			}

		case "set":
			vals := strings.SplitN(arg, "=", 2)
			if len(vals) == 2 {
				// TODO: We cannot parse grub device syntax.
				if vals[0] == "root" {
					continue
				}
				c.variables[vals[0]] = vals[1]
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
