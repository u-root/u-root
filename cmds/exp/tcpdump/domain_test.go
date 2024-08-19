// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"testing"

	"github.com/gopacket/gopacket/layers"
)

func TestDnsData(t *testing.T) {
	tests := []struct {
		name     string
		layer    *layers.DNS
		expected string
	}{
		{
			name: "No error response code, no answers",
			layer: &layers.DNS{
				ID:           1234,
				ResponseCode: layers.DNSResponseCodeNoErr,
				Questions: []layers.DNSQuestion{
					{Type: layers.DNSTypeA, Name: []byte("example.com")},
				},
			},
			expected: "1234 A? example.com (0)",
		},
		{
			name: "Error response code",
			layer: &layers.DNS{
				ID:           5678,
				ResponseCode: layers.DNSResponseCodeNXDomain,
			},
			expected: "5678 Non-Existent Domain (0)",
		},
		{
			name: "answer",
			layer: &layers.DNS{
				ID:           9101,
				ResponseCode: layers.DNSResponseCodeNoErr,
				AA:           true,
				Answers: []layers.DNSResourceRecord{
					{Name: []byte("example.com"), Type: layers.DNSTypeA, Class: layers.DNSClassAny, IP: []byte{192, 0, 2, 1}},
				},
			},
			expected: "9101*<Any, A> (0)",
		},
		{
			name: " 2 answers",
			layer: &layers.DNS{
				ID:           9101,
				ResponseCode: layers.DNSResponseCodeNoErr,
				AA:           true,
				Answers: []layers.DNSResourceRecord{
					{Name: []byte("example.com"), Type: layers.DNSTypeA, Class: layers.DNSClassAny, IP: []byte{192, 0, 2, 1}},
					{Name: []byte("example.io"), Type: layers.DNSTypeA, Class: layers.DNSClassAny, IP: []byte{192, 0, 2, 2}},
				},
			},
			expected: "9101* 2/0/0 <Any, A>, <Any, A> (0)",
		},
		{
			name: "Recursive desired",
			layer: &layers.DNS{
				ID:           1213,
				ResponseCode: layers.DNSResponseCodeNoErr,
				RD:           true,
				Answers: []layers.DNSResourceRecord{
					{Name: []byte("example.com"), Type: layers.DNSTypeA, Class: layers.DNSClassAny, IP: []byte{192, 0, 2, 1}},
				},
			},
			expected: "1213+<Any, A> (0)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := dnsData(tt.layer)
			if result != tt.expected {
				t.Errorf("dnsData() = %v, want %v", result, tt.expected)
			}
		})
	}
}
