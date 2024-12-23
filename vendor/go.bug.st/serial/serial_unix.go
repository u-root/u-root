//
// Copyright 2014-2023 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

//go:build linux || darwin || freebsd || openbsd

package serial

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.bug.st/serial/unixutils"
	"golang.org/x/sys/unix"
)

type unixPort struct {
	handle int

	readTimeout time.Duration
	closeLock   sync.RWMutex
	closeSignal *unixutils.Pipe
	opened      uint32
}

func (port *unixPort) Close() error {
	if !atomic.CompareAndSwapUint32(&port.opened, 1, 0) {
		return nil
	}

	// Close port
	port.releaseExclusiveAccess()
	if err := unix.Close(port.handle); err != nil {
		return err
	}

	if port.closeSignal != nil {
		// Send close signal to all pending reads (if any)
		port.closeSignal.Write([]byte{0})

		// Wait for all readers to complete
		port.closeLock.Lock()
		defer port.closeLock.Unlock()

		// Close signaling pipe
		if err := port.closeSignal.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (port *unixPort) Read(p []byte) (int, error) {
	port.closeLock.RLock()
	defer port.closeLock.RUnlock()
	if atomic.LoadUint32(&port.opened) != 1 {
		return 0, &PortError{code: PortClosed}
	}

	var deadline time.Time
	if port.readTimeout != NoTimeout {
		deadline = time.Now().Add(port.readTimeout)
	}

	fds := unixutils.NewFDSet(port.handle, port.closeSignal.ReadFD())
	for {
		timeout := time.Duration(-1)
		if port.readTimeout != NoTimeout {
			timeout = time.Until(deadline)
			if timeout < 0 {
				// a negative timeout means "no-timeout" in Select(...)
				timeout = 0
			}
		}
		res, err := unixutils.Select(fds, nil, fds, timeout)
		if err == unix.EINTR {
			continue
		}
		if err != nil {
			return 0, err
		}
		if res.IsReadable(port.closeSignal.ReadFD()) {
			return 0, &PortError{code: PortClosed}
		}
		if !res.IsReadable(port.handle) {
			// Timeout happened
			return 0, nil
		}
		n, err := unix.Read(port.handle, p)
		if err == unix.EINTR {
			continue
		}
		// Linux: when the port is disconnected during a read operation
		// the port is left in a "readable with zero-length-data" state.
		// https://stackoverflow.com/a/34945814/1655275
		if n == 0 && err == nil {
			return 0, &PortError{code: PortClosed}
		}
		if n < 0 { // Do not return -1 unix errors
			n = 0
		}
		return n, err
	}
}

func (port *unixPort) Write(p []byte) (n int, err error) {
	n, err = unix.Write(port.handle, p)
	if n < 0 { // Do not return -1 unix errors
		n = 0
	}
	return
}

func (port *unixPort) Break(t time.Duration) error {
	if err := unix.IoctlSetInt(port.handle, ioctlTiocsbrk, 0); err != nil {
		return err
	}

	time.Sleep(t)

	if err := unix.IoctlSetInt(port.handle, ioctlTioccbrk, 0); err != nil {
		return err
	}

	return nil
}

func (port *unixPort) SetMode(mode *Mode) error {
	settings, err := port.getTermSettings()
	if err != nil {
		return err
	}
	if err := setTermSettingsParity(mode.Parity, settings); err != nil {
		return err
	}
	if err := setTermSettingsDataBits(mode.DataBits, settings); err != nil {
		return err
	}
	if err := setTermSettingsStopBits(mode.StopBits, settings); err != nil {
		return err
	}
	requireSpecialBaudrate := false
	if err, special := setTermSettingsBaudrate(mode.BaudRate, settings); err != nil {
		return err
	} else if special {
		requireSpecialBaudrate = true
	}
	if err := port.setTermSettings(settings); err != nil {
		return err
	}
	if requireSpecialBaudrate {
		// MacOSX require this one to be the last operation otherwise an
		// 'Invalid serial port' error is produced.
		if err := port.setSpecialBaudrate(uint32(mode.BaudRate)); err != nil {
			return err
		}
	}
	return nil
}

func (port *unixPort) SetDTR(dtr bool) error {
	status, err := port.getModemBitsStatus()
	if err != nil {
		return err
	}
	if dtr {
		status |= unix.TIOCM_DTR
	} else {
		status &^= unix.TIOCM_DTR
	}
	return port.setModemBitsStatus(status)
}

func (port *unixPort) SetRTS(rts bool) error {
	status, err := port.getModemBitsStatus()
	if err != nil {
		return err
	}
	if rts {
		status |= unix.TIOCM_RTS
	} else {
		status &^= unix.TIOCM_RTS
	}
	return port.setModemBitsStatus(status)
}

func (port *unixPort) SetReadTimeout(timeout time.Duration) error {
	if timeout < 0 && timeout != NoTimeout {
		return &PortError{code: InvalidTimeoutValue}
	}
	port.readTimeout = timeout
	return nil
}

func (port *unixPort) GetModemStatusBits() (*ModemStatusBits, error) {
	status, err := port.getModemBitsStatus()
	if err != nil {
		return nil, err
	}
	return &ModemStatusBits{
		CTS: (status & unix.TIOCM_CTS) != 0,
		DCD: (status & unix.TIOCM_CD) != 0,
		DSR: (status & unix.TIOCM_DSR) != 0,
		RI:  (status & unix.TIOCM_RI) != 0,
	}, nil
}

func nativeOpen(portName string, mode *Mode) (*unixPort, error) {
	h, err := unix.Open(portName, unix.O_RDWR|unix.O_NOCTTY|unix.O_NDELAY, 0)
	if err != nil {
		switch err {
		case unix.EBUSY:
			return nil, &PortError{code: PortBusy}
		case unix.EACCES:
			return nil, &PortError{code: PermissionDenied}
		}
		return nil, err
	}
	port := &unixPort{
		handle:      h,
		opened:      1,
		readTimeout: NoTimeout,
	}

	// Setup serial port
	settings, err := port.getTermSettings()
	if err != nil {
		port.Close()
		return nil, &PortError{code: InvalidSerialPort, causedBy: fmt.Errorf("error getting term settings: %w", err)}
	}

	// Set raw mode
	setRawMode(settings)

	// Explicitly disable RTS/CTS flow control
	setTermSettingsCtsRts(false, settings)

	if port.setTermSettings(settings) != nil {
		port.Close()
		return nil, &PortError{code: InvalidSerialPort, causedBy: fmt.Errorf("error setting term settings: %w", err)}
	}

	if mode.InitialStatusBits != nil {
		status, err := port.getModemBitsStatus()
		if err != nil {
			port.Close()
			return nil, &PortError{code: InvalidSerialPort, causedBy: fmt.Errorf("error getting modem bits status: %w", err)}
		}
		if mode.InitialStatusBits.DTR {
			status |= unix.TIOCM_DTR
		} else {
			status &^= unix.TIOCM_DTR
		}
		if mode.InitialStatusBits.RTS {
			status |= unix.TIOCM_RTS
		} else {
			status &^= unix.TIOCM_RTS
		}
		if err := port.setModemBitsStatus(status); err != nil {
			port.Close()
			return nil, &PortError{code: InvalidSerialPort, causedBy: fmt.Errorf("error setting modem bits status: %w", err)}
		}
	}

	// MacOSX require that this operation is the last one otherwise an
	// 'Invalid serial port' error is returned... don't know why...
	if err := port.SetMode(mode); err != nil {
		port.Close()
		return nil, &PortError{code: InvalidSerialPort, causedBy: fmt.Errorf("error configuring port: %w", err)}
	}

	unix.SetNonblock(h, false)

	port.acquireExclusiveAccess()

	// This pipe is used as a signal to cancel blocking Read
	pipe := &unixutils.Pipe{}
	if err := pipe.Open(); err != nil {
		port.Close()
		return nil, &PortError{code: InvalidSerialPort, causedBy: fmt.Errorf("error opening signaling pipe: %w", err)}
	}
	port.closeSignal = pipe

	return port, nil
}

func nativeGetPortsList() ([]string, error) {
	files, err := ioutil.ReadDir(devFolder)
	if err != nil {
		return nil, err
	}

	ports := make([]string, 0, len(files))
	regex, err := regexp.Compile(regexFilter)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		// Skip folders
		if f.IsDir() {
			continue
		}

		// Keep only devices with the correct name
		if !regex.MatchString(f.Name()) {
			continue
		}

		portName := devFolder + "/" + f.Name()

		// Check if serial port is real or is a placeholder serial port "ttySxx" or "ttyHSxx"
		if strings.HasPrefix(f.Name(), "ttyS") || strings.HasPrefix(f.Name(), "ttyHS") {
			port, err := nativeOpen(portName, &Mode{})
			if err != nil {
				continue
			} else {
				port.Close()
			}
		}

		// Save serial port in the resulting list
		ports = append(ports, portName)
	}

	return ports, nil
}

// termios manipulation functions

func setTermSettingsParity(parity Parity, settings *unix.Termios) error {
	switch parity {
	case NoParity:
		settings.Cflag &^= unix.PARENB
		settings.Cflag &^= unix.PARODD
		settings.Cflag &^= tcCMSPAR
		settings.Iflag &^= unix.INPCK
	case OddParity:
		settings.Cflag |= unix.PARENB
		settings.Cflag |= unix.PARODD
		settings.Cflag &^= tcCMSPAR
		settings.Iflag |= unix.INPCK
	case EvenParity:
		settings.Cflag |= unix.PARENB
		settings.Cflag &^= unix.PARODD
		settings.Cflag &^= tcCMSPAR
		settings.Iflag |= unix.INPCK
	case MarkParity:
		if tcCMSPAR == 0 {
			return &PortError{code: InvalidParity}
		}
		settings.Cflag |= unix.PARENB
		settings.Cflag |= unix.PARODD
		settings.Cflag |= tcCMSPAR
		settings.Iflag |= unix.INPCK
	case SpaceParity:
		if tcCMSPAR == 0 {
			return &PortError{code: InvalidParity}
		}
		settings.Cflag |= unix.PARENB
		settings.Cflag &^= unix.PARODD
		settings.Cflag |= tcCMSPAR
		settings.Iflag |= unix.INPCK
	default:
		return &PortError{code: InvalidParity}
	}
	return nil
}

func setTermSettingsDataBits(bits int, settings *unix.Termios) error {
	databits, ok := databitsMap[bits]
	if !ok {
		return &PortError{code: InvalidDataBits}
	}
	// Remove previous databits setting
	settings.Cflag &^= unix.CSIZE
	// Set requested databits
	settings.Cflag |= databits
	return nil
}

func setTermSettingsStopBits(bits StopBits, settings *unix.Termios) error {
	switch bits {
	case OneStopBit:
		settings.Cflag &^= unix.CSTOPB
	case OnePointFiveStopBits:
		return &PortError{code: InvalidStopBits}
	case TwoStopBits:
		settings.Cflag |= unix.CSTOPB
	default:
		return &PortError{code: InvalidStopBits}
	}
	return nil
}

func setTermSettingsCtsRts(enable bool, settings *unix.Termios) {
	if enable {
		settings.Cflag |= tcCRTSCTS
	} else {
		settings.Cflag &^= tcCRTSCTS
	}
}

func setRawMode(settings *unix.Termios) {
	// Set local mode
	settings.Cflag |= unix.CREAD
	settings.Cflag |= unix.CLOCAL

	// Set raw mode
	settings.Lflag &^= unix.ICANON
	settings.Lflag &^= unix.ECHO
	settings.Lflag &^= unix.ECHOE
	settings.Lflag &^= unix.ECHOK
	settings.Lflag &^= unix.ECHONL
	settings.Lflag &^= unix.ECHOCTL
	settings.Lflag &^= unix.ECHOPRT
	settings.Lflag &^= unix.ECHOKE
	settings.Lflag &^= unix.ISIG
	settings.Lflag &^= unix.IEXTEN

	settings.Iflag &^= unix.IXON
	settings.Iflag &^= unix.IXOFF
	settings.Iflag &^= unix.IXANY
	settings.Iflag &^= unix.INPCK
	settings.Iflag &^= unix.IGNPAR
	settings.Iflag &^= unix.PARMRK
	settings.Iflag &^= unix.ISTRIP
	settings.Iflag &^= unix.IGNBRK
	settings.Iflag &^= unix.BRKINT
	settings.Iflag &^= unix.INLCR
	settings.Iflag &^= unix.IGNCR
	settings.Iflag &^= unix.ICRNL
	settings.Iflag &^= tcIUCLC

	settings.Oflag &^= unix.OPOST

	// Block reads until at least one char is available (no timeout)
	settings.Cc[unix.VMIN] = 1
	settings.Cc[unix.VTIME] = 0
}

// native syscall wrapper functions

func (port *unixPort) getTermSettings() (*unix.Termios, error) {
	return unix.IoctlGetTermios(port.handle, ioctlTcgetattr)
}

func (port *unixPort) setTermSettings(settings *unix.Termios) error {
	return unix.IoctlSetTermios(port.handle, ioctlTcsetattr, settings)
}

func (port *unixPort) getModemBitsStatus() (int, error) {
	return unix.IoctlGetInt(port.handle, unix.TIOCMGET)
}

func (port *unixPort) setModemBitsStatus(status int) error {
	return unix.IoctlSetPointerInt(port.handle, unix.TIOCMSET, status)
}

func (port *unixPort) acquireExclusiveAccess() error {
	return unix.IoctlSetInt(port.handle, unix.TIOCEXCL, 0)
}

func (port *unixPort) releaseExclusiveAccess() error {
	return unix.IoctlSetInt(port.handle, unix.TIOCNXCL, 0)
}
