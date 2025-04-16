package ast

import "fmt"

type ParseError struct {
	Token   *Token
	message string
}

func (err *ParseError) Error() string {
	return fmt.Sprintf("ParseError:  %s\nToken:  '%s'\nTokenType: %s\nLine:   %d - %d\nColumn: %d - %d", err.message, err.Token.Lexeme, err.Token.TokenType.PP(), err.Token.Line.Start, err.Token.Line.End, err.Token.Column.Start, err.Token.Column.End)
}

func NewParseError(token *Token, message string) *ParseError {
	return &ParseError{token, message}
}
