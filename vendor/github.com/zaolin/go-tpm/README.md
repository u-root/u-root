Go-TPM
======

Go-TPM is a Go library that communicates directly with a TPM on Linux. It
marshals and unmarshals buffers directly into and from formats specified in the
TPM spec. This code is under active development; the current version supports
Seal/Unseal, Quote, creating attestation identity keys, and taking ownership of
the TPM.

The examples directory contains some simple examples: creating an AIK, clearing
the TPM (using owner auth), and taking ownership of the TPM.
