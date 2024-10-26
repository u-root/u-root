// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"os"

	"github.com/u-root/u-root/pkg/vpd"
)

// NewGetter returns a new initialized Getter.
func NewGetter() *Getter {
	return &Getter{
		R:   vpd.NewReader(),
		Out: os.Stdout,
	}
}

// Getter is the object that gets VPD variables.
type Getter struct {
	R   *vpd.Reader
	Out io.Writer
}

// Print prints VPD variables. If `key` is an empty string, it will print all
// the variables it finds. A variable can exist as both read-write and
// read-only. In that case a "RO" or "RW" string will also be printed.
func (g *Getter) Print(key string) error {
	allVars := make(map[bool]map[string][]byte)
	if key == "" {
		// read-only variables first
		rovars, err := g.R.GetAll(true)
		if err != nil {
			return fmt.Errorf("failed to read RO variables: %w", err)
		}
		allVars[true] = rovars

		// then read-write variables
		rwvars, err := g.R.GetAll(false)
		if err != nil {
			return fmt.Errorf("failed to read RW variables: %w", err)
		}
		allVars[false] = rwvars
	} else {
		// first the read-only var
		value, errRO := g.R.Get(key, true)
		if errRO == nil {
			allVars[true] = map[string][]byte{key: value}
		}
		// then the read-write var
		value, errRW := g.R.Get(key, false)
		if errRW == nil {
			allVars[false] = map[string][]byte{key: value}
		}

		if len(allVars[true])+len(allVars[false]) == 0 {
			fmt.Fprintf(g.Out, "No variable named '%s' found\n", key)
			// if the variable is simply not set, return without error,
			if os.IsNotExist(errRO) && os.IsNotExist(errRW) {
				return nil
			}
			// otherwise print one or both errors.
			return fmt.Errorf("failed to read variable '%s': RO: %w, RW: %w", key, errRO, errRW)
		}
	}
	for k, v := range allVars[true] {
		fmt.Fprintf(g.Out, "%s(RO) => %s\n", k, v)
	}
	for k, v := range allVars[false] {
		fmt.Fprintf(g.Out, "%s(RW) => %s\n", k, v)
	}
	return nil
}
