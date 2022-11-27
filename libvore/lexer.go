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
	COMMENT

	// literals
	IDENTIFIER
	NUMBER
	STRING

	// misc
	EQUAL
	COLONEQ
	COMMA
	OPENPAREN
	CLOSEPAREN

	// commands
	FIND
	REPLACE
	WITH
	SET
	TO

	// result length
	ALL
	SKIP
	TAKE
	TOP
	LAST

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
	MAYBE
	FEWEST
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
	case COMMENT:
		return "COMMENT"
	case IDENTIFIER:
		return "IDENTIFIER"
	case NUMBER:
		return "NUMBER"
	case STRING:
		return "STRING"
	case EQUAL:
		return "EQUAL"
	case COLONEQ:
		return "COLONEQ"
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
	case WITH:
		return "WITH"
	case SET:
		return "SET"
	case TO:
		return "TO"
	case ALL:
		return "ALL"
	case SKIP:
		return "SKIP"
	case TAKE:
		return "TAKE"
	case TOP:
		return "TOP"
	case LAST:
		return "LAST"
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
	case AND:
		return "AND"
	case EXACTLY:
		return "EXACTLY"
	case MAYBE:
		return "MAYBE"
	case FEWEST:
		return "FEWEST"
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
	offset   int
	line     int
	column   int
	lastRead rune
}

type Lexer struct {
	r           *bufio.Reader
	currentChar rune
	position    *Stack[PositionInfo]
}

func initLexer(r io.Reader) *Lexer {
	stack := NewStack[PositionInfo]()
	stack.Push(PositionInfo{
		offset:   0,
		line:     1,
		column:   1,
		lastRead: rune(0),
	})

	lexer := Lexer{r: bufio.NewReader(r), currentChar: rune(0), position: stack}
	return &lexer
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
		SSTRING_D_ESCAPE
		SSTRING_S_ESCAPE
		SNUMBER
		SEQUAL
		SCOLON
		SCOLONEQ
		SIDENTIFIER
		SCOMMA
		SOPENPAREN
		SCLOSEPAREN
		SCOMMENT
		SDASH
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
		} else if current_state == SCOMMENT {
			if ch == '\n' {
				s.unread_last()
				break
			}
			buf.WriteRune(ch)
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
		} else if ch == '=' && current_state == SSTART {
			buf.WriteRune(ch)
			current_state = SEQUAL
			break
		} else if ch == '=' && current_state == SCOLON {
			buf.WriteRune(ch)
			current_state = SCOLONEQ
			break
		} else if ch == ':' && current_state == SSTART {
			buf.WriteRune(ch)
			current_state = SCOLON
		} else if ch == '-' && (current_state == SSTART || current_state == SDASH) {
			buf.WriteRune(ch)
			if current_state == SSTART {
				current_state = SDASH
			} else if current_state == SDASH {
				current_state = SCOMMENT
			}
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
		} else if ch == '\\' && current_state == SSTRING_DOUBLE {
			current_state = SSTRING_D_ESCAPE
		} else if current_state == SSTRING_D_ESCAPE {
			buf.WriteRune(getEscapedRune(ch))
			current_state = SSTRING_DOUBLE
		} else if current_state == SSTRING_DOUBLE {
			if ch == '"' {
				break
			}
			buf.WriteRune(ch)
		} else if ch == '\'' && current_state == SSTART {
			current_state = SSTRING_SINGLE
		} else if ch == '\\' && current_state == SSTRING_SINGLE {
			current_state = SSTRING_S_ESCAPE
		} else if current_state == SSTRING_S_ESCAPE {
			buf.WriteRune(getEscapedRune(ch))
			current_state = SSTRING_SINGLE
		} else if current_state == SSTRING_SINGLE {
			if ch == '\'' {
				break
			}
			buf.WriteRune(ch)
		} else {
			if current_state != SSTART || unicode.IsDigit(ch) || unicode.IsLetter(ch) || unicode.IsSpace(ch) || ch == '(' || ch == ')' || ch == ',' || ch == ':' || ch == '=' || ch == '"' || ch == '\'' || ch == '-' {
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
	case SSTRING_SINGLE:
		fallthrough
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
		case "with":
			token.tokenType = WITH
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
		case "top":
			token.tokenType = TOP
		case "last":
			token.tokenType = LAST
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
		case "maybe":
			token.tokenType = MAYBE
		case "fewest":
			token.tokenType = FEWEST
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
	case SEQUAL:
		token.tokenType = EQUAL
	case SCOLON:
		token.tokenType = ERROR
	case SCOMMENT:
		token.tokenType = COMMENT
	case SDASH:
		token.tokenType = ERROR
	case SCOLONEQ:
		token.tokenType = COLONEQ
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

func getEscapedRune(ch rune) rune {
	if ch == 'n' {
		return rune(10)
	} else if ch == 't' {
		return rune(9)
	} else if ch == 'r' {
		return rune(13)
	}
	return rune(ch)
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
	posInfo.line = s.get_position().line
	if s.currentChar == '\n' {
		posInfo.line += 1
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
