// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"bytes"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/vishvananda/netlink"
)

func TestParseXfrmProto(t *testing.T) {
	tests := []struct {
		name      string
		tokens    []string
		wantProto netlink.Proto
		wantErr   bool
	}{
		{"Valid ESP", []string{"esp"}, netlink.XFRM_PROTO_ESP, false},
		{"Valid AH", []string{"ah"}, netlink.XFRM_PROTO_AH, false},
		{"Valid COMP", []string{"comp"}, netlink.XFRM_PROTO_COMP, false},
		{"Valid ROUTE2", []string{"route2"}, netlink.XFRM_PROTO_ROUTE2, false},
		{"Valid HAO", []string{"hao"}, netlink.XFRM_PROTO_HAO, false},
		{"Invalid Token", []string{"invalid"}, netlink.XFRM_PROTO_IPSEC_ANY, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cmd{Cursor: -1, Args: tt.tokens}
			gotProto, err := cmd.parseXfrmProto()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseXfrmProto() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotProto != tt.wantProto {
				t.Errorf("parseXfrmProto() = %v, want %v", gotProto, tt.wantProto)
			}
		})
	}
}

func TestParseXfrmMode(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []string
		wantMode netlink.Mode
		wantErr  bool
	}{
		{"Valid Transport", []string{"transport"}, netlink.XFRM_MODE_TRANSPORT, false},
		{"Valid Tunnel", []string{"tunnel"}, netlink.XFRM_MODE_TUNNEL, false},
		{"Valid RO", []string{"ro"}, netlink.XFRM_MODE_ROUTEOPTIMIZATION, false},
		{"Valid In_Trigger", []string{"in_trigger"}, netlink.XFRM_MODE_IN_TRIGGER, false},
		{"Valid Beet", []string{"beet"}, netlink.XFRM_MODE_BEET, false},
		{"Invalid Token", []string{"invalid"}, netlink.XFRM_MODE_MAX, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cmd{Cursor: -1, Args: tt.tokens}
			gotMode, err := cmd.parseXfrmMode()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseXfrmMode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotMode != tt.wantMode {
				t.Errorf("parseXfrmMode() = %v, want %v", gotMode, tt.wantMode)
			}
		})
	}
}

func TestParseXfrmDir(t *testing.T) {
	tests := []struct {
		name    string
		tokens  []string
		wantDir netlink.Dir
		wantErr bool
	}{
		{"Valid In", []string{"in"}, netlink.XFRM_DIR_IN, false},
		{"Valid Out", []string{"out"}, netlink.XFRM_DIR_OUT, false},
		{"Valid Fwd", []string{"fwd"}, netlink.XFRM_DIR_FWD, false},
		{"Invalid Token", []string{"invalid"}, netlink.XFRM_DIR_IN, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cmd{Cursor: -1, Args: tt.tokens}
			gotDir, err := cmd.parseXfrmDir()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseXfrmDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDir != tt.wantDir {
				t.Errorf("parseXfrmDir() = %v, want %v", gotDir, tt.wantDir)
			}
		})
	}
}

func TestParseXfrmAction(t *testing.T) {
	tests := []struct {
		name       string
		tokens     []string
		wantAction netlink.PolicyAction
		wantErr    bool
	}{
		{"Valid Allow", []string{"allow"}, netlink.XFRM_POLICY_ALLOW, false},
		{"Valid Block", []string{"block"}, netlink.XFRM_POLICY_BLOCK, false},
		{"Invalid Token", []string{"invalid"}, netlink.XFRM_POLICY_ALLOW, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cmd{Cursor: -1, Args: tt.tokens}
			gotAction, err := cmd.parseXfrmAction()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseXfrmAction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAction != tt.wantAction {
				t.Errorf("parseXfrmAction() = %v, want %v", gotAction, tt.wantAction)
			}
		})
	}
}

func TestParseXfrmMark(t *testing.T) {
	tests := []struct {
		name    string
		tokens  []string
		want    *netlink.XfrmMark
		wantErr bool
	}{
		{"Valid Mark Only", []string{"1234"}, &netlink.XfrmMark{Value: 1234}, false},
		{"Valid Mark next arg", []string{"1234", "abc"}, &netlink.XfrmMark{Value: 1234}, false},

		{"Valid Mark and Mask", []string{"1234", "mask", "5678"}, &netlink.XfrmMark{Value: 1234, Mask: 5678}, false},
		{"Invalid Mark", []string{"abc"}, nil, true},
		{"Invalid Mask", []string{"1234", "mask", "abc"}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cmd{Cursor: -1, Args: tt.tokens}
			got, err := cmd.parseXfrmMark()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseXfrmMark() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && (got.Value != tt.want.Value || got.Mask != tt.want.Mask) {
				t.Errorf("parseXfrmMark() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseXfrmEncap(t *testing.T) {
	tests := []struct {
		name    string
		tokens  []string
		want    *netlink.XfrmStateEncap
		wantErr bool
	}{
		{
			name:   "Valid espinudp",
			tokens: []string{"espinudp", "1234", "5678", "192.168.1.1"},
			want: &netlink.XfrmStateEncap{
				Type:            netlink.XFRM_ENCAP_ESPINUDP,
				SrcPort:         1234,
				DstPort:         5678,
				OriginalAddress: net.ParseIP("192.168.1.1"),
			},
			wantErr: false,
		},
		{
			name:   "Valid espinudp-nonike",
			tokens: []string{"espinudp-nonike", "1234", "5678", "192.168.1.1"},
			want: &netlink.XfrmStateEncap{
				Type:            netlink.XFRM_ENCAP_ESPINUDP_NONIKE,
				SrcPort:         1234,
				DstPort:         5678,
				OriginalAddress: net.ParseIP("192.168.1.1"),
			},
			wantErr: false,
		},
		{
			name:    "Unsupported espintcp",
			tokens:  []string{"espintcp"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid SrcPort",
			tokens:  []string{"espinudp", "abc"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid DstPort",
			tokens:  []string{"espinudp", "1234", "abc"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid OriginalAddress",
			tokens:  []string{"espinudp", "1234", "5678", "abc"},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cmd{Cursor: -1, Args: tt.tokens}
			got, err := cmd.parseXfrmEncap()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseXfrmEncap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Errorf("parseXfrmEncap() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestParseXfrmLimit(t *testing.T) {
	tests := []struct {
		name    string
		tokens  []string
		want    netlink.XfrmStateLimits
		wantErr bool
	}{
		{
			name:   "all opts",
			tokens: []string{"time-soft", "60", "time-hard", "120", "byte-soft", "1024", "byte-hard", "2048", "packet-soft", "10", "packet-hard", "20", "time-use-soft", "30", "time-use-hard", "40"},
			want: netlink.XfrmStateLimits{
				TimeSoft:    60,
				TimeHard:    120,
				ByteSoft:    1024,
				ByteHard:    2048,
				PacketSoft:  10,
				PacketHard:  20,
				TimeUseSoft: 30,
				TimeUseHard: 40,
			},
			wantErr: false,
		},
		{
			name:    "time-soft invalid",
			tokens:  []string{"time-soft", "abc"},
			wantErr: true,
		},
		{
			name:    "time-hard invalid",
			tokens:  []string{"time-hard", "abc"},
			wantErr: true,
		},
		{
			name:    "time-use-soft invalid",
			tokens:  []string{"time-use-soft", "abc"},
			wantErr: true,
		},
		{
			name:    "time-use-hard invalid",
			tokens:  []string{"time-use-hard", "abc"},
			wantErr: true,
		},
		{
			name:    "byte-soft invalid",
			tokens:  []string{"byte-soft", "abc"},
			wantErr: true,
		},
		{
			name:    "byte-hard invalid",
			tokens:  []string{"byte-hard", "abc"},
			wantErr: true,
		},
		{
			name:    "packet-soft invalid",
			tokens:  []string{"packet-soft", "abc"},
			wantErr: true,
		},
		{
			name:    "packet-hard invalid",
			tokens:  []string{"packet-hard", "abc"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cmd{Cursor: -1, Args: tt.tokens}
			got, err := cmd.parseXfrmLimit()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseXfrmLimit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseXfrmLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrintXfrmMsgExpire(t *testing.T) {
	// Create a mock netlink.XfrmState object with test data
	msg := &netlink.XfrmState{
		Src:          net.ParseIP("192.168.1.1"),
		Dst:          net.ParseIP("192.168.1.2"),
		Proto:        1,
		Spi:          12345,
		Reqid:        67890,
		Mode:         1,
		ReplayWindow: 32,
		Auth: &netlink.XfrmStateAlgo{
			Name:        "hmac(sha256)",
			Key:         []byte("1234567890abcdef"),
			TruncateLen: 128,
		},
		Crypt: &netlink.XfrmStateAlgo{
			Name: "cbc(aes)",
			Key:  []byte("abcdef1234567890"),
		},
		Limits: netlink.XfrmStateLimits{
			ByteSoft:    1000,
			ByteHard:    2000,
			PacketSoft:  100,
			PacketHard:  200,
			TimeSoft:    3600,
			TimeHard:    7200,
			TimeUseSoft: 1800,
			TimeUseHard: 3600,
		},
		Statistics: netlink.XfrmStateStats{
			Bytes:        500,
			Packets:      50,
			AddTime:      100,
			UseTime:      200,
			ReplayWindow: 32,
			Replay:       10,
			Failed:       1,
		},
	}

	// Capture the output
	var buf bytes.Buffer
	printXfrmMsgExpire(&buf, msg)

	// Define the expected output
	expectedOutput := `src 192.168.1.1 dst 192.168.1.2
    proto 1 spi 12345 reqid 67890 mode tunnel
    replay-window 32
    auth-trunc hmac(sha256) 1234567890abcdef 128
    enc cbc(aes) abcdef1234567890
    sel src 192.168.1.1 dst 192.168.1.2
    lifetime config:
      limit: soft (1000)(bytes), hard (2000)(bytes)
      limit: soft (100)(packets), hard (200)(packets)
      expire add: soft 3600(sec), hard 7200(sec)
      expire use: soft 1800(sec), hard 3600(sec)
    lifetime current:
      500(bytes), 50(packets)
      add 100, use 200
    stats:
      replay-window 32 replay 10 failed 1
`

	// Compare the output
	if cmp.Diff(buf.String(), expectedOutput) != "" {
		t.Errorf("printXfrmMsgExpire() mismatch (-want +got):\n%s", cmp.Diff(buf.String(), expectedOutput))
	}
}
