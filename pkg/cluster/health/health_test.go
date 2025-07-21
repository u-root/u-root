// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package health_test

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/u-root/u-root/pkg/cluster/health"
)

func TestList(t *testing.T) {
	for _, tt := range []struct {
		name string
		cmd  string
		args []string
		err  error
	}{
		{name: "bad command", cmd: "/dev/null", args: nil, err: os.ErrPermission},
		{name: "no names", cmd: "echo", args: nil},
		{name: "no names", cmd: "echo", args: []string{"1"}},
		{name: "no names", cmd: "echo", args: []string{"b", "1", "z"}},
	} {
		c := health.NewNodeList(tt.cmd, tt.args...)
		if err := c.Run(); !errors.Is(err, tt.err) {
			t.Errorf("running %s %v: got %v, want %v", tt.cmd, tt.args, err, tt.err)
			continue
		}
		t.Logf("c %v", c)
		if eq := slices.Compare(c.List, tt.args); eq != 0 {
			t.Errorf("compare (%v,%v): got %d, want 0", c.List, tt.args, eq)
			continue
		}
	}
}

// TestGather is a bit challenging, but a good test.
// We write a largely empty JSON to a file in tempdir, read it back,
// hope for the best.
func TestGather(t *testing.T) {
	d := t.TempDir()
	f := filepath.Join(d, "json")
	os.WriteFile(f, []byte(data), 0o666)

	health.V = t.Logf
	for _, tt := range []struct {
		name     string
		cmd      string // to run something on the node
		args     []string
		node     string // the command that is run on the node
		nodeargs []string
		err      error
	}{
		{name: "bad command", cmd: "/dev/null", err: os.ErrPermission},
		{name: "bad JSON", cmd: "echo", args: []string{"a"}},
		{name: "no file", cmd: "cat", node: filepath.Join(d, "k"), args: nil},
	} {
		n := health.NewNodeList("echo", f)
		if err := n.Run(); err != nil {
			t.Errorf("%s:{%v}.Run(): %v != nil", tt.name, n, err)
			continue
		}
		g := n.NewGather(tt.cmd, tt.args...)
		stats, err := g.Run(tt.node, tt.nodeargs...)
		if err != nil {
			t.Errorf("{%v}.Run: got %v, want nil", g, err)
		}

		if len(stats) == 0 {
			t.Errorf("stats: got %d, wanted at least 1", len(stats))
		}
		if len(stats[0].Err) == 0 && tt.err != nil {
			t.Errorf("running %s %s %v: got %v, want %v", tt.cmd, tt.node, tt.args, stats[0].Err, tt.err)
			continue
		}
	}
}
