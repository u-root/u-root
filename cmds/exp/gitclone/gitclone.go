// Copyright 2013-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style license that can be found in the LICENSE file.

// gitclone clones a git repository into a specified directory

package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/src-d/go-git.v4"
)

func main() {
	if len(os.Args) <= 2 {
		log.Fatalf("Usage: gitclone https://github.com/u-root/u-root.git u-root")
	}

	url := os.Args[1]
	dir := os.Args[2]

	fmt.Printf("Cloning '%s' into '%s'...\n", url, dir)
	if _, err := git.PlainClone(dir, false, &git.CloneOptions{URL: url, Progress: os.Stdout}); err != nil {
		log.Fatalf("%v", err)
	}
}
