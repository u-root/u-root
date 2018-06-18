// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import ()

type DummyUpspinService struct {
	Configured	bool
	User	string
	Dir		string
	Store	string
	Seed	string
}

func (us *DummyUpspinService) ToggleFlag() {
  us.Configured = !us.Configured
}

func (us *DummyUpspinService) SetConfig(new UpspinAcctJsonMsg) error {
  us.User = new.User
  us.Dir = new.Dir
  us.Store = new.Store
  us.Seed = new.Seed
	return nil
}

func NewDummyUpspinService() (*DummyUpspinService, error) {
  return &DummyUpspinService {
    Configured:  false,
    User:	 			 "",
    Dir:			 	 "",
    Store:			 "",
    Seed: 			 "",
  }, nil
}
