package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
)

var Nwork int = 1

const Defb = 8192

var offchan chan int64
var zerochan chan int

func copyfile(from, to string, todir bool) bool {
	if todir {
		_, file := path.Split(from)
		to = to + "/" + file
	}

	dirb, err := os.Stat(from)
	if err != nil {
		fmt.Printf("can't stat %s: %v\n", from, err)
		return true
	}
	tob, err := os.Stat(to)
	if err == nil {
		if sameFile(dirb.Sys(), tob.Sys()) {
			fmt.Printf("%s and %s are the same file\n", from, to)
			return true
		}
	}

	if dirb.IsDir() {
		fmt.Printf("%s is a directory\n", from)
		return true
	}

	mode := dirb.Mode() & 0777
	f, err := os.Open(from)
	if err != nil {
		fmt.Printf("can't open %s: %v\n", from, err)
		return true
	}
	defer f.Close()

	t, err := os.OpenFile(to, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		fmt.Printf("can't create %s: %v\n", to, err)
		f.Close()
		return true
	}
	defer t.Close()
	return copy1(f, t, from, to)
}

func copy1(f, t *os.File, from, to string) (ret bool) {
	zerochan <- 0
	fail := make(chan bool, Nwork)

	for i := 0; i < Nwork; i++ {
		go worker(f, t, from, to, fail)
	}
	for i := 0; i < Nwork; i++ {
		end := <-fail
		if end == true {
			ret = true
		}
	}
	return
}

func worker(f, t *os.File, from, to string, fail chan bool) {
	var buf [Defb]byte
	var bp []byte

	l := len(buf)
	bp = buf[0:]
	o := <-offchan
	for {
		n, err := f.ReadAt(bp, o)
		if err != nil && err != io.EOF {
			fmt.Printf("reading %s at %v: %v\n", from, o, err)
			fail <- true
			return
		}
		if n == 0 {
			break
		}

		nb := bp[0:n]
		n, err = t.WriteAt(nb, o)
		if err != nil {
			fmt.Printf("writing %s: %v\n", to, err)
			fail <- true
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
	fail <- false
}

func nextoff() {
	off := int64(0)
	for {
		select {
		case <-zerochan:
			off = 0
		case offchan <- off:
			off += Defb
		}
	}
}

func usage() {
	fmt.Printf("usage: cp [-w workers] from to\n")
	os.Exit(1) // sysfatal
}

var nwork = flag.Int("w", 16, "number of worker goroutines")

func main() {
	todir := false

	flag.Parse()
	Nwork = *nwork
	if flag.NArg() < 2 {
		usage()
	}

	files := flag.Args()
	lf := files[len(files)-1]
	lfdir, err := os.Stat(lf)
	if err == nil {
		todir = lfdir.IsDir()
	}
	if flag.NArg() > 2 && todir == false {
		fmt.Printf("not a directory: %s\n", lf)
		os.Exit(1) // sysfatal
	}

	offchan = make(chan int64, 0)
	zerochan = make(chan int, 0)
	go nextoff()

	failed := false
	for i := 0; i < flag.NArg()-1; i++ {
		if copyfile(files[i], lf, todir) {
			failed = true
		}
	}
	if failed {
		os.Exit(2)
	}
	return
}
