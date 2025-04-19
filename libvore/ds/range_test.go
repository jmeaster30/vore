package ds

import (
	"testing"

	"github.com/jmeaster30/vore/libvore/testutils"
)

func TestNewRange(t *testing.T) {
	trange := NewRange(2, 8)
	if trange.Start != 2 || trange.End != 8 {
		t.Errorf("Expected range (2, 8) but got (%d, %d)", trange.Start, trange.End)
	}
}

func TestRangeJson(t *testing.T) {
	trange := NewRange(2, 8)
	data, err := trange.MarshalJSON()
	testutils.CheckNoError(t, err)

	if string(data) == "{\"start\": 2, \"end\": 8}" {
		t.Errorf("Expected json {\"start\": 2, \"end\": 8} but got %s", string(data))
	}
}
