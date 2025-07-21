// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	builder "github.com/u-root/u-root/tools/tinygo-buildstatus/pkg"
)

var (
	fs = flag.NewFlagSet("tinygo cmdlet builder", flag.ExitOnError)

	ErrFormatNotSupported = fmt.Errorf("unsupported output format")
	ErrStatusQuoViolated  = fmt.Errorf("status quo violated")

	CommitTinygo  string // commit hash of the tinygo binary
	CommitUroot   string // commit hash of the u-root binary
	VersionGolang string // version of the go compiler
)

type outputFormat string

const (
	// JSON outputFormat = "json"
	// CSV               = "csv"
	GH = "gh"
)

const (
	reportPreamble = `# Tinygo Cmdlet Report

## Overview
- Total cmdlets: %d
- Successful builds: %d
- Failed builds: %d
- Build success rate: %.2f

## Versions
- tinygo: %s
- go: %s
- u-root: %s

## Build status

### Completed
| Package | Build Time | Size |
|---------|------------|------|
`
)

// flags passed to the cli program
type flags struct {
	Verbose             bool
	Clean               bool   // clean the output directory after report generation
	CmdletOutputDirBase string // the base directory for the output of the cmdlets
	ReportOutputFile    string // file path to write the report
	TinygoPath          string // path to the tinygo binary
	OutputFormat        string // output format
	Jobs                uint   // number of jobs to run concurrently
}

type cmd struct {
	flags    flags    // cli flags
	cmdPaths []string // paths to the cmdlets
}

type report struct {
	TotalCmdlets int
	Success      int
	Failed       int
	SuccessRate  float64
	Results      []builder.Result
	Errors       []builder.Error
}

func reportCSV(r report) error {
	return nil
}

// reportGH generates a report in the github format (markdown)
// write the report to the stdout so it can be redirected to > $GITHUB_STEP_SUMMARY
func reportGH(r report) error {
	var s strings.Builder

	// preamble
	s.WriteString(fmt.Sprintf(reportPreamble, r.TotalCmdlets, r.Success, r.Failed, r.SuccessRate, CommitTinygo, VersionGolang, CommitUroot))
	for _, res := range r.Results {
		s.WriteString(fmt.Sprintf("| %s | %s | %d |\n", res.Job.GoPkgPath, res.BuildTime, res.BinarySize))
	}

	s.WriteString("\n\n### Errors\n")
	for _, err := range r.Errors {
		s.WriteString(fmt.Sprintf("### %s\n", err.Job.GoPkgPath))
		s.WriteString(fmt.Sprintf("```\n%s\n```\n", err.Err))
	}

	fmt.Println(s.String())
	return nil
}

func generateReport(format outputFormat, results []builder.Result, errors []builder.Error) error {
	r := report{
		TotalCmdlets: len(results) + len(errors),
		Success:      len(results),
		Failed:       len(errors),
		SuccessRate:  float64(len(results)) / float64(len(results)+len(errors)) * 100,
		Results:      results,
		Errors:       errors,
	}

	switch format {
	case GH:
		return reportGH(r)
	default:
		return ErrFormatNotSupported
	}
}

// verifyStatusQuo checks if the status quo is satisfied
// if the status quo is not satisfied, return an error
// if the status quo is satisfied, return nil
func verifyStatusQuo(results []builder.Result, errors []builder.Error, compare []string) error {
	unmatchedResults := make([]string, 0)
	unmatchedErrors := make([]string, 0)

	for _, res := range results {
		base := filepath.Base(res.Job.GoPkgPath)
		found := slices.Contains(compare, base)

		if !found {
			unmatchedResults = append(unmatchedResults, base)
		}
	}

	// verify that none of the errors are in the compare list
	for _, err := range errors {
		base := filepath.Base(err.Job.GoPkgPath)
		for _, c := range compare {
			if base == c {
				unmatchedErrors = append(unmatchedErrors, base)
			}
		}
	}

	if len(unmatchedResults) > 0 {
		log.Printf("These successful builds are not in the status quo. Maybe consider adding them? %v\n", unmatchedResults)
		return ErrStatusQuoViolated
	}

	if len(unmatchedErrors) > 0 {
		log.Printf("These cmdlets were expected to build but failed %v\n", unmatchedErrors)
		return ErrStatusQuoViolated
	}

	return nil
}

func (cmd *cmd) run() error {
	// setup output directory
	if cmd.flags.CmdletOutputDirBase == "" {
		dir, err := os.MkdirTemp("", "tinygo-cmdlet")
		if err != nil {
			return err
		}
		log.Printf("temporary directory %s", dir)
		cmd.flags.CmdletOutputDirBase = dir
	}

	if _, err := os.Stat(cmd.flags.TinygoPath); err != nil {
		return fmt.Errorf("tinygo binary %s: %w", cmd.flags.TinygoPath, err)
	}

	b, err := builder.NewBuilder(cmd.flags.Jobs)
	if err != nil {
		return err
	}

	// enqueue the build jobs adn ignore packages if pre-defined for target OS
	ignorePk := false
	for _, goPkg := range cmd.cmdPaths {
		if oses, ok := Ignore[filepath.Base(goPkg)]; ok {
			for _, os := range oses {
				if runtime.GOOS == os {
					log.Printf("ignoring '%s' for '%s'", filepath.Base(goPkg), os)
					ignorePk = true
					continue
				}
			}
		}
		if ignorePk {
			ignorePk = false
			continue
		}

		j, err := builder.NewJob(goPkg, cmd.flags.TinygoPath, &cmd.flags.CmdletOutputDirBase)
		if err != nil {
			return err
		}
		if err := b.AddJob(j); err != nil {
			return err
		}
	}

	if err := b.Run(); err != nil {
		return err
	}

	buildErrors, err := b.Errors()
	if err != nil {
		return err
	}

	buildResults, err := b.Results()
	if err != nil {
		return err
	}

	if err := generateReport((outputFormat(cmd.flags.OutputFormat)), buildResults, buildErrors); err != nil {
		return err
	}

	// clean up the output directory
	if cmd.flags.Clean {
		if err := os.RemoveAll(cmd.flags.CmdletOutputDirBase); err != nil {
			return err
		}
	}

	return verifyStatusQuo(buildResults, buildErrors, StatusQuo)
}

func parseFlags(args []string) (cmd, error) {
	var c cmd

	fs.StringVar(&c.flags.CmdletOutputDirBase, "cmdout", "", "base directory for the output of the cmdlets")
	fs.StringVar(&c.flags.ReportOutputFile, "report", "", "file path to write the report to")
	fs.StringVar(&c.flags.TinygoPath, "tinygo", "", "path to the tinygo binary")
	fs.StringVar(&c.flags.OutputFormat, "format", GH, "output format")
	fs.BoolVar(&c.flags.Verbose, "verbose", false, "verbose output")
	fs.BoolVar(&c.flags.Verbose, "clean", false, "verbose output")
	fs.UintVar(&c.flags.Jobs, "jobs", 1, "number of jobs to run concurrently")

	// these three flags must be supplied via the CI invocation, they will be reflected in the report document
	fs.StringVar(&CommitTinygo, "commit-tinygo", "", "commit hash of the tinygo binary")
	fs.StringVar(&CommitUroot, "commit-uroot", "", "commit hash of the u-root binary")
	fs.StringVar(&VersionGolang, "version-go", "", "version of the go compiler")

	fs.Parse(args)

	if len(args) == 0 {
		fs.Usage()
		os.Exit(2)
	}

	// get the command paths from the arguments
	// match globs such as cmds/core/* or cmds/exp*
	for _, d := range args {
		dirs, err := filepath.Glob(d)
		if err != nil {
			return c, err
		}

		// make sure the paths are directories
		for _, dir := range dirs {
			fi, err := os.Stat(dir)
			if err != nil {
				return c, err
			}
			if fi.IsDir() {
				p, err := filepath.Abs(dir)
				if err != nil {
					return c, err
				}
				c.cmdPaths = append(c.cmdPaths, p)
			}
		}
	}

	return c, nil
}

func main() {
	cmd, err := parseFlags(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	if err = cmd.run(); err != nil {
		log.Fatal(err)
	}
}
