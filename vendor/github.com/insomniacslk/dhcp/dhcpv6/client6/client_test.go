package client6

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	c := NewClient()
	require.NotNil(t, c)
	require.Equal(t, DefaultReadTimeout, c.ReadTimeout)
	require.Equal(t, DefaultWriteTimeout, c.WriteTimeout)
}
