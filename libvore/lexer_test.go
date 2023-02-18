package libvore

import (
	"strings"
	"testing"
)

func tokenList(t *testing.T, results []*Token, expected []*Token) {
	t.Helper()
	if len(results) != len(expected) {
		t.Errorf("Expected %d results, got %d results\n", len(expected), len(results))
		t.FailNow()
	}

	for i, e := range expected {
		actual := results[i]
		if actual.TokenType != e.TokenType || actual.Lexeme != e.Lexeme {
			t.Logf("Expected token type %d, got %d\nExpected lexeme [%s], got [%s]\n",
				e.TokenType, actual.TokenType,
				e.Lexeme, actual.Lexeme)
		}
		if actual.Offset.Start != e.Offset.Start || actual.Offset.End != e.Offset.End {
			t.Errorf("Expected offset (%d, %d), got offset (%d, %d)\n", e.Offset.Start, e.Offset.End, actual.Offset.Start, actual.Offset.End)
		}
	}
}

func TestLexerBasic(t *testing.T) {
	lexer := initLexer(strings.NewReader("ident 'string1' \"string2\""))
	actual, err := lexer.getTokens()
	checkNoError(t, err)
	tokenList(t, actual, []*Token{
		{IDENTIFIER, NewRange(0, 5), NewRange(1, 1), NewRange(0, 5), "ident"},
		{WS, NewRange(5, 6), NewRange(1, 1), NewRange(5, 6), " "},
		{STRING, NewRange(6, 15), NewRange(1, 1), NewRange(6, 15), "string1"},
		{WS, NewRange(15, 16), NewRange(1, 1), NewRange(15, 16), " "},
		{STRING, NewRange(16, 25), NewRange(1, 1), NewRange(16, 25), "string2"},
		{EOF, NewRange(25, 25), NewRange(1, 1), NewRange(25, 25), ""},
	})
}

func TestCheckUnendingStringError(t *testing.T) {
	lexer := initLexer(strings.NewReader("ident 'testing"))
	tokens, err := lexer.getTokens()

	checkVoreErrorToken(t, err, "LexError", ERROR, "testing", 6, 14, "Unending string")

	if len(tokens) != 0 {
		t.Errorf("Expected no tokens returned on error. Got %d tokens", len(tokens))
	}
}

func TestCheckUnendingBlockCommentError(t *testing.T) {
	lexer := initLexer(strings.NewReader("ident --(test comment"))
	tokens, err := lexer.getTokens()

	checkVoreErrorToken(t, err, "LexError", ERROR, "--(test comment", 6, 21, "Unending block comment")

	if len(tokens) != 0 {
		t.Errorf("Expected no tokens returned on error. Got %d tokens", len(tokens))
	}
}

func TestCheckUnendingBlockCommentError2(t *testing.T) {
	lexer := initLexer(strings.NewReader("ident --(test comment)"))
	tokens, err := lexer.getTokens()

	checkVoreErrorToken(t, err, "LexError", ERROR, "--(test comment)", 6, 22, "Unending block comment")

	if len(tokens) != 0 {
		t.Errorf("Expected no tokens returned on error. Got %d tokens", len(tokens))
	}
}

func TestCheckUnknownToken(t *testing.T) {
	lexer := initLexer(strings.NewReader("ident @"))
	tokens, err := lexer.getTokens()

	checkVoreErrorToken(t, err, "LexError", ERROR, "@", 6, 7, "Unknown token :(")

	if len(tokens) != 0 {
		t.Errorf("Expected no tokens returned on error. Got %d tokens", len(tokens))
	}
}

func ppMatch(t *testing.T, a TokenType, b string) {
	if a.PP() != b {
		t.Errorf("%s != %s", a.PP(), b)
	}
}

func TestTokenTypePP(t *testing.T) {
	ppMatch(t, ERROR, "ERROR")
	ppMatch(t, EOF, "EOF")
	ppMatch(t, WS, "WS")
	ppMatch(t, COMMENT, "COMMENT")
	ppMatch(t, IDENTIFIER, "IDENTIFIER")
	ppMatch(t, NUMBER, "NUMBER")
	ppMatch(t, STRING, "STRING")
	ppMatch(t, EQUAL, "EQUAL")
	ppMatch(t, COLONEQ, "COLONEQ")
	ppMatch(t, COMMA, "COMMA")
	ppMatch(t, OPENPAREN, "OPENPAREN")
	ppMatch(t, CLOSEPAREN, "CLOSEPAREN")
	ppMatch(t, OPENCURLY, "OPENCURLY")
	ppMatch(t, CLOSECURLY, "CLOSECURLY")
	ppMatch(t, FIND, "FIND")
	ppMatch(t, REPLACE, "REPLACE")
	ppMatch(t, WITH, "WITH")
	ppMatch(t, SET, "SET")
	ppMatch(t, TO, "TO")
	ppMatch(t, PATTERN, "PATTERN")
	ppMatch(t, MATCHES, "MATCHES")
	ppMatch(t, TRANSFORM, "TRANSFORM")
	ppMatch(t, ALL, "ALL")
	ppMatch(t, SKIP, "SKIP")
	ppMatch(t, TAKE, "TAKE")
	ppMatch(t, TOP, "TOP")
	ppMatch(t, LAST, "LAST")
	ppMatch(t, ANY, "ANY")
	ppMatch(t, WHITESPACE, "WHITESPACE")
	ppMatch(t, DIGIT, "DIGIT")
	ppMatch(t, UPPER, "UPPER")
	ppMatch(t, LOWER, "LOWER")
	ppMatch(t, LETTER, "LETTER")
	ppMatch(t, LINE, "LINE")
	ppMatch(t, START, "START")
	ppMatch(t, FILE, "FILE")
	ppMatch(t, WORD, "WORD")
	ppMatch(t, END, "END")
	ppMatch(t, BEGIN, "BEGIN")
	ppMatch(t, NOT, "NOT")
	ppMatch(t, AT, "AT")
	ppMatch(t, MOST, "MOST")
	ppMatch(t, LEAST, "LEAST")
	ppMatch(t, BETWEEN, "BETWEEN")
	ppMatch(t, AND, "AND")
	ppMatch(t, EXACTLY, "EXACTLY")
	ppMatch(t, MAYBE, "MAYBE")
	ppMatch(t, FEWEST, "FEWEST")
	ppMatch(t, NAMED, "NAMED")
	ppMatch(t, IN, "IN")
	ppMatch(t, OR, "OR")
	ppMatch(t, IF, "IF")
	ppMatch(t, THEN, "THEN")
	ppMatch(t, ELSE, "ELSE")
	ppMatch(t, DEBUG, "DEBUG")
	ppMatch(t, RETURN, "RETURN")
	ppMatch(t, HEAD, "HEAD")
	ppMatch(t, TAIL, "TAIL")
	ppMatch(t, LOOP, "LOOP")
	ppMatch(t, PLUS, "PLUS")
	ppMatch(t, MINUS, "MINUS")
	ppMatch(t, MULT, "MULT")
	ppMatch(t, DIV, "DIV")
	ppMatch(t, MOD, "MOD")
	ppMatch(t, LESS, "LESS")
	ppMatch(t, GREATER, "GREATER")
	ppMatch(t, LESSEQ, "LESSEQ")
	ppMatch(t, GREATEREQ, "GREATEREQ")
	ppMatch(t, DEQUAL, "DEQUAL")
	ppMatch(t, NEQUAL, "NEQUAL")
	ppMatch(t, CONTINUE, "CONTINUE")
	ppMatch(t, BREAK, "BREAK")
}
