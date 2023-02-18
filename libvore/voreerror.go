package libvore

import "fmt"

type VoreError struct {
	ErrorType string
	Token     *Token
	Message   string
}

func (v *VoreError) Error() string {
	switch v.ErrorType {
	case "LexError":
		return fmt.Sprintf("LEX ERROR: %s\nToken: '%s'\nLine:   %d - %d\nColumn: %d - %d", v.Message, v.Token.Lexeme, v.Token.Line.Start, v.Token.Line.End, v.Token.Column.Start, v.Token.Column.End)
	case "FileError":
		return v.Message
	case "ParseError":
		fmt.Printf("%s\n", v.Message)
		fmt.Printf("%s\n", v.Token.Lexeme)
		return fmt.Sprintf("ERROR:  %s\nToken:  '%s'\nTokenType: %s\nLine:   %d - %d\nColumn: %d - %d", v.Message, v.Token.Lexeme, v.Token.TokenType.PP(), v.Token.Line.Start, v.Token.Line.End, v.Token.Column.Start, v.Token.Column.End)
	case "GenError":
		return v.Message
	}
	return fmt.Sprintf("UNKNOWN ERROR :( %s", v.Message)
}

func NewLexErrorUnknown(t *Token) *VoreError {
	return &VoreError{
		ErrorType: "LexError",
		Token:     t,
		Message:   "Unknown token :(",
	}
}

func NewLexErrorCustomMsg(t *Token, message string) *VoreError {
	return &VoreError{
		ErrorType: "LexError",
		Token:     t,
		Message:   message,
	}
}

func NewParseError(t Token, message string) *VoreError {
	return &VoreError{
		ErrorType: "ParseError",
		Token:     &t,
		Message:   message,
	}
}

func NewGenError(err error) *VoreError {
	return &VoreError{
		ErrorType: "ParseError",
		Token:     nil,
		Message:   err.Error(),
	}
}

func NewFileError(err error) *VoreError {
	return &VoreError{
		ErrorType: "FileError",
		Token:     nil,
		Message:   err.Error(),
	}
}
