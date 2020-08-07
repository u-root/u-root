// +build !windows

// Copyright (c) 2018, Ian Haken. All rights reserved.
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

package main

import (
	"fmt"
	"os"
)

func showUsage() {
	fmt.Println("Usage: ./tpm-sign <sign|verify|generate|extendPcr> [...]")
}

func main() {
	if len(os.Args) < 2 {
		showUsage()
	} else {
		switch os.Args[1] {
		case "sign":
			signAction()
		case "verify":
			verifyAction()
		case "generate":
			generateAction()
		case "extendPcr":
			extendPcrAction()
		default:
			showUsage()
		}
	}
}
