package ds

import (
	"testing"

	"github.com/jmeaster30/vore/libvore/testutils"
)

func TestOptionalNoneHasValue(t *testing.T) {
	value := None[string]()

	testutils.AssertFalse(t, value.HasValue())
}

func TestOptionalNoneHasValueGetValue(t *testing.T) {
	value := None[string]()

	testutils.AssertFalse(t, value.HasValue())
	testutils.MustPanic(t, "Expected None value to panic when reading value.",
		func(t *testing.T) {
			value.GetValue()
		})
}

func TestOptionalSomeHasValueGetValue(t *testing.T) {
	value := Some("hello :)")

	testutils.AssertTrue(t, value.HasValue())

	result := value.GetValue()
	testutils.AssertEqual(t, "hello :)", result)
}

func TestOptionalNoneGetValueOrDefault(t *testing.T) {
	value := None[int]()

	result := value.GetValueOrDefault(5)
	testutils.AssertEqual(t, 5, result)
}

func TestOptionalSomeGetValueOrDefault(t *testing.T) {
	value := Some("hello :)")

	result := value.GetValueOrDefault("oh no :(")
	testutils.AssertEqual(t, "hello :)", result)
}

func TestOptionalEqualNone(t *testing.T) {
	left := None[int]()
	right := None[int]()

	testutils.AssertTrue(t, OptionalEqual(left, right))
}

func TestOptionalEqualNoneAndSome(t *testing.T) {
	left := None[int]()
	right := Some(12)

	testutils.AssertFalse(t, OptionalEqual(left, right))
}

func TestOptionalEqualSomeAndNone(t *testing.T) {
	left := Some(1)
	right := None[int]()

	testutils.AssertFalse(t, OptionalEqual(left, right))
}

func TestOptionalEqualSomeNotEqual(t *testing.T) {
	left := Some(1)
	right := Some(3)

	testutils.AssertFalse(t, OptionalEqual(left, right))
}

func TestOptionalEqualSomeEqual(t *testing.T) {
	left := Some(1)
	right := Some(1)

	testutils.AssertTrue(t, OptionalEqual(left, right))
}

func TestOptionalMapSomeValue(t *testing.T) {
	value := Some(1)
	actual := OptionalMap(value, func(a int) int { return a + 1 })
	testutils.AssertEqual(t, Some(2), actual)
}

func TestOptionalMapNoneValue(t *testing.T) {
	value := None[int]()
	actual := OptionalMap(value, func(a int) int { return a + 1 })
	testutils.AssertEqual(t, None[int](), actual)
}
