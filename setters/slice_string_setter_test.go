package setters

import (
	"reflect"
	"testing"
)

func TestSetValueBySliceOfString(t *testing.T) {
	type testStruct struct {
		name    string
		fn      testFuncSetValueByString
		field   reflect.StructField
		strs    []string
		wantErr bool
	}
	var tests = []testStruct{
		{"full slice", assert([]int{1, 4, 7}), emptyField, []string{"1", "4", "7"}, false},
		{"full array", assert([3]bool{false, true, false}), emptyField, []string{"", "true", "0"}, false},

		{"single int", assert(int(9)), emptyField, []string{"9", "1"}, false},
	}

	for _, tt := range tests {
		v, assertFn := tt.fn()
		if err := SetValueBySliceOfString(v, tt.field, tt.strs); (err != nil) != tt.wantErr {
			t.Errorf("%q. SetValueBySliceOfString() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
		if assertFn != nil {
			assertFn(t, tt.name)
		}
	}

	// type args struct {
	// 	value reflect.Value
	// 	field reflect.StructField
	// 	strs  []string
	// }
	// tests := []struct {
	// 	name    string
	// 	args    args
	// 	wantErr bool
	// }{
	// 	// TODO: Add test cases.
	// }
	// for _, tt := range tests {
	// 	if err := SetValueBySliceOfString(tt.args.value, tt.args.field, tt.args.strs); (err != nil) != tt.wantErr {
	// 		t.Errorf("%q. SetValueBySliceOfString() error = %v, wantErr %v", tt.name, err, tt.wantErr)
	// 	}
	// }
}
