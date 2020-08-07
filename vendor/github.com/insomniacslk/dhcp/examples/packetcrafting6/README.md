# DHCPv6 packet crafting

It is easy to parse, create and manipulate DHCPv6 packets. The `DHCPv6`
interface is the central way to move packets around. Two concrete
implementations, `DHCPv6Message` and `DHCPv6Relay` take care of the
details. The example in [main.go](./main.go) shows how to craft a DHCP6
Solicit message with a custom DUID_LLT, encapsulate it in a Relay message,
and print its details.


The output (slightly modified for readability) is
```
$ go run main.go
2018/11/08 13:56:31 DHCPv6Relay
  messageType=RELAY-FORW
  hopcount=0
  linkaddr=2001:db8::1
  peeraddr=2001:db8::2
  options=[OptRelayMsg{relaymsg=DHCPv6Message(messageType=SOLICIT transactionID=0x9e0242, 4 options)}]

2018/11/08 13:56:31 [12 0 32 1 13 184 0 0 0 0 0 0 0 0 0 0 0 1 32 1 13 184
                     0 0 0 0 0 0 0 0 0 0 0 2 0 9 0 52 1 158 2 66 0 1 0 14
                     0 1 0 1 35 118 253 15 0 250 206 176 12 0 0 6 0 4 0 23
                     0 24 0 8 0 2 0 0 0 3 0 12 250 206 176 12 0 0 14 16 0
                     0 21 24]
```
