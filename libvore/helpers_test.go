package libvore

import (
	"testing"

	"github.com/jmeaster30/vore/libvore/ds"
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

func singleMatch(t *testing.T, results Matches, startOffset int, value string) {
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

func matches(t *testing.T, results Matches, expected []TestMatch) {
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

func checkVoreError(t *testing.T, err error, expectedType string, expectedMessage string) {
	if err == nil {
		t.Errorf("Did not return any error :(")
		t.FailNow()
	}

	if detailedErr, ok := err.(*VoreError); ok {
		if detailedErr.ErrorType != expectedType {
			t.Errorf("Expected %s but got %s", expectedType, detailedErr.ErrorType)
		}
		if detailedErr.Message != expectedMessage {
			t.Errorf("Expected message '%s' but got '%s'", expectedMessage, detailedErr.Message)
		}
	} else {
		t.Errorf("Expected VoreError returned but got some other error. %s", err.Error())
	}
}

func checkVoreErrorToken(t *testing.T,
	err error,
	expectedType string,
	expectedTokenType TokenType,
	expectedLexeme string,
	expectedOffsetStart int,
	expectedOffsetEnd int,
	expectedMessage string,
) {
	if err == nil {
		t.Errorf("Did not return any error :(")
		t.FailNow()
	}

	if detailedErr, ok := err.(*VoreError); ok {
		if detailedErr.ErrorType != expectedType {
			t.Errorf("Expected %s but got %s", expectedType, detailedErr.ErrorType)
		}
		if detailedErr.Token.TokenType != expectedTokenType {
			t.Errorf("Expected tokenType %s but got %s", expectedTokenType.PP(), detailedErr.Token.TokenType.PP())
		}
		if detailedErr.Token.Lexeme != expectedLexeme {
			t.Errorf("Expected lexeme '%s' but got '%s'", expectedLexeme, detailedErr.Token.Lexeme)
		}
		if detailedErr.Token.Offset.Start != expectedOffsetStart && detailedErr.Token.Offset.End != expectedOffsetEnd {
			t.Errorf("Expected range (%d, %d) but got (%d, %d)", expectedOffsetStart, expectedOffsetEnd, detailedErr.Token.Offset.Start, detailedErr.Token.Offset.End)
		}
		if detailedErr.Message != expectedMessage {
			t.Errorf("Expected message '%s' but got '%s'", expectedMessage, detailedErr.Message)
		}
	} else {
		t.Errorf("Expected VoreError returned but got some other error. %s", err.Error())
	}
}
