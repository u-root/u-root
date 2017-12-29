package dhcp4

// OpCode is the BOOTP message type as defined by RFC 2131, Section 2.
//
// Note that the DHCP message type is embedded via OptionDHCPMessageType.
type OpCode uint8

const (
	BootRequest OpCode = 1
	BootReply   OpCode = 2
)

type OptionCode uint8

// Option codes defined by RFC 2132. (Incomplete)
const (
	End                                              OptionCode = 255
	Pad                                              OptionCode = 0
	OptionSubnetMask                                 OptionCode = 1
	OptionTimeOffset                                 OptionCode = 2
	OptionRouters                                    OptionCode = 3
	OptionTimeServers                                OptionCode = 4
	OptionNameServers                                OptionCode = 5
	OptionDomainNameServers                          OptionCode = 6
	OptionLogServers                                 OptionCode = 7
	OptionCookieServers                              OptionCode = 8
	OptionLPRServers                                 OptionCode = 9
	OptionImpressServers                             OptionCode = 10
	OptionResourceLocationServers                    OptionCode = 11
	OptionHostName                                   OptionCode = 12
	OptionBootFileSize                               OptionCode = 13
	OptionMeritDumpFile                              OptionCode = 14
	OptionDomainName                                 OptionCode = 15
	OptionSwapServer                                 OptionCode = 16
	OptionRootPath                                   OptionCode = 17
	OptionExtensionsPath                             OptionCode = 18
	OptionIPForwardingEnableDisable                  OptionCode = 19
	OptionNonLocalSourceRoutingEnableDisable         OptionCode = 20
	OptionPolicyFilter                               OptionCode = 21
	OptionMaximumDatagramReassemblySize              OptionCode = 22
	OptionDefaultIPTimeToLive                        OptionCode = 23
	OptionPathMTUAgingTimeout                        OptionCode = 24
	OptionPathMTUPlateauTable                        OptionCode = 25
	OptionInterfaceMTU                               OptionCode = 26
	OptionAllSubnetsAreLocal                         OptionCode = 27
	OptionBroadcastAddress                           OptionCode = 28
	OptionPerformMaskDiscovery                       OptionCode = 29
	OptionMaskSupplier                               OptionCode = 30
	OptionPerformRouterDiscovery                     OptionCode = 31
	OptionRouterSolicitationAddress                  OptionCode = 32
	OptionStaticRoute                                OptionCode = 33
	OptionTrailerEncapsulation                       OptionCode = 34
	OptionARPCacheTimeout                            OptionCode = 35
	OptionEthernetEncapsulation                      OptionCode = 36
	OptionTCPDefaultTTL                              OptionCode = 37
	OptionTCPKeepaliveInterval                       OptionCode = 38
	OptionTCPKeepaliveGarbage                        OptionCode = 39
	OptionNetworkInformationServiceDomain            OptionCode = 40
	OptionNetworkInformationServers                  OptionCode = 41
	OptionNetworkTimeProtocolServers                 OptionCode = 42
	OptionVendorSpecificInformation                  OptionCode = 43
	OptionNetBIOSOverTCPIPNameServer                 OptionCode = 44
	OptionNetBIOSOverTCPIPDatagramDistributionServer OptionCode = 45
	OptionNetBIOSOverTCPIPNodeType                   OptionCode = 46
	OptionNetBIOSOverTCPIPScope                      OptionCode = 47
	OptionXWindowSystemFontServer                    OptionCode = 48
	OptionXWindowSystemDisplayManager                OptionCode = 49

	// DHCP extensions.
	OptionRequestedIPAddress     OptionCode = 50
	OptionIPAddressLeaseTime     OptionCode = 51
	OptionOverload               OptionCode = 52
	OptionDHCPMessageType        OptionCode = 53
	OptionServerIdentifier       OptionCode = 54
	OptionParameterRequestList   OptionCode = 55
	OptionMessage                OptionCode = 56
	OptionMaximumDHCPMessageSize OptionCode = 57
	OptionRenewalTimeValue       OptionCode = 58
	OptionRebindingTimeValue     OptionCode = 59
	OptionVendorClassIdentifier  OptionCode = 60
	OptionClientIdentifier       OptionCode = 61
	OptionTFTPServerName         OptionCode = 66
	OptionBootFileName           OptionCode = 67
)
