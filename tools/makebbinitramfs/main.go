// make_intramfs creates a CPIO file for booting Linux. Each package becomes a
// symlink to the busybox located at /bin/bb.
package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/uflag"
)

var defaultRamfs = []cpio.Record{
	cpio.Directory("bbin", 0755),
	cpio.Directory("dev", 0755),
	cpio.CharDev("dev/console", 0600, 5, 1),
	cpio.CharDev("dev/tty", 0666, 5, 0),
	cpio.CharDev("dev/null", 0666, 1, 3),
	cpio.CharDev("dev/port", 0640, 1, 4),
	cpio.CharDev("dev/urandom", 0666, 1, 9),
}

var (
	bb          = flag.String("bb", "", "Busybox executable")
	out         = flag.String("out", "", "Output CPIO filename")
	compression = flag.String("compression", "none", "Compression the CPIO (gzip or none)")
	cmdNames    uflag.Strings
)

func init() {
	flag.Var(&cmdNames, "cmd_name", "Command names needing to be symlinked")
}

type flusher interface {
	io.Writer
	Flush() error
}

func flushSafe(f flusher) {
	if err := f.Flush(); err != nil {
		log.Fatal(err)
	}
}

func closeSafe(f io.Closer) {
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	// Open the file
	f, err := os.Create(*out)
	if err != nil {
		log.Fatal(err)
	}
	defer closeSafe(f)
	var bw flusher = bufio.NewWriter(f)
	defer flushSafe(bw)

	// Select compression
	switch *compression {
	case "gzip":
		gz, err := gzip.NewWriterLevel(bw, gzip.BestCompression)
		if err != nil {
			log.Fatal(err)
		}
		defer closeSafe(gz) // adds gzip footer
		bw = gz
	case "none":
	default:
		log.Fatalf("invalid compression method '%s'", *compression)
	}

	// Create a CPIO writer.
	w := cpio.Newc.Writer(bw)

	writeRecord := func(rec cpio.Record) {
		if err := w.WriteRecord(rec); err != nil {
			log.Fatalf("could not write record %q: %v", rec.Name, err)
		}
	}

	// Create default records.
	for _, rec := range defaultRamfs {
		writeRecord(rec)
	}

	// Add the busybox binary.
	data, err := ioutil.ReadFile(*bb)
	if err != nil {
		log.Fatal(err)
	}
	writeRecord(cpio.StaticFile("bbin/bb", string(data), 0755))

	// Create symlinks.
	for _, cmdName := range cmdNames {
		writeRecord(cpio.Symlink(filepath.Join("bbin", cmdName), "bb"))
	}
	writeRecord(cpio.Symlink("init", "bbin/init"))

	// Add trailer.
	if err := cpio.WriteTrailer(w); err != nil {
		log.Fatalf("could not write trailer: %v", err)
	}
}
