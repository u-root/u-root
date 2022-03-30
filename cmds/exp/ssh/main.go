// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// SSH client.
//
// Synopsis:
//     ssh OPTIONS [DEST]
//
// Description:
//     Connects to the specified destination.
//
// Options:
//
// Destination format:
//     [user@]hostname or ssh://[user@]hostname[:port]
package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"os"
	guser "os/user"
	"path/filepath"
	"strings"

	config "github.com/kevinburke/ssh_config"
	sshconfig "github.com/kevinburke/ssh_config"
	"golang.org/x/crypto/ssh"
)

var (
	debug      = flag.Bool("d", false, "enable debug prints")
	keyFile    = flag.String("i", "", "key file")
	configFile = flag.String("F", defaultConfigFile, "config file")

	v = func(string, ...interface{}) {}

	// ssh config file
	cfg *sshconfig.Config
)

// loadConfig loads the SSH config file
func loadConfig(path string) (err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return
	}
	cfg, err = config.Decode(f)
	return
}

func main() {
	flag.Parse()
	if *debug {
		v = log.Printf
	}
	defer cleanup()

	loadConfig(*configFile)

	// Parse out the destination
	user, host, port, err := parseDest(flag.Arg(0))
	if err != nil {
		log.Fatalf("destination parse failed: %v", err)
	}

	// Connect to ssh server
	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// Figure out if there's a keyfile or not
	kf := getKeyFile(host, *keyFile)
	key, err := ioutil.ReadFile(kf)
	if err == nil {
		// The key exists
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			log.Fatalf("ParsePrivateKey %v: %v", kf, err)
		}
		config.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else if err != nil && *keyFile != "" {
		log.Fatalf("Could not read user-specified keyfile %v: %v", kf, err)
	}
	config.Auth = append(config.Auth, ssh.PasswordCallback(readPassword))

	conn, err := ssh.Dial("tcp", net.JoinHostPort(host, port), config)
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer conn.Close()
	// Create a session on that connection
	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("unable to create session: %v", err)
	}
	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	defer session.Close()

	// Set up the terminal
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	if err := raw(); err != nil {
		log.Fatalf("failed to set raw mode: %v", err)
	}
	// Try to figure out the terminal size
	width, height, err := getSize()
	if err != nil {
		log.Fatalf("failed to get terminal size: %v", err)
	}
	// Request pseudo terminal - "xterm" for now, may make this adjustable later.
	if err := session.RequestPty("xterm", height, width, modes); err != nil {
		log.Fatal("request for pseudo terminal failed: ", err)
	}

	// Start shell on remote system
	if err := session.Shell(); err != nil {
		log.Fatal("failed to start shell: ", err)
	}

	// Wait for the session to complete
	session.Wait()
}

// parseDest splits an ssh destination spec into separate user, host, and port fields.
// Example specs:
//		ssh://user@host:port
//		user@host:port
//		user@host
//		host
func parseDest(input string) (user, host, port string, err error) {
	// strip off any preceding ssh://
	input = strings.TrimPrefix(input, "ssh://")
	// try to find user
	i := strings.LastIndex(input, "@")
	if i < 0 {
		var u *guser.User
		u, err = guser.Current()
		if err != nil {
			return
		}
		user = u.Username
	} else {
		user = input[0:i]
		input = input[i+1:]
	}
	if host, port, err = net.SplitHostPort(input); err != nil {
		err = nil
		host = input
		port = "22"
	}
	if host == "" {
		err = errors.New("No host specified")
	}
	return
}

// getKeyFile picks a keyfile if none has been set.
// It will use sshconfig, else use a default.
func getKeyFile(host, kf string) string {
	v("getKeyFile for %q", kf)
	if len(kf) == 0 {
		var err error
		kf, err = cfg.Get(host, "IdentityFile")
		v("key file from config is %q", kf)
		if len(kf) == 0 || err != nil {
			kf = defaultKeyFile
		}
	}
	// The kf will always be non-zero at this point.
	if strings.HasPrefix(kf, "~") {
		kf = filepath.Join(os.Getenv("HOME"), kf[1:])
	}
	v("getKeyFile returns %q", kf)
	// this is a tad annoying, but the config package doesn't handle ~.
	return kf
}
