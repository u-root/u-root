// +build !cgo
// +build windows

package group

import (
	"fmt"
	"runtime"
)

func init() {
	implemented = false
}

func Current() (*Group, error) {
	return nil, fmt.Errorf("group: Current not implemented on %s/%s", runtime.GOOS, runtime.GOARCH)
}

func Lookup(groupname string) (*Group, error) {
	return nil, fmt.Errorf("group: Lookup not implemented on %s/%s", runtime.GOOS, runtime.GOARCH)
}

func LookupId(string) (*Group, error) {
	return nil, fmt.Errorf("group: LookupId not implemented on %s/%s", runtime.GOOS, runtime.GOARCH)
}
