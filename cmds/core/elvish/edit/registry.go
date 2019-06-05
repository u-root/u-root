package edit

import (
	"github.com/u-root/u-root/cmds/core/elvish/eval"
)

// This file contains utilities that facilitates modularization of the editor.

var editorInitFuncs []func(*editor, eval.Ns)

func atEditorInit(f func(*editor, eval.Ns)) {
	editorInitFuncs = append(editorInitFuncs, f)
}
