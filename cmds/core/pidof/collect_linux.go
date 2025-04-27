// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const procPath = "/proc"

func collect(path string) ([]proc, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var procs []proc
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		_, err := strconv.Atoi(name)
		if err != nil {
			continue
		}

		b, err := os.ReadFile(filepath.Join(path, name, "comm"))
		if err != nil {
			continue
		}

		comm := strings.TrimSuffix(string(b), "\n")
		procs = append(procs, proc{comm: comm, pid: name})
	}

	return procs, nil
}
