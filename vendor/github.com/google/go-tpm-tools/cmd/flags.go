package cmd

import (
	"errors"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-tpm-tools/tpm2tools"
	"github.com/google/go-tpm/tpm2"
	"github.com/spf13/cobra"
)

var (
	output   string
	input    string
	nvIndex  uint32
	keyAlgo  = tpm2.AlgRSA
	pcrs     []int
	hashAlgo = tpm2.AlgSHA256
)

type pcrsFlag struct {
	value *[]int
}

func (f *pcrsFlag) Set(val string) error {
	for _, d := range strings.Split(val, ",") {
		pcr, err := strconv.Atoi(d)
		if err != nil {
			return err
		}
		if pcr < 0 || pcr >= tpm2tools.NumPCRs {
			return errors.New("pcr out of range")
		}
		*f.value = append(*f.value, pcr)
	}
	return nil
}

func (f *pcrsFlag) Type() string {
	return "pcrs"
}

func (f *pcrsFlag) String() string {
	return "" // Don't display a default value
}

var algos = map[tpm2.Algorithm]string{
	tpm2.AlgRSA:    "rsa",
	tpm2.AlgECC:    "ecc",
	tpm2.AlgSHA1:   "sha1",
	tpm2.AlgSHA256: "sha256",
	tpm2.AlgSHA384: "sha384",
	tpm2.AlgSHA512: "sha512",
}

type algoFlag struct {
	value   *tpm2.Algorithm
	allowed []tpm2.Algorithm
}

func (f *algoFlag) Set(val string) error {
	present := false
	for _, algo := range f.allowed {
		if algos[algo] == val {
			*f.value = algo
			present = true
		}
	}
	if !present {
		return errors.New("unknown algorithm")
	}
	return nil
}

func (f *algoFlag) Type() string {
	return "algo"
}

func (f *algoFlag) String() string {
	return algos[*f.value]
}

// Allowed gives a string list of the permitted algorithm values for this flag.
func (f *algoFlag) Allowed() string {
	out := make([]string, len(f.allowed))
	for i, a := range f.allowed {
		out[i] = algos[a]
	}
	return strings.Join(out, ", ")
}

// Disable the "help" subcommand (and just use the -h/--help flags).
// This should be called on all commands with subcommands.
// See https://github.com/spf13/cobra/issues/587 for why this is needed.
func hideHelp(cmd *cobra.Command) {
	cmd.SetHelpCommand(&cobra.Command{Hidden: true})
}

// Lets this command specify an output file, for use with dataOutput().
func addOutputFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&output, "output", "",
		"output file (defaults to stdout)")
}

// Lets this command specify an input file, for use with dataInput().
func addInputFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&input, "input", "",
		"input file (defaults to stdin)")
}

// Lets this command specify an NVDATA index, for use with nvIndex.
func addIndexFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().Uint32Var(&nvIndex, "index", 0,
		"NVDATA index, cannot be 0")
}

// Lets this command specify some number of PCR arguments, check if in range.
func addPCRsFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().Var(&pcrsFlag{&pcrs}, "pcrs", "comma separated list of PCR numbers")
}

// Lets this command specify the public key algorithm.
func addPublicKeyAlgoFlag(cmd *cobra.Command) {
	f := algoFlag{&keyAlgo, []tpm2.Algorithm{tpm2.AlgRSA, tpm2.AlgECC}}
	cmd.PersistentFlags().Var(&f, "algo", "public key algorithm: "+f.Allowed())
}

func addHashAlgoFlag(cmd *cobra.Command) {
	f := algoFlag{&hashAlgo, []tpm2.Algorithm{tpm2.AlgSHA1, tpm2.AlgSHA256, tpm2.AlgSHA384, tpm2.AlgSHA512}}
	cmd.PersistentFlags().Var(&f, "hash-algo", "hash algorithm: "+f.Allowed())
}

// alwaysError implements io.ReadWriter by always returning an error
type alwaysError struct {
	error
}

func (ae alwaysError) Write([]byte) (int, error) {
	return 0, ae.error
}

func (ae alwaysError) Read(p []byte) (n int, err error) {
	return 0, ae.error
}

// Handle to output data file. If there is an issue opening the file, the Writer
// returned will return the error upon any call to Write()
func dataOutput() io.Writer {
	if output == "" {
		return os.Stdout
	}

	file, err := os.Create(output)
	if err != nil {
		return alwaysError{err}
	}
	return file
}

// Handle to input data file. If there is an issue opening the file, the Reader
// returned will return the error upon any call to Read()
func dataInput() io.Reader {
	if input == "" {
		return os.Stdin
	}

	file, err := os.Open(input)
	if err != nil {
		return alwaysError{err}
	}
	return file
}

func getSelection() tpm2.PCRSelection {
	return tpm2.PCRSelection{Hash: hashAlgo, PCRs: pcrs}
}

// Load SRK based on tpm2.Algorithm set in the global flag vars.
func getSRK(rwc io.ReadWriter) (*tpm2tools.Key, error) {
	switch keyAlgo {
	case tpm2.AlgRSA:
		return tpm2tools.StorageRootKeyRSA(rwc)
	case tpm2.AlgECC:
		return tpm2tools.StorageRootKeyECC(rwc)
	default:
		panic("unexpected keyAlgo")
	}
}

// Load EK based on tpm2.Algorithm set in the global flag vars.
func getEK(rwc io.ReadWriter) (*tpm2tools.Key, error) {
	switch keyAlgo {
	case tpm2.AlgRSA:
		return tpm2tools.EndorsementKeyRSA(rwc)
	case tpm2.AlgECC:
		return tpm2tools.EndorsementKeyECC(rwc)
	default:
		panic("unexpected keyAlgo")
	}
}
