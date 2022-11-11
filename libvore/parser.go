package libvore

type AstType int

const (
	AstFind AstType = iota
	AstReplace
	AstSet
	AstLoop
	AstBranch
	AstAtom
	AstDec
)

func parse(tokens []*Token) {}
