package libvore

import "fmt"

type AstCommand interface {
	isCmd()
	print()
	generate(state *GenState) (Command, error)
}

type AstExpression interface {
	isExpr()
	print()
	generate(offset int, state *GenState) ([]SearchInstruction, error)
}

type AstLiteral interface {
	isLiteral()
	print()
	generate(offset int, state *GenState) ([]SearchInstruction, error)
}

type AstListable interface {
	isListable()
	print()
	generate(offset int, state *GenState) ([]SearchInstruction, error)
	getMaxSize() int
}

type AstAtom interface {
	isAtom()
	print()
	generate(offset int, state *GenState) ([]SearchInstruction, error)
	generateReplace(offset int, state *GenState) ([]ReplaceInstruction, error)
}

type AstProcessStatement interface {
	isProcessStatement()
	print()
	check(info ProcessTypeInfo) ProcessTypeInfo
	execute(state ProcessState) ProcessState
}

type AstProcessExpression interface {
	isProcessExpr()
	print()
	check(info ProcessTypeInfo) ProcessTypeInfo
	execute(state ProcessState) ProcessState
}

type AstFind struct {
	all  bool
	skip int
	take int
	last int
	body []AstExpression
}

func (f AstFind) isCmd() {}
func (f AstFind) print() {
	fmt.Print("(find")
	if f.all {
		fmt.Print(" all")
	}
	fmt.Printf(" skip %d take %d", f.skip, f.take)
	fmt.Print(" (body")
	for _, expr := range f.body {
		fmt.Print(" ")
		expr.print()
	}
	fmt.Print("))")
}

type AstReplace struct {
	all    bool
	skip   int
	take   int
	last   int
	body   []AstExpression
	result []AstAtom
}

func (r AstReplace) isCmd() {}
func (r AstReplace) print() {
	fmt.Print("(replace")
	if r.all {
		fmt.Print(" all")
	}
	fmt.Printf(" skip %d take %d", r.skip, r.take)
	fmt.Print(" (body")
	for _, expr := range r.body {
		fmt.Print(" ")
		expr.print()
	}
	fmt.Print(") (result")
	for _, expr := range r.result {
		fmt.Print(" ")
		expr.print()
	}
	fmt.Print("))")
}

type AstSet struct {
	id   string
	body AstSetBody
}

func (s AstSet) isCmd() {}
func (s AstSet) print() {
	fmt.Printf("(set %s ", s.id)
	s.body.print()
	fmt.Print(")")
}

type AstSetBody interface {
	generate(state *GenState, id string) (SetCommandBody, error)
	print()
}

type AstSetPattern struct {
	pattern []AstExpression
	body    []AstProcessStatement
}

func (b AstSetPattern) print() {
	fmt.Print("(pattern ")
	for _, val := range b.pattern {
		val.print()
		fmt.Print(" ")
	}
	fmt.Print(") (predicate")
	for _, stmt := range b.body {
		fmt.Print(" ")
		stmt.print()
	}
	fmt.Print(")")
}

type AstSetMatches struct {
	command AstCommand
}

func (b AstSetMatches) print() {
	fmt.Print("(matches ")
	b.command.print()
	fmt.Print(")")
}

type AstSetTransform struct {
	statements []AstProcessStatement
}

func (b AstSetTransform) print() {
	fmt.Print("(transform ")
	for _, stmt := range b.statements {
		fmt.Print(" ")
		stmt.print()
	}
	fmt.Print(")")
}

type AstLoop struct {
	min    int
	max    int
	fewest bool
	body   AstExpression
	name   string
}

func (l AstLoop) isExpr() {}
func (l AstLoop) print() {
	fmt.Printf("(loop min %d max %d fewest %t ", l.min, l.max, l.fewest)
	l.body.print()
	fmt.Print(")")
}

type AstBranch struct {
	left  AstLiteral
	right AstExpression
}

func (b AstBranch) isExpr() {}
func (b AstBranch) print() {
	fmt.Print("(branch ")
	b.left.print()
	fmt.Print(" ")
	b.right.print()
	fmt.Print(")")
}

type AstDec struct {
	name string
	body AstLiteral
}

func (d AstDec) isExpr() {}
func (d AstDec) print() {
	fmt.Printf("(dec '%s' ", d.name)
	d.body.print()
	fmt.Print(")")
}

type AstSub struct {
	name string
	body []AstExpression
}

func (d AstSub) isExpr() {}
func (d AstSub) print() {
	fmt.Printf("(subdec '%s'", d.name)
	for _, expr := range d.body {
		fmt.Print(" ")
		expr.print()
	}
	fmt.Print(")")
}

type AstList struct {
	not      bool
	contents []AstListable
}

func (l AstList) isExpr() {}
func (l AstList) getMaxSize() int {
	max := -1
	for _, c := range l.contents {
		s := c.getMaxSize()
		if s > max {
			max = s
		}
	}
	return max
}

func (l AstList) print() {
	fmt.Print("(in ")
	for _, expr := range l.contents {
		fmt.Print(" ")
		expr.print()
	}
	fmt.Print(")")
}

type AstPrimary struct {
	literal AstLiteral
}

func (s AstPrimary) isExpr() {}
func (s AstPrimary) print() {
	fmt.Print("(primary ")
	s.literal.print()
	fmt.Print(")")
}

type AstRange struct {
	from *AstString
	to   *AstString
}

func (r AstRange) isListable() {}
func (r AstRange) getMaxSize() int {
	//? Can we guarantee that "from" is going to be greater than "to"??
	return len(r.to.value)
}

func (r AstRange) print() {
	fmt.Print("(range ")
	r.from.print()
	fmt.Print(" ")
	r.to.print()
	fmt.Print(")")
}

type AstString struct {
	not      bool
	value    string
	caseless bool
}

func (s AstString) isLiteral()  {}
func (s AstString) isListable() {}
func (s AstString) getMaxSize() int {
	return len(s.value)
}
func (s AstString) isAtom() {}
func (s AstString) print() {
	fmt.Printf("(string '%s')", s.value)
}

type AstSubExpr struct {
	body []AstExpression
}

func (n AstSubExpr) isLiteral() {}
func (n AstSubExpr) print() {
	fmt.Print("(subexpr")
	for _, expr := range n.body {
		fmt.Print(" ")
		expr.print()
	}
	fmt.Print(")")
}

type AstVariable struct {
	name string
}

func (s AstVariable) isLiteral() {}
func (s AstVariable) isAtom()    {}
func (s AstVariable) print() {
	fmt.Printf("(var %s)", s.name)
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
	not       bool
	classType AstCharacterClassType
}

func (c AstCharacterClass) isLiteral()  {}
func (c AstCharacterClass) isListable() {}
func (c AstCharacterClass) getMaxSize() int {
	switch c.classType {
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
	switch c.classType {
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
	name string
	expr AstProcessExpression
}

func (s AstProcessSet) isProcessStatement() {}
func (s AstProcessSet) print() {
	fmt.Printf("(pset '%s' ", s.name)
	s.expr.print()
	fmt.Print(")")
}

type AstProcessReturn struct {
	expr AstProcessExpression
}

func (s AstProcessReturn) isProcessStatement() {}
func (s AstProcessReturn) print() {
	fmt.Print("(return ")
	s.expr.print()
	fmt.Print(")")
}

type AstProcessIf struct {
	condition AstProcessExpression
	trueBody  []AstProcessStatement
	falseBody []AstProcessStatement
}

func (s AstProcessIf) isProcessStatement() {}
func (s AstProcessIf) print() {
	fmt.Print("(if ")
	s.condition.print()
	fmt.Print(" (true")
	for _, expr := range s.trueBody {
		fmt.Print(" ")
		expr.print()
	}
	fmt.Print(") (false")
	for _, expr := range s.falseBody {
		fmt.Print(" ")
		expr.print()
	}
	fmt.Print("))")
}

type AstProcessDebug struct {
	expr AstProcessExpression
}

func (s AstProcessDebug) isProcessStatement() {}
func (s AstProcessDebug) print() {
	fmt.Print("(debug ")
	s.expr.print()
	fmt.Print(")")
}

type AstProcessLoop struct {
	body []AstProcessStatement
}

func (s AstProcessLoop) isProcessStatement() {}
func (s AstProcessLoop) print() {
	fmt.Print("(loop")
	for _, expr := range s.body {
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
	op   TokenType
	expr AstProcessExpression
}

func (e AstProcessUnaryExpression) isProcessExpr() {}
func (e AstProcessUnaryExpression) print() {
	fmt.Printf("(unary %s ", e.op.PP())
	e.expr.print()
	fmt.Print(")")
}

type AstProcessBinaryExpression struct {
	op  TokenType
	lhs AstProcessExpression
	rhs AstProcessExpression
}

func (e AstProcessBinaryExpression) isProcessExpr() {}
func (e AstProcessBinaryExpression) print() {
	fmt.Printf("(binary %s ", e.op.PP())
	e.lhs.print()
	fmt.Print(" ")
	e.rhs.print()
	fmt.Print(")")
}

type AstProcessString struct {
	value string
}

func (e AstProcessString) isProcessExpr() {}
func (e AstProcessString) print() {
	fmt.Printf("(string %s)", e.value)
}

type AstProcessNumber struct {
	value int
}

func (e AstProcessNumber) isProcessExpr() {}
func (e AstProcessNumber) print() {
	fmt.Printf("(number %d)", e.value)
}

type AstProcessBoolean struct {
	value bool
}

func (e AstProcessBoolean) isProcessExpr() {}
func (e AstProcessBoolean) print() {
	fmt.Printf("(boolean %t)", e.value)
}

type AstProcessVariable struct {
	name string
}

func (e AstProcessVariable) isProcessExpr() {}
func (e AstProcessVariable) print() {
	fmt.Printf("(var %s)", e.name)
}
