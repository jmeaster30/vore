package libvore

import "fmt"

type AstCommand interface {
	isCmd()
	print()
}

type AstExpression interface {
	isExpr()
	print()
}

type AstLiteral interface {
	isLiteral()
	print()
}

type AstListable interface {
	isListable()
	print()
}

type AstAtom interface {
	isAtom()
	print()
}

type AstFind struct {
	all  bool
	skip int
	take int
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
	expr AstExpression
}

func (s AstSet) isCmd() {}
func (s AstSet) print() {
	fmt.Printf("(set %s ", s.id)
	s.expr.print()
	fmt.Print(")")
}

type AstLoop struct {
	min    uint64
	max    uint64
	fewest bool
	body   AstLiteral
}

func (l AstLoop) isExpr() {}
func (l AstLoop) print() {
	fmt.Printf("(loop min %d max %d fewest %t ", l.min, l.max, l.fewest)
	l.body.print()
	fmt.Print(")")
}

type AstBranch struct {
	left  AstLiteral
	right AstLiteral
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
	isSubroutine bool
	name         string
	body         AstLiteral
}

func (d AstDec) isExpr() {}
func (d AstDec) print() {
	fmt.Print("(dec ")
	if d.isSubroutine {
		fmt.Print("sub ")
	}
	fmt.Printf("id %s ", d.name)
	d.body.print()
	fmt.Print(")")
}

type AstList struct {
	contents []AstListable
}

func (l AstList) isExpr() {}
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
func (r AstRange) print() {
	fmt.Print("(range ")
	r.from.print()
	fmt.Print(" ")
	r.to.print()
	fmt.Print(")")
}

type AstString struct {
	value string
}

func (s AstString) isLiteral()  {}
func (s AstString) isListable() {}
func (s AstString) isAtom()     {}
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
	ClassLineEnd
	ClassFileEnd
)

type AstCharacterClass struct {
	classType AstCharacterClassType
}

func (c AstCharacterClass) isLiteral() {}
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
	}
	fmt.Printf(")")
}
