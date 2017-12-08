package find

import (
	"testing"
)

func TestSimple(t *testing.T) {
	f, err := New()
	if err != nil {
		t.Fatal(err)
	}
	go f.Find()
	for {
		select {
		case n := <-f.Name:
			t.Logf("%v\n", n)
		case n := <-f.Err:
			t.Logf("%v\n", n)
		default:
		}
	}

}
