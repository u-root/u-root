// Copyright 2012-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// netcat creates arbitrary TCP and UDP connections and listens and sends arbitrary data.
package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"crypto/tls"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/netcat"
	"github.com/u-root/u-root/pkg/uroot/util"
)

var (
	timingOptions       netcat.NetcatTimingOptions
	miscOptions         netcat.NetcatMiscOptions
	outputOptions       netcat.NetcatOutputOptions
	protOptions         netcat.NetcatProtocolOptions
	conmodeOpts         netcat.NetcatConnectModeOptions
	listenModeOpts      netcat.NetcatListenModeOptions
	proxyConfig         netcat.NetcatProxyOptions
	ssl                 netcat.NetcatSSLOptions
	ipv4                bool
	ipv6                bool
	unixSocket          bool
	virtualSocket       bool
	eolCRLF             bool
	execNative          string
	execSh              string
	execLua             string
	looseSourcePointer  uint
	sourcePort          uint
	sourceAddress       string
	listen              bool
	udpSocket           bool
	sctpSocket          bool
	zeroIo              bool
	connectionAllowList []string
	connectionAllowFile string
	connectionDenyList  []string
	connectionDenyFile  string
	proxydns            string
	proxyType           string
	proxyAuthType       string
)

func init() {
	flag.BoolVar(&outputOptions.Verbose, "v", false, "Set verbosity level (can not be used several times)")
	flag.BoolVar(&ipv4, "4", false, "Use IPv4 only")
	flag.BoolVar(&ipv6, "6", false, "Use IPv6 only")

	// TODO: tcp, udp
	flag.BoolVarP(&unixSocket, "unixsock", "U", false, "Use Unix domain sockets only")
	flag.BoolVarP(&virtualSocket, "vsock", "", false, "Use virtual circuit (stream) sockets only")

	// misc::eol
	flag.BoolVarP(&eolCRLF, "crlf", "C", false, "Use CRLF for EOL sequence")
	flag.StringVarP(&execNative, "exec", "e", "", "Executes the given command")                        // EXEC_TYPE_NATIVE
	flag.StringVarP(&execSh, "sh-exec", "c", "", "Executes the given command via /bin/sh")             // EXEC_TYPE_SHELL
	flag.StringVarP(&execLua, "lua-exec", "", "", "Executes the given Lua script (filepath argument)") // EXEC_TYPE_LUA

	flag.StringSliceVar(&conmodeOpts.LooseSourceRouterPoints, "g", []string{}, "Loose source routing hop points (8 max)")
	flag.UintVar(&looseSourcePointer, "G", 4, "Loose source routing hop pointer (<n>)")

	flag.UintVarP(&listenModeOpts.MaxConnections, "max-conns", "m", netcat.DEFAULT_CONNECTION_MAX, "Maximum <n> simultaneous connections")
	flag.UintVarP(&timingOptions.Delay, "delay", "d", 0, "Wait between read/writes")

	flag.StringVarP(&outputOptions.OutFilePath, "output", "o", "", "Dump session data to a file")
	flag.StringVarP(&outputOptions.OutFileHexPath, "hex-dump", "x", "", "Dump session data as hex to a file")
	flag.BoolVarP(&outputOptions.AppendOutput, "append-output", "", false, "Append rather than clobber specified output files")

	flag.UintVarP(&timingOptions.Timeout, "idle-timeout", "I", 0, "Idle read/write timeout")

	flag.UintVarP(&sourcePort, "source-port", "p", netcat.DEFAULT_PORT, "Specify source port to use")
	flag.StringVarP(&sourceAddress, "source", "s", "", "Specify source address to use (doesn't affect -l)")

	flag.BoolVarP(&listen, "listen", "l", false, "Bind and listen for incoming connections")

	flag.BoolVarP(&listenModeOpts.KeepOpen, "keep-open", "k", false, "Accept multiple connections in listen mode")
	flag.BoolVarP(&miscOptions.NoDns, "nodns", "n", false, "Do not resolve hostnames via DNS")
	flag.BoolVarP(&miscOptions.Telnet, "telnet", "t", false, "Answer Telnet negotiations")

	// socket type
	flag.BoolVarP(&udpSocket, "udp", "u", false, "Use UDP instead of default TCP")
	flag.BoolVarP(&sctpSocket, "sctp", "", false, "Use SCTP instead of default TCP")

	flag.UintVarP(&timingOptions.Timeout, "timeout", "w", 0, "Connect timeout")
	flag.BoolVarP(&zeroIo, "", "z", false, "ero-I/O mode, report connection status only")

	flag.BoolVarP(&miscOptions.SendOnly, "send-only", "", false, "Only send data, ignoring received; quit on EOF")
	flag.BoolVarP(&miscOptions.ReceiveOnly, "recv-only", "", false, "Only receive data, never send anything")

	flag.BoolVarP(&miscOptions.NoShutdown, "no-shutdown", "", false, "Continue half-duplex when receiving EOF on stdin")

	flag.StringSliceVarP(&connectionAllowList, "allow", "", nil, "Allow only comma-separated list of IP addresses")
	flag.StringVarP(&connectionAllowFile, "allowfile", "", "", "A file of hosts allowed to connect to Ncat")
	flag.StringSliceVarP(&connectionDenyList, "deny", "", nil, "Deny given hosts from connecting to Ncat")
	flag.StringVarP(&connectionDenyFile, "denyfile", "", "", "A file of hosts denied from connecting to Ncat")

	// Allowlist
	flag.BoolVarP(&listenModeOpts.BrokerMode, "broker", "", false, "Enable Ncat's connection brokering mode")
	flag.BoolVarP(&listenModeOpts.ChatMode, "chat", "", false, "Start a simple Ncat chat server")

	flag.StringVarP(&proxyConfig.Address, "proxy", "", "", "Specify address of host to proxy through (<addr[:port]> )")
	flag.StringVarP(&proxydns, "proxy-dns", "", "", "Specify where to resolve proxy destination")

	flag.StringVarP(&proxyType, "proxy-type", "", "", "Specify proxy type ('http', 'socks4', 'socks5')")
	flag.StringVarP(&proxyAuthType, "proxy-auth", "", "", "Authenticate with HTTP or SOCKS proxy server")

	// ssl
	flag.BoolVarP(&ssl.Enabled, "ssl", "", false, "Connect or listen with SSL")
	flag.StringVarP(&ssl.CertFilePath, "ssl-cert", "", "", "Specify SSL certificate file (PEM) for listening")
	flag.StringVarP(&ssl.KeyFilePath, "ssl-key", "", "", "Specify SSL private key file (PEM) for listening")
	flag.BoolVarP(&ssl.VerifyTrust, "ssl-verify", "", false, "Verify trust and domain name of certificates")
	flag.StringVarP(&ssl.TrustFilePath, "ssl-trustfile", "", "", "PEM file containing trusted SSL certificates")
	flag.StringSliceVarP(&ssl.Ciphers, "ssl-ciphers", "", []string{netcat.DEFAULT_SSL_SUITE_STR}, "Cipherlist containing SSL ciphers to use")
	flag.StringVarP(&ssl.SNI, "ssl-servername", "", "", "Request distinct server name (SNI)")
	flag.StringSliceVarP(&ssl.ALPN, "ssl-alpn", "", nil, "List of protocols to send via ALPN")

	flag.Usage = util.Usage(flag.Usage, netcat.USAGE)
}

func evalParams() (netcat.NetcatConfig, error) {
	flag.Parse()

	// IP Type
	ipType := netcat.DEFAULT_IP_TYPE
	if ipv4 && ipv6 {
		log.Fatal("Cannot specify both IPv4 and IPv6 explicitly")
	}
	if ipv4 {
		ipType = netcat.IP_V4_STRICT
	}

	if ipv6 {
		ipType = netcat.IP_V6_STRICT
	}
	protOptions.IPType = ipType

	// EOL
	eol := netcat.DEFAULT_LF
	if eolCRLF {
		eol = netcat.LINE_FEED_CRLF
	}
	miscOptions.EOL = eol

	// Exec commands
	execs := []string{
		execNative,
		execSh,
		execLua,
	}
	exec, err := netcat.ParseCommands(execs)
	if err != nil {
		return netcat.NetcatConfig{}, err
	}

	// Loose source routing
	if looseSourcePointer%4 != 0 || looseSourcePointer > 28 {
		return netcat.NetcatConfig{}, fmt.Errorf("loose source routing hop pointer must be a multiple of 4 and less than 28")
	}

	// Connection Mode
	conMode := netcat.DEFAULT_CONNECTION_MODE
	if listen {
		conMode = netcat.CONNECTION_MODE_LISTEN
	}

	// Socket Types
	protOptions.SocketType, err = netcat.ParseSocketType(udpSocket, sctpSocket, unixSocket, virtualSocket)
	if err != nil {
		return netcat.NetcatConfig{}, err
	}

	// Access Control: Allowlist and Denylist
	accessControl, err := netcat.ParseAccessControl(connectionAllowFile, connectionAllowList, connectionDenyFile, connectionDenyList)
	if err != nil {
		return netcat.NetcatConfig{}, err
	}

	proxyConfig.DNSType = netcat.ProxyDNSTypeFromString(proxydns)
	proxyConfig.Type = netcat.ProxyTypeFromString(proxyType)
	proxyConfig.Type = netcat.ProxyTypeFromString(proxyAuthType)

	if err := ssl.Verify(); err != nil {
		return netcat.NetcatConfig{}, err
	}

	return netcat.NetcatConfig{
		ConnectionMode:        conMode,
		ConnectionModeOptions: conmodeOpts,
		ListenModeOptions:     listenModeOpts,
		ProtocolOptions:       protOptions,
		Hostname:              sourceAddress,
		Port:                  sourcePort,
		SSLConfig:             ssl,
		ProxyConfig:           proxyConfig,
		AccessControl:         accessControl,
		CommandExec:           exec,
		ZeroIo:                zeroIo,
	}, nil
}

type cmd struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
	config *netcat.NetcatConfig
	args   []string
}

func command(stdin io.Reader, stdout io.Writer, stderr io.Writer, config *netcat.NetcatConfig, args []string) (*cmd, error) {
	return &cmd{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
		config: config,
		args:   args,
	}, nil
}

// From the prepared config generate a network connection that will be used for the netcat command
// TODO: create own NetcatConnection struct
func (c *cmd) connection() (io.ReadWriter, error) {
	network, err := c.config.ProtocolOptions.SocketType.ToGoType(c.config.ProtocolOptions.IPType)
	if err != nil {
		return nil, fmt.Errorf("connection: invalid socket type: %v", err)
	}
	netcat.FLogf(c.config, c.stderr, "network type = %q\n", network)

	if c.config.ConnectionMode == netcat.CONNECTION_MODE_LISTEN && len(c.args) == 0 {
		c.config.Port = netcat.DEFAULT_PORT
	}

	connectHost := c.config.Address()
	// connectHost := c.config.Hostname
	// if connectHost == "" {
	// 	connectHost = c.args[0]
	// }

	// connectPort := c.config.Port
	// if connectPort == 0 {
	// 	connectPort, err := strconv.Atoi(c.args[1])
	// 	if err != nil {
	// 		return nil, fmt.Errorf("connection: %v", err)
	// 	}
	// }

	// connectAddress := connectHost + ":" + connectPort

	switch c.config.ProtocolOptions.SocketType {
	case netcat.SOCKET_TYPE_TCP, netcat.SOCKET_TYPE_UNIX:
		// Listen Mode
		if c.config.ConnectionMode == netcat.CONNECTION_MODE_LISTEN {
			if c.config.SSLConfig.Enabled {
				// TLS
				cer, err := tls.LoadX509KeyPair(c.config.SSLConfig.CertFilePath, c.config.SSLConfig.KeyFilePath)
				if err != nil {
					return nil, fmt.Errorf("connection: %v", err)
				}

				tlsConfig := &tls.Config{Certificates: []tls.Certificate{cer}}
				tlsListen, err := tls.Listen(network, c.config.Address(), tlsConfig)
				if err != nil {
					return nil, fmt.Errorf("connection: %v", err)
				}

				return tlsListen.Accept()
			} else {
				// No TLS
				address := c.config.Address()
				listen, err := net.Listen(network, address)
				if err != nil {
					return nil, err
				}
				netcat.FLogf(c.config, c.stderr, "Listening on %v\n", listen.Addr())
				return listen.Accept()
			}
		} else {
			// Connection Mode
			if c.config.SSLConfig.Enabled {
				// TLS
				tlsConfig := &tls.Config{InsecureSkipVerify: true}
				return tls.Dial(network, connectHost, tlsConfig)
			} else {
				// No TLS
				return net.Dial(network, connectHost)
			}
		}

	case netcat.SOCKET_TYPE_UDP:
		// Listen Mode
		if c.config.ConnectionMode == netcat.CONNECTION_MODE_LISTEN {
			return netcat.NewUdpRemoteConn(network, connectHost, c.stderr, c.config.Output.Verbose)
		} else {
			// Connection Mode
			udpAddr, err := net.ResolveUDPAddr(network, connectHost)
			if err != nil {
				return nil, err
			}
			return net.DialUDP(network, nil, udpAddr)
		}

	// unsupported socket types
	case netcat.SOCKET_TYPE_SCTP, netcat.SOCKET_TYPE_UDP_VSOCK, netcat.SOCKET_TYPE_UDP_UNIX:
		return nil, fmt.Errorf("currently unsupported socket type %q", c.config.ProtocolOptions.SocketType)

	case netcat.SOCKET_TYPE_NONE:
	default:
		return nil, fmt.Errorf("undefined socket type %q", c.config.ProtocolOptions.SocketType)
	}
	return nil, fmt.Errorf("undefined socket type %q", c.config.ProtocolOptions.SocketType)
}

func (c *cmd) run() error {
	// Netcat can operate in 2 modes: connect or listen
	// These modes will automatically be handled by the io.ReadWriter returned by the connection function
	// TODO: implement that for Netcat new? Can we handle all the tls/proxy stuff in the connection function?

	conn, err := c.connection()
	if err != nil {
		return fmt.Errorf("run: %v", err)
	}

	netcat.Logf(c.config, "client is in %q mode\n", c.config.ConnectionMode.String())
	go func() {
		if _, err := io.Copy(conn, c.stdin); err != nil {
			fmt.Fprintf(c.stderr, "run send: %v\n", err)
		}
	}()

	// io.Copy will block until the connection is closed, use a MultiWriter to write to stdout and the output file
	mw := io.MultiWriter(c.stdout, &c.config.Output)
	if n, err := io.Copy(mw, conn); err != nil {
		fmt.Fprintf(c.stderr, "run dump: %v, n = %v\n", err, n)
		return fmt.Errorf("run dump: %v", err)
	}

	netcat.Logf(c.config, "disconnected\n")
	return err
}

func main() {
	config, err := evalParams()
	if err != nil {
		log.Fatalf("error: %v", err)
		flag.Usage()
		os.Exit(1)
	}

	fmt.Printf("config: %+v, args = %+v\n", config, flag.Args())

	c, err := command(os.Stdin, os.Stdout, os.Stderr, &config, flag.Args())
	if err != nil {
		fmt.Printf("error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	if err = c.run(); err != nil {
		log.Fatalf("netcat: %v", err)
	}
}
