package gowalker

import (
	"errors"
	"testing"
)

func TestSliceStringGetValue(t *testing.T) {
	errTest := errors.New("test error")
	errSource := SliceStringsSourceFunc(func(key string) ([]string, bool, error) {
		return nil, false, errTest
	})
	_, _, err := SliceStringGetValue("tag", errSource, emptyField)
	if err != errTest {
		t.Errorf("Wrong error value from source: %v", err)
	}
}

func TestSliceStringsSourceFunc_StringSource(t *testing.T) {
	expect := "test1"
	m := map[string][]string{
		"key": {expect, "test2"},
	}
	source := SliceStringsSourceFunc(func(key string) ([]string, bool, error) {
		return SliceStringsSourceMapStrings(m).Get(key)
	})
	stringSource := source.StringSource()
	s, ok, err := stringSource.Get("key")
	if err != nil {
		t.Errorf("Error on get string from slice source: %v", err)
	}
	if !ok {
		t.Errorf("Wrong exists from slice source")
	}
	if s != expect {
		t.Errorf("Wrong value from source: %v (expect: %v)", s, expect)
	}

	s, ok, err = stringSource.Get("notfound")
	if err != nil {
		t.Errorf("Erro ron get string from slice source: %v", err)
	}
	if ok {
		t.Errorf("Get unexpected value: %v", s)
	}
}
