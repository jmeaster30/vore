package libvore

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
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
	OPENCURLY
	CLOSECURLY
	PLUS
	MINUS
	MULT
	DIV
	LESS
	GREATER
	LESSEQ
	GREATEREQ
	DEQUAL
	NEQUAL
	MOD

	// commands
	FIND
	REPLACE
	WITH
	SET
	TO
	PATTERN
	MATCHES
	TRANSFORM

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
	BEGIN

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
	IF
	THEN
	ELSE
	DEBUG
	RETURN
	HEAD
	TAIL
	LOOP
	BREAK
	CONTINUE
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
	case OPENCURLY:
		return "OPENCURLY"
	case CLOSECURLY:
		return "CLOSECURLY"
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
	case PATTERN:
		return "PATTERN"
	case MATCHES:
		return "MATCHES"
	case TRANSFORM:
		return "TRANSFORM"
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
	case BEGIN:
		return "BEGIN"
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
	case IF:
		return "IF"
	case THEN:
		return "THEN"
	case ELSE:
		return "ELSE"
	case DEBUG:
		return "DEBUG"
	case RETURN:
		return "RETURN"
	case HEAD:
		return "HEAD"
	case TAIL:
		return "TAIL"
	case LOOP:
		return "LOOP"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case MULT:
		return "MULT"
	case DIV:
		return "DIV"
	case MOD:
		return "MOD"
	case LESS:
		return "LESS"
	case GREATER:
		return "GREATER"
	case LESSEQ:
		return "LESSEQ"
	case GREATEREQ:
		return "GREATEREQ"
	case DEQUAL:
		return "DEQUAL"
	case NEQUAL:
		return "NEQUAL"
	case CONTINUE:
		return "CONTINUE"
	case BREAK:
		return "BREAK"
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

func (token Token) print() {
	fmt.Printf("[%s] '%s' \tline: %d, \tstart column: %d, \tend column: %d\n", token.tokenType.pp(), token.lexeme, token.line.Start, token.column.Start, token.column.End)
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
		SEQUAL_1
		SDEQUAL
		SEXCL
		SNEQUAL
		SCOLON
		SCOLONEQ
		SIDENTIFIER
		SCOMMA
		SOPENPAREN
		SCLOSEPAREN
		SOPENCURLY
		SCLOSECURLY
		SCOMMENT
		SCOMMENTSTART
		SBLOCKCOMMENT
		SBLOCKCOMMENTSTARTEND
		SBLOCKCOMMENTENDEND
		SBLOCKCOMMENTFINAL
		SDASH
		SOPERATOR
		SOPERATORSTART
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
		} else if current_state == SBLOCKCOMMENT {
			buf.WriteRune(ch)
			if ch == ')' {
				current_state = SBLOCKCOMMENTSTARTEND
			}
		} else if current_state == SBLOCKCOMMENTSTARTEND && ch == '-' {
			buf.WriteRune(ch)
			current_state = SBLOCKCOMMENTENDEND
		} else if current_state == SBLOCKCOMMENTENDEND && ch == '-' {
			buf.WriteRune(ch)
			current_state = SBLOCKCOMMENTFINAL
			break
		} else if current_state == SBLOCKCOMMENTENDEND || current_state == SBLOCKCOMMENTSTARTEND {
			buf.WriteRune(ch)
			current_state = SBLOCKCOMMENT
		} else if ch == '\\' && current_state == SSTRING_DOUBLE {
			current_state = SSTRING_D_ESCAPE
		} else if current_state == SSTRING_DOUBLE {
			if ch == '"' {
				break
			}
			buf.WriteRune(ch)
		} else if ch == '\\' && current_state == SSTRING_SINGLE {
			current_state = SSTRING_S_ESCAPE
		} else if current_state == SSTRING_SINGLE {
			if ch == '\'' {
				break
			}
			buf.WriteRune(ch)
		} else if ch == '(' && current_state == SCOMMENTSTART {
			buf.WriteRune(ch)
			current_state = SBLOCKCOMMENT
		} else if current_state == SCOMMENTSTART {
			buf.WriteRune(ch)
			current_state = SCOMMENT
		} else if ch == '(' && current_state == SSTART {
			buf.WriteRune(ch)
			current_state = SOPENPAREN
			break
		} else if ch == ')' && current_state == SBLOCKCOMMENT {
			buf.WriteRune(ch)
			break
		} else if ch == ')' && current_state == SSTART {
			buf.WriteRune(ch)
			current_state = SCLOSEPAREN
			break
		} else if ch == '{' && current_state == SSTART {
			buf.WriteRune(ch)
			current_state = SOPENCURLY
			break
		} else if ch == '}' && current_state == SSTART {
			buf.WriteRune(ch)
			current_state = SCLOSECURLY
			break
		} else if ch == ',' && current_state == SSTART {
			buf.WriteRune(ch)
			current_state = SCOMMA
			break
		} else if ch == '!' && current_state == SSTART {
			buf.WriteRune(ch)
			current_state = SEXCL
		} else if ch == '=' && current_state == SEXCL {
			buf.WriteRune(ch)
			current_state = SNEQUAL
			break
		} else if ch == '=' && current_state == SSTART {
			buf.WriteRune(ch)
			current_state = SEQUAL_1
		} else if ch == '=' && current_state == SEQUAL_1 {
			buf.WriteRune(ch)
			current_state = SDEQUAL
			break
		} else if ch == '=' && current_state == SCOLON {
			buf.WriteRune(ch)
			current_state = SCOLONEQ
			break
		} else if ch == '=' && current_state == SOPERATORSTART {
			buf.WriteRune(ch)
			current_state = SOPERATOR
			break
		} else if ch == ':' && current_state == SSTART {
			buf.WriteRune(ch)
			current_state = SCOLON
		} else if ch == '-' && (current_state == SSTART || current_state == SDASH || current_state == SCOMMENTSTART) {
			buf.WriteRune(ch)
			if current_state == SSTART {
				current_state = SDASH
			} else if current_state == SDASH {
				current_state = SCOMMENTSTART
			} else if current_state == SCOMMENTSTART {
				current_state = SCOMMENT
			}
		} else if current_state == SSTART && (ch == '+' || ch == '%' || ch == '*' || ch == '/') {
			buf.WriteRune(ch)
			current_state = SOPERATOR
			break
		} else if current_state == SSTART && (ch == '>' || ch == '<') {
			buf.WriteRune(ch)
			current_state = SOPERATORSTART
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
		} else if current_state == SSTRING_D_ESCAPE {
			if ch == 'x' {
				next_ch := s.read()
				next_next_ch := s.read()
				if IsHex(next_ch) && IsHex(next_next_ch) {
					buf.WriteRune(HexToAscii(next_ch, next_next_ch))
				} else {
					s.unread(2)
					buf.WriteRune('x')
				}
			} else {
				buf.WriteRune(getEscapedRune(ch))
			}
			current_state = SSTRING_DOUBLE
		} else if ch == '\'' && current_state == SSTART {
			current_state = SSTRING_SINGLE
		} else if current_state == SSTRING_S_ESCAPE {
			if ch == 'x' {
				next_ch := s.read()
				next_next_ch := s.read()
				if IsHex(next_ch) && IsHex(next_next_ch) {
					buf.WriteRune(HexToAscii(next_ch, next_next_ch))
				} else {
					s.unread(2)
					buf.WriteRune('x')
				}
			} else {
				buf.WriteRune(getEscapedRune(ch))
			}
			current_state = SSTRING_SINGLE
		} else if current_state == SCOMMENTSTART {
			s.unread_last()
			current_state = SCOMMENT
		} else {
			if current_state != SSTART || unicode.IsDigit(ch) || unicode.IsLetter(ch) || unicode.IsSpace(ch) || ch == '(' || ch == ')' || ch == '{' || ch == '}' || ch == ',' || ch == ':' || ch == '=' || ch == '"' || ch == '\'' || ch == '-' || ch == '+' || ch == '<' || ch == '>' || ch == '*' || ch == '/' || ch == '%' {
				s.unread_last()
			} else {
				buf.WriteRune(ch)
				current_state = SERROR
			}
			break
		}
	}

	token.tokenType = ERROR

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
		case "pattern":
			token.tokenType = PATTERN
		case "matches":
			token.tokenType = MATCHES
		case "transform":
			token.tokenType = TRANSFORM
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
		case "begin":
			token.tokenType = BEGIN
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
		case "if":
			token.tokenType = IF
		case "then":
			token.tokenType = THEN
		case "else":
			token.tokenType = ELSE
		case "debug":
			token.tokenType = DEBUG
		case "return":
			token.tokenType = RETURN
		case "head":
			token.tokenType = HEAD
		case "tail":
			token.tokenType = TAIL
		case "loop":
			token.tokenType = LOOP
		case "continue":
			token.tokenType = CONTINUE
		case "break":
			token.tokenType = BREAK
		}
	case SWHITESPACE:
		token.tokenType = WS
	case SOPENPAREN:
		token.tokenType = OPENPAREN
	case SCLOSEPAREN:
		token.tokenType = CLOSEPAREN
	case SOPENCURLY:
		token.tokenType = OPENCURLY
	case SCLOSECURLY:
		token.tokenType = CLOSECURLY
	case SCOMMA:
		token.tokenType = COMMA
	case SEQUAL_1:
		token.tokenType = EQUAL
	case SDEQUAL:
		token.tokenType = DEQUAL
	case SNEQUAL:
		token.tokenType = NEQUAL
	case SCOLON:
		token.tokenType = ERROR
	case SBLOCKCOMMENTSTARTEND:
		fallthrough
	case SBLOCKCOMMENTENDEND:
		fallthrough
	case SBLOCKCOMMENT:
		token.tokenType = ERROR
	case SBLOCKCOMMENTFINAL:
		fallthrough
	case SCOMMENT:
		token.tokenType = COMMENT
	case SDASH:
		token.tokenType = MINUS
	case SCOLONEQ:
		token.tokenType = COLONEQ
	case SOPERATORSTART:
		fallthrough
	case SOPERATOR:
		token.tokenType = ERROR
		lexeme := buf.String()
		switch lexeme {
		case "+":
			token.tokenType = PLUS
		case "*":
			token.tokenType = MULT
		case "/":
			token.tokenType = DIV
		case "%":
			token.tokenType = MOD
		case "<":
			token.tokenType = LESS
		case ">":
			token.tokenType = GREATER
		case "<=":
			token.tokenType = LESSEQ
		case ">=":
			token.tokenType = GREATEREQ
		}
	case SEND:
		token.tokenType = EOF
	default:
		fmt.Println(current_state)
		fmt.Println(startPosInfo.line)
		fmt.Println(startPosInfo.column)
		endPosInfo := s.get_position()
		fmt.Println(endPosInfo.line)
		fmt.Println(endPosInfo.column)
		panic("Unknown final state")
	}

	endPosInfo := s.get_position()
	token.offset = NewRange(startPosInfo.offset, endPosInfo.offset)
	token.column = NewRange(startPosInfo.column, endPosInfo.column)
	token.line = NewRange(startPosInfo.line, endPosInfo.line)
	token.lexeme = buf.String()
	//fmt.Printf("[%s] '%s' \tline: %d, \tstart column: %d, \tend column: %d\n", token.tokenType.pp(), token.lexeme, token.line.Start, token.column.Start, token.column.End)
	return token
}

func getEscapedRune(ch rune) rune {
	if ch == 'n' {
		return rune(10)
	} else if ch == 't' {
		return rune(9)
	} else if ch == 'r' {
		return rune(13)
	} else if ch == 'a' {
		return rune(7)
	} else if ch == 'b' {
		return rune(8)
	} else if ch == 'f' {
		return rune(12)
	} else if ch == 'v' {
		return rune(11)
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

func IsHex(ch rune) bool {
	return ('0' <= ch && ch <= '9') ||
		('A' <= ch && ch <= 'F') ||
		('a' <= ch && ch <= 'f')
}

func HexToAscii(ch1 rune, ch2 rune) rune {
	input := string(ch1) + string(ch2)
	value, err := strconv.ParseInt(input, 16, 64)
	if err != nil {
		panic("COULDN'T CONVERT")
	}
	//fmt.Printf("FOUND HEX RUNE (%s): %s\n", input, string(rune(value)))
	return rune(value)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
