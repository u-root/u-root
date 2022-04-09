// Copyright 2017-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ntpdate

import (
	"bufio"
	"errors"
	"reflect"
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
			t.Errorf(`len(%v) = %d, want %d`, out, len(out), len(tt.out))
		}

		for i := range out {
			if out[i] != tt.out[i] {
				t.Errorf(`out[%d] = %v, want %v`, i, out[i], tt.out[i])
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
		_, s, err := getTime(tt.servers)
		if err == nil || s != "" {
			t.Errorf(`getTime(%v) = _, %q, %v, want "", not nil`, tt.servers, s, err)
		}
		if match := strings.HasPrefix(err.Error(), tt.err); !match {
			t.Errorf(`strings.HasPrefix(%q, %v) = %t, want true`, err.Error(), tt.err, match)
		}
	}
}

type mockGetterSetter struct {
	getTimeCalls        int
	getTimeArg          []string
	getTimeResult       time.Time
	getTimeResultServer string
	setSystemTimeCalls  int
	setSystemTimeArg    time.Time
	setSystemTimeResult error
	setRTCTimeCalls     int
	setRTCTimeArg       time.Time
	setRTCTimeResult    error
}

func (mgs *mockGetterSetter) GetTime(servers []string) (time.Time, string, error) {
	mgs.getTimeCalls++
	mgs.getTimeArg = servers
	if mgs.getTimeResult.Equal(time.Time{}) {
		return mgs.getTimeResult, "", errors.New("ASPLODE")
	}
	return mgs.getTimeResult, mgs.getTimeResultServer, nil
}

func (mgs *mockGetterSetter) SetSystemTime(t time.Time) error {
	mgs.setSystemTimeCalls++
	mgs.setSystemTimeArg = t
	return mgs.setSystemTimeResult
}

func (mgs *mockGetterSetter) SetRTCTime(t time.Time) error {
	mgs.setRTCTimeCalls++
	mgs.setRTCTimeArg = t
	return mgs.setRTCTimeResult
}

func TestSetTime(t *testing.T) {
	{ // No args, no config, no fallback - fail
		m := &mockGetterSetter{}
		server, offset, err := setTime(nil, "", "", true, m)
		if err == nil {
			t.Fatalf(`setTime(nil, "", "", true, %v) = _, _, %v, want not nil`, m, err)
		}
		if match := strings.Contains(err.Error(), "no servers"); !match {
			t.Errorf(`strings.Contains(%q, "no servers") = %t, want true`, err.Error(), match)
		}
		if server != "" || offset != 0.0 {
			t.Errorf(`setTime(nil, "", "", true, %v) = %q, %f, _, want "", 0.0`, m, server, offset)
		}
		if m.getTimeCalls != 0 || m.setSystemTimeCalls != 0 || m.setRTCTimeCalls != 0 {
			t.Errorf(`m.getTimeCalls, m.setSystemTimeCalls, m.setRTCTimeCalls = %d, %d, %d, want 0, 0, 0`,
				m.getTimeCalls, m.setSystemTimeCalls, m.setRTCTimeCalls)
		}
	}
	{ // Servers from cmd line first, then config, no fallback
		ts := time.Now()
		m := &mockGetterSetter{
			getTimeResult:       ts,
			getTimeResultServer: "foo",
		}
		server, offset, err := setTime([]string{"foo", "bar"}, "testdata/ntp.conf", "unused", false, m)
		if err != nil || server != "foo" || offset == 0.0 {
			t.Errorf(`setTime([]string{"foo", "bar"}, "testdata/ntp.conf", "unused", false, %v) = %q, %f, %v, want "foo", not 0.0, nil`,
				m, server, offset, err)
		}
		if m.getTimeCalls != 1 || m.setSystemTimeCalls != 1 {
			t.Errorf(`m.getTimeCalls, m.setSystemTimeCalls = %d, %d, want 1, 1`, m.getTimeCalls, m.setSystemTimeCalls)
		}
		if match := reflect.DeepEqual([]string{"foo", "bar", "s1", "s2"}, m.getTimeArg); !match {
			t.Errorf(`reflect.DeepEqual([]string{"foo", "bar", "s1", "s2"}, %v) = %t, want true`, m.getTimeArg, match)
		}
		if m.setSystemTimeArg != ts || m.setRTCTimeCalls != 0 {
			t.Errorf(`m.setSystemTimeArg, m.setRTCTimeCalls = %v, %d, want %v, 0`, m.setSystemTimeArg, m.setRTCTimeCalls, ts)
		}
	}
	{ // Servers from config only, no fallback. Also sets RTC.
		ts := time.Now()
		m := &mockGetterSetter{
			getTimeResult:       ts,
			getTimeResultServer: "bar",
		}
		server, offset, err := setTime(nil, "testdata/ntp.conf", "unused", true, m)
		if err != nil || server != "bar" || offset == 0.0 {
			t.Errorf(`setTime(nil, "testdata/ntp.conf", "unused", true, %v) = %q, %f, %v, want "bar", not 0.0, nil`,
				m, server, offset, err)
		}
		if m.getTimeCalls != 1 || m.setSystemTimeCalls != 1 || m.setRTCTimeCalls != 1 {
			t.Errorf(`m.getTimeCalls, m.setSystemTimeCalls, m.setRTCTimeCalls = %d, %d, %d, want 1, 1, 1`,
				m.getTimeCalls, m.setSystemTimeCalls, m.setRTCTimeCalls)
		}
		if match := reflect.DeepEqual([]string{"s1", "s2"}, m.getTimeArg); !match {
			t.Errorf(`reflect.DeepEqual([]string{"s1", "s2"}, %v) = %t, want true`, m.getTimeArg, match)
		}
		if m.setSystemTimeArg != ts || m.setRTCTimeArg != ts {
			t.Errorf(`m.setSystemTimeArg, m.setRTCTimeArg = %v, %v, want %v, %v`, m.setSystemTimeArg, m.setRTCTimeArg, ts, ts)
		}
	}
	{ // Servers from cmdline only, no fallback.
		ts := time.Now()
		m := &mockGetterSetter{
			getTimeResult:       ts,
			getTimeResultServer: "foo",
		}
		server, offset, err := setTime([]string{"foo", "bar"}, "", "unused", false, m)
		if err != nil || server != "foo" || offset == 0.0 {
			t.Errorf(`setTime([]string{"foo", "bar"}, "", "unused", false, %v) = %q, %f, %v, want "foo", not 0.0, nil`,
				m, server, offset, err)
		}
		if m.getTimeCalls != 1 || m.setSystemTimeCalls != 1 {
			t.Errorf(`m.getTimeCalls, m.setSystemTimeCalls = %d, %d, want 1, 1`, m.getTimeCalls, m.setSystemTimeCalls)
		}
		if match := reflect.DeepEqual([]string{"foo", "bar"}, m.getTimeArg); !match {
			t.Errorf(`reflect.DeepEqual([]string{"foo", "bar"}, %v) = %t, want true`, m.getTimeArg, match)
		}
		if m.setSystemTimeArg != ts || m.setRTCTimeCalls != 0 {
			t.Errorf(`m.setSystemTimeArg, m.setRTCTimeCalls = %v, %d, want %v, 0`, m.setSystemTimeArg, m.setRTCTimeCalls, ts)
		}
	}
	{ // Config not found, fallback is used.
		ts := time.Now()
		m := &mockGetterSetter{
			getTimeResult:       ts,
			getTimeResultServer: "HALP",
		}
		server, offset, err := setTime(nil, "testdata/nosuch.conf", "HALP", true, m)
		if err != nil || server != "HALP" || offset == 0.0 {
			t.Errorf(`setTime(nil, "testdata/nosuch.conf", "HALP", true, %v) = %q, %f, %v, want "HALP", not 0.0, nil`,
				m, server, offset, err)
		}
		if m.getTimeCalls != 1 || m.setSystemTimeCalls != 1 || m.setRTCTimeCalls != 1 {
			t.Errorf(`m.getTimeCalls, m.setSystemTimeCalls, m.setRTCTimeCalls = %d, %d, %d, want 1, 1, 1`,
				m.getTimeCalls, m.setSystemTimeCalls, m.setRTCTimeCalls)
		}
		if match := reflect.DeepEqual([]string{"HALP"}, m.getTimeArg); !match {
			t.Errorf(`reflect.DeepEqual([]string{"HALP"}, %v) = %t, want true`, m.getTimeArg, match)
		}
		if m.setSystemTimeArg != ts || m.setRTCTimeArg != ts {
			t.Errorf(`m.setSystemTimeArg, m.setRTCTimeArg = %v, %v, want %v, %v`, m.setSystemTimeArg, m.setRTCTimeArg, ts, ts)
		}
	}
	{ // Get NTP time fails, set not attempted.
		m := &mockGetterSetter{
			getTimeResult: time.Time{},
		}
		server, offset, err := setTime([]string{"foo", "bar"}, "", "unused", true, m)
		if err == nil || server != "" || offset != 0.0 {
			t.Errorf(`setTime([]string{"foo", "bar"}, "", "unused", true,  %v) = %q, %f, %v, want "", 0.0, not nil`,
				m, server, offset, err)
		}
		if match := strings.Contains(err.Error(), "ASPLODE"); !match {
			t.Errorf(`strings.Contains(%q, "ASPLODE") = %t, want true`, err.Error(), match)
		}
		if match := reflect.DeepEqual([]string{"foo", "bar"}, m.getTimeArg); !match {
			t.Errorf(`reflect.DeepEqual([]string{"foo", "bar"}, %v)  = %t, want true`, m.getTimeArg, match)
		}
		if m.getTimeCalls != 1 || m.setSystemTimeCalls != 0 || m.setRTCTimeCalls != 0 {
			t.Errorf(`m.getTimeCalls, m.setSystemTimeCalls, m.setRTCTimeCalls = %d, %d, %d, want 1, 0, 0`,
				m.getTimeCalls, m.setSystemTimeCalls, m.setRTCTimeCalls)
		}
	}
	{ // Set system time fails, set RTC not attempted.
		ts := time.Now()
		m := &mockGetterSetter{
			getTimeResult:       ts,
			getTimeResultServer: "foo",
			setSystemTimeResult: errors.New("ASPLODE"),
		}
		server, offset, err := setTime([]string{"foo", "bar"}, "", "unused", true, m)
		if err == nil || server != "" || offset != 0.0 {
			t.Errorf(`setTime([]string{"foo", "bar"}, "", "unused", true,  %v) = %q, %f, %v, want "", 0.0, not nil`,
				m, server, offset, err)
		}
		if match := strings.Contains(err.Error(), "ASPLODE"); !match {
			t.Errorf(`strings.Contains(%q, "ASPLODE") = %t, want true`, err.Error(), match)
		}
		if match := reflect.DeepEqual([]string{"foo", "bar"}, m.getTimeArg); !match {
			t.Errorf(`reflect.DeepEqual([]string{"foo", "bar"}, %v)  = %t, want true`, m.getTimeArg, match)
		}
		if m.getTimeCalls != 1 || m.setSystemTimeCalls != 1 || m.setRTCTimeCalls != 0 {
			t.Errorf(`m.getTimeCalls, m.setSystemTimeCalls, m.setRTCTimeCalls = %d, %d, %d, want 1, 1, 0`,
				m.getTimeCalls, m.setSystemTimeCalls, m.setRTCTimeCalls)
		}
		if !m.setSystemTimeArg.Equal(ts) {
			t.Error("setSystemTimeArg is not correct")
		}
		if m.setSystemTimeArg != ts {
			t.Errorf(`m.setSystemTimeArg, = %v, want %v`, m.setSystemTimeArg, ts)
		}
	}
	{ // Set RTC time fails.
		ts := time.Now()
		m := &mockGetterSetter{
			getTimeResult:       ts,
			getTimeResultServer: "foo",
			setRTCTimeResult:    errors.New("ASPLODE"),
		}
		server, offset, err := setTime([]string{"foo", "bar"}, "", "unused", true, m)
		if err == nil || server != "" || offset != 0.0 {
			t.Errorf(`setTime([]string{"foo", "bar"}, "", "unused", true,  %v) = %q, %f, %v, want "", 0.0, not nil`,
				m, server, offset, err)
		}
		if match := strings.Contains(err.Error(), "ASPLODE"); !match {
			t.Errorf(`strings.Contains(%q, "ASPLODE") = %t, want true`, err.Error(), match)
		}
		if m.getTimeCalls != 1 || m.setSystemTimeCalls != 1 || m.setRTCTimeCalls != 1 {
			t.Errorf(`m.getTimeCalls, m.setSystemTimeCalls, m.setRTCTimeCalls = %d, %d, %d, want 1, 1, 1`,
				m.getTimeCalls, m.setSystemTimeCalls, m.setRTCTimeCalls)
		}
		if match := reflect.DeepEqual([]string{"foo", "bar"}, m.getTimeArg); !match {
			t.Errorf(`reflect.DeepEqual([]string{"foo", "bar"}, %v)  = %t, want true`, m.getTimeArg, match)
		}
		if m.setSystemTimeArg != ts || m.setRTCTimeArg != ts {
			t.Errorf(`m.setSystemTimeArg, m.setRTCTimeArg = %v, %v, want %v, %v`, m.setSystemTimeArg, m.setRTCTimeArg, ts, ts)
		}
	}
}
