package dhcpv6

import (
	"fmt"
)

// TransactionID is a DHCPv6 Transaction ID defined by RFC 3315, Section 6.
type TransactionID [3]byte

// String prints the transaction ID as a hex value.
func (xid TransactionID) String() string {
	return fmt.Sprintf("0x%x", xid[:])
}

// MessageType represents the kind of DHCPv6 message.
type MessageType uint8

// The DHCPv6 message types defined per RFC 3315, Section 5.3.
const (
	// MessageTypeNone is used internally and is not part of the RFC.
	MessageTypeNone               MessageType = 0
	MessageTypeSolicit            MessageType = 1
	MessageTypeAdvertise          MessageType = 2
	MessageTypeRequest            MessageType = 3
	MessageTypeConfirm            MessageType = 4
	MessageTypeRenew              MessageType = 5
	MessageTypeRebind             MessageType = 6
	MessageTypeReply              MessageType = 7
	MessageTypeRelease            MessageType = 8
	MessageTypeDecline            MessageType = 9
	MessageTypeReconfigure        MessageType = 10
	MessageTypeInformationRequest MessageType = 11
	MessageTypeRelayForward       MessageType = 12
	MessageTypeRelayReply         MessageType = 13
	MessageTypeLeaseQuery         MessageType = 14
	MessageTypeLeaseQueryReply    MessageType = 15
	MessageTypeLeaseQueryDone     MessageType = 16
	MessageTypeLeaseQueryData     MessageType = 17
)

// String prints the message type name.
func (m MessageType) String() string {
	if s, ok := messageTypeToStringMap[m]; ok {
		return s
	}
	return fmt.Sprintf("unknown (%d)", m)
}

// messageTypeToStringMap contains the mapping of MessageTypes to
// human-readable strings.
var messageTypeToStringMap = map[MessageType]string{
	MessageTypeSolicit:            "SOLICIT",
	MessageTypeAdvertise:          "ADVERTISE",
	MessageTypeRequest:            "REQUEST",
	MessageTypeConfirm:            "CONFIRM",
	MessageTypeRenew:              "RENEW",
	MessageTypeRebind:             "REBIND",
	MessageTypeReply:              "REPLY",
	MessageTypeRelease:            "RELEASE",
	MessageTypeDecline:            "DECLINE",
	MessageTypeReconfigure:        "RECONFIGURE",
	MessageTypeInformationRequest: "INFORMATION-REQUEST",
	MessageTypeRelayForward:       "RELAY-FORW",
	MessageTypeRelayReply:         "RELAY-REPL",
	MessageTypeLeaseQuery:         "LEASEQUERY",
	MessageTypeLeaseQueryReply:    "LEASEQUERY-REPLY",
	MessageTypeLeaseQueryDone:     "LEASEQUERY-DONE",
	MessageTypeLeaseQueryData:     "LEASEQUERY-DATA",
}

// OptionCode is a single byte representing the code for a given Option.
type OptionCode uint16

// String returns the option code name.
func (o OptionCode) String() string {
	if s, ok := optionCodeToString[o]; ok {
		return s
	}
	return fmt.Sprintf("unknown (%d)", o)
}

// All DHCPv6 options.
const (
	OptionClientID    OptionCode = 1
	OptionServerID    OptionCode = 2
	OptionIANA        OptionCode = 3
	OptionIATA        OptionCode = 4
	OptionIAAddr      OptionCode = 5
	OptionORO         OptionCode = 6
	OptionPreference  OptionCode = 7
	OptionElapsedTime OptionCode = 8
	OptionRelayMsg    OptionCode = 9
	// skip 10
	OptionAuth                           OptionCode = 11
	OptionUnicast                        OptionCode = 12
	OptionStatusCode                     OptionCode = 13
	OptionRapidCommit                    OptionCode = 14
	OptionUserClass                      OptionCode = 15
	OptionVendorClass                    OptionCode = 16
	OptionVendorOpts                     OptionCode = 17
	OptionInterfaceID                    OptionCode = 18
	OptionReconfMessage                  OptionCode = 19
	OptionReconfAccept                   OptionCode = 20
	OptionSIPServersDomainNameList       OptionCode = 21
	OptionSIPServersIPv6AddressList      OptionCode = 22
	OptionDNSRecursiveNameServer         OptionCode = 23
	OptionDomainSearchList               OptionCode = 24
	OptionIAPD                           OptionCode = 25
	OptionIAPrefix                       OptionCode = 26
	OptionNISServers                     OptionCode = 27
	OptionNISPServers                    OptionCode = 28
	OptionNISDomainName                  OptionCode = 29
	OptionNISPDomainName                 OptionCode = 30
	OptionSNTPServerList                 OptionCode = 31
	OptionInformationRefreshTime         OptionCode = 32
	OptionBCMCSControllerDomainNameList  OptionCode = 33
	OptionBCMCSControllerIPv6AddressList OptionCode = 34
	// skip 35
	OptionGeoConfCivic                            OptionCode = 36
	OptionRemoteID                                OptionCode = 37
	OptionRelayAgentSubscriberID                  OptionCode = 38
	OptionFQDN                                    OptionCode = 39
	OptionPANAAuthenticationAgent                 OptionCode = 40
	OptionNewPOSIXTimezone                        OptionCode = 41
	OptionNewTZDBTimezone                         OptionCode = 42
	OptionEchoRequest                             OptionCode = 43
	OptionLQQuery                                 OptionCode = 44
	OptionClientData                              OptionCode = 45
	OptionCLTTime                                 OptionCode = 46
	OptionLQRelayData                             OptionCode = 47
	OptionLQClientLink                            OptionCode = 48
	OptionMIPv6HomeNetworkIDFQDN                  OptionCode = 49
	OptionMIPv6VisitedHomeNetworkInformation      OptionCode = 50
	OptionLoSTServer                              OptionCode = 51
	OptionCAPWAPAccessControllerAddresses         OptionCode = 52
	OptionRelayID                                 OptionCode = 53
	OptionIPv6AddressMOS                          OptionCode = 54
	OptionIPv6FQDNMOS                             OptionCode = 55
	OptionNTPServer                               OptionCode = 56
	OptionV6AccessDomain                          OptionCode = 57
	OptionSIPUACSList                             OptionCode = 58
	OptionBootfileURL                             OptionCode = 59
	OptionBootfileParam                           OptionCode = 60
	OptionClientArchType                          OptionCode = 61
	OptionNII                                     OptionCode = 62
	OptionGeolocation                             OptionCode = 63
	OptionAFTRName                                OptionCode = 64
	OptionERPLocalDomainName                      OptionCode = 65
	OptionRSOO                                    OptionCode = 66
	OptionPDExclude                               OptionCode = 67
	OptionVirtualSubnetSelection                  OptionCode = 68
	OptionMIPv6IdentifiedHomeNetworkInformation   OptionCode = 69
	OptionMIPv6UnrestrictedHomeNetworkInformation OptionCode = 70
	OptionMIPv6HomeNetworkPrefix                  OptionCode = 71
	OptionMIPv6HomeAgentAddress                   OptionCode = 72
	OptionMIPv6HomeAgentFQDN                      OptionCode = 73
)

// optionCodeToString maps DHCPv6 OptionCodes to human-readable strings.
var optionCodeToString = map[OptionCode]string{
	OptionClientID:                                "OPTION_CLIENTID",
	OptionServerID:                                "OPTION_SERVERID",
	OptionIANA:                                    "OPTION_IA_NA",
	OptionIATA:                                    "OPTION_IA_TA",
	OptionIAAddr:                                  "OPTION_IAADDR",
	OptionORO:                                     "OPTION_ORO",
	OptionPreference:                              "OPTION_PREFERENCE",
	OptionElapsedTime:                             "OPTION_ELAPSED_TIME",
	OptionRelayMsg:                                "OPTION_RELAY_MSG",
	OptionAuth:                                    "OPTION_AUTH",
	OptionUnicast:                                 "OPTION_UNICAST",
	OptionStatusCode:                              "OPTION_STATUS_CODE",
	OptionRapidCommit:                             "OPTION_RAPID_COMMIT",
	OptionUserClass:                               "OPTION_USER_CLASS",
	OptionVendorClass:                             "OPTION_VENDOR_CLASS",
	OptionVendorOpts:                              "OPTION_VENDOR_OPTS",
	OptionInterfaceID:                             "OPTION_INTERFACE_ID",
	OptionReconfMessage:                           "OPTION_RECONF_MSG",
	OptionReconfAccept:                            "OPTION_RECONF_ACCEPT",
	OptionSIPServersDomainNameList:                "SIP Servers Domain Name List",
	OptionSIPServersIPv6AddressList:               "SIP Servers IPv6 Address List",
	OptionDNSRecursiveNameServer:                  "DNS Recursive Name Server",
	OptionDomainSearchList:                        "Domain Search List",
	OptionIAPD:                                    "OPTION_IA_PD",
	OptionIAPrefix:                                "OPTION_IAPREFIX",
	OptionNISServers:                              "OPTION_NIS_SERVERS",
	OptionNISPServers:                             "OPTION_NISP_SERVERS",
	OptionNISDomainName:                           "OPTION_NIS_DOMAIN_NAME",
	OptionNISPDomainName:                          "OPTION_NISP_DOMAIN_NAME",
	OptionSNTPServerList:                          "SNTP Server List",
	OptionInformationRefreshTime:                  "Information Refresh Time",
	OptionBCMCSControllerDomainNameList:           "BCMCS Controller Domain Name List",
	OptionBCMCSControllerIPv6AddressList:          "BCMCS Controller IPv6 Address List",
	OptionGeoConfCivic:                            "OPTION_GEOCONF",
	OptionRemoteID:                                "OPTION_REMOTE_ID",
	OptionRelayAgentSubscriberID:                  "Relay-Agent Subscriber ID",
	OptionFQDN:                                    "FQDN",
	OptionPANAAuthenticationAgent:                 "PANA Authentication Agent",
	OptionNewPOSIXTimezone:                        "OPTION_NEW_POSIX_TIME_ZONE",
	OptionNewTZDBTimezone:                         "OPTION_NEW_TZDB_TIMEZONE",
	OptionEchoRequest:                             "Echo Request",
	OptionLQQuery:                                 "OPTION_LQ_QUERY",
	OptionClientData:                              "OPTION_CLIENT_DATA",
	OptionCLTTime:                                 "OPTION_CLT_TIME",
	OptionLQRelayData:                             "OPTION_LQ_RELAY_DATA",
	OptionLQClientLink:                            "OPTION_LQ_CLIENT_LINK",
	OptionMIPv6HomeNetworkIDFQDN:                  "MIPv6 Home Network ID FQDN",
	OptionMIPv6VisitedHomeNetworkInformation:      "MIPv6 Visited Home Network Information",
	OptionLoSTServer:                              "LoST Server",
	OptionCAPWAPAccessControllerAddresses:         "CAPWAP Access Controller Addresses",
	OptionRelayID:                                 "RELAY_ID",
	OptionIPv6AddressMOS:                          "OPTION-IPv6_Address-MoS",
	OptionIPv6FQDNMOS:                             "OPTION-IPv6-FQDN-MoS",
	OptionNTPServer:                               "OPTION_NTP_SERVER",
	OptionV6AccessDomain:                          "OPTION_V6_ACCESS_DOMAIN",
	OptionSIPUACSList:                             "OPTION_SIP_UA_CS_LIST",
	OptionBootfileURL:                             "OPT_BOOTFILE_URL",
	OptionBootfileParam:                           "OPT_BOOTFILE_PARAM",
	OptionClientArchType:                          "OPTION_CLIENT_ARCH_TYPE",
	OptionNII:                                     "OPTION_NII",
	OptionGeolocation:                             "OPTION_GEOLOCATION",
	OptionAFTRName:                                "OPTION_AFTR_NAME",
	OptionERPLocalDomainName:                      "OPTION_ERP_LOCAL_DOMAIN_NAME",
	OptionRSOO:                                    "OPTION_RSOO",
	OptionPDExclude:                               "OPTION_PD_EXCLUDE",
	OptionVirtualSubnetSelection:                  "Virtual Subnet Selection",
	OptionMIPv6IdentifiedHomeNetworkInformation:   "MIPv6 Identified Home Network Information",
	OptionMIPv6UnrestrictedHomeNetworkInformation: "MIPv6 Unrestricted Home Network Information",
	OptionMIPv6HomeNetworkPrefix:                  "MIPv6 Home Network Prefix",
	OptionMIPv6HomeAgentAddress:                   "MIPv6 Home Agent Address",
	OptionMIPv6HomeAgentFQDN:                      "MIPv6 Home Agent FQDN",
}
