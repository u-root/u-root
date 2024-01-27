//go:build tools

package vmtest

// List u-root commands that need to be in go.mod & go.sum to be buildable as
// dependencies. This way, they aren't eliminated by `go mod tidy`.
//
// But obviously aren't actually importable, since they are main packages.
import (
	_ "github.com/u-root/u-root/cmds/core/cat"
	_ "github.com/u-root/u-root/cmds/core/dhclient"
	_ "github.com/u-root/u-root/cmds/core/false"
	_ "github.com/u-root/u-root/cmds/core/gosh"
	_ "github.com/u-root/u-root/cmds/core/init"
	_ "github.com/u-root/u-root/cmds/core/ip"
	_ "github.com/u-root/u-root/cmds/core/ls"
	_ "github.com/u-root/u-root/cmds/core/shutdown"
	_ "github.com/u-root/u-root/cmds/core/sync"
	_ "github.com/u-root/u-root/cmds/core/wget"
	_ "github.com/u-root/u-root/cmds/exp/pxeserver"
)
