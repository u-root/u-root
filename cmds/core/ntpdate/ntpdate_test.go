// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"strings"
	"testing"
	"time"
)

var configFileTests = []struct {
	config string
	out    []string
}{
	{
		config: "",
		out:    []string{},
	},
	{
		config: "server 127.0.0.1",
		out:    []string{"127.0.0.1"},
	},
	{
		config: "server 127.0.0.1\n",
		out:    []string{"127.0.0.1"},
	},
	{
		config: "servers 127.0.0.1",
		out:    []string{},
	},
	{
		config: "server time.google.com",
		out:    []string{"time.google.com"},
	},
	{
		config: "server 127.0.0.1\n" +
			"server time.google.com",
		out: []string{"127.0.0.1", "time.google.com"},
	},
	{
		config: "servers 127.0.0.1\n" +
			"server time.google.com",
		out: []string{"time.google.com"},
	},
}

func TestConfigParsing(t *testing.T) {
	for _, tt := range configFileTests {
		b := bufio.NewReader(strings.NewReader(tt.config))
		out := parseServers(b)

		if len(out) != len(tt.out) {
			t.Errorf("Different lengths! Expected:\n%v\ngot:\n%v", tt.out, out)
		}

		for i := range out {
			if out[i] != tt.out[i] {
				t.Errorf("Element at index %d differs. expected:\n%v\ngot:\n%v", i, tt.out, out)
			}
		}
	}
}

var getTimeTests = []struct {
	servers []string
	time    time.Time
	err     string
}{
	{
		servers: []string{},
		err:     "unable to get any time from servers",
	},
	{
		servers: []string{"nope.nothing.here"},
		err:     "unable to get any time from servers",
	},
	{
		servers: []string{"nope.nothing.here", "nope.nothing.here2"},
		err:     "unable to get any time from servers",
	},
}

func TestGetNoTime(t *testing.T) {
	for _, tt := range getTimeTests {
		_, err := getTime(tt.servers)

		if err == nil {
			t.Errorf("%v: got nil, want err", tt)
		}
		if !strings.HasPrefix(err.Error(), tt.err) {
			t.Errorf("expected:\n%s\ngot:\n%s", tt.err, err.Error())
		}
	}
}
