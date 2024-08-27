// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package netcat

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

type SSLOptions struct {
	// In connect mode, this option transparently negotiates an SSL session
	// In server mode, this option listens for incoming SSL connections
	// Depending on the Protocol Type, either TLS (TCP) od DTLS (UDP) will be used
	Enabled bool

	CertFilePath  string // Path to the certificate file in PEM format
	KeyFilePath   string // Path to the private key file in PEM format
	VerifyTrust   bool   // In client mode is like --ssl except that it also requires verification of the server certificate. No effect in server mode.
	TrustFilePath string // Verify trust and domain name of certificates

	// List of ciphersuites that Ncat will use when connecting to servers or when accepting SSL connections from clients
	// Syntax is described in the OpenSSL ciphers(1) man page
	Ciphers []string
	SNI     string   // (Server Name Indication) Tell the server the name of the logical server Ncat is contacting
	ALPN    []string // List of protocols to send via the Application-Layer Protocol Negotiation
}

func (s *SSLOptions) GenerateTLSConfiguration() (*tls.Config, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: !s.VerifyTrust,
	}

	if s.CertFilePath == "" || s.KeyFilePath == "" {
		return nil, fmt.Errorf("both  certificate and key file must be provided")
	}

	cer, err := tls.LoadX509KeyPair(s.CertFilePath, s.KeyFilePath)
	if err != nil {
		return nil, fmt.Errorf("connection: %v", err)
	}

	tlsConfig.Certificates = []tls.Certificate{cer}

	if s.VerifyTrust {
		caCert, err := os.ReadFile(s.TrustFilePath)
		if err != nil {
			return nil, fmt.Errorf("cannot read CA certificate: %v", err)
		}
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("cannot append CA certificate to pool")
		}

		tlsConfig.RootCAs = caCertPool
	}

	if s.SNI != "" {
		tlsConfig.ServerName = s.SNI
	}

	if s.ALPN != nil {
		tlsConfig.NextProtos = s.ALPN
	}

	return tlsConfig, nil
}
