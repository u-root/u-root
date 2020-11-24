// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package policy locates and parses a JSON policy file.
package policy

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"

	slaunch "github.com/u-root/u-root/pkg/securelaunch"
	"github.com/u-root/u-root/pkg/securelaunch/config"
	"github.com/u-root/u-root/pkg/securelaunch/eventlog"
	"github.com/u-root/u-root/pkg/securelaunch/launcher"
	"github.com/u-root/u-root/pkg/securelaunch/measurement"
	"github.com/u-root/u-root/pkg/securelaunch/tpm"
)

// Policy describes the policy used to drive the security engine.
//
// The policy is stored as a JSON file.
type Policy struct {
	Config     config.Config
	Collectors []measurement.Collector
	Launcher   launcher.Launcher
	EventLog   eventlog.EventLog
}

// policyBytes is a byte slice to hold a copy of the policy file in memory.
var policyBytes []byte

// pubkeyBytes is a byte slice to hold a copy of the public key file in memory.
var pubkeyBytes []byte

// signatureBytes is a byte slice to hold a copy of the signature file in memory.
var signatureBytes []byte

// Load reads the specified policy file. If pubkey and signature files are
// also set, then load them as well. Note that the pubkey and signature files
// need to be set together (i.e., either they are both set or neither is set).
func Load(policyLocation, pubkeyLocation, signatureLocation string) error {
	var err error

	if policyLocation == "" {
		return fmt.Errorf("policy file must be set")
	}

	if pubkeyLocation != "" && signatureLocation == "" {
		return fmt.Errorf("signature file must be provided with public key file")
	}

	if signatureLocation != "" && pubkeyLocation == "" {
		return fmt.Errorf("public key file must be provided with signature file")
	}

	policyBytes, err = slaunch.ReadFile(policyLocation)
	if err != nil {
		return err
	}

	if pubkeyLocation != "" && signatureLocation != "" {
		pubkeyBytes, err = slaunch.ReadFile(pubkeyLocation)
		if err != nil {
			return err
		}

		signatureBytes, err = slaunch.ReadFile(signatureLocation)
		if err != nil {
			return err
		}
	}

	return nil
}

// VerifyPubkey verifies the public key file against the provided hash.
func VerifyPubkey(hashBytes []byte) error {
	if len(pubkeyBytes) == 0 {
		return fmt.Errorf("public key file not yet loaded or empty")
	}

	if len(hashBytes) == 0 {
		return fmt.Errorf("hash not yet loaded or empty")
	}

	pubkeyHashBytes := tpm.HashReader(bytes.NewReader(pubkeyBytes))
	if !bytes.Equal(pubkeyHashBytes, hashBytes) {
		return fmt.Errorf("public key hash does not match provided")
	}

	return nil
}

// Verify verifies the policy file using the public key and signature.
func Verify() error {
	if len(policyBytes) == 0 {
		return fmt.Errorf("policy file not yet loaded or empty")
	}

	if len(pubkeyBytes) == 0 {
		return fmt.Errorf("public key file not yet loaded or empty")
	}

	if len(signatureBytes) == 0 {
		return fmt.Errorf("signature file not yet loaded or empty")
	}

	// Decode and parse public key.
	pubkeyBlock, _ := pem.Decode(pubkeyBytes)
	if pubkeyBlock == nil || pubkeyBlock.Type != "PUBLIC KEY" {
		return fmt.Errorf("public key is of the wrong type; Pem Type: %v", pubkeyBlock.Type)
	}

	pubkeyParsed, err := x509.ParsePKIXPublicKey(pubkeyBlock.Bytes)
	if err != nil {
		return fmt.Errorf("unable to parse public key: %w", err)
	}

	var pubkey *rsa.PublicKey
	pubkey, ok := pubkeyParsed.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("unable to get public key: %w", err)
	}

	// Hash the policy file.
	hashed := sha256.Sum256(policyBytes)

	// Verify the policy file.
	err = rsa.VerifyPKCS1v15(pubkey, crypto.SHA256, hashed[:], signatureBytes)
	if err != nil {
		return fmt.Errorf("could not verify policy file: %w", err)
	}

	log.Printf("Policy Verified OK")

	return nil
}

// Parse accepts a JSON file as input, parses it into a well defined Policy
// structure and returns a pointer to the Policy structure.
func Parse() (*Policy, error) {
	if len(policyBytes) == 0 {
		return nil, fmt.Errorf("policy file not yet loaded or empty")
	}

	policy := &Policy{}
	var parse struct {
		Config     json.RawMessage   `json:"config"`
		Collectors []json.RawMessage `json:"collectors"`
		Attestor   json.RawMessage   `json:"attestor"`
		Launcher   json.RawMessage   `json:"launcher"`
		EventLog   json.RawMessage   `json:"eventlog"`
	}

	if err := json.Unmarshal(policyBytes, &parse); err != nil {
		log.Printf("policy: Error unmarshalling policy file: %v", err)
		return nil, err
	}

	config.Conf = config.New()
	if len(parse.Config) > 0 {
		if err := json.Unmarshal(parse.Config, &config.Conf); err != nil {
			log.Printf("policy: Error unmarshalling `config` section of policy file: %v", err)
			return nil, err
		}

		log.Printf("policy: Setting measurement PCR to %d", config.Conf.MeasurementPCR)
		measurement.SetPCR(uint32(config.Conf.MeasurementPCR))
	}

	if len(parse.Collectors) > 0 {
		for _, collectors := range parse.Collectors {
			collector, err := measurement.GetCollector(collectors)
			if err != nil {
				log.Printf("policy: Error GetCollector err:c=%s, collector=%v", collectors, collector)
				return nil, err
			}

			policy.Collectors = append(policy.Collectors, collector)
		}
	} else {
		log.Printf("policy: No collectors found, disabling")
		config.Conf.Collectors = false
	}

	if len(parse.Launcher) > 0 {
		if err := json.Unmarshal(parse.Launcher, &policy.Launcher); err != nil {
			log.Printf("policy: Error unmarshalling `launcher` section of policy file: %v", err)
			return nil, err
		}
	}

	if len(parse.EventLog) > 0 {
		if err := json.Unmarshal(parse.EventLog, &policy.EventLog); err != nil {
			log.Printf("policy: Error unmarshalling `eventlog` section of policy file: %v", err)
			return nil, err
		}
	} else {
		log.Printf("policy: No eventlog found, disabling")
		config.Conf.EventLog = false
	}

	return policy, nil
}

// Measure measures the policy file.
func Measure() error {
	if len(policyBytes) == 0 {
		return fmt.Errorf("policy file not yet loaded or empty")
	}

	eventDesc := "File Collector: measured securelaunch policy file"
	if err := measurement.HashBytes(policyBytes, eventDesc); err != nil {
		log.Printf("policy: ERR: could not measure policy file: %v", err)
		return err
	}

	return nil
}
