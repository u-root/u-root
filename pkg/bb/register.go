package bb

import (
	"errors"
	"fmt"
	"os"
)

// ErrNotRegistered is returned by Run if the given command is not registered.
var ErrNotRegistered = errors.New("command not registered")

// Noop is a noop function.
var Noop = func() {}

type bbCmd struct {
	init, main func()
}

var bbcmds = map[string]bbCmd{}

// Register registers an init and main function for name.
func Register(name string, init, main func()) {
	if _, ok := bbcmds[name]; ok {
		panic(fmt.Sprintf("cannot register two commands with name %q", name))
	}
	bbcmds[name] = bbCmd{
		init: init,
		main: main,
	}
}

// Run runs the command with the given name.
//
// If the command's main exits without calling os.Exit, Run will exit with exit
// code 0.
func Run(name string) error {
	if _, ok := bbcmds[name]; !ok {
		return ErrNotRegistered
	}
	bbcmds[name].init()
	bbcmds[name].main()
	os.Exit(0)
	return nil
}
