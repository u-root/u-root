package main

import(
	"testing"
)

func Test_wget(t *testing.T) {

	err := wget("http://example.com")
	if err != nil {
		t.Error(err)
	}
}
