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

func TestFilterXfrmStates(t *testing.T) {
	states := []netlink.XfrmState{
		{
			Src:   net.IPv4(192, 168, 1, 1),
			Dst:   net.IPv4(192, 168, 1, 2),
			Proto: 1,
			Spi:   100,
			Mode:  1,
			Reqid: 1,
		},
		{
			Src:   net.IPv4(192, 168, 1, 3),
			Dst:   net.IPv4(192, 168, 1, 4),
			Proto: 2,
			Spi:   200,
			Mode:  2,
			Reqid: 2,
		},
	}

	tests := []struct {
		name     string
		filter   *netlink.XfrmState
		states   []netlink.XfrmState
		noKeys   bool
		expected string
	}{
		{
			name:     "no filter",
			filter:   nil,
			states:   states,
			expected: "src 192.168.1.1 dst 192.168.1.2\n\tproto 1 spi 0x64 mode tunnel\n\treqid 1\nstatistics: replay-window 0 replay 0 failed 0 bytes 0 packets 0\n\nsrc 192.168.1.3 dst 192.168.1.4\n\tproto 2 spi 0xc8 mode ro\n\treqid 2\nstatistics: replay-window 0 replay 0 failed 0 bytes 0 packets 0\n\n",
		},
		{
			name: "no filter all opts",
			states: []netlink.XfrmState{
				{
					Src:          net.IPv4(192, 168, 1, 1),
					Dst:          net.IPv4(192, 168, 1, 2),
					Proto:        1,
					Spi:          100,
					Mode:         1,
					Reqid:        1,
					ReplayWindow: 1,
					Auth: &netlink.XfrmStateAlgo{
						Name: "hmac(sha256)",
						Key:  []byte("key"),
					},
					Crypt: &netlink.XfrmStateAlgo{
						Name: "cbc(aes)",
						Key:  []byte("key"),
					},
					Encap: &netlink.XfrmStateEncap{
						Type:            1,
						SrcPort:         100,
						DstPort:         200,
						OriginalAddress: net.ParseIP("127.0.0.2"),
					},
					Aead: &netlink.XfrmStateAlgo{
						Name: "rfc4106(gcm(aes))",
						Key:  []byte("key"),
					},
					Mark: &netlink.XfrmMark{
						Value: 1,
						Mask:  2,
					},
					OutputMark: &netlink.XfrmMark{
						Value: 3,
						Mask:  4,
					},
					Limits: netlink.XfrmStateLimits{
						ByteSoft:    1,
						ByteHard:    2,
						PacketSoft:  3,
						PacketHard:  4,
						TimeSoft:    5,
						TimeHard:    6,
						TimeUseSoft: 7,
						TimeUseHard: 8,
					},
				},
			},
			expected: "src 192.168.1.1 dst 192.168.1.2\n" +
				"\tproto 1 spi 0x64 mode tunnel\n" +
				"\treqid 1 replay-window 1\n" +
				"\tauth hmac(sha256) 0x6b6579 24bits\n" +
				"\tenc cbc(aes) 0x6b6579 24bits\n" +
				"\taead rfc4106(gcm(aes)) 0x6b6579 24bits\n" +
				"\tencap type espinudp-non-ike sport 100 dport 200 addr 127.0.0.2\n" +
				"\tmark 1/2\n" +
				"\toutput-mark 3/4\n" +
				"\tsoft-byte-limit 1 hard-byte-limit 2\n" +
				"\tsoft-packet-limit 3 hard-packet-limit 4\n" +
				"\tsoft-add-expires-seconds 5 hard-add-expires-seconds 6\n" +
				"\tsoft-use-expires-seconds 7 hard-use-expires-seconds 8\n" +
				"statistics: replay-window 0 replay 0 failed 0 bytes 0 packets 0\n" +
				"\n",
		},
		{
			name:   "no filter all opts no keys",
			noKeys: true,
			states: []netlink.XfrmState{
				{
					Src:          net.IPv4(192, 168, 1, 1),
					Dst:          net.IPv4(192, 168, 1, 2),
					Proto:        1,
					Spi:          100,
					Mode:         1,
					Reqid:        1,
					ReplayWindow: 1,
					Auth: &netlink.XfrmStateAlgo{
						Name: "hmac(sha256)",
						Key:  []byte("key"),
					},
					Crypt: &netlink.XfrmStateAlgo{
						Name: "cbc(aes)",
						Key:  []byte("key"),
					},
					Encap: &netlink.XfrmStateEncap{
						Type:            1,
						SrcPort:         100,
						DstPort:         200,
						OriginalAddress: net.ParseIP("127.0.0.2"),
					},
					Aead: &netlink.XfrmStateAlgo{
						Name: "rfc4106(gcm(aes))",
						Key:  []byte("key"),
					},
					Mark: &netlink.XfrmMark{
						Value: 1,
						Mask:  2,
					},
					OutputMark: &netlink.XfrmMark{
						Value: 3,
						Mask:  4,
					},
					Limits: netlink.XfrmStateLimits{
						ByteSoft:    1,
						ByteHard:    2,
						PacketSoft:  3,
						PacketHard:  4,
						TimeSoft:    5,
						TimeHard:    6,
						TimeUseSoft: 7,
						TimeUseHard: 8,
					},
				},
			},
			expected: "src 192.168.1.1 dst 192.168.1.2\n" +
				"\tproto 1 spi 0x64 mode tunnel\n" +
				"\treqid 1 replay-window 1\n" +
				"\tauth hmac(sha256) 24bits\n" +
				"\tenc cbc(aes) 24bits\n" +
				"\taead rfc4106(gcm(aes)) 24bits\n" +
				"\tencap type espinudp-non-ike sport 100 dport 200 addr 127.0.0.2\n" +
				"\tmark 1/2\n" +
				"\toutput-mark 3/4\n" +
				"\tsoft-byte-limit 1 hard-byte-limit 2\n" +
				"\tsoft-packet-limit 3 hard-packet-limit 4\n" +
				"\tsoft-add-expires-seconds 5 hard-add-expires-seconds 6\n" +
				"\tsoft-use-expires-seconds 7 hard-use-expires-seconds 8\n" +
				"statistics: replay-window 0 replay 0 failed 0 bytes 0 packets 0\n" +
				"\n",
		},
		{
			name:   "filter by src",
			states: states,
			filter: &netlink.XfrmState{
				Src: net.IPv4(192, 168, 1, 1),
			},
			expected: "src 192.168.1.1 dst 192.168.1.2\n\tproto 1 spi 0x64 mode tunnel\n\treqid 1\nstatistics: replay-window 0 replay 0 failed 0 bytes 0 packets 0\n\n",
		},
		{
			name:   "filter by dst",
			states: states,
			filter: &netlink.XfrmState{
				Dst: net.IPv4(192, 168, 1, 4),
			},
			expected: "src 192.168.1.3 dst 192.168.1.4\n\tproto 2 spi 0xc8 mode ro\n\treqid 2\nstatistics: replay-window 0 replay 0 failed 0 bytes 0 packets 0\n\n",
		},
		{
			name:   "filter by proto",
			states: states,
			filter: &netlink.XfrmState{
				Proto: 1,
			},
			expected: "src 192.168.1.1 dst 192.168.1.2\n\tproto 1 spi 0x64 mode tunnel\n\treqid 1\nstatistics: replay-window 0 replay 0 failed 0 bytes 0 packets 0\n\n",
		},
		{
			name:   "filter by spi",
			states: states,
			filter: &netlink.XfrmState{
				Spi: 200,
			},
			expected: "src 192.168.1.3 dst 192.168.1.4\n\tproto 2 spi 0xc8 mode ro\n\treqid 2\nstatistics: replay-window 0 replay 0 failed 0 bytes 0 packets 0\n\n",
		},
		{
			name:   "filter by mode",
			states: states,
			filter: &netlink.XfrmState{
				Mode: 1,
			},
			expected: "src 192.168.1.1 dst 192.168.1.2\n\tproto 1 spi 0x64 mode tunnel\n\treqid 1\nstatistics: replay-window 0 replay 0 failed 0 bytes 0 packets 0\n\n",
		},
		{
			name:   "filter by reqid",
			states: states,
			filter: &netlink.XfrmState{
				Reqid: 2,
			},
			expected: "src 192.168.1.3 dst 192.168.1.4\n\tproto 2 spi 0xc8 mode ro\n\treqid 2\nstatistics: replay-window 0 replay 0 failed 0 bytes 0 packets 0\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			cmd := cmd{
				Out: &buf,
			}

			cmd.printFilteredXfrmStates(tt.states, tt.filter, tt.noKeys)

			if diff := cmp.Diff(buf.String(), tt.expected); diff != "" {
				t.Errorf("filterXfrmStates() diff:\n%v", diff)
			}
		})
	}
}

func TestParseXfrmStateAddUpdate(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []string
		expected netlink.XfrmState
		wantErr  bool
	}{
		{
			name:   "valid input",
			tokens: []string{"src", "192.168.1.1", "dst", "192.168.1.2", "proto", "esp", "spi", "100", "mode", "tunnel", "enc", "cbc(aes)", "ff", "auth", "hmac(sha256)", "ff", "auth-trunc", "hmac(sha256)", "ff", "2", "aead", "rfc4106(gcm(aes)", "ff", "3", "mark", "1", "reqid", "1", "replay-window", "1", "limit", "time-soft", "100", "encap", "espinudp", "1", "2", "127.0.0.2", "output-mark", "3", "if_id", "1"},
			expected: netlink.XfrmState{
				Src:   net.IPv4(192, 168, 1, 1),
				Dst:   net.IPv4(192, 168, 1, 2),
				Proto: netlink.XFRM_PROTO_ESP,
				Spi:   100,
				Mode:  netlink.XFRM_MODE_TUNNEL,
				Auth: &netlink.XfrmStateAlgo{
					Name:        "hmac(sha256)",
					Key:         []byte{0xff},
					TruncateLen: 2,
				},
				Crypt: &netlink.XfrmStateAlgo{
					Name: "cbc(aes)",
					Key:  []byte{0xff},
				},
				Aead: &netlink.XfrmStateAlgo{
					Name:   "rfc4106(gcm(aes)",
					Key:    []byte{0xff},
					ICVLen: 3,
				},
				Mark: &netlink.XfrmMark{
					Value: 1,
				},
				Reqid:        1,
				ReplayWindow: 1,
				Limits: netlink.XfrmStateLimits{
					TimeSoft: 100,
				},
				Encap: &netlink.XfrmStateEncap{
					Type:            netlink.XFRM_ENCAP_ESPINUDP,
					SrcPort:         1,
					DstPort:         2,
					OriginalAddress: net.ParseIP("127.0.0.2"),
				},
				OutputMark: &netlink.XfrmMark{
					Value: 3,
				},
				Ifid: 1,
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
			name:    "invalid spi",
			tokens:  []string{"spi", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid enc",
			tokens:  []string{"enc", "avc", "y"},
			wantErr: true,
		},
		{
			name:    "invalid mode",
			tokens:  []string{"mode", "y"},
			wantErr: true,
		},
		{
			name:    "invalid auth",
			tokens:  []string{"auth", "abc", "y"},
			wantErr: true,
		},
		{
			name:    "invalid auth-trunc",
			tokens:  []string{"auth-trunc", "abc", "y"},
			wantErr: true,
		},
		{
			name:    "invalid auth-trunc len",
			tokens:  []string{"auth-trunc", "abc", "ff", "y"},
			wantErr: true,
		},
		{
			name:    "invalid aead",
			tokens:  []string{"aead", "abc", "y"},
			wantErr: true,
		},
		{
			name:    "invalid aead len",
			tokens:  []string{"aead", "abc", "ff", "y"},
			wantErr: true,
		},
		{
			name:    "invalid comp",
			tokens:  []string{"comp", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid mark",
			tokens:  []string{"mark", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid reqid",
			tokens:  []string{"reqid", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid replay-window",
			tokens:  []string{"replay-window", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid limit",
			tokens:  []string{"limit", "time-soft", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid limit",
			tokens:  []string{"limit", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid encap",
			tokens:  []string{"encap", "espintcp"},
			wantErr: true,
		},
		{
			name:    "invalid output-mark",
			tokens:  []string{"output-mark", "abc"},
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
			got, err := cmd.parseXfrmStateAddUpdate()
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

func TestParseXfrmStateAllocSPI(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []string
		expected netlink.XfrmState
		wantErr  bool
	}{
		{
			name:   "valid input",
			tokens: []string{"src", "192.168.1.1", "dst", "192.168.1.2", "proto", "esp", "spi", "100", "mode", "tunnel", "mark", "1", "reqid", "1"},
			expected: netlink.XfrmState{
				Src:   net.IPv4(192, 168, 1, 1),
				Dst:   net.IPv4(192, 168, 1, 2),
				Proto: netlink.XFRM_PROTO_ESP,
				Spi:   100,
				Mode:  netlink.XFRM_MODE_TUNNEL,
				Mark: &netlink.XfrmMark{
					Value: 1,
				},
				Reqid: 1,
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
			name:    "invalid spi",
			tokens:  []string{"spi", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid mode",
			tokens:  []string{"mode", "y"},
			wantErr: true,
		},
		{
			name:    "invalid mark",
			tokens:  []string{"mark", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid reqid",
			tokens:  []string{"reqid", "abc"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := cmd{
				Cursor: -1,
				Args:   tt.tokens,
			}
			got, err := cmd.parseXfrmStateAllocSPI()
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

func TestParseXfrmStateDeleteGet(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []string
		expected netlink.XfrmState
		wantErr  bool
	}{
		{
			name:   "valid input",
			tokens: []string{"src", "192.168.1.1", "dst", "192.168.1.2", "proto", "esp", "spi", "100", "mark", "1"},
			expected: netlink.XfrmState{
				Src:   net.IPv4(192, 168, 1, 1),
				Dst:   net.IPv4(192, 168, 1, 2),
				Proto: netlink.XFRM_PROTO_ESP,
				Spi:   100,
				Mark: &netlink.XfrmMark{
					Value: 1,
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
			name:    "invalid spi",
			tokens:  []string{"spi", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid mode",
			tokens:  []string{"mode", "y"},
			wantErr: true,
		},
		{
			name:    "invalid mark",
			tokens:  []string{"mark", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid reqid",
			tokens:  []string{"reqid", "abc"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := cmd{
				Cursor: -1,
				Args:   tt.tokens,
			}
			got, err := cmd.parseXfrmStateDeleteGet()
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

func TestParseXfrmStateListDeleteAll(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []string
		expected netlink.XfrmState
		noKeys   bool
		wantErr  bool
	}{
		{
			name:   "valid input",
			tokens: []string{"src", "192.168.1.1", "dst", "192.168.1.2", "proto", "esp", "spi", "100", "mode", "tunnel", "reqid", "1"},
			expected: netlink.XfrmState{
				Src:   net.IPv4(192, 168, 1, 1),
				Dst:   net.IPv4(192, 168, 1, 2),
				Proto: netlink.XFRM_PROTO_ESP,
				Spi:   100,
				Mode:  netlink.XFRM_MODE_TUNNEL,
				Reqid: 1,
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
			name:    "invalid spi",
			tokens:  []string{"spi", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid mode",
			tokens:  []string{"mode", "y"},
			wantErr: true,
		},
		{
			name:    "invalid mark",
			tokens:  []string{"mark", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid reqid",
			tokens:  []string{"reqid", "abc"},
			wantErr: true,
		},
		{
			name:   "no keys",
			tokens: []string{"nokeys"},
			noKeys: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := cmd{
				Cursor: -1,
				Args:   tt.tokens,
			}
			got, noKeys, err := cmd.parseXfrmStateListDeleteAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseXfrmPolicyAddUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if diff := cmp.Diff(*got, tt.expected); diff != "" {
					t.Errorf("parseXfrmPolicyAddUpdate() mismatch (-want +got):\n%s", diff)
				}

				if noKeys != tt.noKeys {
					t.Errorf("parseXfrmPolicyAddUpdate() noKeys = %v, want %v", noKeys, tt.noKeys)
				}
			}
		})
	}
}
