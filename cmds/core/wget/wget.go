// Copyright 2012-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
	"github.com/u-root/u-root/pkg/uio"
)

var outPath = flag.String("O", "", "output file")
var errEmptyURL = errors.New("empty url")

type command struct {
	url        string
	outputPath string
}

func newCommand(outPath string, url string) *command {
	return &command{
		outputPath: outPath,
		url:        url,
	}
}

func (c *command) run() error {
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
	flag.Parse()
	if err := newCommand(*outPath, flag.Arg(0)).run(); err != nil {
		if errors.Is(err, errEmptyURL) {
			usage()
		}
		log.Fatal(err)
	}
}
