package vals

import (
	"os"
	"testing"

	"github.com/u-root/u-root/cmds/elvish/tt"
)

type xtype int

func TestKind(t *testing.T) {
	tt.Test(t, tt.Fn("Kind", Kind), tt.Table{
		Args(true).Rets("bool"),
		Args("").Rets("string"),
		Args(EmptyList).Rets("list"),
		Args(EmptyMap).Rets("map"),
		Args(xtype(0)).Rets("!!vals.xtype"),

		Args(NewStruct(NewStructDescriptor(), nil)).Rets("map"),
		Args(NewFile(os.Stdin)).Rets("file"),
		Args(NewPipe(os.Stdin, os.Stdout)).Rets("pipe"),
	})
}
