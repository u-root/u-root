package store

import (
	"testing"

	"github.com/u-root/u-root/cmds/core/elvish/store/storedefs"
)

var (
	dirsToAdd  = []string{"/usr/local", "/usr", "/usr/bin", "/usr"}
	black      = map[string]struct{}{"/usr/local": {}}
	wantedDirs = []storedefs.Dir{
		{"/usr", scoreIncrement*scoreDecay*scoreDecay + scoreIncrement},
		{"/usr/bin", scoreIncrement * scoreDecay}}
	dirToDel           = "/usr"
	wantedDirsAfterDel = []storedefs.Dir{
		{"/usr/bin", scoreIncrement * scoreDecay}}
)

func testDir(t *testing.T) {
	tStore := NewCmdHistory()
	for _, path := range dirsToAdd {
		_, err := tStore.Add(path) //, 1)
		if err != nil {
			t.Errorf("tStore.Add(%q) => %v, want <nil>", path, err)
		}
	}

	/* The whole Store interface is a total clusterfuck. Fix later.
	dirs, err := tStore.List(black)
	if err != nil || !reflect.DeepEqual(dirs, wantedDirs) {
		t.Errorf(`tStore.ListDirs() => (%v, %v), want (%v, <nil>)`,
			dirs, err, wantedDirs)
	}

	tStore.DelDir("/usr")
	dirs, err = tStore.List(black)
	if err != nil || !reflect.DeepEqual(dirs, wantedDirsAfterDel) {
		t.Errorf(`After DelDir("/usr"), tStore.ListDirs() => (%v, %v), want (%v, <nil>)`,
			dirs, err, wantedDirsAfterDel)
	}
	*/
}
