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
	"errors"
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
	"github.com/u-root/u-root/pkg/uio"
)

var (
	// ErrDefaultEntryNotFound is returned when the configuration file
	// names a default label that is not part of the configuration.
	ErrDefaultEntryNotFound = errors.New("default variable not set in configuration")
	// ErrInitrdUsedWithoutLinux is returned when an initrd directive is
	// not following a linux directive in the same menu entry
	ErrInitrdUsedWithoutLinux = errors.New("missing linux directive before initrd")
	// ErrModuleUsedWithoutMultiboot is returned when a module directive is
	// not following a multiboot directive in the same menu entry
	ErrModuleUsedWithoutMultiboot = errors.New("missing multiboot directive before module")
)

// Config encapsulates a parsed grub configuration file.
type Config struct {
	// Entries is a map of label name -> label configuration.
	Entries map[string]boot.OSImage

	// DefaultEntry is the default label key to use.
	//
	// If DefaultEntry is non-empty, the label is guaranteed to exist in
	// `Entries`.
	DefaultEntry string
}

// ParseConfigFile parses a grub configuration as specified in
// https://www.gnu.org/software/grub/manual/grub/
//
// Currently, only the linux[16|efi], initrd[16|efi], menuentry and set
// directives are partially supported.
//
// curl.DefaultSchemes is used to fetch any files that must be parsed or
// provided.
//
// `wd` is the default scheme, host, and path for any files named as a
// relative path - e.g. kernel, include, and initramfs paths are requested
// relative to the wd.
func ParseConfigFile(url string, wd *url.URL) (*Config, error) {
	return ParseConfigFileWithSchemes(curl.DefaultSchemes, url, wd)
}

// ParseConfigFileWithSchemes is like ParseConfigFile, but uses the given
// schemes explicitly.
func ParseConfigFileWithSchemes(s curl.Schemes, url string, wd *url.URL) (*Config, error) {
	p := newParserWithSchemes(wd, s)
	if err := p.appendFile(url); err != nil {
		return nil, err
	}
	return p.config, nil
}

type parser struct {
	config *Config
	W      io.Writer

	// parser internals.
	scope    scope
	numEntry int
	curEntry string
	curLabel string
	wd       *url.URL
	schemes  curl.Schemes
}

type scope uint8

const (
	scopeGlobal scope = iota
	scopeEntry
)

// newParserWithSchemes returns a new grub parser using working directory `wd`
// and schemes `s`.
//
// If a path encountered in a configuration file is relative instead of a full
// URL, `wd` is used as the "working directory" of that relative path; the
// resulting URL is roughly `wd.String()/path`.
//
// `s` is used to get files referred to by URLs.
func newParserWithSchemes(wd *url.URL, s curl.Schemes) *parser {
	return &parser{
		config: &Config{
			Entries: make(map[string]boot.OSImage),
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
	if len(config) > 500 {
		// Avoid flooding the console on real systems
		// TODO: do we want to pass a verbose flag or a logger?
		log.Printf("Got config file %s", r)
	} else {
		log.Printf("Got config file %s:\n%s\n", r, string(config))
	}
	return c.append(string(config))
}

func isWhitespace(b byte) bool {
	return b == '\t' || b == '\n' || b == '\v' ||
		b == '\f' || b == '\r' || b == ' '
}

type quote uint8

const (
	unquoted quote = iota
	escape
	singleQuote
	doubleQuote
	doubleQuoteEscape
	comment
)

// Fields splits a grub line and unquote it's components according to
// https://www.gnu.org/software/grub/manual/grub/grub.html#Quoting
// except that the escaping of newline is not supported
func fields(s string) []string {
	var ret []string
	var token []byte

	var context quote
	lastWhiteSpace := true
	for i := range []byte(s) {
		quotes := context != unquoted
		switch context {
		case unquoted:
			switch s[i] {
			case '\\':
				context = escape
				// strip out the quote
				continue
			case '\'':
				context = singleQuote
				// strip out the quote
				continue
			case '"':
				context = doubleQuote
				// strip out the quote
				continue
			case '#':
				if lastWhiteSpace {
					context = comment
					// strip out the rest
					continue
				}
			}

		case escape:
			context = unquoted

		case singleQuote:
			if s[i] == '\'' {
				context = unquoted
				// strip out the quote
				continue
			}

		case doubleQuote:
			switch s[i] {
			case '\\':
				context = doubleQuoteEscape
				// strip out the quote
				continue
			case '"':
				context = unquoted
				// strip out the quote
				continue
			}

		case doubleQuoteEscape:
			switch s[i] {
			case '$', '"', '\\', '\n': // or newline
			default:
				token = append(token, '\\')
			}

			context = doubleQuote

		case comment:
			// should end on newline

			// strip out the rest
			continue

		}

		lastWhiteSpace = isWhitespace(s[i])

		if !isWhitespace(s[i]) || quotes {
			token = append(token, s[i])
		} else if len(token) > 0 {
			ret = append(ret, string(token))
			token = token[:0]
		}
	}

	if len(token) > 0 {
		ret = append(ret, string(token))
	}
	return ret
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

// Append parses `config` and adds the respective configuration to `c`.
func (c *parser) append(config string) error {
	// Here's a shitty parser.
	for _, line := range strings.Split(config, "\n") {
		kv := fields(line)
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
					c.config.DefaultEntry = vals[1]
				}
			}

		case "configfile":
			// TODO test that
			if err := c.appendFile(arg); err != nil {
				return err
			}

		case "menuentry":
			c.scope = scopeEntry
			c.curEntry = strconv.Itoa(c.numEntry)
			c.curLabel = arg
			c.numEntry++

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
			c.config.Entries[c.curEntry] = entry
			c.config.Entries[c.curLabel] = entry

		case "initrd", "initrd16", "initrdefi":
			i, err := c.getFile(arg)
			if err != nil {
				return err
			}
			entry, ok := c.config.Entries[c.curEntry].(*boot.LinuxImage)
			if !ok {
				return ErrInitrdUsedWithoutLinux
			}
			entry.Initrd = i

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
			c.config.Entries[c.curEntry] = entry
			c.config.Entries[c.curLabel] = entry

		case "module":
			// TODO handle --nounzip arguments ? (change parsing)
			m, err := c.getFile(arg)
			if err != nil {
				return err
			}
			entry, ok := c.config.Entries[c.curEntry].(*boot.MultibootImage)
			if !ok {
				return ErrModuleUsedWithoutMultiboot
			}
			// TODO: Lasy tryGzipFilter(m)
			mod := multiboot.Module{
				Module:  m,
				Name:    arg,
				CmdLine: cmdlineQuote(kv[2:]),
			}
			entry.Modules = append(entry.Modules, mod)

		}
	}

	if len(c.config.DefaultEntry) > 0 {
		if _, ok := c.config.Entries[c.config.DefaultEntry]; !ok {
			return ErrDefaultEntryNotFound
		}
	}
	return nil

}
