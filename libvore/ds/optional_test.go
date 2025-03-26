package ds

import (
	"testing"

	"github.com/jmeaster30/vore/libvore/testutils"
)

func TestOptionalNoneHasValue(t *testing.T) {
	value := None[string]()

	if value.HasValue() {
		t.Errorf("None optional had value when it was supposed to have no value.")
	}
}

func TestOptionalNoneHasValueGetValue(t *testing.T) {
	value := None[string]()

	if value.HasValue() {
		t.Errorf("None optional had value when it was supposed to have no value.")
	}

	testutils.MustPanic(t, "Expected None value to panic when reading value.",
		func(t *testing.T) {
			value.GetValue()
		})
}

func TestOptionalSomeHasValueGetValue(t *testing.T) {
	value := Some("hello :)")

	if !value.HasValue() {
		t.Errorf("Some optional has no value when it was supposed to have a value.")
	}

	result := value.GetValue()
	if result != "hello :)" {
		t.Errorf("Some optional returned incorrect value from GetValue(). Expected 'hello :)' but got '%s'", result)
	}
}

func TestOptionalNoneGetValueOrDefault(t *testing.T) {
	value := None[int]()

	result := value.GetValueOrDefault(5)
	if result != 5 {
		t.Errorf("None optional was supposed to return 5 from GetValueOrDefault but actually returned %d.", result)
	}
}

func TestOptionalSomeGetValueOrDefault(t *testing.T) {
	value := Some("hello :)")

	result := value.GetValueOrDefault("oh no :(")
	if result != "hello :)" {
		t.Errorf("Some optional returned incorrect value from GetValueOrDefault(). Expected 'hello :)' but got '%s'", result)
	}
}
