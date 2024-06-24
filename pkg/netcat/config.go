// Copyright 2012-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netcat

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Default values for the netcat command.
func DefaultConfig() NetcatConfig {
	return NetcatConfig{
		Port:           DEFAULT_PORT,
		ConnectionMode: DEFAULT_CONNECTION_MODE,
		ConnectionModeOptions: NetcatConnectModeOptions{
			SourceHost: DEFAULT_IPV4_ADDRESS,
			SourcePort: DEFAULT_SOURCE_PORT,
		},
		ListenModeOptions: NetcatListenModeOptions{
			MaxConnections: DEFAULT_CONNECTION_MAX,
		},
		ProtocolOptions: NetcatProtocolOptions{
			IPType:     DEFAULT_IP_TYPE,
			SocketType: SOCKET_TYPE_TCP,
		},
		SSLConfig: NetcatSSLOptions{
			Ciphers: []string{DEFAULT_SSL_SUITE_STR},
		},
		ProxyConfig: NetcatProxyOptions{
			Type:     PROXY_TYPE_NONE,
			DNSType:  PROXY_DNS_NONE,
			AuthType: PROXY_AUTH_NONE,
		},
		AccessControl: NetcatAccessControlOptions{},
		CommandExec: NetcatExec{
			Type: EXEC_TYPE_NONE,
		},
		Output: NetcatOutputOptions{},
		Misc: NetcatMiscOptions{
			EOL: DEFAULT_LF,
		},
		Timing: NetcatTimingOptions{
			Wait: DEFAULT_WAIT,
		},
	}
}

type NetcatIPType int

const (
	IP_NONE NetcatIPType = iota
	IP_V4
	IP_V6
	IP_V4_V6
	IP_V4_STRICT
	IP_V6_STRICT
)

type NetcatSocketType int

// UDP can be combined with one of {UNIX, VSOCK} or stand alone
const (
	SOCKET_TYPE_TCP NetcatSocketType = iota
	SOCKET_TYPE_UDP
	SOCKET_TYPE_UNIX
	SOCKET_TYPE_VSOCK
	SOCKET_TYPE_SCTP
	SOCKET_TYPE_UDP_VSOCK
	SOCKET_TYPE_UDP_UNIX
	SOCKET_TYPE_NONE
)

func (s NetcatSocketType) String() string {
	return [...]string{
		"tcp",
		"udp",
		"unix",
		"vsock",
		"sctp",
		"udp-vsock",
		"unixgram",
		"none",
	}[s]
}

func (n *NetcatConfig) Network() (string, error) {
	switch n.ProtocolOptions.SocketType {
	case SOCKET_TYPE_TCP:
		switch n.ProtocolOptions.IPType {
		case IP_V4:
			fallthrough
		case IP_V4_STRICT:
			return "tcp4", nil
		case IP_V6_STRICT:
			fallthrough
		case IP_V6:
			return "tcp6", nil
		default:
			return "tcp", nil
		}

	case SOCKET_TYPE_UDP:
		switch n.ProtocolOptions.IPType {
		case IP_V4:
			fallthrough
		case IP_V4_STRICT:
			return "udp4", nil
		case IP_V6_STRICT:
			fallthrough
		case IP_V6:
			return "udp6", nil
		default:
			return "udp", nil
		}

	case SOCKET_TYPE_UNIX:
		return "unix", nil

	case SOCKET_TYPE_UDP_UNIX:
		return "unixgram", nil

	// VSOCK connections don't require a network specification
	case SOCKET_TYPE_VSOCK, SOCKET_TYPE_UDP_VSOCK, SOCKET_TYPE_SCTP:
		return "", nil
	}

	return "", fmt.Errorf("invalid/unimplemented combination of socket and ip type (%v - %v)", n.ProtocolOptions.SocketType, n.ProtocolOptions.IPType)
}

// TODO: should we do this via enabled or can we just pass nil?
// NOTE: The key files have to be of PEM format

type ProxyType int

const (
	PROXY_TYPE_NONE ProxyType = iota
	PROXY_TYPE_HTTP
	PROXY_TYPE_SOCKS4
	PROXY_TYPE_SOCKS5
)

func (p ProxyType) String() string {
	return [...]string{
		"None",
		"HTTP",
		"SOCKS4",
		"SOCKS5",
	}[p]
}

type ProxyPortError struct {
	ProxyType ProxyType
}

func (e *ProxyPortError) Error() string {
	return fmt.Sprintf("ProxyType %s has no default port", ProxyType.String(e.ProxyType))
}

func (p ProxyType) DefaultPort() (uint, error) {
	switch p {
	case PROXY_TYPE_SOCKS5:
		return 1080, nil
	case PROXY_TYPE_HTTP:
		return 3128, nil
	default:
		return 0, &ProxyPortError{ProxyType: p}
	}
}

func ProxyTypeFromString(s string) ProxyType {
	switch strings.ToUpper(s) {
	case "HTTP":
		return PROXY_TYPE_HTTP
	case "SOCKS4":
		return PROXY_TYPE_SOCKS4
	case "SOCKS5":
		return PROXY_TYPE_SOCKS5
	default:
		return PROXY_TYPE_NONE
	}
}

type NetcatConnectTypeOptions struct {
	LooseSourceRouterPoints []string
}

type ProxyAuthType int

const (
	PROXY_AUTH_NONE ProxyAuthType = iota
	PROXY_AUTH_HTTP
	PROXY_AUTH_SOCKS5
)

func ProxyAuthTypeFromString(s string) ProxyAuthType {
	switch strings.ToUpper(s) {
	case "HTTP":
		return PROXY_AUTH_HTTP
	case "SOCKS5":
		return PROXY_AUTH_SOCKS5
	default:
		return PROXY_AUTH_NONE
	}
}

type ProxyDNSType int

const (
	PROXY_DNS_NONE ProxyDNSType = iota
	PROXY_DNS_LOCAL
	PROXY_DNS_REMOTE
	PROXY_DNS_BOTH
)

func ProxyDNSTypeFromString(s string) ProxyDNSType {
	switch strings.ToUpper(s) {
	case "LOCAL":
		return PROXY_DNS_LOCAL
	case "REMOTE":
		return PROXY_DNS_REMOTE
	case "BOTH":
		return PROXY_DNS_BOTH
	default:
		return PROXY_DNS_NONE
	}
}

type NetcatProxyOptions struct {
	Type     ProxyType // If this is none, discard the entire Proxy handling
	Address  string
	DNSType  ProxyDNSType
	Port     uint
	AuthType ProxyAuthType // If this is none, discard the entire ProxyAuth handling
}

type NetcatSSLOptions struct {
	// In connect mode, this option transparently negotiates an SSL session
	// In server mode, this option listens for incoming SSL connections
	// Depending on the Protocol Type, either TLS (TCP) od DTLS (UDP) will be used
	Enabled bool

	CertFilePath  string // Path to the certificate file in PEM format
	KeyFilePath   string // Path to the private key file in PEM format
	VerifyTrust   bool   // In client mode is like --ssl except that it also requires verification of the server certificate. No effect in server mode.
	TrustFilePath string // Verify trust and domain name of certificates

	// List of ciphersuites that Ncat will use when connecting to servers or when accepting SSL connections from clients
	// Syntax is described in the OpenSSL ciphers(1) man page
	Ciphers []string
	SNI     string   // (Server Name Indication) Tell the server the name of the logical server Ncat is contacting
	ALPN    []string // List of protocols to send via the Application-Layer Protocol Negotiation
}

type NetcatConnectionMode int

// Ncat operates in one of two primary modes: connect mode and listen mode.
const (
	CONNECTION_MODE_CONNECT NetcatConnectionMode = iota
	CONNECTION_MODE_LISTEN
)

func (n NetcatConnectionMode) String() string {
	return [...]string{
		"connect",
		"listen",
	}[n]
}

type NetcatProtocolOptions struct {
	IPType     NetcatIPType
	SocketType NetcatSocketType
}

func ParseSocketType(udp, unix, vsock, sctp bool) (NetcatSocketType, error) {
	// tcp ^ (udp || udp && unix || udp && vsock) ^ unix ^ vsock ^ sctp
	if !(udp || unix || vsock || sctp) {
		return SOCKET_TYPE_TCP, nil
	}

	if udp && !(unix || vsock || sctp) {
		return SOCKET_TYPE_UDP, nil
	}

	if udp && (unix != vsock) && !sctp {
		if unix {
			return SOCKET_TYPE_UDP_UNIX, nil
		} else if vsock {
			return SOCKET_TYPE_UDP_VSOCK, nil
		}
	}

	if unix && !(udp || vsock || sctp) {
		return SOCKET_TYPE_UNIX, nil
	}

	if vsock && !(udp || unix || sctp) {
		return SOCKET_TYPE_VSOCK, nil
	}

	if sctp && !(udp || unix || vsock) {
		return SOCKET_TYPE_SCTP, nil
	}

	return SOCKET_TYPE_NONE, fmt.Errorf("invalid socket type combination")
}

// Holds allowed and denied hosts for the connection.
// The map is a string representation of the host is allowed or denied.
// The key is the host and the value is true for allowed and false for denied.
// If the host is not in the map, it is allowed by default.
// The map is populated by the access cli flags, with access flags taking precedence over deny flags.
type NetcatAccessControlOptions struct {
	SetAllowed     bool
	ConnectionList map[string]bool
}

// Check if the hostName is allowed to connect.
func (ac *NetcatAccessControlOptions) IsAllowed(hostNames []string) bool {
	// If atleast one item is in the allowed list, the hostName is only allowed if also mentioned.
	if ac.SetAllowed {
		for _, hostName := range hostNames {
			allowed, ok := ac.ConnectionList[hostName]
			if !ok {
				// If the hostname is not in the list, check the next hostname.
				continue
			}

			// If the host name is allowed, return true.
			return allowed
		}

		// If the host name is not in the list, it is denied by default.
		return false
	}

	// If the host is part of the list and not one entry is set to allowed, it is denied by default.
	for _, hostName := range hostNames {
		// Otherwise check if the host name is denied.
		if _, ok := ac.ConnectionList[hostName]; ok {
			return false
		}
	}

	// If the host is not in the list, it is allowed by default.
	return true
}

func ParseAccessControl(connectionAllowFile string, connectionAllowList []string, connectionDenyFile string, connectionDenyList []string) (NetcatAccessControlOptions, error) {
	accessControl := NetcatAccessControlOptions{}

	accessControl.ConnectionList = make(map[string]bool)

	if connectionDenyFile != "" {
		denyFile, err := os.Open(connectionDenyFile)
		if err != nil {
			log.Fatal(err)
		}
		defer denyFile.Close()

		scanner := bufio.NewScanner(denyFile)
		for scanner.Scan() {
			line := scanner.Text()
			accessControl.ConnectionList[line] = false
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	for _, deniedLine := range connectionDenyList {
		accessControl.ConnectionList[deniedLine] = false
	}

	// allowed hosts are parsed secondly as they take precedence
	if len(connectionAllowList) > 0 || connectionAllowFile != "" {
		accessControl.SetAllowed = true
	}

	if connectionAllowFile != "" {
		allowFile, err := os.Open(connectionAllowFile)
		if err != nil {
			log.Fatal(err)
		}
		defer allowFile.Close()

		scanner := bufio.NewScanner(allowFile)
		for scanner.Scan() {
			line := scanner.Text()
			accessControl.ConnectionList[line] = true
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	for _, allowedLine := range connectionAllowList {
		accessControl.ConnectionList[allowedLine] = true
	}

	return accessControl, nil
}

type NetcatExecType int

const (
	EXEC_TYPE_NATIVE NetcatExecType = iota
	EXEC_TYPE_SHELL
	EXEC_TYPE_LUA
	EXEC_TYPE_NONE // For faster case switching this is appended at the end
)

type NetcatExec struct {
	Type    NetcatExecType
	Command string
}

func ParseCommands(commands []string) (NetcatExec, error) {
	cmds := 0
	last_valid := -1
	for i, e := range commands {
		if e != "" {
			cmds++
			last_valid = i
		}
	}

	// This is a recoverable error, we can just ignore the command
	if last_valid == -1 {
		return NetcatExec{}, nil
	}

	if cmds > 1 {
		return NetcatExec{}, fmt.Errorf("only one of --exec, --sh-exec, and --lua-exec is allowed")
	}

	return NetcatExec{
		Type:    NetcatExecType(last_valid),
		Command: commands[last_valid],
	}, nil
}

// Execute a given command on the host system
// stdout of the command is send to to the connection
// stderr of the command is displayed on stdout of the host
// The host process exits with the exit code of the command unless --keep-open is specified
func (n *NetcatExec) Execute(stdin io.ReadWriter, stdout io.Writer, stderr io.Writer, eol []byte) error {
	var (
		cmd    *exec.Cmd
		buffer bytes.Buffer
	)

	switch n.Type {
	case EXEC_TYPE_NATIVE:
		cmd = exec.Command(n.Command)
	case EXEC_TYPE_SHELL:
		cmd = exec.Command(DEFAULT_SHELL, "-c", n.Command)
	case EXEC_TYPE_LUA:
		return fmt.Errorf("not implemented")
	default:
		return nil
	}

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	scanner := bufio.NewScanner(stdin)
	for scanner.Scan() {
		buffer.WriteString(scanner.Text())
		buffer.Write(eol)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	cmd.Stdin = &buffer

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("exec start: %w", err)
	}

	// Wait waits for the command to exit and waits for any copying to stdin or copying from stdout or stderr to complete.
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("exec wait: %w", err)
	}

	return nil
}

// Since Ncat can be either in connect mode or in listen mode, only one of these
// options structs should be filled at a time.
// TODO: Make this a generic ConnectionModeOptions struct that may have different fields
type NetcatConnectModeOptions struct {
	LooseSourceRouterPoints []string // IPV4_STRICT, Point as IP or Hostname
	LooseSourcePointer      uint     // The argument must be a multiple of 4 and no more than 28
	SourceHost              string   // Address for Ncat to bind to
	SourcePort              string   // Port number for Ncat to bind to
	ZeroIO                  bool     // restrict IO, only report connection
}

// All the modes can be combined with each other, no need to check for mutual exclusivity
type NetcatListenModeOptions struct {
	MaxConnections uint // Maximum number of simultaneous connections accepted by an Ncat instance
	KeepOpen       bool // Accept multiple connections. In this mode there is no way for Ncat to know when its network input is finished, so it will keep running until interrupted
	BrokerMode     bool // Don't echo messages but redirect them to all others. Compatible with other modes
	ChatMode       bool // Enables chat mode, intended for the exchange of text between several users
}

type NetcatTimingOptions struct {
	Delay   time.Duration // Delay interval for lines sent
	Timeout time.Duration // If the idle timeout is reached, the connection is terminated.
	Wait    time.Duration // Fixed timeout for connection attempts.
}

type NetcatOutputOptions struct {
	OutFilePath     string     // Dump session data to a file
	OutFileMutex    sync.Mutex // Mutex for the file
	OutFileHexPath  string     // Dump session data in hex to a file
	OutFileHexMutex sync.Mutex // Mutex for the hex file
	AppendOutput    bool       // Append the resulted output rather than truncating
	Verbose         bool       // TODO: make this adjustable level with -v..v
}

// Write writes the data to the file specified in the options
// If Netcat is not configured to write to a file, it will return 0, nil
// https://go.dev/src/io/io.go
func (n *NetcatOutputOptions) Write(data []byte) (int, error) {
	fileOpts := os.O_CREATE | os.O_WRONLY
	if n.AppendOutput {
		fileOpts |= os.O_APPEND
	}

	if n.OutFilePath != "" {
		n.OutFileMutex.Lock()
		f, err := os.OpenFile(n.OutFilePath, fileOpts, 0o644)
		if err != nil {
			n.OutFileMutex.Unlock()
			return 0, fmt.Errorf("netcat oo open: %w", err)
		}

		_, err = f.Write(data)
		if err != nil {
			n.OutFileMutex.Unlock()
			return 0, fmt.Errorf("netcat oo write: %w", err)
		}
		n.OutFileMutex.Unlock()
	}

	if n.OutFileHexPath != "" {
		n.OutFileHexMutex.Lock()

		f, err := os.OpenFile(n.OutFileHexPath, fileOpts, 0o644)
		if err != nil {
			n.OutFileHexMutex.Unlock()
			return 0, fmt.Errorf("netcat outopt open: %w", err)
		}

		_, err = f.Write([]byte(hex.Dump(data)))
		if err != nil {
			n.OutFileHexMutex.Unlock()
			return 0, fmt.Errorf("netcat outopt write: %w", err)
		}

		n.OutFileHexMutex.Unlock()
	}

	return len(data), nil
}

type NetcatMiscOptions struct {
	EOL         []byte // This can either be a single byte or a sequence of bytes
	ReceiveOnly bool   // Only receive data and will not try to send anything.
	SendOnly    bool   // Ncat will only send data and will ignore anything received
	NoShutdown  bool   // Req: half-duplex mode, Ncat will not invoke shutdown on a socket after seeing EOF on stdin
	NoDNS       bool   // Completely disable hostname resolution across all Ncat options (destination, source address, source routing hops, proxy)
	Telnet      bool   // Handle DO/DONT WILL/WONT Telnet negotiations.
}

// Omit version and help here since they are handled by the flag parser
type NetcatConfig struct {
	Host                  string
	Port                  uint64
	ConnectionMode        NetcatConnectionMode
	ConnectionModeOptions NetcatConnectModeOptions
	ListenModeOptions     NetcatListenModeOptions
	ProtocolOptions       NetcatProtocolOptions
	SSLConfig             NetcatSSLOptions
	ProxyConfig           NetcatProxyOptions
	AccessControl         NetcatAccessControlOptions
	CommandExec           NetcatExec
	Output                NetcatOutputOptions
	Misc                  NetcatMiscOptions
	Timing                NetcatTimingOptions
}

// Return a default address for the provided configuration
func (n *NetcatConfig) Address() (string, error) {
	switch n.ConnectionMode {
	case CONNECTION_MODE_CONNECT:
		if n.Host == "" {
			return "", fmt.Errorf("missing host")
		}

		switch n.ProtocolOptions.SocketType {
		case SOCKET_TYPE_UNIX, SOCKET_TYPE_UDP_UNIX:
			return n.Host, nil

		// unimplemented
		case SOCKET_TYPE_VSOCK, SOCKET_TYPE_UDP_VSOCK, SOCKET_TYPE_SCTP:
			return "nil", fmt.Errorf("currently unsupported socket type %v", n.ProtocolOptions.SocketType)
		default:
			if n.Misc.NoDNS {
				if ip := net.ParseIP(n.Host); ip == nil {
					return "", fmt.Errorf("non-numerical host but DNS resolution is disabled: %v", n.Host)
				}
			}
			return n.Host + ":" + strconv.FormatUint(uint64(n.Port), 10), nil
		}
	case CONNECTION_MODE_LISTEN:
		var address string

		if n.Host != "" {
			address = n.Host
		} else {
			switch n.ProtocolOptions.SocketType {
			case SOCKET_TYPE_TCP | SOCKET_TYPE_UDP:
				switch n.ProtocolOptions.IPType {
				case IP_V4:
					address = DEFAULT_IPV4_ADDRESS
				case IP_V6:
					address = DEFAULT_IPV6_ADDRESS
				}

			case SOCKET_TYPE_UNIX:
				return DEFAULT_UNIX_SOCKET, nil

			// unimplemented
			case SOCKET_TYPE_VSOCK, SOCKET_TYPE_UDP_VSOCK, SOCKET_TYPE_SCTP:
				return "nil", fmt.Errorf("currently unsupported socket type %v", n.ProtocolOptions.SocketType)
			default:
				return "", fmt.Errorf("invalid socket type %v", n.ProtocolOptions.SocketType)
			}
		}

		return address + ":" + strconv.FormatUint(uint64(n.Port), 10), nil
	default:
		return "", fmt.Errorf("invalid connection mode %v", n.ConnectionMode)
	}
}
