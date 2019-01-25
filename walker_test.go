package gowalker

import (
	"errors"
	"reflect"
	"testing"
)

func TestWalkBaseStruct(t *testing.T) {
	var s struct {
		Name string
		ID   *struct {
			Number int
		}
		sex string
	}
	fn := func(value reflect.Value, field reflect.StructField) (bool, error) {
		return false, nil
	}
	expectVisited := map[string]int{
		"Name":   1,
		"ID":     1,
		"Number": 1,
	}
	testWalkStruct(t, &s, WalkerFunc(fn), expectVisited, nil)
}

func TestWalkStructSetted(t *testing.T) {
	var s struct {
		ID *struct {
			Number int
		}
	}
	fn := func(value reflect.Value, field reflect.StructField) (bool, error) {
		if field.Name == "ID" {
			return true, nil
		}
		return false, nil
	}
	expectVisited := map[string]int{
		"ID": 1,
	}
	testWalkStruct(t, &s, WalkerFunc(fn), expectVisited, nil)
}

func TestWalkFuncError(t *testing.T) {
	var s struct {
		Name int
	}
	expectedErr := errors.New("expected error")
	fn := func(value reflect.Value, field reflect.StructField) (bool, error) {
		return false, expectedErr
	}
	expectVisited := map[string]int{
		"Name": 1,
	}
	testWalkStruct(t, &s, WalkerFunc(fn), expectVisited, expectedErr)
}

func testWalkStruct(t *testing.T, value interface{}, w Walker, expectVisited map[string]int, expectedErr error) {
	visitedNames := make(map[string]int)

	err := Walk(value, WalkerFunc(func(value reflect.Value, field reflect.StructField) (bool, error) {
		visitedNames[field.Name]++
		return w.Step(value, field)
	}))
	if err != expectedErr {
		t.Errorf("Not expected error: %v", err)
	}

	if len(visitedNames) != len(expectVisited) {
		t.Errorf("Wrong count of visited names: %#v", visitedNames)
	}
	for k, v := range expectVisited {
		if visitedNames[k] != v {
			t.Errorf("Wrong count of %v: %v (expect: %v)", k, visitedNames[k], v)
		}
	}
}

func TestWalkWrap(t *testing.T) {
	var s struct {
		Document struct {
			ID *struct {
				Number string
			}
		}
	}
	visitedNames := make(map[string]int)
	fn := func(value reflect.Value, field reflect.StructField) (bool, error) {
		visitedNames[field.Name]++
		return false, nil
	}
	expectVisited := map[string]int{
		"Document":           1,
		"Document_ID":        1,
		"Document_ID_Number": 1,
	}

	err := Walk(&s, NewWrapFieldNameWalker(WalkerFunc(fn)))
	if err != nil {
		t.Errorf("Not expected error: %v", err)
	}

	if len(visitedNames) != len(expectVisited) {
		t.Errorf("Wrong count of visited names: %#v", visitedNames)
	}
	for k, v := range expectVisited {
		if visitedNames[k] != v {
			t.Errorf("Wrong count of %v: %v (expect: %v) %#v", k, visitedNames[k], v, visitedNames)
		}
	}
}
