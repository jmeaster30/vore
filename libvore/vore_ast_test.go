package libvore

import (
	"testing"
)

func TestBigTest(t *testing.T) {
	vore, err := Compile(`
		find skip 1 take 1 between 2 and 3 ("hello" or "world") at least 6 "!" at most 9 ":)"
		replace all "helloworld" with "wow!!"
		set a to "okay"
	`)
	if err != nil {
		t.Fail()
	}
	vore.PrintAST()
}
