# DHCPv6 server

A DHCPv6 server requires the user to implement a request handler. Basically the
user has to provide the logic to answer to each packet. The library offers a few
facilities to forge response packets, e.g. `NewAdvertiseFromSolicit`,
`NewReplyFromDHCPv6Message` and so on. Look at the source code to see what's
available.

An example server that will print (but not reply to) the client's request is
shown in [main.go](./main.go)
