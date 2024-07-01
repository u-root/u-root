// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/u-root/u-root/pkg/netcat"
	"github.com/u-root/u-root/pkg/ulog"
)

func TestEvalParams(t *testing.T) {
	defaultConfig := netcat.DefaultConfig()

	hostSet := netcat.DefaultConfig()
	hostSet.Host = "testhost"

	portSet := netcat.DefaultConfig()
	portSet.Host = "testhost"
	portSet.Port = 1234

	ipv4Set := netcat.DefaultConfig()
	ipv4Set.ProtocolOptions.IPType = netcat.IP_V4_STRICT

	ipv6Set := netcat.DefaultConfig()
	ipv6Set.ProtocolOptions.IPType = netcat.IP_V6_STRICT

	execNativeSet := netcat.DefaultConfig()
	execNativeSet.CommandExec.Type = netcat.EXEC_TYPE_NATIVE
	execNativeSet.CommandExec.Command = "testcommand"

	sourcePortSet := netcat.DefaultConfig()
	sourcePortSet.ConnectionModeOptions.SourcePort = "123"

	verboseSet := netcat.DefaultConfig()
	verboseSet.Output.Logger = ulog.Log

	listenMode := netcat.DefaultConfig()
	listenMode.ConnectionMode = netcat.CONNECTION_MODE_LISTEN

	setTimings := netcat.DefaultConfig()
	setTimings.Timing.Wait = 10 * time.Second
	setTimings.Timing.Timeout = 20 * time.Second
	setTimings.Timing.Delay = 30 * time.Second

	sslConfig := netcat.DefaultConfig()
	sslConfig.SSLConfig.CertFilePath = "cert.pem"
	sslConfig.SSLConfig.KeyFilePath = "key.pem"

	clrfConfig := netcat.DefaultConfig()
	clrfConfig.Misc.EOL = netcat.LINE_FEED_CRLF

	chatModeConfig := netcat.DefaultConfig()
	chatModeConfig.ListenModeOptions.ChatMode = true
	chatModeConfig.ListenModeOptions.BrokerMode = true

	// Define test cases
	tests := []struct {
		name       string
		setupFunc  func()
		wantConfig *netcat.Config
		wantErr    bool
	}{
		{
			name: "default",
			setupFunc: func() {
				os.Args = []string{"cmd"}
			},
			wantConfig: &defaultConfig,
			wantErr:    false,
		},
		{
			name: "host set",
			setupFunc: func() {
				os.Args = []string{"cmd", "testhost"}
			},
			wantConfig: &hostSet,
			wantErr:    false,
		},
		{
			name: "port set",
			setupFunc: func() {
				os.Args = []string{"cmd", "testhost", "1234"}
			},
			wantConfig: &portSet,
			wantErr:    false,
		},
		{
			name: "ipv4 set",
			setupFunc: func() {
				os.Args = []string{"cmd", "-4"}
			},
			wantConfig: &ipv4Set,
			wantErr:    false,
		},
		{
			name: "ipv6 set",
			setupFunc: func() {
				os.Args = []string{"cmd", "-6"}
			},
			wantConfig: &ipv6Set,
			wantErr:    false,
		},
		{
			name: "ipv4 & ipv6 set",
			setupFunc: func() {
				os.Args = []string{"cmd", "-4", "-6"}
			},
			wantErr: true,
		},
		{
			name: "exec native",
			setupFunc: func() {
				os.Args = []string{"cmd", "--exec=testcommand"}
			},
			wantConfig: &execNativeSet,
			wantErr:    false,
		},
		{
			name: "loose source pointer false",
			setupFunc: func() {
				os.Args = []string{"cmd", "-G", "3"}
			},
			wantConfig: &execNativeSet,
			wantErr:    true,
		},
		{
			name: "source port set",
			setupFunc: func() {
				os.Args = []string{"cmd", "-p=123"}
			},
			wantConfig: &sourcePortSet,
			wantErr:    false,
		},
		{
			name: "verbose",
			setupFunc: func() {
				os.Args = []string{"cmd", "-v"}
			},
			wantConfig: &verboseSet,
			wantErr:    false,
		},
		{
			name: "listen mode",
			setupFunc: func() {
				os.Args = []string{"cmd", "-l"}
			},
			wantConfig: &listenMode,
			wantErr:    false,
		},
		{
			name: "timings",
			setupFunc: func() {
				os.Args = []string{"cmd", "-i=20s", "--wait=10s", "--delay=30s"}
			},
			wantConfig: &setTimings,
			wantErr:    false,
		},
		{
			name: "proxy",
			setupFunc: func() {
				os.Args = []string{"cmd", "--proxy=proxyhost:1234"}
			},
			wantConfig: &setTimings,
			wantErr:    true,
		},
		{
			name: "ssl missing key file",
			setupFunc: func() {
				os.Args = []string{"cmd", "-ssl", "--ssl-key=key.pem"}
			},
			wantConfig: &sslConfig,
			wantErr:    true,
		},
		{
			name: "ssl missing cert file",
			setupFunc: func() {
				os.Args = []string{"cmd", "-ssl", "--ssl-cert=cert.pem"}
			},
			wantConfig: &sslConfig,
			wantErr:    true,
		},
		{
			name: "crlf",
			setupFunc: func() {
				os.Args = []string{"cmd", "-C"}
			},
			wantConfig: &clrfConfig,
			wantErr:    false,
		},
		{
			name: "invalid port",
			setupFunc: func() {
				os.Args = []string{"cmd", "testhost", "invalid"}
			},
			wantErr: true,
		},
		{
			name: "chat mode",
			setupFunc: func() {
				os.Args = []string{"cmd", "--chat"}
			},
			wantConfig: &chatModeConfig,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flag state

			resetGlobalVars(t)

			tt.setupFunc()

			t.Log(os.Args)
			gotConfig, err := evalParams()
			if err != nil {
				if tt.wantErr {
					return
				}

				t.Errorf("evalParams() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Assert config
			diff := cmp.Diff(gotConfig, tt.wantConfig, cmpopts.IgnoreFields(netcat.Config{}, "Output.OutFileMutex", "Output.OutFileHexMutex", "Output.Logger"))
			if diff != "" {
				t.Errorf("evalParams() diff : %v", diff)
			}
		})
	}
}

func TestCommand(t *testing.T) {
	// Mock inputs
	stdin := bytes.NewBufferString("input data")
	stdout := new(bytes.Buffer)
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

func resetGlobalVars(t *testing.T) {
	t.Helper()
	ipv4 = false
	ipv6 = false
	udpSocket = false
	sctpSocket = false
	unixSocket = false
	virtualSocket = false
	execNative = ""
	execSh = ""
	execLua = ""
	zeroIo = false
	sourcePort = netcat.DEFAULT_SOURCE_PORT
	sourceAddress = ""
	looseSourceRouterPoints = []string{}
	looseSourcePointer = 0
	verbose = false
	outFilePath = ""
	outFileHexPath = ""
	appendOutput = false
	listen = false
	maxConnections = netcat.DEFAULT_CONNECTION_MAX
	keepOpen = false
	brokerMode = false
	chatMode = false
	timingWait = "10s"
	timingDelay = "0s"
	timingTimeout = "0s"
	eolCRLF = false
	noDNS = false
	telnet = false
	sendOnly = false
	receiveOnly = false
	noShutdown = false
	connectionAllowFile = ""
	connectionDenyFile = ""
	connectionAllowList = []string{}
	connectionDenyList = []string{}
	proxyAddress = ""
	proxydns = ""
	proxyType = ""
	proxyAuthType = ""
	sslEnabled = false
	sslCertFilePath = ""
	sslKeyFilePath = ""
	sslVerifyTrust = false
	sslTrustFilePath = ""
	sslCiphers = []string{"ALL", "!aNULL", "!eNULL", "!LOW", "!EXP", "!RC4", "!MD5", "@STRENGTH"}
	sslSNI = ""
	sslALPN = nil
}
