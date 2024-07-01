// Copyright 2012-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// netcat creates arbitrary TCP and UDP connections and listens and sends arbitrary data.
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/netcat"
	"github.com/u-root/u-root/pkg/ulog"
	"github.com/u-root/u-root/pkg/uroot/util"
)

var (
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
	looseSourceRouterPoints []string
	sourcePort              string
	sourceAddress           string
	listen                  bool
	udpSocket               bool
	sctpSocket              bool
	zeroIo                  bool
	connectionAllowList     []string
	connectionAllowFile     string
	connectionDenyList      []string
	connectionDenyFile      string
	proxyAddress            string
	proxydns                string
	proxyType               string
	proxyAuth               string
	maxConnections          uint32
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
	sslCiphers              []string
	sslSNI                  string
	sslALPN                 []string
)

func init() {
	flag.BoolVarP(&ipv4, "ipv4", "4", false, "Use IPv4 only")
	flag.BoolVarP(&ipv6, "ipv6", "6", false, "Use IPv6 only")
	flag.BoolVarP(&udpSocket, "udp", "u", false, "Use UDP instead of default TCP")
	flag.BoolVarP(&sctpSocket, "sctp", "", false, "Use SCTP instead of default TCP")
	flag.BoolVarP(&unixSocket, "unixsock", "U", false, "Use Unix domain sockets only")
	flag.BoolVarP(&virtualSocket, "vsock", "", false, "Use virtual circuit (stream) sockets only")

	// exec
	flag.StringVarP(&execNative, "exec", "e", "", "Executes the given command")                        // EXEC_TYPE_NATIVE
	flag.StringVarP(&execSh, "sh-exec", "c", "", "Executes the given command via /bin/sh")             // EXEC_TYPE_SHELL
	flag.StringVarP(&execLua, "lua-exec", "", "", "Executes the given Lua script (filepath argument)") // EXEC_TYPE_LUA

	// connection mode options
	flag.BoolVarP(&zeroIo, "z", "z", false, "zero-I/O mode, report connection status only")
	flag.StringVarP(&sourcePort, "source-port", "p", netcat.DEFAULT_SOURCE_PORT, "Specify source port to use")
	flag.StringVarP(&sourceAddress, "source", "s", "", "Specify source address to use (doesn't affect -l)")
	flag.StringSliceVarP(&looseSourceRouterPoints, "loose-source-router-points", "g", []string{}, "Loose source routing hop points (8 max)")
	flag.UintVarP(&looseSourcePointer, "loose-source-pointer", "G", 0, "Loose source routing hop pointer (<n>)")

	// output options
	flag.BoolVarP(&verbose, "v", "v", false, "Set verbosity level (can not be used several times)")
	flag.StringVarP(&outFilePath, "output", "o", "", "Dump session data to a file")
	flag.StringVarP(&outFileHexPath, "hex-dump", "x", "", "Dump session data as hex to a file")
	flag.BoolVarP(&appendOutput, "append-output", "", false, "Append rather than clobber specified output files")

	// listen options
	flag.BoolVarP(&listen, "listen", "l", false, "Bind and listen for incoming connections")
	flag.Uint32VarP(&maxConnections, "max-conns", "m", netcat.DEFAULT_CONNECTION_MAX, "Maximum <n> simultaneous connections")
	flag.BoolVarP(&keepOpen, "keep-open", "k", false, "Accept multiple connections in listen mode")
	flag.BoolVarP(&brokerMode, "broker", "", false, "Enable Ncat's connection brokering mode")
	flag.BoolVarP(&chatMode, "chat", "", false, "Start a simple Ncat chat server")

	// timing options
	flag.StringVarP(&timingTimeout, "idle-timeout", "i", "0ms", "Idle read/write timeout")
	flag.StringVarP(&timingDelay, "delay", "d", "0ms", "Wait between read/writes")
	flag.StringVarP(&timingWait, "wait", "w", "10s", "Connect timeout")

	// misc options
	flag.BoolVarP(&eolCRLF, "crlf", "C", false, "Use CRLF for EOL sequence")
	flag.BoolVarP(&noDNS, "nodns", "n", false, "Do not resolve hostnames via DNS")
	flag.BoolVarP(&telnet, "telnet", "t", false, "Answer Telnet negotiations")
	flag.BoolVarP(&sendOnly, "send-only", "", false, "Only send data, ignoring received; quit on EOF")
	flag.BoolVarP(&receiveOnly, "recv-only", "", false, "Only receive data, never send anything")
	flag.BoolVarP(&noShutdown, "no-shutdown", "", false, "Continue half-duplex when receiving EOF on stdin")

	// Allowlist
	flag.StringSliceVarP(&connectionAllowList, "allow", "", nil, "Allow only comma-separated list of IP addresses")
	flag.StringVarP(&connectionAllowFile, "allowfile", "", "", "A file of hosts allowed to connect to Ncat")
	flag.StringSliceVarP(&connectionDenyList, "deny", "", nil, "Deny given hosts from connecting to Ncat")
	flag.StringVarP(&connectionDenyFile, "denyfile", "", "", "A file of hosts denied from connecting to Ncat")

	// proxy
	flag.StringVarP(&proxyAddress, "proxy", "", "", "Specify address of host to proxy through (<addr[:port]> )")
	flag.StringVarP(&proxydns, "proxy-dns", "", "", "Specify where to resolve proxy destination")
	flag.StringVarP(&proxyType, "proxy-type", "", "", "Specify proxy type ('http', 'socks4', 'socks5')")
	flag.StringVarP(&proxyAuth, "proxy-auth", "", "", "Authenticate with HTTP or SOCKS proxy server")

	// ssl
	flag.BoolVarP(&sslEnabled, "ssl", "", false, "Connect or listen with SSL")
	flag.StringVarP(&sslCertFilePath, "ssl-cert", "", "", "Specify SSL certificate file (PEM) for listening")
	flag.StringVarP(&sslKeyFilePath, "ssl-key", "", "", "Specify SSL private key file (PEM) for listening")
	flag.BoolVarP(&sslVerifyTrust, "ssl-verify", "", false, "Verify trust and domain name of certificates")
	flag.StringVarP(&sslTrustFilePath, "ssl-trustfile", "", "", "PEM file containing trusted SSL certificates")
	flag.StringSliceVarP(&sslCiphers, "ssl-ciphers", "", []string{}, "Cipherlist containing SSL ciphers to use")
	flag.StringVarP(&sslSNI, "ssl-servername", "", "", "Request distinct server name (SNI)")
	flag.StringSliceVarP(&sslALPN, "ssl-alpn", "", nil, "List of protocols to send via ALPN")

	flag.Usage = util.Usage(flag.Usage, netcat.USAGE)
}

func evalParams() (*netcat.Config, error) {
	var err error

	config := netcat.DefaultConfig()

	flag.Parse()

	args := flag.Args()
	if len(args) >= 1 {
		config.Host = args[0]
	}

	if len(args) >= 2 {
		port, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid port: %v", err)
		}
		config.Port = port
	}

	// protocol options

	// IP Type
	if ipv4 && ipv6 {
		return nil, fmt.Errorf("cannot specify both IPv4 and IPv6 explicitly")
	}

	if ipv4 {
		config.ProtocolOptions.IPType = netcat.IP_V4_STRICT
	}

	if ipv6 {
		config.ProtocolOptions.IPType = netcat.IP_V6_STRICT
	}

	// Socket Types
	config.ProtocolOptions.SocketType, err = netcat.ParseSocketType(udpSocket, sctpSocket, unixSocket, virtualSocket)
	if err != nil {
		return nil, err
	}

	config.CommandExec, err = netcat.ParseCommands(
		netcat.Exec{
			Type:    netcat.EXEC_TYPE_NATIVE,
			Command: execNative,
		},
		netcat.Exec{
			Type:    netcat.EXEC_TYPE_SHELL,
			Command: execSh,
		},
		netcat.Exec{
			Type:    netcat.EXEC_TYPE_LUA,
			Command: execLua,
		})
	if err != nil {
		return nil, err
	}

	// Loose source routing
	if looseSourcePointer != 0 || len(looseSourceRouterPoints) > 0 {
		return nil, fmt.Errorf("loose source routing is not yet supported")
	}

	config.ConnectionModeOptions.SourceHost = sourceAddress

	if sourcePort != "" {
		config.ConnectionModeOptions.SourcePort = sourcePort
	}

	config.ConnectionModeOptions.ZeroIO = zeroIo

	// OutputOptions
	if verbose {
		config.Output.Logger = ulog.Log
	}

	config.Output.OutFilePath = outFilePath
	config.Output.OutFileHexPath = outFileHexPath
	config.Output.AppendOutput = appendOutput

	// Connection Mode
	if listen {
		config.ConnectionMode = netcat.CONNECTION_MODE_LISTEN
	}

	// Listen Mode Options
	config.ListenModeOptions.MaxConnections = maxConnections
	config.ListenModeOptions.KeepOpen = keepOpen
	config.ListenModeOptions.BrokerMode = brokerMode

	config.ListenModeOptions.ChatMode = chatMode

	// broker-mode is implied by chat-mode
	if chatMode {
		config.ListenModeOptions.BrokerMode = true
	}

	// timing options
	config.Timing.Delay, err = time.ParseDuration(timingDelay)
	if err != nil {
		return nil, fmt.Errorf("invalid delay: %v", err)
	}

	config.Timing.Timeout, err = time.ParseDuration(timingTimeout)
	if err != nil {
		return nil, fmt.Errorf("invalid timeout: %v", err)
	}

	if timingWait != "" {
		config.Timing.Wait, err = time.ParseDuration(timingWait)
		if err != nil {
			return nil, fmt.Errorf("invalid wait: %v", err)
		}
	}

	// Misc Options
	// EOL
	if eolCRLF {
		config.Misc.EOL = netcat.LINE_FEED_CRLF
	}

	config.Misc.NoDNS = noDNS
	config.Misc.Telnet = telnet
	config.Misc.SendOnly = sendOnly
	config.Misc.ReceiveOnly = receiveOnly
	config.Misc.NoShutdown = noShutdown

	// Access Control: Allowlist and Denylist
	config.AccessControl, err = netcat.ParseAccessControl(connectionAllowFile, connectionAllowList, connectionDenyFile, connectionDenyList)
	if err != nil {
		return nil, err
	}

	if proxyAddress != "" && sslEnabled || proxyAddress != "" && sslVerifyTrust {
		return nil, fmt.Errorf("proxy and SSL cannot be used together")
	}

	if (proxyAddress == "" && proxyType != "") || (proxyAddress != "" && proxyType == "") {
		return nil, fmt.Errorf("proxy address and type must be specified together")
	}

	if proxyAddress != "" {
		config.ProxyConfig.Enabled = true
		config.ProxyConfig.Address = proxyAddress
		config.ProxyConfig.Auth = proxyAuth
		config.ProxyConfig.Type = netcat.ProxyTypeFromString(proxyType)
		config.ProxyConfig.DNSType = netcat.ProxyDNSTypeFromString(proxydns)

		if netcat.ProxyTypeFromString(proxyType) != netcat.PROXY_TYPE_SOCKS5 {
			return nil, fmt.Errorf("only SOCKS5 proxy type is supported")
		}

		if config.ProxyConfig.DNSType != netcat.PROXY_DNS_NONE {
			return nil, fmt.Errorf("unsupported proxy DNS type")
		}
	}

	if !reflect.DeepEqual(sslCiphers, []string{}) {
		return nil, fmt.Errorf("selection of ssl-ciphers are not yet supported")
	}

	config.SSLConfig.Enabled = sslEnabled
	config.SSLConfig.CertFilePath = sslCertFilePath
	config.SSLConfig.KeyFilePath = sslKeyFilePath
	config.SSLConfig.VerifyTrust = sslVerifyTrust
	config.SSLConfig.TrustFilePath = sslTrustFilePath
	config.SSLConfig.SNI = sslSNI
	config.SSLConfig.ALPN = sslALPN

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
		return "", "", fmt.Errorf("connection: %v", err)
	}

	address, err := c.config.Address()
	if err != nil {
		return "", "", fmt.Errorf("connection: %v", err)
	}

	return network, address, nil
}

func (c *cmd) run() error {
	network, address, err := c.connection()
	if err != nil {
		return fmt.Errorf("failed to determine connection: %v", err)
	}
	// io.Copy will block until the connection is closed, use a MultiWriter to write to stdout and the output file
	output := io.MultiWriter(c.stdout, &c.config.Output)

	if c.config.ConnectionMode == netcat.CONNECTION_MODE_LISTEN {
		return c.listenMode(netcat.NewConcurrentWriter(output), network, address)
	}

	return c.connectMode(output, network, address)
}

func main() {
	config, err := evalParams()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	c, err := command(os.Stdin, os.Stdout, os.Stderr, config, flag.Args())
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if err = c.run(); err != nil {
		log.Fatalf("error: %v", err)
	}
}
