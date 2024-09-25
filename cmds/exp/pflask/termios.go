// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/u-root/u-root/pkg/termios"
)

func raw() {
	// we don't set raw until the very last, so if they see an issue they can hit ^C
	t, err := termios.GetTermios(1)
	if err != nil {
		log.Fatal(err)
	}
	raw := termios.MakeRaw(t)
	if err = termios.SetTermios(1, raw); err != nil {
		log.Fatal(err)
	}
}
