// Package runtime assembles the Elvish runtime.
package runtime

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/u-root/u-root/cmds/core/elvish/eval"
	"github.com/u-root/u-root/cmds/core/elvish/store/storedefs"
	"github.com/u-root/u-root/cmds/core/elvish/util"
)

var logger = util.GetLogger("[runtime] ")

// InitRuntime initializes the runtime. The caller is responsible for calling
// CleanupRuntime at some point.
func InitRuntime(binpath, sockpath, dbpath string) (*eval.Evaler, string) {
	var dataDir string
	var err error

	// Determine data directory.
	dataDir, err = storedefs.EnsureDataDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "warning: cannot create data directory ~/.elvish")
	}

	// Determine runtime directory.
	_, err = getSecureRunDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot get runtime dir /tmp/elvish-$uid, falling back to data dir ~/.elvish:", err)
	}

	ev := eval.NewEvaler()
	ev.SetLibDir(filepath.Join(dataDir, "lib"))
	return ev, dataDir
}

// CleanupRuntime cleans up the runtime.
func CleanupRuntime(ev *eval.Evaler) {
	ev.Close()
}

var (
	ErrBadOwner      = errors.New("bad owner")
	ErrBadPermission = errors.New("bad permission")
)
