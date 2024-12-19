// Copyright 2014-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	guid "github.com/google/uuid"
	"github.com/u-root/u-root/pkg/efivarfs"
)

var (
	flist   = flag.Bool("list", false, "List all efivars")
	fread   = flag.String("read", "", "Read specified efivar. Variable must be of form -read Name-UUID")
	fdelete = flag.String("delete", "", "Delete specified efivar. Variable must be of form -delete Name-UUID")
	fwrite  = flag.String("write", "", "Write to specified efivar. Variable must be of form -write Name-UUID OR Name\n"+
		"In the later case a UUID is being generated\n"+
		"This command is used with -content to specify the data being written to the efivar.")
	fcontent = flag.String("content", "", "Path to file to write to efivar. Used with -write e.g. -write Foo -content bar.json")
)

func main() {
	flag.Parse()

	if err := runpath(os.Stdout, efivarfs.DefaultVarFS, *flist, *fread, *fdelete, *fwrite, *fcontent); err != nil {
		log.Fatalf("Operation failed: %v", err)
	}
}

func runpath(out io.Writer, p string, list bool, read, remove, write, content string) error {
	e, err := efivarfs.NewPath(p)
	if err != nil {
		return err
	}

	return run(out, e, list, read, remove, write, content)
}

func run(out io.Writer, e efivarfs.EFIVar, list bool, read, remove, write, content string) error {
	if list {
		l, err := efivarfs.SimpleListVariables(e)
		if err != nil {
			return fmt.Errorf("list failed: %w", err)
		}
		for _, s := range l {
			fmt.Fprintln(out, s)
		}
	}

	if read != "" {
		attr, data, err := efivarfs.SimpleReadVariable(e, read)
		if err != nil {
			return fmt.Errorf("read failed: %w", err)
		}
		b, err := io.ReadAll(data)
		if err != nil {
			return fmt.Errorf("reading buffer failed: %w", err)
		}
		fmt.Fprintf(out, "Name: %s, Attributes: %d, Data: %s", read, attr, b)
	}

	if remove != "" {
		if err := efivarfs.SimpleRemoveVariable(e, remove); err != nil {
			return fmt.Errorf("delete failed: %w", err)
		}
	}

	if write != "" {
		if strings.ContainsAny(write, "-") {
			v := strings.SplitN(write, "-", 2)
			if _, err := guid.Parse(v[1]); err != nil {
				return fmt.Errorf("%q malformed: Must be either Name-GUID or just Name: %w", v[1], os.ErrInvalid)
			}
		}
		path, err := filepath.Abs(content)
		if err != nil {
			return fmt.Errorf("could not resolve path: %w", err)
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		if !strings.ContainsAny(write, "-") {
			write = write + "-" + guid.New().String()
		}
		if err = efivarfs.SimpleWriteVariable(e, write, 7, bytes.NewBuffer(b)); err != nil {
			return fmt.Errorf("write failed: %w", err)
		}
	}
	return nil
}
