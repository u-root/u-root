package eval

import (
	"strconv"
	"syscall"

	"github.com/u-root/u-root/cmds/core/elvish/eval/vars"
)

var builtinNs = Ns{
	"pid":   vars.NewRo(strconv.Itoa(syscall.Getpid())),
	"paths": &EnvList{envName: "PATH"},
	"pwd":   PwdVariable{},
}

func addBuiltinFns(fns map[string]interface{}) {
	builtinNs.AddBuiltinFns("", fns)
}
