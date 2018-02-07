dhcp6 [![Build Status](https://travis-ci.org/mdlayher/dhcp6.svg?branch=master)](https://travis-ci.org/mdlayher/dhcp6) [![Coverage Status](https://coveralls.io/repos/mdlayher/dhcp6/badge.svg?branch=master)](https://coveralls.io/r/mdlayher/dhcp6?branch=master) [![GoDoc](http://godoc.org/github.com/mdlayher/dhcp6?status.svg)](http://godoc.org/github.com/mdlayher/dhcp6)
=====

Package `dhcp6` implements a DHCPv6 server, as described in IETF RFC 3315.  MIT Licensed.

At this time, the API is not stable, and may change over time.  The eventual
goal is to implement a server, client, and testing facilities for consumers
of this package.

The design of this package is inspired by Go's `net/http` package.  The Go
standard library is Copyright (c) 2012 The Go Authors. All rights reserved.
The Go license can be found at https://golang.org/LICENSE.
