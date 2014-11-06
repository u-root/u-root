// From docker.io. See the Apache License in this directory
// +build !arm

package netlink

import (
	"math/rand"
)

func randIfrDataByte() int8 {
	return int8(rand.Intn(255))
}
