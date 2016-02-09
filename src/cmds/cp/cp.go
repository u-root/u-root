// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

const buffSize = 8192

var (
	recursive bool
	ask       bool
	force     bool
	verbose   bool
	symlink   bool
	nwork     int
	input     = bufio.NewScanner(os.Stdin)
	offchan   = make(chan int64, 0)
	zerochan  = make(chan int, 0)
)

func init() {
	flag.IntVar(&nwork, "w", runtime.NumCPU(), "number of worker goroutines")
	flag.BoolVar(&recursive, "R", false, "copy file hierarchies")
	flag.BoolVar(&recursive, "r", false, "alias to -R recursive mode")
	flag.BoolVar(&ask, "i", false, "prompt about overwriting file")
	flag.BoolVar(&force, "f", false, "force overwrite files")
	flag.BoolVar(&verbose, "v", false, "verbose copy mode")
	flag.BoolVar(&symlink, "P", false, "don't follow symlinks")
	flag.Parse()
	go nextOff()

}

// ask if the user wants overwrite file
func promptOverwrite(dst string) bool {
	_, err := os.Stat(dst)
	if !os.IsNotExist(err) {
		fmt.Printf("cp: overwrite '%v'? ", dst)
		input.Scan()
		if input.Text()[0] != 'y' {
			return false
		}
	}
	return true
}

// copy src to dst
// todir: if true insert src INTO dir dst
func copyFile(src, dst string, todir bool) error {
	if todir {
		_, file := path.Split(src)
		dst = path.Join(dst, file)
	}

	srcb, err := os.Lstat(src)
	if err != nil {
		return fmt.Errorf("can't stat %v: %v", src, err)
	}

	// don't follow symlinks, copy symlink
	if L := os.ModeSymlink; symlink && srcb.Mode()&L == L {
		linkPath, err := filepath.EvalSymlinks(src)
		if err != nil {
			return fmt.Errorf("can't eval symlink %v: %v", src, err)
		}
		return os.Symlink(linkPath, dst)
	}

	if srcb.IsDir() {
		if recursive {
			return copyDir(src, dst)
		} else {
			return fmt.Errorf("'%v' is a directory, try use recursive option", src)
		}
	}

	dstb, err := os.Stat(dst)
	if err == nil {
		if sameFile(srcb.Sys(), dstb.Sys()) {
			return fmt.Errorf("'%v' and '%v' are the same file", src, dst)
		}
	}
	if ask && !force {
		if !promptOverwrite(dst) {
			return nil
		}
	}

	mode := srcb.Mode() & 0777
	s, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("can't open '%v': %v", src, err)
	}
	defer s.Close()

	d, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("can't create '%v': %v", dst, err)
	}
	defer d.Close()

	return copyOneFile(s, d, src, dst)
}

// copy the content between two files
func copyOneFile(s *os.File, d *os.File, src, dst string) error {
	zerochan <- 0
	fail := make(chan error, nwork)
	for i := 0; i < nwork; i++ {
		go worker(s, d, fail)
	}

	// iterate the errors from channel
	failed := false
	for i := 0; i < nwork; i++ {
		err := <-fail
		if err != nil {
			failed = true
			log.Println(err)
		}
	}

	if verbose {
		fmt.Printf("'%v' -> '%v'\n", src, dst)
	}

	// if some error occurs, returns that error
	if failed {
		return fmt.Errorf("cannot copy the file: '%v'", src)
	}

	return nil
}

// populate dir destination if not exists
// if exists verify is not a dir: return error if is file
// cannot overwrite: dir -> file
func createDir(src, dst string) error {
	dstInfo, err := os.Stat(dst)
	if os.IsNotExist(err) {
		srcInfo, err := os.Stat(src)
		if err != nil {
			return err
		}
		if err := os.Mkdir(dst, srcInfo.Mode()); err != nil {
			return err
		}
		if verbose {
			fmt.Printf("'%v' -> '%v'\n", src, dst)
		}
	} else if !dstInfo.IsDir() {
		return fmt.Errorf("can't overwrite non-dir '%v' with dir '%v'", dst, src)
	}

	return nil
}

// copy file hierarchies
func copyDir(src, dst string) error {
	if err := createDir(src, dst); err != nil {
		return err
	}

	// list files from destination
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return fmt.Errorf("can't list files from '%v': '%v'", src, err)
	}

	// copy recursively the src -> dst
	for _, file := range files {
		fname := file.Name()
		fpath := path.Join(src, fname)
		newDst := path.Join(dst, fname)
		copyFile(fpath, newDst, false)
	}

	return err
}

// concurrent copy, worker routine
func worker(s *os.File, d *os.File, fail chan error) {
	var buf [buffSize]byte
	var bp []byte

	l := len(buf)
	bp = buf[0:]
	o := <-offchan
	for {
		n, err := s.ReadAt(bp, o)
		if err != nil && err != io.EOF {
			fail <- fmt.Errorf("reading %s at %v: %v", s.Name(), o, err)
			return
		}
		if n == 0 {
			break
		}

		nb := bp[0:n]
		n, err = d.WriteAt(nb, o)
		if err != nil {
			fail <- fmt.Errorf("writing %s: %v", d.Name(), err)
			return
		}
		bp = buf[n:]
		o += int64(n)
		l -= n
		if l == 0 {
			l = len(buf)
			bp = buf[0:]
			o = <-offchan
		}
	}
	fail <- nil
}

// handler for next buffers
func nextOff() {
	off := int64(0)
	for {
		select {
		case <-zerochan:
			off = 0
		case offchan <- off:
			off += buffSize
		}
	}
}

func cp(args []string) (lastErr error) {
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
		if err := copyFile(file, to, todir); err != nil {
			log.Printf("cp: %v\n", err)
			lastErr = err
		}
	}

	return
}

func main() {
	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	if err := cp(flag.Args()); err != nil {
		os.Exit(1)
	}

}
