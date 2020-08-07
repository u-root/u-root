package server4

import(
	"log"
	"os"
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/stretchr/testify/require"
)

func TestEmptyLogger(t *testing.T) {
	l := EmptyLogger{}
	msg, err := dhcpv4.New()
	require.Nil(t, err)
	l.Printf("test")
	l.PrintMessage("prefix", msg)
}

func TestShortSummaryLogger(t *testing.T) {
	l := ShortSummaryLogger{
		Printfer: log.New(os.Stderr, "[dhcpv4] ", log.LstdFlags),
	}
	msg, err := dhcpv4.New()
	require.Nil(t, err)
	require.NotNil(t, msg)
	l.Printf("test")
	l.PrintMessage("prefix", msg)
}

func TestDebugLogger(t *testing.T) {
	l := DebugLogger{
		Printfer: log.New(os.Stderr, "[dhcpv4] ", log.LstdFlags),
	}
	msg, err := dhcpv4.New()
	require.Nil(t, err)
	require.NotNil(t, msg)
	l.Printf("test")
	l.PrintMessage("prefix", msg)
}
