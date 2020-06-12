package store

import (
	"sync"

	"github.com/u-root/u-root/cmds/core/elvish/store/storedefs"
)

const (
	scoreDecay     = 0.986 // roughly 0.5^(1/50)
	scoreIncrement = 10
	scorePrecision = 6
)

type Dir struct {
	Score float64
	Name  string
}

type DirHistory struct {
	sync.Mutex
}

// we do not intend to store the directory history in a persisent manner any more.

func init() {
}

// AddDir adds a directory to the directory history.
func (*DirHistory) AddDir(d string, incFactor float64) error {
	return nil
}

// AddDir adds a directory and its score to history.
func (s *DirHistory) AddDirRaw(d string, score float64) error {
	return nil
}

// DelDir deletes a directory record from history.
func (s *DirHistory) DelDir(d string) error {
	return nil
}

// Dirs lists all directories in the directory history whose names are not
// in the blocklist. The results are ordered by scores in descending order.
func (s *DirHistory) Dirs(blocklist map[string]struct{}) ([]storedefs.Dir, error) {
	var dirs []storedefs.Dir
	return dirs, nil
}

type dirList []storedefs.Dir

func (dl dirList) Len() int {
	return len(dl)
}

func (dl dirList) Less(i, j int) bool {
	return dl[i].Score < dl[j].Score
}

func (dl dirList) Swap(i, j int) {
	dl[i], dl[j] = dl[j], dl[i]
}
