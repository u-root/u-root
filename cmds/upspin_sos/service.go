// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	upspinConfigDir = flag.String("configdir", filepath.Join(os.Getenv("HOME"), "upspin"), "path for Upspin config file")
	upspinKeyDir    = flag.String("keydir", filepath.Join(os.Getenv("HOME"), ".ssh"), "path for username directory to hold key files")
)

type UpspinService struct {
	Configured bool
	User       string
	Dir        string
	Store      string
	Seed       string
}

func makeUserDirectories(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return err
		}
		return filepath.Walk(dir, func(name string, info os.FileInfo, err error) error {
			if err == nil {
				err = os.Chown(name, 1000, 1000)
			}
			return err
		})
	}
	return nil
}

func getFileData(path string) map[string]string {
	userData := make(map[string]string)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		// start in unconfigured mode using empty map
		return userData
	}
	// regex for finding key-val separator ": [remote,]" and port ":443"
	splitpoint := regexp.MustCompile("(: )(.*,|)")
	port := regexp.MustCompile("(:443)")
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
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString(fmt.Sprintf("username: %v\n", us.User))
	// hardcoded default server prefix and suffix
	f.WriteString(fmt.Sprintf("dirserver: remote,%v:443\n", us.Dir))
	f.WriteString(fmt.Sprintf("storeserver: remote,%v:443\n", us.Store))
	f.WriteString("packing: ee\n")
	return nil
}

func (us UpspinService) setKeys(path string) error {
	// check if keys are set already
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	args := []string{"keygen", fmt.Sprintf("-secretseed=%v", us.Seed)}
	// if the directory is populated, rotate the keys instead of generating new ones
	// this method of appending args is ugly, but they have to be in this specific order:
	// upspin keygen <-rotate> -secretseed=<seed> path
	if _, err = f.Readdir(1); err == nil {
		args = append(args, "-rotate")
	}

	return exec.Command("upspin", append(args, path)...).Run()
}

func (us *UpspinService) Update() {
	data := getFileData(filepath.Join(*upspinConfigDir, "config"))
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
	makeUserDirectories(*upspinConfigDir)
	if err := us.setFileData(filepath.Join(*upspinConfigDir, "config")); err != nil {
		return err
	}
	fullKeyPath := filepath.Join(*upspinKeyDir, us.User)
	makeUserDirectories(fullKeyPath)
	if err := us.setKeys(fullKeyPath); err != nil {
		return err
	}
	return nil
}

func NewUpspinService() (*UpspinService, error) {
	data := getFileData(filepath.Join(*upspinConfigDir, "config"))
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
