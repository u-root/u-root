## Go Busybox

`bb.go` in this package implements a Go source-to-source transformation on pure
Go code (no cgo).

This AST transformation does the following:

-   Takes a Go command's source files and rewrites them into Go package files
    without global side effects.
-   Writes a `main.go` file with a `main()` that calls into the appropriate Go
    command package based on `argv[0]`.

This allows you to take two Go commands, such as Go implementations of `sl` and
`cowsay` and compile them into one binary.

Which command is invoked is determined by `argv[0]` or `argv[1]` if `argv[0]` is
not recognized. Let's say `bb` is the compiled binary; the following are
equivalent invocations of `sl` and `cowsay`:

```sh
# Make a symlink sl -> bb
ln -s bb sl
./sl -l

# Make a symlink cowsay -> bb
ln -s bb cowsay
./cowsay Haha
```

```sh
./bb sl -l
./bb cowsay Haha
```

### AST Transformation

Principally, the AST transformation moves all global side-effects into callable
package functions. E.g. `main` becomes `Main`, each `init` becomes `InitN`, and
global variable assignments are moved into their own `InitN`.

Then, these `Main` and `Init` functions can be registered with a global map of
commands by name and used when called upon.

Let's say a command `github.com/org/repo/cmds/sl` contains the following
`main.go`:

```go
package main

import (
  "flag"
  "log"
)

var name = flag.String("name", "", "Gimme name")

func init() {
  log.Printf("init")
}

func main() {
  log.Printf("train")
}
```

This would be rewritten to be:

```go
package sl // based on the directory name or bazel-rule go_binary name

import (
  "flag"
  "log"

  // This package holds the global map of commands.
  "github.com/u-root/u-root/pkg/bb"
)

// Type has to be inferred through type checking.
var name string

func Init0() {
  log.Printf("init")
}

func Init1() {
  name = flag.String("name", "", "Gimme name")
}

func Init() {
  // Order is determined by go/types.Info.InitOrder.
  Init0()
  Init1()
}

func Main() {
  log.Printf("main")
}

func init() {
  // Register `sl` as a command.
  bb.Register("sl", Init, Main)
}
```

#### Shortcomings

-   If there is already a function `Main` or `InitN` for some `N`, there may be
    a compilation error.
-   Any packages imported by commands may still have global side-effects
    affecting other commands. Done properly, we would have to rewrite all
    non-standard-library packages as well as commands. This has not been
    necessary to implement so far. It would likely be necessary if two different
    imported packages register the same flag unconditionally globally.

## Generated main

The main file can be generated based on any template Go files, but the default
looks something like the following:

```go
import (
  "os"

  "github.com/u-root/u-root/pkg/bb"

  // Side-effect import registers command with bb.
  _ "github.com/org/repo/cmds/generated/sl"
)

func main() {
  bb.Run(os.Argv[0])
}
```

The default template will use `argv[1]` if `argv[0]` is not in the map.
