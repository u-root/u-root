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

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/netcat"
	"github.com/u-root/u-root/pkg/uroot/util"
)

var errMissingHostnamePort = fmt.Errorf("missing hostname:port")

func parseParams() netcat.NetcatConfig {
	// TODO: create  tt for mutual exclusive flags and reduce it to KNF / DNF and evaluate the flagset
	// With that we don't need to check for mutual exclusive flags but can discard wrong flags at once

	// TODO: make sanity check for network?
	// TODO: add parsing for optional array params that need parsing of multiple values and default empty string
	timingOptions := netcat.NetcatTimingOptions{}
	miscOptions := netcat.NetcatMiscOptions{}
	outputOptions := netcat.NetcatOutputOptions{}

	flag.BoolVar(&outputOptions.Verbose, "v", false, "Set verbosity level (can not be used several times)")
	ipv4 := flag.Bool("4", false, "Use IPv4 only")
	ipv6 := flag.Bool("6", false, "Use IPv6 only")

	protOptions := netcat.NetcatProtocolOptions{}
	ipType := netcat.DEFAULT_IP_TYPE
	if *ipv4 && *ipv6 {
		log.Fatal("Cannot specify both IPv4 and IPv6 explicitly")
	}

	if *ipv4 {
		ipType = netcat.IP_V4_STRICT
	}

	if *ipv6 {
		ipType = netcat.IP_V6_STRICT
	}

	protOptions.IPType = ipType

	// TODO: tcp, udp
	unixSocket := flag.BoolP("unixsock", "U", false, "Use Unix domain sockets only")
	virtualSocket := flag.BoolP("vsock", "", false, "Use virtual circuit (stream) sockets only")

	// misc::eol
	eolCRLF := flag.BoolP("crlf", "C", false, "Use CRLF for EOL sequence")
	eol := netcat.DEFAULT_LF
	if *eolCRLF {
		eol = netcat.LINE_FEED_CRLF
	}
	miscOptions.EOL = eol

	execs := []*string{
		flag.StringP("exec", "e", "", "Executes the given command"),                           // EXEC_TYPE_NATIVE
		flag.StringP("sh-exec", "c", "", "Executes the given command via /bin/sh"),            // EXEC_TYPE_SHELL
		flag.StringP("lua-exec", "", "", "Executes the given Lua script (filepath argument)"), // EXEC_TYPE_LUA
	}

	exec, err := netcat.ParseCommands(execs)
	if err != nil {
		log.Fatal(err)
	}

	conmodeOpts := netcat.NetcatConnectModeOptions{}
	flag.StringSliceVar(&conmodeOpts.LooseSourceRouterPoints, "g", []string{}, "Loose source routing hop points (8 max)")
	looseSourcePointer := flag.Uint("G", 4, "Loose source routing hop pointer (<n>)")
	if *looseSourcePointer%4 != 0 || *looseSourcePointer > 28 {
		log.Fatalf("Loose source routing hop pointer must be a multiple of 4 and less than 28")
	}

	listenModeOpts := netcat.NetcatListenModeOptions{}
	flag.UintVarP(&listenModeOpts.MaxConnections, "max-conns", "m", netcat.DEFAULT_CONNECTION_MAX, "Maximum <n> simultaneous connections")
	flag.UintVarP(&timingOptions.Delay, "delay", "d", 0, "Wait between read/writes")

	// output:: they are not mutual exclusive
	flag.StringVarP(&outputOptions.OutFilePath, "output", "o", "", "Dump session data to a file")
	flag.StringVarP(&outputOptions.OutFileHexPath, "hex-dump", "x", "", "Dump session data as hex to a file")
	flag.BoolVarP(&outputOptions.AppendOutput, "append-output", "", false, "Append rather than clobber specified output files")

	flag.UintVarP(&timingOptions.Timeout, "idle-timeout", "I", 0, "Idle read/write timeout")

	sourcePort := flag.UintP("source-port", "p", 0, "Specify source port to use")
	sourceAddress := flag.StringP("source", "s", "", "Specify source address to use (doesn't affect -l)")

	conMode := netcat.DEFAULT_CONNECTION_MODE
	listen := flag.BoolP("listen", "l", false, "Bind and listen for incoming connections")
	if *listen {
		conMode = netcat.CONNECTION_MODE_LISTEN
	}

	flag.BoolVarP(&listenModeOpts.KeepOpen, "keep-open", "k", false, "Accept multiple connections in listen mode")
	flag.BoolVarP(&miscOptions.NoDns, "nodns", "n", false, "Do not resolve hostnames via DNS")
	flag.BoolVarP(&miscOptions.Telnet, "telnet", "t", false, "Answer Telnet negotiations")

	// socket type
	udpSocket := flag.BoolP("udp", "u", false, "Use UDP instead of default TCP")
	sctpSocket := flag.BoolP("sctp", "", false, "Use SCTP instead of default TCP")
	protOptions.SocketType, err = netcat.ParseSocketType(*udpSocket, *sctpSocket, *unixSocket, *virtualSocket)
	if err != nil {
		log.Fatal(err)
	}

	flag.UintVarP(&timingOptions.Timeout, "timeout", "w", 0, "Connect timeout")
	zeroIo := flag.BoolP("", "z", false, "ero-I/O mode, report connection status only")

	flag.BoolVarP(&miscOptions.SendOnly, "send-only", "", false, "Only send data, ignoring received; quit on EOF")
	flag.BoolVarP(&miscOptions.ReceiveOnly, "recv-only", "", false, "Only receive data, never send anything")

	flag.BoolVarP(&miscOptions.NoShutdown, "no-shutdown", "", false, "Continue half-duplex when receiving EOF on stdin")

	connectionAllowList := flag.StringSliceP("allow", "", nil, "Allow only comma-separated list of IP addresses")
	connectionAllowFile := flag.StringP("allowfile", "", "", "A file of hosts allowed to connect to Ncat")
	connectionDenyList := flag.StringSliceP("deny", "", nil, "Deny given hosts from connecting to Ncat")
	connectionDenyFile := flag.StringP("denyfile", "", "", "A file of hosts denied from connecting to Ncat")

	accessControl, err := netcat.ParseAccessControl(connectionAllowFile, connectionAllowList, connectionDenyFile, connectionDenyList)
	if err != nil {
		log.Fatal(err)
	}

	// Allowlist
	flag.BoolVarP(&listenModeOpts.BrokerMode, "broker", "", false, "Enable Ncat's connection brokering mode")
	flag.BoolVarP(&listenModeOpts.ChatMode, "chat", "", false, "Start a simple Ncat chat server")

	pc := netcat.NetcatProxyConfig{}
	flag.StringVarP(&pc.Address, "proxy", "", "", "Specify address of host to proxy through (<addr[:port]> )")
	flag.StringVarP(&pc.DNSAddress, "proxy-dns", "", "", "Specify where to resolve proxy destination")

	proxyType := flag.StringP("proxy-type", "", "", "Specify proxy type ('http', 'socks4', 'socks5')")
	proxyAuthType := flag.StringP("proxy-auth", "", "", "Authenticate with HTTP or SOCKS proxy server")

	// TODO: do we have to move that elsewhere?
	pc.Type = netcat.ProxyTypeFromString(*proxyType)
	pc.Type = netcat.ProxyTypeFromString(*proxyAuthType)

	// ssl
	ssl := netcat.NetcatSSLConfig{}
	flag.BoolVarP(&ssl.Enabled, "ssl", "", false, "Connect or listen with SSL")
	flag.StringVarP(&ssl.CertFilePath, "ssl-cert", "", "", "Specify SSL certificate file (PEM) for listening")
	flag.StringVarP(&ssl.KeyFilePath, "ssl-key", "", "", "Specify SSL private key file (PEM) for listening")
	flag.BoolVarP(&ssl.VerifyTrust, "ssl-verify", "", false, "Verify trust and domain name of certificates")
	flag.StringVarP(&ssl.TrustFilePath, "ssl-trustfile", "", "", "PEM file containing trusted SSL certificates")
	flag.StringSliceVarP(&ssl.Ciphers, "ssl-ciphers", "", []string{netcat.DEFAULT_SSL_SUITE_STR}, "Cipherlist containing SSL ciphers to use")
	flag.StringVarP(&ssl.SNI, "ssl-servername", "", "", "Request distinct server name (SNI)")
	flag.StringSliceVarP(&ssl.ALPN, "ssl-alpn", "", nil, "List of protocols to send via ALPN")
	flag.Parse()

	// validate the options
	if err := ssl.Verify(); err != nil {
		log.Fatal(err)
	}

	return netcat.NetcatConfig{
		ConnectionMode:        conMode,
		ConnectionModeOptions: conmodeOpts,
		ListenModeOptions:     listenModeOpts,
		ProtocolOptions:       protOptions,
		Hostname:              *sourceAddress,
		Port:                  *sourcePort,
		SSLConfig:             ssl,
		ProxyConfig:           pc,
		AccessControl:         accessControl,
		CommandExec:           exec,
		ZeroIo:                *zeroIo,
	}
}

type cmd struct {
	stdin   io.Reader
	stdout  io.Writer
	stderr  io.Writer
	address string
	config  netcat.NetcatConfig
}

func command(stdin io.Reader, stdout io.Writer, stderr io.Writer, config netcat.NetcatConfig, args []string) (*cmd, error) {
	if len(args) < 1 {
		return nil, errMissingHostnamePort
	}

	return &cmd{
		stdin:   stdin,
		stdout:  stdout,
		stderr:  stderr,
		config:  config,
		address: args[0],
	}, nil
}

func init() {
	flag.Usage = util.Usage(flag.Usage, netcat.Usage)
}

// udpRemoteConn saves raddr from first connection and implement io.ReadWriter
// interface, so io.Copy will work

func (c *cmd) connection() (io.ReadWriter, error) {
	// switch c.network {
	// case "tcp", "tcp4", "tcp6", "unix", "unixpacket":
	// 	if c.listen {
	// 		ln, err := net.Listen(c.network, c.address)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		if c.verbose {
	// 			fmt.Fprintln(c.stderr, "Listening on", ln.Addr())
	// 		}
	// 		return ln.Accept()
	// 	}
	// 	return net.Dial(c.network, c.address)
	// case "udp", "udp4", "udp6":
	// 	addr, err := net.ResolveUDPAddr(c.network, c.address)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if c.listen {
	// 		conn, err := net.ListenUDP(c.network, addr)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		if c.verbose {
	// 			fmt.Fprintln(c.stderr, "Listening on", conn.LocalAddr())
	// 		}
	// 		waitGroup := &sync.WaitGroup{}
	// 		waitGroup.Add(1)
	// 		return &netcat.UdpRemoteConn{conn: conn, wg: waitGroup, once: &sync.Once{}, stderr: c.stderr, verbose: c.verbose}, nil
	// 	}
	// 	return net.DialUDP(c.network, nil, addr)
	// default:
	// 	return nil, fmt.Errorf("unsupported network type %q", c.network)
	// }
	return nil, nil
}

func (c *cmd) run() error {
	conn, err := c.connection()
	if err != nil {
		return err
	}

	go func() {
		if _, err := io.Copy(conn, c.stdin); err != nil {
			fmt.Fprintln(c.stderr, err)
		}
	}()

	if _, err = io.Copy(c.stdout, conn); err != nil {
		fmt.Fprintln(c.stderr, err)
	}

	if c.config.Output.Verbose {
		fmt.Fprintln(c.stderr, "Disconnected")
	}

	return nil
}

func main() {
	config := parseParams()
	fmt.Printf("Config: %v\n", config)

	c, err := command(os.Stdin, os.Stdout, os.Stderr, config, flag.Args())
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}

	if err = c.run(); err != nil {
		log.Fatalf("netcat: %v", err)
	}
}
