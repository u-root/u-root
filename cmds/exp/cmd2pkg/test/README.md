This directory has test commands with interesting dependency trees.

Interesting dependency trees are, for example:

-   mutually and/or diamondly depending Go modules (but no package cycles --
    which is valid):

```
test/mod1/cmd/hello -> test/mod1/pkg/hello
test/mod1/cmd/hello -> test/mod2/pkg/exthello

test/mod2/pkg/exthello -> test/mod1/pkg/hello
test/mod2/pkg/exthello -> test/mod3/pkg/hello
```

-   a dependency whose $GOPATH would not match its Go import path, e.g.
    `github.com/u-root/gobusybox/test/mod4/v2` in `./test/mod4/`.

-   `mod5` depends on `mod6/pkg/mod5hello`, but nothing in `mod5` does. When
    `mod5` and `mod6` commands are in a busybox together, the local version of
    `pkg/mod5hello` should be used.

-   two different modules depending on one third-party module at different
    versions
