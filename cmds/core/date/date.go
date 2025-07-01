// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// date prints the date.
//
// Synopsis:
//
//	date [-u] [+format] | date [-u] [MMDDhhmm[CC]YY[.ss]]
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (r RealClock) Now() time.Time {
	return time.Now()
}

var (
	// default format map from format.go on time lib
	// Help to make format of date with posix compliant
	fmtMap = map[string]string{
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
	flags struct {
		universal bool
		reference string
	}
)

const cmd = "date [-u] [-d FILE] [+format] | date [-u] [-d FILE] [MMDDhhmm[CC]YY[.ss]]"

func init() {
	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = cmd
		defUsage()
	}
	flag.BoolVar(&flags.universal, "u", false, "Coordinated Universal Time (UTC)")
	flag.StringVar(&flags.reference, "r", "", "Display the last modification time of FILE")
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
func dateMap(t time.Time, z *time.Location, format string) string {
	d := t.In(z)
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
			toReplace = dateMap(t, z, "%m/%d/%y")
		case match == "%j":
			// Day of the year as a decimal number [001,366]."
			toReplace = strconv.Itoa(d.YearDay())
		case match == "%n":
			// A <newline>.
			toReplace = "\n"
		case match == "%r":
			// 12-hour clock time [01,12] using the AM/PM notation;
			// in the POSIX locale, this shall be equivalent to %I : %M : %S %p.
			toReplace = dateMap(t, z, "%I:%M:%S %p")
		case match == "%t":
			// A <tab>.
			toReplace = "\t"
		case match == "%T":
			toReplace = dateMap(t, z, "%H:%M:%S")
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
			toReplace = dateMap(t, z, "%m/%d/%y") // TODO: decision algorithm
		case match == "%F":
			// Date yyyy-mm-dd defined by GNU implementation
			toReplace = dateMap(t, z, "%Y-%m-%d")
		case match == "%X":
			// Locale's appropriate time representation.
			toReplace = dateMap(t, z, "%I:%M:%S %p") // TODO: decision algorithm
		default:
			continue
		}

		format = strings.Replace(format, match, toReplace, 1)
		// fmt.Printf("Iteration[%d]: %v\n", iter, format)
	}
	return format
}

func ints(s string, i ...*int) error {
	var err error
	for _, p := range i {
		if *p, err = strconv.Atoi(s[0:2]); err != nil {
			return err
		}
		s = s[2:]
	}
	return nil
}

// getTime gets the desired time as a time.Time.
// It derives it from a unix date command string.
// Some values in the string are optional, namely
// YY and SS. For these values, we use
// time.Now(). For the timezone, we use whatever
// one we are in, or UTC if desired.
func getTime(z *time.Location, s string, clocksource Clock) (t time.Time, err error) {
	var MM, DD, hh, mm int
	// CC is the year / 100, not the "century".
	// i.e. for 2001, CC is 20, not 21.
	YY := clocksource.Now().Year() % 100
	CC := clocksource.Now().Year() / 100
	SS := clocksource.Now().Second()
	if err = ints(s, &MM, &DD, &hh, &mm); err != nil {
		return
	}
	s = s[8:]
	switch len(s) {
	case 0:
	case 2:
		err = ints(s, &YY)
	case 3:
		err = ints(s[1:], &SS)
	case 4:
		err = ints(s, &CC, &YY)
	case 5:
		s = s[0:2] + s[3:]
		err = ints(s, &YY, &SS)
	case 7:
		s = s[0:4] + s[5:]
		err = ints(s, &CC, &YY, &SS)
	default:
		err = fmt.Errorf("optional string is %v instead of [[CC]YY][.ss]", s)
	}

	if err != nil {
		return
	}

	YY = YY + CC*100
	t = time.Date(YY, time.Month(MM), DD, hh, mm, SS, 0, z)
	return
}

func date(t time.Time, z *time.Location) string {
	return t.In(z).Format(time.UnixDate)
}

func run(args []string, univ bool, ref string, clocksource Clock, w io.Writer) error {
	t := clocksource.Now()
	z := time.Local
	if univ {
		z = time.UTC
	}
	if ref != "" {
		stat, err := os.Stat(ref)
		if err != nil {
			return fmt.Errorf("unable to gather stats of file %v", ref)
		}
		t = stat.ModTime()
	}

	switch len(args) {
	case 0:
		fmt.Fprintf(w, "%v\n", date(t, z))
	case 1:
		a0 := args[0]
		if strings.HasPrefix(a0, "+") {
			fmt.Fprintf(w, "%v\n", dateMap(t, z, a0[1:]))
		} else {
			if err := setDate(args[0], z, clocksource); err != nil {
				return fmt.Errorf("%v: %w", a0, err)
			}
		}
	default:
		flag.Usage()
		return nil
	}
	return nil
}

func main() {
	flag.Parse()
	rc := RealClock{}
	if err := run(flag.Args(), flags.universal, flags.reference, rc, os.Stdout); err != nil {
		log.Fatalf("date: %v", err)
	}
}
