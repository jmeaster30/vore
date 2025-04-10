package ast

import "fmt"

type LexError struct {
	token   *Token
	message string
}

func (err *LexError) Error() string {
	return fmt.Sprintf("LEX ERROR: %s\nToken: '%s'\nLine:   %d - %d\nColumn: %d - %d", err.message, err.token.Lexeme, err.token.Line.Start, err.token.Line.End, err.token.Column.Start, err.token.Column.End)
}

func NewLexError(token *Token, message string) *LexError {
	return &LexError{token, message}
}
