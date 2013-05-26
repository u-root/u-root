u-root
======

A universal root. You mount it, and it's mostly Go source with the exception of 5 binaries. 

And that's the interesting part. This set of utilities is all Go, and mostly source.

The /bin should be mounted in a tmpfs. The directory with the source should be in your path.
The bin in ram comes in your path before the directory with the source code.

When you run a command that is not built, you fall through to the command that does a
'go build' of the command, and then execs the command once it is built. From that point on,
when you run the command, you get the one in tmpfs. This is fast.
