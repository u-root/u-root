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
	var fail bool
	f, err := os.Create("asm.go")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal("Closing %v: %v", f.Name(), err)
		}
	}()

	if _, err := fmt.Fprintf(f, "package kexec\nvar purgatories = map[string]*purgatory {\n"); err != nil {
		log.Fatal(err)
	}
	for _, asm := range asms {
		// assemble, then create the hexdump of the elf.
		d, err := ioutil.TempDir("", "kexecgen")
		if err != nil {
			log.Fatalf("%v", err)
		}

		//defer os.RemoveAll(d)
		src := filepath.Join(d, asm.name) + ".S"
		dst := filepath.Join(d, asm.name) + ".o"
		out := filepath.Join(d, asm.name) + ".out"
		if err := ioutil.WriteFile(src, []byte(asm.code), 0666); err != nil {
			log.Fatal(err)
		}

		if out, err := exec.Command(asm.cc[0], append(asm.cc[1:], "-nostdinc", "-nostdlib", "-o", dst, src)...).CombinedOutput(); err != nil {
			log.Printf("%s, %s: %v, %v", src, asm.code, string(out), err)
			continue
		}
		if out, err := exec.Command(asm.ld[0], append(asm.ld[1:], "-nostdinc", "-nostdlib", "-o", out, dst)...).CombinedOutput(); err != nil {
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
		if _, err := fmt.Fprintf(f, "\t\"%s\": &purgatory{\n\tname: \"%s\",\n\thexdump: \n`%s\n`,\n\tcode: %s\n},\n", asm.name, asm.name, string(b), codehex); err != nil {
			log.Fatal(err)
		}
	}
	if fail {
		log.Fatal("There was at least one error")
	}
	if _, err := fmt.Fprintf(f, "}\n"); err != nil {
		log.Printf("Writing final brace: %v", err)
	}
}
