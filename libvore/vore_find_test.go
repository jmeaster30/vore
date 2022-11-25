package libvore

import (
	"testing"
)

type TestMatch struct {
	offset    int
	value     string
	variables []TestVar
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
	if len(results) < len(expected) {
		t.FailNow()
	}
	if len(results) > len(expected) {
		t.Fail()
	}

	for i, e := range expected {
		actual := results[i]
		if actual.value != e.value || actual.offset.Start != e.offset {
			t.Fail()
		}
		if len(actual.variables) != len(e.variables) {
			t.Fail()
		} else {
			for _, exVar := range e.variables {
				if actual.variables[exVar.key] != exVar.value {
					t.Fail()
				}
			}
		}
	}
}

func TestFindString(t *testing.T) {
	vore := Compile("find all 'yay'")
	results := vore.Run("OMG yay :)")
	singleMatch(t, results, 4, "yay")
}

func TestFindDigit(t *testing.T) {
	vore := Compile("find all digit")
	results := vore.Run("please 1234567890 wow")
	matches(t, results, []TestMatch{
		{7, "1", []TestVar{}},
		{8, "2", []TestVar{}},
		{9, "3", []TestVar{}},
		{10, "4", []TestVar{}},
		{11, "5", []TestVar{}},
		{12, "6", []TestVar{}},
		{13, "7", []TestVar{}},
		{14, "8", []TestVar{}},
		{15, "9", []TestVar{}},
		{16, "0", []TestVar{}},
	})
}

func TestFindAtLeast1Digit(t *testing.T) {
	vore := Compile("find all at least 1 digit")
	results := vore.Run("please 1234567890 wow")
	singleMatch(t, results, 7, "1234567890")
}

func TestFindEscapedCharacters(t *testing.T) {
	vore := Compile("find all '\\x77\\x6f\\x77\\x20\\x3B\\x29'")
	results := vore.Run("does this work? wow ;)")
	singleMatch(t, results, 16, "wow ;)")
}

func TestFindWhitespace(t *testing.T) {
	vore := Compile("find all whitespace 'source' whitespace")
	results := vore.Run("you must provide a source for your claims.")
	singleMatch(t, results, 18, " source ")
}

func TestFindLetter(t *testing.T) {
	vore := Compile("find all letter")
	results := vore.Run("345A98(&$(#*%")
	singleMatch(t, results, 3, "A")
}
