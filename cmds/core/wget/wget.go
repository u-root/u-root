// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Wget reads one file from a url and writes to stdout.
//
// Synopsis:
//     wget URL
//
// Description:
//     Returns a non-zero code on failure.
//
// Notes:
//     There are a few differences with GNU wget:
//     - Upon error, the return value is always 1.
//     - The protocol (http/https) is mandatory.
//
// Example:
//     wget -O google.txt http://google.com/
package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net/url"
	"os"
	"path"

	"github.com/u-root/u-root/pkg/curl"
	"github.com/u-root/u-root/pkg/uio"
)

var (
	outPath = flag.String("O", "", "output file")
)

func usage() {
	log.Printf("Usage: %s [ARGS] URL\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	log.SetPrefix("wget: ")

	if flag.Parse(); flag.NArg() != 1 {
		usage()
	}

	argURL := flag.Arg(0)
	if argURL == "" {
		log.Fatalln("Empty URL")
	}

	url, err := url.Parse(argURL)
	if err != nil {
		log.Fatalln(err)
	}

	if *outPath == "" {
		if url.Path != "" && url.Path[len(url.Path)-1] != '/' {
			*outPath = path.Base(url.Path)
		} else {
			*outPath = "index.html"
		}
	}

	schemes := curl.Schemes{
		"tftp": curl.DefaultTFTPClient,
		"http": curl.DefaultHTTPClient,

		// curl.DefaultSchemes doesn't support HTTPS by default.
		"https": curl.DefaultHTTPClient,
		"file":  &curl.LocalFileClient{},
	}

	readerAt, err := schemes.Fetch(context.Background(), url)
	if err != nil {
		log.Fatalf("Failed to download %v: %v", argURL, err)
	}

	w, err := os.Create(*outPath)
	if err != nil {
		log.Fatalf("Failed to create output file %q: %v", *outPath, err)
	}
	defer w.Close()

	if _, err := io.Copy(w, uio.Reader(readerAt)); err != nil {
		log.Fatalf("Failed to read response data: %v", err)
	}
}
