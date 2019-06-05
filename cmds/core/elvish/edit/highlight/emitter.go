package highlight

import (
	"strings"

	"github.com/u-root/u-root/cmds/core/elvish/edit/ui"
	"github.com/u-root/u-root/cmds/core/elvish/eval"
	"github.com/u-root/u-root/cmds/core/elvish/parse"
)

type Emitter struct {
	GoodFormHead func(string) bool
	AddStyling   func(begin, end int, style string)
}

func (e *Emitter) EmitAll(n parse.Node) {
	switch n := n.(type) {
	case *parse.Form:
		e.form(n)
	case *parse.Primary:
		e.primary(n)
	case *parse.Sep:
		e.sep(n)
	}
	for _, child := range n.Children() {
		e.EmitAll(child)
	}
}

func (e *Emitter) form(n *parse.Form) {
	for _, an := range n.Assignments {
		if an.Left != nil && an.Left.Head != nil {
			v := an.Left.Head
			e.AddStyling(v.Begin(), v.End(), styleForGoodVariable.String())
		}
	}
	for _, cn := range n.Vars {
		if len(cn.Indexings) > 0 && cn.Indexings[0].Head != nil {
			v := cn.Indexings[0].Head
			e.AddStyling(v.Begin(), v.End(), styleForGoodVariable.String())
		}
	}
	if n.Head != nil {
		e.formHead(n.Head)
	}
}

func (e *Emitter) formHead(n *parse.Compound) {
	head, err := eval.PurelyEvalCompound(n)
	st := ui.Styles{}
	if err == nil {
		if e.GoodFormHead(head) {
			st = styleForGoodCommand
		} else {
			st = styleForBadCommand
		}
	} else if err != eval.ErrImpure {
		st = styleForBadCommand
	}
	if len(st) > 0 {
		e.AddStyling(n.Begin(), n.End(), st.String())
	}
}

func (e *Emitter) primary(n *parse.Primary) {
	e.AddStyling(n.Begin(), n.End(), styleForPrimary[n.Type].String())
}

func (e *Emitter) sep(n *parse.Sep) {
	septext := n.SourceText()
	switch {
	case strings.TrimSpace(septext) == "":
		// Don't do anything. Whitespaces don't get any styling.
	case strings.HasPrefix(septext, "#"):
		// Comment.
		e.AddStyling(n.Begin(), n.End(), styleForComment.String())
	default:
		e.AddStyling(n.Begin(), n.End(), styleForSep[septext])
	}
}
