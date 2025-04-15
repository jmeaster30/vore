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

type AstCommand interface {
	isCmd()
	Print()
}

type AstExpression interface {
	isExpr()
	Print()
}

type AstLiteral interface {
	isLiteral()
	Print()
}

type AstListable interface {
	isListable()
	Print()
	GetMaxSize() int
}

type AstAtom interface {
	isAtom()
	Print()
}

type AstProcessStatement interface {
	isProcessStatement()
	Print()
}

type AstProcessExpression interface {
	isProcessExpr()
	Print()
}

type AstFind struct {
	All  bool
	Skip int
	Take int
	Last int
	Body []AstExpression
}

func (f AstFind) isCmd() {}
func (f AstFind) Print() {
	fmt.Print("(find")
	if f.All {
		fmt.Print(" all")
	}
	fmt.Printf(" skip %d take %d", f.Skip, f.Take)
	fmt.Print(" (body")
	for _, expr := range f.Body {
		fmt.Print(" ")
		expr.Print()
	}
	fmt.Print("))")
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
func (r AstReplace) Print() {
	fmt.Print("(replace")
	if r.All {
		fmt.Print(" all")
	}
	fmt.Printf(" skip %d take %d", r.Skip, r.Take)
	fmt.Print(" (body")
	for _, expr := range r.Body {
		fmt.Print(" ")
		expr.Print()
	}
	fmt.Print(") (result")
	for _, expr := range r.Result {
		fmt.Print(" ")
		expr.Print()
	}
	fmt.Print("))")
}

type AstSet struct {
	Id   string
	Body AstSetBody
}

func (s AstSet) isCmd() {}
func (s AstSet) Print() {
	fmt.Printf("(set %s ", s.Id)
	s.Body.Print()
	fmt.Print(")")
}

type AstSetBody interface {
	// generate(state *GenState, id string) (SetCommandBody, error)
	Print()
}

type AstSetPattern struct {
	Pattern []AstExpression
	Body    []AstProcessStatement
}

func (b AstSetPattern) Print() {
	fmt.Print("(pattern ")
	for _, val := range b.Pattern {
		val.Print()
		fmt.Print(" ")
	}
	fmt.Print(") (predicate")
	for _, stmt := range b.Body {
		fmt.Print(" ")
		stmt.Print()
	}
	fmt.Print(")")
}

type AstSetMatches struct {
	Command AstCommand
}

func (b AstSetMatches) Print() {
	fmt.Print("(matches ")
	b.Command.Print()
	fmt.Print(")")
}

type AstSetTransform struct {
	Statements []AstProcessStatement
}

func (b AstSetTransform) Print() {
	fmt.Print("(transform ")
	for _, stmt := range b.Statements {
		fmt.Print(" ")
		stmt.Print()
	}
	fmt.Print(")")
}

type AstLoop struct {
	Min    int
	Max    int
	Fewest bool
	Body   AstExpression
	Name   string
}

func (l AstLoop) isExpr() {}
func (l AstLoop) Print() {
	fmt.Printf("(loop min %d max %d fewest %t ", l.Min, l.Max, l.Fewest)
	l.Body.Print()
	fmt.Print(")")
}

type AstBranch struct {
	Left  AstLiteral
	Right AstExpression
}

func (b AstBranch) isExpr() {}
func (b AstBranch) Print() {
	fmt.Print("(branch ")
	b.Left.Print()
	fmt.Print(" ")
	b.Right.Print()
	fmt.Print(")")
}

type AstDec struct {
	Name string
	Body AstLiteral
}

func (d AstDec) isExpr() {}
func (d AstDec) Print() {
	fmt.Printf("(dec '%s' ", d.Name)
	d.Body.Print()
	fmt.Print(")")
}

type AstSub struct {
	Name string
	Body []AstExpression
}

func (d AstSub) isExpr() {}
func (d AstSub) Print() {
	fmt.Printf("(subdec '%s'", d.Name)
	for _, expr := range d.Body {
		fmt.Print(" ")
		expr.Print()
	}
	fmt.Print(")")
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

func (l AstList) Print() {
	fmt.Print("(in ")
	for _, expr := range l.Contents {
		fmt.Print(" ")
		expr.Print()
	}
	fmt.Print(")")
}

type AstPrimary struct {
	Literal AstLiteral
}

func (s AstPrimary) isExpr() {}
func (s AstPrimary) Print() {
	fmt.Print("(primary ")
	s.Literal.Print()
	fmt.Print(")")
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

func (r AstRange) Print() {
	fmt.Print("(range ")
	r.From.Print()
	fmt.Print(" ")
	r.To.Print()
	fmt.Print(")")
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
func (s AstString) Print() {
	fmt.Printf("(string '%s')", s.Value)
}

type AstSubExpr struct {
	Body []AstExpression
}

func (n AstSubExpr) isLiteral() {}
func (n AstSubExpr) Print() {
	fmt.Print("(subexpr")
	for _, expr := range n.Body {
		fmt.Print(" ")
		expr.Print()
	}
	fmt.Print(")")
}

type AstVariable struct {
	Name string
}

func (s AstVariable) isLiteral() {}
func (s AstVariable) isAtom()    {}
func (s AstVariable) Print() {
	fmt.Printf("(var %s)", s.Name)
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

func (c AstCharacterClass) Print() {
	fmt.Printf("(class ")
	switch c.ClassType {
	case ClassAny:
		fmt.Print("any")
	case ClassWhitespace:
		fmt.Print("whitespace")
	case ClassDigit:
		fmt.Print("digit")
	case ClassUpper:
		fmt.Print("upper")
	case ClassLower:
		fmt.Print("lower")
	case ClassLetter:
		fmt.Print("letter")
	case ClassLineStart:
		fmt.Print("line start")
	case ClassFileStart:
		fmt.Print("file start")
	case ClassLineEnd:
		fmt.Print("line end")
	case ClassFileEnd:
		fmt.Print("file end")
	case ClassWordStart:
		fmt.Print("word start")
	case ClassWordEnd:
		fmt.Print("word end")
	case ClassWholeFile:
		fmt.Print("whole file")
	case ClassWholeLine:
		fmt.Print("whole line")
	case ClassWholeWord:
		fmt.Print("whole word")
	}
	fmt.Printf(")")
}

type AstProcessSet struct {
	Name string
	Expr AstProcessExpression
}

func (s AstProcessSet) isProcessStatement() {}
func (s AstProcessSet) Print() {
	fmt.Printf("(pset '%s' ", s.Name)
	s.Expr.Print()
	fmt.Print(")")
}

type AstProcessReturn struct {
	Expr AstProcessExpression
}

func (s AstProcessReturn) isProcessStatement() {}
func (s AstProcessReturn) Print() {
	fmt.Print("(return ")
	s.Expr.Print()
	fmt.Print(")")
}

type AstProcessIf struct {
	Condition AstProcessExpression
	TrueBody  []AstProcessStatement
	FalseBody []AstProcessStatement
}

func (s AstProcessIf) isProcessStatement() {}
func (s AstProcessIf) Print() {
	fmt.Print("(if ")
	s.Condition.Print()
	fmt.Print(" (true")
	for _, expr := range s.TrueBody {
		fmt.Print(" ")
		expr.Print()
	}
	fmt.Print(") (false")
	for _, expr := range s.FalseBody {
		fmt.Print(" ")
		expr.Print()
	}
	fmt.Print("))")
}

type AstProcessDebug struct {
	Expr AstProcessExpression
}

func (s AstProcessDebug) isProcessStatement() {}
func (s AstProcessDebug) Print() {
	fmt.Print("(debug ")
	s.Expr.Print()
	fmt.Print(")")
}

type AstProcessLoop struct {
	Body []AstProcessStatement
}

func (s AstProcessLoop) isProcessStatement() {}
func (s AstProcessLoop) Print() {
	fmt.Print("(loop")
	for _, expr := range s.Body {
		fmt.Print(" ")
		expr.Print()
	}
	fmt.Print(")")
}

type AstProcessContinue struct{}

func (s AstProcessContinue) isProcessStatement() {}
func (s AstProcessContinue) Print() {
	fmt.Print("(continue)")
}

type AstProcessBreak struct{}

func (s AstProcessBreak) isProcessStatement() {}
func (s AstProcessBreak) Print() {
	fmt.Print("(break)")
}

type AstProcessUnaryExpression struct {
	Op   TokenType
	Expr AstProcessExpression
}

func (e AstProcessUnaryExpression) isProcessExpr() {}
func (e AstProcessUnaryExpression) Print() {
	fmt.Printf("(unary %s ", e.Op.PP())
	e.Expr.Print()
	fmt.Print(")")
}

type AstProcessBinaryExpression struct {
	Op  TokenType
	Lhs AstProcessExpression
	Rhs AstProcessExpression
}

func (e AstProcessBinaryExpression) isProcessExpr() {}
func (e AstProcessBinaryExpression) Print() {
	fmt.Printf("(binary %s ", e.Op.PP())
	e.Lhs.Print()
	fmt.Print(" ")
	e.Rhs.Print()
	fmt.Print(")")
}

type AstProcessString struct {
	Value string
}

func (e AstProcessString) isProcessExpr() {}
func (e AstProcessString) Print() {
	fmt.Printf("(string %s)", e.Value)
}

type AstProcessNumber struct {
	Value int
}

func (e AstProcessNumber) isProcessExpr() {}
func (e AstProcessNumber) Print() {
	fmt.Printf("(number %d)", e.Value)
}

type AstProcessBoolean struct {
	Value bool
}

func (e AstProcessBoolean) isProcessExpr() {}
func (e AstProcessBoolean) Print() {
	fmt.Printf("(boolean %t)", e.Value)
}

type AstProcessVariable struct {
	Name string
}

func (e AstProcessVariable) isProcessExpr() {}
func (e AstProcessVariable) Print() {
	fmt.Printf("(var %s)", e.Name)
}
