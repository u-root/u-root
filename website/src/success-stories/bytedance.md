---
title: ByteDance
date: 2023-10-10
layout: 'layouts/case-study.njk'
tags: []
permalink: success-stories/{{ title | slugify }}.html
metaDesc: 'ByteDance greatly improved their boot software stack.'
---

# ByteDance greatly improved their boot software stack

ByteDance[^1] started their Cloud Firmware project in 2020. Their work has proceeded from Version 1 to Version 3. Version 1, which has rolled out to production, is very similar to what Google deployed: ByteDance combines a Linux kernel and u-root, and places it in the UEFI firmware image, to be used as the primary boot software. This work allows them to remove a large amount of UEFI binary blobs, and further allows them to address problems they see in the UEFI software:

- Compared to the Linux community, the UEFI community is NOT active
- Can not fix UEFI issues immediately since key modules are controlled by the Independent BIOS vendor
- The UEFI working model is NOT efficient

U-root allowed ByteDance to greatly improve the boot software stack, since a large amount of C is replaced with a much smaller amount of Go. ByteDance points out that u-root allowed them to “Replace and Enhance Firmware Functions: PXE, HTTP boot, Boot Option Manager, Redfish” as well as “Provide a Diskless Linux Environment with Operation & Maintenance components, and System Stress Tools.” They resolved problems with the UEFI network stack, which in their view is “not powerful, hard to optimize.” They point out that, in their view, the Linux Network Stack is “powerful, independent of firmware vendor” and, further, “there are more linux network experts.”

ByteDance has been using u-root for four years, in a production deployment in their data centers.

[^1]: https://146a55aca6f00848c565-a7635525d40ac1c70300198708936b4e.ssl.cf1.rackcdn.com/images/8a8695ae7c7cfcffcabd2d09f1f8541566381717.pdf
