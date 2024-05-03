// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netcat

import (
	"fmt"
	"io"
	"log"
)

// Logf logs a message if the verbose flag is set.
func Logf(nc *NetcatConfig, format string, args ...interface{}) {
	if nc.Output.Verbose {
		log.Printf(LOG_PREFIX+format, args...)
	}
}

// Logf logs a message if the verbose flag is set.
// The output is written to the provided writer instead of os.stderr
func FLogf(nc *NetcatConfig, w io.Writer, format string, args ...interface{}) {
	if nc.Output.Verbose {
		fmt.Fprintf(w, LOG_PREFIX+format, args...)
	}
}
