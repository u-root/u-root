// Package testevent holds events shared by guest and host.
package testevent

// ErrorEvent is an error.
type ErrorEvent struct {
	Binary string
	Error  string
}
