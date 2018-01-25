// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uroot

import (
	"path/filepath"
	"sync"

	"github.com/u-root/u-root/pkg/golang"
)

// BinaryBuild builds all given packages as separate binaries and includes them
// in the archive.
func BinaryBuild(opts BuildOpts) (ArchiveFiles, error) {
	af := NewArchiveFiles()

	result := make(chan error, len(opts.Packages))
	var wg sync.WaitGroup

	for _, pkg := range opts.Packages {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			result <- opts.Env.Build(
				p,
				filepath.Join(opts.TempDir, "bin", filepath.Base(p)),
				golang.BuildOpts{})
		}(pkg)
	}

	wg.Wait()
	close(result)

	for err := range result {
		if err != nil {
			return ArchiveFiles{}, err
		}
	}

	// Add bin directory to archive.
	if err := af.AddFile(opts.TempDir, ""); err != nil {
		return ArchiveFiles{}, err
	}
	return af, nil
}
