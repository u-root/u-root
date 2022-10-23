// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kmodule

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"testing"
)

var procModsMock = `hid_generic 16384 0 - Live 0x0000000000000000
usbhid 49152 0 - Live 0x0000000000000000
ccm 20480 6 - Live 0x0000000000000000
`

func TestProbeOptions(t *testing.T) {
	// Put modules.dep in the right place
	root := t.TempDir()
	mods := "usr/lib/modules/5.15.0-52-generic"
	moddir := filepath.Join(root, mods)
	dir9p := filepath.Join(moddir, "kernel/net/9p")
	if err := os.MkdirAll(dir9p, 0755); err != nil {
		t.Fatalf("os.Mkdirall(%q, 0755): %v != nil", dir9p, err)
	}
	if err := os.WriteFile(filepath.Join(moddir, "modules.dep"), []byte(modulesDotDep), 0666); err != nil {
		t.Fatalf("Writing modules.Dep: %v != nil", err)
	}
	// Since we are not loading the files, and we are ignoring /proc/modules, the .ko for this test can
	// be empty files.
	for _, n := range []string{"9pnet_virtio.ko", "9pnet.ko"} {
		p := filepath.Join(dir9p, n)
		if err := os.WriteFile(p, []byte{}, 0666); err != nil {
			t.Fatalf("Writing %q: %v != nil", p, err)
		}
	}

	l, err := NewPath("proc.modules")
	if err != nil {
		t.Fatalf(`NewPath("proc.modules"): %v != nil`, err)
	}
	for _, m := range []string{"cpuid", "9pnet_virtio"} {
		o := ProbeOpts{IgnoreProcMods: true, DryRunCB: func(_ string) {}, RootDir: root, KVer: "5.15.0-52-generic"}
		if err := l.ProbeOptions(m, "", o); err != nil {
			t.Errorf("ProbeOptions for %q: %v != nil", m, err)
		}
	}
}

func TestGenLoadedMods(t *testing.T) {
	m := depMap{
		"/lib/modules/6.6.6-generic/kernel/drivers/hid/hid-generic.ko":   &dependency{},
		"/lib/modules/6.6.6-generic/kernel/drivers/hid/usbhid/usbhid.ko": &dependency{},
		"/lib/modules/6.6.6-generic/kernel/crypto/ccm.ko":                &dependency{},
	}
	br := bytes.NewBufferString(procModsMock)
	l, err := NewPath("proc.modules")
	if err != nil {
		t.Fatalf("New(): got %v, want nil", err)
	}
	err = l.genLoadedMods(br, m)
	if err != nil {
		t.Fatalf("fail to l.genLoadedMods: %v\n", err)
	}
	for mod, d := range m {
		if d.state != loaded {
			t.Fatalf("mod %q should have been loaded", path.Base(mod))
		}
	}
}

// bad is an io.Writer that fails badly.
type bad struct{}

func (b *bad) Write([]byte) (int, error) {
	return -1, os.ErrInvalid
}

func TestProbePath(t *testing.T) {
	d := t.TempDir()
	n := filepath.Join(d, "nosuchfile")
	if err := ProbePath(n, "9p", ""); !errors.Is(err, os.ErrNotExist) {
		t.Errorf(`ProbePath(%q, "9p", ""): %v != %v`, n, err, os.ErrNotExist)
	}
	// This is here to test the parts of ProbePath that are used if
	// the first arg is a real file. There are so many possible failures
	// past that point that we ignore the return.
	ProbePath("proc.modules", "9p", "")
}

func TestPretty(t *testing.T) {
	var lsmod = "Module\tSize\tUsed by\ncpuid\t16384\t0\t-\nLine 1 is malformed: \"bogus\"\nvhost_vsock\t\t\t\t24576\t0\t-\nvmw_vsock_virtio_transport_common\t49152\t1\tvhost_vsock,\nvhost\t\t\t\t\t57344\t1\tvhost_vsock,\nvhost_iotlb\t\t\t\t16384\t1\tvhost,\nvsock\t\t\t\t\t53248\t2\tvhost_vsock,vmw_vsock_virtio_transport_common,\nveth\t\t\t\t\t36864\t0\t-\nxt_conntrack\t\t\t\t16384\t2\t-\nxt_MASQUERADE\t\t\t\t20480\t2\t-\nnf_conntrack_netlink\t\t\t57344\t0\t-\nxfrm_user\t\t\t\t49152\t1\t-\nxfrm_algo\t\t\t\t16384\t1\txfrm_user,\nxt_addrtype\t\t\t\t16384\t2\t-\nnft_compat\t\t\t\t20480\t6\t-\nbr_netfilter\t\t\t\t36864\t0\t-\nbridge\t\t\t\t\t307200\t1\tbr_netfilter,\nstp\t\t\t\t\t16384\t1\tbridge,\nllc\t\t\t\t\t16384\t2\tbridge,stp,\ntcp_diag\t\t\t\t16384\t0\t-\ninet_diag\t\t\t\t24576\t1\ttcp_diag,\ntls\t\t\t\t\t118784\t0\t-\nrfkill\t\t\t\t\t32768\t4\t-\noverlay\t\t\t\t\t155648\t0\t-\nnft_chain_nat\t\t\t\t16384\t3\t-\niptable_nat\t\t\t\t16384\t0\t-\nnf_nat\t\t\t\t\t57344\t3\txt_MASQUERADE,nft_chain_nat,iptable_nat,\nnf_conntrack\t\t\t\t180224\t4\txt_conntrack,xt_MASQUERADE,nf_conntrack_netlink,nf_nat,\nnf_defrag_ipv6\t\t\t\t24576\t1\tnf_conntrack,\nnf_defrag_ipv4\t\t\t\t16384\t1\tnf_conntrack,\niptable_filter\t\t\t\t16384\t0\t-\nnf_tables\t\t\t\t278528\t104\tnft_compat,nft_chain_nat,\nlibcrc32c\t\t\t\t16384\t3\tnf_nat,nf_conntrack,nf_tables,\nnfnetlink\t\t\t\t20480\t4\tnf_conntrack_netlink,nft_compat,nf_tables,\ncpufreq_ondemand\t\t\t20480\t0\t-\ncpufreq_userspace\t\t\t20480\t0\t-\ncpufreq_conservative\t\t\t16384\t0\t-\ncpufreq_powersave\t\t\t20480\t0\t-\nsunrpc\t\t\t\t\t659456\t1\t-\nintel_rapl_msr\t\t\t\t20480\t0\t-\nintel_rapl_common\t\t\t28672\t1\tintel_rapl_msr,\nnfit\t\t\t\t\t69632\t0\t-\nbinfmt_misc\t\t\t\t24576\t1\t-\nlibnvdimm\t\t\t\t200704\t1\tnfit,\nnls_ascii\t\t\t\t16384\t1\t-\nnls_cp437\t\t\t\t20480\t1\t-\nvfat\t\t\t\t\t20480\t1\t-\nkvm_intel\t\t\t\t372736\t0\t-\nfat\t\t\t\t\t86016\t1\tvfat,\nkvm\t\t\t\t\t1081344\t1\tkvm_intel,\nirqbypass\t\t\t\t16384\t1\tkvm,\nrapl\t\t\t\t\t20480\t0\t-\npvpanic_mmio\t\t\t\t16384\t0\t-\nsg\t\t\t\t\t40960\t0\t-\npvpanic\t\t\t\t\t16384\t1\tpvpanic_mmio,\nserio_raw\t\t\t\t20480\t0\t-\nevdev\t\t\t\t\t28672\t3\t-\nefi_pstore\t\t\t\t16384\t0\t-\nparport_pc\t\t\t\t40960\t0\t-\nppdev\t\t\t\t\t24576\t0\t-\nlp\t\t\t\t\t20480\t0\t-\nparport\t\t\t\t\t73728\t3\tparport_pc,ppdev,lp,\ntcp_bbr\t\t\t\t\t20480\t0\t-\nfuse\t\t\t\t\t172032\t25\t-\ndrm\t\t\t\t\t622592\t0\t-\nconfigfs\t\t\t\t57344\t1\t-\nqemu_fw_cfg\t\t\t\t20480\t0\t-\nvirtio_rng\t\t\t\t16384\t0\t-\nrng_core\t\t\t\t20480\t2\tvirtio_rng,\nip_tables\t\t\t\t36864\t2\tiptable_nat,iptable_filter,\nx_tables\t\t\t\t57344\t7\txt_conntrack,xt_MASQUERADE,xt_addrtype,nft_compat,iptable_nat,iptable_filter,ip_tables,\nautofs4\t\t\t\t\t53248\t5\t-\next4\t\t\t\t\t954368\t2\t-\ncrc16\t\t\t\t\t16384\t1\text4,\nmbcache\t\t\t\t\t16384\t1\text4,\njbd2\t\t\t\t\t163840\t1\text4,\ncrc32c_generic\t\t\t\t16384\t0\t-\nefivarfs\t\t\t\t16384\t1\t-\ndm_mod\t\t\t\t\t176128\t7\t-\nsd_mod\t\t\t\t\t61440\t4\t-\nt10_pi\t\t\t\t\t16384\t1\tsd_mod,\ncrc64_rocksoft_generic\t\t\t16384\t1\t-\ncrc64_rocksoft\t\t\t\t20480\t1\tt10_pi,\ncrc_t10dif\t\t\t\t20480\t1\tt10_pi,\ncrct10dif_generic\t\t\t16384\t0\t-\ncrc64\t\t\t\t\t20480\t2\tcrc64_rocksoft_generic,crc64_rocksoft,\ncrct10dif_pclmul\t\t\t16384\t1\t-\ncrct10dif_common\t\t\t16384\t3\tcrc_t10dif,crct10dif_generic,crct10dif_pclmul,\ncrc32_pclmul\t\t\t\t16384\t0\t-\ncrc32c_intel\t\t\t\t24576\t4\t-\nghash_clmulni_intel\t\t\t16384\t0\t-\nvirtio_scsi\t\t\t\t24576\t3\t-\nvirtio_net\t\t\t\t69632\t0\t-\nnet_failover\t\t\t\t24576\t1\tvirtio_net,\nscsi_mod\t\t\t\t270336\t3\tsg,sd_mod,virtio_scsi,\nfailover\t\t\t\t16384\t1\tnet_failover,\nscsi_common\t\t\t\t16384\t2\tsg,scsi_mod,\naesni_intel\t\t\t\t380928\t0\t-\ncrypto_simd\t\t\t\t16384\t1\taesni_intel,\nvirtio_pci\t\t\t\t24576\t0\t-\nvirtio_pci_legacy_dev\t\t\t16384\t1\tvirtio_pci,\ncryptd\t\t\t\t\t28672\t2\tghash_clmulni_intel,crypto_simd,\npsmouse\t\t\t\t\t184320\t0\t-\npcspkr\t\t\t\t\t16384\t0\t-\nvirtio_pci_modern_dev\t\t\t20480\t1\tvirtio_pci,\ni2c_piix4\t\t\t\t28672\t0\t-\nbutton\t\t\t\t\t24576\t0\t-\nsha512_ssse3\t\t\t\t49152\t1\t-\nsha512_generic\t\t\t\t16384\t1\tsha512_ssse3,\n"
	l, err := NewPath("proc.modules")
	if err != nil {
		t.Fatalf(`NewPath("proc.modules"): %v != nil`, err)
	}
	var b bytes.Buffer
	if _, err := io.Copy(&b, l); err != nil {
		t.Fatalf("io.copy(&b, LinuxLoader): %v != nil", err)
	}
	var p bytes.Buffer
	if err := Pretty(&p, b.String()); err != nil {
		t.Fatalf("Pretty: %v != nil", err)
	}
	if len(p.String()) != len(lsmod) {
		t.Errorf("len(%d) != len(%d)", len(p.String()), len(lsmod))
	}
	if p.String() != lsmod {
		t.Errorf("lsmod: %q != %q", p.String(), lsmod)

	}

	if err := Pretty(&bad{}, b.String()); !errors.Is(err, os.ErrInvalid) {
		t.Errorf("Pretty(&bad{},...): %v != %v", err, os.ErrInvalid)
	}

}

func TestAll(t *testing.T) {
	l := &LinuxLoader{}
	_ = l.Init([]byte{0}, "")
	// Whether New succeeds or fails is difficult to be certain of.
	// The CI VM may or may not have /proc/modules.
	_, _ = New()
	// Still not supported on CI but call it anyway.
	_ = Probe("9p", "")
}

func TestFileInit(t *testing.T) {
	var errInit = errors.New("we got to l.Init")
	b := &badloader{err: errInit}
	// Try calling FileInit with all 4 file types. If this succeeds, then
	// we should get a return from the load of errInit
	files, err := filepath.Glob("data/*")
	if err != nil {
		t.Skipf("Skipping: no data")
	}
	for _, n := range files {
		f, err := os.Open(n)
		if err != nil {
			t.Errorf("os.Open(%q): %v != nil", n, err)
		}

		if err := b.FileInit(f, "", 0); err != errInit {
			t.Errorf("Loading %q: %v != %v", n, err, errInit)
		}
	}
}
