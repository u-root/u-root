// +build !windows

package main

import (
	"crypto"
)

const (
	srkAuthEnvVar       = "TPM_SRK_AUTH"
	usageAuthEnvVar     = "TPM_USAGE_AUTH"
	migrationAuthEnvVar = "TPM_MIGRATION_AUTH"
)

var hashNames = map[string]crypto.Hash{
	"MD5":       crypto.MD5,
	"SHA1":      crypto.SHA1,
	"SHA224":    crypto.SHA224,
	"SHA256":    crypto.SHA256,
	"SHA384":    crypto.SHA384,
	"SHA512":    crypto.SHA512,
	"MD5SHA1":   crypto.MD5SHA1,
	"RIPEMD160": crypto.RIPEMD160,
}
