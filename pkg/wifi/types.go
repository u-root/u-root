// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wifi

type Option struct {
	Essid     string
	AuthSuite SecProto
}

type WiFi interface {
	Scan() ([]Option, error)
	GetID() (string, error)
	Connect(a ...string) error
}
