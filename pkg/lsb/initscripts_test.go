// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lsb

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestInitScript(t *testing.T) {
	t.Parallel()

	t.Run("Marshal", func(t *testing.T) {
		t.Parallel()
		for name, tc := range map[string]struct {
			input    *InitScript
			expected string
			wantErr  string
		}{
			"success": {
				input: &InitScript{
					Provides:         "example",
					ShortDescription: "An example service",
					Description:      "A longer description of the example service.",
					DefaultStart:     []uint8{2, 3},
					DefaultStop:      []uint8{0, 6},
					RequiredStart:    []string{"network"},
					RequiredStop:     []string{"network"},
					XInteractive:     true,
				},
				expected: `### BEGIN INIT INFO
# Provides: example
# Short-Description: An example service
# Description: A longer description of the example service.
# Default-Start: 2 3
# Default-Stop: 0 6
# Required-Start: network
# Required-Stop: network
# X-Interactive: true
### END INIT INFO
`,
			},
			"empty struct": {
				input:    &InitScript{},
				expected: "### BEGIN INIT INFO\n### END INIT INFO\n",
			},
		} {
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				actual, err := tc.input.Marshal()
				if tc.wantErr != "" {
					if !strings.Contains(err.Error(), tc.wantErr) {
						t.Fatalf(
							"(%v).Marshal() = _, %v, want match for _, %q",
							tc.input,
							err, tc.wantErr,
						)
					}
				} else {
					if err != nil || !cmp.Equal(tc.expected, actual) {
						t.Fatalf(
							"(%v).Marshal() = %v, %v, want %v, nil",
							tc.input, actual, err,
							tc.expected,
						)
					}
				}
			})
		}
	})

	t.Run("Unmarshal", func(t *testing.T) {
		t.Parallel()
		for name, tc := range map[string]struct {
			input    io.Reader
			expected *InitScript
			wantErr  string
		}{
			"success": {
				input: strings.NewReader(`### BEGIN INIT INFO
# Provides: example
# Short-Description: An example service
# Description: A longer description of the example service.
# Default-Start: 2 3
# Default-Stop: 0 6
# Required-Start: network
# Required-Stop: network
# X-Interactive: true
### END INIT INFO`),
				expected: &InitScript{
					Provides:         "example",
					ShortDescription: "An example service",
					Description:      "A longer description of the example service.",
					DefaultStart:     []uint8{2, 3},
					DefaultStop:      []uint8{0, 6},
					RequiredStart:    []string{"network"},
					RequiredStop:     []string{"network"},
					XInteractive:     true,
				},
			},
			"missing start token": {
				input: strings.NewReader(`# Provides: example
# Short-Description: An example service
# Default-Start: 2 3
# Default-Stop: 0 6
### END INIT INFO`),
				expected: &InitScript{
					Provides:     "example",
					DefaultStart: []uint8{2, 3},
					DefaultStop:  []uint8{0, 6},
				},
				wantErr: "lsb block marker missing",
			},
			"missing stop token": {
				input: strings.NewReader(`### BEGIN INIT INFO
# Provides: example
# Default-Start: 2 3
# Default-Stop: 0 6`),
				expected: &InitScript{
					Provides:     "example",
					DefaultStart: []uint8{2, 3},
					DefaultStop:  []uint8{0, 6},
				},
				wantErr: "lsb block marker missing",
			},
			"nil reader": {
				wantErr: "data cannot be nil",
			},
		} {
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				actual := &InitScript{}
				err := actual.Unmarshal(tc.input)
				if tc.wantErr != "" {
					if !strings.Contains(err.Error(), tc.wantErr) {
						t.Fatalf(
							"actual.Unmarshal(%v) = _, %v, want match for _, %q",
							tc.input,
							err, tc.wantErr,
						)
					}
				} else {
					if err != nil {
						t.Fatalf(
							"actual.Unmarshal(%v) = %v, want nil",
							tc.input, err,
						)
					}
					if !reflect.DeepEqual(tc.expected, actual) {
						t.Fatalf(
							"actual = %v, want %v",
							actual, tc.expected,
						)
					}
				}
			})
		}
	})
}
