// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/layers"
)

func TestParseFlags(t *testing.T) {
	tmpDir := t.TempDir()
	tmp, err := os.CreateTemp(tmpDir, "filter.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	tmp.Write([]byte("tcp port 80"))

	tests := []struct {
		name        string
		args        []string
		expectedCmd cmd
		expectedErr bool
	}{
		{
			name: "default with count limit",
			args: []string{"cmd", "-c", "10"},
			expectedCmd: cmd{
				Opts: flags{
					CountPkg:       10,
					SnapshotLength: 262144,
				},
			},
			expectedErr: false,
		},
		{
			name: "help",
			args: []string{"cmd", "-h"},
			expectedCmd: cmd{
				Opts: flags{
					SnapshotLength: 262144,
					Help:           true,
				},
			},
			expectedErr: false,
		},
		{
			name:        "verbose and quiet",
			args:        []string{"cmd", "-v", "-q"},
			expectedErr: true,
		},
		{
			name: "filter file",
			args: []string{"cmd", "-F", tmp.Name()},
			expectedCmd: cmd{
				Opts: flags{
					SnapshotLength: 262144,
					Filter:         "tcp port 80",
					FilterFile:     tmp.Name(),
				},
			},
		},
		{
			name:        "missing filter file",
			args:        []string{"cmd", "-F", "xyz"},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			cmd, err := parseFlags(tt.args, out)
			if (err != nil) != tt.expectedErr {
				t.Errorf("parseFlags() error = %v, expectedErr %v", err, tt.expectedErr)
			}

			if !tt.expectedErr {
				if diff := cmp.Diff(tt.expectedCmd, cmd, cmpopts.IgnoreFields(cmd, "Out", "usage")); diff != "" {
					t.Errorf("parseFlags() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestProcessPacket(t *testing.T) {
	tests := []struct {
		name             string
		packetData       []byte
		opts             flags
		num              int
		lastPkgTimeStamp time.Time
		expectedOutput   string
	}{
		{
			name: "IPv4 TCP packet",
			packetData: []byte{
				// Ethernet header
				0x00, 0x1c, 0x42, 0x00, 0x00, 0x08, 0x00, 0x1c, 0x42, 0x00, 0x00, 0x01, 0x08, 0x00,
				// IPv4 header
				0x45, 0x00, 0x00, 0x3c, 0x1c, 0x46, 0x40, 0x00, 0x40, 0x06, 0xb1, 0xe6, 0xc0, 0xa8, 0x00, 0x68, 0xc0, 0xa8, 0x00, 0x01,
				// TCP header
				0x00, 0x50, 0x00, 0x50, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x50, 0x02, 0x20, 0x00, 0x91, 0x7c, 0x00, 0x00,
			},
			num:              1,
			lastPkgTimeStamp: time.Now(),
			opts:             flags{Number: true, Verbose: true, Device: "eth0", Numerical: true},
			expectedOutput:   "1  00:00:00.000000 IPv4 eth0 (tos 0x0, ttl 64, id 7238, offset 0, flags [DF], proto TCP (6), length 40)\n 192.168.0.104.80 > 192.168.0.1.80: Flags [S], cksum 0x917c, seq 0, ack 0, win 8192, options [], length 0\n",
		},
		{
			name: "IPv4 TCP packet with quiet flag",
			packetData: []byte{
				// Ethernet header
				0x00, 0x1c, 0x42, 0x00, 0x00, 0x08, 0x00, 0x1c, 0x42, 0x00, 0x00, 0x01, 0x08, 0x00,
				// IPv4 header
				0x45, 0x00, 0x00, 0x3c, 0x1c, 0x46, 0x40, 0x00, 0x40, 0x06, 0xb1, 0xe6, 0xc0, 0xa8, 0x00, 0x68, 0xc0, 0xa8, 0x00, 0x01,
				// TCP header
				0x00, 0x50, 0x00, 0x50, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x50, 0x02, 0x20, 0x00, 0x91, 0x7c, 0x00, 0x00,
			},
			num:              1,
			lastPkgTimeStamp: time.Now(),
			opts:             flags{Quiet: true, Device: "eth0", Numerical: true},
			expectedOutput:   "00:00:00.000000 IPv4 eth0 192.168.0.104.80 > 192.168.0.1.80: TCP, length 0\n",
		},
		{
			name: "IPv4 TCP packet with ASCII flag",
			packetData: []byte{
				// Ethernet header
				0x00, 0x1c, 0x42, 0x00, 0x00, 0x08, 0x00, 0x1c, 0x42, 0x00, 0x00, 0x01, 0x08, 0x00,
				// IPv4 header
				0x45, 0x00, 0x00, 0x3c, 0x1c, 0x46, 0x40, 0x00, 0x40, 0x06, 0xb1, 0xe6, 0xc0, 0xa8, 0x00, 0x68, 0xc0, 0xa8, 0x00, 0x01,
				// TCP header
				0x00, 0x50, 0x00, 0x50, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x50, 0x02, 0x20, 0x00, 0x91, 0x7c, 0x00, 0x00,
			},
			num:              1,
			lastPkgTimeStamp: time.Now(),
			opts:             flags{ASCII: true, Device: "eth0", Numerical: true},
			expectedOutput:   "00:00:00.000000 IPv4 eth0 192.168.0.104.80 > 192.168.0.1.80: Flags [S], seq 0, ack 0, win 8192, options [], length 0\n\n",
		},
		{
			name: "IPv4 UDP packet",
			packetData: []byte{
				// Ethernet header
				0x00, 0x1c, 0x42, 0x00, 0x00, 0x08, 0x00, 0x1c, 0x42, 0x00, 0x00, 0x01, 0x08, 0x00,
				// IPv4 header
				0x45, 0x00, 0x00, 0x1c, 0x1c, 0x46, 0x40, 0x00, 0x40, 0x11, 0xb1, 0xe6, 0xc0, 0xa8, 0x00, 0x68, 0xc0, 0xa8, 0x00, 0x01,
				// UDP header
				0x00, 0x35, 0x00, 0x35, 0x00, 0x08, 0x91, 0x7c,
			},
			num:              2,
			lastPkgTimeStamp: time.Now(),
			opts:             flags{Number: true, Verbose: true, Device: "eth0", Numerical: true},
			expectedOutput:   "2  00:00:00.000000 IPv4 eth0 (tos 0x0, ttl 64, id 7238, offset 0, flags [DF], proto UDP (17), length 28)\n 192.168.0.104.53 > 192.168.0.1.53: UDP, length 0\n",
		},
		{
			name: "IPv4 UDPLite packet",
			packetData: []byte{
				// Ethernet header
				0x00, 0x1c, 0x42, 0x00, 0x00, 0x08, 0x00, 0x1c, 0x42, 0x00, 0x00, 0x01, 0x08, 0x00,
				// IPv4 header
				0x45, 0x00, 0x00, 0x1c, 0x1c, 0x46, 0x40, 0x00, 0x40, 0x88, 0xb1, 0xe6, 0xc0, 0xa8, 0x00, 0x68, 0xc0, 0xa8, 0x00, 0x01,
				// UDPLite header
				0x00, 0x35, 0x00, 0x35, 0x00, 0x08, 0x91, 0x7c,
			},
			num:              3,
			lastPkgTimeStamp: time.Now(),
			opts:             flags{Number: true, Device: "eth0", Numerical: true},
			expectedOutput:   "3  00:00:00.000000 IPv4 eth0 192.168.0.104.53 > 192.168.0.1.53: UDPLite, length 0\n",
		},
		{
			name: "IPv4 DNS packet and data with header flag",
			packetData: []byte{
				// Ethernet header
				0x00, 0x1c, 0x42, 0x00, 0x00, 0x08, 0x00, 0x1c, 0x42, 0x00, 0x00, 0x01, 0x08, 0x00,
				// IPv4 header
				0x45, 0x00, 0x00, 0x3c, 0x1c, 0x46, 0x40, 0x00, 0x40, 0x11, 0xb1, 0xe6, 0xc0, 0xa8, 0x00, 0x68, 0xc0, 0xa8, 0x00, 0x01,
				// UDP header
				0x00, 0x35, 0x00, 0x35, 0x00, 0x28, 0x91, 0x7c,
				// DNS header
				0x12, 0x34, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x77, 0x77, 0x77, 0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03, 0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01,
			},
			num:              4,
			lastPkgTimeStamp: time.Now(),
			opts:             flags{DataWithHeader: true, Device: "eth0", Numerical: true},
			expectedOutput:   "00:00:00.000000 IPv4 eth0 192.168.0.104.53 > 192.168.0.1.53: 4660+ A? www.google.com (32)\n0x0000:  001c 4200 0008 001c 4200 0001 0800 4500 \n0x0010:  003c 1c46 4000 4011 b1e6 c0a8 0068 c0a8 \n0x0020:  0001 0035 0035 0028 917c 1234 0100 0001 \n0x0030:  0000 0000 0000 0377 7777 0667 6f6f 676c \n0x0040:  6503 636f 6d00 0001 0001                \n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packet := gopacket.NewPacket(tt.packetData, layers.LayerTypeEthernet, gopacket.Default)

			cmd := &cmd{
				Out:  new(bytes.Buffer),
				Opts: tt.opts,
			}

			cmd.processPacket(packet, tt.num, tt.lastPkgTimeStamp)

			output := cmd.Out.(*bytes.Buffer).String()
			if diff := cmp.Diff(tt.expectedOutput, output); diff != "" {
				t.Errorf("processPacket() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
