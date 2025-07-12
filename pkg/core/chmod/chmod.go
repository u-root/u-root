// Copyright 2016-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package chmod implements the chmod core utility.
package chmod

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/u-root/u-root/pkg/core"
)

const special = 99999

var errBadUsage = errors.New("chmod: chmod [mode] filepath")

// command implements the chmod command.
type command struct {
	core.Base
}

// New creates a new chmod command.
func New() core.Command {
	c := &command{}
	c.Init()
	return c
}

type flags struct {
	recursive bool
	reference string
}

func (c *command) changeMode(path string, mode os.FileMode, octval uint64, mask uint64, operator string) error {
	path = c.ResolvePath(path)

	// A special value for mask means the mode is fully described
	if mask == special {
		if err := os.Chmod(path, mode); err != nil {
			return err
		}
		return nil
	}

	var info os.FileInfo
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	currentMode := info.Mode()

	switch operator {
	case "+":
		// Add permissions
		mode = currentMode | os.FileMode(octval)
	case "-":
		// Remove permissions
		mode = currentMode &^ os.FileMode(octval)
	case "=":
		// Set permissions exactly (within the specified mask)
		mode = (currentMode & os.FileMode(mask)) | os.FileMode(octval)
	}

	if err := os.Chmod(path, mode); err != nil {
		return err
	}
	return nil
}

func (c *command) calculateMode(modeString string) (mode os.FileMode, octval uint64, mask uint64, operator string, err error) {
	octval, err = strconv.ParseUint(modeString, 8, 32)
	if err == nil {
		if octval > 0o777 {
			return mode, octval, mask, operator, fmt.Errorf("%w: invalid octal value %0o. Value should be less than or equal to 0777", strconv.ErrRange, octval)
		}
		// a fully described octal mode was supplied, signal that with a special value for mask
		mask = special
		mode = os.FileMode(octval)
		operator = "="
		return
	}

	// Try with user/group specified first
	reMode := regexp.MustCompile("^([ugoa]+)([-+=])(.*)")
	m := reMode.FindStringSubmatch(modeString)

	// If no match, try without user/group (defaults to 'a' - all)
	if len(m) == 0 {
		reMode = regexp.MustCompile("^([-+=])(.*)")
		m = reMode.FindStringSubmatch(modeString)
		if len(m) > 0 {
			// Insert 'a' as the default user/group
			m = []string{m[0], "a", m[1], m[2]}
		}
	}

	// Test for mode strings with invalid characters.
	// This can't be done in the first regexp: if the match for m[3] is restricted to [rwx]*,
	// `a=9` and `a=` would be indistinguishable: m[3] would be empty.
	// `a=` is a valid (but destructive) operation. Do not turn a typo into that.
	reMode = regexp.MustCompile("^[rwx]*$")
	if len(m) < 4 || !reMode.MatchString(m[3]) {
		return mode, octval, mask, operator, fmt.Errorf("%w:unable to decode mode %q. Please use an octal value or a valid mode string", strconv.ErrSyntax, modeString)
	}

	// m[3] is [rwx]{0,3}
	var octvalDigit uint64
	if strings.Contains(m[3], "r") {
		octvalDigit += 4
	}
	if strings.Contains(m[3], "w") {
		octvalDigit += 2
	}
	if strings.Contains(m[3], "x") {
		octvalDigit++
	}

	// m[2] is [-+=]
	operator = m[2]

	// m[1] is [ugoa]+
	if strings.Contains(m[1], "o") || strings.Contains(m[1], "a") {
		octval += octvalDigit
	}
	if strings.Contains(m[1], "g") || strings.Contains(m[1], "a") {
		octval += octvalDigit << 3
	}
	if strings.Contains(m[1], "u") || strings.Contains(m[1], "a") {
		octval += octvalDigit << 6
	}

	// For "=" operations, we need a mask to preserve unspecified bits
	if operator == "=" {
		mask = 0o777
		if strings.Contains(m[1], "o") || strings.Contains(m[1], "a") {
			mask = mask & 0o770
		}
		if strings.Contains(m[1], "g") || strings.Contains(m[1], "a") {
			mask = mask & 0o707
		}
		if strings.Contains(m[1], "u") || strings.Contains(m[1], "a") {
			mask = mask & 0o077
		}

		// The mode is fully described, signal that with a special value for mask
		if strings.Contains(m[1], "a") {
			mask = special
			mode = os.FileMode(octval)
		}
	}

	return mode, octval, mask, operator, nil
}

func (c *command) run(args []string, f flags) error {
	var mode os.FileMode
	if len(args) < 1 {
		return errBadUsage
	}

	if len(args) < 2 && f.reference == "" {
		return errBadUsage
	}

	var (
		octval, mask uint64
		operator     string
		fileList     []string
	)

	if f.reference != "" {
		refPath := c.ResolvePath(f.reference)
		fi, err := os.Stat(refPath)
		if err != nil {
			return fmt.Errorf("bad reference file: %w", err)
		}
		mask = special
		mode = fi.Mode()
		operator = "="
		fileList = args
	} else {
		var err error
		if mode, octval, mask, operator, err = c.calculateMode(args[0]); err != nil {
			return err
		}
		fileList = args[1:]
	}

	var finalErr error

	for _, name := range fileList {
		if f.recursive {
			err := filepath.Walk(c.ResolvePath(name), func(path string, _ os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				err = c.changeMode(path, mode, octval, mask, operator)
				return err
			})
			if err != nil {
				finalErr = err
				fmt.Fprintln(c.Stderr, err)
			}
		} else {
			err := c.changeMode(name, mode, octval, mask, operator)
			if err != nil {
				finalErr = err
				fmt.Fprintln(c.Stderr, err)
			}
		}
	}
	return finalErr
}

// Run executes the command with a `context.Background()`.
func (c *command) Run(args ...string) error {
	return c.RunContext(context.Background(), args...)
}

// Run executes the command.
func (c *command) RunContext(ctx context.Context, args ...string) error {
	var f flags

	fs := flag.NewFlagSet("chmod", flag.ContinueOnError)
	fs.SetOutput(c.Stderr)

	fs.BoolVar(&f.recursive, "recursive", false, "do changes recursively")
	fs.StringVar(&f.reference, "reference", "", "use mode from reference file")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: chmod MODE FILE...\n\n")
		fmt.Fprintf(fs.Output(), "MODE is a three character octal value or a string like a=rwx\n\n")
		fs.PrintDefaults()
	}

	// Parse arguments manually to handle mode strings that start with - or +
	var parsedArgs []string
	var i int
	for i = 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			i++
			break
		}
		if arg == "-recursive" || arg == "--recursive" {
			f.recursive = true
			continue
		}
		if arg == "-reference" || arg == "--reference" {
			if i+1 < len(args) {
				f.reference = args[i+1]
				i++
				continue
			}
			return fmt.Errorf("flag needs an argument: %s", arg)
		}
		if strings.HasPrefix(arg, "-reference=") || strings.HasPrefix(arg, "--reference=") {
			f.reference = strings.SplitN(arg, "=", 2)[1]
			continue
		}
		// If it starts with - but is not a known flag, treat it as a mode string
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") {
			// Check if it looks like a mode string (contains rwx or is just -)
			if strings.ContainsAny(arg[1:], "rwx") || arg == "-" {
				parsedArgs = append(parsedArgs, arg)
				i++
				break
			}
		}
		// All other arguments are positional
		parsedArgs = append(parsedArgs, arg)
		i++
		break
	}

	// Add remaining arguments
	for ; i < len(args); i++ {
		parsedArgs = append(parsedArgs, args[i])
	}

	if err := c.run(parsedArgs, f); err != nil {
		return err
	}

	return nil
}
