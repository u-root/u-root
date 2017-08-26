// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copy files.
//
// Synopsis:
//     cp [-rRfivwP] FROM... TO
//
// Options:
//     -w n: number of worker goroutines
//     -r: copy file hierarchies
//     -i: prompt about overwriting file
//     -f: force overwrite files
//     -v: verbose copy mode
//     -P: don't follow symlinks
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	//"strconv"
)

// buffSize is the length of buffer during
// the parallel copy using worker function
const buffSize = 8192

var (
	input = bufio.NewReader(os.Stdin)
)

type options struct {
	recursive bool
	ask       bool
	force     bool
	verbose   bool
	symlink   bool
}

// promptOverwrite ask if the user wants overwrite file
// TODO: create a new type of error 

func promptOverwrite(dst string) (bool, error) {
	for {
		fmt.Printf("cp: overwrite %q? (y/n)", dst)
		answer, err := input.ReadString('\n')
		if err != nil {
			return false, err
		}
		switch answer{
		case "y":
			return true, nil
		case "n":
			return false, nil
		}
		fmt.Printf("Please only enter y or n. You entered %s. \n", answer)
	}
	
}

// copyFile copies file between src (source) and dst (destination)
// todir: if true insert src INTO dir dst
func copyFile(src, dst string, todir bool, flags options) error {
	fmt.Printf("in function copy file ")
	if todir {
		file := filepath.Base(src)
		dst = filepath.Join(dst, file)
	}

	srcb, err := os.Lstat(src)
	if err != nil {
		return fmt.Errorf("can't stat %v: %v", src, err)
	}
	//TODO study behavior of CP 
	// don't follow symlinks, copy symlink
	fmt.Printf(" scrb.Mode is %v and L is %v ", srcb.Mode(), os.ModeSymlink)
	if srcb.Mode()&os.ModeType == os.ModeSymlink {
		fmt.Print(" scrb.Mode is %v and modesym is %v ", srcb.Mode(), os.ModeSymlink)
		if flags.symlink{
			linkPath, err := filepath.EvalSymlinks(src)
			if err != nil {
				return fmt.Errorf("can't eval symlink %v: %v", src, err)
			}
			return os.Symlink(linkPath, dst)
		} 		
		
	}

	if srcb.IsDir() {
		if flags.recursive {
			return copyDir(src, dst, flags)
		}
		return fmt.Errorf("%q is a directory, try use recursive option", src)
	}

	dstb, err := os.Stat(dst)
	if !os.IsNotExist(err) {
		if sameFile(srcb.Sys(), dstb.Sys()) {
			return fmt.Errorf("%q and %q are the same file", src, dst)
		}
		if flags.ask && !flags.force {
			overwrite, err := promptOverwrite(dst)
			if err != nil {
				return err
			}
			if !overwrite {
				return nil
			}
		}
	}

	mode := srcb.Mode() & 0777
	s, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("can't open %q: %v", src, err)
	}
	defer s.Close()

	d, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("can't create %q: %v", dst, err)
	}
	defer d.Close()

	_, err = io.Copy(d, s)
	return err
}


// createDir populate dir destination if not exists
// if exists verify is not a dir: return error if is file
// cannot overwrite: dir -> file
func createDir(src, dst string, flags options) error {
	dstInfo, err := os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if err == nil {
		if !dstInfo.IsDir() {
			return fmt.Errorf("can't overwrite non-dir %q with dir %q", dst, src)
		}
		return nil
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if err := os.Mkdir(dst, srcInfo.Mode()); err != nil {
		return err
	}
	if flags.verbose {
		fmt.Printf("%q -> %q\n", src, dst)
	}

	return nil
}

// copyDir copy the file hierarchies
// used at cp when -r or -R flag is true
func copyDir(src, dst string, flags options) error {
	if err := createDir(src, dst, flags); err != nil {
		return err
	}

	// list files from destination
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return fmt.Errorf("can't list files from %q: %q", src, err)
	}

	// copy recursively the src -> dst
	for _, file := range files {
		fname := file.Name()
		fpath := filepath.Join(src, fname)
		newDst := filepath.Join(dst, fname)
		copyFile(fpath, newDst, false, flags)
	}

	return err
}

// cp is a function whose eval the args
// and make decisions for copyfiles
func cp(args []string, flags options) (lastErr error) {
	todir := false
	from, to := args[:len(args)-1], args[len(args)-1]
	toStat, err := os.Stat(to)
	if err == nil {
		todir = toStat.IsDir()
	}
	if flag.NArg() > 2 && todir == false {
		log.Fatalf("is not a directory: %s\n", to)
	}

	for _, file := range from {
		copyFile(file, to, todir, flags)	
		if err != nil {
			log.Printf("cp: %v\n", err)
			lastErr = err
		}
	}

	return err
}

func main() {
	var flags options
	flag.BoolVar(&flags.recursive, "r", false, "alias to -R recursive mode")
	flag.BoolVar(&flags.ask, "i", false, "prompt about overwriting file")
	flag.BoolVar(&flags.force, "f", false, "force overwrite files")
	flag.BoolVar(&flags.verbose, "v", false, "verbose copy mode")
	flag.BoolVar(&flags.symlink, "P", false, "don't follow symlinks")
	flag.Parse()

	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	if err := cp(flag.Args(), flags); err != nil {
		os.Exit(1)
	}

}
