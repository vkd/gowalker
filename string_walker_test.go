package gowalker

import (
	"errors"
	"os"
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

func TestEnvStringSource(t *testing.T) {
	envKey := "TEST_SOURCE_ENV"
	envMyValue := "TEST_SOURCE_VALUE"
	err := os.Setenv(envKey, envMyValue)
	if err != nil {
		t.Errorf("Error on set env: %v", err)
	}
	out, ok, err := EnvStringSource.Get(envKey)
	if err != nil {
		t.Errorf("Error on get env value: %v", err)
	}
	if !ok {
		t.Errorf("Wrong bool value: %v (expect: true)", ok)
	}
	if out != envMyValue {
		t.Errorf("Wrong value from source: %v (expect: %v)", out, envMyValue)
	}
}
