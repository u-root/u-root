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
	"sort"

	flag "github.com/spf13/pflag"
)

const (
	root = "."
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

type FileProperties struct {
	path     string
	byteSize int64
}

func du(w io.Writer, paths []string) {
	var logOutput []FileProperties

	for _, path := range paths {
		fileInfo, err := os.Stat(path)

		if err != nil {
			fmt.Fprintf(w, "failed to access path %v", path)
			continue
		}

		if fileInfo.Mode() == fs.ModeSymlink {
			logOutput = append(logOutput, FileProperties{path, fileInfo.Size()})
			continue
		}

		fsys := os.DirFS(path)
		fileproperties := processFS(w, fsys, path)

		if len(*fileproperties) != 0 {
			logOutput = append(logOutput, *fileproperties...)
		} else {
		}
	}
	printFileProperties(w, logOutput)
}

func processFS(w io.Writer, fsys fs.FS, dirPath string) *[]FileProperties {
	var fileProperties []FileProperties
	var rootProperties FileProperties

	err := fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if d != nil {
				fmt.Fprintf(w, "failed to read dir %v: %v\n", path, err)
				return fs.SkipDir
			} else {
				fmt.Fprintf(w, "failed to walk directory %v: %v\n", path, err)
				return err
			}
		}

		fileInfo, err := d.Info()

		if err != nil {
			fmt.Fprintf(w, "failed to access file infos of %v\n", path)
			return err
		}

		if path == root {
			rootProperties = FileProperties{dirPath, fileInfo.Size()}
			return nil
		}

		if d.IsDir() {
			fileProperties = append(fileProperties, FileProperties{"." + dirPath + "/" + path, fileInfo.Size()})
			return nil
		}

		fileProperties = append(fileProperties, FileProperties{"." + dirPath + "/" + path, fileInfo.Size()})
		return nil

	})

	if err != nil {
		fmt.Fprintf(w, "failed to walk directory %v\n", dirPath)
		return nil
	}

	sort.Slice(fileProperties, func(i, j int) bool {
		return sort.StringsAreSorted([]string{fileProperties[i].path, fileProperties[j].path})
	})

	fileProperties = append(fileProperties, rootProperties)

	return &fileProperties
}

func printFileProperties(w io.Writer, fileProperties []FileProperties) {
	var sizes []string
	var paths []string
	maxLenghtSizes := 1

	for _, properties := range fileProperties {
		size := fmt.Sprintf("%d", properties.byteSize)

		if len(size) > maxLenghtSizes {
			maxLenghtSizes = len(size)
		}
		sizes = append(sizes, size)
		paths = append(paths, fmt.Sprintf("%v", properties.path))
	}

	for idx, _ := range fileProperties {
		sizes[idx] = fmt.Sprintf("%*v", maxLenghtSizes, sizes[idx])
		fmt.Fprintf(w, "%v %v\n", sizes[idx], paths[idx])
	}

}

func main() {
	flag.Parse()
	du(os.Stdout, flag.Args())
}
