// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	ErrNoWorkers           = errors.New("builder: no workers")
	ErrBuilderNoJobs       = errors.New("builder: no jobs to run")
	ErrBuilderInvalidState = errors.New("builder: invalid state")
)

// BuildJob defines a job that can be scheduled by the Builder and
// executed by the work function. It has all necessary information
// to execute the build and gather error logs.
type BuildJob struct {
	GoPkgPath    string   // the path to the go package that is build
	Compiler     string   // compiler used to build the pkg, probably go or tinygo. TODO: refactor to enum
	BuildCommand string   // complete build command issues by builder
	buildDir     *string  // directory where the build binary should be put in the context of go/tinygo this means -o <buildDir>
	env          []string // environment that the execution worker should use additionally
}

// Generate a new BuildJob structure and verify all provided data.
// If the path to the goPkg is invalid, it will return an error.
// TODO: should we enforce that?
func NewBuildJob(goPkgPath string, compiler string, buildDir *string) (BuildJob, error) {
	absPath, err := filepath.Abs(goPkgPath)
	if err != nil {
		return BuildJob{}, fmt.Errorf("buildjob: %w", err)
	}

	if _, err = os.Stat(absPath); err != nil {
		return BuildJob{}, fmt.Errorf("buildjob: file %v not found: %w", absPath, err)
	}

	if _, err := os.Stat(*buildDir); err != nil {
		return BuildJob{}, fmt.Errorf("build dir %v not found", *buildDir)
	}

	// verify it is a buildable go command, aka see if the package is main
	return BuildJob{
		GoPkgPath:    absPath,
		Compiler:     compiler,
		BuildCommand: compiler + " build -tags netgo,purego,noasm,tinygo.enable",
		buildDir:     buildDir,
		env:          nil,
	}, nil
}

// Build result captures information about the build that later
// can be used for statistics that can be gathered after the build.
// BuildResults are always valid and do not contain any error information.
// If errors are encountered during the build process, a BuildError is returned.
type BuildResult struct {
	BuildJob   BuildJob      // copy of the original build job
	BuildTime  time.Duration // time it took to build the pkg
	BinarySize uint64        // size of the binary in bytes
}

// When the build of a command fails, a BuildError is emitted.
type BuildError struct {
	BuildJob BuildJob // copy of the original build job
	BuildErr string   // error message encountered during build
}

func (b *BuildError) String() string {
	return fmt.Sprintf("%v: %v", b.BuildJob.BuildCommand, b.BuildErr)
}

// Implement error interface over BuildError
func (b *BuildError) Error() string {
	return fmt.Sprintf("%v: %v", b.BuildJob.BuildCommand, b.BuildErr)
}

// BuildDelta provides comparative information about two BuildResults.
type BuildDelta struct {
	// size difference of the BuildResults. A negative value means,
	// that the self object is smaller.
	deltaSize int64
	deltaTime time.Duration
}

// Compare struct other with struct self
// The BuildDelta will be documented from the perspective
// of the self object. The goPkgPath of the jobs have to be the same to
// be comparable; this will be ensured.
func (b *BuildResult) Delta(o *BuildResult) (BuildDelta, error) {
	if b.BuildJob.GoPkgPath != o.BuildJob.GoPkgPath {
		return BuildDelta{}, fmt.Errorf(
			"cannot compare packages '%v' and %v",
			b.BuildJob.GoPkgPath,
			o.BuildJob.GoPkgPath,
		)
	}

	if b.BuildJob.Compiler == o.BuildJob.Compiler {
		return BuildDelta{}, fmt.Errorf(
			"cannot compare packages with same compiler '%v'",
			b.BuildJob.Compiler,
		)
	}

	return BuildDelta{
		deltaSize: int64(b.BinarySize) - int64(o.BinarySize),
		deltaTime: b.BuildTime - o.BuildTime,
	}, nil
}

type BuildState int

const (
	Setup   BuildState = iota // pre-running state. jobs can still be enqueued and config can be changed
	Running                   // the Builder is running. No new jobs can be queued
	Stopped                   // the builder was stopped without any errors
	Error                     // the builder encountered and error
)

// The Builder struct manages the entire multi-goroutine build process.
// It starts the worker routines, distributes the BuildJobs, and gathers
// the BuildResults and BuildErrors.
type Builder struct {
	jobQueue    chan BuildJob    // channel for the available build jobs
	resultQueue chan BuildResult // channel for finished BuildResults
	errQueue    chan BuildError  // channel for failed build jobs
	worker      uint             // amount of go work routines
	state       BuildState       // the current state of the builder
	jobs        []BuildJob       // list of jobs that will be added to the jobQueue
	results     []BuildResult    // list of received build results
	errors      []BuildError     // list of received build errors
	logger      *log.Logger      // optional, user provided logger
}

// Generate a new Builder struct with a configuration.
func NewBuilder(worker uint) (Builder, error) {
	if worker == 0 {
		return Builder{}, ErrNoWorkers
	}

	return Builder{
		jobQueue:    make(chan BuildJob, worker),
		resultQueue: make(chan BuildResult, 0xFF),
		errQueue:    make(chan BuildError, 0xFF),
		worker:      worker,
		state:       Setup,
		jobs:        make([]BuildJob, 0),
		results:     make([]BuildResult, 0),
		errors:      make([]BuildError, 0),
		logger:      log.Default(),
	}, nil
}

// Add a new job to the build queue, as long as the builder is till in the Setup state.
func (b *Builder) AddJob(job BuildJob) error {
	if b.state != Setup {
		return fmt.Errorf("job: builder in invalid state, cannot add new job")
	}

	if job.GoPkgPath == "" {
		return fmt.Errorf("job: invalid gopkg, cannot add '%v'", job.GoPkgPath)
	}

	b.jobs = append(b.jobs, job)
	return nil
}

// Provide a custom logger for the builder. By default, the logger will be turned off.
func (b *Builder) SetLogger(logger *log.Logger) {
	b.logger = logger
}

// Start execution of the builder. The builder will be set into the Running state
// and will block all methods modifying the Builders stats.
func (b *Builder) Run() error {
	if b.state != Setup {
		return ErrBuilderInvalidState
	}
	b.state = Running

	if len(b.jobs) == 0 {
		return ErrBuilderNoJobs
	}

	var errBuf bytes.Buffer

	for _, job := range b.jobs {
		// run the build process
		fields := strings.Fields(job.BuildCommand)
		fields = append(fields, "-o", filepath.Join(*job.buildDir, filepath.Base(job.GoPkgPath)))

		c := exec.Command(fields[0], fields[1:]...)
		c.Env = append(os.Environ(), job.env...)
		c.Dir = job.GoPkgPath
		c.Stderr = &errBuf

		t0 := time.Now()
		err := c.Run()
		t1 := time.Now()

		if err != nil {
			// the subprocess failed with a non-zero exit code, so create BuildError
			b.errors = append(b.errors, BuildError{
				BuildJob: job,
				BuildErr: errBuf.String(),
			})
		} else {
			pkgNameToken := strings.Split(job.GoPkgPath, "/")
			binPath := filepath.Join(job.GoPkgPath, pkgNameToken[len(pkgNameToken)-1])

			f, err := os.Stat(binPath)
			if err != nil {
				fmt.Printf("worker: could not find file %v\n", binPath)
				b.errors = append(b.errors, BuildError{
					BuildJob: job,
					BuildErr: err.Error(),
				})
				continue
			}

			if f == nil {
				fmt.Printf("worker: could not find file %v, handle is nil\n", binPath)
				continue
			}

			b.results = append(b.results, BuildResult{
				BuildJob:   job,
				BuildTime:  t1.Sub(t0),
				BinarySize: uint64(f.Size()),
			})
		}
	}

	b.state = Stopped
	return nil
}

// Retrieve the BuildResults from the finished builder.
// If the builder is in an invalid state, return error.
func (b *Builder) Results() ([]BuildResult, error) {
	if b.state != Stopped {
		return nil, fmt.Errorf("results: builder in invalid state")
	}
	return b.results, nil
}

// Retrieve the BuildResults from the finished builder.
// If the builder is in an invalid state, return error.
func (b *Builder) Errors() ([]BuildError, error) {
	if b.state != Stopped {
		return nil, fmt.Errorf("errors: builder in invalid state")
	}
	return b.errors, nil
}
