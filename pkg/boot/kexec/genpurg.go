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

	if _, err := fmt.Fprintf(f, "package kexec\nvar purgatories = map[string]*purgatory {\n\t\"default\": linuxPurgatory,\n"); err != nil {
		log.Fatal(err)
	}
	for _, asm := range asms {
		// assemble, then create the hexdump of the elf.
		d, err := ioutil.TempDir("", "kexecgen")
		if err != nil {
			log.Fatalf("%v", err)
		}

		defer os.RemoveAll(d)
		src := filepath.Join(d, asm.name) + ".S"
		dst := filepath.Join(d, asm.name) + ".."
		if err := ioutil.WriteFile(src, []byte(asm.code), 0666); err != nil {
			log.Fatal(err)
		}

		if out, err := exec.Command(asm.args[0], append(asm.args[1:], "-nostdinc", "-nostdlib", "-o", dst, src)...).CombinedOutput(); err != nil {
			log.Printf("%s, %s: %v, %v", src, asm.code, string(out), err)
			continue
		}
		code, err := ioutil.ReadFile(dst)
		if err != nil {
			log.Fatal(err)
		}
		codehex := "\t[]byte{\n"
		var i int
		for i < len(code) {
			codehex += "\t"
			for _, c := range code[i : i+16] {
				codehex = codehex + fmt.Sprintf("%#02x, ", c)
			}
			codehex += "\n"
			i += len(code[i : i+16])
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
