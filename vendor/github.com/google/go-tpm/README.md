Go-TPM
======

Go-TPM is a Go library that communicates directly with a TPM device on Linux or
Windows machines.

The libraries don't implement the entire spec for neither 1.2 nor 2.0. **If you
need a command that's missing, contributions are welcome!**

Please note that this is not an official Google product.

## Structure

The `tpm` directory contains TPM 1.2 client library. This library is in
["maintenance mode"](#tpm-1.2).

The `tpm2` directory contains TPM 2.0 client library.

The `examples` directory contains some simple examples for both versions of the
spec.

## TPM 1.2

TPM 1.2 support currently has no maintainer. None of the TPM 2.0 maintainers
have expertise on 1.2 either.

As such, TPM 1.2 library is in "maintenance" mode - all PRs with new
functionality or non-critical fixes will be rejected.

**If you'd like to volunteer to maintain the TPM 1.2 library, you can do so via
an [issue](https://github.com/google/go-tpm/issues).** You don't have to be a
Googler to volunteer.
