// Copyright 2013 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
)

func loadVTOC(name string) (*os.File, []file, error) {
	var l int64
	f, err := os.Open(name)
	if err != nil {
		return nil, nil, err
	}
	r := io.LimitReader(f, 8)
	_, err = fmt.Fscanf(r, "%x", &l)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: can't scan vtoc length: %v", name, err)
	}

	r = io.LimitReader(f, l)
	dec := json.NewDecoder(r)

	var vtoc []file
	if err := dec.Decode(&vtoc); err != nil {
		return nil, nil, fmt.Errorf("%s: can't decode: %v", name, err)
	}
	return f, vtoc, nil
}

func buildVTOC(dirs []string) ([]*file, error) {
	var vtoc []*file
	for _, v := range dirs {
		debug("Process %v", v)
		err := filepath.Walk(v, func(name string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			debug("visit %v", name)
			var s syscall.Stat_t
			if err := syscall.Lstat(name, &s); err != nil {
				return fmt.Errorf("%s: %v", name, err)
			}
			f := &file{
				Name:    name,
				Mode:    fi.Mode(),
				ModTime: fi.ModTime(),
				IsDir:   fi.IsDir(),
				Uid:     int(s.Uid),
				Gid:     int(s.Gid),
			}
			switch f.Mode.String()[0] {
			case '-':
				f.Size = fi.Size()
			case 'L':
				f.Link, err = os.Readlink(name)
			case 'D':
				f.Dev = s.Rdev
			}
			vtoc = append(vtoc, f)
			return err
		})
		if err != nil {
			return nil, err
		}
	}

	return vtoc, nil
}

func encodeVTOC(vtoc []*file) (int64, error) {
	return -1, errors.New("not yet")
}

func writeVTOC(f io.Writer, vtoc []*file) (int, error) {
	var outvtoc = make([]*file, len(vtoc))
	// Do a little cleanup
	for i := range vtoc {
		// The path names may be absolute but they will also have been cleaned.
		// If the first char is /, it won't be.
		// N.B. This assumes Unix-style file names. But I don't care about Windows.
		v := *vtoc[i]
		outvtoc[i] = &v
		if outvtoc[i].Name[0] == '/' {
			outvtoc[i].Name = outvtoc[i].Name[1:]
		}
	}

	var v bytes.Buffer
	enc := json.NewEncoder(&v)
	if err := enc.Encode(outvtoc); err != nil {
		return -1, fmt.Errorf("Encoding files: %v", err)
	}

	if _, err := fmt.Fprintf(f, "%07x\n", v.Len()); err != nil {
		return -1, err
	}
	debug("vtoc size is %d", v.Len())
	return f.Write(v.Bytes())
}

func NewVTOCEncoder(opts ...VTOCOpt) error {
	var v vtoc

	for _, o := range opts {
		if err := o(&v); err != nil {
			return err
		}
	}
	return nil
}
