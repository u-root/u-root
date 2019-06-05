// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"reflect"
	"testing"
)

func TestIwconfigRE(t *testing.T) {
	testcases := []struct {
		s   string
		exp bool
	}{
		{"blahblahblah\nlo        no wireless extensions.\n", false},
		{"blahblahblah\nwlp4s0    IEEE 802.11  ESSID:\"stub\"", true},
		{"blahblahblah\n          Mode:Managed  Frequency:5.58 GHz  Access Point: 00:00:00:00:00:00\n", false},
	}
	for _, test := range testcases {
		if out := iwconfigRE.MatchString(test.s); out != test.exp {
			t.Errorf("%s\ngot:%v\nwant:%v", test.s, out, test.exp)
		}
	}
}

func TestParseIwconfig(t *testing.T) {
	testcases := []struct {
		name string
		o    []byte
		exp  []string
	}{
		{
			name: "nil input",
			o:    nil,
			exp:  nil,
		},
		{
			name: "empty string input",
			o:    []byte(""),
			exp:  nil,
		},
		{
			name: "No Wireless in input",
			o: []byte(`
lo        no wireless extensions.

eno1      no wireless extensions.

`),
			exp: nil,
		},
		{
			name: "One (1) Wireless extension",
			o: []byte(`
wlp4s0    IEEE 802.11  ESSID:"stub"
          Mode:Managed  Frequency:5.58 GHz  Access Point: 00:00:00:00:00:00   
          Bit Rate=22 Mb/s   Tx-Power=22 dBm   
          Retry short limit:7   RTS thr:off   Fragment thr:off
          Power Management:on
          Link Quality=27/70  Signal level=-53 dBm  
          Rx invalid nwid:0  Rx invalid crypt:0  Rx invalid frag:0
          Tx excessive retries:0  Invalid misc:0   Missed beacon:0

lo        no wireless extensions.

enp0s31f6  no wireless extensions.

`),
			exp: []string{"wlp4s0"},
		},
		{
			name: "Two (2) Wireless extensions",
			o: []byte(`
wlp4s0    IEEE 802.11  ESSID:"stub"
          Mode:Managed  Frequency:5.58 GHz  Access Point: 00:00:00:00:00:00   
          Bit Rate=22 Mb/s   Tx-Power=22 dBm   
          Retry short limit:7   RTS thr:off   Fragment thr:off
          Power Management:on
          Link Quality=27/70  Signal level=-53 dBm  
          Rx invalid nwid:0  Rx invalid crypt:0  Rx invalid frag:0
          Tx excessive retries:0  Invalid misc:0   Missed beacon:0

wlp4s1    IEEE 802.11  ESSID:"stub"
          Mode:Managed  Frequency:5.58 GHz  Access Point: 00:00:00:00:00:00   
          Bit Rate=22 Mb/s   Tx-Power=22 dBm   
          Retry short limit:7   RTS thr:off   Fragment thr:off
          Power Management:on
          Link Quality=27/70  Signal level=-53 dBm  
          Rx invalid nwid:0  Rx invalid crypt:0  Rx invalid frag:0
          Tx excessive retries:0  Invalid misc:0   Missed beacon:0

lo        no wireless extensions.

enp0s31f6  no wireless extensions.

`),
			exp: []string{"wlp4s0", "wlp4s1"},
		},
	}

	for _, test := range testcases {
		out := parseIwconfig(test.o)
		if !reflect.DeepEqual(out, test.exp) {
			t.Errorf("%v\ngot:%v\nwant:%v", test.name, out, test.exp)
		}
	}
}
