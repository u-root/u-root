// Copyright 2012-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

// netcat creates arbitrary TCP and UDP connections and listens and sends arbitrary data.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/u-root/u-root/pkg/netcat"
	"github.com/u-root/u-root/pkg/ulog"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

type flags struct {
	verbose                 bool
	timingDelay             string
	timingTimeout           string
	timingWait              string
	ipv4                    bool
	ipv6                    bool
	unixSocket              bool
	virtualSocket           bool
	eolCRLF                 bool
	execNative              string
	execSh                  string
	execLua                 string
	looseSourcePointer      uint
	looseSourceRouterPoints string
	sourcePort              string
	sourceAddress           string
	listen                  bool
	udpSocket               bool
	sctpSocket              bool
	zeroIo                  bool
	connectionAllowList     string
	connectionAllowFile     string
	connectionDenyList      string
	connectionDenyFile      string
	proxyAddress            string
	proxydns                string
	proxyType               string
	proxyAuth               string
	maxConnections          uint64
	keepOpen                bool
	noDNS                   bool
	telnet                  bool
	outFilePath             string
	outFileHexPath          string
	appendOutput            bool
	sendOnly                bool
	receiveOnly             bool
	noShutdown              bool
	brokerMode              bool
	chatMode                bool
	sslEnabled              bool
	sslCertFilePath         string
	sslKeyFilePath          string
	sslVerifyTrust          bool
	sslTrustFilePath        string
	sslCiphers              string
	sslSNI                  string
	sslALPN                 string
}

func evalParams(args []string, f flags) (*netcat.Config, error) {
	var err error

	config := netcat.DefaultConfig()

	// Connection Mode
	if f.listen {
		config.ConnectionMode = netcat.CONNECTION_MODE_LISTEN
	}

	// in connect mode the first arg is necessarily host
	// The port is optional as the second argument.
	switch config.ConnectionMode {
	case netcat.CONNECTION_MODE_CONNECT:
		if len(args) < 1 {
			return nil, fmt.Errorf("missing host")
		}

		config.Host = args[0]

		if len(args) >= 2 {
			ports := strings.SplitN(args[1], "-", 2)

			if len(ports) > 2 {
				return nil, fmt.Errorf("too many arguments")
			}

			port, err := strconv.ParseUint(ports[0], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid port: %w", err)
			}

			config.Port = port

			// port-scanning
			if len(ports) == 2 {
				config.ConnectionModeOptions.ScanPorts = true
				config.ConnectionModeOptions.CurrentPort = port
				port, err = strconv.ParseUint(ports[1], 10, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid port: %w", err)
				}
				config.ConnectionModeOptions.EndPort = port
			}

		}

	// If one argument is given in listen mode it is expected to be the port.
	// If two args are given the first arg is the host and the second is the port.
	case netcat.CONNECTION_MODE_LISTEN:
		if len(args) == 1 {
			port, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid port: %w", err)
			}
			config.Port = port
		} else if len(args) >= 2 {
			config.Host = args[0]

			port, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid port: %w", err)
			}
			config.Port = port
		}
	}

	// IP Type
	if f.ipv4 && f.ipv6 {
		return nil, fmt.Errorf("cannot specify both IPv4 and IPv6 explicitly")
	}

	if f.ipv4 {
		config.ProtocolOptions.IPType = netcat.IP_V4_STRICT
	}

	if f.ipv6 {
		config.ProtocolOptions.IPType = netcat.IP_V6_STRICT
	}

	// Socket Types
	config.ProtocolOptions.SocketType, err = netcat.ParseSocketType(f.udpSocket, f.unixSocket, f.virtualSocket, f.sctpSocket)
	if err != nil {
		return nil, err
	}

	config.CommandExec, err = netcat.ParseCommands(
		netcat.Exec{
			Type:    netcat.EXEC_TYPE_NATIVE,
			Command: f.execNative,
		},
		netcat.Exec{
			Type:    netcat.EXEC_TYPE_SHELL,
			Command: f.execSh,
		},
		netcat.Exec{
			Type:    netcat.EXEC_TYPE_LUA,
			Command: f.execLua,
		})
	if err != nil {
		return nil, err
	}

	// Loose source routing
	if f.looseSourcePointer != 0 || len(f.looseSourceRouterPoints) > 0 {
		return nil, fmt.Errorf("loose source routing is not yet supported")
	}

	config.ConnectionModeOptions.SourceHost = f.sourceAddress

	if f.sourcePort != "" {
		config.ConnectionModeOptions.SourcePort = f.sourcePort
	}

	config.ConnectionModeOptions.ZeroIO = f.zeroIo

	// OutputOptions
	if f.verbose {
		config.Output.Logger = ulog.Log
	}

	config.Output.OutFilePath = f.outFilePath
	config.Output.OutFileHexPath = f.outFileHexPath
	config.Output.AppendOutput = f.appendOutput

	// Listen Mode Options
	config.ListenModeOptions.MaxConnections = uint32(f.maxConnections)
	config.ListenModeOptions.KeepOpen = f.keepOpen
	config.ListenModeOptions.BrokerMode = f.brokerMode

	config.ListenModeOptions.ChatMode = f.chatMode

	// broker-mode is implied by chat-mode
	if f.chatMode {
		config.ListenModeOptions.BrokerMode = true
	}

	// timing options
	config.Timing.Delay, err = time.ParseDuration(f.timingDelay)
	if err != nil {
		return nil, fmt.Errorf("invalid delay: %w", err)
	}

	config.Timing.Timeout, err = time.ParseDuration(f.timingTimeout)
	if err != nil {
		return nil, fmt.Errorf("invalid timeout: %w", err)
	}

	if f.timingWait != "" {
		config.Timing.Wait, err = time.ParseDuration(f.timingWait)
		if err != nil {
			return nil, fmt.Errorf("invalid wait: %w", err)
		}
	}

	// Misc Options
	// EOL
	config.Misc.EOL = netcat.DEFAULT_LF

	if f.eolCRLF {
		config.Misc.EOL = netcat.LINE_FEED_CRLF
	}

	config.Misc.NoDNS = f.noDNS
	config.Misc.Telnet = f.telnet
	config.Misc.SendOnly = f.sendOnly
	config.Misc.ReceiveOnly = f.receiveOnly
	config.Misc.NoShutdown = f.noShutdown

	// Access Control: Allowlist and Denylist
	connectionAllowList := strings.Split(f.connectionAllowList, ",")
	if connectionAllowList[0] == "" {
		connectionAllowList = []string{}
	}

	connectionDenyList := strings.Split(f.connectionDenyList, ",")
	if connectionDenyList[0] == "" {
		connectionDenyList = []string{}
	}

	config.AccessControl, err = netcat.ParseAccessControl(f.connectionAllowFile, connectionAllowList, f.connectionDenyFile, connectionDenyList)
	if err != nil {
		return nil, err
	}

	if f.proxyAddress != "" && f.sslEnabled || f.proxyAddress != "" && f.sslVerifyTrust {
		return nil, fmt.Errorf("proxy and SSL cannot be used together")
	}

	if (f.proxyAddress == "" && f.proxyType != "") || (f.proxyAddress != "" && f.proxyType == "") {
		return nil, fmt.Errorf("proxy address and type must be specified together")
	}

	if f.proxyAddress != "" {
		config.ProxyConfig.Enabled = true
		config.ProxyConfig.Address = f.proxyAddress
		config.ProxyConfig.Auth = f.proxyAuth
		config.ProxyConfig.Type = netcat.ProxyTypeFromString(f.proxyType)
		config.ProxyConfig.DNSType = netcat.ProxyDNSTypeFromString(f.proxydns)

		if netcat.ProxyTypeFromString(f.proxyType) != netcat.PROXY_TYPE_SOCKS5 {
			return nil, fmt.Errorf("only SOCKS5 proxy type is supported")
		}

		if config.ProxyConfig.DNSType != netcat.PROXY_DNS_NONE {
			return nil, fmt.Errorf("unsupported proxy DNS type")
		}
	}

	if !reflect.DeepEqual(f.sslCiphers, "") {
		return nil, fmt.Errorf("selection of ssl-ciphers are not yet supported")
	}

	config.SSLConfig.Enabled = f.sslEnabled
	config.SSLConfig.CertFilePath = f.sslCertFilePath
	config.SSLConfig.KeyFilePath = f.sslKeyFilePath
	config.SSLConfig.VerifyTrust = f.sslVerifyTrust
	config.SSLConfig.TrustFilePath = f.sslTrustFilePath
	config.SSLConfig.SNI = f.sslSNI
	alpn := strings.Split(f.sslALPN, ",")
	if alpn[0] != "" {
		config.SSLConfig.ALPN = alpn
	}

	if config.SSLConfig.CertFilePath != "" {
		if _, err := os.Stat(config.SSLConfig.CertFilePath); err != nil {
			return nil, fmt.Errorf("certificate file does not exist")
		}
	}

	if config.SSLConfig.KeyFilePath != "" {
		if _, err := os.Stat(config.SSLConfig.KeyFilePath); err != nil {
			return nil, fmt.Errorf("key file does not exist")
		}
	}

	return &config, nil
}

type cmd struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
	config *netcat.Config
	args   []string
}

func command(stdin io.Reader, stdout io.Writer, stderr io.Writer, config *netcat.Config, args []string) (*cmd, error) {
	return &cmd{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
		config: config,
		args:   args,
	}, nil
}

// From the prepared config generate a network connection that will be used for the netcat command
func (c *cmd) connection() (string, string, error) {
	// check if SSL is available for the selected protocol if enabled
	if c.config.SSLConfig.Enabled || c.config.SSLConfig.VerifyTrust {
		switch s := c.config.ProtocolOptions.SocketType; s {
		case netcat.SOCKET_TYPE_UDP, netcat.SOCKET_TYPE_UNIX, netcat.SOCKET_TYPE_UDP_UNIX, netcat.SOCKET_TYPE_UDP_VSOCK, netcat.SOCKET_TYPE_VSOCK, netcat.SOCKET_TYPE_SCTP:
			return "", "", fmt.Errorf("SSL is not available for %s", s)
		}
	}

	network, err := c.config.ProtocolOptions.Network()
	if err != nil {
		return "", "", fmt.Errorf("connection: %w", err)
	}

	address, err := c.config.Address()
	if err != nil {
		return "", "", fmt.Errorf("connection: %w", err)
	}

	return network, address, nil
}

func run(args []string) error {
	var f flags

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)

	fs.BoolVar(&f.ipv4, "ipv4", false, "Use IPv4 only")
	fs.BoolVar(&f.ipv4, "4", false, "Use IPv4 only (shorthand)")

	fs.BoolVar(&f.ipv6, "ipv6", false, "Use IPv6 only")
	fs.BoolVar(&f.ipv6, "6", false, "Use IPv6 only (shorthand)")

	fs.BoolVar(&f.udpSocket, "udp", false, "Use UDP instead of default TCP")
	fs.BoolVar(&f.udpSocket, "u", false, "Use UDP instead of default TCP (shorthand)")

	fs.BoolVar(&f.sctpSocket, "sctp", false, "Use SCTP instead of default TCP")

	fs.BoolVar(&f.unixSocket, "unixsock", false, "Use Unix domain sockets only")
	fs.BoolVar(&f.unixSocket, "U", false, "Use Unix domain sockets only (shorthand)")

	fs.BoolVar(&f.virtualSocket, "vsock", false, "Use virtual circuit (stream) sockets only")

	// exec
	fs.StringVar(&f.execNative, "exec", "", "Executes the given command")          // EXEC_TYPE_NATIVE
	fs.StringVar(&f.execNative, "e", "", "Executes the given command (shorthand)") // EXEC_TYPE_NATIVE

	fs.StringVar(&f.execSh, "sh-exec", "", "Executes the given command via /bin/sh")       // EXEC_TYPE_SHELL
	fs.StringVar(&f.execSh, "c", "", "Executes the given command via /bin/sh (shorthand)") // EXEC_TYPE_SHELL

	fs.StringVar(&f.execLua, "lua-exec", "", "Executes the given Lua script (filepath argument)") // EXEC_TYPE_LUA

	// connection mode options
	fs.BoolVar(&f.zeroIo, "z", false, "zero-I/O mode, report connection status only")

	fs.StringVar(&f.sourcePort, "source-port", netcat.DEFAULT_SOURCE_PORT, "Specify source port to use")
	fs.StringVar(&f.sourcePort, "p", netcat.DEFAULT_SOURCE_PORT, "Specify source port to use (shorthand)")

	fs.StringVar(&f.sourceAddress, "source", "", "Specify source address to use")
	fs.StringVar(&f.sourceAddress, "s", "", "Specify source address to use (shorthand)")

	fs.StringVar(&f.looseSourceRouterPoints, "g", "", "Loose source routing hop points (8 max)")
	fs.UintVar(&f.looseSourcePointer, "G", 0, "Loose source routing hop pointer (<n>)")

	// output options
	fs.BoolVar(&f.verbose, "verbose", false, "Set verbosity level (can not be used several times)")
	fs.BoolVar(&f.verbose, "v", false, "Set verbosity level (shorthand)")

	fs.StringVar(&f.outFilePath, "output", "", "Dump session data to a file")
	fs.StringVar(&f.outFilePath, "o", "", "Dump session data to a file (shorthand)")

	fs.StringVar(&f.outFileHexPath, "hex-dump", "", "Dump session data as hex to a file")
	fs.StringVar(&f.outFileHexPath, "x", "", "Dump session data as hex to a file (shorthand)")

	fs.BoolVar(&f.appendOutput, "append-output", false, "Append rather than clobber specified output files")

	// listen options
	fs.BoolVar(&f.listen, "listen", false, "Bind and listen for incoming connections")
	fs.BoolVar(&f.listen, "l", false, "Bind and listen for incoming connections (shorthand)")

	fs.Uint64Var(&f.maxConnections, "max-conns", netcat.DEFAULT_CONNECTION_MAX, "Maximum <n> simultaneous connections")
	fs.Uint64Var(&f.maxConnections, "m", netcat.DEFAULT_CONNECTION_MAX, "Maximum <n> simultaneous connections (shorthand)")

	fs.BoolVar(&f.keepOpen, "keep-open", false, "Accept multiple connections in listen mode")
	fs.BoolVar(&f.keepOpen, "k", false, "Accept multiple connections in listen mode (shorthand)")

	fs.BoolVar(&f.brokerMode, "broker", false, "Enable Ncat's connection brokering mode")
	fs.BoolVar(&f.chatMode, "chat", false, "Start a simple Ncat chat server")

	// timing options
	fs.StringVar(&f.timingTimeout, "idle-timeout", "0ms", "Idle read/write timeout")
	fs.StringVar(&f.timingTimeout, "i", "0ms", "Idle read/write timeout (shorthand)")

	fs.StringVar(&f.timingDelay, "delay", "0ms", "Wait between read/writes")
	fs.StringVar(&f.timingDelay, "d", "0ms", "Wait between read/writes (shorthand)")

	fs.StringVar(&f.timingWait, "wait", "10s", "Connect timeout")
	fs.StringVar(&f.timingWait, "w", "10s", "Connect timeout (shorthand)")

	// misc options
	fs.BoolVar(&f.eolCRLF, "crlf", false, "Use CRLF for EOL sequence")
	fs.BoolVar(&f.eolCRLF, "C", false, "Use CRLF for EOL sequence")

	fs.BoolVar(&f.noDNS, "nodns", false, "Do not resolve hostnames via DNS")
	fs.BoolVar(&f.noDNS, "n", false, "Do not resolve hostnames via DNS (shorthand)")

	fs.BoolVar(&f.telnet, "telnet", false, "Answer Telnet negotiations")
	fs.BoolVar(&f.telnet, "t", false, "Answer Telnet negotiations (shorthand)")

	fs.BoolVar(&f.sendOnly, "send-only", false, "Only send data, ignoring received; quit on EOF")
	fs.BoolVar(&f.receiveOnly, "recv-only", false, "Only receive data, never send anything")
	fs.BoolVar(&f.noShutdown, "no-shutdown", false, "Continue half-duplex when receiving EOF on stdin")

	// Allowlist
	fs.StringVar(&f.connectionAllowList, "allow", "", "Allow only comma-separated list of IP addresses")
	fs.StringVar(&f.connectionAllowFile, "allowfile", "", "A file of hosts allowed to connect to Ncat")
	fs.StringVar(&f.connectionDenyList, "deny", "", "Deny given hosts from connecting to Ncat")
	fs.StringVar(&f.connectionDenyFile, "denyfile", "", "A file of hosts denied from connecting to Ncat")

	// proxy
	fs.StringVar(&f.proxyAddress, "proxy", "", "Specify address of host to proxy through (<addr[:port]> )")
	fs.StringVar(&f.proxydns, "proxy-dns", "", "Specify where to resolve proxy destination")
	fs.StringVar(&f.proxyType, "proxy-type", "", "Specify proxy type ('http', 'socks4', 'socks5')")
	fs.StringVar(&f.proxyAuth, "proxy-auth", "", "Authenticate with HTTP or SOCKS proxy server")

	// ssl
	fs.BoolVar(&f.sslEnabled, "ssl", false, "Connect or listen with SSL")
	fs.StringVar(&f.sslCertFilePath, "ssl-cert", "", "Specify SSL certificate file (PEM) for listening")
	fs.StringVar(&f.sslKeyFilePath, "ssl-key", "", "Specify SSL private key file (PEM) for listening")
	fs.BoolVar(&f.sslVerifyTrust, "ssl-verify", false, "Verify trust and domain name of certificates")
	fs.StringVar(&f.sslTrustFilePath, "ssl-trustfile", "", "PEM file containing trusted SSL certificates")
	fs.StringVar(&f.sslCiphers, "ssl-ciphers", "", "Cipherlist containing SSL ciphers to use")
	fs.StringVar(&f.sslSNI, "ssl-servername", "", "Request distinct server name (SNI)")
	fs.StringVar(&f.sslALPN, "ssl-alpn", "", "List of protocols to send via ALPN")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: %s [options] [network address]\nOptions:\n", args[0])
		fs.PrintDefaults()
	}

	fs.Parse(unixflag.ArgsToGoArgs(args[1:]))

	config, err := evalParams(fs.Args(), f)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	c, err := command(os.Stdin, os.Stdout, os.Stderr, config, flag.Args())
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	network, address, err := c.connection()
	if err != nil {
		return fmt.Errorf("failed to determine connection: %w", err)
	}

	// io.Copy will block until the connection is closed, use a MultiWriter to write to stdout and the output file
	output := io.MultiWriter(c.stdout, &c.config.Output)

	if c.config.ConnectionMode == netcat.CONNECTION_MODE_LISTEN {
		return c.listenMode(netcat.NewConcurrentWriter(output), network, address)
	}

	return c.connectMode(output, network, address)
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatalf("error: %v", err)
	}
}
