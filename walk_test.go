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
		sex        string
		StringData testTypeSetStringer
		Data       testType
	}
	expectVisited := map[string]int{
		"Name":       1,
		"Number":     1,
		"StringData": 1,
		"A":          1,
		"B":          1,
		"Data":       1,
	}
	testWalkStruct(t, &s, nil, expectVisited, nil)
}

type testTypeSetStringer struct {
	A string
	B int
}

func (t *testTypeSetStringer) SetString(s string) error { return nil }

type testType struct {
	c int //nolint:unused // ignored by walker
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

func testWalkStruct(t *testing.T, value interface{}, w Walker, expectVisited map[string]int, expectedErr error, opts ...Option) {
	visitedNames := make(map[string]int)

	err := Walk(value, make(Fields, 0, 4), WalkerFunc(func(value reflect.Value, field reflect.StructField, fs Fields) (bool, error) {
		visitedNames[field.Name]++
		if w == nil {
			return false, nil
		}
		return w.Step(value, field, fs)
	}), opts...)
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
	key := StructFieldNamer.Key(fs.Names())
	v.visits[key]++
	return false, nil
}
