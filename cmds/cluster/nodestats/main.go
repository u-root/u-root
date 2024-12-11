// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

// nodestats prints out vital statistics about a node as JSON.
// It currently uses the jaypipes/ghw package, as well as
// files in /sys and /proc.
// Any errors encountered are recorded in the stats struct.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"

	"github.com/jaypipes/ghw"
	"github.com/u-root/u-root/pkg/cluster/health"
)

func node() *health.Stat {
	// Some *packages* write to Stderr instead of returning
	// an error. Further, the error may be a warning.
	// Gather up os.Stderr via a pipe and return it in the health.Stat struct.
	// Any error will be gathered in the errors.Join.
	// It is important to save the old os.Stderr in the event
	// log.Fatal is called in main().
	r, w, errs := os.Pipe()
	f := os.Stderr
	os.Stderr = w
	defer func() {
		os.Stderr = f
	}()

	hn, err := os.Hostname()
	errors.Join(errs, err)

	host, err := ghw.Host()
	errs = errors.Join(errs, err)

	k := health.Kernel{}
	val := reflect.ValueOf(k)
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		n, ok := field.Tag.Lookup("file")
		if !ok {
			errs = errors.Join(errs, fmt.Errorf("%s:%w", field.Name, os.ErrNotExist))
			continue
		}

		dat, err := os.ReadFile(n)
		errs = errors.Join(errs, err)
		reflect.ValueOf(&k).Elem().Field(i).SetString(string(dat))
	}

	// ReadAll would be a bit dangerous in this context.
	// Read a reasonable amount, and record if we did not get
	// it all.
	w.Close()
	var Stderr [65536]byte
	n, err := r.Read(Stderr[:])
	if err != nil && err != io.EOF {
		errs = errors.Join(errs, fmt.Errorf("stderr read %d bytes, got %w but not io.EOF or nil", n, err))
	}

	stats := &health.Stat{Hostname: hn, Info: host, Kernel: k, Stderr: string(Stderr[:n])}
	if errs != nil {
		stats.Err = errs.Error()
	}
	return stats
}

func run(out io.Writer, args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("%v:%w", args[0], os.ErrInvalid)
	}
	stats := node()
	j, err := json.MarshalIndent(stats, "", "\t")
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "%s\n", string(j))
	return nil
}

func main() {
	if err := run(os.Stdout, os.Args); err != nil {
		log.Fatal(err)
	}
}
