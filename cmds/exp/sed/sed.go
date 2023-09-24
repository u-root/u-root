package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

type transform struct {
	from string
	to   string
}

type transforms []transform

func (t *transforms) String() string {
	return fmt.Sprint(*t)
}

func (t *transforms) Set(value string) error {
	transformDefinition := strings.Split(value, "/")

	if len(transformDefinition) != 3 || transformDefinition[0] != "s" {
		return fmt.Errorf("unable to parse transformation. This should be of the form s/old/new")
	}

	*t = append(*t, transform{from: transformDefinition[1], to: transformDefinition[2]})

	return nil
}

type config struct {
	transforms transforms
	inplace    bool
}

var cfg = config{}

func init() {
	flag.Var(&cfg.transforms, "e", "search/replace commands (s/old/new)")
	flag.BoolVar(&cfg.inplace, "i", false, "edit files in place")
}

func transformCopy(cfg config, readStreams []io.ReadCloser, writeStreams []io.WriteCloser) {
	for i := range readStreams {
		r := bufio.NewScanner(readStreams[i])
		w := bufio.NewWriter(writeStreams[i])

		for r.Scan() {
			line := r.Text()
			for _, transform := range cfg.transforms {
				line = strings.ReplaceAll(line, transform.from, transform.to)
			}
			_, err := fmt.Fprintf(w, "%s\n", line)
			if err != nil {
				fmt.Printf("unable to write output: %#v", err)
				os.Exit(1)
			}
		}

		if err := w.Flush(); err != nil {
			panic(err)
		}
	}
}

func createWriteStreams(cfg config, readStreams []io.ReadCloser) []io.WriteCloser {
	l := len(readStreams)
	writeStreams := make([]io.WriteCloser, l)
	for i := range writeStreams {
		if cfg.inplace {
			fiName := readStreams[i].(*os.File).Name()
			fiDir := path.Dir(fiName)
			ftmp, err := os.CreateTemp(fiDir, "sed*.txt")
			if err != nil {
				fmt.Printf("unable to create temp file: %#v\n", err)
				os.Exit(1)
			}
			writeStreams[i] = ftmp
		} else {
			writeStreams[i] = os.Stdout
		}
	}
	return writeStreams
}

func main() {
	flag.Parse()

	var readStreams []io.ReadCloser
	if len(flag.Args()) == 0 {
		readStreams = append(readStreams, os.Stdin)
	} else {
		for _, filename := range flag.Args() {
			fh, err := os.Open(filename)
			if err != nil {
				fmt.Printf("unable to open input stream: %s. %#v\n", filename, err)
				os.Exit(1)
			}
			readStreams = append(readStreams, fh)
		}
	}
	writeStreams := createWriteStreams(cfg, readStreams)
	transformCopy(cfg, readStreams, writeStreams)

	for i := range readStreams {
		fi := readStreams[i]
		fi.Close()
		if cfg.inplace {
			fo := writeStreams[i]
			fo.Close()
			fiName := fi.(*os.File).Name()
			foName := fo.(*os.File).Name()
			err := os.Rename(foName, fiName)
			if err != nil {
				fmt.Printf("uname to rename output file: %#v", err)
			}
		}
	}
}
