// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uroot

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/golang"
)

var (
	builders = map[string]Build{
		"source": sourceBuild,
		"bb":     bbBuild,
	}
	archivers = map[string]Archiver{
		"cpio": CPIOArchiver{},
	}
)

func GetBuilder(name string) (Build, error) {
	build, ok := builders[name]
	if !ok {
		return nil, fmt.Errorf("couldn't find builder %q", name)
	}
	return build, nil
}

func GetArchiver(name string) (Archiver, error) {
	archiver, ok := archivers[name]
	if !ok {
		return nil, fmt.Errorf("couldn't find archival format %q", name)
	}
	return archiver, nil
}

type ArchiveFiles struct {
	// Files is a map of relative archive path -> absolute host file path.
	Files map[string]string

	// Records is a map of relative archive path -> Record to use.
	Records map[string]cpio.Record
}

func NewArchiveFiles() ArchiveFiles {
	return ArchiveFiles{
		Files:   make(map[string]string),
		Records: make(map[string]cpio.Record),
	}
}

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

func (af ArchiveFiles) AddFile(src string, dest string) error {
	if filepath.IsAbs(dest) {
		return fmt.Errorf("cannot add absolute path %q (from %q) to archive", dest, src)
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
		return fmt.Errorf("archive file %q already comes from %q", dest, src)
	}

	af.Files[dest] = src
	return nil
}

func (af ArchiveFiles) AddRecord(r cpio.Record) error {
	if filepath.IsAbs(r.Name) {
		return fmt.Errorf("cannot add absolute path %q to archive", r.Name)
	}

	if _, ok := af.Files[r.Name]; ok {
		return fmt.Errorf("record for %q already exists in archive", r.Name)
	}
	if rr, ok := af.Records[r.Name]; ok {
		if rr.Info == r.Info {
			return nil
		}
		return fmt.Errorf("record for %q already exists", r.Name)
	}

	af.Records[r.Name] = r
	return nil
}

func (af ArchiveFiles) Contains(dest string) bool {
	_, fok := af.Files[dest]
	_, rok := af.Records[dest]
	return fok || rok
}

type BuildOpts struct {
	Env      golang.Environ
	Packages []string
	TempDir  string
}

type Build func(BuildOpts) (ArchiveFiles, error)

type ArchiveOpts struct {
	ArchiveFiles
	OutputFile      *os.File
	BaseArchive     *os.File
	UseExistingInit bool
}

type Archiver interface {
	Archive(ArchiveOpts) error
	DefaultExtension() string
}

func DefaultPackageImports(env golang.Environ) ([]string, error) {
	// Find u-root directory.
	urootDir, err := env.FindPackageDir("github.com/u-root/u-root")
	if err != nil {
		return nil, fmt.Errorf("Couldn't find u-root src directory: %v", err)
	}

	matches, err := filepath.Glob(filepath.Join(urootDir, "cmds/*"))
	if err != nil {
		return nil, fmt.Errorf("couldn't find u-root cmds: %v", err)
	}
	pkgs := make([]string, 0, len(matches))
	for _, match := range matches {
		importPath, err := env.FindPackageByPath(match)
		if err == nil {
			pkgs = append(pkgs, importPath)
		}
	}
	return pkgs, nil
}
