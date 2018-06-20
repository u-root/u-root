// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"strings"
	"regexp"
	"io/ioutil"
)

var (
	upspinConfigDir	= fmt.Sprintf("%v/upspin/config", os.Getenv("HOME"))
)



type DummyUpspinService struct {
	Configured bool
	User       string
	Dir        string
	Store      string
	Seed       string
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
	return &DummyUpspinService{
		Configured: false,
		User:       "",
		Dir:        "",
		Store:      "",
		Seed:       "",
	}, nil
}



type UpspinService struct  {
	Configured bool
	User       string
	Dir        string
	Store      string
	Seed       string
}

func getFileData(path string) map[string]string {
	userData := make(map[string]string)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		// start in unconfigured mode using empty map
		return userData
	}
	// parse file data into map using regex
	reg, err := regexp.Compile("(: )(.*,|)")
	for _, s := range strings.Split(string(b), "\n") {
		keyval := reg.Split(s, -1)
		if len(keyval) == 2 {
			userData[keyval[0]] = keyval[1]
		}
	}
	return userData
}

func (us *UpspinService) SetConfig(new UpspinAcctJsonMsg) error {
	us.User = new.User
	us.Dir = new.Dir
	us.Store = new.Store
	us.Seed = new.Seed
	return nil
}

func (us *UpspinService) Update() {
	data := getFileData(upspinConfigDir)
	us.Configured = true
	us.User       = data["username"]
	us.Dir        = data["dirserver"]
	us.Store      = data["storeserver"]
	us.Seed       = ""
}

func (us *UpspinService) ToggleFlag() {
	us.Configured = !us.Configured
}

func NewUpspinService() (*UpspinService, error) {
	data := getFileData(upspinConfigDir)
	return &UpspinService{
		Configured: true,
		User:       data["username"],
		Dir:        data["dirserver"],
		Store:      data["storeserver"],
		Seed:       "",
	}, nil
}
