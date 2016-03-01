// From docker.io. See the Apache License in this directory
// +build amd64

package netlink

import (
	"math/rand"
)

func randIfrDataByte() int8 {
	return int8(rand.Intn(255))
}
