package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type cmd func(...string) error

var (
	// various commands add themselves to this map as part of
	// their init.
	commands    map[string]cmd
	lpcdebug    = flag.Bool("lpcdebug", true, "Enable lpc debug prints")
	chips       = make(map[string]func(ioport, ioaddr, time.Duration, time.Duration, debugf) ec)
	defaultChip = flag.String("chip", "lpc", "Which chip to use")
	chip        = NewLPC
)

func debug(s string, v ...interface{}) {
	fmt.Printf(s, v...)
}

func main() {
	d := debug
	flag.Parse()
	if !*lpcdebug {
		d = nil
	}
	p, err := NewDevPorts(d)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
	}
	fmt.Printf("p %v\n", p)
	if c, ok := chips[*defaultChip]; !ok {
		fmt.Fprintf(os.Stderr, "Unknown chip %v: Choices: %v\n", *defaultChip, chips)
		os.Exit(1)
	} else {
		chip = c
	}

	ec := chip(p, EC_LPC_ADDR_HOST_CMD, time.Second*10, time.Second*10, d)
	// valid command?
	// TODO: use the command table for real? But what should the type be? interface{}, err?
	a := flag.Args()
	if len(a) == 0 {
		fmt.Printf("usage: ectool command [args]\n")
		os.Exit(1)
	}
	switch a[0] {
	case "info":
		d, err := info(ec)
		fmt.Printf("%v, %v\n", d, err)
	default:
		fmt.Printf("Unknown: %v", a)
	}
}
