// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const proc = "/proc"

var errNotFound = errors.New("pid name not found")

func run(stdout io.Writer, proc string, args []string) error {
	entries, err := os.ReadDir(proc)
	if err != nil {
		return err
	}

	var pids []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		_, err := strconv.Atoi(name)
		if err != nil {
			continue
		}

		b, err := os.ReadFile(filepath.Join(proc, name, "comm"))
		if err != nil {
			continue
		}

		comm := strings.TrimSuffix(string(b), "\n")

		for _, arg := range args {
			if comm == arg {
				pids = append(pids, name)
			}
		}
	}

	if len(pids) == 0 {
		return errNotFound
	}

	fmt.Fprintln(stdout, strings.Join(pids, " "))
	return nil
}

func main() {
	if err := run(os.Stdout, proc, os.Args[1:]); err != nil {
		if errors.Is(err, errNotFound) {
			os.Exit(1)
		}
		log.Fatal(err)
	}
}
