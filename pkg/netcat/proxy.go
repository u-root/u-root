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
)

func (p ProxyType) String() string {
	return [...]string{
		"None",
		"HTTP",
		"SOCKS4",
		"SOCKS5",
	}[p]
}


func (p ProxyType) DefaultPort() (uint, error) {
	switch p {
	case PROXY_TYPE_SOCKS5:
		return 1080, nil
	case PROXY_TYPE_HTTP:
		return 3128, nil
	default:
		return 0, fmt.Errorf("ProxyType %s has no default port", p.String())
	}
}

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

type ProxyAuthType int

const (
	PROXY_AUTH_NONE ProxyAuthType = iota
	PROXY_AUTH_HTTP
	PROXY_AUTH_SOCKS5
)

func ProxyAuthTypeFromString(s string) ProxyAuthType {
	switch strings.ToUpper(s) {
	case "HTTP":
		return PROXY_AUTH_HTTP
	case "SOCKS5":
		return PROXY_AUTH_SOCKS5
	default:
		return PROXY_AUTH_NONE
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
	Type     ProxyType // If this is none, discard the entire Proxy handling
	Address  string
	DNSType  ProxyDNSType
	Port     uint
	AuthType ProxyAuthType // If this is none, discard the entire ProxyAuth handling
}
