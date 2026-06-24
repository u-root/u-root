// SPDX-License-Identifier: MIT
// Copyright 2026 Google LLC

//go:build linux

package ssh9p

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/hugelgupf/p9/fsimpl/localfs"
	"github.com/hugelgupf/p9/p9"
	"github.com/u-root/u-root/pkg/sshstream"
	"golang.org/x/crypto/ssh"
	"golang.org/x/sync/errgroup"
)

func generateTestKeys(t *testing.T, tmpDir string) (string, string) {
	t.Helper()
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	privPath := filepath.Join(tmpDir, "id_rsa")
	privFile, err := os.Create(privPath)
	if err != nil {
		t.Fatalf("Failed to create private key file: %v", err)
	}
	defer privFile.Close()

	privPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	}
	if err := pem.Encode(privFile, privPEM); err != nil {
		t.Fatalf("Failed to write private key: %v", err)
	}

	pub, err := ssh.NewPublicKey(&priv.PublicKey)
	if err != nil {
		t.Fatalf("Failed to create SSH public key: %v", err)
	}

	pubPath := filepath.Join(tmpDir, "id_rsa.pub")
	if err := os.WriteFile(pubPath, ssh.MarshalAuthorizedKey(pub), 0644); err != nil {
		t.Fatalf("Failed to write public key file: %v", err)
	}

	return privPath, pubPath
}

func TestFDFor9PTransportMount(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	// Test TCP connection.
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Listen failed: %v", err)
	}
	defer ln.Close()

	done := make(chan bool)
	go func() {
		c, _ := ln.Accept()
		if c != nil {
			c.Close()
		}
		done <- true
	}()

	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}
	defer conn.Close()

	eg := &errgroup.Group{}
	file, err := fdFor9PTransportMount(ctx, eg, conn)
	if err != nil {
		t.Fatalf("fdFor9PTransportMount failed for TCP: %v", err)
	}
	if file == nil {
		t.Errorf("file should not be nil for TCP")
	} else {
		file.Close()
	}
	<-done

	// Test non-TCP connection (using net.Pipe).
	c1, c2 := net.Pipe()
	defer c1.Close()
	defer c2.Close()

	eg2 := &errgroup.Group{}
	file2, err := fdFor9PTransportMount(ctx, eg2, c1)
	if err != nil {
		t.Fatalf("fdFor9PTransportMount failed for Pipe: %v", err)
	}
	if file2 == nil {
		t.Errorf("file should not be nil for Pipe")
	} else {
		file2.Close()
	}

	// Test context cancellation in fdFor9PTransportMount.
	ctx3, cancel3 := context.WithCancel(t.Context())
	c13, c23 := net.Pipe()
	defer c23.Close()

	eg3 := &errgroup.Group{}
	_, err = fdFor9PTransportMount(ctx3, eg3, c13)
	if err != nil {
		t.Fatalf("fdFor9PTransportMount failed for cancellation test: %v", err)
	}

	cancel3()
	// No easy way to verify userside is closed without more inspection,
	// but we've at least exercised the path.
}

func TestServerServe(t *testing.T) {
	tmpDir := t.TempDir()
	attacher := localfs.Attacher(tmpDir)
	p9s := p9.NewServer(attacher)
	s := NewServer(p9s)

	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Listen failed: %v", err)
	}
	defer ln.Close()

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	serveErr := make(chan error, 1)
	go func() {
		serveErr <- s.Serve(ctx, ln)
	}()

	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}
	defer conn.Close()

	client, err := p9.NewClient(conn)
	if err != nil {
		t.Fatalf("p9.NewClient failed: %v", err)
	}
	defer client.Close()

	// Try to p9.Client.Attach.
	_, err = client.Attach("")
	if err != nil {
		t.Errorf("Attach failed: %v", err)
	}

	cancel()
	err = <-serveErr
	if err != nil && !errors.Is(err, net.ErrClosed) {
		t.Errorf("Serve returned error: %v", err)
	}
}

func TestServerServeSSH(t *testing.T) {
	tmpDir := t.TempDir()
	privPath, pubPath := generateTestKeys(t, tmpDir)

	serverCfg, err := sshstream.NewServerConfig(pubPath, privPath)
	if err != nil {
		t.Fatalf("NewServerConfig failed: %v", err)
	}

	clientCfg, err := sshstream.NewClientConfig(privPath, ssh.InsecureIgnoreHostKey())
	if err != nil {
		t.Fatalf("NewClientConfig failed: %v", err)
	}

	attacher := localfs.Attacher(tmpDir)
	p9s := p9.NewServer(attacher)
	s := NewServer(p9s, WithSSHServer(serverCfg))

	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Listen failed: %v", err)
	}
	defer ln.Close()

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	serveErr := make(chan error, 1)
	go func() {
		serveErr <- s.Serve(ctx, ln)
	}()

	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}
	defer conn.Close()

	// Wrap client side in SSH.
	nc, err := sshstream.NewClient(conn, clientCfg)
	if err != nil {
		t.Fatalf("sshstream.NewClient failed: %v", err)
	}
	defer nc.Close()

	client, err := p9.NewClient(nc)
	if err != nil {
		t.Fatalf("p9.NewClient failed: %v", err)
	}
	defer client.Close()

	// Try to p9.Client.Attach.
	_, err = client.Attach("")
	if err != nil {
		t.Errorf("Attach failed: %v", err)
	}

	cancel()
	err = <-serveErr
	if err != nil && !errors.Is(err, net.ErrClosed) {
		t.Errorf("Serve returned error: %v", err)
	}
}
