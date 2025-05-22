// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// lsfd - list file descriptors
//
// Synopsis:
//
//	lsfd [PID]...

package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	proc   = "/proc"
	format = "%-20s %-5s %s\n"
)

// record store information fetched from map_files, ns and fd folders
// under /proc/${PID}, for now it's only filename
type record struct {
	FileName string
}

func processPID(pidPath string) (string, []record, error) {
	b, err := os.ReadFile(filepath.Join(pidPath, "comm"))
	if err != nil {
		return "", nil, err
	}

	comm := strings.TrimSuffix(string(b), "\n")

	exe, err := os.Readlink(filepath.Join(pidPath, "exe"))
	if err != nil {
		return "", nil, err
	}

	cwd, err := os.Readlink(filepath.Join(pidPath, "cwd"))
	if err != nil {
		return "", nil, err
	}

	root, err := os.Readlink(filepath.Join(pidPath, "root"))
	if err != nil {
		return "", nil, err
	}

	records := []record{{FileName: exe}, {FileName: cwd}, {FileName: root}}
	records = append(records, traverse(filepath.Join(pidPath, "ns"))...)
	records = append(records, traverse(filepath.Join(pidPath, "fd"))...)
	records = append(records, traverse(filepath.Join(pidPath, "map_files"))...)

	return comm, records, nil
}

func run(stdout io.Writer, path string, pids []string) error {
	var headerPrinted bool

	// read all files
	if len(pids) == 0 {
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		for _, e := range entries {
			if !e.IsDir() {
				continue
			}

			_, err := strconv.Atoi(e.Name())
			if err != nil {
				continue
			}

			pid := e.Name()
			comm, records, err := processPID(filepath.Join(proc, pid))
			if err != nil {
				continue
			}

			for _, r := range records {
				if !headerPrinted {
					fmt.Fprintf(stdout, format, "COMMAND NAME", "PID", "NAME")
					headerPrinted = true
				}
				fmt.Fprintf(stdout, format, comm, pid, r.FileName)
			}
		}
	} else {
		for _, pid := range pids {
			_, err := strconv.Atoi(pid)
			if err != nil {
				continue
			}

			comm, records, err := processPID(filepath.Join(proc, pid))
			if err != nil {
				continue
			}

			for _, r := range records {
				if !headerPrinted {
					fmt.Fprintf(stdout, format, "COMMAND NAME", "PID", "NAME")
					headerPrinted = true
				}
				fmt.Fprintf(stdout, format, comm, pid, r.FileName)
			}
		}
	}

	return nil
}

func traverse(path string) []record {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil
	}

	var records []record
	for _, e := range entries {
		if e.Type() == fs.ModeSymlink {
			n, err := os.Readlink(filepath.Join(path, e.Name()))
			if err != nil {
				return nil
			}
			records = append(records, record{FileName: n})
		}
	}

	return records
}

func main() {
	flag.Parse()
	if err := run(os.Stdout, proc, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
