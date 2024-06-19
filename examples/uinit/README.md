# Example User uinit

u-root's built-in init command will call uinit if it exits.

The `uinit.go` in this dir runs a number of commands, including getting an IPv4
address using DHCP, then runs the default shell. When the user exists the shell,
`shutdown halt` is automatically run.
s
To build the example:

```shell
    # make some changes to uinit.go in examples/uinit/uinit.go
    u-root -format=cpio -build=bb -o initramfs.cpio \
        ./cmds/* \
        github.com/u-root/u-root/examples/uinit
```

And then to run it in qemu, with a network:

```shell
    qemu-system-x86_64 \
        -m 4096M \
        -nographic -serial mon:stdio -display none -curses \
        -append "console=ttyS0 " \
        -net nic,vlan=0,model=virtio \
        -net user,vlan=0,hostfwd=tcp::2222-:2222,hostname=u-boot \
        -kernel kernel/vmlinuz-$(KERNEL) \
        -initrd initramfs.cpio
```
