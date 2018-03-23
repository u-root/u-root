// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "github.com/u-root/u-root/pkg/sos"

func main() {
	sos.StartServer(sos.NewSosService())
}
