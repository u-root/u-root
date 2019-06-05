package store

import (
	"encoding/binary"
	"fmt"
	"strings"
	"sync"

	"github.com/u-root/u-root/cmds/core/elvish/store/storedefs"
)

// The numbering is 1-based, not 0-based. Oh well.
type CmdHistory struct {
	sync.Mutex
	cmds map[int]string
	cur  int
	max  int
}

func (history *CmdHistory) Seek(i int) (int, error) {
	history.Lock()
	defer history.Unlock()
	history.cur = history.cur + i
	switch {
	case history.cur < 1:
		history.cur = 1
	}
	if history.cur > history.max {
		history.max = history.cur
	}
	return history.cur, nil
}

// convenience functions
func (history *CmdHistory) Cur() (int, error) {
	return history.Seek(0)
}

func (history *CmdHistory) Prev() (int, error) {
	return history.Seek(-1)
}

func (history *CmdHistory) Next() (int, error) {
	return history.Seek(1)
}

// AddCmd adds a new command to the command history.
func (history *CmdHistory) Add(cmd string) (int, error) {
	history.Lock()
	defer history.Unlock()
	history.cur = len(history.cmds) + 1
	history.cmds[history.cur] = cmd
	if history.cur > history.max {
		history.max = history.cur
	}
	return history.cur, nil
}

// Remove removes a command from command history referenced by
// sequence.
func (history *CmdHistory) Remove(seq int) error {
	history.Lock()
	defer history.Unlock()
	delete(history.cmds, seq)
	return nil
}

// Cmd queries the command history item with the specified sequence number.
func (history *CmdHistory) One(seq int) (string, error) {
	history.Lock()
	defer history.Unlock()
	c, ok := history.cmds[seq]
	if !ok {
		return "", storedefs.ErrNoMatchingCmd
	}
	return c, nil
}

// IterateCmds iterates all the commands in the specified range, and calls the
// callback with the content of each command sequentially.
func (history *CmdHistory) Walk(from, upto int, f func(string) bool) error {
	history.Lock()
	defer history.Unlock()
	for i := from; i < upto; i++ {
		v, ok := history.cmds[i]
		if !ok {
			continue
		}
		if !f(v) {
			return fmt.Errorf("%s fails", v)
		}
	}

	return nil
}

// Cmds returns the contents of all commands within the specified range.
func (history *CmdHistory) List(from, upto int) ([]string, error) {
	history.Lock()
	defer history.Unlock()
	var list []string
	for i := from; i < upto; i++ {
		v, ok := history.cmds[i]
		if !ok {
			continue
		}
		list = append(list, v)
	}
	return list, nil
}

// Search finds the first command after the given sequence number (inclusive)
// with the given prefix.
func (history *CmdHistory) Search(from int, prefix string) (int, string, error) {
	l, _ := history.List(from, history.max)
	for i, v := range l {
		if strings.HasPrefix(v, prefix) {
			return i + from, v, nil
		}
	}

	return 0, "", storedefs.ErrNoMatchingCmd
}

// PrevCmd finds the last command before the given sequence number (exclusive)
// with the given prefix.
func (history *CmdHistory) RSearch(upto int, prefix string) (int, string, error) {
	l, _ := history.List(1, upto)
	for i := range l {
		cmd := l[len(l)-i-1]
		if strings.HasPrefix(cmd, prefix) {
			return upto - i - 1, cmd, nil
		}
	}
	return 0, "", storedefs.ErrNoMatchingCmd
}

func marshalSeq(seq uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(seq))
	return b
}

func unmarshalSeq(key []byte) uint64 {
	return binary.BigEndian.Uint64(key)
}

func NewCmdHistory(s ...string) storedefs.Store {
	c := &CmdHistory{cmds: make(map[int]string)}
	for _, v := range s {
		c.Add(v)
	}
	return c
}
