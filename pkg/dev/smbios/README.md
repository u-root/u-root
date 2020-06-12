# SMBIOS Information Library

TODO: godoc link

## Adding a new type

Types are defined in [DMTF DSP0134](https://www.dmtf.org/dsp/DSP0134). As of
July 2019, the most recent published version is 3.2.0.

Adding a type largely consists of copy-pasting a bunch of information from the
doicument to define a Go struct and its string representation. To make that
easier, a set of tools is provided ina Gist
[here](https://gist.github.com/rojer/4fa173442fb00e24dc7b9d120c2e30af). These
are Python scripts that take chunks of data copy-pasted from this document and
produce a struct, anum or bit field declaration. These are by no means perfect
and do not handle all the cases of weird formatting in the document, and may be
broken completely by future formatting changes (but hopefully can be fixed if
needed).

So, adding a type should look something like this:

*   Create `typeXXX_new_information.go` with boilerplate from some other type
    (copyright, imports).

*   Generate the struct: `tools/gen_struct.py NewInformation >>
    typeXXX_new_information.go`.

*   Generate enums and bitfields if needed.

*   Examine the generated code and tweak it manually if necessart.

*   Run `gofmt` on it: `gofmt typeXXX_new_information.go`.

*   Add new Type and a corresponding constant to [table_type.go](table_type.go).

*   Implement the constructor function, `NewNewInformation` in this case.

*   Most of the parsing can be done by reflection-based parser, see
    `parseStruct`.

*   Add data validation and any additional custom parsing code necessary.

*   Implement the `String()` method for `*NewInformation`.

*   Add getter to the `Info` type, `GetNewInformation`.
