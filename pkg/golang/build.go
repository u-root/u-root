package golang

import (
	"encoding/json"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Environ struct {
	build.Context
}

func Default() Environ {
	return Environ{Context: build.Default}
}

// FindPackageByPath gives the full Go package name for the package in `path`.
//
// This currently assumes that packages are named after the directory they are
// in.
func (c Environ) FindPackageByPath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	p, err := c.Context.ImportDir(abs, 0)
	if err != nil {
		return "", fmt.Errorf("failed to find package in %q: %v", path, err)
	}
	return p.ImportPath, nil
}

// FindPackageDir returns the full path to `pkg` according to the context's Gopaths.
func (c Environ) FindPackageDir(pkg string) (string, error) {
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
	Root       string
	ImportPath string
}

func (c Environ) ListDeps(pkg string) (*ListPackage, error) {
	// The output of this is almost the same as build.Import, except for
	// the dependencies.
	cmd := exec.Command("go", "list", "-json", pkg)
	env := os.Environ()
	env = append(env, c.Env()...)
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
	var cgo int8
	if c.CgoEnabled {
		cgo = 1
	}
	env = append(env, fmt.Sprintf("CGO_ENABLED=%d", cgo))
	return env
}

func (c Environ) String() string {
	return strings.Join(c.Env(), " ")
}

// Optional arguments to Environ.Build.
type BuildOpts struct {
	// ExtraArgs to `go build`.
	ExtraArgs []string
}

// Build compiles `pkg`, writing the executable or object to `binaryPath`.
func (c Environ) Build(pkg string, binaryPath string, opts BuildOpts) error {
	path, err := c.FindPackageDir(pkg)
	if err != nil {
		return err
	}

	args := []string{
		"build",
		"-a", // Force rebuilding of packages.
		"-o", binaryPath,
		"-installsuffix", "uroot",
		"-ldflags", "-s -w", // Strip all symbols.
	}
	if opts.ExtraArgs != nil {
		args = append(args, opts.ExtraArgs...)
	}
	// We always set the working directory, so this is always '.'.
	args = append(args, ".")

	cmd := exec.Command("go", args...)
	cmd.Dir = path
	cmd.Env = append(os.Environ(), c.Env()...)

	if o, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error building go package %v: %v, %v", pkg, string(o), err)
	}
	return nil
}
