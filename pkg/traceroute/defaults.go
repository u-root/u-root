// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traceroute

const (
	MAXDATALEN     = 65000
	DEFAULTDATALEN = 40

	DEFFIRSTHOP  = 1
	MAXHOPS      = 255
	DEFSIMPROBES = 16
	DEFNUMPROBES = 3

	DEFPORT    = 0
	DEFTCPPORT = 80

	DEFWAITSEC    = 5
	DEFHEREFACTOR = 3
	DEFNEARFACTOR = 10
	DEFSENDSECS   = 0

	DEFMODULE = "default"

	IPV4HdrMinLen = 20
	IPV6HdrLen    = 40

	TCPDEFPORT   = 443
	UDPDEFPORT   = 33434
	DEFNUMHOPS   = 20
	DEFNUMTRACES = 3
)
