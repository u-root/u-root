// Copyright 2014-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package libinit

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

// WaitOrphans waits for all remaining processes on the system to exit.
func WaitOrphans() uint {
	var numReaped uint
	for {
		var (
			s unix.WaitStatus
			r unix.Rusage
		)
		p, err := unix.Wait4(-1, &s, 0, &r)
		if p == -1 {
			break
		}
		log.Printf("%v: exited with %v, status %v, rusage %v", p, err, s, r)
		numReaped++
	}
	return numReaped
}

// WithTTYControl turns on controlling the TTY on this command.
func WithTTYControl(ctty bool) CommandModifier {
	return func(c *exec.Cmd) {
		if c.SysProcAttr == nil {
			c.SysProcAttr = &unix.SysProcAttr{}
		}
		c.SysProcAttr.Setctty = ctty
		c.SysProcAttr.Setsid = ctty
	}
}

func WithMultiTTY(mtty bool, openFn func([]string) ([]*os.File, error), ttyNames []string) CommandModifier {
	return func(c *exec.Cmd) {
		if mtty {
			ww, err := openFn(ttyNames)
			if err != nil {
				log.Printf("%q: open devices for multi-TTY output: %v", c.Path, err)
				log.Printf("falling back to default stdout and stderr")
				return
			}

			if len(ww) >= 1 {
				writers := make([]io.Writer, len(ww))
				for i, w := range ww {
					writers[i] = w
				}
				c.Stdout = io.MultiWriter(writers...)
				c.Stderr = io.MultiWriter(writers...)

				// Save this for later use
				for i := 0; i < len(ww); i++ {
					c.Env = append(c.Env, fmt.Sprintf("tty%d=%s", i, ww[i].Name()))
				}
			}
		}
	}
}

// WithCloneFlags adds clone(2) flags to the *exec.Cmd.
func WithCloneFlags(flags uintptr) CommandModifier {
	return func(c *exec.Cmd) {
		if c.SysProcAttr == nil {
			c.SysProcAttr = &unix.SysProcAttr{}
		}
		c.SysProcAttr.Cloneflags = flags
	}
}

func init() {
	osDefault = linuxDefault
}

func linuxDefault(c *exec.Cmd) {
	c.SysProcAttr = &unix.SysProcAttr{
		Setctty: true,
		Setsid:  true,
	}
}

func Raw(r *os.File) error {
	termios, err := unix.IoctlGetTermios(int(r.Fd()), unix.TCGETS)
	if err != nil {
		return err
	}

	termios.Iflag &^= unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP | unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON
	termios.Oflag &^= unix.OPOST
	termios.Lflag &^= unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.IEXTEN
	termios.Cflag &^= unix.CSIZE | unix.PARENB
	termios.Cflag |= unix.CS8
	termios.Cc[unix.VMIN] = 1
	termios.Cc[unix.VTIME] = 0

	if err = unix.IoctlSetTermios(int(r.Fd()), unix.TCSETS, termios); err != nil {
		return err
	}
	if err = syscall.SetNonblock(int(r.Fd()), true); err != nil {
		return err
	}
	return nil
}

func NewPTMS() (*os.File, *os.File, error) {
	p, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	if err != nil {
		time.Sleep(1 * time.Second)
		return nil, nil, err
	}

	// unlock
	var u int32
	// use TIOCSPTLCK with a pointer to zero to clear the lock.
	err = ioctl(p, syscall.TIOCSPTLCK, uintptr(unsafe.Pointer(&u))) //nolint:gosec // Expected unsafe pointer for Syscall call.
	if err != nil {
		return nil, nil, err
	}

	sname, err := ptsname(p)
	if err != nil {
		return nil, nil, err
	}

	t, err := os.OpenFile(sname, os.O_RDWR|syscall.O_NOCTTY|syscall.O_NONBLOCK, 0o620) //nolint:gosec // Expected Open from a variable.
	if err != nil {
		return nil, nil, err
	}

	return p, t, nil
}

// FIX ME: make it not linux-specific
// RunCommands runs commands in sequence.
//
// RunCommands returns how many commands existed and were attempted to run.
//
// commands must refer to absolute paths at the moment.
func RunCommands(debug func(string, ...interface{}), commands ...*exec.Cmd) int {
	var cmdCount int
	for _, cmd := range commands {
		if _, err := os.Stat(cmd.Path); os.IsNotExist(err) {
			debug("%v", err)
			continue
		}

		cmdCount++
		debug("Trying to run %v", cmd)

		// Set up PTM
		m, s, err := NewPTMS()
		if err != nil {
			log.Printf("Error getting PTY: %v", err)
			return 0
		}

		cmd.Stdin = s

		// Launch go routine to copy output from the command to the PTM
		for _, r := range cmd.Env {
			if strings.HasPrefix(r, "tty") {
				tty := strings.Split(r, "=")[1]

				debug("Opening TTY %v", tty)

				// Open the TTY
				t, err := os.OpenFile(tty, os.O_RDWR, 0)
				if err != nil {
					log.Printf("Error opening TTY: %v", err)
					continue
				}

				// Set to Raw mode
				if err := Raw(t); err != nil {
					// Let's still continue if we can't set the TTY to raw mode
					log.Printf("Error setting TTY to raw mode: %v", err)
				}

				go func() {
					for {
						_, err := io.Copy(m, t)
						if err != nil {
							log.Printf("Error copying output from command to PTM: %v", err)
						}
					}
				}()
			}
		}

		if len(cmd.Env) > 0 {
			// clear the environment, otherwise $PATH is not correctly set.
			cmd.Env = nil
		} else {
			// If no tty has been set, lets fall back to os.Stdin
			cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		}

		if err := cmd.Start(); err != nil {
			log.Printf("Error starting %v: %v", cmd, err)
			continue
		}

		// Close the PTM and PTS after the command exits
		go func() {
			cmd.Wait()
			m.Close()
			s.Close()
		}()

		for {
			var s unix.WaitStatus
			var r unix.Rusage
			if p, err := unix.Wait4(-1, &s, 0, &r); p == cmd.Process.Pid {
				debug("Shell exited, exit status %d", s.ExitStatus())
				break
			} else if p != -1 {
				debug("Reaped PID %d, exit status %d", p, s.ExitStatus())
			} else {
				debug("Error from Wait4 for orphaned child: %v", err)
				break
			}
		}
		if err := cmd.Process.Release(); err != nil {
			log.Printf("Error releasing process %v: %v", cmd, err)
		}
	}
	return cmdCount
}

// This comes from the pty package
func ioctl(f *os.File, cmd, ptr uintptr) error {
	return ioctlInner(f.Fd(), cmd, ptr) // Fall back to blocking io.
}

func ioctlInner(fd, cmd, ptr uintptr) error {
	_, _, e := syscall.Syscall(syscall.SYS_IOCTL, fd, cmd, ptr)
	if e != 0 {
		return e
	}
	return nil
}

func ptsname(f *os.File) (string, error) {
	var n uint32
	err := ioctl(f, syscall.TIOCGPTN, uintptr(unsafe.Pointer(&n))) //nolint:gosec // Expected unsafe pointer for Syscall call.
	if err != nil {
		return "", err
	}
	return "/dev/pts/" + strconv.Itoa(int(n)), nil
}
