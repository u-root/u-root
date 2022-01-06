// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

// ResetAllDevices resets all devices found in the path /sys/bus/pci*/devices/reset.
// While this is *techhically* only found on Linux, any system which puts
// that path in the namespace could work. Which, incidentally, demonstrates
// why putting stuff like this in the name space is portable and resilient.
// How could this work on a non-linux system? Pretend we're on plan 9:
// import linuxhost /dev /dev
// note this mount only applies to this process and its children
// run this reset code, and reset the linux devices.
func ResetAllDevices(n ...string) error {
	r, err := NewBusReader(n...)
	if err != nil {
		return err
	}
	d, err := r.Read()
	if err != nil {
		return err
	}
	return d.Reset()
}
