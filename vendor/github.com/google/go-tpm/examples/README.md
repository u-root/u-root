# go-tpm usage examples

This directory contains binaries that show how to use go-tpm.

## Versions

Directories that start with `tpm-` are for TPM 1.x devices.

Directories that start with `tpm2-` are for TPM 2.x devices.

They are not compatible. For example, running `tpm-sign` against a TPM 2.x
device will fail.
