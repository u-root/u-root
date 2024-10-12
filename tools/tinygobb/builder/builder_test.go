// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"path/filepath"
	"testing"
	"time"
)

func TestBuildJob(t *testing.T) {
	dir := t.TempDir()
	tinygobbPathRelative := "../../../tools/tinygobb"
	buildJob, err := NewBuildJob(tinygobbPathRelative, "go", &dir)
	if err != nil {
		t.Fatalf("want nil, got %v\n", err)
	}

	abs, err := filepath.Abs(tinygobbPathRelative)
	if err != nil {
		t.Fatalf("want nil, got %v\n", err)
	}

	if abs != buildJob.goPkgPath {
		t.Fatalf("%v != %v\n", abs, buildJob.goPkgPath)
	}

	// invalid build dir
	invalidPath := "/invalid"
	buildJob, err = NewBuildJob(tinygobbPathRelative, "go", &invalidPath)
	if err == nil {
		t.Fatalf("want err, got nil\n")
	}

}

func TestBuildJobShouldFail(t *testing.T) {
	dir := t.TempDir()
	tinygobbPathRelative := "invalid"
	abs, err := NewBuildJob(tinygobbPathRelative, "go", &dir)
	if err == nil {
		t.Fatalf("want err, got %v (%v)\n", err, abs)
	}
}

func TestDelta(t *testing.T) {
	buildJobGo := BuildJob{
		goPkgPath:    "test",
		compiler:     "go",
		buildCommand: "go build",
	}

	buildJobTinyGo := BuildJob{
		goPkgPath:    "test",
		compiler:     "tinygo",
		buildCommand: "tinygo build",
	}

	b0 := BuildResult{
		buildJob:   buildJobGo,
		buildTime:  time.Duration(time.Second * 2),
		binarySize: 1337,
	}

	b1 := BuildResult{
		buildJob:   buildJobTinyGo,
		buildTime:  time.Duration(time.Second * 2),
		binarySize: 123,
	}

	b2 := b1
	b2.buildJob.goPkgPath = "other"

	delta, err := b0.Delta(&b1)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	wantDeltaSize := int64(b0.binarySize) - int64(b1.binarySize)
	if delta.deltaSize != wantDeltaSize {
		t.Fatalf("want delta size == %v, got %v", wantDeltaSize, delta.deltaSize)
	}

	wantDeltaTime := b0.buildTime - b1.buildTime
	if delta.deltaTime != wantDeltaTime {
		t.Fatalf("want delta size == %v, got %v", wantDeltaSize, delta.deltaTime)
	}

	// delta with itself
	_, err = b0.Delta(&b0)
	if err == nil {
		t.Fatalf("got nil, want err")
	}

	// delta with other pkg
	_, err = b0.Delta(&b2)
	if err == nil {
		t.Fatalf("got nil, want err")
	}
}

func TestBuilder(t *testing.T) {
	worker := 1
	dir := t.TempDir()
	b, err := NewBuilder(uint(worker))
	if err != nil {
		t.Fatalf("wanted nil, got %v", err)
	}

	// 1 worker, 1 job (kind of stupid, but should work)
	b0, err := NewBuildJob("../../../tools/tinygobb", "go", &dir)
	if err != nil {
		t.Fatalf("want nil, go %v", err)
	}

	if err := b.AddJob(b0); err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	if len(b.jobs) != 1 {
		t.Fatalf("|jobQueue| != 1")
	}

	err = b.Run()
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	// 2 worker. 2 jobs

	// 2 worker, 10 jobs
}

func TestBuilderShouldFail(t *testing.T) {
	if _, err := NewBuilder(0); err != ErrNoWorkers {
		t.Fatalf("want %v, got %v", ErrNoWorkers, err)
	}

	worker := 1
	b, err := NewBuilder(uint(worker))
	if err != nil {
		t.Fatalf("wanted nil, got %v", err)
	}

	err = b.Run()
	if err != ErrBuilderNoJobs {
		t.Fatalf("want %v, got %v", ErrBuilderNoJobs, err)
	}

	// state is running or stopped, should not add job
	err = b.AddJob(BuildJob{})
	if err == ErrBuilderInvalidState {
		t.Logf("want err, got %v", err)
	}
}
