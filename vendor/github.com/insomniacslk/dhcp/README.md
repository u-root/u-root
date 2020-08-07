# dhcp
[![Build Status](https://travis-ci.org/insomniacslk/dhcp.svg?branch=master)](https://travis-ci.org/insomniacslk/dhcp)
[![GoDoc](https://godoc.org/github.com/insomniacslk/dhcp?status.svg)](https://godoc.org/github.com/insomniacslk/dhcp)
[![codecov](https://codecov.io/gh/insomniacslk/dhcp/branch/master/graph/badge.svg)](https://codecov.io/gh/insomniacslk/dhcp)
[![Go Report Card](https://goreportcard.com/badge/github.com/insomniacslk/dhcp)](https://goreportcard.com/report/github.com/insomniacslk/dhcp)

DHCPv4 and DHCPv6 decoding/encoding library with client and server code, written in Go.

# How to get the library

The library is split into several parts:
* `dhcpv6`: implementation of DHCPv6 packet, client and server
* `dhcpv4`: implementation of DHCPv4 packet, client and server
* `netboot`: network booting wrappers on top of `dhcpv6` and `dhcpv4`
* `iana`: several IANA constants, and helpers used by `dhcpv6` and `dhcpv4`
* `rfc1035label`: simple implementation of RFC1035 labels, used by `dhcpv6` and
  `dhcpv4`
* `interfaces`, a thin layer of wrappers around network interfaces

You will probably only need `dhcpv6` and/or `dhcpv4` explicitly. The rest is
pulled in automatically if necessary.


So, to get `dhcpv6` and `dhpv4` just run:
```
go get -u github.com/insomniacslk/dhcp/dhcpv{4,6}
```


# Examples

The sections below will illustrate how to use the `dhcpv6` and `dhcpv4`
packages.

* [dhcpv6 client](examples/client6/)
* [dhcpv6 server](examples/server6/)
* [dhcpv6 packet crafting](examples/packetcrafting6)
* TODO dhcpv4 client
* TODO dhcpv4 server
* TODO dhcpv4 packet crafting


See more example code at https://github.com/insomniacslk/exdhcp


# Public projects that use it

* Facebook's DHCP load balancer, `dhcplb`, https://github.com/facebookincubator/dhcplb
* Systemboot, a LinuxBoot distribution that runs as system firmware, https://github.com/systemboot/systemboot
* Router7, a pure-Go router implementation for fiber7 connections, https://github.com/rtr7/router7
* Beats from ElasticSearch, https://github.com/elastic/beats
* Bender from Pinterest, a library for load-testing, https://github.com/pinterest/bender
* FBender from Facebook, a tool for load-testing based on Bender, https://github.com/facebookincubator/fbender
* CoreDHCP, a fast, multithreaded, modular and extensible DHCP server, https://github.com/coredhcp/coredhcp
* u-root, an embeddable root file system, https://github.com/u-root/u-root
* Talos: a modern OS for Kubernetes, https://github.com/talos-systems/talos
