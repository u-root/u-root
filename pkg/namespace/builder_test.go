// Copyright 2020-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package namespace

type args struct {
	ns Namespace
}

type noopNS struct{}

func (m *noopNS) Bind(new string, old string, option mountflag) error        { return nil }
func (m *noopNS) Mount(servername, old, spec string, option mountflag) error { return nil }
func (m *noopNS) Unmount(new string, old string) error                       { return nil }
func (m *noopNS) Import(host string, remotepath string, mountpoint string, options mountflag) error {
	return nil
}
func (m *noopNS) Clear() error            { return nil }
func (m *noopNS) Chdir(path string) error { return nil }
