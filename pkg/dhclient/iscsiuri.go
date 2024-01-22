// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dhclient

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type iscsiURIParser struct {
	rem     []byte
	state   string
	fieldNo iscsiField
	err     error
}

func (i *iscsiURIParser) tok(toks string, err error) []byte {
	n := bytes.Index(i.rem, []byte(toks))
	if n == -1 {
		i.err = err
		return nil
	}

	b := i.rem[:n]
	i.rem = i.rem[n+len(toks):]
	return b
}

func (i *iscsiURIParser) tokAny(toks string, err error) (data []byte, tok byte) {
	n := bytes.IndexAny(i.rem, toks)
	if n == -1 {
		i.err = err
		return nil, 0
	}

	b := i.rem[:n]
	tok = i.rem[n]
	i.rem = i.rem[n+1:]
	return b, tok
}

type iscsiField int

const (
	iscsiMagic  iscsiField = 0
	serverField iscsiField = 1
	protField   iscsiField = 2
	portField   iscsiField = 3
	lunField    iscsiField = 4
	volumeField iscsiField = 5
)

// Format:
//
// iscsi:@"<servername>":"<protocol>":"<port>":"<LUN>":"<targetname>"
//
// @ for now will be ignored. Eventually we would want complete support.
// iscsi:[<username>:<password>[:<reverse>:<password>]@]"<servername>":"<protocol>":"<port>"[:[<iscsi_iface_name>]:[<netdev_name>]]:"<LUN>":"<targetname>"
// "<servername>" may contain an IPv6 address enclosed with [] with an
// arbitrary but bounded number of colons.
//
// "<targetname>" may contain an arbitrary string with an arbitrary number of
// colons.
func ParseISCSIURI(s string) (*net.TCPAddr, string, error) {
	var (
		// port has a default value according to RFC 4173.
		port   = 3260
		ip     net.IP
		volume string
		magic  string
	)
	i := &iscsiURIParser{
		state:   "normal",
		fieldNo: iscsiMagic,
		rem:     []byte(s),
	}
	for i.fieldNo <= volumeField && i.err == nil {
		fno, tok := i.next()
		switch fno {
		case iscsiMagic:
			magic = tok
		case serverField:
			tok = strings.TrimPrefix(tok, "@") // ignore any leading @
			ip = net.ParseIP(tok)
		case protField, lunField:
			// yeah whatever
			continue
		case portField:
			if len(tok) > 0 {
				pv, err := strconv.Atoi(tok)
				if err != nil {
					return nil, "", fmt.Errorf("iSCSI URI %q has invalid port: %w", s, err)
				}
				port = pv
			}
		case volumeField:
			volume = tok
		}
	}
	if i.err != nil {
		return nil, "", fmt.Errorf("iSCSI URI %q failed to parse: %w", s, i.err)
	}
	if magic != "iscsi" {
		return nil, "", fmt.Errorf("iSCSI URI %q is missing iscsi scheme prefix, have %s", s, magic)
	}
	if len(volume) == 0 {
		return nil, "", fmt.Errorf("iSCSI URI %q is missing a volume name", s)
	}
	return &net.TCPAddr{
		IP:   ip,
		Port: port,
	}, volume, nil
}

func (i *iscsiURIParser) next() (iscsiField, string) {
	var val []byte
	switch i.state {
	case "normal":
		val = i.tok(":", fmt.Errorf("fields missing"))
		switch i.fieldNo {
		case iscsiMagic:
			// The next field is an IP or hostname, which may contain other colons.
			i.state = "ip"
		case lunField:
			// The next field is the last field, the volume name,
			// is a free-for-all that may contain as many colons as
			// it wants.
			i.state = "remaining"
		}

	case "ipv6":
		val = i.tok("]:", fmt.Errorf("invalid IPv6 address"))
		i.state = "normal"

	case "remaining":
		val = i.rem
		i.rem = nil

	case "ip":
		var tok byte
		val, tok = i.tokAny("[:", fmt.Errorf("fields missing"))
		switch tok {
		case '[':
			i.state = "ipv6"
			return i.next()
		case ':':
			// IPv4 address is in tok, go back to normal next.
			i.state = "normal"
		}

	default:
		i.err = fmt.Errorf("unrecognized state %s", i.state)
		return -1, ""
	}

	i.fieldNo++
	return i.fieldNo - 1, string(val)
}
