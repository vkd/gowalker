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
	expectVisited := map[string]int{
		"Name":   1,
		"ID":     1,
		"Number": 1,
	}
	testWalkStruct(t, &s, nil, expectVisited, nil)
}

func TestWalkStructSet(t *testing.T) {
	var s struct {
		ID *struct {
			Number int
		}
	}
	fn := func(value reflect.Value, field reflect.StructField, _ Fields) (bool, error) {
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
	fn := func(value reflect.Value, field reflect.StructField, _ Fields) (bool, error) {
		return false, expectedErr
	}
	expectVisited := map[string]int{
		"Name": 1,
	}
	testWalkStruct(t, &s, WalkerFunc(fn), expectVisited, expectedErr)
}

func testWalkStruct(t *testing.T, value interface{}, w Walker, expectVisited map[string]int, expectedErr error) {
	visitedNames := make(map[string]int)

	err := Walk(value, make(Fields, 0, 4), WalkerFunc(func(value reflect.Value, field reflect.StructField, fs Fields) (bool, error) {
		visitedNames[field.Name]++
		if w == nil {
			return false, nil
		}
		return w.Step(value, field, fs)
	}))
	if !errors.Is(err, expectedErr) {
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
	visitedNames := visitedNames{visits: make(map[string]int)}
	expectVisited := map[string]int{
		"Document":           1,
		"Document.ID":        1,
		"Document.ID.Number": 1,
	}

	err := Walk(&s, make(Fields, 0, 3), visitedNames)
	if err != nil {
		t.Errorf("Not expected error: %v", err)
	}

	if len(visitedNames.visits) != len(expectVisited) {
		t.Errorf("Wrong count of visited names: %#v", visitedNames)
	}
	for k, v := range expectVisited {
		if visitedNames.visits[k] != v {
			t.Errorf("Wrong count of %v: %v (expect: %v) %#v", k, visitedNames.visits[k], v, visitedNames)
		}
	}
}

type visitedNames struct {
	visits map[string]int
}

func (v visitedNames) Step(value reflect.Value, field reflect.StructField, fs Fields) (bool, error) {
	key := StructFieldNamer.Key(fs)
	v.visits[key]++
	return false, nil
}
