// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uroot

import (
	"path/filepath"
	"sync"

	"github.com/u-root/u-root/pkg/golang"
)

var BinaryBuilder = Builder{
	Build:            BinaryBuild,
	DefaultBinaryDir: "bin",
}

// BinaryBuild builds all given packages as separate binaries and includes them
// in the archive.
func BinaryBuild(af ArchiveFiles, opts BuildOpts) error {
	result := make(chan error, len(opts.Packages))
	var wg sync.WaitGroup

	for _, pkg := range opts.Packages {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			result <- opts.Env.Build(
				p,
				filepath.Join(opts.TempDir, opts.BinaryDir, filepath.Base(p)),
				golang.BuildOpts{})
		}(pkg)
	}

	wg.Wait()
	close(result)

	for err := range result {
		if err != nil {
			return err
		}
	}

	// Add bin directory to archive.
	return af.AddFile(opts.TempDir, "")
}
