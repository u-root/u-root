// Copyright 2012-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// netcat creates arbitrary TCP and UDP connections and listens and sends arbitrary data.
package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/netcat"
	"github.com/u-root/u-root/pkg/uroot/util"
)

var (
	verbose                 bool
	udp                     bool
	sctp                    bool
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
	proxyAuthType           string
	maxConnections          uint
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
	// protocol options
	flag.BoolVar(&ipv4, "4", false, "Use IPv4 only")
	flag.BoolVar(&ipv6, "6", false, "Use IPv6 only")
	flag.BoolVarP(&udpSocket, "udp", "u", false, "Use UDP instead of default TCP")
	flag.BoolVarP(&sctpSocket, "sctp", "", false, "Use SCTP instead of default TCP")
	flag.BoolVarP(&unixSocket, "unixsock", "U", false, "Use Unix domain sockets only")
	flag.BoolVarP(&virtualSocket, "vsock", "", false, "Use virtual circuit (stream) sockets only")

	// exec
	flag.StringVarP(&execNative, "exec", "e", "", "Executes the given command")                        // EXEC_TYPE_NATIVE
	flag.StringVarP(&execSh, "sh-exec", "c", "", "Executes the given command via /bin/sh")             // EXEC_TYPE_SHELL
	flag.StringVarP(&execLua, "lua-exec", "", "", "Executes the given Lua script (filepath argument)") // EXEC_TYPE_LUA

	// connection mode options
	flag.BoolVarP(&zeroIo, "", "z", false, "ero-I/O mode, report connection status only")
	flag.StringVarP(&sourcePort, "source-port", "p", netcat.DEFAULT_SOURCE_PORT, "Specify source port to use")
	flag.StringVarP(&sourceAddress, "source", "s", "", "Specify source address to use (doesn't affect -l)")
	flag.StringSliceVar(&looseSourceRouterPoints, "g", []string{}, "Loose source routing hop points (8 max)")
	flag.UintVar(&looseSourcePointer, "G", 4, "Loose source routing hop pointer (<n>)")

	// output options
	flag.BoolVar(&verbose, "v", false, "Set verbosity level (can not be used several times)")
	flag.StringVarP(&outFilePath, "output", "o", "", "Dump session data to a file")
	flag.StringVarP(&outFileHexPath, "hex-dump", "x", "", "Dump session data as hex to a file")
	flag.BoolVarP(&appendOutput, "append-output", "", false, "Append rather than clobber specified output files")

	// listen options
	flag.BoolVarP(&listen, "listen", "l", false, "Bind and listen for incoming connections")
	flag.UintVarP(&maxConnections, "max-conns", "m", netcat.DEFAULT_CONNECTION_MAX, "Maximum <n> simultaneous connections")
	flag.BoolVarP(&keepOpen, "keep-open", "k", false, "Accept multiple connections in listen mode")
	flag.BoolVarP(&brokerMode, "broker", "", false, "Enable Ncat's connection brokering mode")
	flag.BoolVarP(&chatMode, "chat", "", false, "Start a simple Ncat chat server")

	// timing options
	flag.StringVarP(&timingWait, "idle-timeout", "i", "0ms", "Idle read/write timeout")
	flag.StringVarP(&timingDelay, "delay", "d", "0ms", "Wait between read/writes")
	flag.StringVarP(&timingTimeout, "timeout", "w", "0ms", "Connect timeout")

	// misc options
	flag.BoolVarP(&eolCRLF, "C", "crlf", false, "Use CRLF for EOL sequence")
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
	// TODO proxy port
	flag.StringVarP(&proxyAddress, "proxy", "", "", "Specify address of host to proxy through (<addr[:port]> )")
	flag.StringVarP(&proxydns, "proxy-dns", "", "", "Specify where to resolve proxy destination")
	flag.StringVarP(&proxyType, "proxy-type", "", "", "Specify proxy type ('http', 'socks4', 'socks5')")
	flag.StringVarP(&proxyAuthType, "proxy-auth", "", "", "Authenticate with HTTP or SOCKS proxy server")

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

func evalParams() (*netcat.NetcatConfig, error) {
	var err error

	config := netcat.DefaultConfig()

	flag.Parse()

	args := flag.Args()
	if len(args) > 1 {
		config.Host = args[0]
	}

	if len(args) > 2 {
		port, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid port: %v", err)
		}
		config.Port = port
	}

	// protocol options

	// IP Type
	if ipv4 && ipv6 {
		log.Fatal("Cannot specify both IPv4 and IPv6 explicitly")
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

	// Exec commands
	execs := []string{
		execNative,
		execSh,
		execLua,
	}
	config.CommandExec, err = netcat.ParseCommands(execs)
	if err != nil {
		return nil, err
	}

	// Loose source routing
	if looseSourcePointer%4 != 0 || looseSourcePointer > 28 {
		return nil, fmt.Errorf("loose source routing hop pointer must be a multiple of 4 and less than 28")
	}

	config.ConnectionModeOptions.SourceHost = sourceAddress

	if sourcePort != "" {
		config.ConnectionModeOptions.SourcePort = sourcePort
	}

	config.ConnectionModeOptions.ZeroIO = zeroIo
	config.ConnectionModeOptions.LooseSourcePointer = looseSourcePointer
	config.ConnectionModeOptions.LooseSourceRouterPoints = looseSourceRouterPoints

	// OutputOptions
	config.Output.Verbose = verbose
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

	// timing options
	config.Timing.Delay, err = time.ParseDuration(timingDelay)
	if err != nil {
		return nil, fmt.Errorf("invalid delay: %v", err)
	}

	config.Timing.Timeout, err = time.ParseDuration(timingTimeout)
	if err != nil {
		return nil, fmt.Errorf("invalid timeout: %v", err)
	}

	config.Timing.Wait, err = time.ParseDuration(timingWait)
	if err != nil {
		return nil, fmt.Errorf("invalid wait: %v", err)
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

	config.ProxyConfig.Address = proxyAddress
	config.ProxyConfig.DNSType = netcat.ProxyDNSTypeFromString(proxydns)
	config.ProxyConfig.Type = netcat.ProxyTypeFromString(proxyType)
	config.ProxyConfig.AuthType = netcat.ProxyAuthTypeFromString(proxyAuthType)

	config.SSLConfig.Enabled = sslEnabled
	config.SSLConfig.CertFilePath = sslCertFilePath
	config.SSLConfig.KeyFilePath = sslKeyFilePath
	config.SSLConfig.VerifyTrust = sslVerifyTrust
	config.SSLConfig.TrustFilePath = sslTrustFilePath
	config.SSLConfig.Ciphers = sslCiphers
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
func (c *cmd) connection() (io.ReadWriter, error) {
	// check if SSL is available for the selected protocol if enabled
	if c.config.SSLConfig.Enabled || c.config.SSLConfig.VerifyTrust {
		switch c.config.ProtocolOptions.SocketType {
		case netcat.SOCKET_TYPE_UDP:
			return nil, fmt.Errorf("SSL is not supported for UDP connections")
		case netcat.SOCKET_TYPE_UNIX:
			return nil, fmt.Errorf("SSL is not supported for UNIX connections")
		case netcat.SOCKET_TYPE_UDP_UNIX:
			return nil, fmt.Errorf("SSL is not supported for UNIX connections")
		case netcat.SOCKET_TYPE_UDP_VSOCK:
			return nil, fmt.Errorf("SSL is not supported for VSOCK connections")
		case netcat.SOCKET_TYPE_VSOCK:
			return nil, fmt.Errorf("SSL is not supported for VSOCK connections")
		case netcat.SOCKET_TYPE_SCTP:
			return nil, fmt.Errorf("SSL is not supported for SCTP connections")
		}
	}

	network, err := c.config.Network()
	if err != nil {
		return nil, fmt.Errorf("connection: %v", err)
	}

	address, err := c.config.Address()
	if err != nil {
		return nil, fmt.Errorf("connection: %v", err)
	}

	if c.config.ConnectionMode == netcat.CONNECTION_MODE_LISTEN {
		return c.listenMode(network, address)
	}
	return c.connectMode(network, address)
}

func (c *cmd) listenMode(network, address string) (io.ReadWriter, error) {
	var (
		err      error
		listener net.Listener
	)

	// If listing mode and Zero-I/O mode are combined the program will block indefinitely
	if c.config.ConnectionModeOptions.ZeroIO {
		for {
			time.Sleep(1 * time.Hour)
		}
	}

	if c.config.Misc.NoDNS {
		return nil, fmt.Errorf("listen: disabling DNS resolution is not supported in listen mode")
	}

	if c.config.ConnectionModeOptions.SourceHost != "" && c.config.ConnectionModeOptions.SourcePort != "" {
		return nil, fmt.Errorf("listen: source host/port cannot be set in listen mode")
	}

	switch c.config.ProtocolOptions.SocketType {
	case netcat.SOCKET_TYPE_TCP:
		if c.config.SSLConfig.Enabled || c.config.SSLConfig.VerifyTrust {
			tlsConfig, err := c.generateTLSConfiguration()
			if err != nil {
				return nil, fmt.Errorf("connection: %v", err)
			}

			listener, err = tls.Listen(network, address, tlsConfig)
			if err != nil {
				return nil, fmt.Errorf("connection: %v", err)
			}

		} else {
			listener, err = net.Listen(network, address)
			if err != nil {
				return nil, err
			}
		}

	case netcat.SOCKET_TYPE_UDP:
		return netcat.NewUDPRemoteConn(network, address, c.stderr, c.config.AccessControl, c.config.Output.Verbose)

	case netcat.SOCKET_TYPE_UNIX:
		listener, err = net.Listen(network, address)
		if err != nil {
			return nil, err
		}

	case netcat.SOCKET_TYPE_UDP_UNIX:
		return netcat.NewUnixgramRemoteConn(network, address, c.stderr, c.config.AccessControl, c.config.Output.Verbose)

	// unsupported socket types
	case netcat.SOCKET_TYPE_SCTP, netcat.SOCKET_TYPE_VSOCK, netcat.SOCKET_TYPE_UDP_VSOCK:
		return nil, fmt.Errorf("currently unsupported socket type %q", c.config.ProtocolOptions.SocketType)

	case netcat.SOCKET_TYPE_NONE:
	default:
		return nil, fmt.Errorf("undefined socket type %q", c.config.ProtocolOptions.SocketType)
	}

	return c.waitOnConnection(listener)
}

func (c *cmd) connectMode(network, address string) (io.ReadWriter, error) {
	var (
		err  error
		conn net.Conn
	)

	dialer := &net.Dialer{
		Timeout: c.config.Timing.Wait,
	}

	switch c.config.ProtocolOptions.SocketType {

	case netcat.SOCKET_TYPE_TCP:
		dialer.LocalAddr, err = net.ResolveTCPAddr(c.config.ConnectionModeOptions.SourceHost, c.config.ConnectionModeOptions.SourcePort)
		if err != nil {
			return nil, fmt.Errorf("connection: failed to resolve source address %v", err)
		}

	case netcat.SOCKET_TYPE_UDP:
		dialer.LocalAddr, err = net.ResolveUDPAddr(c.config.ConnectionModeOptions.SourceHost, c.config.ConnectionModeOptions.SourcePort)
		if err != nil {
			return nil, fmt.Errorf("connection: failed to resolve source address %v", err)
		}

	case netcat.SOCKET_TYPE_UNIX:
		dialer.LocalAddr, err = net.ResolveUnixAddr(c.config.ConnectionModeOptions.SourceHost, c.config.ConnectionModeOptions.SourcePort)
		if err != nil {
			return nil, fmt.Errorf("connection: failed to resolve source address %v", err)
		}

	case netcat.SOCKET_TYPE_UDP_UNIX:
		dialer.LocalAddr, err = net.ResolveUnixAddr(c.config.ConnectionModeOptions.SourceHost, c.config.ConnectionModeOptions.SourcePort)
		if err != nil {
			return nil, fmt.Errorf("connection: failed to resolve source address %v", err)
		}

	// unsupported socket types
	case netcat.SOCKET_TYPE_SCTP, netcat.SOCKET_TYPE_VSOCK, netcat.SOCKET_TYPE_UDP_VSOCK:
		return nil, fmt.Errorf("currently unsupported socket type %q", c.config.ProtocolOptions.SocketType)

	case netcat.SOCKET_TYPE_NONE:
	default:
		return nil, fmt.Errorf("undefined socket type %q", c.config.ProtocolOptions.SocketType)
	}

	// TLS Support
	if c.config.SSLConfig.Enabled || c.config.SSLConfig.VerifyTrust {
		tlsConfig, err := c.generateTLSConfiguration()
		if err != nil {
			return nil, fmt.Errorf("connection: %v", err)
		}

		conn, err = tls.DialWithDialer(dialer, network, address, tlsConfig)
		if err != nil {
			return nil, fmt.Errorf("connection: %v", err)
		}
	} else {
		conn, err = dialer.Dial(network, address)
		if err != nil {
			return nil, fmt.Errorf("connection: %v", err)
		}
	}

	if c.config.Timing.Timeout > 0 {
		conn.SetDeadline(time.Now().Add(c.config.Timing.Timeout))
	}

	return conn, nil
}

// waitOnConnection listens for incoming connections and returns the first connection that is allowed by the access control list.
func (c *cmd) waitOnConnection(listener net.Listener) (net.Conn, error) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			return nil, err
		}

		var remoteAddr string
		switch c.config.ProtocolOptions.SocketType {
		case netcat.SOCKET_TYPE_TCP:
			remoteAddr = conn.RemoteAddr().(*net.TCPAddr).IP.String()
		case netcat.SOCKET_TYPE_UDP:
			remoteAddr = conn.RemoteAddr().(*net.UDPAddr).String()
		case netcat.SOCKET_TYPE_UNIX:
			remoteAddr = conn.RemoteAddr().(*net.UnixAddr).String()
		default:
			return nil, fmt.Errorf("undefined socket type %q", c.config.ProtocolOptions.SocketType)
		}

		if c.config.AccessControl.IsAllowed(remoteAddr) {
			return conn, nil
		}

		conn.Close()
		time.Sleep(50 * time.Millisecond)
	}
}

func (c *cmd) generateTLSConfiguration() (*tls.Config, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: !c.config.SSLConfig.VerifyTrust,
	}

	if c.config.SSLConfig.CertFilePath == "" && c.config.SSLConfig.KeyFilePath != "" || c.config.SSLConfig.CertFilePath != "" && c.config.SSLConfig.KeyFilePath == "" {
		return nil, fmt.Errorf("both  certificate and key file must be provided")
	}

	if c.config.SSLConfig.CertFilePath == "" || c.config.SSLConfig.KeyFilePath == "" {
		cer, err := tls.LoadX509KeyPair(c.config.SSLConfig.CertFilePath, c.config.SSLConfig.KeyFilePath)
		if err != nil {
			return nil, fmt.Errorf("connection: %v", err)
		}

		tlsConfig.Certificates = []tls.Certificate{cer}
	}

	if c.config.SSLConfig.VerifyTrust {
		caCert, err := os.ReadFile(c.config.SSLConfig.TrustFilePath)
		if err != nil {
			return nil, fmt.Errorf("cannot read CA certificate: %v", err)
		}
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("cannot append CA certificate to pool")
		}

		tlsConfig.RootCAs = caCertPool
	}

	if c.config.SSLConfig.SNI != "" {
		tlsConfig.ServerName = c.config.SSLConfig.SNI
	}

	if c.config.SSLConfig.ALPN != nil {
		tlsConfig.NextProtos = c.config.SSLConfig.ALPN
	}

	if len(c.config.SSLConfig.Ciphers) > 0 {
		// Set the cipher suites
	}

	return tlsConfig, nil
}

func (c *cmd) run() error {
	conn, err := c.connection()
	if err != nil {
		return fmt.Errorf("run: %v", err)
	}

	// Return immediately if Zero-I/O mode is enabled and connection is established
	if c.config.ConnectionModeOptions.ZeroIO {
		return nil
	}

	// io.Copy will block until the connection is closed, use a MultiWriter to write to stdout and the output file
	combinedOut := io.MultiWriter(c.stdout, &c.config.Output)

	var wg sync.WaitGroup

	if !c.config.Misc.ReceiveOnly && c.config.ConnectionMode != netcat.CONNECTION_MODE_LISTEN {
		wg.Add(1)

		go func() {
			defer wg.Done()
			c.writeToRemote(conn)
		}()

		// prepare command execution on the server
		if c.config.CommandExec.Type != netcat.EXEC_TYPE_NONE {
			if err := c.config.CommandExec.Execute(conn, io.MultiWriter(conn, combinedOut), c.stderr, c.config.Misc.EOL); err != nil {
				return fmt.Errorf("run command: %v", err)
			}
		}
	}

	// in send-only mode ignore incoming data
	if c.config.Misc.SendOnly {
		return nil
	}

	// read from the connection
	if _, err := io.Copy(combinedOut, conn); err != nil {
		return fmt.Errorf("run dump: %v", err)
	}

	wg.Wait()
	return nil
}

func (c *cmd) writeToRemote(conn io.Writer) {
	eolReader := netcat.NewEOLReader(c.stdin, c.config.Misc.EOL)

	// If the delay is set, read the input line by line in time intervals of the delay duration
	if c.config.Timing.Delay > 0 {
		scanner := bufio.NewScanner(eolReader)

		for scanner.Scan() {
			if _, err := conn.Write([]byte(scanner.Text() + "\n")); err != nil {
				netcat.FLogf(c.config, c.stderr, "run copy: %v", err)
			}

			time.Sleep(c.config.Timing.Delay)
		}

		if err := scanner.Err(); err != nil {
			netcat.FLogf(c.config, c.stderr, "run copy: %v", err)
		}
	} else {
		if _, err := io.Copy(conn, eolReader); err != nil {
			netcat.FLogf(c.config, c.stderr, "run copy: %v", err)
		}
	}

	// do not shutdown the connection if the no-shutdown flag is set
	if c.config.Misc.NoShutdown {
		for {
			time.Sleep(1 * time.Hour)
		}
	}
}

func main() {
	config, err := evalParams()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	c, err := command(os.Stdin, os.Stdout, os.Stderr, config, flag.Args())
	if err != nil {
		fmt.Printf("error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	if err = c.run(); err != nil {
		log.Fatalf("netcat: %v", err)
	}
}
