// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ipxe implements a trivial IPXE config file parser.
package ipxe

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/curl"
	"github.com/u-root/u-root/pkg/ulog"
	"github.com/u-root/uio/uio"
)

// ErrNotIpxeScript is returned when the config file is not an
// ipxe script.
var ErrNotIpxeScript = errors.New("config file is not ipxe as it does not start with #!ipxe")

// parser encapsulates a parsed ipxe configuration file.
//
// We currently only support kernel and initrd commands.
type parser struct {
	bootImage *boot.LinuxImage

	// wd is the current working directory.
	//
	// Relative file paths are interpreted relative to this URL.
	wd *url.URL

	log ulog.Logger

	schemes curl.Schemes
}

// ParseConfig returns a new configuration with the file at URL and default
// schemes.
//
// `s` is used to get files referred to by URLs in the configuration.
func ParseConfig(ctx context.Context, l ulog.Logger, configURL *url.URL, s curl.Schemes) (*boot.LinuxImage, error) {
	c := &parser{
		schemes: s,
		log:     l,
	}
	if err := c.getAndParseFile(ctx, configURL); err != nil {
		return nil, err
	}
	return c.bootImage, nil
}

// getAndParse parses the config file downloaded from `url` and fills in `c`.
func (c *parser) getAndParseFile(ctx context.Context, u *url.URL) error {
	r, err := c.schemes.Fetch(ctx, u)
	if err != nil {
		return err
	}
	data, err := uio.ReadAll(r)
	if err != nil {
		return err
	}
	config := string(data)
	if !strings.HasPrefix(config, "#!ipxe") {
		return ErrNotIpxeScript
	}
	c.log.Printf("Got ipxe config file %s:\n%s\n", r, config)

	// Parent dir of the config file.
	c.wd = &url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
		Path:   path.Dir(u.Path),
	}
	return c.parseIpxe(config)
}

// getFile parses `surl` and returns an io.Reader for the requested url.
func (c *parser) getFile(surl string) (io.ReaderAt, error) {
	u, err := parseURL(surl, c.wd)
	if err != nil {
		return nil, fmt.Errorf("could not parse URL %q: %w", surl, err)
	}
	// Cache content read from http body into a tmpfs file, other
	// than in heap. This cuts down ram consumption and help boot
	// on board with low ram config.
	return uio.NewLazyOpenerAt(surl, func() (io.ReaderAt, error) {
		f, err := os.CreateTemp("", "cache-kernel")
		if err != nil {
			return nil, err
		}
		defer f.Close()
		r, err := c.schemes.LazyFetchWithoutCache(u)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(f, r)
		if err != nil {
			return nil, err
		}
		if err := f.Sync(); err != nil {
			return nil, err
		}
		// Return a read-only copy.
		readOnlyF, err := os.Open(f.Name())
		if err != nil {
			return nil, err
		}
		return readOnlyF, nil
	}), nil
}

func (c *parser) getFileWithoutCache(surl string) (io.Reader, error) {
	u, err := parseURL(surl, c.wd)
	if err != nil {
		return nil, fmt.Errorf("could not parse URL %q: %w", surl, err)
	}
	return c.schemes.LazyFetchWithoutCache(u)
}

func parseURL(name string, wd *url.URL) (*url.URL, error) {
	u, err := url.Parse(name)
	if err != nil {
		return nil, fmt.Errorf("could not parse URL %q: %w", name, err)
	}

	// If it parsed, but it didn't have a Scheme or Host, use the working
	// directory's values.
	if len(u.Scheme) == 0 && wd != nil {
		u.Scheme = wd.Scheme

		if len(u.Host) == 0 {
			// If this is not there, it was likely just a path.
			u.Host = wd.Host

			// Absolute file names don't get the parent
			// directories, just the host and scheme.
			if !path.IsAbs(name) {
				u.Path = path.Join(wd.Path, path.Clean(u.Path))
			}
		}
	}
	return u, nil
}

func (c *parser) createInitrd(initrds []io.Reader) {
	if len(initrds) > 0 {
		c.bootImage.Initrd = boot.CatInitrdsWithFileCache(initrds...)
	}
}

// parseIpxe parses `config` and constructs a BootImage for `c`.
func (c *parser) parseIpxe(config string) error {
	// A trivial ipxe script parser.
	// Currently only supports kernel and initrd commands.
	c.bootImage = &boot.LinuxImage{}

	var initrds []io.Reader
	for _, line := range strings.Split(config, "\n") {
		// Skip blank lines and comment lines.
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '#' {
			continue
		}

		args := strings.Fields(line)
		if len(args) == 0 {
			continue
		}
		cmd := strings.ToLower(args[0])

		switch cmd {
		case "kernel":
			if len(args) > 1 {
				k, err := c.getFile(args[1])
				if err != nil {
					return err
				}
				c.bootImage.Kernel = k
			}

			// Add cmdline if there are any.
			if len(args) > 2 {
				c.bootImage.Cmdline = strings.Join(args[2:], " ")
			}

		case "initrd":
			if len(args) > 1 {
				for _, f := range strings.Split(args[1], ",") {
					i, err := c.getFileWithoutCache(f)
					if err != nil {
						return err
					}
					initrds = append(initrds, i)
				}
			}

		case "boot":
			// Stop parsing at this point, we should go ahead and
			// boot.
			c.createInitrd(initrds)
			return nil

		default:
			c.log.Printf("Ignoring unsupported ipxe cmd: %s", line)
		}
	}

	// EOF - we should go ahead and boot.
	c.createInitrd(initrds)
	return nil
}
