package async

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/client4"
	"github.com/stretchr/testify/require"
)

// server creates a server which responds with a predefined response
func serve(ctx context.Context, addr *net.UDPAddr, response *dhcpv4.DHCPv4) error {
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return err
	}
	go func() {
		defer conn.Close()
		oobdata := []byte{}
		buffer := make([]byte, client4.MaxUDPReceivedPacketSize)
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
				_, err = dhcpv4.FromBytes(buffer[:n])
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
	require.Equal(t, c.ReadTimeout, client4.DefaultReadTimeout)
	require.Equal(t, c.ReadTimeout, client4.DefaultWriteTimeout)
}

func TestOpenInvalidAddrFailes(t *testing.T) {
	c := NewClient()
	err := c.Open(512)
	require.Error(t, err)
}

// This test uses port 15438 so please make sure its not used before running
func TestOpenClose(t *testing.T) {
	c := NewClient()
	addr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:15438")
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
	addr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:15438")
	require.NoError(t, err)
	remote, err := net.ResolveUDPAddr("udp4", "127.0.0.1:15439")
	require.NoError(t, err)
	c.ReadTimeout = 50 * time.Millisecond
	c.WriteTimeout = 50 * time.Millisecond
	c.LocalAddr = addr
	c.RemoteAddr = remote
	err = c.Open(512)
	require.NoError(t, err)
	defer c.Close()
	m, err := dhcpv4.New()
	require.NoError(t, err)
	_, err, timeout := c.Send(m).GetOrTimeout(200)
	require.NoError(t, err)
	require.True(t, timeout)
}

// This test uses ports 15438 and 15439 so please make sure they are not used
// before running
func TestSend(t *testing.T) {
	m, err := dhcpv4.New()
	require.NoError(t, err)
	require.NotNil(t, m)

	c := NewClient()
	addr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:15438")
	require.NoError(t, err)
	remote, err := net.ResolveUDPAddr("udp4", "127.0.0.1:15439")
	require.NoError(t, err)
	c.LocalAddr = addr
	c.RemoteAddr = remote

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = serve(ctx, remote, m)
	require.NoError(t, err)

	err = c.Open(16)
	require.NoError(t, err)
	defer c.Close()

	f := c.Send(m)
	response, err, timeout := f.GetOrTimeout(2000)
	r, ok := response.(*dhcpv4.DHCPv4)
	require.True(t, ok)
	require.False(t, timeout)
	require.NoError(t, err)
	require.Equal(t, m.TransactionID, r.TransactionID)
}
