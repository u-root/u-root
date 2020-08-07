// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package term

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"syscall"
	"testing"
)

var pty *PTY
var tios Termios
var wz *Winsize

// NUMTEST How many combos of window sizes to test
const NUMTEST = 100

// donormfile Creates a standard file.
// Used to test that setting terminal attributes fail on standard files.
func donormfile(input string) (*os.File, error) {
	tdir := os.TempDir()
	pid := os.Getpid()
	nf, err := os.Create(tdir + "/" + "s2test-" + strconv.Itoa(pid))
	if err != nil {
		return nil, err
	}
	os.Remove(tdir + "/" + "s2test-" + strconv.Itoa(pid))
	if _, err = nf.Write([]byte(input)); err != nil {
		return nil, err
	}
	return nf, nil
}

// testraw Checks if we really got all the terminal raw flags set.
func testraw(tr Termios, f string) error {
	if (tr.Iflag & (IGNBRK | BRKINT | PARMRK | ISTRIP | INLCR | IGNCR | ICRNL | IXON)) != 0 {
		return fmt.Errorf("%q Raw failed setting c_iflag , got: %d want: 0", f, tr.Iflag)
	}
	if (tr.Oflag & OPOST) != 0 {
		return fmt.Errorf("%q Raw failed setting Oflag , got: %d want: 0", f, tr.Oflag)
	}
	if (tr.Lflag & (ECHO | ECHONL | ICANON | ISIG | IEXTEN)) != 0 {
		return fmt.Errorf("%q Raw failed setting Lflag, got: %d want: 0", f, tr.Lflag)
	}
	if (tr.Cflag & (PARENB)) != 0 {
		return fmt.Errorf("%q Raw failed setting Cflag , got: %d want: 0", f, tr.Cflag)
	}
	if (tr.Cflag & CSIZE) != CS8 {
		return fmt.Errorf("%q Raw failed setting Cflag CS8, got: %d ", f, tr.Cflag)
	}
	if !(tr.Cc[VMIN] == 1 && tr.Cc[VTIME] == 0) {
		return fmt.Errorf("%q Raw failed setting Cc, got: %d want: 0", f, tr.Cc)
	}
	return nil
}

// testcook confirms that all the flags needed for cooked mode is set.
func testcook(tr Termios, f string) error {
	if tr.Iflag&(BRKINT|IGNPAR|ISTRIP|ICRNL|IXON) != BRKINT+IGNPAR+ISTRIP+ICRNL+IXON {
		return fmt.Errorf("%q Cook failed setting Iflag , got: %d want: %d", f, tr.Iflag, BRKINT+IGNPAR+ISTRIP+ICRNL+IXON)
	}
	if (tr.Oflag & OPOST) != OPOST {
		return fmt.Errorf("%q Cook failed setting Oflag , got: %d want: %d", f, tr.Oflag, OPOST)
	}
	if (tr.Lflag & (ISIG | ICANON)) != ISIG+ICANON {
		return fmt.Errorf("%q Cook failed setting Lflag , got: %d want: %d", f, tr.Lflag, ECHO+ECHONL+ICANON+ISIG+IEXTEN)
	}
	return nil
}

// TestOpenpty Checks if we can open a new PTY.
func TestOpenPTY(t *testing.T) {
	var err error
	if pty, err = OpenPTY(); err != nil {
		t.Fatal("Openpty failed: ", err)
	}
}

// TestClose tests if the PTY Close function handles errors.
func TestClose(t *testing.T) {
	tty, err := OpenPTY()
	if err != nil {
		t.Fatalf("OpenPTY failed: %v", err)
	}
	fd := tty.Slave.Fd()
	// Sneaky syscall Close so we won't just end up with a *File == nil.
	if err = syscall.Close(int(fd)); err != nil {
		t.Fatalf("Closing PTY slave failed: %v", err)
	}
	if err = tty.Close(); err == nil {
		t.Fatal("Close() already closed FD want: fail got: <nil>")
	}
	t.Logf("Closing of half closed PTY resulted in err: %v", err)

	if tty, err = OpenPTY(); err != nil {
		t.Fatalf("OpenPTY failed: %v", err)
	}
	if err = tty.Close(); err != nil {
		t.Fatalf("Closing PTY failed, want: <nil> got: %v", err)
	}
}

// TestIsatty checks Isatty on a standard file and a tty.
func TestIsatty(t *testing.T) {
	f, err := donormfile("TestIsatty")
	if err != nil {
		t.Fatalf("donormfile(\"TestIsatty\") failed: %v", err)
	}
	if Isatty(f) {
		t.Errorf("Isatty for normal file %v got: true want: false", f)
	}
	if !Isatty(pty.Slave) {
		t.Errorf("Isatty for tty file %v got: false want: true", pty.Slave)
	}
}

// TestGetPass checks if reading a password from a TTY works.
func TestGetPass(t *testing.T) {
	f, err := donormfile("TestGetPass")
	if err != nil {
		t.Fatalf("donormfile(\"TestGetPass\") failed: %v", err)
	}
	buf := make([]byte, 512)
	if _, err := GetPass("TestGetPass:", f, buf); err == nil {
		t.Errorf("GetPass(\"TestGetPass\",buf) got: <nil> want: file %s not a tty.", f.Name())
	}
	pty, err := OpenPTY()
	if err != nil {
		t.Fatalf("OpenPTY() failed: %v", err)
	}
	// Readers and writers
	var mu sync.Mutex
	var readbuffer bytes.Buffer
	go func() {
		b := make([]byte, 512)
		for {
			nr, err := pty.Master.Read(b)
			if err != nil {
				break
			}
			mu.Lock()
			readbuffer.Write(b[:nr])
			mu.Unlock()
		}
	}()
	tststring := "SuperSecret\n"
	tstWriter := func(in string) {
		w := []byte(in)
		var err error
		for tot, nr := 0, 0; tot < len(w); tot += nr {
			if nr, err = pty.Master.Write(w[tot:]); err != nil {
				break
			}
		}
	}
	go tstWriter(tststring)
	// Testing with proper PTY
	pass, err := GetPass("TestGetPass:", pty.Slave, buf)
	if err != nil {
		t.Errorf("GetPass(\"TestGetPass:\",pty.Slave,buf) failed: %v", err)
	}
	if string(pass) != tststring[:len(tststring)-1] {
		t.Errorf("GetPass got: %q want: %q", pass, tststring)
	}
	mu.Lock()
	if readbuffer.String() != "TestGetPass:" {
		t.Errorf("GetPass got: %q want: %q", readbuffer.String(), "TestGetPass:")
	}
	readbuffer.Reset()
	mu.Unlock()
	sbuf := buf[:10]
	tststring = "SuperSuperSuperSecret\n"
	go tstWriter(tststring)
	if _, err := GetPass("Pass: ", pty.Slave, sbuf); err == nil {
		t.Errorf("GetPass should fail got: <nil> want: ran out of buffespace")
	}
	// Make sure the buffer was cleared
	for i := 0; i < len(sbuf); i++ {
		if sbuf[i] != 0 {
			t.Errorf("GetPass should clear buffer on errors got: %q want: \"\"", sbuf)
			break
		}
	}
}

// TestGetChar tests out both the GetChar functions.
func TestGetChar(t *testing.T) {
	pty, err := OpenPTY()
	if err != nil {
		t.Fatalf("OpenPTY failed: %v", err)
	}
	tstWrite := func(in string) {
		w := []byte(in)
		var err error
		for tot, nr := 0, 0; tot < len(w); tot += nr {
			if nr, err = pty.Master.Write(w[tot:]); err != nil {
				break
			}
		}
	}
	tstring := "TestTheReader"
	go tstWrite(tstring)
	for idx, c := range []byte(tstring) {
		gc, err := pty.GetChar()
		if err != nil {
			t.Errorf("GetChar failed: %v", err)
		}
		if gc != c {
			t.Errorf("char at idx: %d does not match got: %b want: %b", idx, gc, c)
		}
	}
	tstIface := func(cr io.ByteReader, cmp string) error {
		for idx, c := range []byte(cmp) {
			gc, err := cr.ReadByte()
			if err != nil {
				t.Errorf("ReadByte failed: %v", err)
				return err
			}
			if gc != c {
				t.Errorf("char at idx: %d does not match got: %b want: %b", idx, gc, c)
			}
		}
		return nil
	}
	tstring = "TestTheByteReader"
	go tstWrite(tstring)
	if err := tstIface(pty, tstring); err != nil {
		t.Errorf("ByteReader failed: %v", err)
	}
	pty.Close()
	if _, err := pty.GetChar(); err == nil {
		t.Error("Reading from a closed PTY should fail want: io Error got: nil")
	}
}

// TestPTSName Gets name and tests if it's really a char device.
func TestPTSName(t *testing.T) {
	name, err := pty.PTSName()
	if err != nil {
		t.Error("PTS_name failed: ", err)
	}
	fi, err := os.Stat(name)
	if err != nil {
		t.Errorf("PTS_name failed to open slave named: %q err: %v", name, err)
	}
	mode := fi.Mode()
	if mode&os.ModeCharDevice == 0 {
		t.Errorf("PTS_name failed slave: %q is not a char device", name)
	}
}

// TestWinsz Tests if we can fetch the Terminal size.
// Also sanity checks with a normal file.
func TestWinsz(t *testing.T) {
	nf, err := donormfile("TestGetwinsz")
	if err != nil {
		t.Fatal("TestWinsz Could not create testfile, err: ", err)
	}
	defer nf.Close()
	wz = &tios.Wz
	// So putting some dummy values in to window size and see if they change
	// after reading the PTY size in
	wz.WsRow, wz.WsCol = 255, 255
	if err = tios.Winsz(pty.Slave); err != nil {
		t.Error("Winsz failed: ", err)
	}
	if wz.WsRow == 255 && wz.WsCol == 255 {
		wz.WsRow, wz.WsCol = 123, 124
		if err = tios.Winsz(pty.Slave); err != nil {
			t.Error("Getwinsz failed: ", err)
		}
		if wz.WsRow == 123 && wz.WsCol == 124 {
			t.Error("Winsz is not reading in values")
		}
	}
	if err = tios.Winsz(nf); err == nil {
		t.Error("Should not be able to read Windowsize from non pty")
	}
}

// TestSetwinsz Tests settig the terminal windowsize
func TestSetwinsz(t *testing.T) {
	nf, err := donormfile("TestSetwinsz")
	if err != nil {
		t.Fatal("TestSetwinsz Could not create testfile, err: ", err)
	}
	defer nf.Close()
	wz = &tios.Wz
	var i, y uint16
	for i = 0; i < NUMTEST; i++ {
		for y = 0; y < NUMTEST; y++ {
			wz.WsRow = i
			wz.WsCol = y
			wz.WsXpixel = uint16(rand.Int())
			wz.WsYpixel = uint16(rand.Int())
			xp, yp := wz.WsXpixel, wz.WsYpixel
			if err = tios.Setwinsz(pty.Slave); err != nil {
				t.Errorf("Setwinsz could not set row: %d col: %d err: %v", i, y, err)
				t.Errorf("Setwinsz could not set x: %d y: %d err: %v", xp, yp, err)
			}
			if err = tios.Winsz(pty.Slave); err != nil {
				t.Error("Winsz failed: ", err)
			}
			if !(wz.WsRow == i && wz.WsCol == y && wz.WsXpixel == xp && wz.WsYpixel == yp) {
				t.Errorf("Setwinsz got row: %d col: %d want row: %d col: %d", i, y, wz.WsRow, wz.WsCol)
				t.Errorf("Setwinsz got x: %d y: %d want x: %d y: %d", xp, yp, wz.WsXpixel, wz.WsYpixel)
			}
		}
	}
	if err = tios.Setwinsz(nf); err == nil {
		t.Error("Setwinsz should not work for std. file")
	}
}

// TestTraw Test to set terminal in raw mode.
func TestTraw(t *testing.T) {
	tios.Raw()
	if err := testraw(tios, "TestRaw"); err != nil {
		t.Errorf("TestTraw failed: %v", err)
	}
}

// TestTset Test the Tcsetraw function.
func TestTset(t *testing.T) {
	tios.Raw()
	err := tios.Set(pty.Slave)
	if err != nil {
		t.Fatal("Tset failed: ", err)
	}
	if err := testraw(tios, "TestTset"); err != nil {
		t.Errorf("TestTset failed: %v", err)
	}
	nf, err := donormfile("TestTset")
	if err != nil {
		t.Fatal("Could not create testfile: ", err)
	}
	defer nf.Close()
	err = tios.Set(nf)
	if err == nil {
		t.Error("Should not be able to set attributes on regular file: ", nf.Name())
	}
}

// Test to get the TC attributes.
func TestTattr(t *testing.T) {
	tios.Raw()
	// Set slave to raw
	err := tios.Set(pty.Slave)
	if err != nil {
		t.Fatal("Tattr failed: ", err)
	}
	// Read it back
	tios, err = Attr(pty.Slave)
	if err != nil {
		t.Fatal("Ttattr failed: ", err)
	}
	// Was it RAW?
	if err := testraw(tios, "TestTattr"); err != nil {
		t.Errorf("TestTattr failed: %v", err)
	}
	// Set it to Cooked mode
	tios.Cook()
	if err = tios.Set(pty.Slave); err != nil {
		t.Fatal("Failed to set Cook attributes")
	}
	if err := testcook(tios, "TestTattr"); err != nil {
		t.Errorf("TestTattr failed: %v", err)
	}
	// Should not be able to get/set attributes on non Terminals
	nf, err := donormfile("TestTattr")
	if err != nil {
		t.Fatal("TestTattr failed to create testfile, err: ", err)
	}
	defer nf.Close()
	err = tios.Set(nf)
	if err == nil {
		t.Error("Tattr , should not be able to set attributes on regular file: ", nf.Name())
	}
	tios, err = Attr(nf)
	if err == nil {
		t.Error("Tattr, should not be able to get attributes from regular file: ", nf.Name())
	}
}
