//go:build ignore

package main

type asm struct {
	name string
	cc   []string
	ld   []string
	code string
}

var asms = []asm{
	{
		name: "to32bit_3000",
		cc:   []string{"x86_64-linux-gnu-gcc", "-c", "-nostdlib", "-nostdinc", "-static"},
		ld:   []string{"ld", "-N", "-e entry64", "-Ttext=0x3000"},
		code: `
1: jmp 1b
`,
	},
}
