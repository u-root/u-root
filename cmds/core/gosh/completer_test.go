// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"github.com/knz/bubbline/editline"
)

func TestAutocomplete(t *testing.T) {
	for _, tt := range []struct {
		name        string
		input       string
		completions []string
	}{
		{
			name:        "echo",
			input:       "ech",
			completions: []string{"echo"},
		},
		{
			name:        "cwd",
			input:       "./",
			completions: []string{"completer_test.go"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			val := make([][]rune, 1)
			val[0] = append(val[0], []rune(tt.input)...)

			_, completions := autocomplete(val, 0, len(tt.input))

			if !completionsEqual(0, len(tt.completions), tt.completions, completions) {
				t.Errorf("want: %v, got: %v", tt.completions, completions)
			}
		})
	}
}

func completionsEqual(numCat, numEnt int, want []string, got editline.Completions) bool {
	if got.NumCategories() < numCat {
		return false
	}

	for i := 0; i < numCat; i++ {
		if got.NumEntries(i) < numEnt {
			return false
		}
	}

	for j := 0; j < numCat; j++ {
		for i := 0; i < numEnt; i++ {
			found := false
			for _, entry := range want {
				if entry == got.Entry(j, i).Title() {
					found = true
				}
			}
			if !found {
				return false
			}
		}
	}

	return true
}
