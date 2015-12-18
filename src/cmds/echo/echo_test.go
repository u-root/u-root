package main

import (
	"testing"
)

func Test_echo(t *testing.T) {

	if err:= echo("Simple \ttest"); err != nil {
		t.Error(err)
	}
}
