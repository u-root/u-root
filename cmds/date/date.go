// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Print the date.
//
// Synopsis:
//     date [-u] [+format] | date [-u] [MMDDhhmm[CC]YY[.ss]]
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
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
	cmd   = "date [-u] [+format] | date [-u] [MMDDhhmm[CC]YY[.ss]]"
	z     = time.Local
)

func init() {
	flag.BoolVar(&flags.universal, "u", false, "Coordinated Universal Time (UTC)")
	flag.Usage = func(f func()) func() {
		return func() {
			os.Args[0] = cmd
			f()
		}
	}(flag.Usage)
	flag.Parse()
	if flags.universal {
		z = time.UTC
	}
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
	d := time.Now().In(z)
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
func getTime(s string) (t time.Time, err error) {
	var MM, DD, hh, mm int
	// CC is the year / 100, not the "century".
	// i.e. for 2001, CC is 20, not 21.
	YY := time.Now().Year() % 100
	CC := time.Now().Year() / 100
	SS := time.Now().Second()
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
		err = fmt.Errorf("Optional string is %v instead of [[CC]YY][.ss]", s)
	}

	if err != nil {
		return
	}

	YY = YY + CC*100
	t = time.Date(YY, time.Month(MM), DD, hh, mm, SS, 0, z)
	return
}

func date(z *time.Location) string {
	return time.Now().In(z).Format(time.UnixDate)
}

func main() {
	switch len(flag.Args()) {
	case 0:
		fmt.Printf("%v\n", date(z))
	case 1:
		argv0 := flag.Args()[0]
		if argv0[0] == '+' {
			fmt.Printf("%v\n", dateMap(argv0[1:]))
		} else {
			t, err := getTime(argv0)
			if err != nil {
				log.Fatalf("%v: %v", argv0, err)
			}
			tv := syscall.NsecToTimeval(t.UnixNano())
			if err := syscall.Settimeofday(&tv); err != nil {
				log.Fatalf("%v: %v", argv0, err)
			}
		}
	default:
		flag.Usage()
	}
}
