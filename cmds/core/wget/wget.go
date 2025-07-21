// Copyright 2012-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

// Wget reads one file from a url and writes to stdout.
//
// Synopsis:
//
//	wget URL
//
// Description:
//
//	Returns a non-zero code on failure.
//
// Notes:
//
//	There are a few differences with GNU wget:
//	- Upon error, the return value is always 1.
//	- The protocol (http/https) is mandatory.
//
// Example:
//
//	wget -O google.txt http://google.com/
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/u-root/u-root/pkg/curl"
	"github.com/u-root/uio/uio"

	// To build the dependencies of this package with TinyGo, we need to include
	// the cpuid package, since tinygo does not support the asm code in the
	// cpuid package. The cpuid package will use the tinygo bridge to get the
	// CPU information. For further information see
	// github.com/u-root/cpuid/cpuid_amd64_tinygo_bridge.go
	_ "github.com/u-root/cpuid"
)

var errEmptyURL = errors.New("empty url")

type cmd struct {
	url        string
	outputPath string
}

// flags parses wget flags
// wget is old school, and allows flags after the URL.
// This code does not process the -- flag specified in the
// man page, as the command itself does not seem to either.
func flags(args ...string) (string, string, error) {
	// -- takes priority over everything else.
	// flag package does not allow - as a flag.
	// except, in spite of the docs, wget on linux seems
	// to ignore --. Great. It's a terrible idea anyway,
	// unless you really like creating files that start
	// with -
	// If at some point one wishes to add -- support,
	// the slices package is a good place to start.

	if len(args) == 0 {
		return "", "", errEmptyURL
	}

	f := flag.NewFlagSet(args[0], flag.ContinueOnError)
	outPath := f.String("O", "", "output file")

	if err := f.Parse(args[1:]); err != nil {
		return "", "", err
	}

	if len(f.Args()) == 0 {
		return "", "", errEmptyURL
	}

	URL := f.Args()[0]

	// Now, it is allowed to have switches after the URL,
	// handle following flags
	if err := f.Parse(f.Args()[1:]); err != nil {
		return "", "", err
	}

	return *outPath, URL, nil
}

func command(args ...string) (*cmd, error) {
	outPath, URL, err := flags(args...)
	if err != nil {
		return nil, err
	}

	return &cmd{
		outputPath: outPath,
		url:        URL,
	}, nil
}

func (c *cmd) run() error {
	log.SetPrefix("wget: ")

	if c.url == "" {
		return errEmptyURL
	}

	parsedURL, err := url.Parse(c.url)
	if err != nil {
		return err
	}

	if c.outputPath == "" {
		c.outputPath = defaultOutputPath(parsedURL.Path)
	}
	if c.outputPath == "-" {
		c.outputPath = "/dev/stdout"
	}

	schemes := curl.Schemes{
		"tftp": curl.DefaultTFTPClient,
		"http": curl.DefaultHTTPClient,

		// curl.DefaultSchemes doesn't support HTTPS by default.
		"https": curl.DefaultHTTPClient,
		"file":  &curl.LocalFileClient{},
	}

	reader, err := schemes.FetchWithoutCache(context.Background(), parsedURL)
	if err != nil {
		return fmt.Errorf("failed to download %v: %w", c.url, err)
	}

	if err := uio.ReadIntoFile(reader, c.outputPath); err != nil {
		return err
	}

	return nil
}

func usage() {
	log.Printf("Usage: %s [ARGS] URL\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func defaultOutputPath(urlPath string) string {
	if urlPath == "" || strings.HasSuffix(urlPath, "/") {
		return "index.html"
	}
	return path.Base(urlPath)
}

func main() {
	c, err := command(os.Args...)
	if err == nil {
		err = c.run()
	}
	if err != nil {
		if errors.Is(err, errEmptyURL) {
			usage()
		}
		log.Fatal(err)
	}
}
