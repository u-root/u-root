package shlex

import (
	"fmt"
	"testing"
)

var datanonposix = []struct {
	in  string
	out []string
	err error
}{
	{`This string has an embedded apostrophe, doesn't it?`,
		[]string{
			"This",
			"string",
			"has",
			"an",
			"embedded",
			"apostrophe",
			",",
			"doesn't",
			"it",
			"?",
		},
		nil,
	},
	{"This string has embedded \"double quotes\" and 'single quotes' in it,\nand even \"a 'nested example'\".\n",
		[]string{
			"This",
			"string",
			"has",
			"embedded",
			`"double quotes"`,
			"and",
			`'single quotes'`,
			"in",
			"it",
			",",
			"and",
			"even",
			`"a 'nested example'"`,
			".",
		},
		nil,
	},
	{`Hello world!, こんにちは　世界！`,
		[]string{
			"Hello",
			"world",
			"!",
			",",
			"こんにちは",
			"世界",
			"！",
		},
		nil,
	},
	{`Do"Not"Separate`,
		[]string{`Do"Not"Separate`},
		nil,
	},
	{`"Do"Separate`,
		[]string{`"Do"`, "Separate"},
		nil,
	},
	{`Escaped \e Character not in quotes`,
		[]string{
			"Escaped",
			`\`,
			"e",
			"Character",
			"not",
			"in",
			"quotes",
		},
		nil,
	},
	{`Escaped "\e" Character in double quotes`,
		[]string{
			"Escaped",
			`"\e"`,
			"Character",
			"in",
			"double",
			"quotes",
		},
		nil,
	},
	{`Escaped '\e' Character in single quotes`,
		[]string{
			"Escaped",
			`'\e'`,
			"Character",
			"in",
			"single",
			"quotes",
		},
		nil,
	},
	{`Escaped '\'' \"\'\" single quote`,
		[]string{
			"Escaped",
			`'\'`,
			`' \"\'`,
			`\`,
			`" single quote`,
		},
		ErrNoClosing,
	},
	{`Escaped "\"" \'\"\' double quote`,
		[]string{
			"Escaped",
			`"\"`,
			`" \'\"`,
			`\`,
			`' double quote`,
		},
		ErrNoClosing,
	},
	{`"'Strip extra layer of quotes'"`,
		[]string{`"'Strip extra layer of quotes'"`},
		nil,
	},
}

var dataposix = []struct {
	in  string
	out []string
	err error
}{
	{`This string has an embedded apostrophe, doesn't it?`,
		[]string{
			"This",
			"string",
			"has",
			"an",
			"embedded",
			"apostrophe",
			",",
			"doesnt it?",
		},
		ErrNoClosing,
	},
	{"This string has embedded \"double quotes\" and 'single quotes' in it,\nand even \"a 'nested example'\".\n",
		[]string{
			"This",
			"string",
			"has",
			"embedded",
			`double quotes`,
			"and",
			`single quotes`,
			"in",
			"it",
			",",
			"and",
			"even",
			`a 'nested example'`,
			".",
		},
		nil,
	},
	{`Hello world!, こんにちは　世界！`,
		[]string{
			"Hello",
			"world",
			"!",
			",",
			"こんにちは",
			"世界",
			"！",
		},
		nil,
	},
	{`Do"Not"Separate`,
		[]string{`DoNotSeparate`},
		nil,
	},
	{`"Do"Separate`,
		[]string{"DoSeparate"},
		nil,
	},
	{`Escaped \e Character not in quotes`,
		[]string{
			"Escaped",
			"e",
			"Character",
			"not",
			"in",
			"quotes",
		},
		nil,
	},
	{`Escaped "\e" Character in double quotes`,
		[]string{
			"Escaped",
			`\e`,
			"Character",
			"in",
			"double",
			"quotes",
		},
		nil,
	},
	{`Escaped '\e' Character in single quotes`,
		[]string{
			"Escaped",
			`\e`,
			"Character",
			"in",
			"single",
			"quotes",
		},
		nil,
	},
	{`Escaped '\'' \"\'\" single quote`,
		[]string{
			"Escaped",
			`\ \"\"`,
			"single",
			"quote",
		},
		nil,
	},
	{`Escaped "\"" \'\"\' double quote`,
		[]string{
			"Escaped",
			`"`,
			`'"'`,
			"double",
			"quote",
		},
		nil,
	},
	{`"'Strip extra layer of quotes'"`,
		[]string{`'Strip extra layer of quotes'`},
		nil,
	},
}

func TestSplitNonPOSIX(t *testing.T) {
	testSplit(t, false)
}

func TestSplitPOSIX(t *testing.T) {
	testSplit(t, true)
}

func testSplit(t *testing.T, posix bool) {
	var data []struct {
		in  string
		out []string
		err error
	}
	if posix {
		data = dataposix
	} else {
		data = datanonposix
	}

	for _, d := range data {
		t.Logf("Spliting: `%s'", d.in)

		result, err := NewLexerString(d.in, posix, false).Split()

		// check closing and escaped error
		if err != d.err {
			printToken(result)
			t.Fatalf("Error expected: `%v', but result catched: `%v'",
				d.err, err)
		}

		// check splited number
		if len(result) != len(d.out) {
			printToken(result)
			t.Fatalf("Split expeced: `%d', but result founds: `%d'",
				len(d.out), len(result))
		}

		// check words
		for j, out := range d.out {
			if result[j] != out {
				printToken(result)
				t.Fatalf("Word expeced: `%s', but result founds: `%s' in %d",
					out, result[j], j)
			}
		}
		t.Log("ok")
	}
}

func printToken(s []string) {
	for _, token := range s {
		fmt.Println(token)
	}
}
