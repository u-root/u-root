// Copyright (c) 2018, Google LLC All rights reserved.
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

package tpm2

import (
	"flag"
	"io"
	"testing"
)

var runTPMTests = flag.Bool("use-tbs", false, "Run integration tests against Windows TPM Base Services (TBS). Defaults to false.")

func useDeviceTPM() bool { return *runTPMTests }

func openDeviceTPM(tb testing.TB) io.ReadWriteCloser {
	tb.Helper()
	rw, err := OpenTPM()
	if err != nil {
		tb.Fatalf("Open TPM failed: %s\n", err)
	}
	return rw
}
