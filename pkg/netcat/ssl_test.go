// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package netcat

import (
	"crypto/tls"
	"fmt"
	"os"
	"reflect"
	"testing"
)

const dummyCert = `-----BEGIN CERTIFICATE-----
MIIDazCCAlOgAwIBAgIUWlE1rjpClUu2L8Nn2aa/qV2VjG0wDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yNDA2MjYxMTM0NTZaFw00NDA2
MjExMTM0NTZaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQDR7VuU0pJ1cYsxbUJ/4yIB+lp9WKCsthU9SD3Kwh1b
vNkUn5zPkEHzpL/PDwwzGjx+rto6kze1kyGuqeSy+ID3g8e7wF77kzBGj3Mj0eLW
FE0Y1uLFXvpUQdmg0jsNn445oJWJekFYKgq1V89hGJ5B0159ZpKdRg8q9yA/bcnt
YzaGZMGxr2PfyWPVxQ/kWeyJUGqtdirMAjnqhBbms4bgrzYtyEsSgE3XtgQyMcP2
CvYcWktypyrfiSwQfdVWF0ZzcogRfNXCftO176c+iFdk5V5HmHntFLiD2rFpwL0R
26oKQGAjrIpQELuxkSZHMuNXtnjDLh0hMuuAFyjH2VG5AgMBAAGjUzBRMB0GA1Ud
DgQWBBQMZ6vj8LDrGmxSksgkF0TBWdKiDDAfBgNVHSMEGDAWgBQMZ6vj8LDrGmxS
ksgkF0TBWdKiDDAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQAA
8lQZafcfsdtbO8Az+L2GVZcj2+iHDnzosHCsbRYxl5zpnqhjpcEmNSuKGI2zuJt9
ehUhzbRXylngZGQqunoc4KFe/AtkWkZtQEQenKCraUoPkP5qFPBa737Nq2r2JicS
wR0PDq3Fb5p/V6kinU9YjYs7wrCKPnY75CG+d475qIz+xX6OK0afjuxaOXRc9pud
1JA6aTlXaYlvZ1PjGWCnPg9c2LkOlSSqZSk3ft7wn5TBpBGmhvGBt6WqZZs6DV85
hQ9rtapONW8qvnhv5EE21demLHh8XXY34Aeg/St5hwlEKHIb040QM+Wl7TCw2/Nb
8ITt8qjgwq/18yCDVuoy
-----END CERTIFICATE-----
`

const dummyKey = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDR7VuU0pJ1cYsx
bUJ/4yIB+lp9WKCsthU9SD3Kwh1bvNkUn5zPkEHzpL/PDwwzGjx+rto6kze1kyGu
qeSy+ID3g8e7wF77kzBGj3Mj0eLWFE0Y1uLFXvpUQdmg0jsNn445oJWJekFYKgq1
V89hGJ5B0159ZpKdRg8q9yA/bcntYzaGZMGxr2PfyWPVxQ/kWeyJUGqtdirMAjnq
hBbms4bgrzYtyEsSgE3XtgQyMcP2CvYcWktypyrfiSwQfdVWF0ZzcogRfNXCftO1
76c+iFdk5V5HmHntFLiD2rFpwL0R26oKQGAjrIpQELuxkSZHMuNXtnjDLh0hMuuA
FyjH2VG5AgMBAAECggEADDHMkx2UUmwxGMLvDPzFufWwEf32/3FoVHIA3OlfyTd0
KMWI12na2utkFQQbwlAw2W8Q0DxDDTIpz7qgxWC4JSirjpWDLvwC3uZwWtFTavos
7Fd3Pt3gjspweO4dbhIpseFJLn5Ck3uFubkLG+nRL6O2pnQx6h7qvKU0Y1reUwLI
VJqgfFUBCiXF26T8O+6d/FYiSMdUgOBOmWVTIOMaJOd3KvYEEbqP7Rfa3XQQhR4G
0EcRIC5LLAsigzLH/IsfzMn2jGC7HPnvazir26Dpqu+5teAF8NhOHTVxnO7pDggQ
6H8Vz9ajVFOXxgJsYgzC74TeHr6JEorWVb9yeWpojQKBgQD/2tvcX4sJZh7bpofh
3JMvq8nGGk9Zt5X1IQoHj4rqnWRcimU0mQm7wAjISCGkfazQ7/tJLBcUW1XSBK0W
6AyZSXxGQlXfREmcRBl3po7plvVk57dcxeUSDL6vqmjdKYAOs9Jev56hYiFgMh3u
f4yN416wU5/eodSeSqLmIG+rbwKBgQDSC9TxfrE0rHLJLpO7URU3uWhXURN/ofxG
SVBVQrNmqyTrXLNgMRDJHyfS24r7AuwJLdvJo5nhxiWYshdnzArJrnMbhixPNU/A
aUphEMLOceuggQx/vhD9T4g3HT+q2wkI11aEL8Zq6mBST+UakF1xEreQPtOV6QQA
G3H1iJ9hVwKBgQD1CniNzFfOLacaOZlkgSvaeU4rVGFxDLorZnRDn3+tigZn9whM
4tGGproCj8rgzpioF190yixki8Fa/q2EBcSjPtUuOTQjPDS/3B0EElpHcBQgiyh7
SvFEYz5x4eTDBI8oBaNSqXVVHTXX+sfd9vz3m67Bc6XmxNlsrRDtFF2/MwKBgH0L
3jIHIqghIhTzTa/ujZsnHh8dfWY2oWGWs+SOWQ9+Q/R6s69Ihp21lpfJa+wTyUGN
s5NPeoUW2bsWCykYKDP5Tz3LmwVsz5XVGRrAR7lvyL89FJvYI3UqrAVjvEuTKsXA
rRj0+EMeVUmrltFBsN9oLTAKtxxAJMmLjUSHmZrxAoGAc6IA4bzM646bHDga/FsE
9beh5xaJkN1MYRGaWuoYU9SWex7bYJ/MbsTGEXcChehdyJxZCgrKGrIr0JOoqIQZ
Pd4+B9XLuHcgyHIr4pdYUMcT/PuXywWSCSY2tLixEFspfqCKRXrqOIQ7M4q7JjF1
H4ng16X4yY1OGtXg+MbXeaM=
-----END PRIVATE KEY-----`

func TestGenerateTLSConfigurationExtended(t *testing.T) {
	tmpDir := t.TempDir()

	certFile, err := os.CreateTemp(tmpDir, "cert.pem")
	if err != nil {
		t.Fatalf("Failed to create temp cert file: %v", err)
	}
	defer os.Remove(certFile.Name())

	keyFile, err := os.CreateTemp(tmpDir, "key.pem")
	if err != nil {
		t.Fatalf("Failed to create temp key file: %v", err)
	}
	defer os.Remove(keyFile.Name())

	_, _ = certFile.WriteString(dummyCert)
	_, _ = keyFile.WriteString(dummyKey)

	tests := []struct {
		name      string
		opts      SSLOptions
		wantErr   bool
		checkFunc func(*tls.Config) error // Function to perform additional checks on tls.Config
	}{
		{
			name: "Valid certificate and key",
			opts: SSLOptions{
				Enabled:      true,
				CertFilePath: certFile.Name(),
				KeyFilePath:  keyFile.Name(),
			},

			wantErr: false,
		},
		{
			name: "Missing certificate",
			opts: SSLOptions{
				Enabled:     true,
				KeyFilePath: keyFile.Name(),
			},

			wantErr: true,
		},
		{
			name: "Missing key",
			opts: SSLOptions{
				Enabled:      true,
				CertFilePath: certFile.Name(),
			},

			wantErr: true,
		},
		{
			name: "Fail to parse cert / key",
			opts: SSLOptions{
				Enabled:      true,
				KeyFilePath:  certFile.Name(),
				CertFilePath: keyFile.Name(),
			},

			wantErr: true,
		},
		{
			name: "VerifyTrust with wrong CA certificate",
			opts: SSLOptions{
				Enabled:       true,
				CertFilePath:  certFile.Name(),
				KeyFilePath:   keyFile.Name(),
				VerifyTrust:   true,
				TrustFilePath: keyFile.Name(),
			},
			wantErr: true,
		},
		{
			name: "VerifyTrust with nonexistent CA certificate",
			opts: SSLOptions{
				Enabled:       true,
				CertFilePath:  certFile.Name(),
				KeyFilePath:   keyFile.Name(),
				VerifyTrust:   true,
				TrustFilePath: "nonexistent_ca.pem",
			},
			wantErr: true,
		},
		{
			name: "VerifyTrust with valid CA certificate",
			opts: SSLOptions{
				Enabled:       true,
				CertFilePath:  certFile.Name(),
				KeyFilePath:   keyFile.Name(),
				VerifyTrust:   true,
				TrustFilePath: certFile.Name(),
			},
			wantErr: false,
		},
		{
			name: "SNI is set",
			opts: SSLOptions{
				Enabled:      true,
				CertFilePath: certFile.Name(),
				KeyFilePath:  keyFile.Name(),
				SNI:          "example.com",
			},
			wantErr: false,
			checkFunc: func(cfg *tls.Config) error {
				if cfg.ServerName != "example.com" {
					return fmt.Errorf("expected ServerName to be 'example.com', got %s", cfg.ServerName)
				}
				return nil
			},
		},
		{
			name: "ALPN is set",
			opts: SSLOptions{
				Enabled:      true,
				CertFilePath: certFile.Name(),
				KeyFilePath:  keyFile.Name(),
				ALPN:         []string{"http/1.1", "h2"},
			},
			wantErr: false,
			checkFunc: func(cfg *tls.Config) error {
				if !reflect.DeepEqual(cfg.NextProtos, []string{"http/1.1", "h2"}) {
					return fmt.Errorf("expected NextProtos to be ['http/1.1', 'h2'], got %v", cfg.NextProtos)
				}
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.opts.GenerateTLSConfiguration()
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateTLSConfiguration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.checkFunc != nil {
				if err := tt.checkFunc(got); err != nil {
					t.Errorf("checkFunc failed: %v", err)
				}
			}
		})
	}
}
