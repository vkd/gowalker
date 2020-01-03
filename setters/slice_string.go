package setters

import "reflect"

// SetSliceStrings - set value by slice of strings
//
// Support multi values only: Slice, Array
// other kindes set by string (first value from slice)
func SetSliceStrings(value reflect.Value, field reflect.StructField, strs []string) error {
	switch value.Kind() {
	case reflect.Slice:
		return sliceStringSetter(value, field, strs)
	case reflect.Array:
		return arrayStringSetter(value, field, strs)
	default:
		var str string
		if len(strs) > 0 {
			str = strs[0]
		}
		return SetString(value, field, str)
	}
}
