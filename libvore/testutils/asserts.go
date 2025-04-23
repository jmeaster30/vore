package testutils

import (
	"reflect"
	"testing"
)

func AssertEqual[T any, U any](t *testing.T, expected T, actual U) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %+v (%T) but got %+v (%T)", expected, expected, actual, actual)
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
