// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pogosh

import "testing"

// The positive tests are expected to exit with code 123.
var shellPositiveTests = []struct {
	name string
	in   string
}{
	{
		"Exit",
		`exit 123`,
	},
	{
		"Exit2",
		`exit 123
		exit 0`,
	},
	{
		"Echo",
		`/bin/echo hello
		exit 123
		`,
	},
	/*{
			"If Statement",
			`if true; then
	exit 123
	done`,
		},
		{
			"If-Else Statement",
			`if false; then
		exit 124
	else
		exit 123
	done`,
		},
		{
			"Arithmetic",
			`exit $((2 + 010 + 0x10 + 5 * (1 + 2) + 3 / 2 + 7 % 4 - 8 + (2 << 3) + (7 >> 1) + \
	(1 < 2) + (3 > 4) + (80 <= 80) + (90 >= 90) + (7 == 7) + (8 == 8) + (5 != 3) + (2 & 3) + \
	(3 ^ 2) + (3 | 2) + (3 && 2) + (3 || 0) + (1 ? 2 : 3)))`, // TODO: add difference with 123
		},
		{
			"Arithmetic Assignment",
			`X=5
	Y=$(((X*=5) == 25 ? X : 0))
	exit $((X + Y + 193))`,
		},*/
}

func TestRunPositive(t *testing.T) {
	for _, tt := range shellPositiveTests {
		t.Run(tt.name, func(t *testing.T) {
			state := DefaultState()
			code, err := state.Run(tt.in)

			if err != nil {
				t.Error(err)
			} else {
				if code != 123 {
					t.Errorf("got %v, want 123", code)
				}
			}
		})
	}
}

// The negative tests are expected to return an error.
var shellNegativeTests = []struct {
	name string
	in   string
	err  string
}{
	{
		"Non-ASCII Characters",
		"echo hello\xbd",
		`<pogosh>:1:11: non-ascii character, '\xbd'`,
	},
	/*{
		"Division by Zero",
		"X=0; echo $((15/X))",
		"<pogosh>:1:10: division by zero, '15/0'",
	},*/
}

func TestRunNegative(t *testing.T) {
	for _, tt := range shellNegativeTests {
		t.Run(tt.name, func(t *testing.T) {
			state := DefaultState()
			_, err := state.Run(tt.in)

			errStr := "nil"
			if err != nil {
				errStr = err.Error()
			}

			if errStr != tt.err {
				t.Errorf("got \"%s\", want \"%s\"", errStr, tt.err)
			}
		})
	}
}

/*func ExampleRun() {
	state := DefaultState()
	state.Run(`echo Launching rocket...`)
	state.Run(`ACTION='BLASTOFF!!!'`)
	state.Run(`for T in $(seq 10 -1 1); do echo -n -- "T-$T "; done`)
	state.Run(`echo "$ACTION"`)
	// Output:
	// Launching rocket...
	// 10 9 8 7 6 5 4 3 2 1 BLASTOFF!!!
}

func ExampleRun_parallel() {

}

func ExampleRunInteractive() {
	state := DefaultState()
	state.Prompt = func () { return "> " }
	for {
		code, err := state.RunInteractive()
		if err == nil {
			os.Exit(code)
		}
		fmt.Println("Error:", err)
	}
}

func ExampleExec() {
	DefaultState().Run("example.sh")
	if err != nil {
		fmt.Println("Error:", err)
	}
}*/
