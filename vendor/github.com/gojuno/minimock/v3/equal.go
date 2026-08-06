package minimock

import (
	"context"
	"reflect"
	"unsafe"

	"github.com/davecgh/go-spew/spew"
	"github.com/pmezard/go-difflib/difflib"
)

var dumpConf = spew.ConfigState{
	Indent:                  " ",
	DisablePointerAddresses: true,
	SortKeys:                true,
}

type anyContext struct {
	context.Context
}

var AnyContext = anyContext{}

// Equal returns true if a equals b
func Equal(a, b interface{}) bool {
	if a == nil && b == nil {
		return a == b
	}

	if reflect.TypeOf(a).Kind() == reflect.Struct {
		ap := copyValue(a)
		bp := copyValue(b)

		// for every field in a
		for i := 0; i < reflect.TypeOf(a).NumField(); i++ {
			aFieldValue := unexported(ap.Field(i))
			bFieldValue := unexported(bp.Field(i))

			if checkAnyContext(aFieldValue, bFieldValue) {
				continue
			}

			if !reflect.DeepEqual(aFieldValue, bFieldValue) {
				return false
			}
		}

		return true
	}

	return reflect.DeepEqual(a, b)
}

// Diff returns unified diff of the textual representations of e and a
func Diff(e, a interface{}) string {
	if e == nil || a == nil {
		return ""
	}

	t := reflect.TypeOf(e)
	k := t.Kind()

	if reflect.TypeOf(a) != t {
		return ""
	}

	initialKind := k
	if k == reflect.Ptr {
		t = t.Elem()
		k = t.Kind()
	}

	if k != reflect.Array && k != reflect.Map && k != reflect.Slice && k != reflect.Struct {
		return ""
	}

	if initialKind == reflect.Struct {
		a = setAnyContext(e, a)
	}

	es := dumpConf.Sdump(e)
	as := dumpConf.Sdump(a)

	diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        difflib.SplitLines(es),
		B:        difflib.SplitLines(as),
		Context:  1,
		FromFile: "Expected params",
		ToFile:   "Actual params",
	})

	if err != nil {
		panic(err)
	}

	return "\n\nDiff:\n" + diff
}

func unexported(field reflect.Value) interface{} {
	return unexportedVal(field).Interface()
}

func unexportedVal(field reflect.Value) reflect.Value {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
}
