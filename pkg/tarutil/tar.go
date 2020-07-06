// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tarutil

import (
	"archive/tar"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// passesFilters returns true if the given file passes all filters, false otherwise.
func passesFilters(hdr *tar.Header, filters []Filter) bool {
	for _, filter := range filters {
		if !filter(hdr) {
			return false
		}
	}
	return true
}

// applyToArchive applies function f to all files in the given archive
func applyToArchive(
	tarFile io.Reader, f func(tr *tar.Reader, hdr *tar.Header) error) error {
	tr := tar.NewReader(tarFile)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if err := f(tr, hdr); err != nil {
			return err
		}
	}
	return nil
}

// ListArchive lists the contents of the given tar archive.
func ListArchive(tarFile io.Reader) error {
	return applyToArchive(tarFile, func(tr *tar.Reader, hdr *tar.Header) error {
		fmt.Println(hdr.Name)
		return nil
	})
}

// ExtractDir extracts all the contents of the tar file to the given directory.
func ExtractDir(tarFile io.Reader, dir string) error {
	return ExtractDirFilter(tarFile, dir, nil)
}

// ExtractDirFilter extracts a tar file with the given filter.
func ExtractDirFilter(tarFile io.Reader, dir string, filters []Filter) error {
	fi, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			return fmt.Errorf("could not create directory %s: %v", dir, err)
		}
	} else if err != nil || !fi.IsDir() {
		return fmt.Errorf("could not stat directory %s: %v", dir, err)
	}

	return applyToArchive(tarFile, func(tr *tar.Reader, hdr *tar.Header) error {
		if !passesFilters(hdr, filters) {
			return nil
		}
		return createFileInRoot(hdr, tr, dir)
	})
}

// CreateTar creates a new tar file with all the contents of a directory.
func CreateTar(tarFile io.Writer, files []string) error {
	return CreateTarFilter(tarFile, files, nil)
}

// CreateTarFilter creates a new tar file of the given files, with the given filter.
func CreateTarFilter(tarFile io.Writer, files []string, filters []Filter) error {
	tw := tar.NewWriter(tarFile)
	for _, file := range files {
		err := filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			symlink := ""
			if info.Mode()&os.ModeSymlink != 0 {
				// TODO: symlinks
				return fmt.Errorf("symlinks not yet supported: %q", path)
			}
			hdr, err := tar.FileInfoHeader(info, symlink)
			if err != nil {
				return err
			}
			hdr.Name = path
			if !passesFilters(hdr, filters) {
				return nil
			}
			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}
			switch hdr.Typeflag {
			case tar.TypeLink, tar.TypeSymlink, tar.TypeChar, tar.TypeBlock, tar.TypeDir, tar.TypeFifo:
			default:
				f, err := os.Open(path)
				if err != nil {
					return err
				}
				if _, err := io.Copy(tw, f); err != nil {
					f.Close()
					return err
				}
				f.Close()
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	if err := tw.Close(); err != nil {
		return err
	}
	return nil
}

func createFileInRoot(hdr *tar.Header, r io.Reader, rootDir string) error {
	fi := hdr.FileInfo()
	path := filepath.Clean(filepath.Join(rootDir, hdr.Name))
	if !strings.HasPrefix(path, filepath.Clean(rootDir)) {
		return fmt.Errorf("file outside root directory: %q", path)
	}

	switch fi.Mode() & os.ModeType {
	case os.ModeSymlink:
		// TODO: support symlinks
		return fmt.Errorf("symlinks not yet supported: %q", path)

	case os.FileMode(0):
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		if _, err := io.Copy(f, r); err != nil {
			f.Close()
			return err
		}
		if err := f.Close(); err != nil {
			return err
		}

	case os.ModeDir:
		if err := os.MkdirAll(path, fi.Mode()&os.ModePerm); err != nil {
			return err
		}

	case os.ModeDevice:
		// TODO: support block device
		return fmt.Errorf("block device not yet supported: %q", path)

	case os.ModeCharDevice:
		// TODO: support char device
		return fmt.Errorf("char device not yet supported: %q", path)

	default:
		return fmt.Errorf("%q: Unknown type %#o", path, fi.Mode()&os.ModeType)
	}

	if err := os.Chmod(path, fi.Mode()&os.ModePerm); err != nil {
		return fmt.Errorf("error setting mode %#o on %q: %v",
			fi.Mode()&os.ModePerm, path, err)
	}
	// TODO: also set ownership, etc...
	return nil
}

// Filter is applied to each file while creating or extracting a tar archive.
// The filter can modify the tar.Header struct. If the filter returns false,
// the file is omitted.
type Filter func(hdr *tar.Header) bool

// NoFilter does not filter or modify any files.
func NoFilter(hdr *tar.Header) bool {
	return true
}

// VerboseFilter prints the name of every file.
func VerboseFilter(hdr *tar.Header) bool {
	fmt.Println(hdr.Name)
	return true
}

// VerboseLogFilter logs the name of every file.
func VerboseLogFilter(hdr *tar.Header) bool {
	log.Println(hdr.Name)
	return true
}

// SafeFilter filters out all files which are not regular and not directories.
// It also sets appropriate permissions.
func SafeFilter(hdr *tar.Header) bool {
	if hdr.Typeflag == tar.TypeDir {
		hdr.Mode = 0770
		return true
	}
	if hdr.Typeflag == tar.TypeReg {
		hdr.Mode = 0660
		return true
	}
	return false
}
