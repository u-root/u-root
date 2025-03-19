// Copyright 2012-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netcat

import (
	"fmt"
	"net"
	"strconv"
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
			Ciphers: []string{},
		},
		ProxyConfig: ProxyOptions{
			Type:    DEFAULT_PROXY_TYPE,
			DNSType: PROXY_DNS_NONE,
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

// Since Ncat can be either in connect mode or in listen mode, only one of these
// options structs should be filled at a time.
// TODO: Make this a generic ConnectionModeOptions struct that may have different fields
type ConnectModeOptions struct {
	ScanPorts               bool
	CurrentPort             uint64
	EndPort                 uint64
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
		case SOCKET_TYPE_UDP_VSOCK:
			return "nil", fmt.Errorf("currently unsupported socket type %v", c.ProtocolOptions.SocketType)
		default:
			if c.Misc.NoDNS {
				if ip := net.ParseIP(c.Host); ip == nil {
					return "", fmt.Errorf("non-numerical host but DNS resolution is disabled: %v", c.Host)
				}
			}

			if c.ConnectionModeOptions.CurrentPort != 0 {
				port := c.ConnectionModeOptions.CurrentPort
				return net.JoinHostPort(c.Host, strconv.FormatUint(port, 10)), nil
			}

			return net.JoinHostPort(c.Host, strconv.FormatUint(c.Port, 10)), nil
		}
	case CONNECTION_MODE_LISTEN:
		var address string

		if c.Host != "" {
			address = c.Host
		} else {
			switch c.ProtocolOptions.SocketType {
			case SOCKET_TYPE_TCP, SOCKET_TYPE_UDP:
				switch c.ProtocolOptions.IPType {
				case IP_V6, IP_V6_STRICT:
					address = DEFAULT_IPV6_ADDRESS
				default:
					address = DEFAULT_IPV4_ADDRESS
				}

			case SOCKET_TYPE_UNIX, SOCKET_TYPE_UDP_UNIX:
				return DEFAULT_UNIX_SOCKET, nil

			// unimplemented
			case SOCKET_TYPE_VSOCK, SOCKET_TYPE_UDP_VSOCK:
				return "nil", fmt.Errorf("currently unsupported socket type %v", c.ProtocolOptions.SocketType)
			default:
				return "", fmt.Errorf("invalid socket type %v", c.ProtocolOptions.SocketType)
			}
		}

		if c.ProtocolOptions.SocketType == SOCKET_TYPE_UNIX || c.ProtocolOptions.SocketType == SOCKET_TYPE_UDP_UNIX {
			return address, nil
		}
		return net.JoinHostPort(address, strconv.FormatUint(c.Port, 10)), nil
	default:
		return "", fmt.Errorf("invalid connection mode %v", c.ConnectionMode)
	}
}
