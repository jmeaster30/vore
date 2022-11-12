package libvore

type AstCommand interface {
	isCmd()
}

type AstExpression interface {
	isExpr()
}

type AstLiteral interface {
	isLiteral()
}

type AstListable interface {
	isListable()
}

type AstFind struct {
	all  bool
	skip uint64
	take uint64
	body []*AstExpression
}

func (f AstFind) isCmd() {}

type AstReplace struct {
	all    bool
	skip   uint64
	take   uint64
	body   []*AstExpression
	result []*AstExpression
}

func (r AstReplace) isCmd() {}

type AstSet struct {
	id   string
	expr []*AstExpression
}

func (s AstSet) isCmd() {}

type AstLoop struct {
	min  uint64
	max  uint64
	body *AstLiteral
}

func (l AstLoop) isExpr() {}

type AstBranch struct {
	left  *AstLiteral
	right *AstLiteral
}

func (b AstBranch) isExpr() {}

type AstDec struct {
	isSubroutine bool
	name         string
	body         *AstLiteral
}

func (d AstDec) isExpr() {}

type AstList struct {
	contents []*AstListable
}

func (l AstList) isExpr() {}

type AstPrimary struct {
	literal *AstLiteral
}

func (s AstPrimary) isExpr() {}

type AstRange struct {
	from *AstString
	to   *AstString
}

func (r AstRange) isListable() {}

type AstString struct {
	value string
}

func (s AstString) isLiteral()  {}
func (r AstString) isListable() {}

type AstSubExpr struct {
	body []*AstExpression
}

func (n AstSubExpr) isLiteral() {}

type AstVariable struct {
	name string
}

func (s AstVariable) isLiteral() {}
