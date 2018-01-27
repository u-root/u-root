package uroot

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/u-root/u-root/pkg/cpio"
	"golang.org/x/sys/unix"
)

// ArchiveFiles are host files and records to add to the resulting initramfs.
type ArchiveFiles struct {
	// Files is a map of relative archive path -> absolute host file path.
	Files map[string]string

	// Records is a map of relative archive path -> Record to use.
	//
	// TODO: While the only archive mode is cpio, this will be a
	// cpio.Record. If or when there is another archival mode, we can add a
	// similar uroot.Record type.
	Records map[string]cpio.Record
}

// NewArchiveFiles returns a new archive files map.
func NewArchiveFiles() ArchiveFiles {
	return ArchiveFiles{
		Files:   make(map[string]string),
		Records: make(map[string]cpio.Record),
	}
}

// SortedKeys returns a list of sorted paths in the archive.
func (af ArchiveFiles) SortedKeys() []string {
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

// AddFile adds a host file at `src` into the archive at `dest`.
func (af ArchiveFiles) AddFile(src string, dest string) error {
	src = filepath.Clean(src)
	dest = path.Clean(dest)
	if path.IsAbs(dest) {
		return fmt.Errorf("archive path %q must not be absolute (host source %q)", dest, src)
	}
	if !filepath.IsAbs(src) {
		return fmt.Errorf("source file %q (-> %q) must be absolute", src, dest)
	}

	if _, ok := af.Records[dest]; ok {
		return fmt.Errorf("record for %q already exists in archive", dest)
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
func (af ArchiveFiles) AddRecord(r cpio.Record) error {
	r.Name = path.Clean(r.Name)
	if filepath.IsAbs(r.Name) {
		return fmt.Errorf("record name %q must not be absolute", r.Name)
	}

	if _, ok := af.Files[r.Name]; ok {
		return fmt.Errorf("record for %q already exists in archive", r.Name)
	}
	if rr, ok := af.Records[r.Name]; ok {
		if rr.Info == r.Info {
			return nil
		}
		return fmt.Errorf("record for %q already exists in archive", r.Name)
	}

	af.Records[r.Name] = r
	return nil
}

// Contains returns whether path `dest` is already contained in the archive.
func (af ArchiveFiles) Contains(dest string) bool {
	_, fok := af.Files[dest]
	_, rok := af.Records[dest]
	return fok || rok
}

// Rename renames a file in the archive.
func (af ArchiveFiles) Rename(name string, newname string) {
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
func (af ArchiveFiles) addParent(name string) {
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
func (af ArchiveFiles) fillInParents() {
	for name := range af.Files {
		af.addParent(name)
	}
	for name := range af.Records {
		af.addParent(name)
	}
}

// WriteTo writes all records and files in `af` to `w`.
func (af ArchiveFiles) WriteTo(w ArchiveWriter) error {
	// Add parent directories when not added specifically.
	af.fillInParents()

	// Reproducible builds: Files should be added to the archive in the
	// same order.
	for _, path := range af.SortedKeys() {
		if record, ok := af.Records[path]; ok {
			if err := w.WriteRecord(record); err != nil {
				return err
			}
		}
		if src, ok := af.Files[path]; ok {
			if err := WriteFile(w, src, path); err != nil {
				return err
			}
		}
	}
	return nil
}

// Write uses the given options to determine which files need to be written
// to the output file using the archive format `a` and writes them.
func (opts *ArchiveOpts) Write() error {
	// Add default records.
	for _, rec := range opts.DefaultRecords {
		// Ignore if it doesn't get added. Probably means the user
		// included something for this file or directory already.
		//
		// TODO: ignore only when it already exists in archive.
		opts.ArchiveFiles.AddRecord(rec)
	}

	// Write base archive.
	if opts.BaseArchive != nil {
		transform := cpio.MakeReproducible

		// Rename init to inito if user doesn't want the existing init.
		if !opts.UseExistingInit && opts.Contains("init") {
			transform = func(r cpio.Record) cpio.Record {
				if r.Name == "init" {
					r.Name = "inito"
				}
				return cpio.MakeReproducible(r)
			}
		}
		// If user wants the base archive init, but specified another
		// init, make the other one inito.
		if opts.UseExistingInit && opts.Contains("init") {
			opts.Rename("init", "inito")
		}

		for {
			f, err := opts.BaseArchive.ReadRecord()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			// TODO: ignore only the error where it already exists
			// in archive.
			opts.ArchiveFiles.AddRecord(transform(f))
		}
	}

	if err := opts.ArchiveFiles.WriteTo(opts.OutputFile); err != nil {
		return err
	}
	return opts.OutputFile.Finish()
}

// WriteFile takes the file at `src` on the host system and adds it to the
// archive `w` at path `dest`.
//
// If `src` is a directory, its children will be added to the archive as well.
func WriteFile(w ArchiveWriter, src, dest string) error {
	record, err := cpio.GetRecord(src)
	if err != nil {
		return err
	}

	// Fix the name.
	record.Name = dest
	if err := w.WriteRecord(cpio.MakeReproducible(record)); err != nil {
		return err
	}

	if record.Info.Mode&unix.S_IFMT == unix.S_IFDIR {
		return children(src, func(name string) error {
			return WriteFile(w, filepath.Join(src, name), filepath.Join(dest, name))
		})
	}
	return nil
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
