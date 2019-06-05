package edit

import "github.com/u-root/u-root/cmds/core/elvish/edit/eddefs"

func (ed *editor) SetAction(action eddefs.Action) {
	if ed.nextAction == noAction {
		ed.nextAction = action
	}
}

func (ed *editor) popAction() eddefs.Action {
	action := ed.nextAction
	ed.nextAction = noAction
	return action
}
