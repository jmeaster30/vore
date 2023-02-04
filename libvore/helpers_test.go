package libvore

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
)

type TestMatch struct {
	offset      int
	value       string
	replacement string
	variables   []TestVar
}

type TestVar struct {
	key   string
	value string
}

func singleMatch(t *testing.T, results Matches, startOffset int, value string) {
	t.Helper()
	if len(results) < 1 {
		t.FailNow()
	}
	if len(results) > 1 {
		t.Fail()
	}

	match := results[0]
	if match.value != value || match.offset.Start != startOffset {
		t.FailNow()
	}
}

func matches(t *testing.T, results Matches, expected []TestMatch) {
	t.Helper()
	if len(results) != len(expected) {
		t.Errorf("Expected %d results, got %d results\n", len(expected), len(results))
		t.FailNow()
	}

	for i, e := range expected {
		actual := results[i]
		if actual.value != e.value {
			t.Errorf("Expected value %s, got %s\n", e.value, actual.value)
		}
		if actual.offset.Start != e.offset {
			t.Errorf("Expected offset %d, got %d", e.offset, actual.offset.Start)
		}
		if actual.replacement != e.replacement {
			t.Errorf("Expected replacement %s, got %s\n", e.replacement, actual.replacement)
		}
		if actual.variables.Len() != len(e.variables) {
			t.Errorf("Expected %d variables, got %d variables\n", len(e.variables), actual.variables.Len())
		} else {
			for _, exVar := range e.variables {
				v, prs := actual.variables.Get(exVar.key)
				if prs && v.String().Value != exVar.value {
					t.Errorf("Expected %s, got %s\n", exVar.value, v.String().Value)
				}
			}
		}
	}
}

func checkNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Error(err)
	}
}

func mustPanic(t *testing.T, message string, process func(*testing.T)) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(message)
		}
	}()

	process(t)
}

func pseudo_uuid(t *testing.T) (string, error) {
	t.Helper()
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

func getTestingFilename(t *testing.T, touchFile bool) string {
	filename, err := pseudo_uuid(t)
	checkNoError(t, err)

	fullFilename := filename + ".txt"

	if touchFile {
		if _, err := os.Stat(fullFilename); err == nil {
			t.Errorf("Randomized testing file '%s' already exists. I don't want to overwrite it :(", fullFilename)
		}

		file, err := os.OpenFile(fullFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(0666))
		checkNoError(t, err)

		err = file.Close()
		checkNoError(t, err)
	}

	return fullFilename
}

func removeTestingFile(t *testing.T, filename string) {
	err := os.Remove(filename)
	checkNoError(t, err)
}