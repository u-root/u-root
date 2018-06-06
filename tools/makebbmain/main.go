// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// makebbmain creates a bb main.go source file.
package main

import (
	"flag"
	"go/importer"
	"log"

	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/uroot"
)

var (
	template  = flag.String("template", "github.com/u-root/u-root/pkg/bb/cmd", "bb main.go template package")
	outputDir = flag.String("o", "", "output directory")
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		log.Fatalf("must list bb packages as arguments")
	}

	env := golang.Default()
	importer := importer.For("source", nil)

	if err := uroot.CreateBBMainSource(env, importer, *template, flag.Args(), *outputDir); err != nil {
		log.Fatalf("failed to create bb source file: %v", err)
	}
}
