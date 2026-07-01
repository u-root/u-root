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
	"strings"

	"github.com/u-root/u-root/pkg/pty"
	"github.com/u-root/u-root/pkg/termios"
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
		if mtty && len(ttyNames) > 0 {
			// Save tty names for later use in RunCommands
			// We don't open them here because we need to set them to raw mode
			// and multiplex I/O through a PTY
			for i, name := range ttyNames {
				c.Env = append(c.Env, fmt.Sprintf("tty%d=/dev/%s", i, name))
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

// FIX ME: make it not linux-specific
// RunCommands runs commands in sequence.
//
// RunCommands returns how many commands existed and were attempted to run.
//
// commands must refer to absolute paths at the moment.
func RunCommands(debug func(string, ...any), commands ...*exec.Cmd) int {
	var cmdCount int
	for _, cmd := range commands {
		if _, err := os.Stat(cmd.Path); os.IsNotExist(err) {
			debug("%v", err)
			continue
		}

		cmdCount++
		debug("Trying to run %v", cmd)

		// Collect TTY names from environment
		var ttyNames []string
		var cleanEnv []string
		for _, envVar := range cmd.Env {
			if strings.HasPrefix(envVar, "tty") && strings.Contains(envVar, "=") {
				parts := strings.SplitN(envVar, "=", 2)
				if len(parts) == 2 {
					ttyNames = append(ttyNames, parts[1])
				}
			} else {
				cleanEnv = append(cleanEnv, envVar)
			}
		}

		// If we have TTYs from kernel cmdline, use PTY multiplexing like the shell tool
		if len(ttyNames) > 0 {
			debug("Setting up multi-TTY with %v", ttyNames)

			// Open and configure all TTYs
			var ttys []*os.File
			for _, ttyName := range ttyNames {
				debug("Opening TTY %v", ttyName)
				tty, err := os.OpenFile(ttyName, os.O_RDWR, 0)
				if err != nil {
					debug("Error opening TTY %v: %v", ttyName, err)
					continue
				}

				// Set to raw mode - critical for serial console to work properly
				// Without raw mode, serial has line buffering, echo, and flow control
				if err := termios.MakeRawFile(tty); err != nil {
					debug("Error setting TTY %v to raw mode: %v", ttyName, err)
					// Continue anyway - better to have non-raw than nothing
				}

				ttys = append(ttys, tty)
			}

			if len(ttys) == 0 {
				debug("No TTYs could be opened, falling back to default")
				cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
			} else {
				// Create PTY for the command
				ptmx, pts, err := pty.NewPTMS()
				if err != nil {
					debug("Error creating PTY: %v", err)
					return 0
				}

				// Connect command to PTS
				cmd.Stdin, cmd.Stdout, cmd.Stderr = pts, pts, pts

				// Create channel for clean shutdown
				done := make(chan struct{})

				// Create buffered channels for input from each TTY to prevent input loss
				// due to goroutine scheduling. Direct io.Copy can lose characters.
				inputChans := make([]chan []byte, len(ttys))
				for i := range inputChans {
					inputChans[i] = make(chan []byte, 1024)
				}

				// Read from each TTY into its buffered channel
				for i, tty := range ttys {
					t := tty // capture for goroutine
					ch := inputChans[i]
					go func() {
						buf := make([]byte, 1024)
						for {
							select {
							case <-done:
								close(ch)
								return
							default:
							}
							n, err := t.Read(buf)
							if err != nil {
								if err != io.EOF {
									select {
									case <-done:
										// Shutting down, ignore error
									default:
										debug("TTY read error: %v", err)
									}
								}
								close(ch)
								return
							}
							if n > 0 {
								data := make([]byte, n)
								copy(data, buf[:n])
								select {
								case ch <- data:
								case <-done:
									close(ch)
									return
								}
							}
						}
					}()
				}

				// Multiplex input from all TTY channels to PTM using blocking select
				go func() {
					for {
						// Build select cases dynamically
						var allClosed = true
						for _, ch := range inputChans {
							if ch != nil {
								allClosed = false
								break
							}
						}
						if allClosed {
							return
						}

						// Try to read from any available channel
						for i, ch := range inputChans {
							if ch == nil {
								continue
							}
							select {
							case data, ok := <-ch:
								if !ok {
									inputChans[i] = nil
									continue
								}
								ptmx.Write(data)
							default:
								// Non-blocking per channel, but we'll loop
							}
						}
					}
				}()

				// Multiplex output: PTM → all TTYs (fan-out)
				// Use io.Copy with TeeReader for efficiency like the shell tool
				if len(ttys) == 2 {
					// Optimize for the common case of 2 TTYs
					go io.Copy(ttys[0], io.TeeReader(ptmx, ttys[1]))
				} else {
					// General case: fan out to all TTYs
					go func() {
						buf := make([]byte, 1024)
						for {
							select {
							case <-done:
								return
							default:
							}
							n, err := ptmx.Read(buf)
							if err != nil {
								if err != io.EOF {
									select {
									case <-done:
										// Shutting down, ignore error
									default:
										debug("PTY read error: %v", err)
									}
								}
								return
							}
							if n > 0 {
								// Fan out to all TTYs
								for _, tty := range ttys {
									tty.Write(buf[:n])
								}
							}
						}
					}()
				}

				// Clear tty environment variables and restore clean environment
				cmd.Env = cleanEnv

				if err := cmd.Start(); err != nil {
					log.Printf("Error starting %v: %v", cmd, err)
					// Clean up before continuing
					close(done)
					for _, tty := range ttys {
						tty.Close()
					}
					ptmx.Close()
					pts.Close()
					continue
				}

				// Wait for command and reap orphans
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

				// Clean up after command exits
				close(done)
				for _, tty := range ttys {
					tty.Close()
				}
				ptmx.Close()
				pts.Close()

				if err := cmd.Process.Release(); err != nil {
					log.Printf("Error releasing process %v: %v", cmd, err)
				}

				// We handled this command completely, continue to next
				continue
			}
		} else {
			// No multi-TTY, use default I/O
			cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		}

		// Clear tty environment variables and restore clean environment
		cmd.Env = cleanEnv

		if err := cmd.Start(); err != nil {
			log.Printf("Error starting %v: %v", cmd, err)
			continue
		}

		// Wait for command and reap orphans
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
