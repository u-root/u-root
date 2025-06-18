// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package menu

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/creack/pty"
	"github.com/u-root/u-root/pkg/testutil"
)

var inputDelay = 500 * time.Millisecond

func TestMain(m *testing.M) {
	SetInitialTimeout(inputDelay * 2)
	subsequentTimeout = inputDelay * 2

	os.Exit(m.Run())
}

type testEntry struct {
	mu         sync.Mutex
	label      string
	cmdline    string
	isDefault  bool
	load       error
	loadCalled bool
}

func (d *testEntry) Label() string {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.label
}

func (d *testEntry) String() string {
	return d.Label()
}

func (d *testEntry) Edit(f func(string) string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.cmdline = f(d.cmdline)
}

func (d *testEntry) Load() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.loadCalled = true
	return d.load
}

func (d *testEntry) Exec() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return nil
}

func (d *testEntry) LoadCalled() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.loadCalled
}

func (d *testEntry) IsDefault() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.isDefault
}

type testEntryStringer struct {
	testEntry
}

func (d *testEntryStringer) String() string {
	return d.Label() + " string"
}

func TestExtendedLabel(t *testing.T) {
	for _, tt := range []struct {
		name  string
		entry Entry
		want  string
	}{
		{
			name:  "without stringer",
			entry: &testEntry{label: "label"},
			want:  "label",
		},
		{
			name:  "with stringer",
			entry: &testEntryStringer{testEntry{label: "label"}},
			want:  "label string",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtendedLabel(tt.entry)
			if got != tt.want {
				t.Errorf("ExtendedLabel(%v) = %q; want %q", tt.entry, got, tt.want)
			}
		})
	}
}

var _ = MenuTerminal(&mockTerm{})

type mockTerm struct {
	inputSequence []ReadLine
	readLineCnt   int
}

func (m *mockTerm) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockTerm) Close() error                   { return nil }
func (m *mockTerm) SetPrompt(s string)             {}
func (m *mockTerm) SetEntryCallback(func())        {}
func (m *mockTerm) SetTimeout(time.Duration) error { return nil }
func (m *mockTerm) ReadLine() (string, error) {
	defer func() { m.readLineCnt++ }()
	if m.inputSequence == nil || m.readLineCnt >= len(m.inputSequence) {
		// mimic timeout
		return "", os.ErrDeadlineExceeded
	}
	return m.inputSequence[m.readLineCnt].string, m.inputSequence[m.readLineCnt].error
}

type ReadLine struct {
	string
	error
}

func TestChoose(t *testing.T) {
	// This test takes too long to run for the VM test and doesn't use
	// anything root-specific.
	testutil.SkipIfInVMTest(t)

	for _, tt := range []struct {
		name           string
		userEntry      []ReadLine
		wantedEntry    int
		editingAllowed bool
		expectedCmds   []string
	}{
		{
			name: "just_hit_enter",
			// user just hits enter.
			userEntry:   []ReadLine{{"", nil}},
			wantedEntry: -1, // expect nil
		},
		{
			name:        "hit_nothing",
			userEntry:   []ReadLine{},
			wantedEntry: -1, // expect nil
		},
		{
			name:        "hit_1",
			userEntry:   []ReadLine{{"1", nil}},
			wantedEntry: 1,
		},
		{
			name:        "hit_3",
			userEntry:   []ReadLine{{"3", nil}},
			wantedEntry: 3,
		},
		{
			name:        "out_of_bounds",
			userEntry:   []ReadLine{{"4", nil}},
			wantedEntry: -1, // expect nil
		},
		{
			name:        "not_a_number",
			userEntry:   []ReadLine{{"abc", nil}},
			wantedEntry: -1, // expect nil
		},
		{
			name:           "editing_allowed_override",
			userEntry:      getEditSequence(false, "1", "after"),
			editingAllowed: true,
			expectedCmds:   []string{"after", "before", "before"},
			wantedEntry:    -1, // expect nil
		},
		{
			name:           "editing_allowed_append",
			userEntry:      getEditSequence(true, "2", "after"),
			editingAllowed: true,
			expectedCmds:   []string{"before", "before after", "before"},
			wantedEntry:    -1, // expect nil
		},
		{
			name: "select_after_override",
			userEntry: append(getEditSequence(false, "1", "after"),
				ReadLine{"1", nil}),
			editingAllowed: true,
			expectedCmds:   []string{"after", "before", "before"},
			wantedEntry:    1,
		},
		{
			name:           "editing_not_allowed",
			userEntry:      getEditSequence(true, "1", "after"),
			editingAllowed: false,
			expectedCmds:   []string{"before", "before", "before"},
			wantedEntry:    1, // Edit attempt is parsed as a boot choice
		},
		{
			name:           "edit_fail_reading_1",
			userEntry:      errorOn(1, getEditSequence(false, "1", "after")),
			editingAllowed: true,
			expectedCmds:   []string{"before", "before", "before"},
			wantedEntry:    -1, // expect nil
		},
		{
			name:           "edit_fail_reading_2",
			userEntry:      errorOn(2, getEditSequence(false, "1", "after")),
			editingAllowed: true,
			expectedCmds:   []string{"before", "before", "before"},
			wantedEntry:    -1, // expect nil
		},
		{
			name:           "edit_fail_reading_3",
			userEntry:      errorOn(3, getEditSequence(false, "1", "after")),
			editingAllowed: true,
			expectedCmds:   []string{"before", "before", "before"},
			wantedEntry:    -1, // expect nil
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			entries := []*testEntry{
				{label: "1", cmdline: "before"},
				{label: "2", cmdline: "before"},
				{label: "3", cmdline: "before"},
			}
			var menu []Entry
			for _, e := range entries {
				menu = append(menu, e)
			}

			chosen := make(chan Entry)
			go func() {
				m := &mockTerm{
					inputSequence: tt.userEntry,
				}
				chosen <- Choose(m, tt.editingAllowed, menu...)
			}()

			chosenWant := Entry(nil)
			if tt.wantedEntry > 0 {
				chosenWant = menu[tt.wantedEntry-1] // 1 based index
			}
			if got := <-chosen; got != chosenWant {
				t.Errorf("Choose(%#v, %#v) = %#v, wantedEntry %#v", tt.userEntry, entries, got, tt.wantedEntry)
			}
			// Check for editing
			for i, entry := range entries {
				if i < len(tt.expectedCmds) && entry.cmdline != tt.expectedCmds[i] {
					t.Errorf("Entry %s got cmdline %s, wanted %s", entry.Label(), entry.cmdline, tt.expectedCmds[i])
				}
			}
		})
	}
}

func errorOn(index int, arr []ReadLine) []ReadLine {
	arr[index].error = errors.New("Expected test error")
	return arr
}

func contains(s []string, t string) bool {
	for _, u := range s {
		if u == t {
			return true
		}
	}
	return false
}

func TestShowMenuAndLoadFromFile(t *testing.T) {
	// This test takes too long to run for the VM test and doesn't use
	// anything root-specific.
	testutil.SkipIfInVMTest(t)

	tests := []struct {
		name      string
		entries   []*testEntry
		userEntry []byte

		// calledLabels are the entries for which Do was called.
		calledLabels []string
	}{
		{
			name: "default_entry",
			entries: []*testEntry{
				{label: "1", isDefault: true, load: nil},
				{label: "2", isDefault: true, load: nil},
			},
			// user just hits enter.
			userEntry:    []byte("\r\n"),
			calledLabels: []string{"1"},
		},
		{
			name: "non_default_entry_default",
			entries: []*testEntry{
				{label: "1", isDefault: false, load: nil},
				{label: "2", isDefault: true, load: nil},
				{label: "3", isDefault: true, load: nil},
			},
			// user just hits enter.
			userEntry:    []byte("\r\n"),
			calledLabels: []string{"2"},
		},
		{
			name: "non_default_entry_chosen_but_broken",
			entries: []*testEntry{
				{label: "1", isDefault: false, load: fmt.Errorf("borked")},
				{label: "2", isDefault: true, load: nil},
				{label: "3", isDefault: true, load: nil},
			},
			userEntry:    []byte("1\r\n"),
			calledLabels: []string{"1", "2"},
		},
		{
			name: "last_entry_works",
			entries: []*testEntry{
				{label: "1", isDefault: true, load: fmt.Errorf("foo")},
				{label: "2", isDefault: true, load: fmt.Errorf("bar")},
				{label: "3", isDefault: true, load: nil},
			},
			// user just hits enter.
			userEntry:    []byte("\r\n"),
			calledLabels: []string{"1", "2", "3"},
		},
		{
			name: "indecisive_entry",
			entries: []*testEntry{
				{label: "1", isDefault: true, load: nil},
				{label: "2", isDefault: true, load: nil},
				{label: "3", isDefault: true, load: nil},
			},
			// \x08 is the backspace character
			userEntry:    []byte("1\x082\r\n"),
			calledLabels: []string{"2"},
		},
		{
			name: "timeout_gets_first_default",
			entries: []*testEntry{
				{label: "1", isDefault: true, load: nil},
				{label: "2", isDefault: true, load: nil},
				{label: "3", isDefault: true, load: nil},
			},
			// No input
			userEntry:    []byte{},
			calledLabels: []string{"1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			master, slave, err := pty.Open()
			if err != nil {
				t.Fatalf("%v", err)
			}
			defer master.Close()
			defer slave.Close()

			var entries []Entry
			for _, e := range tt.entries {
				entries = append(entries, e)
			}

			timer := time.NewTimer(initialTimeout * 4)
			entry := make(chan Entry)
			go func() {
				entry <- showMenuAndLoadFromFile(slave, true, entries...)
			}()

			if len(tt.userEntry) > 0 {
				// We have to wait until Choose has actually started trying to read, as
				// ttys are asynchronous.
				//
				// Know a better way? Halp.
				time.Sleep(inputDelay)
				if _, err := master.Write(tt.userEntry); err != nil {
					t.Fatalf("failed to write new-line: %v", err)
				}
			}

			select {
			case <-timer.C:
				t.Errorf("Test %s timed out after %v", tt.name, initialTimeout)
			case got := <-entry:
				if want := tt.calledLabels[len(tt.calledLabels)-1]; got.Label() != want {
					t.Errorf("got label %s wantedEntry label %s", got.Label(), want)
				}

				for _, entry := range tt.entries {
					wantCalled := contains(tt.calledLabels, entry.label)
					if wantCalled != entry.LoadCalled() {
						t.Errorf("Entry %s gotCalled %t, wantCalled %t", entry.Label(), entry.LoadCalled(), wantCalled)
					}
				}
			}
		})
	}
}

func getEditSequence(add bool, bootnum string, cmdline string) []ReadLine {
	var editOpt string
	if add {
		editOpt = "a"
	} else {
		editOpt = "o"
	}

	return []ReadLine{
		{"e", nil},
		{bootnum, nil},
		{editOpt, nil},
		{cmdline, nil},
	}
}
