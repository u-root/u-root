// +build !windows

// Copyright (c) 2014, Google LLC All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tpm

import (
	"io"
	"os"
	"testing"
)

// Skip the test if we can't open the TPM.
func openTPMOrSkip(t *testing.T) io.ReadWriteCloser {
	tpmPath := os.Getenv(tpmPathEnvVar)
	if tpmPath == "" {
		tpmPath = "/dev/tpm0"
	}

	rwc, err := openAndStartupTPM(tpmPath, true)
	if err != nil {
		t.Skipf("Skipping test, since we can't open %s for read/write: %s\n", tpmPath, err)
	}

	return rwc
}
