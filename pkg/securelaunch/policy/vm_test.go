// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package policy

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/hugelgupf/vmtest/govmtest"
	"github.com/hugelgupf/vmtest/guest"
	"github.com/hugelgupf/vmtest/qemu"
	slaunch "github.com/u-root/u-root/pkg/securelaunch"
)

const blk = "sda1"

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

const pubkeyStr = `
-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAmsvG6Bj+cylwEmrUrqKj
W7lvovDDjJxizSrFTizd9gid/V9AT8mbwwoJJp2S2WVdGHPn3hDuFlyQvKTkqQlE
6l1+ALEiAxjiGrKz6c+x21X1thHv01w9/NX/f7K3B1QNnj972k96z6PW0jjafxYe
oc+ylGEqh4GxOlYbcdELdXi261+n+I0CsDJEzLVJnbs7YGW0dPN4UlGIAoVe4O4x
Taz74TBJEODD6GrNLpZeDD2ke0KwZSjMF2lDBX40Oyj/yn10dYUqKvfODNvWYIQg
fRCANXtNC3u89wU0HamJj8agITLOqIxp2UOaBj0qF4lAN4PP+pF6+6r0bngmfy0w
3hnEFVplySGp/S94071hR7zAMd4ZSIACcUQjFvt5BqmHtXgDRviJc4lh8niZ7cg0
ID8PgLjn+wLL5EZKsztWsk4030TeGtZ82eMiczqpsj7/5SWsxRBUraEi5aD6CREG
LIRbeVynOqC8KPXWLWDUvZSgwWUOYykZHYY9qCU4UgWNAx/uOfl4ZHEVK4LtrNI0
+byI7V27o9f5uhz5XJiAjObTAyh31xP2p4Xu6wrCQy1E4i1l9gTsCl6rQAPY0K+A
U8CAmAOKf2E5jEe8itnnsKu6Ii46ndLwJrBGgjsEYRNVln04owEiTsPrUCTPBS3O
RaNetBXso+Hgdg8M25WgCI8CAwEAAQ==
-----END PUBLIC KEY-----
`

const signatureStr = "" +
	"5e344a9ee66b03445d83995f83bace0227d303082e824ff7137d9362e636406a" +
	"ebdfa000b737dd9374d51278753356585aeedd2d8a2c884f2dc8036710f4a24c" +
	"050fd129d2eb1b95ffb4c39cd2de58ad4b0c55d6480eea54d9a0288ba65bdef5" +
	"7e5398dcf2b2510dcde1fb3bb1a7e0c17e00683a0dbfd53ad4170b2ac24741b7" +
	"dfff09ab7b6a960e7f6937ff85b6289f9a1e83a1abeb653ae0aa2682809b6242" +
	"2adb0130cad9ba9db54e0af7be0b5b008dcb4850e049a7b3b238f8feb5bbc91e" +
	"97c55e3becaf7a1fec21aca19a1bb6d283c8c7848686746344ed1481a57154d2" +
	"56c20efdfe0121097734222d02f0529372ab7ab296ec0b244e830c5540270e7f" +
	"4fa729ec35328f63f83b3e6494187c254372261970764fb99438976d904bcffb" +
	"1bb9e0428cc03e50595f07273fc387a4df1611275bc82fb156d9bf0ce8bbd0b7" +
	"f6e008835c9f701e25afd6a531e7c8ab354c3ff9672a5354306091f919fa46bf" +
	"8d3ca668f837ccf7a3adc6841c0c416ecafd299de3c5e03cff28daad9efd5c2a" +
	"2de4526869a209fa6e5703ea6b63d40cb534e1f192458978c475dd155d919fe8" +
	"ff16685d488051f3618b4a3112154a12206e321299b6e0a6a4ddff623a2f4f08" +
	"4d254392e9285aa245242505bb2d94e6d6fd4ada8d12a6ffd3440227e17771e6" +
	"68404e53b87a8a74314f70a8bc710c83459932357d068ad62cd5feeacae6f5c1"

// VM setup:
//
//	 /dev/sda is ../testdata/mbrdisk
//	 it must be copied because there is another test that uses it and
//	 qemu will fail if they run at the same time.
//		  /dev/sda1 is ext4
//		  /dev/sda2 is vfat
//		  /dev/sda3 is fat32
//		  /dev/sda4 is xfs
//
//	  ARM tests will load drives as virtio-blk devices (/dev/vd*)

func TestVM(t *testing.T) {
	qemu.SkipIfNotArch(t, qemu.ArchAMD64)

	govmtest.Run(t, "vmpolicy",
		govmtest.WithPackageToTest("github.com/u-root/u-root/pkg/securelaunch/policy"),
		govmtest.WithQEMUFn(
			qemu.WithVMTimeout(2*time.Minute),
			// CONFIG_ATA_PIIX is required for this option to work.
			qemu.ArbitraryArgs("-hda", "testdata/mbrdisk"),

			// With NVMe devices enabled, kernel crashes when not using q35 machine model.
			qemu.ArbitraryArgs("-machine", "q35"),
		),
	)
}

func TestParse(t *testing.T) {
	guest.SkipIfNotInVM(t)

	slaunch.Debug = t.Logf
	if _, err := slaunch.GetStorageDevice(blk); err != nil {
		t.Skipf("no devices match %v:%v", blk, err)
	}

	policyFile := blk + ":" + "/policy"

	if err := slaunch.WriteFile([]byte(policyStr), policyFile); err != nil {
		t.Fatalf(`WriteFile(policyStr.bytes, tempFile) = %v, not nil`, err)
	}

	if err := Load(policyFile, "", ""); err != nil {
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

func TestVerify(t *testing.T) {
	guest.SkipIfNotInVM(t)

	slaunch.Debug = t.Logf
	if _, err := slaunch.GetStorageDevice(blk); err != nil {
		t.Skipf("no devices match %v:%v", blk, err)
	}

	policyFile := blk + ":" + "/securelaunch.policy"
	pubkeyFile := blk + ":" + "/securelaunch.pubkey"
	signatureFile := blk + ":" + "/securelaunch.sig"

	if err := slaunch.WriteFile([]byte(policyStr), policyFile); err != nil {
		t.Fatalf(`WriteFile(policyStr.bytes, policyFile) = %v, not nil`, err)
	}

	if err := slaunch.WriteFile([]byte(pubkeyStr), pubkeyFile); err != nil {
		t.Fatalf(`WriteFile(pubkeyStr.bytes, pubkeyFile) = %v, not nil`, err)
	}

	signatureBytes, err := hex.DecodeString(signatureStr)
	if err != nil {
		t.Fatalf(`hex.DecodeString(signatureStr) = %v, not nil`, err)
	}
	if err := slaunch.WriteFile(signatureBytes, signatureFile); err != nil {
		t.Fatalf(`WriteFile(signatureBytes, signatureFile) = %v, not nil`, err)
	}

	if err := Load(policyFile, pubkeyFile, signatureFile); err != nil {
		t.Fatalf("Load() = %v, not nil", err)
	}

	if err := Verify(); err != nil {
		t.Fatalf("Verify() = %v, not nil", err)
	}
}
