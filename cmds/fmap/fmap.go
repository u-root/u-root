// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Fmap parses flash maps.
//
// Synopsis:
//     fmap [-s] [FILE]
//
// Description:
//     Return 0 if the flash map is valid and 1 otherwise. Detailed information
//     is printed to stderr. If FILE is not specified, read from stdin.
//
//     This implementation is based off of https://github.com/dhendrix/flashmap.
//
// Options:
//     -s:  print human readable summary
package main

import (
	"flag"
	"log"
	"os"
	"text/template"

	fmap "github.com/u-root/u-root/cmds/fmap/lib"
)

var (
	summary = flag.Bool("s", false, "print human readable summary")
)

// Print human readable summary of the fmap.
func printFMap(f *fmap.FMap) {
	const desc = `Fmap found at {{printf "%#x" .Start}}:
	Signature:  {{printf "%s" .Signature}}
	VerMajor:   {{.VerMajor}}
	VerMinor:   {{.VerMinor}}
	Base:       {{printf "%#x" .Base}}
	Size:       {{printf "%#x" .Size}}
	Name:       {{printf "%s" .Name}}
	NAreas:     {{len .Areas}}
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
	if err := t.Execute(os.Stdout, f); err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	// Choose a reader
	r := os.Stdin
	if flag.NArg() == 1 {
		var err error
		r, err = os.Open(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
	} else if flag.NArg() > 1 {
		log.Fatal("Too many arguments")
	}

	// Read fmap and optionally print summary.
	f := fmap.ReadFMap(r)
	if *summary {
		printFMap(f)
	}
}
