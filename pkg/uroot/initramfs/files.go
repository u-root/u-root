// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package initramfs

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/u-root/u-root/pkg/cpio"
)

// Files are host files and records to add to the resulting initramfs.
type Files struct {
	// Files is a map of relative archive path -> absolute host file path.
	Files map[string]string

	// Records is a map of relative archive path -> Record to use.
	//
	// TODO: While the only archive mode is cpio, this will be a
	// cpio.Record. If or when there is another archival mode, we can add a
	// similar uroot.Record type.
	Records map[string]cpio.Record
}

// NewFiles returns a new archive files map.
func NewFiles() *Files {
	return &Files{
		Files:   make(map[string]string),
		Records: make(map[string]cpio.Record),
	}
}

// sortedKeys returns a list of sorted paths in the archive.
func (af *Files) sortedKeys() []string {
	keys := make([]string, 0, len(af.Files)+len(af.Records))
	for dest := range af.Files {
		keys = append(keys, dest)
	}
	for dest := range af.Records {
		keys = append(keys, dest)
	}
	sort.Sort(sort.StringSlice(keys))
	return keys
}

// AddFile adds a host file at src into the archive at dest.
//
// If src is a directory, it and its children will be added to the archive
// relative to dest.
//
// Duplicate files with identical content will be silently ignored.
func (af *Files) AddFile(src string, dest string) error {
	src = filepath.Clean(src)
	dest = path.Clean(dest)
	if path.IsAbs(dest) {
		return fmt.Errorf("archive path %q must not be absolute (host source %q)", dest, src)
	}

	// We check if it's a directory first. If a directory already exists as
	// a record or file, we want to include its children anyway.
	sInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("Adding %q to archive failed because stat failed: %v", src, err)
	}

	// Recursively add children.
	if sInfo.Mode().IsDir() {
		err := children(src, func(name string) error {
			return af.AddFile(filepath.Join(src, name), filepath.Join(dest, name))
		})
		if err != nil {
			return err
		}

		// Only override an existing directory if all children were
		// added successfully.
		af.Files[dest] = src
		return nil
	}

	if record, ok := af.Records[dest]; ok {
		return fmt.Errorf("record for %q already exists in archive: %v", dest, record)
	}

	if srcFile, ok := af.Files[dest]; ok {
		// Just a duplicate.
		if src == srcFile {
			return nil
		}
		return fmt.Errorf("record for %q already exists in archive (is %q)", dest, src)
	}

	af.Files[dest] = src
	return nil
}

// AddRecord adds a cpio.Record into the archive at `r.Name`.
func (af *Files) AddRecord(r cpio.Record) error {
	r.Name = path.Clean(r.Name)
	if filepath.IsAbs(r.Name) {
		return fmt.Errorf("record name %q must not be absolute", r.Name)
	}

	if src, ok := af.Files[r.Name]; ok {
		return fmt.Errorf("record for %q already exists in archive: file %q", r.Name, src)
	}
	if rr, ok := af.Records[r.Name]; ok {
		if rr.Info == r.Info {
			return nil
		}
		return fmt.Errorf("record for %q already exists in archive: %v", r.Name, rr)
	}

	af.Records[r.Name] = r
	return nil
}

// Contains returns whether path `dest` is already contained in the archive.
func (af *Files) Contains(dest string) bool {
	_, fok := af.Files[dest]
	_, rok := af.Records[dest]
	return fok || rok
}

// Rename renames a file in the archive.
func (af *Files) Rename(name string, newname string) {
	if src, ok := af.Files[name]; ok {
		delete(af.Files, name)
		af.Files[newname] = src
	}
	if record, ok := af.Records[name]; ok {
		delete(af.Records, name)
		record.Name = newname
		af.Records[newname] = record
	}
}

// addParent recursively adds parent directory records for `name`.
func (af *Files) addParent(name string) {
	parent := path.Dir(name)
	if parent == "." {
		return
	}
	if !af.Contains(parent) {
		af.AddRecord(cpio.Directory(parent, 0755))
	}
	af.addParent(parent)
}

// fillInParents adds parent directory records for unparented files in `af`.
func (af *Files) fillInParents() {
	for name := range af.Files {
		af.addParent(name)
	}
	for name := range af.Records {
		af.addParent(name)
	}
}

// WriteTo writes all records and files in `af` to `w`.
func (af *Files) WriteTo(w Writer) error {
	// Add parent directories when not added specifically.
	af.fillInParents()

	// Reproducible builds: Files should be added to the archive in the
	// same order.
	for _, path := range af.sortedKeys() {
		if record, ok := af.Records[path]; ok {
			if err := w.WriteRecord(record); err != nil {
				return err
			}
		}
		if src, ok := af.Files[path]; ok {
			if err := writeFile(w, src, path); err != nil {
				return err
			}
		}
	}
	return nil
}

// writeFile takes the file at `src` on the host system and adds it to the
// archive `w` at path `dest`.
//
// If `src` is a directory, its children will be added to the archive as well.
func writeFile(w Writer, src, dest string) error {
	record, err := cpio.GetFollowedRecord(src)
	if err != nil {
		return err
	}

	// Fix the name.
	record.Name = dest
	return w.WriteRecord(cpio.MakeReproducible(record))
}

// children calls `fn` on all direct children of directory `dir`.
func children(dir string, fn func(name string) error) error {
	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return err
	}

	for _, name := range names {
		if err := fn(name); os.IsNotExist(err) {
			// File was deleted in the meantime.
			continue
		} else if err != nil {
			return err
		}
	}
	return nil
}
