// Copyright (C) 2017 Kale Blankenship. All rights reserved.
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details

package tftp // import "pack.ag/tftp"

import (
	"errors"
	"fmt"
)

var (
	// errBlockSequnce is a sentinel error used internally, never returned to API clients.
	errBlockSequence = errors.New("block sequence error")
	// ErrInvalidURL indicates that the URL passed to Get or Put is invalid.
	ErrInvalidURL = errors.New("invalid URL")
	// ErrInvalidHostIP indicates an empty or invalid host.
	ErrInvalidHostIP = errors.New("invalid host/IP")
	// ErrInvalidFile indicates an empty or invalid file.
	ErrInvalidFile = errors.New("invalid file")
	// ErrSizeNotReceived indicates tsize was not negotiated.
	ErrSizeNotReceived = errors.New("size not received")
	// ErrAddressNotAvailable indicates the server address was requested before
	// the server had been started.
	ErrAddressNotAvailable = errors.New("address not available until server has been started")
	// ErrNoRegisteredHandlers indicates no handlers were registered before starting the server.
	ErrNoRegisteredHandlers = errors.New("no handlers registered")
	// ErrInvalidNetwork indicates that a network other than udp, udp4, or udp6 was configured.
	ErrInvalidNetwork = errors.New("invalid network: must be udp, udp4, or udp6")
	// ErrInvalidBlocksize indicates that a blocksize outside the range 8 to 65464 was configured.
	ErrInvalidBlocksize = errors.New("invalid blocksize: must be between 8 and 65464")
	// ErrInvalidTimeout indicates that a timeout outside the range 1 to 255 was configured.
	ErrInvalidTimeout = errors.New("invalid timeout: must be between 1 and 255")
	// ErrInvalidWindowsize indicates that a windowsize outside the range 1 to 65535 was configured.
	ErrInvalidWindowsize = errors.New("invalid windowsize: must be between 1 and 65535")
	// ErrInvalidMode indicates that a mode other than ModeNetASCII or ModeOctet was configured.
	ErrInvalidMode = errors.New("invalid transfer mode: must be ModeNetASCII or ModeOctet")
	// ErrInvalidRetransmit indicates that the retransmit limit was configured with a negative value.
	ErrInvalidRetransmit = errors.New("invalid retransmit: cannot be negative")
	// ErrMaxRetries indicates that the maximum number of retries has been reached.
	ErrMaxRetries = errors.New("max retries reached")
)

type errUnexpectedDatagram struct {
	dg string //datagram string
}

func (e *errUnexpectedDatagram) Error() string {
	return fmt.Sprintf("unexpected datagram: %s", e.dg)
}

// IsUnexpectedDatagram allows a consumer to check if an error
// is an unexpected datagram.
func IsUnexpectedDatagram(err error) bool {
	err = ErrorCause(err)
	_, ok := err.(*errUnexpectedDatagram)
	return ok
}

type errRemoteError struct {
	dg string
}

func (e *errRemoteError) Error() string {
	return "remote error: " + e.dg
}

// IsRemoteError allows a consumer to check if an error
// was an error by the remote client/server.
func IsRemoteError(err error) bool {
	err = ErrorCause(err)
	_, ok := err.(*errRemoteError)
	return ok
}

type errParsingOption struct {
	option string
	value  string
}

func (e *errParsingOption) Error() string {
	return fmt.Sprintf("error parsing %q for option %q", e.value, e.option)
}

// IsOptionParsingError allows a consumer to check if an error
// was induced during option parsing.
func IsOptionParsingError(err error) bool {
	err = ErrorCause(err)
	_, ok := err.(*errParsingOption)
	return ok
}

// tftpError wraps an error with a context message and is itself and error.
type tftpError struct {
	orig error
	msg  string
}

func (e *tftpError) Error() string {
	return e.msg + ": " + e.orig.Error()
}

// wrapError wraps an error with a contextual message.
//
// This is a simplistic version of github.com/pkg/errors
func wrapError(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &tftpError{orig: err, msg: msg}
}

// ErrorCause extracts the original error from an error wrapped by tftp.
func ErrorCause(err error) error {
	for err != nil {
		tftperr, ok := err.(*tftpError)
		if !ok {
			break
		}
		err = tftperr.orig
	}
	return err
}
