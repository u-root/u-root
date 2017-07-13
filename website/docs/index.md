<center>
  <h3>Behold the Gopher!</h3>
  <img src="img/u-root-logo.png" alt="u-root logo" width=300 />
</center>


# u-root

u-root is an embeddable root file system intended to be placed in a flash device
as part of the firmware image, along with a Linux kernel. Unlike most embedded
root file systems, which consist of large binaries, u-root only has five: an
init program and four Go compiler binaries.


## Setup

On an Ubuntu system, install perquisites and ensure Go is at least version 1.7:

```sh
sudo apt-get install git golang build-essential
go version
```

Set your `GOPATH`:

```sh
export GOPATH="$HOME/go"
```

Clone u-root:

```sh
go get github.com/u-root/u-root
cd "$GOPATH/src/github.com/u-root/u-root"
```

Test u-root inside a chroot:

```sh
go run scripts/ramfs.go -test
```


## Using elvish: a more handy shell

In default, rush is the shell in u-root. Now, thanks to Qi Xiao(\<xiaqqaix@gmail.com\>), u-root users are also able to use a friendly and expressive unix-like shell: __elvish__. Users are free to choose whether to include elvish in u-root or not. Basically, elvish has handy functionalities such as auto completion, command-line existence checks, etc. More info of elvish can be found at: [http://github/elves/elvish](http://github.com/elves/elvish).

If you prefer to use elvish as shell in u-root, here are the instructions:

1. Get project __elvish__:
  ```sh
  go get github.com/elves/elvish
  ```

2. Temporarily, since package `sqlite3` used in elvish has been updated, and its latest
   version includes codes in C (which u-root does not support), users have to
   roll back to last good commit of elvish:
   ```sh
   cd $GOPATH/src/elves/elvish
   git checkout bc5543aef2c493b658d6bd1bb81e3de298de8d2f
   ```

3. Go to u-root repo. If you did `go get github.com/u-root/u-root` before, do:
  ```sh
  cd $GOPATH/src/u-root/u-root
  ```

4. If you prefer to build under bb mode, please do the following command line
   in u-root/u-root/:
   ```sh
   cd ./bb/
   go build .
   CGO_ENABLED=0 ./bb 'src/github.com/u-root/u-root/cmds/[a-z]*' src/github.com/elves/elvish
   ```
   which generates a cpio file, /tmp/initramfs.linux\_amd64.cpio for you to
   start up u-root in qemu.

   If you prefer dynamic buildup mode, do the following command line in u-root/u-root:
   ```sh
   CGO_ENABLED=0 go run scripts/ramfs.go 'src/github.com/u-root/u-root/cmds/[a-z]*' src/github.com/elves/elvish
   ```
   which also generates /tmp/initramfs.linux\_amd64.cpio.

5. Afterwards, users can type command line `elvish` in u-root and start to use elvish as shell.



## Submitting Changes

We use [GitHub Pull Requests](https://github.com/u-root/u-root/pulls) for code
review. Pull requests must receive one approval and pass CI before being merged.

For convenience, it is recommended to use this pre-commit hook:

```sh
ln -s -f ../../scripts/pre-commit .git/hooks/pre-commit
```


## FAQs

### So, why "u-root"?

It's to reflect a universal root, you can mount on every
local and get a userland portable (it's a goal).

### Any publications?

- [USENIX 2015 ATC Paper](https://www.usenix.org/system/files/conference/atc15/atc15-paper-minnich.pdf)
- [USENIX 2015 ATC Talk](https://www.usenix.org/conference/atc15/technical-session/presentation/minnich)


## Community

- [Join the mailing list](https://groups.google.com/forum/#!forum/u-root)
- [Join slack](https://u-root.slack.com/)
- [Checkout the roadmap](https://github.com/u-root/u-root/blob/master/roadmap.md)


## Contributors

* [Ron Minnich](https://github.com/rminnich)
* [Andrew Mirtchovski](https://github.com/mirtchovski)
* [Alexandre Beletti](https://github.com/rhiguita)
* [Manoel Machado](https://github.com/ryukinix)
* [Rafael C. Nunes](https://github.com/rafaelcn)
* [Matheus Pinto Rodrigues](https://github.com/mathgamain)
* [Gan Shun](https://github.com/GanShun)
* [Ryan O'Leary](https://github.com/rjoleary)

