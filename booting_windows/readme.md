# Running Windows over u-root

## Prerequisites
**WARNING**: the scripts will require sudo priviliges. It is highly advised that
you read them carefully before running them!

1. Golang, see `https://golang.org/`.
   Don't forget to have go/bin in your path via `export PATH=${HOME}/go/bin:${PATH}`
1. Packages required to build Linux kernel
1. A functional Windows Server 2019 or Windows 10 **raw** image, asumed to exist at
   `"${WORKSPACE}"/windows.img`. `setup.sh` will create a masking image over it,
   so the original image will not be modified. Windows boot manager,
   bootmgfw.efi is assumed to exist in the 2nd partition of the image. See
   `install_windows.sh` for an example.
1. An environment variable `EFI_WORKSPACE`, where files will be downloaded to or
   otherwise created.
1. kpartx . Install via `sudo apt-get install kpartx`
1. alien. Install via `sudo apt-get install alien`
## Installing the Modified u-root
1.  Install u-root:

    ```shell
    go get github.com/u-root/u-root
    ```

1.  Change the uroot github remote to our modified one:

    ```shell
    pushd ~/go/src/github.com/u-root/u-root
    git remote add oweisse https://github.com/oweisse/u-root  # our revised uroot repo
    git fetch oweisse
    git checkout -b kexec_test oweisse/kexec_test
    go install
    popd
    ```

## Setting up the Kernel.
Setup Linux kernel source tree with our modifications. We modified
kexec_load syscall to launch EFI applications. **Read the script before running!**

The script will:
1. Download an EFI loader image.
1. Extract windows boot-manager from the windows image (see prerequisites above).
1. Clone our forked linux kernel from `https://github.com/oweisse/linux/`
   into $EFI_WORKSPACE.
1. Install prerequisites (sudo required)
1. Build Linux kernel

```
./setup.sh
```

## Running u-Root and booting Windows
The script command line arguments:
1. rebuild_uroot: will also add bootmgfw.efi, extracted from the Windows image
   to the filesystem.
1. rebuild_kernel: Only necessary if you modified the kernel at
   $EFI_WORKSPACE/linux

```
./run_vm.sh rebuild_uroot rebuild_kernel
```

After u-root has loaded, launch Windows bootmanager"
```
pekexec bootmgfw.efi
```

## Attaching gdb
In the following command, `vmlinux` is the Linux kernel we built.
`launch_efi_app` is the function jumping into `bootmgfw.efi` entry point.
Port 1234 is QEMU default port when using the `-s` flag (see `run_vm.sh`).
```
gdb vmlinux -ex "target remote :1234"     \
            -ex "hbreak launch_efi_app"   \
            -ex "layout regs"             \
            -ex "focus next"              \
            -ex "focus next" -ex "c"
```

