# dhcp4

[![CircleCI](https://circleci.com/gh/u-root/dhcp4.svg?style=svg)](https://circleci.com/gh/u-root/dhcp4) [![Go Report Card](https://goreportcard.com/badge/github.com/u-root/dhcp4)](https://goreportcard.com/report/github.com/u-root/dhcp4) [![GoDoc](https://godoc.org/github.com/u-root/dhcp4?status.svg)](https://godoc.org/github.com/u-root/dhcp4) [![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://github.com/u-root/dhcp4/blob/master/LICENSE)

Package `dhcp4` is an IPv4 DHCP library as described in RFC 2131, 2132, and 3396.

It implements encoding and decoding of DHCP messages in `dhcp4`. Option parsing is in the `dhcp4opts` package; a simple client is included in `dhcp4client`. Some day, there may be a server.

If you are already using another IPv4 DHCP library like [krolaw's](https://github.com/krolaw/dhcp4), you can still use `dhcp4opts` to decode options not implemented in krolaw's DHCP library.
