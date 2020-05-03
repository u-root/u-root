// Copyright 2016-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// chmod changes mode bits (e.g. permissions) of a file.
//
// Synopsis:
//     chmod MODE FILE...
//
// Desription:
//     MODE is a three character octal value or a string like a=rwx
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const special = 99999

var (
	recursive bool
	reference string
)

func init() {
	flag.BoolVar(&recursive,
		"R",
		false,
		"do changes recursively")

	flag.BoolVar(&recursive,
		"recursive",
		false,
		"do changes recursively")

	flag.StringVar(&reference,
		"reference",
		"",
		"use mode from reference file")
}

func changeMode(path string, mode os.FileMode, octval uint64, mask uint64) (err error) {
	// A special value for mask means the mode is fully described
	if mask == special {
		return os.Chmod(path, mode)
	}

	var info os.FileInfo
	info, err = os.Stat(path)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	mode = info.Mode() & os.FileMode(mask)
	mode = mode | os.FileMode(octval)

	return os.Chmod(path, mode)
}

func calculateMode(modeString string) (mode os.FileMode, octval uint64, mask uint64) {
	var err error
	octval, err = strconv.ParseUint(modeString, 8, 32)
	if err == nil {
		if octval > 0777 {
			log.Fatalf("Invalid octal value %0o. Value should be less than or equal to 0777.", octval)
		}
		// a fully described octal mode was supplied, signal that with a special value for mask
		mask = special
		mode = os.FileMode(octval)
		return
	}

	reMode := regexp.MustCompile("^([ugoa]+)([-+=])(.*)")
	m := reMode.FindStringSubmatch(modeString)
	// Test for mode strings with invalid characters.
	// This can't be done in the first regexp: if the match for m[3] is restricted to [rwx]*,
	// `a=9` and `a=` would be indistinguishable: m[3] would be empty.
	// `a=` is a valid (but destructive) operation. Do not turn a typo into that.
	reMode = regexp.MustCompile("^[rwx]*$")
	if len(m) < 3 || !reMode.MatchString(m[3]) {
		log.Fatalf("Unable to decode mode %q. Please use an octal value or a valid mode string.", modeString)
	}

	// m[3] is [rwx]{0,3}
	var octvalDigit uint64
	if strings.Contains(m[3], "r") {
		octvalDigit += 4
	}
	if strings.Contains(m[3], "w") {
		octvalDigit += 2
	}
	if strings.Contains(m[3], "x") {
		octvalDigit++
	}

	// m[2] is [-+=]
	var operator = m[2]

	// Use a mask so that we do not overwrite permissions for a user/group that was not specified
	mask = 0777

	// For "-", invert octvalDigit before applying the mask
	if operator == "-" {
		octvalDigit = 7 - octvalDigit
	}

	// m[1] is [ugoa]+
	if strings.Contains(m[1], "o") || strings.Contains(m[1], "a") {
		octval += octvalDigit
		mask = mask & 0770
	}
	if strings.Contains(m[1], "g") || strings.Contains(m[1], "a") {
		octval += octvalDigit << 3
		mask = mask & 0707
	}
	if strings.Contains(m[1], "u") || strings.Contains(m[1], "a") {
		octval += octvalDigit << 6
		mask = mask & 0077
	}

	// For "+" the mask is superfluous, reset it
	if operator == "+" {
		mask = 0777
	}

	// The mode is fully described, signal that with a special value for mask
	if operator == "=" && strings.Contains(m[1], "a") {
		mask = special
		mode = os.FileMode(octval)
	}
	return
}

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Fprintf(os.Stderr, "Usage of %s: [mode] filepath\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	if len(flag.Args()) < 2 && reference == "" {
		fmt.Fprintf(os.Stderr, "Usage of %s: [mode] filepath\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	var mode os.FileMode
	var octval, mask uint64
	var fileList []string

	if reference != "" {
		fi, err := os.Stat(reference)
		if err != nil {
			log.Fatalf("bad reference file: %v", err)

		}
		mask = special
		mode = fi.Mode()
		fileList = flag.Args()
	} else {
		mode, octval, mask = calculateMode(flag.Args()[0])
		fileList = flag.Args()[1:]
	}

	var exitError bool
	for _, name := range fileList {
		if recursive {
			err := filepath.Walk(name, func(path string,
				info os.FileInfo,
				err error) error {
				return changeMode(path, mode, octval, mask)
			})
			if err != nil {
				log.Printf("%v", err)
				exitError = true
			}
		} else {
			if err := changeMode(name, mode, octval, mask); err != nil {
				log.Printf("%v", err)
				exitError = true
			}
		}
	}
	if exitError {
		os.Exit(1)
	}
}
