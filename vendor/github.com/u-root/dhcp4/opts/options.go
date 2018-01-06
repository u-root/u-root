// opts implements Options parsing for DHCPv4 options as described in RFC 2132.
//
// Not all options are currently implemented.
package opts

import (
	"io"
	"net"

	"github.com/u-root/dhcp4"
	"github.com/u-root/dhcp4/util"
)

type SubnetMask net.IPMask

func (s SubnetMask) MarshalBinary() ([]byte, error) {
	return []byte(s[:net.IPv4len]), nil
}

func (s *SubnetMask) UnmarshalBinary(p []byte) error {
	if len(p) < net.IPv4len {
		return io.ErrUnexpectedEOF
	}

	*s = make([]byte, net.IPv4len)
	copy(*s, p[:net.IPv4len])
	return nil
}

// RFC 2132, Section 3.3.
func GetSubnetMask(o dhcp4.Options) (SubnetMask, error) {
	v, err := o.Get(dhcp4.OptionSubnetMask)
	if err != nil {
		return nil, err
	}
	var s SubnetMask
	return s, (&s).UnmarshalBinary(v)
}

type DHCPMessageType uint8

const (
	DHCPDiscover DHCPMessageType = 1
	DHCPOffer    DHCPMessageType = 2
	DHCPRequest  DHCPMessageType = 3
	DHCPDecline  DHCPMessageType = 4
	DHCPACK      DHCPMessageType = 5
	DHCPNAK      DHCPMessageType = 6
	DHCPRelease  DHCPMessageType = 7
	DHCPInform   DHCPMessageType = 8
)

func (d DHCPMessageType) MarshalBinary() ([]byte, error) {
	return []byte{byte(d)}, nil
}

func (d *DHCPMessageType) UnmarshalBinary(p []byte) error {
	if len(p) < 1 {
		return io.ErrUnexpectedEOF
	}

	*d = DHCPMessageType(p[0])
	return nil
}

// RFC 2132, Section 9.4.
func GetDHCPMessageType(o dhcp4.Options) (DHCPMessageType, error) {
	v, err := o.Get(dhcp4.OptionDHCPMessageType)
	if err != nil {
		return 0, err
	}

	var d DHCPMessageType
	return d, (&d).UnmarshalBinary(v)
}

type IP net.IP

func (i IP) MarshalBinary() ([]byte, error) {
	return []byte(i[:net.IPv4len]), nil
}

func (i *IP) UnmarshalBinary(p []byte) error {
	if len(p) < net.IPv4len {
		return io.ErrUnexpectedEOF
	}

	*i = make([]byte, net.IPv4len)
	copy(*i, p[:net.IPv4len])
	return nil
}

func GetIP(code dhcp4.OptionCode, o dhcp4.Options) (IP, error) {
	v, err := o.Get(code)
	if err != nil {
		return nil, err
	}
	var ip IP
	return ip, (&ip).UnmarshalBinary(v)
}

// RFC 2132, Section 3.18.
func GetSwapServer(o dhcp4.Options) (IP, error) {
	return GetIP(dhcp4.OptionSwapServer, o)
}

// RFC 2132, Section 5.3.
func GetBroadcastAddress(o dhcp4.Options) (IP, error) {
	return GetIP(dhcp4.OptionBroadcastAddress, o)
}

// RFC 2132, Section 5.7.
func GetRouterSolicitationAddress(o dhcp4.Options) (IP, error) {
	return GetIP(dhcp4.OptionRouterSolicitationAddress, o)
}

// RFC 2132, Section 9.1.
func GetRequestedIPAddress(o dhcp4.Options) (IP, error) {
	return GetIP(dhcp4.OptionRequestedIPAddress, o)
}

// RFC 2132, Section 9.5.
func GetServerIdentifier(o dhcp4.Options) (IP, error) {
	return GetIP(dhcp4.OptionServerIdentifier, o)
}

type IPs []net.IP

func (i IPs) MarshalBinary() ([]byte, error) {
	b := util.NewBuffer(make([]byte, 0, net.IPv4len*len(i)))
	for _, ip := range i {
		b.WriteBytes(ip.To4())
	}
	return b.Data(), nil
}

func (i *IPs) UnmarshalBinary(p []byte) error {
	b := util.NewBuffer(p)
	if b.Len() == 0 || b.Len()%net.IPv4len != 0 {
		return io.ErrUnexpectedEOF
	}

	*i = make([]net.IP, 0, b.Len()/net.IPv4len)
	for b.Len() > 0 {
		ip := make(net.IP, net.IPv4len)
		b.ReadBytes(ip)
		*i = append(*i, ip)
	}
	return nil
}

func GetIPs(code dhcp4.OptionCode, o dhcp4.Options) (IPs, error) {
	v, err := o.Get(code)
	if err != nil {
		return nil, err
	}

	var i IPs
	return i, (&i).UnmarshalBinary(v)
}

// RFC 2132, Section 3.5.
func GetRouters(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionRouters, o)
}

// RFC 2132, Section 3.6.
func GetTimeServers(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionTimeServers, o)
}

// RFC 2132, Section 3.7.
func GetNameServers(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionNameServers, o)
}

// RFC 2132, Section 3.8.
func GetDomainNameServers(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionDomainNameServers, o)
}

// RFC 2132, Section 3.9.
func GetLogServers(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionLogServers, o)
}

// RFC 2132, Section 3.10.
func GetCookieServers(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionCookieServers, o)
}

// RFC 2132, Section 3.11.
func GetLPRServers(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionLPRServers, o)
}

// RFC 2132, Section 3.12.
func GetImpressServers(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionImpressServers, o)
}

// RFC 2132, Section 3.13.
func GetResourceLocationServers(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionResourceLocationServers, o)
}

// RFC 2132, Section 8.2.
func GetNetworkInformationServers(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionNetworkInformationServers, o)
}

// RFC 2132, Section 8.3.
func GetNetworkTimeProtocolServers(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionNetworkTimeProtocolServers, o)
}

// RFC 2132, Section 8.5.
func GetNetBIOSOverTCPIPNameServer(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionNetBIOSOverTCPIPNameServer, o)
}

// RFC 2132, Section 8.6.
func GetNetBIOSOverTCPIPDatagramDistributionServer(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionNetBIOSOverTCPIPDatagramDistributionServer, o)
}

// RFC 2132, Section 8.9.
func GetXWindowSystemFontServer(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionXWindowSystemFontServer, o)
}

// RFC 2132, Section 8.10.
func GetXWindowSystemDisplayManager(o dhcp4.Options) (IPs, error) {
	return GetIPs(dhcp4.OptionXWindowSystemDisplayManager, o)
}

type String string

func (s String) MarshalBinary() ([]byte, error) {
	return []byte(s), nil
}

func (s *String) UnmarshalBinary(b []byte) error {
	*s = String(string(b))
	return nil
}

func GetString(code dhcp4.OptionCode, o dhcp4.Options) (string, error) {
	v, err := o.Get(code)
	if err != nil {
		return "", err
	}
	var s String
	return string(s), (&s).UnmarshalBinary(v)
}

// RFC 2132, Section 3.14.
func GetHostName(o dhcp4.Options) (string, error) {
	return GetString(dhcp4.OptionHostName, o)
}

// RFC 2132, Section 3.16.
func GetMeritDumpFile(o dhcp4.Options) (string, error) {
	return GetString(dhcp4.OptionMeritDumpFile, o)
}

// RFC 2132, Section 3.17.
func GetDomainName(o dhcp4.Options) (string, error) {
	return GetString(dhcp4.OptionDomainName, o)
}

// RFC 2132, Section 3.19.
func GetRootPath(o dhcp4.Options) (string, error) {
	return GetString(dhcp4.OptionRootPath, o)
}

// RFC 2132, Section 3.20.
func GetExtensionsPath(o dhcp4.Options) (string, error) {
	return GetString(dhcp4.OptionExtensionsPath, o)
}

type OptionCodes []dhcp4.OptionCode

func (o OptionCodes) MarshalBinary() ([]byte, error) {
	b := util.NewBuffer(nil)
	for _, code := range o {
		b.Write8(uint8(code))
	}
	return b.Data(), nil
}

func (o *OptionCodes) UnmarshalBinary(p []byte) error {
	b := util.NewBuffer(p)
	*o = make(OptionCodes, 0, b.Len())
	for b.Len() > 0 {
		*o = append(*o, dhcp4.OptionCode(b.Read8()))
	}
	return nil
}

// RFC 2132, Section 9.8.
func GetParameterRequestList(o dhcp4.Options) (OptionCodes, error) {
	v, err := o.Get(dhcp4.OptionParameterRequestList)
	if err != nil {
		return nil, err
	}
	var oc OptionCodes
	return oc, (&oc).UnmarshalBinary(v)
}

type Uint16 uint16

func (u Uint16) MarshalBinary() ([]byte, error) {
	b := util.NewBuffer(nil)
	b.Write16(uint16(u))
	return b.Data(), nil
}

func (u *Uint16) UnmarshalBinary(p []byte) error {
	b := util.NewBuffer(p)
	if b.Len() < 2 {
		return io.ErrUnexpectedEOF
	}
	*u = Uint16(b.Read16())
	return nil
}

func GetMaximumDHCPMessageSize(o dhcp4.Options) (uint16, error) {
	v, err := o.Get(dhcp4.OptionMaximumDHCPMessageSize)
	if err != nil {
		return 0, err
	}
	var u Uint16
	return uint16(u), (&u).UnmarshalBinary(v)
}
