// Copyright (c) 2018, Google LLC All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package mssim implements the Microsoft simulator TPM2 Transmission Interface
//
// The Microsoft simulator TPM Command Transmission Interface (TCTI) is a
// remote procedure interface donated to the TPM2 Specification by Microsoft.
// Its primary implementation is the tpm_server maintained by IBM.
//
// https://sourceforge.net/projects/ibmswtpm2/
//
// This package implements client code to communicate with server code described
// in the document "TPM2 Specification Part 4: Supporting Routines – Code"
//
// https://trustedcomputinggroup.org/wp-content/uploads/TPM-Rev-2.0-Part-4-Supporting-Routines-01.38-code.pdf
package mssim

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// Constants defined in "D.3.2. Typedefs and Defines"
const (
	tpmSignalPowerOn  uint32 = 1
	tpmSignalPowerOff uint32 = 2
	tpmSendCommand    uint32 = 8
	tpmSignalNVOn     uint32 = 11
	tpmSessionEnd     uint32 = 20
)

// Config holds configuration parameters for connecting to the simulator.
type Config struct {
	// Addresses of the command and platform handlers.
	//
	// Defaults to port 2321 and 2322 on localhost.
	CommandAddress  string
	PlatformAddress string
}

// Open creates connections to the simulator's command and platform ports and
// power cycles the simulator to initialize it.
func Open(config Config) (*Conn, error) {
	cmdAddr := config.CommandAddress
	if cmdAddr == "" {
		cmdAddr = "127.0.0.1:2321"
	}

	platformAddr := config.PlatformAddress
	if platformAddr == "" {
		platformAddr = "127.0.0.1:2322"
	}

	conn, err := net.Dial("tcp", platformAddr)
	if err != nil {
		return nil, fmt.Errorf("dial platform address: %v", err)
	}
	defer conn.Close()

	// Startup the simulator. This order of commands copies IBM's TPM2 tool's
	// "powerup" command line tool and will reset the simulator.
	//
	// https://sourceforge.net/projects/ibmtpm20tss/

	if err := sendPlatformCommand(conn, tpmSignalPowerOff); err != nil {
		return nil, fmt.Errorf("power off platform command failed: %v", err)
	}
	if err := sendPlatformCommand(conn, tpmSignalPowerOn); err != nil {
		return nil, fmt.Errorf("power on platform command failed: %v", err)
	}
	if err := sendPlatformCommand(conn, tpmSignalNVOn); err != nil {
		return nil, fmt.Errorf("nv on platform command failed: %v", err)
	}

	// Gracefully close the connection.
	if err := binary.Write(conn, binary.BigEndian, tpmSessionEnd); err != nil {
		return nil, fmt.Errorf("shutdown platform connection failed: %v", err)
	}

	cmdConn, err := net.Dial("tcp", cmdAddr)
	if err != nil {
		return nil, fmt.Errorf("dial command address: %v", err)
	}

	return &Conn{conn: cmdConn}, nil
}

// sendPlatformCommand sends device management commands to the simulator.
//
// See: "D.4.3.2. PlatformServer()"
func sendPlatformCommand(conn net.Conn, u uint32) error {
	if err := binary.Write(conn, binary.BigEndian, u); err != nil {
		return fmt.Errorf("write platform command: %v", err)
	}

	var rc uint32
	if err := binary.Read(conn, binary.BigEndian, &rc); err != nil {
		return fmt.Errorf("read platform command: %v", err)
	}
	if rc != 0 {
		return fmt.Errorf("unexpected platform command response code: 0x%x", rc)
	}
	return nil
}

// Conn is a Microsoft Simulator client that can be used as a connection for the
// tpm2 package.
type Conn struct {
	// Cached connection
	conn net.Conn

	// Response bytes left over from the previous read.
	prevRead *bytes.Reader
}

// Read a response from the simulator. If the response is longer than the provided
// buffer, the remainder will be cached for the next read.
func (c *Conn) Read(b []byte) (int, error) {
	if c.prevRead != nil && c.prevRead.Len() > 0 {
		return c.prevRead.Read(b)
	}

	// Response frame:
	// - uint32 (size of response)
	// - []byte (response)
	// - uint32 (always 0)
	var respLen uint32
	if err := binary.Read(c.conn, binary.BigEndian, &respLen); err != nil {
		return 0, fmt.Errorf("read MS simulator response header: %v", err)
	}

	resp := make([]byte, int(respLen))
	if _, err := io.ReadFull(c.conn, resp[:]); err != nil {
		return 0, fmt.Errorf("read MS simulator response: %v", err)
	}

	var rc uint32
	if err := binary.Read(c.conn, binary.BigEndian, &rc); err != nil {
		return 0, fmt.Errorf("read MS simulator return code: %v", err)
	}
	if rc != 0 {
		return 0, fmt.Errorf("MS simulator returned invalid return code: 0x%x", rc)
	}

	c.prevRead = bytes.NewReader(resp)
	return c.prevRead.Read(b)
}

// Write a raw command to the simulator. Commands must be written in a single call
// to Write. Commands split over multiple calls will result in multiple framed
// requests.
func (c *Conn) Write(b []byte) (int, error) {
	// See: D.4.3.12. TpmServer()
	buff := &bytes.Buffer{}
	// "send command" flag
	binary.Write(buff, binary.BigEndian, tpmSendCommand)
	// locality 0
	buff.WriteByte(0)
	// size of the command
	binary.Write(buff, binary.BigEndian, uint32(len(b)))
	// raw command
	buff.Write(b)

	if _, err := buff.WriteTo(c.conn); err != nil {
		return 0, fmt.Errorf("write MS simulator command: %v", err)
	}
	return len(b), nil
}

// Close closes any outgoing connections to the TPM simulator.
func (c *Conn) Close() error {
	// See: D.4.3.12. TpmServer()
	// Gracefully close the connection.
	if err := binary.Write(c.conn, binary.BigEndian, tpmSessionEnd); err != nil {
		c.conn.Close()
		return fmt.Errorf("shutdown platform connection failed: %v", err)
	}
	return c.conn.Close()
}
