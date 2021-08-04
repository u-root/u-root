# How to build a UEFI payload

*   Obtain edk2

```shell
git clone --branch uefipayload --recursive https://github.com/linuxboot/edk2 uefipayload
```

*   Follow setup instructions in
    [Get Started with EDK II](https://github.com/tianocore/tianocore.github.io/wiki/Getting-Started-with-EDK-II)

*   build UEFI payload

```shell
make -C BaseTools
source edksetup.sh
build -a X64 -p UefiPayloadPkg/UefiPayloadPkg.dsc -b DEBUG -t GCC5 -D BOOTLOADER=LINUXBOOT
# payload will be in Build/UefiPayloadPkgX64/DEBUG_GCC5/FV/UEFIPAYLOAD.fd
```
