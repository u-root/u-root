// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package expect

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	stdlog "log"

	"google.golang.org/grpc/codes"

	log "github.com/golang/glog"
	term "github.com/google/goterm/term"

	"golang.org/x/crypto/ssh"
)

var (
	// eReg identifies the expect commands in the shell files.
	eReg = regexp.MustCompile(`^# e: (.*)$`)
	// sReg identifies the send commands in the shell files.
	sReg = regexp.MustCompile(`^# s: (.*)$`)
	// sReg identifies the sw/case commands in the shell files.
	cReg = regexp.MustCompile(`^# c: '(.*)' (.*)$`)
)

const expTestData = "./testdata/"

var privateKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpgIBAAKCAQEAnUKOpnWUyxCGy9lyTLK3Crd5BXGG718wIgHzpCvKbopVbYu4
xa9Fk1cjPHp2F/pPwXmIKJuTVpkm/VXcFfABH6cMeszaXqVBhqm6AA0E7y0K0oYZ
GUMMm3sBLPV3ydUHECI2NnEXOCLGysKM6Ht2ZuGxCKdXpquRE1HdLUUIJep31gSO
J4dyQRk12VYHrpTjz1Tzv9prf76vJqYmr6+axeH7I3/9KGnPe1vD0z8NhOwqQONz
DMSjpbSYFiVyRbDVgSi9xq+BeFrASZuwkoHut3tzmIvQLcGIR+LoOeN9mDpXbseQ
y84PXkduF8udrBCIBETekmK7kqi6hwLRCn9/twIDAQABAoIBAQCTeEOvQ5oJpvDR
HpN56ymNCiqZ+TERLhFEAtKIRGxrppufw6O89bToC5HGeAxgReIey6nscqADWFFg
xfBCPjO/i/Y+/fVVReEht+3teEgFRhbc/tVwhBjBgOLEV1hC09rwvTRbb0fX43zJ
zRE4Pfb1WXWbaNngOQkttdoURqTyb+n8zgwx0AUsueSrYxTk1UTF+Jet7g23jRjd
YCCx9qhHez5yif+1LZGIqJD0OKGHr9q+bbOZpy5dqjamuf9ulBvnZkKzcuHf0m/W
Vhf9YI8kOQhPXfztTnZrN5Jg64gGuvJ0sEZucp5hOR4hYkLagOCaUuNIeHG7SsYU
hChWCDphAoGBAM7nr4t714etJnuRE39FG+rylV5K3T2osr2Hwp7wSZfXHZUcnk2N
KSDA9tzeFYX7QxUlc7qNwsLC2WkK97x1WbrNdO8Zn5lmBHfyIS7jxEGQY9htWrgr
sjAaUg/JfHMu/lxNAigXAaGU2VzsTySgB3eWbfaAd/sImaxnVPBajOARAoGBAMKT
NykNYl3zg1pIXGQlu3a0pTv2gRcBnnW7bUAM6b6tDdQZ+5cbu43MfAwhsGsMZ/HL
gKQRJI952olrPa4dEiirxUKqfVVPLnDcUSu6uhvJFpzN2YVdMyPNWc3V9lfIztbu
UvCvupnmeViG9qRoJgbMLIBLBN1oS5MKOL55oqtHAoGBAMX21XZe4qRVHlniQEZo
aELPIe1bMf3Z2FMRfzw1aiSW1R4jiK9o3a4SEuDWuL893jxwXh9jnbJdXklsDgbK
PTVHeZd/672I58Of7vH/SXr13SJp1wAaBt6RgGzMen92uja0E9kp0gy475RCIaNI
XnykeMf+uU1+OBLFt3ZVHS8RAoGBAJr/BpvPK6LHzsTmi6LDY/gFovKHRQH8qiwC
595z6ueXl0J0iDQxRVCJqe9IDu7XbR3yDEGl3kfku69oHDRMuCBp5LNceIaykr4Y
4xhAoOxtXXP/jt1sBsboWDddz+TR8+LG6o8MjUr3i4Z3zJXe2RvlHTX9jJyK7ljt
dZJV9r0VAoGBAI+yBh2oa5D66XztQ2pfKm/B03RYARR0iKvho6Ass1nM1YSvhflt
CkGNYNP5Mwr/VpfTppl0JTl9++gtoAmDgVm19JUYiABhYNmbTq62STJb+LkEL+xK
Jf5UkUJOE8Rf+A1vmI1igjVffSIRLTJC6zOX0JCZMIFKZhyTZsPOuFcm
-----END RSA PRIVATE KEY-----`)

// SSHServer represents one local SSH server.
type SSHServer struct {
	batch []Batcher
	cfg   *ssh.ServerConfig
	l     net.Listener
	sync.Mutex
	opts []Option
	term *term.Termios
}

// New returns a new SSHServer configured for username/password.
func NewSSHServer(user, pass string, b []Batcher, t *term.Termios, opts ...Option) *SSHServer {
	srv := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			if c.User() == user && string(password) == pass {
				return nil, nil
			}
			return nil, errors.New("password rejected for:" + c.User())
		},
	}
	k, _ := ssh.ParsePrivateKey(privateKey)
	srv.AddHostKey(k)
	if t == nil {
		t = &term.Termios{
			Wz: term.Winsize{
				WsCol: 132,
				WsRow: 43,
			},
		}
	}
	return &SSHServer{cfg: srv, batch: b, term: t, opts: opts}
}

// Batcher replaces the batcher.
func (s *SSHServer) Batcher(b []Batcher) {
	s.Lock()
	s.batch = b
	s.Unlock()
}

// Termios replaces the termios.
func (s *SSHServer) Termios(t *term.Termios) {
	if t == nil {
		t = &term.Termios{
			Wz: term.Winsize{
				WsCol: 132,
				WsRow: 43,
			},
		}
	}
	s.term = t
}

// Serve spins up the SSH server and returns the port used.
func (s *SSHServer) Serve() (uint16, error) {
	l, err := net.Listen("tcp", "")
	if err != nil {
		return 0, err
	}
	s.l = l
	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				log.Errorf("Accept failed: %v", err)
				return
			}
			go s.runBatch(c)
		}
	}()
	_, port, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		return 0, err
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		return 0, err
	}
	return uint16(p), nil
}

// Close closes the SSH server write pipe.
func (s *SSHServer) Close() error {
	return s.l.Close()
}

const testTimeout = 20 * time.Second

// RFC 4254 Section 6.2.
type ptyRequestMsg struct {
	Term     string
	Columns  uint32
	Rows     uint32
	Width    uint32
	Height   uint32
	Modelist string
}

func (s *SSHServer) runBatch(conn net.Conn) {
	defer conn.Close()
	_, chs, rq, err := ssh.NewServerConn(conn, s.cfg)
	if err != nil {
		log.Errorf("ssh.NewServerConn failed: %v", err)
		return
	}

	go ssh.DiscardRequests(rq)

	for ch := range chs {
		switch ch.ChannelType() {
		case "session":
			sch, in, err := ch.Accept()
			if err != nil {
				log.Errorf("ch.Accept failed: %v", err)
				return
			}
			for sess := range in {
				switch sess.Type {
				case "dummy":
					if err := sess.Reply(true, nil); err != nil {
						log.Errorf("sess.Reply(%t,nil) failed: %v", true, err)
					}
				case "pty-req":
					ptyReq := ptyRequestMsg{}
					if err := ssh.Unmarshal(sess.Payload, &ptyReq); err != nil {
						if err := sess.Reply(false, nil); err != nil {
							log.Errorf("sess.Reply(%t,nil) failed: %v", false, err)
						}
						log.Errorf("ssh.Unmarshal of PTYRequest failed: %v", err)
						continue
					}
					if ptyReq.Columns != uint32(s.term.Wz.WsCol) || ptyReq.Rows != uint32(s.term.Wz.WsRow) {
						log.Errorf("PTY cols/rows: %d/%d want: %d/%d", ptyReq.Columns, ptyReq.Rows, s.term.Wz.WsCol, s.term.Wz.WsRow)
						if err := sess.Reply(false, nil); err != nil {
							log.Errorf("sess.Reply(%t,nil) failed: %v", false, err)
						}
						continue
					}

					if err := sess.Reply(true, nil); err != nil {
						log.Errorf("sess.Reply(%t,nil) failed: %v", true, err)
					}
				case "shell":
					log.Infof("Shell request coming in")
					resCh := make(chan error)
					defer close(resCh)
					if err := sess.Reply(true, nil); err != nil {
						log.Errorf("sess.Reply(%t,nil) failed: %v", true, err)
					}

					rIn, wIn := io.Pipe()
					rOut, wOut := io.Pipe()
					go io.Copy(sch, rIn)
					go io.Copy(wOut, sch)

					go func() {

						exp, _, err := SpawnGeneric(&GenOptions{
							In:  wIn,
							Out: rOut,
							Wait: func() error {
								return <-resCh
							},
							Close: func() error { return wIn.Close() },
							Check: func() bool { return true },
						}, testTimeout*2, s.opts...)
						s.Lock()
						out, err := exp.ExpectBatch(s.batch, testTimeout*2)
						if err != nil {
							log.Errorf("exp.ExpectBatch(%v) failed: %v, res: %v", s.batch, err, out)
						}
						s.Unlock()
					}()
				default:
					sess.Reply(false, []byte(fmt.Sprint("session type not supported")))
				}
			}
		default:
			ch.Reject(ssh.UnknownChannelType, "channel type not supported")
		}
	}
}

// Tc is an example of implementing custom tag functions.
type Tc uint32

func NewTc() (t Tc) {
	return
}

// Count works like ContinueLog with a counter.
func (t Tc) Count(msg string, s *Status) func() (Tag, *Status) {
	return func() (Tag, *Status) {
		t++
		log.Infof("%d: %s", t, msg)
		return ContinueTag, s
	}
}

// NextLog adds loggin and counting to the Next tag.
func (t Tc) NextLog(msg string) func() (Tag, *Status) {
	return func() (Tag, *Status) {
		t++
		log.Infof("Next %d: %s", t, msg)
		return NextTag, NewStatus(codes.Unimplemented, "Should not matter")
	}
}

func TestBatcher(t *testing.T) {
	tests := []struct {
		name     string
		clt, srv []Batcher
		fail     bool
	}{
		{
			name: "Config mode",
			clt: []Batcher{
				&BExp{`router1>`},
				&BSnd{"conf t\n"},
				&BExp{`\(configure\) router1>`},
			},
			srv: []Batcher{
				&BSnd{`router1> `},
				&BExp{"conf t\n"},
				&BSnd{`(configure) router1> `},
			}}, {
			name: "Login caser",
			clt: []Batcher{
				&BCas{[]Caser{
					&Case{R: regexp.MustCompile(`Login: `), S: "TestUser\n", T: LogContinue("at login prompt", NewStatus(codes.PermissionDenied, "wrong username")), Rt: 1},
					&Case{R: regexp.MustCompile(`Password: `), S: "TestPass\n", T: LogContinue("at password prompt", NewStatus(codes.PermissionDenied, "wrong pass")), Rt: 1},
					&Case{R: regexp.MustCompile(`Permission denied`), T: Fail(NewStatus(codes.PermissionDenied, "login failed"))},
					&Case{R: regexp.MustCompile(`router 1>`), T: OK()},
				}},
			},
			srv: []Batcher{
				&BSnd{"Login: "},
				&BCas{[]Caser{
					&Case{R: regexp.MustCompile("TestUser\n"), S: `Password: `, T: Continue(NewStatus(codes.PermissionDenied, "permission denied")), Rt: 3},
					&Case{R: regexp.MustCompile("TestPass\n"), S: `router 1> `, T: OK()},
				},
				},
			}}, {
			name: "100 Hello World",
			clt: []Batcher{
				&BSnd{`Hello `},
				&BCas{[]Caser{
					&Case{R: regexp.MustCompile("Done"), T: OK()},
					&Case{R: regexp.MustCompile("World"), S: "Hello ", T: NewTc().Count("Hello", NewStatus(codes.OutOfRange, "too many worlds")), Rt: 100},
				}},
			},
			srv: []Batcher{
				&BCas{[]Caser{
					&Case{R: regexp.MustCompile("Hello"), S: "World\n", T: NewTc().Count("World", NewStatus(codes.OK, "too many hellos")), Rt: 99},
				}},
				&BSnd{"Done"},
			},
		}, {
			name: "100 Hello World using Next tag",
			clt: []Batcher{
				&BSnd{`Hello `},
				&BCas{[]Caser{
					&Case{R: regexp.MustCompile("Done"), T: OK()},
					&Case{R: regexp.MustCompile("World"), S: "Hello ", T: NewTc().Count("Hello", NewStatus(codes.OutOfRange, "too many worlds")), Rt: 100},
				}},
			},
			srv: []Batcher{
				&BCas{[]Caser{
					&Case{R: regexp.MustCompile("Hello"), S: "World\n", T: NewTc().NextLog("World"), Rt: 99},
					&Case{R: regexp.MustCompile("Hello"), S: "Done\n", T: LogContinue("Done sent", NewStatus(codes.OK, "Done sent"))},
				}},
			},
		},
	}

	srv := NewSSHServer("test", "test", nil, nil, CheckDuration(40*time.Millisecond))
	port, err := srv.Serve()
	if err != nil {
		t.Fatalf("srv.Serve failed: %v", err)
	}
	defer srv.Close()
	for _, tst := range tests {
		srv.Batcher(tst.srv)
		clt, err := ssh.Dial("tcp", net.JoinHostPort("localhost", strconv.Itoa(int(port))),
			&ssh.ClientConfig{
				User:            "test",
				Auth:            []ssh.AuthMethod{ssh.Password("test")},
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			})
		if err != nil {
			t.Errorf("%s: ssh.Dial failed: %v", tst.name, err)
			continue
		}
		e, _, err := SpawnSSH(clt, testTimeout*2)
		if err != nil {
			t.Errorf("%s: SpawnSSH failed: %v", tst.name, err)
		}
		res, err := e.ExpectBatch(tst.clt, testTimeout*2)
		if got, want := err != nil, tst.fail; got != want {
			t.Errorf("%s: e.ExpectBatch(%v,_) = %v want: %v, res: %q", tst.name, tst.clt, err, want, res)
		}
		if err := clt.Close(); err != nil {
			t.Errorf("%s: clt.Close failed: %v", tst.name, err)
		}
	}
}

var (
	cliMap = map[string]string{
		"show system uptime": `Current time:      1998-10-13 19:45:47 UTC
Time Source:       NTP CLOCK
System booted:     1998-10-12 20:51:41 UTC (22:54:06 ago)
Protocols started: 1998-10-13 19:33:45 UTC (00:12:02 ago)
Last configured:   1998-10-13 19:33:45 UTC (00:12:02 ago) by abc
12:45PM  up 22:54, 2 users, load averages: 0.07, 0.02, 0.01

testuser@testrouter#`,
		"show version": `Cisco IOS Software, 3600 Software (C3660-I-M), Version 12.3(4)T

TAC Support: http://www.cisco.com/tac
Copyright (c) 1986-2003 by Cisco Systems, Inc.
Compiled Thu 18-Sep-03 15:37 by ccai

ROM: System Bootstrap, Version 12.0(6r)T, RELEASE SOFTWARE (fc1)
ROM:

C3660-1 uptime is 1 week, 3 days, 6 hours, 41 minutes
System returned to ROM by power-on
System image file is "slot0:tftpboot/c3660-i-mz.123-4.T"

Cisco 3660 (R527x) processor (revision 1.0) with 57344K/8192K bytes of memory.
Processor board ID JAB055180FF
R527x CPU at 225Mhz, Implementation 40, Rev 10.0, 2048KB L2 Cache

3660 Chassis type: ENTERPRISE
2 FastEthernet interfaces
4 Serial interfaces
DRAM configuration is 64 bits wide with parity disabled.
125K bytes of NVRAM.
16384K bytes of processor board System flash (Read/Write)

Flash card inserted. Reading filesystem...done.
20480K bytes of processor board PCMCIA Slot0 flash (Read/Write)

Configuration register is 0x2102

testrouter#`,
		"show system users": `7:30PM  up 4 days,  2:26, 2 users, load averages: 0.07, 0.02, 0.01
USER     TTY FROM              LOGIN@  IDLE WHAT
root     d0  -                Fri05PM 4days -csh (csh)
blue   p0 level5.company.net 7:30PM     - cli

testuser@testrouter#`,
	}
)

func fakeCli(tMap map[string]string, in io.Reader, out io.Writer) {
	scn := bufio.NewScanner(in)
	for scn.Scan() {
		tst, ok := tMap[scn.Text()]
		if !ok {
			out.Write([]byte(fmt.Sprintf("command: %q not found", scn.Text())))
			continue
		}
		_, err := out.Write([]byte(tst))
		if err != nil {
			log.Warningf("Write of %q failed: %v", tst, err)
			return
		}
	}
}

// ExampleDebugCheck toggles the DebugCheck option.
func ExampleDebugCheck() {
	rIn, wIn := io.Pipe()
	rOut, wOut := io.Pipe()
	rLog, wLog := io.Pipe()
	waitCh := make(chan error)
	defer rIn.Close()
	defer wOut.Close()
	defer wLog.Close()

	go fakeCli(cliMap, rIn, wOut)

	exp, r, err := SpawnGeneric(&GenOptions{
		In:    wIn,
		Out:   rOut,
		Wait:  func() error { return <-waitCh },
		Close: func() error { return wIn.Close() },
		Check: func() bool {
			return true
		}}, -1)
	if err != nil {
		log.Errorf("SpawnGeneric failed: %v", err)
		return
	}
	re := regexp.MustCompile("testrouter#")
	interact := func() {
		for cmd := range cliMap {
			if err := exp.Send(cmd + "\n"); err != nil {
				log.Errorf("exp.Send(%q) failed: %v\n", cmd+"\n", err)
				return
			}
			out, _, err := exp.Expect(re, -1)
			if err != nil {
				log.Errorf("exp.Expect(%v) failed: %v out: %v", re, err, out)
				return
			}
		}
	}

	go func() {
		var last string
		scn := bufio.NewScanner(rLog)
		for scn.Scan() {
			ws := strings.Split(scn.Text(), " ")
			if ws[0] == last {
				continue
			}
			last = ws[0]
			fmt.Println(ws[0])
		}
	}()

	fmt.Println("First round")
	interact()
	fmt.Println("Second round - Debugging enabled")
	prev := exp.Options(DebugCheck(stdlog.New(wLog, "DebugExample ", 0)))
	interact()
	exp.Options(prev)
	fmt.Println("Last round - Previous Check put back")
	interact()

	waitCh <- nil
	exp.Close()
	wOut.Close()

	<-r

	// Output:
	// First round
	// Second round - Debugging enabled
	// DebugExample
	// Last round - Previous Check put back
}

// ExampleChangeCheck changes the check function runtime for an Expect session.
func ExampleChangeCheck() {
	rIn, wIn := io.Pipe()
	rOut, wOut := io.Pipe()
	waitCh := make(chan error)
	outCh := make(chan string)
	defer close(outCh)

	go fakeCli(cliMap, rIn, wOut)
	go func() {
		var last string
		for s := range outCh {
			if s == last {
				continue
			}
			fmt.Println(s)
			last = s
		}
	}()

	exp, r, err := SpawnGeneric(&GenOptions{
		In:    wIn,
		Out:   rOut,
		Wait:  func() error { return <-waitCh },
		Close: func() error { return wIn.Close() },
		Check: func() bool {
			outCh <- "Original check"
			return true
		}}, -1)
	if err != nil {
		fmt.Printf("SpawnGeneric failed: %v\n", err)
		return
	}
	re := regexp.MustCompile("testrouter#")
	interact := func() {
		for cmd := range cliMap {
			if err := exp.Send(cmd + "\n"); err != nil {
				fmt.Printf("exp.Send(%q) failed: %v\n", cmd+"\n", err)
				return
			}
			out, _, err := exp.Expect(re, -1)
			if err != nil {
				fmt.Printf("exp.Expect(%v) failed: %v out: %v", re, err, out)
				return
			}
		}
	}
	interact()
	prev := exp.Options(ChangeCheck(func() bool {
		outCh <- "Replaced check"
		return true
	}))
	interact()
	exp.Options(prev)
	interact()

	waitCh <- nil
	exp.Close()
	wOut.Close()

	<-r
	// Output:
	// Original check
	// Replaced check
	// Original check
}

// ExampleVerbose changes the Verbose and VerboseWriter options.
func ExampleVerbose() {
	rIn, wIn := io.Pipe()
	rOut, wOut := io.Pipe()
	waitCh := make(chan error)
	outCh := make(chan string)
	defer close(outCh)

	go fakeCli(cliMap, rIn, wOut)
	go func() {
		var last string
		for s := range outCh {
			if s == last {
				continue
			}
			fmt.Println(s)
			last = s
		}
	}()

	exp, r, err := SpawnGeneric(&GenOptions{
		In:    wIn,
		Out:   rOut,
		Wait:  func() error { return <-waitCh },
		Close: func() error { return wIn.Close() },
		Check: func() bool {
			return true
		}}, -1, Verbose(true), VerboseWriter(os.Stdout))
	if err != nil {
		fmt.Printf("SpawnGeneric failed: %v\n", err)
		return
	}
	re := regexp.MustCompile("testrouter#")
	var interactCmdSorted []string
	for k := range cliMap {
		interactCmdSorted = append(interactCmdSorted, k)
	}
	sort.Strings(interactCmdSorted)
	interact := func() {
		for _, cmd := range interactCmdSorted {
			if err := exp.Send(cmd + "\n"); err != nil {
				fmt.Printf("exp.Send(%q) failed: %v\n", cmd+"\n", err)
				return
			}
			out, _, err := exp.Expect(re, -1)
			if err != nil {
				fmt.Printf("exp.Expect(%v) failed: %v out: %v", re, err, out)
				return
			}
		}
	}
	interact()

	waitCh <- nil
	exp.Close()
	wOut.Close()

	<-r
	// Output:
	// [34mSent:[39m "show system uptime\n"
	// [32mMatch for RE:[39m "testrouter#" found: ["testrouter#"] Buffer: Current time:      1998-10-13 19:45:47 UTC
	// Time Source:       NTP CLOCK
	// System booted:     1998-10-12 20:51:41 UTC (22:54:06 ago)
	// Protocols started: 1998-10-13 19:33:45 UTC (00:12:02 ago)
	// Last configured:   1998-10-13 19:33:45 UTC (00:12:02 ago) by abc
	// 12:45PM  up 22:54, 2 users, load averages: 0.07, 0.02, 0.01
	//
	// testuser@testrouter#
	// [34mSent:[39m "show system users\n"
	// [32mMatch for RE:[39m "testrouter#" found: ["testrouter#"] Buffer: 7:30PM  up 4 days,  2:26, 2 users, load averages: 0.07, 0.02, 0.01
	// USER     TTY FROM              LOGIN@  IDLE WHAT
	// root     d0  -                Fri05PM 4days -csh (csh)
	// blue   p0 level5.company.net 7:30PM     - cli
	//
	// testuser@testrouter#
	// [34mSent:[39m "show version\n"
	// [32mMatch for RE:[39m "testrouter#" found: ["testrouter#"] Buffer: Cisco IOS Software, 3600 Software (C3660-I-M), Version 12.3(4)T
	//
	// TAC Support: http://www.cisco.com/tac
	// Copyright (c) 1986-2003 by Cisco Systems, Inc.
	// Compiled Thu 18-Sep-03 15:37 by ccai
	//
	// ROM: System Bootstrap, Version 12.0(6r)T, RELEASE SOFTWARE (fc1)
	// ROM:
	//
	// C3660-1 uptime is 1 week, 3 days, 6 hours, 41 minutes
	// System returned to ROM by power-on
	// System image file is "slot0:tftpboot/c3660-i-mz.123-4.T"
	//
	// Cisco 3660 (R527x) processor (revision 1.0) with 57344K/8192K bytes of memory.
	// Processor board ID JAB055180FF
	// R527x CPU at 225Mhz, Implementation 40, Rev 10.0, 2048KB L2 Cache
	//
	// 3660 Chassis type: ENTERPRISE
	// 2 FastEthernet interfaces
	// 4 Serial interfaces
	// DRAM configuration is 64 bits wide with parity disabled.
	// 125K bytes of NVRAM.
	// 16384K bytes of processor board System flash (Read/Write)
	//
	// Flash card inserted. Reading filesystem...done.
	// 20480K bytes of processor board PCMCIA Slot0 flash (Read/Write)
	//
	// Configuration register is 0x2102
	//
	// testrouter#

}

// TestTee tests the Tee option can write to a file.
func TestTee(t *testing.T) {
	// Create a temporary file to tee output.
	f, err := ioutil.TempFile("", "goexpect-tee")
	if err != nil {
		t.Fatalf("Could not create temporary file: %v", err)
	}
	fileName := f.Name()
	defer os.Remove(fileName)

	// Send abcdef to cat 4096 times.
	input := "abcdef\n"
	e, _, err := Spawn("cat", 400*time.Millisecond, Tee(f), CheckDuration(1*time.Millisecond))
	for i := 0; i < 4096; i++ {
		e.Send(input)
		re := regexp.MustCompile(input)
		if _, _, err = e.Expect(re, 400*time.Millisecond); err != nil {
			t.Errorf("Expect(%q) failed: %v", input, err)
		}
	}
	e.Close()

	// Check the tee'd output.
	got, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatalf("Could not read temporary file: %v", err)
	}
	want := strings.Repeat(input, 4096)
	if string(got) != want {
		t.Errorf("tee output mismatch, got: %q want: %q", got, want)
	}
}

// TestTee_SpawnFake tests the Tee option can operate on SpawnFake.
func TestTee_SpawnFake(t *testing.T) {
	// Create a temporary file to tee output.
	f, err := ioutil.TempFile("", "goexpect-tee")
	if err != nil {
		t.Fatalf("Could not create temporary file: %v", err)
	}
	fileName := f.Name()
	defer os.Remove(fileName)

	msg := `
Pretty please don't hack my chassis

router1> `
	srv := []Batcher{
		&BSnd{msg},
	}
	re := regexp.MustCompile("router1>")
	timeout := 2 * time.Second
	exp, endch, err := SpawnFake(srv, timeout, Tee(f))
	if err != nil {
		t.Fatalf("SpawnFake failed: %v", err)
	}
	out, _, err := exp.Expect(re, timeout)
	if err != nil {
		t.Fatalf("Expect(%q,%v), err: %v, out: %q", re.String(), timeout, err, out)
	}
	exp.Close()
	// wait for end
	<-endch

	// Check the tee'd output.
	got, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatalf("Could not read temporary file: %v", err)
	}
	if string(got) != msg {
		t.Errorf("tee output mismatch, got: %q want: %q", got, msg)
	}
}

// TestSpawnGeneric tests out the generic spawn function.
func TestSpawnGeneric(t *testing.T) {
	fr, fw := io.Pipe()
	tests := []struct {
		name  string
		opt   *GenOptions
		check func() bool
		cli   map[string]string
		re    *regexp.Regexp
		fail  bool
	}{{
		name:  "Clean test",
		check: func() bool { return true },
		cli:   cliMap,
		re:    regexp.MustCompile("testrouter#"),
		fail:  false,
	}, {
		name: "Fail check",
		check: func() bool {
			return false
		},
		cli:  cliMap,
		re:   regexp.MustCompile("testrouter#"),
		fail: true,
	}, {
		name: "In nil",
		opt:  &GenOptions{},
		fail: true,
	}, {
		name: "Out nil",
		opt: &GenOptions{
			In: fw,
		},
		fail: true,
	}, {
		name: "Wait nil",
		opt: &GenOptions{
			In:  fw,
			Out: fr,
		},
		fail: true,
	}}

	for _, tst := range tests {
		t.Logf("Running test: %v", tst.name)
		waitCh := make(chan error)
		rIn, wIn := io.Pipe()
		rOut, wOut := io.Pipe()
		if tst.opt == nil {
			tst.opt = &GenOptions{
				In:    wIn,
				Out:   rOut,
				Wait:  func() error { return <-waitCh },
				Close: func() error { return wIn.Close() },
				Check: tst.check}
		}
		go fakeCli(tst.cli, rIn, wOut)
		exp, r, err := SpawnGeneric(tst.opt, -1)
		if err != nil {
			if !tst.fail {
				t.Errorf("test: %v , SpawnGeneric failed: %v", tst.name, err)
			}
			continue
		}
		gotFail := false
		for cmd := range tst.cli {
			err := exp.Send(cmd + "\n")
			if err != nil {
				if tst.fail {
					gotFail = true
					break
				}
				t.Errorf("Send(%q) failed: %v", cmd, err)
				break
			}
			out, _, err := exp.Expect(tst.re, -1)
			if err != nil {
				if tst.fail {
					gotFail = true
					break
				}
				t.Errorf("Expect(%q) failed: %v, out: %q", tst.re, err, out)
				break
			}
		}
		if gotFail != tst.fail {
			t.Errorf("test: %v , failed status mismatch, got: %v want: %v", tst.name, gotFail, tst.fail)
		}
		waitCh <- nil
		exp.Close()
		wOut.Close()

		<-r
	}
}

// TestSendTimeout tests that Send command can fail on timeout.
func TestSendTimeout(t *testing.T) {
	t.Log("Running test: TestSendTimeout")
	rIn, wIn := io.Pipe()
	rOut, wOut := io.Pipe()
	waitCh := make(chan error)
	outCh := make(chan string)
	defer close(outCh)

	go fakeCli(cliMap, rIn, wOut)
	exp, r, err := SpawnGeneric(
		&GenOptions{
			In:    wIn,
			Out:   rOut,
			Wait:  func() error { return <-waitCh },
			Close: func() error { return nil },
			Check: func() bool { return true },
		}, -1, SendTimeout(time.Second))
	if err != nil {
		t.Fatalf("SpawnGeneric(_, %d , SendTimeout(%v)) failed: %v", -1, time.Second, err)
	}

	if err := wIn.Close(); err != nil {
		t.Fatalf("wIn.Close() failed: %v", err)
	}

	
	if err := exp.Send("test" + "\n"); err != nil {
		t.Fatalf("Send(%q) command failed: %v", "test" + "\n", err)
	}

	if err := exp.Send("test" + "\n"); err == nil {
		t.Errorf("Send(%q) = %t want: %t, err: %v", "test" + "\n", (err != nil), true, err)
	}
	waitCh <- nil
	exp.Close()
	wOut.Close()

	<-r
}


// TestSpawnSSHPTY tests the SSHPTY spawner.
func TestSpawnSSHPTY(t *testing.T) {
	tests := []struct {
		name    string
		fail    bool
		srv     []Batcher
		clt     []Batcher
		sshNil  bool
		srvTerm *term.Termios
		cltTerm term.Termios
	}{{
		name:   "sshClient broken",
		fail:   true,
		sshNil: true,
	}, {
		name: "Empty Termios",
		clt: []Batcher{
			&BSnd{"Hello"},
			&BExp{"World"},
		},
		srv: []Batcher{
			&BExp{"Hello"},
			&BSnd{"World"},
		},
	}, {
		name: "Termios mismatch",
		fail: true,
		srvTerm: &term.Termios{
			Wz: term.Winsize{
				WsCol: 120,
				WsRow: 40,
			}},
		cltTerm: term.Termios{
			Wz: term.Winsize{
				WsCol: 240,
				WsRow: 22,
			},
		},
	}}

	srv := NewSSHServer("test", "test", nil, nil)
	port, err := srv.Serve()
	if err != nil {
		t.Fatalf("srv.Serve failed: %v", err)
	}
	defer func() {
		if err := srv.Close(); err != nil {
			t.Errorf("srv.Close failed: %v", err)
		}
	}()

	for _, tst := range tests {
		srv.Batcher(tst.srv)
		srv.Termios(tst.srvTerm)
		var sshClt *ssh.Client
		if !tst.sshNil {
			clt, err := ssh.Dial("tcp", net.JoinHostPort("localhost", strconv.Itoa(int(port))),
				&ssh.ClientConfig{
					User:            "test",
					Auth:            []ssh.AuthMethod{ssh.Password("test")},
					HostKeyCallback: ssh.InsecureIgnoreHostKey(),
				})
			if err != nil {
				t.Errorf("%s: net.Dial(%q) failed: %v", tst.name, net.JoinHostPort("localhost", strconv.Itoa(int(port))), err)
				continue
			}
			sshClt = clt
		}
		e, _, err := SpawnSSHPTY(sshClt, testTimeout*2, tst.cltTerm)
		if got, want := err != nil, tst.fail; got != want {
			t.Errorf("%s: SpawnSSH = %t want: %t, err: %v", tst.name, got, want, err)
			continue
		}
		if err != nil {
			continue
		}
		res, err := e.ExpectBatch(tst.clt, testTimeout*2)
		if err != nil {
			t.Errorf("%s: e.ExpectBatch failed: %v, out: %v", tst.name, err, res)
			continue
		}
		if err := sshClt.Close(); err != nil {
			t.Errorf("%s: clt.Close failed: %v", tst.name, err)
		}
	}
}

// TestOptions tests manipulating options.
func TestOptions(t *testing.T) {
	tests := []struct {
		name  string
		check func() bool
		opts  []Option
		re    *regexp.Regexp
		fail  bool
	}{{
		name:  "No options",
		check: func() bool { return true },
	}, {
		name:  "No check option",
		opts:  []Option{NoCheck()},
		check: func() bool { return false },
	},
	}

	for _, tst := range tests {
		rIn, wIn := io.Pipe()
		rOut, wOut := io.Pipe()
		go fakeCli(cliMap, rIn, wOut)
		waitCh := make(chan error)
		exp, r, err := SpawnGeneric(&GenOptions{
			In:    wIn,
			Out:   rOut,
			Wait:  func() error { return <-waitCh },
			Close: func() error { return wIn.Close() },
			Check: tst.check}, -1, tst.opts...)
		if err != nil {
			t.Errorf("%s: SpawnGeneric failed: %v", tst.name, err)
			continue
		}
		if got, want := exp.Send("\n\n") != nil, tst.fail; got != want {
			t.Errorf("%s: exp.Send(\"\\n\\n\") = %t want: %t", tst.name, got, want)
		}
		waitCh <- nil
		exp.Close()
		wOut.Close()

		<-r
	}
}

// TestSpawn tests out the Spawn function
func TestSpawn(t *testing.T) {
	tests := []struct {
		name   string
		fail   bool
		cmd    string
		cmdErr bool
	}{{
		name: "Spawn non executable fail",
		fail: true,
		cmd:  "/etc/hosts",
	}, {
		name: "Nil return code",
		cmd:  "/bin/true",
	}, {
		name:   "Non nil return code",
		cmd:    "/bin/false",
		cmdErr: true,
	}, {
		name:   "Spawn cat",
		cmd:    "/bin/cat",
		cmdErr: true,
	}}

	for _, tst := range tests {
		e, errCh, err := Spawn(tst.cmd, 8*time.Second)
		if got, want := err != nil, tst.fail; got != want {
			t.Errorf("%s: Spawn(%q) = %t want: %t, err: %v", tst.name, tst.cmd, got, want, err)
			continue
		}
		if err != nil {
			continue
		}
		<-time.After(2 * time.Second)
		if err := e.Close(); err != nil {
			t.Logf("e.Close failed: %v", err)
		}
		res := <-errCh
		if got, want := res != nil, tst.cmdErr; got != want {
			t.Errorf("%s: errCh = %t want: %t, err: %v", tst.name, got, want, err)
			continue
		}
	}
}

// TestSpawnWithArgs tests that arguments with embedded spaces works.
func TestSpawnWithArgs(t *testing.T) {
	args := []string{"echo", "a   b"}
	e, _, err := SpawnWithArgs(args, 400*time.Millisecond)
	if err != nil {
		t.Errorf("Spawn(echo 'a   b') failed: %v", err)
	}

	// Expected to match
	_, _, err = e.Expect(regexp.MustCompile("a   b"), 400*time.Millisecond)
	if err != nil {
		t.Errorf("Expect(a   b) failed: %v", err)
	}

	// Expected to not match
	_, _, err = e.Expect(regexp.MustCompile("a b"), 400*time.Millisecond)
	if err == nil {
		t.Error("Expect(a b) to not match")
	}

	e.Close()
}

// TestExpect tests the Expect function.
func TestExpect(t *testing.T) {
	tests := []struct {
		name    string
		fail    bool
		srv     []Batcher
		timeout time.Duration
		re      *regexp.Regexp
		re2     *regexp.Regexp
	}{{
		name: "Match prompt",
		srv: []Batcher{
			&BSnd{`
Pretty please don't hack my chassis

router1> `},
		},
		re:      regexp.MustCompile("hack"),
		re2:     regexp.MustCompile("router1>"),
		timeout: 2 * time.Second,
	}, {
		name: "Match fail",
		fail: true,
		re:   regexp.MustCompile("router1>"),
		srv: []Batcher{
			&BSnd{`
Welcome

Router42>`},
		},
		timeout: 1 * time.Second,
	}}

	for _, tst := range tests {
		exp, _, err := SpawnFake(tst.srv, tst.timeout, PartialMatch(true))
		if err != nil {
			if !tst.fail {
				t.Errorf("%s: SpawnFake failed: %v", tst.name, err)
			}
			continue
		}
		out, _, err := exp.Expect(tst.re, tst.timeout)
		if got, want := err != nil, tst.fail; got != want {
			t.Errorf("%s: Expect(%q,%v) = %t want: %t , err: %v, out: %q", tst.name, tst.re.String(), tst.timeout, got, want, err, out)
			continue
		}
        out, _, err = exp.Expect(tst.re2, tst.timeout)
        if got, want := err != nil, tst.fail; got != want {
            t.Errorf("%s: Expect(%q,%v) = %t want: %t , err: %v, out: %q", tst.name, tst.re.String(), tst.timeout, got, want, err, out)
            continue
        }
	}
}

// TestScenarios reads and executes the expect/*.sh test scenarios.
func TestScenarios(t *testing.T) {
	//path := runfiles.Path(expTestData)
	files, err := filepath.Glob(expTestData + "/*.sh")
	if err != nil || len(files) == 0 {
		t.Fatalf("filepath.Glob(%q) failed: %v, not testfile found", expTestData+"/*.sh", err)
	}
L1:
	for _, f := range files {
		_, file := filepath.Split(f)
		tst, err := buildTest(file)
		if err != nil {
			t.Errorf("%s: buildTest(%q) failed: %v", file, file, err)
			continue
		}
		// Spawn the testfile
		exp, r, err := Spawn(f, 0)
		if err != nil {
			t.Errorf("%s: Spawn(%q,0) failed: %v", file, file, err)
			continue
		}
		t.Log("Testing scenariofile:", file)
		for _, ts := range tst {
			switch ts.Cmd() {
			case BatchExpect:
				re := regexp.MustCompile(ts.Arg())
				to := ts.Timeout()
				if to == 0 {
					to = 30 * time.Second
				}
				o, _, err := exp.Expect(re, to)
				if err != nil {
					t.Errorf("%s: Expect(%q,%v) failed: %v, out: %q", file, ts.Arg(), to, err, o)
					continue L1
				}
				t.Log("Scenario:", file, "expect:", ts.Arg(), " found")
				if !re.MatchString(o) {
					t.Fatalf("%s: Doublecheck failed re: %q output: %q", file, ts.Arg(), o)
					continue L1
				}
			case BatchSend:
				if err := exp.Send(ts.Arg()); err != nil {
					t.Fatalf("%s: Send(%q) failed: %v", file, ts.Arg(), err)
					continue L1
				}
			case BatchSwitchCase:
				to := ts.Timeout()
				if to == 0 {
					to = 30 * time.Second
				}
				o, _, _, err := exp.ExpectSwitchCase(ts.Cases(), to)
				if err != nil {
					if err.Error() == "process not running" {
						t.Logf("%s: exp.ExpectSwitchCase(%v,%v) failed: %v, process returned: %v", file, ts.Cases(), to, err, <-r)
					}
					t.Errorf("%s: ExpectSwitchCase failed: %v case: %v output: %q", file, err, ts.Cases(), o)
					continue L1
				}
			}
		}
		if err := exp.Close(); err != nil {
			t.Logf("exp.Close failed: %v", err)

		}
	}
}

// TestBatchScenarios runs through the scenarios again , this time as Batchjobs.
func TestBatchScenarios(t *testing.T) {
	//path := runfiles.Path(expTestData)
	files, err := filepath.Glob(expTestData + "/*.sh")
	if err != nil || len(files) == 0 {
		t.Fatalf("filepath.Glob(%q) failed: %v, not testfile found", expTestData+"/*.sh", err)
	}
	for _, f := range files {
		_, file := filepath.Split(f)
		tsts, err := buildTest(file)
		if err != nil {
			t.Errorf("%s: buildTest(%q) failed: %v", f, f, err)
			continue
		}
		batch := []Batcher{}
		for _, tst := range tsts {
			switch tst.Cmd() {
			case BatchExpect:
				batch = append(batch, &BExp{tst.Arg()})
			case BatchSend:
				batch = append(batch, &BSnd{tst.Arg()})
			case BatchSwitchCase:
				batch = append(batch, &BCas{tst.Cases()})
			}
		}
		exp, r, err := Spawn(f, 30*time.Second, Verbose(true))
		if err != nil {
			t.Errorf("%s: Spawn(%q) failed: %v", file, file, err)
			continue
		}
		res, err := exp.ExpectBatch(batch, 30*time.Second)
		if err != nil {
			t.Errorf("%s: ExpectBatch failed: %v, res: %v", file, err, res)
			continue
		}
		exp.Close()
		<-r
	}
}

var tMap map[string][]Batcher

// buildTest Reads the sends and expected outputs from the testfiles eg.
func buildTest(fstring string) ([]Batcher, error) {
	if tMap == nil {
		tMap = make(map[string][]Batcher)
	}
	if tst, ok := tMap[fstring]; ok {
		return tst, nil
	}
	//path := runfiles.Path(expTestData)
	f, err := os.Open(expTestData + "/" + fstring)
	if err != nil {
		return []Batcher{}, err
	}
	defer f.Close()
	scn := bufio.NewScanner(f)
	var (
		etst []Batcher
		// tcases temporary []Caser slice
		tcases []Caser
		// incases toggle to tell if we're currently building a []Caser slice for BatchSwitchCase
		incases bool
	)
	for scn.Scan() {
		ln := scn.Text()
		if err := scn.Err(); err != nil {
			return []Batcher{}, err
		}
		if res := cReg.FindStringSubmatch(ln); res != nil {
			incases = true
			res[2] = strings.Replace(res[2], `\n`, "\n", -1)
			tcases = append(tcases, &Case{regexp.MustCompile(res[1]), res[2], nil, 0})
			continue
		}
		if res := eReg.FindStringSubmatch(ln); res != nil {
			if incases {
				incases = false
				etst = append(etst, &BCas{tcases})
				tcases = []Caser{}
			}
			res[1] = strings.Replace(res[1], `\n`, "\n", -1)
			etst = append(etst, &BExp{res[1]})
			continue
		}
		if res := sReg.FindStringSubmatch(ln); res != nil {
			if incases {
				incases = false
				etst = append(etst, &BCas{tcases})
				tcases = []Caser{}
			}
			res[1] = strings.Replace(res[1], `\n`, "\n", -1)
			etst = append(etst, &BSnd{res[1]})
			continue
		}
	}
	// If c: is the last thing we have we need to tie it up
	if incases {
		etst = append(etst, &BCas{tcases})
	}
	tMap[fstring] = etst
	return etst, nil
}
