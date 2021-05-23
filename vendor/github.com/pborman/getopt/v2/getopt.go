// Copyright 2017 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package getopt (v2) provides traditional getopt processing for implementing
// commands that use traditional command lines.  The standard Go flag package
// cannot be used to write a program that parses flags the way ls or ssh does,
// for example.  Version 2 of this package has a simplified API.
//
// See the github.com/pborman/options package for a simple structure based
// interface to this package.
//
// USAGE
//
// Getopt supports functionality found in both the standard BSD getopt as well
// as (one of the many versions of) the GNU getopt_long.  Being a Go package,
// this package makes common usage easy, but still enables more controlled usage
// if needed.
//
// Typical usage:
//
//	// Declare flags and have getopt return pointers to the values.
//	helpFlag := getopt.Bool('?', "display help")
//	cmdFlag := getopt.StringLong("command", 'c', "default", "the command")
//
//	// Declare flags against existing variables.
//	var {
//		fileName = "/the/default/path"
//		timeout = time.Second * 5
//		verbose bool
//	}
//	func init() {
//		getopt.Flag(&verbose, 'v', "be verbose")
//		getopt.FlagLong(&fileName, "path", 0, "the path")
//		getopt.FlagLong(&timeout, "timeout", 't', "some timeout")
//	}
//
//	func main() {
//		// Parse the program arguments
//		getopt.Parse()
//		// Get the remaining positional parameters
//		args := getopt.Args()
//		...
//
// If you don't want the program to exit on error, use getopt.Getopt:
//
//		err := getopt.Getopt(nil)
//		if err != nil {
//			// code to handle error
//			fmt.Fprintln(os.Stderr, err)
//		}
//
// FLAG SYNTAX
//
// Support is provided for both short (-f) and long (--flag) options.  A single
// option may have both a short and a long name.  Each option may be a flag or a
// value.  A value takes an argument.
//
// Declaring no long names causes this package to process arguments like the
// traditional BSD getopt.
//
// Short flags may be combined into a single parameter.  For example, "-a -b -c"
// may also be expressed "-abc".  Long flags must stand on their own "--alpha
// --beta"
//
// Values require an argument.  For short options the argument may either be
// immediately following the short name or as the next argument.  Only one short
// value may be combined with short flags in a single argument; the short value
// must be after all short flags.  For example, if f is a flag and v is a value,
// then:
//
//  -vvalue    (sets v to "value")
//  -v value   (sets v to "value")
//  -fvvalue   (sets f, and sets v to "value")
//  -fv value  (sets f, and sets v to "value")
//  -vf value  (set v to "f" and value is the first parameter)
//
// For the long value option val:
//
//  --val value (sets val to "value")
//  --val=value (sets val to "value")
//  --valvalue  (invalid option "valvalue")
//
// Values with an optional value only set the value if the value is part of the
// same argument.  In any event, the option count is increased and the option is
// marked as seen.
//
//  -v -f          (sets v and f as being seen)
//  -vvalue -f     (sets v to "value" and sets f)
//  --val -f       (sets v and f as being seen)
//  --val=value -f (sets v to "value" and sets f)
//
// There is no convenience function defined for making the value optional.  The
// SetOptional method must be called on the actual Option.
//
//  v := String("val", 'v', "", "the optional v")
//  Lookup("v").SetOptional()
//
//  var s string
//  FlagLong(&s, "val", 'v', "the optional v).SetOptional()
//
// Parsing continues until the first non-option or "--" is encountered.
//
// The short name "-" can be used, but it either is specified as "-" or as part
// of a group of options, for example "-f-".  If there are no long options
// specified then "--f" could also be used.  If "-" is not declared as an option
// then the single "-" will also terminate the option processing but unlike
// "--", the "-" will be part of the remaining arguments.
//
// ADVANCED USAGE
//
// Normally the parsing is performed by calling the Parse function.  If it is
// important to see the order of the options then the Getopt function should be
// used.  The standard Parse function does the equivalent of:
//
// func Parse() {
//	if err := getopt.Getopt(os.Args, nil); err != nil {
//		fmt.Fprintln(os.Stderr, err)
//		s.usage()
//		os.Exit(1)
//	}
//
// When calling Getopt it is the responsibility of the caller to print any
// errors.
//
// Normally the default option set, CommandLine, is used.  Other option sets may
// be created with New.
//
// After parsing, the sets Args will contain the non-option arguments.  If an
// error is encountered then Args will begin with argument that caused the
// error.
//
// It is valid to call a set's Parse a second time to amen flags or values.  As
// an example:
//
//  var a = getopt.Bool('a', "", "The a flag")
//  var b = getopt.Bool('b', "", "The a flag")
//  var cmd = ""
//
//  var opts = getopt.CommandLine
//
//  opts.Parse(os.Args)
//  if opts.NArgs() > 0 {
//      cmd = opts.Arg(0)
//      opts.Parse(opts.Args())
//  }
//
// If called with set to { "prog", "-a", "cmd", "-b", "arg" } then both and and
// b would be set, cmd would be set to "cmd", and opts.Args() would return {
// "arg" }.
//
// Unless an option type explicitly prohibits it, an option may appear more than
// once in the arguments.  The last value provided to the option is the value.
//
// MANDATORY OPTIONS
//
// An option marked as mandatory and not seen when parsing will cause an error
// to be reported such as: "program: --name is a mandatory option".  An option
// is marked mandatory by using the Mandatory method:
//
//	getopt.FlagLong(&fileName, "path", 0, "the path").Mandatory()
//
// Mandatory options have (required) appended to their help message:
//
//	--path=value    the path (required)
//
// MUTUALLY EXCLUSIVE OPTIONS
//
// Options can be marked as part of a mutually exclusive group.  When two or
// more options in a mutually exclusive group are both seen while parsing then
// an error such as "program: options -a and -b are mutually exclusive" will be
// reported.  Mutually exclusive groups are declared using the SetGroup method:
//
//	getopt.Flag(&a, 'a', "use method A").SetGroup("method")
//	getopt.Flag(&a, 'b', "use method B").SetGroup("method")
//
// A set can have multiple mutually exclusive groups.  Mutually exclusive groups
// are identified with their group name in {}'s appeneded to their help message:
//
//	 -a    use method A {method}
//	 -b    use method B {method}
//
// BUILTIN TYPES
//
// The Flag and FlagLong functions support most standard Go types.  For the
// list, see the description of FlagLong below for a list of supported types.
//
// There are also helper routines to allow single line flag declarations.  These
// types are: Bool, Counter, Duration, Enum, Int16, Int32, Int64, Int, List,
// Signed, String, Uint16, Uint32, Uint64, Uint, and Unsigned.
//
// Each comes in a short and long flavor, e.g., Bool and BoolLong and include
// functions to set the flags on the standard command line or for a specific Set
// of flags.
//
// Except for the Counter, Enum, Signed and Unsigned types, all of these types
// can be declared using Flag and FlagLong by passing in a pointer to the
// appropriate type.
//
// DECLARING NEW FLAG TYPES
//
// A pointer to any type that implements the Value interface may be passed to
// Flag or FlagLong.
//
// VALUEHELP
//
// All non-flag options are created with a "valuehelp" as the last parameter.
// Valuehelp should be 0, 1, or 2 strings.  The first string, if provided, is
// the usage message for the option.  If the second string, if provided, is the
// name to use for the value when displaying the usage.  If not provided the
// term "value" is assumed.
//
// The usage message for the option created with
//
//  StringLong("option", 'o', "defval", "a string of letters")
//
// is
//
//  -o, -option=value
//
//  StringLong("option", 'o', "defval", "a string of letters", "string")
//
// is
//
//  -o, -option=string
package getopt

import (
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

// stderr allows tests to capture output to standard error.
var stderr io.Writer = os.Stderr

// exit allows tests to capture an os.Exit call
var exit = os.Exit

// DisplayWidth is used to determine where to split usage long lines.
var DisplayWidth = 80

// HelpColumn is the maximum column position that help strings start to display
// at.  If the option usage is too long then the help string will be displayed
// on the next line.  For example:
//
//   -a   this is the a flag
//   -u, --under=location
//        the u flag's usage is quite long
var HelpColumn = 20

// PrintUsage prints the usage line and set of options of set S to w.
func (s *Set) PrintUsage(w io.Writer) {
	parts := make([]string, 2, 4)
	parts[0] = "Usage:"
	parts[1] = s.program
	if usage := s.UsageLine(); usage != "" {
		parts = append(parts, usage)
	}
	if s.parameters != "" {
		parts = append(parts, s.parameters)
	}
	fmt.Fprintln(w, strings.Join(parts, " "))
	s.PrintOptions(w)
}

// UsageLine returns the usage line for the set s.  The set's program name and
// parameters, if any, are not included.
func (s *Set) UsageLine() string {
	sort.Sort(s.options)
	flags := ""

	// Build up the list of short flag names and also compute
	// how to display the option in the longer help listing.
	// We also keep track of the longest option usage string
	// that is no more than HelpColumn-3 bytes (at which point
	// we use two lines to display the help).  The three
	// is for the leading space and the two spaces before the
	// help string.
	for _, opt := range s.options {
		if opt.name == "" {
			opt.name = "value"
		}
		if opt.uname == "" {
			opt.uname = opt.usageName()
		}
		if opt.flag && opt.short != 0 && opt.short != '-' {
			flags += string(opt.short)
		}
	}

	var opts []string

	// The short option - is special
	if s.shortOptions['-'] != nil {
		opts = append(opts, "-")
	}

	// If we have a bundle of flags, add them to the list
	if flags != "" {
		opts = append(opts, "-"+flags)
	}

	// Now append all the long options and options that require
	// values.
	for _, opt := range s.options {
		if opt.flag {
			if opt.short != 0 {
				continue
			}
			flags = "--" + opt.long
		} else if opt.short != 0 {
			flags = "-" + string(opt.short) + " " + opt.name
		} else {
			flags = "--" + string(opt.long) + " " + opt.name
		}
		opts = append(opts, flags)
	}
	flags = strings.Join(opts, "] [")
	if flags != "" {
		flags = "[" + flags + "]"
	}
	return flags
}

// PrintOptions prints the list of options in s to w.
func (s *Set) PrintOptions(w io.Writer) {
	sort.Sort(s.options)
	max := 4
	for _, opt := range s.options {
		if opt.name == "" {
			opt.name = "value"
		}
		if opt.uname == "" {
			opt.uname = opt.usageName()
		}
		if max < len(opt.uname) && len(opt.uname) <= HelpColumn-3 {
			max = len(opt.uname)
		}
	}
	// Now print one or more usage lines per option.
	for _, opt := range s.options {
		if opt.uname != "" {
			opt.help = strings.TrimSpace(opt.help)
			if len(opt.help) == 0 && !opt.mandatory && opt.group == "" {
				fmt.Fprintf(w, " %s\n", opt.uname)
				continue
			}
			helpMsg := opt.help

			// If the default value is the known zero value
			// then don't display it.
			def := opt.defval
			switch genericValue(opt.value).(type) {
			case *bool:
				if def == "false" {
					def = ""
				}
			case *int, *int8, *int16, *int32, *int64,
				*uint, *uint8, *uint16, *uint32, *uint64,
				*float32, *float64:
				if def == "0" {
					def = ""
				}
			case *time.Duration:
				if def == "0s" {
					def = ""
				}
			default:
				if opt.flag && def == "false" {
					def = ""
				}
			}
			if def != "" {
				helpMsg += " [" + def + "]"
			}
			if opt.group != "" {
				helpMsg += " {" + opt.group + "}"
			}
			if opt.mandatory {
				helpMsg += " (required)"
			}

			help := strings.Split(helpMsg, "\n")
			// If they did not put in newlines then we will insert
			// them to keep the help messages from wrapping.
			if len(help) == 1 {
				help = breakup(help[0], DisplayWidth-HelpColumn)
			}
			if len(opt.uname) <= max {
				fmt.Fprintf(w, " %-*s  %s\n", max, opt.uname, help[0])
				help = help[1:]
			} else {
				fmt.Fprintf(w, " %s\n", opt.uname)
			}
			for _, s := range help {
				fmt.Fprintf(w, " %-*s  %s\n", max, " ", s)
			}
		}
	}
}

// breakup breaks s up into strings no longer than max bytes.
func breakup(s string, max int) []string {
	var a []string

	for {
		// strip leading spaces
		for len(s) > 0 && s[0] == ' ' {
			s = s[1:]
		}
		// If the option is no longer than the max just return it
		if len(s) <= max {
			if len(s) != 0 {
				a = append(a, s)
			}
			return a
		}
		x := max
		for s[x] != ' ' {
			// the first word is too long?!
			if x == 0 {
				x = max
				for x < len(s) && s[x] != ' ' {
					x++
				}
				if x == len(s) {
					x--
				}
				break
			}
			x--
		}
		for s[x] == ' ' {
			x--
		}
		a = append(a, s[:x+1])
		s = s[x+1:]
	}
}

// Parse uses Getopt to parse args using the options set for s.  The first
// element of args is used to assign the program for s if it is not yet set.  On
// error, Parse displays the error message as well as a usage message on
// standard error and then exits the program.
func (s *Set) Parse(args []string) {
	if err := s.Getopt(args, nil); err != nil {
		fmt.Fprintln(stderr, err)
		s.usage()
		exit(1)
	}
}

// Parse uses Getopt to parse args using the options set for s.  The first
// element of args is used to assign the program for s if it is not yet set.
// Getop calls fn, if not nil, for each option parsed.
//
// Getopt returns nil when all options have been processed (a non-option
// argument was encountered, "--" was encountered, or fn returned false).
//
// On error getopt returns a reference to an InvalidOption (which implements the
// error interface).
func (s *Set) Getopt(args []string, fn func(Option) bool) (err error) {
	s.setState(InProgress)
	defer func() {
		if s.State() == InProgress {
			switch {
			case err != nil:
				s.setState(Failure)
			case len(s.args) == 0:
				s.setState(EndOfArguments)
			default:
				s.setState(Unknown)
			}
		}
	}()

	defer func() {
		if err == nil {
			err = s.checkOptions()
		}
	}()
	if fn == nil {
		fn = func(Option) bool { return true }
	}
	if len(args) == 0 {
		return nil
	}

	if s.program == "" {
		s.program = path.Base(args[0])
	}
	args = args[1:]
Parsing:
	for len(args) > 0 {
		arg := args[0]
		s.args = args
		args = args[1:]

		// end of options?
		if arg == "" || arg[0] != '-' {
			s.setState(EndOfOptions)
			return nil
		}

		if arg == "-" {
			goto ShortParsing
		}

		// explicitly request end of options?
		if arg == "--" {
			s.args = args
			s.setState(DashDash)
			return nil
		}

		// Long option processing
		if len(s.longOptions) > 0 && arg[1] == '-' {
			e := strings.IndexRune(arg, '=')
			var value string
			if e > 0 {
				value = arg[e+1:]
				arg = arg[:e]
			}
			opt := s.longOptions[arg[2:]]
			// If we are processing long options then --f is -f
			// if f is not defined as a long option.
			// This lets you say --f=false
			if opt == nil && len(arg[2:]) == 1 {
				opt = s.shortOptions[rune(arg[2])]
			}
			if opt == nil {
				return unknownOption(arg[2:])
			}
			opt.isLong = true
			// If we require an option and did not have an =
			// then use the next argument as an option.
			if !opt.flag && e < 0 && !opt.optional {
				if len(args) == 0 {
					return missingArg(opt)
				}
				value = args[0]
				args = args[1:]
			}
			opt.count++

			if err := opt.value.Set(value, opt); err != nil {
				return setError(opt, value, err)
			}

			if !fn(opt) {
				s.setState(Terminated)
				return nil
			}
			continue Parsing
		}

		// Short option processing
		arg = arg[1:] // strip -
	ShortParsing:
		for i, c := range arg {
			opt := s.shortOptions[c]
			if opt == nil {
				// In traditional getopt, if - is not registered
				// as an option, a lone - is treated as
				// if there were a -- in front of it.
				if arg == "-" {
					s.setState(Dash)
					return nil
				}
				return unknownOption(c)
			}
			opt.isLong = false
			opt.count++
			var value string
			if !opt.flag {
				value = arg[1+i:]
				if value == "" && !opt.optional {
					if len(args) == 0 {
						return missingArg(opt)
					}
					value = args[0]
					args = args[1:]
				}
			}
			if err := opt.value.Set(value, opt); err != nil {
				return setError(opt, value, err)
			}
			if !fn(opt) {
				s.setState(Terminated)
				return nil
			}
			if !opt.flag {
				continue Parsing
			}
		}
	}
	s.args = []string{}
	return nil
}

func (s *Set) checkOptions() error {
	groups := map[string]Option{}
	for _, opt := range s.options {
		if !opt.Seen() {
			if opt.mandatory {
				return fmt.Errorf("option %s is mandatory", opt.Name())
			}
			continue
		}
		if opt.group == "" {
			continue
		}
		if opt2 := groups[opt.group]; opt2 != nil {
			return fmt.Errorf("options %s and %s are mutually exclusive", opt2.Name(), opt.Name())
		}
		groups[opt.group] = opt
	}
	for _, group := range s.requiredGroups {
		if groups[group] != nil {
			continue
		}
		var flags []string
		for _, opt := range s.options {
			if opt.group == group {
				flags = append(flags, opt.Name())
			}
		}
		return fmt.Errorf("exactly one of the following options must be specified: %s", strings.Join(flags, ", "))
	}
	return nil
}
