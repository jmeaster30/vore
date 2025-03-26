package testutils

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
)

func CheckNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		fmt.Println(err.Error())
		t.Error(err)
		t.FailNow()
	}
}

func MustPanic(t *testing.T, message string, process func(*testing.T)) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(message)
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

func GetTestingFilename(t *testing.T, touchFile bool) string {
	filename, err := PseudoUuid(t)
	CheckNoError(t, err)

	fullFilename := filename + ".txt"

	if touchFile {
		if _, err := os.Stat(fullFilename); err == nil {
			t.Errorf("Randomized testing file '%s' already exists. I don't want to overwrite it :(", fullFilename)
		}

		file, err := os.OpenFile(fullFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(0666))
		CheckNoError(t, err)

		err = file.Close()
		CheckNoError(t, err)
	}

	return fullFilename
}

func RemoveTestingFile(t *testing.T, filename string) {
	err := os.Remove(filename)
	CheckNoError(t, err)
}
