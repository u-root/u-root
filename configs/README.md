# Kernel Configs

These configurations define the core Kconfig options you need for a kernel that
supports the Go runtime, and thus u-root.
See also: <https://go.dev/wiki/MinimumRequirements>

To build a small, flash-ready kernel, you would start from a minimal defconfig,
something like:

```shell
make tinyconfig
cat amd64_config.txt >> .config
make oldconfig
make
```

Or some similar sequence (it has changed over the 15 years of this project).

The exact process is not important, what is important is that you need the options
from these examples to make Go work. For one simple example, Go needs futex and that
is not included in the tinyconfig default.

The files are in several sections, bracketed with comments. The first section includes
what you need for Go; u-root init; serial console; early_printk; an initrd; devtmpfs (used by u-root
init); ELF programs (yes, this is not included by default!); and /proc.

The second section is what you need so the io command will work. 

The third is what you need to build a kernel that can kexec another. The fourth,
for building a kernel that will itself be kexec'd.

Avoid the temptation to add to these files. They are intended to be the absolute smallest
config that will still yield a working u-root initramfs. They have changed very little
since we created them, meaning they can be used across a wide range of kernels.

Append these files to your tinyconfig.

- `amd64_config.txt`: for amd64/x86_64
- `arm_config.txt`: for arm
- `generic_config.txt`: Various architecture agnostic configs used by various
  integration tests.
