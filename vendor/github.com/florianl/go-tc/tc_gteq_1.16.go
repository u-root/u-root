//go:build go1.16
// +build go1.16

package tc

import (
	"io/ioutil"
	"log"
)

func setDummyLogger() *log.Logger {
	return log.New(ioutil.Discard, "", 0)
}
