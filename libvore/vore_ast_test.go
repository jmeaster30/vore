package libvore

import "testing"

func TestBigTest(t *testing.T) {
	vore, err := Compile(`
		find all "yeah"
		find skip 1 take 1 between 2 and 3 (caseless "hello" or "world") in 'a', 'b', 'c' to 'f' at least 6 "!" at most 9 ":)"
		replace all "helloworld" any {whitespace digit} = test upper = yeah yeah (lower letter (line start) file start line end) file end with "wow!!"
		set a to pattern "okay"
	`)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	vore.PrintAST()
}
