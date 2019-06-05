package eval

import (
	"errors"
	"os"

	"github.com/u-root/u-root/cmds/core/elvish/util"
)

// Filesystem.

var ErrStoreNotConnected = errors.New("store not connected")

func init() {
	addBuiltinFns(map[string]interface{}{
		// Directory
		"cd":         cd,
		"tilde-abbr": tildeAbbr,
	})
}

func cd(fm *Frame, args ...string) error {
	var dir string
	switch len(args) {
	case 0:
		dir = mustGetHome("")
	case 1:
		dir = args[0]
	default:
		return ErrArgs
	}

	return fm.Chdir(dir)
}

func dirs(fm *Frame) error {
	return nil
}

func tildeAbbr(path string) string {
	return util.TildeAbbr(path)
}

func isDir(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.Mode().IsDir()
}
