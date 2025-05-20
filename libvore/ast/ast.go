package ast

import (
	"fmt"
	"io"
)

type Ast struct {
	commands []AstCommand
}

func (ast *Ast) Commands() []AstCommand {
	return ast.commands
}

func ParseReader(reader io.Reader) (*Ast, error) {
	lexer := initLexer(reader)

	tokens, lexError := lexer.getTokens()
	if lexError != nil {
		return nil, lexError
	}

	commands, parseError := parse(tokens)
	if parseError != nil {
		return nil, parseError
	}
	return &Ast{commands}, nil
}

type AstNode interface {
	NodeString() string
}

type AstCommand interface {
	AstNode
	isCmd()
}

type AstExpression interface {
	AstNode
	isExpr()
}

type AstLiteral interface {
	AstNode
	isLiteral()
}

type AstListable interface {
	AstNode
	isListable()
	GetMaxSize() int
}

type AstAtom interface {
	AstNode
	isAtom()
}

type AstProcessStatement interface {
	AstNode
	isProcessStatement()
}

type AstProcessExpression interface {
	AstNode
	isProcessExpr()
}

type AstFind struct {
	All  bool
	Skip int
	Take int
	Last int
	Body []AstExpression
}

func (f AstFind) isCmd() {}
func (f AstFind) NodeString() string {
	result := "(find"
	if f.All {
		result += " all"
	}
	result += fmt.Sprintf(" skip %d take %d", f.Skip, f.Take)
	result += " (body"
	for _, expr := range f.Body {
		result += fmt.Sprintf(" %s", expr.NodeString())
	}
	result += "))"
	return result
}

type AstReplace struct {
	All    bool
	Skip   int
	Take   int
	Last   int
	Body   []AstExpression
	Result []AstAtom
}

func (r AstReplace) isCmd() {}
func (r AstReplace) NodeString() string {
	result := "(replace"
	if r.All {
		result = " all"
	}
	result += fmt.Sprintf(" skip %d take %d (body", r.Skip, r.Take)
	for _, expr := range r.Body {
		result += fmt.Sprintf(" %s", expr.NodeString())
	}
	result += ") (result"
	for _, expr := range r.Result {
		result += fmt.Sprintf(" %s", expr.NodeString())
	}
	result += "))"
	return result
}

type AstSet struct {
	Id   string
	Body AstSetBody
}

func (s AstSet) isCmd() {}
func (s AstSet) NodeString() string {
	return fmt.Sprintf("(set %s %s)", s.Id, s.Body.NodeString())
}

type AstSetBody interface {
	// generate(state *GenState, id string) (SetCommandBody, error)
	NodeString() string
}

type AstSetPattern struct {
	Pattern []AstExpression
	Body    []AstProcessStatement
}

func (b AstSetPattern) NodeString() string {
	result := "(pattern "
	for _, val := range b.Pattern {
		result += fmt.Sprintf("%s ", val.NodeString())
	}
	result += ") (predicate"
	for _, stmt := range b.Body {
		result += fmt.Sprintf(" %s", stmt.NodeString())
	}
	result += ")"
	return result
}

type AstSetMatches struct {
	Command AstCommand
}

func (b AstSetMatches) NodeString() string {
	return fmt.Sprintf("(matches %s)", b.Command.NodeString())
}

type AstSetTransform struct {
	Statements []AstProcessStatement
}

func (b AstSetTransform) NodeString() string {
	result := "(transform "
	for _, stmt := range b.Statements {
		result += fmt.Sprintf(" %s", stmt.NodeString())
	}
	result += ")"
	return result
}

type AstLoop struct {
	Min    int
	Max    int
	Fewest bool
	Body   AstExpression
	Name   string
}

func (l AstLoop) isExpr() {}
func (l AstLoop) NodeString() string {
	return fmt.Sprintf("(loop min %d max %d fewest %t %s)", l.Min, l.Max, l.Fewest, l.Body.NodeString())
}

type AstBranch struct {
	Left  AstLiteral
	Right AstExpression
}

func (b AstBranch) isExpr() {}
func (b AstBranch) NodeString() string {
	return fmt.Sprintf("(branch %s %s)", b.Left.NodeString(), b.Right.NodeString())
}

type AstDec struct {
	Name string
	Body AstLiteral
}

func (d AstDec) isExpr() {}
func (d AstDec) NodeString() string {
	return fmt.Sprintf("(dec '%s' %s)", d.Name, d.Body.NodeString())
}

type AstSub struct {
	Name string
	Body []AstExpression
}

func (d AstSub) isExpr() {}
func (d AstSub) NodeString() string {
	result := fmt.Sprintf("(subdec '%s'", d.Name)
	for _, expr := range d.Body {
		result += fmt.Sprintf(" %s", expr.NodeString())
	}
	result += ")"
	return result
}

type AstList struct {
	Not      bool
	Contents []AstListable
}

func (l AstList) isExpr() {}
func (l AstList) GetMaxSize() int {
	max := -1
	for _, c := range l.Contents {
		s := c.GetMaxSize()
		if s > max {
			max = s
		}
	}
	return max
}

func (l AstList) NodeString() string {
	result := "(in "
	for _, expr := range l.Contents {
		result += fmt.Sprintf(" %s", expr.NodeString())
	}
	result += ")"
	return result
}

type AstPrimary struct {
	Literal AstLiteral
}

func (s AstPrimary) isExpr() {}
func (s AstPrimary) NodeString() string {
	return fmt.Sprintf("(primary %s)", s.Literal.NodeString())
}

type AstRange struct {
	From *AstString
	To   *AstString
}

func (r AstRange) isListable() {}
func (r AstRange) GetMaxSize() int {
	//? Can we guarantee that "from" is going to be greater than "to"??
	return len(r.To.Value)
}

func (r AstRange) NodeString() string {
	return fmt.Sprintf("(range %s %s)", r.From.NodeString(), r.To.NodeString())
}

type AstString struct {
	Not      bool
	Value    string
	Caseless bool
}

func (s AstString) isLiteral()  {}
func (s AstString) isListable() {}
func (s AstString) GetMaxSize() int {
	return len(s.Value)
}
func (s AstString) isAtom() {}
func (s AstString) NodeString() string {
	return fmt.Sprintf("(string '%s')", s.Value)
}

type AstSubExpr struct {
	Body []AstExpression
}

func (n AstSubExpr) isLiteral() {}
func (n AstSubExpr) NodeString() string {
	result := "(subexpr"
	for _, expr := range n.Body {
		result += fmt.Sprintf(" %s", expr.NodeString())
	}
	result += ")"
	return result
}

type AstVariable struct {
	Name string
}

func (s AstVariable) isLiteral() {}
func (s AstVariable) isAtom()    {}
func (s AstVariable) NodeString() string {
	return fmt.Sprintf("(var %s)", s.Name)
}

type AstCharacterClassType int

const (
	ClassAny AstCharacterClassType = iota
	ClassWhitespace
	ClassDigit
	ClassUpper
	ClassLower
	ClassLetter
	ClassLineStart
	ClassFileStart
	ClassWordStart
	ClassLineEnd
	ClassFileEnd
	ClassWordEnd
	ClassWholeLine
	ClassWholeFile
	ClassWholeWord
)

func (a AstCharacterClassType) String() string {
	switch a {
	case ClassAny:
		return "ANY"
	case ClassWhitespace:
		return "WS"
	case ClassDigit:
		return "DIGIT"
	case ClassUpper:
		return "UPPER"
	case ClassLower:
		return "LOWER"
	case ClassLetter:
		return "LETTER"
	case ClassLineStart:
		return "LStart"
	case ClassFileStart:
		return "FStart"
	case ClassWordStart:
		return "WStart"
	case ClassLineEnd:
		return "LEnd"
	case ClassFileEnd:
		return "FEnd"
	case ClassWordEnd:
		return "WEnd"
	case ClassWholeLine:
		return "WLine"
	case ClassWholeFile:
		return "WFile"
	case ClassWholeWord:
		return "WWord"
	}
	return "MISSING CHAR CLASS"
}

type AstCharacterClass struct {
	Not       bool
	ClassType AstCharacterClassType
}

func (c AstCharacterClass) isLiteral()  {}
func (c AstCharacterClass) isListable() {}
func (c AstCharacterClass) GetMaxSize() int {
	switch c.ClassType {
	case ClassAny:
		return 1
	case ClassWhitespace:
		return 1
	case ClassDigit:
		return 1
	case ClassUpper:
		return 1
	case ClassLower:
		return 1
	case ClassLetter:
		return 1
	case ClassLineStart:
		return 0
	case ClassFileStart:
		return 0
	case ClassWordStart:
		return 0
	case ClassLineEnd:
		return 0
	case ClassFileEnd:
		return 0
	case ClassWordEnd:
		return 0
	case ClassWholeFile:
		return -1 // TODO i don't know what to do for these
	case ClassWholeLine:
		return -1
	case ClassWholeWord:
		return -1
	}
	panic("shouldn't get here")
}

func (c AstCharacterClass) NodeString() string {
	result := "(class "
	switch c.ClassType {
	case ClassAny:
		result += "any"
	case ClassWhitespace:
		result += "whitespace"
	case ClassDigit:
		result += "digit"
	case ClassUpper:
		result += "upper"
	case ClassLower:
		result += "lower"
	case ClassLetter:
		result += "letter"
	case ClassLineStart:
		result += "line start"
	case ClassFileStart:
		result += "file start"
	case ClassLineEnd:
		result += "line end"
	case ClassFileEnd:
		result += "file end"
	case ClassWordStart:
		result += "word start"
	case ClassWordEnd:
		result += "word end"
	case ClassWholeFile:
		result += "whole file"
	case ClassWholeLine:
		result += "whole line"
	case ClassWholeWord:
		result += "whole word"
	}
	result += ")"
	return result
}

type AstProcessSet struct {
	Name string
	Expr AstProcessExpression
}

func (s AstProcessSet) isProcessStatement() {}
func (s AstProcessSet) NodeString() string {
	return fmt.Sprintf("(pset '%s' %s)", s.Name, s.Expr.NodeString())
}

type AstProcessReturn struct {
	Expr AstProcessExpression
}

func (s AstProcessReturn) isProcessStatement() {}
func (s AstProcessReturn) NodeString() string {
	return fmt.Sprintf("(return %s)", s.Expr.NodeString())
}

type AstProcessIf struct {
	Condition AstProcessExpression
	TrueBody  []AstProcessStatement
	FalseBody []AstProcessStatement
}

func (s AstProcessIf) isProcessStatement() {}
func (s AstProcessIf) NodeString() string {
	result := fmt.Sprintf("(if %s (true", s.Condition.NodeString())
	for _, expr := range s.TrueBody {
		result += fmt.Sprintf(" %s", expr.NodeString())
	}
	result += ") (false"
	for _, expr := range s.FalseBody {
		result += fmt.Sprintf(" %s", expr.NodeString())
	}
	result += "))"
	return result
}

type AstProcessDebug struct {
	Expr AstProcessExpression
}

func (s AstProcessDebug) isProcessStatement() {}
func (s AstProcessDebug) NodeString() string {
	return fmt.Sprintf("(debug %s)", s.Expr.NodeString())
}

type AstProcessLoop struct {
	Body []AstProcessStatement
}

func (s AstProcessLoop) isProcessStatement() {}
func (s AstProcessLoop) NodeString() string {
	result := "(loop"
	for _, expr := range s.Body {
		result += fmt.Sprintf(" %s", expr.NodeString())
	}
	result += ")"
	return result
}

type AstProcessContinue struct{}

func (s AstProcessContinue) isProcessStatement() {}
func (s AstProcessContinue) NodeString() string {
	return "(continue)"
}

type AstProcessBreak struct{}

func (s AstProcessBreak) isProcessStatement() {}
func (s AstProcessBreak) NodeString() string {
	return "(break)"
}

type AstProcessUnaryExpression struct {
	Op   TokenType
	Expr AstProcessExpression
}

func (e AstProcessUnaryExpression) isProcessExpr() {}
func (e AstProcessUnaryExpression) NodeString() string {
	return fmt.Sprintf("(unary %s %s)", e.Op.PP(), e.Expr.NodeString())
}

type AstProcessBinaryExpression struct {
	Op  TokenType
	Lhs AstProcessExpression
	Rhs AstProcessExpression
}

func (e AstProcessBinaryExpression) isProcessExpr() {}
func (e AstProcessBinaryExpression) NodeString() string {
	return fmt.Sprintf("(binary %s %s %s)", e.Op.PP(), e.Lhs.NodeString(), e.Rhs.NodeString())
}

type AstProcessString struct {
	Value string
}

func (e AstProcessString) isProcessExpr() {}
func (e AstProcessString) NodeString() string {
	return fmt.Sprintf("(string %s)", e.Value)
}

type AstProcessNumber struct {
	Value int
}

func (e AstProcessNumber) isProcessExpr() {}
func (e AstProcessNumber) NodeString() string {
	return fmt.Sprintf("(number %d)", e.Value)
}

type AstProcessBoolean struct {
	Value bool
}

func (e AstProcessBoolean) isProcessExpr() {}
func (e AstProcessBoolean) NodeString() string {
	return fmt.Sprintf("(boolean %t)", e.Value)
}

type AstProcessVariable struct {
	Name string
}

func (e AstProcessVariable) isProcessExpr() {}
func (e AstProcessVariable) NodeString() string {
	return fmt.Sprintf("(var %s)", e.Name)
}
