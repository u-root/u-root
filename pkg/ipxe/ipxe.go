// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ipxe implements a trivial IPXE config file parser.
package ipxe

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/uio"
	"github.com/u-root/u-root/pkg/urlfetch"
)

var (
	// ErrNotIpxeScript is returned when the config file is not an
	// ipxe script.
	ErrNotIpxeScript = errors.New("config file is not ipxe as it does not start with #!ipxe")
)

// parser encapsulates a parsed ipxe configuration file.
//
// We currently only support kernel and initrd commands.
type parser struct {
	bootImage *boot.LinuxImage

	schemes urlfetch.Schemes
}

// ParseConfig returns a new  configuration with the file at URL and default
// schemes.
//
// See ParseConfigWithSchemes for more details.
func ParseConfig(configURL *url.URL) (*boot.LinuxImage, error) {
	return ParseConfigWithSchemes(configURL, urlfetch.DefaultSchemes)
}

// ParseConfigWithSchemes returns a new  configuration with the file at URL
// and schemes `s`.
//
// `s` is used to get files referred to by URLs in the configuration.
func ParseConfigWithSchemes(configURL *url.URL, s urlfetch.Schemes) (*boot.LinuxImage, error) {
	c := &parser{
		schemes: s,
	}
	if err := c.getAndParseFile(configURL); err != nil {
		return nil, err
	}
	return c.bootImage, nil
}

// getAndParse parses the config file downloaded from `url` and fills in `c`.
func (c *parser) getAndParseFile(u *url.URL) error {
	r, err := c.schemes.LazyFetch(u)
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
	log.Printf("Got ipxe config file %s:\n%s\n", r, config)
	return c.parseIpxe(config)
}

// getFile parses `surl` and returns an io.Reader for the requested url.
func (c *parser) getFile(surl string) (io.ReaderAt, error) {
	u, err := url.Parse(surl)
	if err != nil {
		return nil, fmt.Errorf("could not parse URL %q: %v", surl, err)
	}
	return c.schemes.LazyFetch(u)
}

// parseIpxe parses `config` and constructs a BootImage for `c`.
func (c *parser) parseIpxe(config string) error {
	// A trivial ipxe script parser.
	// Currently only supports kernel and initrd commands.
	c.bootImage = &boot.LinuxImage{}

	for _, line := range strings.Split(config, "\n") {
		// Skip blank lines and comment lines.
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '#' {
			continue
		}

		args := strings.Fields(line)
		if len(args) <= 1 {
			log.Printf("Ignoring unsupported ipxe cmd: %s", line)
			continue
		}
		cmd := strings.ToLower(args[0])

		switch cmd {
		case "kernel":
			k, err := c.getFile(args[1])
			if err != nil {
				return err
			}
			c.bootImage.Kernel = k

			// Add cmdline if there are any.
			if len(args) > 2 {
				c.bootImage.Cmdline = strings.Join(args[2:], " ")
			}

		case "initrd":
			i, err := c.getFile(args[1])
			if err != nil {
				return err
			}
			c.bootImage.Initrd = i

		default:
			log.Printf("Ignoring unsupported ipxe cmd: %s", line)
		}
	}

	return nil
}
