# Contributing to u-root

We need help with this project. Contributions are very welcome. See the [roadmap](roadmap.md), open [issues](https://github.com/u-root/u-root/issues), and join us in [Slack](CONTRIBUTING.md#communication) to talk about your cool ideas for the project.

## Code of Conduct

Conduct collaboration around u-root in accordance to the [Code of
Conduct](https://github.com/u-root/u-root/wiki/Code-of-Conduct).

## Communication

- [Open Source Firmware Slack team](https://osfw.slack.com),
  channel `#u-root`, sign up [here](https://slack.osfw.dev/)
- [Join the mailing list](https://groups.google.com/forum/#!forum/u-root)

## Bugs

- Please submit issues to https://github.com/u-root/u-root/issues

## Coding Style

The ``u-root`` project aims to follow the standard formatting recommendations
and language idioms set out in the [Effective Go](https://golang.org/doc/effective_go.html)
guide, for example [formatting](https://golang.org/doc/effective_go.html#formatting)
and [names](https://golang.org/doc/effective_go.html#names).

`gofmt` and `staticcheck` are law, although this is not automatically enforced
yet and some housecleaning needs done to achieve that.

We have a few rules not covered by these tools:

- Standard imports are separated from other imports. Example:
    ```
    import (
      "regexp"
      "time"

      dhcp "github.com/krolaw/dhcp4"
    )
    ```

## General Guidelines

We want to implement some of the common commands that exist in upstream projects and elsewhere, but we don't need to copy broken behavior. CLI compatibility with existing implementations isn't required. We can add missing functionality and remove broken behavior from commands as needed.

U-root needs to fit onto small flash storage, (eg. 8 or 16MB SPI). Be cognizant of of how your work is increasing u-root's footprint. The current goal is to keep the BB mode `lzma -9` compressed initramfs image under 3MB.

## Pull Requests

We accept GitHub pull requests.

Fork the project on GitHub, work in your fork and in branches, push
these to your GitHub fork, and when ready, do a GitHub pull requests
against https://github.com/u-root/u-root.

u-root uses Go modules for its dependency management, but still vendors
dependencies in the repository pending module support in the build system.
Please run `go mod tidy` and `go mod vendor` and commit `go.mod`, `go.sum`, and
`vendor/` changes before opening a pull request.

Organize your changes in small and meaningful commits which are easy to review.
Every commit in your pull request needs to be able to build and pass the CI tests.

If the pull request closes an issue please note it as: `"Fixes #NNN"`.

### Patch Format

Well formatted patches aide code review pre-merge and code archaeology in
the future.  The abstract form should be:
```
<component>: Change summary

More detailed explanation of your changes: Why and how.
Wrap it to 72 characters.
See [here] (http://chris.beams.io/posts/git-commit/)
for some more good advices.

Signed-off-by: <contributor@foo.com>
```

An example from this repo:
```
tcz: quiet it down

It had a spurious print that was both annoying and making
boot just a tad slower.

Signed-off-by: Ronald G. Minnich <rminnich@gmail.com>
```

### Developer Sign-Off

For purposes of tracking code-origination, we follow a simple sign-off
process.  If you can attest to the [Developer Certificate of
Origin](https://developercertificate.org/) then you append in each git
commit text a line such as:
```
Signed-off-by: Your Name <username@youremail.com>
```

### Incorporation of Feedback
To not break the conversation history inside the PR avoid force pushes. Instead, push further _'fix up commits'_ to resolve annotations.
Once review is done, do a local rebase to clean up the _'fix up commits'_ and come back to a clean commit history and do a single fore push to the PR.

## Unit Testing Guidelines

### Unit Test Checks

* The [testify](https://github.com/stretchr/testify) package should not be used.
* The [cmp](https://pkg.go.dev/github.com/google/go-cmp/cmp) package is allowed.
* Unit tests in Go should follow the guidelines in this tutorial: https://go.dev/doc/tutorial/add-a-test
  * In particular, the test error should be in the form `Function(...) = ...; want ...`.

For example:

```
if msg != "" || err == nil {
	t.Fatalf(`Hello("") = %q, %v, want "", error`, msg, err)
}
```

* Tests should do all filesystem changes under a temporary directory, either
  created with `ioutil.TempDir` or `testing.T.TempDir`.

### Mocking Dependencies

* Special mocking packages should not be used.
* Prefer to test the real thing where possible.
* Mocking is sometimes necessary. For example:
   * Operations as root.
   * Interacting with special hardware (ex: USB, SPI, PCIe)
   * Modifying machine state (ex: reboot, kexec)
   * Tests which would otherwise be flaky (ex: `time.Sleep`, networking)
* Consider writing an integration test if the program can not be easily run
  directly.
   * Integration tests let you run the command under qemu, which lets you test
     operations with, e.g., virtual hardware.
  * `pkg/mount` contains an example of an integration test run under QEMU.
* Prefer to mock using existing interfaces. For example: `io.Reader`, `fs.FS`
* Avoid mocking using global state. Instead, consider using one of the
  following "dependency injection" idioms:

1. Mocking functions:

```
// The exported function has a meaningful API.
func SetMemAddr(addr, val uint) error {
	return setMemAddr(addr, val, "/dev/mem")
}

// The internal function is called from the unit test. The test can set a
// different `file` argument.
func setMemAddr(addr, val uint, file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	...
}
```

2. Mocking objects:

```
// SPI interface for the underlying calls to the SPI driver.
type SPI interface {
	Transfer(transfers []spidev.Transfer) error
}

// Flash provides operations for SPI flash chips.
type Flash struct {
	// spi is the underlying SPI device.
	spi SPI

  // other fields ...
}

// New creates a new flash device from a SPI interface.
func New(spi SPI) (*Flash, error) {
	f := &Flash{
		spi: spi,
	}

  // initialize other fields ...

  return f
}
```

In the above example, the `flash.New` function takes a `SPI` device which can
be mocked out as follows:

```
f, err := flash.New(spimock.New())
...
```

The `spimock.New()` function returns an implementation of SPI which mocks the
underlying SPI driver. The `Flash` object can be tested without hardware
dependencies.

### VM Tests using QEMU

For code reading or manipulating hardware, it can be reasonable not to mock out
syscalls but to run tests in a virtual environment. In case you need to test
against certain hardware, you can use a QEMU environment via
[vmtest](https://github.com/hugelgupf/vmtest).

In your package, put the setup and corresponding tests in `vm_test.go`.

**IMPORTANT Notes**
* Add `!race` build tag to your `vm_test.go`
* Setup QEMU inside a usual test function
* Make sure the tests assuming this setup are skipped in non-VM test run
* Add your package to `blocklist` in `integration/gotests/gotest_test.go` making
  sure the test doesn't run in the project wide integration tests without the
  proper QEMU setup.

See [/pkg/gpio/gpio_integration_test.go](pkg/gpio) for a simple example, or
[/pkg/mount/block/vm_test.go](pkg/mount).

### Package main

The main function often includes things difficult to test. For example:

1. Process-ending functions such as `log.Fatal` and `os.Exit`. These functions
   also kill the unit test process.
2. Accessing global state such as `os.Args`, `os.Stdin` and `os.Stdout`. It is
   hard to mock out global state cleanly and safely.
   
**Do not use `pkg/testutil` it is deprecated. Instead go with the following:**

The guideline for testing is to factor out everything "difficult" into a
two-line `main` function which remain untested. For example:

```
func run(args []string, stdin io.Reader, stdout io.Writer) error {
	...
}

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
```

### Integration Tests
To test integration with other code/packages without mocking, write integration test in a file `integration_test.go`.

## Code Reviews

Look at the area of code you're modifying, its history, and consider
tagging some of the [maintainers](https://u-root.tk/#contributors) when doing a
pull request in order to instigate some code review.

## Quality Controls

[CircleCI](https://circleci.com/gh/u-root/u-root) is used to test and build commits in a pull request.

See [.circleci/config.yml](.circleci/config.yml) for the CI commands run.
You can use [CircleCI's CLI tool](https://circleci.com/docs/2.0/local-cli/#run-a-job-in-a-container-on-your-machine) to run individual jobs from `.circlecl/config.yml` via Docker, eg. `circleci build --job test`.
