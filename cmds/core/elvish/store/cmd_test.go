package store

import (
	"testing"

	"github.com/u-root/u-root/cmds/core/elvish/store/storedefs"
)

var (
	cmds     = []string{"echo foo", "put bar", "put lorem", "echo bar"}
	searches = []struct {
		next      bool
		seq       int
		prefix    string
		wantedSeq int
		wantedCmd string
		wantedErr error
	}{
		{false, 5, "echo", 4, "echo bar", nil},
		{false, 5, "put", 3, "put lorem", nil},
		{false, 4, "echo", 1, "echo foo", nil},
		{false, 3, "f", 0, "", storedefs.ErrNoMatchingCmd},
		{false, 1, "", 0, "", storedefs.ErrNoMatchingCmd},

		{true, 1, "echo", 1, "echo foo", nil},
		{true, 1, "put", 2, "put bar", nil},
		{true, 2, "echo", 4, "echo bar", nil},
		{true, 4, "put", 0, "", storedefs.ErrNoMatchingCmd},
	}
)

func TestCmd(t *testing.T) {
	tStore := NewCmdHistory()
	startSeq, err := tStore.Next()
	if startSeq != 1 || err != nil {
		t.Errorf("tStore.NextCmdSeq() => (%v, %v), want (1, nil)",
			startSeq, err)
	}
	for i, cmd := range cmds {
		wantSeq := startSeq + i
		seq, err := tStore.Add(cmd)
		if seq != wantSeq || err != nil {
			t.Errorf("tStore.Add(%v) => (%v, %v), want (%v, nil)",
				cmd, seq, err, wantSeq)
		}
	}
	endSeq, err := tStore.Next()
	wantedEndSeq := startSeq + len(cmds)
	if endSeq != wantedEndSeq || err != nil {
		t.Errorf("tStore.Next() => (%v, %v), want (%v, nil)",
			endSeq, err, wantedEndSeq)
	}
	for i, wantedCmd := range cmds {
		seq := i + startSeq
		cmd, err := tStore.One(seq)
		if cmd != wantedCmd || err != nil {
			t.Errorf("tStore.Cmd(%v) => (%v, %v), want (%v, nil)",
				seq, cmd, err, wantedCmd)
		}
	}
	for _, tt := range searches {
		f := tStore.RSearch
		funcname := "tStore.RSearch"
		if tt.next {
			f = tStore.Search
			funcname = "tStore.Search"
		}
		seq, cmd, err := f(tt.seq, tt.prefix)
		if seq != tt.wantedSeq || cmd != tt.wantedCmd || err != tt.wantedErr {
			t.Errorf("%s(%v, %v) in (%v), => (%v, %v, %v), want (%v, %v, %v)",
				funcname, tt.seq, tt.prefix,
				tStore,
				seq, cmd, err,
				tt.wantedSeq, tt.wantedCmd, tt.wantedErr)
		}
	}

	if err := tStore.Remove(1); err != nil {
		t.Error("Failed to remove cmd")
	}
	if seq, err := tStore.One(1); err != storedefs.ErrNoMatchingCmd {
		t.Errorf("Cmd(1) => (%v, %v), want (%v, %v)",
			seq, err, "", storedefs.ErrNoMatchingCmd)
	}
}
