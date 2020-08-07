package async

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/insomniacslk/dhcp/dhcpv6/client6"
	"github.com/stretchr/testify/require"
)

const retries = 5

// solicit creates new solicit based on the mac address
func solicit(input string) (*dhcpv6.Message, error) {
	mac, err := net.ParseMAC(input)
	if err != nil {
		return nil, err
	}
	return dhcpv6.NewSolicit(mac)
}

// server creates a server which responds with a predefined response
func serve(ctx context.Context, addr *net.UDPAddr, response dhcpv6.DHCPv6) error {
	conn, err := net.ListenUDP("udp6", addr)
	if err != nil {
		return err
	}
	go func() {
		defer conn.Close()
		oobdata := []byte{}
		buffer := make([]byte, client6.MaxUDPReceivedPacketSize)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := conn.SetReadDeadline(time.Now().Add(1 * time.Second)); err != nil {
					panic(err)
				}
				n, _, _, src, err := conn.ReadMsgUDP(buffer, oobdata)
				if err != nil {
					continue
				}
				_, err = dhcpv6.FromBytes(buffer[:n])
				if err != nil {
					continue
				}
				if err := conn.SetWriteDeadline(time.Now().Add(1 * time.Second)); err != nil {
					panic(err)
				}
				_, err = conn.WriteTo(response.ToBytes(), src)
				if err != nil {
					continue
				}
			}
		}
	}()
	return nil
}

func TestNewClient(t *testing.T) {
	c := NewClient()
	require.NotNil(t, c)
	require.Equal(t, c.ReadTimeout, client6.DefaultReadTimeout)
	require.Equal(t, c.ReadTimeout, client6.DefaultWriteTimeout)
}

func TestOpenInvalidAddrFailes(t *testing.T) {
	c := NewClient()
	err := c.Open(512)
	require.Error(t, err)
}

// This test uses port 15438 so please make sure its not used before running
func TestOpenClose(t *testing.T) {
	c := NewClient()
	addr, err := net.ResolveUDPAddr("udp6", ":15438")
	require.NoError(t, err)
	c.LocalAddr = addr
	err = c.Open(512)
	require.NoError(t, err)
	defer c.Close()
}

// This test uses ports 15438 and 15439 so please make sure they are not used
// before running
func TestSendTimeout(t *testing.T) {
	c := NewClient()
	addr, err := net.ResolveUDPAddr("udp6", ":15438")
	require.NoError(t, err)
	remote, err := net.ResolveUDPAddr("udp6", ":15439")
	require.NoError(t, err)
	c.ReadTimeout = 50 * time.Millisecond
	c.WriteTimeout = 50 * time.Millisecond
	c.LocalAddr = addr
	c.RemoteAddr = remote
	err = c.Open(512)
	require.NoError(t, err)
	defer c.Close()
	m, err := dhcpv6.NewMessage()
	require.NoError(t, err)
	_, err, timeout := c.Send(m).GetOrTimeout(200)
	require.NoError(t, err)
	require.True(t, timeout)
}

// This test uses ports 15438 and 15439 so please make sure they are not used
// before running
func TestSend(t *testing.T) {
	s, err := solicit("c8:6c:2c:47:96:fd")
	require.NoError(t, err)
	require.NotNil(t, s)

	a, err := dhcpv6.NewAdvertiseFromSolicit(s)
	require.NoError(t, err)
	require.NotNil(t, a)

	c := NewClient()
	addr, err := net.ResolveUDPAddr("udp6", ":15438")
	require.NoError(t, err)
	remote, err := net.ResolveUDPAddr("udp6", ":15439")
	require.NoError(t, err)
	c.LocalAddr = addr
	c.RemoteAddr = remote

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = serve(ctx, remote, a)
	require.NoError(t, err)

	err = c.Open(16)
	require.NoError(t, err)
	defer c.Close()

	f := c.Send(s)

	var passed bool
	for i := 0; i < retries; i++ {
		response, err, timeout := f.GetOrTimeout(1000)
		if timeout {
			continue
		}
		passed = true
		require.NoError(t, err)
		require.Equal(t, a, response)
	}
	require.True(t, passed, "All attempts to TestSend timed out")
}
