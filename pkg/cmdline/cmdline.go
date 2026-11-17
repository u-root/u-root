// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cmdline is parser for kernel command-line args from /proc/cmdline.
//
// It's conformant with
// https://www.kernel.org/doc/html/v4.14/admin-guide/kernel-parameters.html,
// though making 'var_name' and 'var-name' equivalent may need to be done
// separately.
package cmdline

import (
	"io"
	"strings"

	"github.com/u-root/u-root/pkg/shlex"
)

// CmdLine lets people view the raw & parsed /proc/cmdline in one place
type CmdLine struct {
	Raw   string
	AsMap map[string]string
	Err   error
}

// NewCmdLine returns a populated CmdLine struct
func NewCmdLine() *CmdLine {
	return getCmdLine()
}

// FullCmdLine returns the full, raw cmdline string
func FullCmdLine() string {
	return getCmdLine().Raw
}

// parse returns the current command line, trimmed
func parse(cmdlineReader io.Reader) *CmdLine {
	line := &CmdLine{}
	raw, err := io.ReadAll(cmdlineReader)
	line.Err = err
	// This works because string(nil) is ""
	line.Raw = strings.TrimRight(string(raw), "\n")
	line.AsMap = parseToMap(line.Raw)
	return line
}

func dequote(line string) string {
	if len(line) == 0 {
		return line
	}

	quotationMarks := `"'`

	var quote byte
	if len(line) > 1 && strings.ContainsAny(string(line[0]), quotationMarks) {
		quote = line[0]
		line = line[1 : len(line)-1]
	}

	var context []byte
	var newLine []byte
	for _, c := range []byte(line) {
		if c == '\\' {
			context = append(context, c)
		} else if c == quote {
			if len(context) > 0 {
				last := context[len(context)-1]
				if last == c {
					context = context[:len(context)-1]
				} else if last == '\\' {
					// Delete one level of backslash
					newLine = newLine[:len(newLine)-1]
					context = []byte{}
				}
			} else {
				context = append(context, c)
			}
		} else if len(context) > 0 && context[len(context)-1] == '\\' {
			// If backslash is being used to escape something other
			// than "the quote", ignore it.
			context = []byte{}
		}

		newLine = append(newLine, c)
	}
	return string(newLine)
}

// isKernelSpace reports whether the kernel's isspace() (lib/ctype.c) is true
// for c.
func isKernelSpace(c byte) bool {
	switch c {
	case ' ', '\t', '\n', '\v', '\f', '\r', 0xa0:
		return true
	}
	return false
}

// splitParams splits input into parameters the way the kernel's next_arg
// (lib/cmdline.c) does: only '"' quotes, and only outside quotes does a
// kernel space end a parameter.
func splitParams(input string) []string {
	var params []string
	inQuote := false
	start := -1
	for i := 0; i < len(input); i++ {
		c := input[i]
		if c == '"' {
			inQuote = !inQuote
		} else if !inQuote && isKernelSpace(c) {
			if start >= 0 {
				params = append(params, input[start:i])
				start = -1
			}
			continue
		}
		if start < 0 {
			start = i
		}
	}
	if start >= 0 {
		params = append(params, input[start:])
	}
	return params
}

func doParse(input string, handler func(flag, key, canonicalKey, value, trimmedValue string)) {
	for _, flag := range splitParams(input) {
		// kernel variables must allow '-' and '_' to be equivalent in variable
		// names. We will replace dashes with underscores for processing.

		// Split the flag into a key and value, setting value="1" if none
		before, after, ok := strings.Cut(flag, "=")

		if len(flag) == 0 {
			continue
		}
		var key, value string
		if !ok {
			key = flag
			value = "1"
		} else {
			key = before
			value = after
		}
		canonicalKey := strings.Replace(key, "-", "_", -1)
		trimmedValue := dequote(value)

		// Call the user handler
		handler(flag, key, canonicalKey, value, trimmedValue)
	}
}

// parseToMap turns a space-separated kernel commandline into a map
func parseToMap(input string) map[string]string {
	flagMap := make(map[string]string)
	doParse(input, func(flag, key, canonicalKey, value, trimmedValue string) {
		// We store the value twice, once with dash, once with underscores
		// Just in case people check with the wrong method
		flagMap[canonicalKey] = trimmedValue
		flagMap[key] = trimmedValue
	})

	return flagMap
}

// ContainsFlag verifies that the kernel cmdline has a flag set
func (c *CmdLine) ContainsFlag(flag string) bool {
	_, present := c.Flag(flag)
	return present
}

// ContainsFlag verifies that the kernel cmdline has a flag set
func ContainsFlag(flag string) bool {
	return getCmdLine().ContainsFlag(flag)
}

// Flag returns the value of a flag, and whether it was set
func (c *CmdLine) Flag(flag string) (string, bool) {
	canonicalFlag := strings.Replace(flag, "-", "_", -1)
	value, present := c.AsMap[canonicalFlag]
	return value, present
}

// Flag returns the value of a flag, and whether it was set
func Flag(flag string) (string, bool) {
	return getCmdLine().Flag(flag)
}

// getFlagMap gets specified flags as a map
func getFlagMap(flagName string) map[string]string {
	return parseToMap(flagName)
}

// GetInitFlagMap gets the uroot init flags as a map
func (c *CmdLine) GetInitFlagMap() map[string]string {
	initflags, _ := c.Flag("uroot.initflags")
	return getFlagMap(initflags)
}

// GetInitFlagMap gets the uroot init flags as a map
func GetInitFlagMap() map[string]string {
	return getCmdLine().GetInitFlagMap()
}

// GetUinitArgs gets the uinit argvs.
func (c *CmdLine) GetUinitArgs() []string {
	uinitargs, _ := getCmdLine().Flag("uroot.uinitargs")
	return shlex.Argv(uinitargs)
}

// GetUinitArgs gets the uinit argvs.
func GetUinitArgs() []string {
	return getCmdLine().GetUinitArgs()
}

// FlagsForModule gets all flags for a designated module
// and returns them as a space-seperated string designed to be passed to insmod
// Note that similarly to flags, module names with - and _ are treated the same.
func (c *CmdLine) FlagsForModule(name string) string {
	var ret strings.Builder
	flagsAdded := make(map[string]bool) // Ensures duplicate flags aren't both added
	// Module flags come as moduleName.flag in /proc/cmdline
	prefix := strings.Replace(name, "-", "_", -1) + "."
	for flag, val := range c.AsMap {
		canonicalFlag := strings.Replace(flag, "-", "_", -1)
		if !flagsAdded[canonicalFlag] && strings.HasPrefix(canonicalFlag, prefix) {
			flagsAdded[canonicalFlag] = true
			// They are passed to insmod space separated as flag=val
			ret.WriteString(strings.TrimPrefix(canonicalFlag, prefix) + "=" + val + " ")
		}
	}
	return ret.String()
}

// FlagsForModule gets all flags for a designated module
// and returns them as a space-seperated string designed to be passed to insmod
// Note that similarly to flags, module names with - and _ are treated the same.
func FlagsForModule(name string) string {
	return getCmdLine().FlagsForModule(name)
}

// Consoles returns the list of all `console=` values in the kernel command line.
func (c *CmdLine) Consoles() []string {
	consoles := make([]string, 0)
	for part := range strings.FieldsSeq(c.Raw) {
		if after, ok := strings.CutPrefix(part, "console="); ok {
			consoles = append(consoles, strings.Split(after, ",")[0])
		}
	}
	return consoles
}

// Consoles returns the list of all `console=` values in the kernel command line.
func Consoles() []string {
	return getCmdLine().Consoles()
}
