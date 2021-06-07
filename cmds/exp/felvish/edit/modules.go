package edit

import (
	"github.com/u-root/u-root/cmds/exp/felvish/edit/completion"
	"github.com/u-root/u-root/cmds/exp/felvish/edit/history"
	"github.com/u-root/u-root/cmds/exp/felvish/edit/lastcmd"
	"github.com/u-root/u-root/cmds/exp/felvish/edit/location"
	"github.com/u-root/u-root/cmds/exp/felvish/edit/prompt"
	"github.com/u-root/u-root/cmds/exp/felvish/eval"
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
