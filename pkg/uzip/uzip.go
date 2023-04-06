// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package uzip contains utilities for file system->zip and zip->file system conversions.
package uzip

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/u-root/u-root/pkg/upath"
)

// ToZip packs the all files at dir to a zip archive at dest.
func ToZip(dir, dest, comment string) (reterr error) {
	if info, err := os.Stat(dir); err != nil {
		return err
	} else if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", dir)
	}
	archive, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func() {
		if err := archive.Close(); err != nil && reterr == nil {
			reterr = err
		}
	}()

	z := zip.NewWriter(archive)
	defer func() {
		if comment != "" {
			z.SetComment(comment)
		}
		if err := z.Close(); err != nil && reterr == nil {
			reterr = err
		}
	}()

	return writeDir(dir, z)
}

// AppendZip packs the all files at dir to a zip archive at dest.
func AppendZip(dir, dest, comment string) (reterr error) {
	if info, err := os.Stat(dir); err != nil {
		return err
	} else if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", dir)
	}

	archive, err := os.OpenFile(dest, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer func() {
		if err := archive.Close(); err != nil && reterr == nil {
			reterr = err
		}
	}()

	// Go to the end of the file because we are appending.
	end, err := archive.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	z := zip.NewWriter(archive)
	z.SetOffset(end)
	defer func() {
		if comment != "" {
			z.SetComment(comment)
		}
		if err := z.Close(); err != nil && reterr == nil {
			reterr = err
		}
	}()

	return writeDir(dir, z)
}

func writeDir(dir string, z *zip.Writer) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// do not include srcDir into archive
		if info.Name() == filepath.Base(dir) {
			return nil
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// adjust header.Name to preserve folder structure
		header.Name, err = filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := z.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})
}

// Comment returns the comment from the zip file.
func Comment(file string) (string, error) {
	z, err := zip.OpenReader(file)
	if err != nil {
		return "", err
	}
	defer z.Close()
	return z.Comment, nil
}

// FromZip extracts the zip archive at src to dir.
func FromZip(src, dir string) error {
	z, err := zip.OpenReader(src)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	for _, file := range z.File {
		path, err := upath.SafeFilepathJoin(dir, file.Name)
		if err != nil {
			// The behavior is to skip files which are unsafe due to
			// zipslip, but continue extracting everything else.
			log.Printf("Warning: Skipping file %q due to: %v", file.Name, err)
			continue
		}

		if file.FileInfo().IsDir() {
			if err = os.MkdirAll(path, file.Mode()); err != nil {
				return err
			}
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}

		if err = fileReader.Close(); err != nil {
			return err
		}

		if err = targetFile.Close(); err != nil {
			return err
		}
	}

	return nil
}
