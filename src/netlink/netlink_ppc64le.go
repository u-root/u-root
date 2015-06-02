// From docker.io. See the Apache License in this directory
// Wonder what they want to do for arm?
// +build ppc64le

package netlink

import (
	"math/rand"
)

func randIfrDataByte() uint8 {
	return uint8(rand.Intn(255))
}
