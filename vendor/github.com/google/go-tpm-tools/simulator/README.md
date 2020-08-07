# Go bindings to the Microsoft TPM2 Simulator

Microsoft maintains the reference implementation of the TPM2 spec at:
https://github.com/Microsoft/ms-tpm-20-ref/.

The Microsoft code used here is a actually
[a fork of the upstream source](https://github.com/josephlr/ms-tpm-20-ref/tree/google).
It is vendored at `simulator/ms-tpm-20-ref` to maintain compatiblity with
`go get`. Building the simulator requires the OpenSSL headers to be installed.
This can be doen with:
  - Debain based systems (including Ubuntu): `apt install libssl-dev`
  - Red Hat based systems: `yum install openssl-devel`
  - Arch Linux based systems: [`openssl`](https://www.archlinux.org/packages/core/x86_64/openssl/)
    is installed by default (as a dependancy of `base`) and includes the headers.

## Debugging

The simulator provides a useful way to figure out what the TPM is actually doing
when it executes a command. If you compile a test which runs against the
simulator, you can step through the simulator C source to see the exact
operations performed.

To do this:
1. Compile a test as a standalone binary. For example, if you were using a
  `go-tpm-tools/tpm2tools` test (which all run against the simulator), compile
  the test binary named `tpm2tools.test` by running:
    ```bash
    go test -c github.com/google/go-tpm-tools/tpm2tools
    ```
1. Now you can debug the binary using GDB:
    ```bash
    # Load the binary into GDB (fixing any errors/warnings you get)
    gdb ./tpm2tools.test
    # In GDB, set a breakpoint in the funciton you want to use.
    (gdb) break TPM2_CreatePrimary 
    Breakpoint 1 at 0x5d3710: file ./TPMCmd/tpm/src/command/Hierarchy/CreatePrimary.c, line 72.
    # Now you can either run all the tests in the package, or just one.
    # As we want to depug TPM2_CreatePrimary we'll run TestSeal
    (gdb) run -test.run TestSeal
    Starting program: ./tpm2tools.test -test.run TestSeal
    Thread 1 "tpm2tools.test" hit Breakpoint 1, TPM2_CreatePrimary
        at ./TPMCmd/tpm/src/command/Hierarchy/CreatePrimary.c:72
    72	{
    # Go to the next line
    (gdb) n
    81	    newObject = FindEmptyObjectSlot(&out->objectHandle);
    # Step into a function
    (gdb) s
    FindEmptyObjectSlot
        at ./TPMCmd/tpm/src/subsystem/Object.c:266
    266	{
    # Continue until the next breakpoint (or exiting)
    (gdb) c
    Continuing.
    PASS
    [Inferior 1 (process 29395) exited normally]
    ```

## IDE Support

When examining the TPM2 C code, is is often useful to have IDE support for
things like "Go to Definition". To get this working, all your IDE should need
is knowing where the headers are and what `#define` statements to use.

For example, when using [VS Code](https://code.visualstudio.com/) with the
[C/C++ extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode.cpptools),
add the following file to your workplace root at `.vscode/c_cpp_properties.json`:
```json
{
    "configurations": [
        {
            "name": "Linux",
            "includePath": [
                "${workspaceFolder}/**"
            ],
            "defines": [
                "VTPM=NO",
                "SIMULATION=NO",
                "USE_DA_USED=NO",
                "HASH_LIB=Ossl",
                "SYM_LIB=Ossl",
                "MATH_LIB=Ossl"
            ],
            "compilerPath": "/bin/clang",
            "cStandard": "c11",
            "cppStandard": "c++17",
            "intelliSenseMode": "clang-x64"
        }
    ],
    "version": 4
}
```