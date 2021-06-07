package store

import (
	"github.com/u-root/u-root/cmds/exp/felvish/eval"
	"github.com/u-root/u-root/cmds/exp/felvish/store/storedefs"
)

func Ns(s storedefs.Store) eval.Ns {
	return eval.NewNs().AddBuiltinFns("store:", map[string]interface{}{
		"del-dir": s.Remove,
	})
}
