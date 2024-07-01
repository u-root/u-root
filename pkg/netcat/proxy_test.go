// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package netcat

import "testing"

func TestProxyType_String(t *testing.T) {
	tests := []struct {
		name     string
		p        ProxyType
		expected string
	}{
		{name: "Test None", p: PROXY_TYPE_NONE, expected: "None"},
		{name: "Test HTTP", p: PROXY_TYPE_HTTP, expected: "http"},
		{name: "Test SOCKS4", p: PROXY_TYPE_SOCKS4, expected: "socks4"},
		{name: "Test SOCKS5", p: PROXY_TYPE_SOCKS5, expected: "socks5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.String(); got != tt.expected {
				t.Errorf("ProxyType.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestProxyType_DefaultPort(t *testing.T) {
	tests := []struct {
		name        string
		p           ProxyType
		expected    uint
		expectError bool
	}{
		{name: "SOCKS5 Default Port", p: PROXY_TYPE_SOCKS5, expected: 1080, expectError: false},
		{name: "HTTP Default Port", p: PROXY_TYPE_HTTP, expected: 3128, expectError: false},
		{name: "Unsupported Proxy Type", p: PROXY_TYPE_NONE, expected: 0, expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.DefaultPort()
			if (err != nil) != tt.expectError {
				t.Errorf("ProxyType.DefaultPort() error = %v, wantErr %v", err, tt.expectError)
				return
			}
			if got != tt.expected {
				t.Errorf("ProxyType.DefaultPort() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestProxyTypeFromString(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected ProxyType
	}{
		{name: "HTTP lowercase", s: "http", expected: PROXY_TYPE_HTTP},
		{name: "SOCKS4 mixed case", s: "Socks4", expected: PROXY_TYPE_SOCKS4},
		{name: "SOCKS5 uppercase", s: "SOCKS5", expected: PROXY_TYPE_SOCKS5},
		{name: "Invalid type", s: "invalid", expected: PROXY_TYPE_NONE},
		{name: "Empty string", s: "", expected: PROXY_TYPE_NONE},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProxyTypeFromString(tt.s); got != tt.expected {
				t.Errorf("ProxyTypeFromString(%v) = %v, want %v", tt.s, got, tt.expected)
			}
		})
	}
}

func TestProxyDNSTypeFromString(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected ProxyDNSType
	}{
		{name: "LOCAL uppercase", s: "LOCAL", expected: PROXY_DNS_LOCAL},
		{name: "remote lowercase", s: "remote", expected: PROXY_DNS_REMOTE},
		{name: "Both MixedCase", s: "Both", expected: PROXY_DNS_BOTH},
		{name: "invalid type", s: "invalid", expected: PROXY_DNS_NONE},
		{name: "Empty string", s: "", expected: PROXY_DNS_NONE},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProxyDNSTypeFromString(tt.s); got != tt.expected {
				t.Errorf("ProxyDNSTypeFromString(%v) = %v, want %v", tt.s, got, tt.expected)
			}
		})
	}
}
