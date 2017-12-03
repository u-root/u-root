package gzip

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/klauspost/pgzip"
	"github.com/spf13/pflag"
)

type Options struct {
	Blocksize  int
	Level      int
	Processes  int
	Decompress bool
	Force      bool
	Help       bool
	Keep       bool
	License    bool
	Quiet      bool
	Stdin      bool
	Stdout     bool
	Test       bool
	Verbose    bool
	Version    bool
	Suffix     string
}

func (o *Options) ParseArgs() error {
	var levels [10]bool

	pflag.IntVarP(&o.Blocksize, "blocksize", "b", 128, "Set compression block size in KiB")
	pflag.BoolVarP(&o.Decompress, "decompress", "d", false, "Decompress the compressed input")
	pflag.BoolVarP(&o.Force, "force", "f", false, "Force overwrite of output file and compress links")
	pflag.BoolVarP(&o.Help, "help", "h", false, "Display a help screen and quit")
	pflag.BoolVarP(&o.Keep, "keep", "k", false, "Do not delete original file after processing")
	// TODO: implement list option here
	pflag.IntVarP(&o.Processes, "processes", "p", runtime.NumCPU(), "Allow up to n compression threads")
	pflag.BoolVarP(&o.Quiet, "quiet", "q", false, "Print no messages, even on error")
	// TODO: implement recursive option here
	pflag.BoolVarP(&o.Stdout, "stdout", "c", false, "Write all processed output to stdout (won't delete)")
	pflag.StringVarP(&o.Suffix, "suffix", "S", ".gz", "Specify suffix for compression")
	pflag.BoolVarP(&o.Test, "test", "t", false, "Test the integrity of the compressed input")
	pflag.BoolVarP(&o.Verbose, "verbose", "v", false, "Produce more verbose output")
	pflag.BoolVarP(&levels[1], "fast", "1", false, "Compression Level 1")
	pflag.BoolVarP(&levels[2], "two", "2", false, "Compression Level 2")
	pflag.BoolVarP(&levels[3], "three", "3", false, "Compression Level 3")
	pflag.BoolVarP(&levels[4], "four", "4", false, "Compression Level 4")
	pflag.BoolVarP(&levels[5], "five", "5", false, "Compression Level 5")
	pflag.BoolVarP(&levels[6], "six", "6", false, "Compression Level 6")
	pflag.BoolVarP(&levels[7], "seven", "7", false, "Compression Level 7")
	pflag.BoolVarP(&levels[8], "eight", "8", false, "Compression Level 8")
	pflag.BoolVarP(&levels[9], "best", "9", false, "Compression Level 9")

	// Hide the compression Levels 2 - 8 from usage.
	_ = pflag.CommandLine.MarkHidden("two")
	_ = pflag.CommandLine.MarkHidden("three")
	_ = pflag.CommandLine.MarkHidden("four")
	_ = pflag.CommandLine.MarkHidden("five")
	_ = pflag.CommandLine.MarkHidden("six")
	_ = pflag.CommandLine.MarkHidden("seven")
	_ = pflag.CommandLine.MarkHidden("eight")

	pflag.Parse()

	var err error
	o.Level, err = parseLevels(levels)
	if err != nil {
		return err
	}

	return o.validate()
}

// Validate checks options and handles help, version, and license
// output to the user. If forces decompression to be enabled when
// test mode is enabled. It further modifies options if the running
// binary is named gunzip or gzcat to allow for expected behavor for
// those binaries. Finally it checks if there is piped stdin data.
func (o *Options) validate() error {

	if o.Help {
		return &appError{
			level: info,
			msg:   fmt.Sprintf("Usage of %s:\n%s", filepath.Base(os.Args[0]), pflag.CommandLine.FlagUsages()),
		}
	}

	if o.Test {
		o.Decompress = true
	}

	// Support gunzip and gzcat symlinks
	if filepath.Base(os.Args[0]) == "gunzip" {
		o.Decompress = true
	} else if filepath.Base(os.Args[0]) == "gzcat" {
		o.Decompress = true
		o.Stdout = true
	}

	// Stat os.Stdin and ignore errors. stat will be a nil FileInfo if there is an
	// error.
	stat, _ := os.Stdin.Stat()

	// No files passed and arguments and Stdin piped data found.
	// Stdin piped data is ignored if arguments are found.
	if len(pflag.Args()) == 0 && (stat.Mode()&os.ModeNamedPipe) != 0 {
		o.Stdin = true
		// Enable force to ignore suffix checks
		o.Force = true
		// Since there's no filename to derive the output path from, only support
		// outputting to stdout when data is piped from stdin
		o.Stdout = true
	} else if len(pflag.Args()) == 0 {
		// No stdin piped data found and no files passed as arguments
		return &appError{
			level: info,
			msg:   fmt.Sprintf("Usage of %s:\n%s", filepath.Base(os.Args[0]), pflag.CommandLine.FlagUsages()),
		}
	}

	return nil
}

// parseLevels loops through a [10]bool and returns the index of the element
// thats true. If more than one element is true return an error. If no
// element is true, return the constant pgzip.DefaultCompression (-1).
func parseLevels(Levels [10]bool) (int, error) {
	var Level int

	for i, l := range Levels {
		if l && Level != 0 {
			return 0, &appError{level: fatal, msg: "multiple compression Levels specified"}
		} else if l {
			Level = i
		}
	}

	if Level == 0 {
		Level = pgzip.DefaultCompression
	}
	return Level, nil
}
