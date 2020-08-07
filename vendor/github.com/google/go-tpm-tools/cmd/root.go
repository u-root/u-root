// Package cmd contains a CLI to interact with TPM.
package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd is the entrypoint for gotpm.
var RootCmd = &cobra.Command{
	Use: "gotpm",
	Long: `Command line tool for the go-tpm TSS

This tool allows performing TPM2 operations from the command line.
See the per-command documentation for more information.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if quiet && verbose {
			return fmt.Errorf("cannot specify both --quiet and --verbose")
		}
		cmd.SilenceUsage = true
		return nil
	},
}

var (
	quiet   bool
	verbose bool
)

func init() {
	RootCmd.PersistentFlags().BoolVar(&quiet, "quiet", false,
		"print nothing if command is successful")
	RootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false,
		"print additional info to stdout")
	hideHelp(RootCmd)
}

func messageOutput() io.Writer {
	if quiet {
		return ioutil.Discard
	}
	return os.Stdout
}

func debugOutput() io.Writer {
	if verbose {
		return os.Stdout
	}
	return ioutil.Discard
}
