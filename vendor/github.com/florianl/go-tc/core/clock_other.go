//go:build !linux
// +build !linux

package core

func init() {
	clockFactor = 1.0
	tickInUSec = 1.0
}
