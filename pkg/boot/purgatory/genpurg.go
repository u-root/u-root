// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	f, err := os.Create("asm.go")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal("Closing %v: %v", f.Name(), err)
		}
	}()

	if _, err := fmt.Fprintf(f, "// Copyright 2022 the u-root Authors. All rights reserved\n// Use of this source code is governed by a BSD-style\n// license that can be found in the LICENSE file.\n\n"); err != nil {
		log.Fatal(err)
	}
	if _, err := fmt.Fprintf(f, "package purgatory\nvar Purgatories = map[string]*Purgatory {\n"); err != nil {
		log.Fatal(err)
	}
	d, err := ioutil.TempDir("", "kexecgen")
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer os.RemoveAll(d)

	for _, asm := range asms {
		// assemble, then create the hexdump of the elf.
		src := filepath.Join(d, asm.name) + ".S"
		dst := filepath.Join(d, asm.name) + ".o"
		out := filepath.Join(d, asm.name) + ".out"
		if err := ioutil.WriteFile(src, []byte(asm.code), 0o666); err != nil {
			log.Fatal(err)
		}

		if out, err := exec.Command(asm.cc[0], append(asm.cc[1:], "-nostdinc", "-nostdlib", "-o", dst, src)...).CombinedOutput(); err != nil {
			log.Printf("%s, %s: %v, %v", src, asm.code, string(out), err)
			continue
		}
		if out, err := exec.Command(asm.ld[0], append(asm.ld[1:], "-nostdlib", "-o", out, dst)...).CombinedOutput(); err != nil {
			log.Printf("%s, %s: %v, %v", dst, asm.code, string(out), err)
			continue
		}
		code, err := ioutil.ReadFile(out)
		if err != nil {
			log.Fatal(err)
		}
		if len(code) == 0 {
			log.Fatalf("%s: no output: dir %v", asm.name, d)
		}
		codehex := "\t[]byte{\n"
		var i int
		for i < len(code) {
			codehex += "\t"
			for j := 0; j < 16 && i < len(code); j++ {
				codehex = codehex + fmt.Sprintf("%#02x, ", code[i])
				i++
			}
			codehex += "\n"
		}
		codehex += "\n},\t\n"

		var buf bytes.Buffer
		if _, err := io.Copy(hex.Dumper(&buf), bytes.NewBuffer(code)); err != nil {
			log.Fatal(err)
		}
		b := buf.Bytes()
		for i := range b {
			if b[i] == '`' {
				b[i] = '.'
			}
		}
		if _, err := fmt.Fprintf(f, "\t\"%s\": &Purgatory{\n\tName: \"%s\",\n\tHexdump: \n`%s\n`,\n\tCode: %s\n},\n", asm.name, asm.name, string(b), codehex); err != nil {
			log.Fatal(err)
		}
	}
	if _, err := fmt.Fprintf(f, "}\n"); err != nil {
		log.Printf("Writing final brace: %v", err)
	}
}
