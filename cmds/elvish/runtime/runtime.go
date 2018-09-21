// Package runtime assembles the Elvish runtime.
package runtime

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/u-root/u-root/cmds/elvish/eval"
	"github.com/u-root/u-root/cmds/elvish/eval/re"
	"github.com/u-root/u-root/cmds/elvish/eval/str"
	"github.com/u-root/u-root/cmds/elvish/store/storedefs"
	"github.com/u-root/u-root/cmds/elvish/util"
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
	} else {
		if dbpath == "" {
			dbpath = filepath.Join(dataDir, "db")
		}
	}

	// Determine runtime directory.
	runDir, err := getSecureRunDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot get runtime dir /tmp/elvish-$uid, falling back to data dir ~/.elvish:", err)
		runDir = dataDir
	}
	if sockpath == "" {
		sockpath = filepath.Join(runDir, "sock")
	}

	ev := eval.NewEvaler()
	ev.SetLibDir(filepath.Join(dataDir, "lib"))
	ev.InstallModule("re", re.Ns)
	ev.InstallModule("str", str.Ns)
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
