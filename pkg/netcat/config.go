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
	"sync"
	"time"

	"github.com/u-root/u-root/pkg/ulog"
)

// Default values for the netcat command.
func DefaultConfig() Config {
	return Config{
		Port:           DEFAULT_PORT,
		ConnectionMode: DEFAULT_CONNECTION_MODE,
		ConnectionModeOptions: ConnectModeOptions{
			LooseSourceRouterPoints: make([]string, 0),
			SourcePort:              DEFAULT_SOURCE_PORT,
		},
		ListenModeOptions: ListenModeOptions{
			MaxConnections: DEFAULT_CONNECTION_MAX,
		},
		ProtocolOptions: ProtocolOptions{
			IPType:     DEFAULT_IP_TYPE,
			SocketType: SOCKET_TYPE_TCP,
		},
		SSLConfig: SSLOptions{
			Ciphers: DEFAULT_SSL_SUITE_STR,
		},
		ProxyConfig: ProxyOptions{
			Type:     PROXY_TYPE_NONE,
			DNSType:  PROXY_DNS_NONE,
			AuthType: PROXY_AUTH_NONE,
		},
		AccessControl: AccessControlOptions{
			ConnectionList: make(map[string]bool),
		},
		CommandExec: Exec{
			Type: EXEC_TYPE_NONE,
		},
		Output: OutputOptions{
			Logger: ulog.Null,
		},
		Misc: MiscOptions{
			EOL: DEFAULT_LF,
		},
		Timing: TimingOptions{
			Wait: DEFAULT_WAIT,
		},
	}
}

type SSLOptions struct {
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

type ConnectionMode int

// Ncat operates in one of two primary modes: connect mode and listen mode.
const (
	CONNECTION_MODE_CONNECT ConnectionMode = iota
	CONNECTION_MODE_LISTEN
)

func (c ConnectionMode) String() string {
	return [...]string{
		"connect",
		"listen",
	}[c]
}

// Holds allowed and denied hosts for the connection.
// The map is a string representation of the host is allowed or denied.
// The key is the host and the value is true for allowed and false for denied.
// If the host is not in the map, it is allowed by default.
// The map is populated by the access cli flags, with access flags taking precedence over deny flags.
type AccessControlOptions struct {
	SetAllowed     bool
	ConnectionList map[string]bool
}

// Check if the hostName is allowed to connect.
func (ac *AccessControlOptions) IsAllowed(hostNames []string) bool {
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

func ParseAccessControl(connectionAllowFile string, connectionAllowList []string, connectionDenyFile string, connectionDenyList []string) (AccessControlOptions, error) {
	accessControl := AccessControlOptions{}

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

type ExecType int

const (
	EXEC_TYPE_NATIVE ExecType = iota
	EXEC_TYPE_SHELL
	EXEC_TYPE_LUA
	EXEC_TYPE_NONE // For faster case switching this is appended at the end
)

type Exec struct {
	Type    ExecType
	Command string
}

func ParseCommands(commands []string) (Exec, error) {
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
		return Exec{}, nil
	}

	if cmds > 1 {
		return Exec{}, fmt.Errorf("only one of --exec, --sh-exec, and --lua-exec is allowed")
	}

	return Exec{
		Type:    ExecType(last_valid),
		Command: commands[last_valid],
	}, nil
}

// Execute a given command on the host system
// stdout of the command is send to to the connection
// stderr of the command is displayed on stdout of the host
// The host process exits with the exit code of the command unless --keep-open is specified
func (n *Exec) Execute(stdin io.ReadWriter, stdout io.Writer, stderr io.Writer, eol []byte) error {
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
type ConnectModeOptions struct {
	LooseSourceRouterPoints []string // IPV4_STRICT, Point as IP or Hostname
	LooseSourcePointer      uint     // The argument must be a multiple of 4 and no more than 28
	SourceHost              string   // Address for Ncat to bind to
	SourcePort              string   // Port number for Ncat to bind to
	ZeroIO                  bool     // restrict IO, only report connection
}

// All the modes can be combined with each other, no need to check for mutual exclusivity
type ListenModeOptions struct {
	MaxConnections uint32 // Maximum number of simultaneous connections accepted by an Ncat instance
	KeepOpen       bool   // Accept multiple connections. In this mode there is no way for Ncat to know when its network input is finished, so it will keep running until interrupted
	BrokerMode     bool   // Don't echo messages but redirect them to all others. Compatible with other modes
	ChatMode       bool   // Enables chat mode, intended for the exchange of text between several users
}

type TimingOptions struct {
	Delay   time.Duration // Delay interval for lines sent
	Timeout time.Duration // If the idle timeout is reached, the connection is terminated.
	Wait    time.Duration // Fixed timeout for connection attempts.
}

type OutputOptions struct {
	OutFilePath     string      // Dump session data to a file
	OutFileMutex    sync.Mutex  // Mutex for the file
	OutFileHexPath  string      // Dump session data in hex to a file
	OutFileHexMutex sync.Mutex  // Mutex for the hex file
	AppendOutput    bool        // Append the resulted output rather than truncating
	Logger          ulog.Logger // Verbose output
}

// Write writes the data to the file specified in the options
// If Netcat is not configured to write to a file, it will return 0, nil
// https://go.dev/src/io/io.go
func (n *OutputOptions) Write(data []byte) (int, error) {
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

type MiscOptions struct {
	EOL         []byte // This can either be a single byte or a sequence of bytes
	ReceiveOnly bool   // Only receive data and will not try to send anything.
	SendOnly    bool   // Ncat will only send data and will ignore anything received
	NoShutdown  bool   // Req: half-duplex mode, Ncat will not invoke shutdown on a socket after seeing EOF on stdin
	NoDNS       bool   // Completely disable hostname resolution across all Ncat options (destination, source address, source routing hops, proxy)
	Telnet      bool   // Handle DO/DONT WILL/WONT Telnet negotiations.
}

// Omit version and help here since they are handled by the flag parser
type Config struct {
	Host                  string
	Port                  uint64
	ConnectionMode        ConnectionMode
	ConnectionModeOptions ConnectModeOptions
	ListenModeOptions     ListenModeOptions
	ProtocolOptions       ProtocolOptions
	SSLConfig             SSLOptions
	ProxyConfig           ProxyOptions
	AccessControl         AccessControlOptions
	CommandExec           Exec
	Output                OutputOptions
	Misc                  MiscOptions
	Timing                TimingOptions
}

// Return a default address for the provided configuration
func (c *Config) Address() (string, error) {
	switch c.ConnectionMode {
	case CONNECTION_MODE_CONNECT:
		if c.Host == "" {
			return "", fmt.Errorf("missing host")
		}

		switch c.ProtocolOptions.SocketType {
		case SOCKET_TYPE_UNIX, SOCKET_TYPE_UDP_UNIX:
			return c.Host, nil

		// unimplemented
		case SOCKET_TYPE_VSOCK, SOCKET_TYPE_UDP_VSOCK, SOCKET_TYPE_SCTP:
			return "nil", fmt.Errorf("currently unsupported socket type %v", c.ProtocolOptions.SocketType)
		default:
			if c.Misc.NoDNS {
				if ip := net.ParseIP(c.Host); ip == nil {
					return "", fmt.Errorf("non-numerical host but DNS resolution is disabled: %v", c.Host)
				}
			}
			return c.Host + ":" + strconv.FormatUint(uint64(c.Port), 10), nil
		}
	case CONNECTION_MODE_LISTEN:
		var address string

		if c.Host != "" {
			address = c.Host
		} else {
			switch c.ProtocolOptions.SocketType {
			case SOCKET_TYPE_TCP | SOCKET_TYPE_UDP:
				switch c.ProtocolOptions.IPType {
				case IP_V4:
					address = DEFAULT_IPV4_ADDRESS
				case IP_V6:
					address = DEFAULT_IPV6_ADDRESS
				}

			case SOCKET_TYPE_UNIX:
				return DEFAULT_UNIX_SOCKET, nil

			// unimplemented
			case SOCKET_TYPE_VSOCK, SOCKET_TYPE_UDP_VSOCK, SOCKET_TYPE_SCTP:
				return "nil", fmt.Errorf("currently unsupported socket type %v", c.ProtocolOptions.SocketType)
			default:
				return "", fmt.Errorf("invalid socket type %v", c.ProtocolOptions.SocketType)
			}
		}

		return address + ":" + strconv.FormatUint(uint64(c.Port), 10), nil
	default:
		return "", fmt.Errorf("invalid connection mode %v", c.ConnectionMode)
	}
}
