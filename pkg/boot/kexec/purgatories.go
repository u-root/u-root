//go:build ignore

package main

type asm struct {
	name string
	args []string
	code string
}

var asms = []asm{
	{
		name: "to32bit_3000",
		args: []string{"x86_64-linux-gnu-gcc", "-Ttext=0x3000"},
		code: `
1: jmp 1b
`,
	},
}
