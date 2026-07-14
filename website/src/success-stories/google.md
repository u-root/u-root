---
title: Google
date: 2023-11-15
layout: 'layouts/case-study.njk'
tags: []
permalink: success-stories/{{ title | slugify }}.html
metaDesc: 'Google worked with 9elements to bring u-root test coverage up significantly.'
---

# Google worked with 9elements to bring u-root test coverage up significantly

Google deployed u-root-based firmware at scale in its data centers in Dec. 2020, and since that time its use has continued to grow, including support for ARM mainboards. The firmware images are built in seconds, not hours, as with UEFI; they allow Google to remove several megabytes of C-based UEFI software; and they allow Google to apply its world-class expertise in Linux and Go programs to firmware.

The Go toolchain has strong mitigations for supply-chain attacks[^1], always a concern with firmware[^2]. With recent improvements in Go code size, it has become quite easy to fit u-root’s 160 commands in a four MiB xz-compressed binary. Google worked with 9elements to bring u-root test coverage to 75% in 2021.

[^1]: “How Go Mitigates Supply Chain Attacks”, https://go.dev/blog/supply-chain
[^2]: See, e.g, “Binarly Discloses Multiple Firmware Vulnerabilities in Qualcomm and Lenovo ARM-based Devices”, https://www.binarly.io/news/Binarly-Discloses-Multiple-Firmware-Vulnerabilities-in-Qualcomm-and-Lenovo-ARM-based-Devices/index.html
