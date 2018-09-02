// Copyright 2018 the LinuxBoot Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package uuid implements the mixed-endian UUID as implemented by Microsoft.
package uuid

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

const (
	// Size represents number of bytes in a UUID
	Size = 16
	// UExample is a example of a string UUID
	UExample  = "01234567-89AB-CDEF-0123-456789ABCDEF"
	textLen   = len(UExample)
	strFormat = "%02X%02X%02X%02X-%02X%02X-%02X%02X-%02X%02X-%02X%02X%02X%02X%02X%02X"
)

var (
	fields = [...]int{4, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1}
)

// UUID represents a unique identifier.
type UUID [Size]byte

func reverse(b []byte) {
	for i := 0; i < len(b)/2; i++ {
		other := len(b) - i - 1
		b[other], b[i] = b[i], b[other]
	}
}

// Parse parses a uuid string.
func Parse(s string) (*UUID, error) {
	// remove all hyphens to make it easier to parse.
	stripped := strings.Replace(s, "-", "", -1)
	decoded, err := hex.DecodeString(stripped)
	if err != nil {
		return nil, fmt.Errorf("uuid string not correct, need string of the format \n%v\n, got \n%v",
			UExample, s)
	}

	if len(decoded) != Size {
		return nil, fmt.Errorf("uuid string has incorrect length, need string of the format \n%v\n, got \n%v",
			UExample, s)
	}

	u := UUID{}
	i := 0
	copy(u[:], decoded[:])
	// Correct for endianness.
	for _, fieldlen := range fields {
		reverse(u[i : i+fieldlen])
		i += fieldlen
	}
	return &u, nil
}

// MustParse parses a uuid string or panics.
func MustParse(s string) *UUID {
	uuid, err := Parse(s)
	if err != nil {
		log.Fatal(err)
	}
	return uuid
}

func (u UUID) String() string {
	// Not a pointer receiver so we don't have to manually copy.
	i := 0
	// reverse endianness.
	for _, fieldlen := range fields {
		reverse(u[i : i+fieldlen])
		i += fieldlen
	}
	// Convert to []interface{} for easy printing.
	b := make([]interface{}, Size)
	for i := range u[:] {
		b[i] = u[i]
	}
	return fmt.Sprintf(strFormat, b...)
}

// MarshalJSON implements the marshaller interface.
// This allows us to actually read and edit the json file
func (u *UUID) MarshalJSON() ([]byte, error) {
	return []byte(`{"UUID" : "` + u.String() + `"}`), nil
}

// UnmarshalJSON implements the unmarshaller interface.
// This allows us to actually read and edit the json file
func (u *UUID) UnmarshalJSON(b []byte) error {
	j := make(map[string]string)

	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}
	g, err := Parse(j["UUID"])
	if err != nil {
		return err
	}
	copy(u[:], g[:])
	return nil
}
