// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package menu

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/google/goterm/term"
)

type dummyEntry struct {
	mu        sync.Mutex
	label     string
	isDefault bool
	do        error
	called    bool
}

func (d *dummyEntry) Label() string {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.label
}

func (d *dummyEntry) String() string {
	return d.Label()
}

func (d *dummyEntry) Do() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.called = true
	return d.do
}

func (d *dummyEntry) Called() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.called
}

func (d *dummyEntry) IsDefault() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.isDefault
}

func TestChoose(t *testing.T) {
	entry1 := &dummyEntry{label: "1"}
	entry2 := &dummyEntry{label: "2"}
	entry3 := &dummyEntry{label: "3"}

	for _, tt := range []struct {
		name      string
		entries   []Entry
		userEntry []byte
		want      Entry
	}{
		{
			name:    "just_hit_enter",
			entries: []Entry{entry1, entry2, entry3},
			// user just hits enter.
			userEntry: []byte("\r\n"),
			want:      nil,
		},
		{
			name:      "hit_nothing",
			entries:   []Entry{entry1, entry2, entry3},
			userEntry: nil,
			want:      nil,
		},
		{
			name:      "hit_1",
			entries:   []Entry{entry1, entry2, entry3},
			userEntry: []byte("1\r\n"),
			want:      entry1,
		},
		{
			name:      "hit_3",
			entries:   []Entry{entry1, entry2, entry3},
			userEntry: []byte("3\r\n"),
			want:      entry3,
		},
		{
			name:    "tentative_hit_1",
			entries: []Entry{entry1, entry2, entry3},
			// \x08 is the backspace character.
			userEntry: []byte("2\x081\r\n"),
			want:      entry1,
		},
		{
			name:      "out_of_bounds",
			entries:   []Entry{entry1, entry2, entry3},
			userEntry: []byte("4\r\n"),
			want:      nil,
		},
		{
			name:      "not_a_number",
			entries:   []Entry{entry1, entry2, entry3},
			userEntry: []byte("abc\r\n"),
			want:      nil,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pty, err := term.OpenPTY()
			if err != nil {
				t.Fatalf("%v", err)
			}
			defer pty.Close()

			chosen := make(chan Entry)
			go func() {
				chosen <- Choose(pty.Slave, tt.entries...)
			}()

			// Well this sucks.
			//
			// We have to wait until Choose has actually started trying to read, as
			// ttys are asynchronous.
			//
			// Know a better way? Halp.
			time.Sleep(1 * time.Second)

			if tt.userEntry != nil {
				if _, err := pty.Master.Write(tt.userEntry); err != nil {
					t.Fatalf("failed to write new-line: %v", err)
				}
			}

			if got := <-chosen; got != tt.want {
				t.Errorf("Choose(%#v, %#v) = %#v, want %#v", tt.userEntry, tt.entries, got, tt.want)
			}
		})
	}
}

func contains(s []string, t string) bool {
	for _, u := range s {
		if u == t {
			return true
		}
	}
	return false
}

func TestShowMenuAndBoot(t *testing.T) {
	tests := []struct {
		name      string
		entries   []*dummyEntry
		userEntry []byte

		// calledLabels are the entries for which Do was called.
		calledLabels []string
	}{
		{
			name: "default_entry",
			entries: []*dummyEntry{
				{label: "1", isDefault: true, do: errStopTestOnly},
				{label: "2", isDefault: true, do: nil},
			},
			// user just hits enter.
			userEntry:    []byte("\r\n"),
			calledLabels: []string{"1"},
		},
		{
			name: "non_default_entry_default",
			entries: []*dummyEntry{
				{label: "1", isDefault: false, do: errStopTestOnly},
				{label: "2", isDefault: true, do: errStopTestOnly},
				{label: "3", isDefault: true, do: nil},
			},
			// user just hits enter.
			userEntry:    []byte("\r\n"),
			calledLabels: []string{"2"},
		},
		{
			name: "non_default_entry_chosen_but_broken",
			entries: []*dummyEntry{
				{label: "1", isDefault: false, do: fmt.Errorf("borked")},
				{label: "2", isDefault: true, do: errStopTestOnly},
				{label: "3", isDefault: true, do: nil},
			},
			userEntry:    []byte("1\r\n"),
			calledLabels: []string{"1", "2"},
		},
		{
			name: "last_entry_works",
			entries: []*dummyEntry{
				{label: "1", isDefault: true, do: nil},
				{label: "2", isDefault: true, do: nil},
				{label: "3", isDefault: true, do: errStopTestOnly},
			},
			// user just hits enter.
			userEntry:    []byte("\r\n"),
			calledLabels: []string{"1", "2", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pty, err := term.OpenPTY()
			if err != nil {
				t.Fatalf("%v", err)
			}
			defer pty.Close()

			var entries []Entry
			for _, e := range tt.entries {
				entries = append(entries, e)
			}

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				ShowMenuAndBoot(pty.Slave, entries...)
				wg.Done()
			}()

			// Well this sucks.
			//
			// We have to wait until Choose has actually started trying to read, as
			// ttys are asynchronous.
			//
			// Know a better way? Halp.
			time.Sleep(1 * time.Second)

			if tt.userEntry != nil {
				if _, err := pty.Master.Write(tt.userEntry); err != nil {
					t.Fatalf("failed to write new-line: %v", err)
				}
			}

			wg.Wait()

			for _, entry := range tt.entries {
				wantCalled := contains(tt.calledLabels, entry.label)
				if wantCalled != entry.Called() {
					t.Errorf("Entry %s gotCalled %t, wantCalled %t", entry.Label(), entry.Called(), wantCalled)
				}
			}
		})
	}
}
