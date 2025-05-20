package libvore

import (
	"testing"

	"github.com/jmeaster30/vore/libvore/ds"
	"github.com/jmeaster30/vore/libvore/testutils"
)

func TestFindString(t *testing.T) {
	vore, err := Compile("find all 'yay'")
	testutils.CheckNoError(t, err)
	results := vore.Run("OMG yay :)")
	singleMatch(t, results, 4, "yay")
}

func TestFindDigit(t *testing.T) {
	vore, err := Compile("find all digit")
	testutils.CheckNoError(t, err)
	results := vore.Run("please 1234567890 wow")
	matches(t, results, []TestMatch{
		{7, "1", ds.None[string](), []TestVar{}},
		{8, "2", ds.None[string](), []TestVar{}},
		{9, "3", ds.None[string](), []TestVar{}},
		{10, "4", ds.None[string](), []TestVar{}},
		{11, "5", ds.None[string](), []TestVar{}},
		{12, "6", ds.None[string](), []TestVar{}},
		{13, "7", ds.None[string](), []TestVar{}},
		{14, "8", ds.None[string](), []TestVar{}},
		{15, "9", ds.None[string](), []TestVar{}},
		{16, "0", ds.None[string](), []TestVar{}},
	})
}

func TestFindAtLeast1Digit(t *testing.T) {
	vore, err := Compile("find all at least 1 digit")
	testutils.CheckNoError(t, err)
	results := vore.Run("please 1234567890 wow")
	singleMatch(t, results, 7, "1234567890")
}

func TestFindEscapedCharacters(t *testing.T) {
	vore, err := Compile("find all '\\x77\\x6f\\x77\\x20\\x3B\\x29'")
	testutils.CheckNoError(t, err)
	results := vore.Run("does this work? wow ;)")
	singleMatch(t, results, 16, "wow ;)")
}

func TestFindWhitespace(t *testing.T) {
	vore, err := Compile("find all whitespace 'source' whitespace")
	testutils.CheckNoError(t, err)
	results := vore.Run("you must provide a source for your claims.")
	singleMatch(t, results, 18, " source ")
}

func TestFindLetter(t *testing.T) {
	vore, err := Compile("find all letter")
	testutils.CheckNoError(t, err)
	results := vore.Run("345A98(&$(#*%")
	singleMatch(t, results, 3, "A")
}

func TestFindAny(t *testing.T) {
	vore, err := Compile("find all between 3 and 5 any")
	testutils.CheckNoError(t, err)
	results := vore.Run("omg this is cool :)")
	matches(t, results, []TestMatch{
		{0, "omg t", ds.None[string](), []TestVar{}},
		{5, "his i", ds.None[string](), []TestVar{}},
		{10, "s coo", ds.None[string](), []TestVar{}},
		{15, "l :)", ds.None[string](), []TestVar{}},
	})
}

func TestFindAnyFewest(t *testing.T) {
	vore, err := Compile("find all between 3 and 5 any fewest")
	testutils.CheckNoError(t, err)
	results := vore.Run("omg this is")
	matches(t, results, []TestMatch{
		{0, "omg", ds.None[string](), []TestVar{}},
		{3, " th", ds.None[string](), []TestVar{}},
		{6, "is ", ds.None[string](), []TestVar{}},
	})
}

func TestFindFewest(t *testing.T) {
	vore, err := Compile("find all at least 3 letter fewest ' '")
	testutils.CheckNoError(t, err)
	results := vore.Run("oh wow geez nice")
	matches(t, results, []TestMatch{
		{3, "wow ", ds.None[string](), []TestVar{}},
		{7, "geez ", ds.None[string](), []TestVar{}},
	})
}

func TestFindAtLeast3Upper(t *testing.T) {
	vore, err := Compile("find all at least 3 upper")
	testutils.CheckNoError(t, err)
	results := vore.Run("it SHOULD get THIS but THis")
	matches(t, results, []TestMatch{
		{3, "SHOULD", ds.None[string](), []TestVar{}},
		{14, "THIS", ds.None[string](), []TestVar{}},
	})
}

func TestFindAtMost2Lower(t *testing.T) {
	vore, err := Compile("find all at most 2 lower")
	testutils.CheckNoError(t, err)
	results := vore.Run("IT WILL CATCH this AND it WILL GET me")
	matches(t, results, []TestMatch{
		{14, "th", ds.None[string](), []TestVar{}},
		{16, "is", ds.None[string](), []TestVar{}},
		{23, "it", ds.None[string](), []TestVar{}},
		{35, "me", ds.None[string](), []TestVar{}},
	})
}

func TestSkipTest(t *testing.T) {
	vore, err := Compile("find skip 1 take 1 'here'")
	testutils.CheckNoError(t, err)
	results := vore.Run("here >here< here")
	singleMatch(t, results, 6, "here")
}

func TestTopTest(t *testing.T) {
	vore, err := Compile("find top 1 'here'")
	testutils.CheckNoError(t, err)
	results := vore.Run(">here< here here")
	singleMatch(t, results, 1, "here")
}

func TestLastTest(t *testing.T) {
	vore, err := Compile("find last 2 'here'")
	testutils.CheckNoError(t, err)
	results := vore.Run("here >here< >here<")
	matches(t, results, []TestMatch{
		{6, "here", ds.None[string](), []TestVar{}},
		{13, "here", ds.None[string](), []TestVar{}},
	})
}

func TestRecursion1(t *testing.T) {
	vore, err := Compile("find all {'a' maybe mySub 'b'} = mySub")
	testutils.CheckNoError(t, err)
	results := vore.Run("aaaabbbb")
	singleMatch(t, results, 0, "aaaabbbb")
}

func TestRecursion2(t *testing.T) {
	vore, err := Compile("find all {'a' maybe mySub 'b'} = mySub")
	testutils.CheckNoError(t, err)
	results := vore.Run("aabbb")
	singleMatch(t, results, 0, "aabb")
}

func TestRecursion3(t *testing.T) {
	vore, err := Compile("find all {'a' maybe mySub 'b'} = mySub")
	testutils.CheckNoError(t, err)
	results := vore.Run("aaaaab")
	singleMatch(t, results, 4, "ab")
}

func TestOrBranch(t *testing.T) {
	vore, err := Compile("find all 'this' or 'that'")
	testutils.CheckNoError(t, err)
	results := vore.Run("this and that")
	matches(t, results, []TestMatch{
		{0, "this", ds.None[string](), []TestVar{}},
		{9, "that", ds.None[string](), []TestVar{}},
	})
}

func TestInBranch(t *testing.T) {
	vore, err := Compile("find all in 'a', 'b', 'c'")
	testutils.CheckNoError(t, err)
	results := vore.Run("abcdefghijklmnopqrstuvwxyz")
	matches(t, results, []TestMatch{
		{0, "a", ds.None[string](), []TestVar{}},
		{1, "b", ds.None[string](), []TestVar{}},
		{2, "c", ds.None[string](), []TestVar{}},
	})
}

func TestInBranchRange(t *testing.T) {
	vore, err := Compile("find all in 'a' to 'c', 'x' to 'z'")
	testutils.CheckNoError(t, err)
	results := vore.Run("abcdefghijklmnopqrstuvwxyz")
	matches(t, results, []TestMatch{
		{0, "a", ds.None[string](), []TestVar{}},
		{1, "b", ds.None[string](), []TestVar{}},
		{2, "c", ds.None[string](), []TestVar{}},
		{23, "x", ds.None[string](), []TestVar{}},
		{24, "y", ds.None[string](), []TestVar{}},
		{25, "z", ds.None[string](), []TestVar{}},
	})
}

func TestVariables(t *testing.T) {
	vore, err := Compile("find all (at least 1 in 'a' to 'c', 'x' to 'z') = test")
	testutils.CheckNoError(t, err)
	results := vore.Run("abcdefghijklmnopqrstuvwxyz")
	matches(t, results, []TestMatch{
		{0, "abc", ds.None[string](), []TestVar{
			{"test", "abc"},
		}},
		{23, "xyz", ds.None[string](), []TestVar{
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
	testutils.CheckNoError(t, err)
	results := vore.Run(`US123456	lilith
tx555555	martha
FR420420	celeste`)
	matches(t, results, []TestMatch{
		{0, "US123456\tlilith", ds.None[string](), []TestVar{
			{"country", "US"},
			{"department", "123456"},
			{"id", "US123456"},
			{"name", "lilith"},
		}},
		{32, "FR420420\tceleste", ds.None[string](), []TestVar{
			{"country", "FR"},
			{"department", "420420"},
			{"id", "FR420420"},
			{"name", "celeste"},
		}},
	})
}

func TestVariableMatch(t *testing.T) {
	vore, err := Compile("find all 'wow' = wow wow")
	testutils.CheckNoError(t, err)
	results := vore.Run("wow wowwow")
	matches(t, results, []TestMatch{
		{4, "wowwow", ds.None[string](), []TestVar{
			{"wow", "wow"},
		}},
	})
}

func TestReplaceStatement(t *testing.T) {
	vore, err := Compile("replace all 'wow' = wow with '>' wow wow '<'")
	testutils.CheckNoError(t, err)
	results := vore.Run("wow wowwow")
	matches(t, results, []TestMatch{
		{0, "wow", ds.Some(">wowwow<"), []TestVar{
			{"wow", "wow"},
		}},
		{4, "wow", ds.Some(">wowwow<"), []TestVar{
			{"wow", "wow"},
		}},
		{7, "wow", ds.Some(">wowwow<"), []TestVar{
			{"wow", "wow"},
		}},
	})
}

func TestNot(t *testing.T) {
	vore, err := Compile("find all at least 1 not whitespace")
	testutils.CheckNoError(t, err)
	results := vore.Run("this \tfinds all  \nnon-whitespace!")
	matches(t, results, []TestMatch{
		{0, "this", ds.None[string](), []TestVar{}},
		{6, "finds", ds.None[string](), []TestVar{}},
		{12, "all", ds.None[string](), []TestVar{}},
		{18, "non-whitespace!", ds.None[string](), []TestVar{}},
	})
}

func TestNotInBasic(t *testing.T) {
	vore, err := Compile("find all not in 'a' to 'c', 'x' to 'z'")
	testutils.CheckNoError(t, err)
	results := vore.Run("abcdefxyzghi")
	matches(t, results, []TestMatch{
		{3, "d", ds.None[string](), []TestVar{}},
		{4, "e", ds.None[string](), []TestVar{}},
		{5, "f", ds.None[string](), []TestVar{}},
		{9, "g", ds.None[string](), []TestVar{}},
		{10, "h", ds.None[string](), []TestVar{}},
		{11, "i", ds.None[string](), []TestVar{}},
	})
}

func TestNotInInLoop(t *testing.T) {
	vore, err := Compile("find all at least 1 (not in 'a' to 'c', 'x' to 'z')")
	testutils.CheckNoError(t, err)
	results := vore.Run("abcdefxyzghi")
	matches(t, results, []TestMatch{
		{3, "def", ds.None[string](), []TestVar{}},
		{9, "ghi", ds.None[string](), []TestVar{}},
	})
}

func TestBlockComment(t *testing.T) {
	vore, err := Compile("--(find all at least))- 1 (not in 'a' to 'c', 'x' to 'z'))--")
	testutils.CheckNoError(t, err)
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
	testutils.CheckNoError(t, err)
	results := vore.Run("snemail@gmail.com")
	matches(t, results, []TestMatch{
		{0, "snemail@gmail.com", ds.None[string](), []TestVar{}},
	})
}

func TestCSV(t *testing.T) {
	vore, err := CompileFile("../docs/examples/csv.vore")
	testutils.CheckNoError(t, err)
	results := vore.Run(`a, b, c
1, 2, 3
x, y, z`)
	matches(t, results, []TestMatch{
		{0, "a, b, c\n", ds.None[string](), []TestVar{
			{"row", "[ValueHashMap]"}, // TODO check the nested structure
		}},
		{8, "1, 2, 3\n", ds.None[string](), []TestVar{
			{"row", "[ValueHashMap]"}, // TODO check the nested structure
		}},
		{16, "x, y, z", ds.None[string](), []TestVar{
			{"row", "[ValueHashMap]"}, // TODO check the nested structure
		}},
	})
}

func TestCaseless(t *testing.T) {
	vore, err := Compile("find all caseless 'test'")
	testutils.CheckNoError(t, err)
	results := vore.Run(`
		this is a test
		this is a TEST
		this is a Test
		this is a tEsT
	`)
	matches(t, results, []TestMatch{
		{13, "test", ds.None[string](), []TestVar{}},
		{30, "TEST", ds.None[string](), []TestVar{}},
		{47, "Test", ds.None[string](), []TestVar{}},
		{64, "tEsT", ds.None[string](), []TestVar{}},
	})
}

func TestRegexp(t *testing.T) {
	vore, err := Compile("find all @/a+b*/")
	testutils.CheckNoError(t, err)
	results := vore.Run("aaabbb ab a")
	matches(t, results, []TestMatch{
		{0, "aaabbb", ds.None[string](), []TestVar{}},
		{7, "ab", ds.None[string](), []TestVar{}},
		{10, "a", ds.None[string](), []TestVar{}},
	})
}

func TestRegexp2(t *testing.T) {
	vore, err := Compile("find all @/a*?b/")
	testutils.CheckNoError(t, err)
	results := vore.Run("aaabbb ab a")
	matches(t, results, []TestMatch{
		{0, "aaab", ds.None[string](), []TestVar{}},
		{4, "b", ds.None[string](), []TestVar{}},
		{5, "b", ds.None[string](), []TestVar{}},
		{7, "ab", ds.None[string](), []TestVar{}},
	})
}

func TestRegexp3(t *testing.T) {
	vore, err := Compile("find all @/a+?b?/")
	testutils.CheckNoError(t, err)
	results := vore.Run("aaabbb ab a")
	matches(t, results, []TestMatch{
		{0, "a", ds.None[string](), []TestVar{}},
		{1, "a", ds.None[string](), []TestVar{}},
		{2, "ab", ds.None[string](), []TestVar{}},
		{7, "ab", ds.None[string](), []TestVar{}},
		{10, "a", ds.None[string](), []TestVar{}},
	})
}

func TestRegexp4(t *testing.T) {
	vore, err := Compile("find all @/a{4,7}/")
	testutils.CheckNoError(t, err)
	results := vore.Run(`aaaaaaaa
	aaa aaaaaa`)
	matches(t, results, []TestMatch{
		{0, "aaaaaaa", ds.None[string](), []TestVar{}},
		{14, "aaaaaa", ds.None[string](), []TestVar{}},
	})
}

func TestRegexp5(t *testing.T) {
	vore, err := Compile("find all @/a{4,}/")
	testutils.CheckNoError(t, err)
	results := vore.Run(`aaaaaaaa
	aaa aaaaaa`)
	matches(t, results, []TestMatch{
		{0, "aaaaaaaa", ds.None[string](), []TestVar{}},
		{14, "aaaaaa", ds.None[string](), []TestVar{}},
	})
}

func TestRegexp6(t *testing.T) {
	vore, err := Compile("find all @/a{4}/")
	testutils.CheckNoError(t, err)
	results := vore.Run(`aaaaaaaa
	aaa aaaaaa`)
	matches(t, results, []TestMatch{
		{0, "aaaa", ds.None[string](), []TestVar{}},
		{4, "aaaa", ds.None[string](), []TestVar{}},
		{14, "aaaa", ds.None[string](), []TestVar{}},
	})
}

func TestRegexp7(t *testing.T) {
	vore, err := Compile("find all @/a{4,}?/")
	testutils.CheckNoError(t, err)
	results := vore.Run(`aaaaaaaa
	aaa aaaaaa`)
	matches(t, results, []TestMatch{
		{0, "aaaa", ds.None[string](), []TestVar{}},
		{4, "aaaa", ds.None[string](), []TestVar{}},
		{14, "aaaa", ds.None[string](), []TestVar{}},
	})
}

func TestRegexp8(t *testing.T) {
	vore, err := Compile("find all @/.{3}/")
	testutils.CheckNoError(t, err)
	results := vore.Run(`12312312312`)
	matches(t, results, []TestMatch{
		{0, "123", ds.None[string](), []TestVar{}},
		{3, "123", ds.None[string](), []TestVar{}},
		{6, "123", ds.None[string](), []TestVar{}},
	})
}

func TestRegexp9(t *testing.T) {
	vore, err := Compile("find all @/[^]*/")
	testutils.CheckNoError(t, err)
	results := vore.Run(`1231231
	2312`)
	matches(t, results, []TestMatch{
		{0, `1231231
	2312`, ds.None[string](), []TestVar{}},
	})
}

func TestRegexp10(t *testing.T) {
	vore, err := Compile("find all @/[abc]*/")
	testutils.CheckNoError(t, err)
	results := vore.Run(`123aabbcc986`)
	matches(t, results, []TestMatch{
		{3, `aabbcc`, ds.None[string](), []TestVar{}},
	})
}

func TestRegexp11(t *testing.T) {
	vore, err := Compile("find all @/[a-z]{0,2}/")
	testutils.CheckNoError(t, err)
	results := vore.Run("IT WILL CATCH this AND it WILL GET me")
	matches(t, results, []TestMatch{
		{14, "th", ds.None[string](), []TestVar{}},
		{16, "is", ds.None[string](), []TestVar{}},
		{23, "it", ds.None[string](), []TestVar{}},
		{35, "me", ds.None[string](), []TestVar{}},
	})
}

func TestRegexp12(t *testing.T) {
	vore, err := Compile("find all @/[a-]*/")
	testutils.CheckNoError(t, err)
	results := vore.Run(`123aa--a-ac986`)
	matches(t, results, []TestMatch{
		{3, `aa--a-a`, ds.None[string](), []TestVar{}},
	})
}

func TestRegexp13(t *testing.T) {
	vore, err := Compile("find all @/test/")
	testutils.CheckNoError(t, err)
	results := vore.Run("this is a test")
	matches(t, results, []TestMatch{
		{10, "test", ds.None[string](), []TestVar{}},
	})
}

func TestRegexp14(t *testing.T) {
	vore, err := Compile("find all @/a|b/")
	testutils.CheckNoError(t, err)
	results := vore.Run("abc")
	matches(t, results, []TestMatch{
		{0, "a", ds.None[string](), []TestVar{}},
		{1, "b", ds.None[string](), []TestVar{}},
	})
}

func TestRegexp15(t *testing.T) {
	vore, err := Compile("find all @/^test/")
	testutils.CheckNoError(t, err)
	results := vore.Run("test a test")
	matches(t, results, []TestMatch{
		{0, "test", ds.None[string](), []TestVar{}},
	})
}

func TestRegexp16(t *testing.T) {
	vore, err := Compile("find all @/test$/")
	testutils.CheckNoError(t, err)
	results := vore.Run("test a test")
	matches(t, results, []TestMatch{
		{7, "test", ds.None[string](), []TestVar{}},
	})
}

func TestRegexp17(t *testing.T) {
	vore, err := Compile("find all @/[^abc]*/")
	testutils.CheckNoError(t, err)
	results := vore.Run("I really hate the abc's")
	matches(t, results, []TestMatch{
		{0, "I re", ds.None[string](), []TestVar{}},
		{5, "lly h", ds.None[string](), []TestVar{}},
		{11, "te the ", ds.None[string](), []TestVar{}},
		{21, "'s", ds.None[string](), []TestVar{}},
	})
}

func TestRegexp18(t *testing.T) {
	vore, err := Compile("find all @/[]/")
	testutils.CheckNoError(t, err)
	results := vore.Run("This is not a match")
	matches(t, results, []TestMatch{})
}

func TestNotExpressionDeclaration(t *testing.T) {
	vore, err := Compile("find all not letter = wow")
	testutils.CheckNoError(t, err)
	results := vore.Run("123 &abc")
	matches(t, results, []TestMatch{
		{0, "1", ds.None[string](), []TestVar{{"wow", "1"}}},
		{1, "2", ds.None[string](), []TestVar{{"wow", "2"}}},
		{2, "3", ds.None[string](), []TestVar{{"wow", "3"}}},
		{3, " ", ds.None[string](), []TestVar{{"wow", " "}}},
		{4, "&", ds.None[string](), []TestVar{{"wow", "&"}}},
	})
}

func TestRegexp19(t *testing.T) {
	vore, err := Compile("find all @/(?<test>a|b)/ test")
	testutils.CheckNoError(t, err)
	results := vore.Run("aabaccabjjbb")
	matches(t, results, []TestMatch{
		{0, "aa", ds.None[string](), []TestVar{{"test", "a"}}},
		{10, "bb", ds.None[string](), []TestVar{{"test", "b"}}},
	})
}

func TestRegexp20(t *testing.T) {
	vore, err := Compile("find all ('a' or 'b') = test @/\\k<test>/")
	testutils.CheckNoError(t, err)
	results := vore.Run("aabaccabjjbb")
	matches(t, results, []TestMatch{
		{0, "aa", ds.None[string](), []TestVar{{"test", "a"}}},
		{10, "bb", ds.None[string](), []TestVar{{"test", "b"}}},
	})
}

func TestBranchVariable(t *testing.T) {
	vore, err := Compile("find all ('a' or 'b') = test test")
	testutils.CheckNoError(t, err)
	results := vore.Run("aabaccabjjbb")
	matches(t, results, []TestMatch{
		{0, "aa", ds.None[string](), []TestVar{{"test", "a"}}},
		{10, "bb", ds.None[string](), []TestVar{{"test", "b"}}},
	})
}

func TestRegexp21(t *testing.T) {
	vore, err := Compile("find all at least 1 @/\\d/")
	testutils.CheckNoError(t, err)
	results := vore.Run("1234abc567")
	matches(t, results, []TestMatch{
		{0, "1234", ds.None[string](), []TestVar{}},
		{7, "567", ds.None[string](), []TestVar{}},
	})
}

func TestRegexp22(t *testing.T) {
	vore, err := Compile("find all at least 1 @/\\D/")
	testutils.CheckNoError(t, err)
	results := vore.Run("1234abc567")
	matches(t, results, []TestMatch{
		{4, "abc", ds.None[string](), []TestVar{}},
	})
}

func TestRegexp23(t *testing.T) {
	vore, err := Compile("find all at least 1 @/\\s/")
	testutils.CheckNoError(t, err)
	results := vore.Run(`12 34a	bc
567`)
	matches(t, results, []TestMatch{
		{2, " ", ds.None[string](), []TestVar{}},
		{6, "\t", ds.None[string](), []TestVar{}},
		{9, "\n", ds.None[string](), []TestVar{}},
	})
}

func TestRegexp24(t *testing.T) {
	vore, err := Compile("find all at least 1 @/\\S/")
	testutils.CheckNoError(t, err)
	results := vore.Run(`12 34a	bc
567`)
	matches(t, results, []TestMatch{
		{0, "12", ds.None[string](), []TestVar{}},
		{3, "34a", ds.None[string](), []TestVar{}},
		{7, "bc", ds.None[string](), []TestVar{}},
		{10, "567", ds.None[string](), []TestVar{}},
	})
}

func TestRegexp25(t *testing.T) {
	vore, err := Compile("find all @/\\D{0,2}/")
	testutils.CheckNoError(t, err)
	results := vore.Run(`1234abc567`)
	matches(t, results, []TestMatch{
		{4, "ab", ds.None[string](), []TestVar{}},
		{6, "c", ds.None[string](), []TestVar{}},
	})
}

func TestRegexp26(t *testing.T) {
	vore, err := Compile("find all @/\\D{2,}/")
	testutils.CheckNoError(t, err)
	results := vore.Run(`1234abc567`)
	matches(t, results, []TestMatch{
		{4, "abc", ds.None[string](), []TestVar{}},
	})
}

func TestRegexp27(t *testing.T) {
	vore, err := Compile("find all @/\\D{0,2}?/")
	testutils.CheckNoError(t, err)
	results := vore.Run(`1234abc567`)
	matches(t, results, []TestMatch{})
}

func TestRegexp28(t *testing.T) {
	vore, err := Compile("find all @/\\D{2,}?/")
	testutils.CheckNoError(t, err)
	results := vore.Run(`1234abc567`)
	matches(t, results, []TestMatch{
		{4, "ab", ds.None[string](), []TestVar{}},
	})
}

func TestRegexp29(t *testing.T) {
	vore, err := Compile("find all @/(test)\\1/")
	testutils.CheckNoError(t, err)
	results := vore.Run(`testtest`)
	matches(t, results, []TestMatch{
		{0, "testtest", ds.None[string](), []TestVar{
			{"_1", "test"},
		}},
	})
}

func TestRegexp30(t *testing.T) {
	vore, err := Compile("find all @/(t)(e)(s)(t)(e)(x)(p)(r)(e)(s)(s)(i)(o)(n)\\14\\13/")
	testutils.CheckNoError(t, err)
	results := vore.Run(`testexpressionno`)
	matches(t, results, []TestMatch{
		{0, "testexpressionno", ds.None[string](), []TestVar{
			{"_1", "t"},
			{"_2", "e"},
			{"_3", "s"},
			{"_4", "t"},
			{"_5", "e"},
			{"_6", "x"},
			{"_7", "p"},
			{"_8", "r"},
			{"_9", "e"},
			{"_10", "s"},
			{"_11", "s"},
			{"_12", "i"},
			{"_13", "o"},
			{"_14", "n"},
		}},
	})
}

func TestRegexp31(t *testing.T) {
	vore, err := Compile("find all @/(test)\\1test/")
	testutils.CheckNoError(t, err)
	results := vore.Run(`testtesttest`)
	matches(t, results, []TestMatch{
		{0, "testtesttest", ds.None[string](), []TestVar{
			{"_1", "test"},
		}},
	})
}
