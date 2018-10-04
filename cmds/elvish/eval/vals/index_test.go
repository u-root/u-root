package vals

import (
	"testing"

	"github.com/u-root/u-root/cmds/elvish/tt"
)

var (
	li0 = EmptyList
	li4 = MakeList("foo", "bar", "lorem", "ipsum")
	m   = MakeMapFromKV("foo", "bar", "lorem", "ipsum")
)

var indexTests = tt.Table{
	// List indices
	// ============

	// Simple indicies: 0 <= i < n.
	Args(li4, "0").Rets("foo", nil),
	Args(li4, "3").Rets("ipsum", nil),
	Args(li0, "0").Rets(any, anyError),
	Args(li4, "4").Rets(any, anyError),
	Args(li4, "5").Rets(any, anyError),
	// Negative indices: -n <= i < 0.
	Args(li4, "-1").Rets("ipsum", nil),
	Args(li4, "-4").Rets("foo", nil),
	Args(li4, "-5").Rets(any, anyError), // Too negative.
	// Decimal indicies are not allowed even if the value is an integer.
	Args(li4, "0.0").Rets(any, anyError),

	// Slice indicies: 0 <= i <= j <= n.
	Args(li4, "1:3").Rets(eq(MakeList("bar", "lorem")), nil),
	Args(li4, "3:4").Rets(eq(MakeList("ipsum")), nil),
	Args(li4, "4:4").Rets(eq(EmptyList), nil), // i == j == n is allowed.
	// i defaults to 0
	Args(li4, ":2").Rets(eq(MakeList("foo", "bar")), nil),
	Args(li4, ":-1").Rets(eq(MakeList("foo", "bar", "lorem")), nil),
	// j defaults to n
	Args(li4, "3:").Rets(eq(MakeList("ipsum")), nil),
	Args(li4, "-2:").Rets(eq(MakeList("lorem", "ipsum")), nil),
	// Both indices can be omitted.
	Args(li0, ":").Rets(eq(li0), nil),
	Args(li4, ":").Rets(eq(li4), nil),

	// Malformed list indices.
	Args(li4, "a").Rets(any, anyError),
	Args(li4, "1:3:2").Rets(any, anyError),

	// Map indicies
	Args(m, "foo").Rets("bar", nil),
	Args(m, "bad").Rets(any, anyError),
}

func TestIndex(t *testing.T) {
	tt.Test(t, tt.Fn("Index", Index), indexTests)
}
