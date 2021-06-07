package eddefs

import (
	"github.com/u-root/u-root/cmds/exp/felvish/edit/ui"
	"github.com/u-root/u-root/cmds/exp/felvish/eval"
)

// Mode is an editor mode.
type Mode interface {
	ModeLine() ui.Renderer
	Binding(ui.Key) eval.Callable
	Teardown()
}
