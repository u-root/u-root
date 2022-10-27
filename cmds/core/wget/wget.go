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

func usage() {
	log.Printf("Usage: %s [ARGS] URL\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func run(arg string) error {
	log.SetPrefix("wget: ")

	if arg == "" {
		return errEmptyURL
	}

	parsedURL, err := url.Parse(arg)
	if err != nil {
		return err
	}

	if *outPath == "" {
		*outPath = defaultOutputPath(parsedURL.Path)
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
		return fmt.Errorf("failed to download %v: %w", arg, err)
	}

	if err := uio.ReadIntoFile(reader, *outPath); err != nil {
		return err
	}

	return nil
}

func defaultOutputPath(urlPath string) string {
	if urlPath == "" || strings.HasSuffix(urlPath, "/") {
		return "index.html"
	}
	return path.Base(urlPath)
}

func main() {
	flag.Parse()
	if err := run(flag.Arg(0)); err != nil {
		if errors.Is(err, errEmptyURL) {
			usage()
		}
		log.Fatal(err)
	}
}
