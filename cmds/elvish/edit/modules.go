package edit

import (
	"github.com/u-root/u-root/cmds/elvish/edit/completion"
	"github.com/u-root/u-root/cmds/elvish/edit/history"
	"github.com/u-root/u-root/cmds/elvish/edit/lastcmd"
	"github.com/u-root/u-root/cmds/elvish/edit/location"
	"github.com/u-root/u-root/cmds/elvish/edit/prompt"
	"github.com/u-root/u-root/cmds/elvish/eval"
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
