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
var fmtMap = map[string]string{
	"%a": "Mon",
	"%A": "Monday",
	"%b": "Jan",
	"%h": "Jan",
	"%B": "January",
	"%c": time.UnixDate,
	"%d": "02",
	"%e": "_2",
	"%H": "15",
	"%I": "03",
	"%m": "1",
	"%M": "04",
	"%p": "PM",
	"%S": "05",
	"%y": "06",
	"%Y": "2006",
	"%z": "-0700",
	"%Z": "MST",
}

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
		translate, exists := fmtMap[match]
		switch {
		case exists:
			// Values defined by fmtMap
			toReplace = d.Format(translate)
		case match == "%C":
			// Century (a year divided by 100 and truncated to an integer)
			// as a decimal number [00,99].
			toReplace = strconv.Itoa(d.Year() / 100)
		case match == "%D":
			// Date in the format mm/dd/yy.
			toReplace = dateMap("%m/%d/%y")
		case match == "%j":
			// Day of the year as a decimal number [001,366]."
			year, weekYear := d.ISOWeek()
			firstWeekDay := time.Date(year, 1, 1, 1, 1, 1, 1, time.UTC).Weekday()
			weekDay := d.Weekday()
			dayYear := int(weekYear)*7 - (int(firstWeekDay) - 1) + int(weekDay)
			toReplace = strconv.Itoa(dayYear)
		case match == "%n":
			// A <newline>.
			toReplace = "\n"
		case match == "%r":
			// 12-hour clock time [01,12] using the AM/PM notation;
			// in the POSIX locale, this shall be equivalent to %I : %M : %S %p.
			toReplace = dateMap("%I:%M:%S %p")
		case match == "%t":
			// A <tab>.
			toReplace = "\t"
		case match == "%T":
			toReplace = dateMap("%H:%M:%S")
		case match == "%W":
			// Week of the year (Sunday as the first day of the week)
			// as a decimal number [00,53]. All days in a new year preceding
			// the first Sunday shall be considered to be in week 0.
			_, weekYear := d.ISOWeek()
			weekDay := int(d.Weekday())
			isNotSunday := 1
			if weekDay == 0 {
				isNotSunday = 0
			}
			toReplace = strconv.Itoa(weekYear - (isNotSunday))
		case match == "%w":
			// Weekday as a decimal number [0,6] (0=Sunday).
			toReplace = strconv.Itoa(int(d.Weekday()))
		case match == "%V":
			// Week of the year (Monday as the first day of the week)
			// as a decimal number [01,53]. If the week containing January 1
			// has four or more days in the new year, then it shall be
			// considered week 1; otherwise, it shall be the last week
			// of the previous year, and the next week shall be week 1.
			_, weekYear := d.ISOWeek()
			toReplace = strconv.Itoa(int(weekYear))
		case match == "%x":
			// Locale's appropriate date representation.
			toReplace = dateMap("%m/%d/%y") // TODO: decision algorithm
		case match == "%F":
			// Date yyyy-mm-dd defined by GNU implementation
			toReplace = dateMap("%Y-%m-%d")
		case match == "%X":
			// Locale's appropriate time representation.
			toReplace = dateMap("%I:%M:%S %p") // TODO: decision algorithm
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
