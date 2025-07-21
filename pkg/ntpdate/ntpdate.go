// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !windows

package ntpdate

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/beevik/ntp"

	"github.com/u-root/u-root/pkg/rtc"
)

const DefaultNTPConfig = "/etc/ntp.conf"

var Debug = func(string, ...any) {}

func parseServers(r *bufio.Reader) []string {
	var uri []string
	var l string
	var err error

	Debug("Reading config file")
	for err == nil {
		// This handles the case where the last line doesn't end in \n
		l, err = r.ReadString('\n')
		Debug("%v", l)
		if w := strings.Fields(l); len(w) > 1 && w[0] == "server" {
			// We look only for the server lines, we ignore options like iburst
			// TODO(ganshun): figure out what options we want to support.
			uri = append(uri, w[1])
		}
	}

	return uri
}

func getTime(servers []string) (time.Time, string, error) {
	for _, s := range servers {
		Debug("Getting time from %v", s)
		t, err := ntp.Time(s)
		if err == nil {
			// Right now we return on the first valid time.
			// We can implement better heuristics here.
			Debug("Got time %v", t)
			return t, s, nil
		}
		Debug("Error getting time from %s: %v", s, err)
	}

	return time.Time{}, "", fmt.Errorf("unable to get any time from servers %v", servers)
}

// SetTime sets system and optionally RTC time from NTP servers specified in sersers or the config file.
// If successful, returns the server used to set the time and the offset, in seconds.
func SetTime(servers []string, config string, fallback string, setRTC bool) (string, float64, error) {
	return setTime(servers, config, fallback, setRTC, &realGetterSetter{})
}

type timeGetterSetter interface {
	GetTime(servers []string) (time.Time, string, error)
	SetSystemTime(time.Time) error
	SetRTCTime(time.Time) error
}

type realGetterSetter struct{}

func (*realGetterSetter) GetTime(servers []string) (time.Time, string, error) {
	return getTime(servers)
}

func (*realGetterSetter) SetSystemTime(t time.Time) error {
	tv := syscall.NsecToTimeval(t.UnixNano())
	return syscall.Settimeofday(&tv)
}

func (*realGetterSetter) SetRTCTime(t time.Time) error {
	r, err := rtc.OpenRTC()
	if err != nil {
		return fmt.Errorf("unable to open RTC: %w", err)
	}
	defer r.Close()
	return r.Set(t)
}

func setTime(servers []string, config string, fallback string, setRTC bool, gs timeGetterSetter) (string, float64, error) {
	servers = servers[:]

	if config != "" {
		Debug("Reading NTP servers from config file: %v", config)
		f, err := os.Open(config)
		if err == nil {
			defer f.Close()
			configServers := parseServers(bufio.NewReader(f))
			Debug("Found %v servers", len(configServers))
			servers = append(servers, configServers...)
		} else {
			Debug("Unable to open config file: %v", err)
		}
	}

	if len(servers) == 0 && len(fallback) != 0 {
		Debug("No servers provided, falling back to %v", fallback)
		servers = append(servers, fallback)
	}

	if len(servers) == 0 {
		return "", 0, fmt.Errorf("no servers")
	}

	t, server, err := gs.GetTime(servers)
	if err != nil {
		return "", 0, fmt.Errorf("unable to get time: %w", err)
	}

	offset := time.Until(t).Seconds()

	if err = gs.SetSystemTime(t); err != nil {
		return "", 0, fmt.Errorf("unable to set system time: %w", err)
	}
	if setRTC {
		Debug("Setting RTC time...")
		if err = gs.SetRTCTime(t); err != nil {
			return "", 0, fmt.Errorf("unable to set RTC time: %w", err)
		}
	}

	return server, offset, nil
}
