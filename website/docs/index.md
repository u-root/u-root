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


## Submitting Changes

We use [GitHub Pull Requests](https://github.com/u-root/u-root/pulls) for code
review. Pull requests must receive one approval and pass CI before being merged.

For convenience, it is recommended to use this pre-commit hook:

```sh
ln -s -f ../../scripts/pre-commit .git/hooks/pre-commit
```


## How to Create a Good UNIX Interface

One of our goals with this project is learning how to write programs and, more importantly,
the importance of not writing code vs. writing code. 

Here are some 
- [good rules for writing good unix programs](https://lasr.cs.ucla.edu/ficus-members/geoff/interfaces.html)

A good example can be seen in notions of progress bars. The idea of progress bars has come up twice now, from 
two different contributors, once in cp, and once in dd.

What's wrong with adding a progress bar to a program? In short, everything!
Where does the output go? stdout? Then you can not add it to some programs because their output is to stdout. 
So do you add it to stderr? But it's not an error, is it? It's informational. 

Further, once you have added it to one program, how do other programs get progress bars? They'll all end up
getting their own, and they'll all look different. 

Since we're using Go, there's a way to get progress bars and follow the rules of Unix programs:
[expvars](https://golang.org/pkg/expvar/). 

[And there are beautiful viewers for them.](https://github.com/divan/expvarmon).

The basic rule for expvars and u-root should be:
- conditional on a common flag
- off by default

Then we can imagine a program, called progress, which we invoke something like this:

```sh
progress dd if=/dev/zero of=/dev/null
```

Progress can start dd with, e.g., -expvars=true, and dd will then start the service. We can have u-root-wide definitions for default expvars, such as progress, percent, bytes, and so on; and progress can query and display them. 

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

