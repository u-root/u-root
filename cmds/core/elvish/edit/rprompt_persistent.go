package edit

import (
	"github.com/u-root/u-root/cmds/core/elvish/eval"
	"github.com/u-root/u-root/cmds/core/elvish/eval/vars"
)

func init() {
	atEditorInit(func(ed *editor, ns eval.Ns) {
		ed.RpromptPersistent = false
		ns["rprompt-persistent"] = vars.FromPtr(&ed.RpromptPersistent)
	})
}
