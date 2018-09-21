package eddefs

import (
	"github.com/u-root/u-root/cmds/elvish/edit/ui"
	"github.com/u-root/u-root/cmds/elvish/eval"
)

// Mode is an editor mode.
type Mode interface {
	ModeLine() ui.Renderer
	Binding(ui.Key) eval.Callable
	Teardown()
}
