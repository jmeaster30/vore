package libvore

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/jmeaster30/vore/libvore/ds"
	"github.com/jmeaster30/vore/libvore/engine"
	"github.com/jmeaster30/vore/libvore/testutils"
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
		testutils.AssertEqual(t, e.value, actual.Value)
		testutils.AssertEqual(t, e.offset, actual.Offset.Start)
		testutils.AssertEqual(t, e.replacement, actual.Replacement)
		testutils.AssertEqual(t, len(e.variables), len(actual.Variables.Map()))
		for _, expectedVar := range e.variables {
			v, prs := actual.Variables.Get(expectedVar.key)
			testutils.AssertTrue(t, prs)
			testutils.AssertEqual(t, expectedVar.value, v.String())
		}
	}
}

func checkVoreError(t *testing.T, err error, expectedType string, expectedMessage string) {
	t.Helper()
	if err == nil {
		t.Errorf("Did not return any error :(")
		t.FailNow()
	}

	if strings.HasSuffix(expectedType, reflect.TypeOf(err).String()) {
		t.Errorf("Expected %s but got %s", expectedType, reflect.TypeOf(err).String())
	}
	expectedMessageFixed := fmt.Sprintf("%s: %s", expectedType, expectedMessage)
	if !strings.HasPrefix(err.Error(), expectedMessageFixed) {
		t.Errorf("Expected message '%s' but got '%s'", expectedMessageFixed, err.Error())
	}
}
