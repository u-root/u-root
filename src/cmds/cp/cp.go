// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

const DefBuff = 8192
const BarChar = "="
const BarLength = 25

var (
	config struct {
		nwork     int
		recursive bool
		ask       bool
		force     bool
		verbose   bool
		progress  bool
		link      bool // to implement yet
	}
	offchan   chan int64
	zerochan  chan int
	fileschan chan int
	byteschan chan int64
	donebar   chan bool // channel to finish the progress bar
	input     = bufio.NewScanner(os.Stdin)
	cmd       = "cp"
	flags     = "[-w workers] [-Rrifv] [-bar]"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %v %v source destination\n", cmd, flags)
	flag.PrintDefaults()
	os.Exit(1) // sysfatal
}

func init() {
	flag.Usage = usage
	flag.IntVar(&config.nwork, "w", 16, "number of worker goroutines")
	flag.BoolVar(&config.recursive, "R", false, "copy file hierarchies")
	flag.BoolVar(&config.recursive, "r", false, "alias to -R recursive mode")
	flag.BoolVar(&config.ask, "i", false, "prompt about overwriting file")
	flag.BoolVar(&config.force, "f", false, "force overwrite files")
	flag.BoolVar(&config.verbose, "v", false, "verbose copy mode")
	flag.BoolVar(&config.progress, "bar", false, "progress bar on")
	// flag.BoolVar(&config.link, "P", false, "-R: link following all")
	// flag.BoolVar(&config.link, "H", false, "Only copy the link")
	flag.Parse()
}

// because don't exists errors.Newf we use that
func Errorf(formated string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(formated, args...))
}

// ask if the user wants overwrite file
func promptOverwrite(dst string) bool {
	_, err := os.Stat(dst)
	if !os.IsNotExist(err) {
		fmt.Printf("%v: overwrite '%v'? ", cmd, dst)
		input.Scan()
		if input.Text()[0] != 'y' {
			return false
		}
	}

	return true

}

// general copy of that program
// make decisions if need create directory and etc
func copy(src, dst string, todir bool) error {
	if todir {
		_, file := path.Split(src)
		dst = path.Join(dst, file)
	}

	dirb, err := os.Stat(src)
	if err != nil {
		return Errorf("can't stat %v: %v\n", src, err)
	}

	if dirb.IsDir() {
		if config.recursive {
			return copyDir(src, dst)
		} else {
			return Errorf("'%v' is a directory, try use -R option\n", src)
		}
	}

	tob, err := os.Stat(dst)
	if err == nil {
		if sameFile(dirb.Sys(), tob.Sys()) {
			return Errorf("'%v' and '%v' are the same file\n", src, dst)
		}
	}
	if config.ask && !config.force {
		if !promptOverwrite(dst) {
			return nil
		}
	}

	mode := dirb.Mode() & 0777
	s, err := os.Open(src)
	if err != nil {
		return Errorf("can't open '%v': %v\n", src, err)
	}
	defer s.Close()

	d, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return Errorf("can't create '%v': %v\n", dst, err)
	}
	defer d.Close()

	if config.verbose {
		fmt.Printf("'%v' -> '%v'\n", src, dst)
	}

	return copyOneFile(s, d)
}

// copy the content between two files
func copyOneFile(s *os.File, d *os.File) error {
	zerochan <- 0                          // ? i don't understand that channel
	fail := make(chan error, config.nwork) // channel of errors okay

	for i := 0; i < config.nwork; i++ {
		go worker(s, d, fail)
	}

	// iterate the errors from channel
	failed := false
	for i := 0; i < config.nwork; i++ {
		err := <-fail
		if err != nil {
			log.Println(err)
		}
	}

	if config.progress {
		fileschan <- 1
	}

	// if some error occurs, returns this error
	if failed != false {
		return Errorf("can't copy all the file: '%v'", s.Name())
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
		if config.verbose {
			fmt.Printf("'%v' -> '%v'\n", src, dst)
		}
	} else if !dstInfo.IsDir() {
		return Errorf("can't overwrite non-dir '%v' with dir '%v'", dst, src)
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
		return Errorf("can't list files from '%v': '%v'", src, err)
	}

	// copy recursively the src -> dst
	for _, file := range files {
		fname := file.Name()
		fpath := path.Join(src, fname)
		newDst := path.Join(dst, fname)
		copy(fpath, newDst, false)
	}

	return err
}

// concurrent copy, worker routine
func worker(s *os.File, d *os.File, fail chan error) {
	var buf [DefBuff]byte
	var bp []byte

	l := len(buf)
	bp = buf[0:]
	o := <-offchan
	for {
		n, err := s.ReadAt(bp, o)
		if err != nil && err != io.EOF {
			fail <- Errorf("reading %s at %v: %v\n", s.Name(), o, err)
			return
		}
		if n == 0 {
			break
		}

		nb := bp[0:n]
		n, err = d.WriteAt(nb, o)
		if err != nil {
			fail <- Errorf("writing %s: %v\n", d.Name(), err)
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
		if config.progress {
			byteschan <- int64(n)
		}
	}
	fail <- nil
}

func nextOff() {
	off := int64(0)
	for {
		select {
		case <-zerochan:
			off = 0
		case offchan <- off:
			off += DefBuff
		}
	}
}

// backspace on terminal
// to overwrite the last message
func flushStdout(clean int) {
	for i := 0; i < clean; i++ {
		fmt.Printf("%c", '\010')
	}
}

// based on the 0..1 value of status
// create a bar representing
func makebar(status float64) string {
	bar := ""
	for i := 0; i < int(BarLength*status); i++ {
		bar += BarChar
	}

	return bar
}

// some trick to get the unit multiples of bytes
// 0 -> 1 (byte)
// 1 -> 1024 (mByte)
// 2 -> 1024 * 1024 (mByte)
func bitPower(n int) (res int64) {
	res = 1
	for j := 0; j < n; j++ {
		res <<= 10
	}
	return
}

// choose the best unit for the value of bytes
func humanReadble(bytes int64) (weight int64, unit string) {
	const base = 1 << 10

	if bytes > bitPower(4) {
		weight, unit = bitPower(4), "tB"
	} else if bytes > bitPower(3) {
		weight, unit = bitPower(3), "gB"
	} else if bytes > bitPower(2) {
		weight, unit = bitPower(2), "mB"
	} else if bytes > bitPower(1) {
		weight, unit = bitPower(1), "kB"
	} else {
		weight, unit = bitPower(0), "byte"
	}

	return
}

// an progressBar for copying files
func progressBar(maxfiles int, maxbytes int64) {
	byteCount := int64(0)
	var status float64
	fileCount := 0
	w, unit := humanReadble(maxbytes)

	digits := 2
	if unit == "byte" {
		digits = 0
	}

	for {
		select {
		case bytes := <-byteschan:
			byteCount += bytes
		case file := <-fileschan:
			fileCount += file
		}

		if maxbytes != 0 {
			status = float64(byteCount) / float64(maxbytes)
		} else {
			status = float64(fileCount) / float64(maxfiles)
		}

		bar := makebar(status)
		message := fmt.Sprintf(
			"\r %v / %v files [%-*v] %.*f / %.*f %v | %.2f ",
			fileCount, maxfiles, BarLength, bar,
			digits, float32(byteCount/w), digits, float32(maxbytes/w),
			unit, 100*status,
		)
		fmt.Printf(message + "%%")
		flushStdout(len(message))

		if byteCount == maxbytes && fileCount == maxfiles {
			break
		}
	}

	donebar <- true
}

// a-head of copy get the total of files and bytes to copy
// used to make a cool progress bar
func getMaxValues(files []string) (maxfiles int, maxbytes int64) {
	for _, file := range files {
		_ = filepath.Walk(file, func(name string, fi os.FileInfo, err error) error { // don't ignore this!!!
			if err != nil {
				log.Printf("%v: %v\n", name, err)
				return err
			}
			if !fi.IsDir() {
				maxfiles += 1
				maxbytes += fi.Size()
			}

			return nil
		})
	}

	return maxfiles, maxbytes

}

func main() {
	todir := false
	if flag.NArg() < 2 {
		usage()
	}

	files := flag.Args()
	tocopy := files[:len(files)-1]
	lf := files[len(files)-1]
	lfdir, err := os.Stat(lf)
	if err == nil {
		todir = lfdir.IsDir()
	} else if os.IsNotExist(err) {
		todir = true
		cwd, err := os.Getwd()
		if err != nil {
			log.Printf("can't get the current directory: %v", err)
		}
		createDir(cwd, lf)
	} else {
		log.Printf("can't open the %v file: %v", err)
	}
	if flag.NArg() > 2 && todir == false {
		log.Printf("not a directory: %s\n", lf)
		os.Exit(1)
	}

	offchan = make(chan int64, 0)
	zerochan = make(chan int, 0)
	go nextOff()

	if config.progress {
		byteschan = make(chan int64, config.nwork)
		fileschan = make(chan int, 0)
		donebar = make(chan bool, 0)
		maxfiles, maxbytes := getMaxValues(tocopy)
		go progressBar(maxfiles, maxbytes)
		defer fmt.Printf("\n") // on the end print a new line for the progbar
	}

	for _, file := range tocopy {
		err := copy(file, lf, todir)
		if err != nil {
			log.Printf("%v: %v", cmd, err.Error())
			os.Exit(2)
		}
	}

	if config.progress {
		// wait to progress send the message of finish
		<-donebar
	}
}
