// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package policy

import (
	"testing"
	"time"

	"github.com/hugelgupf/vmtest"
	"github.com/hugelgupf/vmtest/guest"
	"github.com/hugelgupf/vmtest/qemu"
	slaunch "github.com/u-root/u-root/pkg/securelaunch"
)

const policyStr = `{
    "collectors": [
      {
        "type": "dmi",
        "events": [
          {
            "label": "BIOS",
            "fields": []
          },
          {
            "label": "System",
            "fields": []
          },
          {
            "label": "Processor",
            "fields": []
          }
        ]
      },
      {
         "type": "files",
         "paths": [ "sda1:/opc/foo" ]
      },
      {
         "type": "storage",
         "paths": [ "sda1" ]
      },
      {
         "type": "cpuid",
         "location": "sda2:/cpuid.txt"
      }
    ],
    "attestor": {},
    "launcher": {
        "type": "kexec",
        "boot entries": {
            "0": {
                "kernel name":"vmlinuz-5.4.17-2036.103.3.el7uek.x86_64",
                "kernel hash":"59c762615cdb09558bcd80d3025d023b436386fd9ab6e09d709418fbbce7680c",
                "initrd name":"initramfs-5.4.17-2036.103.3.el7uek.x86_64.img",
                "initrd hash":"a39a6ba3e35dffd0b91ca0f0dee2a7bfb16a447353746cf83d6dc7139dc99c4a",
                "cmdline":"BOOT_IMAGE=/vmlinuz-5.4.17-2036.103.3.el7uek.x86_64 root=/dev/mapper/ol_ol7--sl-root ro crashkernel=auto rd.luks.uuid=luks-06f28824-6b55-4219-b1a4-69a466af670b rd.lvm.lv=ol_ol7-sl/root rd.lvm.lv=ol_ol7-sl/swap console=ttyS0,115200"
            },
            "1": {
                "kernel name":"vmlinuz-4.14.35-1902.303.5.3.el7uek.x86_64",
                "kernel hash":"fa17071a44c0c185de9f879cddf6823f4d64a0c26604657655dad7c1d2fae39c",
                "initrd name":"initramfs-4.14.35-1902.303.5.3.el7uek.x86_64.img",
                "initrd hash":"c409a5118dacb1c2c71b9dab078ff670f15cae5219475fba902b508dea616187",
                "cmdline":"console=ttyS0,115200n8 earlyprintk=serial,ttyS0,115200n8,keep printk.time=y"
            }
        }
    },
    "eventlog": {
        "type": "file",
        "location": "/evtlog"
    }
}
`

// VM setup:
//
//  /dev/sda is ./testdata/mbrdisk
//	  /dev/sda1 is ext4
//	  /dev/sda2 is vfat
//	  /dev/sda3 is fat32
//	  /dev/sda4 is xfs
//
//   ARM tests will load drives as virtio-blk devices (/dev/vd*)

func TestVM(t *testing.T) {
	vmtest.SkipIfNotArch(t, qemu.ArchAMD64)

	vmtest.RunGoTestsInVM(t, []string{"github.com/u-root/u-root/pkg/securelaunch/policy"},
		vmtest.WithVMOpt(vmtest.WithQEMUFn(
			qemu.WithVMTimeout(2*time.Minute),
			// CONFIG_ATA_PIIX is required for this option to work.
			qemu.ArbitraryArgs("-hda", "../testdata/mbrdisk"),

			// With NVMe devices enabled, kernel crashes when not using q35 machine model.
			qemu.ArbitraryArgs("-machine", "q35"),
		)),
	)
}

func TestParse(t *testing.T) {
	guest.SkipIfNotInVM(t)

	policyFile := "sda1:" + "/policy"

	if err := slaunch.WriteFile([]byte(policyStr), policyFile); err != nil {
		t.Fatalf(`WriteFile(policyStr.bytes, tempFile) = %v, not nil`, err)
	}

	if err := Load(policyFile); err != nil {
		t.Fatalf("Load() = %v, not nil", err)
	}

	policy, err := Parse()
	if err != nil {
		t.Fatalf("Parse() = %v, not nil", err)
	}
	if policy == nil {
		t.Fatalf("no policy file returned")
	}
}
