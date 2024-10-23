// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Scp copies files between hosts on a network.
//
// Synopsis:
//
//	scp [-t|-f] [FILE]
//
// Description:
//
//	If -t is given, decode SCP protocol from stdin and write to FILE.
//	If -f is given, stream FILE over SCP protocol to stdout.
//
// Options:
//
//	-t: Act as the target
//	-f: Act as the source
//	-v: Passed if SCP is verbose, ignored
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

const (
	SUCCESS = 0
)

var (
	isTarget = flag.Bool("t", false, "Act as the target")
	isSource = flag.Bool("f", false, "Act as the source")
	_        = flag.Bool("v", false, "Ignored")
)

func scpSingleSource(w io.Writer, r io.Reader, pth string) error {
	f, err := os.Open(pth)
	if err != nil {
		return err
	}
	defer f.Close()
	s, err := f.Stat()
	if err != nil {
		return err
	}
	filename := path.Base(pth)
	fmt.Fprintf(w, "C0%o %d %s\n", s.Mode(), s.Size(), filename)
	if response(r) != SUCCESS {
		return fmt.Errorf("response was not success")
	}
	_, err = io.Copy(w, f)
	if err != nil {
		return fmt.Errorf("copy error: %w", err)
	}
	reply(w, SUCCESS)

	if response(r) != SUCCESS {
		return fmt.Errorf("response was not success")
	}
	return nil
}

func scpSingleSink(w io.Writer, r io.Reader, path string) error {
	var mode os.FileMode
	var size int64
	filename := ""

	// Ignore the filename, assume it has been provided on the command line.
	// This will not work with directories and recursive copy, but that's not
	// supported right now.
	if _, err := fmt.Fscanf(r, "C0%o %d %s\n", &mode, &size, &filename); err != nil {
		if err == io.ErrUnexpectedEOF {
			return io.EOF
		}
		return fmt.Errorf("fscanf: %w", err)
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, mode)
	if err != nil {
		return fmt.Errorf("open error: %w", err)
	}
	reply(w, SUCCESS)
	defer f.Close()

	_, err = io.CopyN(f, r, size)
	if err != nil {
		return fmt.Errorf("copy error: %w", err)
	}
	if response(r) != SUCCESS {
		return fmt.Errorf("response was not success")
	}
	reply(w, SUCCESS)
	return nil
}

func scpSource(w io.Writer, r io.Reader, path string) error {
	// Sink->Source is started with a response
	if response(r) != SUCCESS {
		return fmt.Errorf("response was not success")
	}
	return scpSingleSource(w, r, path)
}

func scpSink(w io.Writer, r io.Reader, path string) error {
	reply(w, SUCCESS)
	for {
		if err := scpSingleSink(w, r, path); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	return nil
}

func reply(out io.Writer, r byte) {
	out.Write([]byte{r})
}

func response(in io.Reader) byte {
	b := make([]byte, 1)
	in.Read(b)
	return b[0]
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatalf("no file provided")
	}

	if *isSource == *isTarget {
		log.Fatalf("-t or -f needs to be supplied, and not both")
	}

	if *isSource {
		if err := scpSource(os.Stdout, os.Stdin, flag.Args()[0]); err != nil {
			log.Fatalf("scp: %v", err)
		}
	} else if *isTarget {
		if err := scpSink(os.Stdout, os.Stdin, flag.Args()[0]); err != nil {
			log.Fatalf("scp: %v", err)
		}
	}
}
