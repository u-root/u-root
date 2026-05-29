// SPDX-License-Identifier: MIT
// Copyright 2026 Google LLC
//
// Package sshstream provides a byte stream wrapped with SSH over a net.Conn.
// Think of it as TLS, but with SSH as the mechanism for authentication and
// encryption.
//
// Both sides of a net.Conn must participate.  One side acts as the ssh server
// and the other is the client.  Either side of a net.Conn can perform either
// role, but it makes the most sense for whoever Dials to be the client.
package sshstream

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

var (
	sshStreamName = "stream@u-root.org"
)

// Stream implements net.Conn.  It is an ssh.Channel that sits on top of the
// net.Conn, providing io.ReadWriter.
type Stream struct {
	channel ssh.Channel
	sshConn ssh.Conn
	net.Conn
}

// Read implements io.Read().
func (s *Stream) Read(data []byte) (int, error) {
	return s.channel.Read(data)
}

// Write implements io.Write().
func (s *Stream) Write(data []byte) (int, error) {
	return s.channel.Write(data)
}

// Close closes the stream.  Any blocked Read or Write operations will be
// unblocked and return errors.
func (s *Stream) Close() error {
	return errors.Join(s.channel.Close(), s.sshConn.Close(), s.Conn.Close())
}

// NewClient acts as an ssh client and creates an ssh channel on the underlying
// net.Conn c.  The returned net.Conn is a byte stream corresponding to the
// NewServer() on the remote side of c.
func NewClient(c net.Conn, cfg *ssh.ClientConfig) (net.Conn, error) {
	// c, chans, and reqs are all handled by ssh.NewClient(); don't worry
	// about them.
	cIGN, chansIGN, reqsIGN, err := ssh.NewClientConn(c, c.RemoteAddr().String(), cfg)
	if err != nil {
		return nil, fmt.Errorf("new clientconn: %w", err)
	}
	client := ssh.NewClient(cIGN, chansIGN, reqsIGN)

	channel, reqs, err := client.OpenChannel(sshStreamName, []byte{})
	if err != nil {
		return nil, fmt.Errorf("open channel: %w", err)
	}

	go ssh.DiscardRequests(reqs)

	return &Stream{channel: channel, sshConn: client, Conn: c}, nil
}

// Dial is a helper to dial a server and create a NewClient.
func Dial(network, address string, cfg *ssh.ClientConfig) (net.Conn, error) {
	c, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return NewClient(c, cfg)
}

// NewServer acts as an ssh server on the underlying net.Conn c.  The returned
// net.Conn is a byte stream corresponding to the ssh.Channel created by the
// client on the remote side of c.
func NewServer(c net.Conn, cfg *ssh.ServerConfig) (net.Conn, error) {
	server, chans, reqs, err := ssh.NewServerConn(c, cfg)
	if err != nil {
		return nil, fmt.Errorf("new serverconn: %w", err)
	}

	go ssh.DiscardRequests(reqs)

	sockChan := make(chan ssh.Channel)

	go func() {
		defer close(sockChan)
		madeChannel := false
		// Yes, 'chans' is a go chan of ssh.Channels, enjoy!
		for nc := range chans {
			if nc.ChannelType() != sshStreamName {
				nc.Reject(ssh.UnknownChannelType, "must be "+sshStreamName)
				continue
			}
			if madeChannel {
				nc.Reject(ssh.ResourceShortage, "already have a "+sshStreamName)
				continue
			}
			channel, reqs, err := nc.Accept()
			if err != nil {
				log.Printf("sshwrap: channel accept: %v", err)
				continue
			}
			madeChannel = true
			go ssh.DiscardRequests(reqs)
			sockChan <- channel
		}
	}()

	channel := <-sockChan
	if channel == nil {
		return nil, errors.New("sockChan closed early")
	}
	return &Stream{channel: channel, sshConn: server, Conn: c}, nil
}

// NewClientConfig is a helper to make an ssh.ClientConfig.
func NewClientConfig(privateKey string, hkcb ssh.HostKeyCallback) (*ssh.ClientConfig, error) {
	keyBytes, err := os.ReadFile(privateKey)
	if err != nil {
		return nil, fmt.Errorf("key read: %w", err)
	}
	keyParse, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("key parse: %w", err)
	}

	cfg := &ssh.ClientConfig{
		HostKeyCallback: hkcb,
	}
	cfg.Auth = []ssh.AuthMethod{ssh.PublicKeys(keyParse)}
	return cfg, nil
}

// NewServerConfig is a helper to make an ssh.ServerConfig.
func NewServerConfig(authorizedKeys, hostKey string) (*ssh.ServerConfig, error) {
	authBytes, err := os.ReadFile(authorizedKeys)
	if err != nil {
		return nil, fmt.Errorf("authkey read: %w", err)
	}
	authKeys := map[string]bool{}
	for len(authBytes) > 0 {
		k, _, _, rest, err := ssh.ParseAuthorizedKey(authBytes)
		if err != nil {
			return nil, fmt.Errorf("authkey parse: %w", err)
		}

		authKeys[string(k.Marshal())] = true
		authBytes = rest
	}
	cfg := &ssh.ServerConfig{
		PublicKeyCallback: func(cm ssh.ConnMetadata, k ssh.PublicKey) (*ssh.Permissions, error) {
			if authKeys[string(k.Marshal())] {
				return &ssh.Permissions{}, nil
			}
			return nil, fmt.Errorf("unknown key for user %q", cm.User())
		},
	}
	keyBytes, err := os.ReadFile(hostKey)
	if err != nil {
		return nil, fmt.Errorf("key read: %w", err)
	}
	keyParse, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("key parse: %w", err)
	}
	cfg.AddHostKey(keyParse)
	return cfg, nil
}
