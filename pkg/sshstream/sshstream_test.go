// SPDX-License-Identifier: MIT
// Copyright 2026 Google LLC
package sshstream

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"net"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/crypto/ssh"
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

func TestSSHStream(t *testing.T) {
	tmpDir := t.TempDir()
	privPath, pubPath := generateTestKeys(t, tmpDir)

	serverCfg, err := NewServerConfig(pubPath, privPath)
	if err != nil {
		t.Fatalf("NewServerConfig failed: %v", err)
	}

	clientCfg, err := NewClientConfig(privPath, ssh.InsecureIgnoreHostKey())
	if err != nil {
		t.Fatalf("NewClientConfig failed: %v", err)
	}

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Listen failed: %v", err)
	}
	defer listener.Close()

	errChan := make(chan error, 2)
	dataChan := make(chan []byte, 1)

	// Start server
	go func() {
		serverConn, err := listener.Accept()
		if err != nil {
			errChan <- err
			return
		}
		defer serverConn.Close()

		stream, err := NewServer(serverConn, serverCfg)
		if err != nil {
			errChan <- err
			return
		}
		defer stream.Close()

		buf := make([]byte, 1024)
		n, err := stream.Read(buf)
		if err != nil && err != io.EOF {
			errChan <- err
			return
		}
		dataChan <- buf[:n]
	}()

	// Start client
	go func() {
		clientConn, err := net.Dial("tcp", listener.Addr().String())
		if err != nil {
			errChan <- err
			return
		}
		defer clientConn.Close()

		clientStream, err := NewClient(clientConn, clientCfg)
		if err != nil {
			errChan <- err
			return
		}
		defer clientStream.Close()

		_, err = clientStream.Write([]byte("hello world"))
		if err != nil {
			errChan <- err
			return
		}
	}()

	select {
	case err := <-errChan:
		t.Fatalf("Error in goroutine: %v", err)
	case data := <-dataChan:
		// Assumes the server's read grabbed the entire string.
		if !bytes.Equal(data, []byte("hello world")) {
			t.Errorf("Got %q, want %q", string(data), "hello world")
		}
	}
}

func TestSSHStreamRejection(t *testing.T) {
	tmpDir := t.TempDir()
	privPath, pubPath := generateTestKeys(t, tmpDir)

	serverCfg, _ := NewServerConfig(pubPath, privPath)

	otherDir := t.TempDir()
	otherPriv, _ := generateTestKeys(t, otherDir)
	clientCfg, _ := NewClientConfig(otherPriv, ssh.InsecureIgnoreHostKey())

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Listen failed: %v", err)
	}
	defer listener.Close()

	go func() {
		serverConn, err := listener.Accept()
		if err != nil {
			return
		}
		defer serverConn.Close()
		_, _ = NewServer(serverConn, serverCfg)
	}()

	clientConn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}
	defer clientConn.Close()

	_, err = NewClient(clientConn, clientCfg)
	if err == nil {
		t.Fatal("Client should have failed to connect with unauthorized key")
	}
}
