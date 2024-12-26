//
// Copyright 2014-2023 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package serial

import "time"

//go:generate go run golang.org/x/sys/windows/mkwinsyscall -output zsyscall_windows.go syscall_windows.go

// Port is the interface for a serial Port
type Port interface {
	// SetMode sets all parameters of the serial port
	SetMode(mode *Mode) error

	// Stores data received from the serial port into the provided byte array
	// buffer. The function returns the number of bytes read.
	//
	// The Read function blocks until (at least) one byte is received from
	// the serial port or an error occurs.
	Read(p []byte) (n int, err error)

	// Send the content of the data byte array to the serial port.
	// Returns the number of bytes written.
	Write(p []byte) (n int, err error)

	// Wait until all data in the buffer are sent
	Drain() error

	// ResetInputBuffer Purges port read buffer
	ResetInputBuffer() error

	// ResetOutputBuffer Purges port write buffer
	ResetOutputBuffer() error

	// SetDTR sets the modem status bit DataTerminalReady
	SetDTR(dtr bool) error

	// SetRTS sets the modem status bit RequestToSend
	SetRTS(rts bool) error

	// GetModemStatusBits returns a ModemStatusBits structure containing the
	// modem status bits for the serial port (CTS, DSR, etc...)
	GetModemStatusBits() (*ModemStatusBits, error)

	// SetReadTimeout sets the timeout for the Read operation or use serial.NoTimeout
	// to disable read timeout.
	SetReadTimeout(t time.Duration) error

	// Close the serial port
	Close() error

	// Break sends a break for a determined time
	Break(time.Duration) error
}

// NoTimeout should be used as a parameter to SetReadTimeout to disable timeout.
var NoTimeout time.Duration = -1

// ModemStatusBits contains all the modem input status bits for a serial port (CTS, DSR, etc...).
// It can be retrieved with the Port.GetModemStatusBits() method.
type ModemStatusBits struct {
	CTS bool // ClearToSend status
	DSR bool // DataSetReady status
	RI  bool // RingIndicator status
	DCD bool // DataCarrierDetect status
}

// ModemOutputBits contains all the modem output bits for a serial port.
// This is used in the Mode.InitialStatusBits struct to specify the initial status of the bits.
// Note: Linux and MacOSX (and basically all unix-based systems) can not set the status bits
// before opening the port, even if the initial state of the bit is set to false they will go
// anyway to true for a few milliseconds, resulting in a small pulse.
type ModemOutputBits struct {
	RTS bool // ReadyToSend status
	DTR bool // DataTerminalReady status
}

// Open opens the serial port using the specified modes
func Open(portName string, mode *Mode) (Port, error) {
	port, err := nativeOpen(portName, mode)
	if err != nil {
		// Return a nil interface, for which var==nil is true (instead of
		// a nil pointer to a struct that satisfies the interface).
		return nil, err
	}
	return port, err
}

// GetPortsList retrieve the list of available serial ports
func GetPortsList() ([]string, error) {
	return nativeGetPortsList()
}

// Mode describes a serial port configuration.
type Mode struct {
	BaudRate          int              // The serial port bitrate (aka Baudrate)
	DataBits          int              // Size of the character (must be 5, 6, 7 or 8)
	Parity            Parity           // Parity (see Parity type for more info)
	StopBits          StopBits         // Stop bits (see StopBits type for more info)
	InitialStatusBits *ModemOutputBits // Initial output modem bits status (if nil defaults to DTR=true and RTS=true)
}

// Parity describes a serial port parity setting
type Parity int

const (
	// NoParity disable parity control (default)
	NoParity Parity = iota
	// OddParity enable odd-parity check
	OddParity
	// EvenParity enable even-parity check
	EvenParity
	// MarkParity enable mark-parity (always 1) check
	MarkParity
	// SpaceParity enable space-parity (always 0) check
	SpaceParity
)

// StopBits describe a serial port stop bits setting
type StopBits int

const (
	// OneStopBit sets 1 stop bit (default)
	OneStopBit StopBits = iota
	// OnePointFiveStopBits sets 1.5 stop bits
	OnePointFiveStopBits
	// TwoStopBits sets 2 stop bits
	TwoStopBits
)

// PortError is a platform independent error type for serial ports
type PortError struct {
	code     PortErrorCode
	causedBy error
}

// PortErrorCode is a code to easily identify the type of error
type PortErrorCode int

const (
	// PortBusy the serial port is already in used by another process
	PortBusy PortErrorCode = iota
	// PortNotFound the requested port doesn't exist
	PortNotFound
	// InvalidSerialPort the requested port is not a serial port
	InvalidSerialPort
	// PermissionDenied the user doesn't have enough priviledges
	PermissionDenied
	// InvalidSpeed the requested speed is not valid or not supported
	InvalidSpeed
	// InvalidDataBits the number of data bits is not valid or not supported
	InvalidDataBits
	// InvalidParity the selected parity is not valid or not supported
	InvalidParity
	// InvalidStopBits the selected number of stop bits is not valid or not supported
	InvalidStopBits
	// InvalidTimeoutValue the timeout value is not valid or not supported
	InvalidTimeoutValue
	// ErrorEnumeratingPorts an error occurred while listing serial port
	ErrorEnumeratingPorts
	// PortClosed the port has been closed while the operation is in progress
	PortClosed
	// FunctionNotImplemented the requested function is not implemented
	FunctionNotImplemented
)

// EncodedErrorString returns a string explaining the error code
func (e PortError) EncodedErrorString() string {
	switch e.code {
	case PortBusy:
		return "Serial port busy"
	case PortNotFound:
		return "Serial port not found"
	case InvalidSerialPort:
		return "Invalid serial port"
	case PermissionDenied:
		return "Permission denied"
	case InvalidSpeed:
		return "Port speed invalid or not supported"
	case InvalidDataBits:
		return "Port data bits invalid or not supported"
	case InvalidParity:
		return "Port parity invalid or not supported"
	case InvalidStopBits:
		return "Port stop bits invalid or not supported"
	case InvalidTimeoutValue:
		return "Timeout value invalid or not supported"
	case ErrorEnumeratingPorts:
		return "Could not enumerate serial ports"
	case PortClosed:
		return "Port has been closed"
	case FunctionNotImplemented:
		return "Function not implemented"
	default:
		return "Other error"
	}
}

// Error returns the complete error code with details on the cause of the error
func (e PortError) Error() string {
	if e.causedBy != nil {
		return e.EncodedErrorString() + ": " + e.causedBy.Error()
	}
	return e.EncodedErrorString()
}

// Code returns an identifier for the kind of error occurred
func (e PortError) Code() PortErrorCode {
	return e.code
}
