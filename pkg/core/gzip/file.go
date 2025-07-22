// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gzip

import (
	"path/filepath"
	"strings"

	pkggzip "github.com/u-root/u-root/pkg/gzip"
)

// getOutputPath is a helper function to replicate the outputPath method of the File struct
// since it's not exported in the original package.
func getOutputPath(f *pkggzip.File) string {
	if f.Options.Stdout || f.Options.Test {
		return f.Path
	} else if f.Options.Decompress {
		return strings.TrimSuffix(f.Path, f.Options.Suffix)
	}
	return f.Path + f.Options.Suffix
}

// resolveOutputPath resolves the output path relative to the working directory
func resolveOutputPath(g *Gzip, f *pkggzip.File) string {
	outputPath := getOutputPath(f)
	// If the path is already absolute or there's no working directory, return it as is
	if filepath.IsAbs(outputPath) || g.WorkingDir == "" {
		return outputPath
	}
	// Otherwise, resolve it relative to the working directory
	return filepath.Join(g.WorkingDir, outputPath)
}
