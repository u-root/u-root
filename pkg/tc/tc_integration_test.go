// Copyright 2012-20124 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trafficctl_test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/florianl/go-tc"
	"github.com/hugelgupf/vmtest/govmtest"
	"github.com/hugelgupf/vmtest/guest"
	"github.com/hugelgupf/vmtest/qemu"
	trafficctl "github.com/u-root/u-root/pkg/tc"
)

const (
	DummyInterface0 = "eth0"
	DummyInterface1 = "eth1"
)

func TestVM(t *testing.T) {
	govmtest.Run(t, "tc integration",
		govmtest.WithPackageToTest("github.com/u-root/u-root/pkg/tc"),
		govmtest.WithQEMUFn(
			qemu.WithVMTimeout(time.Minute*2),
			qemu.ArbitraryArgs("-nic", fmt.Sprintf("user,id=%s", DummyInterface0)),
			qemu.ArbitraryArgs("-nic", fmt.Sprintf("user,id=%s", DummyInterface1)),
		),
	)
}

func TestQDiscAdd(t *testing.T) {
	guest.SkipIfNotInVM(t)

	rtnl, err := tc.Open(&tc.Config{})
	if err != nil {
		t.Error(err)
	}
	defer rtnl.Close()

	tctl := &trafficctl.Trafficctl{Tc: rtnl}

	for _, tt := range []struct {
		name   string
		args   []string
		err    error
		output string
	}{
		{
			name: "Add Ingress",
			args: []string{
				"dev",
				DummyInterface0,
				"ingress",
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var outbuf bytes.Buffer
			args, err := trafficctl.ParseQDiscArgs(&outbuf, tt.args)
			if !errors.Is(err, tt.err) {
				t.Errorf("ParseQDiscArgs() = %v, not %v", err, tt.err)
			}

			if err := tctl.AddQdisc(&outbuf, args); !errors.Is(err, tt.err) {
				t.Errorf("AddQdisc() = %v, not %v", err, tt.err)
			}
		})
	}
}
