package golang

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Environ struct {
	// GOARCH.
	Arch string

	// GOOS.
	OS string

	// GOPATH.
	Gopaths []string

	// GOROOT.
	Root string
}

func Guess() Environ {
	var env Environ
	if p := os.Getenv("GOPATH"); p != "" {
		env.Gopaths = append(env.Gopaths, strings.Split(p, ":")...)
	}

	if a := os.Getenv("GOARCH"); a != "" {
		env.Arch = a
	} else {
		env.Arch = runtime.GOARCH
	}

	if os := os.Getenv("GOOS"); os != "" {
		env.OS = os
	} else {
		env.OS = runtime.GOOS
	}

	if r := os.Getenv("GOROOT"); r != "" {
		env.Root = r
	} else {
		env.Root = runtime.GOROOT()
	}
	return env
}

// FindPackage returns the full path to `pkg` according to the context's Gopaths.
func (c Environ) FindPackage(pkg string) (string, error) {
	for _, gopath := range c.Gopaths {
		path := filepath.Join(gopath, pkg)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		} else if err != nil {
			return "", fmt.Errorf("error stat(%q): %v", path, err)
		} else {
			return path, nil
		}
	}

	return "", fmt.Errorf("failed to find package %q in gopaths %v", pkg, c.Gopaths)
}

func (c Environ) Env() []string {
	var env []string
	if c.Arch != "" {
		env = append(env, fmt.Sprintf("GOARCH=%s", c.Arch))
	}
	if c.OS != "" {
		env = append(env, fmt.Sprintf("GOOS=%s", c.OS))
	}
	if c.Root != "" {
		env = append(env, fmt.Sprintf("GOROOT=%s", c.Root))
	}
	if len(c.Gopaths) > 0 {
		env = append(env, fmt.Sprintf("GOPATH=%s", strings.Join(c.Gopaths, ":")))
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
