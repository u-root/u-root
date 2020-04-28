// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pogosh

import (
	"reflect"
	"testing"
)

// The positive tests are expected to pass lexing.
var lexerPositiveTests = []struct {
	name string
	in   string
	out  []token
}{
	// Words
	{"Single Word",
		`abc`,
		[]token{{"abc", ttWord}},
	},
	{"Single Character Word",
		`a`,
		[]token{{"a", ttWord}},
	},
	{"Multiple Words",
		`a bc`,
		[]token{{"a", ttWord}, {"bc", ttWord}},
	},

	// Newlines and blanks
	{"Empty input",
		"",
		[]token{},
	},
	{"Blank input",
		"\t  \t  ",
		[]token{},
	},
	{"Various whitespace",
		"  \n\t\n \n",
		[]token{{"\n", ttNewLine}, {"\n", ttNewLine}, {"\n", ttNewLine}},
	},
	{"Whitespace Word Combo",
		" a \nb\t\nc \nd",
		[]token{
			{"a", ttWord}, {"\n", ttNewLine}, {"b", ttWord}, {"\n", ttNewLine},
			{"c", ttWord}, {"\n", ttNewLine}, {"d", ttWord}},
	},

	// Single quotes
	{"Single quotes",
		"'a'",
		[]token{{"'a'", ttWord}},
	},
	{"Single quotes with spaces",
		"'a bc'",
		[]token{{"'a bc'", ttWord}},
	},
	{"Single escape",
		"'\t\\\n'",
		[]token{{"'\t\\\n'", ttWord}},
	},

	// Double quotes
	{"Double Quote",
		`"a"`,
		[]token{{`"a"`, ttWord}},
	},
	{"Double Quotes with spaces",
		`"a bc"`,
		[]token{{`"a bc"`, ttWord}},
	},
	{"Double Quote With Escapes",
		`"\"a\\\$(" ")"`,
		[]token{{`"\"a\\\$("`, ttWord}, {`")"`, ttWord}},
	},
	{"Double Quote With Subexpression",
		`"a$(b)"`,
		[]token{{`"a$(b)"`, ttWord}},
	},
	// TODO: The following two rules aren't a regular language.
	//{"Double Quote With Subexpression Nested",
	//	`"a$(b "c")"`,
	//	[]token{{`"a$(b "c")"`, ttWord}},
	//},
	//{"Double Quote With Subexpression Double Nested",
	//	`"a$(b "c$(d "e")")"`,
	//	[]token{{`"a$(b "c$(d "e")")"`, ttWord}},
	//},

	// Command substitution
	// TODO

	// Line comments
	{"Line Comment",
		`# "comment"`,
		[]token{},
	},
	{"Multiple Line Comments",
		`abc d # comment
efg # comment2`,
		[]token{{"abc", ttWord}, {"d", ttWord}, {"\n", ttNewLine}, {"efg", ttWord}},
	},
	{"Line Comment Continuation",
		"# comment \\\nabc",
		[]token{},
	},

	// Line continuation
	{"Line Continuation",
		"\\\n",
		[]token{},
	},
	{"Line Continuation in Token",
		"ech\\\no abc",
		[]token{{"ech\\\no", ttWord}, {"abc", ttWord}},
	},
	{"Line Continuation in Double Quote",
		"echo \"ab\\\nc\"",
		[]token{{"echo", ttWord}, {"\"ab\\\nc\"", ttWord}},
	},
	{"Line Continuation in Single Quote",
		"echo 'ab\\\nc'",
		[]token{{"echo", ttWord}, {"'ab\\\nc'", ttWord}},
	},

	// Operators
	{"Operator Single",
		">>",
		[]token{{">>", ttDGreat}},
	},
	{"Operators Compact",
		"a&&b||c;;d<<e>>f<&g>&h<>i<<-j>|k",
		[]token{{"a", ttWord}, {"&&", ttAndIf}, {"b", ttWord},
			{"||", ttOrIf}, {"c", ttWord}, {";;", ttDSemi}, {"d", ttWord},
			{"<<", ttDLess}, {"e", ttWord}, {">>", ttDGreat}, {"f", ttWord},
			{"<&", ttLessAnd}, {"g", ttWord}, {">&", ttGreatAnd}, {"h", ttWord},
			{"<>", ttLessGreat}, {"i", ttWord}, {"<<-", ttDLessDash},
			{"j", ttWord}, {">|", ttClobber}, {"k", ttWord}},
	},
	{"Operators Whitespace",
		" a && b || c ;; d << e >> f <& g >& h <> i <<- j >| k ",
		[]token{{"a", ttWord}, {"&&", ttAndIf}, {"b", ttWord},
			{"||", ttOrIf}, {"c", ttWord}, {";;", ttDSemi}, {"d", ttWord},
			{"<<", ttDLess}, {"e", ttWord}, {">>", ttDGreat}, {"f", ttWord},
			{"<&", ttLessAnd}, {"g", ttWord}, {">&", ttGreatAnd}, {"h", ttWord},
			{"<>", ttLessGreat}, {"i", ttWord}, {"<<-", ttDLessDash},
			{"j", ttWord}, {">|", ttClobber}, {"k", ttWord}},
	},
	{"Operators Escaped",
		` a \&\& b \|\| c \;\; d \<\< e \>\> f \<\& g \>\& h \<\> i \<\<\- j \>\| k `,
		[]token{{"a", ttWord}, {`\&\&`, ttWord}, {"b", ttWord},
			{`\|\|`, ttWord}, {"c", ttWord}, {`\;\;`, ttWord}, {"d", ttWord},
			{`\<\<`, ttWord}, {"e", ttWord}, {`\>\>`, ttWord}, {"f", ttWord},
			{`\<\&`, ttWord}, {"g", ttWord}, {`\>\&`, ttWord}, {"h", ttWord},
			{`\<\>`, ttWord}, {"i", ttWord}, {`\<\<\-`, ttWord},
			{"j", ttWord}, {`\>\|`, ttWord}, {"k", ttWord}},
	},
	{"Operator Dash",
		`\<\<- -`,
		[]token{{`\<\<-`, ttWord}, {`-`, ttWord}},
	},
	{"Operator Half Escape",
		`echo \&&`,
		[]token{{"echo", ttWord}, {`\&`, ttWord}, {"&", ttWord}},
	},

	// Examples from POSIX.1-2017
	// TODO: these tests require some work
	/*{"POSIX Example 1",
			`a=1
	set 2
	echo ${a}b-$ab-${1}0-${10}-$10
	`,
			[]string{`a=1`, `set`, `2`, `echo`, `${a}b-$ab-${1}0-${10}-$10`},
		},
		{"POSIX Example 2",
			`foo=asdf
	echo ${foo-bar}xyz}
	foo=
	echo ${foo-bar}xyz}
	unset foo
	echo ${foo-bar}xyz}`,
			[]string{},
		},
		{"POSIX Example 3",
			`# repeat a command 100 times
	x=100
	while [ $x -gt 0 ]
	do
		command x=$(($x-1))
	done`,
			[]string{`x=100`, `while`, `[`, `$x`, `-gt`, `0`, `]`, `do`, `command`, `x=$(($x-1))`, `done`},
		},*/
}

func TestLexerPositive(t *testing.T) {
	for _, tt := range lexerPositiveTests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := tokenize(tt.in)

			if err != nil {
				t.Error(err)
			} else {
				if !reflect.DeepEqual(tokens, tt.out) {
					t.Errorf("got %v, want %v", tokens, tt.out)
				}
			}
		})
	}
}

// Every input is tested with and without a trailing newline
func TestLexerPositiveTrailingNewline(t *testing.T) {
	for _, tt := range lexerPositiveTests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := tokenize(tt.in + "\n")

			if err != nil {
				t.Error(err)
			} else if len(tokens) == 0 || tokens[len(tokens)-1].ttype != ttNewLine {
				t.Errorf("expected %v to have a trailing newline", tokens)
			} else {
				tokens = tokens[:len(tokens)-1]
				if !reflect.DeepEqual(tokens, tt.out) {
					t.Errorf("got %v, want %v", tokens, tt.out)
				}
			}
		})
	}
}
