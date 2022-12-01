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

func checkNoError(t *testing.T, err error) {
	if err != nil {
		t.Fail()
	}
}

func TestFindString(t *testing.T) {
	vore, err := Compile("find all 'yay'")
	checkNoError(t, err)
	results := vore.Run("OMG yay :)")
	singleMatch(t, results, 4, "yay")
}

func TestFindDigit(t *testing.T) {
	vore, err := Compile("find all digit")
	checkNoError(t, err)
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
	vore, err := Compile("find all at least 1 digit")
	checkNoError(t, err)
	results := vore.Run("please 1234567890 wow")
	singleMatch(t, results, 7, "1234567890")
}

func TestFindEscapedCharacters(t *testing.T) {
	vore, err := Compile("find all '\\x77\\x6f\\x77\\x20\\x3B\\x29'")
	checkNoError(t, err)
	results := vore.Run("does this work? wow ;)")
	singleMatch(t, results, 16, "wow ;)")
}

func TestFindWhitespace(t *testing.T) {
	vore, err := Compile("find all whitespace 'source' whitespace")
	checkNoError(t, err)
	results := vore.Run("you must provide a source for your claims.")
	singleMatch(t, results, 18, " source ")
}

func TestFindLetter(t *testing.T) {
	vore, err := Compile("find all letter")
	checkNoError(t, err)
	results := vore.Run("345A98(&$(#*%")
	singleMatch(t, results, 3, "A")
}

func TestFindAny(t *testing.T) {
	vore, err := Compile("find all between 3 and 5 any")
	checkNoError(t, err)
	results := vore.Run("omg this is cool :)")
	matches(t, results, []TestMatch{
		{0, "omg t", []TestVar{}},
		{5, "his i", []TestVar{}},
		{10, "s coo", []TestVar{}},
		{15, "l :)", []TestVar{}},
	})
}

func TestFindAnyFewest(t *testing.T) {
	vore, err := Compile("find all between 3 and 5 any fewest")
	checkNoError(t, err)
	results := vore.Run("omg this is")
	matches(t, results, []TestMatch{
		{0, "omg", []TestVar{}},
		{3, " th", []TestVar{}},
		{6, "is ", []TestVar{}},
	})
}

func TestFindFewest(t *testing.T) {
	vore, err := Compile("find all at least 3 letter fewest ' '")
	checkNoError(t, err)
	results := vore.Run("oh wow geez nice")
	matches(t, results, []TestMatch{
		{3, "wow ", []TestVar{}},
		{7, "geez ", []TestVar{}},
	})
}

func TestFindAtLeast3Upper(t *testing.T) {
	vore, err := Compile("find all at least 3 upper")
	checkNoError(t, err)
	results := vore.Run("it SHOULD get THIS but THis")
	matches(t, results, []TestMatch{
		{3, "SHOULD", []TestVar{}},
		{14, "THIS", []TestVar{}},
	})
}

func TestFindAtMost2Lower(t *testing.T) {
	vore, err := Compile("find all at most 2 lower")
	checkNoError(t, err)
	results := vore.Run("IT WILL CATCH this AND it WILL GET me")
	matches(t, results, []TestMatch{
		{14, "th", []TestVar{}},
		{16, "is", []TestVar{}},
		{23, "it", []TestVar{}},
		{35, "me", []TestVar{}},
	})
}

func TestSkipTest(t *testing.T) {
	vore, err := Compile("find skip 1 take 1 'here'")
	checkNoError(t, err)
	results := vore.Run("here >here< here")
	singleMatch(t, results, 6, "here")
}

func TestTopTest(t *testing.T) {
	vore, err := Compile("find top 1 'here'")
	checkNoError(t, err)
	results := vore.Run(">here< here here")
	singleMatch(t, results, 1, "here")
}

func TestLastTest(t *testing.T) {
	vore, err := Compile("find last 2 'here'")
	checkNoError(t, err)
	results := vore.Run("here >here< >here<")
	matches(t, results, []TestMatch{
		{6, "here", []TestVar{}},
		{13, "here", []TestVar{}},
	})
}

func TestRecursion1(t *testing.T) {
	vore, err := Compile("find all {'a' maybe mySub 'b'} = mySub")
	checkNoError(t, err)
	results := vore.Run("aaaabbbb")
	singleMatch(t, results, 0, "aaaabbbb")
}

func TestRecursion2(t *testing.T) {
	vore, err := Compile("find all {'a' maybe mySub 'b'} = mySub")
	checkNoError(t, err)
	results := vore.Run("aabbb")
	singleMatch(t, results, 0, "aabb")
}

func TestRecursion3(t *testing.T) {
	vore, err := Compile("find all {'a' maybe mySub 'b'} = mySub")
	checkNoError(t, err)
	results := vore.Run("aaaaab")
	singleMatch(t, results, 4, "ab")
}

func TestOrBranch(t *testing.T) {
	vore, err := Compile("find all 'this' or 'that'")
	checkNoError(t, err)
	results := vore.Run("this and that")
	matches(t, results, []TestMatch{
		{0, "this", []TestVar{}},
		{9, "that", []TestVar{}},
	})
}

func TestInBranch(t *testing.T) {
	vore, err := Compile("find all in 'a', 'b', 'c'")
	checkNoError(t, err)
	results := vore.Run("abcdefghijklmnopqrstuvwxyz")
	matches(t, results, []TestMatch{
		{0, "a", []TestVar{}},
		{1, "b", []TestVar{}},
		{2, "c", []TestVar{}},
	})
}

func TestInBranchRange(t *testing.T) {
	vore, err := Compile("find all in 'a' to 'c', 'x' to 'z'")
	checkNoError(t, err)
	results := vore.Run("abcdefghijklmnopqrstuvwxyz")
	matches(t, results, []TestMatch{
		{0, "a", []TestVar{}},
		{1, "b", []TestVar{}},
		{2, "c", []TestVar{}},
		{23, "x", []TestVar{}},
		{24, "y", []TestVar{}},
		{25, "z", []TestVar{}},
	})
}
