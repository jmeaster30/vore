package libvore

import "testing"

func TestCustomPredicates(t *testing.T) {
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
