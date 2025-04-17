package ast

import "fmt"

type ParseError struct {
	token   *Token
	message string
}

func (err *ParseError) Error() string {
	return fmt.Sprintf("ParseError:  %s\nToken:  '%s'\nTokenType: %s\nLine:   %d - %d\nColumn: %d - %d", err.message, err.token.Lexeme, err.token.TokenType.PP(), err.token.Line.Start, err.token.Line.End, err.token.Column.Start, err.token.Column.End)
}

func (err *ParseError) Token() *Token {
	return err.token
}

func (err *ParseError) Message() string {
	return err.message
}

func NewParseError(token *Token, message string) *ParseError {
	return &ParseError{token, message}
}
