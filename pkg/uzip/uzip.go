// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package uzip contains utilities for file system->zip and zip->file system conversions.
package uzip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ToZip packs the all files at dir to a zip archive at dest.
func ToZip(dir, dest string) error {
	info, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !(info.IsDir()) {
		return fmt.Errorf("%s is not a directory", dir)
	}
	archive, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer archive.Close()

	z := zip.NewWriter(archive)
	defer z.Close()

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
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

		// adjust header.Name to preserve folder strulture
		header.Name = strings.TrimPrefix(path, dir)

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

	return err
}

// FromZip extracts the zip archive at src to dir.
func FromZip(src, dir string) error {
	z, err := zip.OpenReader(src)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	for _, file := range z.File {
		path := filepath.Join(dir, file.Name)
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
