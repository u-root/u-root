// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ls

import (
	"flag"
	"fmt"
	"strings"

	"github.com/u-root/u-root/pkg/ls"
)

// addOSSpecificFlags adds OS-specific flags to the flag set.
func (c *command) addOSSpecificFlags(fs *flag.FlagSet, f *flags) {
	fs.BoolVar(&f.final, "p", false, "Print only the final path element of each file name")
}

func (c *command) printFile(stringer ls.Stringer, f file, flags flags) {
	if f.err != nil {
		fmt.Fprintln(c.Stdout, f.err)
		return
	}
	// Hide .files unless -a was given
	if flags.all || !strings.HasPrefix(f.lsfi.Name, ".") {
		// Unless they said -p, we always print the full path
		if !flags.final {
			f.lsfi.Name = f.path
		}
		if flags.classify {
			f.lsfi.Name = f.lsfi.Name + indicator(f.lsfi)
		}
		fmt.Fprintln(c.Stdout, stringer.FileString(f.lsfi))
	}
}
