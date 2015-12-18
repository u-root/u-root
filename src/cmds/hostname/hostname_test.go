package main

import (
	"fmt"
	"testing"
)

func Test_hostname(t *testing.T) {

	if err := hostname(); err != nil {
		t.Error(err)
	}

}
