package highlight

import (
	"github.com/u-root/u-root/cmds/core/elvish/edit/ui"
	"github.com/u-root/u-root/cmds/core/elvish/parse"
)

// Semantically applied styles.
var (
	styleForGoodCommand  = ui.Styles{"green"}
	styleForBadCommand   = ui.Styles{"red"}
	styleForGoodVariable = ui.Styles{"magenta"}
)

// Lexically applied styles.

// ui.Styles for Primary nodes.
var styleForPrimary = map[parse.PrimaryType]ui.Styles{
	parse.Bareword:     {},
	parse.SingleQuoted: {"yellow"},
	parse.DoubleQuoted: {"yellow"},
	parse.Variable:     styleForGoodVariable,
	parse.Wildcard:     {},
	parse.Tilde:        {},
}

var styleForComment = ui.Styles{"cyan"}

// ui.Styles for Sep nodes.
var styleForSep = map[string]string{
	">":  "green",
	">>": "green",
	"<":  "green",
	"?>": "green",
	"|":  "green",

	"(": "bold",
	")": "bold",

	"&": "bold",
}
