// Copyright 2015 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// default format map from format.go on time lib
// Help to make format of date with posix compliant
const (
	stdLongMonth             = "January"
	stdMonth                 = "Jan"
	stdNumMonth              = "1"
	stdZeroMonth             = "01"
	stdLongWeekDay           = "Monday"
	stdWeekDay               = "Mon"
	stdDay                   = "2"
	stdUnderDay              = "_2"
	stdZeroDay               = "02"
	stdHour                  = "15"
	stdHour12                = "3"
	stdZeroHour12            = "03"
	stdMinute                = "4"
	stdZeroMinute            = "04"
	stdSecond                = "5"
	stdZeroSecond            = "05"
	stdLongYear              = "2006"
	stdYear                  = "06"
	stdPM                    = "PM"
	stdpm                    = "pm"
	stdTZ                    = "MST"
	stdISO8601TZ             = "Z0700" // prints Z for UTC
	stdISO8601SecondsTZ      = "Z070000"
	stdISO8601ColonTZ        = "Z07:00" // prints Z for UTC
	stdISO8601ColonSecondsTZ = "Z07:00:00"
	stdNumTZ                 = "-0700" // always numeric
	stdNumSecondsTz          = "-070000"
	stdNumShortTZ            = "-07"    // always numeric
	stdNumColonTZ            = "-07:00" // always numeric
	stdNumColonSecondsTZ     = "-07:00:00"
	stdFracSecond0           = ".0"
	stdFracSecond9           = ".9"
)

var (
	flags struct{ universal bool }
	cmd   = "date [-u] [+format]"
)

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:", cmd)
	flag.PrintDefaults()
	os.Exit(1)
}

func init() {
	flag.BoolVar(&flags.universal, "u", false, "Coordinated Universal Time (UTC)")
	flag.Usage = usage
	flag.Parse()
}

// regex search for +format POSIX patterns
func formatParser(args string) []string {
	pattern := regexp.MustCompile("%[a-zA-Z]")
	match := pattern.FindAll([]byte(args), -1)

	var results []string
	for _, m := range match {
		results = append(results, string(m[:]))
	}

	return results
}

// replace map for the format patterns according POSIX and GNU implementations
func dateMap(format string) string {
	d := time.Now()
	if flags.universal {
		d = d.UTC()
	}
	var toReplace string
	for _, match := range formatParser(format) {
		switch match {
		case "%a":
			// Locale's abbreviated weekday name.
			toReplace = d.Format(stdWeekDay)
		case "%A":
			// Locale's full weekday name.
			toReplace = d.Format(stdLongWeekDay)
		case "%b", "%h":
			// Locale's abbreviated month name. %h is a alias
			toReplace = d.Format(stdMonth)
		case "%B":
			// Locale's full month name.
			toReplace = d.Format(stdLongMonth)
		case "%c":
			// Locale's appropriate date and time representation.
			toReplace = d.Format(time.UnixDate) // change to default: imply -u
		case "%C":
			// Century (a year divided by 100 and truncated to an integer)
			// as a decimal number [00,99].
			toReplace = strconv.Itoa(d.Year() / 100)
		case "%d":
			// Day of the month as a decimal number [01,31].
			toReplace = d.Format(stdZeroDay)
		case "%D":
			// Date in the format mm/dd/yy.
			toReplace = dateMap("%m/%d/%y")
		case "%e":
			// Day of the month as a decimal number [1,31]
			// in a two-digit field with leading <space> character fill.
			toReplace = d.Format(stdUnderDay)
		case "%H":
			// Hour (24-hour clock) as a decimal number [00,23].
			toReplace = d.Format(stdHour)
		case "%I":
			// Hour (12-hour clock) as a decimal number [01,12].
			toReplace = d.Format(stdZeroHour12)
		case "%j":
			// Day of the year as a decimal number [001,366]."
			year, weekYear := d.ISOWeek()
			firstWeekDay := time.Date(year, 1, 1, 1, 1, 1, 1, time.UTC).Weekday()
			weekDay := d.Weekday()
			dayYear := int(weekYear)*7 - (int(firstWeekDay) - 1) + int(weekDay)
			toReplace = strconv.Itoa(dayYear)
		case "%m":
			// Month as a decimal number [01,12].
			toReplace = d.Format(stdNumMonth)
		case "%M":
			// Minute as a decimal number [00,59].
			toReplace = d.Format(stdZeroMinute)
		case "%n":
			// A <newline>.
			toReplace = "\n"
		case "%p":
			// Locale's equivalent of either AM or PM.
			toReplace = d.Format(stdPM)
		case "%r":
			// 12-hour clock time [01,12] using the AM/PM notation;
			// in the POSIX locale, this shall be equivalent to %I : %M : %S %p.
			toReplace = dateMap("%I:%M:%S %p")
		case "%S":
			// Seconds as a decimal number [00,60].
			toReplace = d.Format(stdZeroSecond)
		case "%t":
			// A <tab>.
			toReplace = "\t"
		case "%T":
			toReplace = dateMap("%H:%M:%S")
		case "%W":
			// Week of the year (Sunday as the first day of the week)
			// as a decimal number [00,53]. All days in a new year preceding
			// the first Sunday shall be considered to be in week 0.
			_, weekYear := d.ISOWeek()
			toReplace = strconv.Itoa(int(weekYear))
		case "%w":
			// Weekday as a decimal number [0,6] (0=Sunday).
			toReplace = strconv.Itoa(int(d.Weekday()))
		case "%V":
			// Week of the year (Monday as the first day of the week)
			// as a decimal number [01,53]. If the week containing January 1
			// has four or more days in the new year, then it shall be
			// considered week 1; otherwise, it shall be the last week
			// of the previous year, and the next week shall be week 1.
			toReplace = ":: toImplement ::"
		case "%x":
			// Locale's appropriate date representation.
			toReplace = dateMap("%m/%d/%y") // TODO: decision algorithm
		case "%X":
			// Locale's appropriate time representation.
			toReplace = dateMap("%I:%M:%S %p") // TODO: decision algorithm
		case "%y":
			// Year within century [00,99].
			toReplace = d.Format(stdYear)
		case "%Y":
			// Year with century as a decimal number.
			toReplace = d.Format(stdLongYear)
		case "%z":
			// Defined by GNU implementation: Numeric Timezone
			toReplace = d.Format(stdNumTZ)
		case "%Z":
			// Timezone name, or no characters if no timezone is determinable.
			toReplace = d.Format(stdTZ)
		default:
			continue
		}

		format = strings.Replace(format, match, toReplace, 1)
		// fmt.Printf("Iteration[%d]: %v\n", iter, format)
	}
	return format
}

func date() (string, error) {
	t := time.Now()
	if flags.universal {
		t = t.UTC()
	}
	s := t.Format(time.UnixDate)
	return fmt.Sprintf("%v", s), nil
}

func main() {
	// date without format args
	msg, err := date()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// trick of flag with start '+'
	for _, argv := range flag.Args() {
		if argv[0] == '+' {
			msg = dateMap(argv[1:])
		}
	}

	fmt.Printf("%v\n", msg)

}
