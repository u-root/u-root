package edit

import (
	"src.elv.sh/pkg/cli"
	"src.elv.sh/pkg/cli/mode"
	"src.elv.sh/pkg/cli/tk"
	"src.elv.sh/pkg/eval"
	"src.elv.sh/pkg/parse"
)

//elvdoc:var -instant:binding
//
// Binding for the instant mode.

//elvdoc:fn -instant:start
//
// Starts the instant mode. In instant mode, any text entered at the command
// line is evaluated immediately, with the output displayed.
//
// **WARNING**: Beware of unintended consequences when using destructive
// commands. For example, if you type `sudo rm -rf /tmp/*` in the instant mode,
// Elvish will attempt to evaluate `sudo rm -rf /` when you typed that far.

func initInstant(ed *Editor, ev *eval.Evaler, nb eval.NsBuilder) {
	bindingVar := newBindingVar(emptyBindingsMap)
	bindings := newMapBindings(ed, ev, bindingVar)
	nb.AddNs("-instant",
		eval.NsBuilder{
			"binding": bindingVar,
		}.AddGoFns("<edit:-instant>:", map[string]interface{}{
			"start": func() { instantStart(ed.app, ev, bindings) },
		}).Ns())
}

func instantStart(app cli.App, ev *eval.Evaler, bindings tk.Bindings) {
	execute := func(code string) ([]string, error) {
		outPort, collect, err := eval.StringCapturePort()
		if err != nil {
			return nil, err
		}
		err = ev.Eval(
			parse.Source{Name: "[instant]", Code: code},
			eval.EvalCfg{
				Ports:     []*eval.Port{nil, outPort},
				Interrupt: eval.ListenInterrupts})
		return collect(), err
	}
	w, err := mode.NewInstant(app,
		mode.InstantSpec{Bindings: bindings, Execute: execute})
	if w != nil {
		app.SetAddon(w, false)
		app.Redraw()
	}
	if err != nil {
		app.Notify(err.Error())
	}
}
