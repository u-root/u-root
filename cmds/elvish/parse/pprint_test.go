package parse

import (
	"testing"

	"github.com/u-root/u-root/cmds/elvish/tt"
)

var n = mustParse("ls $x[0]$y[1];echo")

var pprintASTTests = tt.Table{
	tt.Args(n).Rets(
		`Chunk
  Pipeline/Form
    Compound/Indexing/Primary Type=Bareword Value="ls"
    Compound
      Indexing
        Primary Type=Variable Value="x"
        Array/Compound/Indexing/Primary Type=Bareword Value="0"
      Indexing
        Primary Type=Variable Value="y"
        Array/Compound/Indexing/Primary Type=Bareword Value="1"
  Pipeline/Form/Compound/Indexing/Primary Type=Bareword Value="echo"
`),
}

func TestPPrintAST(t *testing.T) {
	tt.Test(t, tt.Fn("PPrintAST", PPrintAST), pprintASTTests)
}

var pprintParseTreeTests = tt.Table{
	tt.Args(n).Rets(
		`Chunk "ls $x[0]$y[1];echo" 0-18
  Pipeline/Form "ls $x[0]$y[1]" 0-13
    Compound/Indexing/Primary "ls" 0-2
    Sep " " 2-3
    Compound "$x[0]$y[1]" 3-13
      Indexing "$x[0]" 3-8
        Primary "$x" 3-5
        Sep "[" 5-6
        Array/Compound/Indexing/Primary "0" 6-7
        Sep "]" 7-8
      Indexing "$y[1]" 8-13
        Primary "$y" 8-10
        Sep "[" 10-11
        Array/Compound/Indexing/Primary "1" 11-12
        Sep "]" 12-13
  Sep ";" 13-14
  Pipeline/Form/Compound/Indexing/Primary "echo" 14-18
`),
}

func TestPprintParseTree(t *testing.T) {
	tt.Test(t, tt.Fn("PPrintParseTree", PPrintParseTree), pprintParseTreeTests)
}

func mustParse(src string) Node {
	n, err := Parse("[test]", src)
	if err != nil {
		panic(err)
	}
	return n
}
