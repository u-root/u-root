// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Fmap parses flash maps.
//
// Synopsis:
//     fmap [-s|-c func|-r i] [FILE]
//
// Description:
//     Return 0 if the flash map is valid and 1 otherwise. Detailed information
//     is printed to stderr. If FILE is not specified, read from stdin.
//
//     This implementation is based off of https://github.com/dhendrix/flashmap.
//
// Options:
//     -c func: print checksum using the given `hash function` (md5|sha1|sha256)
//     -r i: read an area from the flash
//     -s: print human readable summary
//     -u: print human readable usage stats
package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"flag"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"text/template"

	fmap "github.com/u-root/u-root/cmds/fmap/lib"
)

var (
	checksum = flag.String("c", "", "print checksum using the given `hash function` (md5|sha1|sha256)")
	read     = flag.Int("r", -1, "read an area from the flash")
	summary  = flag.Bool("s", false, "print human readable summary")
	usage    = flag.Bool("u", false, "print human readable usage summary")
)

var hashFuncs = map[string](func() hash.Hash){
	"md5":    md5.New,
	"sha1":   sha1.New,
	"sha256": sha256.New,
}

// Print human readable summary of the fmap.
func printFMap(f *fmap.FMap, m *fmap.FMapMetadata) {
	const desc = `Fmap found at {{printf "%#x" .Metadata.Start}}:
	Signature:  {{printf "%s" .Signature}}
	VerMajor:   {{.VerMajor}}
	VerMinor:   {{.VerMinor}}
	Base:       {{printf "%#x" .Base}}
	Size:       {{printf "%#x" .Size}}
	Name:       {{printf "%s" .Name}}
	NAreas:     {{.NAreas}}
{{- range $i, $v := .Areas}}
	Areas[{{$i}}]:
		Offset:  {{printf "%#x" $v.Offset}}
		Size:    {{printf "%#x" $v.Size}}
		Name:    {{printf "%s" $v.Name}}
		Flags:   {{printf "%#x" $v.Flags}} ({{FlagNames $v.Flags}})
{{- end}}
`
	t := template.Must(template.New("desc").
		Funcs(template.FuncMap{"FlagNames": fmap.FlagNames}).
		Parse(desc))
	// Combine the two structs to pass into template.
	combined := struct {
		*fmap.FMap
		Metadata *fmap.FMapMetadata
	}{f, m}
	if err := t.Execute(os.Stdout, combined); err != nil {
		log.Fatal(err)
	}
}

func printUsage(r io.Reader) {
	blockSize := 4 * 1024
	rowLength := 32

	buffer := make([]byte, blockSize)
	fullBlock := bytes.Repeat([]byte{0xff}, blockSize)
	zeroBlock := bytes.Repeat([]byte{0x00}, blockSize)

	fmt.Println("Legend: '.' - full (0xff), '0' - zero (0x00), '#' - mixed")

	var numBlocks, numFull, numZero int
	loop: for {
		fmt.Printf("%#08x: ", numBlocks*blockSize)
		for col := 0; col < rowLength; col++ {
			// Read next block.
			_, err := io.ReadFull(r, buffer)
			if err == io.EOF {
				fmt.Print("\n")
				break loop
			} else if err == io.ErrUnexpectedEOF {
				fmt.Printf("\nWarning: flash is not a multiple of %d", len(buffer))
				break loop
			} else if err != nil {
				log.Fatal(err)
			}
			numBlocks++

			// Analyze block.
			if bytes.Equal(buffer, fullBlock) {
				numFull++
				fmt.Print(".")
			} else if bytes.Equal(buffer, zeroBlock) {
				numZero++
				fmt.Print("0")
			} else {
				fmt.Print("#")
			}
		}
		fmt.Print("\n")
	}

	// Print usage statistics.
	print := func(name string, n int) {
		fmt.Printf("%s %d (%.1f%%)\n", name, n,
			float32(n)/float32(numBlocks)*100)
	}
	print("Blocks:      ", numBlocks)
	print("Full (0xff): ", numFull)
	print("Empty (0x00):", numZero)
	print("Mixed:       ", numBlocks-numFull-numZero)
}

var btoi = map[bool]int{
	false: 0,
	true:  1,
}

func main() {
	flag.Parse()

	// Validate flags
	if btoi[*summary]+btoi[*usage]+btoi[*read >= 0]+btoi[*checksum != ""] > 1 {
		log.Fatal("Only use one flag at a time")
	}
	if *checksum != "" {
		if _, ok := hashFuncs[*checksum]; !ok {
			log.Fatal("Not a valid hash function. Must be one of md5, sha1 or sha256")
		}
	}

	// Choose a reader
	r := os.Stdin
	if flag.NArg() == 1 {
		var err error
		r, err = os.Open(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		defer r.Close()
	} else if flag.NArg() > 1 {
		log.Fatal("Too many arguments")
	}

	// Read fmap.
	f, metadata, err := fmap.ReadFMap(r)
	if err != nil {
		log.Fatal(err)
	}

	switch {
	// Optionally print checksum.
	case *checksum != "":
		checksum, err := f.Checksum(r, hashFuncs[*checksum]())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%x\n", checksum)

	// Optionally print area.
	case *read >= 0:
		areaReader, err := f.ReadArea(r, *read)
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.Copy(os.Stdout, areaReader)
		if err != nil {
			log.Fatal(err)
		}

	// Optionally print summary.
	case *summary:
		printFMap(f, metadata)

	// Optionall print usage.
	case *usage:
		if _, err := r.Seek(0, io.SeekStart); err != nil {
			log.Fatal(err)
		}
		printUsage(r)
	}
}
