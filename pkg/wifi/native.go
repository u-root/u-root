// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wifi

import "fmt"

type NativeWorker struct {
	Interface string
}

func NewNativeWorker(i string) (NativeWorker, error) {
	return NativeWorker{i}, nil
}

func (w *NativeWorker) Scan() ([]Option, error) {
	return nil, fmt.Errorf("Not Yet")
}

func (w *NativeWorker) GetID() (string, error) {
	return "", fmt.Errorf("Not Yet")
}

func (w *NativeWorker) Connect(a ...string) error {
	return fmt.Errorf("Not Yet")
}
