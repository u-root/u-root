# Reimplementation of ncat 

## Status

### Implementation Status

- [x] cli representation
- [x] argument parsing
- [x] TLS
    - [x] TCP
    - [x] UDP
    - [x] TLS Options
      - [x] Enabled
      - [x] CertFilePath
      - [x] KeyFilePath
      - [x] VerifyTrust
      - [x] TrustFilePath
      - [ ] Ciphers 
      - [x] SNI
      - [x] ALPN 
- [x] socket options
- [x] testing
- [x] Listen Options
  - [x] Max Connections
  - [x] KeepOpen
  - [x] BrokerMode 
  - [x] Chat Mode
- [x] AccessControl (allow, deny)
- [x] CommandExec
  - [x] Native
  - [x] Shell
  - [ ] Lua 
- [x] OutputOptions
  - [x] Outfile
  - [x] OutHex
  - [x] AppendOutput
  - [x] Verbose
- [x] Misc
  - [x] EOL character
  - [x] RecvOnly
  - [x] SendOnly
  - [x] NoShutdown
  - [x] NoDns
  - [ ] Telnet 
- [x] Timing Options
  - [x] Delay
  - [x] Timeout
  - [x] Wait
- [x] Proxy
  - [x] Address
  - [x] DNS
  - [x] Type
  - [x] AuthType  
- [x] Connection Mode Options
  - [x] Source Host
  - [x] Source Port
  - [x]  ZERO I/O
  - [ ]  LooseSourceRouterPoints 
  - [ ] LooseSourcePointer 
- [ ] tinygo adjustment
    - [ ] net.Dial is only implemented for "tcp", "tcp4", "udp", and "udp4".
    - [ ] Not more than \<boardconfig\> many connections

### Testing

- [x] unit tests for conversions
- [ ] integration tests 
