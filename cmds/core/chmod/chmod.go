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
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/u-root/u-root/pkg/uroot/util"
)

const (
	special = 99999
	usage   = "chmod: chmod [mode] filepath"
)

var errBadUsage = errors.New(usage)

func init() {
	flag.Usage = util.Usage(flag.Usage, usage)
}

func changeMode(path string, mode os.FileMode, octval uint64, mask uint64) (fs.FileMode, error) {
	// A special value for mask means the mode is fully described
	if mask == special {
		if err := os.Chmod(path, mode); err != nil {
			return 0, err
		}
		return mode, nil
	}

	var info os.FileInfo
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	mode = info.Mode() & os.FileMode(mask)
	mode = mode | os.FileMode(octval)

	if err := os.Chmod(path, mode); err != nil {
		return 0, err
	}
	return mode, nil
}

func calculateMode(modeString string) (mode os.FileMode, octval uint64, mask uint64, err error) {
	octval, err = strconv.ParseUint(modeString, 8, 32)
	if err == nil {
		if octval > 0o777 {
			return mode, octval, mask, fmt.Errorf("%w: invalid octal value %0o. Value should be less than or equal to 0777", strconv.ErrRange, octval)
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
		return mode, octval, mask, fmt.Errorf("%w:unable to decode mode %q. Please use an octal value or a valid mode string", strconv.ErrSyntax, modeString)
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
	operator := m[2]

	// Use a mask so that we do not overwrite permissions for a user/group that was not specified
	mask = 0o777

	// For "-", invert octvalDigit before applying the mask
	if operator == "-" {
		octvalDigit = 7 - octvalDigit
	}

	// m[1] is [ugoa]+
	if strings.Contains(m[1], "o") || strings.Contains(m[1], "a") {
		octval += octvalDigit
		mask = mask & 0o770
	}
	if strings.Contains(m[1], "g") || strings.Contains(m[1], "a") {
		octval += octvalDigit << 3
		mask = mask & 0o707
	}
	if strings.Contains(m[1], "u") || strings.Contains(m[1], "a") {
		octval += octvalDigit << 6
		mask = mask & 0o077
	}

	// For "+" the mask is superfluous, reset it
	if operator == "+" {
		mask = 0o777
	}

	// The mode is fully described, signal that with a special value for mask
	if operator == "=" && strings.Contains(m[1], "a") {
		mask = special
		mode = os.FileMode(octval)
	}
	return mode, octval, mask, nil
}

func chmod(recursive bool, reference string, args ...string) (fs.FileMode, error) {
	var mode os.FileMode
	if len(args) < 1 {
		return mode, errBadUsage
	}

	if len(args) < 2 && reference == "" {
		return mode, errBadUsage
	}

	var (
		err          error
		octval, mask uint64
		fileList     []string
	)

	if reference != "" {
		fi, err := os.Stat(reference)
		if err != nil {
			return 0, fmt.Errorf("bad reference file: %w", err)
		}
		mask = special
		mode = fi.Mode()
		fileList = args
	} else {
		var err error
		if mode, octval, mask, err = calculateMode(args[0]); err != nil {
			return mode, err
		}
		fileList = args[1:]
	}

	for _, name := range fileList {
		if recursive {
			err := filepath.Walk(name, func(path string, info os.FileInfo, err error) error {
				mode, err = changeMode(path, mode, octval, mask)
				return err
			})
			if err != nil {
				return mode, err
			}
		} else {
			if mode, err = changeMode(name, mode, octval, mask); err != nil {
				return mode, err
			}
		}
	}
	return mode, err
}

func main() {
	var (
		recursive = flag.Bool("recursive", false, "do changes recursively")
		reference = flag.String("reference", "", "use mode from reference file")
	)
	flag.Parse()
	if _, err := chmod(*recursive, *reference, flag.Args()...); err != nil {
		log.Fatal(err)
	}
}
