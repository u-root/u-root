// +build !windows,!plan9

package tty

import (
	"os"
	"testing"
	"time"

	"github.com/u-root/u-root/cmds/elvish/edit/ui"
)

// timeout is the longest time the tests wait between writing something on
// the writer and reading it from the reader before declaring that the
// reader has a bug.
const timeoutInterval = 100 * time.Millisecond

func timeout() <-chan time.Time {
	return time.After(timeoutInterval)
}

var (
	theWriter   *os.File
	innerReader *os.File
	theReader   *reader
)

func TestMain(m *testing.M) {
	r, w, err := os.Pipe()
	if err != nil {
		panic("os.Pipe returned error, something is seriously wrong")
	}
	defer r.Close()
	defer w.Close()
	theWriter = w
	innerReader = r
	theReader = newReader(r)
	theReader.Start()
	defer theReader.Stop()

	os.Exit(m.Run())
}

var keyTests = []struct {
	input string
	want  Event
}{
	// Simple graphical key.
	{"x", KeyEvent{'x', 0}},
	{"X", KeyEvent{'X', 0}},
	{" ", KeyEvent{' ', 0}},

	// Ctrl key.
	{"\001", KeyEvent{'A', ui.Ctrl}},
	{"\033", KeyEvent{'[', ui.Ctrl}},

	// Ctrl-ish keys, but not thought as Ctrl keys by our reader.
	{"\n", KeyEvent{'\n', 0}},
	{"\t", KeyEvent{'\t', 0}},
	{"\x7f", KeyEvent{'\x7f', 0}}, // backspace

	// Alt plus simple graphical key.
	{"\033a", KeyEvent{'a', ui.Alt}},
	{"\033[", KeyEvent{'[', ui.Alt}},

	// G3-style key.
	{"\033OA", KeyEvent{ui.Up, 0}},
	{"\033OH", KeyEvent{ui.Home, 0}},

	// CSI-sequence key identified by the ending rune.
	{"\033[A", KeyEvent{ui.Up, 0}},
	{"\033[H", KeyEvent{ui.Home, 0}},
	// Test for all possible modifier
	{"\033[1;2A", KeyEvent{ui.Up, ui.Shift}},

	// CSI-sequence key with one argument, always ending in '~'.
	{"\033[1~", KeyEvent{ui.Home, 0}},
	{"\033[11~", KeyEvent{ui.F1, 0}},

	// CSI-sequence key with three arguments and ending in '~'. The first
	// argument is always 27, the second identifies the modifier and the last
	// identifies the key.
	{"\033[27;4;63~", KeyEvent{';', ui.Shift | ui.Alt}},
}

// FIXME: Fix test flakiness.
func DISABLED_TestKey(t *testing.T) {
	for _, test := range keyTests {
		theWriter.WriteString(test.input)
		select {
		case event := <-theReader.EventChan():
			if event != test.want {
				t.Errorf("Reader reads event %v, want %v", event, test.want)
			}
		case <-timeout():
			t.Errorf("Reader timed out")
		}
	}
}

// TestStopMakesUnderlyingFileAvailable tests that after calling Stop, the
// Reader no longer attempts to read from the underlying file, so it is
// available for use by others.
func TestStopMakesUnderlyingFileAvailable(t *testing.T) {
	theReader.Stop()
	defer theReader.Start()

	s := "lorem ipsum"
	theWriter.WriteString(s)
	gotChan := make(chan string)
	go func() {
		var buf [32]byte
		nr, err := innerReader.Read(buf[:])
		if err != nil {
			t.Errorf("inner.Read returns error: %v", err)
		}
		gotChan <- string(buf[:nr])
	}()
	select {
	case got := <-gotChan:
		if got != s {
			t.Errorf("got %q, want %q", got, s)
		}
	case <-time.After(time.Second):
		t.Error("inner.Read times out")
	}
}

// TestStartAfterStopIndeedStarts tests that calling Start very shortly after
// Stop puts the Reader the correct started state.
// FIXME: Fix test flakiness.
func DISABLED_TestStartAfterStopIndeedStarts(t *testing.T) {
	for i := 0; i < 100; i++ {
		theReader.Stop()
		theReader.Start()

		theWriter.WriteString("a")
		select {
		case event := <-theReader.EventChan():
			wantEvent := KeyEvent(ui.Key{'a', 0})
			if event != wantEvent {
				t.Errorf("After Stop and Start, Reader reads %v, want %v", event, wantEvent)
			}
		case <-timeout():
			t.Errorf("After Stop and Start, Reader timed out")
		}
	}
}
