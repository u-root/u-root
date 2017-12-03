package gzip

import (
	"fmt"
	"io"
)

type errType int

const (
	info     errType = 0
	skipping errType = 1
	warning  errType = 2
	fatal    errType = 3
)

type appError struct {
	msg   string
	level errType
	path  string
}

func (e *appError) Error() string {
	if e.path != "" {
		return e.path + " " + e.msg
	}
	return e.msg
}

func ErrorHandler(err error, stdout io.Writer, stderr io.Writer) int {
	switch err.(type) {
	case *appError:
		switch err.(*appError).level {
		case fatal:
			fmt.Fprintf(stderr, "fatal, %s\n", err)
			return 1
		case warning:
			fmt.Fprintf(stderr, "warning, %s\n", err)
			return 0
		case skipping:
			fmt.Fprintf(stderr, "skipping, %s\n", err)
			return 0
		case info:
			fmt.Fprintf(stdout, "%s\n", err)
			return 0
		}
	case error:
		fmt.Fprintf(stderr, "error, %s\n", err)
		return 1
	}
	return 0
}
