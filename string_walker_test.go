package gowalker

import (
	"errors"
	"testing"
)

func TestStringGetValue(t *testing.T) {
	errTest := errors.New("my test error")
	failSource := StringSourceFunc(func(key string) (string, bool, error) { return "", false, errTest })
	_, _, err := StringGetValue("tag", failSource, emptyField)
	if err != errTest {
		t.Errorf("Wrong error from source: %v", err)
	}
}
