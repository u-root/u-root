<center>
<img src="img/u-root-logo.png" alt="u-root logo" width=300 />
</center>


# u-root

u-root is an embeddable root file system intended to be placed in a flash device
as part of the firmware image, along with a Linux kernel. Unlike most embedded
root file systems, which consist of large binaries, u-root only has five: an
init program and four Go compiler binaries.

## Setup

On an Ubuntu system, install prerequisites and ensure Go is at least version 1.13:

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
```

Generate an initramfs containing u-root Go tools:

```sh
u-root -format=cpio -o initramfs.cpio
```

You can use this initramfs with your favorite Linux kernel in QEMU to try it
out.

More instructions can be found in the repo's
[README.md](https://github.com/u-root/u-root/blob/master/README.md).

## Submitting Changes

We use [GitHub Pull Requests](https://github.com/u-root/u-root/pulls) for code
review. Pull requests must receive one approval and pass CI before being merged.


## FAQs

### So, why "u-root"?

It's to reflect a universal root, you can mount on every
local and get a userland portable (it's a goal).

### Any publications?

- [USENIX 2015 ATC Paper](https://www.usenix.org/system/files/conference/atc15/atc15-paper-minnich.pdf)
- [USENIX 2015 ATC Talk](https://www.usenix.org/conference/atc15/technical-session/presentation/minnich)
- Related: Embedded Linux Conference 2017 LinuxBoot Talk ([YouTube video](https://www.youtube.com/watch?v=iffTJ1vPCSo), [slides](https://schd.ws/hosted_files/osseu17/84/Replace%20UEFI%20with%20Linux.pdf))


## Community

- [Join the mailing list](https://groups.google.com/forum/#!forum/u-root)
- [Join the Open Source Firmware Slack team](https://osfw.slack.com/) (Get an invite [here](https://slack.osfw.dev).)
- [Checkout the roadmap](https://github.com/u-root/u-root/blob/master/roadmap.md)


## Contributors

* [Ron Minnich](https://github.com/rminnich)
* [Andrey Mirtchovski](https://github.com/mirtchovski)
* [Alexandre Beletti](https://github.com/rhiguita)
* [Manoel Machado](https://github.com/ryukinix)
* [Rafael C. Nunes](https://github.com/rafaelcn)
* [Matheus Pinto Rodrigues](https://github.com/mathgamain)
* [Gan Shun Lim](https://github.com/GanShun)
* [Ryan O'Leary](https://github.com/rjoleary)
* [Chris Koch](https://github.com/hugelgupf)
* [Andrea Barberio](https://github.com/insomniacslk)
* [Jean-Marie Verdun](https://github.com/vejmarie)
* [Max Shegai](https://github.com/n-canter)

## Logo

The Go gopher was designed by Renee French.
The u-root logo design is licensed under the Creative Commons 3.0 Attributions license.

The logo is communicating several things:

- u-root has several flavors: firmware and as a root file system

- the gopher at the bottom is a firmware u-root; that gopher brings the machine up, hence the wrench. Its work is also done, so it is resting.

- the other gophers can make more copies of u-root; hence the u-root logo on their chest.

- the highest gopher is showing how u-root is a good root file system for a VM.

- the U itself is a stylized tree, evocative of roots.
