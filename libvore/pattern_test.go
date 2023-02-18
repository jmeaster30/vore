package libvore

import "testing"

func TestDivisibleBy3Check(t *testing.T) {
	vore, err := Compile(`
set divisibleBy3 to pattern
	at least 1 digit
begin
	return match % 3 == 0
end

find all divisibleBy3`)
	checkNoError(t, err)
	results := vore.Run("123 4 6 51 52")
	matches(t, results, []TestMatch{
		{0, "123", "", []TestVar{}},
		{6, "6", "", []TestVar{}},
		{8, "51", "", []TestVar{}},
	})
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
	checkNoError(t, err)
	results := vore.Run("this is a test")
	matches(t, results, []TestMatch{
		{0, "this", ">this<", []TestVar{}},
		{5, "is", ">isis<", []TestVar{}},
		{8, "a", ">aaa<", []TestVar{}},
		{10, "test", ">testtesttesttest<", []TestVar{}},
	})
}
