package program

import (
	"fmt"
	"testing"

	"github.com/u-root/u-root/cmds/core/elvish/program/shell"
)

var findProgramTests = []struct {
	args    []string
	checker func(Program) bool
}{
	{[]string{"-help"}, isShowHelp},
	{[]string{"-version"}, isShowVersion},
	{[]string{"-buildinfo"}, isShowBuildInfo},
	{[]string{}, isShell},
	{[]string{"-c", "echo"}, func(p Program) bool {
		return p.(*shell.Shell).Cmd
	}},
	{[]string{"-compileonly"}, func(p Program) bool {
		return p.(*shell.Shell).CompileOnly
	}},

	{[]string{"-bin", "/elvish"}, func(p Program) bool {
		return p.(*shell.Shell).BinPath == "/elvish"
	}},
	{[]string{"-db", "/db"}, func(p Program) bool {
		return p.(*shell.Shell).DbPath == "/db"
	}},
}

func isShowHelp(p Program) bool         { _, ok := p.(ShowHelp); return ok }
func isShowCorrectUsage(p Program) bool { _, ok := p.(ShowCorrectUsage); return ok }
func isShowVersion(p Program) bool      { _, ok := p.(ShowVersion); return ok }
func isShowBuildInfo(p Program) bool    { _, ok := p.(ShowBuildInfo); return ok }
func isShell(p Program) bool            { _, ok := p.(*shell.Shell); return ok }

func TestFindProgram(t *testing.T) {
	for i, test := range findProgramTests {
		f := parse(test.args)
		p := FindProgram(f)
		if !test.checker(p) {
			t.Errorf("test #%d (args = %q) failed", i, test.args)
		}
	}
}

func parse(args []string) *flagSet {
	f := newFlagSet()
	err := f.Parse(args)
	if err != nil {
		panic(fmt.Sprintln("bad flags in test", args))
	}
	return f
}
