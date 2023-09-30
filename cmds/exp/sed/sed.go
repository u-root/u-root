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
	ftmp     os.File
	filename string
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
				// Replaces only the first occurrence which golang stdlib does not offer!
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
	const MinSearchExprParts = 3
	const MaxSearchExprParts = 4
	var ts transforms
	for _, e := range expr {
		switch e[:1] {
		case "s":
			separator := e[1:2]
			eParts := strings.Split(e, separator)
			if len(eParts) < MinSearchExprParts || len(eParts) > MaxSearchExprParts {
				return nil, fmt.Errorf("unable to parse transformation. This should be of the form s/old/new/")
			}
			global := false
			if len(eParts) == MaxSearchExprParts {
				global = strings.Contains(eParts[3], "g")
			}
			ts = append(ts, transform{from: regexp.MustCompile(eParts[1]), to: eParts[2], global: global})
		default:
			return nil, fmt.Errorf("unsupported sed expression")
		}
	}
	return ts, nil
}

type params struct {
	expr    []string
	inplace bool
}

type sedCommand struct {
	rc io.Reader
	wc io.Writer
	ts transforms
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

func command(stdin io.ReadCloser, stdout, stderr io.Writer, p params, args []string) *cmd {
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
func (c *cmd) sed(sc *sedCommand) error {
	r := sc.rc
	for i := range sc.ts {
		// r = t.run(r)
		r = sc.ts[i].run(r)
	}
	_, err := io.Copy(sc.wc, r)
	if err != nil {
		return err
	}
	return nil
}

func (c *cmd) run() error {
	defer c.stdout.Flush()
	var ts transforms
	ts, err := ts.parse(c.params.expr)
	if err != nil {
		return fmt.Errorf("error parsing expressions: %w", err)
	}

	if len(c.args) == 0 {
		sc := &sedCommand{c.stdin, c.stdout, ts}
		return c.sed(sc)
	}
	for i := range c.args {
		fi, err := os.Open(c.args[i])
		if err != nil {
			return fmt.Errorf("unable to open input file: %v", err)
		}
		if c.inplace {
			fo, err := newTmpWriter(fi.Name())
			if err != nil {
				return fmt.Errorf("unable to open output file: %v", err)
			}
			err = c.sed(&sedCommand{fi, fo, ts})
			if err != nil {
				return err
			}
			fi.Close()
			fo.Close()
		} else {
			fo := c.stdout
			err = c.sed(&sedCommand{fi, fo, ts})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
