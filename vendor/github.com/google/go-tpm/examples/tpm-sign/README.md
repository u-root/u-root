# tpm-sign

This example shows how you can generate keys inside the TPM and use them for signature/verification operations. This utility supports `sign`, `verify`, `generate`, and `extendPcr` actions. Use `./tpm-sign <action> --help` for advanced usage of each action.

## Basic Usage
The following snippet shows how you can generate a key, sign data with it, and verify the signature.

```
$ ./tpm-sign generate
Writing keyblob to keyblob
Writing public key to publickey
$ echo test_data | ./tpm-sign sign
Writing signature to sig.data
$ echo test_data | ./tpm-sign verify
Signature valid.
```

## Binding against PCRs
This example shows how you can generate a key that is bound against PCR values.

```
$ ./tpm-sign extendPcr --reset --pcr 16
$ ./tpm-sign generate --pcrs 0,16
Writing keyblob to keyblob
Writing public key to publickey
$ echo test_data | ./tpm-sign sign
Writing signature to sig.data
$ echo test_measurement | ./tpm-sign extendPcr --pcr 16
$ echo test_data | ./tpm-sign sign
Could not perform sign operation: tpm: the named PCR value does not match the current PCR value
```

