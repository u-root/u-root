// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// List modules currently loaded in the Linux kernel
//
// Synopsis:
//	lsmod
//
// Description:
//	lsmod is a clone of lsmod(8)
//
// Author:
//     Roland Kammerer <dev.rck@gmail.com>
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("/proc/modules")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fmt.Println("Module                  Size  Used by")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), " ")
		name, size, used, usedBy := s[0], s[1], s[2], s[3]
		final := fmt.Sprintf("%-19s %8s  %s", name, size, used)
		if usedBy != "-" {
			usedBy = usedBy[:len(usedBy)-1]
			final += fmt.Sprintf(" %s", usedBy)
		}
		fmt.Println(final)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
