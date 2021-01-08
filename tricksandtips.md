# u-root Tricks and Tips

Here are some things we do that you might find useful.

## Pipe u-root output to stdout

One thing to keep in mind is that on many systems, now, standard
out is a named file you can write to. Further, you can 
create a cpio that has no Linux initramfs content. Finally, you can
merge cpios to form a more comprehensive cpio.

## inventory the cpio archive without creating it.
```
u-root -o /dev/stdout | cpio -ivt
```

## Make a pure u-root cpio with no other files.

```
u-root -base /dev/null
```

## Make an image that is BOTH dynamic and busybox

Well, it *almost* works ... help wanted.

Because busybox binaries go in /bbin
and source binaries go in /ubin. This is useful because
you can get a faster boot but still have a Go toolchain

```
u-root -base /dev/null -o /tmp/busybox.cpio
u-root -base /tmp/busybox.cpio -build source -o combined.cpio
```

To see what that looks like:
```
cpio -ivt < combined.cpio
```

You can use testramfs to run it and watch it fork-bomb itself
on init. So close! Working on it.
