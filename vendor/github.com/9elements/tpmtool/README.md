tpmtool is a tool for TPM interaction and disk encryption. It is written in pure Go.

# Basic Features

-   Supports TPM 1.2 and 2.0 with [Go TSS](https://github.com/google/go-tpm).
-   Higher TPM abstraction layer (TSPI) is implemented.
-   Written in pure Go.
-   TPM states are derived by Linux sysfs.
-   Automatic TSS selection based on TPM version.
-   TPM1 & TPM2 event log parser
-   **Currently only TSPI for TPM specification 1.2 is available.**

## Core Features

-   Shows the TPM status.

```bash
TPM Manufacturer:          STMicroelectronics
TPM spec:                  1.2
TPM owned:                 true
TPM activated:             true
TPM enabled:               true
TPM temporary deactivated: false
```

-   Dumps Endorsement Key into a file and shows the fingerprint.
-   Takes ownership of the TPM.
-   Clears ownership of the TPM.
-   Resets TPM lock in case of active bruteforce detection.
-   Sealing/Unsealing credentials with custom/current set of PCRs.
-   Resealing of credentials using a sealing configuration for PCR pre-calculation
-   List and read PCRs
-   Measures a file into given PCR index.
-   Dump TPM eventlog from OS or custom eventlog binary file input.
-   Cryptsetup:
    -   Format device and seal credential.
    -   Open device by sealed credential.
    -   Close device.
    -   Measure device luks header into a given PCR.

# Package Availability

[![Packaging status](https://repology.org/badge/vertical-allrepos/tpmtool.svg)](https://repology.org/metapackage/tpmtool)

# Dependencies

-   [cryptsetup](https://gitlab.com/cryptsetup/cryptsetup) binary is required for the disk commands.

# PCR pre-calculation

PCR pre-calculation is an important feature to reseal credential in case of PCR changes e.g. kernel/firmware update.

Usage:

```bash
tpmtool crypt reseal sealing.yml sealed-key.file
```

Example sealing configuration:

```yaml
---
pcr0:
  - method: measure
    filepaths:
      - /boot/kernel
      - /boot/initramfs
  - method: extend
    hashes:
      - 8dad1c80be028384f26b929b7e7e251fbe3c1d5
pcr1:
  - method: dynamic
pcr2:
  - method: static
    hash: c3018af653e2f1a16118dd8bab2f409fbc82aa9f
pcr3:
  - method: log
    firmware: UEFI
```

## Calculation methods

Every PCR can contain different calculation methods. The static and dynamic method are standalone and can't be used with other methods.

### Static

Overwrites and sets the PCR hash you define.

**method:** static

**hash:** 8dad1c80be028384f26b929b7e7e251fbe3c1d5 (string type)

### Dynamic

Gets the current PCR of the TPM. Overwrites and sets the hash.

**method:** dynamic

### Extend

Extends a hash into the current PCR.

**method:** extend

**hashes:** [ 8dad1c80be028384f26b929b7e7e251fbe3c1d5, c3018af653e2f1a16118dd8bab2f409fbc82aa9f ] \(array type)

### FimwareLog

Uses the existing firmware log for PCR pre-calculation.

**method:** log

**firmware:** BIOS (enum type){UEFI, BIOS}

### Measure

Measures a file into the current PCR.

**method:** measure

**filepaths:** [ /foo/bash, /test/foo ] \(array type)

### Luks

Measures a LUKS header of a device into the current PCR.

**method:** luks

**devicepath:** /dev/sda (string type)
