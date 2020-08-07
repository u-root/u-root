// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iscsinl

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/sh"
	"github.com/u-root/u-root/pkg/testutil"
)

func TestMain(m *testing.M) {
	if os.Getuid() == 0 {
		if err := sh.RunWithLogs("dhclient", "-ipv6=false"); err != nil {
			log.Fatalf("could not configure network for tests: %v", err)
		}
	}

	os.Exit(m.Run())
}

func TestMountIscsi(t *testing.T) {
	testutil.SkipIfNotRoot(t)

	devices, err := MountIscsi(
		WithInitiator(os.Getenv("INITIATOR_NAME")),
		WithTarget(fmt.Sprintf("%s:%s", os.Getenv("TGT_SERVER"), os.Getenv("TGT_PORT")), os.Getenv("TGT_VOLUME")),
		WithDigests("None"),
	)
	if err != nil {
		log.Println(err)
		t.Error(err)
	}

	for _, device := range devices {
		// Make the mountpoint, and mount it.
		if err := os.MkdirAll("/mp", 0755); err != nil {
			t.Fatal(err)
		}

		mp, err := mount.TryMount(fmt.Sprintf("/dev/%s1", device), "/mp", 0)
		if err != nil {
			t.Fatal(err)
		}

		if err := sh.RunWithLogs("ls", "-l", "/mp"); err != nil {
			t.Error(err)
		}

		if err := mp.Unmount(0); err != nil {
			t.Error(err)
		}
	}
}
