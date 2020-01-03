package setters

import (
	"reflect"
	"testing"
	"time"
	"unsafe"
)

var emptyField = reflect.StructField{}

type assertFunc func(t *testing.T, name string)

type testFuncSetValueByString func() (reflect.Value, assertFunc)

func TestSetValueByString(t *testing.T) {
	type testStruct struct {
		name    string
		fn      testFuncSetValueByString
		field   reflect.StructField
		str     string
		wantErr bool
	}
	var tests = []testStruct{
		{"Invalid Kind", func() (reflect.Value, assertFunc) { return reflect.Value{}, nil }, emptyField, "", true},
		{"time.Time layout", assert(time.Date(2020, 11, 6, 0, 0, 0, 0, time.UTC)), reflect.StructField{Tag: `time_layout:"02-2006-01"`}, "06-2020-11", false},
		{"not setted ptr", func() (reflect.Value, assertFunc) {
			var i *int
			return reflect.ValueOf(i), nil
		}, emptyField, "", true},
		{"ptr", func() (reflect.Value, assertFunc) {
			var i *int
			var expect int = 6
			return reflect.ValueOf(&i).Elem(), makeAssertFunc(expect, reflect.ValueOf(&expect))
		}, emptyField, "", false},
	}

	// positive emptyField tests
	for _, tt := range []struct {
		v   interface{}
		str string
	}{
		{true, "true"},

		{int(7), "7"},
		{int8(7), "7"},
		{int16(7), "7"},
		{int32(7), "7"},
		{int64(7), "7"},

		{uint(7), "7"},
		{uint8(7), "7"},
		{uint16(7), "7"},
		{uint32(7), "7"},
		{uint64(7), "7"},

		{float32(7), "7"},
		{float64(7), "7"},

		// time.Duration
		{5 * time.Second, "5s"},

		// array
		{[3]int{7, 0, 0}, "7"},
		{[3]string{"test", "", ""}, "test"},

		// slice
		{[]int{7}, "7"},

		// map
		{map[string]int{"one": 1, "two": 2}, `{"one": 1, "two": 2}`},

		// struct
		{struct{ Name string }{Name: "mike"}, `{"Name": "mike"}`},

		// time.Time
		{time.Date(2020, 1, 12, 14, 0, 5, 0, time.UTC), "2020-01-12T14:00:05Z"},
	} {
		name := reflect.TypeOf(tt.v).String()
		tests = append(tests, testStruct{"positive: " + name, assert(tt.v), emptyField, tt.str, false})
	}

	// empty string tests
	for _, tt := range []struct {
		v interface{}
	}{
		{false},
		{int(0)},
		{uint(0)},
		{float32(0)},
	} {
		name := reflect.TypeOf(tt.v).String()
		tests = append(tests, testStruct{"empty string: " + name, assert(tt.v), emptyField, "", false})
	}

	// negative emptyField tests
	for _, tt := range []struct {
		v   interface{}
		str string
	}{
		{false, "wrongValue"},
		{int(0), "wrongValue"},
		{uint(0), "wrongValue"},
		{float32(0), "wrongValue"},
		{time.Duration(0), "wrongValue"},
		{[]int{}, "wrongValue"},
		{time.Time{}, "wrongValue"},
	} {
		name := reflect.TypeOf(tt.v).String()
		tests = append(tests, testStruct{"wrong: " + name, assert(tt.v), emptyField, tt.str, true})
	}

	// not support tests
	for _, tt := range []struct {
		v interface{}
	}{
		{complex64(0)},
		{complex128(0)},
		{chan struct{}(nil)},
		{unsafe.Pointer(nil)},
	} {
		name := reflect.TypeOf(tt.v).String()
		tests = append(tests, testStruct{"not support: " + name, assert(tt.v), emptyField, "", true})
	}

	// TESTS
	for _, tt := range tests {
		v, assertFn := tt.fn()
		if err := SetString(v, tt.field, tt.str); (err != nil) != tt.wantErr {
			t.Errorf("%q. SetValueByString() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
		if assertFn != nil {
			assertFn(t, tt.name)
		}
	}
}

func assert(expect interface{}) testFuncSetValueByString {
	return assertDefaultZero(expect)
}

func assertDefaultZero(expect interface{}) testFuncSetValueByString {
	return func() (reflect.Value, assertFunc) {
		tExpect := reflect.TypeOf(expect)
		value := reflect.New(tExpect)
		return value.Elem(), makeAssertFunc(expect, value)
	}
}

func makeAssertFunc(expect interface{}, valuePtr reflect.Value) assertFunc {
	return func(t *testing.T, name string) {
		val := valuePtr.Elem().Interface()
		if !isEquals(expect, val) {
			t.Errorf("%q: wrong %s value: %v (expect: %v)", name, reflect.TypeOf(expect).String(), val, expect)
		}
	}
}

func isEquals(expect, val interface{}) bool {
	switch reflect.ValueOf(expect).Kind() {
	case reflect.Map:
		return isMapEquals(reflect.ValueOf(expect), reflect.ValueOf(val))
	case reflect.Slice:
		return isSliceEquals(reflect.ValueOf(expect), reflect.ValueOf(val))
	default:
		return val == expect
	}
}

func isMapEquals(expect reflect.Value, val reflect.Value) bool {
	if expect.Len() != val.Len() {
		return false
	}
	for _, key := range expect.MapKeys() {
		v := val.MapIndex(key)
		if !v.IsValid() {
			return false
		}
		if !isEquals(expect.MapIndex(key).Interface(), v.Interface()) {
			return false
		}
	}
	return true
}

func isSliceEquals(expect, val reflect.Value) bool {
	if expect.Len() != val.Len() {
		return false
	}
	for i := 0; i < expect.Len(); i++ {
		e := expect.Index(i).Interface()
		v := val.Index(i).Interface()
		if !isEquals(e, v) {
			return false
		}
	}
	return true
}
