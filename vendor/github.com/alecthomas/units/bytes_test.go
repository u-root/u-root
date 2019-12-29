package units

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase2BytesString(t *testing.T) {
	assert.Equal(t, Base2Bytes(0).String(), "0B")
	assert.Equal(t, Base2Bytes(1025).String(), "1KiB1B")
	assert.Equal(t, Base2Bytes(1048577).String(), "1MiB1B")
}

func TestParseBase2Bytes(t *testing.T) {
	n, err := ParseBase2Bytes("0B")
	assert.NoError(t, err)
	assert.Equal(t, 0, int(n))
	n, err = ParseBase2Bytes("1kB")
	assert.Error(t, err)
	n, err = ParseBase2Bytes("1KB")
	assert.NoError(t, err)
	assert.Equal(t, 1024, int(n))
	n, err = ParseBase2Bytes("1MB1KB25B")
	assert.NoError(t, err)
	assert.Equal(t, 1049625, int(n))
	n, err = ParseBase2Bytes("1.5MB")
	assert.NoError(t, err)
	assert.Equal(t, 1572864, int(n))

	n, err = ParseBase2Bytes("1kiB")
	assert.Error(t, err)
	n, err = ParseBase2Bytes("1KiB")
	assert.NoError(t, err)
	assert.Equal(t, 1024, int(n))
	n, err = ParseBase2Bytes("1MiB1KiB25B")
	assert.NoError(t, err)
	assert.Equal(t, 1049625, int(n))
	n, err = ParseBase2Bytes("1.5MiB")
	assert.NoError(t, err)
	assert.Equal(t, 1572864, int(n))
}

func TestMetricBytesString(t *testing.T) {
	assert.Equal(t, MetricBytes(0).String(), "0B")
	// TODO: SI standard prefix is lowercase "kB"
	assert.Equal(t, MetricBytes(1001).String(), "1KB1B")
	assert.Equal(t, MetricBytes(1001025).String(), "1MB1KB25B")
}

func TestParseMetricBytes(t *testing.T) {
	n, err := ParseMetricBytes("0B")
	assert.NoError(t, err)
	assert.Equal(t, 0, int(n))
	n, err = ParseMetricBytes("1kB")
	assert.NoError(t, err)
	assert.Equal(t, 1000, int(n))
	n, err = ParseMetricBytes("1KB1B")
	assert.NoError(t, err)
	assert.Equal(t, 1001, int(n))
	n, err = ParseMetricBytes("1MB1KB25B")
	assert.NoError(t, err)
	assert.Equal(t, 1001025, int(n))
	n, err = ParseMetricBytes("1.5MB")
	assert.NoError(t, err)
	assert.Equal(t, 1500000, int(n))
}

func TestParseStrictBytes(t *testing.T) {
	n, err := ParseStrictBytes("0B")
	assert.NoError(t, err)
	assert.Equal(t, 0, int(n))

	n, err = ParseStrictBytes("1kiB")
	assert.Error(t, err)
	n, err = ParseStrictBytes("1KiB")
	assert.NoError(t, err)
	assert.Equal(t, 1024, int(n))
	n, err = ParseStrictBytes("1MiB1KiB25B")
	assert.NoError(t, err)
	assert.Equal(t, 1049625, int(n))
	n, err = ParseStrictBytes("1.5MiB")
	assert.NoError(t, err)
	assert.Equal(t, 1572864, int(n))

	n, err = ParseStrictBytes("0B")
	assert.NoError(t, err)
	assert.Equal(t, 0, int(n))
	n, err = ParseStrictBytes("1kB")
	assert.NoError(t, err)
	assert.Equal(t, 1000, int(n))
	n, err = ParseStrictBytes("1KB1B")
	assert.NoError(t, err)
	assert.Equal(t, 1001, int(n))
	n, err = ParseStrictBytes("1MB1KB25B")
	assert.NoError(t, err)
	assert.Equal(t, 1001025, int(n))
	n, err = ParseStrictBytes("1.5MB")
	assert.NoError(t, err)
	assert.Equal(t, 1500000, int(n))
}
