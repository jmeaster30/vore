package libvore

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type TokenType int

const (
	// special tokens
	ERROR TokenType = iota
	EOF
	WS

	// literals
	IDENTIFIER
	NUMBER
	STRING

	// misc
	COMMA
	OPENPAREN
	CLOSEPAREN

	// commands
	FIND
	REPLACE
	SET
	TO

	// result length
	ALL
	SKIP
	TAKE

	// classes
	ANY
	WHITESPACE
	DIGIT
	UPPER
	LOWER
	LETTER
	LINE
	FILE
	START
	END

	// keywords
	NOT
	AT
	LEAST
	MOST
	BETWEEN
	AND
	EXACTLY
	IN
	OR
)

func (t TokenType) pp() string {
	switch t {
	case ERROR:
		return "ERROR"
	case EOF:
		return "EOF"
	case WS:
		return "WS"
	case IDENTIFIER:
		return "IDENTIFIER"
	case NUMBER:
		return "NUMBER"
	case STRING:
		return "STRING"
	case COMMA:
		return "COMMA"
	case OPENPAREN:
		return "OPENPAREN"
	case CLOSEPAREN:
		return "CLOSEPAREN"
	case FIND:
		return "FIND"
	case REPLACE:
		return "REPLACE"
	case ALL:
		return "ALL"
	case SKIP:
		return "SKIP"
	case TAKE:
		return "TAKE"
	case ANY:
		return "ANY"
	case WHITESPACE:
		return "WHITESPACE"
	case DIGIT:
		return "DIGIT"
	case UPPER:
		return "UPPER"
	case LOWER:
		return "LOWER"
	case LETTER:
		return "LETTER"
	case LINE:
		return "LINE"
	case START:
		return "START"
	case FILE:
		return "FILE"
	case END:
		return "END"
	case NOT:
		return "NOT"
	case AT:
		return "AT"
	case MOST:
		return "MOST"
	case LEAST:
		return "LEAST"
	case BETWEEN:
		return "BETWEEN"
	case EXACTLY:
		return "EXACTLY"
	case IN:
		return "IN"
	case OR:
		return "OR"
	default:
		panic("UNKNOWN TOKEN TYPE")
	}
}

type Token struct {
	tokenType TokenType
	offset    *Range
	line      *Range
	column    *Range
	lexeme    string
}

type PositionInfo struct {
	offset   uint64
	line     uint64
	column   uint64
	lastRead rune
}

type Lexer struct {
	r           *bufio.Reader
	currentChar rune
	position    Stack[PositionInfo]
}

func initLexer(r io.Reader) *Lexer {
	lexer := &Lexer{r: bufio.NewReader(r), currentChar: rune(0), position: NewStack[PositionInfo]()}
	lexer.position.Push(PositionInfo{
		offset:   0,
		line:     1,
		column:   0,
		lastRead: rune(0),
	})
	return lexer
}

func (s *Lexer) getTokens() []*Token {
	tokens := []*Token{}
	for {
		token := s.getNextToken()
		tokens = append(tokens, token)
		if token.tokenType == EOF {
			break
		}
	}
	return tokens
}

func (s *Lexer) getNextToken() *Token {
	type TokenState int

	const (
		SSTART TokenState = iota
		SWHITESPACE
		SSTRING_DOUBLE
		SSTRING_SINGLE
		SNUMBER
		SIDENTIFIER
		SCOMMA
		SOPENPAREN
		SCLOSEPAREN
		SERROR
		SEND
	)

	current_state := SSTART
	token := &Token{}
	startPosInfo := s.get_position()

	var buf bytes.Buffer

	for {
		ch := s.read()
		//fmt.Printf("%d - %d\n", int(ch), current_state)
		if ch == 0 && current_state == SSTART {
			current_state = SEND
			break
		} else if ch == 0 {
			s.unread_last()
			break
		} else if ch == '(' && current_state == SSTART {
			buf.WriteRune(ch)
			current_state = SOPENPAREN
			break
		} else if ch == ')' && current_state == SSTART {
			buf.WriteRune(ch)
			current_state = SCLOSEPAREN
			break
		} else if ch == ',' && current_state == SSTART {
			buf.WriteRune(ch)
			current_state = SCOMMA
			break
		} else if unicode.IsSpace(ch) {
			if current_state == SSTART || current_state == SWHITESPACE {
				current_state = SWHITESPACE
				buf.WriteRune(ch)
			} else {
				s.unread_last()
				break
			}
		} else if unicode.IsDigit(ch) && (current_state == SNUMBER || current_state == SSTART) {
			current_state = SNUMBER
			buf.WriteRune(ch)
		} else if unicode.IsLetter(ch) && current_state == SSTART {
			current_state = SIDENTIFIER
			buf.WriteRune(ch)
		} else if (unicode.IsDigit(ch) || unicode.IsLetter(ch)) && current_state == SIDENTIFIER {
			current_state = SIDENTIFIER
			buf.WriteRune(ch)
		} else if ch == '"' && current_state == SSTART {
			current_state = SSTRING_DOUBLE
			buf.WriteRune(ch)
		} else if current_state == SSTRING_DOUBLE {
			buf.WriteRune(ch)
			if ch == '"' {
				break
			}
		} else {
			if current_state != SSTART || unicode.IsDigit(ch) || unicode.IsLetter(ch) || unicode.IsSpace(ch) || ch == '(' || ch == ')' || ch == ',' {
				s.unread_last()
			} else {
				buf.WriteRune(ch)
				current_state = SERROR
			}
			break
		}
	}

	switch current_state {
	case SERROR:
		token.tokenType = ERROR
	case SSTRING_DOUBLE:
		token.tokenType = STRING
	case SNUMBER:
		token.tokenType = NUMBER
	case SIDENTIFIER:
		token.tokenType = IDENTIFIER
		lexeme := strings.ToLower(buf.String())
		switch lexeme {
		case "find":
			token.tokenType = FIND
		case "replace":
			token.tokenType = REPLACE
		case "set":
			token.tokenType = SET
		case "to":
			token.tokenType = TO
		case "all":
			token.tokenType = ALL
		case "skip":
			token.tokenType = SKIP
		case "take":
			token.tokenType = TAKE
		case "any":
			token.tokenType = ANY
		case "whitespace":
			token.tokenType = WHITESPACE
		case "digit":
			token.tokenType = DIGIT
		case "upper":
			token.tokenType = UPPER
		case "lower":
			token.tokenType = LOWER
		case "letter":
			token.tokenType = LETTER
		case "line":
			token.tokenType = LINE
		case "file":
			token.tokenType = FILE
		case "start":
			token.tokenType = START
		case "end":
			token.tokenType = END
		case "not":
			token.tokenType = NOT
		case "at":
			token.tokenType = AT
		case "least":
			token.tokenType = LEAST
		case "most":
			token.tokenType = MOST
		case "between":
			token.tokenType = BETWEEN
		case "and":
			token.tokenType = AND
		case "exactly":
			token.tokenType = EXACTLY
		case "in":
			token.tokenType = IN
		case "or":
			token.tokenType = OR
		}
	case SWHITESPACE:
		token.tokenType = WS
	case SOPENPAREN:
		token.tokenType = OPENPAREN
	case SCLOSEPAREN:
		token.tokenType = CLOSEPAREN
	case SCOMMA:
		token.tokenType = COMMA
	case SEND:
		token.tokenType = EOF
	default:
		fmt.Println(current_state)
		panic("Unknown final state")
	}

	endPosInfo := s.get_position()
	token.offset = NewRange(startPosInfo.offset, endPosInfo.offset)
	token.column = NewRange(startPosInfo.column, endPosInfo.column)
	token.line = NewRange(startPosInfo.line, endPosInfo.line)
	token.lexeme = buf.String()
	return token
}

func (s *Lexer) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return rune(0)
	}
	posInfo := PositionInfo{}
	posInfo.lastRead = s.currentChar
	posInfo.offset = s.get_position().offset + 1
	posInfo.column = s.get_position().column + 1
	if s.currentChar == '\n' {
		posInfo.line = s.get_position().line + 1
		posInfo.column = 1
	}
	s.position.Push(posInfo)
	s.currentChar = ch
	return ch
}

func (s *Lexer) get_position() *PositionInfo {
	return s.position.Peek()
}

func (s *Lexer) unread_last() {
	s.unread(1)
}

func (s *Lexer) unread(amount uint64) {
	if amount >= s.position.Size() {
		panic("You can't pop that much!!!")
	}
	var lastPopped *PositionInfo
	for i := uint64(0); i < amount; i++ {
		_ = s.r.UnreadRune()
		lastPopped = s.position.Pop()
	}
	s.currentChar = lastPopped.lastRead
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
