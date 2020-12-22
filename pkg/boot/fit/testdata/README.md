# How to create the test fit image

* install u-boot tools
```shell
apt-get install u-boot-tools
```

* create dummy kernel and initramfs
```shell
head -c 1000 /dev/urandom > $tmp_dir/dummy_kernel
head -c 1000 /dev/urandom > $tmp_dir/dummy_initramfs.cpio
```

*  create itb file
```shell
mkimage -f $tmp_dir/config.its fitimage.itb
```
