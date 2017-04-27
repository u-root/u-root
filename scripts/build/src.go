// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

import (
	"fmt"
	gobuild "go/build"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sync"
)

func init() {
	builders["src"] = srcBuilder{}
}

type srcBuilder struct {
}

type srcDstPair struct {
	src, dst string
}

// Generate u-root files with on-the-fly compilation. This includes the Go
// toolchain.
func (b srcBuilder) generate(config Config) ([]file, error) {
	// For parallism, store files in a chan and convert to a slice afterwards.
	fileChan := make(chan file)
	wg := sync.WaitGroup{}

	// Create a temporary directory for the output of "go build".
	tempDir, err := ioutil.TempDir("", "uroot")
	if err != nil {
		return nil, err
	}
	// It is a bit strange because even though we delete intermediate file
	// before they are read by the next stage, since they were opened before
	// they are deleted, they can still be read. This is behaviour we want.
	defer os.RemoveAll(tempDir)

	// Read all go source files of the selected packages along with all the
	// dependent source files.
	wg.Add(1)
	go func() {
		defer wg.Done()
		files, err := listGoFiles(config)
		if err != nil {
			log.Fatalf("%v", err)
		}
		for _, f := range files {
			data, err := os.Open(f.src)
			if err != nil {
				log.Fatalf("unable to open %q: %v", f.src, err)
			}
			fileChan <- file{
				path: f.dst,
				data: data,
				mode: os.FileMode(0444),
			}
		}
	}()

	// Compile the five binaries needed for the Go toolchain: init, go,
	// compile, link and asm.
	toolDir := fmt.Sprintf("go/pkg/tool/%v_%v", gobuild.Default.GOOS, gobuild.Default.GOARCH)
	for _, v := range []srcDstPair{
		{path.Join(gobuild.Default.GOPATH, "src/github.com/u-root/u-root/cmds/init"), "init"},
		{path.Join(gobuild.Default.GOROOT, "src/cmd/go"), "go/bin/go"},
		{path.Join(gobuild.Default.GOROOT, "src/cmd/compile"), path.Join(toolDir, "compile")},
		{path.Join(gobuild.Default.GOROOT, "src/cmd/link"), path.Join(toolDir, "link")},
		{path.Join(gobuild.Default.GOROOT, "src/cmd/asm"), path.Join(toolDir, "asm")},
	} {
		wg.Add(1)
		go func(v srcDstPair) {
			defer wg.Done()
			outPath := path.Join(tempDir, v.dst)
			if err := buildBinary(v.src, outPath); err != nil {
				log.Fatalf("failed building binary %q: %v", v.src, err)
			}
			data, err := os.Open(outPath)
			if err != nil {
				log.Fatalf("unable to read %q: %v", outPath, err)
			}
			fileChan <- file{
				path: v.dst,
				data: data,
				mode: os.FileMode(0555),
			}
		}(v)
	}

	// Once all goroutines are complete, close the channel.
	go func() {
		wg.Wait() // TODO: add a progress bar ;)
		close(fileChan)
	}()

	// Copy the files from the channel into a slice.
	files := []file{}
	for f := range fileChan {
		files = append(files, f)
	}
	return files, nil
}

// Build a Go binary.
func buildBinary(srcPkgPath, outPath string) error {
	buildArgs := []string{
		"build",

		// rebuild all from scratch
		"-a",

		// build in a separate directory with the "uroot" suffix
		"-installsuffix=uroot",

		// strip symbols (huge space savings)
		"-ldflags=-s -w",

		// output binary location
		"-o", outPath,
	}

	cmd := exec.Command("go", buildArgs...)
	cmd.Dir = srcPkgPath
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// listGoFiles determines the list of Go source files for inclusion.
func listGoFiles(config Config) ([]srcDstPair, error) {
	// Perform a breadth-first-search to find all dependencies.
	unvisited := make([]string, len(config.Packages))
	copy(unvisited, config.Packages)
	parents := make([]string, len(config.Packages))
	pkgSet := make(map[string]bool)
	files := []srcDstPair{}

	for len(unvisited) != 0 {
		// pop
		pkgName := unvisited[len(unvisited)-1]
		unvisited = unvisited[:len(unvisited)-1]

		// pop
		parent := parents[len(parents)-1]
		parents = parents[:len(parents)-1]

		// Skip fake packages.
		if pkgName == "C" {
			continue
		}

		// skip if already visited
		if pkgSet[pkgName] {
			continue
		}
		pkgSet[pkgName] = true

		// Parent package is needed to determine correct vendor directory.
		if parent == "" {
			// TODO: Might break if there are multiple GOPATHs. This can be
			// achieved with one more for loop!
			parent = filepath.SplitList(gobuild.Default.GOPATH)[0]
		}

		// urootPath is needed to find the proper vendor directory.
		p, err := gobuild.Import(pkgName, parent, 0)
		if err != nil {
			log.Printf("warning: cannot find package %q: %v", pkgName, err)
			continue
		}

		// Decide where the files will be mapped.
		for _, v := range append(append(p.GoFiles, p.SFiles...), p.HFiles...) {
			if p.Goroot {
				// Copy GOROOT files to "go/src".
				files = append(files, srcDstPair{
					path.Join(p.Dir, v),
					path.Join("go/src", p.ImportPath, v),
				})
			} else {
				// Copy GOPATH files to "src".
				files = append(files, srcDstPair{
					path.Join(p.Dir, v),
					path.Join("src", p.ImportPath, v),
				})
			}
		}

		// Check transitive dependencies.
		for _, import_ := range p.Imports {
			unvisited = append(unvisited, import_)
			parents = append(parents, p.Dir)
		}
	}

	// There are also some additional header files needed by the assembler
	// (even when cgo is disabled).
	matches, err := filepath.Glob(path.Join(gobuild.Default.GOROOT, "pkg/include/*"))
	if err != nil {
		return nil, err
	}
	for _, v := range matches {
		files = append(files, srcDstPair{v, path.Join("go/pkg/include", filepath.Base(v))})
	}

	return files, nil
}
