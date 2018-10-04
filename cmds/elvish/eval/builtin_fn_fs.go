package eval

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/u-root/u-root/cmds/elvish/eval/vals"
	"github.com/u-root/u-root/cmds/elvish/util"
)

// Filesystem.

var ErrStoreNotConnected = errors.New("store not connected")

func init() {
	addBuiltinFns(map[string]interface{}{
		// Directory
		"cd":          cd,
		"dir-history": dirs,

		// Path
		"path-abs":      filepath.Abs,
		"path-base":     filepath.Base,
		"path-clean":    filepath.Clean,
		"path-dir":      filepath.Dir,
		"path-ext":      filepath.Ext,
		"eval-symlinks": filepath.EvalSymlinks,
		"tilde-abbr":    tildeAbbr,

		// File types
		"-is-dir": isDir,
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

var dirDescriptor = vals.NewStructDescriptor("path", "score")

func newDirStruct(path string, score float64) *vals.Struct {
	return vals.NewStruct(dirDescriptor,
		[]interface{}{path, vals.FromGo(score)})
}

func dirs(fm *Frame) error {
	//	out := fm.ports[1].Chan
	//	for _, dir := range dirs {
	//		out <- newDirStruct(dir.Path, dir.Score)
	//	}
	return nil
}

func tildeAbbr(path string) string {
	return util.TildeAbbr(path)
}

func isDir(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.Mode().IsDir()
}
