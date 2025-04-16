package libvore

import (
	"reflect"
	"testing"

	"github.com/jmeaster30/vore/libvore/ds"
	"github.com/jmeaster30/vore/libvore/engine"
)

type TestMatch struct {
	offset      int
	value       string
	replacement ds.Optional[string]
	variables   []TestVar
}

type TestVar struct {
	key   string
	value string
}

func singleMatch(t *testing.T, results engine.Matches, startOffset int, value string) {
	t.Helper()
	if len(results) < 1 {
		t.FailNow()
	}
	if len(results) > 1 {
		t.Fail()
	}

	match := results[0]
	if match.Value != value || match.Offset.Start != startOffset {
		t.FailNow()
	}
}

func matches(t *testing.T, results engine.Matches, expected []TestMatch) {
	t.Helper()
	if len(results) != len(expected) {
		t.Errorf("Expected %d results, got %d results\n", len(expected), len(results))
		t.FailNow()
	}

	for i, e := range expected {
		actual := results[i]
		if actual.Value != e.value {
			t.Errorf("Expected value %s, got %s\n", e.value, actual.Value)
		}
		if actual.Offset.Start != e.offset {
			t.Errorf("Expected offset %d, got %d", e.offset, actual.Offset.Start)
		}
		if actual.Replacement != e.replacement {
			t.Errorf("Expected replacement %s, got %s\n", e.replacement.GetValueOrDefault("NONE OPTIONAL VALUE"), actual.Replacement.GetValueOrDefault("NONE OPTIONAL VALUE"))
		}
		if actual.Variables.Len() != len(e.variables) {
			t.Errorf("Expected %d variables, got %d variables\n", len(e.variables), actual.Variables.Len())
		} else {
			for _, exVar := range e.variables {
				v, prs := actual.Variables.Get(exVar.key)
				if prs && v.String().Value != exVar.value {
					t.Errorf("Expected %s, got %s\n", exVar.value, v.String().Value)
				}
			}
		}
	}
}

func checkVoreError(t *testing.T, err error, expectedType reflect.Type, expectedMessage string) {
	if err == nil {
		t.Errorf("Did not return any error :(")
		t.FailNow()
	}

	if reflect.TypeOf(err) != expectedType {
		t.Errorf("Expected %s but got %s", expectedType, reflect.TypeOf(err))
	}
	if err.Error() != expectedMessage {
		t.Errorf("Expected message '%s' but got '%s'", expectedMessage, err.Error())
	}
}
