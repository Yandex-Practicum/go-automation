package common

import "reflect"

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	rv := reflect.ValueOf(i)
	if !rv.IsValid() {
		return true
	}

	switch rv.Kind() {
	case reflect.Ptr, reflect.Map:
		return rv.IsNil()
	}
	return false
}
