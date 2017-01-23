// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Fmap parses flash maps.
//
// Synopsis:
//     fmap [-s|-r i] [FILE]
//
// Description:
//     Return 0 if the flash map is valid and 1 otherwise. Detailed information
//     is printed to stderr. If FILE is not specified, read from stdin.
//
//     This implementation is based off of https://github.com/dhendrix/flashmap.
//
// Options:
//     -s: print human readable summary
//     -r i: read an area from the flash
package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"text/template"

	fmap "github.com/u-root/u-root/cmds/fmap/lib"
)

var (
	summary = flag.Bool("s", false, "print human readable summary")
	read    = flag.Int("r", -1, "read an area from the flash")
)

// Print human readable summary of the fmap.
func printFMap(f *fmap.FMap, m *fmap.FMapMetadata) {
	const desc = `Fmap found at {{printf "%#x" .Metadata.Start}}:
	Signature:  {{printf "%s" .Signature}}
	VerMajor:   {{.VerMajor}}
	VerMinor:   {{.VerMinor}}
	Base:       {{printf "%#x" .Base}}
	Size:       {{printf "%#x" .Size}}
	Name:       {{printf "%s" .Name}}
	NAreas:     {{.NAreas}}
{{- range $i, $v := .Areas}}
	Areas[{{$i}}]:
		Offset:  {{printf "%#x" $v.Offset}}
		Size:    {{printf "%#x" $v.Size}}
		Name:    {{printf "%s" $v.Name}}
		Flags:   {{printf "%#x" $v.Flags}} ({{FlagNames $v.Flags}})
{{- end}}
`
	t := template.Must(template.New("desc").
		Funcs(template.FuncMap{"FlagNames": fmap.FlagNames}).
		Parse(desc))
	// Combine the two structs to pass into template.
	combined := struct {
		*fmap.FMap
		Metadata *fmap.FMapMetadata
	}{f, m}
	if err := t.Execute(os.Stdout, combined); err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	if *summary && *read >= 0 {
		log.Fatal("Both flags cannot be used at once")
	}

	// Choose a reader
	r := os.Stdin
	if flag.NArg() == 1 {
		var err error
		r, err = os.Open(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		defer r.Close()
	} else if flag.NArg() > 1 {
		log.Fatal("Too many arguments")
	}

	// Read fmap.
	f, metadata, err := fmap.ReadFMap(r)
	if err != nil {
		log.Fatal(err)
	}
	// Optionally print area.
	if *read >= 0 {
		areaReader, err := f.ReadArea(r, *read)
		if err != nil {
			log.Fatal(err)
		}
		_, err = bufio.NewWriter(os.Stdout).ReadFrom(areaReader)
		if err != nil {
			log.Fatal(err)
		}
	}
	// Optionally print summar.
	if *summary {
		printFMap(f, metadata)
	}
}
