package checker

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type someError struct{}

func (s someError) Error() string {
	return "some error"
}

func TestMakeFunctionThatReturnsErrorReturnTypedNilVariadic(t *testing.T) {
	var typedNil *someError
	fn, err := makeFunctionThatReturnsError(func(e error, errs ...error) error { return e }, typedNil)
	require.NoError(t, err)
	assert.Nil(t, fn())
}

func TestMakeFunctionThatReturnsErrorReturnUntypedNilVariadic(t *testing.T) {
	fn, err := makeFunctionThatReturnsError(func(e error, errs ...error) error { return e }, nil)
	require.NoError(t, err)
	assert.Nil(t, fn())
}

func TestMakeFunctionThatReturnsErrorReturnErr(t *testing.T) {
	fn, err := makeFunctionThatReturnsError(func(e error) error { return e }, errors.New("err"))
	require.NoError(t, err)
	assert.Error(t, fn())
}
