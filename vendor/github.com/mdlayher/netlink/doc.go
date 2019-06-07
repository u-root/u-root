// Package netlink provides low-level access to Linux netlink sockets.
//
// If you have any questions or you'd like some guidance, please join us on
// Gophers Slack (https://invite.slack.golangbridge.org) in the #networking
// channel!
//
//
// Debugging
//
// This package supports rudimentary netlink connection debugging support.
// To enable this, run your binary with the NLDEBUG environment variable set.
// Debugging information will be output to stderr with a prefix of "nl:".
//
// To use the debugging defaults, use:
//
//   $ NLDEBUG=1 ./nlctl
//
// To configure individual aspects of the debugger, pass key/value options such
// as:
//
//   $ NLDEBUG=level=1 ./nlctl
//
// Available key/value debugger options include:
//
//   level=N: specify the debugging level (only "1" is currently supported)
package netlink

//go:generate dot netlink.dot -T svg -o netlink.svg
