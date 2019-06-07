// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checker

import "fmt"

// foreground colours
const (
	ColorBlack   = "\x1b[0;30m"
	ColorRed     = "\x1b[0;31m"
	ColorGreen   = "\x1b[0;32m"
	ColorYellow  = "\x1b[0;33m"
	ColorBlue    = "\x1b[0;34m"
	ColorMagenta = "\x1b[0;35m"
	ColorGrey    = "\x1b[0;36m"
	ColorNone    = "\x1b[0m"
)

func colorize(col, f string, a ...interface{}) string {
	return col + fmt.Sprintf(f, a...) + ColorNone
}

func red(format string, args ...interface{}) string {
	return colorize(ColorRed, format, args...)
}

func green(format string, args ...interface{}) string {
	return colorize(ColorGreen, format, args...)
}

func yellow(format string, args ...interface{}) string {
	return colorize(ColorYellow, format, args...)
}

func blue(format string, args ...interface{}) string {
	return colorize(ColorBlue, format, args...)
}

func magenta(format string, args ...interface{}) string {
	return colorize(ColorMagenta, format, args...)
}

func grey(format string, args ...interface{}) string {
	return colorize(ColorGrey, format, args...)
}
