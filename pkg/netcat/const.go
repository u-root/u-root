// Copyright 2012-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netcat

import "time"

// Default values for the netcat command.
// These values were taken from the original netcat project.
const (
	DEFAULT_PORT            = 31337
	DEFAULT_SOURCE_PORT     = "31337"
	DEFAULT_IP_TYPE         = IP_V4_V6
	DEFAULT_CONNECTION_MODE = CONNECTION_MODE_CONNECT
	DEFAULT_CONNECTION_MAX  = 100
	DEFAULT_SHELL           = "/bin/sh"
	DEFAULT_UNIX_SOCKET     = "/tmp/netcat.sock"
	DEFAULT_IPV4_ADDRESS    = "0.0.0.0"
	DEFAULT_IPV6_ADDRESS    = "::"
	DEFAULT_WAIT            = time.Duration(10) * time.Second
)

var (
	DEFAULT_LF            = LINE_FEED_LF
	DEFAULT_SSL_SUITE_STR = []string{"ALL", "!aNULL", "!eNULL", "!LOW", "!EXP", "!RC4", "!MD5", "@STRENGTH"}
	LINE_FEED_LF          = []byte{0xa}        // \n
	LINE_FEED_CRLF        = []byte{0x0d, 0x0a} // \r\n
)
