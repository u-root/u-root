// Package group allows user group lookups by name or id.
package group

import (
	"strconv"
)

var implemented = true // set to false by lookup_stub.go's init

// Group represents a user group.
type Group struct {
	Gid  string // group id
	Name string
	Members []string
}

type UnknownGroupIdError int

func (e UnknownGroupIdError) Error() string {
	return "group: unknown groupid " + strconv.Itoa(int(e))
}

type UnknownGroupError string

func (e UnknownGroupError) Error() string {
	return "group: unknown group " + string(e)
}
