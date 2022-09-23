// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kmodule

var modulesDotDep = `kernel/arch/x86/events/amd/amd-uncore.ko:
kernel/arch/x86/events/intel/intel-cstate.ko:
kernel/arch/x86/events/rapl.ko:
kernel/arch/x86/kernel/cpu/mce/mce-inject.ko:
kernel/arch/x86/kernel/msr.ko:
kernel/arch/x86/kernel/cpuid.ko:
kernel/arch/x86/crypto/twofish-x86_64.ko: kernel/crypto/twofish_common.ko
kernel/arch/x86/crypto/twofish-x86_64-3way.ko: kernel/arch/x86/crypto/twofish-x86_64.ko kernel/crypto/twofish_common.ko
kernel/arch/x86/crypto/twofish-avx-x86_64.ko: kernel/crypto/crypto_simd.ko kernel/crypto/cryptd.ko kernel/arch/x86/crypto/twofish-x86_64-3way.ko kernel/arch/x86/crypto/twofish-x86_64.ko kernel/crypto/twofish_common.ko
kernel/arch/x86/crypto/serpent-sse2-x86_64.ko: kernel/crypto/serpent_generic.ko kernel/crypto/crypto_simd.ko kernel/crypto/cryptd.ko
kernel/arch/x86/crypto/serpent-avx-x86_64.ko: kernel/crypto/serpent_generic.ko kernel/crypto/crypto_simd.ko kernel/crypto/cryptd.ko
kernel/arch/x86/crypto/serpent-avx2.ko: kernel/arch/x86/crypto/serpent-avx-x86_64.ko kernel/crypto/serpent_generic.ko kernel/crypto/crypto_simd.ko kernel/crypto/cryptd.ko
kernel/arch/x86/crypto/des3_ede-x86_64.ko: kernel/lib/crypto/libdes.ko
kernel/arch/x86/crypto/camellia-x86_64.ko:
kernel/arch/x86/crypto/camellia-aesni-avx-x86_64.ko: kernel/arch/x86/crypto/camellia-x86_64.ko kernel/crypto/crypto_simd.ko kernel/crypto/cryptd.ko
kernel/arch/x86/crypto/camellia-aesni-avx2.ko: kernel/arch/x86/crypto/camellia-aesni-avx-x86_64.ko kernel/arch/x86/crypto/camellia-x86_64.ko kernel/crypto/crypto_simd.ko kernel/crypto/cryptd.ko
kernel/arch/x86/crypto/blowfish-x86_64.ko: kernel/crypto/blowfish_common.ko
kernel/arch/x86/crypto/cast5-avx-x86_64.ko: kernel/crypto/cast5_generic.ko kernel/crypto/cast_common.ko kernel/crypto/crypto_simd.ko kernel/crypto/cryptd.ko
kernel/arch/x86/crypto/cast6-avx-x86_64.ko: kernel/crypto/cast6_generic.ko kernel/crypto/cast_common.ko kernel/crypto/crypto_simd.ko kernel/crypto/cryptd.ko
kernel/arch/x86/crypto/aegis128-aesni.ko: kernel/crypto/crypto_simd.ko kernel/crypto/cryptd.ko
kernel/arch/x86/crypto/chacha-x86_64.ko: kernel/lib/crypto/libchacha.ko
kernel/arch/x86/crypto/aesni-intel.ko: kernel/crypto/crypto_simd.ko kernel/crypto/cryptd.ko
kernel/arch/x86/crypto/sha1-ssse3.ko:
kernel/arch/x86/crypto/sha256-ssse3.ko:
kernel/arch/x86/crypto/sha512-ssse3.ko:
kernel/arch/x86/crypto/blake2s-x86_64.ko:
kernel/arch/x86/crypto/ghash-clmulni-intel.ko: kernel/crypto/cryptd.ko
kernel/arch/x86/crypto/crc32-pclmul.ko:
kernel/arch/x86/crypto/crct10dif-pclmul.ko:
kernel/arch/x86/crypto/poly1305-x86_64.ko:
kernel/arch/x86/crypto/nhpoly1305-sse2.ko: kernel/crypto/nhpoly1305.ko kernel/lib/crypto/libpoly1305.ko
kernel/arch/x86/crypto/nhpoly1305-avx2.ko: kernel/crypto/nhpoly1305.ko kernel/lib/crypto/libpoly1305.ko
kernel/arch/x86/crypto/curve25519-x86_64.ko: kernel/lib/crypto/libcurve25519-generic.ko
kernel/arch/x86/crypto/sm4-aesni-avx-x86_64.ko: kernel/lib/crypto/libsm4.ko kernel/crypto/crypto_simd.ko kernel/crypto/cryptd.ko
kernel/arch/x86/crypto/sm4-aesni-avx2-x86_64.ko: kernel/arch/x86/crypto/sm4-aesni-avx-x86_64.ko kernel/lib/crypto/libsm4.ko kernel/crypto/crypto_simd.ko kernel/crypto/cryptd.ko
kernel/arch/x86/platform/atom/punit_atom_debug.ko:
kernel/arch/x86/kvm/kvm.ko:
kernel/arch/x86/kvm/kvm-intel.ko: kernel/arch/x86/kvm/kvm.ko
kernel/arch/x86/kvm/kvm-amd.ko: kernel/drivers/crypto/ccp/ccp.ko kernel/arch/x86/kvm/kvm.ko
kernel/kernel/kheaders.ko:
kernel/mm/hwpoison-inject.ko:
kernel/mm/z3fold.ko:
kernel/fs/nfs_common/nfs_acl.ko: kernel/net/sunrpc/sunrpc.ko
kernel/fs/nfs_common/grace.ko:
kernel/fs/quota/quota_v1.ko:
kernel/fs/quota/quota_v2.ko: kernel/fs/quota/quota_tree.ko
kernel/fs/quota/quota_tree.ko:
kernel/fs/fat/msdos.ko:
kernel/fs/nls/nls_cp737.ko:
kernel/fs/nls/nls_cp775.ko:
kernel/fs/nls/nls_cp850.ko:
kernel/fs/nls/nls_cp852.ko:
kernel/fs/nls/nls_cp855.ko:
kernel/fs/nls/nls_cp857.ko:
kernel/fs/nls/nls_cp860.ko:
kernel/fs/nls/nls_cp861.ko:
kernel/fs/nls/nls_cp862.ko:
kernel/fs/nls/nls_cp863.ko:
kernel/fs/nls/nls_cp864.ko:
kernel/fs/nls/nls_cp865.ko:
kernel/fs/nls/nls_cp866.ko:
kernel/fs/nls/nls_cp869.ko:
kernel/fs/nls/nls_cp874.ko:
kernel/fs/nls/nls_cp932.ko:
kernel/fs/nls/nls_euc-jp.ko:
kernel/fs/nls/nls_cp936.ko:
kernel/fs/nls/nls_cp949.ko:
kernel/fs/nls/nls_cp950.ko:
kernel/fs/nls/nls_cp1250.ko:
kernel/fs/nls/nls_cp1251.ko:
kernel/fs/nls/nls_ascii.ko:
kernel/fs/nls/nls_iso8859-1.ko:
kernel/fs/nls/nls_iso8859-2.ko:
kernel/fs/nls/nls_iso8859-3.ko:
kernel/fs/nls/nls_iso8859-4.ko:
kernel/fs/nls/nls_iso8859-5.ko:
kernel/fs/nls/nls_iso8859-6.ko:
kernel/fs/nls/nls_iso8859-7.ko:
kernel/fs/nls/nls_cp1255.ko:
kernel/fs/nls/nls_iso8859-9.ko:
kernel/fs/nls/nls_iso8859-13.ko:
kernel/fs/nls/nls_iso8859-14.ko:
kernel/fs/nls/nls_iso8859-15.ko:
kernel/fs/nls/nls_koi8-r.ko:
kernel/fs/nls/nls_koi8-u.ko:
kernel/fs/nls/nls_koi8-ru.ko:
kernel/fs/nls/nls_utf8.ko:
kernel/fs/nls/mac-celtic.ko:
kernel/fs/nls/mac-centeuro.ko:
kernel/fs/nls/mac-croatian.ko:
kernel/fs/nls/mac-cyrillic.ko:
kernel/fs/nls/mac-gaelic.ko:
kernel/fs/nls/mac-greek.ko:
kernel/fs/nls/mac-iceland.ko:
kernel/fs/nls/mac-inuit.ko:
kernel/fs/nls/mac-romanian.ko:
kernel/fs/nls/mac-roman.ko:
kernel/fs/nls/mac-turkish.ko:
kernel/fs/fuse/cuse.ko:
kernel/fs/fuse/virtiofs.ko:
kernel/fs/pstore/ramoops.ko: kernel/lib/reed_solomon/reed_solomon.ko
kernel/fs/pstore/pstore_zone.ko:
kernel/fs/pstore/pstore_blk.ko: kernel/fs/pstore/pstore_zone.ko
kernel/fs/binfmt_misc.ko:
kernel/fs/dlm/dlm.ko:
kernel/fs/netfs/netfs.ko:
kernel/fs/fscache/fscache.ko: kernel/fs/netfs/netfs.ko
kernel/fs/reiserfs/reiserfs.ko:
kernel/fs/cramfs/cramfs.ko: kernel/drivers/mtd/mtd.ko
kernel/fs/coda/coda.ko:
kernel/fs/minix/minix.ko:
kernel/fs/exfat/exfat.ko:
kernel/fs/bfs/bfs.ko:
kernel/fs/isofs/isofs.ko:
kernel/fs/hfsplus/hfsplus.ko:
kernel/fs/hfs/hfs.ko:
kernel/fs/freevxfs/freevxfs.ko:
kernel/fs/nfs/nfs.ko: kernel/fs/lockd/lockd.ko kernel/fs/nfs_common/grace.ko kernel/net/sunrpc/sunrpc.ko kernel/fs/fscache/fscache.ko kernel/fs/netfs/netfs.ko
kernel/fs/nfs/nfsv2.ko: kernel/fs/nfs/nfs.ko kernel/fs/lockd/lockd.ko kernel/fs/nfs_common/grace.ko kernel/net/sunrpc/sunrpc.ko kernel/fs/fscache/fscache.ko kernel/fs/netfs/netfs.ko
kernel/fs/nfs/nfsv3.ko: kernel/fs/nfs_common/nfs_acl.ko kernel/fs/nfs/nfs.ko kernel/fs/lockd/lockd.ko kernel/fs/nfs_common/grace.ko kernel/net/sunrpc/sunrpc.ko kernel/fs/fscache/fscache.ko kernel/fs/netfs/netfs.ko
kernel/fs/nfs/nfsv4.ko: kernel/fs/nfs/nfs.ko kernel/fs/lockd/lockd.ko kernel/fs/nfs_common/grace.ko kernel/net/sunrpc/sunrpc.ko kernel/fs/fscache/fscache.ko kernel/fs/netfs/netfs.ko
kernel/fs/nfs/filelayout/nfs_layout_nfsv41_files.ko: kernel/fs/nfs/nfsv4.ko kernel/fs/nfs/nfs.ko kernel/fs/lockd/lockd.ko kernel/fs/nfs_common/grace.ko kernel/net/sunrpc/sunrpc.ko kernel/fs/fscache/fscache.ko kernel/fs/netfs/netfs.ko
kernel/fs/nfs/blocklayout/blocklayoutdriver.ko: kernel/fs/nfs/nfsv4.ko kernel/fs/nfs/nfs.ko kernel/fs/lockd/lockd.ko kernel/fs/nfs_common/grace.ko kernel/net/sunrpc/sunrpc.ko kernel/fs/fscache/fscache.ko kernel/fs/netfs/netfs.ko
kernel/fs/nfs/flexfilelayout/nfs_layout_flexfiles.ko: kernel/fs/nfs/nfsv4.ko kernel/fs/nfs/nfs.ko kernel/fs/lockd/lockd.ko kernel/fs/nfs_common/grace.ko kernel/net/sunrpc/sunrpc.ko kernel/fs/fscache/fscache.ko kernel/fs/netfs/netfs.ko
kernel/fs/nfsd/nfsd.ko: kernel/net/sunrpc/auth_gss/auth_rpcgss.ko kernel/fs/nfs_common/nfs_acl.ko kernel/fs/lockd/lockd.ko kernel/fs/nfs_common/grace.ko kernel/net/sunrpc/sunrpc.ko
kernel/fs/lockd/lockd.ko: kernel/fs/nfs_common/grace.ko kernel/net/sunrpc/sunrpc.ko
kernel/fs/sysv/sysv.ko:
kernel/fs/smbfs_common/cifs_arc4.ko:
kernel/fs/smbfs_common/cifs_md4.ko:
kernel/fs/cifs/cifs.ko: kernel/fs/smbfs_common/cifs_arc4.ko kernel/fs/smbfs_common/cifs_md4.ko kernel/fs/fscache/fscache.ko kernel/fs/netfs/netfs.ko
kernel/fs/ksmbd/ksmbd.ko: kernel/drivers/infiniband/core/rdma_cm.ko kernel/drivers/infiniband/core/iw_cm.ko kernel/drivers/infiniband/core/ib_cm.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/fs/hpfs/hpfs.ko:
kernel/fs/ntfs/ntfs.ko:
kernel/fs/ntfs3/ntfs3.ko:
kernel/fs/ufs/ufs.ko:
kernel/fs/efs/efs.ko:
kernel/fs/jffs2/jffs2.ko: kernel/drivers/mtd/mtd.ko
kernel/fs/ubifs/ubifs.ko: kernel/drivers/mtd/ubi/ubi.ko kernel/drivers/mtd/mtd.ko
kernel/fs/affs/affs.ko:
kernel/fs/romfs/romfs.ko:
kernel/fs/qnx4/qnx4.ko:
kernel/fs/qnx6/qnx6.ko:
kernel/fs/autofs/autofs4.ko:
kernel/fs/adfs/adfs.ko:
kernel/fs/overlayfs/overlay.ko:
kernel/fs/orangefs/orangefs.ko:
kernel/fs/udf/udf.ko: kernel/lib/crc-itu-t.ko
kernel/fs/omfs/omfs.ko: kernel/lib/crc-itu-t.ko
kernel/fs/jfs/jfs.ko:
kernel/fs/xfs/xfs.ko: kernel/lib/libcrc32c.ko
kernel/fs/9p/9p.ko: kernel/net/9p/9pnet.ko kernel/fs/fscache/fscache.ko kernel/fs/netfs/netfs.ko
kernel/fs/afs/kafs.ko: kernel/net/rxrpc/rxrpc.ko kernel/net/ipv6/ip6_udp_tunnel.ko kernel/net/ipv4/udp_tunnel.ko kernel/fs/fscache/fscache.ko kernel/fs/netfs/netfs.ko
kernel/fs/nilfs2/nilfs2.ko:
kernel/fs/befs/befs.ko:
kernel/fs/cachefiles/cachefiles.ko: kernel/fs/fscache/fscache.ko kernel/fs/netfs/netfs.ko
kernel/fs/ocfs2/ocfs2.ko: kernel/fs/ocfs2/cluster/ocfs2_nodemanager.ko kernel/fs/ocfs2/ocfs2_stackglue.ko kernel/fs/quota/quota_tree.ko
kernel/fs/ocfs2/ocfs2_stackglue.ko:
kernel/fs/ocfs2/ocfs2_stack_o2cb.ko: kernel/fs/ocfs2/dlm/ocfs2_dlm.ko kernel/fs/ocfs2/cluster/ocfs2_nodemanager.ko kernel/fs/ocfs2/ocfs2_stackglue.ko
kernel/fs/ocfs2/ocfs2_stack_user.ko: kernel/fs/dlm/dlm.ko kernel/fs/ocfs2/ocfs2_stackglue.ko
kernel/fs/ocfs2/dlmfs/ocfs2_dlmfs.ko: kernel/fs/ocfs2/cluster/ocfs2_nodemanager.ko kernel/fs/ocfs2/ocfs2_stackglue.ko
kernel/fs/ocfs2/cluster/ocfs2_nodemanager.ko:
kernel/fs/ocfs2/dlm/ocfs2_dlm.ko: kernel/fs/ocfs2/cluster/ocfs2_nodemanager.ko
kernel/fs/btrfs/btrfs.ko: kernel/crypto/xor.ko kernel/lib/zstd/zstd_compress.ko kernel/lib/raid6/raid6_pq.ko kernel/lib/libcrc32c.ko
kernel/fs/gfs2/gfs2.ko: kernel/fs/dlm/dlm.ko kernel/lib/libcrc32c.ko
kernel/fs/f2fs/f2fs.ko: kernel/lib/lz4/lz4hc_compress.ko kernel/lib/lz4/lz4_compress.ko kernel/lib/zstd/zstd_compress.ko
kernel/fs/ceph/ceph.ko: kernel/net/ceph/libceph.ko kernel/lib/libcrc32c.ko kernel/fs/fscache/fscache.ko kernel/fs/netfs/netfs.ko
kernel/fs/erofs/erofs.ko: kernel/lib/libcrc32c.ko
kernel/fs/vboxsf/vboxsf.ko: kernel/drivers/virt/vboxguest/vboxguest.ko
kernel/fs/zonefs/zonefs.ko:
kernel/fs/shiftfs.ko:
kernel/crypto/asymmetric_keys/asym_tpm.ko:
kernel/crypto/asymmetric_keys/pkcs8_key_parser.ko:
kernel/crypto/asymmetric_keys/pkcs7_test_key.ko:
kernel/crypto/asymmetric_keys/tpm_key_parser.ko: kernel/crypto/asymmetric_keys/asym_tpm.ko
kernel/crypto/crypto_engine.ko:
kernel/crypto/echainiv.ko:
kernel/crypto/sm2_generic.ko: kernel/crypto/sm3_generic.ko
kernel/crypto/ecdsa_generic.ko: kernel/crypto/ecc.ko
kernel/crypto/crypto_user.ko:
kernel/crypto/cmac.ko:
kernel/crypto/vmac.ko:
kernel/crypto/xcbc.ko:
kernel/crypto/md4.ko:
kernel/crypto/rmd160.ko:
kernel/crypto/sha3_generic.ko:
kernel/crypto/sm3_generic.ko:
kernel/crypto/streebog_generic.ko:
kernel/crypto/wp512.ko:
kernel/crypto/blake2b_generic.ko:
kernel/crypto/blake2s_generic.ko:
kernel/crypto/cfb.ko:
kernel/crypto/pcbc.ko:
kernel/crypto/lrw.ko:
kernel/crypto/keywrap.ko:
kernel/crypto/adiantum.ko: kernel/lib/crypto/libpoly1305.ko
kernel/crypto/nhpoly1305.ko: kernel/lib/crypto/libpoly1305.ko
kernel/crypto/ccm.ko:
kernel/crypto/chacha20poly1305.ko:
kernel/crypto/aegis128.ko:
kernel/crypto/pcrypt.ko:
kernel/crypto/cryptd.ko:
kernel/crypto/des_generic.ko: kernel/lib/crypto/libdes.ko
kernel/crypto/fcrypt.ko:
kernel/crypto/blowfish_generic.ko: kernel/crypto/blowfish_common.ko
kernel/crypto/blowfish_common.ko:
kernel/crypto/twofish_generic.ko: kernel/crypto/twofish_common.ko
kernel/crypto/twofish_common.ko:
kernel/crypto/serpent_generic.ko:
kernel/crypto/sm4_generic.ko: kernel/lib/crypto/libsm4.ko
kernel/crypto/aes_ti.ko:
kernel/crypto/camellia_generic.ko:
kernel/crypto/cast_common.ko:
kernel/crypto/cast5_generic.ko: kernel/crypto/cast_common.ko
kernel/crypto/cast6_generic.ko: kernel/crypto/cast_common.ko
kernel/crypto/chacha_generic.ko: kernel/lib/crypto/libchacha.ko
kernel/crypto/poly1305_generic.ko: kernel/lib/crypto/libpoly1305.ko
kernel/crypto/michael_mic.ko:
kernel/crypto/crc32_generic.ko:
kernel/crypto/authenc.ko:
kernel/crypto/authencesn.ko: kernel/crypto/authenc.ko
kernel/crypto/lz4.ko: kernel/lib/lz4/lz4_compress.ko
kernel/crypto/lz4hc.ko: kernel/lib/lz4/lz4hc_compress.ko
kernel/crypto/xxhash_generic.ko:
kernel/crypto/842.ko: kernel/lib/842/842_decompress.ko kernel/lib/842/842_compress.ko
kernel/crypto/ansi_cprng.ko:
kernel/crypto/tcrypt.ko:
kernel/crypto/af_alg.ko:
kernel/crypto/algif_hash.ko: kernel/crypto/af_alg.ko
kernel/crypto/algif_skcipher.ko: kernel/crypto/af_alg.ko
kernel/crypto/algif_rng.ko: kernel/crypto/af_alg.ko
kernel/crypto/algif_aead.ko: kernel/crypto/af_alg.ko
kernel/crypto/zstd.ko: kernel/lib/zstd/zstd_compress.ko
kernel/crypto/ofb.ko:
kernel/crypto/ecc.ko:
kernel/crypto/essiv.ko: kernel/crypto/authenc.ko
kernel/crypto/curve25519-generic.ko: kernel/lib/crypto/libcurve25519-generic.ko
kernel/crypto/ecdh_generic.ko: kernel/crypto/ecc.ko
kernel/crypto/ecrdsa_generic.ko: kernel/crypto/ecc.ko
kernel/crypto/xor.ko:
kernel/crypto/async_tx/async_tx.ko:
kernel/crypto/async_tx/async_memcpy.ko: kernel/crypto/async_tx/async_tx.ko
kernel/crypto/async_tx/async_xor.ko: kernel/crypto/async_tx/async_tx.ko kernel/crypto/xor.ko
kernel/crypto/async_tx/async_pq.ko: kernel/crypto/async_tx/async_xor.ko kernel/crypto/async_tx/async_tx.ko kernel/crypto/xor.ko kernel/lib/raid6/raid6_pq.ko
kernel/crypto/async_tx/async_raid6_recov.ko: kernel/crypto/async_tx/async_memcpy.ko kernel/crypto/async_tx/async_pq.ko kernel/crypto/async_tx/async_xor.ko kernel/crypto/async_tx/async_tx.ko kernel/crypto/xor.ko kernel/lib/raid6/raid6_pq.ko
kernel/crypto/crypto_simd.ko: kernel/crypto/cryptd.ko
kernel/block/kyber-iosched.ko:
kernel/block/bfq.ko:
kernel/lib/math/cordic.ko:
kernel/lib/crypto/libchacha.ko:
kernel/lib/crypto/libarc4.ko:
kernel/lib/crypto/libchacha20poly1305.ko: kernel/arch/x86/crypto/chacha-x86_64.ko kernel/arch/x86/crypto/poly1305-x86_64.ko kernel/lib/crypto/libchacha.ko
kernel/lib/crypto/libcurve25519-generic.ko:
kernel/lib/crypto/libcurve25519.ko:
kernel/lib/crypto/libdes.ko:
kernel/lib/crypto/libpoly1305.ko:
kernel/lib/crypto/libsm4.ko:
kernel/lib/lz4/lz4_compress.ko:
kernel/lib/lz4/lz4hc_compress.ko:
kernel/lib/zstd/zstd_compress.ko:
kernel/lib/xz/xz_dec_test.ko:
kernel/lib/test_bpf.ko:
kernel/lib/test_blackhole_dev.ko:
kernel/lib/crc-itu-t.ko:
kernel/lib/crc64.ko:
kernel/lib/crc4.ko:
kernel/lib/crc7.ko:
kernel/lib/libcrc32c.ko:
kernel/lib/crc8.ko:
kernel/lib/842/842_compress.ko:
kernel/lib/842/842_decompress.ko:
kernel/lib/reed_solomon/reed_solomon.ko:
kernel/lib/bch.ko:
kernel/lib/raid6/raid6_pq.ko:
kernel/lib/ts_kmp.ko:
kernel/lib/ts_bm.ko:
kernel/lib/ts_fsm.ko:
kernel/lib/notifier-error-inject.ko:
kernel/lib/pm-notifier-error-inject.ko: kernel/lib/notifier-error-inject.ko
kernel/lib/memory-notifier-error-inject.ko: kernel/lib/notifier-error-inject.ko
kernel/lib/lru_cache.ko:
kernel/lib/parman.ko:
kernel/lib/objagg.ko:
kernel/drivers/irqchip/irq-madera.ko:
kernel/drivers/bus/mhi/core/mhi.ko:
kernel/drivers/bus/mhi/mhi_pci_generic.ko: kernel/drivers/bus/mhi/core/mhi.ko
kernel/drivers/phy/broadcom/phy-bcm-kona-usb2.ko:
kernel/drivers/phy/intel/phy-intel-lgm-emmc.ko:
kernel/drivers/phy/marvell/phy-pxa-28nm-hsic.ko:
kernel/drivers/phy/marvell/phy-pxa-28nm-usb2.ko:
kernel/drivers/phy/motorola/phy-cpcap-usb.ko: kernel/drivers/usb/musb/musb_hdrc.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/phy/qualcomm/phy-qcom-usb-hs.ko: kernel/drivers/usb/common/ulpi.ko
kernel/drivers/phy/qualcomm/phy-qcom-usb-hsic.ko: kernel/drivers/usb/common/ulpi.ko
kernel/drivers/phy/samsung/phy-exynos-usb2.ko:
kernel/drivers/phy/ti/phy-tusb1210.ko: kernel/drivers/usb/common/ulpi.ko
kernel/drivers/phy/phy-can-transceiver.ko:
kernel/drivers/phy/phy-lgm-usb.ko:
kernel/drivers/pinctrl/intel/pinctrl-lynxpoint.ko:
kernel/drivers/pinctrl/intel/pinctrl-alderlake.ko:
kernel/drivers/pinctrl/intel/pinctrl-broxton.ko:
kernel/drivers/pinctrl/intel/pinctrl-cannonlake.ko:
kernel/drivers/pinctrl/intel/pinctrl-cedarfork.ko:
kernel/drivers/pinctrl/intel/pinctrl-denverton.ko:
kernel/drivers/pinctrl/intel/pinctrl-elkhartlake.ko:
kernel/drivers/pinctrl/intel/pinctrl-emmitsburg.ko:
kernel/drivers/pinctrl/intel/pinctrl-geminilake.ko:
kernel/drivers/pinctrl/intel/pinctrl-icelake.ko:
kernel/drivers/pinctrl/intel/pinctrl-jasperlake.ko:
kernel/drivers/pinctrl/intel/pinctrl-lakefield.ko:
kernel/drivers/pinctrl/intel/pinctrl-lewisburg.ko:
kernel/drivers/pinctrl/intel/pinctrl-sunrisepoint.ko:
kernel/drivers/pinctrl/intel/pinctrl-tigerlake.ko:
kernel/drivers/pinctrl/cirrus/pinctrl-madera.ko:
kernel/drivers/pinctrl/pinctrl-da9062.ko:
kernel/drivers/pinctrl/pinctrl-mcp23s08_i2c.ko: kernel/drivers/pinctrl/pinctrl-mcp23s08.ko
kernel/drivers/pinctrl/pinctrl-mcp23s08_spi.ko: kernel/drivers/pinctrl/pinctrl-mcp23s08.ko
kernel/drivers/pinctrl/pinctrl-mcp23s08.ko:
kernel/drivers/gpio/gpio-generic.ko:
kernel/drivers/gpio/gpio-104-dio-48e.ko:
kernel/drivers/gpio/gpio-104-idi-48.ko:
kernel/drivers/gpio/gpio-104-idio-16.ko:
kernel/drivers/gpio/gpio-aaeon.ko: kernel/drivers/platform/x86/asus-wmi.ko kernel/drivers/input/sparse-keymap.ko kernel/drivers/acpi/video.ko kernel/drivers/acpi/platform_profile.ko kernel/drivers/platform/x86/wmi.ko
kernel/drivers/gpio/gpio-adp5520.ko:
kernel/drivers/gpio/gpio-adp5588.ko:
kernel/drivers/gpio/gpio-aggregator.ko:
kernel/drivers/gpio/gpio-amd8111.ko:
kernel/drivers/gpio/gpio-amd-fch.ko:
kernel/drivers/gpio/gpio-amdpt.ko: kernel/drivers/gpio/gpio-generic.ko
kernel/drivers/gpio/gpio-arizona.ko:
kernel/drivers/gpio/gpio-bd9571mwv.ko:
kernel/drivers/gpio/gpio-da9052.ko:
kernel/drivers/gpio/gpio-da9055.ko:
kernel/drivers/gpio/gpio-dln2.ko: kernel/drivers/mfd/dln2.ko
kernel/drivers/gpio/gpio-dwapb.ko: kernel/drivers/gpio/gpio-generic.ko
kernel/drivers/gpio/gpio-exar.ko:
kernel/drivers/gpio/gpio-f7188x.ko:
kernel/drivers/gpio/gpio-gpio-mm.ko:
kernel/drivers/gpio/gpio-ich.ko:
kernel/drivers/gpio/gpio-it87.ko:
kernel/drivers/gpio/gpio-janz-ttl.ko:
kernel/drivers/gpio/gpio-kempld.ko: kernel/drivers/mfd/kempld-core.ko
kernel/drivers/gpio/gpio-ljca.ko: kernel/drivers/mfd/ljca.ko
kernel/drivers/gpio/gpio-lp3943.ko: kernel/drivers/mfd/lp3943.ko
kernel/drivers/gpio/gpio-lp873x.ko:
kernel/drivers/gpio/gpio-madera.ko:
kernel/drivers/gpio/gpio-max3191x.ko: kernel/lib/crc8.ko
kernel/drivers/gpio/gpio-max7300.ko: kernel/drivers/gpio/gpio-max730x.ko
kernel/drivers/gpio/gpio-max7301.ko: kernel/drivers/gpio/gpio-max730x.ko
kernel/drivers/gpio/gpio-max730x.ko:
kernel/drivers/gpio/gpio-max732x.ko:
kernel/drivers/gpio/gpio-mb86s7x.ko:
kernel/drivers/gpio/gpio-mc33880.ko:
kernel/drivers/gpio/gpio-menz127.ko: kernel/drivers/mcb/mcb.ko kernel/drivers/gpio/gpio-generic.ko
kernel/drivers/gpio/gpio-ml-ioh.ko:
kernel/drivers/gpio/gpio-pca953x.ko:
kernel/drivers/gpio/gpio-pca9570.ko:
kernel/drivers/gpio/gpio-pcf857x.ko:
kernel/drivers/gpio/gpio-pcie-idio-24.ko:
kernel/drivers/gpio/gpio-pci-idio-16.ko:
kernel/drivers/gpio/gpio-pisosr.ko:
kernel/drivers/gpio/gpio-rdc321x.ko:
kernel/drivers/gpio/gpio-sch311x.ko:
kernel/drivers/gpio/gpio-sch.ko:
kernel/drivers/gpio/gpio-siox.ko: kernel/drivers/siox/siox-core.ko
kernel/drivers/gpio/gpio-tpic2810.ko:
kernel/drivers/gpio/gpio-tps65086.ko:
kernel/drivers/gpio/gpio-tps65912.ko:
kernel/drivers/gpio/gpio-tqmx86.ko:
kernel/drivers/gpio/gpio-twl4030.ko:
kernel/drivers/gpio/gpio-twl6040.ko:
kernel/drivers/gpio/gpio-ucb1400.ko:
kernel/drivers/gpio/gpio-viperboard.ko:
kernel/drivers/gpio/gpio-virtio.ko:
kernel/drivers/gpio/gpio-vx855.ko:
kernel/drivers/gpio/gpio-wcove.ko:
kernel/drivers/gpio/gpio-winbond.ko:
kernel/drivers/gpio/gpio-wm831x.ko:
kernel/drivers/gpio/gpio-wm8350.ko:
kernel/drivers/gpio/gpio-wm8994.ko:
kernel/drivers/gpio/gpio-ws16c48.ko:
kernel/drivers/gpio/gpio-xra1403.ko:
kernel/drivers/pwm/pwm-cros-ec.ko:
kernel/drivers/pwm/pwm-dwc.ko:
kernel/drivers/pwm/pwm-iqs620a.ko:
kernel/drivers/pwm/pwm-lp3943.ko: kernel/drivers/mfd/lp3943.ko
kernel/drivers/pwm/pwm-pca9685.ko:
kernel/drivers/pwm/pwm-twl.ko:
kernel/drivers/pwm/pwm-twl-led.ko:
kernel/drivers/pci/hotplug/cpcihp_zt5550.ko:
kernel/drivers/pci/hotplug/cpcihp_generic.ko:
kernel/drivers/pci/hotplug/acpiphp_ibm.ko:
kernel/drivers/pci/endpoint/functions/pci-epf-ntb.ko:
kernel/drivers/pci/controller/pci-hyperv.ko: kernel/drivers/pci/controller/pci-hyperv-intf.ko kernel/drivers/hv/hv_vmbus.ko
kernel/drivers/pci/controller/pci-hyperv-intf.ko:
kernel/drivers/pci/controller/vmd.ko:
kernel/drivers/pci/switch/switchtec.ko:
kernel/drivers/pci/pci-stub.ko:
kernel/drivers/pci/pci-pf-stub.ko:
kernel/drivers/pci/xen-pcifront.ko:
kernel/drivers/rapidio/switches/tsi57x.ko:
kernel/drivers/rapidio/switches/idtcps.ko:
kernel/drivers/rapidio/switches/tsi568.ko:
kernel/drivers/rapidio/switches/idt_gen2.ko:
kernel/drivers/rapidio/switches/idt_gen3.ko:
kernel/drivers/rapidio/devices/tsi721_mport.ko:
kernel/drivers/rapidio/devices/rio_mport_cdev.ko:
kernel/drivers/rapidio/rio-scan.ko:
kernel/drivers/rapidio/rio_cm.ko:
kernel/drivers/video/backlight/ams369fg06.ko: kernel/drivers/video/backlight/lcd.ko
kernel/drivers/video/backlight/lcd.ko:
kernel/drivers/video/backlight/hx8357.ko:
kernel/drivers/video/backlight/ili922x.ko: kernel/drivers/video/backlight/lcd.ko
kernel/drivers/video/backlight/ili9320.ko: kernel/drivers/video/backlight/lcd.ko
kernel/drivers/video/backlight/l4f00242t03.ko: kernel/drivers/video/backlight/lcd.ko
kernel/drivers/video/backlight/lms283gf05.ko: kernel/drivers/video/backlight/lcd.ko
kernel/drivers/video/backlight/lms501kf03.ko: kernel/drivers/video/backlight/lcd.ko
kernel/drivers/video/backlight/ltv350qv.ko: kernel/drivers/video/backlight/lcd.ko
kernel/drivers/video/backlight/otm3225a.ko: kernel/drivers/video/backlight/lcd.ko
kernel/drivers/video/backlight/platform_lcd.ko: kernel/drivers/video/backlight/lcd.ko
kernel/drivers/video/backlight/tdo24m.ko: kernel/drivers/video/backlight/lcd.ko
kernel/drivers/video/backlight/vgg2432a4.ko: kernel/drivers/video/backlight/ili9320.ko kernel/drivers/video/backlight/lcd.ko
kernel/drivers/video/backlight/88pm860x_bl.ko:
kernel/drivers/video/backlight/aat2870_bl.ko:
kernel/drivers/video/backlight/adp5520_bl.ko:
kernel/drivers/video/backlight/adp8860_bl.ko:
kernel/drivers/video/backlight/adp8870_bl.ko:
kernel/drivers/video/backlight/apple_bl.ko:
kernel/drivers/video/backlight/as3711_bl.ko:
kernel/drivers/video/backlight/bd6107.ko:
kernel/drivers/video/backlight/cr_bllcd.ko: kernel/drivers/video/backlight/lcd.ko
kernel/drivers/video/backlight/da903x_bl.ko:
kernel/drivers/video/backlight/da9052_bl.ko:
kernel/drivers/video/backlight/gpio_backlight.ko:
kernel/drivers/video/backlight/ktd253-backlight.ko:
kernel/drivers/video/backlight/lm3533_bl.ko: kernel/drivers/mfd/lm3533-ctrlbank.ko kernel/drivers/mfd/lm3533-core.ko
kernel/drivers/video/backlight/lm3630a_bl.ko:
kernel/drivers/video/backlight/lm3639_bl.ko:
kernel/drivers/video/backlight/lp855x_bl.ko:
kernel/drivers/video/backlight/lp8788_bl.ko:
kernel/drivers/video/backlight/lv5207lp.ko:
kernel/drivers/video/backlight/max8925_bl.ko:
kernel/drivers/video/backlight/pandora_bl.ko:
kernel/drivers/video/backlight/pcf50633-backlight.ko: kernel/drivers/mfd/pcf50633.ko
kernel/drivers/video/backlight/pwm_bl.ko:
kernel/drivers/video/backlight/qcom-wled.ko:
kernel/drivers/video/backlight/rt4831-backlight.ko:
kernel/drivers/video/backlight/kb3886_bl.ko:
kernel/drivers/video/backlight/sky81452-backlight.ko:
kernel/drivers/video/backlight/wm831x_bl.ko:
kernel/drivers/video/backlight/arcxcnn_bl.ko:
kernel/drivers/video/backlight/rave-sp-backlight.ko: kernel/drivers/mfd/rave-sp.ko
kernel/drivers/video/fbdev/core/sysfillrect.ko:
kernel/drivers/video/fbdev/core/syscopyarea.ko:
kernel/drivers/video/fbdev/core/sysimgblt.ko:
kernel/drivers/video/fbdev/core/fb_sys_fops.ko:
kernel/drivers/video/fbdev/core/svgalib.ko:
kernel/drivers/video/fbdev/core/fb_ddc.ko:
kernel/drivers/video/fbdev/arcfb.ko: kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/video/fbdev/cyber2000fb.ko: kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/video/fbdev/pm2fb.ko:
kernel/drivers/video/fbdev/pm3fb.ko:
kernel/drivers/video/fbdev/i740fb.ko: kernel/drivers/video/fbdev/core/fb_ddc.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/video/fbdev/matrox/matroxfb_base.ko: kernel/drivers/video/fbdev/matrox/matroxfb_g450.ko kernel/drivers/video/fbdev/matrox/matroxfb_Ti3026.ko kernel/drivers/video/fbdev/matrox/matroxfb_accel.ko kernel/drivers/video/fbdev/matrox/matroxfb_DAC1064.ko kernel/drivers/video/fbdev/matrox/g450_pll.ko kernel/drivers/video/fbdev/matrox/matroxfb_misc.ko
kernel/drivers/video/fbdev/matrox/matroxfb_accel.ko:
kernel/drivers/video/fbdev/matrox/matroxfb_DAC1064.ko: kernel/drivers/video/fbdev/matrox/g450_pll.ko kernel/drivers/video/fbdev/matrox/matroxfb_misc.ko
kernel/drivers/video/fbdev/matrox/matroxfb_Ti3026.ko: kernel/drivers/video/fbdev/matrox/matroxfb_misc.ko
kernel/drivers/video/fbdev/matrox/matroxfb_misc.ko:
kernel/drivers/video/fbdev/matrox/g450_pll.ko: kernel/drivers/video/fbdev/matrox/matroxfb_misc.ko
kernel/drivers/video/fbdev/matrox/matroxfb_g450.ko: kernel/drivers/video/fbdev/matrox/g450_pll.ko kernel/drivers/video/fbdev/matrox/matroxfb_misc.ko
kernel/drivers/video/fbdev/matrox/matroxfb_crtc2.ko: kernel/drivers/video/fbdev/matrox/matroxfb_base.ko kernel/drivers/video/fbdev/matrox/matroxfb_g450.ko kernel/drivers/video/fbdev/matrox/matroxfb_Ti3026.ko kernel/drivers/video/fbdev/matrox/matroxfb_accel.ko kernel/drivers/video/fbdev/matrox/matroxfb_DAC1064.ko kernel/drivers/video/fbdev/matrox/g450_pll.ko kernel/drivers/video/fbdev/matrox/matroxfb_misc.ko
kernel/drivers/video/fbdev/matrox/i2c-matroxfb.ko: kernel/drivers/video/fbdev/matrox/matroxfb_base.ko kernel/drivers/video/fbdev/matrox/matroxfb_g450.ko kernel/drivers/video/fbdev/matrox/matroxfb_Ti3026.ko kernel/drivers/video/fbdev/matrox/matroxfb_accel.ko kernel/drivers/video/fbdev/matrox/matroxfb_DAC1064.ko kernel/drivers/video/fbdev/matrox/g450_pll.ko kernel/drivers/video/fbdev/matrox/matroxfb_misc.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/video/fbdev/matrox/matroxfb_maven.ko: kernel/drivers/video/fbdev/matrox/matroxfb_misc.ko
kernel/drivers/video/fbdev/riva/rivafb.ko: kernel/drivers/video/vgastate.ko kernel/drivers/video/fbdev/core/fb_ddc.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/video/fbdev/nvidia/nvidiafb.ko: kernel/drivers/video/vgastate.ko kernel/drivers/video/fbdev/core/fb_ddc.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/video/fbdev/aty/atyfb.ko:
kernel/drivers/video/fbdev/aty/aty128fb.ko:
kernel/drivers/video/fbdev/aty/radeonfb.ko: kernel/drivers/video/fbdev/core/fb_ddc.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/video/fbdev/macmodes.ko:
kernel/drivers/video/fbdev/sis/sisfb.ko:
kernel/drivers/video/fbdev/via/viafb.ko: kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/video/fbdev/kyro/kyrofb.ko:
kernel/drivers/video/fbdev/savage/savagefb.ko: kernel/drivers/video/vgastate.ko kernel/drivers/video/fbdev/core/fb_ddc.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/video/fbdev/neofb.ko: kernel/drivers/video/vgastate.ko
kernel/drivers/video/fbdev/tdfxfb.ko:
kernel/drivers/video/fbdev/vt8623fb.ko: kernel/drivers/video/fbdev/core/svgalib.ko kernel/drivers/video/vgastate.ko
kernel/drivers/video/fbdev/tridentfb.ko: kernel/drivers/video/fbdev/core/fb_ddc.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/video/fbdev/vermilion/vmlfb.ko:
kernel/drivers/video/fbdev/vermilion/crvml.ko: kernel/drivers/video/fbdev/vermilion/vmlfb.ko
kernel/drivers/video/fbdev/s3fb.ko: kernel/drivers/video/fbdev/core/svgalib.ko kernel/drivers/video/vgastate.ko kernel/drivers/video/fbdev/core/fb_ddc.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/video/fbdev/arkfb.ko: kernel/drivers/video/fbdev/core/svgalib.ko kernel/drivers/video/vgastate.ko
kernel/drivers/video/fbdev/hecubafb.ko: kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/video/fbdev/n411.ko:
kernel/drivers/video/fbdev/hgafb.ko:
kernel/drivers/video/fbdev/sstfb.ko:
kernel/drivers/video/fbdev/cirrusfb.ko:
kernel/drivers/video/fbdev/metronomefb.ko: kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/video/fbdev/s1d13xxxfb.ko:
kernel/drivers/video/fbdev/sm501fb.ko: kernel/drivers/mfd/sm501.ko
kernel/drivers/video/fbdev/udlfb.ko: kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/video/fbdev/smscufx.ko: kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/video/fbdev/xen-fbfront.ko: kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/video/fbdev/carminefb.ko:
kernel/drivers/video/fbdev/mb862xx/mb862xxfb.ko:
kernel/drivers/video/fbdev/hyperv_fb.ko: kernel/drivers/hv/hv_vmbus.ko
kernel/drivers/video/fbdev/ocfb.ko:
kernel/drivers/video/fbdev/sm712fb.ko:
kernel/drivers/video/fbdev/uvesafb.ko:
kernel/drivers/video/fbdev/vga16fb.ko: kernel/drivers/video/vgastate.ko
kernel/drivers/video/fbdev/ssd1307fb.ko: kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/video/fbdev/simplefb.ko:
kernel/drivers/video/vgastate.ko:
kernel/drivers/char/ipmi/ipmi_msghandler.ko:
kernel/drivers/char/ipmi/ipmi_devintf.ko: kernel/drivers/char/ipmi/ipmi_msghandler.ko
kernel/drivers/char/ipmi/ipmi_si.ko: kernel/drivers/char/ipmi/ipmi_msghandler.ko
kernel/drivers/char/ipmi/ipmi_ssif.ko: kernel/drivers/char/ipmi/ipmi_msghandler.ko
kernel/drivers/char/ipmi/ipmi_watchdog.ko: kernel/drivers/char/ipmi/ipmi_msghandler.ko
kernel/drivers/char/ipmi/ipmi_poweroff.ko: kernel/drivers/char/ipmi/ipmi_msghandler.ko
kernel/drivers/acpi/apei/einj.ko:
kernel/drivers/acpi/dptf/dptf_power.ko:
kernel/drivers/acpi/dptf/dptf_pch_fivr.ko:
kernel/drivers/acpi/acpi_ipmi.ko: kernel/drivers/char/ipmi/ipmi_msghandler.ko
kernel/drivers/acpi/video.ko:
kernel/drivers/acpi/acpi_tad.ko:
kernel/drivers/acpi/platform_profile.ko:
kernel/drivers/acpi/nfit/nfit.ko:
kernel/drivers/acpi/sbshc.ko:
kernel/drivers/acpi/sbs.ko: kernel/drivers/acpi/sbshc.ko
kernel/drivers/acpi/ec_sys.ko:
kernel/drivers/acpi/acpi_pad.ko:
kernel/drivers/acpi/acpi_extlog.ko:
kernel/drivers/acpi/acpi_configfs.ko:
kernel/drivers/clk/xilinx/xlnx_vcu.ko:
kernel/drivers/clk/clk-cdce706.ko:
kernel/drivers/clk/clk-cs2000-cp.ko:
kernel/drivers/clk/clk-lmk04832.ko:
kernel/drivers/clk/clk-max9485.ko:
kernel/drivers/clk/clk-palmas.ko:
kernel/drivers/clk/clk-pwm.ko:
kernel/drivers/clk/clk-si5341.ko:
kernel/drivers/clk/clk-si5351.ko:
kernel/drivers/clk/clk-si544.ko:
kernel/drivers/clk/clk-twl6040.ko:
kernel/drivers/clk/clk-wm831x.ko:
kernel/drivers/dma/idxd/idxd.ko: kernel/drivers/dma/idxd/idxd_bus.ko
kernel/drivers/dma/idxd/idxd_bus.ko:
kernel/drivers/dma/qcom/hdma_mgmt.ko:
kernel/drivers/dma/qcom/hdma.ko:
kernel/drivers/dma/altera-msgdma.ko:
kernel/drivers/dma/ptdma/ptdma.ko:
kernel/drivers/dma/dw/dw_dmac_core.ko:
kernel/drivers/dma/dw/dw_dmac.ko: kernel/drivers/dma/dw/dw_dmac_core.ko
kernel/drivers/dma/dw/dw_dmac_pci.ko: kernel/drivers/dma/dw/dw_dmac_core.ko
kernel/drivers/dma/dw-edma/dw-edma.ko:
kernel/drivers/dma/dw-edma/dw-edma-pcie.ko: kernel/drivers/dma/dw-edma/dw-edma.ko
kernel/drivers/dma/idma64.ko:
kernel/drivers/dma/ioat/ioatdma.ko: kernel/drivers/dca/dca.ko
kernel/drivers/dma/plx_dma.ko:
kernel/drivers/dma/sf-pdma/sf-pdma.ko:
kernel/drivers/soc/qcom/qmi_helpers.ko:
kernel/drivers/virtio/virtio_input.ko:
kernel/drivers/virtio/virtio_vdpa.ko: kernel/drivers/vdpa/vdpa.ko
kernel/drivers/virtio/virtio_mem.ko:
kernel/drivers/virtio/virtio_dma_buf.ko:
kernel/drivers/xen/xen-evtchn.ko:
kernel/drivers/xen/xen-gntdev.ko:
kernel/drivers/xen/xen-gntalloc.ko:
kernel/drivers/xen/xenfs/xenfs.ko: kernel/drivers/xen/xen-privcmd.ko
kernel/drivers/xen/xen-pciback/xen-pciback.ko:
kernel/drivers/xen/xen-privcmd.ko:
kernel/drivers/xen/xen-scsiback.ko: kernel/drivers/target/target_core_mod.ko
kernel/drivers/xen/pvcalls-front.ko:
kernel/drivers/xen/xen-front-pgdir-shbuf.ko:
kernel/drivers/regulator/fixed.ko:
kernel/drivers/regulator/virtual.ko:
kernel/drivers/regulator/userspace-consumer.ko:
kernel/drivers/regulator/88pg86x.ko:
kernel/drivers/regulator/88pm800-regulator.ko:
kernel/drivers/regulator/88pm8607.ko:
kernel/drivers/regulator/aat2870-regulator.ko:
kernel/drivers/regulator/act8865-regulator.ko:
kernel/drivers/regulator/ad5398.ko:
kernel/drivers/regulator/arizona-ldo1.ko:
kernel/drivers/regulator/arizona-micsupp.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/drivers/regulator/as3711-regulator.ko:
kernel/drivers/regulator/atc260x-regulator.ko:
kernel/drivers/regulator/axp20x-regulator.ko:
kernel/drivers/regulator/bcm590xx-regulator.ko:
kernel/drivers/regulator/bd9571mwv-regulator.ko:
kernel/drivers/regulator/da903x-regulator.ko:
kernel/drivers/regulator/da9052-regulator.ko:
kernel/drivers/regulator/da9055-regulator.ko:
kernel/drivers/regulator/da9062-regulator.ko:
kernel/drivers/regulator/da9210-regulator.ko:
kernel/drivers/regulator/da9211-regulator.ko:
kernel/drivers/regulator/fan53555.ko:
kernel/drivers/regulator/gpio-regulator.ko:
kernel/drivers/regulator/isl6271a-regulator.ko:
kernel/drivers/regulator/isl9305.ko:
kernel/drivers/regulator/lm363x-regulator.ko:
kernel/drivers/regulator/lp3971.ko:
kernel/drivers/regulator/lp3972.ko:
kernel/drivers/regulator/lp872x.ko:
kernel/drivers/regulator/lp8788-buck.ko:
kernel/drivers/regulator/lp8788-ldo.ko:
kernel/drivers/regulator/lp8755.ko:
kernel/drivers/regulator/ltc3589.ko:
kernel/drivers/regulator/ltc3676.ko:
kernel/drivers/regulator/max14577-regulator.ko:
kernel/drivers/regulator/max1586.ko:
kernel/drivers/regulator/max8649.ko:
kernel/drivers/regulator/max8660.ko:
kernel/drivers/regulator/max8893.ko:
kernel/drivers/regulator/max8907-regulator.ko:
kernel/drivers/regulator/max8925-regulator.ko:
kernel/drivers/regulator/max8952.ko:
kernel/drivers/regulator/max8997-regulator.ko:
kernel/drivers/regulator/max8998.ko:
kernel/drivers/regulator/max77693-regulator.ko:
kernel/drivers/regulator/max77826-regulator.ko:
kernel/drivers/regulator/mc13783-regulator.ko: kernel/drivers/regulator/mc13xxx-regulator-core.ko kernel/drivers/mfd/mc13xxx-core.ko
kernel/drivers/regulator/mc13892-regulator.ko: kernel/drivers/regulator/mc13xxx-regulator-core.ko kernel/drivers/mfd/mc13xxx-core.ko
kernel/drivers/regulator/mc13xxx-regulator-core.ko: kernel/drivers/mfd/mc13xxx-core.ko
kernel/drivers/regulator/mp8859.ko:
kernel/drivers/regulator/mt6311-regulator.ko:
kernel/drivers/regulator/mt6315-regulator.ko: kernel/drivers/base/regmap/regmap-spmi.ko kernel/drivers/spmi/spmi.ko
kernel/drivers/regulator/mt6323-regulator.ko:
kernel/drivers/regulator/mt6358-regulator.ko:
kernel/drivers/regulator/mt6359-regulator.ko:
kernel/drivers/regulator/mt6360-regulator.ko:
kernel/drivers/regulator/mt6397-regulator.ko:
kernel/drivers/regulator/qcom-labibb-regulator.ko:
kernel/drivers/regulator/qcom_spmi-regulator.ko:
kernel/drivers/regulator/qcom_usb_vbus-regulator.ko:
kernel/drivers/regulator/palmas-regulator.ko:
kernel/drivers/regulator/pca9450-regulator.ko:
kernel/drivers/regulator/pv88060-regulator.ko:
kernel/drivers/regulator/pv88080-regulator.ko:
kernel/drivers/regulator/pv88090-regulator.ko:
kernel/drivers/regulator/pwm-regulator.ko:
kernel/drivers/regulator/tps51632-regulator.ko:
kernel/drivers/regulator/pcap-regulator.ko:
kernel/drivers/regulator/pcf50633-regulator.ko:
kernel/drivers/regulator/rpi-panel-attiny-regulator.ko:
kernel/drivers/regulator/rc5t583-regulator.ko:
kernel/drivers/regulator/rt4801-regulator.ko:
kernel/drivers/regulator/rt4831-regulator.ko:
kernel/drivers/regulator/rt5033-regulator.ko:
kernel/drivers/regulator/rt6160-regulator.ko:
kernel/drivers/regulator/rt6245-regulator.ko:
kernel/drivers/regulator/rtmv20-regulator.ko:
kernel/drivers/regulator/rtq2134-regulator.ko:
kernel/drivers/regulator/rtq6752-regulator.ko:
kernel/drivers/regulator/sky81452-regulator.ko:
kernel/drivers/regulator/slg51000-regulator.ko:
kernel/drivers/regulator/tps6105x-regulator.ko:
kernel/drivers/regulator/tps62360-regulator.ko:
kernel/drivers/regulator/tps65023-regulator.ko:
kernel/drivers/regulator/tps6507x-regulator.ko:
kernel/drivers/regulator/tps65086-regulator.ko:
kernel/drivers/regulator/tps65090-regulator.ko:
kernel/drivers/regulator/tps6524x-regulator.ko:
kernel/drivers/regulator/tps6586x-regulator.ko:
kernel/drivers/regulator/tps65910-regulator.ko:
kernel/drivers/regulator/tps65912-regulator.ko:
kernel/drivers/regulator/tps80031-regulator.ko:
kernel/drivers/regulator/tps65132-regulator.ko:
kernel/drivers/regulator/twl-regulator.ko:
kernel/drivers/regulator/twl6030-regulator.ko:
kernel/drivers/regulator/wm831x-dcdc.ko:
kernel/drivers/regulator/wm831x-isink.ko:
kernel/drivers/regulator/wm831x-ldo.ko:
kernel/drivers/regulator/wm8350-regulator.ko:
kernel/drivers/regulator/wm8400-regulator.ko:
kernel/drivers/regulator/wm8994-regulator.ko:
kernel/drivers/reset/reset-ti-syscon.ko:
kernel/drivers/tty/serial/8250/8250_exar.ko:
kernel/drivers/tty/serial/8250/serial_cs.ko: kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/tty/serial/8250/8250_men_mcb.ko: kernel/drivers/mcb/mcb.ko
kernel/drivers/tty/serial/8250/8250_dw.ko:
kernel/drivers/tty/serial/8250/8250_lpss.ko: kernel/drivers/dma/dw/dw_dmac_core.ko
kernel/drivers/tty/serial/bcm63xx_uart.ko:
kernel/drivers/tty/serial/max3100.ko:
kernel/drivers/tty/serial/sc16is7xx.ko:
kernel/drivers/tty/serial/jsm/jsm.ko:
kernel/drivers/tty/serial/uartlite.ko:
kernel/drivers/tty/serial/altera_uart.ko:
kernel/drivers/tty/serial/altera_jtaguart.ko:
kernel/drivers/tty/serial/lantiq.ko:
kernel/drivers/tty/serial/arc_uart.ko:
kernel/drivers/tty/serial/rp2.ko:
kernel/drivers/tty/serial/fsl_lpuart.ko:
kernel/drivers/tty/serial/fsl_linflexuart.ko:
kernel/drivers/tty/serial/men_z135_uart.ko: kernel/drivers/mcb/mcb.ko
kernel/drivers/tty/serial/sprd_serial.ko:
kernel/drivers/tty/ipwireless/ipwireless.ko: kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/tty/n_hdlc.ko:
kernel/drivers/tty/n_gsm.ko:
kernel/drivers/tty/moxa.ko:
kernel/drivers/tty/mxser.ko:
kernel/drivers/tty/nozomi.ko:
kernel/drivers/tty/ttynull.ko:
kernel/drivers/tty/synclink_gt.ko: kernel/drivers/net/wan/hdlc.ko
kernel/drivers/char/hw_random/timeriomem-rng.ko:
kernel/drivers/char/hw_random/intel-rng.ko:
kernel/drivers/char/hw_random/amd-rng.ko:
kernel/drivers/char/hw_random/ba431-rng.ko:
kernel/drivers/char/hw_random/via-rng.ko:
kernel/drivers/char/hw_random/virtio-rng.ko:
kernel/drivers/char/hw_random/xiphera-trng.ko:
kernel/drivers/char/agp/sis-agp.ko:
kernel/drivers/char/tpm/tpm_tis_spi.ko:
kernel/drivers/char/tpm/tpm_tis_i2c_cr50.ko:
kernel/drivers/char/tpm/tpm_i2c_atmel.ko:
kernel/drivers/char/tpm/tpm_i2c_infineon.ko:
kernel/drivers/char/tpm/tpm_i2c_nuvoton.ko:
kernel/drivers/char/tpm/tpm_nsc.ko:
kernel/drivers/char/tpm/tpm_atmel.ko:
kernel/drivers/char/tpm/tpm_infineon.ko:
kernel/drivers/char/tpm/st33zp24/tpm_st33zp24.ko:
kernel/drivers/char/tpm/st33zp24/tpm_st33zp24_i2c.ko: kernel/drivers/char/tpm/st33zp24/tpm_st33zp24.ko
kernel/drivers/char/tpm/st33zp24/tpm_st33zp24_spi.ko: kernel/drivers/char/tpm/st33zp24/tpm_st33zp24.ko
kernel/drivers/char/tpm/xen-tpmfront.ko:
kernel/drivers/char/tpm/tpm_vtpm_proxy.ko:
kernel/drivers/char/uv_mmtimer.ko:
kernel/drivers/char/lp.ko: kernel/drivers/parport/parport.ko
kernel/drivers/char/applicom.ko:
kernel/drivers/char/nvram.ko:
kernel/drivers/char/ppdev.ko: kernel/drivers/parport/parport.ko
kernel/drivers/char/tlclk.ko:
kernel/drivers/char/mwave/mwave.ko:
kernel/drivers/char/pcmcia/synclink_cs.ko: kernel/drivers/net/wan/hdlc.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/char/pcmcia/cm4000_cs.ko: kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/char/pcmcia/cm4040_cs.ko: kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/char/pcmcia/scr24x_cs.ko: kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/char/hangcheck-timer.ko:
kernel/drivers/char/xillybus/xillybus_class.ko:
kernel/drivers/char/xillybus/xillybus_core.ko: kernel/drivers/char/xillybus/xillybus_class.ko
kernel/drivers/char/xillybus/xillybus_pcie.ko: kernel/drivers/char/xillybus/xillybus_core.ko kernel/drivers/char/xillybus/xillybus_class.ko
kernel/drivers/char/xillybus/xillyusb.ko: kernel/drivers/char/xillybus/xillybus_class.ko
kernel/drivers/iommu/amd/iommu_v2.ko:
kernel/drivers/base/regmap/regmap-slimbus.ko: kernel/drivers/slimbus/slimbus.ko
kernel/drivers/base/regmap/regmap-spmi.ko: kernel/drivers/spmi/spmi.ko
kernel/drivers/base/regmap/regmap-w1.ko: kernel/drivers/w1/wire.ko
kernel/drivers/base/regmap/regmap-sdw.ko: kernel/drivers/soundwire/soundwire-bus.ko
kernel/drivers/base/regmap/regmap-sdw-mbq.ko: kernel/drivers/soundwire/soundwire-bus.ko
kernel/drivers/base/regmap/regmap-sccb.ko:
kernel/drivers/base/regmap/regmap-i3c.ko: kernel/drivers/i3c/i3c.ko
kernel/drivers/base/regmap/regmap-spi-avmm.ko:
kernel/drivers/block/rnbd/rnbd-client.ko: kernel/drivers/infiniband/ulp/rtrs/rtrs-client.ko kernel/drivers/infiniband/ulp/rtrs/rtrs-core.ko kernel/drivers/infiniband/core/rdma_cm.ko kernel/drivers/infiniband/core/iw_cm.ko kernel/drivers/infiniband/core/ib_cm.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/block/rnbd/rnbd-server.ko: kernel/drivers/infiniband/ulp/rtrs/rtrs-server.ko kernel/drivers/infiniband/ulp/rtrs/rtrs-core.ko kernel/drivers/infiniband/core/rdma_cm.ko kernel/drivers/infiniband/core/iw_cm.ko kernel/drivers/infiniband/core/ib_cm.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/block/floppy.ko:
kernel/drivers/block/brd.ko:
kernel/drivers/block/pktcdvd.ko:
kernel/drivers/block/nbd.ko:
kernel/drivers/block/cryptoloop.ko:
kernel/drivers/block/virtio_blk.ko:
kernel/drivers/block/sx8.ko:
kernel/drivers/block/xen-blkback/xen-blkback.ko:
kernel/drivers/block/drbd/drbd.ko: kernel/lib/lru_cache.ko kernel/lib/libcrc32c.ko
kernel/drivers/block/rbd.ko: kernel/net/ceph/libceph.ko kernel/lib/libcrc32c.ko
kernel/drivers/block/mtip32xx/mtip32xx.ko:
kernel/drivers/block/rsxx/rsxx.ko:
kernel/drivers/block/zram/zram.ko:
kernel/drivers/block/null_blk/null_blk.ko:
kernel/drivers/misc/eeprom/at24.ko:
kernel/drivers/misc/eeprom/at25.ko:
kernel/drivers/misc/eeprom/eeprom.ko:
kernel/drivers/misc/eeprom/max6875.ko:
kernel/drivers/misc/eeprom/eeprom_93cx6.ko:
kernel/drivers/misc/eeprom/eeprom_93xx46.ko:
kernel/drivers/misc/eeprom/idt_89hpesx.ko:
kernel/drivers/misc/eeprom/ee1004.ko:
kernel/drivers/misc/cb710/cb710.ko:
kernel/drivers/misc/ti-st/st_drv.ko:
kernel/drivers/misc/lis3lv02d/lis3lv02d.ko:
kernel/drivers/misc/lis3lv02d/lis3lv02d_i2c.ko: kernel/drivers/misc/lis3lv02d/lis3lv02d.ko
kernel/drivers/misc/cardreader/alcor_pci.ko:
kernel/drivers/misc/cardreader/rtsx_pci.ko:
kernel/drivers/misc/cardreader/rtsx_usb.ko:
kernel/drivers/misc/pvpanic/pvpanic.ko:
kernel/drivers/misc/pvpanic/pvpanic-mmio.ko: kernel/drivers/misc/pvpanic/pvpanic.ko
kernel/drivers/misc/pvpanic/pvpanic-pci.ko: kernel/drivers/misc/pvpanic/pvpanic.ko
kernel/drivers/misc/ibmasm/ibmasm.ko:
kernel/drivers/misc/ad525x_dpot.ko:
kernel/drivers/misc/ad525x_dpot-i2c.ko: kernel/drivers/misc/ad525x_dpot.ko
kernel/drivers/misc/ad525x_dpot-spi.ko: kernel/drivers/misc/ad525x_dpot.ko
kernel/drivers/misc/dummy-irq.ko:
kernel/drivers/misc/ics932s401.ko:
kernel/drivers/misc/tifm_core.ko:
kernel/drivers/misc/tifm_7xx1.ko: kernel/drivers/misc/tifm_core.ko
kernel/drivers/misc/phantom.ko:
kernel/drivers/misc/bh1770glc.ko:
kernel/drivers/misc/apds990x.ko:
kernel/drivers/misc/enclosure.ko:
kernel/drivers/misc/sgi-xp/xp.ko: kernel/drivers/misc/sgi-gru/gru.ko
kernel/drivers/misc/sgi-xp/xpc.ko: kernel/drivers/misc/sgi-xp/xp.ko kernel/drivers/misc/sgi-gru/gru.ko
kernel/drivers/misc/sgi-xp/xpnet.ko: kernel/drivers/misc/sgi-xp/xp.ko kernel/drivers/misc/sgi-gru/gru.ko
kernel/drivers/misc/sgi-gru/gru.ko:
kernel/drivers/misc/hpilo.ko:
kernel/drivers/misc/apds9802als.ko:
kernel/drivers/misc/isl29003.ko:
kernel/drivers/misc/isl29020.ko:
kernel/drivers/misc/tsl2550.ko:
kernel/drivers/misc/ds1682.ko:
kernel/drivers/misc/c2port/core.ko:
kernel/drivers/misc/c2port/c2port-duramar2150.ko: kernel/drivers/misc/c2port/core.ko
kernel/drivers/misc/hmc6352.ko:
kernel/drivers/misc/vmw_balloon.ko: kernel/drivers/misc/vmw_vmci/vmw_vmci.ko
kernel/drivers/misc/altera-stapl/altera-stapl.ko:
kernel/drivers/misc/mei/mei.ko:
kernel/drivers/misc/mei/mei-me.ko: kernel/drivers/misc/mei/mei.ko
kernel/drivers/misc/mei/mei-txe.ko: kernel/drivers/misc/mei/mei.ko
kernel/drivers/misc/mei/mei-vsc.ko: kernel/drivers/misc/mei/mei.ko
kernel/drivers/misc/mei/hdcp/mei_hdcp.ko: kernel/drivers/misc/mei/mei.ko
kernel/drivers/misc/ivsc/intel_vsc.ko:
kernel/drivers/misc/ivsc/mei_csi.ko: kernel/drivers/misc/ivsc/intel_vsc.ko kernel/drivers/misc/mei/mei.ko
kernel/drivers/misc/ivsc/mei_ace.ko: kernel/drivers/misc/ivsc/intel_vsc.ko kernel/drivers/misc/mei/mei.ko
kernel/drivers/misc/ivsc/mei_pse.ko: kernel/drivers/misc/mei/mei.ko
kernel/drivers/misc/ivsc/mei_ace_debug.ko: kernel/drivers/misc/mei/mei.ko
kernel/drivers/misc/vmw_vmci/vmw_vmci.ko:
kernel/drivers/misc/lattice-ecp3-config.ko:
kernel/drivers/misc/genwqe/genwqe_card.ko: kernel/lib/crc-itu-t.ko
kernel/drivers/misc/echo/echo.ko:
kernel/drivers/misc/dw-xdata-pcie.ko:
kernel/drivers/misc/bcm-vk/bcm_vk.ko:
kernel/drivers/misc/habanalabs/habanalabs.ko:
kernel/drivers/misc/uacce/uacce.ko:
kernel/drivers/misc/xilinx_sdfec.ko:
kernel/drivers/mfd/88pm800.ko: kernel/drivers/mfd/88pm80x.ko
kernel/drivers/mfd/88pm80x.ko:
kernel/drivers/mfd/88pm805.ko: kernel/drivers/mfd/88pm80x.ko
kernel/drivers/mfd/sm501.ko:
kernel/drivers/mfd/bcm590xx.ko:
kernel/drivers/mfd/bd9571mwv.ko:
kernel/drivers/mfd/cros_ec_dev.ko:
kernel/drivers/mfd/htc-pasic3.ko:
kernel/drivers/mfd/lp873x.ko:
kernel/drivers/mfd/ti_am335x_tscadc.ko:
kernel/drivers/mfd/tqmx86.ko:
kernel/drivers/mfd/arizona.ko:
kernel/drivers/mfd/arizona-i2c.ko: kernel/drivers/mfd/arizona.ko
kernel/drivers/mfd/arizona-spi.ko: kernel/drivers/mfd/arizona.ko
kernel/drivers/mfd/wcd934x.ko: kernel/drivers/base/regmap/regmap-slimbus.ko kernel/drivers/slimbus/slimbus.ko
kernel/drivers/mfd/wm8994.ko:
kernel/drivers/mfd/madera.ko:
kernel/drivers/mfd/madera-i2c.ko: kernel/drivers/mfd/madera.ko
kernel/drivers/mfd/madera-spi.ko: kernel/drivers/mfd/madera.ko
kernel/drivers/mfd/tps6105x.ko:
kernel/drivers/mfd/tps65010.ko:
kernel/drivers/mfd/tps6507x.ko:
kernel/drivers/mfd/tps65086.ko:
kernel/drivers/mfd/mc13xxx-core.ko:
kernel/drivers/mfd/mc13xxx-spi.ko: kernel/drivers/mfd/mc13xxx-core.ko
kernel/drivers/mfd/mc13xxx-i2c.ko: kernel/drivers/mfd/mc13xxx-core.ko
kernel/drivers/mfd/ucb1400_core.ko: kernel/sound/ac97_bus.ko
kernel/drivers/mfd/axp20x.ko:
kernel/drivers/mfd/axp20x-i2c.ko: kernel/drivers/mfd/axp20x.ko
kernel/drivers/mfd/lp3943.ko:
kernel/drivers/mfd/ti-lmu.ko:
kernel/drivers/mfd/da9062-core.ko:
kernel/drivers/mfd/da9150-core.ko:
kernel/drivers/mfd/max8907.ko:
kernel/drivers/mfd/mp2629.ko:
kernel/drivers/mfd/pcf50633.ko:
kernel/drivers/mfd/pcf50633-adc.ko: kernel/drivers/mfd/pcf50633.ko
kernel/drivers/mfd/pcf50633-gpio.ko: kernel/drivers/mfd/pcf50633.ko
kernel/drivers/mfd/kempld-core.ko:
kernel/drivers/mfd/intel_quark_i2c_gpio.ko:
kernel/drivers/mfd/lpc_sch.ko:
kernel/drivers/mfd/lpc_ich.ko:
kernel/drivers/mfd/rdc321x-southbridge.ko:
kernel/drivers/mfd/janz-cmodio.ko:
kernel/drivers/mfd/vx855.ko:
kernel/drivers/mfd/wl1273-core.ko:
kernel/drivers/mfd/si476x-core.ko:
kernel/drivers/mfd/intel-lpss.ko:
kernel/drivers/mfd/intel-lpss-pci.ko: kernel/drivers/mfd/intel-lpss.ko
kernel/drivers/mfd/intel-lpss-acpi.ko: kernel/drivers/mfd/intel-lpss.ko
kernel/drivers/mfd/intel_pmc_bxt.ko:
kernel/drivers/mfd/intel_pmt.ko:
kernel/drivers/mfd/viperboard.ko:
kernel/drivers/mfd/lm3533-core.ko:
kernel/drivers/mfd/lm3533-ctrlbank.ko: kernel/drivers/mfd/lm3533-core.ko
kernel/drivers/mfd/retu-mfd.ko:
kernel/drivers/mfd/iqs62x.ko:
kernel/drivers/mfd/menf21bmc.ko:
kernel/drivers/mfd/dln2.ko:
kernel/drivers/mfd/rt4831.ko:
kernel/drivers/mfd/rt5033.ko:
kernel/drivers/mfd/sky81452.ko:
kernel/drivers/mfd/intel_soc_pmic_bxtwc.ko:
kernel/drivers/mfd/intel_soc_pmic_chtdc_ti.ko:
kernel/drivers/mfd/mt6360-core.ko: kernel/lib/crc8.ko
kernel/drivers/mfd/mt6397.ko:
kernel/drivers/mfd/intel_soc_pmic_mrfld.ko:
kernel/drivers/mfd/rave-sp.ko:
kernel/drivers/mfd/intel-m10-bmc.ko: kernel/drivers/base/regmap/regmap-spi-avmm.ko
kernel/drivers/mfd/atc260x-core.ko:
kernel/drivers/mfd/atc260x-i2c.ko: kernel/drivers/mfd/atc260x-core.ko
kernel/drivers/mfd/mfd-aaeon.ko: kernel/drivers/platform/x86/asus-wmi.ko kernel/drivers/input/sparse-keymap.ko kernel/drivers/acpi/video.ko kernel/drivers/acpi/platform_profile.ko kernel/drivers/platform/x86/wmi.ko
kernel/drivers/mfd/ljca.ko:
kernel/drivers/nfc/fdp/fdp.ko: kernel/net/nfc/nci/nci.ko kernel/net/nfc/nfc.ko
kernel/drivers/nfc/fdp/fdp_i2c.ko: kernel/drivers/nfc/fdp/fdp.ko kernel/net/nfc/nci/nci.ko kernel/net/nfc/nfc.ko
kernel/drivers/nfc/pn544/pn544.ko: kernel/net/nfc/hci/hci.ko kernel/net/nfc/nfc.ko
kernel/drivers/nfc/pn544/pn544_i2c.ko: kernel/drivers/nfc/pn544/pn544.ko kernel/net/nfc/hci/hci.ko kernel/net/nfc/nfc.ko
kernel/drivers/nfc/pn544/pn544_mei.ko: kernel/drivers/nfc/mei_phy.ko kernel/drivers/nfc/pn544/pn544.ko kernel/net/nfc/hci/hci.ko kernel/net/nfc/nfc.ko kernel/drivers/misc/mei/mei.ko
kernel/drivers/nfc/microread/microread.ko: kernel/net/nfc/hci/hci.ko kernel/net/nfc/nfc.ko
kernel/drivers/nfc/microread/microread_i2c.ko: kernel/drivers/nfc/microread/microread.ko kernel/net/nfc/hci/hci.ko kernel/net/nfc/nfc.ko
kernel/drivers/nfc/microread/microread_mei.ko: kernel/drivers/nfc/microread/microread.ko kernel/drivers/nfc/mei_phy.ko kernel/net/nfc/hci/hci.ko kernel/net/nfc/nfc.ko kernel/drivers/misc/mei/mei.ko
kernel/drivers/nfc/pn533/pn533.ko: kernel/net/nfc/nfc.ko
kernel/drivers/nfc/pn533/pn533_usb.ko: kernel/drivers/nfc/pn533/pn533.ko kernel/net/nfc/nfc.ko
kernel/drivers/nfc/pn533/pn533_i2c.ko: kernel/drivers/nfc/pn533/pn533.ko kernel/net/nfc/nfc.ko
kernel/drivers/nfc/pn533/pn532_uart.ko: kernel/drivers/nfc/pn533/pn533.ko kernel/net/nfc/nfc.ko
kernel/drivers/nfc/mei_phy.ko: kernel/net/nfc/hci/hci.ko kernel/net/nfc/nfc.ko kernel/drivers/misc/mei/mei.ko
kernel/drivers/nfc/nfcsim.ko: kernel/net/nfc/nfc_digital.ko kernel/net/nfc/nfc.ko kernel/lib/crc-itu-t.ko
kernel/drivers/nfc/port100.ko: kernel/net/nfc/nfc_digital.ko kernel/net/nfc/nfc.ko kernel/lib/crc-itu-t.ko
kernel/drivers/nfc/nfcmrvl/nfcmrvl.ko: kernel/net/nfc/nci/nci.ko kernel/net/nfc/nfc.ko
kernel/drivers/nfc/nfcmrvl/nfcmrvl_usb.ko: kernel/drivers/nfc/nfcmrvl/nfcmrvl.ko kernel/net/nfc/nci/nci.ko kernel/net/nfc/nfc.ko
kernel/drivers/nfc/nfcmrvl/nfcmrvl_uart.ko: kernel/net/nfc/nci/nci_uart.ko kernel/drivers/nfc/nfcmrvl/nfcmrvl.ko kernel/net/nfc/nci/nci.ko kernel/net/nfc/nfc.ko
kernel/drivers/nfc/nfcmrvl/nfcmrvl_i2c.ko: kernel/drivers/nfc/nfcmrvl/nfcmrvl.ko kernel/net/nfc/nci/nci.ko kernel/net/nfc/nfc.ko
kernel/drivers/nfc/nfcmrvl/nfcmrvl_spi.ko: kernel/net/nfc/nci/nci_spi.ko kernel/drivers/nfc/nfcmrvl/nfcmrvl.ko kernel/net/nfc/nci/nci.ko kernel/net/nfc/nfc.ko
kernel/drivers/nfc/trf7970a.ko: kernel/net/nfc/nfc_digital.ko kernel/net/nfc/nfc.ko kernel/lib/crc-itu-t.ko
kernel/drivers/nfc/st21nfca/st21nfca_hci.ko: kernel/net/nfc/hci/hci.ko kernel/net/nfc/nfc.ko
kernel/drivers/nfc/st21nfca/st21nfca_i2c.ko: kernel/drivers/nfc/st21nfca/st21nfca_hci.ko kernel/net/nfc/hci/hci.ko kernel/net/nfc/nfc.ko
kernel/drivers/nfc/st-nci/st-nci.ko: kernel/net/nfc/nci/nci.ko kernel/net/nfc/nfc.ko
kernel/drivers/nfc/st-nci/st-nci_i2c.ko: kernel/drivers/nfc/st-nci/st-nci.ko kernel/net/nfc/nci/nci.ko kernel/net/nfc/nfc.ko
kernel/drivers/nfc/st-nci/st-nci_spi.ko: kernel/drivers/nfc/st-nci/st-nci.ko kernel/net/nfc/nci/nci.ko kernel/net/nfc/nfc.ko
kernel/drivers/nfc/nxp-nci/nxp-nci.ko: kernel/net/nfc/nci/nci.ko kernel/net/nfc/nfc.ko
kernel/drivers/nfc/nxp-nci/nxp-nci_i2c.ko: kernel/drivers/nfc/nxp-nci/nxp-nci.ko kernel/net/nfc/nci/nci.ko kernel/net/nfc/nfc.ko
kernel/drivers/nfc/s3fwrn5/s3fwrn5.ko: kernel/net/nfc/nci/nci.ko kernel/net/nfc/nfc.ko
kernel/drivers/nfc/s3fwrn5/s3fwrn5_i2c.ko: kernel/drivers/nfc/s3fwrn5/s3fwrn5.ko kernel/net/nfc/nci/nci.ko kernel/net/nfc/nfc.ko
kernel/drivers/nfc/s3fwrn5/s3fwrn82_uart.ko: kernel/drivers/nfc/s3fwrn5/s3fwrn5.ko kernel/net/nfc/nci/nci.ko kernel/net/nfc/nfc.ko
kernel/drivers/nfc/st95hf/st95hf.ko: kernel/net/nfc/nfc_digital.ko kernel/net/nfc/nfc.ko kernel/lib/crc-itu-t.ko
kernel/drivers/nfc/virtual_ncidev.ko: kernel/net/nfc/nci/nci.ko kernel/net/nfc/nfc.ko
kernel/drivers/nvdimm/nd_pmem.ko: kernel/drivers/nvdimm/nd_btt.ko
kernel/drivers/nvdimm/nd_btt.ko:
kernel/drivers/nvdimm/nd_blk.ko: kernel/drivers/nvdimm/nd_btt.ko
kernel/drivers/nvdimm/virtio_pmem.ko: kernel/drivers/nvdimm/nd_virtio.ko
kernel/drivers/nvdimm/nd_virtio.ko:
kernel/drivers/dax/pmem/dax_pmem.ko: kernel/drivers/dax/pmem/dax_pmem_core.ko
kernel/drivers/dax/pmem/dax_pmem_core.ko:
kernel/drivers/dax/pmem/dax_pmem_compat.ko: kernel/drivers/dax/device_dax.ko kernel/drivers/dax/pmem/dax_pmem_core.ko
kernel/drivers/dax/hmem/dax_hmem.ko:
kernel/drivers/dax/device_dax.ko:
kernel/drivers/dax/kmem.ko:
kernel/drivers/macintosh/mac_hid.ko:
kernel/drivers/scsi/device_handler/scsi_dh_rdac.ko:
kernel/drivers/scsi/device_handler/scsi_dh_hp_sw.ko:
kernel/drivers/scsi/device_handler/scsi_dh_emc.ko:
kernel/drivers/scsi/device_handler/scsi_dh_alua.ko:
kernel/drivers/scsi/megaraid/megaraid_mm.ko:
kernel/drivers/scsi/megaraid/megaraid_mbox.ko: kernel/drivers/scsi/megaraid/megaraid_mm.ko
kernel/drivers/scsi/megaraid/megaraid_sas.ko:
kernel/drivers/scsi/pcmcia/qlogic_cs.ko: kernel/drivers/scsi/qlogicfas408.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/scsi/pcmcia/fdomain_cs.ko: kernel/drivers/scsi/fdomain.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/scsi/pcmcia/aha152x_cs.ko: kernel/drivers/scsi/scsi_transport_spi.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/scsi/pcmcia/sym53c500_cs.ko: kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/scsi/raid_class.ko:
kernel/drivers/scsi/scsi_transport_spi.ko:
kernel/drivers/scsi/scsi_transport_fc.ko:
kernel/drivers/scsi/scsi_transport_iscsi.ko:
kernel/drivers/scsi/scsi_transport_sas.ko:
kernel/drivers/scsi/libsas/libsas.ko: kernel/drivers/scsi/scsi_transport_sas.ko
kernel/drivers/scsi/scsi_transport_srp.ko:
kernel/drivers/scsi/libfc/libfc.ko: kernel/drivers/scsi/scsi_transport_fc.ko
kernel/drivers/scsi/fcoe/fcoe.ko: kernel/drivers/scsi/fcoe/libfcoe.ko kernel/drivers/scsi/libfc/libfc.ko kernel/drivers/scsi/scsi_transport_fc.ko
kernel/drivers/scsi/fcoe/libfcoe.ko: kernel/drivers/scsi/libfc/libfc.ko kernel/drivers/scsi/scsi_transport_fc.ko
kernel/drivers/scsi/fnic/fnic.ko: kernel/drivers/scsi/fcoe/libfcoe.ko kernel/drivers/scsi/libfc/libfc.ko kernel/drivers/scsi/scsi_transport_fc.ko
kernel/drivers/scsi/snic/snic.ko:
kernel/drivers/scsi/bnx2fc/bnx2fc.ko: kernel/drivers/net/ethernet/broadcom/cnic.ko kernel/drivers/uio/uio.ko kernel/drivers/scsi/fcoe/libfcoe.ko kernel/drivers/scsi/libfc/libfc.ko kernel/drivers/scsi/scsi_transport_fc.ko
kernel/drivers/scsi/qedf/qedf.ko: kernel/drivers/net/ethernet/qlogic/qed/qed.ko kernel/drivers/scsi/fcoe/libfcoe.ko kernel/drivers/scsi/libfc/libfc.ko kernel/drivers/scsi/scsi_transport_fc.ko kernel/lib/crc8.ko
kernel/drivers/scsi/libiscsi.ko: kernel/drivers/scsi/scsi_transport_iscsi.ko
kernel/drivers/scsi/libiscsi_tcp.ko: kernel/drivers/scsi/libiscsi.ko kernel/drivers/scsi/scsi_transport_iscsi.ko
kernel/drivers/scsi/iscsi_tcp.ko: kernel/drivers/scsi/libiscsi_tcp.ko kernel/drivers/scsi/libiscsi.ko kernel/drivers/scsi/scsi_transport_iscsi.ko
kernel/drivers/scsi/iscsi_boot_sysfs.ko:
kernel/drivers/scsi/53c700.ko: kernel/drivers/scsi/scsi_transport_spi.ko
kernel/drivers/scsi/sim710.ko: kernel/drivers/scsi/53c700.ko kernel/drivers/scsi/scsi_transport_spi.ko
kernel/drivers/scsi/advansys.ko:
kernel/drivers/scsi/BusLogic.ko:
kernel/drivers/scsi/dpt_i2o.ko:
kernel/drivers/scsi/arcmsr/arcmsr.ko:
kernel/drivers/scsi/aha1740.ko:
kernel/drivers/scsi/aic7xxx/aic7xxx.ko: kernel/drivers/scsi/scsi_transport_spi.ko
kernel/drivers/scsi/aic7xxx/aic79xx.ko: kernel/drivers/scsi/scsi_transport_spi.ko
kernel/drivers/scsi/aacraid/aacraid.ko:
kernel/drivers/scsi/aic94xx/aic94xx.ko: kernel/drivers/scsi/libsas/libsas.ko kernel/drivers/scsi/scsi_transport_sas.ko
kernel/drivers/scsi/pm8001/pm80xx.ko: kernel/drivers/scsi/libsas/libsas.ko kernel/drivers/scsi/scsi_transport_sas.ko
kernel/drivers/scsi/isci/isci.ko: kernel/drivers/scsi/libsas/libsas.ko kernel/drivers/scsi/scsi_transport_sas.ko
kernel/drivers/scsi/ips.ko:
kernel/drivers/scsi/fdomain.ko:
kernel/drivers/scsi/fdomain_pci.ko: kernel/drivers/scsi/fdomain.ko
kernel/drivers/scsi/qlogicfas408.ko:
kernel/drivers/scsi/qla1280.ko:
kernel/drivers/scsi/qla2xxx/qla2xxx.ko: kernel/drivers/nvme/host/nvme-fc.ko kernel/drivers/nvme/host/nvme-fabrics.ko kernel/drivers/nvme/host/nvme-core.ko kernel/drivers/scsi/scsi_transport_fc.ko
kernel/drivers/scsi/qla2xxx/tcm_qla2xxx.ko: kernel/drivers/scsi/qla2xxx/qla2xxx.ko kernel/drivers/nvme/host/nvme-fc.ko kernel/drivers/nvme/host/nvme-fabrics.ko kernel/drivers/nvme/host/nvme-core.ko kernel/drivers/scsi/scsi_transport_fc.ko kernel/drivers/target/target_core_mod.ko
kernel/drivers/scsi/qla4xxx/qla4xxx.ko: kernel/drivers/scsi/iscsi_boot_sysfs.ko kernel/drivers/scsi/libiscsi.ko kernel/drivers/scsi/scsi_transport_iscsi.ko
kernel/drivers/scsi/lpfc/lpfc.ko: kernel/drivers/nvme/target/nvmet-fc.ko kernel/drivers/nvme/target/nvmet.ko kernel/drivers/nvme/host/nvme-fc.ko kernel/drivers/nvme/host/nvme-fabrics.ko kernel/drivers/nvme/host/nvme-core.ko kernel/drivers/scsi/scsi_transport_fc.ko
kernel/drivers/scsi/elx/efct.ko: kernel/drivers/scsi/scsi_transport_fc.ko kernel/drivers/target/target_core_mod.ko
kernel/drivers/scsi/bfa/bfa.ko: kernel/drivers/scsi/scsi_transport_fc.ko
kernel/drivers/scsi/csiostor/csiostor.ko: kernel/drivers/scsi/scsi_transport_fc.ko
kernel/drivers/scsi/dmx3191d.ko: kernel/drivers/scsi/scsi_transport_spi.ko
kernel/drivers/scsi/hpsa.ko: kernel/drivers/scsi/scsi_transport_sas.ko
kernel/drivers/scsi/smartpqi/smartpqi.ko: kernel/drivers/scsi/scsi_transport_sas.ko
kernel/drivers/scsi/sym53c8xx_2/sym53c8xx.ko: kernel/drivers/scsi/scsi_transport_spi.ko
kernel/drivers/scsi/dc395x.ko: kernel/drivers/scsi/scsi_transport_spi.ko
kernel/drivers/scsi/esp_scsi.ko: kernel/drivers/scsi/scsi_transport_spi.ko
kernel/drivers/scsi/am53c974.ko: kernel/drivers/scsi/esp_scsi.ko kernel/drivers/scsi/scsi_transport_spi.ko
kernel/drivers/scsi/megaraid.ko:
kernel/drivers/scsi/mpt3sas/mpt3sas.ko: kernel/drivers/scsi/raid_class.ko kernel/drivers/scsi/scsi_transport_sas.ko
kernel/drivers/scsi/mpi3mr/mpi3mr.ko:
kernel/drivers/scsi/ufs/ufshcd-core.ko:
kernel/drivers/scsi/ufs/tc-dwc-g210-pci.ko: kernel/drivers/scsi/ufs/tc-dwc-g210.ko kernel/drivers/scsi/ufs/ufshcd-dwc.ko kernel/drivers/scsi/ufs/ufshcd-core.ko
kernel/drivers/scsi/ufs/ufshcd-dwc.ko: kernel/drivers/scsi/ufs/ufshcd-core.ko
kernel/drivers/scsi/ufs/tc-dwc-g210.ko: kernel/drivers/scsi/ufs/ufshcd-dwc.ko kernel/drivers/scsi/ufs/ufshcd-core.ko
kernel/drivers/scsi/ufs/tc-dwc-g210-pltfrm.ko: kernel/drivers/scsi/ufs/ufshcd-pltfrm.ko kernel/drivers/scsi/ufs/tc-dwc-g210.ko kernel/drivers/scsi/ufs/ufshcd-dwc.ko kernel/drivers/scsi/ufs/ufshcd-core.ko
kernel/drivers/scsi/ufs/cdns-pltfrm.ko: kernel/drivers/scsi/ufs/ufshcd-pltfrm.ko kernel/drivers/scsi/ufs/ufshcd-core.ko
kernel/drivers/scsi/ufs/ufshcd-pci.ko: kernel/drivers/scsi/ufs/ufshcd-core.ko
kernel/drivers/scsi/ufs/ufshcd-pltfrm.ko: kernel/drivers/scsi/ufs/ufshcd-core.ko
kernel/drivers/scsi/atp870u.ko:
kernel/drivers/scsi/initio.ko:
kernel/drivers/scsi/a100u2w.ko:
kernel/drivers/scsi/myrb.ko: kernel/drivers/scsi/raid_class.ko
kernel/drivers/scsi/myrs.ko: kernel/drivers/scsi/raid_class.ko
kernel/drivers/scsi/3w-xxxx.ko:
kernel/drivers/scsi/3w-9xxx.ko:
kernel/drivers/scsi/3w-sas.ko:
kernel/drivers/scsi/ppa.ko: kernel/drivers/parport/parport.ko
kernel/drivers/scsi/imm.ko: kernel/drivers/parport/parport.ko
kernel/drivers/scsi/ipr.ko:
kernel/drivers/scsi/hptiop.ko:
kernel/drivers/scsi/stex.ko:
kernel/drivers/scsi/mvsas/mvsas.ko: kernel/drivers/scsi/libsas/libsas.ko kernel/drivers/scsi/scsi_transport_sas.ko
kernel/drivers/scsi/mvumi.ko:
kernel/drivers/scsi/cxgbi/libcxgbi.ko: kernel/drivers/net/ethernet/chelsio/libcxgb/libcxgb.ko kernel/drivers/scsi/libiscsi_tcp.ko kernel/drivers/scsi/libiscsi.ko kernel/drivers/scsi/scsi_transport_iscsi.ko
kernel/drivers/scsi/cxgbi/cxgb3i/cxgb3i.ko: kernel/drivers/net/ethernet/chelsio/cxgb3/cxgb3.ko kernel/drivers/net/mdio.ko kernel/drivers/scsi/cxgbi/libcxgbi.ko kernel/drivers/net/ethernet/chelsio/libcxgb/libcxgb.ko kernel/drivers/scsi/libiscsi_tcp.ko kernel/drivers/scsi/libiscsi.ko kernel/drivers/scsi/scsi_transport_iscsi.ko
kernel/drivers/scsi/cxgbi/cxgb4i/cxgb4i.ko: kernel/drivers/net/ethernet/chelsio/cxgb4/cxgb4.ko kernel/net/tls/tls.ko kernel/drivers/scsi/cxgbi/libcxgbi.ko kernel/drivers/net/ethernet/chelsio/libcxgb/libcxgb.ko kernel/drivers/scsi/libiscsi_tcp.ko kernel/drivers/scsi/libiscsi.ko kernel/drivers/scsi/scsi_transport_iscsi.ko
kernel/drivers/scsi/bnx2i/bnx2i.ko: kernel/drivers/scsi/libiscsi.ko kernel/drivers/scsi/scsi_transport_iscsi.ko kernel/drivers/net/ethernet/broadcom/cnic.ko kernel/drivers/uio/uio.ko
kernel/drivers/scsi/qedi/qedi.ko: kernel/drivers/scsi/iscsi_boot_sysfs.ko kernel/drivers/scsi/libiscsi.ko kernel/drivers/scsi/scsi_transport_iscsi.ko kernel/drivers/net/ethernet/qlogic/qed/qed.ko kernel/drivers/uio/uio.ko kernel/lib/crc8.ko
kernel/drivers/scsi/be2iscsi/be2iscsi.ko: kernel/drivers/scsi/iscsi_boot_sysfs.ko kernel/drivers/scsi/libiscsi.ko kernel/drivers/scsi/scsi_transport_iscsi.ko
kernel/drivers/scsi/esas2r/esas2r.ko:
kernel/drivers/scsi/pmcraid.ko:
kernel/drivers/scsi/virtio_scsi.ko:
kernel/drivers/scsi/vmw_pvscsi.ko:
kernel/drivers/scsi/xen-scsifront.ko:
kernel/drivers/scsi/hv_storvsc.ko: kernel/drivers/scsi/scsi_transport_fc.ko kernel/drivers/hv/hv_vmbus.ko
kernel/drivers/scsi/wd719x.ko: kernel/drivers/misc/eeprom/eeprom_93cx6.ko
kernel/drivers/scsi/st.ko:
kernel/drivers/scsi/ch.ko:
kernel/drivers/scsi/ses.ko: kernel/drivers/misc/enclosure.ko kernel/drivers/scsi/scsi_transport_sas.ko
kernel/drivers/scsi/scsi_debug.ko:
kernel/drivers/nvme/host/nvme-core.ko:
kernel/drivers/nvme/host/nvme.ko: kernel/drivers/nvme/host/nvme-core.ko
kernel/drivers/nvme/host/nvme-fabrics.ko: kernel/drivers/nvme/host/nvme-core.ko
kernel/drivers/nvme/host/nvme-rdma.ko: kernel/drivers/nvme/host/nvme-fabrics.ko kernel/drivers/nvme/host/nvme-core.ko kernel/drivers/infiniband/core/rdma_cm.ko kernel/drivers/infiniband/core/iw_cm.ko kernel/drivers/infiniband/core/ib_cm.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/nvme/host/nvme-fc.ko: kernel/drivers/nvme/host/nvme-fabrics.ko kernel/drivers/nvme/host/nvme-core.ko
kernel/drivers/nvme/host/nvme-tcp.ko: kernel/drivers/nvme/host/nvme-fabrics.ko kernel/drivers/nvme/host/nvme-core.ko
kernel/drivers/nvme/target/nvmet.ko: kernel/drivers/nvme/host/nvme-core.ko
kernel/drivers/nvme/target/nvme-loop.ko: kernel/drivers/nvme/target/nvmet.ko kernel/drivers/nvme/host/nvme-fabrics.ko kernel/drivers/nvme/host/nvme-core.ko
kernel/drivers/nvme/target/nvmet-rdma.ko: kernel/drivers/nvme/target/nvmet.ko kernel/drivers/nvme/host/nvme-core.ko kernel/drivers/infiniband/core/rdma_cm.ko kernel/drivers/infiniband/core/iw_cm.ko kernel/drivers/infiniband/core/ib_cm.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/nvme/target/nvmet-fc.ko: kernel/drivers/nvme/target/nvmet.ko kernel/drivers/nvme/host/nvme-core.ko
kernel/drivers/nvme/target/nvmet-tcp.ko: kernel/drivers/nvme/target/nvmet.ko kernel/drivers/nvme/host/nvme-core.ko
kernel/drivers/ata/ahci.ko: kernel/drivers/ata/libahci.ko
kernel/drivers/ata/libahci.ko:
kernel/drivers/ata/acard-ahci.ko: kernel/drivers/ata/libahci.ko
kernel/drivers/ata/ahci_platform.ko: kernel/drivers/ata/libahci_platform.ko kernel/drivers/ata/libahci.ko
kernel/drivers/ata/libahci_platform.ko: kernel/drivers/ata/libahci.ko
kernel/drivers/ata/sata_inic162x.ko:
kernel/drivers/ata/sata_sil24.ko:
kernel/drivers/ata/sata_dwc_460ex.ko: kernel/drivers/dma/dw/dw_dmac_core.ko
kernel/drivers/ata/pdc_adma.ko:
kernel/drivers/ata/sata_qstor.ko:
kernel/drivers/ata/sata_sx4.ko:
kernel/drivers/ata/sata_mv.ko:
kernel/drivers/ata/sata_nv.ko:
kernel/drivers/ata/sata_promise.ko:
kernel/drivers/ata/sata_sil.ko:
kernel/drivers/ata/sata_sis.ko:
kernel/drivers/ata/sata_svw.ko:
kernel/drivers/ata/sata_uli.ko:
kernel/drivers/ata/sata_via.ko:
kernel/drivers/ata/sata_vsc.ko:
kernel/drivers/ata/pata_ali.ko:
kernel/drivers/ata/pata_amd.ko:
kernel/drivers/ata/pata_artop.ko:
kernel/drivers/ata/pata_atiixp.ko:
kernel/drivers/ata/pata_atp867x.ko:
kernel/drivers/ata/pata_cmd64x.ko:
kernel/drivers/ata/pata_cypress.ko:
kernel/drivers/ata/pata_efar.ko:
kernel/drivers/ata/pata_hpt366.ko:
kernel/drivers/ata/pata_hpt37x.ko:
kernel/drivers/ata/pata_hpt3x2n.ko:
kernel/drivers/ata/pata_hpt3x3.ko:
kernel/drivers/ata/pata_it8213.ko:
kernel/drivers/ata/pata_it821x.ko:
kernel/drivers/ata/pata_jmicron.ko:
kernel/drivers/ata/pata_marvell.ko:
kernel/drivers/ata/pata_netcell.ko:
kernel/drivers/ata/pata_ninja32.ko:
kernel/drivers/ata/pata_ns87415.ko:
kernel/drivers/ata/pata_oldpiix.ko:
kernel/drivers/ata/pata_optidma.ko:
kernel/drivers/ata/pata_pdc2027x.ko:
kernel/drivers/ata/pata_pdc202xx_old.ko:
kernel/drivers/ata/pata_radisys.ko:
kernel/drivers/ata/pata_rdc.ko:
kernel/drivers/ata/pata_sch.ko:
kernel/drivers/ata/pata_serverworks.ko:
kernel/drivers/ata/pata_sil680.ko:
kernel/drivers/ata/pata_piccolo.ko:
kernel/drivers/ata/pata_triflex.ko:
kernel/drivers/ata/pata_via.ko:
kernel/drivers/ata/pata_sl82c105.ko:
kernel/drivers/ata/pata_cmd640.ko:
kernel/drivers/ata/pata_mpiix.ko:
kernel/drivers/ata/pata_ns87410.ko:
kernel/drivers/ata/pata_opti.ko:
kernel/drivers/ata/pata_pcmcia.ko: kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/ata/pata_platform.ko:
kernel/drivers/ata/pata_rz1000.ko:
kernel/drivers/ata/pata_acpi.ko:
kernel/drivers/ata/pata_legacy.ko:
kernel/drivers/gpu/drm/i2c/ch7006.ko: kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/i2c/sil164.ko: kernel/drivers/gpu/drm/drm.ko
kernel/drivers/gpu/drm/i2c/tda998x.ko: kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/i2c/tda9950.ko: kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/gpu/drm/panel/panel-raspberrypi-touchscreen.ko: kernel/drivers/gpu/drm/drm.ko
kernel/drivers/gpu/drm/panel/panel-widechips-ws2401.ko: kernel/drivers/gpu/drm/drm_mipi_dbi.ko kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/bridge/analogix/analogix-anx78xx.ko: kernel/drivers/gpu/drm/bridge/analogix/analogix_dp.ko kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/bridge/analogix/analogix_dp.ko: kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/tiny/bochs.ko: kernel/drivers/gpu/drm/drm_vram_helper.ko kernel/drivers/gpu/drm/drm_ttm_helper.ko kernel/drivers/gpu/drm/ttm/ttm.ko kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/tiny/cirrus.ko: kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/tiny/gm12u320.ko: kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/tiny/simpledrm.ko: kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/tiny/hx8357d.ko: kernel/drivers/gpu/drm/drm_mipi_dbi.ko kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/tiny/ili9225.ko: kernel/drivers/gpu/drm/drm_mipi_dbi.ko kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/tiny/ili9341.ko: kernel/drivers/gpu/drm/drm_mipi_dbi.ko kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/tiny/ili9486.ko: kernel/drivers/gpu/drm/drm_mipi_dbi.ko kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/tiny/mi0283qt.ko: kernel/drivers/gpu/drm/drm_mipi_dbi.ko kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/tiny/repaper.ko: kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/tiny/st7586.ko: kernel/drivers/gpu/drm/drm_mipi_dbi.ko kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/tiny/st7735r.ko: kernel/drivers/gpu/drm/drm_mipi_dbi.ko kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/xen/drm_xen_front.ko: kernel/drivers/xen/xen-front-pgdir-shbuf.ko kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/gud/gud.ko: kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko kernel/lib/lz4/lz4_compress.ko
kernel/drivers/gpu/drm/drm_vram_helper.ko: kernel/drivers/gpu/drm/drm_ttm_helper.ko kernel/drivers/gpu/drm/ttm/ttm.ko kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/drm_ttm_helper.ko: kernel/drivers/gpu/drm/ttm/ttm.ko kernel/drivers/gpu/drm/drm.ko
kernel/drivers/gpu/drm/drm_kms_helper.ko: kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/drm.ko:
kernel/drivers/gpu/drm/drm_mipi_dbi.ko: kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/ttm/ttm.ko: kernel/drivers/gpu/drm/drm.ko
kernel/drivers/gpu/drm/scheduler/gpu-sched.ko: kernel/drivers/gpu/drm/drm.ko
kernel/drivers/gpu/drm/radeon/radeon.ko: kernel/drivers/gpu/drm/drm_ttm_helper.ko kernel/drivers/gpu/drm/ttm/ttm.ko kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/amd/amdgpu/amdgpu.ko: kernel/drivers/iommu/amd/iommu_v2.ko kernel/drivers/gpu/drm/scheduler/gpu-sched.ko kernel/drivers/gpu/drm/drm_ttm_helper.ko kernel/drivers/gpu/drm/ttm/ttm.ko kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/i915/i915.ko: kernel/drivers/gpu/drm/ttm/ttm.ko kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko kernel/drivers/acpi/video.ko
kernel/drivers/gpu/drm/i915/gvt/kvmgt.ko: kernel/drivers/vfio/mdev/mdev.ko kernel/drivers/gpu/drm/i915/i915.ko kernel/drivers/gpu/drm/ttm/ttm.ko kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko kernel/drivers/acpi/video.ko kernel/arch/x86/kvm/kvm.ko
kernel/drivers/gpu/drm/mgag200/mgag200.ko: kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/vmwgfx/vmwgfx.ko: kernel/drivers/gpu/drm/ttm/ttm.ko kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/vgem/vgem.ko: kernel/drivers/gpu/drm/drm.ko
kernel/drivers/gpu/drm/vkms/vkms.ko: kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/nouveau/nouveau.ko: kernel/drivers/platform/x86/mxm-wmi.ko kernel/drivers/gpu/drm/drm_ttm_helper.ko kernel/drivers/gpu/drm/ttm/ttm.ko kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko kernel/drivers/acpi/video.ko kernel/drivers/platform/x86/wmi.ko
kernel/drivers/gpu/drm/gma500/gma500_gfx.ko: kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko kernel/drivers/acpi/video.ko
kernel/drivers/gpu/drm/udl/udl.ko: kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/ast/ast.ko: kernel/drivers/gpu/drm/drm_vram_helper.ko kernel/drivers/gpu/drm/drm_ttm_helper.ko kernel/drivers/gpu/drm/ttm/ttm.ko kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/qxl/qxl.ko: kernel/drivers/gpu/drm/drm_ttm_helper.ko kernel/drivers/gpu/drm/ttm/ttm.ko kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/virtio/virtio-gpu.ko: kernel/drivers/virtio/virtio_dma_buf.ko kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/vboxvideo/vboxvideo.ko: kernel/drivers/gpu/drm/drm_vram_helper.ko kernel/drivers/gpu/drm/drm_ttm_helper.ko kernel/drivers/gpu/drm/ttm/ttm.ko kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/gpu/drm/hyperv/hyperv_drm.ko: kernel/drivers/gpu/drm/drm_kms_helper.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/gpu/drm/drm.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko kernel/drivers/hv/hv_vmbus.ko
kernel/drivers/spi/spi-mux.ko: kernel/drivers/mux/mux-core.ko
kernel/drivers/spi/spidev.ko:
kernel/drivers/spi/spi-loopback-test.ko:
kernel/drivers/spi/spi-altera-platform.ko: kernel/drivers/spi/spi-altera-core.ko
kernel/drivers/spi/spi-altera-core.ko:
kernel/drivers/spi/spi-altera-dfl.ko: kernel/drivers/fpga/dfl.ko kernel/drivers/fpga/fpga-region.ko kernel/drivers/fpga/fpga-bridge.ko kernel/drivers/fpga/fpga-mgr.ko kernel/drivers/spi/spi-altera-core.ko
kernel/drivers/spi/spi-axi-spi-engine.ko:
kernel/drivers/spi/spi-bitbang.ko:
kernel/drivers/spi/spi-butterfly.ko: kernel/drivers/spi/spi-bitbang.ko kernel/drivers/parport/parport.ko
kernel/drivers/spi/spi-cadence.ko:
kernel/drivers/spi/spi-dln2.ko: kernel/drivers/mfd/dln2.ko
kernel/drivers/spi/spi-dw.ko:
kernel/drivers/spi/spi-dw-mmio.ko: kernel/drivers/spi/spi-dw.ko
kernel/drivers/spi/spi-dw-pci.ko: kernel/drivers/spi/spi-dw.ko
kernel/drivers/spi/spi-gpio.ko: kernel/drivers/spi/spi-bitbang.ko
kernel/drivers/spi/spi-lantiq-ssc.ko:
kernel/drivers/spi/spi-ljca.ko: kernel/drivers/mfd/ljca.ko
kernel/drivers/spi/spi-lm70llp.ko: kernel/drivers/spi/spi-bitbang.ko kernel/drivers/parport/parport.ko
kernel/drivers/spi/spi-mxic.ko:
kernel/drivers/spi/spi-nxp-fspi.ko:
kernel/drivers/spi/spi-oc-tiny.ko: kernel/drivers/spi/spi-bitbang.ko
kernel/drivers/spi/spi-pxa2xx-platform.ko:
kernel/drivers/spi/spi-pxa2xx-pci.ko:
kernel/drivers/spi/spi-sc18is602.ko:
kernel/drivers/spi/spi-sifive.ko:
kernel/drivers/spi/spi-tle62x0.ko:
kernel/drivers/spi/spi-xcomm.ko:
kernel/drivers/spi/spi-zynqmp-gqspi.ko:
kernel/drivers/spi/spi-amd.ko:
kernel/drivers/spi/spi-slave-time.ko:
kernel/drivers/spi/spi-slave-system-control.ko:
kernel/drivers/net/phy/phylink.ko:
kernel/drivers/net/phy/sfp.ko: kernel/drivers/net/mdio/mdio-i2c.ko
kernel/drivers/net/phy/adin.ko:
kernel/drivers/net/phy/amd.ko:
kernel/drivers/net/phy/aquantia.ko:
kernel/drivers/net/phy/at803x.ko:
kernel/drivers/net/phy/ax88796b.ko:
kernel/drivers/net/phy/bcm54140.ko: kernel/drivers/net/phy/bcm-phy-lib.ko
kernel/drivers/net/phy/bcm7xxx.ko: kernel/drivers/net/phy/bcm-phy-lib.ko
kernel/drivers/net/phy/bcm87xx.ko:
kernel/drivers/net/phy/bcm-phy-lib.ko:
kernel/drivers/net/phy/broadcom.ko: kernel/drivers/net/phy/bcm-phy-lib.ko
kernel/drivers/net/phy/cicada.ko:
kernel/drivers/net/phy/cortina.ko:
kernel/drivers/net/phy/davicom.ko:
kernel/drivers/net/phy/dp83640.ko:
kernel/drivers/net/phy/dp83822.ko:
kernel/drivers/net/phy/dp83848.ko:
kernel/drivers/net/phy/dp83867.ko:
kernel/drivers/net/phy/dp83869.ko:
kernel/drivers/net/phy/dp83tc811.ko:
kernel/drivers/net/phy/icplus.ko:
kernel/drivers/net/phy/intel-xway.ko:
kernel/drivers/net/phy/et1011c.ko:
kernel/drivers/net/phy/lxt.ko:
kernel/drivers/net/phy/marvell10g.ko:
kernel/drivers/net/phy/marvell.ko:
kernel/drivers/net/phy/marvell-88x2222.ko:
kernel/drivers/net/phy/mxl-gpy.ko:
kernel/drivers/net/phy/mediatek-ge.ko:
kernel/drivers/net/phy/spi_ks8995.ko:
kernel/drivers/net/phy/micrel.ko:
kernel/drivers/net/phy/microchip.ko:
kernel/drivers/net/phy/microchip_t1.ko:
kernel/drivers/net/phy/mscc/mscc.ko: kernel/drivers/net/macsec.ko
kernel/drivers/net/phy/motorcomm.ko:
kernel/drivers/net/phy/national.ko:
kernel/drivers/net/phy/nxp-c45-tja11xx.ko:
kernel/drivers/net/phy/nxp-tja11xx.ko:
kernel/drivers/net/phy/qsemi.ko:
kernel/drivers/net/phy/realtek.ko:
kernel/drivers/net/phy/uPD60620.ko:
kernel/drivers/net/phy/rockchip.ko:
kernel/drivers/net/phy/smsc.ko:
kernel/drivers/net/phy/ste10Xp.ko:
kernel/drivers/net/phy/teranetics.ko:
kernel/drivers/net/phy/vitesse.ko:
kernel/drivers/net/phy/xilinx_gmii2rgmii.ko:
kernel/drivers/net/mdio/mdio-bcm-unimac.ko:
kernel/drivers/net/mdio/mdio-bitbang.ko:
kernel/drivers/net/mdio/mdio-cavium.ko:
kernel/drivers/net/mdio/mdio-gpio.ko: kernel/drivers/net/mdio/mdio-bitbang.ko
kernel/drivers/net/mdio/mdio-i2c.ko:
kernel/drivers/net/mdio/mdio-mscc-miim.ko:
kernel/drivers/net/mdio/mdio-mvusb.ko:
kernel/drivers/net/mdio/mdio-thunder.ko:
kernel/drivers/net/pcs/pcs_xpcs.ko:
kernel/drivers/net/pcs/pcs-lynx.ko: kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/ethernet/3com/3c509.ko:
kernel/drivers/net/ethernet/3com/3c589_cs.ko: kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/net/ethernet/3com/3c574_cs.ko: kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/net/ethernet/3com/3c59x.ko: kernel/drivers/net/mii.ko
kernel/drivers/net/ethernet/3com/typhoon.ko:
kernel/drivers/net/ethernet/8390/ne2k-pci.ko: kernel/drivers/net/ethernet/8390/8390.ko
kernel/drivers/net/ethernet/8390/8390.ko:
kernel/drivers/net/ethernet/8390/axnet_cs.ko: kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/net/ethernet/8390/pcnet_cs.ko: kernel/drivers/net/ethernet/8390/8390.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/net/ethernet/adaptec/starfire.ko: kernel/drivers/net/mii.ko
kernel/drivers/net/ethernet/agere/et131x.ko:
kernel/drivers/net/ethernet/alacritech/slicoss.ko:
kernel/drivers/net/ethernet/alteon/acenic.ko:
kernel/drivers/net/ethernet/amazon/ena/ena.ko:
kernel/drivers/net/ethernet/amd/amd8111e.ko: kernel/drivers/net/mii.ko
kernel/drivers/net/ethernet/amd/nmclan_cs.ko: kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/net/ethernet/amd/pcnet32.ko: kernel/drivers/net/mii.ko
kernel/drivers/net/ethernet/amd/xgbe/amd-xgbe.ko:
kernel/drivers/net/ethernet/aquantia/atlantic/atlantic.ko: kernel/drivers/net/macsec.ko
kernel/drivers/net/ethernet/atheros/atlx/atl1.ko: kernel/drivers/net/mii.ko
kernel/drivers/net/ethernet/atheros/atlx/atl2.ko:
kernel/drivers/net/ethernet/atheros/atl1e/atl1e.ko:
kernel/drivers/net/ethernet/atheros/atl1c/atl1c.ko:
kernel/drivers/net/ethernet/atheros/alx/alx.ko: kernel/drivers/net/mdio.ko
kernel/drivers/net/ethernet/cadence/macb.ko: kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/ethernet/cadence/macb_pci.ko:
kernel/drivers/net/ethernet/broadcom/b44.ko: kernel/drivers/ssb/ssb.ko kernel/drivers/net/mii.ko
kernel/drivers/net/ethernet/broadcom/genet/genet.ko:
kernel/drivers/net/ethernet/broadcom/bnx2.ko:
kernel/drivers/net/ethernet/broadcom/cnic.ko: kernel/drivers/uio/uio.ko
kernel/drivers/net/ethernet/broadcom/bnx2x/bnx2x.ko: kernel/drivers/net/mdio.ko kernel/lib/libcrc32c.ko
kernel/drivers/net/ethernet/broadcom/tg3.ko:
kernel/drivers/net/ethernet/broadcom/bcmsysport.ko:
kernel/drivers/net/ethernet/broadcom/bnxt/bnxt_en.ko:
kernel/drivers/net/ethernet/brocade/bna/bna.ko:
kernel/drivers/net/ethernet/cavium/common/cavium_ptp.ko:
kernel/drivers/net/ethernet/cavium/thunder/thunder_xcv.ko:
kernel/drivers/net/ethernet/cavium/thunder/thunder_bgx.ko: kernel/drivers/net/ethernet/cavium/thunder/thunder_xcv.ko
kernel/drivers/net/ethernet/cavium/thunder/nicpf.ko: kernel/drivers/net/ethernet/cavium/thunder/thunder_bgx.ko kernel/drivers/net/ethernet/cavium/thunder/thunder_xcv.ko
kernel/drivers/net/ethernet/cavium/thunder/nicvf.ko: kernel/drivers/net/ethernet/cavium/common/cavium_ptp.ko
kernel/drivers/net/ethernet/cavium/liquidio/liquidio.ko:
kernel/drivers/net/ethernet/cavium/liquidio/liquidio_vf.ko:
kernel/drivers/net/ethernet/chelsio/inline_crypto/ch_ipsec/ch_ipsec.ko: kernel/drivers/net/ethernet/chelsio/cxgb4/cxgb4.ko kernel/net/tls/tls.ko
kernel/drivers/net/ethernet/chelsio/inline_crypto/ch_ktls/ch_ktls.ko: kernel/drivers/net/ethernet/chelsio/cxgb4/cxgb4.ko kernel/net/tls/tls.ko
kernel/drivers/net/ethernet/chelsio/cxgb/cxgb.ko: kernel/drivers/net/mdio.ko
kernel/drivers/net/ethernet/chelsio/cxgb3/cxgb3.ko: kernel/drivers/net/mdio.ko
kernel/drivers/net/ethernet/chelsio/cxgb4/cxgb4.ko: kernel/net/tls/tls.ko
kernel/drivers/net/ethernet/chelsio/cxgb4vf/cxgb4vf.ko:
kernel/drivers/net/ethernet/chelsio/libcxgb/libcxgb.ko:
kernel/drivers/net/ethernet/cisco/enic/enic.ko:
kernel/drivers/net/ethernet/dec/tulip/xircom_cb.ko:
kernel/drivers/net/ethernet/dec/tulip/dmfe.ko:
kernel/drivers/net/ethernet/dec/tulip/winbond-840.ko: kernel/drivers/net/mii.ko
kernel/drivers/net/ethernet/dec/tulip/de2104x.ko:
kernel/drivers/net/ethernet/dec/tulip/tulip.ko:
kernel/drivers/net/ethernet/dec/tulip/de4x5.ko:
kernel/drivers/net/ethernet/dec/tulip/uli526x.ko:
kernel/drivers/net/ethernet/dlink/dl2k.ko:
kernel/drivers/net/ethernet/dlink/sundance.ko: kernel/drivers/net/mii.ko
kernel/drivers/net/ethernet/emulex/benet/be2net.ko:
kernel/drivers/net/ethernet/fujitsu/fmvj18x_cs.ko: kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/net/ethernet/google/gve/gve.ko:
kernel/drivers/net/ethernet/huawei/hinic/hinic.ko:
kernel/drivers/net/ethernet/intel/e100.ko: kernel/drivers/net/mii.ko
kernel/drivers/net/ethernet/intel/e1000/e1000.ko:
kernel/drivers/net/ethernet/intel/e1000e/e1000e.ko:
kernel/drivers/net/ethernet/intel/igb/igb.ko: kernel/drivers/dca/dca.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/net/ethernet/intel/igc/igc.ko:
kernel/drivers/net/ethernet/intel/igbvf/igbvf.ko:
kernel/drivers/net/ethernet/intel/ixgbe/ixgbe.ko: kernel/net/xfrm/xfrm_algo.ko kernel/drivers/net/mdio.ko kernel/drivers/dca/dca.ko
kernel/drivers/net/ethernet/intel/ixgbevf/ixgbevf.ko:
kernel/drivers/net/ethernet/intel/i40e/i40e.ko:
kernel/drivers/net/ethernet/intel/ixgb/ixgb.ko:
kernel/drivers/net/ethernet/intel/iavf/iavf.ko:
kernel/drivers/net/ethernet/intel/fm10k/fm10k.ko:
kernel/drivers/net/ethernet/intel/ice/ice.ko:
kernel/drivers/net/ethernet/microsoft/mana/mana.ko:
kernel/drivers/net/ethernet/marvell/prestera/prestera.ko: kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko
kernel/drivers/net/ethernet/marvell/prestera/prestera_pci.ko: kernel/drivers/net/ethernet/marvell/prestera/prestera.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko
kernel/drivers/net/ethernet/marvell/mvmdio.ko:
kernel/drivers/net/ethernet/marvell/skge.ko:
kernel/drivers/net/ethernet/marvell/sky2.ko:
kernel/drivers/net/ethernet/mellanox/mlx4/mlx4_core.ko:
kernel/drivers/net/ethernet/mellanox/mlx4/mlx4_en.ko: kernel/drivers/net/ethernet/mellanox/mlx4/mlx4_core.ko
kernel/drivers/net/ethernet/mellanox/mlx5/core/mlx5_core.ko: kernel/drivers/net/ethernet/mellanox/mlxfw/mlxfw.ko kernel/net/psample/psample.ko kernel/net/tls/tls.ko kernel/drivers/pci/controller/pci-hyperv-intf.ko
kernel/drivers/net/ethernet/mellanox/mlxsw/mlxsw_core.ko: kernel/drivers/net/ethernet/mellanox/mlxfw/mlxfw.ko
kernel/drivers/net/ethernet/mellanox/mlxsw/mlxsw_pci.ko: kernel/drivers/net/ethernet/mellanox/mlxsw/mlxsw_core.ko kernel/drivers/net/ethernet/mellanox/mlxfw/mlxfw.ko
kernel/drivers/net/ethernet/mellanox/mlxsw/mlxsw_i2c.ko: kernel/drivers/net/ethernet/mellanox/mlxsw/mlxsw_core.ko kernel/drivers/net/ethernet/mellanox/mlxfw/mlxfw.ko
kernel/drivers/net/ethernet/mellanox/mlxsw/mlxsw_spectrum.ko: kernel/drivers/net/ethernet/mellanox/mlxsw/mlxsw_pci.ko kernel/drivers/net/ethernet/mellanox/mlxsw/mlxsw_core.ko kernel/drivers/net/ethernet/mellanox/mlxfw/mlxfw.ko kernel/drivers/net/vxlan.ko kernel/net/ipv6/ip6_tunnel.ko kernel/net/ipv6/tunnel6.ko kernel/lib/objagg.ko kernel/net/psample/psample.ko kernel/lib/parman.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/net/ipv6/ip6_udp_tunnel.ko kernel/net/ipv4/udp_tunnel.ko
kernel/drivers/net/ethernet/mellanox/mlxsw/mlxsw_minimal.ko: kernel/drivers/net/ethernet/mellanox/mlxsw/mlxsw_i2c.ko kernel/drivers/net/ethernet/mellanox/mlxsw/mlxsw_core.ko kernel/drivers/net/ethernet/mellanox/mlxfw/mlxfw.ko
kernel/drivers/net/ethernet/mellanox/mlxfw/mlxfw.ko:
kernel/drivers/net/ethernet/micrel/ks8842.ko:
kernel/drivers/net/ethernet/micrel/ks8851_common.ko: kernel/drivers/net/mii.ko kernel/drivers/misc/eeprom/eeprom_93cx6.ko
kernel/drivers/net/ethernet/micrel/ks8851_spi.ko: kernel/drivers/net/ethernet/micrel/ks8851_common.ko kernel/drivers/net/mii.ko kernel/drivers/misc/eeprom/eeprom_93cx6.ko
kernel/drivers/net/ethernet/micrel/ks8851_par.ko: kernel/drivers/net/ethernet/micrel/ks8851_common.ko kernel/drivers/net/mii.ko kernel/drivers/misc/eeprom/eeprom_93cx6.ko
kernel/drivers/net/ethernet/micrel/ksz884x.ko: kernel/drivers/net/mii.ko
kernel/drivers/net/ethernet/microchip/enc28j60.ko:
kernel/drivers/net/ethernet/microchip/encx24j600.ko: kernel/drivers/net/ethernet/microchip/encx24j600-regmap.ko
kernel/drivers/net/ethernet/microchip/encx24j600-regmap.ko:
kernel/drivers/net/ethernet/microchip/lan743x.ko:
kernel/drivers/net/ethernet/mscc/mscc_ocelot_switch_lib.ko:
kernel/drivers/net/ethernet/myricom/myri10ge/myri10ge.ko: kernel/drivers/dca/dca.ko
kernel/drivers/net/ethernet/natsemi/natsemi.ko:
kernel/drivers/net/ethernet/natsemi/ns83820.ko:
kernel/drivers/net/ethernet/neterion/s2io.ko:
kernel/drivers/net/ethernet/neterion/vxge/vxge.ko:
kernel/drivers/net/ethernet/netronome/nfp/nfp.ko: kernel/net/tls/tls.ko
kernel/drivers/net/ethernet/ni/nixge.ko:
kernel/drivers/net/ethernet/nvidia/forcedeth.ko:
kernel/drivers/net/ethernet/packetengines/hamachi.ko: kernel/drivers/net/mii.ko
kernel/drivers/net/ethernet/packetengines/yellowfin.ko:
kernel/drivers/net/ethernet/qlogic/qla3xxx.ko:
kernel/drivers/net/ethernet/qlogic/qlcnic/qlcnic.ko:
kernel/drivers/net/ethernet/qlogic/netxen/netxen_nic.ko:
kernel/drivers/net/ethernet/qlogic/qed/qed.ko: kernel/lib/crc8.ko
kernel/drivers/net/ethernet/qlogic/qede/qede.ko: kernel/drivers/net/ethernet/qlogic/qed/qed.ko kernel/lib/crc8.ko
kernel/drivers/net/ethernet/qualcomm/emac/qcom-emac.ko:
kernel/drivers/net/ethernet/qualcomm/rmnet/rmnet.ko:
kernel/drivers/net/ethernet/realtek/8139cp.ko: kernel/drivers/net/mii.ko
kernel/drivers/net/ethernet/realtek/8139too.ko: kernel/drivers/net/mii.ko
kernel/drivers/net/ethernet/realtek/atp.ko:
kernel/drivers/net/ethernet/realtek/r8169.ko:
kernel/drivers/net/ethernet/rdc/r6040.ko:
kernel/drivers/net/ethernet/rocker/rocker.ko:
kernel/drivers/net/ethernet/samsung/sxgbe/samsung-sxgbe.ko:
kernel/drivers/net/ethernet/silan/sc92031.ko:
kernel/drivers/net/ethernet/sis/sis190.ko: kernel/drivers/net/mii.ko
kernel/drivers/net/ethernet/sis/sis900.ko: kernel/drivers/net/mii.ko
kernel/drivers/net/ethernet/sfc/sfc.ko: kernel/drivers/net/mdio.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/net/ethernet/sfc/falcon/sfc-falcon.ko: kernel/drivers/net/mdio.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/net/ethernet/smsc/smc91c92_cs.ko: kernel/drivers/net/mii.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/net/ethernet/smsc/epic100.ko: kernel/drivers/net/mii.ko
kernel/drivers/net/ethernet/smsc/smsc9420.ko:
kernel/drivers/net/ethernet/smsc/smsc911x.ko:
kernel/drivers/net/ethernet/stmicro/stmmac/stmmac.ko: kernel/drivers/net/pcs/pcs_xpcs.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/ethernet/stmicro/stmmac/stmmac-platform.ko: kernel/drivers/net/ethernet/stmicro/stmmac/stmmac.ko kernel/drivers/net/pcs/pcs_xpcs.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/ethernet/stmicro/stmmac/dwmac-generic.ko: kernel/drivers/net/ethernet/stmicro/stmmac/stmmac-platform.ko kernel/drivers/net/ethernet/stmicro/stmmac/stmmac.ko kernel/drivers/net/pcs/pcs_xpcs.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/ethernet/stmicro/stmmac/stmmac-pci.ko: kernel/drivers/net/ethernet/stmicro/stmmac/stmmac.ko kernel/drivers/net/pcs/pcs_xpcs.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/ethernet/stmicro/stmmac/dwmac-intel.ko: kernel/drivers/net/ethernet/stmicro/stmmac/stmmac.ko kernel/drivers/net/pcs/pcs_xpcs.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/ethernet/stmicro/stmmac/dwmac-loongson.ko: kernel/drivers/net/ethernet/stmicro/stmmac/stmmac.ko kernel/drivers/net/pcs/pcs_xpcs.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/ethernet/sun/sunhme.ko:
kernel/drivers/net/ethernet/sun/sungem.ko: kernel/drivers/net/sungem_phy.ko
kernel/drivers/net/ethernet/sun/cassini.ko:
kernel/drivers/net/ethernet/sun/niu.ko:
kernel/drivers/net/ethernet/tehuti/tehuti.ko:
kernel/drivers/net/ethernet/ti/tlan.ko:
kernel/drivers/net/ethernet/via/via-rhine.ko: kernel/drivers/net/mii.ko
kernel/drivers/net/ethernet/via/via-velocity.ko:
kernel/drivers/net/ethernet/wiznet/w5100.ko:
kernel/drivers/net/ethernet/wiznet/w5100-spi.ko: kernel/drivers/net/ethernet/wiznet/w5100.ko
kernel/drivers/net/ethernet/wiznet/w5300.ko:
kernel/drivers/net/ethernet/xilinx/ll_temac.ko:
kernel/drivers/net/ethernet/xilinx/xilinx_emaclite.ko:
kernel/drivers/net/ethernet/xilinx/xilinx_emac.ko: kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/ethernet/xircom/xirc2ps_cs.ko: kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/net/ethernet/synopsys/dwc-xlgmac.ko:
kernel/drivers/net/ethernet/pensando/ionic/ionic.ko:
kernel/drivers/net/ethernet/altera/altera_tse.ko:
kernel/drivers/net/ethernet/ec_bhf.ko:
kernel/drivers/net/ethernet/dnet.ko:
kernel/drivers/net/ethernet/jme.ko: kernel/drivers/net/mii.ko
kernel/drivers/net/ethernet/fealnx.ko: kernel/drivers/net/mii.ko
kernel/drivers/net/ethernet/ethoc.ko:
kernel/drivers/net/fddi/defxx.ko:
kernel/drivers/net/fddi/skfp/skfp.ko:
kernel/drivers/net/hamradio/mkiss.ko: kernel/net/ax25/ax25.ko
kernel/drivers/net/hamradio/6pack.ko: kernel/net/ax25/ax25.ko
kernel/drivers/net/hamradio/yam.ko: kernel/net/ax25/ax25.ko
kernel/drivers/net/hamradio/bpqether.ko: kernel/net/ax25/ax25.ko
kernel/drivers/net/hamradio/baycom_ser_fdx.ko: kernel/drivers/net/hamradio/hdlcdrv.ko kernel/net/ax25/ax25.ko
kernel/drivers/net/hamradio/hdlcdrv.ko: kernel/net/ax25/ax25.ko
kernel/drivers/net/hamradio/baycom_ser_hdx.ko: kernel/drivers/net/hamradio/hdlcdrv.ko kernel/net/ax25/ax25.ko
kernel/drivers/net/hamradio/baycom_par.ko: kernel/drivers/net/hamradio/hdlcdrv.ko kernel/net/ax25/ax25.ko kernel/drivers/parport/parport.ko
kernel/drivers/net/ppp/ppp_async.ko:
kernel/drivers/net/ppp/bsd_comp.ko:
kernel/drivers/net/ppp/ppp_deflate.ko:
kernel/drivers/net/ppp/ppp_mppe.ko: kernel/lib/crypto/libarc4.ko
kernel/drivers/net/ppp/ppp_synctty.ko:
kernel/drivers/net/ppp/pppox.ko:
kernel/drivers/net/ppp/pppoe.ko: kernel/drivers/net/ppp/pppox.ko
kernel/drivers/net/ppp/pptp.ko: kernel/net/ipv4/gre.ko kernel/drivers/net/ppp/pppox.ko
kernel/drivers/net/slip/slip.ko:
kernel/drivers/net/wan/hdlc.ko:
kernel/drivers/net/wan/hdlc_raw.ko: kernel/drivers/net/wan/hdlc.ko
kernel/drivers/net/wan/hdlc_raw_eth.ko: kernel/drivers/net/wan/hdlc.ko
kernel/drivers/net/wan/hdlc_cisco.ko: kernel/drivers/net/wan/hdlc.ko
kernel/drivers/net/wan/hdlc_fr.ko: kernel/drivers/net/wan/hdlc.ko
kernel/drivers/net/wan/hdlc_ppp.ko: kernel/drivers/net/wan/hdlc.ko
kernel/drivers/net/wan/hdlc_x25.ko: kernel/net/lapb/lapb.ko kernel/drivers/net/wan/hdlc.ko
kernel/drivers/net/wan/farsync.ko: kernel/drivers/net/wan/hdlc.ko
kernel/drivers/net/wan/lmc/lmc.ko: kernel/drivers/net/wan/hdlc.ko
kernel/drivers/net/wan/lapbether.ko: kernel/net/lapb/lapb.ko
kernel/drivers/net/wan/wanxl.ko: kernel/drivers/net/wan/hdlc.ko
kernel/drivers/net/wan/pci200syn.ko: kernel/drivers/net/wan/hdlc.ko
kernel/drivers/net/wan/pc300too.ko: kernel/drivers/net/wan/hdlc.ko
kernel/drivers/net/wireless/admtek/adm8211.ko: kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko kernel/drivers/misc/eeprom/eeprom_93cx6.ko
kernel/drivers/net/wireless/ath/ath5k/ath5k.ko: kernel/drivers/net/wireless/ath/ath.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/ath/ath9k/ath9k.ko: kernel/drivers/net/wireless/ath/ath9k/ath9k_common.ko kernel/drivers/net/wireless/ath/ath9k/ath9k_hw.ko kernel/drivers/net/wireless/ath/ath.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/ath/ath9k/ath9k_hw.ko: kernel/drivers/net/wireless/ath/ath.ko kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/ath/ath9k/ath9k_common.ko: kernel/drivers/net/wireless/ath/ath9k/ath9k_hw.ko kernel/drivers/net/wireless/ath/ath.ko kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/ath/ath9k/ath9k_htc.ko: kernel/drivers/net/wireless/ath/ath9k/ath9k_common.ko kernel/drivers/net/wireless/ath/ath9k/ath9k_hw.ko kernel/drivers/net/wireless/ath/ath.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/ath/ath9k/ath9k_pci_owl_loader.ko:
kernel/drivers/net/wireless/ath/carl9170/carl9170.ko: kernel/drivers/net/wireless/ath/ath.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/ath/ath6kl/ath6kl_core.ko: kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/ath/ath6kl/ath6kl_sdio.ko: kernel/drivers/net/wireless/ath/ath6kl/ath6kl_core.ko kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/ath/ath6kl/ath6kl_usb.ko: kernel/drivers/net/wireless/ath/ath6kl/ath6kl_core.ko kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/ath/ar5523/ar5523.ko: kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/ath/wil6210/wil6210.ko: kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/ath/ath10k/ath10k_core.ko: kernel/drivers/net/wireless/ath/ath.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/ath/ath10k/ath10k_pci.ko: kernel/drivers/net/wireless/ath/ath10k/ath10k_core.ko kernel/drivers/net/wireless/ath/ath.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/ath/ath10k/ath10k_sdio.ko: kernel/drivers/net/wireless/ath/ath10k/ath10k_core.ko kernel/drivers/net/wireless/ath/ath.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/ath/ath10k/ath10k_usb.ko: kernel/drivers/net/wireless/ath/ath10k/ath10k_core.ko kernel/drivers/net/wireless/ath/ath.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/ath/wcn36xx/wcn36xx.ko: kernel/drivers/rpmsg/rpmsg_core.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/ath/ath11k/ath11k.ko: kernel/drivers/soc/qcom/qmi_helpers.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/ath/ath11k/ath11k_ahb.ko: kernel/drivers/net/wireless/ath/ath11k/ath11k.ko kernel/drivers/soc/qcom/qmi_helpers.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/ath/ath11k/ath11k_pci.ko: kernel/drivers/net/wireless/ath/ath11k/ath11k.ko kernel/drivers/soc/qcom/qmi_helpers.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko kernel/drivers/bus/mhi/core/mhi.ko
kernel/drivers/net/wireless/ath/ath.ko: kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/atmel/atmel.ko: kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/atmel/atmel_pci.ko: kernel/drivers/net/wireless/atmel/atmel.ko kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/atmel/atmel_cs.ko: kernel/drivers/net/wireless/atmel/atmel.ko kernel/net/wireless/cfg80211.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/net/wireless/atmel/at76c50x-usb.ko: kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/broadcom/b43/b43.ko: kernel/lib/math/cordic.ko kernel/drivers/bcma/bcma.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko kernel/drivers/ssb/ssb.ko
kernel/drivers/net/wireless/broadcom/b43legacy/b43legacy.ko: kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko kernel/drivers/ssb/ssb.ko
kernel/drivers/net/wireless/broadcom/brcm80211/brcmutil/brcmutil.ko:
kernel/drivers/net/wireless/broadcom/brcm80211/brcmfmac/brcmfmac.ko: kernel/drivers/net/wireless/broadcom/brcm80211/brcmutil/brcmutil.ko kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/broadcom/brcm80211/brcmsmac/brcmsmac.ko: kernel/drivers/net/wireless/broadcom/brcm80211/brcmutil/brcmutil.ko kernel/lib/math/cordic.ko kernel/drivers/bcma/bcma.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/cisco/airo.ko: kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/cisco/airo_cs.ko: kernel/drivers/net/wireless/cisco/airo.ko kernel/net/wireless/cfg80211.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/net/wireless/intel/ipw2x00/ipw2100.ko: kernel/drivers/net/wireless/intel/ipw2x00/libipw.ko kernel/net/wireless/lib80211.ko kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/intel/ipw2x00/ipw2200.ko: kernel/drivers/net/wireless/intel/ipw2x00/libipw.ko kernel/net/wireless/lib80211.ko kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/intel/ipw2x00/libipw.ko: kernel/net/wireless/lib80211.ko kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/intel/iwlegacy/iwlegacy.ko: kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/intel/iwlegacy/iwl4965.ko: kernel/drivers/net/wireless/intel/iwlegacy/iwlegacy.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/intel/iwlegacy/iwl3945.ko: kernel/drivers/net/wireless/intel/iwlegacy/iwlegacy.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/intel/iwlwifi/iwlwifi.ko: kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/intel/iwlwifi/dvm/iwldvm.ko: kernel/drivers/net/wireless/intel/iwlwifi/iwlwifi.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/intel/iwlwifi/mvm/iwlmvm.ko: kernel/drivers/net/wireless/intel/iwlwifi/iwlwifi.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/intersil/hostap/hostap.ko: kernel/net/wireless/lib80211.ko
kernel/drivers/net/wireless/intersil/hostap/hostap_cs.ko: kernel/drivers/net/wireless/intersil/hostap/hostap.ko kernel/net/wireless/lib80211.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/net/wireless/intersil/hostap/hostap_plx.ko: kernel/drivers/net/wireless/intersil/hostap/hostap.ko kernel/net/wireless/lib80211.ko
kernel/drivers/net/wireless/intersil/hostap/hostap_pci.ko: kernel/drivers/net/wireless/intersil/hostap/hostap.ko kernel/net/wireless/lib80211.ko
kernel/drivers/net/wireless/intersil/orinoco/orinoco.ko: kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/intersil/orinoco/orinoco_cs.ko: kernel/drivers/net/wireless/intersil/orinoco/orinoco.ko kernel/net/wireless/cfg80211.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/net/wireless/intersil/orinoco/orinoco_plx.ko: kernel/drivers/net/wireless/intersil/orinoco/orinoco.ko kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/intersil/orinoco/orinoco_tmd.ko: kernel/drivers/net/wireless/intersil/orinoco/orinoco.ko kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/intersil/orinoco/orinoco_nortel.ko: kernel/drivers/net/wireless/intersil/orinoco/orinoco.ko kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/intersil/orinoco/spectrum_cs.ko: kernel/drivers/net/wireless/intersil/orinoco/orinoco.ko kernel/net/wireless/cfg80211.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/net/wireless/intersil/orinoco/orinoco_usb.ko: kernel/drivers/net/wireless/intersil/orinoco/orinoco.ko kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/intersil/p54/p54common.ko: kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/intersil/p54/p54usb.ko: kernel/drivers/net/wireless/intersil/p54/p54common.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/intersil/p54/p54pci.ko: kernel/drivers/net/wireless/intersil/p54/p54common.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/intersil/p54/p54spi.ko: kernel/drivers/net/wireless/intersil/p54/p54common.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/marvell/libertas/libertas.ko: kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/marvell/libertas/usb8xxx.ko: kernel/drivers/net/wireless/marvell/libertas/libertas.ko kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/marvell/libertas/libertas_cs.ko: kernel/drivers/net/wireless/marvell/libertas/libertas.ko kernel/net/wireless/cfg80211.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/net/wireless/marvell/libertas/libertas_sdio.ko: kernel/drivers/net/wireless/marvell/libertas/libertas.ko kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/marvell/libertas/libertas_spi.ko: kernel/drivers/net/wireless/marvell/libertas/libertas.ko kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/marvell/libertas_tf/libertas_tf.ko: kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/marvell/libertas_tf/libertas_tf_usb.ko: kernel/drivers/net/wireless/marvell/libertas_tf/libertas_tf.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/marvell/mwifiex/mwifiex.ko: kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/marvell/mwifiex/mwifiex_sdio.ko: kernel/drivers/net/wireless/marvell/mwifiex/mwifiex.ko kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/marvell/mwifiex/mwifiex_pcie.ko: kernel/drivers/net/wireless/marvell/mwifiex/mwifiex.ko kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/marvell/mwifiex/mwifiex_usb.ko: kernel/drivers/net/wireless/marvell/mwifiex/mwifiex.ko kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/marvell/mwl8k.ko: kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/mediatek/mt7601u/mt7601u.ko: kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/mediatek/mt76/mt76.ko: kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/mediatek/mt76/mt76-usb.ko: kernel/drivers/net/wireless/mediatek/mt76/mt76.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/mediatek/mt76/mt76-sdio.ko: kernel/drivers/net/wireless/mediatek/mt76/mt76.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/mediatek/mt76/mt76x02-lib.ko: kernel/drivers/net/wireless/mediatek/mt76/mt76.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/mediatek/mt76/mt76x02-usb.ko: kernel/drivers/net/wireless/mediatek/mt76/mt76-usb.ko kernel/drivers/net/wireless/mediatek/mt76/mt76x02-lib.ko kernel/drivers/net/wireless/mediatek/mt76/mt76.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/mediatek/mt76/mt76-connac-lib.ko: kernel/drivers/net/wireless/mediatek/mt76/mt76.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/mediatek/mt76/mt76x0/mt76x0u.ko: kernel/drivers/net/wireless/mediatek/mt76/mt76x0/mt76x0-common.ko kernel/drivers/net/wireless/mediatek/mt76/mt76x02-usb.ko kernel/drivers/net/wireless/mediatek/mt76/mt76-usb.ko kernel/drivers/net/wireless/mediatek/mt76/mt76x02-lib.ko kernel/drivers/net/wireless/mediatek/mt76/mt76.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/mediatek/mt76/mt76x0/mt76x0e.ko: kernel/drivers/net/wireless/mediatek/mt76/mt76x0/mt76x0-common.ko kernel/drivers/net/wireless/mediatek/mt76/mt76x02-lib.ko kernel/drivers/net/wireless/mediatek/mt76/mt76.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/mediatek/mt76/mt76x0/mt76x0-common.ko: kernel/drivers/net/wireless/mediatek/mt76/mt76x02-lib.ko kernel/drivers/net/wireless/mediatek/mt76/mt76.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/mediatek/mt76/mt76x2/mt76x2-common.ko: kernel/drivers/net/wireless/mediatek/mt76/mt76x02-lib.ko kernel/drivers/net/wireless/mediatek/mt76/mt76.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/mediatek/mt76/mt76x2/mt76x2e.ko: kernel/drivers/net/wireless/mediatek/mt76/mt76x2/mt76x2-common.ko kernel/drivers/net/wireless/mediatek/mt76/mt76x02-lib.ko kernel/drivers/net/wireless/mediatek/mt76/mt76.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/mediatek/mt76/mt76x2/mt76x2u.ko: kernel/drivers/net/wireless/mediatek/mt76/mt76x2/mt76x2-common.ko kernel/drivers/net/wireless/mediatek/mt76/mt76x02-usb.ko kernel/drivers/net/wireless/mediatek/mt76/mt76-usb.ko kernel/drivers/net/wireless/mediatek/mt76/mt76x02-lib.ko kernel/drivers/net/wireless/mediatek/mt76/mt76.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/mediatek/mt76/mt7603/mt7603e.ko: kernel/drivers/net/wireless/mediatek/mt76/mt76.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/mediatek/mt76/mt7615/mt7615-common.ko: kernel/drivers/net/wireless/mediatek/mt76/mt76-connac-lib.ko kernel/drivers/net/wireless/mediatek/mt76/mt76.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/mediatek/mt76/mt7615/mt7615e.ko: kernel/drivers/net/wireless/mediatek/mt76/mt7615/mt7615-common.ko kernel/drivers/net/wireless/mediatek/mt76/mt76-connac-lib.ko kernel/drivers/net/wireless/mediatek/mt76/mt76.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/mediatek/mt76/mt7615/mt7663-usb-sdio-common.ko: kernel/drivers/net/wireless/mediatek/mt76/mt7615/mt7615-common.ko kernel/drivers/net/wireless/mediatek/mt76/mt76-connac-lib.ko kernel/drivers/net/wireless/mediatek/mt76/mt76.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/mediatek/mt76/mt7615/mt7663u.ko: kernel/drivers/net/wireless/mediatek/mt76/mt7615/mt7663-usb-sdio-common.ko kernel/drivers/net/wireless/mediatek/mt76/mt7615/mt7615-common.ko kernel/drivers/net/wireless/mediatek/mt76/mt76-connac-lib.ko kernel/drivers/net/wireless/mediatek/mt76/mt76-usb.ko kernel/drivers/net/wireless/mediatek/mt76/mt76.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/mediatek/mt76/mt7615/mt7663s.ko: kernel/drivers/net/wireless/mediatek/mt76/mt76-sdio.ko kernel/drivers/net/wireless/mediatek/mt76/mt7615/mt7663-usb-sdio-common.ko kernel/drivers/net/wireless/mediatek/mt76/mt7615/mt7615-common.ko kernel/drivers/net/wireless/mediatek/mt76/mt76-connac-lib.ko kernel/drivers/net/wireless/mediatek/mt76/mt76.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/mediatek/mt76/mt7915/mt7915e.ko: kernel/drivers/net/wireless/mediatek/mt76/mt76.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/mediatek/mt76/mt7921/mt7921e.ko: kernel/drivers/net/wireless/mediatek/mt76/mt76-connac-lib.ko kernel/drivers/net/wireless/mediatek/mt76/mt76.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/microchip/wilc1000/wilc1000.ko: kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/microchip/wilc1000/wilc1000-sdio.ko: kernel/drivers/net/wireless/microchip/wilc1000/wilc1000.ko kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/microchip/wilc1000/wilc1000-spi.ko: kernel/lib/crc7.ko kernel/drivers/net/wireless/microchip/wilc1000/wilc1000.ko kernel/net/wireless/cfg80211.ko kernel/lib/crc-itu-t.ko
kernel/drivers/net/wireless/ralink/rt2x00/rt2x00lib.ko: kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/ralink/rt2x00/rt2x00mmio.ko: kernel/drivers/net/wireless/ralink/rt2x00/rt2x00lib.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/ralink/rt2x00/rt2x00pci.ko: kernel/drivers/net/wireless/ralink/rt2x00/rt2x00lib.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/ralink/rt2x00/rt2x00usb.ko: kernel/drivers/net/wireless/ralink/rt2x00/rt2x00lib.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/ralink/rt2x00/rt2800lib.ko: kernel/drivers/net/wireless/ralink/rt2x00/rt2x00lib.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/ralink/rt2x00/rt2800mmio.ko: kernel/drivers/net/wireless/ralink/rt2x00/rt2800lib.ko kernel/drivers/net/wireless/ralink/rt2x00/rt2x00mmio.ko kernel/drivers/net/wireless/ralink/rt2x00/rt2x00lib.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/ralink/rt2x00/rt2400pci.ko: kernel/drivers/net/wireless/ralink/rt2x00/rt2x00pci.ko kernel/drivers/net/wireless/ralink/rt2x00/rt2x00mmio.ko kernel/drivers/net/wireless/ralink/rt2x00/rt2x00lib.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko kernel/drivers/misc/eeprom/eeprom_93cx6.ko
kernel/drivers/net/wireless/ralink/rt2x00/rt2500pci.ko: kernel/drivers/net/wireless/ralink/rt2x00/rt2x00pci.ko kernel/drivers/net/wireless/ralink/rt2x00/rt2x00mmio.ko kernel/drivers/net/wireless/ralink/rt2x00/rt2x00lib.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko kernel/drivers/misc/eeprom/eeprom_93cx6.ko
kernel/drivers/net/wireless/ralink/rt2x00/rt61pci.ko: kernel/drivers/net/wireless/ralink/rt2x00/rt2x00pci.ko kernel/drivers/net/wireless/ralink/rt2x00/rt2x00mmio.ko kernel/drivers/net/wireless/ralink/rt2x00/rt2x00lib.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko kernel/drivers/misc/eeprom/eeprom_93cx6.ko kernel/lib/crc-itu-t.ko
kernel/drivers/net/wireless/ralink/rt2x00/rt2800pci.ko: kernel/drivers/net/wireless/ralink/rt2x00/rt2800mmio.ko kernel/drivers/net/wireless/ralink/rt2x00/rt2800lib.ko kernel/drivers/net/wireless/ralink/rt2x00/rt2x00pci.ko kernel/drivers/net/wireless/ralink/rt2x00/rt2x00mmio.ko kernel/drivers/net/wireless/ralink/rt2x00/rt2x00lib.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko kernel/drivers/misc/eeprom/eeprom_93cx6.ko
kernel/drivers/net/wireless/ralink/rt2x00/rt2500usb.ko: kernel/drivers/net/wireless/ralink/rt2x00/rt2x00usb.ko kernel/drivers/net/wireless/ralink/rt2x00/rt2x00lib.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/ralink/rt2x00/rt73usb.ko: kernel/drivers/net/wireless/ralink/rt2x00/rt2x00usb.ko kernel/drivers/net/wireless/ralink/rt2x00/rt2x00lib.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko kernel/lib/crc-itu-t.ko
kernel/drivers/net/wireless/ralink/rt2x00/rt2800usb.ko: kernel/drivers/net/wireless/ralink/rt2x00/rt2x00usb.ko kernel/drivers/net/wireless/ralink/rt2x00/rt2800lib.ko kernel/drivers/net/wireless/ralink/rt2x00/rt2x00lib.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/realtek/rtl818x/rtl8180/rtl818x_pci.ko: kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko kernel/drivers/misc/eeprom/eeprom_93cx6.ko
kernel/drivers/net/wireless/realtek/rtl818x/rtl8187/rtl8187.ko: kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko kernel/drivers/misc/eeprom/eeprom_93cx6.ko
kernel/drivers/net/wireless/realtek/rtlwifi/rtlwifi.ko: kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/realtek/rtlwifi/rtl_pci.ko: kernel/drivers/net/wireless/realtek/rtlwifi/rtlwifi.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/realtek/rtlwifi/rtl_usb.ko: kernel/drivers/net/wireless/realtek/rtlwifi/rtlwifi.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/realtek/rtlwifi/rtl8192c/rtl8192c-common.ko: kernel/drivers/net/wireless/realtek/rtlwifi/rtlwifi.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/realtek/rtlwifi/rtl8192ce/rtl8192ce.ko: kernel/drivers/net/wireless/realtek/rtlwifi/rtl_pci.ko kernel/drivers/net/wireless/realtek/rtlwifi/rtl8192c/rtl8192c-common.ko kernel/drivers/net/wireless/realtek/rtlwifi/rtlwifi.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/realtek/rtlwifi/rtl8192cu/rtl8192cu.ko: kernel/drivers/net/wireless/realtek/rtlwifi/rtl_usb.ko kernel/drivers/net/wireless/realtek/rtlwifi/rtl8192c/rtl8192c-common.ko kernel/drivers/net/wireless/realtek/rtlwifi/rtlwifi.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/realtek/rtlwifi/rtl8192se/rtl8192se.ko: kernel/drivers/net/wireless/realtek/rtlwifi/rtl_pci.ko kernel/drivers/net/wireless/realtek/rtlwifi/rtlwifi.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/realtek/rtlwifi/rtl8192de/rtl8192de.ko: kernel/drivers/net/wireless/realtek/rtlwifi/rtl_pci.ko kernel/drivers/net/wireless/realtek/rtlwifi/rtlwifi.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/realtek/rtlwifi/rtl8723ae/rtl8723ae.ko: kernel/drivers/net/wireless/realtek/rtlwifi/btcoexist/btcoexist.ko kernel/drivers/net/wireless/realtek/rtlwifi/rtl8723com/rtl8723-common.ko kernel/drivers/net/wireless/realtek/rtlwifi/rtl_pci.ko kernel/drivers/net/wireless/realtek/rtlwifi/rtlwifi.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/realtek/rtlwifi/rtl8723be/rtl8723be.ko: kernel/drivers/net/wireless/realtek/rtlwifi/btcoexist/btcoexist.ko kernel/drivers/net/wireless/realtek/rtlwifi/rtl8723com/rtl8723-common.ko kernel/drivers/net/wireless/realtek/rtlwifi/rtl_pci.ko kernel/drivers/net/wireless/realtek/rtlwifi/rtlwifi.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/realtek/rtlwifi/rtl8188ee/rtl8188ee.ko: kernel/drivers/net/wireless/realtek/rtlwifi/rtl_pci.ko kernel/drivers/net/wireless/realtek/rtlwifi/rtlwifi.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/realtek/rtlwifi/btcoexist/btcoexist.ko: kernel/drivers/net/wireless/realtek/rtlwifi/rtlwifi.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/realtek/rtlwifi/rtl8723com/rtl8723-common.ko: kernel/drivers/net/wireless/realtek/rtlwifi/rtlwifi.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/realtek/rtlwifi/rtl8821ae/rtl8821ae.ko: kernel/drivers/net/wireless/realtek/rtlwifi/btcoexist/btcoexist.ko kernel/drivers/net/wireless/realtek/rtlwifi/rtl_pci.ko kernel/drivers/net/wireless/realtek/rtlwifi/rtlwifi.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/realtek/rtlwifi/rtl8192ee/rtl8192ee.ko: kernel/drivers/net/wireless/realtek/rtlwifi/btcoexist/btcoexist.ko kernel/drivers/net/wireless/realtek/rtlwifi/rtl_pci.ko kernel/drivers/net/wireless/realtek/rtlwifi/rtlwifi.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/realtek/rtl8xxxu/rtl8xxxu.ko: kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/realtek/rtw88/rtw88_core.ko: kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/realtek/rtw88/rtw88_8822b.ko: kernel/drivers/net/wireless/realtek/rtw88/rtw88_core.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/realtek/rtw88/rtw88_8822be.ko: kernel/drivers/net/wireless/realtek/rtw88/rtw88_8822b.ko kernel/drivers/net/wireless/realtek/rtw88/rtw88_pci.ko kernel/drivers/net/wireless/realtek/rtw88/rtw88_core.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/realtek/rtw88/rtw88_8822c.ko: kernel/drivers/net/wireless/realtek/rtw88/rtw88_core.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/realtek/rtw88/rtw88_8822ce.ko: kernel/drivers/net/wireless/realtek/rtw88/rtw88_8822c.ko kernel/drivers/net/wireless/realtek/rtw88/rtw88_pci.ko kernel/drivers/net/wireless/realtek/rtw88/rtw88_core.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/realtek/rtw88/rtw88_8723d.ko: kernel/drivers/net/wireless/realtek/rtw88/rtw88_core.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/realtek/rtw88/rtw88_8723de.ko: kernel/drivers/net/wireless/realtek/rtw88/rtw88_8723d.ko kernel/drivers/net/wireless/realtek/rtw88/rtw88_pci.ko kernel/drivers/net/wireless/realtek/rtw88/rtw88_core.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/realtek/rtw88/rtw88_8821c.ko: kernel/drivers/net/wireless/realtek/rtw88/rtw88_core.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/realtek/rtw88/rtw88_8821ce.ko: kernel/drivers/net/wireless/realtek/rtw88/rtw88_8821c.ko kernel/drivers/net/wireless/realtek/rtw88/rtw88_pci.ko kernel/drivers/net/wireless/realtek/rtw88/rtw88_core.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/realtek/rtw88/rtw88_pci.ko: kernel/drivers/net/wireless/realtek/rtw88/rtw88_core.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/realtek/rtw89/rtw89_core.ko: kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/realtek/rtw89/rtw89_pci.ko: kernel/drivers/net/wireless/realtek/rtw89/rtw89_core.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/rsi/rsi_91x.ko: kernel/drivers/bluetooth/btrsi.ko kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko kernel/crypto/ecc.ko
kernel/drivers/net/wireless/rsi/rsi_sdio.ko: kernel/drivers/net/wireless/rsi/rsi_91x.ko kernel/drivers/bluetooth/btrsi.ko kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko kernel/crypto/ecc.ko
kernel/drivers/net/wireless/rsi/rsi_usb.ko: kernel/drivers/net/wireless/rsi/rsi_91x.ko kernel/drivers/bluetooth/btrsi.ko kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko kernel/crypto/ecc.ko
kernel/drivers/net/wireless/st/cw1200/cw1200_core.ko: kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/st/cw1200/cw1200_wlan_sdio.ko: kernel/drivers/net/wireless/st/cw1200/cw1200_core.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/st/cw1200/cw1200_wlan_spi.ko: kernel/drivers/net/wireless/st/cw1200/cw1200_core.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/ti/wlcore/wlcore.ko: kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/ti/wlcore/wlcore_sdio.ko:
kernel/drivers/net/wireless/ti/wl12xx/wl12xx.ko: kernel/drivers/net/wireless/ti/wlcore/wlcore.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/ti/wl1251/wl1251.ko: kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/ti/wl1251/wl1251_spi.ko: kernel/drivers/net/wireless/ti/wl1251/wl1251.ko kernel/lib/crc7.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/ti/wl1251/wl1251_sdio.ko: kernel/drivers/net/wireless/ti/wl1251/wl1251.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/ti/wl18xx/wl18xx.ko: kernel/drivers/net/wireless/ti/wlcore/wlcore.ko kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/zydas/zd1211rw/zd1211rw.ko: kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/zydas/zd1201.ko: kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/quantenna/qtnfmac/qtnfmac.ko: kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/quantenna/qtnfmac/qtnfmac_pcie.ko: kernel/drivers/net/wireless/quantenna/qtnfmac/qtnfmac.ko kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wireless/ray_cs.ko: kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/net/wireless/wl3501_cs.ko: kernel/net/wireless/cfg80211.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/net/wireless/rndis_wlan.ko: kernel/drivers/net/usb/rndis_host.ko kernel/drivers/net/usb/cdc_ether.ko kernel/drivers/net/usb/usbnet.ko kernel/net/wireless/cfg80211.ko kernel/drivers/net/mii.ko
kernel/drivers/net/wireless/mac80211_hwsim.ko: kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/net/wireless/virt_wifi.ko: kernel/net/wireless/cfg80211.ko
kernel/drivers/net/wwan/wwan_hwsim.ko:
kernel/drivers/net/wwan/mhi_wwan_ctrl.ko: kernel/drivers/bus/mhi/core/mhi.ko
kernel/drivers/net/wwan/mhi_wwan_mbim.ko: kernel/drivers/bus/mhi/core/mhi.ko
kernel/drivers/net/wwan/rpmsg_wwan_ctrl.ko: kernel/drivers/rpmsg/rpmsg_core.ko
kernel/drivers/net/wwan/iosm/iosm.ko:
kernel/drivers/net/bonding/bonding.ko: kernel/net/tls/tls.ko
kernel/drivers/net/ipvlan/ipvlan.ko:
kernel/drivers/net/ipvlan/ipvtap.ko: kernel/drivers/net/ipvlan/ipvlan.ko kernel/drivers/net/tap.ko
kernel/drivers/net/dummy.ko:
kernel/drivers/net/wireguard/wireguard.ko: kernel/arch/x86/crypto/curve25519-x86_64.ko kernel/lib/crypto/libchacha20poly1305.ko kernel/arch/x86/crypto/chacha-x86_64.ko kernel/arch/x86/crypto/poly1305-x86_64.ko kernel/lib/crypto/libcurve25519-generic.ko kernel/lib/crypto/libchacha.ko kernel/net/ipv6/ip6_udp_tunnel.ko kernel/net/ipv4/udp_tunnel.ko
kernel/drivers/net/eql.ko:
kernel/drivers/net/ifb.ko:
kernel/drivers/net/macsec.ko:
kernel/drivers/net/macvlan.ko:
kernel/drivers/net/macvtap.ko: kernel/drivers/net/macvlan.ko kernel/drivers/net/tap.ko
kernel/drivers/net/mii.ko:
kernel/drivers/net/mdio.ko:
kernel/drivers/net/netconsole.ko:
kernel/drivers/net/rionet.ko:
kernel/drivers/net/team/team.ko:
kernel/drivers/net/team/team_mode_broadcast.ko: kernel/drivers/net/team/team.ko
kernel/drivers/net/team/team_mode_roundrobin.ko: kernel/drivers/net/team/team.ko
kernel/drivers/net/team/team_mode_random.ko: kernel/drivers/net/team/team.ko
kernel/drivers/net/team/team_mode_activebackup.ko: kernel/drivers/net/team/team.ko
kernel/drivers/net/team/team_mode_loadbalance.ko: kernel/drivers/net/team/team.ko
kernel/drivers/net/tap.ko:
kernel/drivers/net/veth.ko:
kernel/drivers/net/virtio_net.ko: kernel/drivers/net/net_failover.ko kernel/net/core/failover.ko
kernel/drivers/net/vxlan.ko: kernel/net/ipv6/ip6_udp_tunnel.ko kernel/net/ipv4/udp_tunnel.ko
kernel/drivers/net/geneve.ko: kernel/net/ipv6/ip6_udp_tunnel.ko kernel/net/ipv4/udp_tunnel.ko
kernel/drivers/net/bareudp.ko: kernel/net/ipv6/ip6_udp_tunnel.ko kernel/net/ipv4/udp_tunnel.ko
kernel/drivers/net/gtp.ko: kernel/net/ipv4/udp_tunnel.ko
kernel/drivers/net/nlmon.ko:
kernel/drivers/net/vrf.ko:
kernel/drivers/net/vsockmon.ko: kernel/net/vmw_vsock/vsock.ko
kernel/drivers/net/mhi_net.ko: kernel/drivers/bus/mhi/core/mhi.ko
kernel/drivers/net/arcnet/arcnet.ko:
kernel/drivers/net/arcnet/rfc1201.ko: kernel/drivers/net/arcnet/arcnet.ko
kernel/drivers/net/arcnet/rfc1051.ko: kernel/drivers/net/arcnet/arcnet.ko
kernel/drivers/net/arcnet/arc-rawmode.ko: kernel/drivers/net/arcnet/arcnet.ko
kernel/drivers/net/arcnet/capmode.ko: kernel/drivers/net/arcnet/arcnet.ko
kernel/drivers/net/arcnet/com90xx.ko: kernel/drivers/net/arcnet/arcnet.ko
kernel/drivers/net/arcnet/com90io.ko: kernel/drivers/net/arcnet/arcnet.ko
kernel/drivers/net/arcnet/arc-rimi.ko: kernel/drivers/net/arcnet/arcnet.ko
kernel/drivers/net/arcnet/com20020.ko: kernel/drivers/net/arcnet/arcnet.ko
kernel/drivers/net/arcnet/com20020-pci.ko: kernel/drivers/net/arcnet/com20020.ko kernel/drivers/net/arcnet/arcnet.ko
kernel/drivers/net/arcnet/com20020_cs.ko: kernel/drivers/net/arcnet/com20020.ko kernel/drivers/net/arcnet/arcnet.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/net/appletalk/ltpc.ko: kernel/net/appletalk/appletalk.ko kernel/net/802/psnap.ko kernel/net/llc/llc.ko
kernel/drivers/net/caif/caif_serial.ko:
kernel/drivers/net/caif/caif_virtio.ko: kernel/drivers/vhost/vringh.ko kernel/drivers/vhost/vhost_iotlb.ko
kernel/drivers/net/can/dev/can-dev.ko:
kernel/drivers/net/can/spi/mcp251xfd/mcp251xfd.ko: kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/spi/hi311x.ko: kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/spi/mcp251x.ko: kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/usb/usb_8dev.ko: kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/usb/ems_usb.ko: kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/usb/esd_usb2.ko: kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/usb/etas_es58x/etas_es58x.ko: kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/usb/gs_usb.ko: kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/usb/kvaser_usb/kvaser_usb.ko: kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/usb/mcba_usb.ko: kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/usb/peak_usb/peak_usb.ko: kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/usb/ucan.ko: kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/softing/softing.ko: kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/softing/softing_cs.ko: kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/net/can/vcan.ko:
kernel/drivers/net/can/vxcan.ko:
kernel/drivers/net/can/slcan.ko:
kernel/drivers/net/can/cc770/cc770.ko: kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/cc770/cc770_isa.ko: kernel/drivers/net/can/cc770/cc770.ko kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/cc770/cc770_platform.ko: kernel/drivers/net/can/cc770/cc770.ko kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/c_can/c_can.ko: kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/c_can/c_can_platform.ko: kernel/drivers/net/can/c_can/c_can.ko kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/c_can/c_can_pci.ko: kernel/drivers/net/can/c_can/c_can.ko kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/ifi_canfd/ifi_canfd.ko: kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/janz-ican3.ko: kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/kvaser_pciefd.ko: kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/m_can/m_can.ko: kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/m_can/m_can_pci.ko: kernel/drivers/net/can/m_can/m_can.ko kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/m_can/m_can_platform.ko: kernel/drivers/net/can/m_can/m_can.ko kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/m_can/tcan4x5x.ko: kernel/drivers/net/can/m_can/m_can.ko kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/peak_canfd/peak_pciefd.ko: kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/sja1000/ems_pci.ko: kernel/drivers/net/can/sja1000/sja1000.ko kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/sja1000/ems_pcmcia.ko: kernel/drivers/net/can/sja1000/sja1000.ko kernel/drivers/net/can/dev/can-dev.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/net/can/sja1000/f81601.ko: kernel/drivers/net/can/sja1000/sja1000.ko kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/sja1000/kvaser_pci.ko: kernel/drivers/net/can/sja1000/sja1000.ko kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/sja1000/peak_pci.ko: kernel/drivers/net/can/sja1000/sja1000.ko kernel/drivers/net/can/dev/can-dev.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/net/can/sja1000/peak_pcmcia.ko: kernel/drivers/net/can/sja1000/sja1000.ko kernel/drivers/net/can/dev/can-dev.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/net/can/sja1000/plx_pci.ko: kernel/drivers/net/can/sja1000/sja1000.ko kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/sja1000/sja1000.ko: kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/sja1000/sja1000_isa.ko: kernel/drivers/net/can/sja1000/sja1000.ko kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/can/sja1000/sja1000_platform.ko: kernel/drivers/net/can/sja1000/sja1000.ko kernel/drivers/net/can/dev/can-dev.ko
kernel/drivers/net/dsa/b53/b53_common.ko: kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/b53/b53_spi.ko: kernel/drivers/net/dsa/b53/b53_common.ko kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/b53/b53_mdio.ko: kernel/drivers/net/dsa/b53/b53_common.ko kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/b53/b53_mmap.ko: kernel/drivers/net/dsa/b53/b53_common.ko kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/b53/b53_srab.ko: kernel/drivers/net/dsa/b53/b53_serdes.ko kernel/drivers/net/dsa/b53/b53_common.ko kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/b53/b53_serdes.ko:
kernel/drivers/net/dsa/hirschmann/hellcreek_sw.ko: kernel/net/sched/sch_taprio.ko kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/microchip/ksz_common.ko: kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/microchip/ksz9477.ko: kernel/drivers/net/dsa/microchip/ksz_common.ko kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/microchip/ksz9477_i2c.ko: kernel/drivers/net/dsa/microchip/ksz9477.ko kernel/drivers/net/dsa/microchip/ksz_common.ko kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/microchip/ksz9477_spi.ko: kernel/drivers/net/dsa/microchip/ksz9477.ko kernel/drivers/net/dsa/microchip/ksz_common.ko kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/microchip/ksz8795.ko: kernel/drivers/net/dsa/microchip/ksz_common.ko kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/microchip/ksz8795_spi.ko: kernel/drivers/net/dsa/microchip/ksz8795.ko kernel/drivers/net/dsa/microchip/ksz_common.ko kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/microchip/ksz8863_smi.ko: kernel/drivers/net/dsa/microchip/ksz8795.ko kernel/drivers/net/dsa/microchip/ksz_common.ko kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/mv88e6xxx/mv88e6xxx.ko: kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/ocelot/mscc_seville.ko: kernel/drivers/net/pcs/pcs-lynx.ko kernel/drivers/net/ethernet/mscc/mscc_ocelot_switch_lib.ko kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/qca/ar9331.ko: kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/sja1105/sja1105.ko: kernel/net/sched/sch_taprio.ko kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/drivers/net/pcs/pcs_xpcs.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/xrs700x/xrs700x.ko: kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/xrs700x/xrs700x_i2c.ko: kernel/drivers/net/dsa/xrs700x/xrs700x.ko kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/xrs700x/xrs700x_mdio.ko: kernel/drivers/net/dsa/xrs700x/xrs700x.ko kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/bcm-sf2.ko: kernel/drivers/net/dsa/b53/b53_common.ko kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/lantiq_gswip.ko: kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/mt7530.ko: kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/mv88e6060.ko: kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/qca8k.ko: kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/realtek-smi.ko: kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/lan9303-core.ko: kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/lan9303_i2c.ko: kernel/drivers/net/dsa/lan9303-core.ko kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/lan9303_mdio.ko: kernel/drivers/net/dsa/lan9303-core.ko kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/vitesse-vsc73xx-core.ko: kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/vitesse-vsc73xx-platform.ko: kernel/drivers/net/dsa/vitesse-vsc73xx-core.ko kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/dsa/vitesse-vsc73xx-spi.ko: kernel/drivers/net/dsa/vitesse-vsc73xx-core.ko kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/drivers/net/plip/plip.ko: kernel/drivers/parport/parport.ko
kernel/drivers/net/sb1000.ko:
kernel/drivers/net/sungem_phy.ko:
kernel/drivers/net/ieee802154/fakelb.ko: kernel/net/mac802154/mac802154.ko kernel/net/ieee802154/ieee802154.ko
kernel/drivers/net/ieee802154/at86rf230.ko: kernel/net/mac802154/mac802154.ko kernel/net/ieee802154/ieee802154.ko
kernel/drivers/net/ieee802154/mrf24j40.ko: kernel/net/mac802154/mac802154.ko kernel/net/ieee802154/ieee802154.ko
kernel/drivers/net/ieee802154/cc2520.ko: kernel/net/mac802154/mac802154.ko kernel/net/ieee802154/ieee802154.ko
kernel/drivers/net/ieee802154/atusb.ko: kernel/net/mac802154/mac802154.ko kernel/net/ieee802154/ieee802154.ko
kernel/drivers/net/ieee802154/adf7242.ko: kernel/net/mac802154/mac802154.ko kernel/net/ieee802154/ieee802154.ko
kernel/drivers/net/ieee802154/ca8210.ko: kernel/net/mac802154/mac802154.ko kernel/net/ieee802154/ieee802154.ko
kernel/drivers/net/ieee802154/mcr20a.ko: kernel/net/mac802154/mac802154.ko kernel/net/ieee802154/ieee802154.ko
kernel/drivers/net/ieee802154/mac802154_hwsim.ko: kernel/net/mac802154/mac802154.ko kernel/net/ieee802154/ieee802154.ko
kernel/drivers/net/vmxnet3/vmxnet3.ko:
kernel/drivers/net/xen-netback/xen-netback.ko:
kernel/drivers/net/usb/catc.ko:
kernel/drivers/net/usb/kaweth.ko:
kernel/drivers/net/usb/pegasus.ko: kernel/drivers/net/mii.ko
kernel/drivers/net/usb/rtl8150.ko:
kernel/drivers/net/usb/r8152.ko: kernel/drivers/net/mii.ko
kernel/drivers/net/usb/hso.ko:
kernel/drivers/net/usb/lan78xx.ko:
kernel/drivers/net/usb/asix.ko: kernel/drivers/net/usb/usbnet.ko kernel/drivers/net/mii.ko
kernel/drivers/net/usb/ax88179_178a.ko: kernel/drivers/net/usb/usbnet.ko kernel/drivers/net/mii.ko
kernel/drivers/net/usb/cdc_ether.ko: kernel/drivers/net/usb/usbnet.ko kernel/drivers/net/mii.ko
kernel/drivers/net/usb/cdc_eem.ko: kernel/drivers/net/usb/usbnet.ko kernel/drivers/net/mii.ko
kernel/drivers/net/usb/dm9601.ko: kernel/drivers/net/usb/usbnet.ko kernel/drivers/net/mii.ko
kernel/drivers/net/usb/sr9700.ko: kernel/drivers/net/usb/usbnet.ko kernel/drivers/net/mii.ko
kernel/drivers/net/usb/sr9800.ko: kernel/drivers/net/usb/usbnet.ko kernel/drivers/net/mii.ko
kernel/drivers/net/usb/smsc75xx.ko: kernel/drivers/net/usb/usbnet.ko kernel/drivers/net/mii.ko
kernel/drivers/net/usb/smsc95xx.ko: kernel/drivers/net/usb/usbnet.ko kernel/drivers/net/mii.ko
kernel/drivers/net/usb/gl620a.ko: kernel/drivers/net/usb/usbnet.ko kernel/drivers/net/mii.ko
kernel/drivers/net/usb/net1080.ko: kernel/drivers/net/usb/usbnet.ko kernel/drivers/net/mii.ko
kernel/drivers/net/usb/plusb.ko: kernel/drivers/net/usb/usbnet.ko kernel/drivers/net/mii.ko
kernel/drivers/net/usb/rndis_host.ko: kernel/drivers/net/usb/cdc_ether.ko kernel/drivers/net/usb/usbnet.ko kernel/drivers/net/mii.ko
kernel/drivers/net/usb/cdc_subset.ko: kernel/drivers/net/usb/usbnet.ko kernel/drivers/net/mii.ko
kernel/drivers/net/usb/zaurus.ko: kernel/drivers/net/usb/cdc_ether.ko kernel/drivers/net/usb/usbnet.ko kernel/drivers/net/mii.ko
kernel/drivers/net/usb/mcs7830.ko: kernel/drivers/net/usb/usbnet.ko kernel/drivers/net/mii.ko
kernel/drivers/net/usb/usbnet.ko: kernel/drivers/net/mii.ko
kernel/drivers/net/usb/int51x1.ko: kernel/drivers/net/usb/usbnet.ko kernel/drivers/net/mii.ko
kernel/drivers/net/usb/cdc-phonet.ko: kernel/net/phonet/phonet.ko
kernel/drivers/net/usb/kalmia.ko: kernel/drivers/net/usb/usbnet.ko kernel/drivers/net/mii.ko
kernel/drivers/net/usb/ipheth.ko:
kernel/drivers/net/usb/sierra_net.ko: kernel/drivers/net/usb/usbnet.ko kernel/drivers/net/mii.ko
kernel/drivers/net/usb/cx82310_eth.ko: kernel/drivers/net/usb/usbnet.ko kernel/drivers/net/mii.ko
kernel/drivers/net/usb/cdc_ncm.ko: kernel/drivers/net/usb/cdc_ether.ko kernel/drivers/net/usb/usbnet.ko kernel/drivers/net/mii.ko
kernel/drivers/net/usb/huawei_cdc_ncm.ko: kernel/drivers/usb/class/cdc-wdm.ko kernel/drivers/net/usb/cdc_ncm.ko kernel/drivers/net/usb/cdc_ether.ko kernel/drivers/net/usb/usbnet.ko kernel/drivers/net/mii.ko
kernel/drivers/net/usb/lg-vl600.ko: kernel/drivers/net/usb/cdc_ether.ko kernel/drivers/net/usb/usbnet.ko kernel/drivers/net/mii.ko
kernel/drivers/net/usb/qmi_wwan.ko: kernel/drivers/usb/class/cdc-wdm.ko kernel/drivers/net/usb/usbnet.ko kernel/drivers/net/mii.ko
kernel/drivers/net/usb/cdc_mbim.ko: kernel/drivers/usb/class/cdc-wdm.ko kernel/drivers/net/usb/cdc_ncm.ko kernel/drivers/net/usb/cdc_ether.ko kernel/drivers/net/usb/usbnet.ko kernel/drivers/net/mii.ko
kernel/drivers/net/usb/ch9200.ko: kernel/drivers/net/usb/usbnet.ko kernel/drivers/net/mii.ko
kernel/drivers/net/usb/aqc111.ko: kernel/drivers/net/usb/usbnet.ko kernel/drivers/net/mii.ko
kernel/drivers/net/usb/r8153_ecm.ko: kernel/drivers/net/usb/r8152.ko kernel/drivers/net/usb/cdc_ether.ko kernel/drivers/net/usb/usbnet.ko kernel/drivers/net/mii.ko
kernel/drivers/net/hyperv/hv_netvsc.ko: kernel/drivers/hv/hv_vmbus.ko
kernel/drivers/net/ntb_netdev.ko: kernel/drivers/ntb/ntb_transport.ko kernel/drivers/ntb/ntb.ko
kernel/drivers/net/fjes/fjes.ko:
kernel/drivers/net/thunderbolt-net.ko: kernel/drivers/thunderbolt/thunderbolt.ko
kernel/drivers/net/netdevsim/netdevsim.ko: kernel/net/psample/psample.ko
kernel/drivers/net/net_failover.ko: kernel/net/core/failover.ko
kernel/drivers/message/fusion/mptbase.ko:
kernel/drivers/message/fusion/mptscsih.ko: kernel/drivers/message/fusion/mptbase.ko
kernel/drivers/message/fusion/mptspi.ko: kernel/drivers/message/fusion/mptscsih.ko kernel/drivers/message/fusion/mptbase.ko kernel/drivers/scsi/scsi_transport_spi.ko
kernel/drivers/message/fusion/mptfc.ko: kernel/drivers/message/fusion/mptscsih.ko kernel/drivers/message/fusion/mptbase.ko kernel/drivers/scsi/scsi_transport_fc.ko
kernel/drivers/message/fusion/mptsas.ko: kernel/drivers/message/fusion/mptscsih.ko kernel/drivers/message/fusion/mptbase.ko kernel/drivers/scsi/scsi_transport_sas.ko
kernel/drivers/message/fusion/mptctl.ko: kernel/drivers/message/fusion/mptbase.ko
kernel/drivers/message/fusion/mptlan.ko: kernel/drivers/message/fusion/mptbase.ko
kernel/drivers/firewire/firewire-core.ko: kernel/lib/crc-itu-t.ko
kernel/drivers/firewire/firewire-ohci.ko: kernel/drivers/firewire/firewire-core.ko kernel/lib/crc-itu-t.ko
kernel/drivers/firewire/firewire-sbp2.ko: kernel/drivers/firewire/firewire-core.ko kernel/lib/crc-itu-t.ko
kernel/drivers/firewire/firewire-net.ko: kernel/drivers/firewire/firewire-core.ko kernel/lib/crc-itu-t.ko
kernel/drivers/firewire/nosy.ko:
kernel/drivers/vfio/mdev/mdev.ko:
kernel/drivers/auxdisplay/charlcd.ko:
kernel/drivers/auxdisplay/hd44780_common.ko: kernel/drivers/auxdisplay/charlcd.ko
kernel/drivers/auxdisplay/ks0108.ko: kernel/drivers/parport/parport.ko
kernel/drivers/auxdisplay/cfag12864b.ko: kernel/drivers/auxdisplay/ks0108.ko kernel/drivers/parport/parport.ko
kernel/drivers/auxdisplay/cfag12864bfb.ko: kernel/drivers/auxdisplay/cfag12864b.ko kernel/drivers/auxdisplay/ks0108.ko kernel/drivers/parport/parport.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/auxdisplay/img-ascii-lcd.ko:
kernel/drivers/auxdisplay/hd44780.ko: kernel/drivers/auxdisplay/hd44780_common.ko kernel/drivers/auxdisplay/charlcd.ko
kernel/drivers/auxdisplay/panel.ko: kernel/drivers/auxdisplay/hd44780_common.ko kernel/drivers/auxdisplay/charlcd.ko kernel/drivers/parport/parport.ko
kernel/drivers/auxdisplay/lcd2s.ko: kernel/drivers/auxdisplay/charlcd.ko
kernel/drivers/usb/common/usb-conn-gpio.ko:
kernel/drivers/usb/common/ulpi.ko:
kernel/drivers/usb/core/ledtrig-usbport.ko:
kernel/drivers/usb/phy/phy-generic.ko:
kernel/drivers/usb/phy/phy-tahvo.ko: kernel/drivers/mfd/retu-mfd.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/phy/phy-gpio-vbus-usb.ko: kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/phy/phy-isp1301.ko:
kernel/drivers/usb/dwc2/dwc2_pci.ko: kernel/drivers/usb/phy/phy-generic.ko
kernel/drivers/usb/host/oxu210hp-hcd.ko:
kernel/drivers/usb/host/isp116x-hcd.ko:
kernel/drivers/usb/host/xhci-pci.ko: kernel/drivers/usb/host/xhci-pci-renesas.ko
kernel/drivers/usb/host/xhci-pci-renesas.ko:
kernel/drivers/usb/host/xhci-plat-hcd.ko:
kernel/drivers/usb/host/sl811-hcd.ko:
kernel/drivers/usb/host/sl811_cs.ko: kernel/drivers/usb/host/sl811-hcd.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/usb/host/u132-hcd.ko: kernel/drivers/usb/misc/ftdi-elan.ko
kernel/drivers/usb/host/r8a66597-hcd.ko:
kernel/drivers/usb/host/fsl-mph-dr-of.ko:
kernel/drivers/usb/host/ehci-fsl.ko:
kernel/drivers/usb/host/bcma-hcd.ko: kernel/drivers/bcma/bcma.ko
kernel/drivers/usb/host/ssb-hcd.ko: kernel/drivers/ssb/ssb.ko
kernel/drivers/usb/host/fotg210-hcd.ko:
kernel/drivers/usb/host/max3421-hcd.ko:
kernel/drivers/usb/storage/uas.ko: kernel/drivers/usb/storage/usb-storage.ko
kernel/drivers/usb/storage/usb-storage.ko:
kernel/drivers/usb/storage/ums-alauda.ko: kernel/drivers/usb/storage/usb-storage.ko
kernel/drivers/usb/storage/ums-cypress.ko: kernel/drivers/usb/storage/usb-storage.ko
kernel/drivers/usb/storage/ums-datafab.ko: kernel/drivers/usb/storage/usb-storage.ko
kernel/drivers/usb/storage/ums-eneub6250.ko: kernel/drivers/usb/storage/usb-storage.ko
kernel/drivers/usb/storage/ums-freecom.ko: kernel/drivers/usb/storage/usb-storage.ko
kernel/drivers/usb/storage/ums-isd200.ko: kernel/drivers/usb/storage/usb-storage.ko
kernel/drivers/usb/storage/ums-jumpshot.ko: kernel/drivers/usb/storage/usb-storage.ko
kernel/drivers/usb/storage/ums-karma.ko: kernel/drivers/usb/storage/usb-storage.ko
kernel/drivers/usb/storage/ums-onetouch.ko: kernel/drivers/usb/storage/usb-storage.ko
kernel/drivers/usb/storage/ums-realtek.ko: kernel/drivers/usb/storage/usb-storage.ko
kernel/drivers/usb/storage/ums-sddr09.ko: kernel/drivers/usb/storage/usb-storage.ko
kernel/drivers/usb/storage/ums-sddr55.ko: kernel/drivers/usb/storage/usb-storage.ko
kernel/drivers/usb/storage/ums-usbat.ko: kernel/drivers/usb/storage/usb-storage.ko
kernel/drivers/usb/misc/adutux.ko:
kernel/drivers/usb/misc/appledisplay.ko:
kernel/drivers/usb/misc/cypress_cy7c63.ko:
kernel/drivers/usb/misc/cytherm.ko:
kernel/drivers/usb/misc/emi26.ko:
kernel/drivers/usb/misc/emi62.ko:
kernel/drivers/usb/misc/ezusb.ko:
kernel/drivers/usb/misc/ftdi-elan.ko:
kernel/drivers/usb/misc/apple-mfi-fastcharge.ko:
kernel/drivers/usb/misc/idmouse.ko:
kernel/drivers/usb/misc/iowarrior.ko:
kernel/drivers/usb/misc/isight_firmware.ko:
kernel/drivers/usb/misc/usblcd.ko:
kernel/drivers/usb/misc/ldusb.ko:
kernel/drivers/usb/misc/legousbtower.ko:
kernel/drivers/usb/misc/usbtest.ko:
kernel/drivers/usb/misc/ehset.ko:
kernel/drivers/usb/misc/trancevibrator.ko:
kernel/drivers/usb/misc/uss720.ko: kernel/drivers/parport/parport.ko
kernel/drivers/usb/misc/usbsevseg.ko:
kernel/drivers/usb/misc/yurex.ko:
kernel/drivers/usb/misc/usb251xb.ko:
kernel/drivers/usb/misc/usb3503.ko:
kernel/drivers/usb/misc/usb4604.ko:
kernel/drivers/usb/misc/chaoskey.ko:
kernel/drivers/usb/misc/sisusbvga/sisusbvga.ko:
kernel/drivers/usb/misc/lvstest.ko:
kernel/drivers/usb/roles/intel-xhci-usb-role-switch.ko:
kernel/drivers/usb/dwc3/dwc3.ko: kernel/drivers/usb/common/ulpi.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/dwc3/dwc3-pci.ko:
kernel/drivers/usb/dwc3/dwc3-haps.ko:
kernel/drivers/usb/isp1760/isp1760.ko: kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/cdns3/cdns-usb-common.ko:
kernel/drivers/usb/cdns3/cdns3.ko: kernel/drivers/usb/cdns3/cdns-usb-common.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/cdns3/cdns3-pci-wrap.ko:
kernel/drivers/usb/cdns3/cdnsp-udc-pci.ko: kernel/drivers/usb/cdns3/cdns-usb-common.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/mon/usbmon.ko:
kernel/drivers/usb/c67x00/c67x00.ko:
kernel/drivers/usb/class/cdc-acm.ko:
kernel/drivers/usb/class/usblp.ko:
kernel/drivers/usb/class/cdc-wdm.ko:
kernel/drivers/usb/class/usbtmc.ko:
kernel/drivers/usb/image/mdc800.ko:
kernel/drivers/usb/image/microtek.ko:
kernel/drivers/usb/serial/usbserial.ko:
kernel/drivers/usb/serial/aircable.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/ark3116.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/belkin_sa.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/ch341.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/cp210x.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/cyberjack.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/cypress_m8.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/usb_debug.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/digi_acceleport.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/io_edgeport.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/io_ti.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/empeg.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/f81232.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/f81534.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/ftdi_sio.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/garmin_gps.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/ipaq.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/ipw.ko: kernel/drivers/usb/serial/usb_wwan.ko kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/ir-usb.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/iuu_phoenix.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/keyspan.ko: kernel/drivers/usb/misc/ezusb.ko kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/keyspan_pda.ko: kernel/drivers/usb/misc/ezusb.ko kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/kl5kusb105.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/kobil_sct.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/mct_u232.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/metro-usb.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/mos7720.ko: kernel/drivers/usb/serial/usbserial.ko kernel/drivers/parport/parport.ko
kernel/drivers/usb/serial/mos7840.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/mxuport.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/navman.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/omninet.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/opticon.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/option.ko: kernel/drivers/usb/serial/usb_wwan.ko kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/oti6858.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/pl2303.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/qcaux.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/qcserial.ko: kernel/drivers/usb/serial/usb_wwan.ko kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/quatech2.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/safe_serial.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/sierra.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/usb-serial-simple.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/spcp8x5.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/ssu100.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/symbolserial.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/usb_wwan.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/ti_usb_3410_5052.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/upd78f0730.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/visor.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/wishbone-serial.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/whiteheat.ko: kernel/drivers/usb/misc/ezusb.ko kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/xr_serial.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/serial/xsens_mt.ko: kernel/drivers/usb/serial/usbserial.ko
kernel/drivers/usb/atm/cxacru.ko: kernel/drivers/usb/atm/usbatm.ko kernel/net/atm/atm.ko
kernel/drivers/usb/atm/speedtch.ko: kernel/drivers/usb/atm/usbatm.ko kernel/net/atm/atm.ko
kernel/drivers/usb/atm/ueagle-atm.ko: kernel/drivers/usb/atm/usbatm.ko kernel/net/atm/atm.ko
kernel/drivers/usb/atm/usbatm.ko: kernel/net/atm/atm.ko
kernel/drivers/usb/atm/xusbatm.ko: kernel/drivers/usb/atm/usbatm.ko kernel/net/atm/atm.ko
kernel/drivers/usb/musb/musb_hdrc.ko: kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/chipidea/ci_hdrc.ko: kernel/drivers/usb/common/ulpi.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/chipidea/ci_hdrc_usb2.ko: kernel/drivers/usb/chipidea/ci_hdrc.ko kernel/drivers/usb/common/ulpi.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/chipidea/ci_hdrc_msm.ko: kernel/drivers/usb/chipidea/ci_hdrc.ko kernel/drivers/usb/common/ulpi.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/chipidea/ci_hdrc_pci.ko: kernel/drivers/usb/chipidea/ci_hdrc.ko kernel/drivers/usb/phy/phy-generic.ko kernel/drivers/usb/common/ulpi.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/libcomposite.ko: kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/udc/udc-core.ko:
kernel/drivers/usb/gadget/udc/net2272.ko: kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/udc/net2280.ko: kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/udc/snps_udc_core.ko: kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/udc/amd5536udc_pci.ko: kernel/drivers/usb/gadget/udc/snps_udc_core.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/udc/pxa27x_udc.ko: kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/udc/goku_udc.ko: kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/udc/r8a66597-udc.ko: kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/udc/pch_udc.ko: kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/udc/mv_udc.ko: kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/udc/fotg210-udc.ko: kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/udc/mv_u3d_core.ko: kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/udc/gr_udc.ko: kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/udc/bdc/bdc.ko: kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/udc/max3420_udc.ko: kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/function/usb_f_acm.ko: kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/function/u_serial.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/function/usb_f_ss_lb.ko: kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/function/u_serial.ko: kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/function/usb_f_serial.ko: kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/function/u_serial.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/function/usb_f_obex.ko: kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/function/u_serial.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/function/u_ether.ko: kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/function/usb_f_ncm.ko: kernel/drivers/usb/gadget/function/u_ether.ko kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/function/usb_f_ecm.ko: kernel/drivers/usb/gadget/function/u_ether.ko kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/function/usb_f_phonet.ko: kernel/drivers/usb/gadget/function/u_ether.ko kernel/drivers/usb/gadget/libcomposite.ko kernel/net/phonet/phonet.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/function/usb_f_eem.ko: kernel/drivers/usb/gadget/function/u_ether.ko kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/function/usb_f_ecm_subset.ko: kernel/drivers/usb/gadget/function/u_ether.ko kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/function/usb_f_rndis.ko: kernel/drivers/usb/gadget/function/u_ether.ko kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/function/usb_f_mass_storage.ko: kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/function/usb_f_fs.ko: kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/function/u_audio.ko: kernel/drivers/usb/gadget/libcomposite.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/function/usb_f_uac1.ko: kernel/drivers/usb/gadget/function/u_audio.ko kernel/drivers/usb/gadget/libcomposite.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/function/usb_f_uac1_legacy.ko: kernel/drivers/usb/gadget/libcomposite.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/function/usb_f_uac2.ko: kernel/drivers/usb/gadget/function/u_audio.ko kernel/drivers/usb/gadget/libcomposite.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/function/usb_f_uvc.ko: kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-dma-sg.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/function/usb_f_midi.ko: kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/drivers/usb/gadget/libcomposite.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/function/usb_f_hid.ko: kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/function/usb_f_printer.ko: kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/function/usb_f_tcm.ko: kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/target/target_core_mod.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/legacy/g_zero.ko: kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/legacy/g_audio.ko: kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/legacy/g_ether.ko: kernel/drivers/usb/gadget/function/usb_f_rndis.ko kernel/drivers/usb/gadget/function/u_ether.ko kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/legacy/gadgetfs.ko: kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/legacy/g_ffs.ko: kernel/drivers/usb/gadget/function/usb_f_fs.ko kernel/drivers/usb/gadget/function/usb_f_rndis.ko kernel/drivers/usb/gadget/function/u_ether.ko kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/legacy/g_mass_storage.ko: kernel/drivers/usb/gadget/function/usb_f_mass_storage.ko kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/legacy/g_serial.ko: kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/legacy/g_printer.ko: kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/legacy/g_midi.ko: kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/legacy/g_cdc.ko: kernel/drivers/usb/gadget/function/u_ether.ko kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/legacy/g_hid.ko: kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/legacy/g_dbgp.ko: kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/function/u_serial.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/legacy/g_nokia.ko: kernel/drivers/usb/gadget/function/usb_f_mass_storage.ko kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/legacy/g_webcam.ko: kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/legacy/g_ncm.ko: kernel/drivers/usb/gadget/function/u_ether.ko kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/legacy/g_acm_ms.ko: kernel/drivers/usb/gadget/function/usb_f_mass_storage.ko kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/legacy/tcm_usb_gadget.ko: kernel/drivers/usb/gadget/libcomposite.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/gadget/legacy/raw_gadget.ko: kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/usbip/usbip-core.ko:
kernel/drivers/usb/usbip/vhci-hcd.ko: kernel/drivers/usb/usbip/usbip-core.ko
kernel/drivers/usb/usbip/usbip-host.ko: kernel/drivers/usb/usbip/usbip-core.ko
kernel/drivers/usb/usbip/usbip-vudc.ko: kernel/drivers/usb/usbip/usbip-core.ko kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/usb/typec/typec.ko:
kernel/drivers/usb/typec/altmodes/typec_displayport.ko: kernel/drivers/usb/typec/typec.ko
kernel/drivers/usb/typec/altmodes/typec_nvidia.ko: kernel/drivers/usb/typec/altmodes/typec_displayport.ko kernel/drivers/usb/typec/typec.ko
kernel/drivers/usb/typec/tcpm/tcpm.ko: kernel/drivers/usb/typec/typec.ko
kernel/drivers/usb/typec/tcpm/fusb302.ko: kernel/drivers/usb/typec/tcpm/tcpm.ko kernel/drivers/usb/typec/typec.ko
kernel/drivers/usb/typec/tcpm/tcpci.ko: kernel/drivers/usb/typec/tcpm/tcpm.ko kernel/drivers/usb/typec/typec.ko
kernel/drivers/usb/typec/tcpm/tcpci_rt1711h.ko: kernel/drivers/usb/typec/tcpm/tcpci.ko kernel/drivers/usb/typec/tcpm/tcpm.ko kernel/drivers/usb/typec/typec.ko
kernel/drivers/usb/typec/tcpm/tcpci_mt6360.ko: kernel/drivers/usb/typec/tcpm/tcpci.ko kernel/drivers/usb/typec/tcpm/tcpm.ko kernel/drivers/usb/typec/typec.ko
kernel/drivers/usb/typec/tcpm/tcpci_maxim.ko: kernel/drivers/usb/typec/tcpm/tcpci.ko kernel/drivers/usb/typec/tcpm/tcpm.ko kernel/drivers/usb/typec/typec.ko
kernel/drivers/usb/typec/ucsi/typec_ucsi.ko: kernel/drivers/usb/typec/typec.ko
kernel/drivers/usb/typec/ucsi/ucsi_acpi.ko: kernel/drivers/usb/typec/ucsi/typec_ucsi.ko kernel/drivers/usb/typec/typec.ko
kernel/drivers/usb/typec/ucsi/ucsi_ccg.ko: kernel/drivers/usb/typec/ucsi/typec_ucsi.ko kernel/drivers/usb/typec/typec.ko
kernel/drivers/usb/typec/tipd/tps6598x.ko: kernel/drivers/usb/typec/typec.ko
kernel/drivers/usb/typec/hd3ss3220.ko: kernel/drivers/usb/typec/typec.ko
kernel/drivers/usb/typec/stusb160x.ko: kernel/drivers/usb/typec/typec.ko
kernel/drivers/usb/typec/mux/pi3usb30532.ko: kernel/drivers/usb/typec/typec.ko
kernel/drivers/usb/typec/mux/intel_pmc_mux.ko: kernel/drivers/usb/typec/typec.ko
kernel/drivers/input/serio/parkbd.ko: kernel/drivers/parport/parport.ko
kernel/drivers/input/serio/serport.ko:
kernel/drivers/input/serio/ct82c710.ko:
kernel/drivers/input/serio/pcips2.ko:
kernel/drivers/input/serio/ps2mult.ko:
kernel/drivers/input/serio/serio_raw.ko:
kernel/drivers/input/serio/altera_ps2.ko:
kernel/drivers/input/serio/arc_ps2.ko:
kernel/drivers/input/serio/hyperv-keyboard.ko: kernel/drivers/hv/hv_vmbus.ko
kernel/drivers/input/serio/ps2-gpio.ko:
kernel/drivers/input/serio/userio.ko:
kernel/drivers/input/keyboard/adc-keys.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/input/keyboard/adp5520-keys.ko:
kernel/drivers/input/keyboard/adp5588-keys.ko:
kernel/drivers/input/keyboard/adp5589-keys.ko:
kernel/drivers/input/keyboard/applespi.ko:
kernel/drivers/input/keyboard/cros_ec_keyb.ko: kernel/drivers/input/matrix-keymap.ko
kernel/drivers/input/keyboard/dlink-dir685-touchkeys.ko:
kernel/drivers/input/keyboard/gpio_keys.ko:
kernel/drivers/input/keyboard/gpio_keys_polled.ko:
kernel/drivers/input/keyboard/tca6416-keypad.ko:
kernel/drivers/input/keyboard/tca8418_keypad.ko: kernel/drivers/input/matrix-keymap.ko
kernel/drivers/input/keyboard/iqs62x-keys.ko: kernel/drivers/mfd/iqs62x.ko
kernel/drivers/input/keyboard/lkkbd.ko:
kernel/drivers/input/keyboard/lm8323.ko:
kernel/drivers/input/keyboard/lm8333.ko: kernel/drivers/input/matrix-keymap.ko
kernel/drivers/input/keyboard/matrix_keypad.ko: kernel/drivers/input/matrix-keymap.ko
kernel/drivers/input/keyboard/max7359_keypad.ko: kernel/drivers/input/matrix-keymap.ko
kernel/drivers/input/keyboard/mcs_touchkey.ko:
kernel/drivers/input/keyboard/mpr121_touchkey.ko:
kernel/drivers/input/keyboard/mtk-pmic-keys.ko:
kernel/drivers/input/keyboard/newtonkbd.ko:
kernel/drivers/input/keyboard/opencores-kbd.ko:
kernel/drivers/input/keyboard/qt1050.ko:
kernel/drivers/input/keyboard/qt1070.ko:
kernel/drivers/input/keyboard/qt2160.ko:
kernel/drivers/input/keyboard/samsung-keypad.ko: kernel/drivers/input/matrix-keymap.ko
kernel/drivers/input/keyboard/stowaway.ko:
kernel/drivers/input/keyboard/sunkbd.ko:
kernel/drivers/input/keyboard/tm2-touchkey.ko:
kernel/drivers/input/keyboard/twl4030_keypad.ko: kernel/drivers/input/matrix-keymap.ko
kernel/drivers/input/keyboard/xtkbd.ko:
kernel/drivers/input/mouse/appletouch.ko:
kernel/drivers/input/mouse/bcm5974.ko:
kernel/drivers/input/mouse/cyapatp.ko: kernel/lib/crc-itu-t.ko
kernel/drivers/input/mouse/elan_i2c.ko:
kernel/drivers/input/mouse/gpio_mouse.ko:
kernel/drivers/input/mouse/psmouse.ko:
kernel/drivers/input/mouse/sermouse.ko:
kernel/drivers/input/mouse/synaptics_i2c.ko:
kernel/drivers/input/mouse/synaptics_usb.ko:
kernel/drivers/input/mouse/vsxxxaa.ko:
kernel/drivers/input/joystick/a3d.ko: kernel/drivers/input/gameport/gameport.ko
kernel/drivers/input/joystick/adc-joystick.ko: kernel/drivers/iio/buffer/industrialio-buffer-cb.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/input/joystick/adi.ko: kernel/drivers/input/gameport/gameport.ko
kernel/drivers/input/joystick/as5011.ko:
kernel/drivers/input/joystick/analog.ko: kernel/drivers/input/gameport/gameport.ko
kernel/drivers/input/joystick/cobra.ko: kernel/drivers/input/gameport/gameport.ko
kernel/drivers/input/joystick/db9.ko: kernel/drivers/parport/parport.ko
kernel/drivers/input/joystick/fsia6b.ko:
kernel/drivers/input/joystick/gamecon.ko: kernel/drivers/input/ff-memless.ko kernel/drivers/parport/parport.ko
kernel/drivers/input/joystick/gf2k.ko: kernel/drivers/input/gameport/gameport.ko
kernel/drivers/input/joystick/grip.ko: kernel/drivers/input/gameport/gameport.ko
kernel/drivers/input/joystick/grip_mp.ko: kernel/drivers/input/gameport/gameport.ko
kernel/drivers/input/joystick/guillemot.ko: kernel/drivers/input/gameport/gameport.ko
kernel/drivers/input/joystick/iforce/iforce.ko:
kernel/drivers/input/joystick/iforce/iforce-serio.ko: kernel/drivers/input/joystick/iforce/iforce.ko
kernel/drivers/input/joystick/iforce/iforce-usb.ko: kernel/drivers/input/joystick/iforce/iforce.ko
kernel/drivers/input/joystick/interact.ko: kernel/drivers/input/gameport/gameport.ko
kernel/drivers/input/joystick/joydump.ko: kernel/drivers/input/gameport/gameport.ko
kernel/drivers/input/joystick/magellan.ko:
kernel/drivers/input/joystick/psxpad-spi.ko: kernel/drivers/input/ff-memless.ko
kernel/drivers/input/joystick/pxrc.ko:
kernel/drivers/input/joystick/qwiic-joystick.ko:
kernel/drivers/input/joystick/sidewinder.ko: kernel/drivers/input/gameport/gameport.ko
kernel/drivers/input/joystick/spaceball.ko:
kernel/drivers/input/joystick/spaceorb.ko:
kernel/drivers/input/joystick/stinger.ko:
kernel/drivers/input/joystick/tmdc.ko: kernel/drivers/input/gameport/gameport.ko
kernel/drivers/input/joystick/turbografx.ko: kernel/drivers/parport/parport.ko
kernel/drivers/input/joystick/twidjoy.ko:
kernel/drivers/input/joystick/warrior.ko:
kernel/drivers/input/joystick/walkera0701.ko: kernel/drivers/parport/parport.ko
kernel/drivers/input/joystick/xpad.ko: kernel/drivers/input/ff-memless.ko
kernel/drivers/input/joystick/zhenhua.ko:
kernel/drivers/input/tablet/acecad.ko:
kernel/drivers/input/tablet/aiptek.ko:
kernel/drivers/input/tablet/hanwang.ko:
kernel/drivers/input/tablet/kbtab.ko:
kernel/drivers/input/tablet/pegasus_notetaker.ko:
kernel/drivers/input/tablet/wacom_serial4.ko:
kernel/drivers/input/touchscreen/88pm860x-ts.ko:
kernel/drivers/input/touchscreen/ad7877.ko:
kernel/drivers/input/touchscreen/ad7879.ko:
kernel/drivers/input/touchscreen/ad7879-i2c.ko: kernel/drivers/input/touchscreen/ad7879.ko
kernel/drivers/input/touchscreen/ad7879-spi.ko: kernel/drivers/input/touchscreen/ad7879.ko
kernel/drivers/input/touchscreen/resistive-adc-touch.ko: kernel/drivers/iio/buffer/industrialio-buffer-cb.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/input/touchscreen/ads7846.ko:
kernel/drivers/input/touchscreen/atmel_mxt_ts.ko: kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/input/touchscreen/auo-pixcir-ts.ko:
kernel/drivers/input/touchscreen/bu21013_ts.ko:
kernel/drivers/input/touchscreen/bu21029_ts.ko:
kernel/drivers/input/touchscreen/chipone_icn8505.ko:
kernel/drivers/input/touchscreen/cy8ctma140.ko:
kernel/drivers/input/touchscreen/cy8ctmg110_ts.ko:
kernel/drivers/input/touchscreen/cyttsp_core.ko:
kernel/drivers/input/touchscreen/cyttsp_i2c.ko: kernel/drivers/input/touchscreen/cyttsp_i2c_common.ko kernel/drivers/input/touchscreen/cyttsp_core.ko
kernel/drivers/input/touchscreen/cyttsp_i2c_common.ko:
kernel/drivers/input/touchscreen/cyttsp_spi.ko: kernel/drivers/input/touchscreen/cyttsp_core.ko
kernel/drivers/input/touchscreen/cyttsp4_core.ko:
kernel/drivers/input/touchscreen/cyttsp4_i2c.ko: kernel/drivers/input/touchscreen/cyttsp4_core.ko kernel/drivers/input/touchscreen/cyttsp_i2c_common.ko
kernel/drivers/input/touchscreen/cyttsp4_spi.ko: kernel/drivers/input/touchscreen/cyttsp4_core.ko
kernel/drivers/input/touchscreen/da9034-ts.ko:
kernel/drivers/input/touchscreen/da9052_tsi.ko:
kernel/drivers/input/touchscreen/dynapro.ko:
kernel/drivers/input/touchscreen/edt-ft5x06.ko:
kernel/drivers/input/touchscreen/hampshire.ko:
kernel/drivers/input/touchscreen/hycon-hy46xx.ko:
kernel/drivers/input/touchscreen/gunze.ko:
kernel/drivers/input/touchscreen/eeti_ts.ko:
kernel/drivers/input/touchscreen/ektf2127.ko:
kernel/drivers/input/touchscreen/elo.ko:
kernel/drivers/input/touchscreen/egalax_ts_serial.ko:
kernel/drivers/input/touchscreen/exc3000.ko:
kernel/drivers/input/touchscreen/fujitsu_ts.ko:
kernel/drivers/input/touchscreen/goodix.ko:
kernel/drivers/input/touchscreen/hideep.ko:
kernel/drivers/input/touchscreen/ili210x.ko:
kernel/drivers/input/touchscreen/ilitek_ts_i2c.ko:
kernel/drivers/input/touchscreen/inexio.ko:
kernel/drivers/input/touchscreen/max11801_ts.ko:
kernel/drivers/input/touchscreen/mc13783_ts.ko: kernel/drivers/mfd/mc13xxx-core.ko
kernel/drivers/input/touchscreen/mcs5000_ts.ko:
kernel/drivers/input/touchscreen/melfas_mip4.ko:
kernel/drivers/input/touchscreen/mms114.ko:
kernel/drivers/input/touchscreen/msg2638.ko:
kernel/drivers/input/touchscreen/mtouch.ko:
kernel/drivers/input/touchscreen/mk712.ko:
kernel/drivers/input/touchscreen/usbtouchscreen.ko:
kernel/drivers/input/touchscreen/pcap_ts.ko:
kernel/drivers/input/touchscreen/penmount.ko:
kernel/drivers/input/touchscreen/pixcir_i2c_ts.ko:
kernel/drivers/input/touchscreen/raydium_i2c_ts.ko:
kernel/drivers/input/touchscreen/s6sy761.ko:
kernel/drivers/input/touchscreen/silead.ko:
kernel/drivers/input/touchscreen/sis_i2c.ko: kernel/lib/crc-itu-t.ko
kernel/drivers/input/touchscreen/st1232.ko:
kernel/drivers/input/touchscreen/stmfts.ko:
kernel/drivers/input/touchscreen/sur40.ko: kernel/drivers/media/common/videobuf2/videobuf2-dma-sg.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/input/touchscreen/surface3_spi.ko:
kernel/drivers/input/touchscreen/ti_am335x_tsc.ko: kernel/drivers/mfd/ti_am335x_tscadc.ko
kernel/drivers/input/touchscreen/touchit213.ko:
kernel/drivers/input/touchscreen/touchright.ko:
kernel/drivers/input/touchscreen/touchwin.ko:
kernel/drivers/input/touchscreen/tsc40.ko:
kernel/drivers/input/touchscreen/tsc200x-core.ko:
kernel/drivers/input/touchscreen/tsc2004.ko: kernel/drivers/input/touchscreen/tsc200x-core.ko
kernel/drivers/input/touchscreen/tsc2005.ko: kernel/drivers/input/touchscreen/tsc200x-core.ko
kernel/drivers/input/touchscreen/tsc2007.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/input/touchscreen/ucb1400_ts.ko: kernel/drivers/mfd/ucb1400_core.ko kernel/sound/ac97_bus.ko
kernel/drivers/input/touchscreen/wacom_w8001.ko:
kernel/drivers/input/touchscreen/wacom_i2c.ko:
kernel/drivers/input/touchscreen/wdt87xx_i2c.ko:
kernel/drivers/input/touchscreen/wm831x-ts.ko:
kernel/drivers/input/touchscreen/wm97xx-ts.ko:
kernel/drivers/input/touchscreen/sx8654.ko:
kernel/drivers/input/touchscreen/tps6507x-ts.ko:
kernel/drivers/input/touchscreen/zet6223.ko:
kernel/drivers/input/touchscreen/zforce_ts.ko:
kernel/drivers/input/touchscreen/rohm_bu21023.ko:
kernel/drivers/input/touchscreen/iqs5xx.ko:
kernel/drivers/input/touchscreen/zinitix.ko:
kernel/drivers/input/misc/88pm860x_onkey.ko:
kernel/drivers/input/misc/88pm80x_onkey.ko:
kernel/drivers/input/misc/ad714x.ko:
kernel/drivers/input/misc/ad714x-i2c.ko: kernel/drivers/input/misc/ad714x.ko
kernel/drivers/input/misc/ad714x-spi.ko: kernel/drivers/input/misc/ad714x.ko
kernel/drivers/input/misc/adxl34x.ko:
kernel/drivers/input/misc/adxl34x-i2c.ko: kernel/drivers/input/misc/adxl34x.ko
kernel/drivers/input/misc/adxl34x-spi.ko: kernel/drivers/input/misc/adxl34x.ko
kernel/drivers/input/misc/apanel.ko:
kernel/drivers/input/misc/arizona-haptics.ko: kernel/drivers/input/ff-memless.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/drivers/input/misc/atc260x-onkey.ko:
kernel/drivers/input/misc/ati_remote2.ko:
kernel/drivers/input/misc/atlas_btns.ko:
kernel/drivers/input/misc/bma150.ko:
kernel/drivers/input/misc/cm109.ko:
kernel/drivers/input/misc/cma3000_d0x.ko:
kernel/drivers/input/misc/cma3000_d0x_i2c.ko: kernel/drivers/input/misc/cma3000_d0x.ko
kernel/drivers/input/misc/da7280.ko:
kernel/drivers/input/misc/da9052_onkey.ko:
kernel/drivers/input/misc/da9055_onkey.ko:
kernel/drivers/input/misc/da9063_onkey.ko:
kernel/drivers/input/misc/e3x0-button.ko:
kernel/drivers/input/misc/drv260x.ko: kernel/drivers/input/ff-memless.ko
kernel/drivers/input/misc/drv2665.ko: kernel/drivers/input/ff-memless.ko
kernel/drivers/input/misc/drv2667.ko: kernel/drivers/input/ff-memless.ko
kernel/drivers/input/misc/gpio-beeper.ko:
kernel/drivers/input/misc/gpio_decoder.ko:
kernel/drivers/input/misc/gpio-vibra.ko: kernel/drivers/input/ff-memless.ko
kernel/drivers/input/misc/ims-pcu.ko:
kernel/drivers/input/misc/iqs269a.ko:
kernel/drivers/input/misc/iqs626a.ko:
kernel/drivers/input/misc/keyspan_remote.ko:
kernel/drivers/input/misc/kxtj9.ko:
kernel/drivers/input/misc/max77693-haptic.ko: kernel/drivers/input/ff-memless.ko
kernel/drivers/input/misc/max8925_onkey.ko:
kernel/drivers/input/misc/max8997_haptic.ko: kernel/drivers/input/ff-memless.ko
kernel/drivers/input/misc/mc13783-pwrbutton.ko: kernel/drivers/mfd/mc13xxx-core.ko
kernel/drivers/input/misc/mma8450.ko:
kernel/drivers/input/misc/palmas-pwrbutton.ko:
kernel/drivers/input/misc/pcap_keys.ko:
kernel/drivers/input/misc/pcf50633-input.ko: kernel/drivers/mfd/pcf50633.ko
kernel/drivers/input/misc/pcf8574_keypad.ko:
kernel/drivers/input/misc/pcspkr.ko:
kernel/drivers/input/misc/powermate.ko:
kernel/drivers/input/misc/pwm-beeper.ko:
kernel/drivers/input/misc/pwm-vibra.ko: kernel/drivers/input/ff-memless.ko
kernel/drivers/input/misc/rave-sp-pwrbutton.ko: kernel/drivers/mfd/rave-sp.ko
kernel/drivers/input/misc/regulator-haptic.ko: kernel/drivers/input/ff-memless.ko
kernel/drivers/input/misc/retu-pwrbutton.ko: kernel/drivers/mfd/retu-mfd.ko
kernel/drivers/input/misc/axp20x-pek.ko:
kernel/drivers/input/misc/rotary_encoder.ko:
kernel/drivers/input/misc/soc_button_array.ko:
kernel/drivers/input/misc/twl4030-pwrbutton.ko:
kernel/drivers/input/misc/twl4030-vibra.ko: kernel/drivers/input/ff-memless.ko
kernel/drivers/input/misc/twl6040-vibra.ko:
kernel/drivers/input/misc/wm831x-on.ko:
kernel/drivers/input/misc/xen-kbdfront.ko:
kernel/drivers/input/misc/yealink.ko:
kernel/drivers/input/misc/ideapad_slidebar.ko:
kernel/drivers/input/ff-memless.ko:
kernel/drivers/input/sparse-keymap.ko:
kernel/drivers/input/matrix-keymap.ko:
kernel/drivers/input/input-leds.ko:
kernel/drivers/input/joydev.ko:
kernel/drivers/input/evbug.ko:
kernel/drivers/input/rmi4/rmi_core.ko: kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/input/rmi4/rmi_i2c.ko: kernel/drivers/input/rmi4/rmi_core.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/input/rmi4/rmi_spi.ko: kernel/drivers/input/rmi4/rmi_core.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/input/rmi4/rmi_smbus.ko: kernel/drivers/input/rmi4/rmi_core.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/rtc/rtc-88pm80x.ko:
kernel/drivers/rtc/rtc-88pm860x.ko:
kernel/drivers/rtc/rtc-ab-b5ze-s3.ko:
kernel/drivers/rtc/rtc-ab-eoz9.ko:
kernel/drivers/rtc/rtc-abx80x.ko:
kernel/drivers/rtc/rtc-bq32k.ko:
kernel/drivers/rtc/rtc-bq4802.ko:
kernel/drivers/rtc/rtc-cros-ec.ko:
kernel/drivers/rtc/rtc-da9052.ko:
kernel/drivers/rtc/rtc-da9055.ko:
kernel/drivers/rtc/rtc-da9063.ko:
kernel/drivers/rtc/rtc-ds1286.ko:
kernel/drivers/rtc/rtc-ds1302.ko:
kernel/drivers/rtc/rtc-ds1305.ko:
kernel/drivers/rtc/rtc-ds1307.ko:
kernel/drivers/rtc/rtc-ds1343.ko:
kernel/drivers/rtc/rtc-ds1347.ko:
kernel/drivers/rtc/rtc-ds1374.ko:
kernel/drivers/rtc/rtc-ds1390.ko:
kernel/drivers/rtc/rtc-ds1511.ko:
kernel/drivers/rtc/rtc-ds1553.ko:
kernel/drivers/rtc/rtc-ds1672.ko:
kernel/drivers/rtc/rtc-ds1685.ko:
kernel/drivers/rtc/rtc-ds1742.ko:
kernel/drivers/rtc/rtc-ds2404.ko:
kernel/drivers/rtc/rtc-ds3232.ko:
kernel/drivers/rtc/rtc-em3027.ko:
kernel/drivers/rtc/rtc-fm3130.ko:
kernel/drivers/rtc/rtc-ftrtc010.ko:
kernel/drivers/rtc/rtc-goldfish.ko:
kernel/drivers/rtc/rtc-hid-sensor-time.ko: kernel/drivers/iio/common/hid-sensors/hid-sensor-iio-common.ko kernel/drivers/hid/hid-sensor-hub.ko kernel/drivers/hid/hid.ko
kernel/drivers/rtc/rtc-isl12022.ko:
kernel/drivers/rtc/rtc-isl1208.ko:
kernel/drivers/rtc/rtc-lp8788.ko:
kernel/drivers/rtc/rtc-m41t80.ko:
kernel/drivers/rtc/rtc-m41t93.ko:
kernel/drivers/rtc/rtc-m41t94.ko:
kernel/drivers/rtc/rtc-m48t35.ko:
kernel/drivers/rtc/rtc-m48t59.ko:
kernel/drivers/rtc/rtc-m48t86.ko:
kernel/drivers/rtc/rtc-max6900.ko:
kernel/drivers/rtc/rtc-max6902.ko:
kernel/drivers/rtc/rtc-max6916.ko:
kernel/drivers/rtc/rtc-max8907.ko:
kernel/drivers/rtc/rtc-max8925.ko:
kernel/drivers/rtc/rtc-max8997.ko:
kernel/drivers/rtc/rtc-max8998.ko:
kernel/drivers/rtc/rtc-mc13xxx.ko: kernel/drivers/mfd/mc13xxx-core.ko
kernel/drivers/rtc/rtc-mcp795.ko:
kernel/drivers/rtc/rtc-msm6242.ko:
kernel/drivers/rtc/rtc-mt6397.ko:
kernel/drivers/rtc/rtc-palmas.ko:
kernel/drivers/rtc/rtc-pcap.ko:
kernel/drivers/rtc/rtc-pcf2123.ko:
kernel/drivers/rtc/rtc-pcf2127.ko:
kernel/drivers/rtc/rtc-pcf50633.ko: kernel/drivers/mfd/pcf50633.ko
kernel/drivers/rtc/rtc-pcf85063.ko:
kernel/drivers/rtc/rtc-pcf8523.ko:
kernel/drivers/rtc/rtc-pcf85363.ko:
kernel/drivers/rtc/rtc-pcf8563.ko:
kernel/drivers/rtc/rtc-pcf8583.ko:
kernel/drivers/rtc/rtc-r9701.ko:
kernel/drivers/rtc/rtc-rc5t583.ko:
kernel/drivers/rtc/rtc-rp5c01.ko:
kernel/drivers/rtc/rtc-rs5c348.ko:
kernel/drivers/rtc/rtc-rs5c372.ko:
kernel/drivers/rtc/rtc-rv3028.ko:
kernel/drivers/rtc/rtc-rv3029c2.ko:
kernel/drivers/rtc/rtc-rv3032.ko:
kernel/drivers/rtc/rtc-rv8803.ko:
kernel/drivers/rtc/rtc-rx4581.ko:
kernel/drivers/rtc/rtc-rx6110.ko:
kernel/drivers/rtc/rtc-rx8010.ko:
kernel/drivers/rtc/rtc-rx8025.ko:
kernel/drivers/rtc/rtc-rx8581.ko:
kernel/drivers/rtc/rtc-s35390a.ko:
kernel/drivers/rtc/rtc-sd3078.ko:
kernel/drivers/rtc/rtc-stk17ta8.ko:
kernel/drivers/rtc/rtc-tps6586x.ko:
kernel/drivers/rtc/rtc-tps65910.ko:
kernel/drivers/rtc/rtc-tps80031.ko:
kernel/drivers/rtc/rtc-v3020.ko:
kernel/drivers/rtc/rtc-wilco-ec.ko: kernel/drivers/platform/chrome/wilco_ec/wilco_ec.ko kernel/drivers/platform/chrome/cros_ec_lpcs.ko kernel/drivers/platform/chrome/cros_ec.ko
kernel/drivers/rtc/rtc-wm831x.ko:
kernel/drivers/rtc/rtc-wm8350.ko:
kernel/drivers/rtc/rtc-x1205.ko:
kernel/drivers/i2c/algos/i2c-algo-bit.ko:
kernel/drivers/i2c/algos/i2c-algo-pca.ko:
kernel/drivers/i2c/busses/i2c-scmi.ko:
kernel/drivers/i2c/busses/i2c-ali1535.ko:
kernel/drivers/i2c/busses/i2c-ali1563.ko:
kernel/drivers/i2c/busses/i2c-ali15x3.ko:
kernel/drivers/i2c/busses/i2c-amd756.ko:
kernel/drivers/i2c/busses/i2c-amd756-s4882.ko: kernel/drivers/i2c/busses/i2c-amd756.ko
kernel/drivers/i2c/busses/i2c-amd8111.ko:
kernel/drivers/i2c/busses/i2c-cht-wc.ko:
kernel/drivers/i2c/busses/i2c-i801.ko: kernel/drivers/i2c/i2c-smbus.ko
kernel/drivers/i2c/busses/i2c-isch.ko:
kernel/drivers/i2c/busses/i2c-ismt.ko:
kernel/drivers/i2c/busses/i2c-nforce2.ko:
kernel/drivers/i2c/busses/i2c-nforce2-s4985.ko: kernel/drivers/i2c/busses/i2c-nforce2.ko
kernel/drivers/i2c/busses/i2c-nvidia-gpu.ko:
kernel/drivers/i2c/busses/i2c-piix4.ko:
kernel/drivers/i2c/busses/i2c-sis5595.ko:
kernel/drivers/i2c/busses/i2c-sis630.ko:
kernel/drivers/i2c/busses/i2c-sis96x.ko:
kernel/drivers/i2c/busses/i2c-via.ko: kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/i2c/busses/i2c-viapro.ko:
kernel/drivers/i2c/busses/i2c-amd-mp2-pci.ko:
kernel/drivers/i2c/busses/i2c-amd-mp2-plat.ko: kernel/drivers/i2c/busses/i2c-amd-mp2-pci.ko
kernel/drivers/i2c/busses/i2c-cbus-gpio.ko:
kernel/drivers/i2c/busses/i2c-designware-pci.ko:
kernel/drivers/i2c/busses/i2c-gpio.ko: kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/i2c/busses/i2c-kempld.ko: kernel/drivers/mfd/kempld-core.ko
kernel/drivers/i2c/busses/i2c-ocores.ko:
kernel/drivers/i2c/busses/i2c-pca-platform.ko: kernel/drivers/i2c/algos/i2c-algo-pca.ko
kernel/drivers/i2c/busses/i2c-simtec.ko: kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/i2c/busses/i2c-xiic.ko:
kernel/drivers/i2c/busses/i2c-diolan-u2c.ko:
kernel/drivers/i2c/busses/i2c-dln2.ko: kernel/drivers/mfd/dln2.ko
kernel/drivers/i2c/busses/i2c-cp2615.ko:
kernel/drivers/i2c/busses/i2c-parport.ko: kernel/drivers/i2c/i2c-smbus.ko kernel/drivers/parport/parport.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/i2c/busses/i2c-robotfuzz-osif.ko:
kernel/drivers/i2c/busses/i2c-taos-evm.ko:
kernel/drivers/i2c/busses/i2c-tiny-usb.ko:
kernel/drivers/i2c/busses/i2c-viperboard.ko:
kernel/drivers/i2c/busses/i2c-cros-ec-tunnel.ko:
kernel/drivers/i2c/busses/i2c-ljca.ko: kernel/drivers/mfd/ljca.ko
kernel/drivers/i2c/busses/i2c-mlxcpld.ko:
kernel/drivers/i2c/busses/i2c-virtio.ko:
kernel/drivers/i2c/muxes/i2c-mux-gpio.ko: kernel/drivers/i2c/i2c-mux.ko
kernel/drivers/i2c/muxes/i2c-mux-ltc4306.ko: kernel/drivers/i2c/i2c-mux.ko
kernel/drivers/i2c/muxes/i2c-mux-mlxcpld.ko: kernel/drivers/i2c/i2c-mux.ko
kernel/drivers/i2c/muxes/i2c-mux-pca9541.ko: kernel/drivers/i2c/i2c-mux.ko
kernel/drivers/i2c/muxes/i2c-mux-pca954x.ko: kernel/drivers/i2c/i2c-mux.ko
kernel/drivers/i2c/muxes/i2c-mux-reg.ko: kernel/drivers/i2c/i2c-mux.ko
kernel/drivers/i2c/i2c-smbus.ko:
kernel/drivers/i2c/i2c-mux.ko:
kernel/drivers/i2c/i2c-stub.ko:
kernel/drivers/i3c/i3c.ko:
kernel/drivers/i3c/master/i3c-master-cdns.ko: kernel/drivers/i3c/i3c.ko
kernel/drivers/i3c/master/dw-i3c-master.ko: kernel/drivers/i3c/i3c.ko
kernel/drivers/i3c/master/svc-i3c-master.ko: kernel/drivers/i3c/i3c.ko
kernel/drivers/i3c/master/mipi-i3c-hci/mipi-i3c-hci.ko: kernel/drivers/i3c/i3c.ko
kernel/drivers/media/i2c/msp3400.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ccs/ccs.ko: kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/i2c/ccs-pll.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/et8ek8/et8ek8.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/cx25840/cx25840.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/m5mols/m5mols.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/aptina-pll.ko:
kernel/drivers/media/i2c/tvaudio.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/tda7432.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/saa6588.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/tda9840.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/tda1997x.ko: kernel/drivers/media/v4l2-core/v4l2-dv-timings.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/drivers/media/i2c/tea6415c.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/tea6420.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/saa7110.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/saa7115.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/saa717x.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/saa7127.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/saa7185.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/saa6752hs.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ad5820.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ak7375.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/dw9714.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/dw9768.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/dw9807-vcm.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/adv7170.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/adv7175.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/adv7180.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/adv7183.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/adv7343.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/adv7393.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/adv7604.ko: kernel/drivers/media/v4l2-core/v4l2-dv-timings.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/i2c/adv7842.ko: kernel/drivers/media/v4l2-core/v4l2-dv-timings.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/i2c/ad9389b.ko: kernel/drivers/media/v4l2-core/v4l2-dv-timings.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/adv7511-v4l2.ko: kernel/drivers/media/v4l2-core/v4l2-dv-timings.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/i2c/vpx3220.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/vs6624.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/bt819.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/bt856.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/bt866.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ks0127.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ths7303.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ths8200.ko: kernel/drivers/media/v4l2-core/v4l2-dv-timings.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/tvp5150.ko: kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/tvp514x.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/tvp7002.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/tw2804.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/tw9903.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/tw9906.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/tw9910.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/cs3308.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/cs5345.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/cs53l32a.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/m52790.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/tlv320aic23b.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/uda1342.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/wm8775.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/wm8739.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/vp27smpx.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/sony-btf-mpx.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/upd64031a.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/upd64083.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ov02a10.ko: kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ov2640.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ov2680.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ov2685.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ov2740.ko: kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ov5647.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ov5648.ko: kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ov5670.ko: kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ov5675.ko: kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ov5695.ko: kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ov6650.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ov7251.ko: kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ov7640.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ov7670.ko: kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ov772x.ko: kernel/drivers/base/regmap/regmap-sccb.ko kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ov7740.ko: kernel/drivers/base/regmap/regmap-sccb.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ov8856.ko: kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ov8865.ko: kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ov9640.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ov9650.ko: kernel/drivers/base/regmap/regmap-sccb.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ov9734.ko: kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ov13858.ko: kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/mt9m001.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/mt9m032.ko: kernel/drivers/media/i2c/aptina-pll.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/mt9m111.ko: kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/mt9p031.ko: kernel/drivers/media/i2c/aptina-pll.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/mt9t001.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/mt9t112.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/mt9v011.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/mt9v032.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/mt9v111.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/sr030pc30.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/noon010pc30.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/rj54n1cb0c.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/s5k6aa.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/s5k6a3.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/s5k4ecgx.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/s5k5baf.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/s5c73m3/s5c73m3.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/adp1653.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/lm3560.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/lm3646.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ccs-pll.ko:
kernel/drivers/media/i2c/ak881x.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ir-kbd-i2c.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/i2c/video-i2c.ko: kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ml86v7667.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ov2659.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/tc358743.ko: kernel/drivers/media/v4l2-core/v4l2-dv-timings.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/i2c/hi556.ko: kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/imx208.ko: kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/imx214.ko: kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/imx219.ko: kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/imx258.ko: kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/imx274.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/imx290.ko: kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/imx319.ko: kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/imx355.ko: kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/max9271.ko:
kernel/drivers/media/i2c/rdacm20.ko: kernel/drivers/media/i2c/max9271.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/rdacm21.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/st-mipid02.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/max2175.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/hm11b1.ko: kernel/drivers/media/i2c/power_ctrl_logic.ko kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ov01a1s.ko: kernel/drivers/media/i2c/power_ctrl_logic.ko kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/i2c/ov01a10.ko: kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/misc/ivsc/intel_vsc.ko
kernel/drivers/media/i2c/ov02c10.ko: kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/misc/ivsc/intel_vsc.ko
kernel/drivers/media/i2c/power_ctrl_logic.ko:
kernel/drivers/media/tuners/tuner-xc2028.ko:
kernel/drivers/media/tuners/tuner-simple.ko: kernel/drivers/media/tuners/tuner-types.ko
kernel/drivers/media/tuners/tuner-types.ko:
kernel/drivers/media/tuners/mt20xx.ko:
kernel/drivers/media/tuners/tda8290.ko:
kernel/drivers/media/tuners/tea5767.ko:
kernel/drivers/media/tuners/tea5761.ko:
kernel/drivers/media/tuners/tda9887.ko:
kernel/drivers/media/tuners/tda827x.ko:
kernel/drivers/media/tuners/tda18271.ko:
kernel/drivers/media/tuners/xc5000.ko:
kernel/drivers/media/tuners/xc4000.ko:
kernel/drivers/media/tuners/msi001.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/tuners/mt2060.ko:
kernel/drivers/media/tuners/mt2063.ko:
kernel/drivers/media/tuners/mt2266.ko:
kernel/drivers/media/tuners/qt1010.ko:
kernel/drivers/media/tuners/mt2131.ko:
kernel/drivers/media/tuners/mxl5005s.ko:
kernel/drivers/media/tuners/mxl5007t.ko:
kernel/drivers/media/tuners/mc44s803.ko:
kernel/drivers/media/tuners/max2165.ko:
kernel/drivers/media/tuners/tda18218.ko:
kernel/drivers/media/tuners/tda18212.ko:
kernel/drivers/media/tuners/e4000.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/tuners/fc2580.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/tuners/tua9001.ko:
kernel/drivers/media/tuners/si2157.ko: kernel/drivers/media/mc/mc.ko
kernel/drivers/media/tuners/fc0011.ko:
kernel/drivers/media/tuners/fc0012.ko:
kernel/drivers/media/tuners/fc0013.ko:
kernel/drivers/media/tuners/it913x.ko:
kernel/drivers/media/tuners/r820t.ko:
kernel/drivers/media/tuners/mxl301rf.ko:
kernel/drivers/media/tuners/qm1d1c0042.ko:
kernel/drivers/media/tuners/qm1d1b0004.ko:
kernel/drivers/media/tuners/m88rs6000t.ko:
kernel/drivers/media/tuners/tda18250.ko:
kernel/drivers/media/rc/keymaps/rc-adstech-dvb-t-pci.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-alink-dtu-m.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-anysee.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-apac-viewcomp.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-astrometa-t2hybrid.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-asus-pc39.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-asus-ps3-100.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-ati-tv-wonder-hd-600.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-ati-x10.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-avermedia-a16d.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-avermedia.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-avermedia-cardbus.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-avermedia-dvbt.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-avermedia-m135a.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-avermedia-m733a-rm-k6.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-avermedia-rm-ks.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-avertv-303.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-azurewave-ad-tu700.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-beelink-gs1.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-behold.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-behold-columbus.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-budget-ci-old.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-cinergy-1400.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-cinergy.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-ct-90405.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-d680-dmb.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-delock-61959.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-dib0700-nec.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-dib0700-rc5.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-digitalnow-tinytwin.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-digittrade.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-dm1105-nec.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-dntv-live-dvb-t.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-dntv-live-dvbt-pro.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-dtt200u.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-dvbsky.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-dvico-mce.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-dvico-portable.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-em-terratec.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-encore-enltv2.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-encore-enltv.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-encore-enltv-fm53.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-evga-indtube.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-eztv.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-flydvb.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-flyvideo.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-fusionhdtv-mce.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-gadmei-rm008z.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-geekbox.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-genius-tvgo-a11mce.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-gotview7135.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-hisi-poplar.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-hisi-tv-demo.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-imon-mce.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-imon-pad.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-imon-rsc.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-iodata-bctv7e.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-it913x-v1.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-it913x-v2.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-kaiomy.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-khadas.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-khamsin.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-kworld-315u.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-kworld-pc150u.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-kworld-plus-tv-analog.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-leadtek-y04g0051.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-lme2510.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-manli.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-mecool-kii-pro.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-mecool-kiii-pro.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-medion-x10.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-medion-x10-digitainer.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-medion-x10-or2x.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-minix-neo.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-msi-digivox-ii.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-msi-digivox-iii.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-msi-tvanywhere.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-msi-tvanywhere-plus.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-nebula.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-nec-terratec-cinergy-xs.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-norwood.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-npgtech.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-odroid.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-pctv-sedna.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-pine64.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-pinnacle-color.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-pinnacle-grey.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-pinnacle-pctv-hd.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-pixelview.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-pixelview-mk12.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-pixelview-002t.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-pixelview-new.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-powercolor-real-angel.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-proteus-2309.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-purpletv.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-pv951.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-hauppauge.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-rc6-mce.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-real-audio-220-32-keys.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-reddo.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-snapstream-firefly.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-streamzap.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-tanix-tx3mini.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-tanix-tx5max.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-tbs-nec.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-technisat-ts35.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-technisat-usb2.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-terratec-cinergy-c-pci.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-terratec-cinergy-s2-hd.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-terratec-cinergy-xs.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-terratec-slim.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-terratec-slim-2.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-tevii-nec.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-tivo.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-total-media-in-hand.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-total-media-in-hand-02.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-trekstor.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-tt-1500.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-twinhan-dtv-cab-ci.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-twinhan1027.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-vega-s9x.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-videomate-m1f.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-videomate-s350.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-videomate-tv-pvr.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-videostrong-kii-pro.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-wetek-hub.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-wetek-play2.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-winfast.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-winfast-usbii-deluxe.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-su3000.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-xbox-360.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-xbox-dvd.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-x96max.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/keymaps/rc-zx-irdec.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/rc-core.ko:
kernel/drivers/media/rc/ir-nec-decoder.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/ir-rc5-decoder.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/ir-rc6-decoder.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/ir-jvc-decoder.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/ir-sony-decoder.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/ir-sanyo-decoder.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/ir-sharp-decoder.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/ir-mce_kbd-decoder.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/ir-xmp-decoder.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/ir-imon-decoder.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/ir-rcmm-decoder.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/ati_remote.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/imon.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/imon_raw.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/ite-cir.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/mceusb.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/fintek-cir.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/nuvoton-cir.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/ene_ir.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/redrat3.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/streamzap.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/winbond-cir.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/rc-loopback.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/igorplugusb.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/iguanair.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/ttusbir.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/serial_ir.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/sir_ir.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/xbox_remote.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/rc/ir_toy.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/common/b2c2/b2c2-flexcop.ko: kernel/drivers/media/dvb-frontends/s5h1420.ko kernel/drivers/media/dvb-frontends/cx24113.ko kernel/drivers/media/dvb-frontends/cx24123.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/common/saa7146/saa7146.ko:
kernel/drivers/media/common/saa7146/saa7146_vv.ko: kernel/drivers/media/v4l2-core/videobuf-dma-sg.ko kernel/drivers/media/v4l2-core/videobuf-core.ko kernel/drivers/media/common/saa7146/saa7146.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/common/siano/smsmdtv.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/common/siano/smsdvb.ko: kernel/drivers/media/common/siano/smsmdtv.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/common/v4l2-tpg/v4l2-tpg.ko:
kernel/drivers/media/common/videobuf2/videobuf2-common.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko: kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/common/videobuf2/videobuf2-memops.ko: kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko: kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/common/videobuf2/videobuf2-dma-contig.ko: kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/common/videobuf2/videobuf2-dma-sg.ko: kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/common/videobuf2/videobuf2-dvb.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/common/cx2341x.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/common/tveeprom.ko:
kernel/drivers/media/common/cypress_firmware.ko:
kernel/drivers/media/common/ttpci-eeprom.ko:
kernel/drivers/media/platform/cadence/cdns-csi2rx.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/platform/cadence/cdns-csi2tx.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/platform/aspeed-video.ko: kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/platform/via-camera.ko: kernel/drivers/video/fbdev/via/viafb.ko kernel/drivers/media/common/videobuf2/videobuf2-dma-sg.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/media/platform/marvell-ccic/cafe_ccic.ko: kernel/drivers/media/platform/marvell-ccic/mcam-core.ko kernel/drivers/media/common/videobuf2/videobuf2-dma-contig.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-dma-sg.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/platform/marvell-ccic/mcam-core.ko: kernel/drivers/media/common/videobuf2/videobuf2-dma-contig.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-dma-sg.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/platform/m2m-deinterlace.ko: kernel/drivers/media/v4l2-core/v4l2-mem2mem.ko kernel/drivers/media/common/videobuf2/videobuf2-dma-contig.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/pci/ttpci/budget-core.ko: kernel/drivers/media/common/ttpci-eeprom.ko kernel/drivers/media/common/saa7146/saa7146.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/pci/ttpci/budget.ko: kernel/drivers/media/pci/ttpci/budget-core.ko kernel/drivers/media/common/ttpci-eeprom.ko kernel/drivers/media/common/saa7146/saa7146.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/pci/ttpci/budget-av.ko: kernel/drivers/media/common/saa7146/saa7146_vv.ko kernel/drivers/media/v4l2-core/videobuf-dma-sg.ko kernel/drivers/media/v4l2-core/videobuf-core.ko kernel/drivers/media/pci/ttpci/budget-core.ko kernel/drivers/media/common/ttpci-eeprom.ko kernel/drivers/media/common/saa7146/saa7146.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/pci/ttpci/budget-ci.ko: kernel/drivers/media/pci/ttpci/budget-core.ko kernel/drivers/media/common/ttpci-eeprom.ko kernel/drivers/media/common/saa7146/saa7146.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/pci/b2c2/b2c2-flexcop-pci.ko: kernel/drivers/media/common/b2c2/b2c2-flexcop.ko kernel/drivers/media/dvb-frontends/s5h1420.ko kernel/drivers/media/dvb-frontends/cx24113.ko kernel/drivers/media/dvb-frontends/cx24123.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/pci/pluto2/pluto2.ko: kernel/drivers/media/dvb-frontends/tda1004x.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/media/pci/dm1105/dm1105.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/media/pci/pt1/earth-pt1.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/pci/pt3/earth-pt3.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/pci/mantis/mantis_core.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/pci/mantis/mantis.ko: kernel/drivers/media/pci/mantis/mantis_core.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/pci/mantis/hopper.ko: kernel/drivers/media/pci/mantis/mantis_core.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/pci/ngene/ngene.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/pci/ddbridge/ddbridge.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/pci/ddbridge/ddbridge-dummy-fe.ko:
kernel/drivers/media/pci/saa7146/mxb.ko: kernel/drivers/media/common/saa7146/saa7146_vv.ko kernel/drivers/media/v4l2-core/videobuf-dma-sg.ko kernel/drivers/media/v4l2-core/videobuf-core.ko kernel/drivers/media/common/saa7146/saa7146.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/pci/saa7146/hexium_orion.ko: kernel/drivers/media/common/saa7146/saa7146_vv.ko kernel/drivers/media/v4l2-core/videobuf-dma-sg.ko kernel/drivers/media/v4l2-core/videobuf-core.ko kernel/drivers/media/common/saa7146/saa7146.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/pci/saa7146/hexium_gemini.ko: kernel/drivers/media/common/saa7146/saa7146_vv.ko kernel/drivers/media/v4l2-core/videobuf-dma-sg.ko kernel/drivers/media/v4l2-core/videobuf-core.ko kernel/drivers/media/common/saa7146/saa7146.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/pci/smipcie/smipcie.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/media/pci/netup_unidvb/netup-unidvb.ko: kernel/drivers/media/common/videobuf2/videobuf2-dvb.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/pci/intel/ipu3/ipu3-cio2.ko: kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/common/videobuf2/videobuf2-dma-sg.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/pci/intel/ipu6/intel-ipu6.ko:
kernel/drivers/media/pci/intel/ipu6/intel-ipu6-isys.ko: kernel/drivers/media/pci/intel/ipu6/intel-ipu6.ko kernel/drivers/media/common/videobuf2/videobuf2-dma-contig.ko kernel/drivers/media/v4l2-core/v4l2-fwnode.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/pci/intel/ipu6/intel-ipu6-psys.ko: kernel/drivers/media/pci/intel/ipu6/intel-ipu6.ko
kernel/drivers/media/pci/ivtv/ivtv.ko: kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/common/cx2341x.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/media/pci/ivtv/ivtv-alsa.ko: kernel/drivers/media/pci/ivtv/ivtv.ko kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/common/cx2341x.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/media/pci/ivtv/ivtvfb.ko: kernel/drivers/media/pci/ivtv/ivtv.ko kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/common/cx2341x.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/media/pci/cx18/cx18.ko: kernel/drivers/media/v4l2-core/videobuf-vmalloc.ko kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/common/cx2341x.ko kernel/drivers/media/v4l2-core/videobuf-core.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/media/pci/cx18/cx18-alsa.ko: kernel/drivers/media/pci/cx18/cx18.ko kernel/drivers/media/v4l2-core/videobuf-vmalloc.ko kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/common/cx2341x.ko kernel/drivers/media/v4l2-core/videobuf-core.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/media/pci/cx23885/cx23885.ko: kernel/drivers/media/pci/cx23885/altera-ci.ko kernel/drivers/media/tuners/tda18271.ko kernel/drivers/misc/altera-stapl/altera-stapl.ko kernel/drivers/media/dvb-frontends/m88ds3103.ko kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/common/cx2341x.ko kernel/drivers/media/common/videobuf2/videobuf2-dvb.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/i2c/i2c-mux.ko kernel/drivers/media/common/videobuf2/videobuf2-dma-sg.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/drivers/media/pci/cx23885/altera-ci.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/pci/cx25821/cx25821.ko: kernel/drivers/media/common/videobuf2/videobuf2-dma-sg.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/pci/cx25821/cx25821-alsa.ko: kernel/drivers/media/pci/cx25821/cx25821.ko kernel/drivers/media/common/videobuf2/videobuf2-dma-sg.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/drivers/media/pci/cx88/cx88xx.ko: kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/media/pci/cx88/cx8800.ko: kernel/drivers/media/pci/cx88/cx88xx.ko kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/common/videobuf2/videobuf2-dma-sg.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/media/pci/cx88/cx8802.ko: kernel/drivers/media/pci/cx88/cx88xx.ko kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/media/pci/cx88/cx88-alsa.ko: kernel/drivers/media/pci/cx88/cx88xx.ko kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/media/pci/cx88/cx88-blackbird.ko: kernel/drivers/media/pci/cx88/cx8802.ko kernel/drivers/media/pci/cx88/cx8800.ko kernel/drivers/media/pci/cx88/cx88xx.ko kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/common/cx2341x.ko kernel/drivers/media/common/videobuf2/videobuf2-dma-sg.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/media/pci/cx88/cx88-dvb.ko: kernel/drivers/media/pci/cx88/cx88-vp3054-i2c.ko kernel/drivers/media/pci/cx88/cx8802.ko kernel/drivers/media/pci/cx88/cx88xx.ko kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/common/videobuf2/videobuf2-dvb.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/common/videobuf2/videobuf2-dma-sg.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/media/pci/cx88/cx88-vp3054-i2c.ko: kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/media/pci/bt8xx/bttv.ko: kernel/drivers/media/radio/tea575x.ko kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/v4l2-core/videobuf-dma-sg.ko kernel/drivers/media/v4l2-core/videobuf-core.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/media/pci/bt8xx/bt878.ko: kernel/drivers/media/pci/bt8xx/bttv.ko kernel/drivers/media/radio/tea575x.ko kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/v4l2-core/videobuf-dma-sg.ko kernel/drivers/media/v4l2-core/videobuf-core.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/media/pci/bt8xx/dvb-bt8xx.ko: kernel/drivers/media/pci/bt8xx/bt878.ko kernel/drivers/media/pci/bt8xx/bttv.ko kernel/drivers/media/radio/tea575x.ko kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/v4l2-core/videobuf-dma-sg.ko kernel/drivers/media/v4l2-core/videobuf-core.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/media/pci/bt8xx/dst.ko: kernel/drivers/media/pci/bt8xx/bt878.ko kernel/drivers/media/pci/bt8xx/bttv.ko kernel/drivers/media/radio/tea575x.ko kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/v4l2-core/videobuf-dma-sg.ko kernel/drivers/media/v4l2-core/videobuf-core.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/media/pci/bt8xx/dst_ca.ko: kernel/drivers/media/pci/bt8xx/dst.ko kernel/drivers/media/pci/bt8xx/bt878.ko kernel/drivers/media/pci/bt8xx/bttv.ko kernel/drivers/media/radio/tea575x.ko kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/v4l2-core/videobuf-dma-sg.ko kernel/drivers/media/v4l2-core/videobuf-core.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/media/pci/saa7134/saa7134.ko: kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/common/videobuf2/videobuf2-dma-sg.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/pci/saa7134/saa7134-empress.ko: kernel/drivers/media/pci/saa7134/saa7134.ko kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/common/videobuf2/videobuf2-dma-sg.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/pci/saa7134/saa7134-go7007.ko: kernel/drivers/media/usb/go7007/go7007.ko kernel/drivers/media/pci/saa7134/saa7134.ko kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-dma-sg.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/drivers/media/pci/saa7134/saa7134-alsa.ko: kernel/drivers/media/pci/saa7134/saa7134.ko kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/common/videobuf2/videobuf2-dma-sg.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/drivers/media/pci/saa7134/saa7134-dvb.ko: kernel/drivers/media/pci/saa7134/saa7134.ko kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/common/videobuf2/videobuf2-dvb.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/common/videobuf2/videobuf2-dma-sg.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/pci/saa7164/saa7164.ko: kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/pci/tw68/tw68.ko: kernel/drivers/media/common/videobuf2/videobuf2-dma-sg.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/pci/tw686x/tw686x.ko: kernel/drivers/media/common/videobuf2/videobuf2-dma-contig.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-dma-sg.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/drivers/media/pci/dt3155/dt3155.ko: kernel/drivers/media/common/videobuf2/videobuf2-dma-contig.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/pci/meye/meye.ko: kernel/drivers/platform/x86/sony-laptop.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/acpi/video.ko
kernel/drivers/media/pci/solo6x10/solo6x10.ko: kernel/drivers/media/common/videobuf2/videobuf2-dma-contig.ko kernel/drivers/media/common/videobuf2/videobuf2-dma-sg.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/drivers/media/pci/cobalt/cobalt.ko: kernel/drivers/mtd/chips/chipreg.ko kernel/drivers/media/v4l2-core/v4l2-dv-timings.ko kernel/drivers/media/common/videobuf2/videobuf2-dma-sg.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/media/pci/tw5864/tw5864.ko: kernel/drivers/media/common/videobuf2/videobuf2-dma-contig.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/ttusb-dec/ttusb_dec.ko: kernel/drivers/media/usb/ttusb-dec/ttusbdecfe.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/ttusb-dec/ttusbdecfe.ko:
kernel/drivers/media/usb/ttusb-budget/dvb-ttusb-budget.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/dvb-usb/dvb-usb.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb/dvb-usb-vp7045.ko: kernel/drivers/media/usb/dvb-usb/dvb-usb.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb/dvb-usb-vp702x.ko: kernel/drivers/media/usb/dvb-usb/dvb-usb.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb/dvb-usb-gp8psk.ko: kernel/drivers/media/usb/dvb-usb/dvb-usb.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb/dvb-usb-dtt200u.ko: kernel/drivers/media/usb/dvb-usb/dvb-usb.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb/dvb-usb-dibusb-common.ko: kernel/drivers/media/usb/dvb-usb/dvb-usb.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb/dvb-usb-dibusb-mc-common.ko: kernel/drivers/media/usb/dvb-usb/dvb-usb-dibusb-common.ko kernel/drivers/media/dvb-frontends/dib3000mc.ko kernel/drivers/media/dvb-frontends/dibx000_common.ko kernel/drivers/media/usb/dvb-usb/dvb-usb.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb/dvb-usb-a800.ko: kernel/drivers/media/usb/dvb-usb/dvb-usb-dibusb-mc-common.ko kernel/drivers/media/usb/dvb-usb/dvb-usb-dibusb-common.ko kernel/drivers/media/dvb-frontends/dib3000mc.ko kernel/drivers/media/dvb-frontends/dibx000_common.ko kernel/drivers/media/usb/dvb-usb/dvb-usb.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb/dvb-usb-dibusb-mb.ko: kernel/drivers/media/usb/dvb-usb/dvb-usb-dibusb-common.ko kernel/drivers/media/usb/dvb-usb/dvb-usb.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb/dvb-usb-dibusb-mc.ko: kernel/drivers/media/usb/dvb-usb/dvb-usb-dibusb-mc-common.ko kernel/drivers/media/usb/dvb-usb/dvb-usb-dibusb-common.ko kernel/drivers/media/dvb-frontends/dib3000mc.ko kernel/drivers/media/dvb-frontends/dibx000_common.ko kernel/drivers/media/usb/dvb-usb/dvb-usb.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb/dvb-usb-nova-t-usb2.ko: kernel/drivers/media/usb/dvb-usb/dvb-usb-dibusb-mc-common.ko kernel/drivers/media/usb/dvb-usb/dvb-usb-dibusb-common.ko kernel/drivers/media/dvb-frontends/dib3000mc.ko kernel/drivers/media/dvb-frontends/dibx000_common.ko kernel/drivers/media/usb/dvb-usb/dvb-usb.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb/dvb-usb-umt-010.ko: kernel/drivers/media/usb/dvb-usb/dvb-usb-dibusb-common.ko kernel/drivers/media/usb/dvb-usb/dvb-usb.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb/dvb-usb-m920x.ko: kernel/drivers/media/usb/dvb-usb/dvb-usb.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb/dvb-usb-digitv.ko: kernel/drivers/media/usb/dvb-usb/dvb-usb.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb/dvb-usb-cxusb.ko: kernel/drivers/media/dvb-frontends/dib0070.ko kernel/drivers/media/usb/dvb-usb/dvb-usb.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb/dvb-usb-ttusb2.ko: kernel/drivers/media/usb/dvb-usb/dvb-usb.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb/dvb-usb-dib0700.ko: kernel/drivers/media/dvb-frontends/dib9000.ko kernel/drivers/media/dvb-frontends/dib7000m.ko kernel/drivers/media/dvb-frontends/dib0090.ko kernel/drivers/media/dvb-frontends/dib0070.ko kernel/drivers/media/dvb-frontends/dib3000mc.ko kernel/drivers/media/dvb-frontends/dibx000_common.ko kernel/drivers/media/usb/dvb-usb/dvb-usb.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb/dvb-usb-opera.ko: kernel/drivers/media/usb/dvb-usb/dvb-usb.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb/dvb-usb-af9005.ko: kernel/drivers/media/usb/dvb-usb/dvb-usb.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb/dvb-usb-af9005-remote.ko:
kernel/drivers/media/usb/dvb-usb/dvb-usb-pctv452e.ko: kernel/drivers/media/usb/dvb-usb/dvb-usb.ko kernel/drivers/media/common/ttpci-eeprom.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb/dvb-usb-dw2102.ko: kernel/drivers/media/usb/dvb-usb/dvb-usb.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb/dvb-usb-dtv5100.ko: kernel/drivers/media/usb/dvb-usb/dvb-usb.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb/dvb-usb-cinergyT2.ko: kernel/drivers/media/usb/dvb-usb/dvb-usb.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb/dvb-usb-az6027.ko: kernel/drivers/media/dvb-frontends/stb0899.ko kernel/drivers/media/dvb-frontends/stb6100.ko kernel/drivers/media/usb/dvb-usb/dvb-usb.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb/dvb-usb-technisat-usb2.ko: kernel/drivers/media/usb/dvb-usb/dvb-usb.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb-v2/dvb_usb_v2.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb-v2/dvb-usb-af9015.ko: kernel/drivers/media/usb/dvb-usb-v2/dvb_usb_v2.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb-v2/dvb-usb-af9035.ko: kernel/drivers/media/usb/dvb-usb-v2/dvb_usb_v2.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb-v2/dvb-usb-anysee.ko: kernel/drivers/media/usb/dvb-usb-v2/dvb_usb_v2.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb-v2/dvb-usb-au6610.ko: kernel/drivers/media/usb/dvb-usb-v2/dvb_usb_v2.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb-v2/dvb-usb-az6007.ko: kernel/drivers/media/common/cypress_firmware.ko kernel/drivers/media/usb/dvb-usb-v2/dvb_usb_v2.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb-v2/dvb-usb-ce6230.ko: kernel/drivers/media/usb/dvb-usb-v2/dvb_usb_v2.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb-v2/dvb-usb-ec168.ko: kernel/drivers/media/usb/dvb-usb-v2/dvb_usb_v2.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb-v2/dvb-usb-lmedm04.ko: kernel/drivers/media/usb/dvb-usb-v2/dvb_usb_v2.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb-v2/dvb-usb-gl861.ko: kernel/drivers/media/usb/dvb-usb-v2/dvb_usb_v2.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb-v2/dvb-usb-mxl111sf.ko: kernel/drivers/media/usb/dvb-usb-v2/dvb_usb_v2.ko kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb-v2/mxl111sf-demod.ko:
kernel/drivers/media/usb/dvb-usb-v2/mxl111sf-tuner.ko:
kernel/drivers/media/usb/dvb-usb-v2/dvb-usb-rtl28xxu.ko: kernel/drivers/media/usb/dvb-usb-v2/dvb_usb_v2.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb-v2/dvb-usb-dvbsky.ko: kernel/drivers/media/usb/dvb-usb-v2/dvb_usb_v2.ko kernel/drivers/media/dvb-frontends/m88ds3103.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/i2c/i2c-mux.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/dvb-usb-v2/zd1301.ko: kernel/drivers/media/dvb-frontends/zd1301_demod.ko kernel/drivers/media/usb/dvb-usb-v2/dvb_usb_v2.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/siano/smsusb.ko: kernel/drivers/media/common/siano/smsmdtv.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/b2c2/b2c2-flexcop-usb.ko: kernel/drivers/media/common/b2c2/b2c2-flexcop.ko kernel/drivers/media/dvb-frontends/s5h1420.ko kernel/drivers/media/dvb-frontends/cx24113.ko kernel/drivers/media/dvb-frontends/cx24123.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/zr364xx/zr364xx.ko: kernel/drivers/media/v4l2-core/videobuf-vmalloc.ko kernel/drivers/media/v4l2-core/videobuf-core.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/stkwebcam/stkwebcam.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/s2255/s2255drv.ko: kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/uvc/uvcvideo.ko: kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_main.ko: kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_benq.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_conex.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_cpia1.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_dtcs033.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_etoms.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_finepix.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_jeilinj.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_jl2005bcd.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_kinect.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_konica.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_mars.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_mr97310a.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_nw80x.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_ov519.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_ov534.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_ov534_9.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_pac207.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_pac7302.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_pac7311.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_se401.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_sn9c2028.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_sn9c20x.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_sonixb.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_sonixj.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_spca500.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_spca501.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_spca505.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_spca506.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_spca508.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_spca561.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_spca1528.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_sq905.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_sq905c.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_sq930x.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_sunplus.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_stk014.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_stk1135.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_stv0680.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_t613.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_topro.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_touptek.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_tv8532.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_vc032x.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_vicam.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_xirlink_cit.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gspca_zc3xx.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/m5602/gspca_m5602.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/stv06xx/gspca_stv06xx.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/gspca/gl860/gspca_gl860.ko: kernel/drivers/media/usb/gspca/gspca_main.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/pwc/pwc.ko: kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/airspy/airspy.ko: kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/hackrf/hackrf.ko: kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/msi2500/msi2500.ko: kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/cpia2/cpia2.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/au0828/au0828.ko: kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/hdpvr/hdpvr.ko: kernel/drivers/media/v4l2-core/v4l2-dv-timings.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/pvrusb2/pvrusb2.ko: kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/common/cx2341x.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/stk1160/stk1160.ko: kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/cx231xx/cx231xx.ko: kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/common/cx2341x.ko kernel/drivers/i2c/i2c-mux.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/cx231xx/cx231xx-alsa.ko: kernel/drivers/media/usb/cx231xx/cx231xx.ko kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/common/cx2341x.ko kernel/drivers/i2c/i2c-mux.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/drivers/media/usb/cx231xx/cx231xx-dvb.ko: kernel/drivers/media/usb/cx231xx/cx231xx.ko kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/common/cx2341x.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/i2c/i2c-mux.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/tm6000/tm6000.ko: kernel/drivers/media/v4l2-core/videobuf-vmalloc.ko kernel/drivers/media/v4l2-core/videobuf-core.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/tm6000/tm6000-alsa.ko: kernel/drivers/media/usb/tm6000/tm6000.ko kernel/drivers/media/v4l2-core/videobuf-vmalloc.ko kernel/drivers/media/v4l2-core/videobuf-core.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/drivers/media/usb/tm6000/tm6000-dvb.ko: kernel/drivers/media/usb/tm6000/tm6000.ko kernel/drivers/media/v4l2-core/videobuf-vmalloc.ko kernel/drivers/media/v4l2-core/videobuf-core.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/em28xx/em28xx.ko: kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/em28xx/em28xx-v4l.ko: kernel/drivers/media/usb/em28xx/em28xx.ko kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/em28xx/em28xx-alsa.ko: kernel/drivers/media/usb/em28xx/em28xx.ko kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/drivers/media/usb/em28xx/em28xx-dvb.ko: kernel/drivers/media/usb/em28xx/em28xx.ko kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/em28xx/em28xx-rc.ko: kernel/drivers/media/usb/em28xx/em28xx.ko kernel/drivers/media/common/tveeprom.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/usb/usbtv/usbtv.ko: kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/drivers/media/usb/go7007/go7007.ko: kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/drivers/media/usb/go7007/go7007-usb.ko: kernel/drivers/media/usb/go7007/go7007.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/drivers/media/usb/go7007/go7007-loader.ko: kernel/drivers/media/common/cypress_firmware.ko
kernel/drivers/media/usb/go7007/s2250.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/usb/as102/dvb-as102.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/mmc/siano/smssdio.ko: kernel/drivers/media/common/siano/smsmdtv.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/firewire/firedtv.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko kernel/drivers/firewire/firewire-core.ko kernel/lib/crc-itu-t.ko
kernel/drivers/media/spi/gs1662.ko: kernel/drivers/media/v4l2-core/v4l2-dv-timings.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/spi/cxd2880-spi.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/test-drivers/vimc/vimc.ko: kernel/drivers/media/common/v4l2-tpg/v4l2-tpg.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/test-drivers/vivid/vivid.ko: kernel/drivers/media/common/v4l2-tpg/v4l2-tpg.ko kernel/drivers/media/common/videobuf2/videobuf2-dma-contig.ko kernel/drivers/media/v4l2-core/v4l2-dv-timings.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/test-drivers/vim2m.ko: kernel/drivers/media/v4l2-core/v4l2-mem2mem.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/test-drivers/vicodec/vicodec.ko: kernel/drivers/media/v4l2-core/v4l2-mem2mem.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/dvb-frontends/dvb-pll.ko:
kernel/drivers/media/dvb-frontends/stv0299.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/dvb-frontends/stb0899.ko:
kernel/drivers/media/dvb-frontends/stb6100.ko:
kernel/drivers/media/dvb-frontends/cx22700.ko:
kernel/drivers/media/dvb-frontends/s5h1432.ko:
kernel/drivers/media/dvb-frontends/cx24110.ko:
kernel/drivers/media/dvb-frontends/tda8083.ko:
kernel/drivers/media/dvb-frontends/l64781.ko:
kernel/drivers/media/dvb-frontends/dib3000mb.ko:
kernel/drivers/media/dvb-frontends/dib3000mc.ko: kernel/drivers/media/dvb-frontends/dibx000_common.ko
kernel/drivers/media/dvb-frontends/dibx000_common.ko:
kernel/drivers/media/dvb-frontends/dib7000m.ko: kernel/drivers/media/dvb-frontends/dibx000_common.ko
kernel/drivers/media/dvb-frontends/dib7000p.ko: kernel/drivers/media/dvb-frontends/dibx000_common.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/dvb-frontends/dib8000.ko: kernel/drivers/media/dvb-frontends/dibx000_common.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/dvb-frontends/dib9000.ko: kernel/drivers/media/dvb-frontends/dibx000_common.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/dvb-frontends/mt312.ko:
kernel/drivers/media/dvb-frontends/ves1820.ko:
kernel/drivers/media/dvb-frontends/ves1x93.ko:
kernel/drivers/media/dvb-frontends/tda1004x.ko:
kernel/drivers/media/dvb-frontends/sp887x.ko:
kernel/drivers/media/dvb-frontends/nxt6000.ko:
kernel/drivers/media/dvb-frontends/mt352.ko:
kernel/drivers/media/dvb-frontends/zl10036.ko:
kernel/drivers/media/dvb-frontends/zl10039.ko:
kernel/drivers/media/dvb-frontends/zl10353.ko:
kernel/drivers/media/dvb-frontends/cx22702.ko:
kernel/drivers/media/dvb-frontends/drxd.ko:
kernel/drivers/media/dvb-frontends/tda10021.ko:
kernel/drivers/media/dvb-frontends/tda10023.ko:
kernel/drivers/media/dvb-frontends/stv0297.ko:
kernel/drivers/media/dvb-frontends/nxt200x.ko:
kernel/drivers/media/dvb-frontends/or51211.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/dvb-frontends/or51132.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/dvb-frontends/bcm3510.ko:
kernel/drivers/media/dvb-frontends/s5h1420.ko:
kernel/drivers/media/dvb-frontends/lgdt330x.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/dvb-frontends/lgdt3305.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/dvb-frontends/lgdt3306a.ko: kernel/drivers/i2c/i2c-mux.ko
kernel/drivers/media/dvb-frontends/mxl692.ko:
kernel/drivers/media/dvb-frontends/lg2160.ko:
kernel/drivers/media/dvb-frontends/cx24123.ko:
kernel/drivers/media/dvb-frontends/lnbh25.ko:
kernel/drivers/media/dvb-frontends/lnbh29.ko:
kernel/drivers/media/dvb-frontends/lnbp21.ko:
kernel/drivers/media/dvb-frontends/lnbp22.ko:
kernel/drivers/media/dvb-frontends/isl6405.ko:
kernel/drivers/media/dvb-frontends/isl6421.ko:
kernel/drivers/media/dvb-frontends/tda10086.ko:
kernel/drivers/media/dvb-frontends/tda826x.ko:
kernel/drivers/media/dvb-frontends/tda8261.ko:
kernel/drivers/media/dvb-frontends/dib0070.ko:
kernel/drivers/media/dvb-frontends/dib0090.ko:
kernel/drivers/media/dvb-frontends/tua6100.ko:
kernel/drivers/media/dvb-frontends/s5h1409.ko:
kernel/drivers/media/dvb-frontends/itd1000.ko:
kernel/drivers/media/dvb-frontends/au8522_common.ko:
kernel/drivers/media/dvb-frontends/au8522_dig.ko: kernel/drivers/media/dvb-frontends/au8522_common.ko
kernel/drivers/media/dvb-frontends/au8522_decoder.ko: kernel/drivers/media/dvb-frontends/au8522_common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/dvb-frontends/tda10048.ko:
kernel/drivers/media/dvb-frontends/cx24113.ko:
kernel/drivers/media/dvb-frontends/s5h1411.ko:
kernel/drivers/media/dvb-frontends/lgs8gl5.ko:
kernel/drivers/media/dvb-frontends/tda665x.ko:
kernel/drivers/media/dvb-frontends/lgs8gxx.ko:
kernel/drivers/media/dvb-frontends/atbm8830.ko:
kernel/drivers/media/dvb-frontends/dvb_dummy_fe.ko:
kernel/drivers/media/dvb-frontends/af9013.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/i2c/i2c-mux.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/dvb-frontends/cx24116.ko:
kernel/drivers/media/dvb-frontends/cx24117.ko:
kernel/drivers/media/dvb-frontends/cx24120.ko:
kernel/drivers/media/dvb-frontends/si21xx.ko:
kernel/drivers/media/dvb-frontends/si2168.ko: kernel/drivers/i2c/i2c-mux.ko
kernel/drivers/media/dvb-frontends/stv0288.ko:
kernel/drivers/media/dvb-frontends/stb6000.ko:
kernel/drivers/media/dvb-frontends/s921.ko:
kernel/drivers/media/dvb-frontends/stv6110.ko:
kernel/drivers/media/dvb-frontends/stv0900.ko:
kernel/drivers/media/dvb-frontends/stv090x.ko:
kernel/drivers/media/dvb-frontends/stv6110x.ko:
kernel/drivers/media/dvb-frontends/m88ds3103.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/i2c/i2c-mux.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/dvb-frontends/mn88472.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/dvb-frontends/mn88473.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/dvb-frontends/isl6423.ko:
kernel/drivers/media/dvb-frontends/ec100.ko:
kernel/drivers/media/dvb-frontends/ds3000.ko:
kernel/drivers/media/dvb-frontends/ts2020.ko:
kernel/drivers/media/dvb-frontends/mb86a16.ko:
kernel/drivers/media/dvb-frontends/drx39xyj/drx39xyj.ko:
kernel/drivers/media/dvb-frontends/mb86a20s.ko:
kernel/drivers/media/dvb-frontends/ix2505v.ko:
kernel/drivers/media/dvb-frontends/stv0367.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/dvb-frontends/cxd2820r.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/dvb-frontends/cxd2841er.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/dvb-frontends/drxk.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/dvb-frontends/tda18271c2dd.ko:
kernel/drivers/media/dvb-frontends/stv0910.ko:
kernel/drivers/media/dvb-frontends/stv6111.ko:
kernel/drivers/media/dvb-frontends/mxl5xx.ko:
kernel/drivers/media/dvb-frontends/si2165.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/dvb-frontends/a8293.ko:
kernel/drivers/media/dvb-frontends/sp2.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/dvb-frontends/tda10071.ko:
kernel/drivers/media/dvb-frontends/rtl2830.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/i2c/i2c-mux.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/dvb-frontends/rtl2832.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/i2c/i2c-mux.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/dvb-frontends/rtl2832_sdr.ko: kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/dvb-frontends/m88rs2000.ko:
kernel/drivers/media/dvb-frontends/af9033.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/dvb-frontends/as102_fe.ko:
kernel/drivers/media/dvb-frontends/gp8psk-fe.ko:
kernel/drivers/media/dvb-frontends/tc90522.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/dvb-frontends/mn88443x.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/dvb-frontends/horus3a.ko:
kernel/drivers/media/dvb-frontends/ascot2e.ko:
kernel/drivers/media/dvb-frontends/helene.ko:
kernel/drivers/media/dvb-frontends/zd1301_demod.ko:
kernel/drivers/media/dvb-frontends/cxd2099.ko:
kernel/drivers/media/dvb-frontends/cxd2880/cxd2880.ko: kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/mc/mc.ko:
kernel/drivers/media/v4l2-core/videodev.ko: kernel/drivers/media/mc/mc.ko
kernel/drivers/media/v4l2-core/v4l2-fwnode.ko: kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/v4l2-core/v4l2-async.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/v4l2-core/v4l2-dv-timings.ko:
kernel/drivers/media/v4l2-core/tuner.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/v4l2-core/v4l2-mem2mem.ko: kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/v4l2-core/v4l2-flash-led-class.ko: kernel/drivers/leds/led-class-flash.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/v4l2-core/videobuf-core.ko:
kernel/drivers/media/v4l2-core/videobuf-dma-sg.ko: kernel/drivers/media/v4l2-core/videobuf-core.ko
kernel/drivers/media/v4l2-core/videobuf-vmalloc.ko: kernel/drivers/media/v4l2-core/videobuf-core.ko
kernel/drivers/media/dvb-core/dvb-core.ko: kernel/drivers/media/mc/mc.ko
kernel/drivers/media/cec/core/cec.ko: kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/cec/i2c/ch7322.ko: kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/cec/platform/cros-ec/cros-ec-cec.ko: kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/cec/platform/seco/seco-cec.ko: kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/cec/usb/pulse8/pulse8-cec.ko: kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/cec/usb/rainshadow/rainshadow-cec.ko: kernel/drivers/media/cec/core/cec.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/media/radio/radio-maxiradio.ko: kernel/drivers/media/radio/tea575x.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/radio/radio-shark.ko: kernel/drivers/media/radio/tea575x.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/radio/shark2.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/radio/radio-si476x.ko: kernel/drivers/mfd/si476x-core.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/radio/dsbr100.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/radio/si470x/radio-si470x-common.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/radio/si470x/radio-si470x-usb.ko: kernel/drivers/media/radio/si470x/radio-si470x-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/radio/si470x/radio-si470x-i2c.ko: kernel/drivers/media/radio/si470x/radio-si470x-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/radio/si4713/si4713.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/radio/si4713/radio-usb-si4713.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/radio/si4713/radio-platform-si4713.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/radio/radio-mr800.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/radio/radio-keene.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/radio/radio-ma901.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/radio/radio-tea5764.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/radio/saa7706h.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/radio/tef6862.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/radio/radio-wl1273.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/radio/wl128x/fm_drv.ko: kernel/drivers/misc/ti-st/st_drv.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/radio/tea575x.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/media/radio/radio-raremono.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/pps/clients/pps-ldisc.ko:
kernel/drivers/pps/clients/pps_parport.ko: kernel/drivers/parport/parport.ko
kernel/drivers/pps/clients/pps-gpio.ko:
kernel/drivers/ptp/ptp_ines.ko:
kernel/drivers/ptp/ptp_kvm.ko:
kernel/drivers/ptp/ptp_clockmatrix.ko:
kernel/drivers/ptp/ptp_idt82p33.ko:
kernel/drivers/ptp/ptp_vmw.ko:
kernel/drivers/ptp/ptp_ocp.ko: kernel/drivers/mtd/mtd.ko
kernel/drivers/power/reset/atc260x-poweroff.ko:
kernel/drivers/power/supply/generic-adc-battery.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/power/supply/pda_power.ko:
kernel/drivers/power/supply/axp20x_usb_power.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/power/supply/max8925_power.ko:
kernel/drivers/power/supply/wm831x_backup.ko:
kernel/drivers/power/supply/wm831x_power.ko:
kernel/drivers/power/supply/wm8350_power.ko:
kernel/drivers/power/supply/test_power.ko:
kernel/drivers/power/supply/88pm860x_battery.ko:
kernel/drivers/power/supply/adp5061.ko:
kernel/drivers/power/supply/axp20x_battery.ko:
kernel/drivers/power/supply/axp20x_ac_power.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/power/supply/cw2015_battery.ko:
kernel/drivers/power/supply/ds2760_battery.ko: kernel/drivers/w1/wire.ko
kernel/drivers/power/supply/ds2780_battery.ko: kernel/drivers/w1/slaves/w1_ds2780.ko kernel/drivers/w1/wire.ko
kernel/drivers/power/supply/ds2781_battery.ko: kernel/drivers/w1/slaves/w1_ds2781.ko kernel/drivers/w1/wire.ko
kernel/drivers/power/supply/ds2782_battery.ko:
kernel/drivers/power/supply/ltc2941-battery-gauge.ko:
kernel/drivers/power/supply/goldfish_battery.ko:
kernel/drivers/power/supply/sbs-battery.ko:
kernel/drivers/power/supply/sbs-charger.ko:
kernel/drivers/power/supply/sbs-manager.ko: kernel/drivers/i2c/i2c-mux.ko
kernel/drivers/power/supply/bq27xxx_battery.ko:
kernel/drivers/power/supply/bq27xxx_battery_i2c.ko: kernel/drivers/power/supply/bq27xxx_battery.ko
kernel/drivers/power/supply/bq27xxx_battery_hdq.ko: kernel/drivers/power/supply/bq27xxx_battery.ko kernel/drivers/w1/wire.ko
kernel/drivers/power/supply/da9030_battery.ko:
kernel/drivers/power/supply/da9052-battery.ko:
kernel/drivers/power/supply/da9150-charger.ko: kernel/drivers/mfd/da9150-core.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/power/supply/da9150-fg.ko: kernel/drivers/mfd/da9150-core.ko
kernel/drivers/power/supply/max17040_battery.ko:
kernel/drivers/power/supply/max17042_battery.ko:
kernel/drivers/power/supply/max1721x_battery.ko: kernel/drivers/base/regmap/regmap-w1.ko kernel/drivers/w1/wire.ko
kernel/drivers/power/supply/rt5033_battery.ko:
kernel/drivers/power/supply/rt9455_charger.ko:
kernel/drivers/power/supply/twl4030_madc_battery.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/power/supply/88pm860x_charger.ko:
kernel/drivers/power/supply/pcf50633-charger.ko: kernel/drivers/mfd/pcf50633.ko
kernel/drivers/power/supply/rx51_battery.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/power/supply/isp1704_charger.ko: kernel/drivers/usb/gadget/udc/udc-core.ko
kernel/drivers/power/supply/max8903_charger.ko:
kernel/drivers/power/supply/twl4030_charger.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/power/supply/lp8727_charger.ko:
kernel/drivers/power/supply/lp8788-charger.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/power/supply/gpio-charger.ko:
kernel/drivers/power/supply/lt3651-charger.ko:
kernel/drivers/power/supply/ltc4162-l-charger.ko:
kernel/drivers/power/supply/max14577_charger.ko:
kernel/drivers/power/supply/max77693_charger.ko:
kernel/drivers/power/supply/max8997_charger.ko:
kernel/drivers/power/supply/max8998_charger.ko:
kernel/drivers/power/supply/mp2629_charger.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/power/supply/mt6360_charger.ko:
kernel/drivers/power/supply/bq2415x_charger.ko:
kernel/drivers/power/supply/bq24190_charger.ko:
kernel/drivers/power/supply/bq24257_charger.ko:
kernel/drivers/power/supply/bq24735-charger.ko:
kernel/drivers/power/supply/bq2515x_charger.ko:
kernel/drivers/power/supply/bq25890_charger.ko:
kernel/drivers/power/supply/bq25980_charger.ko:
kernel/drivers/power/supply/bq256xx_charger.ko:
kernel/drivers/power/supply/smb347-charger.ko:
kernel/drivers/power/supply/tps65090-charger.ko:
kernel/drivers/power/supply/axp288_fuel_gauge.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/power/supply/axp288_charger.ko:
kernel/drivers/power/supply/cros_usbpd-charger.ko: kernel/drivers/platform/chrome/cros_usbpd_notify.ko
kernel/drivers/power/supply/cros_peripheral_charger.ko:
kernel/drivers/power/supply/bd99954-charger.ko:
kernel/drivers/power/supply/wilco-charger.ko: kernel/drivers/platform/chrome/wilco_ec/wilco_ec.ko kernel/drivers/platform/chrome/cros_ec_lpcs.ko kernel/drivers/platform/chrome/cros_ec.ko
kernel/drivers/power/supply/surface_battery.ko: kernel/drivers/platform/surface/aggregator/surface_aggregator.ko
kernel/drivers/power/supply/surface_charger.ko: kernel/drivers/platform/surface/aggregator/surface_aggregator.ko
kernel/drivers/hwmon/hwmon-vid.ko:
kernel/drivers/hwmon/acpi_power_meter.ko:
kernel/drivers/hwmon/asus_atk0110.ko:
kernel/drivers/hwmon/asb100.ko: kernel/drivers/hwmon/hwmon-vid.ko
kernel/drivers/hwmon/w83627hf.ko: kernel/drivers/hwmon/hwmon-vid.ko
kernel/drivers/hwmon/w83773g.ko:
kernel/drivers/hwmon/w83792d.ko:
kernel/drivers/hwmon/w83793.ko: kernel/drivers/hwmon/hwmon-vid.ko
kernel/drivers/hwmon/w83795.ko:
kernel/drivers/hwmon/w83781d.ko: kernel/drivers/hwmon/hwmon-vid.ko
kernel/drivers/hwmon/w83791d.ko: kernel/drivers/hwmon/hwmon-vid.ko
kernel/drivers/hwmon/abituguru.ko:
kernel/drivers/hwmon/abituguru3.ko:
kernel/drivers/hwmon/ad7314.ko:
kernel/drivers/hwmon/ad7414.ko:
kernel/drivers/hwmon/ad7418.ko:
kernel/drivers/hwmon/adc128d818.ko:
kernel/drivers/hwmon/adcxx.ko:
kernel/drivers/hwmon/adm1021.ko:
kernel/drivers/hwmon/adm1025.ko: kernel/drivers/hwmon/hwmon-vid.ko
kernel/drivers/hwmon/adm1026.ko: kernel/drivers/hwmon/hwmon-vid.ko
kernel/drivers/hwmon/adm1029.ko:
kernel/drivers/hwmon/adm1031.ko:
kernel/drivers/hwmon/adm1177.ko:
kernel/drivers/hwmon/adm9240.ko: kernel/drivers/hwmon/hwmon-vid.ko
kernel/drivers/hwmon/ads7828.ko:
kernel/drivers/hwmon/ads7871.ko:
kernel/drivers/hwmon/adt7x10.ko:
kernel/drivers/hwmon/adt7310.ko: kernel/drivers/hwmon/adt7x10.ko
kernel/drivers/hwmon/adt7410.ko: kernel/drivers/hwmon/adt7x10.ko
kernel/drivers/hwmon/adt7411.ko:
kernel/drivers/hwmon/adt7462.ko:
kernel/drivers/hwmon/adt7470.ko:
kernel/drivers/hwmon/adt7475.ko: kernel/drivers/hwmon/hwmon-vid.ko
kernel/drivers/hwmon/aht10.ko:
kernel/drivers/hwmon/applesmc.ko:
kernel/drivers/hwmon/aquacomputer_d5next.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hwmon/as370-hwmon.ko:
kernel/drivers/hwmon/asc7621.ko:
kernel/drivers/hwmon/aspeed-pwm-tacho.ko:
kernel/drivers/hwmon/atxp1.ko: kernel/drivers/hwmon/hwmon-vid.ko
kernel/drivers/hwmon/axi-fan-control.ko:
kernel/drivers/hwmon/coretemp.ko:
kernel/drivers/hwmon/corsair-cpro.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hwmon/corsair-psu.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hwmon/da9052-hwmon.ko:
kernel/drivers/hwmon/da9055-hwmon.ko:
kernel/drivers/hwmon/dell-smm-hwmon.ko:
kernel/drivers/hwmon/dme1737.ko: kernel/drivers/hwmon/hwmon-vid.ko
kernel/drivers/hwmon/drivetemp.ko:
kernel/drivers/hwmon/ds620.ko:
kernel/drivers/hwmon/ds1621.ko:
kernel/drivers/hwmon/emc1403.ko:
kernel/drivers/hwmon/emc2103.ko:
kernel/drivers/hwmon/emc6w201.ko:
kernel/drivers/hwmon/f71805f.ko:
kernel/drivers/hwmon/f71882fg.ko:
kernel/drivers/hwmon/f75375s.ko:
kernel/drivers/hwmon/fam15h_power.ko:
kernel/drivers/hwmon/fschmd.ko:
kernel/drivers/hwmon/ftsteutates.ko:
kernel/drivers/hwmon/g760a.ko:
kernel/drivers/hwmon/g762.ko:
kernel/drivers/hwmon/gl518sm.ko:
kernel/drivers/hwmon/gl520sm.ko: kernel/drivers/hwmon/hwmon-vid.ko
kernel/drivers/hwmon/hih6130.ko:
kernel/drivers/hwmon/hwmon-aaeon.ko: kernel/drivers/platform/x86/asus-wmi.ko kernel/drivers/input/sparse-keymap.ko kernel/drivers/acpi/video.ko kernel/drivers/acpi/platform_profile.ko kernel/drivers/platform/x86/wmi.ko
kernel/drivers/hwmon/i5500_temp.ko:
kernel/drivers/hwmon/i5k_amb.ko:
kernel/drivers/hwmon/ibmaem.ko: kernel/drivers/char/ipmi/ipmi_msghandler.ko
kernel/drivers/hwmon/ibmpex.ko: kernel/drivers/char/ipmi/ipmi_msghandler.ko
kernel/drivers/hwmon/iio_hwmon.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/hwmon/ina209.ko:
kernel/drivers/hwmon/ina2xx.ko:
kernel/drivers/hwmon/ina3221.ko:
kernel/drivers/hwmon/intel-m10-bmc-hwmon.ko:
kernel/drivers/hwmon/it87.ko: kernel/drivers/hwmon/hwmon-vid.ko
kernel/drivers/hwmon/jc42.ko:
kernel/drivers/hwmon/k8temp.ko:
kernel/drivers/hwmon/k10temp.ko:
kernel/drivers/hwmon/lineage-pem.ko:
kernel/drivers/hwmon/lm63.ko:
kernel/drivers/hwmon/lm70.ko:
kernel/drivers/hwmon/lm73.ko:
kernel/drivers/hwmon/lm75.ko:
kernel/drivers/hwmon/lm77.ko:
kernel/drivers/hwmon/lm78.ko: kernel/drivers/hwmon/hwmon-vid.ko
kernel/drivers/hwmon/lm80.ko:
kernel/drivers/hwmon/lm83.ko:
kernel/drivers/hwmon/lm85.ko: kernel/drivers/hwmon/hwmon-vid.ko
kernel/drivers/hwmon/lm87.ko: kernel/drivers/hwmon/hwmon-vid.ko
kernel/drivers/hwmon/lm90.ko:
kernel/drivers/hwmon/lm92.ko:
kernel/drivers/hwmon/lm93.ko: kernel/drivers/hwmon/hwmon-vid.ko
kernel/drivers/hwmon/lm95234.ko:
kernel/drivers/hwmon/lm95241.ko:
kernel/drivers/hwmon/lm95245.ko:
kernel/drivers/hwmon/ltc2945.ko:
kernel/drivers/hwmon/ltc2947-core.ko:
kernel/drivers/hwmon/ltc2947-i2c.ko: kernel/drivers/hwmon/ltc2947-core.ko
kernel/drivers/hwmon/ltc2947-spi.ko: kernel/drivers/hwmon/ltc2947-core.ko
kernel/drivers/hwmon/ltc2990.ko:
kernel/drivers/hwmon/ltc2992.ko:
kernel/drivers/hwmon/ltc4151.ko:
kernel/drivers/hwmon/ltc4215.ko:
kernel/drivers/hwmon/ltc4222.ko:
kernel/drivers/hwmon/ltc4245.ko:
kernel/drivers/hwmon/ltc4260.ko:
kernel/drivers/hwmon/ltc4261.ko:
kernel/drivers/hwmon/max1111.ko:
kernel/drivers/hwmon/max127.ko:
kernel/drivers/hwmon/max16065.ko:
kernel/drivers/hwmon/max1619.ko:
kernel/drivers/hwmon/max1668.ko:
kernel/drivers/hwmon/max197.ko:
kernel/drivers/hwmon/max31722.ko:
kernel/drivers/hwmon/max31730.ko:
kernel/drivers/hwmon/max6621.ko:
kernel/drivers/hwmon/max6639.ko:
kernel/drivers/hwmon/max6642.ko:
kernel/drivers/hwmon/max6650.ko:
kernel/drivers/hwmon/max6697.ko:
kernel/drivers/hwmon/max31790.ko:
kernel/drivers/hwmon/mc13783-adc.ko: kernel/drivers/mfd/mc13xxx-core.ko
kernel/drivers/hwmon/mcp3021.ko:
kernel/drivers/hwmon/tc654.ko:
kernel/drivers/hwmon/tps23861.ko:
kernel/drivers/hwmon/mlxreg-fan.ko:
kernel/drivers/hwmon/menf21bmc_hwmon.ko:
kernel/drivers/hwmon/mr75203.ko:
kernel/drivers/hwmon/nct6683.ko:
kernel/drivers/hwmon/nct6775.ko: kernel/drivers/hwmon/hwmon-vid.ko
kernel/drivers/hwmon/nct7802.ko:
kernel/drivers/hwmon/nct7904.ko:
kernel/drivers/hwmon/npcm750-pwm-fan.ko:
kernel/drivers/hwmon/ntc_thermistor.ko:
kernel/drivers/hwmon/nzxt-kraken2.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hwmon/pc87360.ko: kernel/drivers/hwmon/hwmon-vid.ko
kernel/drivers/hwmon/pc87427.ko:
kernel/drivers/hwmon/pcf8591.ko:
kernel/drivers/hwmon/powr1220.ko:
kernel/drivers/hwmon/sbtsi_temp.ko:
kernel/drivers/hwmon/sbrmi.ko:
kernel/drivers/hwmon/sch56xx-common.ko:
kernel/drivers/hwmon/sch5627.ko: kernel/drivers/hwmon/sch56xx-common.ko
kernel/drivers/hwmon/sch5636.ko: kernel/drivers/hwmon/sch56xx-common.ko
kernel/drivers/hwmon/sht15.ko:
kernel/drivers/hwmon/sht21.ko:
kernel/drivers/hwmon/sht3x.ko: kernel/lib/crc8.ko
kernel/drivers/hwmon/sht4x.ko: kernel/lib/crc8.ko
kernel/drivers/hwmon/shtc1.ko:
kernel/drivers/hwmon/sis5595.ko:
kernel/drivers/hwmon/smm665.ko:
kernel/drivers/hwmon/smsc47b397.ko:
kernel/drivers/hwmon/smsc47m1.ko:
kernel/drivers/hwmon/smsc47m192.ko: kernel/drivers/hwmon/hwmon-vid.ko
kernel/drivers/hwmon/stts751.ko:
kernel/drivers/hwmon/amc6821.ko:
kernel/drivers/hwmon/tc74.ko:
kernel/drivers/hwmon/thmc50.ko:
kernel/drivers/hwmon/tmp102.ko:
kernel/drivers/hwmon/tmp103.ko:
kernel/drivers/hwmon/tmp108.ko:
kernel/drivers/hwmon/tmp401.ko:
kernel/drivers/hwmon/tmp421.ko:
kernel/drivers/hwmon/tmp513.ko:
kernel/drivers/hwmon/via-cputemp.ko: kernel/drivers/hwmon/hwmon-vid.ko
kernel/drivers/hwmon/via686a.ko:
kernel/drivers/hwmon/vt1211.ko: kernel/drivers/hwmon/hwmon-vid.ko
kernel/drivers/hwmon/vt8231.ko:
kernel/drivers/hwmon/w83627ehf.ko: kernel/drivers/hwmon/hwmon-vid.ko
kernel/drivers/hwmon/w83l785ts.ko:
kernel/drivers/hwmon/w83l786ng.ko:
kernel/drivers/hwmon/wm831x-hwmon.ko:
kernel/drivers/hwmon/wm8350-hwmon.ko:
kernel/drivers/hwmon/xgene-hwmon.ko:
kernel/drivers/hwmon/pmbus/pmbus_core.ko:
kernel/drivers/hwmon/pmbus/pmbus.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/adm1266.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko kernel/lib/crc8.ko
kernel/drivers/hwmon/pmbus/adm1275.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/bel-pfe.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/bpa-rs600.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/fsp-3y.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/ibm-cffps.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/dps920ab.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/inspur-ipsps.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/ir35221.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/ir36021.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/ir38064.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/irps5401.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/isl68137.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/lm25066.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/ltc2978.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/ltc3815.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/max15301.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/max16064.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/max16601.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/max20730.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/max20751.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/max31785.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/max34440.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/max8688.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/mp2888.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/mp2975.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/pm6764tr.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/pxe1610.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/q54sj108a2.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/stpddc60.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/tps40422.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/tps53679.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/ucd9000.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/ucd9200.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/xdpe12284.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/zl6100.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/hwmon/pmbus/pim4328.ko: kernel/drivers/hwmon/pmbus/pmbus_core.ko
kernel/drivers/thermal/intel/intel_powerclamp.ko:
kernel/drivers/thermal/intel/x86_pkg_temp_thermal.ko:
kernel/drivers/thermal/intel/intel_soc_dts_iosf.ko:
kernel/drivers/thermal/intel/intel_soc_dts_thermal.ko: kernel/drivers/thermal/intel/intel_soc_dts_iosf.ko
kernel/drivers/thermal/intel/int340x_thermal/int3400_thermal.ko: kernel/drivers/thermal/intel/int340x_thermal/acpi_thermal_rel.ko
kernel/drivers/thermal/intel/int340x_thermal/int340x_thermal_zone.ko:
kernel/drivers/thermal/intel/int340x_thermal/int3402_thermal.ko: kernel/drivers/thermal/intel/int340x_thermal/int340x_thermal_zone.ko
kernel/drivers/thermal/intel/int340x_thermal/int3403_thermal.ko: kernel/drivers/thermal/intel/int340x_thermal/int340x_thermal_zone.ko
kernel/drivers/thermal/intel/int340x_thermal/processor_thermal_device.ko: kernel/drivers/thermal/intel/int340x_thermal/processor_thermal_rfim.ko kernel/drivers/thermal/intel/int340x_thermal/processor_thermal_mbox.ko kernel/drivers/thermal/intel/int340x_thermal/processor_thermal_rapl.ko kernel/drivers/powercap/intel_rapl_common.ko kernel/drivers/thermal/intel/int340x_thermal/int340x_thermal_zone.ko
kernel/drivers/thermal/intel/int340x_thermal/int3401_thermal.ko: kernel/drivers/thermal/intel/int340x_thermal/processor_thermal_device.ko kernel/drivers/thermal/intel/int340x_thermal/processor_thermal_rfim.ko kernel/drivers/thermal/intel/int340x_thermal/processor_thermal_mbox.ko kernel/drivers/thermal/intel/int340x_thermal/processor_thermal_rapl.ko kernel/drivers/powercap/intel_rapl_common.ko kernel/drivers/thermal/intel/int340x_thermal/int340x_thermal_zone.ko
kernel/drivers/thermal/intel/int340x_thermal/processor_thermal_device_pci_legacy.ko: kernel/drivers/thermal/intel/int340x_thermal/processor_thermal_device.ko kernel/drivers/thermal/intel/int340x_thermal/processor_thermal_rfim.ko kernel/drivers/thermal/intel/int340x_thermal/processor_thermal_mbox.ko kernel/drivers/thermal/intel/int340x_thermal/processor_thermal_rapl.ko kernel/drivers/powercap/intel_rapl_common.ko kernel/drivers/thermal/intel/int340x_thermal/int340x_thermal_zone.ko kernel/drivers/thermal/intel/intel_soc_dts_iosf.ko
kernel/drivers/thermal/intel/int340x_thermal/processor_thermal_device_pci.ko: kernel/drivers/thermal/intel/int340x_thermal/processor_thermal_device.ko kernel/drivers/thermal/intel/int340x_thermal/processor_thermal_rfim.ko kernel/drivers/thermal/intel/int340x_thermal/processor_thermal_mbox.ko kernel/drivers/thermal/intel/int340x_thermal/processor_thermal_rapl.ko kernel/drivers/powercap/intel_rapl_common.ko kernel/drivers/thermal/intel/int340x_thermal/int340x_thermal_zone.ko
kernel/drivers/thermal/intel/int340x_thermal/processor_thermal_rapl.ko: kernel/drivers/powercap/intel_rapl_common.ko
kernel/drivers/thermal/intel/int340x_thermal/processor_thermal_rfim.ko: kernel/drivers/thermal/intel/int340x_thermal/processor_thermal_mbox.ko
kernel/drivers/thermal/intel/int340x_thermal/processor_thermal_mbox.ko:
kernel/drivers/thermal/intel/int340x_thermal/int3406_thermal.ko: kernel/drivers/acpi/video.ko
kernel/drivers/thermal/intel/int340x_thermal/acpi_thermal_rel.ko:
kernel/drivers/thermal/intel/intel_bxt_pmic_thermal.ko:
kernel/drivers/thermal/intel/intel_pch_thermal.ko:
kernel/drivers/thermal/intel/intel_tcc_cooling.ko:
kernel/drivers/thermal/intel/intel_menlow.ko:
kernel/drivers/thermal/thermal-generic-adc.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/watchdog/pretimeout_panic.ko:
kernel/drivers/watchdog/pcwd_pci.ko:
kernel/drivers/watchdog/wdt_pci.ko:
kernel/drivers/watchdog/pcwd_usb.ko:
kernel/drivers/watchdog/cadence_wdt.ko:
kernel/drivers/watchdog/twl4030_wdt.ko:
kernel/drivers/watchdog/dw_wdt.ko:
kernel/drivers/watchdog/retu_wdt.ko: kernel/drivers/mfd/retu-mfd.ko
kernel/drivers/watchdog/acquirewdt.ko:
kernel/drivers/watchdog/advantechwdt.ko:
kernel/drivers/watchdog/alim1535_wdt.ko:
kernel/drivers/watchdog/alim7101_wdt.ko:
kernel/drivers/watchdog/ebc-c384_wdt.ko:
kernel/drivers/watchdog/f71808e_wdt.ko:
kernel/drivers/watchdog/sp5100_tco.ko:
kernel/drivers/watchdog/sbc_fitpc2_wdt.ko:
kernel/drivers/watchdog/eurotechwdt.ko:
kernel/drivers/watchdog/ib700wdt.ko:
kernel/drivers/watchdog/ibmasr.ko:
kernel/drivers/watchdog/wafer5823wdt.ko:
kernel/drivers/watchdog/i6300esb.ko:
kernel/drivers/watchdog/ie6xx_wdt.ko:
kernel/drivers/watchdog/iTCO_wdt.ko: kernel/drivers/mfd/intel_pmc_bxt.ko kernel/drivers/watchdog/iTCO_vendor_support.ko
kernel/drivers/watchdog/iTCO_vendor_support.ko:
kernel/drivers/watchdog/it8712f_wdt.ko:
kernel/drivers/watchdog/it87_wdt.ko:
kernel/drivers/watchdog/hpwdt.ko:
kernel/drivers/watchdog/kempld_wdt.ko: kernel/drivers/mfd/kempld-core.ko
kernel/drivers/watchdog/sc1200wdt.ko:
kernel/drivers/watchdog/pc87413_wdt.ko:
kernel/drivers/watchdog/nv_tco.ko:
kernel/drivers/watchdog/sbc60xxwdt.ko:
kernel/drivers/watchdog/cpu5wdt.ko:
kernel/drivers/watchdog/sch311x_wdt.ko:
kernel/drivers/watchdog/smsc37b787_wdt.ko:
kernel/drivers/watchdog/tqmx86_wdt.ko:
kernel/drivers/watchdog/via_wdt.ko:
kernel/drivers/watchdog/w83627hf_wdt.ko:
kernel/drivers/watchdog/w83877f_wdt.ko:
kernel/drivers/watchdog/w83977f_wdt.ko:
kernel/drivers/watchdog/machzwd.ko:
kernel/drivers/watchdog/sbc_epx_c3.ko:
kernel/drivers/watchdog/mei_wdt.ko: kernel/drivers/misc/mei/mei.ko
kernel/drivers/watchdog/ni903x_wdt.ko:
kernel/drivers/watchdog/nic7018_wdt.ko:
kernel/drivers/watchdog/mlx_wdt.ko:
kernel/drivers/watchdog/wdt_aaeon.ko: kernel/drivers/platform/x86/asus-wmi.ko kernel/drivers/input/sparse-keymap.ko kernel/drivers/acpi/video.ko kernel/drivers/acpi/platform_profile.ko kernel/drivers/platform/x86/wmi.ko
kernel/drivers/watchdog/of_xilinx_wdt.ko:
kernel/drivers/watchdog/mena21_wdt.ko:
kernel/drivers/watchdog/xen_wdt.ko:
kernel/drivers/watchdog/da9052_wdt.ko:
kernel/drivers/watchdog/da9055_wdt.ko:
kernel/drivers/watchdog/da9062_wdt.ko:
kernel/drivers/watchdog/da9063_wdt.ko:
kernel/drivers/watchdog/wdat_wdt.ko:
kernel/drivers/watchdog/wm831x_wdt.ko:
kernel/drivers/watchdog/wm8350_wdt.ko:
kernel/drivers/watchdog/max63xx_wdt.ko:
kernel/drivers/watchdog/ziirave_wdt.ko:
kernel/drivers/watchdog/softdog.ko:
kernel/drivers/watchdog/menf21bmc_wdt.ko:
kernel/drivers/watchdog/menz69_wdt.ko: kernel/drivers/mcb/mcb.ko
kernel/drivers/watchdog/rave-sp-wdt.ko: kernel/drivers/mfd/rave-sp.ko
kernel/drivers/md/linear.ko:
kernel/drivers/md/raid0.ko:
kernel/drivers/md/raid1.ko:
kernel/drivers/md/raid10.ko:
kernel/drivers/md/raid456.ko: kernel/crypto/async_tx/async_raid6_recov.ko kernel/crypto/async_tx/async_memcpy.ko kernel/crypto/async_tx/async_pq.ko kernel/crypto/async_tx/async_xor.ko kernel/crypto/async_tx/async_tx.ko kernel/crypto/xor.ko kernel/lib/raid6/raid6_pq.ko kernel/lib/libcrc32c.ko
kernel/drivers/md/multipath.ko:
kernel/drivers/md/faulty.ko:
kernel/drivers/md/md-cluster.ko: kernel/fs/dlm/dlm.ko
kernel/drivers/md/bcache/bcache.ko: kernel/lib/crc64.ko
kernel/drivers/md/dm-unstripe.ko:
kernel/drivers/md/dm-bufio.ko:
kernel/drivers/md/dm-bio-prison.ko:
kernel/drivers/md/dm-crypt.ko:
kernel/drivers/md/dm-delay.ko:
kernel/drivers/md/dm-flakey.ko:
kernel/drivers/md/dm-multipath.ko:
kernel/drivers/md/dm-round-robin.ko: kernel/drivers/md/dm-multipath.ko
kernel/drivers/md/dm-queue-length.ko: kernel/drivers/md/dm-multipath.ko
kernel/drivers/md/dm-service-time.ko: kernel/drivers/md/dm-multipath.ko
kernel/drivers/md/dm-historical-service-time.ko: kernel/drivers/md/dm-multipath.ko
kernel/drivers/md/dm-io-affinity.ko: kernel/drivers/md/dm-multipath.ko
kernel/drivers/md/dm-switch.ko:
kernel/drivers/md/dm-snapshot.ko: kernel/drivers/md/dm-bufio.ko
kernel/drivers/md/persistent-data/dm-persistent-data.ko: kernel/drivers/md/dm-bufio.ko kernel/lib/libcrc32c.ko
kernel/drivers/md/dm-mirror.ko: kernel/drivers/md/dm-region-hash.ko kernel/drivers/md/dm-log.ko
kernel/drivers/md/dm-log.ko:
kernel/drivers/md/dm-region-hash.ko: kernel/drivers/md/dm-log.ko
kernel/drivers/md/dm-log-userspace.ko: kernel/drivers/md/dm-log.ko
kernel/drivers/md/dm-zero.ko:
kernel/drivers/md/dm-raid.ko: kernel/drivers/md/raid456.ko kernel/crypto/async_tx/async_raid6_recov.ko kernel/crypto/async_tx/async_memcpy.ko kernel/crypto/async_tx/async_pq.ko kernel/crypto/async_tx/async_xor.ko kernel/crypto/async_tx/async_tx.ko kernel/crypto/xor.ko kernel/lib/raid6/raid6_pq.ko kernel/lib/libcrc32c.ko
kernel/drivers/md/dm-thin-pool.ko: kernel/drivers/md/persistent-data/dm-persistent-data.ko kernel/drivers/md/dm-bio-prison.ko kernel/drivers/md/dm-bufio.ko kernel/lib/libcrc32c.ko
kernel/drivers/md/dm-verity.ko: kernel/drivers/md/dm-bufio.ko
kernel/drivers/md/dm-cache.ko: kernel/drivers/md/persistent-data/dm-persistent-data.ko kernel/drivers/md/dm-bio-prison.ko kernel/drivers/md/dm-bufio.ko kernel/lib/libcrc32c.ko
kernel/drivers/md/dm-cache-smq.ko: kernel/drivers/md/dm-cache.ko kernel/drivers/md/persistent-data/dm-persistent-data.ko kernel/drivers/md/dm-bio-prison.ko kernel/drivers/md/dm-bufio.ko kernel/lib/libcrc32c.ko
kernel/drivers/md/dm-ebs.ko: kernel/drivers/md/dm-bufio.ko
kernel/drivers/md/dm-era.ko: kernel/drivers/md/persistent-data/dm-persistent-data.ko kernel/drivers/md/dm-bufio.ko kernel/lib/libcrc32c.ko
kernel/drivers/md/dm-clone.ko: kernel/drivers/md/persistent-data/dm-persistent-data.ko kernel/drivers/md/dm-bufio.ko kernel/lib/libcrc32c.ko
kernel/drivers/md/dm-log-writes.ko:
kernel/drivers/md/dm-integrity.ko: kernel/crypto/async_tx/async_xor.ko kernel/crypto/async_tx/async_tx.ko kernel/drivers/md/dm-bufio.ko kernel/crypto/xor.ko
kernel/drivers/md/dm-zoned.ko:
kernel/drivers/md/dm-writecache.ko:
kernel/drivers/accessibility/speakup/speakup_acntsa.ko: kernel/drivers/accessibility/speakup/speakup.ko
kernel/drivers/accessibility/speakup/speakup_apollo.ko: kernel/drivers/accessibility/speakup/speakup.ko
kernel/drivers/accessibility/speakup/speakup_audptr.ko: kernel/drivers/accessibility/speakup/speakup.ko
kernel/drivers/accessibility/speakup/speakup_bns.ko: kernel/drivers/accessibility/speakup/speakup.ko
kernel/drivers/accessibility/speakup/speakup_dectlk.ko: kernel/drivers/accessibility/speakup/speakup.ko
kernel/drivers/accessibility/speakup/speakup_decext.ko: kernel/drivers/accessibility/speakup/speakup.ko
kernel/drivers/accessibility/speakup/speakup_ltlk.ko: kernel/drivers/accessibility/speakup/speakup.ko
kernel/drivers/accessibility/speakup/speakup_soft.ko: kernel/drivers/accessibility/speakup/speakup.ko
kernel/drivers/accessibility/speakup/speakup_spkout.ko: kernel/drivers/accessibility/speakup/speakup.ko
kernel/drivers/accessibility/speakup/speakup_txprt.ko: kernel/drivers/accessibility/speakup/speakup.ko
kernel/drivers/accessibility/speakup/speakup_dummy.ko: kernel/drivers/accessibility/speakup/speakup.ko
kernel/drivers/accessibility/speakup/speakup.ko:
kernel/drivers/isdn/hardware/mISDN/hfcpci.ko: kernel/drivers/isdn/mISDN/mISDN_core.ko
kernel/drivers/isdn/hardware/mISDN/hfcmulti.ko: kernel/drivers/isdn/mISDN/mISDN_core.ko
kernel/drivers/isdn/hardware/mISDN/hfcsusb.ko: kernel/drivers/isdn/mISDN/mISDN_core.ko
kernel/drivers/isdn/hardware/mISDN/avmfritz.ko: kernel/drivers/isdn/hardware/mISDN/mISDNipac.ko kernel/drivers/isdn/mISDN/mISDN_core.ko
kernel/drivers/isdn/hardware/mISDN/speedfax.ko: kernel/drivers/isdn/hardware/mISDN/mISDNisar.ko kernel/drivers/isdn/hardware/mISDN/mISDNipac.ko kernel/drivers/isdn/mISDN/mISDN_core.ko
kernel/drivers/isdn/hardware/mISDN/mISDNinfineon.ko: kernel/drivers/isdn/hardware/mISDN/mISDNipac.ko kernel/drivers/isdn/mISDN/mISDN_core.ko
kernel/drivers/isdn/hardware/mISDN/w6692.ko: kernel/drivers/isdn/mISDN/mISDN_core.ko
kernel/drivers/isdn/hardware/mISDN/netjet.ko: kernel/drivers/isdn/hardware/mISDN/isdnhdlc.ko kernel/drivers/isdn/hardware/mISDN/mISDNipac.ko kernel/drivers/isdn/mISDN/mISDN_core.ko
kernel/drivers/isdn/hardware/mISDN/mISDNipac.ko: kernel/drivers/isdn/mISDN/mISDN_core.ko
kernel/drivers/isdn/hardware/mISDN/mISDNisar.ko: kernel/drivers/isdn/mISDN/mISDN_core.ko
kernel/drivers/isdn/hardware/mISDN/isdnhdlc.ko:
kernel/drivers/isdn/capi/kernelcapi.ko:
kernel/drivers/isdn/mISDN/mISDN_core.ko:
kernel/drivers/isdn/mISDN/mISDN_dsp.ko: kernel/drivers/isdn/mISDN/mISDN_core.ko
kernel/drivers/isdn/mISDN/l1oip.ko: kernel/drivers/isdn/mISDN/mISDN_core.ko
kernel/drivers/edac/edac_mce_amd.ko:
kernel/drivers/edac/i5000_edac.ko:
kernel/drivers/edac/i5100_edac.ko:
kernel/drivers/edac/i5400_edac.ko:
kernel/drivers/edac/i7300_edac.ko:
kernel/drivers/edac/i7core_edac.ko:
kernel/drivers/edac/sb_edac.ko:
kernel/drivers/edac/pnd2_edac.ko:
kernel/drivers/edac/igen6_edac.ko:
kernel/drivers/edac/e752x_edac.ko:
kernel/drivers/edac/i82975x_edac.ko:
kernel/drivers/edac/i3000_edac.ko:
kernel/drivers/edac/i3200_edac.ko:
kernel/drivers/edac/ie31200_edac.ko:
kernel/drivers/edac/x38_edac.ko:
kernel/drivers/edac/amd64_edac.ko: kernel/drivers/edac/edac_mce_amd.ko
kernel/drivers/edac/skx_edac.ko: kernel/drivers/acpi/nfit/nfit.ko
kernel/drivers/edac/i10nm_edac.ko: kernel/drivers/acpi/nfit/nfit.ko
kernel/drivers/cpufreq/speedstep-lib.ko:
kernel/drivers/cpufreq/p4-clockmod.ko: kernel/drivers/cpufreq/speedstep-lib.ko
kernel/drivers/cpufreq/amd_freq_sensitivity.ko:
kernel/drivers/cpuidle/cpuidle-haltpoll.ko:
kernel/drivers/mmc/core/mmc_block.ko:
kernel/drivers/mmc/core/sdio_uart.ko:
kernel/drivers/mmc/host/sdhci.ko:
kernel/drivers/mmc/host/sdhci-pci.ko: kernel/drivers/mmc/host/cqhci.ko kernel/drivers/mmc/host/sdhci.ko
kernel/drivers/mmc/host/sdhci-acpi.ko: kernel/drivers/mmc/host/sdhci.ko
kernel/drivers/mmc/host/sdhci_f_sdh30.ko: kernel/drivers/mmc/host/sdhci-pltfm.ko kernel/drivers/mmc/host/sdhci.ko
kernel/drivers/mmc/host/wbsd.ko:
kernel/drivers/mmc/host/alcor.ko: kernel/drivers/misc/cardreader/alcor_pci.ko
kernel/drivers/mmc/host/mtk-sd.ko: kernel/drivers/mmc/host/cqhci.ko
kernel/drivers/mmc/host/tifm_sd.ko: kernel/drivers/misc/tifm_core.ko
kernel/drivers/mmc/host/mmc_spi.ko: kernel/drivers/mmc/host/of_mmc_spi.ko kernel/lib/crc7.ko kernel/lib/crc-itu-t.ko
kernel/drivers/mmc/host/of_mmc_spi.ko:
kernel/drivers/mmc/host/sdricoh_cs.ko: kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/mmc/host/cb710-mmc.ko: kernel/drivers/misc/cb710/cb710.ko
kernel/drivers/mmc/host/via-sdmmc.ko:
kernel/drivers/mmc/host/vub300.ko:
kernel/drivers/mmc/host/ushc.ko:
kernel/drivers/mmc/host/usdhi6rol0.ko:
kernel/drivers/mmc/host/toshsd.ko:
kernel/drivers/mmc/host/rtsx_pci_sdmmc.ko: kernel/drivers/misc/cardreader/rtsx_pci.ko
kernel/drivers/mmc/host/rtsx_usb_sdmmc.ko: kernel/drivers/misc/cardreader/rtsx_usb.ko
kernel/drivers/mmc/host/sdhci-pltfm.ko: kernel/drivers/mmc/host/sdhci.ko
kernel/drivers/mmc/host/cqhci.ko:
kernel/drivers/mmc/host/sdhci-xenon-driver.ko: kernel/drivers/mmc/host/sdhci-pltfm.ko kernel/drivers/mmc/host/sdhci.ko
kernel/drivers/leds/trigger/ledtrig-timer.ko:
kernel/drivers/leds/trigger/ledtrig-oneshot.ko:
kernel/drivers/leds/trigger/ledtrig-heartbeat.ko:
kernel/drivers/leds/trigger/ledtrig-backlight.ko:
kernel/drivers/leds/trigger/ledtrig-gpio.ko:
kernel/drivers/leds/trigger/ledtrig-activity.ko:
kernel/drivers/leds/trigger/ledtrig-default-on.ko:
kernel/drivers/leds/trigger/ledtrig-transient.ko:
kernel/drivers/leds/trigger/ledtrig-camera.ko:
kernel/drivers/leds/trigger/ledtrig-netdev.ko:
kernel/drivers/leds/trigger/ledtrig-pattern.ko:
kernel/drivers/leds/trigger/ledtrig-audio.ko:
kernel/drivers/leds/trigger/ledtrig-tty.ko:
kernel/drivers/leds/led-class-flash.ko:
kernel/drivers/leds/led-class-multicolor.ko:
kernel/drivers/leds/leds-88pm860x.ko:
kernel/drivers/leds/leds-aaeon.ko: kernel/drivers/platform/x86/asus-wmi.ko kernel/drivers/input/sparse-keymap.ko kernel/drivers/acpi/video.ko kernel/drivers/acpi/platform_profile.ko kernel/drivers/platform/x86/wmi.ko
kernel/drivers/leds/leds-adp5520.ko:
kernel/drivers/leds/leds-apu.ko:
kernel/drivers/leds/leds-bd2802.ko:
kernel/drivers/leds/leds-blinkm.ko:
kernel/drivers/leds/leds-clevo-mail.ko:
kernel/drivers/leds/leds-da903x.ko:
kernel/drivers/leds/leds-da9052.ko:
kernel/drivers/leds/leds-gpio.ko:
kernel/drivers/leds/leds-ss4200.ko:
kernel/drivers/leds/leds-lm3530.ko:
kernel/drivers/leds/leds-lm3532.ko:
kernel/drivers/leds/leds-lm3533.ko: kernel/drivers/mfd/lm3533-ctrlbank.ko kernel/drivers/mfd/lm3533-core.ko
kernel/drivers/leds/leds-lm355x.ko:
kernel/drivers/leds/leds-lm36274.ko: kernel/drivers/leds/leds-ti-lmu-common.ko
kernel/drivers/leds/leds-lm3642.ko:
kernel/drivers/leds/leds-lp3944.ko:
kernel/drivers/leds/leds-lp3952.ko:
kernel/drivers/leds/leds-lp50xx.ko: kernel/drivers/leds/led-class-multicolor.ko
kernel/drivers/leds/leds-lp8788.ko:
kernel/drivers/leds/leds-lt3593.ko:
kernel/drivers/leds/leds-max8997.ko:
kernel/drivers/leds/leds-mc13783.ko: kernel/drivers/mfd/mc13xxx-core.ko
kernel/drivers/leds/leds-menf21bmc.ko:
kernel/drivers/leds/leds-mlxcpld.ko:
kernel/drivers/leds/leds-mlxreg.ko:
kernel/drivers/leds/leds-mt6323.ko:
kernel/drivers/leds/leds-nic78bx.ko:
kernel/drivers/leds/leds-pca9532.ko:
kernel/drivers/leds/leds-pca955x.ko:
kernel/drivers/leds/leds-pca963x.ko:
kernel/drivers/leds/leds-pwm.ko:
kernel/drivers/leds/leds-regulator.ko:
kernel/drivers/leds/leds-tca6507.ko:
kernel/drivers/leds/leds-ti-lmu-common.ko:
kernel/drivers/leds/leds-tlc591xx.ko:
kernel/drivers/leds/leds-tps6105x.ko:
kernel/drivers/leds/leds-wm831x-status.ko:
kernel/drivers/leds/leds-wm8350.ko:
kernel/drivers/leds/leds-dac124s085.ko:
kernel/drivers/leds/uleds.ko:
kernel/drivers/leds/flash/leds-as3645a.ko: kernel/drivers/media/v4l2-core/v4l2-flash-led-class.ko kernel/drivers/leds/led-class-flash.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/leds/flash/leds-lm3601x.ko: kernel/drivers/leds/led-class-flash.ko
kernel/drivers/leds/flash/leds-rt8515.ko: kernel/drivers/media/v4l2-core/v4l2-flash-led-class.ko kernel/drivers/leds/led-class-flash.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/leds/flash/leds-sgm3140.ko: kernel/drivers/media/v4l2-core/v4l2-flash-led-class.ko kernel/drivers/leds/led-class-flash.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/firmware/efi/efi-pstore.ko:
kernel/drivers/firmware/efi/efibc.ko:
kernel/drivers/firmware/efi/test/efi_test.ko:
kernel/drivers/firmware/efi/capsule-loader.ko:
kernel/drivers/firmware/dmi-sysfs.ko:
kernel/drivers/firmware/iscsi_ibft.ko: kernel/drivers/scsi/iscsi_boot_sysfs.ko
kernel/drivers/firmware/qemu_fw_cfg.ko:
kernel/drivers/crypto/ccp/ccp.ko:
kernel/drivers/crypto/ccp/ccp-crypto.ko: kernel/drivers/crypto/ccp/ccp.ko
kernel/drivers/crypto/atmel-i2c.ko:
kernel/drivers/crypto/atmel-ecc.ko: kernel/drivers/crypto/atmel-i2c.ko kernel/crypto/ecdh_generic.ko kernel/crypto/ecc.ko
kernel/drivers/crypto/atmel-sha204a.ko: kernel/drivers/crypto/atmel-i2c.ko
kernel/drivers/crypto/chelsio/chcr.ko: kernel/drivers/net/ethernet/chelsio/cxgb4/cxgb4.ko kernel/net/tls/tls.ko kernel/crypto/authenc.ko
kernel/drivers/crypto/cavium/nitrox/n5pf.ko:
kernel/drivers/crypto/padlock-aes.ko:
kernel/drivers/crypto/padlock-sha.ko:
kernel/drivers/crypto/qat/qat_common/intel_qat.ko: kernel/crypto/authenc.ko
kernel/drivers/crypto/qat/qat_dh895xcc/qat_dh895xcc.ko: kernel/drivers/crypto/qat/qat_common/intel_qat.ko kernel/crypto/authenc.ko
kernel/drivers/crypto/qat/qat_c3xxx/qat_c3xxx.ko: kernel/drivers/crypto/qat/qat_common/intel_qat.ko kernel/crypto/authenc.ko
kernel/drivers/crypto/qat/qat_c62x/qat_c62x.ko: kernel/drivers/crypto/qat/qat_common/intel_qat.ko kernel/crypto/authenc.ko
kernel/drivers/crypto/qat/qat_4xxx/qat_4xxx.ko: kernel/drivers/crypto/qat/qat_common/intel_qat.ko kernel/crypto/authenc.ko
kernel/drivers/crypto/qat/qat_dh895xccvf/qat_dh895xccvf.ko: kernel/drivers/crypto/qat/qat_common/intel_qat.ko kernel/crypto/authenc.ko
kernel/drivers/crypto/qat/qat_c3xxxvf/qat_c3xxxvf.ko: kernel/drivers/crypto/qat/qat_common/intel_qat.ko kernel/crypto/authenc.ko
kernel/drivers/crypto/qat/qat_c62xvf/qat_c62xvf.ko: kernel/drivers/crypto/qat/qat_common/intel_qat.ko kernel/crypto/authenc.ko
kernel/drivers/crypto/virtio/virtio_crypto.ko: kernel/crypto/crypto_engine.ko
kernel/drivers/crypto/inside-secure/crypto_safexcel.ko: kernel/crypto/authenc.ko kernel/lib/crypto/libdes.ko
kernel/drivers/crypto/amlogic/amlogic-gxl-crypto.ko: kernel/crypto/crypto_engine.ko
kernel/drivers/staging/media/atomisp/i2c/ov5693/atomisp-ov5693.ko: kernel/drivers/staging/media/atomisp/pci/atomisp_gmin_platform.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/staging/media/atomisp/i2c/atomisp-mt9m114.ko: kernel/drivers/staging/media/atomisp/pci/atomisp_gmin_platform.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/staging/media/atomisp/i2c/atomisp-gc2235.ko: kernel/drivers/staging/media/atomisp/pci/atomisp_gmin_platform.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/staging/media/atomisp/i2c/atomisp-ov2722.ko: kernel/drivers/staging/media/atomisp/pci/atomisp_gmin_platform.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/staging/media/atomisp/i2c/atomisp-ov2680.ko: kernel/drivers/staging/media/atomisp/pci/atomisp_gmin_platform.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/staging/media/atomisp/i2c/atomisp-gc0310.ko: kernel/drivers/staging/media/atomisp/pci/atomisp_gmin_platform.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/staging/media/atomisp/i2c/atomisp-libmsrlisthelper.ko:
kernel/drivers/staging/media/atomisp/i2c/atomisp-lm3554.ko: kernel/drivers/staging/media/atomisp/pci/atomisp_gmin_platform.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/staging/media/atomisp/atomisp.ko: kernel/drivers/staging/media/atomisp/pci/atomisp_gmin_platform.ko kernel/drivers/media/v4l2-core/videobuf-vmalloc.ko kernel/drivers/media/v4l2-core/videobuf-core.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/staging/media/atomisp/pci/atomisp_gmin_platform.ko:
kernel/drivers/staging/media/ipu3/ipu3-imgu.ko: kernel/drivers/media/common/videobuf2/videobuf2-dma-sg.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/staging/media/zoran/zr36067.ko: kernel/drivers/staging/media/zoran/videocodec.ko kernel/drivers/media/common/videobuf2/videobuf2-dma-contig.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/staging/media/zoran/videocodec.ko:
kernel/drivers/staging/media/zoran/zr36050.ko: kernel/drivers/staging/media/zoran/videocodec.ko
kernel/drivers/staging/media/zoran/zr36016.ko: kernel/drivers/staging/media/zoran/videocodec.ko
kernel/drivers/staging/media/zoran/zr36060.ko: kernel/drivers/staging/media/zoran/videocodec.ko
kernel/drivers/staging/media/av7110/budget-patch.ko: kernel/drivers/media/pci/ttpci/budget-core.ko kernel/drivers/media/common/ttpci-eeprom.ko kernel/drivers/media/common/saa7146/saa7146.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/staging/media/av7110/dvb-ttpci.ko: kernel/drivers/media/common/saa7146/saa7146_vv.ko kernel/drivers/media/v4l2-core/videobuf-dma-sg.ko kernel/drivers/media/v4l2-core/videobuf-core.ko kernel/drivers/media/common/ttpci-eeprom.ko kernel/drivers/media/common/saa7146/saa7146.ko kernel/drivers/media/dvb-core/dvb-core.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/drivers/media/rc/rc-core.ko
kernel/drivers/staging/media/av7110/sp8870.ko:
kernel/drivers/staging/vme/devices/vme_user.ko:
kernel/drivers/staging/android/ashmem_linux.ko:
kernel/drivers/staging/unisys/visornic/visornic.ko: kernel/drivers/visorbus/visorbus.ko
kernel/drivers/staging/unisys/visorinput/visorinput.ko: kernel/drivers/visorbus/visorbus.ko
kernel/drivers/staging/unisys/visorhba/visorhba.ko: kernel/drivers/visorbus/visorbus.ko
kernel/drivers/staging/wlan-ng/prism2_usb.ko: kernel/net/wireless/cfg80211.ko
kernel/drivers/staging/rtl8192u/r8192u_usb.ko: kernel/lib/crypto/libarc4.ko
kernel/drivers/staging/rtl8192e/rtllib.ko: kernel/net/wireless/lib80211.ko
kernel/drivers/staging/rtl8192e/rtllib_crypt_ccmp.ko: kernel/net/wireless/lib80211.ko
kernel/drivers/staging/rtl8192e/rtllib_crypt_tkip.ko: kernel/net/wireless/lib80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/staging/rtl8192e/rtllib_crypt_wep.ko: kernel/net/wireless/lib80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/staging/rtl8192e/rtl8192e/r8192e_pci.ko: kernel/drivers/staging/rtl8192e/rtllib.ko kernel/net/wireless/lib80211.ko
kernel/drivers/staging/rtl8723bs/r8723bs.ko: kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/staging/rtl8712/r8712u.ko: kernel/net/wireless/cfg80211.ko
kernel/drivers/staging/r8188eu/r8188eu.ko:
kernel/drivers/staging/rts5208/rts5208.ko:
kernel/drivers/staging/vt6655/vt6655_stage.ko: kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/staging/vt6656/vt6656_stage.ko: kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/staging/iio/accel/adis16203.ko: kernel/drivers/iio/imu/adis_lib.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/staging/iio/accel/adis16240.ko: kernel/drivers/iio/imu/adis_lib.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/staging/iio/adc/ad7816.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/staging/iio/adc/ad7280a.ko: kernel/drivers/iio/industrialio.ko kernel/lib/crc8.ko
kernel/drivers/staging/iio/addac/adt7316.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/staging/iio/addac/adt7316-spi.ko: kernel/drivers/staging/iio/addac/adt7316.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/staging/iio/addac/adt7316-i2c.ko: kernel/drivers/staging/iio/addac/adt7316.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/staging/iio/cdc/ad7746.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/staging/iio/frequency/ad9832.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/staging/iio/frequency/ad9834.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/staging/iio/impedance-analyzer/ad5933.ko: kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/staging/iio/meter/ade7854.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/staging/iio/meter/ade7854-i2c.ko: kernel/drivers/staging/iio/meter/ade7854.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/staging/iio/meter/ade7854-spi.ko: kernel/drivers/staging/iio/meter/ade7854.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/staging/iio/resolver/ad2s1210.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/staging/sm750fb/sm750fb.ko:
kernel/drivers/staging/gdm724x/gdmulte.ko:
kernel/drivers/staging/gdm724x/gdmtty.ko:
kernel/drivers/staging/fwserial/firewire-serial.ko: kernel/drivers/firewire/firewire-core.ko kernel/lib/crc-itu-t.ko
kernel/drivers/staging/gs_fpgaboot/gs_fpga.ko:
kernel/drivers/staging/fbtft/fbtft.ko: kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_agm1264k-fl.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_bd663474.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_hx8340bn.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_hx8347d.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_hx8353d.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_hx8357d.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_ili9163.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_ili9320.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_ili9325.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_ili9340.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_ili9341.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_ili9481.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_ili9486.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_pcd8544.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_ra8875.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_s6d02a1.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_s6d1121.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_seps525.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_sh1106.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_ssd1289.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_ssd1305.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_ssd1306.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_ssd1325.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_ssd1331.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_ssd1351.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_st7735r.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_st7789v.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_tinylcd.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_tls8204.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_uc1611.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_uc1701.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_upd161704.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/fbtft/fb_watterott.ko: kernel/drivers/staging/fbtft/fbtft.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko
kernel/drivers/staging/most/net/most_net.ko: kernel/drivers/most/most_core.ko
kernel/drivers/staging/most/video/most_video.ko: kernel/drivers/most/most_core.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/staging/most/i2c/most_i2c.ko: kernel/drivers/most/most_core.ko
kernel/drivers/staging/ks7010/ks7010.ko:
kernel/drivers/staging/greybus/gb-bootrom.ko: kernel/drivers/greybus/greybus.ko
kernel/drivers/staging/greybus/gb-firmware.ko: kernel/drivers/staging/greybus/gb-spilib.ko kernel/drivers/greybus/greybus.ko
kernel/drivers/staging/greybus/gb-spilib.ko: kernel/drivers/greybus/greybus.ko
kernel/drivers/staging/greybus/gb-hid.ko: kernel/drivers/greybus/greybus.ko kernel/drivers/hid/hid.ko
kernel/drivers/staging/greybus/gb-light.ko: kernel/drivers/greybus/greybus.ko kernel/drivers/media/v4l2-core/v4l2-flash-led-class.ko kernel/drivers/leds/led-class-flash.ko kernel/drivers/media/v4l2-core/v4l2-async.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/staging/greybus/gb-log.ko: kernel/drivers/greybus/greybus.ko
kernel/drivers/staging/greybus/gb-loopback.ko: kernel/drivers/greybus/greybus.ko
kernel/drivers/staging/greybus/gb-power-supply.ko: kernel/drivers/greybus/greybus.ko
kernel/drivers/staging/greybus/gb-raw.ko: kernel/drivers/greybus/greybus.ko
kernel/drivers/staging/greybus/gb-vibrator.ko: kernel/drivers/greybus/greybus.ko
kernel/drivers/staging/greybus/gb-audio-codec.ko: kernel/drivers/staging/greybus/gb-audio-gb.ko kernel/drivers/staging/greybus/gb-audio-apbridgea.ko kernel/drivers/greybus/greybus.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/drivers/staging/greybus/gb-audio-module.ko: kernel/drivers/staging/greybus/gb-audio-manager.ko kernel/drivers/staging/greybus/gb-audio-codec.ko kernel/drivers/staging/greybus/gb-audio-gb.ko kernel/drivers/staging/greybus/gb-audio-apbridgea.ko kernel/drivers/greybus/greybus.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/drivers/staging/greybus/gb-audio-gb.ko: kernel/drivers/greybus/greybus.ko
kernel/drivers/staging/greybus/gb-audio-apbridgea.ko: kernel/drivers/greybus/greybus.ko
kernel/drivers/staging/greybus/gb-audio-manager.ko:
kernel/drivers/staging/greybus/gb-gbphy.ko: kernel/drivers/greybus/greybus.ko
kernel/drivers/staging/greybus/gb-gpio.ko: kernel/drivers/staging/greybus/gb-gbphy.ko kernel/drivers/greybus/greybus.ko
kernel/drivers/staging/greybus/gb-i2c.ko: kernel/drivers/staging/greybus/gb-gbphy.ko kernel/drivers/greybus/greybus.ko
kernel/drivers/staging/greybus/gb-pwm.ko: kernel/drivers/staging/greybus/gb-gbphy.ko kernel/drivers/greybus/greybus.ko
kernel/drivers/staging/greybus/gb-sdio.ko: kernel/drivers/staging/greybus/gb-gbphy.ko kernel/drivers/greybus/greybus.ko
kernel/drivers/staging/greybus/gb-spi.ko: kernel/drivers/staging/greybus/gb-gbphy.ko kernel/drivers/staging/greybus/gb-spilib.ko kernel/drivers/greybus/greybus.ko
kernel/drivers/staging/greybus/gb-uart.ko: kernel/drivers/staging/greybus/gb-gbphy.ko kernel/drivers/greybus/greybus.ko
kernel/drivers/staging/greybus/gb-usb.ko: kernel/drivers/staging/greybus/gb-gbphy.ko kernel/drivers/greybus/greybus.ko
kernel/drivers/staging/pi433/pi433.ko:
kernel/drivers/staging/fieldbus/fieldbus_dev.ko:
kernel/drivers/staging/qlge/qlge.ko:
kernel/drivers/staging/wfx/wfx.ko: kernel/net/mac80211/mac80211.ko kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/drivers/platform/x86/dell/alienware-wmi.ko: kernel/drivers/platform/x86/wmi.ko
kernel/drivers/platform/x86/dell/dcdbas.ko:
kernel/drivers/platform/x86/dell/dell-laptop.ko: kernel/drivers/platform/x86/dell/dell-wmi.ko kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/drivers/platform/x86/dell/dell-smbios.ko kernel/drivers/platform/x86/dell/dcdbas.ko kernel/drivers/platform/x86/dell/dell-wmi-descriptor.ko kernel/drivers/input/sparse-keymap.ko kernel/drivers/acpi/video.ko kernel/drivers/platform/x86/wmi.ko
kernel/drivers/platform/x86/dell/dell-rbtn.ko:
kernel/drivers/platform/x86/dell/dell_rbu.ko:
kernel/drivers/platform/x86/dell/dell-smbios.ko: kernel/drivers/platform/x86/dell/dcdbas.ko kernel/drivers/platform/x86/dell/dell-wmi-descriptor.ko kernel/drivers/platform/x86/wmi.ko
kernel/drivers/platform/x86/dell/dell-smo8800.ko:
kernel/drivers/platform/x86/dell/dell-uart-backlight.ko: kernel/drivers/acpi/video.ko
kernel/drivers/platform/x86/dell/dell-wmi.ko: kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/drivers/platform/x86/dell/dell-smbios.ko kernel/drivers/platform/x86/dell/dcdbas.ko kernel/drivers/platform/x86/dell/dell-wmi-descriptor.ko kernel/drivers/input/sparse-keymap.ko kernel/drivers/acpi/video.ko kernel/drivers/platform/x86/wmi.ko
kernel/drivers/platform/x86/dell/dell-wmi-aio.ko: kernel/drivers/input/sparse-keymap.ko kernel/drivers/platform/x86/wmi.ko
kernel/drivers/platform/x86/dell/dell-wmi-descriptor.ko: kernel/drivers/platform/x86/wmi.ko
kernel/drivers/platform/x86/dell/dell-wmi-led.ko: kernel/drivers/platform/x86/wmi.ko
kernel/drivers/platform/x86/dell/dell-wmi-sysman/dell-wmi-sysman.ko: kernel/drivers/platform/x86/firmware_attributes_class.ko kernel/drivers/platform/x86/wmi.ko
kernel/drivers/platform/x86/intel/atomisp2/intel_atomisp2_led.ko:
kernel/drivers/platform/x86/intel/wmi/intel-wmi-sbl-fw-update.ko: kernel/drivers/platform/x86/wmi.ko
kernel/drivers/platform/x86/intel/wmi/intel-wmi-thunderbolt.ko: kernel/drivers/platform/x86/wmi.ko
kernel/drivers/platform/x86/intel/int1092/intel_sar.ko:
kernel/drivers/platform/x86/intel/int33fe/intel_cht_int33fe.ko:
kernel/drivers/platform/x86/intel/int3472/intel_skl_int3472.ko:
kernel/drivers/platform/x86/intel/pmt/pmt_class.ko:
kernel/drivers/platform/x86/intel/pmt/pmt_telemetry.ko: kernel/drivers/platform/x86/intel/pmt/pmt_class.ko
kernel/drivers/platform/x86/intel/pmt/pmt_crashlog.ko: kernel/drivers/platform/x86/intel/pmt/pmt_class.ko
kernel/drivers/platform/x86/intel/speed_select_if/isst_if_common.ko:
kernel/drivers/platform/x86/intel/speed_select_if/isst_if_mmio.ko: kernel/drivers/platform/x86/intel/speed_select_if/isst_if_common.ko
kernel/drivers/platform/x86/intel/speed_select_if/isst_if_mbox_pci.ko: kernel/drivers/platform/x86/intel/speed_select_if/isst_if_common.ko
kernel/drivers/platform/x86/intel/speed_select_if/isst_if_mbox_msr.ko: kernel/drivers/platform/x86/intel/speed_select_if/isst_if_common.ko
kernel/drivers/platform/x86/intel/telemetry/intel_telemetry_core.ko:
kernel/drivers/platform/x86/intel/telemetry/intel_telemetry_pltdrv.ko: kernel/drivers/platform/x86/intel/intel_punit_ipc.ko kernel/drivers/platform/x86/intel/telemetry/intel_telemetry_core.ko
kernel/drivers/platform/x86/intel/telemetry/intel_telemetry_debugfs.ko: kernel/drivers/platform/x86/intel/telemetry/intel_telemetry_core.ko kernel/drivers/mfd/intel_pmc_bxt.ko
kernel/drivers/platform/x86/intel/intel-hid.ko: kernel/drivers/input/sparse-keymap.ko
kernel/drivers/platform/x86/intel/intel-vbtn.ko: kernel/drivers/input/sparse-keymap.ko
kernel/drivers/platform/x86/intel/intel_int0002_vgpio.ko:
kernel/drivers/platform/x86/intel/intel_oaktrail.ko: kernel/drivers/acpi/video.ko
kernel/drivers/platform/x86/intel/intel_bxtwc_tmu.ko:
kernel/drivers/platform/x86/intel/intel_chtdc_ti_pwrbtn.ko:
kernel/drivers/platform/x86/intel/intel_mrfld_pwrbtn.ko:
kernel/drivers/platform/x86/intel/intel_punit_ipc.ko:
kernel/drivers/platform/x86/intel/intel-rst.ko:
kernel/drivers/platform/x86/intel/intel-smartconnect.ko:
kernel/drivers/platform/x86/intel/intel-uncore-frequency.ko:
kernel/drivers/platform/x86/wmi.ko:
kernel/drivers/platform/x86/wmi-bmof.ko: kernel/drivers/platform/x86/wmi.ko
kernel/drivers/platform/x86/huawei-wmi.ko: kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/drivers/input/sparse-keymap.ko kernel/drivers/platform/x86/wmi.ko
kernel/drivers/platform/x86/mxm-wmi.ko: kernel/drivers/platform/x86/wmi.ko
kernel/drivers/platform/x86/nvidia-wmi-ec-backlight.ko: kernel/drivers/platform/x86/wmi.ko
kernel/drivers/platform/x86/peaq-wmi.ko: kernel/drivers/platform/x86/wmi.ko
kernel/drivers/platform/x86/xiaomi-wmi.ko: kernel/drivers/platform/x86/wmi.ko
kernel/drivers/platform/x86/gigabyte-wmi.ko: kernel/drivers/platform/x86/wmi.ko
kernel/drivers/platform/x86/acerhdf.ko:
kernel/drivers/platform/x86/acer-wireless.ko:
kernel/drivers/platform/x86/acer-wmi.ko: kernel/drivers/input/sparse-keymap.ko kernel/drivers/acpi/video.ko kernel/drivers/platform/x86/wmi.ko
kernel/drivers/platform/x86/amd-pmc.ko:
kernel/drivers/platform/x86/adv_swbutton.ko:
kernel/drivers/platform/x86/apple-gmux.ko: kernel/drivers/video/backlight/apple_bl.ko kernel/drivers/acpi/video.ko
kernel/drivers/platform/x86/asus-laptop.ko: kernel/drivers/input/sparse-keymap.ko kernel/drivers/acpi/video.ko
kernel/drivers/platform/x86/asus-wireless.ko:
kernel/drivers/platform/x86/asus-wmi.ko: kernel/drivers/input/sparse-keymap.ko kernel/drivers/acpi/video.ko kernel/drivers/acpi/platform_profile.ko kernel/drivers/platform/x86/wmi.ko
kernel/drivers/platform/x86/asus-nb-wmi.ko: kernel/drivers/platform/x86/asus-wmi.ko kernel/drivers/input/sparse-keymap.ko kernel/drivers/acpi/video.ko kernel/drivers/acpi/platform_profile.ko kernel/drivers/platform/x86/wmi.ko
kernel/drivers/platform/x86/eeepc-laptop.ko: kernel/drivers/input/sparse-keymap.ko kernel/drivers/acpi/video.ko
kernel/drivers/platform/x86/eeepc-wmi.ko: kernel/drivers/platform/x86/asus-wmi.ko kernel/drivers/input/sparse-keymap.ko kernel/drivers/acpi/video.ko kernel/drivers/acpi/platform_profile.ko kernel/drivers/platform/x86/wmi.ko
kernel/drivers/platform/x86/meraki-mx100.ko:
kernel/drivers/platform/x86/amilo-rfkill.ko:
kernel/drivers/platform/x86/fujitsu-laptop.ko: kernel/drivers/input/sparse-keymap.ko kernel/drivers/acpi/video.ko
kernel/drivers/platform/x86/fujitsu-tablet.ko:
kernel/drivers/platform/x86/gpd-pocket-fan.ko:
kernel/drivers/platform/x86/hp_accel.ko: kernel/drivers/misc/lis3lv02d/lis3lv02d.ko
kernel/drivers/platform/x86/hp-wmi.ko: kernel/drivers/input/sparse-keymap.ko kernel/drivers/acpi/platform_profile.ko kernel/drivers/platform/x86/wmi.ko
kernel/drivers/platform/x86/uv_sysfs.ko:
kernel/drivers/platform/x86/ibm_rtl.ko:
kernel/drivers/platform/x86/ideapad-laptop.ko: kernel/drivers/input/sparse-keymap.ko kernel/drivers/acpi/video.ko kernel/drivers/acpi/platform_profile.ko kernel/drivers/platform/x86/wmi.ko
kernel/drivers/platform/x86/hdaps.ko:
kernel/drivers/platform/x86/thinkpad_acpi.ko: kernel/drivers/char/nvram.ko kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko kernel/drivers/acpi/video.ko kernel/drivers/acpi/platform_profile.ko
kernel/drivers/platform/x86/think-lmi.ko: kernel/drivers/platform/x86/firmware_attributes_class.ko kernel/drivers/platform/x86/wmi.ko
kernel/drivers/platform/x86/msi-laptop.ko: kernel/drivers/input/sparse-keymap.ko kernel/drivers/acpi/video.ko
kernel/drivers/platform/x86/msi-wmi.ko: kernel/drivers/input/sparse-keymap.ko kernel/drivers/acpi/video.ko kernel/drivers/platform/x86/wmi.ko
kernel/drivers/platform/x86/pcengines-apuv2.ko:
kernel/drivers/platform/x86/samsung-laptop.ko: kernel/drivers/acpi/video.ko
kernel/drivers/platform/x86/samsung-q10.ko:
kernel/drivers/platform/x86/toshiba_bluetooth.ko:
kernel/drivers/platform/x86/toshiba_haps.ko:
kernel/drivers/platform/x86/toshiba_acpi.ko: kernel/drivers/iio/industrialio.ko kernel/drivers/input/sparse-keymap.ko kernel/drivers/acpi/video.ko kernel/drivers/platform/x86/wmi.ko
kernel/drivers/platform/x86/classmate-laptop.ko:
kernel/drivers/platform/x86/compal-laptop.ko: kernel/drivers/acpi/video.ko
kernel/drivers/platform/x86/lg-laptop.ko: kernel/drivers/input/sparse-keymap.ko kernel/drivers/platform/x86/wmi.ko
kernel/drivers/platform/x86/panasonic-laptop.ko: kernel/drivers/input/sparse-keymap.ko kernel/drivers/acpi/video.ko
kernel/drivers/platform/x86/sony-laptop.ko: kernel/drivers/acpi/video.ko
kernel/drivers/platform/x86/system76_acpi.ko:
kernel/drivers/platform/x86/topstar-laptop.ko: kernel/drivers/input/sparse-keymap.ko
kernel/drivers/platform/x86/firmware_attributes_class.ko:
kernel/drivers/platform/x86/serial-multi-instantiate.ko:
kernel/drivers/platform/x86/mlx-platform.ko:
kernel/drivers/platform/x86/wireless-hotkey.ko:
kernel/drivers/platform/x86/intel_ips.ko:
kernel/drivers/platform/x86/intel_scu_pltdrv.ko:
kernel/drivers/platform/x86/intel_scu_ipcutil.ko:
kernel/drivers/platform/mellanox/mlxreg-hotplug.ko:
kernel/drivers/platform/mellanox/mlxreg-io.ko:
kernel/drivers/platform/chrome/chromeos_laptop.ko:
kernel/drivers/platform/chrome/chromeos_pstore.ko:
kernel/drivers/platform/chrome/chromeos_tbmc.ko:
kernel/drivers/platform/chrome/cros_ec.ko:
kernel/drivers/platform/chrome/cros_ec_i2c.ko: kernel/drivers/platform/chrome/cros_ec.ko
kernel/drivers/platform/chrome/cros_ec_ishtp.ko: kernel/drivers/hid/intel-ish-hid/intel-ishtp.ko kernel/drivers/platform/chrome/cros_ec.ko
kernel/drivers/platform/chrome/cros_ec_spi.ko: kernel/drivers/platform/chrome/cros_ec.ko
kernel/drivers/platform/chrome/cros_ec_typec.ko: kernel/drivers/platform/chrome/cros_usbpd_notify.ko kernel/drivers/usb/typec/typec.ko
kernel/drivers/platform/chrome/cros_ec_lpcs.ko: kernel/drivers/platform/chrome/cros_ec.ko
kernel/drivers/platform/chrome/cros_kbd_led_backlight.ko:
kernel/drivers/platform/chrome/cros_ec_chardev.ko:
kernel/drivers/platform/chrome/cros_ec_lightbar.ko:
kernel/drivers/platform/chrome/cros_ec_debugfs.ko:
kernel/drivers/platform/chrome/cros-ec-sensorhub.ko:
kernel/drivers/platform/chrome/cros_ec_sysfs.ko:
kernel/drivers/platform/chrome/cros_usbpd_logger.ko:
kernel/drivers/platform/chrome/cros_usbpd_notify.ko:
kernel/drivers/platform/chrome/wilco_ec/wilco_ec.ko: kernel/drivers/platform/chrome/cros_ec_lpcs.ko kernel/drivers/platform/chrome/cros_ec.ko
kernel/drivers/platform/chrome/wilco_ec/wilco_ec_debugfs.ko: kernel/drivers/platform/chrome/wilco_ec/wilco_ec.ko kernel/drivers/platform/chrome/cros_ec_lpcs.ko kernel/drivers/platform/chrome/cros_ec.ko
kernel/drivers/platform/chrome/wilco_ec/wilco_ec_events.ko:
kernel/drivers/platform/chrome/wilco_ec/wilco_ec_telem.ko: kernel/drivers/platform/chrome/wilco_ec/wilco_ec.ko kernel/drivers/platform/chrome/cros_ec_lpcs.ko kernel/drivers/platform/chrome/cros_ec.ko
kernel/drivers/platform/surface/surface3-wmi.ko: kernel/drivers/platform/x86/wmi.ko
kernel/drivers/platform/surface/surface3_button.ko:
kernel/drivers/platform/surface/surface3_power.ko:
kernel/drivers/platform/surface/surface_acpi_notify.ko: kernel/drivers/platform/surface/aggregator/surface_aggregator.ko
kernel/drivers/platform/surface/aggregator/surface_aggregator.ko:
kernel/drivers/platform/surface/surface_aggregator_cdev.ko: kernel/drivers/platform/surface/aggregator/surface_aggregator.ko
kernel/drivers/platform/surface/surface_aggregator_registry.ko: kernel/drivers/platform/surface/aggregator/surface_aggregator.ko
kernel/drivers/platform/surface/surface_dtx.ko: kernel/drivers/platform/surface/aggregator/surface_aggregator.ko
kernel/drivers/platform/surface/surface_gpe.ko:
kernel/drivers/platform/surface/surface_hotplug.ko:
kernel/drivers/platform/surface/surface_platform_profile.ko: kernel/drivers/platform/surface/aggregator/surface_aggregator.ko kernel/drivers/acpi/platform_profile.ko
kernel/drivers/platform/surface/surfacepro3_button.ko:
kernel/drivers/mailbox/mailbox-altera.ko:
kernel/drivers/virt/vboxguest/vboxguest.ko:
kernel/drivers/virt/nitro_enclaves/nitro_enclaves.ko:
kernel/drivers/virt/acrn/acrn.ko:
kernel/drivers/hv/hv_vmbus.ko:
kernel/drivers/hv/hv_utils.ko: kernel/drivers/hv/hv_vmbus.ko
kernel/drivers/hv/hv_balloon.ko: kernel/drivers/hv/hv_vmbus.ko
kernel/drivers/extcon/extcon-adc-jack.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/extcon/extcon-axp288.ko:
kernel/drivers/extcon/extcon-fsa9480.ko:
kernel/drivers/extcon/extcon-gpio.ko:
kernel/drivers/extcon/extcon-intel-int3496.ko:
kernel/drivers/extcon/extcon-intel-cht-wc.ko:
kernel/drivers/extcon/extcon-intel-mrfld.ko:
kernel/drivers/extcon/extcon-max14577.ko:
kernel/drivers/extcon/extcon-max3355.ko:
kernel/drivers/extcon/extcon-max77693.ko:
kernel/drivers/extcon/extcon-max77843.ko:
kernel/drivers/extcon/extcon-max8997.ko:
kernel/drivers/extcon/extcon-palmas.ko:
kernel/drivers/extcon/extcon-ptn5150.ko:
kernel/drivers/extcon/extcon-rt8973a.ko:
kernel/drivers/extcon/extcon-sm5502.ko:
kernel/drivers/extcon/extcon-usb-gpio.ko:
kernel/drivers/extcon/extcon-usbc-cros-ec.ko:
kernel/drivers/extcon/extcon-usbc-tusb320.ko:
kernel/drivers/memory/dfl-emif.ko: kernel/drivers/fpga/dfl.ko kernel/drivers/fpga/fpga-region.ko kernel/drivers/fpga/fpga-bridge.ko kernel/drivers/fpga/fpga-mgr.ko
kernel/drivers/vme/bridges/vme_ca91cx42.ko:
kernel/drivers/vme/bridges/vme_tsi148.ko:
kernel/drivers/vme/bridges/vme_fake.ko:
kernel/drivers/vme/boards/vme_vmivme7805.ko:
kernel/drivers/powercap/intel_rapl_common.ko:
kernel/drivers/powercap/intel_rapl_msr.ko: kernel/drivers/powercap/intel_rapl_common.ko
kernel/drivers/hwtracing/intel_th/intel_th.ko:
kernel/drivers/hwtracing/intel_th/intel_th_pci.ko: kernel/drivers/hwtracing/intel_th/intel_th.ko
kernel/drivers/hwtracing/intel_th/intel_th_acpi.ko: kernel/drivers/hwtracing/intel_th/intel_th.ko
kernel/drivers/hwtracing/intel_th/intel_th_gth.ko: kernel/drivers/hwtracing/intel_th/intel_th.ko
kernel/drivers/hwtracing/intel_th/intel_th_sth.ko: kernel/drivers/hwtracing/stm/stm_core.ko kernel/drivers/hwtracing/intel_th/intel_th.ko
kernel/drivers/hwtracing/intel_th/intel_th_msu.ko: kernel/drivers/hwtracing/intel_th/intel_th.ko
kernel/drivers/hwtracing/intel_th/intel_th_pti.ko: kernel/drivers/hwtracing/intel_th/intel_th.ko
kernel/drivers/hwtracing/intel_th/intel_th_msu_sink.ko: kernel/drivers/hwtracing/intel_th/intel_th_msu.ko kernel/drivers/hwtracing/intel_th/intel_th.ko
kernel/drivers/android/binder_linux.ko:
kernel/drivers/nvmem/nvmem_qcom-spmi-sdam.ko:
kernel/drivers/nvmem/nvmem-rave-sp-eeprom.ko:
kernel/drivers/nvmem/nvmem-rmem.ko:
kernel/drivers/vdpa/mlx5/mlx5_vdpa.ko: kernel/drivers/net/ethernet/mellanox/mlx5/core/mlx5_core.ko kernel/drivers/vhost/vringh.ko kernel/drivers/vhost/vhost_iotlb.ko kernel/drivers/net/ethernet/mellanox/mlxfw/mlxfw.ko kernel/net/psample/psample.ko kernel/net/tls/tls.ko kernel/drivers/vdpa/vdpa.ko kernel/drivers/pci/controller/pci-hyperv-intf.ko
kernel/drivers/vdpa/vdpa.ko:
kernel/drivers/vdpa/vdpa_sim/vdpa_sim.ko: kernel/drivers/vhost/vringh.ko kernel/drivers/vhost/vhost_iotlb.ko kernel/drivers/vdpa/vdpa.ko
kernel/drivers/vdpa/vdpa_sim/vdpa_sim_net.ko: kernel/drivers/vdpa/vdpa_sim/vdpa_sim.ko kernel/drivers/vhost/vringh.ko kernel/drivers/vhost/vhost_iotlb.ko kernel/drivers/vdpa/vdpa.ko
kernel/drivers/vdpa/vdpa_sim/vdpa_sim_blk.ko: kernel/drivers/vdpa/vdpa_sim/vdpa_sim.ko kernel/drivers/vhost/vringh.ko kernel/drivers/vhost/vhost_iotlb.ko kernel/drivers/vdpa/vdpa.ko
kernel/drivers/vdpa/vdpa_user/vduse.ko: kernel/drivers/vhost/vhost_iotlb.ko kernel/drivers/vdpa/vdpa.ko
kernel/drivers/vdpa/ifcvf/ifcvf.ko: kernel/drivers/vdpa/vdpa.ko
kernel/drivers/vdpa/virtio_pci/vp_vdpa.ko: kernel/drivers/vdpa/vdpa.ko
kernel/drivers/parport/parport.ko:
kernel/drivers/parport/parport_pc.ko: kernel/drivers/parport/parport.ko
kernel/drivers/parport/parport_serial.ko: kernel/drivers/parport/parport_pc.ko kernel/drivers/parport/parport.ko
kernel/drivers/parport/parport_cs.ko: kernel/drivers/parport/parport_pc.ko kernel/drivers/parport/parport.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/parport/parport_ax88796.ko: kernel/drivers/parport/parport.ko
kernel/drivers/cxl/core/cxl_core.ko:
kernel/drivers/cxl/cxl_pci.ko: kernel/drivers/cxl/core/cxl_core.ko
kernel/drivers/cxl/cxl_acpi.ko: kernel/drivers/cxl/core/cxl_core.ko
kernel/drivers/cxl/cxl_pmem.ko: kernel/drivers/cxl/core/cxl_core.ko
kernel/drivers/video/fbdev/intelfb/intelfb.ko: kernel/drivers/i2c/algos/i2c-algo-bit.ko
kernel/drivers/target/target_core_mod.ko:
kernel/drivers/target/target_core_iblock.ko: kernel/drivers/target/target_core_mod.ko
kernel/drivers/target/target_core_file.ko: kernel/drivers/target/target_core_mod.ko
kernel/drivers/target/target_core_pscsi.ko: kernel/drivers/target/target_core_mod.ko
kernel/drivers/target/target_core_user.ko: kernel/drivers/uio/uio.ko kernel/drivers/target/target_core_mod.ko
kernel/drivers/target/loopback/tcm_loop.ko: kernel/drivers/target/target_core_mod.ko
kernel/drivers/target/tcm_fc/tcm_fc.ko: kernel/drivers/scsi/libfc/libfc.ko kernel/drivers/scsi/scsi_transport_fc.ko kernel/drivers/target/target_core_mod.ko
kernel/drivers/target/iscsi/iscsi_target_mod.ko: kernel/drivers/target/target_core_mod.ko
kernel/drivers/target/iscsi/cxgbit/cxgbit.ko: kernel/drivers/target/iscsi/iscsi_target_mod.ko kernel/drivers/net/ethernet/chelsio/cxgb4/cxgb4.ko kernel/net/tls/tls.ko kernel/drivers/net/ethernet/chelsio/libcxgb/libcxgb.ko kernel/drivers/target/target_core_mod.ko
kernel/drivers/target/sbp/sbp_target.ko: kernel/drivers/firewire/firewire-core.ko kernel/drivers/target/target_core_mod.ko kernel/lib/crc-itu-t.ko
kernel/drivers/mtd/parsers/ar7part.ko: kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/parsers/cmdlinepart.ko: kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/parsers/redboot.ko: kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/chips/chipreg.ko:
kernel/drivers/mtd/chips/cfi_probe.ko: kernel/drivers/mtd/chips/cfi_util.ko kernel/drivers/mtd/chips/gen_probe.ko kernel/drivers/mtd/chips/chipreg.ko
kernel/drivers/mtd/chips/cfi_util.ko:
kernel/drivers/mtd/chips/cfi_cmdset_0020.ko: kernel/drivers/mtd/chips/cfi_util.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/chips/cfi_cmdset_0002.ko: kernel/drivers/mtd/chips/cfi_util.ko
kernel/drivers/mtd/chips/cfi_cmdset_0001.ko: kernel/drivers/mtd/chips/cfi_util.ko
kernel/drivers/mtd/chips/gen_probe.ko:
kernel/drivers/mtd/chips/jedec_probe.ko: kernel/drivers/mtd/chips/cfi_util.ko kernel/drivers/mtd/chips/gen_probe.ko kernel/drivers/mtd/chips/chipreg.ko
kernel/drivers/mtd/chips/map_ram.ko: kernel/drivers/mtd/chips/chipreg.ko
kernel/drivers/mtd/chips/map_rom.ko: kernel/drivers/mtd/chips/chipreg.ko
kernel/drivers/mtd/chips/map_absent.ko: kernel/drivers/mtd/chips/chipreg.ko
kernel/drivers/mtd/lpddr/qinfo_probe.ko: kernel/drivers/mtd/lpddr/lpddr_cmds.ko kernel/drivers/mtd/chips/chipreg.ko
kernel/drivers/mtd/lpddr/lpddr_cmds.ko:
kernel/drivers/mtd/maps/map_funcs.ko:
kernel/drivers/mtd/maps/l440gx.ko: kernel/drivers/mtd/maps/map_funcs.ko kernel/drivers/mtd/chips/chipreg.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/maps/amd76xrom.ko: kernel/drivers/mtd/maps/map_funcs.ko kernel/drivers/mtd/chips/chipreg.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/maps/esb2rom.ko: kernel/drivers/mtd/maps/map_funcs.ko kernel/drivers/mtd/chips/chipreg.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/maps/ichxrom.ko: kernel/drivers/mtd/maps/map_funcs.ko kernel/drivers/mtd/chips/chipreg.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/maps/ck804xrom.ko: kernel/drivers/mtd/maps/map_funcs.ko kernel/drivers/mtd/chips/chipreg.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/maps/physmap.ko: kernel/drivers/mtd/maps/map_funcs.ko kernel/drivers/mtd/chips/chipreg.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/maps/pcmciamtd.ko: kernel/drivers/mtd/chips/chipreg.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/maps/sbc_gxx.ko: kernel/drivers/mtd/chips/chipreg.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/maps/pci.ko: kernel/drivers/mtd/chips/chipreg.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/maps/nettel.ko: kernel/drivers/mtd/maps/map_funcs.ko kernel/drivers/mtd/chips/chipreg.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/maps/scb2_flash.ko: kernel/drivers/mtd/maps/map_funcs.ko kernel/drivers/mtd/chips/chipreg.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/maps/plat-ram.ko: kernel/drivers/mtd/maps/map_funcs.ko kernel/drivers/mtd/chips/chipreg.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/maps/intel_vr_nor.ko: kernel/drivers/mtd/maps/map_funcs.ko kernel/drivers/mtd/chips/chipreg.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/devices/slram.ko: kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/devices/phram.ko: kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/devices/pmc551.ko: kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/devices/mtdram.ko: kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/devices/block2mtd.ko: kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/devices/mtd_dataflash.ko: kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/devices/mchp23k256.ko: kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/devices/mchp48l640.ko: kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/devices/sst25l.ko: kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/nand/onenand/onenand.ko: kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/nand/onenand/generic.ko: kernel/drivers/mtd/nand/onenand/onenand.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/nand/raw/nand.ko: kernel/drivers/mtd/nand/nandcore.ko kernel/lib/bch.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/nand/raw/sm_common.ko: kernel/drivers/mtd/nand/raw/nand.ko kernel/drivers/mtd/nand/nandcore.ko kernel/lib/bch.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/nand/raw/cafe_nand.ko: kernel/drivers/mtd/nand/raw/nand.ko kernel/drivers/mtd/nand/nandcore.ko kernel/lib/bch.ko kernel/drivers/mtd/mtd.ko kernel/lib/reed_solomon/reed_solomon.ko
kernel/drivers/mtd/nand/raw/denali.ko: kernel/drivers/mtd/nand/raw/nand.ko kernel/drivers/mtd/nand/nandcore.ko kernel/lib/bch.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/nand/raw/denali_pci.ko: kernel/drivers/mtd/nand/raw/denali.ko kernel/drivers/mtd/nand/raw/nand.ko kernel/drivers/mtd/nand/nandcore.ko kernel/lib/bch.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/nand/raw/diskonchip.ko: kernel/drivers/mtd/nand/raw/nand.ko kernel/drivers/mtd/nand/nandcore.ko kernel/lib/bch.ko kernel/drivers/mtd/mtd.ko kernel/lib/reed_solomon/reed_solomon.ko
kernel/drivers/mtd/nand/raw/nandsim.ko: kernel/drivers/mtd/nand/raw/nand.ko kernel/drivers/mtd/nand/nandcore.ko kernel/lib/bch.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/nand/raw/gpio.ko: kernel/drivers/mtd/nand/raw/nand.ko kernel/drivers/mtd/nand/nandcore.ko kernel/lib/bch.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/nand/raw/plat_nand.ko: kernel/drivers/mtd/nand/raw/nand.ko kernel/drivers/mtd/nand/nandcore.ko kernel/lib/bch.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/nand/raw/r852.ko: kernel/drivers/mtd/nand/raw/sm_common.ko kernel/drivers/mtd/nand/raw/nand.ko kernel/drivers/mtd/nand/nandcore.ko kernel/lib/bch.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/nand/raw/mxic_nand.ko: kernel/drivers/mtd/nand/raw/nand.ko kernel/drivers/mtd/nand/nandcore.ko kernel/lib/bch.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/nand/raw/arasan-nand-controller.ko: kernel/drivers/mtd/nand/raw/nand.ko kernel/drivers/mtd/nand/nandcore.ko kernel/lib/bch.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/nand/spi/spinand.ko: kernel/drivers/mtd/nand/nandcore.ko kernel/lib/bch.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/nand/nandcore.ko: kernel/lib/bch.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/mtd.ko:
kernel/drivers/mtd/mtd_blkdevs.ko: kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/mtdblock.ko: kernel/drivers/mtd/mtd_blkdevs.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/mtdblock_ro.ko: kernel/drivers/mtd/mtd_blkdevs.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/ftl.ko: kernel/drivers/mtd/mtd_blkdevs.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/nftl.ko: kernel/drivers/mtd/mtd_blkdevs.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/inftl.ko: kernel/drivers/mtd/mtd_blkdevs.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/rfd_ftl.ko: kernel/drivers/mtd/mtd_blkdevs.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/ssfdc.ko: kernel/drivers/mtd/mtd_blkdevs.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/sm_ftl.ko: kernel/drivers/mtd/mtd_blkdevs.ko kernel/drivers/mtd/nand/nandcore.ko kernel/lib/bch.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/mtdoops.ko: kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/mtdpstore.ko: kernel/fs/pstore/pstore_blk.ko kernel/fs/pstore/pstore_zone.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/mtdswap.ko: kernel/drivers/mtd/mtd_blkdevs.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/spi-nor/spi-nor.ko: kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/ubi/ubi.ko: kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/ubi/gluebi.ko: kernel/drivers/mtd/ubi/ubi.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/mtd/hyperbus/hyperbus-core.ko: kernel/drivers/mtd/chips/chipreg.ko kernel/drivers/mtd/mtd.ko
kernel/drivers/spmi/spmi.ko:
kernel/drivers/spmi/hisi-spmi-controller.ko: kernel/drivers/spmi/spmi.ko
kernel/drivers/hsi/clients/hsi_char.ko: kernel/drivers/hsi/hsi.ko
kernel/drivers/hsi/hsi.ko:
kernel/drivers/slimbus/slimbus.ko:
kernel/drivers/slimbus/slim-qcom-ctrl.ko: kernel/drivers/slimbus/slimbus.ko
kernel/drivers/atm/zatm.ko: kernel/drivers/atm/uPD98402.ko kernel/net/atm/atm.ko
kernel/drivers/atm/uPD98402.ko: kernel/net/atm/atm.ko
kernel/drivers/atm/nicstar.ko: kernel/net/atm/atm.ko
kernel/drivers/atm/ambassador.ko: kernel/net/atm/atm.ko
kernel/drivers/atm/horizon.ko: kernel/net/atm/atm.ko
kernel/drivers/atm/iphase.ko: kernel/drivers/atm/suni.ko kernel/net/atm/atm.ko
kernel/drivers/atm/suni.ko: kernel/net/atm/atm.ko
kernel/drivers/atm/fore_200e.ko: kernel/net/atm/atm.ko
kernel/drivers/atm/eni.ko: kernel/drivers/atm/suni.ko kernel/net/atm/atm.ko
kernel/drivers/atm/idt77252.ko: kernel/drivers/atm/suni.ko kernel/net/atm/atm.ko
kernel/drivers/atm/solos-pci.ko: kernel/net/atm/atm.ko
kernel/drivers/atm/adummy.ko: kernel/net/atm/atm.ko
kernel/drivers/atm/atmtcp.ko: kernel/net/atm/atm.ko
kernel/drivers/atm/firestream.ko: kernel/net/atm/atm.ko
kernel/drivers/atm/lanai.ko: kernel/net/atm/atm.ko
kernel/drivers/atm/he.ko: kernel/drivers/atm/suni.ko kernel/net/atm/atm.ko
kernel/drivers/uio/uio.ko:
kernel/drivers/uio/uio_cif.ko: kernel/drivers/uio/uio.ko
kernel/drivers/uio/uio_pdrv_genirq.ko: kernel/drivers/uio/uio.ko
kernel/drivers/uio/uio_dmem_genirq.ko: kernel/drivers/uio/uio.ko
kernel/drivers/uio/uio_aec.ko: kernel/drivers/uio/uio.ko
kernel/drivers/uio/uio_sercos3.ko: kernel/drivers/uio/uio.ko
kernel/drivers/uio/uio_pci_generic.ko: kernel/drivers/uio/uio.ko
kernel/drivers/uio/uio_netx.ko: kernel/drivers/uio/uio.ko
kernel/drivers/uio/uio_pruss.ko: kernel/drivers/uio/uio.ko
kernel/drivers/uio/uio_mf624.ko: kernel/drivers/uio/uio.ko
kernel/drivers/uio/uio_hv_generic.ko: kernel/drivers/uio/uio.ko kernel/drivers/hv/hv_vmbus.ko
kernel/drivers/uio/uio_dfl.ko: kernel/drivers/fpga/dfl.ko kernel/drivers/fpga/fpga-region.ko kernel/drivers/fpga/fpga-bridge.ko kernel/drivers/fpga/fpga-mgr.ko kernel/drivers/uio/uio.ko
kernel/drivers/pcmcia/pcmcia_core.ko:
kernel/drivers/pcmcia/pcmcia.ko: kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/pcmcia/pcmcia_rsrc.ko: kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/pcmcia/yenta_socket.ko: kernel/drivers/pcmcia/pcmcia_rsrc.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/pcmcia/pd6729.ko: kernel/drivers/pcmcia/pcmcia_rsrc.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/pcmcia/i82092.ko: kernel/drivers/pcmcia/pcmcia_rsrc.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/block/aoe/aoe.ko:
kernel/drivers/block/paride/paride.ko: kernel/drivers/parport/parport.ko
kernel/drivers/block/paride/aten.ko: kernel/drivers/block/paride/paride.ko kernel/drivers/parport/parport.ko
kernel/drivers/block/paride/bpck.ko: kernel/drivers/block/paride/paride.ko kernel/drivers/parport/parport.ko
kernel/drivers/block/paride/comm.ko: kernel/drivers/block/paride/paride.ko kernel/drivers/parport/parport.ko
kernel/drivers/block/paride/dstr.ko: kernel/drivers/block/paride/paride.ko kernel/drivers/parport/parport.ko
kernel/drivers/block/paride/kbic.ko: kernel/drivers/block/paride/paride.ko kernel/drivers/parport/parport.ko
kernel/drivers/block/paride/epat.ko: kernel/drivers/block/paride/paride.ko kernel/drivers/parport/parport.ko
kernel/drivers/block/paride/epia.ko: kernel/drivers/block/paride/paride.ko kernel/drivers/parport/parport.ko
kernel/drivers/block/paride/frpw.ko: kernel/drivers/block/paride/paride.ko kernel/drivers/parport/parport.ko
kernel/drivers/block/paride/friq.ko: kernel/drivers/block/paride/paride.ko kernel/drivers/parport/parport.ko
kernel/drivers/block/paride/fit2.ko: kernel/drivers/block/paride/paride.ko kernel/drivers/parport/parport.ko
kernel/drivers/block/paride/fit3.ko: kernel/drivers/block/paride/paride.ko kernel/drivers/parport/parport.ko
kernel/drivers/block/paride/on20.ko: kernel/drivers/block/paride/paride.ko kernel/drivers/parport/parport.ko
kernel/drivers/block/paride/on26.ko: kernel/drivers/block/paride/paride.ko kernel/drivers/parport/parport.ko
kernel/drivers/block/paride/ktti.ko: kernel/drivers/block/paride/paride.ko kernel/drivers/parport/parport.ko
kernel/drivers/block/paride/pd.ko: kernel/drivers/block/paride/paride.ko kernel/drivers/parport/parport.ko
kernel/drivers/block/paride/pcd.ko: kernel/drivers/block/paride/paride.ko kernel/drivers/parport/parport.ko
kernel/drivers/block/paride/pf.ko: kernel/drivers/block/paride/paride.ko kernel/drivers/parport/parport.ko
kernel/drivers/block/paride/pt.ko: kernel/drivers/block/paride/paride.ko kernel/drivers/parport/parport.ko
kernel/drivers/block/paride/pg.ko: kernel/drivers/block/paride/paride.ko kernel/drivers/parport/parport.ko
kernel/drivers/input/gameport/gameport.ko:
kernel/drivers/input/gameport/emu10k1-gp.ko: kernel/drivers/input/gameport/gameport.ko
kernel/drivers/input/gameport/fm801-gp.ko: kernel/drivers/input/gameport/gameport.ko
kernel/drivers/input/gameport/lightning.ko: kernel/drivers/input/gameport/gameport.ko
kernel/drivers/input/gameport/ns558.ko: kernel/drivers/input/gameport/gameport.ko
kernel/drivers/w1/masters/matrox_w1.ko: kernel/drivers/w1/wire.ko
kernel/drivers/w1/masters/ds2490.ko: kernel/drivers/w1/wire.ko
kernel/drivers/w1/masters/ds2482.ko: kernel/drivers/w1/wire.ko
kernel/drivers/w1/masters/ds1wm.ko: kernel/drivers/w1/wire.ko
kernel/drivers/w1/masters/w1-gpio.ko: kernel/drivers/w1/wire.ko
kernel/drivers/w1/masters/sgi_w1.ko: kernel/drivers/w1/wire.ko
kernel/drivers/w1/slaves/w1_therm.ko: kernel/drivers/w1/wire.ko
kernel/drivers/w1/slaves/w1_smem.ko: kernel/drivers/w1/wire.ko
kernel/drivers/w1/slaves/w1_ds2405.ko: kernel/drivers/w1/wire.ko
kernel/drivers/w1/slaves/w1_ds2408.ko: kernel/drivers/w1/wire.ko
kernel/drivers/w1/slaves/w1_ds2413.ko: kernel/drivers/w1/wire.ko
kernel/drivers/w1/slaves/w1_ds2406.ko: kernel/drivers/w1/wire.ko
kernel/drivers/w1/slaves/w1_ds2423.ko: kernel/drivers/w1/wire.ko
kernel/drivers/w1/slaves/w1_ds2430.ko: kernel/drivers/w1/wire.ko
kernel/drivers/w1/slaves/w1_ds2431.ko: kernel/drivers/w1/wire.ko
kernel/drivers/w1/slaves/w1_ds2805.ko: kernel/drivers/w1/wire.ko
kernel/drivers/w1/slaves/w1_ds2433.ko: kernel/drivers/w1/wire.ko
kernel/drivers/w1/slaves/w1_ds2438.ko: kernel/drivers/w1/wire.ko
kernel/drivers/w1/slaves/w1_ds250x.ko: kernel/drivers/w1/wire.ko
kernel/drivers/w1/slaves/w1_ds2780.ko: kernel/drivers/w1/wire.ko
kernel/drivers/w1/slaves/w1_ds2781.ko: kernel/drivers/w1/wire.ko
kernel/drivers/w1/slaves/w1_ds28e04.ko: kernel/drivers/w1/wire.ko
kernel/drivers/w1/slaves/w1_ds28e17.ko: kernel/drivers/w1/wire.ko
kernel/drivers/w1/wire.ko:
kernel/drivers/bluetooth/hci_vhci.ko: kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/crypto/ecc.ko
kernel/drivers/bluetooth/hci_uart.ko: kernel/drivers/bluetooth/btqca.ko kernel/drivers/bluetooth/btrtl.ko kernel/drivers/bluetooth/btbcm.ko kernel/drivers/bluetooth/btintel.ko kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/crypto/ecc.ko
kernel/drivers/bluetooth/bcm203x.ko: kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/crypto/ecc.ko
kernel/drivers/bluetooth/bpa10x.ko: kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/crypto/ecc.ko
kernel/drivers/bluetooth/bfusb.ko: kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/crypto/ecc.ko
kernel/drivers/bluetooth/dtl1_cs.ko: kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko kernel/crypto/ecc.ko
kernel/drivers/bluetooth/bt3c_cs.ko: kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko kernel/crypto/ecc.ko
kernel/drivers/bluetooth/bluecard_cs.ko: kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko kernel/crypto/ecc.ko
kernel/drivers/bluetooth/btusb.ko: kernel/drivers/bluetooth/btrtl.ko kernel/drivers/bluetooth/btbcm.ko kernel/drivers/bluetooth/btintel.ko kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/crypto/ecc.ko
kernel/drivers/bluetooth/btsdio.ko: kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/crypto/ecc.ko
kernel/drivers/bluetooth/btintel.ko: kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/crypto/ecc.ko
kernel/drivers/bluetooth/ath3k.ko: kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/crypto/ecc.ko
kernel/drivers/bluetooth/btmrvl.ko: kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/crypto/ecc.ko
kernel/drivers/bluetooth/btmrvl_sdio.ko: kernel/drivers/bluetooth/btmrvl.ko kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/crypto/ecc.ko
kernel/drivers/bluetooth/btmtksdio.ko: kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/crypto/ecc.ko
kernel/drivers/bluetooth/btmtkuart.ko: kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/crypto/ecc.ko
kernel/drivers/bluetooth/btbcm.ko: kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/crypto/ecc.ko
kernel/drivers/bluetooth/btrtl.ko: kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/crypto/ecc.ko
kernel/drivers/bluetooth/btqca.ko: kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/crypto/ecc.ko
kernel/drivers/bluetooth/virtio_bt.ko: kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/crypto/ecc.ko
kernel/drivers/bluetooth/hci_nokia.ko: kernel/drivers/bluetooth/hci_uart.ko kernel/drivers/bluetooth/btqca.ko kernel/drivers/bluetooth/btrtl.ko kernel/drivers/bluetooth/btbcm.ko kernel/drivers/bluetooth/btintel.ko kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/crypto/ecc.ko
kernel/drivers/bluetooth/btrsi.ko: kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/crypto/ecc.ko
kernel/drivers/memstick/core/memstick.ko:
kernel/drivers/memstick/core/ms_block.ko: kernel/drivers/memstick/core/memstick.ko
kernel/drivers/memstick/core/mspro_block.ko: kernel/drivers/memstick/core/memstick.ko
kernel/drivers/memstick/host/tifm_ms.ko: kernel/drivers/memstick/core/memstick.ko kernel/drivers/misc/tifm_core.ko
kernel/drivers/memstick/host/jmb38x_ms.ko: kernel/drivers/memstick/core/memstick.ko
kernel/drivers/memstick/host/r592.ko: kernel/drivers/memstick/core/memstick.ko
kernel/drivers/memstick/host/rtsx_pci_ms.ko: kernel/drivers/memstick/core/memstick.ko kernel/drivers/misc/cardreader/rtsx_pci.ko
kernel/drivers/memstick/host/rtsx_usb_ms.ko: kernel/drivers/memstick/core/memstick.ko kernel/drivers/misc/cardreader/rtsx_usb.ko
kernel/drivers/infiniband/core/ib_core.ko:
kernel/drivers/infiniband/core/ib_cm.ko: kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/core/iw_cm.ko: kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/core/rdma_cm.ko: kernel/drivers/infiniband/core/iw_cm.ko kernel/drivers/infiniband/core/ib_cm.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/core/ib_umad.ko: kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/core/ib_uverbs.ko: kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/core/rdma_ucm.ko: kernel/drivers/infiniband/core/ib_uverbs.ko kernel/drivers/infiniband/core/rdma_cm.ko kernel/drivers/infiniband/core/iw_cm.ko kernel/drivers/infiniband/core/ib_cm.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/hw/mthca/ib_mthca.ko: kernel/drivers/infiniband/core/ib_uverbs.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/hw/qib/ib_qib.ko: kernel/drivers/infiniband/sw/rdmavt/rdmavt.ko kernel/drivers/infiniband/core/ib_uverbs.ko kernel/drivers/dca/dca.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/hw/cxgb4/iw_cxgb4.ko: kernel/drivers/infiniband/core/ib_uverbs.ko kernel/drivers/net/ethernet/chelsio/cxgb4/cxgb4.ko kernel/net/tls/tls.ko kernel/drivers/net/ethernet/chelsio/libcxgb/libcxgb.ko kernel/drivers/infiniband/core/rdma_cm.ko kernel/drivers/infiniband/core/iw_cm.ko kernel/drivers/infiniband/core/ib_cm.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/hw/efa/efa.ko: kernel/drivers/infiniband/core/ib_uverbs.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/hw/irdma/irdma.ko: kernel/drivers/net/ethernet/intel/i40e/i40e.ko kernel/drivers/net/ethernet/intel/ice/ice.ko kernel/drivers/infiniband/core/ib_uverbs.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/hw/mlx4/mlx4_ib.ko: kernel/drivers/infiniband/core/ib_uverbs.ko kernel/drivers/net/ethernet/mellanox/mlx4/mlx4_core.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/hw/mlx5/mlx5_ib.ko: kernel/drivers/infiniband/core/ib_uverbs.ko kernel/drivers/net/ethernet/mellanox/mlx5/core/mlx5_core.ko kernel/drivers/net/ethernet/mellanox/mlxfw/mlxfw.ko kernel/net/psample/psample.ko kernel/net/tls/tls.ko kernel/drivers/pci/controller/pci-hyperv-intf.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/hw/ocrdma/ocrdma.ko: kernel/drivers/net/ethernet/emulex/benet/be2net.ko kernel/drivers/infiniband/core/ib_uverbs.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/hw/vmw_pvrdma/vmw_pvrdma.ko: kernel/drivers/infiniband/core/ib_uverbs.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/hw/usnic/usnic_verbs.ko: kernel/drivers/net/ethernet/cisco/enic/enic.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/hw/hfi1/hfi1.ko: kernel/drivers/infiniband/sw/rdmavt/rdmavt.ko kernel/drivers/infiniband/core/ib_uverbs.ko kernel/drivers/i2c/algos/i2c-algo-bit.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/hw/qedr/qedr.ko: kernel/drivers/net/ethernet/qlogic/qede/qede.ko kernel/drivers/infiniband/core/ib_uverbs.ko kernel/drivers/net/ethernet/qlogic/qed/qed.ko kernel/lib/crc8.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/hw/bnxt_re/bnxt_re.ko: kernel/drivers/net/ethernet/broadcom/bnxt/bnxt_en.ko kernel/drivers/infiniband/core/ib_uverbs.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/ulp/ipoib/ib_ipoib.ko: kernel/drivers/infiniband/core/ib_cm.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/ulp/srp/ib_srp.ko: kernel/drivers/scsi/scsi_transport_srp.ko kernel/drivers/infiniband/core/rdma_cm.ko kernel/drivers/infiniband/core/iw_cm.ko kernel/drivers/infiniband/core/ib_cm.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/ulp/srpt/ib_srpt.ko: kernel/drivers/target/target_core_mod.ko kernel/drivers/infiniband/core/rdma_cm.ko kernel/drivers/infiniband/core/iw_cm.ko kernel/drivers/infiniband/core/ib_cm.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/ulp/iser/ib_iser.ko: kernel/drivers/scsi/libiscsi.ko kernel/drivers/scsi/scsi_transport_iscsi.ko kernel/drivers/infiniband/core/rdma_cm.ko kernel/drivers/infiniband/core/iw_cm.ko kernel/drivers/infiniband/core/ib_cm.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/ulp/isert/ib_isert.ko: kernel/drivers/target/iscsi/iscsi_target_mod.ko kernel/drivers/target/target_core_mod.ko kernel/drivers/infiniband/core/rdma_cm.ko kernel/drivers/infiniband/core/iw_cm.ko kernel/drivers/infiniband/core/ib_cm.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/ulp/opa_vnic/opa_vnic.ko: kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/ulp/rtrs/rtrs-core.ko: kernel/drivers/infiniband/core/rdma_cm.ko kernel/drivers/infiniband/core/iw_cm.ko kernel/drivers/infiniband/core/ib_cm.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/ulp/rtrs/rtrs-client.ko: kernel/drivers/infiniband/ulp/rtrs/rtrs-core.ko kernel/drivers/infiniband/core/rdma_cm.ko kernel/drivers/infiniband/core/iw_cm.ko kernel/drivers/infiniband/core/ib_cm.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/ulp/rtrs/rtrs-server.ko: kernel/drivers/infiniband/ulp/rtrs/rtrs-core.ko kernel/drivers/infiniband/core/rdma_cm.ko kernel/drivers/infiniband/core/iw_cm.ko kernel/drivers/infiniband/core/ib_cm.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/sw/rdmavt/rdmavt.ko: kernel/drivers/infiniband/core/ib_uverbs.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/sw/rxe/rdma_rxe.ko: kernel/drivers/infiniband/core/ib_uverbs.ko kernel/net/ipv6/ip6_udp_tunnel.ko kernel/net/ipv4/udp_tunnel.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/infiniband/sw/siw/siw.ko: kernel/lib/libcrc32c.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/drivers/dca/dca.ko:
kernel/drivers/hid/hid.ko:
kernel/drivers/hid/uhid.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-generic.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-a4tech.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-accutouch.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-alps.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-axff.ko: kernel/drivers/hid/hid.ko kernel/drivers/input/ff-memless.ko
kernel/drivers/hid/hid-apple.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-appleir.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-creative-sb0540.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-asus.ko: kernel/drivers/hid/usbhid/usbhid.ko kernel/drivers/hid/hid.ko kernel/drivers/platform/x86/asus-wmi.ko kernel/drivers/input/sparse-keymap.ko kernel/drivers/acpi/video.ko kernel/drivers/acpi/platform_profile.ko kernel/drivers/platform/x86/wmi.ko
kernel/drivers/hid/hid-aureal.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-belkin.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-betopff.ko: kernel/drivers/hid/hid.ko kernel/drivers/input/ff-memless.ko
kernel/drivers/hid/hid-bigbenff.ko: kernel/drivers/hid/hid.ko kernel/drivers/input/ff-memless.ko
kernel/drivers/hid/hid-cherry.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-chicony.ko: kernel/drivers/hid/usbhid/usbhid.ko kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-cmedia.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-corsair.ko: kernel/drivers/hid/usbhid/usbhid.ko kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-cougar.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-cp2112.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-cypress.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-dr.ko: kernel/drivers/hid/hid.ko kernel/drivers/input/ff-memless.ko
kernel/drivers/hid/hid-emsff.ko: kernel/drivers/hid/hid.ko kernel/drivers/input/ff-memless.ko
kernel/drivers/hid/hid-elan.ko: kernel/drivers/hid/usbhid/usbhid.ko kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-elecom.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-elo.ko: kernel/drivers/hid/usbhid/usbhid.ko kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-ezkey.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-ft260.ko: kernel/drivers/hid/usbhid/usbhid.ko kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-gembird.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-gfrm.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-glorious.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-google-hammer.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-vivaldi.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-gt683r.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-gyration.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-holtek-kbd.ko: kernel/drivers/hid/usbhid/usbhid.ko kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-holtek-mouse.ko: kernel/drivers/hid/usbhid/usbhid.ko kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-holtekff.ko: kernel/drivers/hid/hid.ko kernel/drivers/input/ff-memless.ko
kernel/drivers/hid/hid-hyperv.ko: kernel/drivers/hid/hid.ko kernel/drivers/hv/hv_vmbus.ko
kernel/drivers/hid/hid-icade.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-ite.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-jabra.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-kensington.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-keytouch.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-kye.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-lcpower.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-lenovo.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-logitech.ko: kernel/drivers/hid/usbhid/usbhid.ko kernel/drivers/hid/hid.ko kernel/drivers/input/ff-memless.ko
kernel/drivers/hid/hid-lg-g15.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-logitech-dj.ko: kernel/drivers/hid/usbhid/usbhid.ko kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-logitech-hidpp.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-macally.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-magicmouse.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-maltron.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-mcp2221.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-mf.ko: kernel/drivers/hid/hid.ko kernel/drivers/input/ff-memless.ko
kernel/drivers/hid/hid-microsoft.ko: kernel/drivers/hid/hid.ko kernel/drivers/input/ff-memless.ko
kernel/drivers/hid/hid-monterey.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-multitouch.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-nti.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-ntrig.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-ortek.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-prodikeys.ko: kernel/drivers/hid/usbhid/usbhid.ko kernel/drivers/hid/hid.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/drivers/hid/hid-pl.ko: kernel/drivers/hid/hid.ko kernel/drivers/input/ff-memless.ko
kernel/drivers/hid/hid-penmount.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-petalynx.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-picolcd.ko: kernel/drivers/hid/hid.ko kernel/drivers/media/rc/rc-core.ko kernel/drivers/video/fbdev/core/fb_sys_fops.ko kernel/drivers/video/fbdev/core/syscopyarea.ko kernel/drivers/video/fbdev/core/sysfillrect.ko kernel/drivers/video/fbdev/core/sysimgblt.ko kernel/drivers/video/backlight/lcd.ko
kernel/drivers/hid/hid-plantronics.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-playstation.ko: kernel/drivers/hid/hid.ko kernel/drivers/input/ff-memless.ko
kernel/drivers/hid/hid-primax.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-redragon.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-retrode.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-roccat.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-roccat-common.ko:
kernel/drivers/hid/hid-roccat-arvo.ko: kernel/drivers/hid/hid-roccat.ko kernel/drivers/hid/hid-roccat-common.ko kernel/drivers/hid/usbhid/usbhid.ko kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-roccat-isku.ko: kernel/drivers/hid/hid-roccat.ko kernel/drivers/hid/hid-roccat-common.ko kernel/drivers/hid/usbhid/usbhid.ko kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-roccat-kone.ko: kernel/drivers/hid/hid-roccat.ko kernel/drivers/hid/usbhid/usbhid.ko kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-roccat-koneplus.ko: kernel/drivers/hid/hid-roccat.ko kernel/drivers/hid/hid-roccat-common.ko kernel/drivers/hid/usbhid/usbhid.ko kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-roccat-konepure.ko: kernel/drivers/hid/hid-roccat.ko kernel/drivers/hid/hid-roccat-common.ko kernel/drivers/hid/usbhid/usbhid.ko kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-roccat-kovaplus.ko: kernel/drivers/hid/hid-roccat.ko kernel/drivers/hid/hid-roccat-common.ko kernel/drivers/hid/usbhid/usbhid.ko kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-roccat-lua.ko: kernel/drivers/hid/hid-roccat-common.ko kernel/drivers/hid/usbhid/usbhid.ko kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-roccat-pyra.ko: kernel/drivers/hid/hid-roccat.ko kernel/drivers/hid/hid-roccat-common.ko kernel/drivers/hid/usbhid/usbhid.ko kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-roccat-ryos.ko: kernel/drivers/hid/hid-roccat.ko kernel/drivers/hid/hid-roccat-common.ko kernel/drivers/hid/usbhid/usbhid.ko kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-roccat-savu.ko: kernel/drivers/hid/hid-roccat.ko kernel/drivers/hid/hid-roccat-common.ko kernel/drivers/hid/usbhid/usbhid.ko kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-rmi.ko: kernel/drivers/hid/hid.ko kernel/drivers/input/rmi4/rmi_core.ko kernel/drivers/media/common/videobuf2/videobuf2-vmalloc.ko kernel/drivers/media/common/videobuf2/videobuf2-memops.ko kernel/drivers/media/common/videobuf2/videobuf2-v4l2.ko kernel/drivers/media/common/videobuf2/videobuf2-common.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/drivers/hid/hid-saitek.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-samsung.ko: kernel/drivers/hid/usbhid/usbhid.ko kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-semitek.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-sjoy.ko: kernel/drivers/hid/hid.ko kernel/drivers/input/ff-memless.ko
kernel/drivers/hid/hid-sony.ko: kernel/drivers/hid/usbhid/usbhid.ko kernel/drivers/hid/hid.ko kernel/drivers/input/ff-memless.ko
kernel/drivers/hid/hid-speedlink.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-steam.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-steelseries.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-sunplus.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-gaff.ko: kernel/drivers/hid/hid.ko kernel/drivers/input/ff-memless.ko
kernel/drivers/hid/hid-tmff.ko: kernel/drivers/hid/hid.ko kernel/drivers/input/ff-memless.ko
kernel/drivers/hid/hid-thrustmaster.ko: kernel/drivers/hid/usbhid/usbhid.ko kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-tivo.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-topseed.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-twinhan.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-u2fzero.ko: kernel/drivers/hid/usbhid/usbhid.ko kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-uclogic.ko: kernel/drivers/hid/usbhid/usbhid.ko kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-udraw-ps3.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-led.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-xinmo.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-zpff.ko: kernel/drivers/hid/hid.ko kernel/drivers/input/ff-memless.ko
kernel/drivers/hid/hid-zydacron.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-viewsonic.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/wacom.ko: kernel/drivers/hid/usbhid/usbhid.ko kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-waltop.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-wiimote.ko: kernel/drivers/hid/hid.ko kernel/drivers/input/ff-memless.ko
kernel/drivers/hid/hid-sensor-hub.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/hid-sensor-custom.ko: kernel/drivers/hid/hid-sensor-hub.ko kernel/drivers/hid/hid.ko
kernel/drivers/hid/usbhid/usbhid.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/usbhid/usbkbd.ko:
kernel/drivers/hid/usbhid/usbmouse.ko:
kernel/drivers/hid/i2c-hid/i2c-hid.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/i2c-hid/i2c-hid-acpi.ko: kernel/drivers/hid/i2c-hid/i2c-hid.ko kernel/drivers/hid/hid.ko
kernel/drivers/hid/intel-ish-hid/intel-ishtp.ko:
kernel/drivers/hid/intel-ish-hid/intel-ish-ipc.ko: kernel/drivers/hid/intel-ish-hid/intel-ishtp.ko
kernel/drivers/hid/intel-ish-hid/intel-ishtp-hid.ko: kernel/drivers/hid/intel-ish-hid/intel-ishtp.ko kernel/drivers/hid/hid.ko
kernel/drivers/hid/intel-ish-hid/intel-ishtp-loader.ko: kernel/drivers/hid/intel-ish-hid/intel-ishtp.ko
kernel/drivers/hid/amd-sfh-hid/amd_sfh.ko: kernel/drivers/hid/hid.ko
kernel/drivers/hid/surface-hid/surface_hid_core.ko: kernel/drivers/platform/surface/aggregator/surface_aggregator.ko kernel/drivers/hid/hid.ko
kernel/drivers/hid/surface-hid/surface_hid.ko: kernel/drivers/hid/surface-hid/surface_hid_core.ko kernel/drivers/platform/surface/aggregator/surface_aggregator.ko kernel/drivers/hid/hid.ko
kernel/drivers/hid/surface-hid/surface_kbd.ko: kernel/drivers/hid/surface-hid/surface_hid_core.ko kernel/drivers/platform/surface/aggregator/surface_aggregator.ko kernel/drivers/hid/hid.ko
kernel/drivers/ssb/ssb.ko:
kernel/drivers/bcma/bcma.ko:
kernel/drivers/vhost/vhost_net.ko: kernel/drivers/vhost/vhost.ko kernel/drivers/vhost/vhost_iotlb.ko kernel/drivers/net/tap.ko
kernel/drivers/vhost/vhost_scsi.ko: kernel/drivers/vhost/vhost.ko kernel/drivers/vhost/vhost_iotlb.ko kernel/drivers/target/target_core_mod.ko
kernel/drivers/vhost/vhost_vsock.ko: kernel/net/vmw_vsock/vmw_vsock_virtio_transport_common.ko kernel/drivers/vhost/vhost.ko kernel/drivers/vhost/vhost_iotlb.ko kernel/net/vmw_vsock/vsock.ko
kernel/drivers/vhost/vringh.ko: kernel/drivers/vhost/vhost_iotlb.ko
kernel/drivers/vhost/vhost_vdpa.ko: kernel/drivers/vhost/vhost.ko kernel/drivers/vhost/vhost_iotlb.ko kernel/drivers/vdpa/vdpa.ko
kernel/drivers/vhost/vhost.ko: kernel/drivers/vhost/vhost_iotlb.ko
kernel/drivers/vhost/vhost_iotlb.ko:
kernel/drivers/greybus/greybus.ko:
kernel/drivers/greybus/gb-es2.ko: kernel/drivers/greybus/greybus.ko
kernel/drivers/comedi/comedi_pci.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/comedi_pcmcia.ko: kernel/drivers/comedi/comedi.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/comedi/comedi_usb.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/comedi.ko:
kernel/drivers/comedi/kcomedilib/kcomedilib.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/comedi_8254.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/comedi_isadma.ko:
kernel/drivers/comedi/drivers/comedi_bond.ko: kernel/drivers/comedi/kcomedilib/kcomedilib.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/comedi_test.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/comedi_parport.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/amplc_dio200.ko: kernel/drivers/comedi/drivers/amplc_dio200_common.ko kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/amplc_pc236.ko: kernel/drivers/comedi/drivers/amplc_pc236_common.ko kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/amplc_pc263.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/pcl711.ko: kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/pcl724.ko: kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/pcl726.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/pcl730.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/pcl812.ko: kernel/drivers/comedi/drivers/comedi_isadma.ko kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/pcl816.ko: kernel/drivers/comedi/drivers/comedi_isadma.ko kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/pcl818.ko: kernel/drivers/comedi/drivers/comedi_isadma.ko kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/pcm3724.ko: kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/rti800.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/rti802.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/dac02.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/das16m1.ko: kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/das08_isa.ko: kernel/drivers/comedi/drivers/das08.ko kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/das16.ko: kernel/drivers/comedi/drivers/comedi_isadma.ko kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/das800.ko: kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/das1800.ko: kernel/drivers/comedi/drivers/comedi_isadma.ko kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/das6402.ko: kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/dt2801.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/dt2811.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/dt2814.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/dt2815.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/dt2817.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/dt282x.ko: kernel/drivers/comedi/drivers/comedi_isadma.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/dmm32at.ko: kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/fl512.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/aio_aio12_8.ko: kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/aio_iiro_16.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/ii_pci20kc.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/c6xdigio.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/mpc624.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/adq12b.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/ni_at_a2150.ko: kernel/drivers/comedi/drivers/comedi_isadma.ko kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/ni_at_ao.ko: kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/ni_atmio.ko: kernel/drivers/comedi/drivers/ni_routing.ko kernel/drivers/comedi/drivers/ni_tio.ko kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/ni_atmio16d.ko: kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/ni_labpc.ko: kernel/drivers/comedi/drivers/ni_labpc_common.ko kernel/drivers/comedi/drivers/ni_labpc_isadma.ko kernel/drivers/comedi/drivers/comedi_isadma.ko kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/pcmad.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/pcmda12.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/pcmmio.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/pcmuio.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/multiq3.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/s526.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/8255_pci.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/addi_watchdog.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/addi_apci_1032.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/addi_apci_1500.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/addi_apci_1516.ko: kernel/drivers/comedi/drivers/addi_watchdog.ko kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/addi_apci_1564.ko: kernel/drivers/comedi/drivers/addi_watchdog.ko kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/addi_apci_16xx.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/addi_apci_2032.ko: kernel/drivers/comedi/drivers/addi_watchdog.ko kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/addi_apci_2200.ko: kernel/drivers/comedi/drivers/addi_watchdog.ko kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/addi_apci_3120.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/addi_apci_3501.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/addi_apci_3xxx.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/adl_pci6208.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/adl_pci7x3x.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/adl_pci8164.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/adl_pci9111.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/adl_pci9118.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/adv_pci1710.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/adv_pci1720.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/adv_pci1723.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/adv_pci1724.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/adv_pci1760.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/adv_pci_dio.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/amplc_dio200_pci.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/drivers/amplc_dio200_common.ko kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/amplc_pci236.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/drivers/amplc_pc236_common.ko kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/amplc_pci263.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/amplc_pci224.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/amplc_pci230.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/contec_pci_dio.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/das08_pci.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/drivers/das08.ko kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/dt3000.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/dyna_pci10xx.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/gsc_hpdi.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/icp_multi.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/daqboard2000.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/jr3_pci.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/ke_counter.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/cb_pcidas64.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/cb_pcidas.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/cb_pcidda.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/cb_pcimdas.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/cb_pcimdda.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/me4000.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/me_daq.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/ni_6527.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/ni_65xx.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/ni_660x.ko: kernel/drivers/comedi/drivers/ni_tiocmd.ko kernel/drivers/comedi/drivers/mite.ko kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/drivers/ni_routing.ko kernel/drivers/comedi/drivers/ni_tio.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/ni_670x.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/ni_labpc_pci.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/drivers/ni_labpc_common.ko kernel/drivers/comedi/drivers/ni_labpc_isadma.ko kernel/drivers/comedi/drivers/comedi_isadma.ko kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/ni_pcidio.ko: kernel/drivers/comedi/drivers/mite.ko kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/ni_pcimio.ko: kernel/drivers/comedi/drivers/ni_tiocmd.ko kernel/drivers/comedi/drivers/mite.ko kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/drivers/ni_routing.ko kernel/drivers/comedi/drivers/ni_tio.ko kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/rtd520.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/s626.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/mf6x4.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/cb_das16_cs.ko: kernel/drivers/comedi/comedi_pcmcia.ko kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/comedi/drivers/das08_cs.ko: kernel/drivers/comedi/comedi_pcmcia.ko kernel/drivers/comedi/drivers/das08.ko kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/comedi/drivers/ni_daq_700.ko: kernel/drivers/comedi/comedi_pcmcia.ko kernel/drivers/comedi/comedi.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/comedi/drivers/ni_daq_dio24.ko: kernel/drivers/comedi/comedi_pcmcia.ko kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/comedi.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/comedi/drivers/ni_labpc_cs.ko: kernel/drivers/comedi/comedi_pcmcia.ko kernel/drivers/comedi/drivers/ni_labpc_common.ko kernel/drivers/comedi/drivers/ni_labpc_isadma.ko kernel/drivers/comedi/drivers/comedi_isadma.ko kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/comedi/drivers/ni_mio_cs.ko: kernel/drivers/comedi/comedi_pcmcia.ko kernel/drivers/comedi/drivers/ni_routing.ko kernel/drivers/comedi/drivers/ni_tio.ko kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/comedi.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/comedi/drivers/quatech_daqp_cs.ko: kernel/drivers/comedi/comedi_pcmcia.ko kernel/drivers/comedi/comedi.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko
kernel/drivers/comedi/drivers/dt9812.ko: kernel/drivers/comedi/comedi_usb.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/ni_usb6501.ko: kernel/drivers/comedi/comedi_usb.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/usbdux.ko: kernel/drivers/comedi/comedi_usb.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/usbduxfast.ko: kernel/drivers/comedi/comedi_usb.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/usbduxsigma.ko: kernel/drivers/comedi/comedi_usb.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/vmk80xx.ko: kernel/drivers/comedi/comedi_usb.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/mite.ko: kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/ni_tio.ko:
kernel/drivers/comedi/drivers/ni_tiocmd.ko: kernel/drivers/comedi/drivers/mite.ko kernel/drivers/comedi/comedi_pci.ko kernel/drivers/comedi/drivers/ni_routing.ko kernel/drivers/comedi/drivers/ni_tio.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/ni_routing.ko:
kernel/drivers/comedi/drivers/ni_labpc_common.ko: kernel/drivers/comedi/drivers/ni_labpc_isadma.ko kernel/drivers/comedi/drivers/comedi_isadma.ko kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/ni_labpc_isadma.ko: kernel/drivers/comedi/drivers/comedi_isadma.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/comedi_8255.ko: kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/8255.ko: kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/amplc_dio200_common.ko: kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/amplc_pc236_common.ko: kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/das08.ko: kernel/drivers/comedi/drivers/comedi_8255.ko kernel/drivers/comedi/drivers/comedi_8254.ko kernel/drivers/comedi/comedi.ko
kernel/drivers/comedi/drivers/tests/comedi_example_test.ko:
kernel/drivers/comedi/drivers/tests/ni_routes_test.ko: kernel/drivers/comedi/drivers/ni_routing.ko
kernel/drivers/rpmsg/rpmsg_core.ko:
kernel/drivers/rpmsg/rpmsg_char.ko: kernel/drivers/rpmsg/rpmsg_core.ko
kernel/drivers/rpmsg/rpmsg_ns.ko: kernel/drivers/rpmsg/rpmsg_core.ko
kernel/drivers/rpmsg/qcom_glink.ko: kernel/drivers/rpmsg/rpmsg_core.ko
kernel/drivers/rpmsg/qcom_glink_rpm.ko: kernel/drivers/rpmsg/qcom_glink.ko kernel/drivers/rpmsg/rpmsg_core.ko
kernel/drivers/rpmsg/virtio_rpmsg_bus.ko: kernel/drivers/rpmsg/rpmsg_ns.ko kernel/drivers/rpmsg/rpmsg_core.ko
kernel/drivers/soundwire/soundwire-bus.ko:
kernel/drivers/soundwire/soundwire-generic-allocation.ko: kernel/drivers/soundwire/soundwire-bus.ko
kernel/drivers/soundwire/soundwire-cadence.ko: kernel/drivers/soundwire/soundwire-bus.ko
kernel/drivers/soundwire/soundwire-intel.ko: kernel/drivers/soundwire/soundwire-generic-allocation.ko kernel/drivers/soundwire/soundwire-cadence.ko kernel/drivers/soundwire/soundwire-bus.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/drivers/soundwire/soundwire-qcom.ko: kernel/drivers/soundwire/soundwire-bus.ko kernel/drivers/slimbus/slimbus.ko
kernel/drivers/iio/accel/adis16201.ko: kernel/drivers/iio/imu/adis_lib.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/adis16209.ko: kernel/drivers/iio/imu/adis_lib.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/adxl372.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/adxl372_i2c.ko: kernel/drivers/iio/accel/adxl372.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/adxl372_spi.ko: kernel/drivers/iio/accel/adxl372.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/bma220_spi.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/bma400_core.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/bma400_i2c.ko: kernel/drivers/iio/accel/bma400_core.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/bma400_spi.ko: kernel/drivers/iio/accel/bma400_core.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/bmc150-accel-core.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/bmc150-accel-i2c.ko: kernel/drivers/iio/accel/bmc150-accel-core.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/bmc150-accel-spi.ko: kernel/drivers/iio/accel/bmc150-accel-core.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/bmi088-accel-core.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/bmi088-accel-spi.ko: kernel/drivers/iio/accel/bmi088-accel-core.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/da280.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/da311.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/dmard09.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/dmard10.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/fxls8962af-core.ko: kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/fxls8962af-i2c.ko: kernel/drivers/iio/accel/fxls8962af-core.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/fxls8962af-spi.ko: kernel/drivers/iio/accel/fxls8962af-core.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/hid-sensor-accel-3d.ko: kernel/drivers/iio/common/hid-sensors/hid-sensor-trigger.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/common/hid-sensors/hid-sensor-iio-common.ko kernel/drivers/hid/hid-sensor-hub.ko kernel/drivers/hid/hid.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/kxcjk-1013.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/kxsd9.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/kxsd9-spi.ko: kernel/drivers/iio/accel/kxsd9.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/kxsd9-i2c.ko: kernel/drivers/iio/accel/kxsd9.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/mc3230.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/mma7455_core.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/mma7455_i2c.ko: kernel/drivers/iio/accel/mma7455_core.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/mma7455_spi.ko: kernel/drivers/iio/accel/mma7455_core.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/mma7660.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/mma8452.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/mma9551_core.ko:
kernel/drivers/iio/accel/mma9551.ko: kernel/drivers/iio/accel/mma9551_core.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/mma9553.ko: kernel/drivers/iio/accel/mma9551_core.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/mxc4005.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/mxc6255.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/sca3000.ko: kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/sca3300.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko kernel/lib/crc8.ko
kernel/drivers/iio/accel/stk8312.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/stk8ba50.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/cros_ec_accel_legacy.ko: kernel/drivers/iio/common/cros_ec_sensors/cros_ec_sensors_core.ko kernel/drivers/platform/chrome/cros-ec-sensorhub.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/ssp_accel_sensor.ko: kernel/drivers/iio/common/ssp_sensors/ssp_iio.ko kernel/drivers/iio/common/ssp_sensors/sensorhub.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/st_accel.ko: kernel/drivers/iio/common/st_sensors/st_sensors.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/st_accel_i2c.ko: kernel/drivers/iio/common/st_sensors/st_sensors_i2c.ko kernel/drivers/iio/accel/st_accel.ko kernel/drivers/iio/common/st_sensors/st_sensors.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/accel/st_accel_spi.ko: kernel/drivers/iio/common/st_sensors/st_sensors_spi.ko kernel/drivers/iio/accel/st_accel.ko kernel/drivers/iio/common/st_sensors/st_sensors.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ad_sigma_delta.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ad7091r5.ko: kernel/drivers/iio/adc/ad7091r-base.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ad7091r-base.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ad7124.ko:
kernel/drivers/iio/adc/ad7192.ko: kernel/drivers/iio/adc/ad_sigma_delta.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ad7266.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ad7291.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ad7292.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ad7298.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ad7923.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ad7476.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ad7606_par.ko: kernel/drivers/iio/adc/ad7606.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ad7606_spi.ko: kernel/drivers/iio/adc/ad7606.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ad7606.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ad7766.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ad7768-1.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ad7780.ko: kernel/drivers/iio/adc/ad_sigma_delta.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ad7791.ko: kernel/drivers/iio/adc/ad_sigma_delta.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ad7793.ko: kernel/drivers/iio/adc/ad_sigma_delta.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ad7887.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ad7949.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ad799x.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/axp20x_adc.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/axp288_adc.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/cc10001_adc.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/da9150-gpadc.ko: kernel/drivers/mfd/da9150-core.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/dln2-adc.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko kernel/drivers/mfd/dln2.ko
kernel/drivers/iio/adc/hi8435.ko: kernel/drivers/iio/industrialio-triggered-event.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/hx711.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ina2xx-adc.ko: kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/intel_mrfld_adc.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/lp8788_adc.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ltc2471.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ltc2485.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ltc2496.ko: kernel/drivers/iio/adc/ltc2497-core.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ltc2497-core.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ltc2497.ko: kernel/drivers/iio/adc/ltc2497-core.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/max1027.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/max11100.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/max1118.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/max1241.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/max1363.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/max9611.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/mcp320x.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/mcp3422.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/mcp3911.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/mt6360-adc.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/men_z188_adc.ko: kernel/drivers/iio/industrialio.ko kernel/drivers/mcb/mcb.ko
kernel/drivers/iio/adc/mp2629_adc.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/nau7802.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/palmas_gpadc.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/qcom-spmi-adc5.ko: kernel/drivers/iio/adc/qcom-vadc-common.ko
kernel/drivers/iio/adc/qcom-spmi-iadc.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/qcom-vadc-common.ko:
kernel/drivers/iio/adc/qcom-spmi-vadc.ko:
kernel/drivers/iio/adc/stx104.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ti-adc081c.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ti-adc0832.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ti-adc084s021.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ti-adc12138.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ti-adc108s102.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ti-adc128s052.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ti-adc161s626.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ti-ads1015.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ti-ads7950.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ti-ads131e08.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ti_am335x_adc.ko: kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/mfd/ti_am335x_tscadc.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ti-tlc4541.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/ti-tsc2046.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/twl4030-madc.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/twl6030-gpadc.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/viperboard_adc.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/adc/xilinx-xadc.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/amplifiers/ad8366.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/amplifiers/hmc425a.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/buffer/industrialio-buffer-cb.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/buffer/industrialio-buffer-dma.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/buffer/industrialio-buffer-dmaengine.ko: kernel/drivers/iio/buffer/industrialio-buffer-dma.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/buffer/industrialio-hw-consumer.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko: kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/buffer/kfifo_buf.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/cdc/ad7150.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/chemical/atlas-sensor.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/chemical/atlas-ezo-sensor.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/chemical/bme680_core.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/chemical/bme680_i2c.ko: kernel/drivers/iio/chemical/bme680_core.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/chemical/bme680_spi.ko: kernel/drivers/iio/chemical/bme680_core.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/chemical/ccs811.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/chemical/ams-iaq-core.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/chemical/pms7003.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/chemical/scd30_core.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/chemical/scd30_i2c.ko: kernel/drivers/iio/chemical/scd30_core.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko kernel/lib/crc8.ko
kernel/drivers/iio/chemical/scd30_serial.ko: kernel/drivers/iio/chemical/scd30_core.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/chemical/sgp30.ko: kernel/drivers/iio/industrialio.ko kernel/lib/crc8.ko
kernel/drivers/iio/chemical/sgp40.ko: kernel/drivers/iio/industrialio.ko kernel/lib/crc8.ko
kernel/drivers/iio/chemical/sps30.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/chemical/sps30_i2c.ko: kernel/drivers/iio/chemical/sps30.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko kernel/lib/crc8.ko
kernel/drivers/iio/chemical/sps30_serial.ko: kernel/drivers/iio/chemical/sps30.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/chemical/vz89x.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/common/cros_ec_sensors/cros_ec_sensors_core.ko: kernel/drivers/platform/chrome/cros-ec-sensorhub.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/common/cros_ec_sensors/cros_ec_sensors.ko: kernel/drivers/iio/common/cros_ec_sensors/cros_ec_sensors_core.ko kernel/drivers/platform/chrome/cros-ec-sensorhub.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/common/cros_ec_sensors/cros_ec_lid_angle.ko: kernel/drivers/iio/common/cros_ec_sensors/cros_ec_sensors_core.ko kernel/drivers/platform/chrome/cros-ec-sensorhub.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/common/hid-sensors/hid-sensor-iio-common.ko: kernel/drivers/hid/hid-sensor-hub.ko kernel/drivers/hid/hid.ko
kernel/drivers/iio/common/hid-sensors/hid-sensor-trigger.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/common/hid-sensors/hid-sensor-iio-common.ko kernel/drivers/hid/hid-sensor-hub.ko kernel/drivers/hid/hid.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/common/ms_sensors/ms_sensors_i2c.ko:
kernel/drivers/iio/common/ssp_sensors/sensorhub.ko:
kernel/drivers/iio/common/ssp_sensors/ssp_iio.ko: kernel/drivers/iio/common/ssp_sensors/sensorhub.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/common/st_sensors/st_sensors_i2c.ko:
kernel/drivers/iio/common/st_sensors/st_sensors_spi.ko:
kernel/drivers/iio/common/st_sensors/st_sensors.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ad5360.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ad5380.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ad5421.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ad5624r_spi.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ad5064.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ad5504.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ad5446.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ad5449.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ad5592r-base.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ad5592r.ko: kernel/drivers/iio/dac/ad5592r-base.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ad5593r.ko: kernel/drivers/iio/dac/ad5592r-base.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ad5755.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ad5758.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ad5761.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ad5764.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ad5766.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ad5770r.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ad5791.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ad5686.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ad5686-spi.ko: kernel/drivers/iio/dac/ad5686.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ad5696-i2c.ko: kernel/drivers/iio/dac/ad5686.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ad7303.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ad8801.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/cio-dac.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ds4424.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ltc1660.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ltc2632.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/m62332.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/max517.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/mcp4725.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/mcp4922.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ti-dac082s085.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ti-dac5571.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ti-dac7311.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dac/ti-dac7612.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/dummy/iio_dummy.ko: kernel/drivers/iio/industrialio-sw-device.ko kernel/drivers/iio/industrialio-configfs.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/gyro/adis16080.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/gyro/adis16130.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/gyro/adis16136.ko: kernel/drivers/iio/imu/adis_lib.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/gyro/adis16260.ko: kernel/drivers/iio/imu/adis_lib.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/gyro/adxrs290.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/gyro/adxrs450.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/gyro/bmg160_core.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/gyro/bmg160_i2c.ko: kernel/drivers/iio/gyro/bmg160_core.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/gyro/bmg160_spi.ko: kernel/drivers/iio/gyro/bmg160_core.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/gyro/fxas21002c_core.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/gyro/fxas21002c_i2c.ko: kernel/drivers/iio/gyro/fxas21002c_core.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/gyro/fxas21002c_spi.ko: kernel/drivers/iio/gyro/fxas21002c_core.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/gyro/hid-sensor-gyro-3d.ko: kernel/drivers/iio/common/hid-sensors/hid-sensor-trigger.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/common/hid-sensors/hid-sensor-iio-common.ko kernel/drivers/hid/hid-sensor-hub.ko kernel/drivers/hid/hid.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/gyro/mpu3050.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/i2c/i2c-mux.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/gyro/itg3200.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/gyro/ssp_gyro_sensor.ko: kernel/drivers/iio/common/ssp_sensors/ssp_iio.ko kernel/drivers/iio/common/ssp_sensors/sensorhub.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/gyro/st_gyro.ko: kernel/drivers/iio/common/st_sensors/st_sensors.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/gyro/st_gyro_i2c.ko: kernel/drivers/iio/gyro/st_gyro.ko kernel/drivers/iio/common/st_sensors/st_sensors_i2c.ko kernel/drivers/iio/common/st_sensors/st_sensors.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/gyro/st_gyro_spi.ko: kernel/drivers/iio/gyro/st_gyro.ko kernel/drivers/iio/common/st_sensors/st_sensors_spi.ko kernel/drivers/iio/common/st_sensors/st_sensors.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/frequency/ad9523.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/frequency/adf4350.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/frequency/adf4371.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/health/afe4403.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/health/afe4404.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/health/max30100.ko: kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/health/max30102.ko: kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/humidity/am2315.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/humidity/dht11.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/humidity/hdc100x.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/humidity/hdc2010.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/humidity/hid-sensor-humidity.ko: kernel/drivers/iio/common/hid-sensors/hid-sensor-trigger.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/common/hid-sensors/hid-sensor-iio-common.ko kernel/drivers/hid/hid-sensor-hub.ko kernel/drivers/hid/hid.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/humidity/hts221.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/humidity/hts221_i2c.ko: kernel/drivers/iio/humidity/hts221.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/humidity/hts221_spi.ko: kernel/drivers/iio/humidity/hts221.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/humidity/htu21.ko: kernel/drivers/iio/common/ms_sensors/ms_sensors_i2c.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/humidity/si7005.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/humidity/si7020.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/imu/bmi160/bmi160_core.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/imu/bmi160/bmi160_i2c.ko: kernel/drivers/iio/imu/bmi160/bmi160_core.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/imu/bmi160/bmi160_spi.ko: kernel/drivers/iio/imu/bmi160/bmi160_core.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/imu/inv_icm42600/inv-icm42600.ko: kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/imu/inv_icm42600/inv-icm42600-i2c.ko: kernel/drivers/iio/imu/inv_icm42600/inv-icm42600.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/imu/inv_icm42600/inv-icm42600-spi.ko: kernel/drivers/iio/imu/inv_icm42600/inv-icm42600.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/imu/inv_mpu6050/inv-mpu6050.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/imu/inv_mpu6050/inv-mpu6050-i2c.ko: kernel/drivers/iio/imu/inv_mpu6050/inv-mpu6050.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/i2c/i2c-mux.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/imu/inv_mpu6050/inv-mpu6050-spi.ko: kernel/drivers/iio/imu/inv_mpu6050/inv-mpu6050.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/imu/st_lsm6dsx/st_lsm6dsx.ko: kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/imu/st_lsm6dsx/st_lsm6dsx_i2c.ko: kernel/drivers/iio/imu/st_lsm6dsx/st_lsm6dsx.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/imu/st_lsm6dsx/st_lsm6dsx_spi.ko: kernel/drivers/iio/imu/st_lsm6dsx/st_lsm6dsx.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/imu/st_lsm6dsx/st_lsm6dsx_i3c.ko: kernel/drivers/base/regmap/regmap-i3c.ko kernel/drivers/iio/imu/st_lsm6dsx/st_lsm6dsx.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/i3c/i3c.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/imu/st_lsm9ds0/st_lsm9ds0.ko: kernel/drivers/iio/magnetometer/st_magn.ko kernel/drivers/iio/accel/st_accel.ko kernel/drivers/iio/common/st_sensors/st_sensors.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/imu/st_lsm9ds0/st_lsm9ds0_i2c.ko: kernel/drivers/iio/imu/st_lsm9ds0/st_lsm9ds0.ko kernel/drivers/iio/magnetometer/st_magn.ko kernel/drivers/iio/accel/st_accel.ko kernel/drivers/iio/common/st_sensors/st_sensors.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/imu/st_lsm9ds0/st_lsm9ds0_spi.ko: kernel/drivers/iio/imu/st_lsm9ds0/st_lsm9ds0.ko kernel/drivers/iio/magnetometer/st_magn.ko kernel/drivers/iio/accel/st_accel.ko kernel/drivers/iio/common/st_sensors/st_sensors.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/imu/adis16400.ko: kernel/drivers/iio/imu/adis_lib.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/imu/adis16460.ko: kernel/drivers/iio/imu/adis_lib.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/imu/adis16475.ko: kernel/drivers/iio/imu/adis_lib.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/imu/adis16480.ko: kernel/drivers/iio/imu/adis_lib.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/imu/adis_lib.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/imu/fxos8700_core.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/imu/fxos8700_i2c.ko: kernel/drivers/iio/imu/fxos8700_core.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/imu/fxos8700_spi.ko: kernel/drivers/iio/imu/fxos8700_core.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/imu/kmx61.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/acpi-als.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/adjd_s311.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/adux1020.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/al3010.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/al3320a.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/apds9300.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/apds9960.ko: kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/as73211.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/bh1750.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/bh1780.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/cm32181.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/cm3232.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/cm3323.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/cm36651.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/cros_ec_light_prox.ko: kernel/drivers/iio/common/cros_ec_sensors/cros_ec_sensors_core.ko kernel/drivers/platform/chrome/cros-ec-sensorhub.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/gp2ap002.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/gp2ap020a00f.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/hid-sensor-als.ko: kernel/drivers/iio/common/hid-sensors/hid-sensor-trigger.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/common/hid-sensors/hid-sensor-iio-common.ko kernel/drivers/hid/hid-sensor-hub.ko kernel/drivers/hid/hid.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/hid-sensor-prox.ko: kernel/drivers/iio/common/hid-sensors/hid-sensor-trigger.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/common/hid-sensors/hid-sensor-iio-common.ko kernel/drivers/hid/hid-sensor-hub.ko kernel/drivers/hid/hid.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/iqs621-als.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/isl29018.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/isl29028.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/isl29125.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/jsa1212.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/lm3533-als.ko: kernel/drivers/iio/industrialio.ko kernel/drivers/mfd/lm3533-core.ko
kernel/drivers/iio/light/ltr501.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/lv0104cs.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/max44000.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/max44009.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/noa1305.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/opt3001.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/pa12203001.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/rpr0521.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/tsl2563.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/si1133.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/si1145.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/stk3310.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/st_uvis25_core.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/st_uvis25_i2c.ko: kernel/drivers/iio/light/st_uvis25_core.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/st_uvis25_spi.ko: kernel/drivers/iio/light/st_uvis25_core.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/tcs3414.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/tcs3472.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/tsl2583.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/tsl2591.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/tsl2772.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/tsl4531.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/us5182d.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/vcnl4000.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/vcnl4035.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/veml6030.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/veml6070.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/vl6180.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/light/zopt2201.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/magnetometer/ak8975.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/magnetometer/bmc150_magn.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/magnetometer/bmc150_magn_i2c.ko: kernel/drivers/iio/magnetometer/bmc150_magn.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/magnetometer/bmc150_magn_spi.ko: kernel/drivers/iio/magnetometer/bmc150_magn.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/magnetometer/mag3110.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/magnetometer/hid-sensor-magn-3d.ko: kernel/drivers/iio/common/hid-sensors/hid-sensor-trigger.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/common/hid-sensors/hid-sensor-iio-common.ko kernel/drivers/hid/hid-sensor-hub.ko kernel/drivers/hid/hid.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/magnetometer/mmc35240.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/magnetometer/st_magn.ko: kernel/drivers/iio/common/st_sensors/st_sensors.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/magnetometer/st_magn_i2c.ko: kernel/drivers/iio/magnetometer/st_magn.ko kernel/drivers/iio/common/st_sensors/st_sensors_i2c.ko kernel/drivers/iio/common/st_sensors/st_sensors.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/magnetometer/st_magn_spi.ko: kernel/drivers/iio/magnetometer/st_magn.ko kernel/drivers/iio/common/st_sensors/st_sensors_spi.ko kernel/drivers/iio/common/st_sensors/st_sensors.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/magnetometer/hmc5843_core.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/magnetometer/hmc5843_i2c.ko: kernel/drivers/iio/magnetometer/hmc5843_core.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/magnetometer/hmc5843_spi.ko: kernel/drivers/iio/magnetometer/hmc5843_core.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/magnetometer/rm3100-core.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/magnetometer/rm3100-i2c.ko: kernel/drivers/iio/magnetometer/rm3100-core.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/magnetometer/rm3100-spi.ko: kernel/drivers/iio/magnetometer/rm3100-core.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/magnetometer/yamaha-yas530.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/orientation/hid-sensor-incl-3d.ko: kernel/drivers/iio/common/hid-sensors/hid-sensor-trigger.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/common/hid-sensors/hid-sensor-iio-common.ko kernel/drivers/hid/hid-sensor-hub.ko kernel/drivers/hid/hid.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/orientation/hid-sensor-rotation.ko: kernel/drivers/iio/common/hid-sensors/hid-sensor-trigger.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/common/hid-sensors/hid-sensor-iio-common.ko kernel/drivers/hid/hid-sensor-hub.ko kernel/drivers/hid/hid.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/position/hid-sensor-custom-intel-hinge.ko: kernel/drivers/iio/common/hid-sensors/hid-sensor-trigger.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/common/hid-sensors/hid-sensor-iio-common.ko kernel/drivers/hid/hid-sensor-hub.ko kernel/drivers/hid/hid.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/position/iqs624-pos.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/potentiometer/ad5110.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/potentiometer/ad5272.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/potentiometer/ds1803.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/potentiometer/max5432.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/potentiometer/max5481.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/potentiometer/max5487.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/potentiometer/mcp4018.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/potentiometer/mcp4131.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/potentiometer/mcp4531.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/potentiometer/mcp41010.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/potentiometer/tpl0102.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/potentiostat/lmp91000.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/buffer/industrialio-buffer-cb.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/pressure/abp060mg.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/pressure/bmp280.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/pressure/bmp280-i2c.ko: kernel/drivers/iio/pressure/bmp280.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/pressure/bmp280-spi.ko: kernel/drivers/iio/pressure/bmp280.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/pressure/dlhl60d.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/pressure/dps310.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/pressure/cros_ec_baro.ko: kernel/drivers/iio/common/cros_ec_sensors/cros_ec_sensors_core.ko kernel/drivers/platform/chrome/cros-ec-sensorhub.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/pressure/hid-sensor-press.ko: kernel/drivers/iio/common/hid-sensors/hid-sensor-trigger.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/common/hid-sensors/hid-sensor-iio-common.ko kernel/drivers/hid/hid-sensor-hub.ko kernel/drivers/hid/hid.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/pressure/hp03.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/pressure/icp10100.ko: kernel/drivers/iio/industrialio.ko kernel/lib/crc8.ko
kernel/drivers/iio/pressure/mpl115.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/pressure/mpl115_i2c.ko: kernel/drivers/iio/pressure/mpl115.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/pressure/mpl115_spi.ko: kernel/drivers/iio/pressure/mpl115.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/pressure/mpl3115.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/pressure/ms5611_core.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/pressure/ms5611_i2c.ko: kernel/drivers/iio/pressure/ms5611_core.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/pressure/ms5611_spi.ko: kernel/drivers/iio/pressure/ms5611_core.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/pressure/ms5637.ko: kernel/drivers/iio/common/ms_sensors/ms_sensors_i2c.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/pressure/st_pressure.ko: kernel/drivers/iio/common/st_sensors/st_sensors.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/pressure/t5403.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/pressure/hp206c.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/pressure/zpa2326.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/pressure/zpa2326_i2c.ko: kernel/drivers/iio/pressure/zpa2326.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/pressure/zpa2326_spi.ko: kernel/drivers/iio/pressure/zpa2326.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/pressure/st_pressure_i2c.ko: kernel/drivers/iio/pressure/st_pressure.ko kernel/drivers/iio/common/st_sensors/st_sensors_i2c.ko kernel/drivers/iio/common/st_sensors/st_sensors.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/pressure/st_pressure_spi.ko: kernel/drivers/iio/pressure/st_pressure.ko kernel/drivers/iio/common/st_sensors/st_sensors_spi.ko kernel/drivers/iio/common/st_sensors/st_sensors.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/proximity/as3935.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/proximity/cros_ec_mkbp_proximity.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/proximity/isl29501.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/proximity/pulsedlight-lidar-lite-v2.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/proximity/mb1232.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/proximity/ping.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/proximity/rfd77402.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/proximity/srf04.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/proximity/srf08.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/proximity/sx9310.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/proximity/sx9500.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/proximity/vcnl3020.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/proximity/vl53l0x-i2c.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/resolver/ad2s90.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/resolver/ad2s1200.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/temperature/iqs620at-temp.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/temperature/ltc2983.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/temperature/hid-sensor-temperature.ko: kernel/drivers/iio/common/hid-sensors/hid-sensor-trigger.ko kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/common/hid-sensors/hid-sensor-iio-common.ko kernel/drivers/hid/hid-sensor-hub.ko kernel/drivers/hid/hid.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/temperature/maxim_thermocouple.ko: kernel/drivers/iio/buffer/industrialio-triggered-buffer.ko kernel/drivers/iio/buffer/kfifo_buf.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/temperature/max31856.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/temperature/mlx90614.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/temperature/mlx90632.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/temperature/tmp006.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/temperature/tmp007.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/temperature/tmp117.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/temperature/tsys01.ko: kernel/drivers/iio/common/ms_sensors/ms_sensors_i2c.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/temperature/tsys02d.ko: kernel/drivers/iio/common/ms_sensors/ms_sensors_i2c.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/trigger/iio-trig-hrtimer.ko: kernel/drivers/iio/industrialio-sw-trigger.ko kernel/drivers/iio/industrialio-configfs.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/trigger/iio-trig-interrupt.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/trigger/iio-trig-sysfs.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/trigger/iio-trig-loop.ko: kernel/drivers/iio/industrialio-sw-trigger.ko kernel/drivers/iio/industrialio-configfs.ko kernel/drivers/iio/industrialio.ko
kernel/drivers/iio/industrialio.ko:
kernel/drivers/iio/industrialio-configfs.ko:
kernel/drivers/iio/industrialio-sw-device.ko: kernel/drivers/iio/industrialio-configfs.ko
kernel/drivers/iio/industrialio-sw-trigger.ko: kernel/drivers/iio/industrialio-configfs.ko
kernel/drivers/iio/industrialio-triggered-event.ko: kernel/drivers/iio/industrialio.ko
kernel/drivers/ipack/devices/ipoctal.ko: kernel/drivers/ipack/ipack.ko
kernel/drivers/ipack/carriers/tpci200.ko: kernel/drivers/ipack/ipack.ko
kernel/drivers/ipack/ipack.ko:
kernel/drivers/ntb/ntb.ko:
kernel/drivers/ntb/hw/idt/ntb_hw_idt.ko: kernel/drivers/ntb/ntb.ko
kernel/drivers/ntb/hw/intel/ntb_hw_intel.ko: kernel/drivers/ntb/ntb.ko
kernel/drivers/ntb/hw/epf/ntb_hw_epf.ko: kernel/drivers/ntb/ntb.ko
kernel/drivers/ntb/hw/mscc/ntb_hw_switchtec.ko: kernel/drivers/pci/switch/switchtec.ko kernel/drivers/ntb/ntb.ko
kernel/drivers/ntb/test/ntb_pingpong.ko: kernel/drivers/ntb/ntb.ko
kernel/drivers/ntb/test/ntb_tool.ko: kernel/drivers/ntb/ntb.ko
kernel/drivers/ntb/test/ntb_perf.ko: kernel/drivers/ntb/ntb.ko
kernel/drivers/ntb/ntb_transport.ko: kernel/drivers/ntb/ntb.ko
kernel/drivers/mcb/mcb.ko:
kernel/drivers/mcb/mcb-pci.ko: kernel/drivers/mcb/mcb.ko
kernel/drivers/mcb/mcb-lpc.ko: kernel/drivers/mcb/mcb.ko
kernel/drivers/thunderbolt/thunderbolt.ko:
kernel/drivers/hwtracing/stm/stm_core.ko:
kernel/drivers/hwtracing/stm/stm_p_basic.ko: kernel/drivers/hwtracing/stm/stm_core.ko
kernel/drivers/hwtracing/stm/stm_p_sys-t.ko: kernel/drivers/hwtracing/stm/stm_core.ko
kernel/drivers/hwtracing/stm/dummy_stm.ko: kernel/drivers/hwtracing/stm/stm_core.ko
kernel/drivers/hwtracing/stm/stm_console.ko: kernel/drivers/hwtracing/stm/stm_core.ko
kernel/drivers/hwtracing/stm/stm_heartbeat.ko: kernel/drivers/hwtracing/stm/stm_core.ko
kernel/drivers/hwtracing/stm/stm_ftrace.ko: kernel/drivers/hwtracing/stm/stm_core.ko
kernel/drivers/fpga/fpga-mgr.ko:
kernel/drivers/fpga/altera-cvp.ko: kernel/drivers/fpga/fpga-mgr.ko
kernel/drivers/fpga/altera-ps-spi.ko: kernel/drivers/fpga/fpga-mgr.ko
kernel/drivers/fpga/machxo2-spi.ko: kernel/drivers/fpga/fpga-mgr.ko
kernel/drivers/fpga/xilinx-spi.ko: kernel/drivers/fpga/fpga-mgr.ko
kernel/drivers/fpga/altera-pr-ip-core.ko: kernel/drivers/fpga/fpga-mgr.ko
kernel/drivers/fpga/fpga-bridge.ko:
kernel/drivers/fpga/altera-freeze-bridge.ko: kernel/drivers/fpga/fpga-bridge.ko
kernel/drivers/fpga/xilinx-pr-decoupler.ko: kernel/drivers/fpga/fpga-bridge.ko
kernel/drivers/fpga/fpga-region.ko: kernel/drivers/fpga/fpga-bridge.ko kernel/drivers/fpga/fpga-mgr.ko
kernel/drivers/fpga/dfl.ko: kernel/drivers/fpga/fpga-region.ko kernel/drivers/fpga/fpga-bridge.ko kernel/drivers/fpga/fpga-mgr.ko
kernel/drivers/fpga/dfl-fme.ko: kernel/drivers/fpga/dfl.ko kernel/drivers/fpga/fpga-region.ko kernel/drivers/fpga/fpga-bridge.ko kernel/drivers/fpga/fpga-mgr.ko
kernel/drivers/fpga/dfl-fme-mgr.ko: kernel/drivers/fpga/fpga-mgr.ko
kernel/drivers/fpga/dfl-fme-br.ko: kernel/drivers/fpga/dfl.ko kernel/drivers/fpga/fpga-region.ko kernel/drivers/fpga/fpga-bridge.ko kernel/drivers/fpga/fpga-mgr.ko
kernel/drivers/fpga/dfl-fme-region.ko: kernel/drivers/fpga/fpga-region.ko kernel/drivers/fpga/fpga-bridge.ko kernel/drivers/fpga/fpga-mgr.ko
kernel/drivers/fpga/dfl-afu.ko: kernel/drivers/fpga/dfl.ko kernel/drivers/fpga/fpga-region.ko kernel/drivers/fpga/fpga-bridge.ko kernel/drivers/fpga/fpga-mgr.ko
kernel/drivers/fpga/dfl-n3000-nios.ko: kernel/drivers/fpga/dfl.ko kernel/drivers/fpga/fpga-region.ko kernel/drivers/fpga/fpga-bridge.ko kernel/drivers/fpga/fpga-mgr.ko
kernel/drivers/fpga/dfl-pci.ko: kernel/drivers/fpga/dfl.ko kernel/drivers/fpga/fpga-region.ko kernel/drivers/fpga/fpga-bridge.ko kernel/drivers/fpga/fpga-mgr.ko
kernel/drivers/tee/tee.ko:
kernel/drivers/tee/amdtee/amdtee.ko: kernel/drivers/tee/tee.ko kernel/drivers/crypto/ccp/ccp.ko
kernel/drivers/mux/mux-core.ko:
kernel/drivers/mux/mux-adg792a.ko: kernel/drivers/mux/mux-core.ko
kernel/drivers/mux/mux-adgs1408.ko: kernel/drivers/mux/mux-core.ko
kernel/drivers/mux/mux-gpio.ko: kernel/drivers/mux/mux-core.ko
kernel/drivers/visorbus/visorbus.ko:
kernel/drivers/siox/siox-core.ko:
kernel/drivers/siox/siox-bus-gpio.ko: kernel/drivers/siox/siox-core.ko
kernel/drivers/gnss/gnss.ko:
kernel/drivers/gnss/gnss-serial.ko: kernel/drivers/gnss/gnss.ko
kernel/drivers/gnss/gnss-mtk.ko: kernel/drivers/gnss/gnss-serial.ko kernel/drivers/gnss/gnss.ko
kernel/drivers/gnss/gnss-sirf.ko: kernel/drivers/gnss/gnss.ko
kernel/drivers/gnss/gnss-ubx.ko: kernel/drivers/gnss/gnss-serial.ko kernel/drivers/gnss/gnss.ko
kernel/drivers/counter/counter.ko:
kernel/drivers/counter/104-quad-8.ko: kernel/drivers/counter/counter.ko
kernel/drivers/counter/interrupt-cnt.ko: kernel/drivers/counter/counter.ko
kernel/drivers/counter/intel-qep.ko: kernel/drivers/counter/counter.ko
kernel/drivers/most/most_core.ko:
kernel/drivers/most/most_usb.ko: kernel/drivers/most/most_core.ko
kernel/drivers/most/most_cdev.ko: kernel/drivers/most/most_core.ko
kernel/drivers/most/most_snd.ko: kernel/drivers/most/most_core.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soundcore.ko:
kernel/sound/core/oss/snd-mixer-oss.ko: kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/core/snd.ko: kernel/sound/soundcore.ko
kernel/sound/core/snd-ctl-led.ko: kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/core/snd-hwdep.ko: kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/core/snd-timer.ko: kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/core/snd-hrtimer.ko: kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/core/snd-pcm.ko: kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/core/snd-pcm-dmaengine.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/core/snd-seq-device.ko: kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/core/snd-rawmidi.ko: kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/core/seq/snd-seq.ko: kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/core/seq/snd-seq-dummy.ko: kernel/sound/core/seq/snd-seq.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/core/seq/snd-seq-midi.ko: kernel/sound/core/seq/snd-seq-midi-event.ko kernel/sound/core/seq/snd-seq.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/core/seq/snd-seq-midi-emul.ko: kernel/sound/core/seq/snd-seq.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/core/seq/snd-seq-midi-event.ko: kernel/sound/core/seq/snd-seq.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/core/seq/snd-seq-virmidi.ko: kernel/sound/core/seq/snd-seq-midi-event.ko kernel/sound/core/seq/snd-seq.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/core/snd-compress.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/i2c/other/snd-ak4117.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/i2c/other/snd-ak4xxx-adda.ko: kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/i2c/other/snd-ak4114.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/i2c/other/snd-ak4113.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/i2c/other/snd-pt2258.ko: kernel/sound/i2c/snd-i2c.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/i2c/snd-cs8427.ko: kernel/sound/i2c/snd-i2c.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/i2c/snd-i2c.ko: kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/drivers/snd-dummy.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/drivers/snd-aloop.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/drivers/snd-virmidi.ko: kernel/sound/core/seq/snd-seq-virmidi.ko kernel/sound/core/seq/snd-seq-midi-event.ko kernel/sound/core/seq/snd-seq.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/drivers/snd-serial-u16550.ko: kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/drivers/snd-mtpav.ko: kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/drivers/snd-mts64.ko: kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/drivers/parport/parport.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/drivers/snd-portman2x4.ko: kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/drivers/parport/parport.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/drivers/opl3/snd-opl3-lib.ko: kernel/sound/core/snd-hwdep.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/drivers/opl3/snd-opl3-synth.ko: kernel/sound/core/seq/snd-seq-midi-emul.ko kernel/sound/drivers/opl3/snd-opl3-lib.ko kernel/sound/core/snd-hwdep.ko kernel/sound/core/seq/snd-seq.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/drivers/mpu401/snd-mpu401-uart.ko: kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/drivers/mpu401/snd-mpu401.ko: kernel/sound/drivers/mpu401/snd-mpu401-uart.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/drivers/vx/snd-vx-lib.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/drivers/pcsp/snd-pcsp.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/isa/sb/snd-sb-common.ko: kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/snd-ad1889.ko: kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/snd-als300.ko: kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/snd-als4000.ko: kernel/sound/isa/sb/snd-sb-common.ko kernel/sound/drivers/mpu401/snd-mpu401-uart.ko kernel/sound/drivers/opl3/snd-opl3-lib.ko kernel/sound/core/snd-hwdep.ko kernel/drivers/input/gameport/gameport.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/snd-atiixp.ko: kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/snd-atiixp-modem.ko: kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/snd-azt3328.ko: kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/sound/drivers/mpu401/snd-mpu401-uart.ko kernel/sound/drivers/opl3/snd-opl3-lib.ko kernel/sound/core/snd-hwdep.ko kernel/drivers/input/gameport/gameport.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/snd-bt87x.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/snd-cmipci.ko: kernel/sound/drivers/mpu401/snd-mpu401-uart.ko kernel/sound/drivers/opl3/snd-opl3-lib.ko kernel/sound/core/snd-hwdep.ko kernel/drivers/input/gameport/gameport.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/snd-cs4281.ko: kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/sound/drivers/opl3/snd-opl3-lib.ko kernel/sound/core/snd-hwdep.ko kernel/drivers/input/gameport/gameport.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/snd-ens1370.ko: kernel/drivers/input/gameport/gameport.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/snd-ens1371.ko: kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/drivers/input/gameport/gameport.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/snd-es1938.ko: kernel/sound/drivers/mpu401/snd-mpu401-uart.ko kernel/sound/drivers/opl3/snd-opl3-lib.ko kernel/sound/core/snd-hwdep.ko kernel/drivers/input/gameport/gameport.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/snd-es1968.ko: kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/sound/drivers/mpu401/snd-mpu401-uart.ko kernel/drivers/media/radio/tea575x.ko kernel/drivers/input/gameport/gameport.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/snd-fm801.ko: kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/sound/drivers/mpu401/snd-mpu401-uart.ko kernel/sound/drivers/opl3/snd-opl3-lib.ko kernel/sound/core/snd-hwdep.ko kernel/drivers/media/radio/tea575x.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/snd-intel8x0.ko: kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/snd-intel8x0m.ko: kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/snd-maestro3.ko: kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/snd-rme32.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/snd-rme96.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/snd-sonicvibes.ko: kernel/sound/drivers/mpu401/snd-mpu401-uart.ko kernel/sound/drivers/opl3/snd-opl3-lib.ko kernel/sound/core/snd-hwdep.ko kernel/drivers/input/gameport/gameport.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/snd-via82xx.ko: kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/sound/drivers/mpu401/snd-mpu401-uart.ko kernel/drivers/input/gameport/gameport.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/snd-via82xx-modem.ko: kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/ac97/snd-ac97-codec.ko: kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/ali5451/snd-ali5451.ko: kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/asihpi/snd-asihpi.ko: kernel/sound/core/snd-hwdep.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/au88x0/snd-au8810.ko: kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/sound/drivers/mpu401/snd-mpu401-uart.ko kernel/drivers/input/gameport/gameport.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/au88x0/snd-au8820.ko: kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/sound/drivers/mpu401/snd-mpu401-uart.ko kernel/drivers/input/gameport/gameport.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/au88x0/snd-au8830.ko: kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/sound/drivers/mpu401/snd-mpu401-uart.ko kernel/drivers/input/gameport/gameport.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/aw2/snd-aw2.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/ctxfi/snd-ctxfi.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/ca0106/snd-ca0106.ko: kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/cs46xx/snd-cs46xx.ko: kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/drivers/input/gameport/gameport.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/lola/snd-lola.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/lx6464es/snd-lx6464es.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/echoaudio/snd-darla20.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/echoaudio/snd-gina20.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/echoaudio/snd-layla20.ko: kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/echoaudio/snd-darla24.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/echoaudio/snd-gina24.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/echoaudio/snd-layla24.ko: kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/echoaudio/snd-mona.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/echoaudio/snd-mia.ko: kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/echoaudio/snd-echo3g.ko: kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/echoaudio/snd-indigo.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/echoaudio/snd-indigoio.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/echoaudio/snd-indigodj.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/echoaudio/snd-indigoiox.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/echoaudio/snd-indigodjx.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/emu10k1/snd-emu10k1.ko: kernel/sound/synth/snd-util-mem.ko kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/sound/core/snd-hwdep.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/emu10k1/snd-emu10k1-synth.ko: kernel/sound/synth/emux/snd-emux-synth.ko kernel/sound/pci/emu10k1/snd-emu10k1.ko kernel/sound/synth/snd-util-mem.ko kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/sound/core/seq/snd-seq-midi-emul.ko kernel/sound/core/snd-hwdep.ko kernel/sound/core/seq/snd-seq-virmidi.ko kernel/sound/core/seq/snd-seq-midi-event.ko kernel/sound/core/seq/snd-seq.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/emu10k1/snd-emu10k1x.ko: kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/hda/snd-hda-codec.ko: kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/hda/snd-hda-codec-generic.ko: kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/hda/snd-hda-codec-realtek.ko: kernel/sound/pci/hda/snd-hda-codec-generic.ko kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/hda/snd-hda-codec-cmedia.ko: kernel/sound/pci/hda/snd-hda-codec-generic.ko kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/hda/snd-hda-codec-analog.ko: kernel/sound/pci/hda/snd-hda-codec-generic.ko kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/hda/snd-hda-codec-idt.ko: kernel/sound/pci/hda/snd-hda-codec-generic.ko kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/hda/snd-hda-codec-si3054.ko: kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/hda/snd-hda-codec-cirrus.ko: kernel/sound/pci/hda/snd-hda-codec-generic.ko kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/hda/snd-hda-codec-cs8409.ko: kernel/sound/pci/hda/snd-hda-codec-generic.ko kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/hda/snd-hda-codec-ca0110.ko: kernel/sound/pci/hda/snd-hda-codec-generic.ko kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/hda/snd-hda-codec-ca0132.ko: kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/hda/snd-hda-codec-conexant.ko: kernel/sound/pci/hda/snd-hda-codec-generic.ko kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/hda/snd-hda-codec-via.ko: kernel/sound/pci/hda/snd-hda-codec-generic.ko kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/hda/snd-hda-codec-hdmi.ko: kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/hda/snd-hda-scodec-cs35l41.ko: kernel/sound/soc/codecs/snd-soc-cs35l41-lib.ko
kernel/sound/pci/hda/snd-hda-scodec-cs35l41-i2c.ko: kernel/sound/pci/hda/snd-hda-scodec-cs35l41.ko kernel/sound/soc/codecs/snd-soc-cs35l41-lib.ko
kernel/sound/pci/hda/snd-hda-scodec-cs35l41-spi.ko: kernel/sound/pci/hda/snd-hda-scodec-cs35l41.ko kernel/sound/soc/codecs/snd-soc-cs35l41-lib.ko
kernel/sound/pci/hda/snd-hda-intel.ko: kernel/sound/hda/snd-intel-dspcfg.ko kernel/sound/hda/snd-intel-sdw-acpi.ko kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/ice1712/snd-ice1712.ko: kernel/sound/i2c/snd-cs8427.ko kernel/sound/i2c/snd-i2c.ko kernel/sound/pci/ice1712/snd-ice17xx-ak4xxx.ko kernel/sound/i2c/other/snd-ak4xxx-adda.ko kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/sound/drivers/mpu401/snd-mpu401-uart.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/ice1712/snd-ice17xx-ak4xxx.ko: kernel/sound/i2c/other/snd-ak4xxx-adda.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/ice1712/snd-ice1724.ko: kernel/sound/i2c/other/snd-ak4113.ko kernel/sound/i2c/other/snd-pt2258.ko kernel/sound/i2c/other/snd-ak4114.ko kernel/sound/i2c/snd-i2c.ko kernel/sound/pci/ice1712/snd-ice17xx-ak4xxx.ko kernel/sound/i2c/other/snd-ak4xxx-adda.ko kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/korg1212/snd-korg1212.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/mixart/snd-mixart.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/nm256/snd-nm256.ko: kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/oxygen/snd-oxygen-lib.ko: kernel/sound/drivers/mpu401/snd-mpu401-uart.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/oxygen/snd-oxygen.ko: kernel/sound/pci/oxygen/snd-oxygen-lib.ko kernel/sound/drivers/mpu401/snd-mpu401-uart.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/oxygen/snd-virtuoso.ko: kernel/sound/pci/oxygen/snd-oxygen-lib.ko kernel/sound/drivers/mpu401/snd-mpu401-uart.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/pcxhr/snd-pcxhr.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/riptide/snd-riptide.ko: kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/sound/drivers/mpu401/snd-mpu401-uart.ko kernel/sound/drivers/opl3/snd-opl3-lib.ko kernel/sound/core/snd-hwdep.ko kernel/drivers/input/gameport/gameport.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/rme9652/snd-rme9652.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/rme9652/snd-hdsp.ko: kernel/sound/core/snd-hwdep.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/rme9652/snd-hdspm.ko: kernel/sound/core/snd-hwdep.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/trident/snd-trident.ko: kernel/sound/synth/snd-util-mem.ko kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/sound/drivers/mpu401/snd-mpu401-uart.ko kernel/drivers/input/gameport/gameport.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/ymfpci/snd-ymfpci.ko: kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/sound/drivers/mpu401/snd-mpu401-uart.ko kernel/sound/drivers/opl3/snd-opl3-lib.ko kernel/sound/core/snd-hwdep.ko kernel/drivers/input/gameport/gameport.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pci/vx222/snd-vx222.ko: kernel/sound/drivers/vx/snd-vx-lib.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/synth/snd-util-mem.ko:
kernel/sound/synth/emux/snd-emux-synth.ko: kernel/sound/synth/snd-util-mem.ko kernel/sound/core/seq/snd-seq-midi-emul.ko kernel/sound/core/snd-hwdep.ko kernel/sound/core/seq/snd-seq-virmidi.ko kernel/sound/core/seq/snd-seq-midi-event.ko kernel/sound/core/seq/snd-seq.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/usb/snd-usb-audio.ko: kernel/sound/usb/snd-usbmidi-lib.ko kernel/sound/core/snd-hwdep.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/drivers/media/mc/mc.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/usb/snd-usbmidi-lib.ko: kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/usb/misc/snd-ua101.ko: kernel/sound/usb/snd-usbmidi-lib.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/usb/usx2y/snd-usb-usx2y.ko: kernel/sound/usb/snd-usbmidi-lib.ko kernel/sound/core/snd-hwdep.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/usb/usx2y/snd-usb-us122l.ko: kernel/sound/usb/snd-usbmidi-lib.ko kernel/sound/core/snd-hwdep.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/usb/caiaq/snd-usb-caiaq.ko: kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/usb/6fire/snd-usb-6fire.ko: kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/usb/hiface/snd-usb-hiface.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/usb/bcd2000/snd-bcd2000.ko: kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/usb/line6/snd-usb-line6.ko: kernel/sound/core/snd-hwdep.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/usb/line6/snd-usb-pod.ko: kernel/sound/usb/line6/snd-usb-line6.ko kernel/sound/core/snd-hwdep.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/usb/line6/snd-usb-podhd.ko: kernel/sound/usb/line6/snd-usb-line6.ko kernel/sound/core/snd-hwdep.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/usb/line6/snd-usb-toneport.ko: kernel/sound/usb/line6/snd-usb-line6.ko kernel/sound/core/snd-hwdep.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/usb/line6/snd-usb-variax.ko: kernel/sound/usb/line6/snd-usb-line6.ko kernel/sound/core/snd-hwdep.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/firewire/snd-firewire-lib.ko: kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/drivers/firewire/firewire-core.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko kernel/lib/crc-itu-t.ko
kernel/sound/firewire/dice/snd-dice.ko: kernel/sound/firewire/snd-firewire-lib.ko kernel/sound/core/snd-hwdep.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/drivers/firewire/firewire-core.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko kernel/lib/crc-itu-t.ko
kernel/sound/firewire/oxfw/snd-oxfw.ko: kernel/sound/firewire/snd-firewire-lib.ko kernel/sound/core/snd-hwdep.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/drivers/firewire/firewire-core.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko kernel/lib/crc-itu-t.ko
kernel/sound/firewire/snd-isight.ko: kernel/sound/firewire/snd-firewire-lib.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/drivers/firewire/firewire-core.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko kernel/lib/crc-itu-t.ko
kernel/sound/firewire/fireworks/snd-fireworks.ko: kernel/sound/firewire/snd-firewire-lib.ko kernel/sound/core/snd-hwdep.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/drivers/firewire/firewire-core.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko kernel/lib/crc-itu-t.ko
kernel/sound/firewire/bebob/snd-bebob.ko: kernel/sound/firewire/snd-firewire-lib.ko kernel/sound/core/snd-hwdep.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/drivers/firewire/firewire-core.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko kernel/lib/crc-itu-t.ko
kernel/sound/firewire/digi00x/snd-firewire-digi00x.ko: kernel/sound/firewire/snd-firewire-lib.ko kernel/sound/core/snd-hwdep.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/drivers/firewire/firewire-core.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko kernel/lib/crc-itu-t.ko
kernel/sound/firewire/tascam/snd-firewire-tascam.ko: kernel/sound/firewire/snd-firewire-lib.ko kernel/sound/core/snd-hwdep.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/drivers/firewire/firewire-core.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko kernel/lib/crc-itu-t.ko
kernel/sound/firewire/motu/snd-firewire-motu.ko: kernel/sound/firewire/snd-firewire-lib.ko kernel/sound/core/snd-hwdep.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/drivers/firewire/firewire-core.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko kernel/lib/crc-itu-t.ko
kernel/sound/firewire/fireface/snd-fireface.ko: kernel/sound/firewire/snd-firewire-lib.ko kernel/sound/core/snd-hwdep.ko kernel/sound/core/snd-rawmidi.ko kernel/sound/core/snd-seq-device.ko kernel/drivers/firewire/firewire-core.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko kernel/lib/crc-itu-t.ko
kernel/sound/pcmcia/vx/snd-vxpocket.ko: kernel/sound/drivers/vx/snd-vx-lib.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/pcmcia/pdaudiocf/snd-pdaudiocf.ko: kernel/sound/i2c/other/snd-ak4117.ko kernel/drivers/pcmcia/pcmcia.ko kernel/drivers/pcmcia/pcmcia_core.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/snd-soc-acpi.ko:
kernel/sound/soc/snd-soc-core.ko: kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-ac97.ko: kernel/sound/pci/ac97/snd-ac97-codec.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-adau-utils.ko:
kernel/sound/soc/codecs/snd-soc-adau1372.ko: kernel/sound/soc/codecs/snd-soc-adau-utils.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-adau1372-i2c.ko: kernel/sound/soc/codecs/snd-soc-adau1372.ko kernel/sound/soc/codecs/snd-soc-adau-utils.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-adau1372-spi.ko: kernel/sound/soc/codecs/snd-soc-adau1372.ko kernel/sound/soc/codecs/snd-soc-adau-utils.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-adau1701.ko: kernel/sound/soc/codecs/snd-soc-sigmadsp-i2c.ko kernel/sound/soc/codecs/snd-soc-sigmadsp.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-adau17x1.ko: kernel/sound/soc/codecs/snd-soc-sigmadsp-regmap.ko kernel/sound/soc/codecs/snd-soc-sigmadsp.ko kernel/sound/soc/codecs/snd-soc-adau-utils.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-adau1761.ko: kernel/sound/soc/codecs/snd-soc-adau17x1.ko kernel/sound/soc/codecs/snd-soc-sigmadsp-regmap.ko kernel/sound/soc/codecs/snd-soc-sigmadsp.ko kernel/sound/soc/codecs/snd-soc-adau-utils.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-adau1761-i2c.ko: kernel/sound/soc/codecs/snd-soc-adau1761.ko kernel/sound/soc/codecs/snd-soc-adau17x1.ko kernel/sound/soc/codecs/snd-soc-sigmadsp-regmap.ko kernel/sound/soc/codecs/snd-soc-sigmadsp.ko kernel/sound/soc/codecs/snd-soc-adau-utils.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-adau1761-spi.ko: kernel/sound/soc/codecs/snd-soc-adau1761.ko kernel/sound/soc/codecs/snd-soc-adau17x1.ko kernel/sound/soc/codecs/snd-soc-sigmadsp-regmap.ko kernel/sound/soc/codecs/snd-soc-sigmadsp.ko kernel/sound/soc/codecs/snd-soc-adau-utils.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-adau7002.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-adau7118.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-adau7118-i2c.ko: kernel/sound/soc/codecs/snd-soc-adau7118.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-adau7118-hw.ko: kernel/sound/soc/codecs/snd-soc-adau7118.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-ak4104.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-ak4118.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-ak4458.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-ak4554.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-ak4613.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-ak4642.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-ak5386.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-ak5558.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-alc5623.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-arizona.ko: kernel/drivers/mfd/arizona.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-bd28623.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-bt-sco.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-cros-ec-codec.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-cs35l32.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-cs35l33.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-cs35l34.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-cs35l35.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-cs35l36.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-cs35l41.ko: kernel/sound/soc/codecs/snd-soc-wm-adsp.ko kernel/sound/soc/codecs/snd-soc-cs35l41-lib.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-cs35l41-lib.ko:
kernel/sound/soc/codecs/snd-soc-cs35l41-spi.ko: kernel/sound/soc/codecs/snd-soc-cs35l41.ko kernel/sound/soc/codecs/snd-soc-wm-adsp.ko kernel/sound/soc/codecs/snd-soc-cs35l41-lib.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-cs35l41-i2c.ko: kernel/sound/soc/codecs/snd-soc-cs35l41.ko kernel/sound/soc/codecs/snd-soc-wm-adsp.ko kernel/sound/soc/codecs/snd-soc-cs35l41-lib.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-cs42l42.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-cs42l51.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-cs42l51-i2c.ko: kernel/sound/soc/codecs/snd-soc-cs42l51.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-cs42l52.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-cs42l56.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-cs42l73.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-cs4234.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-cs4265.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-cs4270.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-cs4271.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-cs4271-i2c.ko: kernel/sound/soc/codecs/snd-soc-cs4271.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-cs4271-spi.ko: kernel/sound/soc/codecs/snd-soc-cs4271.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-cs42xx8.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-cs42xx8-i2c.ko: kernel/sound/soc/codecs/snd-soc-cs42xx8.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-cs43130.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-cs4341.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-cs4349.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-cs53l30.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-cx2072x.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-da7213.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-da7219.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-dmic.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-es7134.ko:
kernel/sound/soc/codecs/snd-soc-es7241.ko:
kernel/sound/soc/codecs/snd-soc-es8316.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-es8328.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-es8328-i2c.ko: kernel/sound/soc/codecs/snd-soc-es8328.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-es8328-spi.ko: kernel/sound/soc/codecs/snd-soc-es8328.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-gtm601.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-hdac-hdmi.ko: kernel/sound/hda/ext/snd-hda-ext-core.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-hdac-hda.ko: kernel/sound/hda/ext/snd-hda-ext-core.ko kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-ics43432.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-inno-rk3036.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-max9759.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-max98088.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-max98090.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-max98357a.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-max9867.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-max98927.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-max98373.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-max98373-i2c.ko: kernel/sound/soc/codecs/snd-soc-max98373.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-max98373-sdw.ko: kernel/drivers/base/regmap/regmap-sdw.ko kernel/sound/soc/codecs/snd-soc-max98373.ko kernel/drivers/soundwire/soundwire-bus.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-max98390.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-max9860.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-msm8916-analog.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-msm8916-digital.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-mt6351.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-mt6358.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-mt6660.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-nau8315.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-nau8540.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-nau8810.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-nau8822.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-nau8824.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-nau8825.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-hdmi-codec.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-pcm1681.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-pcm179x-codec.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-pcm1789-i2c.ko: kernel/sound/soc/codecs/snd-soc-pcm1789-codec.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-pcm1789-codec.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-pcm179x-i2c.ko: kernel/sound/soc/codecs/snd-soc-pcm179x-codec.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-pcm179x-spi.ko: kernel/sound/soc/codecs/snd-soc-pcm179x-codec.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-pcm186x.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-pcm186x-i2c.ko: kernel/sound/soc/codecs/snd-soc-pcm186x.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-pcm186x-spi.ko: kernel/sound/soc/codecs/snd-soc-pcm186x.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-pcm3060.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-pcm3060-i2c.ko: kernel/sound/soc/codecs/snd-soc-pcm3060.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-pcm3060-spi.ko: kernel/sound/soc/codecs/snd-soc-pcm3060.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-pcm3168a.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-pcm3168a-i2c.ko: kernel/sound/soc/codecs/snd-soc-pcm3168a.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-pcm3168a-spi.ko: kernel/sound/soc/codecs/snd-soc-pcm3168a.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-pcm5102a.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-pcm512x.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-pcm512x-i2c.ko: kernel/sound/soc/codecs/snd-soc-pcm512x.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-pcm512x-spi.ko: kernel/sound/soc/codecs/snd-soc-pcm512x.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rk3328.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rl6231.ko:
kernel/sound/soc/codecs/snd-soc-rl6347a.ko:
kernel/sound/soc/codecs/snd-soc-rt1011.ko: kernel/sound/soc/codecs/snd-soc-rl6231.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rt1015.ko: kernel/sound/soc/codecs/snd-soc-rl6231.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rt1015p.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rt1308.ko: kernel/sound/soc/codecs/snd-soc-rl6231.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rt1308-sdw.ko: kernel/drivers/base/regmap/regmap-sdw.ko kernel/drivers/soundwire/soundwire-bus.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rt1316-sdw.ko: kernel/drivers/base/regmap/regmap-sdw.ko kernel/drivers/soundwire/soundwire-bus.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rt286.ko: kernel/sound/soc/codecs/snd-soc-rl6347a.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rt298.ko: kernel/sound/soc/codecs/snd-soc-rl6347a.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rt5514.ko: kernel/sound/soc/codecs/snd-soc-rt5514-spi.ko kernel/sound/soc/codecs/snd-soc-rl6231.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rt5514-spi.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rt5616.ko: kernel/sound/soc/codecs/snd-soc-rl6231.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rt5631.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rt5640.ko: kernel/sound/soc/codecs/snd-soc-rl6231.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rt5645.ko: kernel/sound/soc/codecs/snd-soc-rl6231.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rt5651.ko: kernel/sound/soc/codecs/snd-soc-rl6231.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rt5659.ko: kernel/sound/soc/codecs/snd-soc-rl6231.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rt5660.ko: kernel/sound/soc/codecs/snd-soc-rl6231.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rt5663.ko: kernel/sound/soc/codecs/snd-soc-rl6231.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rt5670.ko: kernel/sound/soc/codecs/snd-soc-rl6231.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rt5677.ko: kernel/sound/soc/codecs/snd-soc-rt5677-spi.ko kernel/sound/soc/codecs/snd-soc-rl6231.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rt5677-spi.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rt5682.ko: kernel/sound/soc/codecs/snd-soc-rl6231.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rt5682-i2c.ko: kernel/sound/soc/codecs/snd-soc-rt5682.ko kernel/sound/soc/codecs/snd-soc-rl6231.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rt5682-sdw.ko: kernel/sound/soc/codecs/snd-soc-rt5682.ko kernel/sound/soc/codecs/snd-soc-rl6231.ko kernel/drivers/base/regmap/regmap-sdw.ko kernel/drivers/soundwire/soundwire-bus.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rt700.ko: kernel/drivers/base/regmap/regmap-sdw.ko kernel/drivers/soundwire/soundwire-bus.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rt711.ko: kernel/drivers/base/regmap/regmap-sdw.ko kernel/drivers/soundwire/soundwire-bus.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rt711-sdca.ko: kernel/drivers/base/regmap/regmap-sdw-mbq.ko kernel/drivers/base/regmap/regmap-sdw.ko kernel/drivers/soundwire/soundwire-bus.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rt715.ko: kernel/drivers/base/regmap/regmap-sdw.ko kernel/drivers/soundwire/soundwire-bus.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-rt715-sdca.ko: kernel/drivers/base/regmap/regmap-sdw-mbq.ko kernel/drivers/base/regmap/regmap-sdw.ko kernel/drivers/soundwire/soundwire-bus.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-sdw-mockup.ko: kernel/drivers/soundwire/soundwire-bus.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-sgtl5000.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-sigmadsp.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-sigmadsp-i2c.ko: kernel/sound/soc/codecs/snd-soc-sigmadsp.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-sigmadsp-regmap.ko: kernel/sound/soc/codecs/snd-soc-sigmadsp.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-si476x.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-spdif-rx.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-spdif-tx.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-ssm2305.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-ssm2518.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-ssm2602.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-ssm2602-spi.ko: kernel/sound/soc/codecs/snd-soc-ssm2602.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-ssm2602-i2c.ko: kernel/sound/soc/codecs/snd-soc-ssm2602.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-ssm4567.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-sta32x.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-sta350.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-sti-sas.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-tas2552.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-tas2562.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-tas2764.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-tas5086.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-tas571x.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-tas5720.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-tas6424.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-tda7419.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-tas2770.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-tfa9879.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-tfa989x.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-tlv320aic23.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-tlv320aic23-i2c.ko: kernel/sound/soc/codecs/snd-soc-tlv320aic23.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-tlv320aic23-spi.ko: kernel/sound/soc/codecs/snd-soc-tlv320aic23.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-tlv320aic31xx.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-tlv320aic32x4.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-tlv320aic32x4-i2c.ko: kernel/sound/soc/codecs/snd-soc-tlv320aic32x4.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-tlv320aic32x4-spi.ko: kernel/sound/soc/codecs/snd-soc-tlv320aic32x4.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-tlv320aic3x.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-tlv320aic3x-i2c.ko: kernel/sound/soc/codecs/snd-soc-tlv320aic3x.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-tlv320aic3x-spi.ko: kernel/sound/soc/codecs/snd-soc-tlv320aic3x.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-tlv320adcx140.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-tscs42xx.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-tscs454.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-ts3a227e.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-uda1334.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wcd-mbhc.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wcd9335.ko: kernel/drivers/slimbus/slimbus.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wcd934x.ko: kernel/sound/soc/codecs/snd-soc-wcd-mbhc.ko kernel/drivers/slimbus/slimbus.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wcd938x.ko: kernel/sound/soc/codecs/snd-soc-wcd938x-sdw.ko kernel/sound/soc/codecs/snd-soc-wcd-mbhc.ko kernel/drivers/base/regmap/regmap-sdw.ko kernel/drivers/soundwire/soundwire-bus.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wcd938x-sdw.ko: kernel/drivers/soundwire/soundwire-bus.ko
kernel/sound/soc/codecs/snd-soc-wm5102.ko: kernel/sound/soc/codecs/snd-soc-arizona.ko kernel/sound/soc/codecs/snd-soc-wm-adsp.ko kernel/drivers/mfd/arizona.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wm8510.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wm8523.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wm8524.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wm8580.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wm8711.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wm8728.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wm8731.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wm8737.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wm8741.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wm8750.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wm8753.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wm8770.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wm8776.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wm8782.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wm8804.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wm8804-i2c.ko: kernel/sound/soc/codecs/snd-soc-wm8804.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wm8804-spi.ko: kernel/sound/soc/codecs/snd-soc-wm8804.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wm8903.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wm8904.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wm8960.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wm8962.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wm8974.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wm8978.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wm8985.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wm-adsp.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-wsa881x.ko: kernel/drivers/base/regmap/regmap-sdw.ko kernel/drivers/soundwire/soundwire-bus.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-zl38060.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-max98504.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-simple-amplifier.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-tpa6130a2.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-lpass-wsa-macro.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-lpass-va-macro.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-lpass-rx-macro.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-lpass-tx-macro.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/codecs/snd-soc-simple-mux.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/generic/snd-soc-simple-card-utils.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/generic/snd-soc-simple-card.ko: kernel/sound/soc/generic/snd-soc-simple-card-utils.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/adi/snd-soc-adi-axi-i2s.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/adi/snd-soc-adi-axi-spdif.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/amd/acp_audio_dma.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/amd/snd-soc-acp-da7219mx98357-mach.ko: kernel/sound/soc/codecs/snd-soc-da7219.ko kernel/sound/soc/amd/acp_audio_dma.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/amd/snd-soc-acp-rt5645-mach.ko: kernel/sound/soc/codecs/snd-soc-rt5645.ko kernel/sound/soc/codecs/snd-soc-rl6231.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/amd/raven/snd-pci-acp3x.ko:
kernel/sound/soc/amd/raven/snd-acp3x-pcm-dma.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/amd/raven/snd-acp3x-i2s.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/amd/snd-soc-acp-rt5682-mach.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/amd/renoir/snd-rn-pci-acp3x.ko:
kernel/sound/soc/amd/renoir/snd-acp3x-pdm-dma.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/amd/renoir/snd-acp3x-rn.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/amd/vangogh/snd-pci-acp5x.ko:
kernel/sound/soc/amd/vangogh/snd-acp5x-i2s.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/amd/vangogh/snd-acp5x-pcm-dma.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/amd/yc/snd-pci-acp6x.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/amd/yc/snd-acp6x-pdm-dma.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/amd/yc/snd-soc-acp6x-mach.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/bcm/snd-soc-63xx.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/dwc/designware_i2s.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/fsl/snd-soc-fsl-audmix.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/fsl/snd-soc-fsl-asrc.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/fsl/snd-soc-fsl-sai.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/fsl/snd-soc-fsl-ssi.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/fsl/snd-soc-fsl-spdif.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/fsl/snd-soc-fsl-esai.ko:
kernel/sound/soc/fsl/snd-soc-fsl-micfil.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/fsl/snd-soc-fsl-mqs.ko:
kernel/sound/soc/fsl/snd-soc-fsl-easrc.ko:
kernel/sound/soc/fsl/snd-soc-fsl-xcvr.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/fsl/snd-soc-fsl-rpmsg.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/fsl/snd-soc-imx-audmux.ko:
kernel/sound/soc/hisilicon/hi6210-i2s.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/img/img-i2s-in.ko:
kernel/sound/soc/img/img-i2s-out.ko:
kernel/sound/soc/img/img-parallel-out.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/img/img-spdif-in.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/img/img-spdif-out.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/img/pistachio-internal-dac.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/common/snd-soc-sst-dsp.ko:
kernel/sound/soc/intel/common/snd-soc-sst-ipc.ko:
kernel/sound/soc/intel/common/snd-soc-acpi-intel-match.ko: kernel/sound/soc/snd-soc-acpi.ko
kernel/sound/soc/intel/atom/snd-soc-sst-atom-hifi2-platform.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/atom/sst/snd-intel-sst-core.ko: kernel/sound/soc/intel/atom/snd-soc-sst-atom-hifi2-platform.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/atom/sst/snd-intel-sst-pci.ko: kernel/sound/soc/intel/atom/sst/snd-intel-sst-core.ko kernel/sound/soc/intel/atom/snd-soc-sst-atom-hifi2-platform.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/atom/sst/snd-intel-sst-acpi.ko: kernel/sound/soc/intel/common/snd-soc-acpi-intel-match.ko kernel/sound/soc/snd-soc-acpi.ko kernel/sound/soc/intel/atom/sst/snd-intel-sst-core.ko kernel/sound/soc/intel/atom/snd-soc-sst-atom-hifi2-platform.ko kernel/sound/hda/snd-intel-dspcfg.ko kernel/sound/hda/snd-intel-sdw-acpi.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/catpt/snd-soc-catpt.ko: kernel/sound/soc/intel/common/snd-soc-acpi-intel-match.ko kernel/sound/soc/snd-soc-acpi.ko kernel/sound/hda/snd-intel-dspcfg.ko kernel/sound/hda/snd-intel-sdw-acpi.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko kernel/drivers/dma/dw/dw_dmac_core.ko
kernel/sound/soc/intel/skylake/snd-soc-skl.ko: kernel/sound/soc/codecs/snd-soc-hdac-hda.ko kernel/sound/hda/ext/snd-hda-ext-core.ko kernel/sound/soc/intel/common/snd-soc-sst-ipc.ko kernel/sound/soc/intel/common/snd-soc-sst-dsp.ko kernel/sound/soc/intel/common/snd-soc-acpi-intel-match.ko kernel/sound/soc/snd-soc-acpi.ko kernel/sound/hda/snd-intel-dspcfg.ko kernel/sound/hda/snd-intel-sdw-acpi.ko kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/skylake/snd-soc-skl-ssp-clk.ko: kernel/sound/soc/intel/skylake/snd-soc-skl.ko kernel/sound/soc/codecs/snd-soc-hdac-hda.ko kernel/sound/hda/ext/snd-hda-ext-core.ko kernel/sound/soc/intel/common/snd-soc-sst-ipc.ko kernel/sound/soc/intel/common/snd-soc-sst-dsp.ko kernel/sound/soc/intel/common/snd-soc-acpi-intel-match.ko kernel/sound/soc/snd-soc-acpi.ko kernel/sound/hda/snd-intel-dspcfg.ko kernel/sound/hda/snd-intel-sdw-acpi.ko kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-sof_rt5682.ko: kernel/sound/soc/intel/boards/snd-soc-intel-hda-dsp-common.ko kernel/sound/soc/sof/snd-sof.ko kernel/sound/soc/intel/boards/snd-soc-intel-sof-maxim-common.ko kernel/sound/soc/codecs/snd-soc-hdac-hdmi.ko kernel/sound/hda/ext/snd-hda-ext-core.ko kernel/sound/soc/codecs/snd-soc-rt5682.ko kernel/sound/soc/codecs/snd-soc-rl6231.ko kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-sof_cs42l42.ko: kernel/sound/soc/intel/boards/snd-soc-intel-hda-dsp-common.ko kernel/sound/soc/sof/snd-sof.ko kernel/sound/soc/intel/boards/snd-soc-intel-sof-maxim-common.ko kernel/sound/soc/codecs/snd-soc-hdac-hdmi.ko kernel/sound/hda/ext/snd-hda-ext-core.ko kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-sst-haswell.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-sst-bxt-da7219_max98357a.ko: kernel/sound/soc/intel/boards/snd-soc-intel-hda-dsp-common.ko kernel/sound/soc/codecs/snd-soc-hdac-hdmi.ko kernel/sound/hda/ext/snd-hda-ext-core.ko kernel/sound/soc/codecs/snd-soc-da7219.ko kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-sst-bxt-rt298.ko: kernel/sound/soc/codecs/snd-soc-rt298.ko kernel/sound/soc/codecs/snd-soc-rl6347a.ko kernel/sound/soc/intel/boards/snd-soc-intel-hda-dsp-common.ko kernel/sound/soc/codecs/snd-soc-hdac-hdmi.ko kernel/sound/hda/ext/snd-hda-ext-core.ko kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-sst-sof-pcm512x.ko: kernel/sound/soc/intel/boards/snd-soc-intel-hda-dsp-common.ko kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-sst-sof-wm8804.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-sst-glk-rt5682_max98357a.ko: kernel/sound/soc/intel/boards/snd-soc-intel-hda-dsp-common.ko kernel/sound/soc/codecs/snd-soc-hdac-hdmi.ko kernel/sound/hda/ext/snd-hda-ext-core.ko kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-sst-broadwell.ko: kernel/sound/soc/codecs/snd-soc-rt286.ko kernel/sound/soc/codecs/snd-soc-rl6347a.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-sst-bdw-rt5650-mach.ko: kernel/sound/soc/codecs/snd-soc-rt5645.ko kernel/sound/soc/codecs/snd-soc-rl6231.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-sst-bdw-rt5677-mach.ko: kernel/sound/soc/codecs/snd-soc-rt5677.ko kernel/sound/soc/codecs/snd-soc-rt5677-spi.ko kernel/sound/soc/codecs/snd-soc-rl6231.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-sst-bytcr-rt5640.ko: kernel/sound/soc/codecs/snd-soc-rt5640.ko kernel/sound/soc/snd-soc-acpi.ko kernel/sound/soc/codecs/snd-soc-rl6231.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-sst-bytcr-rt5651.ko: kernel/sound/soc/snd-soc-acpi.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-sst-bytcr-wm5102.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-sst-cht-bsw-rt5672.ko: kernel/sound/soc/codecs/snd-soc-rt5670.ko kernel/sound/soc/codecs/snd-soc-rl6231.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-sst-cht-bsw-rt5645.ko: kernel/sound/soc/snd-soc-acpi.ko kernel/sound/soc/codecs/snd-soc-rt5645.ko kernel/sound/soc/codecs/snd-soc-rl6231.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-sst-cht-bsw-max98090_ti.ko: kernel/sound/soc/codecs/snd-soc-ts3a227e.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-sst-cht-bsw-nau8824.ko: kernel/sound/soc/codecs/snd-soc-nau8824.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-sst-byt-cht-cx2072x.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-sst-byt-cht-da7213.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-sst-byt-cht-es8316.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-cml_rt1011_rt5682.ko: kernel/sound/soc/intel/boards/snd-soc-intel-hda-dsp-common.ko kernel/sound/soc/codecs/snd-soc-hdac-hdmi.ko kernel/sound/hda/ext/snd-hda-ext-core.ko kernel/sound/soc/codecs/snd-soc-rt5682.ko kernel/sound/soc/codecs/snd-soc-rl6231.ko kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-kbl_da7219_max98357a.ko: kernel/sound/soc/codecs/snd-soc-hdac-hdmi.ko kernel/sound/hda/ext/snd-hda-ext-core.ko kernel/sound/soc/codecs/snd-soc-da7219.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-kbl_da7219_max98927.ko: kernel/sound/soc/codecs/snd-soc-hdac-hdmi.ko kernel/sound/hda/ext/snd-hda-ext-core.ko kernel/sound/soc/codecs/snd-soc-da7219.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-kbl_rt5663_max98927.ko: kernel/sound/soc/codecs/snd-soc-rt5663.ko kernel/sound/soc/codecs/snd-soc-hdac-hdmi.ko kernel/sound/hda/ext/snd-hda-ext-core.ko kernel/sound/soc/codecs/snd-soc-rl6231.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-kbl_rt5663_rt5514_max98927.ko: kernel/sound/soc/codecs/snd-soc-rt5663.ko kernel/sound/soc/codecs/snd-soc-hdac-hdmi.ko kernel/sound/hda/ext/snd-hda-ext-core.ko kernel/sound/soc/codecs/snd-soc-rl6231.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-kbl_rt5660.ko: kernel/sound/soc/codecs/snd-soc-hdac-hdmi.ko kernel/sound/hda/ext/snd-hda-ext-core.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-skl_rt286.ko: kernel/sound/soc/codecs/snd-soc-rt286.ko kernel/sound/soc/codecs/snd-soc-rl6347a.ko kernel/sound/soc/codecs/snd-soc-hdac-hdmi.ko kernel/sound/hda/ext/snd-hda-ext-core.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-skl_nau88l25_max98357a.ko: kernel/sound/soc/codecs/snd-soc-nau8825.ko kernel/sound/soc/codecs/snd-soc-hdac-hdmi.ko kernel/sound/hda/ext/snd-hda-ext-core.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-skl_nau88l25_ssm4567.ko: kernel/sound/soc/codecs/snd-soc-nau8825.ko kernel/sound/soc/codecs/snd-soc-hdac-hdmi.ko kernel/sound/hda/ext/snd-hda-ext-core.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-skl_hda_dsp.ko: kernel/sound/soc/intel/boards/snd-soc-intel-hda-dsp-common.ko kernel/sound/soc/codecs/snd-soc-hdac-hdmi.ko kernel/sound/hda/ext/snd-hda-ext-core.ko kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-sof_da7219_max98373.ko: kernel/sound/soc/intel/boards/snd-soc-intel-hda-dsp-common.ko kernel/sound/soc/codecs/snd-soc-da7219.ko kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-ehl-rt5660.ko: kernel/sound/soc/intel/boards/snd-soc-intel-hda-dsp-common.ko kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-sof-sdw.ko: kernel/sound/soc/intel/boards/snd-soc-intel-hda-dsp-common.ko kernel/sound/soc/intel/boards/snd-soc-intel-sof-maxim-common.ko kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/drivers/soundwire/soundwire-bus.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-intel-hda-dsp-common.ko: kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/intel/boards/snd-soc-intel-sof-maxim-common.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/sof/intel/snd-sof-intel-atom.ko: kernel/sound/soc/sof/snd-sof.ko kernel/sound/soc/snd-soc-acpi.ko kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/sof/intel/snd-sof-acpi-intel-byt.ko: kernel/sound/soc/sof/intel/snd-sof-intel-ipc.ko kernel/sound/soc/sof/snd-sof-acpi.ko kernel/sound/soc/sof/intel/snd-sof-intel-atom.ko kernel/sound/soc/sof/xtensa/snd-sof-xtensa-dsp.ko kernel/sound/soc/sof/snd-sof.ko kernel/sound/soc/intel/common/snd-soc-acpi-intel-match.ko kernel/sound/soc/snd-soc-acpi.ko kernel/sound/hda/snd-intel-dspcfg.ko kernel/sound/hda/snd-intel-sdw-acpi.ko kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/sof/intel/snd-sof-acpi-intel-bdw.ko: kernel/sound/soc/sof/intel/snd-sof-intel-ipc.ko kernel/sound/soc/sof/snd-sof-acpi.ko kernel/sound/soc/sof/xtensa/snd-sof-xtensa-dsp.ko kernel/sound/soc/sof/snd-sof.ko kernel/sound/soc/intel/common/snd-soc-acpi-intel-match.ko kernel/sound/soc/snd-soc-acpi.ko kernel/sound/hda/snd-intel-dspcfg.ko kernel/sound/hda/snd-intel-sdw-acpi.ko kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/sof/intel/snd-sof-intel-ipc.ko: kernel/sound/soc/sof/snd-sof.ko kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/sof/intel/snd-sof-intel-hda-common.ko: kernel/drivers/soundwire/soundwire-intel.ko kernel/drivers/soundwire/soundwire-generic-allocation.ko kernel/drivers/soundwire/soundwire-cadence.ko kernel/sound/soc/sof/intel/snd-sof-intel-hda.ko kernel/sound/soc/sof/snd-sof-pci.ko kernel/sound/soc/sof/xtensa/snd-sof-xtensa-dsp.ko kernel/sound/soc/sof/snd-sof.ko kernel/sound/soc/codecs/snd-soc-hdac-hda.ko kernel/sound/hda/ext/snd-hda-ext-core.ko kernel/sound/soc/intel/common/snd-soc-acpi-intel-match.ko kernel/sound/soc/snd-soc-acpi.ko kernel/sound/hda/snd-intel-dspcfg.ko kernel/sound/hda/snd-intel-sdw-acpi.ko kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/drivers/soundwire/soundwire-bus.ko kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/sof/intel/snd-sof-intel-hda.ko: kernel/sound/hda/ext/snd-hda-ext-core.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/sof/intel/snd-sof-pci-intel-tng.ko: kernel/sound/soc/sof/snd-sof-pci.ko kernel/sound/soc/sof/intel/snd-sof-intel-ipc.ko kernel/sound/soc/sof/intel/snd-sof-intel-atom.ko kernel/sound/soc/sof/xtensa/snd-sof-xtensa-dsp.ko kernel/sound/soc/sof/snd-sof.ko kernel/sound/soc/snd-soc-acpi.ko kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/sof/intel/snd-sof-pci-intel-apl.ko: kernel/sound/soc/sof/intel/snd-sof-intel-hda-common.ko kernel/drivers/soundwire/soundwire-intel.ko kernel/drivers/soundwire/soundwire-generic-allocation.ko kernel/drivers/soundwire/soundwire-cadence.ko kernel/sound/soc/sof/intel/snd-sof-intel-hda.ko kernel/sound/soc/sof/snd-sof-pci.ko kernel/sound/soc/sof/xtensa/snd-sof-xtensa-dsp.ko kernel/sound/soc/sof/snd-sof.ko kernel/sound/soc/codecs/snd-soc-hdac-hda.ko kernel/sound/hda/ext/snd-hda-ext-core.ko kernel/sound/soc/intel/common/snd-soc-acpi-intel-match.ko kernel/sound/soc/snd-soc-acpi.ko kernel/sound/hda/snd-intel-dspcfg.ko kernel/sound/hda/snd-intel-sdw-acpi.ko kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/drivers/soundwire/soundwire-bus.ko kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/sof/intel/snd-sof-pci-intel-cnl.ko: kernel/sound/soc/sof/intel/snd-sof-intel-hda-common.ko kernel/drivers/soundwire/soundwire-intel.ko kernel/drivers/soundwire/soundwire-generic-allocation.ko kernel/drivers/soundwire/soundwire-cadence.ko kernel/sound/soc/sof/intel/snd-sof-intel-hda.ko kernel/sound/soc/sof/snd-sof-pci.ko kernel/sound/soc/sof/xtensa/snd-sof-xtensa-dsp.ko kernel/sound/soc/sof/snd-sof.ko kernel/sound/soc/codecs/snd-soc-hdac-hda.ko kernel/sound/hda/ext/snd-hda-ext-core.ko kernel/sound/soc/intel/common/snd-soc-acpi-intel-match.ko kernel/sound/soc/snd-soc-acpi.ko kernel/sound/hda/snd-intel-dspcfg.ko kernel/sound/hda/snd-intel-sdw-acpi.ko kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/drivers/soundwire/soundwire-bus.ko kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/sof/intel/snd-sof-pci-intel-icl.ko: kernel/sound/soc/sof/intel/snd-sof-intel-hda-common.ko kernel/drivers/soundwire/soundwire-intel.ko kernel/drivers/soundwire/soundwire-generic-allocation.ko kernel/drivers/soundwire/soundwire-cadence.ko kernel/sound/soc/sof/intel/snd-sof-intel-hda.ko kernel/sound/soc/sof/snd-sof-pci.ko kernel/sound/soc/sof/xtensa/snd-sof-xtensa-dsp.ko kernel/sound/soc/sof/snd-sof.ko kernel/sound/soc/codecs/snd-soc-hdac-hda.ko kernel/sound/hda/ext/snd-hda-ext-core.ko kernel/sound/soc/intel/common/snd-soc-acpi-intel-match.ko kernel/sound/soc/snd-soc-acpi.ko kernel/sound/hda/snd-intel-dspcfg.ko kernel/sound/hda/snd-intel-sdw-acpi.ko kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/drivers/soundwire/soundwire-bus.ko kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/sof/intel/snd-sof-pci-intel-tgl.ko: kernel/sound/soc/sof/intel/snd-sof-intel-hda-common.ko kernel/drivers/soundwire/soundwire-intel.ko kernel/drivers/soundwire/soundwire-generic-allocation.ko kernel/drivers/soundwire/soundwire-cadence.ko kernel/sound/soc/sof/intel/snd-sof-intel-hda.ko kernel/sound/soc/sof/snd-sof-pci.ko kernel/sound/soc/sof/xtensa/snd-sof-xtensa-dsp.ko kernel/sound/soc/sof/snd-sof.ko kernel/sound/soc/codecs/snd-soc-hdac-hda.ko kernel/sound/hda/ext/snd-hda-ext-core.ko kernel/sound/soc/intel/common/snd-soc-acpi-intel-match.ko kernel/sound/soc/snd-soc-acpi.ko kernel/sound/hda/snd-intel-dspcfg.ko kernel/sound/hda/snd-intel-sdw-acpi.ko kernel/sound/pci/hda/snd-hda-codec.ko kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-hwdep.ko kernel/drivers/soundwire/soundwire-bus.ko kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/sof/snd-sof.ko: kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/sof/snd-sof-acpi.ko: kernel/sound/soc/sof/snd-sof.ko kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/sof/snd-sof-pci.ko: kernel/sound/soc/sof/snd-sof.ko kernel/drivers/leds/trigger/ledtrig-audio.ko kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/sof/xtensa/snd-sof-xtensa-dsp.ko:
kernel/sound/soc/xilinx/snd-soc-xlnx-i2s.ko:
kernel/sound/soc/xilinx/snd-soc-xlnx-formatter-pcm.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/soc/xilinx/snd-soc-xlnx-spdif.ko:
kernel/sound/soc/xtensa/snd-soc-xtfpga-i2s.ko: kernel/sound/soc/snd-soc-core.ko kernel/sound/core/snd-compress.ko kernel/sound/ac97_bus.ko kernel/sound/core/snd-pcm-dmaengine.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/hda/snd-hda-core.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/hda/ext/snd-hda-ext-core.ko: kernel/sound/hda/snd-hda-core.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/hda/snd-intel-dspcfg.ko: kernel/sound/hda/snd-intel-sdw-acpi.ko
kernel/sound/hda/snd-intel-sdw-acpi.ko:
kernel/sound/x86/snd-hdmi-lpe-audio.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/xen/snd_xen_front.ko: kernel/drivers/xen/xen-front-pgdir-shbuf.ko kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/virtio/virtio_snd.ko: kernel/sound/core/snd-pcm.ko kernel/sound/core/snd-timer.ko kernel/sound/core/snd.ko kernel/sound/soundcore.ko
kernel/sound/ac97_bus.ko:
kernel/ubuntu/ubuntu-host/ubuntu-host.ko:
kernel/samples/trace_printk/trace-printk.ko:
kernel/samples/ftrace/ftrace-direct.ko:
kernel/samples/ftrace/ftrace-direct-too.ko:
kernel/samples/ftrace/ftrace-direct-modify.ko:
kernel/samples/ftrace/sample-trace-array.ko:
kernel/net/core/pktgen.ko:
kernel/net/core/failover.ko:
kernel/net/802/p8022.ko: kernel/net/llc/llc.ko
kernel/net/802/psnap.ko: kernel/net/llc/llc.ko
kernel/net/802/stp.ko: kernel/net/llc/llc.ko
kernel/net/802/garp.ko: kernel/net/802/stp.ko kernel/net/llc/llc.ko
kernel/net/802/mrp.ko:
kernel/net/sched/act_police.ko:
kernel/net/sched/act_gact.ko:
kernel/net/sched/act_mirred.ko:
kernel/net/sched/act_sample.ko: kernel/net/psample/psample.ko
kernel/net/sched/act_ipt.ko: kernel/net/netfilter/x_tables.ko
kernel/net/sched/act_nat.ko:
kernel/net/sched/act_pedit.ko:
kernel/net/sched/act_simple.ko:
kernel/net/sched/act_skbedit.ko:
kernel/net/sched/act_csum.ko: kernel/lib/libcrc32c.ko
kernel/net/sched/act_mpls.ko:
kernel/net/sched/act_vlan.ko:
kernel/net/sched/act_bpf.ko:
kernel/net/sched/act_connmark.ko: kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/sched/act_ctinfo.ko: kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/sched/act_skbmod.ko:
kernel/net/sched/act_tunnel_key.ko:
kernel/net/sched/act_ct.ko: kernel/net/netfilter/nf_flow_table.ko kernel/net/netfilter/nf_nat.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/sched/act_gate.ko:
kernel/net/sched/sch_cbq.ko:
kernel/net/sched/sch_htb.ko:
kernel/net/sched/sch_hfsc.ko:
kernel/net/sched/sch_red.ko:
kernel/net/sched/sch_gred.ko:
kernel/net/sched/sch_ingress.ko:
kernel/net/sched/sch_dsmark.ko:
kernel/net/sched/sch_sfb.ko:
kernel/net/sched/sch_sfq.ko:
kernel/net/sched/sch_tbf.ko:
kernel/net/sched/sch_teql.ko:
kernel/net/sched/sch_prio.ko:
kernel/net/sched/sch_multiq.ko:
kernel/net/sched/sch_atm.ko:
kernel/net/sched/sch_netem.ko:
kernel/net/sched/sch_drr.ko:
kernel/net/sched/sch_plug.ko:
kernel/net/sched/sch_ets.ko:
kernel/net/sched/sch_mqprio.ko:
kernel/net/sched/sch_skbprio.ko:
kernel/net/sched/sch_choke.ko:
kernel/net/sched/sch_qfq.ko:
kernel/net/sched/sch_codel.ko:
kernel/net/sched/sch_fq_codel.ko:
kernel/net/sched/sch_cake.ko:
kernel/net/sched/sch_fq.ko:
kernel/net/sched/sch_hhf.ko:
kernel/net/sched/sch_pie.ko:
kernel/net/sched/sch_fq_pie.ko: kernel/net/sched/sch_pie.ko
kernel/net/sched/sch_cbs.ko:
kernel/net/sched/sch_etf.ko:
kernel/net/sched/sch_taprio.ko:
kernel/net/sched/cls_u32.ko:
kernel/net/sched/cls_route.ko:
kernel/net/sched/cls_fw.ko:
kernel/net/sched/cls_rsvp.ko:
kernel/net/sched/cls_tcindex.ko:
kernel/net/sched/cls_rsvp6.ko:
kernel/net/sched/cls_basic.ko:
kernel/net/sched/cls_flow.ko:
kernel/net/sched/cls_cgroup.ko:
kernel/net/sched/cls_bpf.ko:
kernel/net/sched/cls_flower.ko:
kernel/net/sched/cls_matchall.ko:
kernel/net/sched/em_cmp.ko:
kernel/net/sched/em_nbyte.ko:
kernel/net/sched/em_u32.ko:
kernel/net/sched/em_meta.ko:
kernel/net/sched/em_text.ko:
kernel/net/sched/em_canid.ko:
kernel/net/sched/em_ipset.ko: kernel/net/netfilter/ipset/ip_set.ko kernel/net/netfilter/nfnetlink.ko
kernel/net/sched/em_ipt.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netlink/netlink_diag.ko:
kernel/net/netfilter/nfnetlink.ko:
kernel/net/netfilter/nfnetlink_acct.ko: kernel/net/netfilter/nfnetlink.ko
kernel/net/netfilter/nfnetlink_queue.ko: kernel/net/netfilter/nfnetlink.ko
kernel/net/netfilter/nfnetlink_log.ko: kernel/net/netfilter/nfnetlink.ko
kernel/net/netfilter/nfnetlink_osf.ko: kernel/net/netfilter/nfnetlink.ko
kernel/net/netfilter/nfnetlink_hook.ko: kernel/net/netfilter/nfnetlink.ko
kernel/net/netfilter/nf_conntrack.ko: kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nf_conntrack_netlink.ko: kernel/net/netfilter/nfnetlink.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nfnetlink_cttimeout.ko: kernel/net/netfilter/nfnetlink.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nfnetlink_cthelper.ko: kernel/net/netfilter/nfnetlink.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nf_conntrack_amanda.ko: kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nf_conntrack_ftp.ko: kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nf_conntrack_h323.ko: kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nf_conntrack_irc.ko: kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nf_conntrack_broadcast.ko: kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nf_conntrack_netbios_ns.ko: kernel/net/netfilter/nf_conntrack_broadcast.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nf_conntrack_snmp.ko: kernel/net/netfilter/nf_conntrack_broadcast.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nf_conntrack_pptp.ko: kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nf_conntrack_sane.ko: kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nf_conntrack_sip.ko: kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nf_conntrack_tftp.ko: kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nf_log_syslog.ko:
kernel/net/netfilter/nf_nat.ko: kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nf_nat_amanda.ko: kernel/net/netfilter/nf_conntrack_amanda.ko kernel/net/netfilter/nf_nat.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nf_nat_ftp.ko: kernel/net/netfilter/nf_conntrack_ftp.ko kernel/net/netfilter/nf_nat.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nf_nat_irc.ko: kernel/net/netfilter/nf_conntrack_irc.ko kernel/net/netfilter/nf_nat.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nf_nat_sip.ko: kernel/net/netfilter/nf_conntrack_sip.ko kernel/net/netfilter/nf_nat.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nf_nat_tftp.ko: kernel/net/netfilter/nf_conntrack_tftp.ko kernel/net/netfilter/nf_nat.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nf_synproxy_core.ko: kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nf_conncount.ko: kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nf_dup_netdev.ko:
kernel/net/netfilter/nf_tables.ko: kernel/net/netfilter/nfnetlink.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_compat.ko: kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/net/netfilter/x_tables.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_connlimit.ko: kernel/net/netfilter/nf_conncount.ko kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_numgen.ko: kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_ct.ko: kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_flow_offload.ko: kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/net/netfilter/nf_flow_table.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_limit.ko: kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_nat.ko: kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/net/netfilter/nf_nat.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_objref.ko: kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_queue.ko: kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_quota.ko: kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_reject.ko: kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_reject_inet.ko: kernel/net/ipv4/netfilter/nf_reject_ipv4.ko kernel/net/ipv6/netfilter/nf_reject_ipv6.ko kernel/net/netfilter/nft_reject.ko kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_reject_netdev.ko: kernel/net/ipv4/netfilter/nf_reject_ipv4.ko kernel/net/ipv6/netfilter/nf_reject_ipv6.ko kernel/net/netfilter/nft_reject.ko kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_tunnel.ko: kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_counter.ko: kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_log.ko: kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_masq.ko: kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/net/netfilter/nf_nat.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_redir.ko: kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/net/netfilter/nf_nat.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_hash.ko: kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_fib.ko: kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_fib_inet.ko: kernel/net/ipv4/netfilter/nft_fib_ipv4.ko kernel/net/ipv6/netfilter/nft_fib_ipv6.ko kernel/net/netfilter/nft_fib.ko kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_fib_netdev.ko: kernel/net/ipv4/netfilter/nft_fib_ipv4.ko kernel/net/ipv6/netfilter/nft_fib_ipv6.ko kernel/net/netfilter/nft_fib.ko kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_socket.ko: kernel/net/ipv4/netfilter/nf_socket_ipv4.ko kernel/net/ipv6/netfilter/nf_socket_ipv6.ko kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_osf.ko: kernel/net/netfilter/nfnetlink_osf.ko kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_tproxy.ko: kernel/net/ipv6/netfilter/nf_tproxy_ipv6.ko kernel/net/ipv4/netfilter/nf_tproxy_ipv4.ko kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_xfrm.ko: kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_synproxy.ko: kernel/net/netfilter/nf_synproxy_core.ko kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_chain_nat.ko: kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/net/netfilter/nf_nat.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_dup_netdev.ko: kernel/net/netfilter/nf_dup_netdev.ko kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nft_fwd_netdev.ko: kernel/net/netfilter/nf_dup_netdev.ko kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nf_flow_table.ko: kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/nf_flow_table_inet.ko: kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/net/netfilter/nf_flow_table.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/x_tables.ko:
kernel/net/netfilter/xt_tcpudp.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_mark.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_connmark.ko: kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/net/netfilter/x_tables.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/xt_set.ko: kernel/net/netfilter/ipset/ip_set.ko kernel/net/netfilter/nfnetlink.ko kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_nat.ko: kernel/net/netfilter/nf_nat.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/net/netfilter/x_tables.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/xt_AUDIT.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_CHECKSUM.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_CLASSIFY.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_CONNSECMARK.ko: kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/net/netfilter/x_tables.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/xt_CT.ko: kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/net/netfilter/x_tables.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/xt_DSCP.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_HL.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_HMARK.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_LED.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_LOG.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_NETMAP.ko: kernel/net/netfilter/nf_nat.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/net/netfilter/x_tables.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/xt_NFLOG.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_NFQUEUE.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_RATEEST.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_REDIRECT.ko: kernel/net/netfilter/nf_nat.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/net/netfilter/x_tables.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/xt_MASQUERADE.ko: kernel/net/netfilter/nf_nat.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/net/netfilter/x_tables.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/xt_SECMARK.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_TPROXY.ko: kernel/net/ipv6/netfilter/nf_tproxy_ipv6.ko kernel/net/ipv4/netfilter/nf_tproxy_ipv4.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_TCPMSS.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_TCPOPTSTRIP.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_TEE.ko: kernel/net/ipv6/netfilter/nf_dup_ipv6.ko kernel/net/ipv4/netfilter/nf_dup_ipv4.ko kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_TRACE.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_IDLETIMER.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_addrtype.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_bpf.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_cluster.ko: kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/net/netfilter/x_tables.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/xt_comment.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_connbytes.ko: kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/net/netfilter/x_tables.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/xt_connlabel.ko: kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/net/netfilter/x_tables.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/xt_connlimit.ko: kernel/net/netfilter/nf_conncount.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/net/netfilter/x_tables.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/xt_conntrack.ko: kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/net/netfilter/x_tables.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/xt_cpu.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_dccp.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_devgroup.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_dscp.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_ecn.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_esp.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_hashlimit.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_helper.ko: kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/net/netfilter/x_tables.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/xt_hl.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_ipcomp.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_iprange.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_ipvs.ko: kernel/net/netfilter/ipvs/ip_vs.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/net/netfilter/x_tables.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/xt_l2tp.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_length.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_limit.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_mac.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_multiport.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_nfacct.ko: kernel/net/netfilter/nfnetlink_acct.ko kernel/net/netfilter/nfnetlink.ko kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_osf.ko: kernel/net/netfilter/nfnetlink_osf.ko kernel/net/netfilter/nfnetlink.ko kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_owner.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_cgroup.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_physdev.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_pkttype.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_policy.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_quota.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_rateest.ko: kernel/net/netfilter/xt_RATEEST.ko kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_realm.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_recent.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_sctp.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_socket.ko: kernel/net/ipv4/netfilter/nf_socket_ipv4.ko kernel/net/ipv6/netfilter/nf_socket_ipv6.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_state.ko: kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/net/netfilter/x_tables.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/xt_statistic.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_string.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_tcpmss.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_time.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/xt_u32.ko: kernel/net/netfilter/x_tables.ko
kernel/net/netfilter/ipset/ip_set.ko: kernel/net/netfilter/nfnetlink.ko
kernel/net/netfilter/ipset/ip_set_bitmap_ip.ko: kernel/net/netfilter/ipset/ip_set.ko kernel/net/netfilter/nfnetlink.ko
kernel/net/netfilter/ipset/ip_set_bitmap_ipmac.ko: kernel/net/netfilter/ipset/ip_set.ko kernel/net/netfilter/nfnetlink.ko
kernel/net/netfilter/ipset/ip_set_bitmap_port.ko: kernel/net/netfilter/ipset/ip_set.ko kernel/net/netfilter/nfnetlink.ko
kernel/net/netfilter/ipset/ip_set_hash_ip.ko: kernel/net/netfilter/ipset/ip_set.ko kernel/net/netfilter/nfnetlink.ko
kernel/net/netfilter/ipset/ip_set_hash_ipmac.ko: kernel/net/netfilter/ipset/ip_set.ko kernel/net/netfilter/nfnetlink.ko
kernel/net/netfilter/ipset/ip_set_hash_ipmark.ko: kernel/net/netfilter/ipset/ip_set.ko kernel/net/netfilter/nfnetlink.ko
kernel/net/netfilter/ipset/ip_set_hash_ipport.ko: kernel/net/netfilter/ipset/ip_set.ko kernel/net/netfilter/nfnetlink.ko
kernel/net/netfilter/ipset/ip_set_hash_ipportip.ko: kernel/net/netfilter/ipset/ip_set.ko kernel/net/netfilter/nfnetlink.ko
kernel/net/netfilter/ipset/ip_set_hash_ipportnet.ko: kernel/net/netfilter/ipset/ip_set.ko kernel/net/netfilter/nfnetlink.ko
kernel/net/netfilter/ipset/ip_set_hash_mac.ko: kernel/net/netfilter/ipset/ip_set.ko kernel/net/netfilter/nfnetlink.ko
kernel/net/netfilter/ipset/ip_set_hash_net.ko: kernel/net/netfilter/ipset/ip_set.ko kernel/net/netfilter/nfnetlink.ko
kernel/net/netfilter/ipset/ip_set_hash_netport.ko: kernel/net/netfilter/ipset/ip_set.ko kernel/net/netfilter/nfnetlink.ko
kernel/net/netfilter/ipset/ip_set_hash_netiface.ko: kernel/net/netfilter/ipset/ip_set.ko kernel/net/netfilter/nfnetlink.ko
kernel/net/netfilter/ipset/ip_set_hash_netnet.ko: kernel/net/netfilter/ipset/ip_set.ko kernel/net/netfilter/nfnetlink.ko
kernel/net/netfilter/ipset/ip_set_hash_netportnet.ko: kernel/net/netfilter/ipset/ip_set.ko kernel/net/netfilter/nfnetlink.ko
kernel/net/netfilter/ipset/ip_set_list_set.ko: kernel/net/netfilter/ipset/ip_set.ko kernel/net/netfilter/nfnetlink.ko
kernel/net/netfilter/ipvs/ip_vs.ko: kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/ipvs/ip_vs_rr.ko: kernel/net/netfilter/ipvs/ip_vs.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/ipvs/ip_vs_wrr.ko: kernel/net/netfilter/ipvs/ip_vs.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/ipvs/ip_vs_lc.ko: kernel/net/netfilter/ipvs/ip_vs.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/ipvs/ip_vs_wlc.ko: kernel/net/netfilter/ipvs/ip_vs.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/ipvs/ip_vs_fo.ko: kernel/net/netfilter/ipvs/ip_vs.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/ipvs/ip_vs_ovf.ko: kernel/net/netfilter/ipvs/ip_vs.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/ipvs/ip_vs_lblc.ko: kernel/net/netfilter/ipvs/ip_vs.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/ipvs/ip_vs_lblcr.ko: kernel/net/netfilter/ipvs/ip_vs.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/ipvs/ip_vs_dh.ko: kernel/net/netfilter/ipvs/ip_vs.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/ipvs/ip_vs_sh.ko: kernel/net/netfilter/ipvs/ip_vs.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/ipvs/ip_vs_mh.ko: kernel/net/netfilter/ipvs/ip_vs.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/ipvs/ip_vs_sed.ko: kernel/net/netfilter/ipvs/ip_vs.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/ipvs/ip_vs_nq.ko: kernel/net/netfilter/ipvs/ip_vs.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/ipvs/ip_vs_twos.ko: kernel/net/netfilter/ipvs/ip_vs.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/ipvs/ip_vs_ftp.ko: kernel/net/netfilter/ipvs/ip_vs.ko kernel/net/netfilter/nf_nat.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/netfilter/ipvs/ip_vs_pe_sip.ko: kernel/net/netfilter/ipvs/ip_vs.ko kernel/net/netfilter/nf_conntrack_sip.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko:
kernel/net/ipv4/netfilter/nf_socket_ipv4.ko:
kernel/net/ipv4/netfilter/nf_tproxy_ipv4.ko:
kernel/net/ipv4/netfilter/nf_reject_ipv4.ko:
kernel/net/ipv4/netfilter/nf_nat_h323.ko: kernel/net/netfilter/nf_conntrack_h323.ko kernel/net/netfilter/nf_nat.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/ipv4/netfilter/nf_nat_pptp.ko: kernel/net/netfilter/nf_conntrack_pptp.ko kernel/net/netfilter/nf_nat.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/ipv4/netfilter/nf_nat_snmp_basic.ko: kernel/net/netfilter/nf_conntrack_snmp.ko kernel/net/netfilter/nf_conntrack_broadcast.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/ipv4/netfilter/nft_reject_ipv4.ko: kernel/net/ipv4/netfilter/nf_reject_ipv4.ko kernel/net/netfilter/nft_reject.ko kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/lib/libcrc32c.ko
kernel/net/ipv4/netfilter/nft_fib_ipv4.ko: kernel/net/netfilter/nft_fib.ko kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/lib/libcrc32c.ko
kernel/net/ipv4/netfilter/nft_dup_ipv4.ko: kernel/net/ipv4/netfilter/nf_dup_ipv4.ko kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/lib/libcrc32c.ko
kernel/net/ipv4/netfilter/nf_flow_table_ipv4.ko: kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/net/netfilter/nf_flow_table.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/ipv4/netfilter/ip_tables.ko: kernel/net/netfilter/x_tables.ko
kernel/net/ipv4/netfilter/iptable_filter.ko: kernel/net/ipv4/netfilter/ip_tables.ko kernel/net/netfilter/x_tables.ko
kernel/net/ipv4/netfilter/iptable_mangle.ko: kernel/net/ipv4/netfilter/ip_tables.ko kernel/net/netfilter/x_tables.ko
kernel/net/ipv4/netfilter/iptable_nat.ko: kernel/net/ipv4/netfilter/ip_tables.ko kernel/net/netfilter/nf_nat.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/net/netfilter/x_tables.ko kernel/lib/libcrc32c.ko
kernel/net/ipv4/netfilter/iptable_raw.ko: kernel/net/ipv4/netfilter/ip_tables.ko kernel/net/netfilter/x_tables.ko
kernel/net/ipv4/netfilter/iptable_security.ko: kernel/net/ipv4/netfilter/ip_tables.ko kernel/net/netfilter/x_tables.ko
kernel/net/ipv4/netfilter/ipt_ah.ko: kernel/net/netfilter/x_tables.ko
kernel/net/ipv4/netfilter/ipt_rpfilter.ko: kernel/net/netfilter/x_tables.ko
kernel/net/ipv4/netfilter/ipt_CLUSTERIP.ko: kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/net/netfilter/x_tables.ko kernel/lib/libcrc32c.ko
kernel/net/ipv4/netfilter/ipt_ECN.ko: kernel/net/netfilter/x_tables.ko
kernel/net/ipv4/netfilter/ipt_REJECT.ko: kernel/net/ipv4/netfilter/nf_reject_ipv4.ko kernel/net/netfilter/x_tables.ko
kernel/net/ipv4/netfilter/ipt_SYNPROXY.ko: kernel/net/netfilter/nf_synproxy_core.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/net/netfilter/x_tables.ko kernel/lib/libcrc32c.ko
kernel/net/ipv4/netfilter/arp_tables.ko: kernel/net/netfilter/x_tables.ko
kernel/net/ipv4/netfilter/arpt_mangle.ko: kernel/net/netfilter/x_tables.ko
kernel/net/ipv4/netfilter/arptable_filter.ko: kernel/net/ipv4/netfilter/arp_tables.ko kernel/net/netfilter/x_tables.ko
kernel/net/ipv4/netfilter/nf_dup_ipv4.ko:
kernel/net/ipv4/ip_tunnel.ko:
kernel/net/ipv4/ipip.ko: kernel/net/ipv4/tunnel4.ko kernel/net/ipv4/ip_tunnel.ko
kernel/net/ipv4/fou.ko: kernel/net/ipv4/ip_tunnel.ko kernel/net/ipv6/ip6_udp_tunnel.ko kernel/net/ipv4/udp_tunnel.ko
kernel/net/ipv4/gre.ko:
kernel/net/ipv4/ip_gre.ko: kernel/net/ipv4/ip_tunnel.ko kernel/net/ipv4/gre.ko
kernel/net/ipv4/udp_tunnel.ko:
kernel/net/ipv4/ip_vti.ko: kernel/net/ipv4/tunnel4.ko kernel/net/ipv4/ip_tunnel.ko
kernel/net/ipv4/ah4.ko: kernel/net/xfrm/xfrm_algo.ko
kernel/net/ipv4/esp4.ko: kernel/net/xfrm/xfrm_algo.ko
kernel/net/ipv4/esp4_offload.ko: kernel/net/ipv4/esp4.ko kernel/net/xfrm/xfrm_algo.ko
kernel/net/ipv4/ipcomp.ko: kernel/net/xfrm/xfrm_ipcomp.ko kernel/net/xfrm/xfrm_algo.ko
kernel/net/ipv4/xfrm4_tunnel.ko: kernel/net/ipv4/tunnel4.ko
kernel/net/ipv4/tunnel4.ko:
kernel/net/ipv4/inet_diag.ko:
kernel/net/ipv4/tcp_diag.ko: kernel/net/ipv4/inet_diag.ko
kernel/net/ipv4/udp_diag.ko: kernel/net/ipv4/inet_diag.ko
kernel/net/ipv4/raw_diag.ko: kernel/net/ipv4/inet_diag.ko
kernel/net/ipv4/tcp_bbr.ko:
kernel/net/ipv4/tcp_bic.ko:
kernel/net/ipv4/tcp_cdg.ko:
kernel/net/ipv4/tcp_dctcp.ko:
kernel/net/ipv4/tcp_westwood.ko:
kernel/net/ipv4/tcp_highspeed.ko:
kernel/net/ipv4/tcp_hybla.ko:
kernel/net/ipv4/tcp_htcp.ko:
kernel/net/ipv4/tcp_vegas.ko:
kernel/net/ipv4/tcp_nv.ko:
kernel/net/ipv4/tcp_veno.ko:
kernel/net/ipv4/tcp_scalable.ko:
kernel/net/ipv4/tcp_lp.ko:
kernel/net/ipv4/tcp_yeah.ko: kernel/net/ipv4/tcp_vegas.ko
kernel/net/ipv4/tcp_illinois.ko:
kernel/net/xfrm/xfrm_algo.ko:
kernel/net/xfrm/xfrm_user.ko: kernel/net/xfrm/xfrm_algo.ko
kernel/net/xfrm/xfrm_compat.ko: kernel/net/xfrm/xfrm_user.ko kernel/net/xfrm/xfrm_algo.ko
kernel/net/xfrm/xfrm_ipcomp.ko: kernel/net/xfrm/xfrm_algo.ko
kernel/net/xfrm/xfrm_interface.ko: kernel/net/ipv6/xfrm6_tunnel.ko kernel/net/ipv4/tunnel4.ko kernel/net/ipv6/tunnel6.ko
kernel/net/unix/unix_diag.ko:
kernel/net/ipv6/netfilter/ip6_tables.ko: kernel/net/netfilter/x_tables.ko
kernel/net/ipv6/netfilter/ip6table_filter.ko: kernel/net/ipv6/netfilter/ip6_tables.ko kernel/net/netfilter/x_tables.ko
kernel/net/ipv6/netfilter/ip6table_mangle.ko: kernel/net/ipv6/netfilter/ip6_tables.ko kernel/net/netfilter/x_tables.ko
kernel/net/ipv6/netfilter/ip6table_raw.ko: kernel/net/ipv6/netfilter/ip6_tables.ko kernel/net/netfilter/x_tables.ko
kernel/net/ipv6/netfilter/ip6table_security.ko: kernel/net/ipv6/netfilter/ip6_tables.ko kernel/net/netfilter/x_tables.ko
kernel/net/ipv6/netfilter/ip6table_nat.ko: kernel/net/ipv6/netfilter/ip6_tables.ko kernel/net/netfilter/nf_nat.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/net/netfilter/x_tables.ko kernel/lib/libcrc32c.ko
kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko:
kernel/net/ipv6/netfilter/nf_socket_ipv6.ko:
kernel/net/ipv6/netfilter/nf_tproxy_ipv6.ko:
kernel/net/ipv6/netfilter/nf_reject_ipv6.ko:
kernel/net/ipv6/netfilter/nf_dup_ipv6.ko:
kernel/net/ipv6/netfilter/nft_reject_ipv6.ko: kernel/net/ipv6/netfilter/nf_reject_ipv6.ko kernel/net/netfilter/nft_reject.ko kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/lib/libcrc32c.ko
kernel/net/ipv6/netfilter/nft_dup_ipv6.ko: kernel/net/ipv6/netfilter/nf_dup_ipv6.ko kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/lib/libcrc32c.ko
kernel/net/ipv6/netfilter/nft_fib_ipv6.ko: kernel/net/netfilter/nft_fib.ko kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/lib/libcrc32c.ko
kernel/net/ipv6/netfilter/nf_flow_table_ipv6.ko: kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/net/netfilter/nf_flow_table.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/ipv6/netfilter/ip6t_ah.ko: kernel/net/netfilter/x_tables.ko
kernel/net/ipv6/netfilter/ip6t_eui64.ko: kernel/net/netfilter/x_tables.ko
kernel/net/ipv6/netfilter/ip6t_frag.ko: kernel/net/netfilter/x_tables.ko
kernel/net/ipv6/netfilter/ip6t_ipv6header.ko: kernel/net/netfilter/x_tables.ko
kernel/net/ipv6/netfilter/ip6t_mh.ko: kernel/net/netfilter/x_tables.ko
kernel/net/ipv6/netfilter/ip6t_hbh.ko: kernel/net/netfilter/x_tables.ko
kernel/net/ipv6/netfilter/ip6t_rpfilter.ko: kernel/net/netfilter/x_tables.ko
kernel/net/ipv6/netfilter/ip6t_rt.ko: kernel/net/netfilter/x_tables.ko
kernel/net/ipv6/netfilter/ip6t_srh.ko: kernel/net/netfilter/x_tables.ko
kernel/net/ipv6/netfilter/ip6t_NPT.ko: kernel/net/netfilter/x_tables.ko
kernel/net/ipv6/netfilter/ip6t_REJECT.ko: kernel/net/ipv6/netfilter/nf_reject_ipv6.ko kernel/net/netfilter/x_tables.ko
kernel/net/ipv6/netfilter/ip6t_SYNPROXY.ko: kernel/net/netfilter/nf_synproxy_core.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/net/netfilter/x_tables.ko kernel/lib/libcrc32c.ko
kernel/net/ipv6/ah6.ko: kernel/net/xfrm/xfrm_algo.ko
kernel/net/ipv6/esp6.ko: kernel/net/xfrm/xfrm_algo.ko
kernel/net/ipv6/esp6_offload.ko: kernel/net/ipv6/esp6.ko kernel/net/xfrm/xfrm_algo.ko
kernel/net/ipv6/ipcomp6.ko: kernel/net/ipv6/xfrm6_tunnel.ko kernel/net/xfrm/xfrm_ipcomp.ko kernel/net/ipv6/tunnel6.ko kernel/net/xfrm/xfrm_algo.ko
kernel/net/ipv6/xfrm6_tunnel.ko: kernel/net/ipv6/tunnel6.ko
kernel/net/ipv6/tunnel6.ko:
kernel/net/ipv6/mip6.ko:
kernel/net/ipv6/ila/ila.ko:
kernel/net/ipv6/ip6_vti.ko: kernel/net/ipv6/xfrm6_tunnel.ko kernel/net/ipv6/ip6_tunnel.ko kernel/net/ipv6/tunnel6.ko
kernel/net/ipv6/sit.ko: kernel/net/ipv4/tunnel4.ko kernel/net/ipv4/ip_tunnel.ko
kernel/net/ipv6/ip6_tunnel.ko: kernel/net/ipv6/tunnel6.ko
kernel/net/ipv6/ip6_gre.ko: kernel/net/ipv4/gre.ko kernel/net/ipv6/ip6_tunnel.ko kernel/net/ipv6/tunnel6.ko
kernel/net/ipv6/fou6.ko: kernel/net/ipv4/fou.ko kernel/net/ipv4/ip_tunnel.ko kernel/net/ipv6/ip6_tunnel.ko kernel/net/ipv6/tunnel6.ko kernel/net/ipv6/ip6_udp_tunnel.ko kernel/net/ipv4/udp_tunnel.ko
kernel/net/ipv6/ip6_udp_tunnel.ko:
kernel/net/bpfilter/bpfilter.ko:
kernel/net/packet/af_packet_diag.ko:
kernel/net/8021q/8021q.ko: kernel/net/802/garp.ko kernel/net/802/mrp.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko
kernel/net/wireless/cfg80211.ko:
kernel/net/wireless/lib80211.ko:
kernel/net/wireless/lib80211_crypt_wep.ko: kernel/net/wireless/lib80211.ko kernel/lib/crypto/libarc4.ko
kernel/net/wireless/lib80211_crypt_ccmp.ko: kernel/net/wireless/lib80211.ko
kernel/net/wireless/lib80211_crypt_tkip.ko: kernel/net/wireless/lib80211.ko kernel/lib/crypto/libarc4.ko
kernel/net/rfkill/rfkill-gpio.ko:
kernel/net/mpls/mpls_gso.ko:
kernel/net/mpls/mpls_router.ko: kernel/net/ipv4/ip_tunnel.ko
kernel/net/mpls/mpls_iptunnel.ko: kernel/net/mpls/mpls_router.ko kernel/net/ipv4/ip_tunnel.ko
kernel/net/xdp/xsk_diag.ko:
kernel/net/mptcp/mptcp_diag.ko: kernel/net/ipv4/inet_diag.ko
kernel/net/llc/llc.ko:
kernel/net/llc/llc2.ko: kernel/net/llc/llc.ko
kernel/net/tls/tls.ko:
kernel/net/key/af_key.ko: kernel/net/xfrm/xfrm_algo.ko
kernel/net/bridge/netfilter/nft_meta_bridge.ko: kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/lib/libcrc32c.ko
kernel/net/bridge/netfilter/nft_reject_bridge.ko: kernel/net/ipv4/netfilter/nf_reject_ipv4.ko kernel/net/ipv6/netfilter/nf_reject_ipv6.ko kernel/net/netfilter/nft_reject.ko kernel/net/netfilter/nf_tables.ko kernel/net/netfilter/nfnetlink.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/lib/libcrc32c.ko
kernel/net/bridge/netfilter/nf_conntrack_bridge.ko: kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/lib/libcrc32c.ko
kernel/net/bridge/netfilter/ebtables.ko: kernel/net/netfilter/x_tables.ko
kernel/net/bridge/netfilter/ebtable_broute.ko: kernel/net/bridge/netfilter/ebtables.ko kernel/net/netfilter/x_tables.ko
kernel/net/bridge/netfilter/ebtable_filter.ko: kernel/net/bridge/netfilter/ebtables.ko kernel/net/netfilter/x_tables.ko
kernel/net/bridge/netfilter/ebtable_nat.ko: kernel/net/bridge/netfilter/ebtables.ko kernel/net/netfilter/x_tables.ko
kernel/net/bridge/netfilter/ebt_802_3.ko: kernel/net/netfilter/x_tables.ko
kernel/net/bridge/netfilter/ebt_among.ko: kernel/net/netfilter/x_tables.ko
kernel/net/bridge/netfilter/ebt_arp.ko: kernel/net/netfilter/x_tables.ko
kernel/net/bridge/netfilter/ebt_ip.ko: kernel/net/netfilter/x_tables.ko
kernel/net/bridge/netfilter/ebt_ip6.ko: kernel/net/netfilter/x_tables.ko
kernel/net/bridge/netfilter/ebt_limit.ko: kernel/net/netfilter/x_tables.ko
kernel/net/bridge/netfilter/ebt_mark_m.ko: kernel/net/netfilter/x_tables.ko
kernel/net/bridge/netfilter/ebt_pkttype.ko: kernel/net/netfilter/x_tables.ko
kernel/net/bridge/netfilter/ebt_stp.ko: kernel/net/netfilter/x_tables.ko
kernel/net/bridge/netfilter/ebt_vlan.ko: kernel/net/netfilter/x_tables.ko
kernel/net/bridge/netfilter/ebt_arpreply.ko: kernel/net/netfilter/x_tables.ko
kernel/net/bridge/netfilter/ebt_mark.ko: kernel/net/netfilter/x_tables.ko
kernel/net/bridge/netfilter/ebt_dnat.ko: kernel/net/netfilter/x_tables.ko
kernel/net/bridge/netfilter/ebt_redirect.ko: kernel/net/netfilter/x_tables.ko
kernel/net/bridge/netfilter/ebt_snat.ko: kernel/net/netfilter/x_tables.ko
kernel/net/bridge/netfilter/ebt_log.ko: kernel/net/netfilter/x_tables.ko
kernel/net/bridge/netfilter/ebt_nflog.ko: kernel/net/netfilter/x_tables.ko
kernel/net/bridge/bridge.ko: kernel/net/802/stp.ko kernel/net/llc/llc.ko
kernel/net/bridge/br_netfilter.ko: kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko
kernel/net/dsa/dsa_core.ko: kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/net/dsa/tag_ar9331.ko: kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/net/dsa/tag_brcm.ko: kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/net/dsa/tag_dsa.ko: kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/net/dsa/tag_gswip.ko: kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/net/dsa/tag_hellcreek.ko: kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/net/dsa/tag_ksz.ko: kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/net/dsa/tag_rtl4_a.ko: kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/net/dsa/tag_lan9303.ko: kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/net/dsa/tag_mtk.ko: kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/net/dsa/tag_ocelot.ko: kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/net/dsa/tag_ocelot_8021q.ko: kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/net/dsa/tag_qca.ko: kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/net/dsa/tag_sja1105.ko: kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/net/dsa/tag_trailer.ko: kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/net/dsa/tag_xrs700x.ko: kernel/net/dsa/dsa_core.ko kernel/net/hsr/hsr.ko kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/drivers/net/phy/phylink.ko
kernel/net/appletalk/appletalk.ko: kernel/net/802/psnap.ko kernel/net/llc/llc.ko
kernel/net/x25/x25.ko:
kernel/net/lapb/lapb.ko:
kernel/net/netrom/netrom.ko: kernel/net/ax25/ax25.ko
kernel/net/rose/rose.ko: kernel/net/ax25/ax25.ko
kernel/net/ax25/ax25.ko:
kernel/net/can/can.ko:
kernel/net/can/can-raw.ko: kernel/net/can/can.ko
kernel/net/can/can-bcm.ko: kernel/net/can/can.ko
kernel/net/can/can-gw.ko: kernel/net/can/can.ko
kernel/net/can/j1939/can-j1939.ko: kernel/net/can/can.ko
kernel/net/can/can-isotp.ko: kernel/net/can/can.ko
kernel/net/bluetooth/bluetooth.ko: kernel/crypto/ecdh_generic.ko kernel/crypto/ecc.ko
kernel/net/bluetooth/rfcomm/rfcomm.ko: kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/crypto/ecc.ko
kernel/net/bluetooth/bnep/bnep.ko: kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/crypto/ecc.ko
kernel/net/bluetooth/cmtp/cmtp.ko: kernel/drivers/isdn/capi/kernelcapi.ko kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/crypto/ecc.ko
kernel/net/bluetooth/hidp/hidp.ko: kernel/drivers/hid/hid.ko kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/crypto/ecc.ko
kernel/net/bluetooth/bluetooth_6lowpan.ko: kernel/net/6lowpan/6lowpan.ko kernel/net/bluetooth/bluetooth.ko kernel/crypto/ecdh_generic.ko kernel/crypto/ecc.ko
kernel/net/sunrpc/sunrpc.ko:
kernel/net/sunrpc/auth_gss/auth_rpcgss.ko: kernel/net/sunrpc/sunrpc.ko
kernel/net/sunrpc/auth_gss/rpcsec_gss_krb5.ko: kernel/net/sunrpc/auth_gss/auth_rpcgss.ko kernel/net/sunrpc/sunrpc.ko
kernel/net/sunrpc/xprtrdma/rpcrdma.ko: kernel/drivers/infiniband/core/rdma_cm.ko kernel/drivers/infiniband/core/iw_cm.ko kernel/drivers/infiniband/core/ib_cm.ko kernel/drivers/infiniband/core/ib_core.ko kernel/net/sunrpc/sunrpc.ko
kernel/net/rxrpc/rxrpc.ko: kernel/net/ipv6/ip6_udp_tunnel.ko kernel/net/ipv4/udp_tunnel.ko
kernel/net/kcm/kcm.ko:
kernel/net/atm/atm.ko:
kernel/net/atm/clip.ko: kernel/net/atm/atm.ko
kernel/net/atm/br2684.ko: kernel/net/atm/atm.ko
kernel/net/atm/lec.ko: kernel/net/atm/atm.ko
kernel/net/atm/mpoa.ko: kernel/net/atm/atm.ko
kernel/net/atm/pppoatm.ko: kernel/net/atm/atm.ko
kernel/net/l2tp/l2tp_core.ko: kernel/net/ipv6/ip6_udp_tunnel.ko kernel/net/ipv4/udp_tunnel.ko
kernel/net/l2tp/l2tp_ppp.ko: kernel/net/l2tp/l2tp_netlink.ko kernel/net/l2tp/l2tp_core.ko kernel/drivers/net/ppp/pppox.ko kernel/net/ipv6/ip6_udp_tunnel.ko kernel/net/ipv4/udp_tunnel.ko
kernel/net/l2tp/l2tp_ip.ko: kernel/net/l2tp/l2tp_core.ko kernel/net/ipv6/ip6_udp_tunnel.ko kernel/net/ipv4/udp_tunnel.ko
kernel/net/l2tp/l2tp_netlink.ko: kernel/net/l2tp/l2tp_core.ko kernel/net/ipv6/ip6_udp_tunnel.ko kernel/net/ipv4/udp_tunnel.ko
kernel/net/l2tp/l2tp_eth.ko: kernel/net/l2tp/l2tp_netlink.ko kernel/net/l2tp/l2tp_core.ko kernel/net/ipv6/ip6_udp_tunnel.ko kernel/net/ipv4/udp_tunnel.ko
kernel/net/l2tp/l2tp_debugfs.ko: kernel/net/l2tp/l2tp_core.ko kernel/net/ipv6/ip6_udp_tunnel.ko kernel/net/ipv4/udp_tunnel.ko
kernel/net/l2tp/l2tp_ip6.ko: kernel/net/l2tp/l2tp_ip.ko kernel/net/l2tp/l2tp_core.ko kernel/net/ipv6/ip6_udp_tunnel.ko kernel/net/ipv4/udp_tunnel.ko
kernel/net/decnet/netfilter/dn_rtmsg.ko:
kernel/net/decnet/decnet.ko:
kernel/net/phonet/phonet.ko:
kernel/net/phonet/pn_pep.ko: kernel/net/phonet/phonet.ko
kernel/net/dccp/dccp.ko:
kernel/net/dccp/dccp_ipv4.ko: kernel/net/dccp/dccp.ko
kernel/net/dccp/dccp_ipv6.ko: kernel/net/dccp/dccp_ipv4.ko kernel/net/dccp/dccp.ko
kernel/net/dccp/dccp_diag.ko: kernel/net/dccp/dccp.ko kernel/net/ipv4/inet_diag.ko
kernel/net/sctp/sctp.ko: kernel/net/ipv6/ip6_udp_tunnel.ko kernel/net/ipv4/udp_tunnel.ko kernel/lib/libcrc32c.ko
kernel/net/sctp/sctp_diag.ko: kernel/net/sctp/sctp.ko kernel/net/ipv4/inet_diag.ko kernel/net/ipv6/ip6_udp_tunnel.ko kernel/net/ipv4/udp_tunnel.ko kernel/lib/libcrc32c.ko
kernel/net/rds/rds.ko:
kernel/net/rds/rds_rdma.ko: kernel/net/rds/rds.ko kernel/drivers/infiniband/core/rdma_cm.ko kernel/drivers/infiniband/core/iw_cm.ko kernel/drivers/infiniband/core/ib_cm.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/net/rds/rds_tcp.ko: kernel/net/rds/rds.ko
kernel/net/mac80211/mac80211.ko: kernel/net/wireless/cfg80211.ko kernel/lib/crypto/libarc4.ko
kernel/net/tipc/tipc.ko: kernel/net/ipv6/ip6_udp_tunnel.ko kernel/net/ipv4/udp_tunnel.ko
kernel/net/tipc/diag.ko: kernel/net/tipc/tipc.ko kernel/net/ipv6/ip6_udp_tunnel.ko kernel/net/ipv4/udp_tunnel.ko
kernel/net/smc/smc.ko: kernel/drivers/infiniband/core/ib_core.ko
kernel/net/smc/smc_diag.ko: kernel/net/smc/smc.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/net/9p/9pnet.ko:
kernel/net/9p/9pnet_xen.ko: kernel/net/9p/9pnet.ko
kernel/net/9p/9pnet_virtio.ko: kernel/net/9p/9pnet.ko
kernel/net/9p/9pnet_rdma.ko: kernel/net/9p/9pnet.ko kernel/drivers/infiniband/core/rdma_cm.ko kernel/drivers/infiniband/core/iw_cm.ko kernel/drivers/infiniband/core/ib_cm.ko kernel/drivers/infiniband/core/ib_core.ko
kernel/net/caif/caif.ko:
kernel/net/caif/chnl_net.ko: kernel/net/caif/caif.ko
kernel/net/caif/caif_socket.ko: kernel/net/caif/caif.ko
kernel/net/caif/caif_usb.ko: kernel/net/caif/caif.ko
kernel/net/6lowpan/6lowpan.ko:
kernel/net/6lowpan/nhc_dest.ko: kernel/net/6lowpan/6lowpan.ko
kernel/net/6lowpan/nhc_fragment.ko: kernel/net/6lowpan/6lowpan.ko
kernel/net/6lowpan/nhc_hop.ko: kernel/net/6lowpan/6lowpan.ko
kernel/net/6lowpan/nhc_ipv6.ko: kernel/net/6lowpan/6lowpan.ko
kernel/net/6lowpan/nhc_mobility.ko: kernel/net/6lowpan/6lowpan.ko
kernel/net/6lowpan/nhc_routing.ko: kernel/net/6lowpan/6lowpan.ko
kernel/net/6lowpan/nhc_udp.ko: kernel/net/6lowpan/6lowpan.ko
kernel/net/ieee802154/6lowpan/ieee802154_6lowpan.ko: kernel/net/6lowpan/6lowpan.ko kernel/net/ieee802154/ieee802154.ko
kernel/net/ieee802154/ieee802154.ko:
kernel/net/ieee802154/ieee802154_socket.ko: kernel/net/ieee802154/ieee802154.ko
kernel/net/mac802154/mac802154.ko: kernel/net/ieee802154/ieee802154.ko
kernel/net/ceph/libceph.ko: kernel/lib/libcrc32c.ko
kernel/net/batman-adv/batman-adv.ko: kernel/net/bridge/bridge.ko kernel/net/802/stp.ko kernel/net/llc/llc.ko kernel/lib/libcrc32c.ko
kernel/net/nfc/nfc.ko:
kernel/net/nfc/nci/nci.ko: kernel/net/nfc/nfc.ko
kernel/net/nfc/nci/nci_spi.ko:
kernel/net/nfc/nci/nci_uart.ko:
kernel/net/nfc/hci/hci.ko: kernel/net/nfc/nfc.ko
kernel/net/nfc/nfc_digital.ko: kernel/net/nfc/nfc.ko kernel/lib/crc-itu-t.ko
kernel/net/psample/psample.ko:
kernel/net/ife/ife.ko:
kernel/net/openvswitch/openvswitch.ko: kernel/net/nsh/nsh.ko kernel/net/netfilter/nf_conncount.ko kernel/net/netfilter/nf_nat.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/lib/libcrc32c.ko
kernel/net/openvswitch/vport-vxlan.ko: kernel/net/openvswitch/openvswitch.ko kernel/net/nsh/nsh.ko kernel/net/netfilter/nf_conncount.ko kernel/net/netfilter/nf_nat.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/drivers/net/vxlan.ko kernel/net/ipv6/ip6_udp_tunnel.ko kernel/net/ipv4/udp_tunnel.ko kernel/lib/libcrc32c.ko
kernel/net/openvswitch/vport-geneve.ko: kernel/drivers/net/geneve.ko kernel/net/openvswitch/openvswitch.ko kernel/net/nsh/nsh.ko kernel/net/netfilter/nf_conncount.ko kernel/net/netfilter/nf_nat.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/net/ipv6/ip6_udp_tunnel.ko kernel/net/ipv4/udp_tunnel.ko kernel/lib/libcrc32c.ko
kernel/net/openvswitch/vport-gre.ko: kernel/net/ipv4/ip_gre.ko kernel/net/openvswitch/openvswitch.ko kernel/net/nsh/nsh.ko kernel/net/ipv4/ip_tunnel.ko kernel/net/netfilter/nf_conncount.ko kernel/net/netfilter/nf_nat.ko kernel/net/netfilter/nf_conntrack.ko kernel/net/ipv6/netfilter/nf_defrag_ipv6.ko kernel/net/ipv4/netfilter/nf_defrag_ipv4.ko kernel/net/ipv4/gre.ko kernel/lib/libcrc32c.ko
kernel/net/vmw_vsock/vsock.ko:
kernel/net/vmw_vsock/vsock_diag.ko: kernel/net/vmw_vsock/vsock.ko
kernel/net/vmw_vsock/vmw_vsock_vmci_transport.ko: kernel/net/vmw_vsock/vsock.ko kernel/drivers/misc/vmw_vmci/vmw_vmci.ko
kernel/net/vmw_vsock/vmw_vsock_virtio_transport.ko: kernel/net/vmw_vsock/vmw_vsock_virtio_transport_common.ko kernel/net/vmw_vsock/vsock.ko
kernel/net/vmw_vsock/vmw_vsock_virtio_transport_common.ko: kernel/net/vmw_vsock/vsock.ko
kernel/net/vmw_vsock/hv_sock.ko: kernel/net/vmw_vsock/vsock.ko kernel/drivers/hv/hv_vmbus.ko
kernel/net/vmw_vsock/vsock_loopback.ko: kernel/net/vmw_vsock/vmw_vsock_virtio_transport_common.ko kernel/net/vmw_vsock/vsock.ko
kernel/net/nsh/nsh.ko:
kernel/net/hsr/hsr.ko:
kernel/net/qrtr/qrtr.ko: kernel/net/qrtr/ns.ko
kernel/net/qrtr/ns.ko:
kernel/net/qrtr/qrtr-smd.ko: kernel/net/qrtr/qrtr.ko kernel/net/qrtr/ns.ko kernel/drivers/rpmsg/rpmsg_core.ko
kernel/net/qrtr/qrtr-tun.ko: kernel/net/qrtr/qrtr.ko kernel/net/qrtr/ns.ko
kernel/net/qrtr/qrtr-mhi.ko: kernel/net/qrtr/qrtr.ko kernel/net/qrtr/ns.ko kernel/drivers/bus/mhi/core/mhi.ko
kernel/net/mctp/mctp.ko:
kernel/zfs/zzstd.ko: kernel/zfs/spl.ko
kernel/zfs/zunicode.ko:
kernel/zfs/znvpair.ko: kernel/zfs/spl.ko
kernel/zfs/spl.ko:
kernel/zfs/zcommon.ko: kernel/zfs/znvpair.ko kernel/zfs/spl.ko
kernel/v4l2loopback/v4l2loopback.ko: kernel/drivers/media/v4l2-core/videodev.ko kernel/drivers/media/mc/mc.ko
kernel/zfs/icp.ko: kernel/zfs/zcommon.ko kernel/zfs/znvpair.ko kernel/zfs/spl.ko
kernel/zfs/zavl.ko: kernel/zfs/spl.ko
kernel/zfs/zlua.ko:
kernel/zfs/zfs.ko: kernel/zfs/zunicode.ko kernel/zfs/zzstd.ko kernel/zfs/zlua.ko kernel/zfs/zavl.ko kernel/zfs/icp.ko kernel/zfs/zcommon.ko kernel/zfs/znvpair.ko kernel/zfs/spl.ko
`
