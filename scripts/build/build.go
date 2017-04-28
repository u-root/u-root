// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

import (
	"errors"
	"log"
	"os"
	"path"
	"sort"
)

// Generators register themselves in init() functions.
var (
	buildGenerators   = map[string]buildGenerator{}
	archiveGenerators = map[string]archiveGenerator{}
)

// Uniq sorts and remove duplicates from a slice of strings.
func Uniq(s []string) []string {
	set := make(map[string]bool, len(s))
	for _, v := range s {
		set[v] = true
	}
	slice := make([]string, 0, len(set))
	for k := range set {
		slice = append(slice, k)
	}
	sort.Strings(slice)
	return slice
}

// Build a u-root archive and optionally run it.
func Build(config Config) error {
	// Select the build generators.
	bGens := []buildGenerator{}
	for _, buildFormat := range config.BuildFormats {
		bGen, ok := buildGenerators[buildFormat]
		if !ok {
			return errors.New("invalid build generator")
		}
		bGens = append(bGens, bGen)
	}

	// Select the archive generator.
	aGen, ok := archiveGenerators[config.ArchiveFormat]
	if !ok {
		return errors.New("invalid archive generator")
	}

	// Generate the files.
	files := []file{}
	for _, bGen := range bGens {
		moreFiles, err := bGen.generate(config)
		if err != nil {
			return err
		}
		files = append(files, moreFiles...)
	}

	// Create directory entries where necessary.
	// Ex: banana/carrot/chicken.txt, creates directories for banana and carrot.
	dirSet := map[string]bool{}
	for _, f := range files {
		dir := f.path
		if !f.mode.IsDir() {
			for path.Dir(dir) != dir {
				dir = path.Dir(dir)
				// Set to false if unset, means the directory needs creating
				dirSet[dir] = dirSet[dir]
			}
		} else {
			// Set to true, meaning the directory already exists
			dirSet[dir] = true
		}
	}
	for k, v := range dirSet {
		if !v {
			files = append(files, file{
				path: k,
				mode: 0755 | os.ModeDir,
			})
		}
	}

	// Sort the files by path. This is requires for reproducible builds.
	sort.Slice(files, func(i, j int) bool {
		if files[i].path == files[j].path {
			log.Printf("warning: multiple files named %q", files[i].path)
		}
		return files[i].path < files[j].path
	})

	// Generate the archive.
	if err := aGen.generate(config, files); err != nil {
		return err
	}
	log.Printf("build complete: %q", config.OutputPath)

	// Optionally run the archive.
	if config.Run {
		return aGen.run(config)
	}
	return nil
}
