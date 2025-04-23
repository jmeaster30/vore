package testutils

import (
	"fmt"
	"math/rand"
	"testing"
)

func CheckNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func MustPanic(t *testing.T, message string, process func(*testing.T)) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Error(message)
		}
	}()

	process(t)
}

func PseudoUuid(t *testing.T) (string, error) {
	t.Helper()
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}
