// Copyright 2024-2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package netcat

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

func (s *SSLOptions) GenerateTLSConfiguration() (*tls.Config, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: !s.VerifyTrust,
	}

	if s.CertFilePath == "" || s.KeyFilePath == "" {
		return nil, fmt.Errorf("both  certificate and key file must be provided")
	}

	cer, err := tls.LoadX509KeyPair(s.CertFilePath, s.KeyFilePath)
	if err != nil {
		return nil, fmt.Errorf("connection: %w", err)
	}

	tlsConfig.Certificates = []tls.Certificate{cer}

	if s.VerifyTrust {
		caCert, err := os.ReadFile(s.TrustFilePath)
		if err != nil {
			return nil, fmt.Errorf("cannot read CA certificate: %w", err)
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
