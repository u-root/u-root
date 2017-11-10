package golang

import (
	"encoding/json"
	"fmt"
	"go/build"
	"os"
	"os/exec"
)

type Environ struct {
	build.Context
}

func Default() Environ {
	return Environ{Context: build.Default}
}

// FindPackage returns the full path to `pkg` according to the context's Gopaths.
func (c Environ) FindPackage(pkg string) (string, error) {
	p, err := c.Context.Import(pkg, "", 0)
	if err != nil {
		return "", fmt.Errorf("failed to find package %q in gopath %q", pkg, c.Context.GOPATH)
	}
	return p.Dir, nil
}

func (c Environ) ListPackage(pkg string) (*build.Package, error) {
	return c.Context.Import(pkg, "", 0)
}

type ListPackage struct {
	Dir        string
	Deps       []string
	GoFiles    []string
	SFiles     []string
	HFiles     []string
	Goroot     bool
	ImportPath string
}

type ListOpts struct {
	CGO bool
}

func (c Environ) ListDeps(pkg string, opts ListOpts) (*ListPackage, error) {
	// The output of this is almost the same as build.Import, except for
	// the dependencies.
	cmd := exec.Command("go", "list", "-json", pkg)
	env := os.Environ()
	env = append(env, c.Env()...)
	if !opts.CGO {
		env = append(env, "CGO_ENABLED=0")
	}
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var p ListPackage
	if err := json.Unmarshal(out, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (c Environ) Env() []string {
	var env []string
	if c.GOARCH != "" {
		env = append(env, fmt.Sprintf("GOARCH=%s", c.GOARCH))
	}
	if c.GOOS != "" {
		env = append(env, fmt.Sprintf("GOOS=%s", c.GOOS))
	}
	if c.GOROOT != "" {
		env = append(env, fmt.Sprintf("GOROOT=%s", c.GOROOT))
	}
	if c.GOPATH != "" {
		env = append(env, fmt.Sprintf("GOPATH=%s", c.GOPATH))
	}
	return env
}

// Optional arguments to Environ.Build.
type BuildOpts struct {
	// Whether to enable CGO for the build.
	CGO bool

	// ExtraArgs to `go build`.
	ExtraArgs []string
}

// Build compiles `pkg`, writing the executable or object to `binaryPath`.
func (c Environ) Build(pkg string, binaryPath string, opts BuildOpts) error {
	path, err := c.FindPackage(pkg)
	if err != nil {
		return err
	}

	args := []string{
		"build",
		"-a", // Force rebuilding of packages.
		"-o", binaryPath,
		"-installsuffix", "cgo",
		"-ldflags", "-s -w", // Strip all symbols.
	}
	if opts.ExtraArgs != nil {
		args = append(args, opts.ExtraArgs...)
	}
	// We always set the working directory, so this is always '.'.
	args = append(args, ".")

	var env []string
	env = os.Environ()
	if !opts.CGO {
		env = append(env, "CGO_ENABLED=0")
	}
	env = append(env, c.Env()...)

	cmd := exec.Command("go", args...)
	cmd.Dir = path
	cmd.Env = env

	if o, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error building go package %v: %v, %v", pkg, string(o), err)
	}
	return nil
}
