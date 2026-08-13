# Contributing to go-pcap

Thank you for helping improve native Go packet capture and tcpdump-style cBPF
filtering. Small, well-scoped changes with clear packet-layout semantics are
the easiest to review and maintain.

## Development environment

- Go 1.23 or later. Use the toolchain requested by `go.mod` when it specifies
  one.
- Linux is required for the live integration suite. Ordinary builds and unit
  tests can run without capture privileges.
- tcpdump 4.99.x is optional. It is useful for online compatibility checks and
  deliberately regenerating the checked-in cBPF golden fixtures.

Clone the repository and run `go mod download` before the first full test run
if your module cache is empty.

## Common commands

| Command | Purpose |
| --- | --- |
| `make build` | Build the `pcap` command-line utility. |
| `make test` | Run the Go test suite. |
| `make lint` | Run golangci-lint. |
| `make fmt-check` | Check Go imports/formatting, the `interface{} → any` rewrite, and tracked Shell scripts without modifying files. |
| `make vet` | Run `go vet ./...`. |
| `make bench` | Run repeatable filter benchmarks. |
| `make integration` | Run the loopback capture entry point; it skips without capture privilege. |

The live suite needs root or `CAP_NET_RAW`. To execute it rather than receive a
clean skip, use:

```sh
make build
sudo -E env "PATH=$PATH" bash test/run.sh
```

## Test expectations

Before opening a pull request, the minimum expected verification is:

```sh
go build ./... && go vet ./... && make test
```

Run `make fmt-check` and `make lint` for Go changes. Add focused unit or VM
tests for behavior changes, and run the privileged integration suite when a
change can affect live capture, filter installation, or loopback behavior.

The tcpdump decision-equivalence fixtures are embedded from
`filter/testdata/tcpdump-4.99.0-libpcap-1.10.0/`. Do not move them. Refresh a
fixture only for a reviewed semantic change; the fixture README describes the
generation command and provenance requirements.

## Commits and pull requests

Keep commits small and organized by logical unit. Separate a package-local
change from its caller, documentation, or workflow updates when that makes the
history easier to review. Commit messages should be concise and explain why
the change is needed, not only what files changed.

Describe the packet layout for filter changes: whether input begins at an
Ethernet header or an IP header, the expected behavior for L2-only predicates,
and the tests that prove it. Do not mix generated fixture refreshes or
unrelated formatting with a semantic compiler change.

## Further reading

- [Architecture and compiler internals](docs/concepts/architecture.md)
- [Adding a filter primitive](docs/contributing/new-primitive.md)
- [Testing](docs/contributing/testing.md)

## License

This project is licensed under Apache-2.0. New source files must include the
`Copyright 2026 The HuaTuo Authors` Apache-2.0 header used by the files in
`filter/`. Preserve applicable third-party notices when modifying derived
code.
