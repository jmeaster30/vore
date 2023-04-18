package libvore

import "testing"

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
		{7, "1", None[string](), []TestVar{}},
		{8, "2", None[string](), []TestVar{}},
		{9, "3", None[string](), []TestVar{}},
		{10, "4", None[string](), []TestVar{}},
		{11, "5", None[string](), []TestVar{}},
		{12, "6", None[string](), []TestVar{}},
		{13, "7", None[string](), []TestVar{}},
		{14, "8", None[string](), []TestVar{}},
		{15, "9", None[string](), []TestVar{}},
		{16, "0", None[string](), []TestVar{}},
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
		{0, "omg t", None[string](), []TestVar{}},
		{5, "his i", None[string](), []TestVar{}},
		{10, "s coo", None[string](), []TestVar{}},
		{15, "l :)", None[string](), []TestVar{}},
	})
}

func TestFindAnyFewest(t *testing.T) {
	vore, err := Compile("find all between 3 and 5 any fewest")
	checkNoError(t, err)
	results := vore.Run("omg this is")
	matches(t, results, []TestMatch{
		{0, "omg", None[string](), []TestVar{}},
		{3, " th", None[string](), []TestVar{}},
		{6, "is ", None[string](), []TestVar{}},
	})
}

func TestFindFewest(t *testing.T) {
	vore, err := Compile("find all at least 3 letter fewest ' '")
	checkNoError(t, err)
	results := vore.Run("oh wow geez nice")
	matches(t, results, []TestMatch{
		{3, "wow ", None[string](), []TestVar{}},
		{7, "geez ", None[string](), []TestVar{}},
	})
}

func TestFindAtLeast3Upper(t *testing.T) {
	vore, err := Compile("find all at least 3 upper")
	checkNoError(t, err)
	results := vore.Run("it SHOULD get THIS but THis")
	matches(t, results, []TestMatch{
		{3, "SHOULD", None[string](), []TestVar{}},
		{14, "THIS", None[string](), []TestVar{}},
	})
}

func TestFindAtMost2Lower(t *testing.T) {
	vore, err := Compile("find all at most 2 lower")
	checkNoError(t, err)
	results := vore.Run("IT WILL CATCH this AND it WILL GET me")
	matches(t, results, []TestMatch{
		{14, "th", None[string](), []TestVar{}},
		{16, "is", None[string](), []TestVar{}},
		{23, "it", None[string](), []TestVar{}},
		{35, "me", None[string](), []TestVar{}},
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
		{6, "here", None[string](), []TestVar{}},
		{13, "here", None[string](), []TestVar{}},
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
		{0, "this", None[string](), []TestVar{}},
		{9, "that", None[string](), []TestVar{}},
	})
}

func TestInBranch(t *testing.T) {
	vore, err := Compile("find all in 'a', 'b', 'c'")
	checkNoError(t, err)
	results := vore.Run("abcdefghijklmnopqrstuvwxyz")
	matches(t, results, []TestMatch{
		{0, "a", None[string](), []TestVar{}},
		{1, "b", None[string](), []TestVar{}},
		{2, "c", None[string](), []TestVar{}},
	})
}

func TestInBranchRange(t *testing.T) {
	vore, err := Compile("find all in 'a' to 'c', 'x' to 'z'")
	checkNoError(t, err)
	results := vore.Run("abcdefghijklmnopqrstuvwxyz")
	matches(t, results, []TestMatch{
		{0, "a", None[string](), []TestVar{}},
		{1, "b", None[string](), []TestVar{}},
		{2, "c", None[string](), []TestVar{}},
		{23, "x", None[string](), []TestVar{}},
		{24, "y", None[string](), []TestVar{}},
		{25, "z", None[string](), []TestVar{}},
	})
}

func TestVariables(t *testing.T) {
	vore, err := Compile("find all (at least 1 in 'a' to 'c', 'x' to 'z') = test")
	checkNoError(t, err)
	results := vore.Run("abcdefghijklmnopqrstuvwxyz")
	matches(t, results, []TestMatch{
		{0, "abc", None[string](), []TestVar{
			{"test", "abc"},
		}},
		{23, "xyz", None[string](), []TestVar{
			{"test", "xyz"},
		}},
	})
}

func TestNameIdMatches(t *testing.T) {
	vore, err := Compile(`
		find all
			line start (
				(exactly 2 upper) = country
				(at least 6 digit) = department
			) = id
			'\t'
			(at least 1 any fewest) = name
			line end`)
	checkNoError(t, err)
	results := vore.Run(`US123456	lilith
tx555555	martha
FR420420	celeste`)
	matches(t, results, []TestMatch{
		{0, "US123456\tlilith", None[string](), []TestVar{
			{"country", "US"},
			{"department", "123456"},
			{"id", "US123456"},
			{"name", "lilith"},
		}},
		{32, "FR420420\tceleste", None[string](), []TestVar{
			{"country", "FR"},
			{"department", "420420"},
			{"id", "FR420420"},
			{"name", "celeste"},
		}},
	})
}

func TestVariableMatch(t *testing.T) {
	vore, err := Compile("find all 'wow' = wow wow")
	checkNoError(t, err)
	results := vore.Run("wow wowwow")
	matches(t, results, []TestMatch{
		{4, "wowwow", None[string](), []TestVar{
			{"wow", "wow"},
		}},
	})
}

func TestReplaceStatement(t *testing.T) {
	vore, err := Compile("replace all 'wow' = wow with '>' wow wow '<'")
	checkNoError(t, err)
	results := vore.Run("wow wowwow")
	matches(t, results, []TestMatch{
		{0, "wow", Some(">wowwow<"), []TestVar{
			{"wow", "wow"},
		}},
		{4, "wow", Some(">wowwow<"), []TestVar{
			{"wow", "wow"},
		}},
		{7, "wow", Some(">wowwow<"), []TestVar{
			{"wow", "wow"},
		}},
	})
}

func TestNot(t *testing.T) {
	vore, err := Compile("find all at least 1 not whitespace")
	checkNoError(t, err)
	results := vore.Run("this \tfinds all  \nnon-whitespace!")
	matches(t, results, []TestMatch{
		{0, "this", None[string](), []TestVar{}},
		{6, "finds", None[string](), []TestVar{}},
		{12, "all", None[string](), []TestVar{}},
		{18, "non-whitespace!", None[string](), []TestVar{}},
	})
}

func TestNotInBasic(t *testing.T) {
	vore, err := Compile("find all not in 'a' to 'c', 'x' to 'z'")
	checkNoError(t, err)
	results := vore.Run("abcdefxyzghi")
	matches(t, results, []TestMatch{
		{3, "d", None[string](), []TestVar{}},
		{4, "e", None[string](), []TestVar{}},
		{5, "f", None[string](), []TestVar{}},
		{9, "g", None[string](), []TestVar{}},
		{10, "h", None[string](), []TestVar{}},
		{11, "i", None[string](), []TestVar{}},
	})
}

func TestNotInInLoop(t *testing.T) {
	vore, err := Compile("find all at least 1 (not in 'a' to 'c', 'x' to 'z')")
	checkNoError(t, err)
	results := vore.Run("abcdefxyzghi")
	matches(t, results, []TestMatch{
		{3, "def", None[string](), []TestVar{}},
		{9, "ghi", None[string](), []TestVar{}},
	})
}

func TestBlockComment(t *testing.T) {
	vore, err := Compile("--(find all at least))- 1 (not in 'a' to 'c', 'x' to 'z'))--")
	checkNoError(t, err)
	results := vore.Run("oh wow a test!")
	matches(t, results, []TestMatch{})
}

func TestEmail(t *testing.T) {
	vore, err := Compile(`
set localPart to pattern
  in letter, digit, "!", "#", "$", "%",
    "&", "'", "*", "+", "/", "=", "?",
    "^", "_", "{", "|", "}", "~", "-"  -- it is a little long to write but "verbose" is in the name

set hexPart1 to pattern
  in "\x01" to "\x08", "\x0b", "\x0C",
    "\x0e" to "\x1f", "\x21",
    "\x23" to "\x5b", "\x5d" to "\x7f"

set hexPart2 to pattern
  in "\x01" to "\x09", "\x0b", "\x0C",
    "\x0e" to "\x7f"

set ld to pattern
  in letter, digit

set ldd to pattern
  in letter, digit, "-"

find all
  (at least 1 localPart at least 0 ("." at least 1 localPart))
  or
  ('"' at least 0 (hexPart1 or ('\\' hexPart2)) '"')
  "@"
  (at least 1 (ld maybe (at least 0 ldd ld) '.') ld maybe (at least 0 ldd ld))
  or
  ("["
    exactly 3 (("25" in "0" to "5") or (("2" in "0" to "4" digit) or (maybe ("0" or "1") digit maybe digit)) ".")
    (("25" in "0" to "5") or (("2" in "0" to "4" digit) or (maybe ("0" or "1") digit maybe digit)))
      or (maybe (at least 0 ldd ld) ":" at least 1 (hexPart1 or ("\\" hexPart2)))
  "]")`)
	checkNoError(t, err)
	results := vore.Run("jhneasterday09@gmail.com")
	matches(t, results, []TestMatch{
		{0, "jhneasterday09@gmail.com", None[string](), []TestVar{}},
	})
}

func TestCSV(t *testing.T) {
	vore, err := CompileFile("../docs/examples/csv.vore")
	checkNoError(t, err)
	results := vore.Run(`a, b, c
1, 2, 3
x, y, z`)
	matches(t, results, []TestMatch{
		{0, "a, b, c\n", None[string](), []TestVar{
			{"row", "[ValueHashMap]"}, // TODO check the nested structure
		}},
		{8, "1, 2, 3\n", None[string](), []TestVar{
			{"row", "[ValueHashMap]"}, // TODO check the nested structure
		}},
		{16, "x, y, z", None[string](), []TestVar{
			{"row", "[ValueHashMap]"}, // TODO check the nested structure
		}},
	})
}

func TestCaseless(t *testing.T) {
	vore, err := Compile("find all caseless 'test'")
	checkNoError(t, err)
	results := vore.Run(`
		this is a test
		this is a TEST
		this is a Test
		this is a tEsT
	`)
	matches(t, results, []TestMatch{
		{13, "test", None[string](), []TestVar{}},
		{30, "TEST", None[string](), []TestVar{}},
		{47, "Test", None[string](), []TestVar{}},
		{64, "tEsT", None[string](), []TestVar{}},
	})
}

func TestRegexp(t *testing.T) {
	vore, err := Compile("find all @/a+b*/")
	checkNoError(t, err)
	results := vore.Run("aaabbb ab a")
	matches(t, results, []TestMatch{
		{0, "aaabbb", None[string](), []TestVar{}},
		{7, "ab", None[string](), []TestVar{}},
		{10, "a", None[string](), []TestVar{}},
	})
}

func TestRegexp2(t *testing.T) {
	vore, err := Compile("find all @/a*?b/")
	checkNoError(t, err)
	results := vore.Run("aaabbb ab a")
	matches(t, results, []TestMatch{
		{0, "aaab", None[string](), []TestVar{}},
		{4, "b", None[string](), []TestVar{}},
		{5, "b", None[string](), []TestVar{}},
		{7, "ab", None[string](), []TestVar{}},
	})
}

func TestRegexp3(t *testing.T) {
	vore, err := Compile("find all @/a+?b?/")
	checkNoError(t, err)
	results := vore.Run("aaabbb ab a")
	matches(t, results, []TestMatch{
		{0, "a", None[string](), []TestVar{}},
		{1, "a", None[string](), []TestVar{}},
		{2, "ab", None[string](), []TestVar{}},
		{7, "ab", None[string](), []TestVar{}},
		{10, "a", None[string](), []TestVar{}},
	})
}

func TestRegexp4(t *testing.T) {
	vore, err := Compile("find all @/a{4,7}/")
	checkNoError(t, err)
	results := vore.Run(`aaaaaaaa
	aaa aaaaaa`)
	matches(t, results, []TestMatch{
		{0, "aaaaaaa", None[string](), []TestVar{}},
		{14, "aaaaaa", None[string](), []TestVar{}},
	})
}

func TestRegexp5(t *testing.T) {
	vore, err := Compile("find all @/a{4,}/")
	checkNoError(t, err)
	results := vore.Run(`aaaaaaaa
	aaa aaaaaa`)
	matches(t, results, []TestMatch{
		{0, "aaaaaaaa", None[string](), []TestVar{}},
		{14, "aaaaaa", None[string](), []TestVar{}},
	})
}

func TestRegexp6(t *testing.T) {
	vore, err := Compile("find all @/a{4}/")
	checkNoError(t, err)
	results := vore.Run(`aaaaaaaa
	aaa aaaaaa`)
	matches(t, results, []TestMatch{
		{0, "aaaa", None[string](), []TestVar{}},
		{4, "aaaa", None[string](), []TestVar{}},
		{14, "aaaa", None[string](), []TestVar{}},
	})
}

func TestRegexp7(t *testing.T) {
	vore, err := Compile("find all @/a{4,}?/")
	checkNoError(t, err)
	results := vore.Run(`aaaaaaaa
	aaa aaaaaa`)
	matches(t, results, []TestMatch{
		{0, "aaaa", None[string](), []TestVar{}},
		{4, "aaaa", None[string](), []TestVar{}},
		{14, "aaaa", None[string](), []TestVar{}},
	})
}

func TestRegexp8(t *testing.T) {
	vore, err := Compile("find all @/.{3}/")
	checkNoError(t, err)
	results := vore.Run(`12312312312`)
	matches(t, results, []TestMatch{
		{0, "123", None[string](), []TestVar{}},
		{3, "123", None[string](), []TestVar{}},
		{6, "123", None[string](), []TestVar{}},
	})
}
