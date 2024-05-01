// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// lsmod list currently loaded Linux kernel modules.
//
// Synopsis:
//
//	lsmod
//
// Description:
//
//	lsmod is a clone of lsmod(8)
//
// Author:
//
//	Roland Kammerer <dev.rck@gmail.com>
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func run(stdout io.Writer, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintln(stdout, "Module                  Size  Used by")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), " ")
		name, size, used, usedBy := s[0], s[1], s[2], s[3]
		final := fmt.Sprintf("%-19s %8s  %s", name, size, used)
		if usedBy != "-" {
			usedBy = usedBy[:len(usedBy)-1]
			final += fmt.Sprintf(" %s", usedBy)
		}
		fmt.Fprintln(stdout, final)
	}

	return scanner.Err()
}

func main() {
	if err := run(os.Stdout, "/proc/modules"); err != nil {
		log.Fatal(err)
	}
}
