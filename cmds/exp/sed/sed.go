// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// By Ahmed Kamal <email.ahmedkamal@googlemail.com>

// sed edits file contents using regular expressions.
//
// Synopsis:
//
//	sed [-ie] [FILE]...
//
// Options:
//
//	-i, --in-place        in-place file edit
//	-e, --expression      edits to perform (only s/foo/bar style supported)

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	flag "github.com/spf13/pflag"
)

// Writes to a temp file first
// Renames to the final file upon closing
type tmpWriter struct {
	filename string
	ftmp     os.File
}

func newTmpWriter(filename string) (*tmpWriter, error) {
	ftmp, err := os.CreateTemp("/tmp", ".sed*.txt")
	if err != nil {
		return nil, fmt.Errorf("unable to create temp file: %w", err)
	}
	return &tmpWriter{filename: filename, ftmp: *ftmp}, nil
}

func (tw *tmpWriter) Write(b []byte) (int, error) {
	return tw.ftmp.Write(b)
}

func (tw *tmpWriter) Close() error {
	err := tw.ftmp.Close()
	if err != nil {
		return err
	}
	os.Rename(tw.ftmp.Name(), tw.filename)
	return nil
}

type transform struct {
	from   *regexp.Regexp
	to     string
	global bool
}

func (t *transform) run(input io.Reader) io.Reader {
	pr, pw := io.Pipe()
	scanner := bufio.NewScanner(input)
	go func() {
		defer pw.Close()
		for scanner.Scan() {
			line := scanner.Text()
			var replaced string
			if t.global {
				replaced = t.from.ReplaceAllString(line, t.to)
			} else {
				counter := 0
				// Replaces only the first occurence which golang stdlib does not offer!
				replaced = t.from.ReplaceAllStringFunc(line, func(value string) string {
					if counter == 1 {
						return value
					}
					counter++
					return t.from.ReplaceAllString(value, t.to)
				})
			}
			pw.Write([]byte(replaced))
			pw.Write([]byte("\n")) // We unconditionally send \n after every line which could be wrong on the last line!
		}
	}()
	return pr
}

type transforms []transform

func (t *transforms) String() string {
	return fmt.Sprint(*t)
}

func (t *transforms) parse(expr []string) (transforms, error) {
	var ts transforms
	for _, e := range expr {
		switch e[:1] {
		case "s":
			separator := e[1:2]
			e_parts := strings.Split(e, separator)
			if len(e_parts) < 3 || len(e_parts) > 4 {
				return nil, fmt.Errorf("unable to parse transformation. This should be of the form s/old/new/")
			}
			global := false
			if len(e_parts) == 4 {
				global = strings.Contains(e_parts[3], "g")
			}
			ts = append(ts, transform{from: regexp.MustCompile(e_parts[1]), to: e_parts[2], global: global})
		}
	}
	return ts, nil
}

type params struct {
	expr    []string
	inplace bool
}

type sedCommand struct {
	rc   io.ReadCloser
	wc   io.WriteCloser
	name string
}

func parseParams() params {
	p := params{}
	var e []string
	flag.StringArrayVarP(&p.expr, "expression", "e", e, "Expression to execute (only search/replace currently supported)")
	flag.BoolVarP(&p.inplace, "in-place", "i", false, "Perform file edits in-place")
	flag.Parse()

	return p
}

func main() {

	if err := command(os.Stdin, os.Stdout, os.Stderr, parseParams(), flag.Args()).run(); err != nil {
		if err != nil {
			log.Fatal(err)
		}
	}
}

// cmd contains the actually business logic of sed
type cmd struct {
	stdin  io.ReadCloser
	stdout *bufio.Writer
	stderr io.Writer
	args   []string
	params
}

func command(stdin io.ReadCloser, stdout io.Writer, stderr io.Writer, p params, args []string) *cmd {
	return &cmd{
		stdin:  stdin,
		stdout: bufio.NewWriter(stdout),
		stderr: stderr,
		params: p,
		args:   args,
	}
}

// sed reads data from the os.File embedded in sedCommand.
// It matches each line against the re and prints the matching result
// If we are only looking for a match, we exit as soon as the condition is met.
// "match" means result of re.Match == match flag.
func (c *cmd) sed(sc *sedCommand, re *regexp.Regexp) error {
	// r := bufio.NewScanner(sc.rc)
	// defer sc.rc.Close()
	// for r.Scan() {
	// 	line := r.Text()
	// }
	// c.stdout.Flush()
	io.Copy(sc.wc, sc.rc)
	return nil
}

func (c *cmd) run() error {
	defer c.stdout.Flush()
	// sed business logic here
	// var ts transforms
	if c.params.expr != nil {

	}
	sc := &sedCommand{rc: os.Stdin, wc: os.Stdout, name: "foo"}
	c.sed(sc, regexp.MustCompile(`\d+`))
	return nil
}
