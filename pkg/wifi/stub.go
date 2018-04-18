// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wifi

var _ = WiFi(&StubWorker{})

type StubWorker struct {
	Options []Option
	ID      string
}

func (w *StubWorker) Scan() ([]Option, error) {
	return w.Options, nil
}

func (w *StubWorker) GetID() (string, error) {
	return w.ID, nil
}

func (*StubWorker) Connect(a ...string) error {
	return nil
}

func NewStubWorker(id string, options ...Option) (WiFi, error) {
	return &StubWorker{ID: id, Options: options}, nil
}
