# How to create the test fit image

* install u-boot tools
```shell
apt-get install u-boot-tools
```

* create dummy kernel and initramfs
```shell
IMAGE_SIZE=100
yes "k0" | tr -d '\n' | head -c $IMAGE_SIZE > $tmp_dir/dummy_kernel_0
yes "k1" | tr -d '\n' | head -c $IMAGE_SIZE > $tmp_dir/dummy_kernel_1
yes "i0" | tr -d '\n' | head -c $IMAGE_SIZE > $tmp_dir/dummy_initramfs_0.cpio
```

* create signatures to embed
```shell
key0=042FF30F752685F2
key1=9ED18B2103E33767
gpg --import ./key0 ./key1
gpg --default-key $key0 --output $tmp_dir/key0_initram0_pgp.sig --detach-sig $tmp_dir/dummy_initramfs_0.cpio
gpg --default-key $key1 --output $tmp_dir/key1_initram0_pgp.sig --detach-sig $tmp_dir/dummy_initramfs_0.cpio
gpg --default-key $key0 --output $tmp_dir/key0_kernel0_pgp.sig --detach-sig $tmp_dir/dummy_kernel_0
```

* optional clean up keyring
```shell
gpg --delete-secret-key $key0 $key1
gpg --delete-key $key0 $key1
```

*  create itb file
```shell
mkimage -f fitimage.its fitimage.itb
```
