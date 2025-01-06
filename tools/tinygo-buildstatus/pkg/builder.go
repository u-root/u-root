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

// Job defines a job that can be scheduled by the Builder and
// executed by the work function. It has all necessary information
// to execute the build and gather error logs.
type Job struct {
	GoPkgPath string   // the path to the go package that is build
	Compiler  string   // compiler used to build the pkg, probably go or tinygo. TODO: refactor to enum
	Command   string   // complete build command issues by builder
	buildDir  *string  // directory where the build binary should be put in the context of go/tinygo this means -o <buildDir>
	env       []string // environment that the execution worker should use additionally
}

// Generate a new BuildJob structure and verify all provided data.
// If the path to the goPkg is invalid, it will return an error.
// TODO: should we enforce that?
func NewJob(goPkgPath string, compiler string, buildDir *string) (Job, error) {
	absPath, err := filepath.Abs(goPkgPath)
	if err != nil {
		return Job{}, fmt.Errorf("job: %w", err)
	}

	if _, err = os.Stat(absPath); err != nil {
		return Job{}, fmt.Errorf("buildjob: file %v not found: %w", absPath, err)
	}

	if _, err := os.Stat(*buildDir); err != nil {
		return Job{}, fmt.Errorf("build dir %v not found", *buildDir)
	}

	// verify it is a buildable go command, aka see if the package is main
	return Job{
		GoPkgPath: absPath,
		Compiler:  compiler,
		Command:   compiler + " build -tags netgo,purego,noasm,tinygo.enable",
		buildDir:  buildDir,
		env:       nil,
	}, nil
}

// Build result captures information about the build that later
// can be used for statistics that can be gathered after the build.
// BuildResults are always valid and do not contain any error information.
// If errors are encountered during the build process, a BuildError is returned.
type Result struct {
	Job        Job           // copy of the original build job
	BuildTime  time.Duration // time it took to build the pkg
	BinarySize uint64        // size of the binary in bytes
}

// When the build of a command fails, a Error is emitted.
type Error struct {
	Job Job    // copy of the original build job
	Err string // error message encountered during build
}

func (b *Error) String() string {
	return fmt.Sprintf("%v: %v", b.Job.Command, b.Err)
}

// Implement error interface over BuildError
func (b *Error) Error() string {
	return fmt.Sprintf("%v: %v", b.Job.Command, b.Err)
}

// delta provides comparative information about two BuildResults.
type delta struct {
	// size difference of the BuildResults. A negative value means,
	// that the self object is smaller.
	deltaSize int64
	deltaTime time.Duration
}

// Compare struct other with struct self
// The delta will be documented from the perspective
// of the self object. The goPkgPath of the jobs have to be the same to
// be comparable; this will be ensured.
func (b *Result) delta(o *Result) (delta, error) {
	if b.Job.GoPkgPath != o.Job.GoPkgPath {
		return delta{}, fmt.Errorf(
			"cannot compare packages '%v' and %v",
			b.Job.GoPkgPath,
			o.Job.GoPkgPath,
		)
	}

	if b.Job.Compiler == o.Job.Compiler {
		return delta{}, fmt.Errorf(
			"cannot compare packages with same compiler '%v'",
			b.Job.Compiler,
		)
	}

	return delta{
		deltaSize: int64(b.BinarySize) - int64(o.BinarySize),
		deltaTime: b.BuildTime - o.BuildTime,
	}, nil
}

type BuildState int

const (
	Setup   BuildState = iota // pre-running state. jobs can still be enqueued and config can be changed
	Running                   // the Builder is running. No new jobs can be queued
	Stopped                   // the builder was stopped without any errors
	Err                       // the builder encountered and error
)

// The Builder struct manages the entire multi-goroutine build process.
// It starts the worker routines, distributes the BuildJobs, and gathers
// the BuildResults and BuildErrors.
type Builder struct {
	jobQueue    chan Job    // channel for the available build jobs
	resultQueue chan Result // channel for finished BuildResults
	errQueue    chan Error  // channel for failed build jobs
	worker      uint        // amount of go work routines
	state       BuildState  // the current state of the builder
	jobs        []Job       // list of jobs that will be added to the jobQueue
	results     []Result    // list of received build results
	errors      []Error     // list of received build errors
	logger      *log.Logger // optional, user provided logger
}

// Generate a new Builder struct with a configuration.
func NewBuilder(worker uint) (Builder, error) {
	if worker == 0 {
		return Builder{}, ErrNoWorkers
	}

	return Builder{
		jobQueue:    make(chan Job, worker),
		resultQueue: make(chan Result, 0xFF),
		errQueue:    make(chan Error, 0xFF),
		worker:      worker,
		state:       Setup,
		jobs:        make([]Job, 0),
		results:     make([]Result, 0),
		errors:      make([]Error, 0),
		logger:      log.Default(),
	}, nil
}

// Add a new job to the build queue, as long as the builder is till in the Setup state.
func (b *Builder) AddJob(job Job) error {
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
		fields := strings.Fields(job.Command)
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
			b.errors = append(b.errors, Error{
				Job: job,
				Err: errBuf.String(),
			})
		} else {
			pkgNameToken := strings.Split(job.GoPkgPath, "/")
			binPath := filepath.Join(*job.buildDir, pkgNameToken[len(pkgNameToken)-1])

			f, err := os.Stat(binPath)
			if err != nil {
				fmt.Printf("worker: could not find file %v\n", binPath)
				b.errors = append(b.errors, Error{
					Job: job,
					Err: err.Error(),
				})
				continue
			}

			if f == nil {
				fmt.Printf("worker: could not find file %v, handle is nil\n", binPath)
				continue
			}

			b.results = append(b.results, Result{
				Job:        job,
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
func (b *Builder) Results() ([]Result, error) {
	if b.state != Stopped {
		return nil, fmt.Errorf("results: builder in invalid state")
	}
	return b.results, nil
}

// Retrieve the BuildResults from the finished builder.
// If the builder is in an invalid state, return error.
func (b *Builder) Errors() ([]Error, error) {
	if b.state != Stopped {
		return nil, fmt.Errorf("errors: builder in invalid state")
	}
	return b.errors, nil
}
