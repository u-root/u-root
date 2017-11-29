/*
Command testpaths generates files with one path per line.  The first argument is the output directory and the second
must be one of 'standard', 'branchFactor', 'segmentCount', or 'segmentSize'.  Each additional (integer) argument
generates a .txt file by the same name.
*/
package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dgryski/go-metro"
)

// A writePathFn should write one path per line to w.
type writePathFn func(w io.Writer, arg int64)

var writePathFns = map[string]writePathFn{
	// Writes count routes of a 'standard' form.
	"standard": func(w io.Writer, count int64) {
		for i := int64(0); i < count; i++ {
			path := generateRoute(uint64(i))
			if _, err := io.WriteString(w, path); err != nil {
				log.Fatalf("failed to write path #%d %q: %s", i, path, err)
			}
			if _, err := w.Write([]byte{'\n'}); err != nil {
				log.Fatal("failed to write new line:", err)
			}
		}
	},

	// Writes count routes with different values for the same component.
	"branchFactor": func(w io.Writer, count int64) {
		for i := int64(0); i < count; i++ {
			sum := md5.Sum([]byte{byte(i >> 16), byte(i), byte(i >> 24), byte(i >> 8)})
			if _, err := fmt.Fprintf(w, "/%x", sum); err != nil {
				log.Fatalf("failed to write path #%d: %s", i, err)
			}
			if _, err := w.Write([]byte{'\n'}); err != nil {
				log.Fatal("failed to write new line:", err)
			}
		}
	},

	// Writes a route with count segments.
	"segmentCount": func(w io.Writer, count int64) {
		path := strings.Repeat("/test", int(count))
		if _, err := io.WriteString(w, path); err != nil {
			log.Fatalf("failed to write path %d", count)
		}
	},

	// Writes a route with 5 counts sized segments.
	"segmentSize": func(w io.Writer, size int64) {
		path := strings.Repeat("/"+strings.Repeat("a", int(size)), 5)
		if _, err := io.WriteString(w, path); err != nil {
			log.Fatalf("failed to write path %d", size)
		}
	},
}

func main() {
	if len(os.Args) < 4 {
		log.Fatal("too few arguments")
	}

	dir := os.Args[1]
	name := os.Args[2]
	nameDir := filepath.Join(dir, name)
	if err := os.MkdirAll(nameDir, 0777); err != nil {
		log.Fatal("failed to create diretory:", err)
	}

	writePath, ok := writePathFns[name]
	if !ok {
		log.Fatal("unrecognized name:", name)
	}

	for _, arg := range os.Args[3:] {
		routes, err := strconv.ParseInt(arg, 10, 0)
		if err != nil {
			log.Fatal("failed to parse argument:", err)
		}

		createFile(filepath.Join(nameDir, arg+".txt"), routes, writePath)
	}
}

func createFile(name string, arg int64, writePath writePathFn) {
	f, err := os.Create(name)
	if err != nil {
		log.Fatal("failed to create file:", err)
	}
	defer f.Close()
	writePath(f, arg)
}

const (
	seed1 = 33928
	seed2 = 9274035
)

// generateRoute generates a deterministic test route based on u.
// Routes are 32 bytes printed as hex, and split into 1, 2, 4, or 8 parts.
// ~1/5 of the parts after the first are parameters. ~1/2 of suffix parameters are catch all.
// Example: /61aab75e6134555e/:param/20c98ae9da239a03/*wildcard
//
// Verified conflict free up to: 1,000,000
func generateRoute(u uint64) string {
	s := make([]byte, 0, 32)
	h1, h2 := metro.Hash128([]byte{byte(u >> 24), byte(u >> 8), byte(u), byte(u >> 16)}, seed1)
	h3, h4 := metro.Hash128([]byte{byte(u), byte(u >> 16), byte(u >> 8), byte(u >> 24)}, seed2)
	for _, h := range []uint64{h1, h2, h3, h4} {
		var i uint
		for ; i < 64; i += 8 {
			s = append(s, byte(h>>i))
		}
	}

	parts := 1 << (u % 4)  // [1,2,4,8]
	each := len(s) / parts // [32,16,8,4]

	var b bytes.Buffer
	for i := 0; i < parts; i++ {
		b.WriteRune('/')

		// if not first, 1/5 chance of param
		if i > 0 && metro.Hash64(b.Bytes(), seed1)%5 == 0 {
			// if last, 1/2 chance of *
			if i == parts-1 && u%2 == 0 {
				b.WriteString(`*wildcard`)
			} else {
				b.WriteString(`:param`)
			}
		} else {
			nameStart := i * each
			fmt.Fprintf(&b, "%x", s[nameStart:nameStart+each])
		}
	}

	return b.String()
}
