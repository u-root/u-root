// Package build contains build information.
//
// Build information should be set during compilation by passing
// -ldflags "-X github.com/u-root/u-root/cmds/core/elvish/build.Var=value" to "go build" or
// "go get".
package build

var (
	Version = "unknown"
	Builder = "unknown"
)
