package cmd

import (
	"fmt"
	"io/ioutil"

	pb "github.com/golang/protobuf/proto"
	"github.com/spf13/cobra"

	tpmpb "github.com/google/go-tpm-tools/proto"
	"github.com/google/go-tpm-tools/tpm2tools"
	"github.com/google/go-tpm/tpm2"
)

var sealCmd = &cobra.Command{
	Use:   "seal",
	Short: "Seal some data to the TPM",
	Long: `Encrypt the input data using the TPM

TPMs support a "sealing" operation that allows some secret data to be encrypted
by a particular TPM. This data can only be decrypted by the same TPM that did
the encryption.

Optionally (using the --pcrs flag), this decryption can be furthur restricted to
only work if certain Platform Control Registers (PCRs) are in the correct state.
This allows a key (i.e. a disk encryption key) to be bound to specific machine
state (like Secure Boot).`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		rwc, err := openTpm()
		if err != nil {
			return err
		}
		defer rwc.Close()

		fmt.Fprintln(debugOutput(), "Loading SRK")
		srk, err := getSRK(rwc)
		if err != nil {
			return err
		}
		defer srk.Close()

		fmt.Fprintln(debugOutput(), "Reading sealed data")
		secret, err := ioutil.ReadAll(dataInput())
		if err != nil {
			return err
		}

		sel := getSelection()
		fmt.Fprintf(debugOutput(), "Sealing to PCRs: %v\n", sel.PCRs)
		var sOpt tpm2tools.SealOpt
		if len(sel.PCRs) > 0 {
			sOpt = tpm2tools.SealCurrent{PCRSelection: sel}
		}
		sealed, err := srk.Seal(secret, sOpt)
		if err != nil {
			return fmt.Errorf("sealing data: %v", err)
		}

		fmt.Fprintln(debugOutput(), "Writing sealed data")
		if err := pb.MarshalText(dataOutput(), sealed); err != nil {
			return err
		}
		fmt.Fprintf(debugOutput(), "Sealed data to PCRs: %v\n", sel.PCRs)
		return nil
	},
}

var unsealCmd = &cobra.Command{
	Use:   "unseal",
	Short: "Unseal some data previously sealed to the TPM",
	Long: `Decrypt the input data using the TPM

The opposite of "gotpm seal". This takes in some sealed input and decrypts it
using the TPM. This operation will fail if used on a different TPM, or if the
Platform Control Registers (PCRs) are in the incorrect state.

All the necessary data to decrypt the sealed input is present in the input blob.
We do not need to specify the PCRs used for unsealing.

We do support an optional "certification" process. A list of PCRs may be
provided with --pcrs, and the unwrapping will fail if the PCR values when
sealing differ from the current PCR values. This allows for verification of the
machine state when sealing took place.
`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		rwc, err := openTpm()
		if err != nil {
			return err
		}
		defer rwc.Close()

		fmt.Fprintln(debugOutput(), "Reading sealed data")
		data, err := ioutil.ReadAll(dataInput())
		if err != nil {
			return err
		}
		var sealed tpmpb.SealedBytes
		if err := pb.UnmarshalText(string(data), &sealed); err != nil {
			return err
		}

		fmt.Fprintln(debugOutput(), "Loading SRK")
		keyAlgo = tpm2.Algorithm(sealed.GetSrk())
		srk, err := getSRK(rwc)
		if err != nil {
			return err
		}
		defer srk.Close()

		fmt.Fprintln(debugOutput(), "Unsealing data")

		certifySel := tpm2.PCRSelection{Hash: tpm2tools.CertifyHashAlgTpm, PCRs: pcrs}
		var cOpt tpm2tools.CertifyOpt
		if len(certifySel.PCRs) > 0 {
			cOpt = tpm2tools.CertifyCurrent{PCRSelection: certifySel}
		}
		secret, err := srk.Unseal(&sealed, cOpt)
		if err != nil {
			return fmt.Errorf("unsealing data: %v", err)
		}

		fmt.Fprintln(debugOutput(), "Writing secret data")
		if _, err := dataOutput().Write(secret); err != nil {
			return fmt.Errorf("writing secret data: %v", err)
		}
		fmt.Fprintln(debugOutput(), "Unsealed data using TPM")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(sealCmd)
	RootCmd.AddCommand(unsealCmd)
	addInputFlag(sealCmd)
	addInputFlag(unsealCmd)
	addOutputFlag(sealCmd)
	addOutputFlag(unsealCmd)
	// PCRs and hash algorithm only used for sealing
	addPCRsFlag(sealCmd)
	addHashAlgoFlag(sealCmd)
	addPCRsFlag(unsealCmd)
	addPublicKeyAlgoFlag(sealCmd)
}
