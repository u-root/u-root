package edit

import (
	"math"

	"github.com/u-root/u-root/cmds/core/elvish/eval"
	"github.com/u-root/u-root/cmds/core/elvish/eval/vars"
	"github.com/u-root/u-root/cmds/core/elvish/util"
)

func init() {
	atEditorInit(func(ed *editor, ns eval.Ns) {
		ed.maxHeight = math.Inf(1)
		ns["max-height"] = vars.FromPtr(&ed.maxHeight)
	})
}

func maxHeightToInt(h float64) int {
	if math.IsInf(h, 1) {
		return util.MaxInt
	}
	return int(h)
}
