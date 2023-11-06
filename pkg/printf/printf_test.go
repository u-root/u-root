package printf_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/u-root/u-root/pkg/printf"
)

func TestPrintfBasic(t *testing.T) {

	args := func(args ...string) []string {
		return args
	}

	type testCase struct {
		a   []string
		e   string
		err error
	}

	cases := []testCase{
		{a: nil, err: printf.ErrNotEnoughArguments},
		{a: args("%j"), err: printf.ErrInvalidDirective},
		{a: args("hello"), e: "hello"},
		{a: args(`hello\n`), e: "hello\n"},
		{a: args(`hello\c there`), e: "hello"},
		{a: args(`\"hello\"`), e: `"hello"`},
		{a: args(`\\hello\\`), e: `\hello\`},
		{a: args(`hello\a`), e: "hello\a"},
		{a: args(`hello\b`), e: "hello\b"},
		// {a: args(`hello\e`), e: "hello\e"}, // TODO: figure this test out
		{a: args(`hello\f`), e: "hello\f"},
		{a: args(`hello\n`), e: "hello\n"},
		{a: args(`hello\r`), e: "hello\r"},
		{a: args(`hello\v`), e: "hello\v"},
		{a: args(`hello\123`), e: "helloS"},
		{a: args(`hello\x48`), e: "helloH"},
		{a: args(`hello\u123z`), e: "helloģz"},
		{a: args(`hello\u1230`), e: "helloሰ"},
		{a: args(`hello %%`), e: "hello %"},
		{a: args(`hello %s`, `\u1230`), e: `hello \u1230`},
		{a: args(`hello %b`, `\u1230`), e: "hello ሰ"},
	}

	for i, v := range cases {
		o := &bytes.Buffer{}
		_, err := printf.Fprintf(o, v.a...)
		ans := o.String()
		if v.err != nil {
			if err == nil {
				t.Errorf("case %d: exected err %s, got nil", i, v.err)
			}
			if !errors.Is(err, v.err) {
				t.Errorf("case %d: exected err %s, got %s", i, v.err, err)
			}
			continue
		}
		if err != nil {
			t.Errorf("case %d: exected err nil, got %s", i, err)
		}
		if v.e != ans {
			t.Errorf("case %d: exected '%s', got '%s'", i, v.e, ans)
		}
	}

}
