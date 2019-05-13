// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ipxe aims to implement a trivial IPXE configuration handler.
package ipxe

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/pxe"
	"github.com/u-root/u-root/pkg/uio"
)

var (
	// ErrNotIpxeScript is returned when the config file is not an
	// ipxe script.
	ErrNotIpxeScript = errors.New("config file is not ipxe as it does not start with #!ipxe")
)

// Config encapsulates a parsed ipxe configuration file.
//
// We currently only support kernel and initrd commands.
type Config struct {
	BootImage *boot.LinuxImage

	schemes pxe.Schemes
}

// NewConfig returns a new IPXE configuration with the file at URL and default
// schemes.
//
// See NewConfigWithSchemes for more details.
func NewConfig(configURL *url.URL) (*Config, error) {
	return NewConfigWithSchemes(configURL, pxe.DefaultSchemes)
}

// NewConfigWithSchemes returns a new IPXE configuration with the file at URL
// and schemes `s`.
//
// `s` is used to get files referred to by URLs in the configuration.
func NewConfigWithSchemes(configURL *url.URL, s pxe.Schemes) (*Config, error) {
	c := &Config{
		schemes: s,
	}
	if err := c.getAndParseFile(configURL); err != nil {
		return nil, err
	}
	return c, nil
}

// getAndParse parses the config file downloaded from `url` and fills in `c`.
func (c *Config) getAndParseFile(u *url.URL) error {
	r, err := c.schemes.LazyGetFile(u)
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
func (c *Config) getFile(surl string) (io.ReaderAt, error) {
	u, err := url.Parse(surl)
	if err != nil {
		return nil, fmt.Errorf("could not parse URL %q: %v", surl, err)
	}

	return c.schemes.LazyGetFile(u)
}

// parseIpxe parses `config` and constructs a BootImage for `c`.
func (c *Config) parseIpxe(config string) error {
	// A trivial ipxe script parser.
	// Currently only supports kernel and initrd commands.
	c.BootImage = &boot.LinuxImage{}

	for _, line := range strings.Split(config, "\n") {
		// Skip blank lines and comment lines.
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '#' {
			continue
		}

		args := strings.Fields(line)
		if len(args) <= 1 {
			log.Printf("Ignoring unsupported ipxe cmd: %s\n", line)
			continue
		}
		cmd := strings.ToLower(args[0])

		switch cmd {
		case "kernel":
			k, err := c.getFile(args[1])
			if err != nil {
				return err
			}
			c.BootImage.Kernel = k

			// Add cmdline if there are any.
			if len(args) > 2 {
				c.BootImage.Cmdline = strings.Join(args[2:], " ")
			}

		case "initrd":
			i, err := c.getFile(args[1])
			if err != nil {
				return err
			}
			c.BootImage.Initrd = i

		default:
			log.Printf("Ignoring unsupported ipxe cmd: %s\n", line)
		}
	}

	return nil
}
