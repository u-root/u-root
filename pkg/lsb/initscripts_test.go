// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lsb

import (
	"reflect"
	"testing"
)

func TestMarshalUnmarshal(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		input    *InitScript
		expected string
		wantErr  bool
	}{
		"simple script": {
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
			wantErr: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := tc.input.Marshal()
			if tc.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Error("expected error, got nil")
				}
				if !reflect.DeepEqual(result, tc.expected) {
					t.Errorf("expected %+#v, got %+#v", tc.expected, result)
				}

				// Test Unmarshal
				unmarshalled := &InitScript{}
				err := unmarshalled.Unmarshal(result)
				if err != nil {
					t.Error("expected error, got nil")
				}
				if !reflect.DeepEqual(result, tc.expected) {
					t.Errorf("expected %+#v, got %+#v", tc.expected, result)
				}
			}
		})
	}
}
