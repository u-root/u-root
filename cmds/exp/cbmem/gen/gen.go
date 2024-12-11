// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"
)

var apu2 = flag.Bool("apu2", false, "Use hardcoded values for an APU2 -- very iffy")

// This is a known set of offsets for the APU2.
// It is very unlikely to be widely applicable but it
// is useful for a bootstrap. Not recommended for wide use.
func genAPU2(n string) ([]byte, error) {
	// mmap(NULL, 4096, PROT_READ, MAP_SHARED, 3, 0) = 0x7fa13d64f000
	// mmap(NULL, 1320, PROT_READ, MAP_SHARED, 3, 0) = 0x7fa13d64e000
	// mmap(NULL, 4096, PROT_READ, MAP_SHARED, 3, 0x77fae000) = 0x7fa13d64d000
	// mmap(NULL, 392, PROT_READ, MAP_SHARED, 3, 0x77fae000) = 0x7fa13d64c000
	// mmap(NULL, 8, PROT_READ, MAP_SHARED, 3, 0x77fdf000) = 0x7fa13d64f000
	// mmap(NULL, 65, PROT_READ, MAP_SHARED, 3, 0x77fdf000) = 0x7fa13d64f000
	f, err := os.Open("/dev/mem")
	if err != nil {
		return nil, err
	}
	out := bytes.NewBuffer([]byte("package main\nvar apu2 = []seg {\n"))
	for _, m := range []struct {
		offset int64
		size   int64
	}{
		{size: 4096, offset: 0},
		{size: 1320, offset: 0},
		{size: 4096, offset: 0x77fae000},
		{size: 392, offset: 0x77fae000},
		{size: 8, offset: 0x77fdf000},
		{size: 65, offset: 0x77fdf000},
	} {
		b, err := syscall.Mmap(int(f.Fd()), m.offset, int(m.size), syscall.PROT_READ, syscall.MAP_SHARED)
		if err != nil {
			return nil, fmt.Errorf("mmap %d bytes at %#x: %w", m.size, m.offset, err)
		}
		fmt.Fprintf(out, "{off: %#x, dat:[]byte{\n", m.offset)
		for i := 0; i < len(b); i += 8 {
			fmt.Fprintf(out, "\t/*%#08x*/\t", m.offset+int64(i))
			for j := 0; j < 8 && i+j < len(b); j++ {
				c := b[i+j]
				fmt.Fprintf(out, "%#02x/*%q*/, ", c, c)
			}
			fmt.Fprintf(out, "\n")
		}
		fmt.Fprintf(out, "\t},},\n")

	}
	fmt.Fprintf(out, "}\n")
	return out.Bytes(), nil
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatal("Usage: gen [-apu2] outputfile.go")
	}
	if *apu2 {
		s, err := genAPU2("apu2")
		if err != nil {
			log.Fatal(err)
		}
		if err := os.WriteFile("apu2.go", []byte(s), 0o644); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	// This will only be useful when we have some working dump code.
	// a := flag.Args()
	// if len(a) != 1 {
	// 	log.Fatal("usage: %s platform-name")
	// }
	// h, err := os.Hostname()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// gen(h)
	log.Fatal("not yet")
}
