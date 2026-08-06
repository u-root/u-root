package minimock

import (
	"context"
	"reflect"
)

func copyValue(value interface{}) reflect.Value {
	rValue := reflect.ValueOf(value)
	newValue := reflect.New(rValue.Type()).Elem()
	newValue.Set(rValue)

	return newValue
}

func checkAnyContext(eFieldValue, aFieldValue interface{}) bool {
	_, ok := eFieldValue.(anyContext)
	if ok {
		if ctx, ok := aFieldValue.(context.Context); ok && ctx != nil {
			return true
		}
	}

	return false
}

func setAnyContext(src, dst interface{}) interface{} {
	srcp := copyValue(src)
	dstp := copyValue(dst)

	for i := 0; i < srcp.NumField(); i++ {
		srcFieldValue := unexportedVal(srcp.Field(i))
		dstFieldValue := unexportedVal(dstp.Field(i))

		if checkAnyContext(srcFieldValue.Interface(), dstFieldValue.Interface()) {
			// we set context field to anyContext because
			// we don't want to display diff between two contexts in case
			// of anyContext
			dstFieldValue.Set(srcFieldValue)
		}
	}

	return dstp.Interface()
}
