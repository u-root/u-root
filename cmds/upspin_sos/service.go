// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	upspinConfigDir = fmt.Sprintf("%v/upspin", os.Getenv("HOME"))
	upspinKeyDir    = fmt.Sprintf("%v/.ssh", os.Getenv("HOME"))
)

type UpspinService struct {
	Configured bool
	User       string
	Dir        string
	Store      string
	Seed       string
}

func getFileData(path string) map[string]string {
	userData := make(map[string]string)
	b, err := ioutil.ReadFile(fmt.Sprintf("%v/config", path))
	if err != nil {
		// start in unconfigured mode using empty map
		return userData
	}
	// regex for finding key-val separator ": [remote,]" and port ":443"
	splitpoint, err := regexp.Compile("(: )(.*,|)")
	port, err := regexp.Compile("(:443)")
	for _, s := range strings.Split(string(b), "\n") {
		s := port.ReplaceAllString(s, "")
		keyval := splitpoint.Split(s, -1)
		if len(keyval) == 2 {
			userData[keyval[0]] = keyval[1]
		}
	}
	return userData
}

func (us UpspinService) setFileData(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0777); err != nil {
			return err
		}
	}
	f, err := os.Create(fmt.Sprintf("%v/config", path))
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString(fmt.Sprintf("username: %v\n", us.User))
	// hardcoded default server prefix and suffix
	f.WriteString(fmt.Sprintf("dirserver: remote,%v:443\n", us.Dir))
	f.WriteString(fmt.Sprintf("storeserver: remote,%v:443\n", us.Store))
	// hardcoded packing security
	f.WriteString("packing: ee")
	return nil
}

func (us UpspinService) setKeys(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0777); err != nil {
			return err
		}
	}
	err := exec.Command("upspin", "keygen", fmt.Sprintf("-secretseed=%v", us.Seed), path).Start()
	if err != nil {
		return err
	}
	return nil
}

func (us *UpspinService) Update() {
	data := getFileData(upspinConfigDir)
	us.User = data["username"]
	us.Dir = data["dirserver"]
	us.Store = data["storeserver"]
}

func (us *UpspinService) ToggleFlag() {
	us.Configured = !us.Configured
}

func (us *UpspinService) SetConfig(new UpspinAcctJsonMsg) error {
	us.User = new.User
	us.Dir = new.Dir
	us.Store = new.Store
	us.Seed = new.Seed
	if err := us.setFileData(upspinConfigDir); err != nil {
		return err
	}
	if err := us.setKeys(fmt.Sprintf("%v/%v", upspinKeyDir, us.User)); err != nil {
		return err
	}
	return nil
}

func NewUpspinService() (*UpspinService, error) {
	data := getFileData(upspinConfigDir)
	config := false
	if len(data) > 0 {
		config = true
	}
	return &UpspinService{
		Configured: config,
		User:       data["username"],
		Dir:        data["dirserver"],
		Store:      data["storeserver"],
		Seed:       "",
	}, nil
}
