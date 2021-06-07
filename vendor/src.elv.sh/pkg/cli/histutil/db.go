package histutil

import (
	"strings"

	"src.elv.sh/pkg/store"
)

// DB is the interface of the storage database.
type DB interface {
	NextCmdSeq() (int, error)
	AddCmd(cmd string) (int, error)
	CmdsWithSeq(from, upto int) ([]store.Cmd, error)
	PrevCmd(upto int, prefix string) (store.Cmd, error)
	NextCmd(from int, prefix string) (store.Cmd, error)
}

// FaultyInMemoryDB is an in-memory DB implementation that can be injected
// one-off errors. It is useful in tests.
type FaultyInMemoryDB interface {
	DB
	// SetOneOffError causes the next operation on the database to return the
	// given error.
	SetOneOffError(err error)
}

// NewFaultyInMemoryDB creates a new FaultyInMemoryDB with the given commands.
func NewFaultyInMemoryDB(cmds ...string) FaultyInMemoryDB {
	return &testDB{cmds: cmds}
}

// Implementation of FaultyInMemoryDB.
type testDB struct {
	cmds        []string
	oneOffError error
}

func (s *testDB) SetOneOffError(err error) {
	s.oneOffError = err
}

func (s *testDB) error() error {
	err := s.oneOffError
	s.oneOffError = nil
	return err
}

func (s *testDB) NextCmdSeq() (int, error) {
	return len(s.cmds), s.error()
}

func (s *testDB) AddCmd(cmd string) (int, error) {
	if s.oneOffError != nil {
		return -1, s.error()
	}
	s.cmds = append(s.cmds, cmd)
	return len(s.cmds) - 1, nil
}

func (s *testDB) CmdsWithSeq(from, upto int) ([]store.Cmd, error) {
	if err := s.error(); err != nil {
		return nil, err
	}
	if from < 0 {
		from = 0
	}
	if upto < 0 || upto > len(s.cmds) {
		upto = len(s.cmds)
	}
	var cmds []store.Cmd
	for i := from; i < upto; i++ {
		cmds = append(cmds, store.Cmd{Text: s.cmds[i], Seq: i})
	}
	return cmds, nil
}

func (s *testDB) PrevCmd(upto int, prefix string) (store.Cmd, error) {
	if s.oneOffError != nil {
		return store.Cmd{}, s.error()
	}
	if upto < 0 || upto > len(s.cmds) {
		upto = len(s.cmds)
	}
	for i := upto - 1; i >= 0; i-- {
		if strings.HasPrefix(s.cmds[i], prefix) {
			return store.Cmd{Text: s.cmds[i], Seq: i}, nil
		}
	}
	return store.Cmd{}, store.ErrNoMatchingCmd
}

func (s *testDB) NextCmd(from int, prefix string) (store.Cmd, error) {
	if s.oneOffError != nil {
		return store.Cmd{}, s.error()
	}
	if from < 0 {
		from = 0
	}
	for i := from; i < len(s.cmds); i++ {
		if strings.HasPrefix(s.cmds[i], prefix) {
			return store.Cmd{Text: s.cmds[i], Seq: i}, nil
		}
	}
	return store.Cmd{}, store.ErrNoMatchingCmd
}
