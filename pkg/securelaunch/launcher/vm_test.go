// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package launcher

import (
	"encoding/hex"
	"path/filepath"
	"testing"
	"time"

	"github.com/hugelgupf/vmtest/govmtest"
	"github.com/hugelgupf/vmtest/guest"
	"github.com/hugelgupf/vmtest/qemu"
	"github.com/u-root/u-root/pkg/core/cp"
	slaunch "github.com/u-root/u-root/pkg/securelaunch"
)

const kernelStr = "" +
	"3ee77a9d7add191b33b34e577fdd4bc3c67a9b96964c1c32c4e0f186c076ccce" +
	"8c4f656296a81d2f5a80c9c100e53e70c5109aca2d1fbaceb0ba86cc6657b071" +
	"2297a0a6f802f8d9812b81d7cf0b663feb923affa4dd209ac334aa45fe9f33c9" +
	"531298bae5a825361d4407e5bf4c9dffde9f6043ca99787a90abae9ae13d10b9" +
	"49d3bd5a44cb46597ca84158552c9395942e544ce270fd87e3af7e12a3a6a2c3" +
	"7220571a68ec2dd1c122e2300afe96c95288fcbc0d07933dd1e8c89b9dae45fe" +
	"bafa534137d3ed603ab0b244c42c324b962bff956e29531150cd5d0d54d01076" +
	"12f34ac78dc380428190634c780efe1595d9ba976a764710e137f77f0751dcd5" +
	"4d307a5b7ca7fe3fe77fd63e480f4a215eb0149afe683800d65722436837ea8f" +
	"caa5cf92d78d414e113dc22df1428b76a4dd56d6555c92054f9b069b9260eea4" +
	"372d1ab35bd2ac7444b8cc9f969cc3c3629a178f81a2ce94a1c10fd3700d6777" +
	"2f21cd888026a51287efa83fde0e3aa3940161369fdc15225b95f63ae64954ff" +
	"fe9c22418e5928c0cbf57e93a210cb8e38c1c47bc6ed6cd1367b3671fb69c6c7" +
	"f58a6262a94a309d5d5f3e27d0f29be2bd41b9a046053a0898f2a90899299297" +
	"4a1afbcc021f06b4196307f7e1eb6ec04d3d7cf907acaa656f07c96d0484cecf" +
	"8b57361597c1646556d1f69371977019d4e0030daac15412da94e92dc3e6002a"

const kernelHashStr = "00b01a7861b9b50e0f46be18f6cb71fb1776c5ede6f7d2d73561487b85495000"

const initrdStr = "" +
	"1c1a49bab0b3a2518035d9da63235e19b06bf7147eff209adad78665ecfeb1d1" +
	"8c63fb6c5a6c6f527d1b482ce50d341d01857e7ac4c3e2fcb3b92f1056823b7a" +
	"ddd2a904df9b30823f2406bac56ff83668caa01ab957ee89c59db8aa84673280" +
	"0f852e4a6fa9fef0d8c2b7b9994b1b48aa679524ecc6c5264defecce62848208" +
	"c4edefc234e64a21ff7ea3b89c4630e8b285533bbe24da531b47f1e47d780006" +
	"f9512dab04373def1adaa480b238e173ff63ecc3f49ec2e4e75f0a3e6307ae9c" +
	"502d6f9699b7411deb1067cbd116877c3194cb707ab7590ee00113954b7e9383" +
	"39f21a7362ef9f6c205efdc60ddc92ceab9bfa975634681c22035de68c080495" +
	"2bc0491cb39cba16aa01c5d04b95e2d3e2d5d427540f91dcf85ec3ee68a31754" +
	"846cdbb6e4be913b73595f5febdc55778ead54f6a3d1cf685880853d21469aff" +
	"44650373b5f23cdc3ee86bd0872c05ca6f0ca3236f0825986350b44a23922c01" +
	"50823da31431eeac6f9d857eabd734be83600ffa4567550add5360e2497419c4" +
	"6b37beb75866648204cc45bf547a5a4ff7312ffd38e708c563b9919e949ca7a0" +
	"9e655205ccef8c3556f12d789f4b5d06f4eb0e8e50a0ca5e422ecd06996ea8c2" +
	"f7da3547b62e8660024f3467d35d6203e33912a9256500caf80f205448f09a28" +
	"0890039d73dc373ae5ceeea6889b4579cba367a5fa5e210b2b605553ced7f369"

const initrdHashStr = "0853abe9bab31dbaf174a66f7e1215d2fc0dee37a8ad2b6d215749dab9d2073b"

var bootEntries = make(map[string]BootEntry)

func TestSecureLaunchLauncherVM(t *testing.T) {
	qemu.SkipIfNotArch(t, qemu.ArchAMD64)

	d := t.TempDir()
	mbrdisk := filepath.Join(d, "mbrdisk")
	if err := cp.Default.Copy("testdata/mbrdisk", mbrdisk); err != nil {
		t.Fatalf("copying testdata/mbrdisk to %q:got %v, want nil", mbrdisk, err)
	}

	govmtest.Run(t, "vm",
		govmtest.WithPackageToTest("github.com/u-root/u-root/pkg/securelaunch/launcher"),
		govmtest.WithQEMUFn(
			qemu.WithVMTimeout(2*time.Minute),
			// CONFIG_ATA_PIIX is required for this option to work.
			qemu.ArbitraryArgs("-hda", mbrdisk),

			// With NVMe devices enabled, kernel crashes when not using q35 machine model.
			qemu.ArbitraryArgs("-machine", "q35"),
		),
	)
}

func TestSecureLaunchLauncherMatchBootEntry(t *testing.T) {
	guest.SkipIfNotInVM(t)

	const blk = "sda1"

	slaunch.Debug = t.Logf
	if _, err := slaunch.GetStorageDevice("sda1"); err != nil {
		t.Skipf("no devices match %v:%v", blk, err)
	}

	kernelFile := blk + ":" + "/vmlinux"
	initrdFile := blk + ":" + "/initrd"

	kernelBytes, err := hex.DecodeString(kernelStr)
	if err != nil {
		t.Fatalf(`hex.DecodeString(kernelStr) = %v, not nil`, err)
	}
	if err := slaunch.WriteFile(kernelBytes, kernelFile); err != nil {
		t.Fatalf(`WriteFile(kernelBytes, kernelFile) = %v, not nil`, err)
	}

	initrdBytes, err := hex.DecodeString(initrdStr)
	if err != nil {
		t.Fatalf(`hex.DecodeString(initrdStr) = %v, not nil`, err)
	}
	if err := slaunch.WriteFile(initrdBytes, initrdFile); err != nil {
		t.Fatalf(`WriteFile(initrdBytes, initrdFile) = %v, not nil`, err)
	}

	entryName := "entry_name"
	bootEntries[entryName] = BootEntry{
		KernelName: kernelFile,
		KernelHash: kernelHashStr,
		InitrdName: initrdFile,
		InitrdHash: initrdHashStr,
	}

	if err := MatchBootEntry(entryName, bootEntries); err != nil {
		t.Fatalf(`MatchBootEntry() = %v, not nil`, err)
	}
}
