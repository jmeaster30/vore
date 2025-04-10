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
	//for _, token := range tokens {
	//	fmt.Printf("[%s] '%s' \tline: %d, \tstart column: %d, \tend column: %d\n", token.tokenType.pp(), token.lexeme, token.line.Start, token.column.Start, token.column.End)
	//}

	commands, parseError := parse(tokens)
	if parseError != nil {
		return nil, parseError
	}
	return &Ast{commands}, nil
}

type AstCommand interface {
	isCmd()
	print()
	// generate(state *GenState) (Command, error)
}

type AstExpression interface {
	isExpr()
	print()
	// generate(offset int, state *GenState) ([]SearchInstruction, error)
}

type AstLiteral interface {
	isLiteral()
	print()
	// generate(offset int, state *GenState) ([]SearchInstruction, error)
}

type AstListable interface {
	isListable()
	print()
	// generate(offset int, state *GenState) ([]SearchInstruction, error)
	GetMaxSize() int
}

type AstAtom interface {
	isAtom()
	print()
	// generate(offset int, state *GenState) ([]SearchInstruction, error)
	// generateReplace(offset int, state *GenState) ([]ReplaceInstruction, error)
}

type AstProcessStatement interface {
	isProcessStatement()
	print()
	// check(info ProcessTypeInfo) ProcessTypeInfo
	// execute(state ProcessState) ProcessState
}

type AstProcessExpression interface {
	isProcessExpr()
	print()
	// check(info ProcessTypeInfo) ProcessTypeInfo
	// execute(state ProcessState) ProcessState
}

type AstFind struct {
	All  bool
	Skip int
	Take int
	Last int
	Body []AstExpression
}

func (f AstFind) isCmd() {}
func (f AstFind) print() {
	fmt.Print("(find")
	if f.All {
		fmt.Print(" all")
	}
	fmt.Printf(" skip %d take %d", f.Skip, f.Take)
	fmt.Print(" (body")
	for _, expr := range f.Body {
		fmt.Print(" ")
		expr.print()
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
func (r AstReplace) print() {
	fmt.Print("(replace")
	if r.All {
		fmt.Print(" all")
	}
	fmt.Printf(" skip %d take %d", r.Skip, r.Take)
	fmt.Print(" (body")
	for _, expr := range r.Body {
		fmt.Print(" ")
		expr.print()
	}
	fmt.Print(") (result")
	for _, expr := range r.Result {
		fmt.Print(" ")
		expr.print()
	}
	fmt.Print("))")
}

type AstSet struct {
	Id   string
	Body AstSetBody
}

func (s AstSet) isCmd() {}
func (s AstSet) print() {
	fmt.Printf("(set %s ", s.Id)
	s.Body.print()
	fmt.Print(")")
}

type AstSetBody interface {
	// generate(state *GenState, id string) (SetCommandBody, error)
	print()
}

type AstSetPattern struct {
	Pattern []AstExpression
	Body    []AstProcessStatement
}

func (b AstSetPattern) print() {
	fmt.Print("(pattern ")
	for _, val := range b.Pattern {
		val.print()
		fmt.Print(" ")
	}
	fmt.Print(") (predicate")
	for _, stmt := range b.Body {
		fmt.Print(" ")
		stmt.print()
	}
	fmt.Print(")")
}

type AstSetMatches struct {
	Command AstCommand
}

func (b AstSetMatches) print() {
	fmt.Print("(matches ")
	b.Command.print()
	fmt.Print(")")
}

type AstSetTransform struct {
	Statements []AstProcessStatement
}

func (b AstSetTransform) print() {
	fmt.Print("(transform ")
	for _, stmt := range b.Statements {
		fmt.Print(" ")
		stmt.print()
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
func (l AstLoop) print() {
	fmt.Printf("(loop min %d max %d fewest %t ", l.Min, l.Max, l.Fewest)
	l.Body.print()
	fmt.Print(")")
}

type AstBranch struct {
	Left  AstLiteral
	Right AstExpression
}

func (b AstBranch) isExpr() {}
func (b AstBranch) print() {
	fmt.Print("(branch ")
	b.Left.print()
	fmt.Print(" ")
	b.Right.print()
	fmt.Print(")")
}

type AstDec struct {
	Name string
	Body AstLiteral
}

func (d AstDec) isExpr() {}
func (d AstDec) print() {
	fmt.Printf("(dec '%s' ", d.Name)
	d.Body.print()
	fmt.Print(")")
}

type AstSub struct {
	Name string
	Body []AstExpression
}

func (d AstSub) isExpr() {}
func (d AstSub) print() {
	fmt.Printf("(subdec '%s'", d.Name)
	for _, expr := range d.Body {
		fmt.Print(" ")
		expr.print()
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

func (l AstList) print() {
	fmt.Print("(in ")
	for _, expr := range l.Contents {
		fmt.Print(" ")
		expr.print()
	}
	fmt.Print(")")
}

type AstPrimary struct {
	Literal AstLiteral
}

func (s AstPrimary) isExpr() {}
func (s AstPrimary) print() {
	fmt.Print("(primary ")
	s.Literal.print()
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

func (r AstRange) print() {
	fmt.Print("(range ")
	r.From.print()
	fmt.Print(" ")
	r.To.print()
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
func (s AstString) print() {
	fmt.Printf("(string '%s')", s.Value)
}

type AstSubExpr struct {
	Body []AstExpression
}

func (n AstSubExpr) isLiteral() {}
func (n AstSubExpr) print() {
	fmt.Print("(subexpr")
	for _, expr := range n.Body {
		fmt.Print(" ")
		expr.print()
	}
	fmt.Print(")")
}

type AstVariable struct {
	Name string
}

func (s AstVariable) isLiteral() {}
func (s AstVariable) isAtom()    {}
func (s AstVariable) print() {
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

func (c AstCharacterClass) print() {
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
func (s AstProcessSet) print() {
	fmt.Printf("(pset '%s' ", s.Name)
	s.Expr.print()
	fmt.Print(")")
}

type AstProcessReturn struct {
	Expr AstProcessExpression
}

func (s AstProcessReturn) isProcessStatement() {}
func (s AstProcessReturn) print() {
	fmt.Print("(return ")
	s.Expr.print()
	fmt.Print(")")
}

type AstProcessIf struct {
	Condition AstProcessExpression
	TrueBody  []AstProcessStatement
	FalseBody []AstProcessStatement
}

func (s AstProcessIf) isProcessStatement() {}
func (s AstProcessIf) print() {
	fmt.Print("(if ")
	s.Condition.print()
	fmt.Print(" (true")
	for _, expr := range s.TrueBody {
		fmt.Print(" ")
		expr.print()
	}
	fmt.Print(") (false")
	for _, expr := range s.FalseBody {
		fmt.Print(" ")
		expr.print()
	}
	fmt.Print("))")
}

type AstProcessDebug struct {
	Expr AstProcessExpression
}

func (s AstProcessDebug) isProcessStatement() {}
func (s AstProcessDebug) print() {
	fmt.Print("(debug ")
	s.Expr.print()
	fmt.Print(")")
}

type AstProcessLoop struct {
	Body []AstProcessStatement
}

func (s AstProcessLoop) isProcessStatement() {}
func (s AstProcessLoop) print() {
	fmt.Print("(loop")
	for _, expr := range s.Body {
		fmt.Print(" ")
		expr.print()
	}
	fmt.Print(")")
}

type AstProcessContinue struct{}

func (s AstProcessContinue) isProcessStatement() {}
func (s AstProcessContinue) print() {
	fmt.Print("(continue)")
}

type AstProcessBreak struct{}

func (s AstProcessBreak) isProcessStatement() {}
func (s AstProcessBreak) print() {
	fmt.Print("(break)")
}

type AstProcessUnaryExpression struct {
	Op   TokenType
	Expr AstProcessExpression
}

func (e AstProcessUnaryExpression) isProcessExpr() {}
func (e AstProcessUnaryExpression) print() {
	fmt.Printf("(unary %s ", e.Op.PP())
	e.Expr.print()
	fmt.Print(")")
}

type AstProcessBinaryExpression struct {
	Op  TokenType
	Lhs AstProcessExpression
	Rhs AstProcessExpression
}

func (e AstProcessBinaryExpression) isProcessExpr() {}
func (e AstProcessBinaryExpression) print() {
	fmt.Printf("(binary %s ", e.Op.PP())
	e.Lhs.print()
	fmt.Print(" ")
	e.Rhs.print()
	fmt.Print(")")
}

type AstProcessString struct {
	Value string
}

func (e AstProcessString) isProcessExpr() {}
func (e AstProcessString) print() {
	fmt.Printf("(string %s)", e.Value)
}

type AstProcessNumber struct {
	Value int
}

func (e AstProcessNumber) isProcessExpr() {}
func (e AstProcessNumber) print() {
	fmt.Printf("(number %d)", e.Value)
}

type AstProcessBoolean struct {
	Value bool
}

func (e AstProcessBoolean) isProcessExpr() {}
func (e AstProcessBoolean) print() {
	fmt.Printf("(boolean %t)", e.Value)
}

type AstProcessVariable struct {
	Name string
}

func (e AstProcessVariable) isProcessExpr() {}
func (e AstProcessVariable) print() {
	fmt.Printf("(var %s)", e.Name)
}
