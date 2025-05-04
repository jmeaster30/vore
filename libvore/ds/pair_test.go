package ds

import (
	"testing"

	"github.com/jmeaster30/vore/libvore/testutils"
)

func TestNewPair(t *testing.T) {
	pair := NewPair(12, "value")
	testutils.AssertEqual(t, 12, pair.left)
	testutils.AssertEqual(t, "value", pair.right)
}

func TestPairLeftRight(t *testing.T) {
	pair := NewPair("mine", 2345)
	testutils.AssertEqual(t, pair.left, pair.Left())
	testutils.AssertEqual(t, pair.right, pair.Right())
}

func TestPairFlip(t *testing.T) {
	pair := NewPair(1.2, ":3")
	flipped := pair.Flip()
	testutils.AssertEqual(t, pair.left, flipped.right)
	testutils.AssertEqual(t, pair.right, flipped.left)
}

func TestPairValues(t *testing.T) {
	pair := NewPair(1111, 2222)
	left, right := pair.Values()
	testutils.AssertEqual(t, pair.left, left)
	testutils.AssertEqual(t, pair.right, right)
}
