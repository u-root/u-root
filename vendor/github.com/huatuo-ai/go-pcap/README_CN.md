<p align="center">
  <img src="docs/img/go-pcap-logo-compact.png" alt="Go-Pcap 鲸头鹳字标 logo" width="600" />
</p>

<p align="center">
  <strong>原生 Go 抓包、tcpdump 风格 cBPF 编译、无 CGO 交叉构建</strong>
</p>

<p align="center">
  <a href="https://github.com/huatuo-ai/go-pcap/stargazers"><img src="https://img.shields.io/github/stars/huatuo-ai/go-pcap?style=social" alt="GitHub Stars" /></a>
  <a href="https://github.com/huatuo-ai/go-pcap/issues"><img src="https://img.shields.io/github/issues/huatuo-ai/go-pcap" alt="GitHub Issues" /></a>
  <a href="./LICENSE"><img src="https://img.shields.io/badge/License-Apache_2.0-green" alt="Apache 2.0 License" /></a>
  <a href="./CONTRIBUTING.md"><img src="https://img.shields.io/badge/PRs-welcome-brightgreen" alt="PRs Welcome" /></a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-Native_Packet_Capture-00ADD8" alt="原生 Go 抓包" />
  <img src="https://img.shields.io/badge/cBPF-tcpdump_style_filters-0B3C4A" alt="tcpdump 风格 cBPF 过滤器" />
  <img src="https://img.shields.io/badge/CGO-Zero_dependency-blue" alt="无 CGO 构建" />
</p>

<p align="center">
  <a href="./README.md"><strong>English</strong></a> ·
  <a href="./docs/index.md"><strong>Documentation</strong></a> ·
  <a href="./examples/README.md"><strong>Examples</strong></a>
</p>

---

## 什么是 Go-Pcap

**Go-Pcap** 是一个原生 Go 的抓包库和 tcpdump 风格 cBPF 过滤编译器。它提供类似 libpcap 的抓包接口且不依赖 CGO，便于使用 `CGO_ENABLED=0` 构建和交叉编译。

![go-pcap 演示](demo.gif)

## 核心功能

- **原生 Go 抓包**：提供类似 libpcap 的抓包 API，无需 CGO。
- **tcpdump 风格过滤器**：将常用的协议、主机、网络、端口和逻辑表达式编译为 cBPF。
- **Ethernet 与裸 IP**：支持 Ethernet（`EN10MB`）和裸 IP（`RAW`）报文布局。
- **协议覆盖**：支持 IPv4、IPv6、ARP/RARP、TCP、UDP、ICMP、ICMP6、IGMP、PIM、ESP、AH、VRRP、VLAN 和 MPLS 流量。
- **抓包 CLI**：输出 tcpdump 风格摘要，并支持常用的 tcpdump 显示选项。
- **交叉编译**：使用 `CGO_ENABLED=0` 为支持的 Linux 和 macOS/Darwin 目标构建。

## 快速开始

在 Go 模块中安装 Go-Pcap：

```sh
go get github.com/huatuo-ai/go-pcap@latest
```

库的使用方法见[示例](examples/README.md)，详细说明见[文档](docs/index.md)。

## 过滤语法

| 能力 | 示例 |
| :--- | :--- |
| 协议与端口 | `tcp and port 443` |
| 方向与端口范围 | `src portrange 1000-2000` |
| IPv6 | `ip6 and udp and port 53` |
| 报文字段 | `tcp[tcpflags] & (tcp-syn\|tcp-ack) == (tcp-syn\|tcp-ack)` |
| 封装协议 | `vlan 100 and tcp port 443`、`mpls and ip` |
| 逻辑表达式 | `tcp and port 80 or udp`、`not (tcp or udp)` |

更多信息请参阅[过滤语法指南](docs/guides/filter-language.md)。

## 平台支持

| 平台 | 抓包支持 | 说明 |
| :--- | :--- | :--- |
| Linux | 支持 | 使用 AF_PACKET，需要抓包权限 |
| macOS/Darwin | 支持 | 需要抓包权限 |

## 文档

更多信息请参阅 [Go-Pcap 文档](docs/index.md)。

## 贡献

欢迎提交 issue 和 PR，尤其是新的协议、链路类型、兼容性用例和性能改进。

本地检查和 PR 约定见 [CONTRIBUTING.md](CONTRIBUTING.md)。[项目文档](docs/index.md)包含[架构](docs/concepts/architecture.md)、[编译器内部实现](docs/concepts/compiler-internals.md)和[新增过滤原语](docs/contributing/new-primitive.md)等深入指南。

## 许可证

Go-Pcap 基于 [packetcap/go-pcap](https://github.com/packetcap/go-pcap) 演进，并依据 [Apache License 2.0](LICENSE) 开源。
