package eddefs

import (
	"github.com/u-root/u-root/cmds/core/elvish/edit/ui"
)

type ListingProvider interface {
	Len() int
	Show(i int) (string, ui.Styled)
	Filter(filter string) int
	Accept(i int, ed Editor)
	ModeTitle(int) string
}
