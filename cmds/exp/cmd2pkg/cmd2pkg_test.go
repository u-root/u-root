// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/ulog"
)

var echo = `package main
import "fmt"
func main(){
fmt.Printf("hi\n")
}
`

var echopkg = `
package main

import "echo/xxx/echo"
func main() {
	echo.Main()
}
`

func TestCmd2Pkg(t *testing.T) {
	if _, err := exec.LookPath("go"); err != nil {
		t.Skipf("no go compiler: %v", err)
	}

	dir := t.TempDir()
	t.Chdir(dir)

	ef := filepath.Join(dir, "echo.go")
	eb := filepath.Join(dir, "echo")

	if err := os.WriteFile(ef, []byte(echo), 0666); err != nil {
		t.Fatalf("writing %s: got %v, want nil", ef, err)
	}

	// Make sure it builds ...
	if out, err := exec.Command("go", "build", "-o", eb, ef).CombinedOutput(); err != nil {
		t.Fatalf("building %s: got (%s, %v), want nil", ef, string(out), err)
	}

	echoOut := "hi\n"

	if out, err := exec.Command(eb, "echo").CombinedOutput(); err != nil || string(out) != echoOut {
		t.Fatalf("running %s: got (%s, %v), want nil", eb, string(out), err)
	}

	if out, err := exec.Command("sh", "-c", "go mod init echo && go mod tidy").CombinedOutput(); err != nil {
		t.Fatalf("go mod init and tidy %s: got (%s, %v), want nil", dir, string(out), err)
	}

	if err := (&command{dir: filepath.Join("xxx"), l: ulog.Log}).execute("."); err != nil {
		t.Fatalf("execute cmd at %v to package: got %v, want nil", dir, err)
	}

	// Try to build it
	if err := os.WriteFile(ef, []byte(echopkg), 0666); err != nil {
		t.Fatalf("writing %s: got %v, want nil", ef, err)
	}

	// Make sure it builds ...
	if out, err := exec.Command("go", "build", "-o", eb, ef).CombinedOutput(); err != nil {
		t.Errorf("building %s: got (%s, %v), want nil", ef, string(out), err)
	}

	// and should still have the same output
	if out, err := exec.Command(eb, "echo").CombinedOutput(); err != nil || string(out) != echoOut {
		t.Fatalf("running %s: got (%s, %v), want nil", eb, string(out), err)
	}

}
