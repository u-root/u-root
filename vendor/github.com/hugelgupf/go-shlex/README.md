# go-shlex

[![CircleCI](https://circleci.com/gh/hugelgupf/go-shlex.svg?style=svg)](https://circleci.com/gh/hugelgupf/go-shlex)
[![Go Report Card](https://goreportcard.com/badge/github.com/hugelgupf/go-shlex)](https://goreportcard.com/report/github.com/hugelgupf/go-shlex)
[![GoDoc](https://godoc.org/github.com/hugelgupf/go-shlex?status.svg)](https://godoc.org/github.com/hugelgupf/go-shlex)

go-shlex is a POSIX command-line shell-like argument parser.

### Differences

-   [anmitsu/go-shlex](https://github.com/anmitsu/go-shlex): anmitsu does not
    support comments (#) and double-quoted dollar ($) and backtick (`)
    characters.

-   [google/shlex](https://github.com/google/shlex): google does not support
    Unicode spaces and double-quoted newlines (\n) and backslashes (\\). google
    also stops parsing upon error, while we (and anmitsu) will return partial
    results.
