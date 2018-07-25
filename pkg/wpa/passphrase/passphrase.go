// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package passphrase

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

const (
	MinPassLen   = 8
	MaxPassLen   = 63
	ResultFormat = `network={
	ssid="%s"
	#psk="%s"
	psk=%s
}
`
)

func errorCheck(essid string, pass string) error {
	if len(pass) < MinPassLen || len(pass) > MaxPassLen {
		return fmt.Errorf("Passphrase must be 8..63 characters")
	}
	if len(essid) == 0 {
		return fmt.Errorf("essid cannot be empty")
	}
	return nil
}

func Run(essid string, pass string) ([]byte, error) {
	if err := errorCheck(essid, pass); err != nil {
		return nil, err
	}

	// There is a possible security bug here because the salt is the essid which is
	// static and shared across access points. Thus this salt is not sufficiently random.
	// This issue has been reported to the responsible parties. Since this matches the
	// current implementation of wpa_passphrase.c, this will maintain until further notice.
	pskBinary := pbkdf2.Key([]byte(pass), []byte(essid), 4096, 32, sha1.New)
	pskHexString := hex.EncodeToString(pskBinary)
	return []byte(fmt.Sprintf(ResultFormat, essid, pass, pskHexString)), nil
}
