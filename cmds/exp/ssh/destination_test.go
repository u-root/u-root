package main

import (
	guser "os/user"
	"testing"
)

type d struct {
	input string
	user  string
	host  string
	port  string
}

var tests = []d{
	d{"example.org", "", "example.org", "22"},
	d{"foo@example.org", "foo", "example.org", "22"},
	d{"foo@example.org", "foo", "example.org", "22"},
	d{"ssh://192.168.0.2:23", "", "192.168.0.2", "23"},
	d{"ssh://x@example.org", "x", "example.org", "22"},
}

func TestParseDest(t *testing.T) {
	for _, x := range tests {
		if x.user == "" {
			var u *guser.User
			u, _ = guser.Current()
			x.user = u.Username
		}
		user, host, port, err := parseDest(x.input)
		if err != nil {
			t.Fatal(err)
		}
		if user != x.user || host != x.host || port != x.port {
			t.Fatalf("failed on %v: got %v, %v, %v", x, user, host, port)
		}
	}
}
