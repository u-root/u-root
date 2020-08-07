# Zero Touch Provisioning (ZTP) DHCPv6 Parsing for Network Hardware Vendors

## Currently Supported Vendors For DHCPv6 ZTP
 - Arista
 - ZPE

## Why Do We Need This?
Many network hardware vendors support features that allow network devices to provision themselves with proper supporting automation/tools. Network devices can rely on DHCP and other methods to gather bootfile info, IPs, etc. DHCPv6 Vendor options provides us Vendor Name, Make, Model, and Serial Number data. This data can be used to uniquely identify individual network devices at provisioning time and can be used by tooling to make decisions necessary to correctly and reliably provision a network device.

For more details on a large-scale ZTP deployment, check out how this is done at Facebook, [Scaling Backbone Networks Through Zero Touch Provisioning](https://code.fb.com/networking-traffic/scaling-the-facebook-backbone-through-zero-touch-provisioning/).


### Example Data
Vendor specific data is commonly in a delimiter separated format containing Vendor Name, Model, Make, and Serial Number. This of course will vary per vendor and there could be more or less data.
Vendor;Model;Version;SerialNumber
`Arista;DCS-7060;01.011;ZZZ00000000`
