package dhcpv6

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var ErrDummy = errors.New("dummy error")

type MatchingAddressTestSuite struct {
	suite.Suite
	m mock.Mock

	ips   []net.IP
	addrs []net.Addr
}

func (s *MatchingAddressTestSuite) InterfaceAddresses(name string) ([]net.Addr, error) {
	args := s.m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if ifaddrs, ok := args.Get(0).([]net.Addr); ok {
		return ifaddrs, args.Error(1)
	}
	panic(fmt.Sprintf("assert: arguments: InterfaceAddresses(0) failed because object wasn't correct type: %v", args.Get(0)))
}

func (s *MatchingAddressTestSuite) Match(ip net.IP) bool {
	args := s.m.Called(ip)
	return args.Bool(0)
}

func (s *MatchingAddressTestSuite) SetupTest() {
	InterfaceAddresses = s.InterfaceAddresses
	s.ips = []net.IP{
		net.ParseIP("2401:db00:3020:70e1:face:0:7e:0"),
		net.ParseIP("2803:6080:890c:847e::1"),
		net.ParseIP("fe80::4a57:ddff:fe04:d8e9"),
	}
	s.addrs = []net.Addr{}
	for _, ip := range s.ips {
		s.addrs = append(s.addrs, &net.IPNet{IP: ip})
	}
}

func (s *MatchingAddressTestSuite) TestGetMatchingAddr() {
	// Check if error from InterfaceAddresses immediately returns error
	s.m.On("InterfaceAddresses", "eth0").Return(nil, ErrDummy).Once()
	_, err := getMatchingAddr("eth0", s.Match)
	s.Assert().Equal(ErrDummy, err)
	s.m.AssertExpectations(s.T())
	// Check if the looping is stopped after finding a matching address
	s.m.On("InterfaceAddresses", "eth0").Return(s.addrs, nil).Once()
	s.m.On("Match", s.ips[0]).Return(false).Once()
	s.m.On("Match", s.ips[1]).Return(true).Once()
	ip, err := getMatchingAddr("eth0", s.Match)
	s.Require().NoError(err)
	s.Assert().Equal(s.ips[1], ip)
	s.m.AssertExpectations(s.T())
	// Check if the looping skips not matching addresses
	s.m.On("InterfaceAddresses", "eth0").Return(s.addrs, nil).Once()
	s.m.On("Match", s.ips[0]).Return(false).Once()
	s.m.On("Match", s.ips[1]).Return(false).Once()
	s.m.On("Match", s.ips[2]).Return(true).Once()
	ip, err = getMatchingAddr("eth0", s.Match)
	s.Require().NoError(err)
	s.Assert().Equal(s.ips[2], ip)
	s.m.AssertExpectations(s.T())
	// Check if the error is returned if no matching address is found
	s.m.On("InterfaceAddresses", "eth0").Return(s.addrs, nil).Once()
	s.m.On("Match", s.ips[0]).Return(false).Once()
	s.m.On("Match", s.ips[1]).Return(false).Once()
	s.m.On("Match", s.ips[2]).Return(false).Once()
	_, err = getMatchingAddr("eth0", s.Match)
	s.Assert().EqualError(err, "no matching address found for interface eth0")
	s.m.AssertExpectations(s.T())
}

func (s *MatchingAddressTestSuite) TestGetLinkLocalAddr() {
	s.m.On("InterfaceAddresses", "eth0").Return(s.addrs, nil).Once()
	ip, err := GetLinkLocalAddr("eth0")
	s.Require().NoError(err)
	s.Assert().Equal(s.ips[2], ip)
	s.m.AssertExpectations(s.T())
}

func (s *MatchingAddressTestSuite) TestGetGlobalAddr() {
	s.m.On("InterfaceAddresses", "eth0").Return(s.addrs, nil).Once()
	ip, err := GetGlobalAddr("eth0")
	s.Require().NoError(err)
	s.Assert().Equal(s.ips[0], ip)
	s.m.AssertExpectations(s.T())
}

func TestMatchingAddressTestSuite(t *testing.T) {
	suite.Run(t, new(MatchingAddressTestSuite))
}

func Test_ExtractMAC(t *testing.T) {
	//SOLICIT message wrapped in Relay-Forw
	var relayForwBytesDuidUUID = []byte{
		0x0c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xfe, 0x80, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x26, 0x8a, 0x07, 0xff, 0xfe, 0x56,
		0xdc, 0xa4, 0x00, 0x12, 0x00, 0x06, 0x24, 0x8a,
		0x07, 0x56, 0xdc, 0xa4, 0x00, 0x09, 0x00, 0x5a,
		0x06, 0x7d, 0x9b, 0xca, 0x00, 0x01, 0x00, 0x12,
		0x00, 0x04, 0xb7, 0xfd, 0x0a, 0x8c, 0x1b, 0x14,
		0x10, 0xaa, 0xeb, 0x0a, 0x5b, 0x3f, 0xe8, 0x9d,
		0x0f, 0x56, 0x00, 0x06, 0x00, 0x0a, 0x00, 0x17,
		0x00, 0x18, 0x00, 0x17, 0x00, 0x18, 0x00, 0x01,
		0x00, 0x08, 0x00, 0x02, 0xff, 0xff, 0x00, 0x03,
		0x00, 0x28, 0x07, 0x56, 0xdc, 0xa4, 0x00, 0x00,
		0x0e, 0x10, 0x00, 0x00, 0x15, 0x18, 0x00, 0x05,
		0x00, 0x18, 0x26, 0x20, 0x01, 0x0d, 0xc0, 0x82,
		0x90, 0x63, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xaf, 0xa0, 0x00, 0x00, 0x1c, 0x20, 0x00, 0x00,
		0x1d, 0x4c}
	packet, err := FromBytes(relayForwBytesDuidUUID)
	require.NoError(t, err)
	mac, err := ExtractMAC(packet)
	require.NoError(t, err)
	require.Equal(t, mac.String(), "24:8a:07:56:dc:a4")

	// MAC extracted from DUID
	duid := Duid{
		Type:          DUID_LL,
		HwType:        iana.HWTypeEthernet,
		LinkLayerAddr: []byte{0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa},
	}
	solicit, err := NewMessage(WithClientID(duid))
	require.NoError(t, err)
	relay, err := EncapsulateRelay(solicit, MessageTypeRelayForward, net.IPv6zero, net.IPv6zero)
	require.NoError(t, err)
	mac, err = ExtractMAC(relay)
	require.NoError(t, err)
	require.Equal(t, mac.String(), "aa:aa:aa:aa:aa:aa")

	// no client ID
	solicit, err = NewMessage()
	require.NoError(t, err)
	_, err = ExtractMAC(solicit)
	require.Error(t, err)

	// DUID is not DuidLL or DuidLLT
	duid = Duid{}
	solicit, err = NewMessage(WithClientID(duid))
	require.NoError(t, err)
	_, err = ExtractMAC(solicit)
	require.Error(t, err)
}
