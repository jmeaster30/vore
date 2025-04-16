package ast

import "fmt"

type LexError struct {
	Token   *Token
	message string
}

func (err *LexError) Error() string {
	return fmt.Sprintf("LexError: %s\nToken: '%s'\nLine:   %d - %d\nColumn: %d - %d", err.message, err.Token.Lexeme, err.Token.Line.Start, err.Token.Line.End, err.Token.Column.Start, err.Token.Column.End)
}

func NewLexError(token *Token, message string) *LexError {
	return &LexError{token, message}
}
