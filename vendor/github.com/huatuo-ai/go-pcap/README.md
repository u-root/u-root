<p align="center">
  <img src="docs/img/go-pcap-logo-compact.png" alt="Go-Pcap shoebill wordmark logo" width="600" />
</p>

<p align="center">
  <strong>Native Go Packet Capture, tcpdump-style cBPF Compilation, CGO-free Cross Builds</strong>
</p>

<p align="center">
  <a href="https://github.com/huatuo-ai/go-pcap/stargazers"><img src="https://img.shields.io/github/stars/huatuo-ai/go-pcap?style=social" alt="GitHub Stars" /></a>
  <a href="https://github.com/huatuo-ai/go-pcap/issues"><img src="https://img.shields.io/github/issues/huatuo-ai/go-pcap" alt="GitHub Issues" /></a>
  <a href="./LICENSE"><img src="https://img.shields.io/badge/License-Apache_2.0-green" alt="Apache 2.0 License" /></a>
  <a href="./CONTRIBUTING.md"><img src="https://img.shields.io/badge/PRs-welcome-brightgreen" alt="PRs Welcome" /></a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-Native_Packet_Capture-00ADD8" alt="Native Go packet capture" />
  <img src="https://img.shields.io/badge/cBPF-tcpdump_style_filters-0B3C4A" alt="tcpdump-style cBPF filters" />
  <img src="https://img.shields.io/badge/CGO-Zero_dependency-blue" alt="CGO-free builds" />
</p>

<p align="center">
  <a href="./README_CN.md"><strong>中文文档</strong></a> ·
  <a href="./docs/index.md"><strong>Documentation</strong></a> ·
  <a href="./examples/README.md"><strong>Examples</strong></a>
</p>

---

## What is Go-Pcap

**Go-Pcap** is a native Go packet-capture library and tcpdump-style cBPF filter compiler. It provides a libpcap-like capture surface without CGO, making `CGO_ENABLED=0` builds and cross-compilation straightforward.

![go-pcap demo](demo.gif)

## Key Features

- **Native Go Packet Capture**: Provides a libpcap-like capture API without requiring CGO.
- **tcpdump-style Filters**: Compiles common protocol, host, network, port, and logical expressions to cBPF.
- **Ethernet and Raw IP**: Supports both Ethernet (`EN10MB`) and raw IP (`RAW`) packet layouts.
- **Protocol Coverage**: Supports IPv4, IPv6, ARP/RARP, TCP, UDP, ICMP, ICMP6, IGMP, PIM, ESP, AH, VRRP, VLAN, and MPLS traffic.
- **Packet Capture CLI**: Prints tcpdump-style summaries and supports commonly used tcpdump display options.
- **Cross Compilation**: Builds with `CGO_ENABLED=0` for supported Linux and macOS/Darwin targets.

## Getting Started

Install Go-Pcap in your Go module:

```sh
go get github.com/huatuo-ai/go-pcap@latest
```

See the [examples](examples/README.md) for library usage and the [documentation](docs/index.md) for detailed guides.

## Filter Language

| Capability | Example |
| :--- | :--- |
| Protocol and port | `tcp and port 443` |
| Direction and range | `src portrange 1000-2000` |
| IPv6 | `ip6 and udp and port 53` |
| Packet fields | `tcp[tcpflags] & (tcp-syn\|tcp-ack) == (tcp-syn\|tcp-ack)` |
| Encapsulation | `vlan 100 and tcp port 443`, `mpls and ip` |
| Logical expressions | `tcp and port 80 or udp`, `not (tcp or udp)` |

See the [filter language guide](docs/guides/filter-language.md) for more details.

## Platform Support

| Platform | Capture Support | Notes |
| :--- | :--- | :--- |
| Linux | Supported | Uses AF_PACKET and requires packet-capture privileges |
| macOS/Darwin | Supported | Requires packet-capture privileges |

## Documentation

For more information, visit the [Go-Pcap documentation](docs/index.md).

## Contributing

Issues and pull requests are welcome, especially for new protocol support, link types, compatibility cases, and performance work.

See [CONTRIBUTING.md](CONTRIBUTING.md) for local checks and pull-request expectations. The [documentation](docs/index.md) includes deeper guides for the [architecture](docs/concepts/architecture.md), [compiler internals](docs/concepts/compiler-internals.md), and [new filter primitives](docs/contributing/new-primitive.md).

## License

Go-Pcap is derived from [packetcap/go-pcap](https://github.com/packetcap/go-pcap) and is open source under the [Apache License 2.0](LICENSE).
