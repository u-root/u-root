// Copyright 2012-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netcat

import (
	"os"
	"strings"
	"testing"
)

func TestAccessControlOptions_IsAllowed(t *testing.T) {
	tests := []struct {
		name           string
		ac             AccessControlOptions
		hostNames      []string
		expectedResult bool
	}{
		{
			name: "Allowed host in list",
			ac: AccessControlOptions{
				SetAllowed:     true,
				ConnectionList: map[string]bool{"example.com": true},
			},
			hostNames:      []string{"example.com"},
			expectedResult: true,
		},
		{
			name: "Denied host in list",
			ac: AccessControlOptions{
				SetAllowed: true,
				ConnectionList: map[string]bool{
					"another-example.com": true,
					"example.com":         false,
				},
			},
			hostNames:      []string{"example.com"},
			expectedResult: false,
		},
		{
			name: "Host not in list, denied by default",
			ac: AccessControlOptions{
				SetAllowed:     true,
				ConnectionList: map[string]bool{"example.com": true},
			},
			hostNames:      []string{"notlisted.com"},
			expectedResult: false,
		},
		{
			name: "Host not in list, allowed by default when SetAllowed is false",
			ac: AccessControlOptions{
				SetAllowed:     false,
				ConnectionList: map[string]bool{"example.com": false},
			},
			hostNames:      []string{"notlisted.com"},
			expectedResult: true,
		},
		{
			name: "Denied host in list when SetAllowed is false",
			ac: AccessControlOptions{
				SetAllowed:     false,
				ConnectionList: map[string]bool{"example.com": false},
			},
			hostNames:      []string{"example.com"},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ac.IsAllowed(tt.hostNames); got != tt.expectedResult {
				t.Errorf("AccessControlOptions.IsAllowed() = %v, want %v", got, tt.expectedResult)
			}
		})
	}
}

func TestParseAccessControl(t *testing.T) {
	// Setup temporary files for testing
	allowFileContent := []string{"allowed.com"}
	denyFileContent := []string{"denied.com"}

	allowFileName := createTempFileWithContent(t, allowFileContent)
	defer os.Remove(allowFileName)

	denyFileName := createTempFileWithContent(t, denyFileContent)
	defer os.Remove(denyFileName)

	tests := []struct {
		name                   string
		connectionAllowFile    string
		connectionAllowList    []string
		connectionDenyFile     string
		connectionDenyList     []string
		expectedConnectionList map[string]bool
		expectError            bool
	}{
		{
			name:                   "Allow and Deny Files",
			connectionAllowFile:    allowFileName,
			connectionAllowList:    []string{},
			connectionDenyFile:     denyFileName,
			connectionDenyList:     []string{},
			expectedConnectionList: map[string]bool{"allowed.com": true, "denied.com": false},
			expectError:            false,
		},
		{
			name:                   "Allow and Deny Lists",
			connectionAllowFile:    "",
			connectionAllowList:    []string{"listallowed.com"},
			connectionDenyFile:     "",
			connectionDenyList:     []string{"listdenied.com"},
			expectedConnectionList: map[string]bool{"listallowed.com": true, "listdenied.com": false},
			expectError:            false,
		},
		{
			name:                   "Combined Allow File and List",
			connectionAllowFile:    allowFileName,
			connectionAllowList:    []string{"listallowed.com"},
			connectionDenyFile:     "",
			connectionDenyList:     []string{},
			expectedConnectionList: map[string]bool{"allowed.com": true, "listallowed.com": true},
			expectError:            false,
		},
		{
			name:                   "Combined Deny File and List",
			connectionAllowFile:    "",
			connectionAllowList:    []string{},
			connectionDenyFile:     denyFileName,
			connectionDenyList:     []string{"listdenied.com"},
			expectedConnectionList: map[string]bool{"denied.com": false, "listdenied.com": false},
			expectError:            false,
		},
		{
			name:                   "Allow List Overrides Deny File",
			connectionAllowFile:    "",
			connectionAllowList:    []string{"denied.com"},
			connectionDenyFile:     denyFileName,
			connectionDenyList:     []string{},
			expectedConnectionList: map[string]bool{"denied.com": true},
			expectError:            false,
		},
		{
			name:                   "Empty Allow and Deny Lists",
			connectionAllowFile:    "",
			connectionAllowList:    []string{},
			connectionDenyFile:     "",
			connectionDenyList:     []string{},
			expectedConnectionList: map[string]bool{},
			expectError:            false,
		},
		{
			name:                   "Nonexistent Deny Files",
			connectionAllowFile:    allowFileName,
			connectionAllowList:    []string{},
			connectionDenyFile:     "nonexistent_deny_file",
			connectionDenyList:     []string{},
			expectedConnectionList: map[string]bool{},
			expectError:            true, // Expect error due to nonexistent files
		},
		{
			name:                   "Nonexistent Allow Files",
			connectionAllowFile:    "nonexistent_allow_file",
			connectionAllowList:    []string{},
			connectionDenyFile:     denyFileName,
			connectionDenyList:     []string{},
			expectedConnectionList: map[string]bool{},
			expectError:            true, // Expect error due to nonexistent files
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac, err := ParseAccessControl(tt.connectionAllowFile, tt.connectionAllowList, tt.connectionDenyFile, tt.connectionDenyList)
			if (err != nil) != tt.expectError {
				t.Fatalf("ParseAccessControl() error = %v, wantErr %v", err, tt.expectError)
			}

			if !tt.expectError {
				if len(ac.ConnectionList) != len(tt.expectedConnectionList) {
					t.Errorf("Expected %d items in the ConnectionList, got %d", len(tt.expectedConnectionList), len(ac.ConnectionList))
				}

				for host, allowed := range tt.expectedConnectionList {
					if ac.ConnectionList[host] != allowed {
						t.Errorf("Expected %v for host %s, got %v", allowed, host, ac.ConnectionList[host])
					}
				}
			}
		})
	}
}

func createTempFileWithContent(t *testing.T, content []string) string {
	t.Helper()
	tmpDir := t.TempDir()

	tmpFile, err := os.CreateTemp(tmpDir, "example")
	if err != nil {
		t.Fatal(err)
	}
	defer tmpFile.Close()

	contentStr := strings.Join(content, "\n")
	_, err = tmpFile.WriteString(contentStr)
	if err != nil {
		t.Fatal(err)
	}

	return tmpFile.Name()
}
