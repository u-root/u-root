// Package dhcp4opts implements Options parsing for DHCPv4 options as described in RFC 2132.
//
// Not all options are currently implemented.
package dhcp4opts

import (
	"github.com/u-root/dhcp4"
)

// GetSubnetMask returns the subnet mask of `o`.
//
// The subnet mask option is defined by RFC 2132, Section 3.3.
func GetSubnetMask(o dhcp4.Options) (SubnetMask, error) {
	v, err := o.Get(dhcp4.OptionSubnetMask)
	if err != nil {
		return nil, err
	}
	var s SubnetMask
	return s, (&s).UnmarshalBinary(v)
}

// GetRouters returns the list of router IPs in `o`.
//
// The router option is defined by RFC 2132, Section 3.5.
func GetRouters(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionRouters, o)
}

// GetTimeServers returns the list of time server IPs in `o`.
//
// The time server option is defined by RFC 2132, Section 3.6.
func GetTimeServers(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionTimeServers, o)
}

// GetNameServers returns the list of IEN 116 name server IPs in `o`.
//
// The name server option is defined by RFC 2132, Section 3.7.
func GetNameServers(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionNameServers, o)
}

// GetDomainNameServers returns the list of DNS server IPs in `o`.
//
// The domain name server option is defined by RFC 2132, Section 3.8.
func GetDomainNameServers(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionDomainNameServers, o)
}

// GetLogServers returns the list of MIT-LCS UDP log server IPs in `o`.
//
// The log server option is defined by RFC 2132, Section 3.9.
func GetLogServers(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionLogServers, o)
}

// GetCookieServers returns the list of RFC 865 cookie server IPs in `o`.
//
// The cookie server option is defined by RFC 2132, Section 3.10.
func GetCookieServers(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionCookieServers, o)
}

// GetLPRServers returns the list of RFC 1179 line printer server IPs in `o`.
//
// The LPR server option is defined by RFC 2132, Section 3.11.
func GetLPRServers(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionLPRServers, o)
}

// GetImpressServers returns the list of Imagen Impress server IPs in `o`.
//
// The impress server option is defined by RFC 2132, Section 3.12.
func GetImpressServers(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionImpressServers, o)
}

// GetResourceLocationServers returns the list of RFC 887 Resource Location
// server IPs in `o`.
//
// The resource location server option is defined by RFC 2132, Section 3.13.
func GetResourceLocationServers(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionResourceLocationServers, o)
}

// GetHostName returns the host name in `o`.
//
// The host name option is defined by RFC 2132, Section 3.14.
func GetHostName(o dhcp4.Options) (string, error) {
	return GetString(dhcp4.OptionHostName, o)
}

// GetMeritDumpFile returns the path name to be used for client crash's core
// dumps.
//
// The merit dump file is defined by RFC 2132, Section 3.16.
func GetMeritDumpFile(o dhcp4.Options) (string, error) {
	return GetString(dhcp4.OptionMeritDumpFile, o)
}

// GetDomainName returns the domain name that should be used with DNS resolvers
// in `o`.
//
// The domain name option is defined by RFC 2132, Section 3.17.
func GetDomainName(o dhcp4.Options) (string, error) {
	return GetString(dhcp4.OptionDomainName, o)
}

// GetSwapServer returns the swap server IP of `o`.
//
// The swap server option is defined by RFC 2132, Section 3.18.
func GetSwapServer(o dhcp4.Options) (IP, error) {
	return GetIP(dhcp4.OptionSwapServer, o)
}

// GetRootPath returns the disk's root path name in `o`.
//
// The root path option is defined by RFC 2132, Section 3.19.
func GetRootPath(o dhcp4.Options) (string, error) {
	return GetString(dhcp4.OptionRootPath, o)
}

// GetExtensionsPath returns the extension path name in `o`.
//
// The extension path option is defined by RFC 2132, Section 3.20.
func GetExtensionsPath(o dhcp4.Options) (string, error) {
	return GetString(dhcp4.OptionExtensionsPath, o)
}

// GetBroadcastAddress returns the client's subnet broadcast address of `o`.
//
// The broadcast address option is defined by RFC 2132, Section 5.3.
func GetBroadcastAddress(o dhcp4.Options) (IP, error) {
	return GetIP(dhcp4.OptionBroadcastAddress, o)
}

// GetRouterSolicitationAddress returns the router solicitation IP of `o`.
//
// The router solicitation address option is defined by RFC 2132, Section 5.7.
func GetRouterSolicitationAddress(o dhcp4.Options) (IP, error) {
	return GetIP(dhcp4.OptionRouterSolicitationAddress, o)
}

// GetNetworkInformationServers returns the list of NI server IPs in `o`.
//
// The network information server option is defined by RFC 2132, Section 8.2.
func GetNetworkInformationServers(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionNetworkInformationServers, o)
}

// GetNetworkTimeProtocolServers returns the list of NTP server IPs in `o`.
//
// The network time protocol server option is defined by RFC 2132, Section 8.3.
func GetNetworkTimeProtocolServers(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionNetworkTimeProtocolServers, o)
}

// GetNBNServer returns the list of NetBIOS over TCP/IP name server IPs in `o`.
//
// The NetBIOS over TCP/IP name server option is defined by RFC 2132, Section
// 8.5.
func GetNBNServer(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionNetBIOSOverTCPIPNameServer, o)
}

// GetNBDDServer returns the list of NetBIOS over TCP/IP Datagram Distribution
// server IPs in `o`.
//
// The NetBIOS over TCP/IP Datagram Distribution Server option is defined by
// RFC 2132, Section 8.6.
func GetNBDDServer(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionNetBIOSOverTCPIPDatagramDistributionServer, o)
}

// GetXWindowSystemFontServer returns the list of X window system font server
// IPs in `o`.
//
// The X window system font server option is defined by RFC 2132, Section 8.9.
func GetXWindowSystemFontServer(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionXWindowSystemFontServer, o)
}

// GetXWindowSystemDisplayManager returns the list of X window system display
// manager server IPs in `o`.
//
// The X window system display manager option is defined by RFC 2132, Section
// 8.10.
func GetXWindowSystemDisplayManager(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionXWindowSystemDisplayManager, o)
}

// GetRequestedIPAddress returns the client's requested IP in `o`.
//
// The requested IP address option is defined by RFC 2132, Section 9.1.
func GetRequestedIPAddress(o dhcp4.Options) (IP, error) {
	return GetIP(dhcp4.OptionRequestedIPAddress, o)
}

// GetServerIdentifier returns the server's identifier IP in `o`.
//
// The server identifier option is defined by RFC 2132, Section 9.5.
func GetServerIdentifier(o dhcp4.Options) (IP, error) {
	return GetIP(dhcp4.OptionServerIdentifier, o)
}

// GetDHCPMessageType returns the DHCP message type of `o`.
//
// The DHCP message type option is defined by RFC 2132, Section 9.6.
func GetDHCPMessageType(o dhcp4.Options) (DHCPMessageType, error) {
	v, err := o.Get(dhcp4.OptionDHCPMessageType)
	if err != nil {
		return 0, err
	}

	var d DHCPMessageType
	return d, (&d).UnmarshalBinary(v)
}

// GetParameterRequestList returns the list of requested DHCP option codes in
// `o`.
//
// The parameter request list option is defined by RFC 2132, Section 9.8.
func GetParameterRequestList(o dhcp4.Options) (OptionCodes, error) {
	v, err := o.Get(dhcp4.OptionParameterRequestList)
	if err != nil {
		return nil, err
	}
	var oc OptionCodes
	return oc, (&oc).UnmarshalBinary(v)
}

// GetMaximumDHCPMessageSize returns the maximum DHCP message size of `o`.
//
// The maximum DHCP message size option is defined by RFC 2132, Section 9.10.
func GetMaximumDHCPMessageSize(o dhcp4.Options) (uint16, error) {
	v, err := o.Get(dhcp4.OptionMaximumDHCPMessageSize)
	if err != nil {
		return 0, err
	}
	var u Uint16
	return uint16(u), (&u).UnmarshalBinary(v)
}
