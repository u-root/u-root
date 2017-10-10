# For maintainers only

# We use the github hub tool for code review.

## Setup your u-root Github Repository

Follow the instructions for using github. At some point, you'll need to fork github.com/u-root/u-root, then
get it via go get or gitclone. In any event, your u-root repo should end up in
$GOPATH/src/github.com/u-root/u-root

## Keep an eye on github PR's and provide reviews

# We use govendor for maintaining dependencies.
``u-root`` uses [govendor](https://github.com/kardianos/govendor) for its dependency management.

## To manage dependencies

### Add new dependencies

  - Edit your code to import foo/bar
  - Run `govendor add +external` from the top level

### Remove dependencies

  - Run `govendor remove foo/bar`

### Update dependencies

  - Run `govendor remove +vendor`
  - Run `govendor add +external`

# Style Guide

In [CONTRIBUTING.md](CONTRIBUTING.md) we say `gofmt` and `golint` are law,
but that's not enforced (yet) in automation.

# Maintainers

* [Ron Minnich](https://github.com/rminnich)
* [Andrew Mirtchovski](https://github.com/mirtchovski)
* [Alexandre Beletti](https://github.com/rhiguita)
* [Manoel Machado](https://github.com/ryukinix)
* [Rafael C. Nunes](https://github.com/rafaelcn)
* [Matheus Pinto Rodrigues](https://github.com/mathgamain)
* [Gan Shun](https://github.com/GanShun)
* [Ryan O'Leary](https://github.com/rjoleary)
