package printf

import (
	"bytes"
	"reflect"
	"testing"
)

func TestReadFormat(t *testing.T) {

	type testcase struct {
		str    string
		expect *format
	}
	cases := []testcase{
		{"0*.*f", &format{padZero: true, width: -1, precision: -1, length: 0, specifier: 'f'}},
	}
	for _, c := range cases {
		fr := bytes.NewBuffer([]byte(c.str))
		f, err := readFormat(fr)
		if err != nil {
			t.Errorf("read format: %s", err)
		}
		if !reflect.DeepEqual(f, c.expect) {
			t.Errorf("want %+v got %+v", f, c.expect)
		}
	}
}
