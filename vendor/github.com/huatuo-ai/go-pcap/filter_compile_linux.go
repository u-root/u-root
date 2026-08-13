// Copyright 2026 The HuaTuo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build linux

package pcap

import (
	"golang.org/x/net/bpf"

	"github.com/huatuo-ai/go-pcap/filter"
)

func compileLiveFilter(expr string) ([]bpf.Instruction, error) {
	return filter.CompileWithOptions(expr, filter.CompileOptions{
		LinkType: filter.LinkTypeEthernet,
		Target:   filter.CompileTargetLinuxSocket,
	})
}
