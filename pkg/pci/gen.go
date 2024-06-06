// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"io/ioutil"
	"log"
	"os"
)

var (
	pciidspath = [...]string{"/usr/share/misc/pci.ids"}
	code       = `// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci
var pciids = []byte(
`
)

func main() {
	for _, p := range pciidspath {
		b, err := ioutil.ReadFile(p)
		if err != nil {
			log.Printf("%v did not work, keep looking", err)
			continue
		}
		code = code + "`" + string(b) + "`)\n"
		err = ioutil.WriteFile("pciids.go", []byte(code), 0644)
		if err != nil {
			log.Fatalf("Unable to write pciids.go: %v", err)
		}
		os.Exit(0)
	}
	log.Fatalf("Unable to find a pci.ids file in these paths: %v", pciidspath)
}
