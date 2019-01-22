// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// esxiboot executes ESXi kernel over the running kernel.
//
// Synopsis:
//     esxiboot --config <config>
//
// Description:
//     Loads and executes ESXi kernel.
//
// Options:
//     --config=FILE or -c=FILE: set the ESXi config
//
// The config file has the following syntax:
//
// kernel=PATH
// kernelopt=OPTS
// modules=MOD1 [ARGS] --- MOD2 [ARGS] --- ...
//
// Lines starting with '#' are ignored.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/u-root/u-root/pkg/kexec"
	"github.com/u-root/u-root/pkg/multiboot"
)

var cfg = flag.StringP("config", "c", "", "Set the ESXi config")

const (
	kernel  = "kernel"
	args    = "kernelopt"
	modules = "modules"

	comment = '#'
	sep     = "---"
)

type options struct {
	kernel  string
	args    string
	modules []string
}

func parse(fname string) (options, error) {
	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var opt options

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 || line[0] == comment {
			continue
		}

		tokens := strings.SplitN(line, "=", 2)
		if len(tokens) != 2 {
			return opt, fmt.Errorf("bad line %q", line)
		}
		key := strings.TrimSpace(tokens[0])
		val := strings.TrimSpace(tokens[1])
		switch key {
		case kernel:
			opt.kernel = val
		case args:
			opt.args = val
		case modules:
			for _, tok := range strings.Split(val, sep) {
				tok = strings.TrimSpace(tok)
				opt.modules = append(opt.modules, tok)
			}
		}
	}

	err = scanner.Err()
	return opt, err
}

func main() {
	flag.Parse()
	if *cfg == "" {
		log.Fatalf("Config cannot be empty")
	}

	opts, err := parse(*cfg)
	if err != nil {
		log.Fatalf("Cannot parse config %v: %v", *cfg, err)
	}

	p, err := os.Executable()
	if err != nil {
		log.Fatalf("Cannot find current executable path: %v", err)
	}
	trampoline, err := filepath.EvalSymlinks(p)
	if err != nil {
		log.Fatalf("Cannot eval symlinks for %v: %v", p, err)
	}

	m := multiboot.New(opts.kernel, opts.args, trampoline, opts.modules)
	if err := m.Load(false); err != nil {
		log.Fatalf("Load failed: %v", err)
	}

	if err := kexec.Load(m.EntryPoint, m.Segments(), 0); err != nil {
		log.Fatalf("kexec.Load() error: %v", err)
	}
	if err := kexec.Reboot(); err != nil {
		log.Fatalf("kexec.Reboot() error: %v", err)
	}
}
