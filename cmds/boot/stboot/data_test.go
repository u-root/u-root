package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearchBootballFiles(t *testing.T) {
	dir := "testdata/datapartition/bootballs"
	ret, err := searchBootballFiles(dir)
	t.Log(ret)
	require.NotEmpty(t, ret)
	require.NoError(t, err)
}
