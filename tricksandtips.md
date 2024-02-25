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
