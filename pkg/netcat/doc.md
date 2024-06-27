# Reimplementation of ncat 

## Status

### Implementation

Some netcat functionality may be outdated and not really necessary. They are marked using a (?)

- [x] cli representation
- [x] argument parsing
- [ ] TLS
    - [x] TCP
    - [x] UDP
    - [ ] TLS Options
      - [x] Enabled
      - [x] CertFilePath
      - [x] KeyFilePath
      - [x] VerifyTrust
      - [x] TrustFilePath
      - [ ] Ciphers (?)
      - [x] SNI
      - [x] ALPN 
- [x] socket options (?)
- [x] testing
- [ ] Listen Options
  - [x] Max Connections
  - [x] KeepOpen
  - [ ] BrokerMode (?)
  - [ ] Chat Mode (?)
- [x] AccessControl (allow, deny)
- [x] CommandExec
  - [x] Native
  - [x] Shell
  - [ ] Lua (?)
- [x] OutputOptions
  - [x] Outfile
  - [x] OutHex
  - [x] AppendOutput
  - [x] Verbose
- [ ] Misc
  - [x] EOL character
  - [x] RecvOnly
  - [x] SendOnly
  - [x] NoShutdown
  - [x] NoDns
  - [ ] Telnet (?)
- [x] Timing Options
  - [x] Delay
  - [x] Timeout
  - [x] Wait
- [ ] Proxy (?)
  - [ ] Address
  - [ ] DNS
  - [ ] Type
  - [ ] AuthType  
- [ ] Connection Mode Options
  - [x] Source Host
  - [x] Source Port
  - [x]  ZERO I/O
  - [ ]  LooseSourceRouterPoints (?)
  - [ ] LooseSourcePointer (?)
- [ ] tinygo adjustment
    - [ ] net.Dial is only implemented for "tcp", "tcp4", "udp", and "udp4".
    - [ ] Not more than \<boardconfig\> many connections

### Testing

- [x] unit tests for conversions
- [ ] integration tests 
- [ ] `vmtest` integration test setup
