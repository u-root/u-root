// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Run with `go run checklicenses.go`. This script has one drawback:
// - It does not correct the licenses; it simply outputs a list of files which
//   do not conform and returns 1 if the list is non-empty.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var absPath = flag.Bool("a", false, "Print absolute paths")

const uroot = "$GOPATH/src/github.com/u-root/u-root"

// The first few lines of every go file is expected to contain this license.
var license = regexp.MustCompile(
	`^// Copyright [\d\-, ]+ the u-root Authors\. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file\.
`)

type rule struct {
	*regexp.Regexp
	invert bool
}

func accept(s string) rule {
	return rule{
		regexp.MustCompile("^" + s + "$"),
		false,
	}
}

func reject(s string) rule {
	return rule{
		regexp.MustCompile("^" + s + "$"),
		true,
	}
}

// A file is checked iff all the accepts and none of the rejects match.
var rules = []rule{
	accept(`.*\.go`),
	reject(`/vendor/.*`),      // Various authors
	reject(`/cmds/dhcp/.*`),   // Graham King
	reject(`/cmds/ectool/.*`), // Chromium authors
	reject(`/cmds/ldd/.*`),    // Go authors
	reject(`/cmds/ping/.*`),   // Go authors
	reject(`/netlink/.*`),     // Docker (Apache license)
}

func main() {
	flag.Parse()
	uroot := os.ExpandEnv(uroot)
	incorrect := []string{}

	// Walk u-root tree.
	err := filepath.Walk(uroot, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		// Test rules
		trimmedPath := strings.TrimPrefix(path, uroot)
		for _, r := range rules {
			if r.MatchString(trimmedPath) == r.invert {
				return nil
			}
		}
		// Read from the file.
		r, err := os.Open(path)
		if err != nil {
			return err
		}
		defer r.Close()
		contents, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}
		if !license.Match(contents) {
			p := trimmedPath
			if *absPath {
				p = path
			}
			incorrect = append(incorrect, p)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Print files with incorrect licenses.
	if len(incorrect) > 0 {
		fmt.Println(strings.Join(incorrect, "\n"))
		os.Exit(1)
	}
}
