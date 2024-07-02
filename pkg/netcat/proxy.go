// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package netcat

import (
	"fmt"
	"strings"
)

type ProxyType int

const (
	PROXY_TYPE_NONE ProxyType = iota
	PROXY_TYPE_HTTP
	PROXY_TYPE_SOCKS4
	PROXY_TYPE_SOCKS5
	DEFAULT_PROXY_TYPE = PROXY_TYPE_NONE
)

func ProxyTypeFromString(s string) ProxyType {
	switch strings.ToUpper(s) {
	case "HTTP":
		return PROXY_TYPE_HTTP
	case "SOCKS4":
		return PROXY_TYPE_SOCKS4
	case "SOCKS5":
		return PROXY_TYPE_SOCKS5
	default:
		return PROXY_TYPE_NONE
	}
}

func (p ProxyType) String() string {
	return [...]string{
		"None",
		"http",
		"socks4",
		"socks5",
	}[p]
}

func (p ProxyType) DefaultPort() (uint, error) {
	switch p {
	case PROXY_TYPE_SOCKS4, PROXY_TYPE_SOCKS5:
		return 1080, nil
	case PROXY_TYPE_HTTP:
		return 3128, nil
	default:
		return 0, fmt.Errorf("ProxyType %s has no default port", p.String())
	}
}

type ProxyDNSType int

const (
	PROXY_DNS_NONE ProxyDNSType = iota
	PROXY_DNS_LOCAL
	PROXY_DNS_REMOTE
	PROXY_DNS_BOTH
)

func ProxyDNSTypeFromString(s string) ProxyDNSType {
	switch strings.ToUpper(s) {
	case "LOCAL":
		return PROXY_DNS_LOCAL
	case "REMOTE":
		return PROXY_DNS_REMOTE
	case "BOTH":
		return PROXY_DNS_BOTH
	default:
		return PROXY_DNS_NONE
	}
}

type ProxyOptions struct {
	Enabled bool
	Type    ProxyType
	Address string // if address is empty, discard the entire proxy handling
	Auth    string
	DNSType ProxyDNSType
}
