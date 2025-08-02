// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

// SSH client.
//
// Synopsis:
//
//	ssh OPTIONS [DEST]
//
// Description:
//
//	Connects to the specified destination.
//
// Options:
//
// Destination format:
//
//	[user@]hostname or ssh://[user@]hostname[:port]
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	guser "os/user"
	"path/filepath"
	"strings"

	sshconfig "github.com/kevinburke/ssh_config"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"golang.org/x/term"
)

var (
	flags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	debug      = flags.Bool("d", false, "enable debug prints")
	keyFile    = flags.String("i", "", "key file")
	configFile = flags.String("F", defaultConfigFile, "config file")

	v = func(string, ...any) {}

	// ssh config file
	cfg *sshconfig.Config

	errInvalidArgs = errors.New("invalid command-line arguments")
)

// loadConfig loads the SSH config file
func loadConfig(path string) (err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		if os.IsNotExist(err) {
			err = nil
			cfg = &sshconfig.Config{}
		}
		return
	}
	cfg, err = sshconfig.Decode(f)
	return
}

func main() {
	if err := run(os.Args, os.Stdin, os.Stdout, os.Stderr); err != nil {
		log.Fatalf("%v", err)
	}
}

func knownHosts() (ssh.HostKeyCallback, error) {
	etc, err := filepath.Glob("/etc/*/ssh_known_hosts")
	if err != nil {
		return nil, err
	}
	if home, ok := os.LookupEnv("HOME"); ok {
		etc = append(etc, filepath.Join(home, ".ssh", "known_hosts"))
	}
	return knownhosts.New(etc...)
}

// we demand that stdin be a proper os.File because we need to be able to put it in raw mode
func run(osArgs []string, stdin *os.File, stdout io.Writer, stderr io.Writer) error {
	flags.SetOutput(stderr)
	flags.Parse(osArgs[1:])
	if *debug {
		v = log.Printf
	}
	defer cleanup(stdin)

	// Check if they're given appropriate arguments
	args := flags.Args()
	var dest string
	if len(args) >= 1 {
		dest = args[0]
		args = args[1:]
	} else {
		fmt.Fprintf(stderr, "usage: %v [flags] [user@]dest[:port] [command]\n", osArgs[0])
		flags.PrintDefaults()
		return errInvalidArgs
	}

	// Read the config file (if any)
	if err := loadConfig(*configFile); err != nil {
		return fmt.Errorf("config parse failed: %w", err)
	}

	// Parse out the destination
	user, host, port, err := parseDest(dest)
	if err != nil {
		return fmt.Errorf("destination parse failed: %w", err)
	}

	cb, err := knownHosts()
	if err != nil {
		return fmt.Errorf("known hosts:%w", err)
	}
	// Build a client config with appropriate auth methods
	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: cb,
	}
	// Figure out if there's a keyfile or not
	kf := getKeyFile(host, *keyFile)
	key, err := os.ReadFile(kf)
	if err == nil {
		// The key exists
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return fmt.Errorf("ParsePrivateKey %v: %w", kf, err)
		}
		config.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else if err != nil && *keyFile != "" {
		return fmt.Errorf("could not read user-specified keyfile %v: %w", kf, err)
	}
	v("Config: %+v\n", config)
	if term.IsTerminal(int(stdin.Fd())) {
		pwReader := func() (string, error) {
			return readPassword(stdin, stdout)
		}
		config.Auth = append(config.Auth, ssh.PasswordCallback(pwReader))
	}

	// Now connect to the server
	conn, err := ssh.Dial("tcp", net.JoinHostPort(host, port), config)
	if err != nil {
		return fmt.Errorf("unable to connect: %w", err)
	}
	defer conn.Close()
	// Create a session on that connection
	session, err := conn.NewSession()
	if err != nil {
		return fmt.Errorf("unable to create session: %w", err)
	}
	session.Stdin = stdin
	session.Stdout = stdout
	session.Stderr = stderr
	defer session.Close()

	if len(args) > 0 {
		// run the command
		if err := session.Run(strings.Join(args, " ")); err != nil {
			return fmt.Errorf("failed to run command: %w", err)
		}
	} else {
		// Set up the terminal
		modes := ssh.TerminalModes{
			ssh.ECHO:          1,     // disable echoing
			ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
			ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
		}
		if term.IsTerminal(int(stdin.Fd())) {
			if err := raw(stdin); err != nil {
				// throw a notice but continue
				log.Printf("failed to set raw mode: %v", err)
			}
			// Try to figure out the terminal size
			width, height, err := getSize(stdin)
			if err != nil {
				return fmt.Errorf("failed to get terminal size: %w", err)
			}
			// Request pseudo terminal - "xterm" for now, may make this adjustable later.
			if err := session.RequestPty("xterm", height, width, modes); err != nil {
				log.Print("request for pseudo terminal failed: ", err)
			}
		}
		// Start shell on remote system
		if err := session.Shell(); err != nil {
			log.Fatal("failed to start shell: ", err)
		}
		// Wait for the session to complete
		session.Wait()
	}
	return nil
}

// parseDest splits an ssh destination spec into separate user, host, and port fields.
// Example specs:
//
//	ssh://user@host:port
//	user@host:port
//	user@host
//	host
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
		err = errors.New("no host specified")
	}
	return
}

// getKeyFile picks a keyfile if none has been set.
// It will use sshconfig, else use a default.
// The kf parameter is a user-specified key file. We pass it
// here so it can be re-written if it contains a ~
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
