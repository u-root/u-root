// Copyright 2014-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

	if err := run(*flist, *fread, *fdelete, *fwrite, *fcontent); err != nil {
		log.Fatalf("Operation failed: %v", err)
	}
}

func run(list bool, read, delete, write, content string) error {
	if list {
		l, err := efivarfs.SimpleListVariables()
		if err != nil {
			return fmt.Errorf("list failed: %v", err)
		}
		for _, s := range l {
			log.Println(s)
		}
	}

	if read != "" {
		attr, data, err := efivarfs.SimpleReadVariable(read)
		if err != nil {
			return fmt.Errorf("read failed: %v", err)
		}
		b, err := io.ReadAll(data)
		if err != nil {
			return fmt.Errorf("reading buffer failed: %v", err)
		}
		log.Printf("Name: %s, Attributes: %d, Data: %s", read, attr, b)
	}

	if delete != "" {
		if err := efivarfs.SimpleRemoveVariable(delete); err != nil {
			return fmt.Errorf("delete failed: %v", err)
		}
	}

	if write != "" {
		if strings.ContainsAny(write, "-") {
			v := strings.SplitN(write, "-", 2)
			if _, err := guid.Parse(v[1]); err != nil {
				return fmt.Errorf("var name malformed: Must be either Name-GUID or just Name")
			}
		}
		path, err := filepath.Abs(content)
		if err != nil {
			return fmt.Errorf("could not resolve path: %v", err)
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file: %v", err)
		}
		if !strings.ContainsAny(write, "-") {
			write = write + "-" + guid.New().String()
		}
		if err = efivarfs.SimpleWriteVariable(write, 7, bytes.NewBuffer(b)); err != nil {
			return fmt.Errorf("write failed: %v", err)
		}
	}
	return nil
}
