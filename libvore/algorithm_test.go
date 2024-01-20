package libvore

import (
	"reflect"
	"strings"
	"testing"
)

func TestMax1(t *testing.T) {
	max := Max(-1, 10)
	if max != 10 {
		t.Errorf("Max was supposed to be 10 but was %d", max)
	}
}

func TestMax2(t *testing.T) {
	max := Max(10, -3)
	if max != 10 {
		t.Errorf("Max was supposed to be 10 but was %d", max)
	}
}

func TestMax3(t *testing.T) {
	max := Max(-1, -5)
	if max != -1 {
		t.Errorf("Max was supposed to be -1 but was %d", max)
	}
}

func TestMax4(t *testing.T) {
	max := Max(20, 100)
	if max != 100 {
		t.Errorf("Max was supposed to be 100 but was %d", max)
	}
}

func TestMin1(t *testing.T) {
	max := Min(-1, 10)
	if max != -1 {
		t.Errorf("Min was supposed to be -1 but was %d", max)
	}
}

func TestMin2(t *testing.T) {
	max := Min(10, -3)
	if max != -3 {
		t.Errorf("Min was supposed to be -3 but was %d", max)
	}
}

func TestMin3(t *testing.T) {
	max := Min(-1, -5)
	if max != -5 {
		t.Errorf("Min was supposed to be -5 but was %d", max)
	}
}

func TestMin4(t *testing.T) {
	max := Min(20, 100)
	if max != 20 {
		t.Errorf("Min was supposed to be 20 but was %d", max)
	}
}

func TestSplitKeep(t *testing.T) {
	result := SplitKeep("abcaabc", "a")
	if !reflect.DeepEqual(result, []string{"a", "bc", "a", "a", "bc"}) {
		t.Errorf("SplitKeep was supposed to be [a, bc, a, a, bc] but was [%s]", strings.Join(result, ", "))
	}
}
