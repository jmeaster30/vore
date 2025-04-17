package ast

import "fmt"

type LexError struct {
	token   *Token
	message string
}

func (err *LexError) Error() string {
	return fmt.Sprintf("LexError: %s\nToken: '%s'\nLine:   %d - %d\nColumn: %d - %d", err.message, err.token.Lexeme, err.token.Line.Start, err.token.Line.End, err.token.Column.Start, err.token.Column.End)
}

func (err *LexError) Token() *Token {
	return err.token
}

func (err *LexError) Message() string {
	return err.message
}

func NewLexError(token *Token, message string) *LexError {
	return &LexError{token, message}
}
