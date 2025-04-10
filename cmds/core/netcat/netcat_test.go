// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/u-root/u-root/pkg/netcat"
	"github.com/u-root/u-root/pkg/ulog"
)

func TestEvalParams(t *testing.T) {
	baseFlags := flags{
		sourcePort:     netcat.DEFAULT_SOURCE_PORT,
		maxConnections: netcat.DEFAULT_CONNECTION_MAX,
		timingTimeout:  "0ms",
		timingDelay:    "0ms",
		timingWait:     "10s",
	}

	hostSet := netcat.DefaultConfig()
	hostSet.Host = "testhost"

	hostSetListen := netcat.DefaultConfig()
	hostSetListen.Host = "testhost"
	hostSetListen.ConnectionMode = netcat.CONNECTION_MODE_LISTEN

	portScan := netcat.DefaultConfig()
	portScan.Host = "testhost"
	portScan.Port = 1234
	portScan.ConnectionModeOptions.ScanPorts = true
	portScan.ConnectionModeOptions.CurrentPort = 1234
	portScan.ConnectionModeOptions.EndPort = 1345

	portSet := netcat.DefaultConfig()
	portSet.Host = "testhost"
	portSet.Port = 1234

	onlyPortSet := netcat.DefaultConfig()
	onlyPortSet.ConnectionMode = netcat.CONNECTION_MODE_LISTEN
	onlyPortSet.Port = 1234

	portSetListen := netcat.DefaultConfig()
	portSetListen.ConnectionMode = netcat.CONNECTION_MODE_LISTEN
	portSetListen.Host = "testhost"
	portSetListen.Port = 1234

	sourcePortSetListen := netcat.DefaultConfig()
	sourcePortSetListen.ConnectionMode = netcat.CONNECTION_MODE_LISTEN
	sourcePortSetListen.Host = "testhost"
	sourcePortSetListen.Port = 3333
	sourcePortSetListen.ConnectionModeOptions.SourcePort = "3333"

	ipv4Set := netcat.DefaultConfig()
	ipv4Set.Host = "testhost"
	ipv4Set.ProtocolOptions.IPType = netcat.IP_V4_STRICT

	ipv6Set := netcat.DefaultConfig()
	ipv6Set.Host = "testhost"
	ipv6Set.ProtocolOptions.IPType = netcat.IP_V6_STRICT

	execNativeSet := netcat.DefaultConfig()
	execNativeSet.Host = "testhost"
	execNativeSet.CommandExec.Type = netcat.EXEC_TYPE_NATIVE
	execNativeSet.CommandExec.Command = "testcommand"

	sourcePortSet := netcat.DefaultConfig()
	sourcePortSet.Host = "testhost"
	sourcePortSet.ConnectionModeOptions.SourcePort = "123"

	verboseSet := netcat.DefaultConfig()
	verboseSet.Host = "testhost"
	verboseSet.Output.Logger = ulog.Log

	listenMode := netcat.DefaultConfig()
	listenMode.ConnectionMode = netcat.CONNECTION_MODE_LISTEN

	setTimings := netcat.DefaultConfig()
	setTimings.Host = "testhost"
	setTimings.Timing.Wait = 10 * time.Second
	setTimings.Timing.Timeout = 20 * time.Second
	setTimings.Timing.Delay = 30 * time.Second

	sslConfig := netcat.DefaultConfig()
	sslConfig.Host = "testhost"
	sslConfig.SSLConfig.CertFilePath = "cert.pem"
	sslConfig.SSLConfig.KeyFilePath = "key.pem"

	clrfConfig := netcat.DefaultConfig()
	clrfConfig.Host = "testhost"
	clrfConfig.Misc.EOL = netcat.LINE_FEED_CRLF

	chatModeConfig := netcat.DefaultConfig()
	chatModeConfig.ConnectionMode = netcat.CONNECTION_MODE_LISTEN
	chatModeConfig.ListenModeOptions.ChatMode = true
	chatModeConfig.ListenModeOptions.BrokerMode = true

	// Define test cases
	tests := []struct {
		name       string
		args       []string
		modify     func(flags) flags
		wantConfig *netcat.Config
		wantErr    bool
	}{
		{
			name: "host set",
			args: []string{"testhost"},
			modify: func(f flags) flags {
				return f
			},
			wantConfig: &hostSet,
			wantErr:    false,
		},
		{
			name: "host set in listen mode",
			args: []string{"testhost"},
			modify: func(f flags) flags {
				f.listen = true
				return f
			},
			wantConfig: &hostSetListen,
			wantErr:    false,
		},
		{
			name:       "port set",
			args:       []string{"testhost", "1234"},
			modify:     func(f flags) flags { return f },
			wantConfig: &portSet,
			wantErr:    false,
		},
		{
			name: "port set in listen mode",
			args: []string{"testhost", "1234"},
			modify: func(f flags) flags {
				f.listen = true
				f.sourcePort = ""
				return f
			},
			wantConfig: &portSetListen,
			wantErr:    false,
		},
		{
			name: "source port set in listen mode",
			args: []string{"testhost"},
			modify: func(f flags) flags {
				f.listen = true
				f.sourcePort = "3333"
				return f
			},
			wantConfig: &sourcePortSetListen,
			wantErr:    false,
		},
		{
			name: "port in listen mode",
			args: []string{"1234"},
			modify: func(f flags) flags {
				f.listen = true
				return f
			},
			wantConfig: &onlyPortSet,
		},
		{
			name: "invalid source port set in listen mode",
			args: []string{"testhost"},
			modify: func(f flags) flags {
				f.listen = true
				f.sourcePort = "abc"
				return f
			},
			wantErr: true,
		},
		{
			name: "port set invalid",
			args: []string{"testhost", "1234-"},
			modify: func(f flags) flags {
				f.listen = true
				return f
			},
			wantConfig: &hostSet,
			wantErr:    true,
		},
		{
			name: "port set invalid connection mode",
			args: []string{"testhost", "1234--"},
			modify: func(f flags) flags {
				return f
			},
			wantErr: true,
		},
		{
			name: "ipv4 set",
			args: []string{"testhost"},
			modify: func(f flags) flags {
				f.ipv4 = true
				return f
			},
			wantConfig: &ipv4Set,
			wantErr:    false,
		},
		{
			name: "ipv6 set",
			args: []string{"testhost"},
			modify: func(f flags) flags {
				f.ipv6 = true
				return f
			},
			wantConfig: &ipv6Set,
			wantErr:    false,
		},
		{
			name: "ipv4 & ipv6 set",
			args: []string{"testhost"},
			modify: func(f flags) flags {
				f.ipv4 = true
				f.ipv6 = true
				return f
			},
			wantErr: true,
		},
		{
			name: "exec native",
			args: []string{"testhost"},
			modify: func(f flags) flags {
				f.execNative = "testcommand"
				return f
			},
			wantConfig: &execNativeSet,
			wantErr:    false,
		},
		{
			name: "source port set",
			args: []string{"testhost"},
			modify: func(f flags) flags {
				f.sourcePort = "123"
				return f
			},
			wantConfig: &sourcePortSet,
			wantErr:    false,
		},
		{
			name: "verbose",
			args: []string{"testhost"},
			modify: func(f flags) flags {
				f.verbose = true
				return f
			},
			wantConfig: &verboseSet,
			wantErr:    false,
		},
		{
			name: "listen mode",
			args: []string{},
			modify: func(f flags) flags {
				f.listen = true
				return f
			},
			wantConfig: &listenMode,
			wantErr:    false,
		},
		{
			name: "timings",
			args: []string{"testhost"},
			modify: func(f flags) flags {
				f.timingTimeout = "20s"
				f.timingWait = "10s"
				f.timingDelay = "30s"
				return f
			},
			wantConfig: &setTimings,
			wantErr:    false,
		},
		{
			name: "proxy",
			args: []string{"testhost"},
			modify: func(f flags) flags {
				f.proxyAddress = "proxyhost:1234"
				return f
			},
			wantConfig: &setTimings,
			wantErr:    true,
		},
		{
			name: "ssl missing key file",
			args: []string{"testhost"},
			modify: func(f flags) flags {
				f.sslEnabled = true
				f.sslKeyFilePath = "cert.pem"
				return f
			},
			wantConfig: &sslConfig,
			wantErr:    true,
		},
		{
			name: "ssl missing cert file",
			args: []string{"testhost"},
			modify: func(f flags) flags {
				f.sslEnabled = true
				f.sslCertFilePath = "cert.pem"
				return f
			},
			wantConfig: &sslConfig,
			wantErr:    true,
		},
		{
			name: "crlf",
			args: []string{"testhost"},
			modify: func(f flags) flags {
				f.eolCRLF = true
				return f
			},
			wantConfig: &clrfConfig,
			wantErr:    false,
		},
		{
			name: "invalid proxy",
			args: []string{"testhost"},
			modify: func(f flags) flags {
				f.proxyAddress = "aa"
				return f
			},
			wantErr: true,
		},
		{
			name: "invalid proxy type",
			args: []string{"testhost"},
			modify: func(f flags) flags {
				f.proxyAddress = "proxyhost:1234"
				f.proxyType = "socks4"
				return f
			},
			wantErr: true,
		},
		{
			name: "invalid proxy dns type",
			modify: func(f flags) flags {
				f.proxyAddress = "proxyhost:1234"
				f.proxyType = "socks5"
				f.proxydns = "both"
				return f
			},
			wantErr: true,
		},
		{
			name: "chat mode",
			modify: func(f flags) flags {
				f.listen = true
				f.chatMode = true
				return f
			},
			wantConfig: &chatModeConfig,
		},
		{
			name: "invalid keep-open mode",
			modify: func(f flags) flags {
				f.udpSocket = true
				f.keepOpen = true
				return f
			},
			wantErr: true,
		},
		{
			name: "invalid broker mode",
			modify: func(f flags) flags {
				f.udpSocket = true
				f.brokerMode = true
				return f
			},
			wantErr: true,
		},
		{
			name: "invalid chat mode",
			modify: func(f flags) flags {
				f.udpSocket = true
				f.chatMode = true
				return f
			},
			wantErr: true,
		},
		{
			name:       "port scan",
			args:       []string{"testhost", "1234-1345"},
			modify:     func(f flags) flags { return f },
			wantConfig: &portScan,
			wantErr:    false,
		},
		{
			name:    "invalid port scan",
			args:    []string{"testhost", "12aa34-1345"},
			modify:  func(f flags) flags { return f },
			wantErr: true,
		},
		{
			name:    "invalid port scan array",
			args:    []string{"testhost", "1234-1345-333"},
			modify:  func(f flags) flags { return f },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Modify flags
			flags := tt.modify(baseFlags)
			gotConfig, err := evalParams(tt.args, flags)
			if err != nil {
				if tt.wantErr {
					return
				}

				t.Errorf("evalParams() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Assert config
			if !tt.wantErr {
				diff := cmp.Diff(gotConfig, tt.wantConfig, cmpopts.IgnoreFields(netcat.Config{}, "Output.OutFileMutex", "Output.OutFileHexMutex", "Output.Logger"))
				if diff != "" {
					t.Errorf("evalParams() diff : %v", diff)
				}
			}
		})
	}
}

type closableBuffer struct {
	bytes.Buffer
}

func (cb *closableBuffer) Close() error {
	return nil
}

func TestCommand(t *testing.T) {
	// Mock inputs
	stdin := bytes.NewBufferString("input data")
	stdout := new(closableBuffer)
	stderr := new(bytes.Buffer)
	config := &netcat.Config{} // Assuming Config is a struct within the netcat package
	args := []string{"arg1", "arg2"}

	// Expected cmd struct
	expectedCmd := &cmd{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
		config: config,
		args:   args,
	}

	// Call the function
	resultCmd, err := command(stdin, stdout, stderr, config, args)
	// Verify no error is returned
	if err != nil {
		t.Errorf("command() error = %v, wantErr %v", err, nil)
	}

	// Verify the result
	if !reflect.DeepEqual(resultCmd, expectedCmd) {
		t.Errorf("command() = %v, want %v", resultCmd, expectedCmd)
	}
}
