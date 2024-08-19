// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"bytes"
	"net"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/vishvananda/netlink"
)

func TestParseXfrmPolicyTmpl(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []string
		expected netlink.XfrmPolicyTmpl
		wantErr  bool
	}{
		{
			name:   "valid input",
			tokens: []string{"src", "192.168.1.1", "dst", "192.168.1.2", "proto", "esp", "spi", "12345", "mode", "tunnel", "reqid", "67890", "level", "use"},
			expected: netlink.XfrmPolicyTmpl{
				Src:      net.ParseIP("192.168.1.1"),
				Dst:      net.ParseIP("192.168.1.2"),
				Proto:    netlink.XFRM_PROTO_ESP,
				Spi:      12345,
				Mode:     netlink.XFRM_MODE_TUNNEL,
				Reqid:    67890,
				Optional: 1,
			},
			wantErr: false,
		},
		{
			name:    "invalid src",
			tokens:  []string{"src", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid dst",
			tokens:  []string{"dst", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid spi",
			tokens:  []string{"spi", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid mode",
			tokens:  []string{"mode", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid proto",
			tokens:  []string{"proto", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid reqid",
			tokens:  []string{"reqid", "abc"},
			wantErr: true,
		},
		{
			name:   "required level",
			tokens: []string{"level", "required"},
			expected: netlink.XfrmPolicyTmpl{
				Optional: 0,
			},
		},
		{
			name:    "invalid level",
			tokens:  []string{"level", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid arg",
			tokens:  []string{"arg", "abc"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			cmd := cmd{
				Cursor: -1,
				Args:   tt.tokens,
				Out:    out,
			}

			got, err := cmd.parseXfrmPolicyTmpl()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseXfrmPolicyTmpl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !reflect.DeepEqual(*got, tt.expected) {
					t.Errorf("parseXfrmPolicyTmpl() = %v, expected %v", got, tt.expected)
				}
			}
		})
	}
}

func TestParseXfrmPolicyAddUpdate(t *testing.T) {
	_, src, err := net.ParseCIDR("192.168.1.1/24")
	if err != nil {
		t.Fatal(err)
	}

	_, dst, err := net.ParseCIDR("192.168.1.2/24")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		tokens   []string
		expected netlink.XfrmPolicy
		wantErr  bool
	}{
		{
			name:   "valid input",
			tokens: []string{"src", "192.168.1.1/24", "dst", "192.168.1.2/24", "proto", "esp", "sport", "1234", "dport", "5678", "dir", "in", "mark", "1", "index", "1", "action", "allow", "priority", "1", "if_id", "1", "tmpl", "proto", "esp"},
			expected: netlink.XfrmPolicy{
				Src:      src,
				Dst:      dst,
				Proto:    netlink.XFRM_PROTO_ESP,
				SrcPort:  1234,
				DstPort:  5678,
				Dir:      netlink.XFRM_DIR_IN,
				Mark:     &netlink.XfrmMark{Value: 1},
				Index:    1,
				Action:   netlink.XFRM_POLICY_ALLOW,
				Priority: 1,
				Ifid:     1,
				Tmpls: []netlink.XfrmPolicyTmpl{
					{
						Proto: netlink.XFRM_PROTO_ESP,
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "invalid arg",
			tokens:  []string{"dsty", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid dst",
			tokens:  []string{"dst", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid src",
			tokens:  []string{"src", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid proto",
			tokens:  []string{"proto", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid sport",
			tokens:  []string{"sport", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid dport",
			tokens:  []string{"dport", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid dir",
			tokens:  []string{"dir", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid mark",
			tokens:  []string{"mark", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid index",
			tokens:  []string{"index", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid action",
			tokens:  []string{"action", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid priority",
			tokens:  []string{"priority", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid if_id",
			tokens:  []string{"if_id", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid tmpl",
			tokens:  []string{"tmpl", "abc"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := cmd{
				Cursor: -1,
				Args:   tt.tokens,
			}
			got, err := cmd.parseXfrmPolicyAddUpdate()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseXfrmPolicyAddUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if diff := cmp.Diff(*got, tt.expected); diff != "" {
					t.Errorf("parseXfrmPolicyAddUpdate() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestParseXfrmPolicyDeleteGet(t *testing.T) {
	_, src, err := net.ParseCIDR("192.168.1.1/24")
	if err != nil {
		t.Fatal(err)
	}

	_, dst, err := net.ParseCIDR("192.168.1.2/24")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		tokens   []string
		expected netlink.XfrmPolicy
		wantErr  bool
	}{
		{
			name:   "valid input",
			tokens: []string{"src", "192.168.1.1/24", "dst", "192.168.1.2/24", "proto", "esp", "sport", "1234", "dport", "5678", "dir", "in", "mark", "1", "if_id", "1"},
			expected: netlink.XfrmPolicy{
				Src:     src,
				Dst:     dst,
				Proto:   netlink.XFRM_PROTO_ESP,
				SrcPort: 1234,
				DstPort: 5678,
				Dir:     netlink.XFRM_DIR_IN,
				Mark:    &netlink.XfrmMark{Value: 1},
				Action:  netlink.XFRM_POLICY_ALLOW,
				Ifid:    1,
			},
			wantErr: false,
		},
		{
			name:   "index input",
			tokens: []string{"index", "2"},
			expected: netlink.XfrmPolicy{
				Index: 2,
			},
			wantErr: false,
		},
		{
			name:    "index and sport",
			tokens:  []string{"index", "2", "sport", "2"},
			wantErr: true,
		},
		{
			name:    "invalid arg",
			tokens:  []string{"dsty", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid dst",
			tokens:  []string{"dst", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid src",
			tokens:  []string{"src", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid proto",
			tokens:  []string{"proto", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid sport",
			tokens:  []string{"sport", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid dport",
			tokens:  []string{"dport", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid dir",
			tokens:  []string{"dir", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid mark",
			tokens:  []string{"mark", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid index",
			tokens:  []string{"index", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid if_id",
			tokens:  []string{"if_id", "abc"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := cmd{
				Cursor: -1,
				Args:   tt.tokens,
			}
			got, err := cmd.parseXfrmPolicyDeleteGet()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseXfrmPolicyAddUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if diff := cmp.Diff(*got, tt.expected); diff != "" {
					t.Errorf("parseXfrmPolicyAddUpdate() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestParseXfrmPolicyListDeleteAll(t *testing.T) {
	_, src, err := net.ParseCIDR("192.168.1.1/24")
	if err != nil {
		t.Fatal(err)
	}

	_, dst, err := net.ParseCIDR("192.168.1.2/24")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		tokens   []string
		expected netlink.XfrmPolicy
		wantErr  bool
	}{
		{
			name:   "valid input",
			tokens: []string{"src", "192.168.1.1/24", "dst", "192.168.1.2/24", "proto", "esp", "sport", "1234", "dport", "5678", "dir", "in", "mark", "1", "if_id", "1"},
			expected: netlink.XfrmPolicy{
				Src:     src,
				Dst:     dst,
				Proto:   netlink.XFRM_PROTO_ESP,
				SrcPort: 1234,
				DstPort: 5678,
				Dir:     netlink.XFRM_DIR_IN,
				Mark:    &netlink.XfrmMark{Value: 1},
				Action:  netlink.XFRM_POLICY_ALLOW,
				Ifid:    1,
			},
			wantErr: false,
		},
		{
			name:   "index input",
			tokens: []string{"index", "2"},
			expected: netlink.XfrmPolicy{
				Index: 2,
			},
			wantErr: false,
		},
		{
			name:    "invalid arg",
			tokens:  []string{"dsty", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid dst",
			tokens:  []string{"dst", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid src",
			tokens:  []string{"src", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid proto",
			tokens:  []string{"proto", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid sport",
			tokens:  []string{"sport", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid dport",
			tokens:  []string{"dport", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid dir",
			tokens:  []string{"dir", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid mark",
			tokens:  []string{"mark", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid index",
			tokens:  []string{"index", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid if_id",
			tokens:  []string{"if_id", "abc"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := cmd{
				Cursor: -1,
				Args:   tt.tokens,
			}
			got, err := cmd.parseXfrmPolicyListDeleteAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseXfrmPolicyAddUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if diff := cmp.Diff(*got, tt.expected); diff != "" {
					t.Errorf("parseXfrmPolicyAddUpdate() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestPrintXfrmPolicy(t *testing.T) {
	_, src, err := net.ParseCIDR("192.168.1.1/24")
	if err != nil {
		t.Fatal(err)
	}

	_, dst, err := net.ParseCIDR("192.168.1.2/24")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name     string
		policy   netlink.XfrmPolicy
		expected string
	}{
		{
			name: "valid policy with mark and tmpl",
			policy: netlink.XfrmPolicy{
				Src:      src,
				Dst:      dst,
				Dir:      1,
				Priority: 10,
				Proto:    1,
				SrcPort:  1234,
				DstPort:  5678,
				Action:   1,
				Ifid:     1,
				Mark:     &netlink.XfrmMark{Value: 1, Mask: 0xFFFFFFFF},
				Tmpls: []netlink.XfrmPolicyTmpl{
					{
						Src:   net.ParseIP("192.168.1.1"),
						Dst:   net.ParseIP("192.168.1.2"),
						Proto: 1,
						Reqid: 1,
						Mode:  1,
						Spi:   1,
					},
				},
			},
			expected: "src 192.168.1.0/24 dst 192.168.1.0/24\n\tdir out priority 10\n\tproto 1 sport 1234 dport 5678\n\taction block if_id 1\n\tmark 1/ffffffff\n\ttmpl src 192.168.1.1 dst 192.168.1.2\n\t\tproto 1 reqid 1 mode tunnel spi 1\n",
		},
		{
			name: "policy without mark and tmpl",
			policy: netlink.XfrmPolicy{
				Src:      src,
				Dst:      dst,
				Dir:      1,
				Priority: 10,
				Proto:    1,
				SrcPort:  1234,
				DstPort:  5678,
				Action:   1,
				Ifid:     1,
			},
			expected: "src 192.168.1.0/24 dst 192.168.1.0/24\n\tdir out priority 10\n\tproto 1 sport 1234 dport 5678\n\taction block if_id 1\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			printXfrmPolicy(&buf, tt.policy)
			if diff := cmp.Diff(tt.expected, buf.String()); diff != "" {
				t.Errorf("printXfrmPolicy() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPrintFilteredXfrmPolicies(t *testing.T) {
	_, src, err := net.ParseCIDR("192.168.1.1/24")
	if err != nil {
		t.Fatal(err)
	}

	_, dst, err := net.ParseCIDR("192.168.1.2/24")
	if err != nil {
		t.Fatal(err)
	}
	_, src2, err := net.ParseCIDR("192.168.1.3/24")
	if err != nil {
		t.Fatal(err)
	}

	_, dst2, err := net.ParseCIDR("192.168.1.4/24")
	if err != nil {
		t.Fatal(err)
	}

	policies := []netlink.XfrmPolicy{
		{
			Src:      src,
			Dst:      dst,
			Dir:      1,
			Priority: 10,
			Proto:    1,
			SrcPort:  1234,
			DstPort:  5678,
			Action:   1,
			Ifid:     1,
			Mark:     &netlink.XfrmMark{Value: 1, Mask: 0xFFFFFFFF},
			Tmpls: []netlink.XfrmPolicyTmpl{
				{
					Src:   net.ParseIP("192.168.1.1"),
					Dst:   net.ParseIP("192.168.1.2"),
					Proto: 1,
					Reqid: 1,
					Mode:  1,
					Spi:   1,
				},
			},
		},
		{
			Src:      src2,
			Dst:      dst2,
			Dir:      1,
			Priority: 10,
			Proto:    1,
			SrcPort:  1234,
			DstPort:  5678,
			Action:   1,
			Ifid:     1,
		},
	}

	tests := []struct {
		name     string
		filter   *netlink.XfrmPolicy
		expected string
	}{
		{
			name:     "no filter",
			filter:   nil,
			expected: "src 192.168.1.0/24 dst 192.168.1.0/24\n\tdir out priority 10\n\tproto 1 sport 1234 dport 5678\n\taction block if_id 1\n\tmark 1/ffffffff\n\ttmpl src 192.168.1.1 dst 192.168.1.2\n\t\tproto 1 reqid 1 mode tunnel spi 1\n\nsrc 192.168.1.0/24 dst 192.168.1.0/24\n\tdir out priority 10\n\tproto 1 sport 1234 dport 5678\n\taction block if_id 1\n\n",
		},
		{
			name: "filter by src",
			filter: &netlink.XfrmPolicy{
				Src: src,
			},
			expected: "",
		},
		{
			name: "filter by dst",
			filter: &netlink.XfrmPolicy{
				Dst: dst,
			},
			expected: "",
		},
		{
			name: "filter by proto",
			filter: &netlink.XfrmPolicy{
				Proto: 2,
			},
			expected: "",
		},
		{
			name: "filter by dir",
			filter: &netlink.XfrmPolicy{
				Dir: 2,
			},
			expected: "",
		},
		{
			name: "filter by mark",
			filter: &netlink.XfrmPolicy{
				Mark: &netlink.XfrmMark{Value: 2, Mask: 0xFFFFFFA},
			},
			expected: "",
		},
		{
			name: "filter by index",
			filter: &netlink.XfrmPolicy{
				Index: 2,
			},
			expected: "",
		},
		{
			name: "filter by ifid",
			filter: &netlink.XfrmPolicy{
				Ifid: 2,
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			printFilteredXfrmPolicies(&buf, policies, tt.filter)

			if diff := cmp.Diff(tt.expected, buf.String()); diff != "" {
				t.Errorf("printFilteredXfrmPolicies() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
