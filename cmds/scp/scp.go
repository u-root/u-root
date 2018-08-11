// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Scp copies files between hosts on a network.
//
// Synopsis:
//     scp [-t|-f] [FILE]
//
// Description:
//     If -t is given, decode SCP protocol from stdin and write to FILE.
//     If -f is given, stream FILE over SCP protocol to stdout.
//
// Options:
//     -t: Act as the target
//     -f: Act as the source
//     -v: Passed if SCP is verbose, ignored
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	SUCCESS = 0
)

var (
	isTarget = flag.Bool("t", false, "Act as the target")
	isSource = flag.Bool("f", false, "Act as the source")
	_        = flag.Bool("v", false, "Ignored")
)

func scpSource(w io.Writer, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	s, err := f.Stat()
	if err != nil {
		return err
	}
	w.Write([]byte(fmt.Sprintf("C0%o %d %s\n", s.Mode(), s.Size(), path)))
	if response() != SUCCESS {
		log.Fatalf("response was not success")
	}
	_, err = io.Copy(w, f)
	if err != nil {
		log.Fatalf("copy error: %v", err)
	}
	reply(SUCCESS)

	if response() != SUCCESS {
		log.Fatalf("response was not success")
	}
	return nil
}

func scpSink(r io.Reader, path string) error {
	var mode os.FileMode
	var size int64
	filename := ""

	if _, err := fmt.Fscanf(r, "C0%o %d %s\n", &mode, &size, &filename); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, mode)
	if err != nil {
		log.Fatalf("open error: %v", err)
	}
	reply(SUCCESS)
	defer f.Close()

	_, err = io.CopyN(f, r, size)
	if err != nil {
		log.Fatalf("copy error: %v", err)
	}
	reply(SUCCESS)
	return nil
}

func reply(r byte) {
	os.Stdout.Write([]byte{r})
}

func response() byte {
	b := make([]byte, 1)
	os.Stdout.Read(b)
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
		// Sink->Source is started with a response
		if response() != SUCCESS {
			log.Fatalf("response was not success")
		}
		if err := scpSource(os.Stdout, flag.Args()[0]); err != nil {
			log.Fatalf("scp: %v", err)
		}
	} else if *isTarget {
		// Sink->Source starts with a response
		reply(SUCCESS)
		for {
			if err := scpSink(os.Stdin, flag.Args()[0]); err != nil {
				log.Fatalf("scp: %v", err)
			}
		}
	}
}
