// Copyright 2012-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netcat

const (
	DEFAULT_PORT            = 31337
	DEFAULT_IP_TYPE         = IP_V4_V6
	DEFAULT_CONNECTION_MODE = CONNECTION_MODE_CONNECT
	DEFAULT_CONNECTION_MAX  = 100
	DEFAULT_SSL_SUITE_STR   = "ALL:!aNULL:!eNULL:!LOW:!EXP:!RC4:!MD5:@STRENGTH"
)

var (
	DEFAULT_LF     = LINE_FEED_LF
	LINE_FEED_LF   = []byte{0xa}        // \n
	LINE_FEED_CRLF = []byte{0x0d, 0x0a} // \r\n
)

const (
	Usage      = "netcat [go-style network address]"
	LOG_PREFIX = "netcat: "
)
