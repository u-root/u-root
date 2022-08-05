// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// du - estimate and display disk usage of files
//
// Synopsis:
//     du [OPTIONS] [FILE]...
//
// Options:
//     -a:               write count of all files, not just directories
// 	   -B=SIZE:          scale sizes to SIZE before printing on console
//     -c                display grand total
//     -h:               print sizes in human readable format
//     -S:               for directories, do not include sizes of subdirectories
//     -s:               display only total for each directory
//     -t:               show time of last modification of any file or directory
package main

//TODOS: parse files to array, printout for every file, recursive calls, flags
import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	flag "github.com/spf13/pflag"
)

var (
	all            = flag.BoolP("all", "a", false, "write count of all files, not just directories")
	block_size     = flag.StringP("block-size", "B", "FIX", "scale sizes to SIZE before printing on console")
	human_readable = flag.BoolP("human-readable", "h", false, "print sizes in human readable format")
	separate_dirs  = flag.BoolP("separate-dirs", "S", false, "for directories, do not include sizes of subdirectories")
	summarize      = flag.BoolP("summarize", "s", false, "display only total for  each directory")
	timeF          = flag.Bool("time", false, "show time of last modification of any file or directory")
	total          = flag.BoolP("total", "c", false, "display grand total")
)

func du(w io.Writer, paths []string) {
	fileProperties := make(map[string]int64)
	dirProperties := make(map[string]int64)

	for _, path := range paths {
		fileInfo, err := os.Stat(path)

		if err != nil {
			fmt.Fprintf(w, "failed to access path %v", path)
			continue
		}

		if fileInfo.Mode() == fs.ModeSymlink {
			fileProperties[path] = fileInfo.Size()
			continue
		}

		filePS, dirPS := processPath(w, path)

		for path, size := range *filePS {
			fileProperties[path] = size
		}

		for path, size := range *dirPS {
			dirProperties[path] = size
		}

	}
	//sort
	//update folder size

	for path, size := range dirProperties {
		fileProperties[path] = size
	}

	//sort
	var sortedFiles []string
	for fileSize := range fileProperties {
		sortedFiles = append(sortedFiles, fileSize)
	}
	sort.Strings(sortedFiles)

	printFileProperties(w, fileProperties, sortedFiles)
}

func processPath(w io.Writer, dirPath string) (*map[string]int64, *map[string]int64) {
	fileProperties := make(map[string]int64)
	dirProperties := make(map[string]int64)

	err := filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(w, "failed to walk file %v: %v\n", path, err)
			return err
		}

		if info.IsDir() {
			dirProperties[path] = info.Size() >> 10
			return nil
		}

		fileProperties[path] = info.Size() >> 10
		return nil

	})

	if err != nil {
		fmt.Fprintf(w, "failed to walk file %v\n", dirPath)
		return nil, nil
	}

	return &fileProperties, &dirProperties
}

func printFileProperties(w io.Writer, fileProperties map[string]int64, sortedPaths []string) {
	var sizes []string
	maxLenghtSizes := 1

	for _, sortedPath := range sortedPaths {
		size := fmt.Sprintf("%d", fileProperties[sortedPath])

		if len(size) > maxLenghtSizes {
			maxLenghtSizes = len(size)
		}
		sizes = append(sizes, size)
	}

	idx := 0
	for range fileProperties {
		sizes[idx] = fmt.Sprintf("%*v", maxLenghtSizes, sizes[idx])
		fmt.Fprintf(w, "%v %v\n", sizes[idx], sortedPaths[idx])
		idx++
	}

}

func main() {
	flag.Parse()
	// du(os.Stdout, flag.Args())
	du(os.Stdout, []string{"/home/l1x0r/Pictures"})

}
