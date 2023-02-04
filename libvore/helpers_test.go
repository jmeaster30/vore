package libvore

import "testing"

type TestMatch struct {
	offset      int
	value       string
	replacement string
	variables   []TestVar
}

type TestVar struct {
	key   string
	value string
}

func singleMatch(t *testing.T, results Matches, startOffset int, value string) {
	t.Helper()
	if len(results) < 1 {
		t.FailNow()
	}
	if len(results) > 1 {
		t.Fail()
	}

	match := results[0]
	if match.value != value || match.offset.Start != startOffset {
		t.FailNow()
	}
}

func matches(t *testing.T, results Matches, expected []TestMatch) {
	t.Helper()
	if len(results) != len(expected) {
		t.Errorf("Expected %d results, got %d results\n", len(expected), len(results))
		t.FailNow()
	}

	for i, e := range expected {
		actual := results[i]
		if actual.value != e.value {
			t.Logf("Expected value %s, got %s\n", e.value, actual.value)
		}
		if actual.offset.Start != e.offset {
			t.Logf("Expected offset %d, got %d", e.offset, actual.offset.Start)
		}
		if actual.replacement != e.replacement {
			t.Logf("Expected replacement %s, got %s\n", e.replacement, actual.replacement)
		}
		if actual.variables.Len() != len(e.variables) {
			t.Errorf("Expected %d variables, got %d variables\n", len(e.variables), actual.variables.Len())
		} else {
			for _, exVar := range e.variables {
				v, prs := actual.variables.Get(exVar.key)
				if prs && v.String().Value != exVar.value {
					t.Errorf("Expected %s, got %s\n", exVar.value, v.String().Value)
				}
			}
		}
	}
}

func checkNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Error(err)
	}
}

func mustPanic(t *testing.T, process func(*testing.T)) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	process(t)
}
