// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ntpdate

import (
	"bufio"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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
		_, s, err := getTime(tt.servers)
		if err == nil {
			t.Errorf("%v: got nil, want err", tt)
		}
		if !strings.HasPrefix(err.Error(), tt.err) {
			t.Errorf("expected:\n%s\ngot:\n%s", tt.err, err.Error())
		}
		if s != "" {
			t.Errorf("%v: got %q, want empty", tt, s)
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
		require.Error(t, err)
		require.Contains(t, err.Error(), "no servers")
		require.Equal(t, "", server)
		require.Equal(t, 0.0, offset)
		require.Equal(t, 0, m.getTimeCalls)
		require.Equal(t, 0, m.setSystemTimeCalls)
		require.Equal(t, 0, m.setRTCTimeCalls)
	}
	{ // Servers from cmd line first, then config, no fallback
		ts := time.Now()
		m := &mockGetterSetter{
			getTimeResult:       ts,
			getTimeResultServer: "foo",
		}
		server, offset, err := setTime([]string{"foo", "bar"}, "testdata/ntp.conf", "unused", false, m)
		require.Equal(t, "foo", server)
		require.NotEqual(t, 0.0, offset)
		require.NoError(t, err)
		require.Equal(t, 1, m.getTimeCalls)
		require.Equal(t, []string{"foo", "bar", "s1", "s2"}, m.getTimeArg)
		require.Equal(t, 1, m.setSystemTimeCalls)
		require.Equal(t, ts, m.setSystemTimeArg)
		require.Equal(t, 0, m.setRTCTimeCalls)
	}
	{ // Servers from config only, no fallback. Also sets RTC.
		ts := time.Now()
		m := &mockGetterSetter{
			getTimeResult:       ts,
			getTimeResultServer: "bar",
		}
		server, offset, err := setTime(nil, "testdata/ntp.conf", "unused", true, m)
		require.Equal(t, "bar", server)
		require.NotEqual(t, 0.0, offset)
		require.NoError(t, err)
		require.Equal(t, 1, m.getTimeCalls)
		require.Equal(t, []string{"s1", "s2"}, m.getTimeArg)
		require.Equal(t, 1, m.setSystemTimeCalls)
		require.Equal(t, ts, m.setSystemTimeArg)
		require.Equal(t, 1, m.setRTCTimeCalls)
		require.Equal(t, ts, m.setRTCTimeArg)
	}
	{ // Servers from cmdline only, no fallback.
		ts := time.Now()
		m := &mockGetterSetter{
			getTimeResult:       ts,
			getTimeResultServer: "foo",
		}
		server, offset, err := setTime([]string{"foo", "bar"}, "", "unused", false, m)
		require.Equal(t, "foo", server)
		require.NotEqual(t, 0.0, offset)
		require.NoError(t, err)
		require.Equal(t, 1, m.getTimeCalls)
		require.Equal(t, []string{"foo", "bar"}, m.getTimeArg)
		require.Equal(t, 1, m.setSystemTimeCalls)
		require.Equal(t, ts, m.setSystemTimeArg)
		require.Equal(t, 0, m.setRTCTimeCalls)
	}
	{ // Config not found, fallback is used.
		ts := time.Now()
		m := &mockGetterSetter{
			getTimeResult:       ts,
			getTimeResultServer: "HALP",
		}
		server, offset, err := setTime(nil, "testdata/nosuch.conf", "HALP", true, m)
		require.Equal(t, "HALP", server)
		require.NotEqual(t, 0.0, offset)
		require.NoError(t, err)
		require.Equal(t, 1, m.getTimeCalls)
		require.Equal(t, []string{"HALP"}, m.getTimeArg)
		require.Equal(t, 1, m.setSystemTimeCalls)
		require.Equal(t, ts, m.setSystemTimeArg)
		require.Equal(t, 1, m.setRTCTimeCalls)
		require.Equal(t, ts, m.setRTCTimeArg)
	}
	{ // Get NTP time fails, set not attempted.
		m := &mockGetterSetter{
			getTimeResult: time.Time{},
		}
		server, offset, err := setTime([]string{"foo", "bar"}, "", "unused", true, m)
		require.Equal(t, "", server)
		require.Equal(t, 0.0, offset)
		require.Error(t, err)
		require.Contains(t, err.Error(), "ASPLODE")
		require.Equal(t, 1, m.getTimeCalls)
		require.Equal(t, []string{"foo", "bar"}, m.getTimeArg)
		require.Equal(t, 0, m.setSystemTimeCalls)
		require.Equal(t, 0, m.setRTCTimeCalls)
	}
	{ // Set system time fails, set RTC not attempted.
		ts := time.Now()
		m := &mockGetterSetter{
			getTimeResult:       ts,
			getTimeResultServer: "foo",
			setSystemTimeResult: errors.New("ASPLODE"),
		}
		server, offset, err := setTime([]string{"foo", "bar"}, "", "unused", true, m)
		require.Equal(t, "", server)
		require.Equal(t, 0.0, offset)
		require.Error(t, err)
		require.Contains(t, err.Error(), "ASPLODE")
		require.Equal(t, 1, m.getTimeCalls)
		require.Equal(t, []string{"foo", "bar"}, m.getTimeArg)
		require.Equal(t, 1, m.setSystemTimeCalls)
		require.Equal(t, ts, m.setSystemTimeArg)
		require.Equal(t, 0, m.setRTCTimeCalls)
	}
	{ // Set RTC time fails.
		ts := time.Now()
		m := &mockGetterSetter{
			getTimeResult:       ts,
			getTimeResultServer: "foo",
			setRTCTimeResult:    errors.New("ASPLODE"),
		}
		server, offset, err := setTime([]string{"foo", "bar"}, "", "unused", true, m)
		require.Equal(t, "", server)
		require.Equal(t, 0.0, offset)
		require.Error(t, err)
		require.Contains(t, err.Error(), "ASPLODE")
		require.Equal(t, 1, m.getTimeCalls)
		require.Equal(t, []string{"foo", "bar"}, m.getTimeArg)
		require.Equal(t, 1, m.setSystemTimeCalls)
		require.Equal(t, ts, m.setSystemTimeArg)
		require.Equal(t, 1, m.setRTCTimeCalls)
		require.Equal(t, ts, m.setRTCTimeArg)
	}
}
