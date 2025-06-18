// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !windows && !plan9

package ls

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	humanize "github.com/dustin/go-humanize"
	"golang.org/x/sys/unix"
)

// FileString implements Stringer.FileString.
func (ls LongStringer) FileString(fi FileInfo) string {
	// Golang's FileMode.String() is almost sufficient, except we would
	// rather use b and c for devices.
	replacer := strings.NewReplacer("Dc", "c", "D", "b")

	// Ex: crw-rw-rw-  root  root  1, 3  Feb 6 09:31  null
	pattern := "%[1]s\t%[2]s\t%[3]s\t%[4]d, %[5]d\t%[7]v\t%[8]s"
	if fi.Mode&os.ModeDevice == 0 && fi.Mode&os.ModeCharDevice == 0 {
		// Ex: -rw-rw----  myuser  myuser  1256  Feb 6 09:31  recipes.txt
		pattern = "%[1]s\t%[2]s\t%[3]s\t%[6]s\t%[7]v\t%[8]s"
	}

	var size string
	if ls.Human {
		size = humanize.Bytes(uint64(fi.Size))
	} else {
		size = strconv.FormatInt(fi.Size, 10)
	}

	s := fmt.Sprintf(pattern,
		replacer.Replace(fi.Mode.String()),
		lookupUserName(fi.UID),
		lookupGroupName(fi.GID),
		unix.Major(fi.Rdev),
		unix.Minor(fi.Rdev),
		size,
		fi.MTime.Format("Jan _2 15:04"),
		ls.Name.FileString(fi))

	if fi.Mode&os.ModeType == os.ModeSymlink {
		s += fmt.Sprintf(" -> %v", fi.SymlinkTarget)
	}
	return s
}
