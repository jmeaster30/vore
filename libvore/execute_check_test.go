package libvore

import (
	"reflect"
	"testing"

	"github.com/jmeaster30/vore/libvore/bytecode"
	"github.com/jmeaster30/vore/libvore/ds"
	"github.com/jmeaster30/vore/libvore/testutils"
)

func TestDivisibleBy3Check(t *testing.T) {
	vore, err := Compile(`
set divisibleBy3 to pattern
	at least 1 digit
begin
	return match % 3 == 0
end

find all divisibleBy3`)
	testutils.CheckNoError(t, err)
	results := vore.Run("123 4 6 51 52")
	matches(t, results, []TestMatch{
		{0, "123", ds.None[string](), []TestVar{}},
		{6, "6", ds.None[string](), []TestVar{}},
		{8, "51", ds.None[string](), []TestVar{}},
	})
}

func TestTrueLiteralIfStmt(t *testing.T) {
	vore, err := Compile(`
set check to function
	if true then
		return 'oh yeah'
	end
	return match
end

replace all 'test' with check`)
	testutils.CheckNoError(t, err)
	results := vore.Run("this is a test")
	matches(t, results, []TestMatch{
		{10, "test", ds.Some("oh yeah"), []TestVar{}},
	})
}

func TestFalseLiteralIfStmt(t *testing.T) {
	vore, err := Compile(`
set check to function
begin
	if false then
		return 'oh yeah'
	end
	return match
end

replace all 'test' with check`)
	testutils.CheckNoError(t, err)
	results := vore.Run("this is a test")
	matches(t, results, []TestMatch{
		{10, "test", ds.Some("test"), []TestVar{}},
	})
}

func TestErroredPredicate(t *testing.T) {
	vore, err := Compile(`
set divisibleBy3 to pattern
	at least 1 digit
begin
	return match
end

find all divisibleBy3`)
	checkVoreError(t, err, reflect.TypeFor[bytecode.GenError](), "Since we are in the predicate of a pattern, return values must be a boolean")
	if vore != nil {
		t.Errorf("Expected vore to be nil but it was not")
	}
}

func TestErroredTransform(t *testing.T) {
	vore, err := Compile(`
set foo to transform
	return 1 == 1
end

replace all "bar" with foo`)
	checkVoreError(t, err, reflect.TypeFor[bytecode.GenError](), "Since we are in a transform function, return values must be a string or a number")
	if vore != nil {
		t.Errorf("Expected vore to be nil but it was not")
	}
}

func TestBreakOutsideLoopError(t *testing.T) {
	vore, err := Compile(`
set foo to transform
	if 1 == 1 then
		break
	end
end

replace all "bar" with foo`)
	checkVoreError(t, err, reflect.TypeFor[bytecode.GenError](), "Cannot use 'break' outside of a loop.")
	if vore != nil {
		t.Errorf("Expected vore to be nil but it was not")
	}
}

func TestElseBranchBreakOutsideLoopError(t *testing.T) {
	vore, err := Compile(`
set foo to transform
	if 1 == 1 then
		debug ":)"
	else
		break
	end
end

replace all "bar" with foo`)
	checkVoreError(t, err, reflect.TypeFor[bytecode.GenError](), "Cannot use 'break' outside of a loop.")
	if vore != nil {
		t.Errorf("Expected vore to be nil but it was not")
	}
}

func TestContinueOutsideLoopError(t *testing.T) {
	vore, err := Compile(`
set foo to transform
	if 1 == 1 then
		continue
	end
end

replace all "bar" with foo`)
	checkVoreError(t, err, reflect.TypeFor[bytecode.GenError](), "Cannot use 'continue' outside of a loop.")
	if vore != nil {
		t.Errorf("Expected vore to be nil but it was not")
	}
}

func TestNestedIfStatementConditionError(t *testing.T) {
	vore, err := Compile(`
set foo to transform
	loop
		if "a" == "a" then
			if "wow" then
				break
			end
		end
	end
end

replace all "bar" with foo`)
	checkVoreError(t, err, reflect.TypeFor[bytecode.GenError](), "Condition of an if statement must be a boolean.")
	if vore != nil {
		t.Errorf("Expected vore to be nil but it was not")
	}
}

func TestIfStatementConditionError(t *testing.T) {
	vore, err := Compile(`
set foo to transform
	if "wow" then
		break
	end
end

replace all "bar" with foo`)
	checkVoreError(t, err, reflect.TypeFor[bytecode.GenError](), "Condition of an if statement must be a boolean.")
	if vore != nil {
		t.Errorf("Expected vore to be nil but it was not")
	}
}

func TestUndefinedOperator(t *testing.T) {
	vore, err := Compile(`
	set foo to transform
		return "test" / "this"
	end
	
	replace all "bar" with foo`)
	checkVoreError(t, err, reflect.TypeFor[bytecode.GenError](), "Operator not defined for type.")
	if vore != nil {
		t.Errorf("Expected vore to be nil but it was not")
	}
}

func TestRepeatMatchByMatchNumber(t *testing.T) {
	vore, err := Compile(`
set matchRepeater to transform
	set result to ""
	set index to 0
	loop
		if index >= matchNumber then
			break
		end
		set result to result + match
    set index to index + 1
	end
	return result
end

replace all word start at least 1 any fewest word end with ">" matchRepeater "<"`)
	testutils.CheckNoError(t, err)
	results := vore.Run("this is a test")
	matches(t, results, []TestMatch{
		{0, "this", ds.Some(">this<"), []TestVar{}},
		{5, "is", ds.Some(">isis<"), []TestVar{}},
		{8, "a", ds.Some(">aaa<"), []TestVar{}},
		{10, "test", ds.Some(">testtesttesttest<"), []TestVar{}},
	})
}

func TestTheDarknessInsideMe(t *testing.T) {
	vore, err := Compile("replace all 'hello' with 'goodbye'")
	testutils.CheckNoError(t, err)
	results := vore.Run("this is it. hello world")
	matches(t, results, []TestMatch{
		{12, "hello", ds.Some("goodbye"), []TestVar{}},
	})
}
