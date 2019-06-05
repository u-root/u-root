package edit

import (
	"github.com/u-root/u-root/cmds/core/elvish/edit/completion"
	"github.com/u-root/u-root/cmds/core/elvish/edit/history"
	"github.com/u-root/u-root/cmds/core/elvish/edit/lastcmd"
	"github.com/u-root/u-root/cmds/core/elvish/edit/location"
	"github.com/u-root/u-root/cmds/core/elvish/edit/prompt"
	"github.com/u-root/u-root/cmds/core/elvish/eval"
)

func init() {
	atEditorInit(func(ed *editor, ns eval.Ns) {
		location.Init(ed, ns)
		lastcmd.Init(ed, ns)
		history.Init(ed, ns)
		completion.Init(ed, ns)
		prompt.Init(ed, ns)
	})
}
