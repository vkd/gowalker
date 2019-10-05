package setters

import "reflect"

// SetValueBySliceOfString - set value by slice of strings
//
// Support multi values only: Slice, Array
// other kindes setted by string (first value from slice)
func SetValueBySliceOfString(value reflect.Value, field reflect.StructField, strs []string) error {
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
		return SetValueByString(value, field, str)
	}
}