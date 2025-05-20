package testutils

import (
	"reflect"
	"testing"
)

func AssertTrue(t *testing.T, actual bool) {
	t.Helper()
	if !actual {
		t.Error("Expected true but got false :(")
		t.FailNow()
	}
}

func AssertFalse(t *testing.T, actual bool) {
	t.Helper()
	if actual {
		t.Error("Expected false but got true :(")
		t.FailNow()
	}
}

func AssertEqual[T any, U any](t *testing.T, expected T, actual U) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %+v (%T) but got %+v (%T)", expected, expected, actual, actual)
		t.FailNow()
	}
}

func AssertEqualLabel[T any, U any](t *testing.T, label string, expected T, actual U) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %s %+v (%T) but got %+v (%T)", label, expected, expected, actual, actual)
		t.FailNow()
	}
}

func AssertLength[T any](t *testing.T, expectedLength int, array []T) {
	t.Helper()
	if expectedLength != len(array) {
		t.Errorf("Expected a length of %d but got %d", expectedLength, len(array))
		t.FailNow()
	}
}
