package ast

import (
	"fmt"
	"strings"
	"testing"
)

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

	var erroredToken *Token
	var erri any = err
	switch detailedErr := erri.(type) {
	case *LexError:
		if expectedType != "LexError" {
			t.Errorf("Expected %s but got %s", expectedType, "LexError")
		}
		erroredToken = detailedErr.Token()
	case *ParseError:
		if expectedType != "ParseError" {
			t.Errorf("Expected %s but got %s", expectedType, "ParseError")
		}
		erroredToken = detailedErr.Token()
	default:
		t.Errorf("Expected %s but got %T", expectedType, erri)
	}

	if erroredToken.TokenType != expectedTokenType {
		t.Errorf("Expected tokenType %s but got %s", expectedTokenType.PP(), erroredToken.TokenType.PP())
	}
	if erroredToken.Lexeme != expectedLexeme {
		t.Errorf("Expected lexeme '%s' but got '%s'", expectedLexeme, erroredToken.Lexeme)
	}
	if erroredToken.Offset.Start != expectedOffsetStart && erroredToken.Offset.End != expectedOffsetEnd {
		t.Errorf("Expected range (%d, %d) but got (%d, %d)", expectedOffsetStart, expectedOffsetEnd, erroredToken.Offset.Start, erroredToken.Offset.End)
	}
	expectedMessageFixed := fmt.Sprintf("%s: %s", expectedType, expectedMessage)
	if !strings.HasPrefix(err.Error(), expectedMessageFixed) {
		t.Errorf("Expected message '%s' but got '%s'", expectedMessageFixed, err.Error())
	}
}
