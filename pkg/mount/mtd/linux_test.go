// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mtd

import (
	"bytes"
	"os"
	"testing"
)

var testString = "This is just a test string"

func TestOpen(t *testing.T) {
	m, err := NewDev(DevName)
	if err != nil {
		tmpDir := t.TempDir()
		testmtd, err := os.CreateTemp(tmpDir, "testmtd")
		if err != nil {
			t.Errorf(`os.Create(tmpDir, "testmtd")=file, %q, want file, nil`, err)
		}
		DevName = testmtd.Name()
		m, err = NewDev(DevName)
		if err != nil {
			t.Fatal(err)
		}
	}

	defer m.Close()

	if _, err := m.QueueWriteAt([]byte(testString), 0); err != nil {
		t.Errorf("m.QueueWrite([]byte(testString), 0)=-,%q, want _,nil", err)
	}
	buf := make([]byte, 26)
	if _, err := m.ReadAt(buf, 0); err != nil {
		t.Errorf("m.ReadAt([]byte(), 0)=-,%q, want _,nil", err)
	}
	if !bytes.Equal(buf, []byte(testString)) {
		t.Errorf("bytes.Equal(buf, []byte(testString))=false, want true")
	}
	if m.Name() != DevName {
		t.Errorf("want %s == %s, want m.DevName() == testmtd.Name()", m.Name(), DevName)
	}
}
