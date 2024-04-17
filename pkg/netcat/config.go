// Copyright 2012-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netcat

import (
	"fmt"
	"log"
	"os"
	"strings"
)

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

// TODO: combine ipv4 and ipv6 within that
func (s NetcatSocketType) ToGoType(i NetcatIPType) (string, error) {
	switch i {
	case IP_V4:
		if s == SOCKET_TYPE_TCP {
			return "tcp4", nil
		}
		if s == SOCKET_TYPE_UDP {
			return "udp4", nil
		}
	case IP_V6:
		if s == SOCKET_TYPE_TCP {
			return "tcp6", nil
		}
		if s == SOCKET_TYPE_UDP {
			return "udp6", nil
		}
	case IP_V4_V6:
		if s == SOCKET_TYPE_TCP {
			return "tcp", nil
		}
		if s == SOCKET_TYPE_UDP {
			return "udp", nil
		}
	}
	return "", fmt.Errorf("invalid/unimplemented combination of socket and ip type")
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

type NetcatProxyConfig struct {
	Type     ProxyType // If this is none, discard the entire Proxy handling
	Address  string
	DNSType  ProxyDNSType
	Port     uint
	AuthType ProxyAuthType // If this is none, discard the entire ProxyAuth handling
}

type NetcatSSLConfig struct {
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

func (config *NetcatSSLConfig) Verify() error {
	// If it is enabled but any other value is set it's invalid
	if config.Enabled && (config.CertFilePath != "" ||
		config.KeyFilePath != "" ||
		config.VerifyTrust ||
		config.TrustFilePath != "" ||
		len(config.Ciphers) > 0 && (len(config.Ciphers) == 1 && config.Ciphers[0] != DEFAULT_SSL_SUITE_STR) ||
		config.SNI != "" ||
		len(config.ALPN) > 0) {
		return fmt.Errorf("ssl is enabled but other ssl options are set")
	}

	// check if provided files exist
	if config.CertFilePath != "" {
		if _, err := os.Stat(config.CertFilePath); err != nil {
			return fmt.Errorf("certificate file does not exist")
		}
	}

	if config.KeyFilePath != "" {
		if _, err := os.Stat(config.KeyFilePath); err != nil {
			return fmt.Errorf("key file does not exist")
		}
	}

	return nil
}

// Verify checks if the configuration is valid and can be used
func (config NetcatConfig) Verify() error {
	return nil
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

// TODO: review policy if we want to merge cli and file or if they are mutual exclusive
// For now this is mutual exclusive
type NetcatAccessControlOptions struct {
	ConnectionAllowList []string
	ConnectionDenyList  []string
}

func ParseAccessControl(connectionAllowFile *string, connectionAllowList *[]string, connectionDenyFile *string, connectionDenyList *[]string) (NetcatAccessControlOptions, error) {
	accessControl := NetcatAccessControlOptions{}

	// Allowlist
	if *connectionAllowFile != "" && *connectionAllowList != nil {
		log.Fatal("Cannot specify both allowlist and allowfile")
	}
	if *connectionAllowFile != "" {
		data, err := os.ReadFile(*connectionAllowFile)
		if err != nil {
			log.Fatal(err)

		}
		accessControl.ConnectionAllowList = strings.Split(string(data), ",")
	}
	if *connectionAllowList != nil {
		accessControl.ConnectionAllowList = *connectionAllowList
	}

	// Denylist
	if *connectionDenyFile != "" && *connectionDenyList != nil {
		log.Fatal("Cannot specify both denylist and denyfile")
	}
	if *connectionDenyFile != "" {
		data, err := os.ReadFile(*connectionDenyFile)
		if err != nil {
			log.Fatal(err)

		}
		accessControl.ConnectionDenyList = strings.Split(string(data), ",")
	}
	if *connectionDenyList != nil {
		accessControl.ConnectionDenyList = *connectionDenyList
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

func ParseCommands(commands []*string) (NetcatExec, error) {
	cmds := 0
	last_valid := -1
	for i, e := range commands {
		if *e != "" {
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
		Command: *commands[last_valid],
	}, nil
}

// Since Ncat can be either in connect mode or in listen mode, only one of these
// options structs should be filled at a time.
// TODO: Make this a generic ConnectionModeOptions struct that may have different fields
type NetcatConnectModeOptions struct {
	LooseSourceRouterPoints []string // IPV4_STRICT, Point as IP or Hostname
	LooseSourcePointer      uint     // The argument must be a multiple of 4 and no more than 28.
	Host                    string   // Address for Ncat to bind to
	Port                    uint     // Port number for Ncat to bind to.
}

// All the modes can be combined with each other, no need to check for mutual exclusivity
type NetcatListenModeOptions struct {
	MaxConnections uint // Maximum number of simultaneous connections accepted by an Ncat instance
	KeepOpen       bool // Accept multiple connections. In this mode there is no way for Ncat to know when its network input is finished, so it will keep running until interrupted
	BrokerMode     bool // Don't echo messages but redirect them to all others. Compatible with other modes
	ChatMode       bool // Enables chat mode, intended for the exchange of text between several users
}

type NetcatTimingOptions struct {
	Delay   uint // Delay interval for lines sent
	Timeout uint // If the idle timeout is reached, the connection is terminated.
	Wait    uint // Fixed timeout for connection attempts.
}

type NetcatOutputOptions struct {
	OutFilePath    string // Dump session data to a file
	OutFileHexPath string //  Dump session data in hex to a file.
	AppendOutput   bool   // Append the resulted output rather than truncating
	Verbose        bool   // TOOD: make this adjustable level with -v..v
}

type NetcatMiscOptions struct {
	EOL         []byte // This can either be a single byte or a sequence of bytes
	ReceiveOnly bool   // Only receive data and will not try to send anything.
	SendOnly    bool   // Ncat will only send data and will ignore anything received
	NoShutdown  bool   // Req: half-duplex mode, Ncat will not invoke shutdown on a socket after seeing EOF on stdin
	NoDns       bool   // Completely disable hostname resolution across all Ncat options (destination, source address, source routing hops, proxy)
	Telnet      bool   // Handle DO/DONT WILL/WONT Telnet negotiations.
}

// Omit version and help here since they are handled by the flag parser
type NetcatConfig struct {
	ConnectionMode        NetcatConnectionMode
	ConnectionModeOptions NetcatConnectModeOptions
	ListenModeOptions     NetcatListenModeOptions
	ProtocolOptions       NetcatProtocolOptions
	SocketType            NetcatSocketType
	SSLConfig             NetcatSSLConfig
	ProxyConfig           NetcatProxyConfig
	AccessControl         NetcatAccessControlOptions
	CommandExec           NetcatExec
	Output                NetcatOutputOptions
	Misc                  NetcatMiscOptions
	Timing                NetcatTimingOptions
	Hostname              string
	Port                  uint
	ZeroIo                bool
}
