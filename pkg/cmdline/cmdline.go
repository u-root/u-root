// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cmdline is a parser for Linux kernel command-line args.
//
// cmdline can parse command-line args from /proc/cmdline.
//
// It's conformant with
// https://www.kernel.org/doc/html/v4.14/admin-guide/kernel-parameters.html,
// though making 'var_name' and 'var-name' equivalent may need to be done
// separately.
package cmdline

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"unicode"
)

func canon(flag string) string {
	return strings.Replace(flag, "-", "_", -1)
}

// Cmdline represents a kernel command-line arg string.
type Cmdline struct {
	raw   string
	asMap map[string]string
}

// NewCmdline returns an empty Cmdline object.
func NewCmdline() *Cmdline {
	return &Cmdline{
		asMap: make(map[string]string),
	}
}

var (
	hostOnce    sync.Once
	hostCmdline *Cmdline
	hostErr     error
)

// Parse parses s into a Cmdline object according to the Linux kernel
// rules for kernel parameters.
func Parse(s string) *Cmdline {
	return &Cmdline{
		raw:   s,
		asMap: parseToMap(s, true),
	}
}

// HostCmdline parses the host's /proc/cmdline.
func HostCmdline() (*Cmdline, error) {
	hostOnce.Do(func() {
		b, err := ioutil.ReadFile("/proc/cmdline")
		if err != nil {
			hostErr = err
			return
		}
		content := strings.TrimRight(string(b), "\n")
		hostCmdline = Parse(content)
	})
	return hostCmdline, hostErr
}

// Copy returns an identical deep copy of c.
func (c *Cmdline) Copy() *Cmdline {
	d := &Cmdline{
		raw:   c.raw,
		asMap: make(map[string]string),
	}
	for k, v := range c.asMap {
		d.asMap[k] = v
	}
	return d
}

// Reuse reads given flags from the host system (if they exist) and inserts them to c.
func (c *Cmdline) Reuse(flag ...string) {
	for _, f := range flag {
		if value, ok := Flag(f); ok {
			// TODO: Not quite ok because of quoting.
			c.Append(fmt.Sprintf("%s=%q", f, value))
		}
	}
}

// Remove removes the given flags from the kernel parameters.
func (c *Cmdline) Remove(flag ...string) {
	// kernel variables must allow '-' and '_' to be equivalent in variable
	// names.
	for i, v := range flag {
		flag[i] = canon(v)
	}

	for _, f := range flag {
		if _, ok := c.asMap[f]; ok {
			delete(c.asMap, f)
		}
	}
	c.raw = removeFilter(c.raw, flag)
}

func stringsContain(s []string, q string) bool {
	for _, r := range s {
		if r == q {
			return true
		}
	}
	return false
}

// RemoveFilter filters out variable for a given space-separated kernel commandline
func removeFilter(input string, flags []string) string {
	var newCl []string
	doParse(input, func(flag, key, value, trimmedValue string) {
		if !stringsContain(flags, canon(key)) {
			newCl = append(newCl, flag)
		}
	})
	return strings.Join(newCl, " ")
}

// Append appends values to the kernel params, and overrides earlier values.
func (c *Cmdline) Append(s string) {
	if len(c.raw) == 0 {
		c.raw = s
	} else {
		c.raw = c.raw + " " + s
	}
	// Appending overrides earlier values.
	asMap := parseToMap(s, true)
	for k, v := range asMap {
		c.asMap[k] = v
	}
}

// Prepend prepends values to the kernel params, and values do not override
// existing flags.
func (c *Cmdline) Prepend(s string) {
	if len(c.raw) == 0 {
		c.raw = s
	} else {
		c.raw = s + " " + c.raw
	}
	// Later values have priorities, so prepending should check c.asMap
	// before adding values.
	asMap := parseToMap(s, true)
	for k, v := range asMap {
		if _, ok := c.asMap[k]; !ok {
			c.asMap[k] = v
		}
	}
}

// String returns the full kernel parameter string.
func (c *Cmdline) String() string {
	return c.raw
}

// Flag returns the value corresponding to the given key in the kernel
// parameter string c.
func (c *Cmdline) Flag(flag string) (string, bool) {
	if c == nil {
		return "", false
	}
	s, ok := c.asMap[canon(flag)]
	return s, ok
}

// ContainsFlag verifies that the cmdline has a flag set.
func (c *Cmdline) Contains(flag string) bool {
	_, present := c.Flag(flag)
	return present
}

// ContainsFlag verifies that the host kernel cmdline has a flag set.
func ContainsFlag(flag string) bool {
	h, _ := HostCmdline()
	return h.Contains(flag)
}

// Flag returns the host kernel cmdline value for flag.
func Flag(flag string) (string, bool) {
	h, _ := HostCmdline()
	return h.Flag(flag)
}

// GetInitFlagMap gets the init flags as a map
func GetInitFlagMap() map[string]string {
	initflags, _ := Flag("uroot.initflags")
	return parseToMap(initflags, false)
}

// GetUinitArgs gets the uinit argvs.
func GetUinitArgs() []string {
	uinitargs, _ := Flag("uroot.uinitargs")
	return strings.Fields(uinitargs)
}

// FlagsForModule gets all flags for a designated module and returns them as a
// space-seperated string designed to be passed to insmod.
//
// Note that similarly to flags, module names with - and _ are treated the same.
func (c *Cmdline) FlagsForModule(name string) string {
	if c == nil {
		return ""
	}
	var params []string
	flagsAdded := make(map[string]bool) // Ensures duplicate flags aren't both added

	// Module flags come as moduleName.flag in /proc/cmdline
	prefix := canon(name) + "."

	for flag, val := range c.asMap {
		cf := canon(flag)
		if !flagsAdded[cf] && strings.HasPrefix(cf, prefix) {
			flagsAdded[cf] = true
			// They are passed to insmod space seperated as flag=val
			params = append(params, strings.TrimPrefix(cf, prefix)+"="+val)
		}
	}
	return strings.Join(params, " ")
}

func doParse(input string, handler func(flag, key, value, trimmedValue string)) {
	lastQuote := rune(0)
	quotedFieldsCheck := func(c rune) bool {
		switch {
		case c == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastQuote = c
			return false
		default:
			return unicode.IsSpace(c)
		}
	}

	for _, flag := range strings.FieldsFunc(string(input), quotedFieldsCheck) {
		// kernel variables must allow '-' and '_' to be equivalent in variable
		// names. We will replace dashes with underscores for processing.

		// Split the flag into a key and value, setting value="1" if none
		split := strings.Index(flag, "=")

		if len(flag) == 0 {
			continue
		}
		var key, value string
		if split == -1 {
			key = flag
			value = "1"
		} else {
			key = flag[:split]
			value = flag[split+1:]
		}
		trimmedValue := strings.Trim(value, "\"'")

		// Call the user handler
		handler(flag, key, value, trimmedValue)
	}
}

// parseToMap turns a space-separated kernel commandline into a map
func parseToMap(input string, canonical bool) map[string]string {
	flagMap := make(map[string]string)
	doParse(input, func(flag, key, value, trimmedValue string) {
		if canonical {
			flagMap[canon(key)] = trimmedValue
		} else {
			flagMap[key] = trimmedValue
		}
	})
	return flagMap
}
