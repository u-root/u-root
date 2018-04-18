// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wifi

import (
	"fmt"
	"syscall"
	"unsafe"
)

type NativeWorker struct {
	Interface string
	FD        int
	Range     IWRange
}

func NewNativeWorker(i string) (WiFi, error) {
	s, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_IP)
	if err != nil {
		return nil, err
	}
	return &NativeWorker{FD: s, Interface: i}, nil
}

func (w *NativeWorker) Scan() ([]Option, error) {
	return nil, fmt.Errorf("Not Yet")
}

func (w *NativeWorker) GetID() (string, error) {
	return "", fmt.Errorf("Not Yet")
}

func (w *NativeWorker) Connect(a ...string) error {
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(w.FD), SIOCGIWRANGE, uintptr(unsafe.Pointer(&w.Range)))
	return err
}
